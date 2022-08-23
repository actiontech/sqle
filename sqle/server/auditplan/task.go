package auditplan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/oracle"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/format"
	"github.com/sirupsen/logrus"
)

var errNoSQLInAuditPlan = errors.New(errors.DataConflict, fmt.Errorf("there is no SQLs in audit plan"))
var errNoSQLNeedToBeAudited = errors.New(errors.DataConflict, fmt.Errorf("there is no SQLs need to be audited in audit plan"))

type Task interface {
	Start() error
	Stop() error
	Audit() (*model.AuditPlanReportV2, error)
	FullSyncSQLs([]*SQL) error
	PartialSyncSQLs([]*SQL) error
	GetSQLs(map[string]interface{}) ([]Head, []map[string] /* head name */ string, uint64, error)
}

type Head struct {
	Name string
	Desc string
	Type string
}

type SQL struct {
	SQLContent  string
	Fingerprint string
	Info        map[string]interface{}
}

func NewTask(entry *logrus.Entry, ap *model.AuditPlan) Task {
	entry = entry.WithField("name", ap.Name)
	switch ap.Type {
	case TypeMySQLSchemaMeta:
		return NewSchemaMetaTask(entry, ap)
	case TypeOracleTopSQL:
		return NewOracleTopSQLTask(entry, ap)
	case TypeTiDBAuditLog:
		return NewTiDBAuditLogTask(entry, ap)
	default:
		return NewDefaultTask(entry, ap)
	}
}

type baseTask struct {
	ap *model.AuditPlan
	// persist is a database handle which store AuditPlan.
	persist *model.Storage
	logger  *logrus.Entry
}

func newBaseTask(entry *logrus.Entry, ap *model.AuditPlan) *baseTask {
	log.NewEntry()
	return &baseTask{
		ap:      ap,
		persist: model.GetStorage(),
		logger:  entry,
	}
}

func (at *baseTask) Start() error {
	return nil
}

func (at *baseTask) Stop() error {
	return nil
}

func (at *baseTask) audit(task *model.Task) (*model.AuditPlanReportV2, error) {
	auditPlanSQLs, err := at.persist.GetAuditPlanSQLs(at.ap.Name)
	if err != nil {
		return nil, err
	}

	if len(auditPlanSQLs) == 0 {
		return nil, errNoSQLInAuditPlan
	}

	filteredSqls, err := filterSQLsByPeriod(at.ap.Params, auditPlanSQLs)
	if err != nil {
		return nil, err
	}

	if len(filteredSqls) == 0 {
		return nil, errNoSQLNeedToBeAudited
	}

	for i, sql := range filteredSqls {
		task.ExecuteSQLs = append(task.ExecuteSQLs, &model.ExecuteSQL{
			BaseSQL: model.BaseSQL{
				Number:  uint(i),
				Content: sql.SQLContent,
			},
		})
	}

	err = server.Audit(at.logger, task, at.ap.RuleTemplateName)
	if err != nil {
		return nil, err
	}

	auditPlanReport := &model.AuditPlanReportV2{
		AuditPlanID: at.ap.ID,
		PassRate:    task.PassRate,
		Score:       task.Score,
		AuditLevel:  task.AuditLevel,
	}
	for i, executeSQL := range task.ExecuteSQLs {
		auditPlanReport.AuditPlanReportSQLs = append(auditPlanReport.AuditPlanReportSQLs, &model.AuditPlanReportSQLV2{
			SQL:         executeSQL.Content,
			Number:      uint(i + 1),
			AuditResult: executeSQL.AuditResult,
		})
	}
	err = at.persist.Save(auditPlanReport)
	if err != nil {
		return nil, err
	}
	return auditPlanReport, nil
}

