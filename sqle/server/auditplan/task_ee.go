//go:build enterprise
// +build enterprise

package auditplan

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server"

	"github.com/percona/go-mysql/query"
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

	inst, _, err := at.persist.GetInstanceByNameAndProjectID(at.ap.InstanceName, at.ap.ProjectId)
	if err != nil {
		at.logger.Warnf("get instance fail, error: %v", err)
		return
	}

	if !driver.GetPluginManager().IsOptionalModuleEnabled(inst.DbType, driverV2.OptionalModuleQuery) {
		at.logger.Warnf("can not do this task, %v", driver.NewErrPluginAPINotImplement(driverV2.OptionalModuleQuery))
		return
	}

	plugin, err := driver.GetPluginManager().OpenPlugin(at.logger, inst.DbType, &driverV2.Config{
		DSN: &driverV2.DSN{
			Host:             inst.Host,
			Port:             inst.Port,
			User:             inst.User,
			Password:         inst.Password,
			AdditionalParams: inst.AdditionalParams,
		},
	})
	if err != nil {
		at.logger.Warnf("get plugin failed, error: %v", err)
		return
	}
	defer plugin.Close(context.Background())

	sql := at.getCollectSQL()
	if sql == "" {
		at.logger.Warnf("unknown metric of interest")
		return
	}
	err = at.collect(plugin, sql)
	if err != nil {
		at.logger.Warnf("collect failed, error: %v", err)
		return
	}
}

func (at *OBMySQLTopSQLTask) collect(p driver.Plugin, sql string) error {
	result, err := p.Query(context.Background(), sql, &driverV2.QueryConf{TimeOutSecond: 5})
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

	return at.persist.OverrideAuditPlanSQLs(at.ap.ID, convertSQLsToModelSQLs(sqls))
}

