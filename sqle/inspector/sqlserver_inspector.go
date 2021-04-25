package inspector

import (
	"actiontech.cloud/sqle/sqle/sqle/errors"
	"actiontech.cloud/sqle/sqle/sqle/model"
	"actiontech.cloud/sqle/sqle/sqle/sqlserver/SqlserverProto"
	"actiontech.cloud/sqle/sqle/sqle/sqlserverClient"
	"fmt"

	"github.com/pingcap/parser/ast"
	"github.com/sirupsen/logrus"
)

// Inspect implements Inspector interface for SQL Server.
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

func (i *SqlserverInspect) ParseSqlType() error {
	for _, commitSql := range i.Task.ExecuteSQLs {
		nodes, err := i.ParseSql(commitSql.Content)
		if err != nil {
			return err
		}
		i.addNodeCounter(nodes)
	}
	return nil
}

func (i *SqlserverInspect) addNodeCounter(nodes []ast.Node) {
	for _, node := range nodes {
		switch stmt := node.(type) {
		case sqlserverClient.SqlServerNode:
			if stmt.IsDDLStmt() {
				i.counterDDL += 1
			} else if stmt.IsDMLStmt() {
				i.counterDML += 1
			} else if stmt.IsProcedureStmt() {
				i.counterProcedure += 1
			} else if stmt.IsFunctionStmt() {
				i.counterFunction += 1
			}
		}
	}
}

func (i *SqlserverInspect) Add(sql *model.BaseSQL, action func(sql *model.BaseSQL) error) error {
	nodes, err := i.ParseSql(sql.Content)
	if err != nil {
		return err
	}
	i.addNodeCounter(nodes)

	sql.Stmts = nodes
	i.SqlArray = append(i.SqlArray, sql)
	i.SqlAction = append(i.SqlAction, action)
	return nil
}

func (i *SqlserverInspect) Advise(rules []model.Rule) error {
	i.Logger().Info("start advise sql")

	sqls := []string{}
	for _, commitSql := range i.Task.ExecuteSQLs {
		sqls = append(sqls, commitSql.Content)
	}
	ruleNames := []string{}
	for _, rule := range rules {
		ruleNames = append(ruleNames, rule.Name)
	}
	meta := sqlserverClient.GetSqlserverMeta(i.Task.Instance.User, i.Task.Instance.Password, i.Task.Instance.Host, i.Task.Instance.Port, i.Task.Schema, "")
	ddlContextSqls := []*SqlserverProto.DDLContext{}
	for _, task := range i.RelateTasks {
		sqls := []string{}
		for _, commitSql := range task.ExecuteSQLs {
			sqls = append(sqls, commitSql.Content)
		}
		ddlContextSql := &SqlserverProto.DDLContext{
			Sqls: sqls,
		}
		ddlContextSqls = append(ddlContextSqls, ddlContextSql)
	}

	out, err := sqlserverClient.GetClient().Advise(sqls, ruleNames, meta, ddlContextSqls)
	if err != nil {
		i.Logger().Errorf("advise t-sql from ms grpc server failed, error: %v", err)
		return err
	} else {
		i.Logger().Info("advise sql finish")
	}

	results := out.GetResults()
	for idx, commitSql := range i.Task.ExecuteSQLs {
		result := results[idx]
		stmt := sqlserverClient.NewSqlServerStmt(commitSql.Content, result.IsDDL, result.IsDML, result.IsProcedure, result.IsFunction)
		if stmt.IsDDLStmt() {
			i.counterDDL += 1
		} else if stmt.IsDMLStmt() {
			i.counterDML += 1
		} else if stmt.IsProcedureStmt() {
			i.counterProcedure += 1
		} else if stmt.IsFunctionStmt() {
			i.counterFunction += 1
		}
		commitSql.AuditLevel = result.AdviseLevel
		commitSql.AuditResult = result.AdviseResultMessage
		commitSql.AuditStatus = model.SQLAuditStatusFinished
	}
	i.HasInvalidSql = out.BaseValidatorFailed

	if i.SqlType() == model.SQL_TYPE_MULTI {
		return errors.SQL_STMT_CONFLICT_ERROR
	}
	if i.SqlType() == model.SQL_TYPE_PROCEDURE_FUNCTION_MULTI {
		return errors.SQL_STMT_PROCEUDRE_FUNCTION_ERROR
	}

	return err
}

func (i *SqlserverInspect) GenerateAllRollbackSql() ([]*model.RollbackSQL, error) {
	i.Logger().Info("start generate rollback sql")

	var meta = sqlserverClient.GetSqlserverMeta(i.Task.Instance.User, i.Task.Instance.Password, i.Task.Instance.Host, i.Task.Instance.Port, i.Task.Schema, "")
	sqls, err := sqlserverClient.GetClient().GenerateAllRollbackSql(i.Task.ExecuteSQLs, &SqlserverProto.Config{DMLRollbackMaxRows: i.config.DMLRollbackMaxRows}, meta)
	if err != nil {
		i.Logger().Errorf("generage t-sql rollback sqls error: %v", err)
	} else {
		i.Logger().Info("generate rollback sql finish")
	}
	if err != nil {
		return nil, err
	}

	// update reason of no rollback sql
	if len(i.Task.ExecuteSQLs) != len(sqls) {
		return nil, fmt.Errorf("don't match sql rollback result")
	}
	rollbackSqls := []*model.RollbackSQL{}
	for idx, val := range sqls {
		executeSQL := i.Task.ExecuteSQLs[idx]
		if val.Sql != "" {
			rollbackSqls = append(rollbackSqls, &model.RollbackSQL{
				BaseSQL: model.BaseSQL{
					Content: val.Sql,
				},
				ExecuteSQLId: executeSQL.ID,
			})
		}
		if val.ErrMsg != "" {
			result := newInspectResults()
			if executeSQL.AuditResult != "" {
				result.add(executeSQL.AuditLevel, executeSQL.AuditResult)
			}
			result.add(model.RULE_LEVEL_NOTICE, val.ErrMsg)
			executeSQL.AuditLevel = result.level()
			executeSQL.AuditResult = result.message()
		}
	}

	return i.Inspect.GetAllRollbackSqlReversed(rollbackSqls), err
}

func (i *SqlserverInspect) GetProcedureFunctionBackupSql(sql string) ([]string, error) {
	i.Logger().Info("start get procedure & function backup sql")

	var meta = sqlserverClient.GetSqlserverMeta(i.Task.Instance.User, i.Task.Instance.Password, i.Task.Instance.Host, i.Task.Instance.Port, i.Task.Schema, "")
	backupSql, err := sqlserverClient.GetClient().GetProcedureFunctionBackupSql(sql, meta)
	if err != nil {
		i.Logger().Errorf("get procedure & function backup sql error: %v", err)
		return nil, err
	}
	i.Logger().Info("get procedure & function backup sql finish")

	return backupSql, nil
}
