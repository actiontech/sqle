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
	var task *model.Task
	if at.ap.InstanceName == "" {
		task = &model.Task{
			DBType: at.ap.DBType,
		}
	} else {
		instance, _, err := dms.GetInstanceInProjectByName(context.Background(), string(at.ap.ProjectId), at.ap.InstanceName)
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
	SELECT sql_text,db,TIME_TO_SEC(query_time) AS query_time, start_time
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
			fpInfo := fmt.Sprintf(`{"counter":%v,"last_receive_timestamp":"%v","schema":"%v","average_query_time":%d,"start_time":"%v"}`,
				fp.counter, now.Format(time.RFC3339), fp.schema, fp.queryTimeSeconds, fp.startTime)
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
}

type sqlFromSlowLogs []*sqlFromSlowLog

type sqlFingerprintInfo struct {
	lastSql               string
	lastSqlSchema         string
	sqlCount              int
	totalQueryTimeSeconds int
	startTime             string
}

func (s *sqlFingerprintInfo) queryTime() int {
	return s.totalQueryTimeSeconds / s.sqlCount
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
		} else {
			sqlInfosMap[fp] = &sqlFingerprintInfo{
				sqlCount:              1,
				lastSql:               sqlItem.sql,
				lastSqlSchema:         sqlItem.schema,
				totalQueryTimeSeconds: sqlItem.queryTimeSeconds,
				startTime:             sqlItem.startTime,
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

	task := &model.Task{DBType: at.ap.DBType}

	if at.ap.InstanceName != "" {
		instance, _, err := dms.GetInstanceInProjectByName(context.Background(), string(at.ap.ProjectId), at.ap.InstanceName)
		if err != nil {
			return nil, err
		}
		task.Instance = instance
		task.Schema = at.ap.InstanceDatabase
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
	}) {
		return "", fmt.Errorf("invalid indicator: %v", indicator)
	}
	return indicator, nil
}

// ref: https://www.ibm.com/docs/zh/db2/11.1?topic=views-snap-get-dyn-sql-dynsql-snapshot
func (at *DB2TopSQLTask) collectSQL() (string, error) {
	sql := `
SELECT 
    stmt_text,   
	num_executions,   
    real(total_exec_time)*1000+DECIMAL(real(total_exec_time_ms)/1000,18,3) AS total_elapsed_time,   
	DECIMAL((real(total_exec_time)*1000+real(total_exec_time_ms)/1000)/real(num_executions),18,3) AS avg_elapsed_time_ms,   
    DECIMAL((real(total_sys_cpu_time)*1000+real(total_sys_cpu_time_ms)/1000+real(total_usr_cpu_time)*1000+real(total_usr_cpu_time_ms)/1000)/real(num_executions),18,3) as avg_cpu_time_ms    
FROM sysibmadm.snapdyn_sql     
WHERE stmt_text !='' AND num_executions > 0    
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
	DB2IndicatorTotalElapsedTime   = "total_elapsed_time"  // 总执行时间
	DB2IndicatorAverageElapsedTime = "avg_elapsed_time_ms" // 平均执行时间
	DB2IndicatorAverageCPUTime     = "avg_cpu_time_ms"     // 平均 CPU 时间
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
				DatabaseName:     at.ap.InstanceName,
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
			Desc: "平均 CPU 时间(ms)",
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
	task := &model.Task{
		DBType: at.ap.DBType,
	}
	if at.ap.InstanceName != "" {
		instance, _, err := dms.GetInstanceInProjectByName(context.Background(), string(at.ap.ProjectId), at.ap.InstanceName)
		if err != nil {
			return nil, err
		}
		task.Instance = instance
		task.Schema = at.ap.InstanceDatabase
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
