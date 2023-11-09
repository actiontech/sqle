package index

import (
	"log"
	"testing"

	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/stretchr/testify/assert"
)

func TestTraverse(t *testing.T) {

	sql := `
		SELECT a.column1, MAX(DISTINCT b.column2),MIN(a.column1),current_timestamp,ABS(a.column2)
		FROM table1 a
		JOIN table2 b ON a.common_column = b.common_column
		WHERE a.column3 LIKE '%some_value%' AND a.column2 LIKE '_any_value'
		GROUP BY a.column1;
	`
	// sql := `SELECT *
	// FROM table1 a
	// JOIN table2 b ON a.id = b.id JOIN (SELECT * FROM employees) e
	// WHERE a.city LIKE '%some_value%';`

	node, err := util.ParseOneSql(sql)
	assert.NoError(t, err)
	visitor := &IndexOptimizer{
		CurrentTableName: "table1",
		OriginNode:       node,
	}
	node.Accept(visitor)
	for _, result := range visitor.Results {
		log.Println(*result)
	}
}
