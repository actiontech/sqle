package inspector

import (
	"github.com/pingcap/tidb/ast"
	"github.com/sirupsen/logrus"
	"sqle/model"
	"sqle/sqlserver/SqlserverProto"
	"sqle/sqlserverClient"
)

type SqlserverInspect struct {
	*Inspect
}

func NeSqlserverInspect(entry *logrus.Entry, ctx *Context, task *model.Task, relateTasks []model.Task,
	rules map[string]model.Rule) Inspector {
	return &SqlserverInspect{
		Inspect: NewInspect(entry, ctx, task, relateTasks, rules),
	}
}

func (i *SqlserverInspect) ParseSql(sql string) ([]ast.Node, error) {
	sqls, err := sqlserverClient.GetClient().ParseSql(sql)
	if err != nil {
		i.Logger().Errorf("parser t-sql from ms grpc server failed, error: %v", err)
	}
	return sqls, err
}

func (i *SqlserverInspect) Add(sql *model.Sql, action func(sql *model.Sql) error) error {
	nodes, err := i.ParseSql(sql.Content)
	if err != nil {
		return err
	}
	for _, node := range nodes {
		switch stmt := node.(type) {
		case sqlserverClient.SqlServerNode:
			if stmt.IsDDLStmt() {
				i.counterDDL += 1
			} else if stmt.IsDMLStmt() {
				i.counterDML += 1
			}
		}
	}

	sql.Stmts = nodes
	i.SqlArray = append(i.SqlArray, sql)
	i.SqlAction = append(i.SqlAction, action)
	return nil
}

func (i *SqlserverInspect) Advise(rules []model.Rule) error {
	i.Logger().Info("start advise sql")
	var meta = sqlserverClient.GetSqlserverMeta(i.Task.Instance.User, i.Task.Instance.Password, i.Task.Instance.Host, i.Task.Instance.Port, i.Task.Schema, "")
	err := sqlserverClient.GetClient().Advise(i.RelateTasks, i.Task.CommitSqls, rules, meta)
	if err != nil {
		i.Logger().Errorf("advise t-sql from ms grpc server failed, error: %v", err)
	} else {
		i.Logger().Info("advise sql finish")
	}
	return err
}

func (i *SqlserverInspect) GenerateAllRollbackSql() ([]*model.RollbackSql, error) {
	i.Logger().Info("start generate rollback sql")

	var meta = sqlserverClient.GetSqlserverMeta(i.Task.Instance.User, i.Task.Instance.Password, i.Task.Instance.Host, i.Task.Instance.Port, i.Task.Schema, "")
	rollbackSqls, err := sqlserverClient.GetClient().GenerateAllRollbackSql(i.Task.CommitSqls, &SqlserverProto.Config{DMLRollbackMaxRows: i.config.DMLRollbackMaxRows}, meta)
	if err != nil {
		i.Logger().Errorf("generage t-sql rollback sqls error: %v", err)
	} else {
		i.Logger().Info("generate rollback sql finish")
	}
	return i.Inspect.GetAllRollbackSql(rollbackSqls), err
}
