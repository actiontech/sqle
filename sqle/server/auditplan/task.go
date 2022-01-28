package auditplan

import (
	"fmt"
	"sync"
	"time"

	"github.com/actiontech/sqle/sqle/log"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/sirupsen/logrus"
)

var errNoSQLInAuditPlan = errors.New(errors.DataConflict, fmt.Errorf("there is no SQLs in audit plan"))

type Task interface {
	Start() error
	Stop() error
	Audit() (*model.AuditPlanReportV2, error)
}

func NewTask(entry *logrus.Entry, ap *model.AuditPlan) Task {
	entry = entry.WithField("name", ap.Name)
	switch ap.Type {
	case TypeMySQLSchemaMeta:
		return NewSchemaMetaTask(entry, ap)
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

	for i, sql := range auditPlanSQLs {
		task.ExecuteSQLs = append(task.ExecuteSQLs, &model.ExecuteSQL{
			BaseSQL: model.BaseSQL{
				Number:  uint(i),
				Content: sql.LastSQL,
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

type runnerTask struct {
	*baseTask
	sync.WaitGroup
	isStarted bool
	cancel    chan struct{}
	runnerFn  func(chan struct{})
}

func newRunnerTask(entry *logrus.Entry, ap *model.AuditPlan) *runnerTask {
	return &runnerTask{
		newBaseTask(entry, ap),
		sync.WaitGroup{},
		false,
		make(chan struct{}),
		func(cancel chan struct{}) { // default runner
			<-cancel
		},
	}
}

func (at *runnerTask) Start() error {
	if at.isStarted {
		return nil
	}
	go func() {
		at.WaitGroup.Add(1)
		at.isStarted = true
		at.logger.Infof("start task")
		at.runnerFn(at.cancel)
		at.WaitGroup.Done()
	}()
	return nil
}

func (at *runnerTask) Stop() error {
	if !at.isStarted {
		return nil
	}
	at.cancel <- struct{}{}
	at.isStarted = false
	at.WaitGroup.Wait()
	at.logger.Infof("stop task")
	return nil
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

type SchemaMetaTask struct {
	*runnerTask
}

func NewSchemaMetaTask(entry *logrus.Entry, ap *model.AuditPlan) *SchemaMetaTask {
	runnerTask := newRunnerTask(entry, ap)
	task := &SchemaMetaTask{
		runnerTask,
	}
	task.runnerFn = task.runner
	return task
}

func (at *SchemaMetaTask) runner(cancel chan struct{}) {
	interval := at.ap.Params.GetParam("collect_interval_minute").Int()
	if interval == 0 {
		interval = 60
	}
	collectView := at.ap.Params.GetParam("collect_view").Bool()
	at.do(collectView)
	tk := time.Tick(time.Duration(interval) * time.Minute)
	select {
	case <-cancel:
		return
	case <-tk:
		at.logger.Infof("tick %s", at.ap.Name)
		at.do(collectView)
	}
}

func (at *SchemaMetaTask) do(CollectView bool) {
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
	if CollectView {
		views, err = db.ShowSchemaViews(at.ap.InstanceDatabase)
		if err != nil {
			at.logger.Errorf("get schema view fail, error: %v", err)
			return
		}
	}
	sqls := make([]*model.AuditPlanSQL, 0, len(tables)+len(views))
	for _, table := range tables {
		sql, err := db.ShowCreateTable(table)
		if err != nil {
			at.logger.Errorf("show create table fail, error: %v", err)
			return
		}
		sqls = append(sqls, &model.AuditPlanSQL{
			LastSQL:              sql,
			Fingerprint:          sql,
			Counter:              1,
			LastReceiveTimestamp: time.Now().String(),
		})
	}
	for _, view := range views {
		sql, err := db.ShowCreateView(view)
		if err != nil {
			at.logger.Errorf("show create table fail, error: %v", err)
			return
		}
		sqls = append(sqls, &model.AuditPlanSQL{
			LastSQL:              sql,
			Fingerprint:          sql,
			Counter:              1,
			LastReceiveTimestamp: time.Now().String(),
		})
	}
	if len(sqls) > 0 {
		err = at.persist.OverrideAuditPlanSQLs(at.ap.Name, sqls)
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
