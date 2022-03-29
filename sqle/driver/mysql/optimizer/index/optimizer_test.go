package index

import (
	"context"
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/format"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var testLogger = logrus.New()

func TestOptimizer_Optimize(t *testing.T) {
	entry := testLogger.WithFields(logrus.Fields{"test": "optimizer"})

	type databaseMock struct {
		expectQuery string
		rows        [][]string
	}

	explainHead := []string{"id", "table", "type"}
	showTableStatusHead := []string{"Name", "Rows"}
	cardinalityHead := []string{"cardinality"}
	showGlobalVariableHead := []string{"Variable_name", "Value"}

	var optimizerTests = []struct {
		SQL             string
		databaseMocks   []databaseMock
		optimizerOption []optimizerOption

		// output
		output []*OptimizeResult
	}{
		{"select 1", []databaseMock{}, nil, nil},
		{
			"select * from exist_tb_1 where id = 1",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "exist_tb_1", "const"}}},
			},
			nil,
			nil,
		},
		{
			"select * from exist_tb_1 as t where id = 1",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "t", "const"}}},
			},
			nil,
			nil,
		},
		{
			"select * from exist_tb_3 where v1 = 1",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "exist_tb_3", executor.ExplainRecordAccessTypeAll}}},
			},
			nil,
			[]*OptimizeResult{{"exist_tb_3", []string{"v1"}, ""}},
		},
		{
			"select * from exist_tb_3 where v1 = 1",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "exist_tb_3", executor.ExplainRecordAccessTypeIndex}}},
			},
			nil,
			[]*OptimizeResult{{"exist_tb_3", []string{"v1"}, ""}},
		},
		{
			"select * from exist_tb_3 where v1 = 1 and v2 = 2 and v3 > 3",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "exist_tb_3", executor.ExplainRecordAccessTypeIndex}}},
				{"show table status", [][]string{showTableStatusHead, {"exist_tb_3", "1000"}}},
				{"select count(distinct `v1`)", [][]string{cardinalityHead, {"100"}}},
				{"select count(distinct `v2`)", [][]string{cardinalityHead, {"101"}}},
			},
			nil,
			[]*OptimizeResult{{"exist_tb_3", []string{"v2", "v1"}, ""}},
		},

		{
			"select * from exist_tb_3 where v1 = 1 and v2 = 2 and v3 > 3",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "exist_tb_3", executor.ExplainRecordAccessTypeIndex}}},
				{"show table status", [][]string{showTableStatusHead, {"exist_tb_3", "1000"}}},
				{"select count(distinct `v1`)", [][]string{cardinalityHead, {"101"}}},
				{"select count(distinct `v2`)", [][]string{cardinalityHead, {"100"}}},
			},
			nil,
			[]*OptimizeResult{{"exist_tb_3", []string{"v1", "v2"}, ""}},
		},
		{
			"select v1,v2,v3 from exist_tb_3 where v2 = 1 and v1 = 2 and v3 > 3",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "exist_tb_3", executor.ExplainRecordAccessTypeIndex}}},
				{"show table status", [][]string{showTableStatusHead, {"exist_tb_3", "1000"}}},
				{"select count(distinct `v2`)", [][]string{cardinalityHead, {"102"}}},
				{"select count(distinct `v1`)", [][]string{cardinalityHead, {"101"}}},
				{"select count(distinct `v3`)", [][]string{cardinalityHead, {"100"}}},
			},
			nil,
			[]*OptimizeResult{{"exist_tb_3", []string{"v2", "v1", "v3"}, ""}},
		},
		{
			"select v1,v2,v3 from exist_tb_3 where v2 = 1 and v1 = 2 and v3 > 3",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "exist_tb_3", executor.ExplainRecordAccessTypeIndex}}},
				{"show table status", [][]string{showTableStatusHead, {"exist_tb_3", "1000"}}},
				{"select count(distinct `v2`)", [][]string{cardinalityHead, {"102"}}},
				{"select count(distinct `v1`)", [][]string{cardinalityHead, {"101"}}},
			},
			[]optimizerOption{WithCompositeIndexMaxColumn(2)},
			[]*OptimizeResult{{"exist_tb_3", []string{"v2", "v1"}, ""}},
		},
		// multi table, single select
		{
			"select * from exist_tb_1 join exist_tb_2 on exist_tb_1.v1 = exist_tb_2.v1",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "exist_tb_1", executor.ExplainRecordAccessTypeAll}, {"1", "exist_tb_2", executor.ExplainRecordAccessTypeAll}}},
			},
			nil,
			[]*OptimizeResult{{"exist_tb_2", []string{"v1"}, ""}},
		},
		{
			"select * from exist_tb_1 join exist_tb_2 on exist_tb_1.v1 = exist_tb_2.v1",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "exist_tb_1", executor.ExplainRecordAccessTypeAll}, {"1", "exist_tb_2", "ref"}}},
			},
			nil,
			nil,
		},
		{
			"select * from exist_tb_1 join exist_tb_2 using(v1)",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "exist_tb_1", executor.ExplainRecordAccessTypeAll}, {"1", "exist_tb_2", executor.ExplainRecordAccessTypeAll}}},
			},
			nil,
			[]*OptimizeResult{{"exist_tb_2", []string{"v1"}, ""}},
		},
		// will not give advice when join without condition
		{
			"select * from exist_tb_1 join exist_tb_2",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "exist_tb_1", executor.ExplainRecordAccessTypeAll}, {"1", "exist_tb_2", "ref"}}},
			},
			nil,
			nil,
		},
		{
			"select * from exist_tb_1 cross join exist_tb_2",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "exist_tb_1", executor.ExplainRecordAccessTypeAll}, {"1", "exist_tb_2", "ref"}}},
			},
			nil,
			nil,
		},
		{
			"select * from exist_tb_1, exist_tb_2",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "exist_tb_1", executor.ExplainRecordAccessTypeAll}, {"1", "exist_tb_2", "ref"}}},
			},
			nil,
			nil,
		},
		// sub-queries
		{
			"select * from (select v1,v2 from exist_tb_2 where v1 = 2) as t1",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "exist_tb_2", executor.ExplainRecordAccessTypeIndex}}},
				{"show table status", [][]string{showTableStatusHead, {"exist_tb_2", "1000"}}},
				{"select count(distinct `v1`)", [][]string{cardinalityHead, {"100"}}},
				{"select count(distinct `v2`)", [][]string{cardinalityHead, {"101"}}},
			},
			nil,
			[]*OptimizeResult{{"exist_tb_2", []string{"v2", "v1"}, ""}},
		},
		{
			"select * from exist_tb_2 where left(v3, 5) = 'hello'",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "exist_tb_2", executor.ExplainRecordAccessTypeIndex}}},
				{"SHOW GLOBAL VARIABLES", [][]string{showGlobalVariableHead, {"version", "5.6.12"}}},
			},
			nil,
			nil,
		},
		{
			"select * from exist_tb_2 where left(v3, 5) = 'hello'",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "exist_tb_2", executor.ExplainRecordAccessTypeIndex}}},
				{"SHOW GLOBAL VARIABLES", [][]string{showGlobalVariableHead, {"version", "5.7.3"}}},
			},
			nil,
			[]*OptimizeResult{{"exist_tb_2", []string{"LEFT(`v3`, 5)"}, ""}},
		},
		{
			"select * from exist_tb_2 where left(v3, 5) = 'hello'",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "exist_tb_2", executor.ExplainRecordAccessTypeIndex}}},
				{"SHOW GLOBAL VARIABLES", [][]string{showGlobalVariableHead, {"version", "8.0.14"}}},
			},
			nil,
			[]*OptimizeResult{{"exist_tb_2", []string{"LEFT(`v3`, 5)"}, ""}},
		},
		{
			"select * from exist_tb_2 where v3 like 'mike%'",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "exist_tb_2", executor.ExplainRecordAccessTypeIndex}}},
			},
			nil,
			[]*OptimizeResult{{"exist_tb_2", []string{"v3"}, ""}},
		},
		{
			"select * from exist_tb_2 where v3 like '_mike%'",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "exist_tb_2", executor.ExplainRecordAccessTypeIndex}}},
			},
			nil,
			nil,
		},
		{
			"select * from exist_tb_2 where v3 like '%mike%'",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "exist_tb_2", executor.ExplainRecordAccessTypeIndex}}},
			},
			nil,
			nil,
		},
		{
			"select * from exist_tb_2 where v3 like '%mike%' and v1 = 1",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "exist_tb_2", executor.ExplainRecordAccessTypeIndex}}},
			},
			nil,
			[]*OptimizeResult{{"exist_tb_2", []string{"v1"}, ""}},
		},

		{
			"select max(v3) from exist_tb_3",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "exist_tb_3", executor.ExplainRecordAccessTypeAll}}},
			},
			nil,
			[]*OptimizeResult{{"exist_tb_3", []string{"v3"}, ""}},
		},

		{
			"select min(v3) from exist_tb_3",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "exist_tb_3", executor.ExplainRecordAccessTypeAll}}},
			},
			nil,
			[]*OptimizeResult{{"exist_tb_3", []string{"v3"}, ""}},
		},
		{
			"select sum(v3) from exist_tb_3",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "exist_tb_3", executor.ExplainRecordAccessTypeAll}}},
			},
			nil,
			nil,
		},

		{
			"select v1, v2 from EXIST_TB_5 where v1 = '1'",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "EXIST_TB_5", executor.ExplainRecordAccessTypeAll}}},
				{"show table status", [][]string{showTableStatusHead, {"EXIST_TB_5", "10000000"}}},
			},
			nil,
			[]*OptimizeResult{{"EXIST_TB_5", []string{"v1", "v2"}, ""}},
		},
		{
			"select * from EXIST_TB_5 join exist_tb_3 on EXIST_TB_5.v1 = exist_tb_3.v1",
			[]databaseMock{
				{"EXPLAIN", [][]string{explainHead, {"1", "exist_tb_3", executor.ExplainRecordAccessTypeAll}, {"1", "EXIST_TB_5", executor.ExplainRecordAccessTypeAll}}},
			},
			nil,
			[]*OptimizeResult{{"EXIST_TB_5", []string{"v1"}, ""}},
		},
	}
	for i, tt := range optimizerTests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			ss, err := parser.New().ParseOneStmt(tt.SQL, "", "")
			assert.NoError(t, err)
			e, mocker, err := executor.NewMockExecutor()
			assert.NoError(t, err)

			for _, mock := range tt.databaseMocks {
				e := mocker.ExpectQuery(regexp.QuoteMeta(mock.expectQuery))
				rows := sqlmock.NewRows(mock.rows[0])
				for _, row := range mock.rows[1:] {
					var rowI []driver.Value
					for _, v := range row {
						rowI = append(rowI, v)
					}
					rows.AddRow(rowI...)
				}
				e.WillReturnRows(rows)
			}

			o := NewOptimizer(entry, session.NewMockContext(e), tt.optimizerOption...)
			fmt.Println("sqle:", ss)
			optimizeResults, err := o.Optimize(context.TODO(), ss.(*ast.SelectStmt))
			assert.NoError(t, err)
			assert.Equal(t, len(tt.output), len(optimizeResults))
			for i, want := range tt.output {
				assert.Equal(t, want.TableName, optimizeResults[i].TableName)
				assert.Equal(t, want.IndexedColumns, optimizeResults[i].IndexedColumns)
			}
			mocker.MatchExpectationsInOrder(true)
			assert.NoError(t, mocker.ExpectationsWereMet())
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
		{"select * from t1 as t2", map[string]string{"t2": "SELECT * FROM t1 AS t2", "t1": "SELECT * FROM t1 AS t2"}, nil},
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

func TestCanOptimize(t *testing.T) {
	logger := testLogger.WithField("test", "test_can_optimize")
	tests := []struct {
		sql    string
		expect bool
	}{
		{"select 1", false},
		{"select * from t1", false},
		{"select * from exist_tb_1", true},
		{"select * from t1, t2", false},
		{"select * from t1 join t2", false},
		{"select * from t1 cross join t2", false},
		{"select * from t1 inner join t2", false},
		{"select * from exist_tb_1, exist_tb_2", true},
		{"select * from exist_tb_1 join exist_tb_2", true},
		{"select * from exist_tb_1 cross join exist_tb_2", true},
		{"select * from exist_tb_1 inner join exist_tb_2", true},
		{"select * from t1, exist_tb_2", false},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			e, _, err := executor.NewMockExecutor()
			assert.NoError(t, err)
			n, err := util.ParseOneSql(tt.sql)
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, CanOptimize(logger, session.NewMockContext(e), n))
		})
	}
}
