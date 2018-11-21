package inspector

import (
	"database/sql/driver"
	"fmt"
	"github.com/pingcap/tidb/ast"
	"sqle/model"
)

func (i *Inspector) Commit() error {
	err := i.prepare()
	if err != nil {
		return err
	}
	defer i.closeDbConn()
	if i.isDMLStmt {
		return i.commitDML()
	}
	return i.commitDDL()
}

func (i *Inspector) commitDDL() error {
	conn, err := i.getDbConn()
	if err != nil {
		return err
	}
	for n, node := range i.SqlStmt {
		sql := i.SqlArray[n]
		var schema string
		var table string

		switch stmt:=node.(type) {
		case *ast.CreateTableStmt:
			schema = i.getSchemaName(stmt.Table)
			table = stmt.Table.Name.String()
		case *ast.AlterTableStmt:
			schema = i.getSchemaName(stmt.Table)
			table = stmt.Table.Name.String()
		default:

		}
		err := conn.Db.ExecDDL(sql.Sql,schema,table)
		if err != nil {
			sql.ExecStatus = model.TASK_ACTION_ERROR
			sql.ExecResult = err.Error()
		} else {
			sql.ExecStatus = model.TASK_ACTION_DONE
			sql.ExecResult = "ok"
		}
		i.updateSchemaCtx(node)
	}
	return nil
}

func (i *Inspector) commitDML() error {
	conn, err := i.getDbConn()
	if err != nil {
		return err
	}
	for _, sql := range i.SqlArray {
		var result driver.Result
		var err error
		var a int64
		sql.StartBinlogFile, sql.StartBinlogPos, err = conn.FetchMasterBinlogPos()
		if err != nil {
			goto ERROR
		}
		result, err = conn.Db.Exec(sql.Sql)
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
		continue

	ERROR:
		sql.ExecStatus = model.TASK_ACTION_ERROR
		sql.ExecResult = err.Error()
	}
	return nil
}
