package mysql

import (
	"database/sql/driver"
	"fmt"

	"actiontech.cloud/sqle/sqle/sqle/model"

	"github.com/labstack/gommon/log"
)

func (i *Inspect) CommitDDL(sql *model.BaseSQL) error {
	conn, err := i.getDbConn()
	if err != nil {
		return err
	}
	_, err = conn.Db.Exec(sql.Content)
	if err != nil {
		sql.ExecStatus = model.SQLExecuteStatusFailed
		sql.ExecResult = err.Error()
	} else {
		sql.ExecStatus = model.SQLExecuteStatusSucceeded
		sql.ExecResult = "ok"
	}
	return nil
}

func (i *Inspect) CommitDMLs(sqls []*model.BaseSQL) error {
	var retErr error
	var conn *Executor
	var startBinlogFile, endBinlogFile string
	var startBinlogPos, endBinlogPos int64
	var results []driver.Result
	qs := make([]string, 0, len(sqls))
	for _, sql := range sqls {
		qs = append(qs, sql.Content)
	}

	defer func() {
		for idx, sql := range sqls {
			if retErr != nil {
				sql.ExecStatus = model.SQLExecuteStatusFailed
				sql.ExecResult = retErr.Error()
				continue
			}
			sql.StartBinlogFile = startBinlogFile
			sql.StartBinlogPos = startBinlogPos
			rowAffects, err := results[idx].RowsAffected()
			if err != nil {
				log.Warnf("get rows affect failed, error: %v", err)
			}
			sql.RowAffects = rowAffects
			sql.ExecStatus = model.SQLExecuteStatusSucceeded
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
