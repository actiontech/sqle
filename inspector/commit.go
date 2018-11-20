package inspector

import (
	"database/sql/driver"
	"fmt"
	"sqle/model"
)

func (i *Inspector) Commit() error {
	err := i.Advise()
	if err != nil {
		return err
	}
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
	for _, sql := range i.SqlArray {
		_, err := conn.Db.Exec(sql.Sql)
		if err != nil {
			sql.ExecStatus = model.TASK_ACTION_ERROR
			sql.ExecResult = err.Error()
		} else {
			sql.ExecStatus = model.TASK_ACTION_DONE
			sql.ExecResult = "ok"
		}
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
