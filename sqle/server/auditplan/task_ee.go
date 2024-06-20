//go:build enterprise
// +build enterprise

package auditplan

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/common"
	"github.com/actiontech/sqle/sqle/log"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"
	dry "github.com/ungerik/go-dry"

	"github.com/actiontech/sqle/sqle/driver/mysql/util"

	"github.com/sirupsen/logrus"
)

const (
	OBMySQLIndicatorCPUTime     = "cpu_time"
	OBMySQLIndicatorIOWait      = "io_wait"
	OBMySQLIndicatorElapsedTime = "elapsed_time"
	SlowLogQueryNums            = 1000
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

	inst, _, err := dms.GetInstanceInProjectByName(context.Background(), string(at.ap.ProjectId), at.ap.InstanceName)
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

func (at *OBMySQLTopSQLTask) Audit() (*AuditResultResp, error) {
	task, err := getTaskWithInstanceByAuditPlan(at.ap, at.persist)
	if err != nil {
		return nil, err
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
WHERE %v != ''
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
			OBMySQLSQLKeySQLText,
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
WHERE %v != ''
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
			OBMySQLSQLKeySQLText,
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
WHERE %v != ''
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
			OBMySQLSQLKeySQLText,
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

func (at *SlowLogTask) FullSyncSQLs(sqls []*SQL) error {
	if at.ap.Params.GetParam(paramKeySlowLogCollectInput).Int() == slowlogCollectInputTable {
		return at.sqlCollector.FullSyncSQLs(sqls)
	}
	return at.baseTask.FullSyncSQLs(sqls)
}

func (at *SlowLogTask) PartialSyncSQLs(sqls []*SQL) error {
	if at.ap.Params.GetParam(paramKeySlowLogCollectInput).Int() == slowlogCollectInputTable {
		return at.sqlCollector.PartialSyncSQLs(sqls)
	}
	return at.persist.UpdateSlowLogAuditPlanSQLs(at.ap.ID, convertSQLsToModelSQLs(sqls))
}

func (at *SlowLogTask) collectorDo() {

	if at.ap.Params.GetParam(paramKeySlowLogCollectInput).Int() != slowlogCollectInputTable {
		return
	}

	if at.ap.InstanceName == "" {
		at.logger.Warnf("instance is not configured")
		return
	}

	instance, _, err := dms.GetInstanceInProjectByName(context.Background(), string(at.ap.ProjectId), at.ap.InstanceName)
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

	queryStartTime, err := at.persist.GetLatestStartTimeAuditPlanSQL(at.ap.ID)
	if err != nil {
		at.logger.Errorf("get start time failed, error: %v", err)
		return
	}
	querySQL := `
	SELECT sql_text,db,TIME_TO_SEC(query_time) AS query_time, start_time, rows_examined
	FROM mysql.slow_log
	WHERE sql_text != ''
		AND db NOT IN ('information_schema','performance_schema','mysql','sys')
	`

	sqlInfos := []*sqlFromSlowLog{}

	for {
		extraCondition := fmt.Sprintf(" AND start_time>'%s' ORDER BY start_time LIMIT %d", queryStartTime, SlowLogQueryNums)
		execQuerySQL := querySQL + extraCondition

		res, err := db.Db.Query(execQuerySQL)

		if err != nil {
			at.logger.Errorf("query slow log failed, error: %v", err)
			break
		}

		for i := range res {
			sqlInfo := &sqlFromSlowLog{
				sql:       res[i]["sql_text"].String,
				schema:    res[i]["db"].String,
				startTime: res[i]["start_time"].String,
			}
			queryTime, err := strconv.Atoi(res[i]["query_time"].String)
			if err != nil {
				at.logger.Warnf("unexpected data format: %v, ", res[i]["query_time"].String)
				continue
			}
			sqlInfo.queryTimeSeconds = queryTime
			sqlInfo.rowExamined, err = strconv.ParseFloat(res[i]["rows_examined"].String, 64)
			if err != nil {
				at.logger.Warnf("unexpected data format: %v, ", res[i]["rows_examined"].String)
				continue
			}

			sqlInfos = append(sqlInfos, sqlInfo)
		}

		if len(res) < SlowLogQueryNums {
			break
		}

		queryStartTime = res[len(res)-1]["start_time"].String

		time.Sleep(500 * time.Millisecond)
	}

	if len(sqlInfos) == 0 {
		return
	}
	sqlFingerprintInfos, err := sqlFromSlowLogs(sqlInfos).mergeByFingerprint(at.logger)
	if err != nil {
		at.logger.Errorf("merge finger sqls failed, error: %v", err)
		return
	}

	auditPlanSQLs := make([]*model.AuditPlanSQLV2, len(sqlFingerprintInfos))
	{
		now := time.Now()
		for i := range sqlFingerprintInfos {
			fp := sqlFingerprintInfos[i]
			fpInfo := fmt.Sprintf(`{"counter":%v,"last_receive_timestamp":"%v","schema":"%v","average_query_time":%d,"start_time":"%v","row_examined_avg":%v}`,
				fp.counter, now.Format(time.RFC3339), fp.schema, fp.queryTimeSeconds, fp.startTime, fp.rowExaminedAvg)
			auditPlanSQLs[i] = &model.AuditPlanSQLV2{
				Fingerprint: fp.fingerprint,
				SQLContent:  fp.sql,
				Info:        []byte(fpInfo),
				Schema:      fp.schema,
			}
		}
	}

	if err = at.persist.UpdateSlowLogCollectAuditPlanSQLs(at.ap.ID, auditPlanSQLs); err != nil {
		at.logger.Errorf("save mysql slow log to storage failed, error: %v", err)
		return
	}
}

type sqlFromSlowLog struct {
	sql              string
	schema           string
	queryTimeSeconds int
	startTime        string
	rowExamined      float64
}

type sqlFromSlowLogs []*sqlFromSlowLog

type sqlFingerprintInfo struct {
	lastSql               string
	lastSqlSchema         string
	sqlCount              int
	totalQueryTimeSeconds int
	startTime             string
	totalExaminedRows     float64
}

func (s *sqlFingerprintInfo) queryTime() int {
	return s.totalQueryTimeSeconds / s.sqlCount
}

func (s *sqlFingerprintInfo) rowExaminedAvg() float64 {
	return s.totalExaminedRows / float64(s.sqlCount)
}

func (s sqlFromSlowLogs) mergeByFingerprint(entry *logrus.Entry) ([]sqlInfo, error) {

	sqlInfos := []sqlInfo{}
	sqlInfosMap := map[string] /*sql fingerprint*/ *sqlFingerprintInfo{}

	for i := range s {
		sqlItem := s[i]
		fp, err := util.Fingerprint(sqlItem.sql, true)
		if err != nil {
			entry.Warnf("get sql finger print failed, err: %v", err)
		}
		if fp == "" {
			continue
		}

		if sqlInfosMap[fp] != nil {
			sqlInfosMap[fp].lastSql = sqlItem.sql
			sqlInfosMap[fp].lastSqlSchema = sqlItem.schema
			sqlInfosMap[fp].sqlCount++
			sqlInfosMap[fp].totalQueryTimeSeconds += sqlItem.queryTimeSeconds
			sqlInfosMap[fp].startTime = sqlItem.startTime
			sqlInfosMap[fp].totalExaminedRows += sqlItem.rowExamined
		} else {
			sqlInfosMap[fp] = &sqlFingerprintInfo{
				sqlCount:              1,
				lastSql:               sqlItem.sql,
				lastSqlSchema:         sqlItem.schema,
				totalQueryTimeSeconds: sqlItem.queryTimeSeconds,
				startTime:             sqlItem.startTime,
				totalExaminedRows:     sqlItem.rowExamined,
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
			sqlInfos[i].startTime = sqlInfo.startTime
			sqlInfos[i].rowExaminedAvg = utils.Round(sqlInfo.rowExaminedAvg(), 6)
		}
	}

	return sqlInfos, nil
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
			Desc: "SQL",
			Type: "sql",
		},
		{
			Name: "counter",
			Desc: "数量",
		},
		{
			Name: "last_receive_timestamp",
			Desc: "最后匹配时间",
		},
		{
			Name: "average_query_time",
			Desc: "平均执行时间",
		},
		{
			Name: "max_query_time",
			Desc: "最长执行时间",
		},
		{
			Name: "row_examined_avg",
			Desc: "平均扫描行数",
		},
		{
			Name: "db_user",
			Desc: "用户",
		},
		{
			Name: "schema",
			Desc: "Schema",
		},
	}
	rows := make([]map[string]string, 0, len(auditPlanSQLs))
	for _, sql := range auditPlanSQLs {
		var info = struct {
			Counter              uint64   `json:"counter"`
			LastReceiveTimestamp string   `json:"last_receive_timestamp"`
			AverageQueryTime     *float64 `json:"query_time_avg"`
			MaxQueryTime         *float64 `json:"query_time_max"`
			RowExaminedAvg       *float64 `json:"row_examined_avg"`
			DBUser               string   `json:"db_user"`
		}{}
		err := json.Unmarshal(sql.Info, &info)
		if err != nil {
			return nil, nil, 0, err
		}
		row := map[string]string{
			"sql":                    sql.SQLContent,
			"fingerprint":            sql.Fingerprint,
			"counter":                strconv.FormatUint(info.Counter, 10),
			"last_receive_timestamp": info.LastReceiveTimestamp,
			"db_user":                info.DBUser,
			"schema":                 sql.Schema,
		}

		if info.RowExaminedAvg != nil {
			row["row_examined_avg"] = fmt.Sprintf("%.6f", *info.RowExaminedAvg)
		}
		// 兼容之前没有平均执行时间和最长执行时间的数据，没有数据的时候不会在前端显示0.00000导致误解
		if info.AverageQueryTime != nil {
			row["average_query_time"] = fmt.Sprintf("%.6f", *info.AverageQueryTime)
		}
		if info.MaxQueryTime != nil {
			row["max_query_time"] = fmt.Sprintf("%.6f", *info.MaxQueryTime)
		}
		rows = append(rows, row)
	}
	return head, rows, count, nil
}

// HACK: slow SQLs may be executed in different Schemas.
// Before auditing sql, we need to insert a Schema switching statement.
// And need to manually execute server.ReplenishTaskStatistics() to recalculate
// real task object score
func (at *SlowLogTask) Audit() (*AuditResultResp, error) {
	return auditWithSchema(at.logger, at.persist, at.ap)
}

type DB2TopSQLTask struct {
	*sqlCollector
}

func NewDB2TopSQLTask(entry *logrus.Entry, ap *model.AuditPlan) Task {
	task := &DB2TopSQLTask{
		sqlCollector: newSQLCollector(entry, ap),
	}
	task.do = task.collectorDo
	return task
}

func (at *DB2TopSQLTask) Audit() (*AuditResultResp, error) {
	task, err := getTaskWithInstanceByAuditPlan(at.ap, at.persist)
	if err != nil {
		return nil, err
	}
	return at.baseTask.audit(task)
}

func (at *DB2TopSQLTask) indicator() (string, error) {
	indicator := at.ap.Params.GetParam(paramKeyIndicator).String()
	if indicator == "" {
		return DB2IndicatorAverageElapsedTime, nil
	}

	if !dry.StringInSlice(indicator, []string{
		DB2IndicatorNumExecutions,
		DB2IndicatorTotalElapsedTime,
		DB2IndicatorAverageElapsedTime,
		DB2IndicatorAverageCPUTime,
		DB2IndicatorLockWaitTime,
		DB2IndicatorLockWaitNum,
		DB2IndicatorSQLWaitTime,
		DB2IndicatorTotalActTime,
	}) {
		return "", fmt.Errorf("invalid indicator: %v", indicator)
	}
	return indicator, nil
}

// ref: https://www.ibm.com/docs/zh/db2/11.1?topic=views-snap-get-dyn-sql-dynsql-snapshot
func (at *DB2TopSQLTask) collectSQL() (string, error) {
	// SET TENANT ?是用于设置当前会话的租户标识的语句
	// MON_GET_PKG_CACHE_STMT表函数中可能会存在两次`SET TENANT ?`。因为指纹唯一性导致存表失败，所以过滤对`SET TENANT ?`语句的采集
	sql := `
	SELECT
	num_executions,
	total_act_time,
	total_act_wait_time,
	lock_wait_time,
	lock_waits,
	stmt_exec_time,
	total_cpu_time / NUM_EXEC_WITH_METRICS AS avg_cpu_time,
	STMT_EXEC_TIME/NUM_EXECUTIONS AS avg_elapsed_time_ms,
	substr(stmt_text, 1, 300) AS stmt_text
	FROM
	TABLE(MON_GET_PKG_CACHE_STMT(NULL, NULL, NULL, -2)) T
	WHERE
	upper(stmt_text) NOT LIKE '%%MON_GET_PKG_CACHE_STMT%%' AND NUM_EXEC_WITH_METRICS <> 0 AND upper(stmt_text) <> 'SET TENANT ?'
ORDER BY %s DESC   
`
	indicator, err := at.indicator()
	if err != nil {
		return "", err
	}

	sql = fmt.Sprintf(sql, indicator)

	// limit top N
	{
		topN := at.ap.Params.GetParam(paramKeyTopN).Int()
		if topN == 0 {
			topN = 10
		}
		sql = fmt.Sprintf(`%v FETCH FIRST %d ROWS ONLY `, sql, topN)
	}

	return sql, nil
}

const (
	DB2IndicatorNumExecutions      = "num_executions"      // 执行次数
	DB2IndicatorTotalElapsedTime   = "stmt_exec_time"      // 总执行时间
	DB2IndicatorAverageElapsedTime = "avg_elapsed_time_ms" // 平均执行时间
	DB2IndicatorAverageCPUTime     = "avg_cpu_time"        // 平均 CPU 时间
	DB2IndicatorLockWaitTime       = "lock_wait_time"      // 锁等待时间
	DB2IndicatorLockWaitNum        = "lock_waits"          // 锁等待次数
	DB2IndicatorSQLWaitTime        = "total_act_wait_time" // 活动等待总时间
	DB2IndicatorTotalActTime       = "total_act_time"      // 活动总时间
)

func (at *DB2TopSQLTask) collectorDo() {

	if at.ap.InstanceName == "" {
		at.logger.Warn("instance is not configured")
		return
	}

	inst, _, err := dms.GetInstanceInProjectByName(context.Background(), string(at.ap.ProjectId), at.ap.InstanceName)
	if err != nil {
		at.logger.Warnf("get instance fail, error: %v", err)
		return
	}

	if !driver.GetPluginManager().
		IsOptionalModuleEnabled(inst.DbType, driverV2.OptionalModuleQuery) {
		at.logger.Warnf("can not do this task, %v",
			driver.NewErrPluginAPINotImplement(driverV2.OptionalModuleQuery))
		return
	}

	plugin, err := driver.GetPluginManager().OpenPlugin(
		at.logger, inst.DbType, &driverV2.Config{
			DSN: &driverV2.DSN{
				Host:             inst.Host,
				Port:             inst.Port,
				User:             inst.User,
				Password:         inst.Password,
				AdditionalParams: inst.AdditionalParams,
				DatabaseName:     at.ap.InstanceDatabase,
			},
		})
	if err != nil {
		at.logger.Warnf("get plugin failed, error: %v", err)
		return
	}
	defer plugin.Close(context.Background())

	sql, err := at.collectSQL()
	if err != nil {
		at.logger.Warnf("generate collect sql failed, error: %v", err)
		return
	}

	result, err := plugin.Query(context.Background(), sql,
		&driverV2.QueryConf{TimeOutSecond: 10})
	if err != nil {
		at.logger.Warnf("collect failed, error: %v", err)
		return
	}

	if len(result.Rows) == 0 {
		at.logger.Infof("db2 top sql audit_plan(%v) collected no statement", at.ap.ID)
		return
	}

	at.logger.Infof("db2 top sql audit_plan(%v) collected %v statements", at.ap.ID, len(result.Rows))

	sqls := []*SQL{}
	for i := range result.Rows {
		row := result.Rows[i]
		s := &SQL{Info: make(map[string]interface{}, 0)}
		for j := range row.Values {
			if strings.ToLower(result.Column[j].Key) == "stmt_text" {
				s.SQLContent = row.Values[j].Value
				s.Fingerprint = row.Values[j].Value
			} else {
				s.Info[strings.ToLower(result.Column[j].Key)] = row.Values[j].Value
				s.Info["last_receive_timestamp"] = time.Now().Format(time.RFC3339)
			}
		}
		sqls = append(sqls, s)
	}

	if err := at.persist.OverrideAuditPlanSQLs(at.ap.ID, convertSQLsToModelSQLs(sqls)); err != nil {
		at.logger.Errorf("save db2 top sql to storage failed, error: %v", err)
		return
	}
	return
}

func (at *DB2TopSQLTask) getSQLHead() []Head {
	return []Head{
		{
			Name: "sql",
			Desc: "SQL语句",
			Type: "sql",
		},
		{
			Name: DB2IndicatorTotalElapsedTime,
			Desc: "总执行时间(ms)",
		},
		{
			Name: DB2IndicatorAverageElapsedTime,
			Desc: "平均执行时间(ms)",
		},
		{
			Name: DB2IndicatorNumExecutions,
			Desc: "执行次数",
		},
		{
			Name: DB2IndicatorAverageCPUTime,
			Desc: "平均 CPU 时间(μs)",
		},
		{
			Name: DB2IndicatorLockWaitTime,
			Desc: "锁等待时间(ms)",
		},
		{
			Name: DB2IndicatorLockWaitNum,
			Desc: "锁等待次数",
		},
		{
			Name: DB2IndicatorSQLWaitTime,
			Desc: "活动等待总时间(ms)",
		},
		{
			Name: DB2IndicatorTotalActTime,
			Desc: "活动总时间(ms)",
		},
	}
}

func (at *DB2TopSQLTask) GetSQLs(args map[string]interface{}) ([]Head, []map[string] /* head name */ string, uint64, error) {
	auditPlanSQLs, count, err := at.persist.GetAuditPlanSQLsByReq(args)
	if err != nil {
		return nil, nil, 0, err
	}
	result := []map[string]string{}

	for i := range auditPlanSQLs {
		mp := map[string]string{"sql": auditPlanSQLs[i].SQLContent}

		origin, err := auditPlanSQLs[i].Info.OriginValue()
		if err != nil {
			return nil, nil, 0, err
		}
		for k := range origin {
			mp[k] = fmt.Sprintf("%v", origin[k])
		}
		result = append(result, mp)
	}
	return at.getSQLHead(), result, count, nil
}

type DB2SchemaMetaTask struct {
	*sqlCollector
}

func NewDB2SchemaMetaTask(entry *logrus.Entry, ap *model.AuditPlan) Task {
	sqlCollector := newSQLCollector(entry, ap)
	task := &DB2SchemaMetaTask{
		sqlCollector,
	}
	sqlCollector.do = task.collectorDo
	return task
}

func (at *DB2SchemaMetaTask) Audit() (*AuditResultResp, error) {
	task, err := getTaskWithInstanceByAuditPlan(at.ap, at.persist)
	if err != nil {
		return nil, err
	}
	return at.baseTask.audit(task)
}

func (at *DB2SchemaMetaTask) GetSQLs(args map[string]interface{}) ([]Head, []map[string] /* head name */ string, uint64, error) {
	auditPlanSQLs, count, err := at.persist.GetAuditPlanSQLsByReq(args)
	if err != nil {
		return nil, nil, count, err
	}
	head, rows := buildSchemaMetaSQLsResult(auditPlanSQLs)
	return head, rows, count, nil
}

func (at *DB2SchemaMetaTask) isSchemaValid(plugin driver.Plugin) (bool, error) {
	schemasFromInst, err := plugin.Schemas(context.Background())
	if err != nil {
		return false, fmt.Errorf("get schemas from db2 failed, error: %v", err)
	}
	return utils.StringsContains(schemasFromInst, at.ap.InstanceDatabase), nil
}

func (at *DB2SchemaMetaTask) collectorDo() {
	if at.ap.InstanceName == "" {
		at.logger.Warnf("instance is not configured")
		return
	}
	if at.ap.InstanceDatabase == "" {
		at.logger.Warnf("instance schema is not configured")
		return
	}
	instance, _, err := dms.GetInstanceInProjectByName(context.Background(), string(at.ap.ProjectId), at.ap.InstanceName)
	if err != nil {
		return
	}

	pluginMgr := driver.GetPluginManager()
	if !pluginMgr.IsOptionalModuleEnabled(instance.DbType, driverV2.OptionalModuleQuery) {
		at.logger.Errorf("collect DB2 schema meta failed: %v", driver.NewErrPluginAPINotImplement(driverV2.OptionalModuleQuery))
		return
	}
	plugin, err := pluginMgr.OpenPlugin(at.logger, instance.DbType, &driverV2.Config{DSN: &driverV2.DSN{
		Host:             instance.Host,
		Port:             instance.Port,
		User:             instance.User,
		Password:         instance.Password,
		AdditionalParams: instance.AdditionalParams,
		DatabaseName:     at.ap.InstanceDatabase,
	}})
	if err != nil {
		at.logger.Errorf("connect to instance fail, error: %v", err)
		return
	}
	tempVariableName := fmt.Sprintf("%v.sqle_get_ddl_token", at.ap.InstanceDatabase)
	valIsCreated := false
	defer func() {
		if valIsCreated {
			_, err = plugin.Exec(context.Background(), fmt.Sprintf(`DROP VARIABLE %v`, tempVariableName))
			if err != nil {
				at.logger.Errorf("drop variable failed, error: %v", err)
			}
		}

		plugin.Close(context.Background())
	}()

	if valid, err := at.isSchemaValid(plugin); err != nil {
		at.logger.Errorf("check schema failed: %v", err)
		return
	} else if !valid {
		at.logger.Errorf("schema [%v] doesn't exist in db2 instance", at.ap.InstanceDatabase)
		return
	}

	tables, err := at.getTablesFromSchema(context.Background(), plugin, at.ap.InstanceDatabase)
	if err != nil {
		at.logger.Errorf("get tables from schema [%v] failed, error: %v", at.ap.InstanceDatabase, err)
		return
	}

	var views []string
	if at.ap.Params.GetParam("collect_view").Bool() {
		views, err = at.getViewsFromSchema(context.Background(), plugin, at.ap.InstanceDatabase)
		if err != nil {
			at.logger.Errorf("get views from schema [%v] failed, error: %v", at.ap.InstanceDatabase, err)
			return
		}
	}

	var sqls []string
	_, err = plugin.Exec(context.Background(), fmt.Sprintf(`CREATE OR REPLACE VARIABLE %v integer`, tempVariableName))
	if err != nil {
		at.logger.Errorf("create variable failed, error: %v", err)
		return
	}
	valIsCreated = true
	for _, table := range tables {
		sql := fmt.Sprintf(`CALL SYSPROC.DB2LK_GENERATE_DDL('-t %v.%v -e',%v)`, at.ap.InstanceDatabase, table, tempVariableName)
		_, err = plugin.Exec(context.Background(), sql)
		if err != nil {
			at.logger.Errorf("generate ddl failed, sql: %s, error: %v", sql, err)
			continue
		}
		result, err := plugin.Query(context.Background(), fmt.Sprintf(`
SELECT VARCHAR(SQL_STMT,2000) AS CREATE_TABLE_DDL FROM SYSTOOLS.DB2LOOK_INFO WHERE OP_TOKEN = %v AND OBJ_TYPE = 'TABLE' ORDER BY OP_SEQUENCE ASC
`, tempVariableName), &driverV2.QueryConf{TimeOutSecond: 10})
		if err != nil {
			at.logger.Errorf("get create table ddl for table [%v] failed, error: %v", table, err)
			continue
		}

		if len(result.Column) != 1 || result.Column[0].Key != "CREATE_TABLE_DDL" || len(result.Rows) != 1 {
			at.logger.Errorf("parse create table ddl records  for table [%v] failed", table)
			continue
		}
		createTableDDL := ""
		for _, value := range result.Rows[0].Values {
			createTableDDL = value.Value
		}
		if createTableDDL == "" {
			at.logger.Errorf("get empty create table ddl for table %v.%v", at.ap.InstanceDatabase, table)
			continue
		}
		sqls = append(sqls, createTableDDL)
	}

	for _, view := range views {
		sql := fmt.Sprintf(`call SYSPROC.DB2LK_GENERATE_DDL('-v %v.%v -e',%v)`, at.ap.InstanceDatabase, view, tempVariableName)
		_, err = plugin.Exec(context.Background(), sql)
		if err != nil {

			at.logger.Errorf("generate ddl failed, sql: %s, error: %v", sql, err)
			continue
		}
		result, err := plugin.Query(context.Background(), fmt.Sprintf(`
SELECT VARCHAR(SQL_STMT,2000) AS CREATE_VIEW_DDL FROM SYSTOOLS.DB2LOOK_INFO WHERE OP_TOKEN = %v AND OBJ_TYPE = 'VIEW' ORDER BY OP_SEQUENCE ASC
`, tempVariableName), &driverV2.QueryConf{TimeOutSecond: 10})
		if err != nil {
			at.logger.Errorf("get create view ddl for view [%v] failed, error: %v", view, err)
			continue
		}

		if len(result.Column) != 1 || result.Column[0].Key != "CREATE_VIEW_DDL" || len(result.Rows) != 1 {
			at.logger.Errorf("parse create view ddl records  for view [%v] failed", view)
			continue
		}
		createViewDDL := ""
		for _, value := range result.Rows[0].Values {
			createViewDDL = value.Value
		}
		if createViewDDL == "" {
			at.logger.Errorf("get empty create view ddl for view %v.%v", at.ap.InstanceDatabase, view)
			continue
		}
		sqls = append(sqls, createViewDDL)
	}
	if len(sqls) > 0 {
		err = at.persist.OverrideAuditPlanSQLs(at.ap.ID, convertRawSQLToModelSQLs(sqls, at.ap.InstanceDatabase))
		if err != nil {
			at.logger.Errorf("save schema meta to storage fail, error: %v", err)
		}
	}
}

func (at *DB2SchemaMetaTask) getTablesFromSchema(ctx context.Context, plugin driver.Plugin, schema string) (tables []string, err error) {
	res, err := plugin.Query(ctx, fmt.Sprintf(`SELECT TABNAME FROM SYSCAT.TABLES WHERE TABSCHEMA = '%v' AND TYPE = 'T'`, schema), &driverV2.QueryConf{TimeOutSecond: 10})
	if err != nil {
		return nil, fmt.Errorf("query sql failed: %v", err)
	}
	if len(res.Column) != 1 || res.Column[0].Key != "TABNAME" {
		return nil, fmt.Errorf("parse query results failed")
	}
	for _, row := range res.Rows {
		for _, value := range row.Values {
			tables = append(tables, value.Value)
			continue
		}
	}
	return tables, nil
}

func (at *DB2SchemaMetaTask) getViewsFromSchema(ctx context.Context, plugin driver.Plugin, schema string) (views []string, err error) {
	res, err := plugin.Query(ctx, fmt.Sprintf(`SELECT VIEWNAME FROM SYSCAT.VIEWS WHERE VIEWSCHEMA = '%v'`, schema), &driverV2.QueryConf{TimeOutSecond: 10})
	if err != nil {
		return nil, fmt.Errorf("query sql failed: %v", err)
	}
	if len(res.Column) != 1 || res.Column[0].Key != "VIEWNAME" {
		return nil, fmt.Errorf("parse query results failed")
	}
	for _, row := range res.Rows {
		for _, value := range row.Values {
			views = append(views, value.Value)
			continue
		}
	}
	return views, nil
}

type DynPerformanceDmColumns struct {
	SQLFullText      string  `json:"sql_fulltext"`
	Executions       float64 `json:"executions"`
	TotalExecTime    float64 `json:"total_exec_time"`
	AverageExecTime  float64 `json:"average_exec_time"`
	CPUTime          float64 `json:"cpu_time"`
	PhyReadPageCnt   float64 `json:"phy_read_page_cnt"`
	LogicReadPageCnt float64 `json:"logic_read_page_cnt"`
}

// Dm Top SQL
const (
	DynPerformanceViewDmTpl = `
SELECT
    sql_fulltext,
    executions,
    total_exec_time,
    average_exec_time,
    cpu_time,
    phy_read_page_cnt,
    logic_read_page_cnt
FROM (
    SELECT
        SQL_TXT AS sql_fulltext,
        COUNT(*) AS executions,
        SUM(EXEC_TIME) AS total_exec_time,
        SUM(EXEC_TIME) / COUNT(*) OVER () AS average_exec_time,
        (SUM(EXEC_TIME) - SUM(PARSE_TIME) - SUM(IO_WAIT_TIME)) AS cpu_time,
        SUM(PHY_READ_CNT) AS phy_read_page_cnt,
        SUM(LOGIC_READ_CNT) AS logic_read_page_cnt,
        ROW_NUMBER() OVER (ORDER BY COUNT(*) DESC) AS row_num
    FROM V$SQL_STAT_HISTORY
    GROUP BY SQL_TXT
) t WHERE executions > 0 AND row_num <= %v ORDER BY %v DESC`
	DmTopSQLMetricExecutions       = "executions"
	DmTopSQLMetricTotalExecTime    = "total_exec_time"
	DmTopSQLMetricAverageExecTime  = "average_exec_time"
	DmTopSQLMetricCPUTime          = "cpu_time"
	DmTopSQLMetricPhyReadPageCnt   = "phy_read_page_cnt"
	DmTopSQLMetricLogicReadPageCnt = "logic_read_page_cnt"
)

// DmTopSQLTask implement the Task interface.
//
// DmTopSQLTask is a loop task which collect Top SQL from DM instance.
type DmTopSQLTask struct {
	*sqlCollector
}

func NewDmTopSQLTask(entry *logrus.Entry, ap *model.AuditPlan) Task {
	task := &DmTopSQLTask{
		sqlCollector: newSQLCollector(entry, ap),
	}
	task.sqlCollector.do = task.collectorDo
	return task
}

func (at *DmTopSQLTask) collectorDo() {
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	inst, _, err := dms.GetInstanceInProjectByName(ctx, string(at.ap.ProjectId), at.ap.InstanceName)
	if err != nil {
		at.logger.Errorf("query instance fail by projectId=%s and instanceName=%s, error: %v",
			string(at.ap.ProjectId), at.ap.InstanceName, err)
		return
	}

	sqls, err := queryTopSQLsForDm(inst, at.ap.InstanceDatabase, at.ap.Params.GetParam("order_by_column").String(),
		at.ap.Params.GetParam("top_n").Int())
	if err != nil {
		at.logger.Errorf("query top sql fail, error: %v", err)
		return
	}

	if len(sqls) == 0 {
		at.logger.Info("sql result count is 0")
		return
	}

	apSQLs := make([]*SQL, 0, len(sqls))
	for _, sql := range sqls {
		apSQLs = append(apSQLs, &SQL{
			SQLContent:  sql.SQLFullText,
			Fingerprint: sql.SQLFullText,
			Info: map[string]interface{}{
				DmTopSQLMetricExecutions:       sql.Executions,
				DmTopSQLMetricTotalExecTime:    sql.TotalExecTime,
				DmTopSQLMetricAverageExecTime:  sql.AverageExecTime,
				DmTopSQLMetricCPUTime:          sql.CPUTime,
				DmTopSQLMetricPhyReadPageCnt:   sql.PhyReadPageCnt,
				DmTopSQLMetricLogicReadPageCnt: sql.LogicReadPageCnt,
			},
		})
	}
	err = at.persist.OverrideAuditPlanSQLs(at.ap.ID, convertSQLsToModelSQLs(apSQLs))
	if err != nil {
		at.logger.Errorf("save top sql to storage fail, error: %v", err)
	}
}

func queryTopSQLsForDm(inst *model.Instance, database string, orderBy string, topN int) ([]*DynPerformanceDmColumns, error) {
	plugin, err := common.NewDriverManagerWithoutAudit(log.NewEntry(), inst, database)
	if err != nil {
		return nil, err
	}
	defer plugin.Close(context.TODO())

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	sql := fmt.Sprintf(DynPerformanceViewDmTpl, topN, orderBy)
	result, err := plugin.Query(ctx, sql, &driverV2.QueryConf{TimeOutSecond: 120})
	if err != nil {
		return nil, err
	}
	var ret []*DynPerformanceDmColumns
	rows := result.Rows
	for _, row := range rows {
		values := row.Values
		if len(values) < 7 {
			continue
		}
		executions, err := strconv.ParseFloat(values[1].Value, 64)
		if err != nil {
			return nil, err
		}
		totalExecTime, err := strconv.ParseFloat(values[2].Value, 64)
		if err != nil {
			return nil, err
		}
		averageExecTime, err := strconv.ParseFloat(values[3].Value, 64)
		if err != nil {
			return nil, err
		}
		cpuTime, err := strconv.ParseFloat(values[4].Value, 64)
		if err != nil {
			return nil, err
		}
		phyReadPageCnt, err := strconv.ParseFloat(values[5].Value, 64)
		if err != nil {
			return nil, err
		}
		logicReadPageCnt, err := strconv.ParseFloat(values[6].Value, 64)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &DynPerformanceDmColumns{
			SQLFullText:      values[0].Value,
			Executions:       executions,
			TotalExecTime:    totalExecTime,
			AverageExecTime:  averageExecTime,
			CPUTime:          cpuTime,
			PhyReadPageCnt:   phyReadPageCnt,
			LogicReadPageCnt: logicReadPageCnt,
		})
	}
	return ret, nil
}

func (at *DmTopSQLTask) Audit() (*AuditResultResp, error) {
	task, err := getTaskWithInstanceByAuditPlan(at.ap, at.persist)
	if err != nil {
		return nil, err
	}
	return at.baseTask.audit(task)
}

func (at *DmTopSQLTask) GetSQLs(args map[string]interface{}) ([]Head, []map[string] /* head name */ string, uint64, error) {
	auditPlanSQLs, count, err := at.persist.GetAuditPlanSQLsByReq(args)
	if err != nil {
		return nil, nil, count, err
	}
	heads := []Head{
		{
			Name: "sql",
			Desc: "SQL语句",
			Type: "sql",
		},
		{
			Name: DmTopSQLMetricExecutions,
			Desc: "总执行次数",
		},
		{
			Name: DmTopSQLMetricTotalExecTime,
			Desc: "总执行时间(s)",
		},
		{
			Name: DmTopSQLMetricAverageExecTime,
			Desc: "平均执行时间(s)",
		},
		{
			Name: DmTopSQLMetricCPUTime,
			Desc: "CPU时间占用(s)",
		},
		{
			Name: DmTopSQLMetricPhyReadPageCnt,
			Desc: "物理读页数",
		},
		{
			Name: DmTopSQLMetricLogicReadPageCnt,
			Desc: "逻辑读页数",
		},
	}
	rows := make([]map[string]string, 0, len(auditPlanSQLs))
	for _, sql := range auditPlanSQLs {
		info := &DynPerformanceDmColumns{}
		if err := json.Unmarshal(sql.Info, info); err != nil {
			return nil, nil, 0, err
		}
		rows = append(rows, map[string]string{
			"sql":                          sql.SQLContent,
			DmTopSQLMetricExecutions:       strconv.Itoa(int(info.Executions)),
			DmTopSQLMetricTotalExecTime:    fmt.Sprintf("%v", utils.Round(float64(info.TotalExecTime)/1000, 3)),   //视图中时间单位是毫秒，所以除以1000得到秒
			DmTopSQLMetricAverageExecTime:  fmt.Sprintf("%v", utils.Round(float64(info.AverageExecTime)/1000, 3)), //视图中时间单位是毫秒，所以除以1000得到秒
			DmTopSQLMetricCPUTime:          fmt.Sprintf("%v", utils.Round(float64(info.CPUTime)/1000, 3)),         //视图中时间单位是毫秒，所以除以1000得到秒
			DmTopSQLMetricPhyReadPageCnt:   strconv.Itoa(int(info.PhyReadPageCnt)),
			DmTopSQLMetricLogicReadPageCnt: strconv.Itoa(int(info.LogicReadPageCnt)),
		})
	}
	return heads, rows, count, nil
}

/*
查询OceanBase版本的SQL语句

 1. 根据文档查询到的支持范围包括：3.2.3-4.3.1(最新)
 2. 测试时使用OceanBase版本3.1.5，3.1.5也支持该用法

参考链接：https://www.oceanbase.com/quicksearch?q=OB_VERSION
*/
const GetOceanbaseVersionSQL string = "SELECT OB_VERSION() FROM DUAL"
const OceanBaseVersion4_0_0 string = "4.0.0"

// ob for oracle 是从4.1.0开始适配的，因此默认版本设定为4.1.0
const DefaultOBForOracleVersion string = "4.1.0"

func getOceanBaseVersion(ctx context.Context, inst *model.Instance, database string) (string, error) {
	plugin, err := common.NewDriverManagerWithoutAudit(log.NewEntry(), inst, database)
	if err != nil {
		return "", err
	}
	defer plugin.Close(context.TODO())
	versionResult, err := plugin.Query(ctx, GetOceanbaseVersionSQL, &driverV2.QueryConf{TimeOutSecond: 20})
	if err != nil {
		return "", err
	}
	if len(versionResult.Column) == 1 && len(versionResult.Rows) == 1 {
		return versionResult.Rows[0].Values[0].Value, nil
	}
	return "", fmt.Errorf("unexpected result of ob version")
}

// ObForOracleTopSQLTask implement the Task interface.
//
// ObForOracleTopSQLTask is a loop task which collect Top SQL from oracle instance.
type ObForOracleTopSQLTask struct {
	*sqlCollector
	obVersion string
}

func NewObForOracleTopSQLTask(entry *logrus.Entry, ap *model.AuditPlan) Task {
	task := &ObForOracleTopSQLTask{
		sqlCollector: newSQLCollector(entry, ap),
	}
	task.sqlCollector.do = task.collectorDo
	return task
}

func (at *ObForOracleTopSQLTask) collectorDo() {
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	inst, _, err := dms.GetInstanceInProjectByName(ctx, string(at.ap.ProjectId), at.ap.InstanceName)
	if err != nil {
		at.logger.Errorf("query instance fail by projectId=%s and instanceName=%s, error: %v",
			string(at.ap.ProjectId), at.ap.InstanceName, err)
		return
	}

	if at.obVersion == "" {
		at.obVersion, err = getOceanBaseVersion(ctx, inst, at.ap.InstanceDatabase)
		if err != nil {
			log.Logger().Errorf("get ocean base version failed, use default version %v, error is %v", DefaultOBForOracleVersion, err)
			at.obVersion = DefaultOBForOracleVersion
		}
	}

	sqls, err := queryTopSQLs(inst, at.ap.InstanceDatabase, at.obVersion,
		at.ap.Params.GetParam("order_by_column").String(),
		at.ap.Params.GetParam("top_n").Int())
	if err != nil {
		at.logger.Errorf("query top sql fail, error: %v", err)
		return
	}

	if len(sqls) > 0 {
		apSQLs := make([]*SQL, 0, len(sqls))
		for _, sql := range sqls {
			apSQLs = append(apSQLs, &SQL{
				SQLContent:  sql.SQLFullText,
				Fingerprint: sql.SQLFullText,
				Info: map[string]interface{}{
					DynPerformanceViewObForOracleColumnExecutions:     sql.Executions,
					DynPerformanceViewObForOracleColumnElapsedTime:    sql.ElapsedTime,
					DynPerformanceViewObForOracleColumnCPUTime:        sql.CPUTime,
					DynPerformanceViewObForOracleColumnDiskReads:      sql.DiskReads,
					DynPerformanceViewObForOracleColumnBufferGets:     sql.BufferGets,
					DynPerformanceViewObForOracleColumnUserIOWaitTime: sql.UserIOWaitTime,
				},
			})
		}
		err = at.persist.OverrideAuditPlanSQLs(at.ap.ID, convertSQLsToModelSQLs(apSQLs))
		if err != nil {
			at.logger.Errorf("save top sql to storage fail, error: %v", err)
		}
	}
}

// 从OceanBase4.0.0版本开始，GV$PLAN_CACHE_PLAN_STAT视图名称调整为GV$OB_PLAN_CACHE_PLAN_STAT
// 参考：https://www.oceanbase.com/docs/common-oceanbase-database-cn-1000000000820360
const PlanCacheViewNameOBV4After string = "GV$OB_PLAN_CACHE_PLAN_STAT"
const PlanCacheViewNameOBV4Before string = "GV$PLAN_CACHE_PLAN_STAT"

func getPlanCacheViewNameByObVersion(obVersion string) string {
	var viewName string = PlanCacheViewNameOBV4After
	isLess, err := utils.IsVersionLessThan(obVersion, OceanBaseVersion4_0_0)
	if err != nil {
		log.Logger().Errorf("compare ocean base version failed, use default view name for top sql %v, error is %v", viewName, err)
		isLess = false
	}
	if isLess {
		viewName = PlanCacheViewNameOBV4Before
	}
	return viewName
}

func queryTopSQLs(inst *model.Instance, database string, obVersion string, orderBy string, topN int) ([]*DynPerformanceObForOracleColumns, error) {
	plugin, err := common.NewDriverManagerWithoutAudit(log.NewEntry(), inst, database)
	if err != nil {
		return nil, err
	}
	defer plugin.Close(context.TODO())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	sql := fmt.Sprintf(
		DynPerformanceViewObForOracleTpl,
		// 通过视图PlanCacheView查询OecanBase for Oracle的Top SQL
		getPlanCacheViewNameByObVersion(obVersion), orderBy, topN,
	)
	result, err := plugin.Query(ctx, sql, &driverV2.QueryConf{TimeOutSecond: 20})
	if err != nil {
		return nil, err
	}
	var ret []*DynPerformanceObForOracleColumns
	rows := result.Rows
	for _, row := range rows {
		values := row.Values
		if len(values) == 0 {
			continue
		}
		executions, err := strconv.ParseFloat(values[1].Value, 64)
		if err != nil {
			return nil, err
		}
		elapsedTime, err := strconv.ParseFloat(values[2].Value, 64)
		if err != nil {
			return nil, err
		}
		cpuTime, err := strconv.ParseFloat(values[3].Value, 64)
		if err != nil {
			return nil, err
		}
		diskReads, err := strconv.ParseFloat(values[4].Value, 64)
		if err != nil {
			return nil, err
		}
		bufferGets, err := strconv.ParseFloat(values[5].Value, 64)
		if err != nil {
			return nil, err
		}
		userIoWaitTime, err := strconv.ParseFloat(values[6].Value, 64)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &DynPerformanceObForOracleColumns{
			SQLFullText:    values[0].Value,
			Executions:     executions,
			ElapsedTime:    elapsedTime,
			CPUTime:        cpuTime,
			DiskReads:      diskReads,
			BufferGets:     bufferGets,
			UserIOWaitTime: userIoWaitTime,
		})
	}
	return ret, nil
}