func filterSQLsByPeriod(params params.Params, sqls []*model.AuditPlanSQLV2) (filteredSqls []*model.AuditPlanSQLV2, err error) {
	period := params.GetParam(paramKeyAuditSQLsScrappedInLastPeriodMinute).Int()
	if period <= 0 {
		return sqls, nil
	}

	t := time.Now()
	minus := -1
	startTime := t.Add(time.Minute * time.Duration(minus*period))
	for _, sql := range sqls {
		var info = struct {
			LastReceiveTimestamp time.Time `json:"last_receive_timestamp"`
		}{}
		err := json.Unmarshal(sql.Info, &info)
		if err != nil {
			return nil, fmt.Errorf("parse last_receive_timestamp failed: %v", err)
		}

		if info.LastReceiveTimestamp.Before(startTime) {
			continue
		}
		newSql := *sql
		filteredSqls = append(filteredSqls, &newSql)
	}
	return filteredSqls, nil
}

type sqlCollector struct {
	*baseTask
	sync.WaitGroup
	isStarted bool
	cancel    chan struct{}
	do        func()
}

func newSQLCollector(entry *logrus.Entry, ap *model.AuditPlan) *sqlCollector {
	return &sqlCollector{
		newBaseTask(entry, ap),
		sync.WaitGroup{},
		false,
		make(chan struct{}),
		func() { // default
			entry.Warn("sql collector do nothing")
		},
	}
}

func (at *sqlCollector) Start() error {
	if at.isStarted {
		return nil
	}
	at.WaitGroup.Add(1)
	go func() {
		at.isStarted = true
		at.logger.Infof("start task")
		at.loop(at.cancel)
		at.WaitGroup.Done()
	}()
	return nil
}

func (at *sqlCollector) Stop() error {
	if !at.isStarted {
		return nil
	}
	at.cancel <- struct{}{}
	at.isStarted = false
	at.WaitGroup.Wait()
	at.logger.Infof("stop task")
	return nil
}

func (at *sqlCollector) FullSyncSQLs(sqls []*SQL) error {
	at.logger.Warnf("someone try to sync sql to audit plan(%v), but sql should collected by task itself", at.ap.Name)
	return nil
}

func (at *sqlCollector) PartialSyncSQLs(sqls []*SQL) error {
	at.logger.Warnf("someone try to sync sql to audit plan(%v), but sql should collected by task itself", at.ap.Name)
	return nil
}

func (at *sqlCollector) loop(cancel chan struct{}) {
	interval := at.ap.Params.GetParam(paramKeyCollectIntervalMinute).Int()
	if interval == 0 {
		interval = 60
	}
	at.do()

	tk := time.NewTicker(time.Duration(interval) * time.Minute)
	for {
		select {
		case <-cancel:
			tk.Stop()
			return
		case <-tk.C:
			at.logger.Infof("tick %s", at.ap.Name)
			at.do()
		}
	}
}

type DefaultTask struct {
	*baseTask
}

func NewDefaultTask(entry *logrus.Entry, ap *model.AuditPlan) *DefaultTask {
	return &DefaultTask{newBaseTask(entry, ap)}
}

func (at *DefaultTask) Audit() (*model.AuditPlanReportV2, error) {
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

func convertSQLsToModelSQLs(sqls []*SQL) []*model.AuditPlanSQLV2 {
	as := make([]*model.AuditPlanSQLV2, len(sqls))
	for i, sql := range sqls {
		data, _ := json.Marshal(sql.Info)
		as[i] = &model.AuditPlanSQLV2{
			Fingerprint: sql.Fingerprint,
			SQLContent:  sql.SQLContent,
			Info:        data,
		}
	}
	return as
}

func convertRawSQLToModelSQLs(sqls []string) []*model.AuditPlanSQLV2 {
	as := make([]*model.AuditPlanSQLV2, len(sqls))
	for i, sql := range sqls {
		as[i] = &model.AuditPlanSQLV2{
			Fingerprint: sql,
			SQLContent:  sql,
		}
	}
	return as
}

func (at *baseTask) FullSyncSQLs(sqls []*SQL) error {
	return at.persist.OverrideAuditPlanSQLs(at.ap.Name, convertSQLsToModelSQLs(sqls))
}

func (at *baseTask) PartialSyncSQLs(sqls []*SQL) error {
	return at.persist.UpdateDefaultAuditPlanSQLs(at.ap.Name, convertSQLsToModelSQLs(sqls))
}

func (at *baseTask) GetSQLs(args map[string]interface{}) ([]Head, []map[string] /* head name */ string, uint64, error) {
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
	}
	rows := make([]map[string]string, 0, len(auditPlanSQLs))
	for _, sql := range auditPlanSQLs {
		var info = struct {
			Counter              uint64 `json:"counter"`
			LastReceiveTimestamp string `json:"last_receive_timestamp"`
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
		})
	}
	return head, rows, count, nil
}

