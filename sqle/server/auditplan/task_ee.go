//go:build enterprise
// +build enterprise

package auditplan

import (
	"context"
	"fmt"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/sirupsen/logrus"
)

const (
	OBMySQLIndicatorCPUTime     = "cpu_time"
	OBMySQLIndicatorIOWait      = "io_wait"
	OBMySQLIndicatorElapsedTime = "elapsed_time"
)

type OBMySQLTopSQLTask struct {
	*sqlCollector
}

func NewOBMySQLTopSQLTask(entry *logrus.Entry, ap *model.AuditPlan) Task {
	task := &OBMySQLTopSQLTask{
		sqlCollector: newSQLCollector(entry, ap),
	}
	task.sqlCollector.do = task.collectorDo
	return task
}

func (at *OBMySQLTopSQLTask) collectorDo() {
	select {
	case <-at.cancel:
		at.logger.Info("cancel task")
		return
	default:
	}

	if at.ap.InstanceName == "" {
		at.logger.Warnf("instance is not configured")
		return
	}

	inst, _, err := at.persist.GetInstanceByName(at.ap.InstanceName)
	if err != nil {
		at.logger.Warnf("get instance fail, error: %v", err)
		return
	}

	m, err := driver.NewDriverManger(at.logger, at.ap.DBType, &driver.Config{
		DSN: &driver.DSN{
			Host:             inst.Host,
			Port:             inst.Port,
			User:             inst.User,
			Password:         inst.Password,
			AdditionalParams: inst.AdditionalParams,
		},
	})

	if err != nil {
		at.logger.Warnf("get driver manager failed")
		return
	}
	defer m.Close(context.Background())

	queryDriver, err := m.GetSQLQueryDriver()
	if err != nil {
		at.logger.Warnf("get sql query driver failed")
		return
	}

	sql := at.getCollectSQL()
	if sql == "" {
		at.logger.Warnf("unknown metric of interest")
		return
	}
	err = at.collect(queryDriver, sql)
	if err != nil {
		at.logger.Warnf("collect failed, error: %v", err)
		return
	}
}

func (at *OBMySQLTopSQLTask) collect(queryDriver driver.SQLQueryDriver, sql string) error {
	result, err := queryDriver.Query(context.Background(), sql, &driver.QueryConf{TimeOutSecond: 5})
	if err != nil {
		return err
	}
	if len(result.Column) <= 0 {
		return nil
	}

	sqlTextIndex := 0
	for i, param := range result.Column {
		if param.String() == OBMySQLSQLKeySQLText {
			sqlTextIndex = i
			break
		}
	}

	sqls := []*SQL{}
	for _, row := range result.Rows {
		s := &SQL{
			Info: map[string]interface{}{},
		}
		for i, value := range row.Values {
			if i == sqlTextIndex {
				s.SQLContent = value.Value
				s.Fingerprint = value.Value
			} else {
				s.Info[result.Column[i].String()] = value.Value
			}
		}
		sqls = append(sqls, s)
	}

	return at.persist.OverrideAuditPlanSQLs(at.ap.Name, convertSQLsToModelSQLs(sqls))
}

func (at *OBMySQLTopSQLTask) Audit() (*model.AuditPlanReportV2, error) {
	var task *model.Task
	if at.ap.InstanceName == "" {
		task = &model.Task{
			DBType: at.ap.DBType,
		}
	} else {
		instance, _, err := at.persist.GetInstanceByName(at.ap.InstanceName)
		if err != nil {
			return nil, err
		}
		task = &model.Task{
			Instance: instance,
			Schema:   at.ap.InstanceDatabase,
			DBType:   at.ap.DBType,
		}
	}
	return at.baseTask.audit(task)
}

