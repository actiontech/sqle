package auditplan

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/actiontech/sqle/sqle/log"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
	pkgerrors "github.com/actiontech/sqle/sqle/pkg/errors"
	"github.com/actiontech/sqle/sqle/pkg/oracle"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/sirupsen/logrus"
)

var errNoSQLInAuditPlan = errors.New(errors.DataConflict, fmt.Errorf("there is no SQLs in audit plan"))

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

func NewTask(entry *logrus.Entry, ap *model.AuditPlan) (Task, error) {
	entry = entry.WithField("name", ap.Name)
	switch ap.Type {
	case TypeMySQLSchemaMeta:
		return NewSchemaMetaTask(entry, ap), nil
	case TypeOracleTopSQL:
		return NewOracleTopSQLTask(entry, ap)
	default:
		return NewDefaultTask(entry, ap), nil
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

	for i, sql := range auditPlanSQLs {
		task.ExecuteSQLs = append(task.ExecuteSQLs, &model.ExecuteSQL{
			BaseSQL: model.BaseSQL{
				Number:  uint(i),
				Content: sql.SQLContent,
			},
		})
	}

	err = server.Audit(at.logger, task)
	if err != nil {
		return nil, err
	}

	auditPlanReport := &model.AuditPlanReportV2{AuditPlanID: at.ap.ID}
	for _, executeSQL := range task.ExecuteSQLs {
		auditPlanReport.AuditPlanReportSQLs = append(auditPlanReport.AuditPlanReportSQLs, &model.AuditPlanReportSQLV2{
			SQL:         executeSQL.Content,
			AuditResult: executeSQL.AuditResult,
		})
	}

	err = at.persist.Save(auditPlanReport)
	if err != nil {
		return nil, err
	}
	return auditPlanReport, nil
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
		at.looper(at.cancel)
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

func (at *sqlCollector) looper(cancel chan struct{}) {
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
		Host:         instance.Host,
		Port:         instance.Port,
		User:         instance.User,
		Password:     instance.Password,
		DatabaseName: at.ap.InstanceDatabase},
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
		sql, err := db.ShowCreateTable(table)
		if err != nil {
			at.logger.Errorf("show create table fail, error: %v", err)
			return
		}
		sqls = append(sqls, sql)
	}
	for _, view := range views {
		sql, err := db.ShowCreateView(view)
		if err != nil {
			at.logger.Errorf("show create table fail, error: %v", err)
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

	db *oracle.DB
}

func NewOracleTopSQLTask(entry *logrus.Entry, ap *model.AuditPlan) (*OracleTopSQLTask, error) {
	inst := ap.Instance
	if inst == nil {
		return nil, fmt.Errorf("instance is not configured")
	}

	dsn := &oracle.DSN{
		Host:        inst.Host,
		Port:        inst.Port,
		User:        inst.User,
		Password:    inst.Password,
		ServiceName: ap.InstanceDatabase,
	}
	db, err := oracle.NewDB(dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to instance fail, error: %v", err)
	}

	task := &OracleTopSQLTask{
		sqlCollector: newSQLCollector(entry, ap),
		db:           db,
	}
	task.sqlCollector.do = task.collectorDo
	return task, nil
}

func (at *OracleTopSQLTask) collectorDo() {
	inst := at.ap.Instance
	if inst == nil {
		at.logger.Warnf("instance is not configured")
		return
	}

	select {
	case <-at.cancel:
		at.logger.Info("cancel task")
		return
	default:
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sqls, err := at.db.QueryTopSQLs(ctx, at.ap.Params.GetParam("top_n").Int())
	if err != nil {
		at.logger.Errorf("query top sql fail, error: %v", err)
		return
	}
	if len(sqls) > 0 {
		err = at.persist.OverrideAuditPlanSQLs(at.ap.Name, convertRawSQLToModelSQLs(sqls))
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
	}
	rows := make([]map[string]string, 0, len(auditPlanSQLs))
	for _, sql := range auditPlanSQLs {
		rows = append(rows, map[string]string{
			"sql": sql.SQLContent,
		})
	}
	return heads, rows, count, nil
}

func (at *OracleTopSQLTask) Stop() error {
	var errs []error
	if err := at.sqlCollector.Stop(); err != nil {
		errs = append(errs, err)
	}
	if err := at.db.Close(); err != nil {
		errs = append(errs, err)
	}
	return pkgerrors.Combine(errs...)
}