type SchemaMetaTask struct {
	*sqlCollector
}

func NewSchemaMetaTask(entry *logrus.Entry, ap *model.AuditPlan) *SchemaMetaTask {
	sqlCollector := newSQLCollector(entry, ap)
	task := &SchemaMetaTask{
		sqlCollector,
	}
	sqlCollector.do = task.collectorDo
	return task
}

func (at *SchemaMetaTask) collectorDo() {
	if at.ap.InstanceName == "" {
		at.logger.Warnf("instance is not configured")
		return
	}
	if at.ap.InstanceDatabase == "" {
		at.logger.Warnf("instance schema is not configured")
		return
	}
	instance, _, err := at.persist.GetInstanceByName(at.ap.InstanceName)
	if err != nil {
		return
	}
	db, err := executor.NewExecutor(at.logger, &driver.DSN{
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

	tables, err := db.ShowSchemaTables(at.ap.InstanceDatabase)
	if err != nil {
		at.logger.Errorf("get schema table fail, error: %v", err)
		return
	}
	var views []string
	if at.ap.Params.GetParam("collect_view").Bool() {
		views, err = db.ShowSchemaViews(at.ap.InstanceDatabase)
		if err != nil {
			at.logger.Errorf("get schema view fail, error: %v", err)
			return
		}
	}
	sqls := make([]string, 0, len(tables)+len(views))
	for _, table := range tables {
		sql, err := db.ShowCreateTable(utils.SupplementalQuotationMarks(at.ap.InstanceDatabase), utils.SupplementalQuotationMarks(table))
		if err != nil {
			at.logger.Errorf("show create table fail, error: %v", err)
			return
		}
		sqls = append(sqls, sql)
	}
	for _, view := range views {
		sql, err := db.ShowCreateView(utils.SupplementalQuotationMarks(view))
		if err != nil {
			at.logger.Errorf("show create view fail, error: %v", err)
			return
		}
		sqls = append(sqls, sql)
	}
	if len(sqls) > 0 {
		err = at.persist.OverrideAuditPlanSQLs(at.ap.Name, convertRawSQLToModelSQLs(sqls))
		if err != nil {
			at.logger.Errorf("save schema meta to storage fail, error: %v", err)
		}
	}
}

func (at *SchemaMetaTask) Audit() (*model.AuditPlanReportV2, error) {
	task := &model.Task{
		DBType: at.ap.DBType,
	}
	return at.baseTask.audit(task)
}

func (at *SchemaMetaTask) GetSQLs(args map[string]interface{}) ([]Head, []map[string] /* head name */ string, uint64, error) {
	auditPlanSQLs, count, err := at.persist.GetAuditPlanSQLsByReq(args)
	if err != nil {
		return nil, nil, count, err
	}
	head := []Head{
		{
			Name: "sql",
			Desc: "SQL语句",
			Type: "sql",
		},
	}
	rows := make([]map[string]string, 0, len(auditPlanSQLs))
	for _, sql := range auditPlanSQLs {
		rows = append(rows, map[string]string{
			"sql": sql.SQLContent,
		})
	}
	return head, rows, count, nil
}

// OracleTopSQLTask implement the Task interface.
//
// OracleTopSQLTask is a loop task which collect Top SQL from oracle instance.
type OracleTopSQLTask struct {
	*sqlCollector
}

func NewOracleTopSQLTask(entry *logrus.Entry, ap *model.AuditPlan) *OracleTopSQLTask {
	task := &OracleTopSQLTask{
		sqlCollector: newSQLCollector(entry, ap),
	}
	task.sqlCollector.do = task.collectorDo
	return task
}

func (at *OracleTopSQLTask) collectorDo() {
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
	// This depends on: https://github.com/actiontech/sqle-oracle-plugin.
	// If your Oracle db plugin does not implement the parameter `service_name`,
	// you can only use the default service name `XE`.
	// TODO: using DB plugin to query SQL.
	serviceName := inst.AdditionalParams.GetParam("service_name").String()
	dsn := &oracle.DSN{
		Host:        inst.Host,
		Port:        inst.Port,
		User:        inst.User,
		Password:    inst.Password,
		ServiceName: serviceName,
	}
	db, err := oracle.NewDB(dsn)
	if err != nil {
		at.logger.Errorf("connect to instance fail, error: %v", err)
		return
	}
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sqls, err := db.QueryTopSQLs(ctx, at.ap.Params.GetParam("top_n").Int(), at.ap.Params.GetParam("order_by_column").String())
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
					oracle.DynPerformanceViewSQLAreaColumnExecutions:     sql.Executions,
					oracle.DynPerformanceViewSQLAreaColumnElapsedTime:    sql.ElapsedTime,
					oracle.DynPerformanceViewSQLAreaColumnCPUTime:        sql.CPUTime,
					oracle.DynPerformanceViewSQLAreaColumnDiskReads:      sql.DiskReads,
					oracle.DynPerformanceViewSQLAreaColumnBufferGets:     sql.BufferGets,
					oracle.DynPerformanceViewSQLAreaColumnUserIOWaitTime: sql.UserIOWaitTime,
				},
			})
		}

		err = at.persist.OverrideAuditPlanSQLs(at.ap.Name, convertSQLsToModelSQLs(apSQLs))
		if err != nil {
			at.logger.Errorf("save top sql to storage fail, error: %v", err)
		}
	}
}