func (at *ObForOracleTopSQLTask) Audit() (*AuditResultResp, error) {
	task := &model.Task{
		DBType: at.ap.DBType,
	}
	return at.baseTask.audit(task)
}

func (at *ObForOracleTopSQLTask) GetSQLs(args map[string]interface{}) ([]Head, []map[string] /* head name */ string, uint64, error) {
	auditPlanSQLs, count, err := at.persist.GetAuditPlanSQLsByReq(args)
	if err != nil {
		return nil, nil, count, err
	}
	heads := []Head{
		{
			Name: "sql",
			Desc: "SQL语句",
			Type: "sql",
		},
		{
			Name: DynPerformanceViewObForOracleColumnExecutions,
			Desc: "总执行次数",
		},
		{
			Name: DynPerformanceViewObForOracleColumnElapsedTime,
			Desc: "执行时间(s)",
		},
		{
			Name: DynPerformanceViewObForOracleColumnCPUTime,
			Desc: "CPU消耗时间(s)",
		},
		{
			Name: DynPerformanceViewObForOracleColumnDiskReads,
			Desc: "物理读次数",
		},
		{
			Name: DynPerformanceViewObForOracleColumnBufferGets,
			Desc: "逻辑读次数",
		},
		{
			Name: DynPerformanceViewObForOracleColumnUserIOWaitTime,
			Desc: "I/O等待时间(s)",
		},
	}
	rows := make([]map[string]string, 0, len(auditPlanSQLs))
	for _, sql := range auditPlanSQLs {
		info := &DynPerformanceObForOracleColumns{}
		if err := json.Unmarshal(sql.Info, info); err != nil {
			return nil, nil, 0, err
		}
		rows = append(rows, map[string]string{
			"sql": sql.SQLContent,
			DynPerformanceViewObForOracleColumnExecutions:     strconv.Itoa(int(info.Executions)),
			DynPerformanceViewObForOracleColumnElapsedTime:    fmt.Sprintf("%v", utils.Round(float64(info.ElapsedTime)/1000/1000, 3)), //视图中时间单位是微秒，所以除以1000000得到秒
			DynPerformanceViewObForOracleColumnCPUTime:        fmt.Sprintf("%v", utils.Round(float64(info.CPUTime)/1000/1000, 3)),     //视图中时间单位是微秒，所以除以1000000得到秒
			DynPerformanceViewObForOracleColumnDiskReads:      strconv.Itoa(int(info.DiskReads)),
			DynPerformanceViewObForOracleColumnBufferGets:     strconv.Itoa(int(info.BufferGets)),
			DynPerformanceViewObForOracleColumnUserIOWaitTime: fmt.Sprintf("%v", utils.Round(float64(info.UserIOWaitTime)/1000/1000, 3)), //视图中时间单位是微秒，所以除以1000000得到秒
		})
	}
	return heads, rows, count, nil
}