func (at *OBMySQLTopSQLTask) Audit() (*model.AuditPlanReportV2, error) {
	var task *model.Task
	if at.ap.InstanceName == "" {
		task = &model.Task{
			DBType: at.ap.DBType,
		}
	} else {
		instance, _, err := at.persist.GetInstanceByNameAndProjectID(at.ap.InstanceName, at.ap.ProjectId)
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

type SlowLogTask struct {
	*sqlCollector
}

func NewSlowLogTask(entry *logrus.Entry, ap *model.AuditPlan) Task {
	t := &SlowLogTask{newSQLCollector(entry, ap)}
	t.do = t.collectorDo

	return t
}

const (
	slowlogCollectInputLogFile = 0 // FILE: mysql-slow.log
	slowlogCollectInputTable   = 1 // TABLE: mysql.slow_log
)

func (at *SlowLogTask) collectorDo() {

	if at.ap.Params.GetParam(paramKeySlowLogCollectInput).Int() != slowlogCollectInputTable {
		return
	}

	if at.ap.InstanceName == "" {
		at.logger.Warnf("instance is not configured")
		return
	}

	instance, _, err := at.persist.GetInstanceByNameAndProjectID(at.ap.InstanceName, at.ap.ProjectId)
	if err != nil {
		return
	}

	db, err := executor.NewExecutor(at.logger, &driverV2.DSN{
		Host:             instance.Host,
		Port:             instance.Port,
		User:             instance.User,
		Password:         instance.Password,
		AdditionalParams: instance.AdditionalParams,
		DatabaseName:     at.ap.InstanceDatabase,
	},
		at.ap.InstanceDatabase)
	if err != nil {
		at.logger.Errorf("connect to instance fail, error: %v", err)
		return
	}
	defer db.Db.Close()

	res, err := db.Db.Query(`
SELECT sql_text,db,TIME_TO_SEC(query_time) AS query_time
FROM mysql.slow_log
WHERE sql_text != ''
    AND db NOT IN ('information_schema','performance_schema','mysql','sys')
`)

	if err != nil {
		at.logger.Errorf("query slow log failed, error: %v", err)
		return
	}

	if len(res) == 0 {
		return
	}

	sqls := make([]*sqlFromSlowLog, len(res))

	for i := range res {
		sqls[i] = &sqlFromSlowLog{
			sql:    res[i]["sql_text"].String,
			schema: res[i]["db"].String,
		}
		queryTime, err := strconv.Atoi(res[i]["query_time"].String)
		if err != nil {
			at.logger.Warnf("unexpected data format: %v, ", res[i]["query_time"].String)
			continue
		}
		sqls[i].queryTimeSeconds = queryTime
	}

	sqlFingerprintInfos := sqlFromSlowLogs(sqls).mergeByFingerprint()

	auditPlanSQLs := make([]*model.AuditPlanSQLV2, len(sqlFingerprintInfos))
	{
		now := time.Now()
		for i := range sqlFingerprintInfos {
			fp := sqlFingerprintInfos[i]
			fpInfo := fmt.Sprintf(`{"counter":%v,"last_receive_timestamp":"%v","schema":"%v","average_query_time":"%d"}`,
				fp.counter, now.Format(time.RFC3339), fp.schema, fp.queryTimeSeconds)
			auditPlanSQLs[i] = &model.AuditPlanSQLV2{
				Fingerprint: fp.fingerprint,
				SQLContent:  fp.sql,
				Info:        []byte(fpInfo),
			}
		}
	}

	if err = at.persist.OverrideAuditPlanSQLs(at.ap.ID, auditPlanSQLs); err != nil {
		at.logger.Errorf("save mysql slow log to storage failed, error: %v", err)
		return
	}
}

type sqlFromSlowLog struct {
	sql              string
	schema           string
	queryTimeSeconds int
}

type sqlFromSlowLogs []*sqlFromSlowLog

type sqlFingerprintInfo struct {
	lastSql               string
	lastSqlSchema         string
	sqlCount              int
	totalQueryTimeSeconds int
}

func (s *sqlFingerprintInfo) queryTime() int {
	return s.totalQueryTimeSeconds / s.sqlCount
}

func (s sqlFromSlowLogs) mergeByFingerprint() []sqlInfo {

	sqlInfos := []sqlInfo{}
	sqlInfosMap := map[string] /*sql fingerprint*/ *sqlFingerprintInfo{}

	for i := range s {
		sqlItem := s[i]
		fp := query.Fingerprint(sqlItem.sql)
		if sqlInfosMap[fp] != nil {
			sqlInfosMap[fp].lastSql = sqlItem.sql
			sqlInfosMap[fp].lastSqlSchema = sqlItem.schema
			sqlInfosMap[fp].sqlCount++
			sqlInfosMap[fp].totalQueryTimeSeconds += sqlItem.queryTimeSeconds
		} else {
			sqlInfosMap[fp] = &sqlFingerprintInfo{
				sqlCount:              1,
				lastSql:               sqlItem.sql,
				lastSqlSchema:         sqlItem.schema,
				totalQueryTimeSeconds: sqlItem.queryTimeSeconds,
			}
			sqlInfos = append(sqlInfos, sqlInfo{fingerprint: fp})
		}
	}

	for i := range sqlInfos {
		fp := sqlInfos[i].fingerprint
		sqlInfo := sqlInfosMap[fp]
		if sqlInfo != nil {
			sqlInfos[i].counter = sqlInfo.sqlCount
			sqlInfos[i].sql = sqlInfo.lastSql
			sqlInfos[i].schema = sqlInfo.lastSqlSchema
			sqlInfos[i].queryTimeSeconds = sqlInfo.queryTime()
		}

	}

	return sqlInfos
}

func (at *SlowLogTask) GetSQLs(args map[string]interface{}) (
	[]Head, []map[string] /* head name */ string, uint64, error) {

	auditPlanSQLs, count, err := at.persist.GetAuditPlanSQLsByReq(args)
	if err != nil {
		return nil, nil, count, err
	}
	head := []Head{
		{
			Name: "fingerprint",
			Desc: "SQL指纹",
			Type: "sql",
		},
		{
			Name: "sql",
			Desc: "最后一次匹配到该指纹的语句",
			Type: "sql",
		},
		{
			Name: "counter",
			Desc: "匹配到该指纹的语句数量",
		},
		{
			Name: "last_receive_timestamp",
			Desc: "最后一次匹配到该指纹的时间",
		},
		{
			Name: "average_query_time",
			Desc: "平均执行时间（秒）",
		},
	}
	rows := make([]map[string]string, 0, len(auditPlanSQLs))
	for _, sql := range auditPlanSQLs {
		var info = struct {
			Counter              uint64 `json:"counter"`
			LastReceiveTimestamp string `json:"last_receive_timestamp"`
			AverageQueryTime     string `json:"average_query_time"`
		}{}
		err := json.Unmarshal(sql.Info, &info)
		if err != nil {
			return nil, nil, 0, err
		}
		rows = append(rows, map[string]string{
			"sql":                    sql.SQLContent,
			"fingerprint":            sql.Fingerprint,
			"counter":                strconv.FormatUint(info.Counter, 10),
			"last_receive_timestamp": info.LastReceiveTimestamp,
			"average_query_time":     info.AverageQueryTime,
		})
	}
	return head, rows, count, nil
}

// HACK: slow SQLs may be executed in different Schemas.
// Before auditing sql, we need to insert a Schema switching statement.
// And need to manually execute server.ReplenishTaskStatistics() to recalculate
// real task object score
func (at *SlowLogTask) Audit() (*model.AuditPlanReportV2, error) {
	return auditWithSchema(at.logger, at.persist, at.ap)
}