func (at *OracleTopSQLTask) Audit() (*model.AuditPlanReportV2, error) {
	task := &model.Task{
		DBType: at.ap.DBType,
	}
	return at.baseTask.audit(task)
}

func (at *OracleTopSQLTask) GetSQLs(args map[string]interface{}) ([]Head, []map[string] /* head name */ string, uint64, error) {
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
			Name: oracle.DynPerformanceViewSQLAreaColumnExecutions,
			Desc: "总执行次数",
		},
		{
			Name: oracle.DynPerformanceViewSQLAreaColumnElapsedTime,
			Desc: "执行时间(s)",
		},
		{
			Name: oracle.DynPerformanceViewSQLAreaColumnCPUTime,
			Desc: "CPU消耗时间(s)",
		},
		{
			Name: oracle.DynPerformanceViewSQLAreaColumnDiskReads,
			Desc: "物理读",
		},
		{
			Name: oracle.DynPerformanceViewSQLAreaColumnBufferGets,
			Desc: "逻辑读",
		},
		{
			Name: oracle.DynPerformanceViewSQLAreaColumnUserIOWaitTime,
			Desc: "I/O等待时间(s)",
		},
	}
	rows := make([]map[string]string, 0, len(auditPlanSQLs))
	for _, sql := range auditPlanSQLs {
		info := &oracle.DynPerformanceSQLArea{}
		if err := json.Unmarshal(sql.Info, info); err != nil {
			return nil, nil, 0, err
		}
		rows = append(rows, map[string]string{
			"sql": sql.SQLContent,
			oracle.DynPerformanceViewSQLAreaColumnExecutions:     strconv.FormatInt(info.Executions, 10),
			oracle.DynPerformanceViewSQLAreaColumnElapsedTime:    fmt.Sprintf("%v", utils.Round(float64(info.ElapsedTime)/1000/1000, 3)),
			oracle.DynPerformanceViewSQLAreaColumnCPUTime:        fmt.Sprintf("%v", utils.Round(float64(info.CPUTime)/1000/1000, 3)),
			oracle.DynPerformanceViewSQLAreaColumnDiskReads:      strconv.FormatInt(info.DiskReads, 10),
			oracle.DynPerformanceViewSQLAreaColumnBufferGets:     strconv.FormatInt(info.BufferGets, 10),
			oracle.DynPerformanceViewSQLAreaColumnUserIOWaitTime: fmt.Sprintf("%v", utils.Round(float64(info.UserIOWaitTime)/1000/1000, 3)),
		})
	}
	return heads, rows, count, nil
}

type TiDBAuditLogTask struct {
	*DefaultTask
}