type DynPerformanceObForOracleColumns struct {
	SQLFullText    string  `json:"sql_fulltext"`
	Executions     float64 `json:"executions"`
	ElapsedTime    float64 `json:"elapsed_time"`
	CPUTime        float64 `json:"cpu_time"`
	DiskReads      float64 `json:"disk_reads"`
	BufferGets     float64 `json:"buffer_gets"`
	UserIOWaitTime float64 `json:"user_io_wait_time"`
}

const (
	DynPerformanceViewObForOracleTpl = `
SELECT
    t1.sql_fulltext as sql_fulltext,
    sum(t1.EXECUTIONS) as executions,
    sum(t1.ELAPSED_TIME) as elapsed_time,
    sum(t1.CPU_TIME) as cpu_time,
    sum(t1.DISK_READS) as disk_reads,
    sum(t1.BUFFERS_GETS) as buffer_gets,
    sum(t1.USER_IO_WAIT_TIME) as user_io_wait_time
FROM 
	(
		SELECT
			to_char(QUERY_SQL) as sql_fulltext,
			EXECUTIONS,
			ELAPSED_TIME,
			CPU_TIME,
			DISK_READS,
			BUFFERS_GETS,
			USER_IO_WAIT_TIME
		FROM 
			%v
		WHERE
			to_char(QUERY_SQL) != 'null'
	) t1 
GROUP BY 
    t1.sql_fulltext
ORDER BY 
    %v DESC
FETCH FIRST %v ROWS ONLY
`
	DynPerformanceViewObForOracleColumnExecutions     = "executions"
	DynPerformanceViewObForOracleColumnElapsedTime    = "elapsed_time"
	DynPerformanceViewObForOracleColumnCPUTime        = "cpu_time"
	DynPerformanceViewObForOracleColumnDiskReads      = "disk_reads"
	DynPerformanceViewObForOracleColumnBufferGets     = "buffer_gets"
	DynPerformanceViewObForOracleColumnUserIOWaitTime = "user_io_wait_time"
)