const (
	// 通用采集项
	OBMySQLSQLKeySQLText            = "sql_text"
	OBMySQLSQLInfoKeyFirstRequest   = "first_request"
	OBMySQLSQLInfoKeyExecCount      = "exec_count"
	OBMySQLSQLInfoKeyLastRequest    = "last_request"
	OBMySQLSQLInfoKeyAverageElapsed = "average_elapsed"

	// OBMySQLIndicatorElapsedTime 对应采集项
	OBMySQLSQLInfoKeyMaxElapsed = "max_elapsed"

	// OBMySQLIndicatorCPUTime 对应采集项
	OBMySQLSQLInfoKeyAverageCPU = "average_cpu"

	// OBMySQLIndicatorIOWait 对应采集项
	OBMySQLSQLInfoKeyAverageIOWait = "average_io_wait"
	OBMySQLSQLInfoKeyDiskRead      = "disk_read"
	OBMySQLSQLInfoKeyBufferRead    = "buffer_read"
)

func (at *OBMySQLTopSQLTask) GetSQLs(args map[string]interface{}) ([]Head, []map[string] /* head name */ string, uint64, error) {
	auditPlanSQLs, count, err := at.persist.GetAuditPlanSQLsByReq(args)
	if err != nil {
		return nil, nil, 0, err
	}
	result := []map[string]string{}
	for _, planSQL := range auditPlanSQLs {
		mp := map[string]string{
			OBMySQLSQLKeySQLText: planSQL.SQLContent,
		}

		origin, err := planSQL.Info.OriginValue()
		if err != nil {
			return nil, nil, 0, err
		}
		for k, v := range origin {
			mp[k] = fmt.Sprintf("%v", v)
		}
		result = append(result, mp)
	}
	return at.getHead(), result, count, nil
}

func (at *OBMySQLTopSQLTask) getCollectSQL() string {
	topN := at.ap.Params.GetParam(paramKeyTopN).Int()

	switch at.ap.Params.GetParam(paramKeyIndicator).String() {
	case OBMySQLIndicatorElapsedTime:
		return fmt.Sprintf(`
SELECT
    SQL_TEXT AS %v, 
    EXECUTIONS AS %v, 
    CEIL(AVG_EXE_USEC/1000) AS %v, 
    CEIL(SLOWEST_EXE_USEC/1000) AS %v, 
    FROM_UNIXTIME(TIME_TO_USEC(FIRST_LOAD_TIME)/1000000) AS %v,
    FROM_UNIXTIME(TIME_TO_USEC(LAST_ACTIVE_TIME)/1000000) AS %v
FROM
    OCEANBASE.GV$SQL
GROUP BY
    SQL_ID
ORDER BY
    %v
DESC
LIMIT %v
`, OBMySQLSQLKeySQLText,
			OBMySQLSQLInfoKeyExecCount,
			OBMySQLSQLInfoKeyAverageElapsed,
			OBMySQLSQLInfoKeyMaxElapsed,
			OBMySQLSQLInfoKeyFirstRequest,
			OBMySQLSQLInfoKeyLastRequest,
			OBMySQLSQLInfoKeyMaxElapsed,
			topN)

	case OBMySQLIndicatorCPUTime:
		return fmt.Sprintf(`
SELECT
    SQL_TEXT AS %v, 
    EXECUTIONS AS %v, 
    CEIL(AVG_EXE_USEC/1000) AS %v,
    CEIL(CPU_TIME/EXECUTIONS/1000) AS %v, 
    FROM_UNIXTIME(TIME_TO_USEC(FIRST_LOAD_TIME)/1000000) AS %v,
    FROM_UNIXTIME(TIME_TO_USEC(LAST_ACTIVE_TIME)/1000000) AS %v
FROM
    OCEANBASE.GV$SQL
GROUP BY
    SQL_ID
ORDER BY
    %v
DESC
LIMIT %v
`, OBMySQLSQLKeySQLText,
			OBMySQLSQLInfoKeyExecCount,
			OBMySQLSQLInfoKeyAverageElapsed,
			OBMySQLSQLInfoKeyAverageCPU,
			OBMySQLSQLInfoKeyFirstRequest,
			OBMySQLSQLInfoKeyLastRequest,
			OBMySQLSQLInfoKeyAverageCPU,
			topN,
		)

	case OBMySQLIndicatorIOWait:
		return fmt.Sprintf(`
SELECT
    SQL_TEXT AS %v, 
    EXECUTIONS AS %v, 
    CEIL(USER_IO_WAIT_TIME/EXECUTIONS/1000) AS %v, 
    CEIL(BUFFER_GETS/EXECUTIONS) AS %v,
    CEIL(DISK_READS/EXECUTIONS) AS %v,
    FROM_UNIXTIME(TIME_TO_USEC(FIRST_LOAD_TIME)/1000000) AS %v,
    FROM_UNIXTIME(TIME_TO_USEC(LAST_ACTIVE_TIME)/1000000) AS %v
FROM
    OCEANBASE.GV$SQL
GROUP BY
    SQL_ID
ORDER BY
    %v
DESC
LIMIT %v
`, OBMySQLSQLKeySQLText,
			OBMySQLSQLInfoKeyExecCount,
			OBMySQLSQLInfoKeyAverageIOWait,
			OBMySQLSQLInfoKeyBufferRead,
			OBMySQLSQLInfoKeyDiskRead,
			OBMySQLSQLInfoKeyFirstRequest,
			OBMySQLSQLInfoKeyLastRequest,
			OBMySQLSQLInfoKeyAverageIOWait,
			topN,
		)

	default:
		return ""
	}
}