func NewTiDBAuditLogTask(entry *logrus.Entry, ap *model.AuditPlan) *TiDBAuditLogTask {
	return &TiDBAuditLogTask{NewDefaultTask(entry, ap)}
}

func (at *TiDBAuditLogTask) Audit() (*model.AuditPlanReportV2, error) {
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

	auditPlanSQLs, err := at.persist.GetAuditPlanSQLs(at.ap.Name)
	if err != nil {
		return nil, err
	}

	if len(auditPlanSQLs) == 0 {
		return nil, errNoSQLInAuditPlan
	}

	filteredSqls, err := filterSQLsByPeriod(at.ap.Params, auditPlanSQLs)
	if err != nil {
		return nil, err
	}

	if len(filteredSqls) == 0 {
		return nil, errNoSQLNeedToBeAudited
	}

	for i, sql := range filteredSqls {
		schema := ""
		info, _ := sql.Info.OriginValue()
		if schemaStr, ok := info[server.AuditSchema].(string); ok {
			schema = schemaStr
		}

		task.ExecuteSQLs = append(task.ExecuteSQLs, &model.ExecuteSQL{
			BaseSQL: model.BaseSQL{
				Number:  uint(i),
				Content: sql.SQLContent,
				Schema:  schema,
			},
		})
	}

	err = server.HookAudit(at.logger, task, &TiDBAuditHook{}, at.ap.RuleTemplateName)
	if err != nil {
		return nil, err
	}

	auditPlanReport := &model.AuditPlanReportV2{
		AuditPlanID: at.ap.ID,
		PassRate:    task.PassRate,
		Score:       task.Score,
		AuditLevel:  task.AuditLevel,
	}
	for i, executeSQL := range task.ExecuteSQLs {
		auditPlanReport.AuditPlanReportSQLs = append(auditPlanReport.AuditPlanReportSQLs, &model.AuditPlanReportSQLV2{
			SQL:         executeSQL.Content,
			Number:      uint(i + 1),
			AuditResult: executeSQL.AuditResult,
		})
	}
	err = at.persist.Save(auditPlanReport)
	if err != nil {
		return nil, err
	}
	return auditPlanReport, nil
}

// 审核前填充上缺失的schema, 审核后还原被审核SQL, 并添加注释说明sql在哪个库执行的
type TiDBAuditHook struct {
	originalSQL string
}

func (t *TiDBAuditHook) BeforeAudit(sql *model.ExecuteSQL) {
	if sql.Schema == "" {
		return
	}
	t.originalSQL = sql.Content
	newSQL, err := tidbCompletionSchema(sql.Content, sql.Schema)
	if err != nil {
		return
	}
	sql.Content = newSQL
}

func (t *TiDBAuditHook) AfterAudit(sql *model.ExecuteSQL) {
	if sql.Schema == "" {
		return
	}
	sql.Content = fmt.Sprintf("%v -- current schema: %v", t.originalSQL, sql.Schema)
}

// 填充sql缺失的schema
func tidbCompletionSchema(sql, schema string) (string, error) {
	stmts, _, err := parser.New().PerfectParse(sql, "", "")
	if err != nil {
		return "", err
	}
	if len(stmts) != 1 {
		return "", parser.ErrSyntax
	}

	stmts[0].Accept(&completionSchemaVisitor{schema: schema})
	buf := new(bytes.Buffer)
	restoreCtx := format.NewRestoreCtx(format.DefaultRestoreFlags, buf)
	err = stmts[0].Restore(restoreCtx)
	return buf.String(), err
}

// completionSchemaVisitor implements ast.Visitor interface.
type completionSchemaVisitor struct {
	schema string
}

func (g *completionSchemaVisitor) Enter(n ast.Node) (node ast.Node, skipChildren bool) {
	if stmt, ok := n.(*ast.TableName); ok {
		if stmt.Schema.L == "" {
			stmt.Schema.L = strings.ToLower(g.schema)
			stmt.Schema.O = g.schema
		}
	}
	return n, false
}

func (g *completionSchemaVisitor) Leave(n ast.Node) (node ast.Node, ok bool) {
	return n, true
}