type DynPerformancePgColumns struct {
	SQLFullText    string  `json:"sql_fulltext"`
	Executions     float64 `json:"executions"`
	ElapsedTime    float64 `json:"elapsed_time"`
	DiskReads      float64 `json:"disk_reads"`
	BufferGets     float64 `json:"buffer_gets"`
	UserIOWaitTime float64 `json:"user_io_wait_time"`
}

const (
	DynPerformanceViewPgTpl = `
SELECT query as sql_fulltext,
sum(calls) as executions,
sum(total_exec_time) AS elapsed_time,
sum(shared_blks_read) AS disk_reads, -- 表示从共享缓冲区中读取的块数。这个值表示数据库系统从磁盘或其他存储介质中读取的数据块数量，而不是从内存中读取的数据。
sum(shared_blks_hit) AS buffer_gets, -- 表示从共享缓冲区中命中的块数。这个值表示数据库系统从内存中读取的数据块数量，而不是从磁盘或其他存储介质中读取的数据。
sum(blk_read_time) as user_io_wait_time
FROM pg_stat_statements
WHERE calls > 0
AND query <> '<insufficient privilege>' -- 过滤包含"<insufficient privilege>"的SQL语句 https://github.com/actiontech/sqle-ee/issues/1586
group by query
ORDER BY %v DESC limit %v`
	DynPerformanceViewPgSQLColumnExecutions     = "executions"
	DynPerformanceViewPgSQLColumnElapsedTime    = "elapsed_time"
	DynPerformanceViewPgSQLColumnDiskReads      = "disk_reads"
	DynPerformanceViewPgSQLColumnBufferGets     = "buffer_gets"
	DynPerformanceViewPgSQLColumnUserIOWaitTime = "user_io_wait_time"
)

