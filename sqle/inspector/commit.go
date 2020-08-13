package inspector

import (
	"database/sql/driver"
	"fmt"
	"actiontech.cloud/universe/sqle/v3/sqle/executor"
	"actiontech.cloud/universe/sqle/v3/sqle/model"

	"github.com/labstack/gommon/log"
	"github.com/pingcap/parser/ast"
)

func (i *Inspect) CommitDDL(sql *model.Sql) error {
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

func (i *Inspect) CommitDMLs(sqls []*model.Sql) error {
	if i.Task.Instance.DbType == model.DB_TYPE_MYCAT {
		return i.commitMycatDMLs(sqls)
	}

	var retErr error
	var conn *executor.Executor
	var startBinlogFile, endBinlogFile string
	var startBinlogPos, endBinlogPos int64
	var results []driver.Result

	qs := []string{}
	sqlToQsIdxes := make([][]int, len(sqls), len(sqls))
	qsIndex := 0
	for sqlIdx, sql := range sqls {
		qsIdxes := []int{}
		for _, stmt := range sql.Stmts {
			qs = append(qs, stmt.Text())
			qsIdxes = append(qsIdxes, qsIndex)
			qsIndex += 1
		}
		sqlToQsIdxes[sqlIdx] = qsIdxes
	}
	defer func() {
		for sqlIdx, sql := range sqls {
			if retErr != nil {
				sql.ExecStatus = model.TASK_ACTION_ERROR
				sql.ExecResult = retErr.Error()
				continue
			}
			sql.StartBinlogFile = startBinlogFile
			sql.StartBinlogPos = startBinlogPos
			for _, qsIdx := range sqlToQsIdxes[sqlIdx] {
				rowAffects, err := results[qsIdx].RowsAffected()
				if err != nil {
					log.Warnf("get rows affect failed, error: %v", err)
					continue
				}
				sql.RowAffects += rowAffects
			}
			sql.ExecStatus = model.TASK_ACTION_DONE
			sql.ExecResult = "ok"
			sql.EndBinlogFile = endBinlogFile
			sql.EndBinlogPos = endBinlogPos
		}
	}()

	conn, retErr = i.getDbConn()
	if retErr != nil {
		return retErr
	}

	startBinlogFile, startBinlogPos, retErr = conn.FetchMasterBinlogPos()
	if retErr != nil {
		return retErr
	}

	results, err := conn.Db.Transact(qs...)
	if err != nil {
		retErr = err
	} else if len(results) != len(qs) {
		retErr = fmt.Errorf("number of transaction result does not match number of sqls")
	} else {
		endBinlogFile, endBinlogPos, _ = conn.FetchMasterBinlogPos()
	}
	return retErr
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

func (i *Inspect) commitMycatDMLs(sqls []*model.Sql) error {
	conn, err := i.getDbConn()
	if err != nil {
		for _, sql := range sqls {
			sql.ExecStatus = model.TASK_ACTION_ERROR
			sql.ExecResult = err.Error()
		}
		return err
	}

	for _, sql := range sqls {
		startBinlogFile, startBinlogPos, err := conn.FetchMasterBinlogPos()
		if err != nil {
			sql.ExecStatus = model.TASK_ACTION_ERROR
			sql.ExecResult = err.Error()
			return err
		}
		sql.StartBinlogFile = startBinlogFile
		sql.StartBinlogPos = startBinlogPos

		for _, stmt := range sql.Stmts {
			query := stmt.Text()
			result, err := conn.Db.Exec(query)
			if err != nil {
				sql.ExecStatus = model.TASK_ACTION_ERROR
				sql.ExecResult = err.Error()
				return err
			}
			rowAffect, err := result.RowsAffected()
			if err != nil {
				log.Warnf("get rows affect failed, error: %v", err)
			} else {
				sql.RowAffects += rowAffect
			}
		}
		sql.ExecStatus = model.TASK_ACTION_DONE
		sql.ExecStatus = "ok"
		sql.EndBinlogFile, sql.EndBinlogPos, _ = conn.FetchMasterBinlogPos()
	}
	return nil
}
