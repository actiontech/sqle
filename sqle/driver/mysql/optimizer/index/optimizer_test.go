package index

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/pingcap/parser/format"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func mockExplain(mocker sqlmock.Sqlmock, rows [][]string) {
	e := mocker.ExpectQuery(regexp.QuoteMeta("EXPLAIN"))
	r := sqlmock.NewRows([]string{"id", "table", "type"})
	for _, row := range rows {
		e.WillReturnRows(r.AddRow(row[0], row[1], row[2]))
	}
}

var testLogger = logrus.New()

func TestOptimizer_Optimize(t *testing.T) {
	entry := testLogger.WithFields(logrus.Fields{"test": "optimizer"})
	var optimizerTests = []struct {
		SQL             string
		explain         [][]string
		optimizerOption []optimizerOption

		output []*OptimizeResult
	}{
		// single table, single select
		{"select 1", [][]string{}, nil, nil},
		{"select * from t where id = 1", [][]string{{"1", "t", "const"}}, nil, nil},
		{"select * from t where c1 = 1", [][]string{{"1", "t", executor.ExplainRecordAccessTypeAll}}, nil, []*OptimizeResult{{"t", []string{"c1"}, ""}}},
		{"select * from t where c1 = 1", [][]string{{"1", "t", executor.ExplainRecordAccessTypeIndex}}, nil, []*OptimizeResult{{"t", []string{"c1"}, ""}}},
		{"select * from t where c1 = 1 and c2 = 2 and c3 > 3", [][]string{{"1", "t", executor.ExplainRecordAccessTypeIndex}}, nil, []*OptimizeResult{{"t", []string{"c1", "c2"}, ""}}},
		{"select c1,c2,c3 from t where c2 = 1 and c1 = 2 and c3 > 3", [][]string{{"1", "t", executor.ExplainRecordAccessTypeIndex}}, nil, []*OptimizeResult{{"t", []string{"c2", "c1", "c3"}, ""}}},
		{"select c1,c2,c3,c4 from t where c2 = 1 and c1 = 2 and c3 > 3", [][]string{{"1", "t", executor.ExplainRecordAccessTypeIndex}}, []optimizerOption{WithCompositeIndexMaxColumn(4)}, []*OptimizeResult{{"t", []string{"c2", "c1", "c3", "c4"}, ""}}},

		// multi table, single select
		{"select * from t1 join t2 on t1.c1 = t2.c1", [][]string{{"1", "t1", executor.ExplainRecordAccessTypeAll}, {"1", "t2", executor.ExplainRecordAccessTypeAll}}, nil, []*OptimizeResult{{"t2", []string{"c1"}, ""}}},
		{"select * from t1 join t2 on t1.c1 = t2.c1", [][]string{{"1", "t2", executor.ExplainRecordAccessTypeAll}, {"1", "t1", executor.ExplainRecordAccessTypeAll}}, nil, []*OptimizeResult{{"t1", []string{"c1"}, ""}}},
		{"select * from t1 join t2 on t1.c1 = t2.c1", [][]string{{"1", "t1", executor.ExplainRecordAccessTypeAll}, {"1", "t2", "ref"}}, nil, nil},

		// subqueries
		{"select * from (select c1,c2 from t2 where c1 = 2) as t1", [][]string{{"1", "t2", executor.ExplainRecordAccessTypeIndex}}, nil, []*OptimizeResult{{"t2", []string{"c1", "c2"}, ""}}},
	}
	for i, tt := range optimizerTests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			ss, err := parser.New().ParseOneStmt(tt.SQL, "", "")
			assert.NoError(t, err)
			e, mocker, err := executor.NewMockExecutor()
			assert.NoError(t, err)

			mockExplain(mocker, tt.explain)
			o := NewOptimizer(entry, session.NewMockContext(e), tt.optimizerOption...)

			gots, err := o.Optimize(context.TODO(), ss.(*ast.SelectStmt))
			assert.NoError(t, err)
			for i, got := range gots {
				assert.Equal(t, tt.output[i].TableName, got.TableName)
				assert.Equal(t, tt.output[i].IndexedColumns, got.IndexedColumns)
			}
			mocker.MatchExpectationsInOrder(true)
		})
	}
}