// PostgreSQLTopSQLTask implement the Task interface.
//
// PostgreSQLTopSQLTask is a loop task which collect Top SQL from oracle instance.
type PostgreSQLTopSQLTask struct {
	*sqlCollector
}

func NewPostgreSQLTopSQLTask(entry *logrus.Entry, ap *model.AuditPlan) Task {
	task := &PostgreSQLTopSQLTask{
		sqlCollector: newSQLCollector(entry, ap),
	}
	task.sqlCollector.do = task.collectorDo
	return task
}

func (at *PostgreSQLTopSQLTask) collectorDo() {
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

	// 超时2分钟
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	inst, _, err := dms.GetInstanceInProjectByName(ctx, string(at.ap.ProjectId), at.ap.InstanceName)
	if err != nil {
		at.logger.Errorf("query instance fail by projectId=%s and instanceName=%s, error: %v",
			string(at.ap.ProjectId), at.ap.InstanceName, err)
		return
	}

	sqls, err := queryTopSQLsForPg(inst, at.ap.InstanceDatabase, at.ap.Params.GetParam("order_by_column").String(),
		at.ap.Params.GetParam("top_n").Int())
	if err != nil {
		at.logger.Errorf("query top sql fail, error: %v", err)
		return
	}

	if len(sqls) == 0 {
		at.logger.Info("sql result count is 0")
		return
	}

	apSQLs := make([]*SQL, 0, len(sqls))
	for _, sql := range sqls {
		apSQLs = append(apSQLs, &SQL{
			SQLContent:  sql.SQLFullText,
			Fingerprint: sql.SQLFullText,
			Info: map[string]interface{}{
				DynPerformanceViewPgSQLColumnExecutions:     sql.Executions,
				DynPerformanceViewPgSQLColumnElapsedTime:    sql.ElapsedTime,
				DynPerformanceViewPgSQLColumnDiskReads:      sql.DiskReads,
				DynPerformanceViewPgSQLColumnBufferGets:     sql.BufferGets,
				DynPerformanceViewPgSQLColumnUserIOWaitTime: sql.UserIOWaitTime,
			},
		})
	}
	err = at.persist.OverrideAuditPlanSQLs(at.ap.ID, convertSQLsToModelSQLs(apSQLs))
	if err != nil {
		at.logger.Errorf("save top sql to storage fail, error: %v", err)
	}
}

