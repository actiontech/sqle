package auditplan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/oracle"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/actiontech/sqle/sqle/utils"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	rds20140815 "github.com/alibabacloud-go/rds-20140815/v2/client"
	_util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/rds"
	"github.com/percona/go-mysql/query"
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
	Audit() (*AuditResultResp, error)
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
	Schema      string
	Info        map[string]interface{}
}

func NewTask(entry *logrus.Entry, ap *model.AuditPlan) Task {
	entry = entry.WithField("name", ap.Name)

	meta, err := GetMeta(ap.Type)
	if err != nil || meta.CreateTask == nil {
		return NewDefaultTask(entry, ap)
	}

	return meta.CreateTask(entry, ap)
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

type AuditResultResp struct {
	AuditPlanID  uint64
	Task         *model.Task
	FilteredSqls []*model.AuditPlanSQLV2
}

func (at *baseTask) audit(task *model.Task) (*AuditResultResp, error) {
	auditPlanSQLs, err := at.persist.GetAuditPlanSQLs(at.ap.ID)
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

	err = server.Audit(at.logger, task, &at.ap.ProjectId, at.ap.RuleTemplateName)
	if err != nil {
		return nil, err
	}

	return &AuditResultResp{
		AuditPlanID:  uint64(at.ap.ID),
		Task:         task,
		FilteredSqls: filteredSqls,
	}, nil
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

	loopInterval func() time.Duration
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
		func() time.Duration {
			interval := ap.Params.GetParam(paramKeyCollectIntervalMinute).Int()
			if interval == 0 {
				interval = 60
			}
			return time.Minute * time.Duration(interval)
		},
	}
}

