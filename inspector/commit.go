package inspector

import (
	"database/sql/driver"
	"fmt"
	"github.com/pingcap/tidb/ast"
	"sqle/model"
)

func (i *Inspector) CommitAll() error {
	err := i.Prepare()
	if err != nil {
		return err
	}
	defer i.closeDbConn()
	for i.NextCommitSql() {
		err := i.Commit()
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Inspector) Commit() error {
	if i.isDMLStmt {
		return i.commitDML()
	} else {
		return i.commitDDL()
	}
}

func (i *Inspector) commitDDL() error {
	conn, err := i.getDbConn()
	if err != nil {
		return err
	}
	if i.Instance.DbType == model.DB_TYPE_MYCAT {
		return i.commitMycatDDL()
	}
	sql := i.GetCommitSql()
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

func (i *Inspector) commitDML() error {
	conn, err := i.getDbConn()
	if err != nil {
		return err
	}
	sql := i.GetCommitSql()
	var result driver.Result
	var a int64

	sql.StartBinlogFile, sql.StartBinlogPos, err = conn.FetchMasterBinlogPos()
	if err != nil {
		goto ERROR
	}
	result, err = conn.Db.Exec(sql.Content)
	if err != nil {
		goto ERROR
	}
	a, err = result.RowsAffected()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("row_affect: ", a)
	}

	sql.RowAffects, _ = result.RowsAffected()
	sql.ExecStatus = model.TASK_ACTION_DONE
	sql.ExecResult = "ok"
	// if sql has commit success, ignore error for check status.
	sql.EndBinlogFile, sql.EndBinlogPos, _ = conn.FetchMasterBinlogPos()
ERROR:
	sql.ExecStatus = model.TASK_ACTION_ERROR
	sql.ExecResult = err.Error()
	return err
}

func (i *Inspector) commitMycatDDL() error {
	conn, err := i.getDbConn()
	if err != nil {
		return err
	}
	node := i.GetSqlStmt()
	sql := i.GetCommitSql()
	i.updateSchemaCtx(node)

	var schema string
	var table string

	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		schema = i.getSchemaName(stmt.Table)
		table = stmt.Table.Name.String()
	case *ast.AlterTableStmt:
		schema = i.getSchemaName(stmt.Table)
		table = stmt.Table.Name.String()
	case *ast.CreateIndexStmt:
		schema = i.getSchemaName(stmt.Table)
		table = stmt.Table.Name.String()
	case *ast.UseStmt:
		goto DONE
	default:
	}
	err = conn.Db.ExecDDL(ReplaceTableName(node), schema, table)

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

//func (i *Inspector) Rollback(sql *model.RollbackSql) error {
//	nodes, err := parseSql(i.Instance.DbType, sql.Sql)
//	if err != nil {
//		return err
//	}
//	for _, node := range nodes {
//		switch node.(type) {
//		case ast.DDLNode:
//
//		}
//	}
//}
