package mysql

import (
	"testing"

	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/stretchr/testify/assert"
)

// TestCTESQLTypeThreeStates covers AC-006: sql_type follows outer CTE semantics.
// CTE+DML uses DELETE (UPDATE likewise); CTE+INSERT remains Unparsed in current parser and is out of scope.
func TestCTESQLTypeThreeStates(t *testing.T) {
	cases := []struct {
		name string
		sql  string
		want string
	}{
		{
			name: "cte_plus_select_dql",
			sql:  `WITH cte (id) AS (SELECT 1) SELECT * FROM cte`,
			want: driverV2.SQLTypeDQL,
		},
		{
			name: "cte_plus_delete_dml",
			sql:  `WITH cte AS (SELECT id FROM exist_tb_1 WHERE id = 1) DELETE FROM exist_tb_1 WHERE id IN (SELECT id FROM cte)`,
			want: driverV2.SQLTypeDML,
		},
		{
			name: "real_ddl",
			sql:  `CREATE TABLE t_ac006 (id INT)`,
			want: driverV2.SQLTypeDDL,
		},
	}

	i := &MysqlDriverImpl{}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			stmts, err := util.ParseSql(c.sql)
			assert.NoError(t, err)
			if !assert.NotEmpty(t, stmts) {
				return
			}
			assert.Equal(t, c.want, i.assertSQLType(stmts[0]))
		})
	}
}