func (at *sqlCollector) Start() error {
	if at.isStarted {
		return nil
	}
	interval := at.loopInterval()

	at.WaitGroup.Add(1)
	go func() {
		at.isStarted = true
		at.logger.Infof("start task")
		at.loop(at.cancel, interval)
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

func (at *sqlCollector) loop(cancel chan struct{}, interval time.Duration) {
	at.do()
	if interval == 0 {
		at.logger.Warnf("task(%v) loop interval can not be zero", at.ap.Name)
		return
	}

	tk := time.NewTicker(interval)
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

func NewDefaultTask(entry *logrus.Entry, ap *model.AuditPlan) Task {
	return &DefaultTask{newBaseTask(entry, ap)}
}

func (at *DefaultTask) Audit() (*AuditResultResp, error) {
	var task *model.Task
	if at.ap.InstanceName == "" {
		task = &model.Task{
			DBType: at.ap.DBType,
		}
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		defer cancel()
		instance, _, err := dms.GetInstanceInProjectByName(ctx, string(at.ap.ProjectId), at.ap.InstanceName)
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
			Schema:      sql.Schema,
			Info:        data,
		}
	}
	return as
}

func convertRawSQLToModelSQLs(sqls []string, schema string) []*model.AuditPlanSQLV2 {
	as := make([]*model.AuditPlanSQLV2, len(sqls))
	for i, sql := range sqls {
		as[i] = &model.AuditPlanSQLV2{
			Fingerprint: sql,
			SQLContent:  sql,
			Schema:      schema,
		}
	}
	return as
}

func (at *baseTask) FullSyncSQLs(sqls []*SQL) error {
	return at.persist.OverrideAuditPlanSQLs(at.ap.ID, convertSQLsToModelSQLs(sqls))
}

func (at *baseTask) PartialSyncSQLs(sqls []*SQL) error {
	return at.persist.UpdateDefaultAuditPlanSQLs(at.ap.ID, convertSQLsToModelSQLs(sqls))
}

func (at *baseTask) GetSQLs(args map[string]interface{}) ([]Head, []map[string] /* head name */ string, uint64, error) {
	return baseTaskGetSQLs(args, at.persist)
}

func baseTaskGetSQLs(args map[string]interface{}, persist *model.Storage) ([]Head, []map[string] /* head name */ string, uint64, error) {
	auditPlanSQLs, count, err := persist.GetAuditPlanSQLsByReq(args)
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

func getTaskWithInstanceByAuditPlan(ap *model.AuditPlan, persist *model.Storage) (*model.Task, error) {
	var task *model.Task
	if ap.InstanceName == "" {
		task = &model.Task{
			DBType: ap.DBType,
		}
	} else {
		instance, _, err := dms.GetInstanceInProjectByName(context.TODO(), string(ap.ProjectId), ap.InstanceName)
		if err != nil {
			return nil, err
		}
		task = &model.Task{
			Instance: instance,
			Schema:   ap.InstanceDatabase,
			DBType:   ap.DBType,
		}
	}
	return task, nil
}

type SchemaMetaTask struct {
	*sqlCollector
}

func NewSchemaMetaTask(entry *logrus.Entry, ap *model.AuditPlan) Task {
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	instance, _, err := dms.GetInstanceInProjectByName(ctx, string(at.ap.ProjectId), at.ap.InstanceName)
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
		err = at.persist.OverrideAuditPlanSQLs(at.ap.ID, convertRawSQLToModelSQLs(sqls, at.ap.InstanceDatabase))
		if err != nil {
			at.logger.Errorf("save schema meta to storage fail, error: %v", err)
		}
	}
}

func (at *SchemaMetaTask) Audit() (*AuditResultResp, error) {
	task, err := getTaskWithInstanceByAuditPlan(at.ap, at.persist)
	if err != nil {
		return nil, err
	}
	return at.baseTask.audit(task)
}

func (at *SchemaMetaTask) GetSQLs(args map[string]interface{}) ([]Head, []map[string] /* head name */ string, uint64, error) {
	auditPlanSQLs, count, err := at.persist.GetAuditPlanSQLsByReq(args)
	if err != nil {
		return nil, nil, count, err
	}
	head, rows := buildSchemaMetaSQLsResult(auditPlanSQLs)
	return head, rows, count, nil
}

func buildSchemaMetaSQLsResult(auditPlanSQLs []*model.AuditPlanSQLListDetail) ([]Head, []map[string] /* head name */ string) {
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
	return head, rows
}

// OracleTopSQLTask implement the Task interface.
//
// OracleTopSQLTask is a loop task which collect Top SQL from oracle instance.
type OracleTopSQLTask struct {
	*sqlCollector
}

func NewOracleTopSQLTask(entry *logrus.Entry, ap *model.AuditPlan) Task {
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	inst, _, err := dms.GetInstanceInProjectByName(ctx, string(at.ap.ProjectId), at.ap.InstanceName)

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

	ctx, cancel = context.WithCancel(context.Background())
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

		err = at.persist.OverrideAuditPlanSQLs(at.ap.ID, convertSQLsToModelSQLs(apSQLs))
		if err != nil {
			at.logger.Errorf("save top sql to storage fail, error: %v", err)
		}
	}
}

func (at *OracleTopSQLTask) Audit() (*AuditResultResp, error) {
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

func NewTiDBAuditLogTask(entry *logrus.Entry, ap *model.AuditPlan) Task {
	return &TiDBAuditLogTask{NewDefaultTask(entry, ap).(*DefaultTask)}
}

func (at *TiDBAuditLogTask) Audit() (*AuditResultResp, error) {
	var task *model.Task
	if at.ap.InstanceName == "" {
		task = &model.Task{
			DBType: at.ap.DBType,
		}
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		defer cancel()

		instance, _, err := dms.GetInstanceInProjectByName(ctx, string(at.ap.ProjectId), at.ap.InstanceName)
		if err != nil {
			return nil, err
		}
		task = &model.Task{
			Instance: instance,
			Schema:   at.ap.InstanceDatabase,
			DBType:   at.ap.DBType,
		}
	}

	auditPlanSQLs, err := at.persist.GetAuditPlanSQLs(at.ap.ID)
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

		executeSQL := &model.ExecuteSQL{
			BaseSQL: model.BaseSQL{
				Number:  uint(i),
				Content: sql.SQLContent,
				Schema:  schema,
			},
		}

		task.ExecuteSQLs = append(task.ExecuteSQLs, executeSQL)
	}

	hook := &TiDBAuditHook{
		originalSQLs: map[*model.ExecuteSQL]string{},
	}
	err = server.HookAudit(at.logger, task, hook, &at.ap.ProjectId, at.ap.RuleTemplateName)
	if err != nil {
		return nil, err
	}

	return &AuditResultResp{
		AuditPlanID:  uint64(at.ap.ID),
		Task:         task,
		FilteredSqls: filteredSqls,
	}, nil
}

// 审核前填充上缺失的schema, 审核后还原被审核SQL, 并添加注释说明sql在哪个库执行的
type TiDBAuditHook struct {
	originalSQLs map[*model.ExecuteSQL]string
}

func (t *TiDBAuditHook) BeforeAudit(sql *model.ExecuteSQL) {
	if sql.Schema == "" {
		return
	}
	t.originalSQLs[sql] = sql.Content
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
	if o, ok := t.originalSQLs[sql]; ok {
		sql.Content = fmt.Sprintf("%v -- current schema: %v", o, sql.Schema)
	}
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

// aliRdsMySQLTask implement the Task interface.
//
// aliRdsMySQLTask is a loop task which collect slow log from ali rds MySQL instance.
type aliRdsMySQLTask struct {
	*sqlCollector
	lastEndTime *time.Time
	pullLogs    func(client *rds20140815.Client, DBInstanId string, startTime, endTime time.Time, pageSize, pageNum int32) (sqls []SqlFromAliCloud, err error)
}

func (at *aliRdsMySQLTask) collectorDo() {
	if at.ap.InstanceName == "" {
		at.logger.Warnf("instance is not configured")
		return
	}

	rdsDBInstanceId := at.ap.Params.GetParam(paramKeyDBInstanceId).String()
	if rdsDBInstanceId == "" {
		at.logger.Warnf("rds DB instance ID is not configured")
		return
	}

	accessKeyId := at.ap.Params.GetParam(paramKeyAccessKeyId).String()
	if accessKeyId == "" {
		at.logger.Warnf("Access Key ID is not configured")
		return
	}

	accessKeySecret := at.ap.Params.GetParam(paramKeyAccessKeySecret).String()
	if accessKeySecret == "" {
		at.logger.Warnf("Access Key Secret is not configured")
		return
	}

	rdsPath := at.ap.Params.GetParam(paramKeyRdsPath).String()
	if rdsPath == "" {
		at.logger.Warnf("RDS Path is not configured")
		return
	}

	firstScrapInLastHours := at.ap.Params.GetParam(paramKeyFirstSqlsScrappedInLastPeriodHours).Int()
	if firstScrapInLastHours == 0 {
		firstScrapInLastHours = 24
	}
	theMaxSupportedDays := 31 // 支持往前查看慢日志的最大天数
	hoursDuringADay := 24
	if firstScrapInLastHours > theMaxSupportedDays*hoursDuringADay {
		at.logger.Warnf("Can not get slow logs from so early time. firstScrapInLastHours=%v", firstScrapInLastHours)
		return
	}

	client, err := at.CreateClient(rdsPath, tea.String(accessKeyId), tea.String(accessKeySecret))
	if err != nil {
		at.logger.Warnf("create client for polardb mysql failed: %v", err)
		return
	}

	pageSize := 100
	now := time.Now().UTC()
	var startTime time.Time
	if at.isFirstScrap() {
		startTime = now.Add(time.Duration(-1*firstScrapInLastHours) * time.Hour)
	} else {
		startTime = *at.lastEndTime
	}
	var pageNum int32 = 1
	slowSqls := []SqlFromAliCloud{}
	for {
		newSlowSqls, err := at.pullLogs(client, rdsDBInstanceId, startTime, now, int32(pageSize), pageNum)
		if err != nil {
			at.logger.Warnf("pull rds logs failed: %v", err)
			return
		}
		filteredNewSlowSqls := at.filterSlowSqlsByExecutionTime(newSlowSqls, startTime)
		slowSqls = append(slowSqls, filteredNewSlowSqls...)

		if len(newSlowSqls) < pageSize {
			break
		}
		pageNum++
	}

	mergedSlowSqls := mergeSQLsByFingerprint(slowSqls)
	if len(mergedSlowSqls) > 0 {
		if at.isFirstScrap() {
			err = at.persist.OverrideAuditPlanSQLs(at.ap.ID, at.convertSQLInfosToModelSQLs(mergedSlowSqls, now))
			if err != nil {
				at.logger.Errorf("save sqls to storage fail, error: %v", err)
				return
			}
		} else {
			err = at.persist.UpdateDefaultAuditPlanSQLs(at.ap.ID, at.convertSQLInfosToModelSQLs(mergedSlowSqls, now))
			if err != nil {
				at.logger.Errorf("save sqls to storage fail, error: %v", err)
				return
			}
		}
	}

	// update lastEndTime
	// 查询的起始时间为上一次查询到的最后一条慢语句的开始执行时间
	if len(slowSqls) > 0 {
		lastSlowSql := slowSqls[len(slowSqls)-1]
		at.lastEndTime = &lastSlowSql.executionStartTime
	}
}

// 因为查询的起始时间为上一次查询到的最后一条慢语句的executionStartTime（精确到秒），而查询起始时间只能精确到分钟，所以有可能还是会查询到上一次查询过的慢语句，需要将其过滤掉
func (at *aliRdsMySQLTask) filterSlowSqlsByExecutionTime(slowSqls []SqlFromAliCloud, executionTime time.Time) (res []SqlFromAliCloud) {
	for _, sql := range slowSqls {
		if !sql.executionStartTime.After(executionTime) {
			continue
		}
		res = append(res, sql)
	}
	return
}

func (at *aliRdsMySQLTask) isFirstScrap() bool {
	return at.lastEndTime == nil
}

type sqlInfo struct {
	counter          int
	fingerprint      string
	sql              string
	schema           string
	queryTimeSeconds int
	startTime        string
}

func mergeSQLsByFingerprint(sqls []SqlFromAliCloud) []sqlInfo {
	sqlInfos := []sqlInfo{}

	counter := map[string]int /*slice subscript*/ {}
	for _, sql := range sqls {
		fp := query.Fingerprint(sql.sql)
		if index, exist := counter[fp]; exist {
			sqlInfos[index].counter += 1
			sqlInfos[index].fingerprint = fp
			sqlInfos[index].sql = sql.sql
			sqlInfos[index].schema = sql.schema
		} else {
			sqlInfos = append(sqlInfos, sqlInfo{
				counter:     1,
				fingerprint: fp,
				sql:         sql.sql,
				schema:      sql.schema,
			})
			counter[fp] = len(sqlInfos) - 1
		}

	}
	return sqlInfos
}

func (at *aliRdsMySQLTask) Audit() (*AuditResultResp, error) {
	task := &model.Task{
		DBType: at.ap.DBType,
	}
	return at.baseTask.audit(task)
}

func (at *aliRdsMySQLTask) GetSQLs(args map[string]interface{}) ([]Head, []map[string] /* head name */ string, uint64, error) {
	return baseTaskGetSQLs(args, at.persist)
}

func (at *aliRdsMySQLTask) CreateClient(rdsPath string, accessKeyId *string, accessKeySecret *string) (_result *rds20140815.Client, _err error) {
	config := &openapi.Config{
		// 您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String(rdsPath)
	_result, _err = rds20140815.NewClient(config)
	return _result, _err
}

type SqlFromAliCloud struct {
	sql                string
	executionStartTime time.Time
	schema             string
}

func (at *aliRdsMySQLTask) convertSQLInfosToModelSQLs(sqls []sqlInfo, now time.Time) []*model.AuditPlanSQLV2 {
	return convertRawSlowSQLWitchFromSqlInfo(sqls, now)
}

func convertRawSlowSQLWitchFromSqlInfo(sqls []sqlInfo, now time.Time) []*model.AuditPlanSQLV2 {
	as := make([]*model.AuditPlanSQLV2, len(sqls))
	for i, sql := range sqls {
		modelInfo := fmt.Sprintf(`{"counter":%v,"last_receive_timestamp":"%v","schema":"%v"}`, sql.counter, now.Format(time.RFC3339), sql.schema)
		as[i] = &model.AuditPlanSQLV2{
			Fingerprint: sql.fingerprint,
			SQLContent:  sql.sql,
			Info:        []byte(modelInfo),
			Schema:      sql.schema,
		}
	}
	return as
}

type AliRdsMySQLSlowLogTask struct {
	*aliRdsMySQLTask
}

func NewAliRdsMySQLSlowLogTask(entry *logrus.Entry, ap *model.AuditPlan) Task {
	sqlCollector := newSQLCollector(entry, ap)
	a := &AliRdsMySQLSlowLogTask{}
	task := &aliRdsMySQLTask{
		sqlCollector: sqlCollector,
		lastEndTime:  nil,
		pullLogs:     a.pullSlowLogs,
	}
	sqlCollector.do = task.collectorDo
	a.aliRdsMySQLTask = task
	return a
}

// 查询内容范围是开始时间的0s到设置结束时间的0s，所以结束时间点的慢日志是查询不到的
// startTime和endTime对应的是慢语句的开始执行时间
func (at *AliRdsMySQLSlowLogTask) pullSlowLogs(client *rds20140815.Client, DBInstanId string, startTime, endTime time.Time, pageSize, pageNum int32) (sqls []SqlFromAliCloud, err error) {
	describeSlowLogRecordsRequest := &rds20140815.DescribeSlowLogRecordsRequest{
		DBInstanceId: tea.String(DBInstanId),
		StartTime:    tea.String(startTime.Format("2006-01-02T15:04Z")),
		EndTime:      tea.String(endTime.Format("2006-01-02T15:04Z")),
		PageSize:     tea.Int32(pageSize),
		PageNumber:   tea.Int32(pageNum),
	}

	runtime := &_util.RuntimeOptions{}
	response := &rds20140815.DescribeSlowLogRecordsResponse{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()

		var err error
		response, err = client.DescribeSlowLogRecordsWithOptions(describeSlowLogRecordsRequest, runtime)
		if err != nil {
			return err
		}
		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		errMsg := _util.AssertAsString(error.Message)
		return nil, fmt.Errorf("get slow log failed: %v", *errMsg)
	}

	sqls = make([]SqlFromAliCloud, len(response.Body.Items.SQLSlowRecord))
	for i, slowRecord := range response.Body.Items.SQLSlowRecord {
		execStartTime, err := time.Parse("2006-01-02T15:04:05Z", utils.NvlString(slowRecord.ExecutionStartTime))
		if err != nil {
			return nil, fmt.Errorf("parse execution-start-time failed: %v", err)
		}
		sqls[i] = SqlFromAliCloud{
			sql:                utils.NvlString(slowRecord.SQLText),
			executionStartTime: execStartTime,
		}
	}
	return sqls, nil
}

type AliRdsMySQLAuditLogTask struct {
	*aliRdsMySQLTask
}

func NewAliRdsMySQLAuditLogTask(entry *logrus.Entry, ap *model.AuditPlan) Task {
	sqlCollector := newSQLCollector(entry, ap)
	a := &AliRdsMySQLAuditLogTask{}
	task := &aliRdsMySQLTask{
		sqlCollector: sqlCollector,
		lastEndTime:  nil,
		pullLogs:     a.pullAuditLogs,
	}
	sqlCollector.do = task.collectorDo
	a.aliRdsMySQLTask = task
	return a
}

// 查询内容范围是开始时间的0s到设置结束时间的0s，所以结束时间点的慢日志是查询不到的
// startTime和endTime对应的是慢语句的开始执行时间
func (at *AliRdsMySQLAuditLogTask) pullAuditLogs(client *rds20140815.Client, DBInstanId string, startTime, endTime time.Time, pageSize, pageNum int32) (sqls []SqlFromAliCloud, err error) {
	describeSQLLogRecordsRequest := &rds20140815.DescribeSQLLogRecordsRequest{
		ClientToken:  tea.String(time.Now().String()),
		DBInstanceId: tea.String(DBInstanId),
		StartTime:    tea.String(startTime.Format("2006-01-02T15:04:05Z")),
		EndTime:      tea.String(endTime.Format("2006-01-02T15:04:05Z")),
		PageSize:     tea.Int32(pageSize),
		PageNumber:   tea.Int32(pageNum),
	}
	runtime := &_util.RuntimeOptions{}
	response := &rds20140815.DescribeSQLLogRecordsResponse{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()

		var err error
		response, err = client.DescribeSQLLogRecordsWithOptions(describeSQLLogRecordsRequest, runtime)
		if err != nil {
			return err
		}
		return nil
	}()
	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		errMsg := _util.AssertAsString(error.Message)
		return nil, fmt.Errorf("get audit log failed: %v", *errMsg)
	}

	sqls = make([]SqlFromAliCloud, len(response.Body.Items.SQLRecord))
	for i, slowRecord := range response.Body.Items.SQLRecord {
		execStartTime, err := time.Parse("2006-01-02T15:04:05Z", utils.NvlString(slowRecord.ExecuteTime))
		if err != nil {
			return nil, fmt.Errorf("parse execution-start-time failed: %v", err)
		}
		sqls[i] = SqlFromAliCloud{
			sql:                utils.NvlString(slowRecord.SQLText),
			executionStartTime: execStartTime,
		}
	}
	return sqls, nil
}

type MySQLProcesslistTask struct {
	*sqlCollector
}

func (at *MySQLProcesslistTask) processlistSQL() string {
	sql := `
SELECT DISTINCT db,time,info
FROM information_schema.processlist
WHERE ID != connection_id() AND info != '' AND db NOT IN ('information_schema','performance_schema','mysql','sys')
%v
`
	whereSqlMinSecond := ""
	{
		sqlMinSecond := at.ap.Params.GetParam(paramKeySQLMinSecond).Int()
		if sqlMinSecond > 0 {
			whereSqlMinSecond = fmt.Sprintf("AND TIME > %d", sqlMinSecond)
		}
	}
	sql = fmt.Sprintf(sql, whereSqlMinSecond)
	return sql
}

func (at *MySQLProcesslistTask) Audit() (*AuditResultResp, error) {
	return auditWithSchema(at.logger, at.persist, at.ap)
}

// HACK:  SQLs may be executed in different Schemas.
// Before auditing sql, we need to insert a SCHEMA SWITCHING statement.
// And need to manually execute server.ReplenishTaskStatistics() to recalculate
// real task object score
func auditWithSchema(l *logrus.Entry, persist *model.Storage, ap *model.AuditPlan) (
	*AuditResultResp, error) {

	auditPlanSQLs, err := persist.GetAuditPlanSQLs(ap.ID)
	if err != nil {
		return nil, err
	}

	if len(auditPlanSQLs) == 0 {
		return nil, errNoSQLInAuditPlan
	}

	filteredSqls, err := filterSQLsByPeriod(ap.Params, auditPlanSQLs)
	if err != nil {
		return nil, err
	}

	if len(filteredSqls) == 0 {
		return nil, errNoSQLNeedToBeAudited
	}

	task := &model.Task{Instance: ap.Instance, DBType: ap.DBType, Schema: ap.InstanceDatabase}
	vTask := &model.Task{Instance: ap.Instance, DBType: ap.DBType, Schema: ap.InstanceDatabase}

	for i, sql := range filteredSqls {
		sqlItem := &model.ExecuteSQL{
			BaseSQL: model.BaseSQL{
				Number:  uint(i),
				Content: sql.SQLContent,
				Schema:  sql.Schema,
			},
		}
		{
			info := struct {
				AuditSchema string `json:"AuditSchema"`
			}{}
			err := json.Unmarshal(sql.Info, &info)
			if err != nil {
				return nil, fmt.Errorf("parse schema failed: %v", err)
			}
			if info.AuditSchema != "" {
				vTask.ExecuteSQLs = append(vTask.ExecuteSQLs, &model.ExecuteSQL{
					BaseSQL: model.BaseSQL{
						Content: fmt.Sprintf("USE %s;", info.AuditSchema),
					},
				})
			}
		}
		task.ExecuteSQLs = append(task.ExecuteSQLs, sqlItem)
		vTask.ExecuteSQLs = append(vTask.ExecuteSQLs, sqlItem) // vTask is a copy of task for schema switch

	}

	err = server.Audit(l, vTask, &ap.ProjectId, ap.RuleTemplateName)
	if err != nil {
		return nil, err
	}
	// replenish task statistics manually for real task
	server.ReplenishTaskStatistics(task)

	return &AuditResultResp{
		AuditPlanID:  uint64(ap.ID),
		Task:         task,
		FilteredSqls: filteredSqls,
	}, nil
}

func NewMySQLProcesslistTask(entry *logrus.Entry, ap *model.AuditPlan) Task {
	sqlCollector := newSQLCollector(entry, ap)
	task := &MySQLProcesslistTask{
		sqlCollector,
	}
	sqlCollector.do = task.collectorDo
	sqlCollector.loopInterval = func() time.Duration {
		interval := ap.Params.GetParam(paramKeyCollectIntervalSecond).Int()
		if interval == 0 {
			interval = 60
		}
		return time.Second * time.Duration(interval)
	}

	return task
}

func (at *MySQLProcesslistTask) collectorDo() {
	if at.ap.InstanceName == "" {
		at.logger.Warnf("instance is not configured")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	instance, _, err := dms.GetInstanceInProjectByName(ctx, string(at.ap.ProjectId), at.ap.InstanceName)

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

	res, err := db.Db.Query(at.processlistSQL())
	if err != nil {
		at.logger.Errorf("show processlist failed, error: %v", err)
		return
	}

	if len(res) == 0 {
		return
	}

	sqls := make([]SqlFromAliCloud, len(res))
	for i := range res {
		sqls[i] = SqlFromAliCloud{sql: res[i]["info"].String, schema: res[i]["db"].String}
	}

	sqlInfos := mergeSQLsByFingerprint(sqls)

	if len(sqlInfos) > 0 {
		err = at.persist.UpdateDefaultAuditPlanSQLs(at.ap.ID,
			convertRawSlowSQLWitchFromSqlInfo(sqlInfos, time.Now()))
		if err != nil {
			at.logger.Errorf("save processlist to storage fail, error: %v", err)
			return
		}
	}
}

type BaiduRdsMySQLSlowLogTask struct {
	*baiduRdsMySQLTask
}

type SqlFromBaiduCloud struct {
	sql                string
	executionStartTime time.Time
	schema             string
}

type baiduRdsMySQLTask struct {
	*sqlCollector
	lastEndTime *time.Time
	pullLogs    func(client *rds.Client, DBInstanceId string, startTime, endTime time.Time, pageSize, pageNum int32) (sqlList []SqlFromBaiduCloud, err error)
}

func (bt *baiduRdsMySQLTask) Audit() (*AuditResultResp, error) {
	task := &model.Task{
		DBType: bt.ap.DBType,
	}
	return bt.baseTask.audit(task)
}

func (bt *baiduRdsMySQLTask) collectorDo() {
	if bt.ap.InstanceName == "" {
		bt.logger.Warnf("instance is not configured")
		return
	}

	rdsDBInstanceId := bt.ap.Params.GetParam(paramKeyDBInstanceId).String()
	if rdsDBInstanceId == "" {
		bt.logger.Warnf("db instance id is not configured")
		return
	}

	accessKeyID := bt.ap.Params.GetParam(paramKeyAccessKeyId).String()
	if accessKeyID == "" {
		bt.logger.Warnf("access key id is not configured")
		return
	}

	accessKeySecret := bt.ap.Params.GetParam(paramKeyAccessKeySecret).String()
	if accessKeySecret == "" {
		bt.logger.Warnf("access key secret is not configured")
		return
	}

	rdsPath := bt.ap.Params.GetParam(paramKeyRdsPath).String()
	if rdsPath == "" {
		bt.logger.Warnf("rds path is not configured")
		return
	}

	firstScrapInLastHours := bt.ap.Params.GetParam(paramKeyFirstSqlsScrappedInLastPeriodHours).Int()
	if firstScrapInLastHours == 0 {
		firstScrapInLastHours = 24
	}

	theMaxSupportedDays := 7 // 支持往前查看慢日志的最大天数
	hoursDuringADay := 24
	if firstScrapInLastHours > theMaxSupportedDays*hoursDuringADay {
		bt.logger.Warnf("Can not get slow logs from so early time. firstScrapInLastHours=%v", firstScrapInLastHours)
		return
	}

	client, err := bt.CreateClient(rdsPath, accessKeyID, accessKeySecret)
	if err != nil {
		bt.logger.Warnf("create client failed, error: %v", err)
		return
	}

	// use UTC time
	// https://cloud.baidu.com/doc/RDS/s/Ejwvz0uoq
	now := time.Now().UTC()
	var startTime time.Time
	if bt.isFirstScrap() {
		startTime = now.Add(time.Duration(-1*firstScrapInLastHours) * time.Hour)
	} else {
		startTime = *bt.lastEndTime
	}

	var pageNum int32 = 1
	pageSize := 100
	//nolint:staticcheck
	slowSqlList := make([]SqlFromBaiduCloud, 0)
	for {
		// 每次取100条
		slowSqlListFromBaiduRds, err := bt.pullLogs(client, rdsDBInstanceId, startTime, now, int32(pageSize), pageNum)
		if err != nil {
			bt.logger.Warnf("pull slow logs failed, error: %v", err)
			return
		}

		slowSqlList = append(slowSqlList, slowSqlListFromBaiduRds...)
		if len(slowSqlListFromBaiduRds) < pageSize {
			break
		}
		pageNum++
	}

	mergedSlowSqlList := mergerSqlListByFingerprint(slowSqlList)
	if len(mergedSlowSqlList) > 0 {
		if bt.isFirstScrap() {
			err = bt.persist.OverrideAuditPlanSQLs(bt.ap.ID, bt.convertSQLInfoListToModelSQLList(mergedSlowSqlList, now))
			if err != nil {
				bt.logger.Warnf("save slow logs to storage failed, error: %v", err)
				return
			}
		} else {
			err = bt.persist.UpdateDefaultAuditPlanSQLs(bt.ap.ID, bt.convertSQLInfoListToModelSQLList(mergedSlowSqlList, now))
			if err != nil {
				bt.logger.Warnf("save slow logs to storage failed, error: %v", err)
				return
			}
		}
	}

	if len(slowSqlList) > 0 {
		lastSlowSql := slowSqlList[len(slowSqlList)-1]
		bt.lastEndTime = &lastSlowSql.executionStartTime
	}
}

func mergerSqlListByFingerprint(sqlList []SqlFromBaiduCloud) []sqlInfo {
	sqlInfoList := make([]sqlInfo, 0)
	counter := map[string] /*finger*/ int /*slice subscript*/ {}
	for _, sql := range sqlList {
		fp := query.Fingerprint(sql.sql)
		if index, exist := counter[fp]; exist {
			sqlInfoList[index].counter += 1
			sqlInfoList[index].fingerprint = fp
			sqlInfoList[index].sql = sql.sql
			sqlInfoList[index].schema = sql.schema
		} else {
			sqlInfoList = append(sqlInfoList, sqlInfo{
				counter:     1,
				fingerprint: fp,
				sql:         sql.sql,
				schema:      sql.schema,
			})
			counter[fp] = len(sqlInfoList) - 1
		}
	}
	return sqlInfoList
}

func (bt *baiduRdsMySQLTask) isFirstScrap() bool {
	return bt.lastEndTime == nil
}

func (bt *baiduRdsMySQLTask) convertSQLInfoListToModelSQLList(list []sqlInfo, now time.Time) []*model.AuditPlanSQLV2 {
	return convertRawSlowSQLWitchFromSqlInfo(list, now)
}

func (bt *baiduRdsMySQLTask) CreateClient(rdsPath string, accessKeyID string, accessKeySecret string) (*rds.Client, error) {
	client, err := rds.NewClient(accessKeyID, accessKeySecret, rdsPath)
	if err != nil {
		return nil, err
	}
	return client, nil
}

type slowLogReq struct {
	StartTime string   `json:"startTime"`
	EndTime   string   `json:"endTime"`
	PageNo    int      `json:"pageNo"`
	PageSize  int      `json:"pageSize"`
	DbName    []string `json:"dbName"`
}

type slowLogDetail struct {
	InstanceId   string  `json:"instanceId"`
	UserName     string  `json:"userName"`
	DbName       string  `json:"dbName"`
	HostIp       string  `json:"hostIp"`
	QueryTime    float64 `json:"queryTime"`
	LockTime     float64 `json:"lockTime"`
	RowsExamined int     `json:"rowsExamined"`
	RowsSent     int     `json:"rowsSent"`
	Sql          string  `json:"sql"`
	ExecuteTime  string  `json:"executeTime"`
}

type slowLogDetailResp struct {
	Count    int             `json:"count"`
	SlowLogs []slowLogDetail `json:"slowLogs"`
}

func (t BaiduRdsMySQLSlowLogTask) pullSlowLogs(client *rds.Client, instanceID string, startTime time.Time, endTime time.Time, size int32, num int32) (sqlList []SqlFromBaiduCloud, err error) {
	req := slowLogReq{
		StartTime: startTime.Format("2006-01-02T15:04:05Z"),
		EndTime:   endTime.Format("2006-01-02T15:04:05Z"),
		PageNo:    int(num),
		PageSize:  int(size),
		DbName:    nil,
	}

	uri := getRdsUriWithInstanceId(instanceID)

	resp := new(slowLogDetailResp)
	err = bce.NewRequestBuilder(client).
		WithMethod(http.MethodPost).
		WithURL(uri).
		WithBody(&req).
		WithResult(&resp).
		Do()
	if err != nil {
		return nil, err
	}

	sqlList = make([]SqlFromBaiduCloud, len(resp.SlowLogs))
	for i, slowLog := range resp.SlowLogs {
		execStartTime, err := time.Parse("2006-01-02 15:04:05", slowLog.ExecuteTime)
		if err != nil {
			return nil, err
		}

		sqlList[i] = SqlFromBaiduCloud{
			sql:                slowLog.Sql,
			executionStartTime: execStartTime,
		}
	}

	return sqlList, nil
}

func getRdsUriWithInstanceId(instanceId string) string {
	return rds.URI_PREFIX + rds.REQUEST_RDS_URL + "/" + instanceId + "/slowlogs/details"
}

func NewBaiduRdsMySQLSlowLogTask(entry *logrus.Entry, ap *model.AuditPlan) Task {
	sqlCollector := newSQLCollector(entry, ap)
	b := &BaiduRdsMySQLSlowLogTask{}
	task := &baiduRdsMySQLTask{
		sqlCollector: sqlCollector,
		lastEndTime:  nil,
		pullLogs:     b.pullSlowLogs,
	}
	sqlCollector.do = task.collectorDo
	b.baiduRdsMySQLTask = task

	return b
}