func TestOptimizer_parseSelectStmt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input string
		sel   map[string] /*table name*/ string /*select SQL*/
		join  map[string] /*table name*/ string /*join on column*/
	}{
		// single select(single table)
		{"select 1", nil, nil},
		{"select * from t1", map[string]string{"t1": "SELECT * FROM t1"}, nil},
		{"select * from t1 as t2", map[string]string{"t2": "SELECT * FROM t1 AS t2"}, nil},
		// single select(multi table/join)
		{"select * from t1 join t2 on t1.id = t2.id", nil, map[string]string{"t1": "id", "t2": "id"}},
		{"select * from t1 left join t2 on t1.id = t2.id", nil, map[string]string{"t1": "id", "t2": "id"}},
		{"select * from t1 right join t2 on t1.id = t2.id", nil, map[string]string{"t1": "id", "t2": "id"}},
		{"select * from t1 as t1_alias join t2 as t2_alias on t1_alias.id = t2_alias.id", nil, map[string]string{"t1_alias": "id", "t2_alias": "id"}},
		// multi select
		{"select * from (select * from t1) as t2", map[string]string{"t2": "SELECT * FROM (SELECT * FROM (t1)) AS t2", "t1": "SELECT * FROM t1"}, nil},
		{"select * from t1 where id = (select * from t2)", map[string]string{"t1": "SELECT * FROM t1 WHERE id=(SELECT * FROM t2)", "t2": "SELECT * FROM t2"}, nil},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			stmt, err := parser.New().ParseOneStmt(tt.input, "", "")
			assert.NoError(t, err)

			o := Optimizer{tables: map[string]*tableInSelect{}}
			o.parseSelectStmt(stmt.(*ast.SelectStmt))
			for n, tbl := range o.tables {
				if tbl.singleTableSel == nil {
					c, ok := tt.join[n]
					assert.True(t, ok)
					assert.Equal(t, c, tbl.joinOnColumn)
				} else {
					var buf strings.Builder
					assert.NoError(t, tbl.singleTableSel.Restore(format.NewRestoreCtx(0, &buf)))
					assert.Equal(t, tt.sel[n], buf.String())
				}
			}
		})
	}
}

func Test_removeDrivingTable(t *testing.T) {
	tests := []struct {
		input  []*executor.ExplainRecord
		output []*executor.ExplainRecord
	}{
		{[]*executor.ExplainRecord{}, []*executor.ExplainRecord{}},
		{[]*executor.ExplainRecord{{Id: "1", Table: "t1"}}, []*executor.ExplainRecord{{Id: "1", Table: "t1"}}},
		{[]*executor.ExplainRecord{{Id: "1", Table: "t1"}, {Id: "1", Table: "t2"}}, []*executor.ExplainRecord{{Id: "1", Table: "t2"}}},
		{[]*executor.ExplainRecord{{Id: "1", Table: "t1"}, {Id: "1", Table: "t2"}, {Id: "2", Table: "t3"}}, []*executor.ExplainRecord{{Id: "1", Table: "t2"}, {Id: "2", Table: "t3"}}},
		{[]*executor.ExplainRecord{{Id: "1", Table: "t1"}, {Id: "1", Table: "t2"}, {Id: "2", Table: "t3"}, {Id: "3", Table: "t4"}}, []*executor.ExplainRecord{{Id: "1", Table: "t2"}, {Id: "2", Table: "t3"}, {Id: "3", Table: "t4"}}},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got := removeDrivingTable(tt.input)
			for i, g := range got {
				assert.Equal(t, tt.output[i].Id, g.Id)
				assert.Equal(t, tt.output[i].Table, g.Table)
			}
		})
	}
}

func TestOptimizer_needIndex(t *testing.T) {
	tests := []struct {
		tableName   string
		indexColumn []string
		want        bool
	}{
		{"exist_tb_1", []string{"v2", "v1"}, true},
		{"exist_tb_3", []string{"v1", "v2", "v3"}, true},

		{"exist_tb_1", []string{"id"}, false},
		{"exist_tb_1", []string{"v1", "v2"}, false},
		{"exist_tb_1", []string{"v1"}, false},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			e, _, err := executor.NewMockExecutor()
			assert.NoError(t, err)

			o := NewOptimizer(testLogger.WithField("test", "test"), session.NewMockContext(e))
			mockSelect := fmt.Sprintf("select * from %s", tt.tableName)
			stmt, err := parser.New().ParseOneStmt(mockSelect, "", "")
			assert.NoError(t, err)
			o.tables[tt.tableName] = &tableInSelect{singleTableSel: stmt.(*ast.SelectStmt)}
			got, err := o.needIndex(tt.tableName, tt.indexColumn...)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