func (at *OBMySQLTopSQLTask) getHead() []Head {
	switch at.ap.Params.GetParam(paramKeyIndicator).String() {
	case OBMySQLIndicatorElapsedTime:
		return []Head{
			{
				Name: OBMySQLSQLKeySQLText,
				Desc: "SQL指纹",
				Type: "sql",
			}, {
				Name: OBMySQLSQLInfoKeyExecCount,
				Desc: "执行次数",
			}, {
				Name: OBMySQLSQLInfoKeyAverageElapsed,
				Desc: "平均执行时间(毫秒)",
			}, {
				Name: OBMySQLSQLInfoKeyMaxElapsed,
				Desc: "最长执行时间(毫秒)",
			}, {
				Name: OBMySQLSQLInfoKeyFirstRequest,
				Desc: "首次执行时间",
			}, {
				Name: OBMySQLSQLInfoKeyLastRequest,
				Desc: "最后执行时间",
			},
		}
	case OBMySQLIndicatorIOWait:
		return []Head{
			{
				Name: OBMySQLSQLKeySQLText,
				Desc: "SQL指纹",
				Type: "sql",
			}, {
				Name: OBMySQLSQLInfoKeyExecCount,
				Desc: "执行次数",
			}, {
				Name: OBMySQLSQLInfoKeyAverageIOWait,
				Desc: "平均IO等待时间(毫秒)",
			}, {
				Name: OBMySQLSQLInfoKeyBufferRead,
				Desc: "平均逻辑读次数",
			}, {
				Name: OBMySQLSQLInfoKeyDiskRead,
				Desc: "平均物理读次数",
			}, {
				Name: OBMySQLSQLInfoKeyFirstRequest,
				Desc: "首次执行时间",
			}, {
				Name: OBMySQLSQLInfoKeyLastRequest,
				Desc: "最后执行时间",
			},
		}
	case OBMySQLIndicatorCPUTime:
		return []Head{
			{
				Name: OBMySQLSQLKeySQLText,
				Desc: "SQL指纹",
				Type: "sql",
			}, {
				Name: OBMySQLSQLInfoKeyExecCount,
				Desc: "执行次数",
			}, {
				Name: OBMySQLSQLInfoKeyAverageCPU,
				Desc: "平均CPU时间(毫秒)",
			}, {
				Name: OBMySQLSQLInfoKeyAverageElapsed,
				Desc: "SQL平均执行时间(毫秒)",
			}, {
				Name: OBMySQLSQLInfoKeyFirstRequest,
				Desc: "首次执行时间",
			}, {
				Name: OBMySQLSQLInfoKeyLastRequest,
				Desc: "最后执行时间",
			},
		}
	}
	return []Head{}
}