func queryTopSQLsForPg(inst *model.Instance, database string, orderBy string, topN int) ([]*DynPerformancePgColumns, error) {
	plugin, err := common.NewDriverManagerWithoutAudit(log.NewEntry(), inst, database)
	if err != nil {
		return nil, err
	}
	defer plugin.Close(context.TODO())

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	sql := fmt.Sprintf(DynPerformanceViewPgTpl, orderBy, topN)
	result, err := plugin.Query(ctx, sql, &driverV2.QueryConf{TimeOutSecond: 120})
	if err != nil {
		return nil, err
	}
	var ret []*DynPerformancePgColumns
	rows := result.Rows
	for _, row := range rows {
		values := row.Values
		if len(values) < 6 {
			continue
		}
		executions, err := strconv.ParseFloat(values[1].Value, 64)
		if err != nil {
			return nil, err
		}
		elapsedTime, err := strconv.ParseFloat(values[2].Value, 64)
		if err != nil {
			return nil, err
		}
		diskReads, err := strconv.ParseFloat(values[3].Value, 64)
		if err != nil {
			return nil, err
		}
		bufferGets, err := strconv.ParseFloat(values[4].Value, 64)
		if err != nil {
			return nil, err
		}
		userIoWaitTime, err := strconv.ParseFloat(values[5].Value, 64)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &DynPerformancePgColumns{
			SQLFullText:    values[0].Value,
			Executions:     executions,
			ElapsedTime:    elapsedTime,
			DiskReads:      diskReads,
			BufferGets:     bufferGets,
			UserIOWaitTime: userIoWaitTime,
		})
	}
	return ret, nil
}

