package inspector

import (
	"database/sql/driver"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/pingcap/tidb/ast"
	"sqle/model"
)

func (i *Inspect) CommitAll() error {
	for _, commitSql := range i.Task.CommitSqls {
		currentSql := commitSql
		err := i.Add(&currentSql.Sql, func(sql *model.Sql) error {
			err := i.Commit(sql)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return i.Do()
}

func (i *Inspect) RollbackAll(sql *model.RollbackSql) error {
	for _, rollbackSql := range i.Task.RollbackSqls {
		currentSql := rollbackSql
		err := i.Add(&currentSql.Sql, func(sql *model.Sql) error {
			err := i.Commit(sql)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return i.Do()
}

func (i *Inspect) Commit(sql *model.Sql) error {
	if i.SqlType() == model.SQL_TYPE_DDL {
		return i.commitDDL(sql)
	} else {
		return i.commitDML(sql)
	}
}

func (i *Inspect) commitDDL(sql *model.Sql) error {
	conn, err := i.getDbConn()
	if err != nil {
		return err
	}
	if i.Task.Instance.DbType == model.DB_TYPE_MYCAT {
		return i.commitMycatDDL(sql)
	}
	_, err = conn.Db.Exec(sql.Content)
	if err != nil {
		sql.ExecStatus = model.TASK_ACTION_ERROR
		sql.ExecResult = err.Error()
	} else {
		sql.ExecStatus = model.TASK_ACTION_DONE
		sql.ExecResult = "ok"
	}
	return nil
}

func (i *Inspect) commitDML(sql *model.Sql) error {
	var err error
	var result driver.Result
	var rowAffect int64
	var qs []string

	conn, err := i.getDbConn()
	if err != nil {
		return err
	}

	sql.StartBinlogFile, sql.StartBinlogPos, err = conn.FetchMasterBinlogPos()
	if err != nil {
		goto ERROR
	}

	qs = make([]string, 0, len(sql.Stmts))
	for _, stmt := range sql.Stmts {
		qs = append(qs, stmt.Text())
	}

	if len(qs) > 1 && i.Task.Instance.DbType != model.DB_TYPE_MYCAT {
		err = conn.Db.Transact(qs...)
		if err != nil {
			goto ERROR
		}
	} else {
		for _, query := range qs {
			result, err = conn.Db.Exec(query)
			if err != nil {
				goto ERROR
			}
			rowAffect, err = result.RowsAffected()
			if err != nil {
				log.Warnf("get rows affect failed, error: %s", err)
			} else {
				sql.RowAffects += rowAffect
			}
		}
	}

	sql.ExecStatus = model.TASK_ACTION_DONE
	sql.ExecResult = "ok"
	// if sql has commit success, ignore error for check status.
	sql.EndBinlogFile, sql.EndBinlogPos, _ = conn.FetchMasterBinlogPos()
	return nil
ERROR:
	sql.ExecStatus = model.TASK_ACTION_ERROR
	sql.ExecResult = err.Error()
	return err
}

func (i *Inspect) commitMycatDDL(sql *model.Sql) error {
	conn, err := i.getDbConn()
	if err != nil {
		return err
	}
	var schemaName string
	var tableName string

	stmt := sql.Stmts[0]
	query := stmt.Text()
	switch stmt := stmt.(type) {
	case *ast.CreateTableStmt:
		schemaName = stmt.Table.Schema.String()
		tableName = stmt.Table.Name.String()
	case *ast.AlterTableStmt:
		schemaName = stmt.Table.Schema.String()
		tableName = stmt.Table.Name.String()
	case *ast.DropTableStmt:
		if stmt.Tables == nil || len(stmt.Tables) == 0 {
			goto DONE
		}
		if len(stmt.Tables) > 1 {
			err = fmt.Errorf("don't support multi drop table in diff schema on mycat")
			goto DONE
		}
		table := stmt.Tables[0]
		schemaName = table.Schema.String()
		tableName = table.Name.String()
	case *ast.CreateIndexStmt:
		schemaName = stmt.Table.Schema.String()
		tableName = stmt.Table.Name.String()
	case *ast.DropIndexStmt:
		schemaName = stmt.Table.Schema.String()
		tableName = stmt.Table.Name.String()
	case *ast.DropDatabaseStmt:
		err = fmt.Errorf("don't support drop database on mycat")
		goto DONE
	case *ast.UseStmt:
		goto DONE
	default:
	}
	if schemaName != "" {
		query = replaceTableName(query, schemaName, tableName)
	}
	// if no schema name in table name, use default schema name
	if schemaName == "" {
		schemaName = i.Ctx.currentSchema
	}
	err = conn.Db.ExecDDL(query, schemaName, tableName)

DONE:
	if err != nil {
		sql.ExecStatus = model.TASK_ACTION_ERROR
		sql.ExecResult = err.Error()
	} else {
		sql.ExecStatus = model.TASK_ACTION_DONE
		sql.ExecResult = "ok"
	}
	return nil
}
