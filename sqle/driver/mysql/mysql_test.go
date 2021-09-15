package mysql

import (
	"context"
	"testing"

	"github.com/actiontech/sqle/sqle/model"
	"github.com/stretchr/testify/assert"
)

func TestInspect_Parse(t *testing.T) {
	nodes, err := DefaultMysqlInspect().Parse(context.TODO(), `
use test_db;
create trigger my_trigger before insert on t1 for each row insert into t2(id, c1) values(1, '2');
create table t1(id int);
	`)
	assert.NoError(t, err)
	for _, node := range nodes {
		assert.Equal(t, node.Type, model.SQLTypeDDL)
	}

	nodes, err = DefaultMysqlInspect().Parse(context.TODO(), "select * from t1")
	assert.NoError(t, err)
	assert.Len(t, nodes, 1)
	assert.Equal(t, nodes[0].Type, model.SQLTypeDML)
}