func (at *PostgreSQLTopSQLTask) Audit() (*AuditResultResp, error) {
	task, err := getTaskWithInstanceByAuditPlan(at.ap, at.persist)
	if err != nil {
		return nil, err
	}
	return at.baseTask.audit(task)
}

func (at *PostgreSQLTopSQLTask) GetSQLs(args map[string]interface{}) ([]Head, []map[string] /* head name */ string, uint64, error) {
	auditPlanSQLs, count, err := at.persist.GetAuditPlanSQLsByReq(args)
	if err != nil {
		return nil, nil, count, err
	}
	heads := []Head{
		{
			Name: "sql",
			Desc: "SQL语句",
			Type: "sql",
		},
		{
			Name: DynPerformanceViewPgSQLColumnExecutions,
			Desc: "总执行次数",
		},
		{
			Name: DynPerformanceViewPgSQLColumnElapsedTime,
			Desc: "执行时间(s)",
		},
		{
			Name: DynPerformanceViewPgSQLColumnDiskReads,
			Desc: "物理读块数",
		},
		{
			Name: DynPerformanceViewPgSQLColumnBufferGets,
			Desc: "逻辑读块数",
		},
		{
			Name: DynPerformanceViewPgSQLColumnUserIOWaitTime,
			Desc: "I/O等待时间(s)",
		},
	}
	rows := make([]map[string]string, 0, len(auditPlanSQLs))
	for _, sql := range auditPlanSQLs {
		info := &DynPerformancePgColumns{}
		if err := json.Unmarshal(sql.Info, info); err != nil {
			return nil, nil, 0, err
		}
		rows = append(rows, map[string]string{
			"sql":                                       sql.SQLContent,
			DynPerformanceViewPgSQLColumnExecutions:     strconv.Itoa(int(info.Executions)),
			DynPerformanceViewPgSQLColumnElapsedTime:    fmt.Sprintf("%v", utils.Round(float64(info.ElapsedTime)/1000, 3)), //视图中时间单位是毫秒，所以除以1000得到秒
			DynPerformanceViewPgSQLColumnDiskReads:      strconv.Itoa(int(info.DiskReads)),
			DynPerformanceViewPgSQLColumnBufferGets:     strconv.Itoa(int(info.BufferGets)),
			DynPerformanceViewPgSQLColumnUserIOWaitTime: fmt.Sprintf("%v", utils.Round(float64(info.UserIOWaitTime)/1000, 3)), //视图中时间单位是毫秒，所以除以1000得到秒
		})
	}
	return heads, rows, count, nil
}
