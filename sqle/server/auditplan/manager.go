package auditplan

import (
	"context"
	"fmt"
	"time"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/notification"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type AuditPlan struct {
	ID                  uint
	ProjectId           string
	Name                string
	CronExpression      string
	DBType              string
	Token               string
	InstanceID          string
	CreateUserID        string
	Type                string
	RuleTemplateName    string
	Params              params.Params
	InstanceAuditPlanId uint

	Instance *model.Instance
}

// Deprecated
func ConvertModelToAuditPlan(a *model.AuditPlan) *AuditPlan {
	return &AuditPlan{
		ID:               a.ID,
		ProjectId:        string(a.ProjectId),
		Name:             a.Name,
		CronExpression:   a.CronExpression,
		DBType:           a.DBType,
		Token:            a.Token,
		InstanceID:       a.InstanceName,
		CreateUserID:     a.CreateUserID,
		Type:             a.Type,
		RuleTemplateName: a.RuleTemplateName,
		Params:           a.Params,
		Instance:         a.Instance,
	}
}

func ConvertModelToAuditPlanV2(a *model.AuditPlanDetail) *AuditPlan {
	return &AuditPlan{
		ID:                  a.ID,
		ProjectId:           a.ProjectId,
		DBType:              a.DBType,
		Token:               a.Token,
		InstanceID:          a.InstanceID,
		CreateUserID:        a.CreateUserID,
		Type:                a.Type,
		RuleTemplateName:    a.RuleTemplateName,
		Params:              a.Params,
		InstanceAuditPlanId: a.InstanceAuditPlanID,
		Instance:            a.Instance,
	}
}

var ErrAuditPlanNotExist = errors.New(errors.DataNotExist, fmt.Errorf("audit plan not exist"))
var ErrAuditPlanExisted = errors.New(errors.DataExist, fmt.Errorf("audit plan existed"))

func auditV2(auditPlanId uint, task Task) (*model.AuditPlanReportV2, error) {
	auditResultResp, err := task.Audit()
	if err != nil {
		return nil, err
	}

	taskResp := auditResultResp.Task
	auditPlanReport := &model.AuditPlanReportV2{
		AuditPlanID: uint(auditResultResp.AuditPlanID),
		PassRate:    taskResp.PassRate,
		Score:       taskResp.Score,
		AuditLevel:  taskResp.AuditLevel,
	}

	for i, executeSQL := range taskResp.ExecuteSQLs {
		auditPlanReport.AuditPlanReportSQLs = append(auditPlanReport.AuditPlanReportSQLs, &model.AuditPlanReportSQLV2{
			SQL:          executeSQL.Content,
			Number:       uint(i + 1),
			AuditResults: executeSQL.AuditResults,
			Schema:       executeSQL.Schema,
		})
	}

	s := model.GetStorage()
	err = s.Save(auditPlanReport)
	if err != nil {
		return nil, err
	}

	return auditPlanReport, notification.NotifyAuditPlan(auditPlanId, auditPlanReport)
}

func Audit(entry *logrus.Entry, ap *AuditPlan) (*model.AuditPlanReportV2, error) {
	task := NewTask(entry, ap)
	return auditV2(ap.ID, task)
}

func UploadSQLsV2(entry *logrus.Entry, ap *AuditPlan, sqls []*SQL) error {
	task := NewTask(entry, ap)
	return task.FullSyncSQLs(sqls)
}

func UploadSQLs(entry *logrus.Entry, ap *AuditPlan, sqls []*SQL, isPartialSync bool) error {
	task := NewTask(entry, ap)
	if isPartialSync {
		return task.PartialSyncSQLs(sqls)
	} else {
		return task.FullSyncSQLs(sqls)
	}
}

// todo: 弃用
func GetSQLs(entry *logrus.Entry, ap *AuditPlan, args map[string]interface{}) ([]Head, []map[string] /* head name */ string, uint64, error) {
	return []Head{}, []map[string]string{}, 0, nil
}

func GetSQLHead(ap *AuditPlan, persist *model.Storage) ([]Head, error) {
	meta, err := GetMeta(ap.Type)
	if err != nil {
		return nil, err
	}
	return meta.Handler.Head(ap), nil
}

func GetSQLFilterMeta(ctx context.Context, ap *AuditPlan, persist *model.Storage) ([]FilterMeta, error) {
	meta, err := GetMeta(ap.Type)
	if err != nil {
		return nil, err
	}
	return meta.Handler.Filters(ctx, log.NewEntry(), ap, persist), nil
}

func GetSQLData(ctx context.Context, ap *AuditPlan, persist *model.Storage, filters []Filter, orderBy string, isAsc bool, limit, offset int) ([]map[string] /* head name */ string, uint64, error) {
	meta, err := GetMeta(ap.Type)
	if err != nil {
		return nil, 0, err
	}
	return meta.Handler.GetSQLData(ctx, ap, persist, filters, orderBy, isAsc, limit, offset)
}

func init() {
	server.OnlyRunOnLeaderJobs = append(server.OnlyRunOnLeaderJobs, NewManager, NewAuditPlanHandlerJob, NewAuditPlanAggregateSQLJob)
}

func NewManager(entry *logrus.Entry) server.ServerJob {
	now := time.Now()
	manager := &Manager{
		scheduler: &scheduler{
			cron:     cron.New(),
			entryIDs: make(map[uint]cron.EntryID),
		},
		persist:      model.GetStorage(),
		logger:       entry.WithField("job", "audit_plan"),
		tasks:        map[uint]Task{},
		lastSyncTime: &now,
		exitCh:       make(chan struct{}),
		doneCh:       make(chan struct{}),
	}
	return manager
}

// Manager is the struct managing the persistent AuditPlans.
type Manager struct {
	// v3.2407 Deprecated
	scheduler *scheduler

	// persist is a database handle which store AuditPlan.
	persist *model.Storage

	logger *logrus.Entry

	tasks map[uint] /* audit plan id*/ Task

	lastSyncTime   *time.Time
	isFullSyncDone bool

	exitCh chan struct{}
	doneCh chan struct{}
}

func (mgr *Manager) Start() {
	// 删除定时审核扫描任务定时任务，后续采集结束会直接审核保存
	// mgr.scheduler.start()
	mgr.logger.Infoln("audit plan manager started")

	go func() {
		tick := time.NewTicker(5 * time.Second)
		defer tick.Stop()
		for {
			select {
			case <-tick.C:
				err := mgr.sync()
				if err != nil {
					mgr.logger.Errorf("sync audit plan task failed, error: %v", err)
				}
			case <-mgr.exitCh:
				mgr.doneCh <- struct{}{}
				return
			}
		}
	}()
}

func (mgr *Manager) sync() error {
	// 全量同步智能扫描任务，仅需成功做一次
	if !mgr.isFullSyncDone {
		aps, err := dms.ListActiveAuditPlansWithInstanceV2(mgr.persist.ListActiveAuditPlanDetail)
		if err != nil {
			return err
		}
		mgr.isFullSyncDone = true
		for _, v := range aps {
			ap := v

			err := mgr.startAuditPlan(ConvertModelToAuditPlanV2(ap))
			if err != nil {
				mgr.logger.WithField("name", ap.Type).Errorf("start audit task failed, error: %v", err)
			}
		}
	}
	// 增量同步智能扫描任务，根据数据库记录的更新时间筛选，更新后将下次筛选的时间为上一次记录的最晚的更新时间。
	aps, err := mgr.persist.GetLatestAuditPlanRecordsV2(*mgr.lastSyncTime)
	if err != nil {
		return err
	}

	for _, v := range aps {
		ap := v
		err := mgr.syncTask(ap.ID)
		if err != nil {
			mgr.logger.WithField("id", ap.ID).Errorf("sync audit task failed, error: %v", err)
		}
		mgr.lastSyncTime = &ap.UpdatedAt
	}
	return nil
}

func (mgr *Manager) Stop() {
	mgr.exitCh <- struct{}{}
	<-mgr.doneCh

	for name := range mgr.tasks {
		err := mgr.deleteAuditPlan(name)
		if err != nil {
			mgr.logger.WithField("name", name).Errorf("stop audit task failed, error: %v", err)
		}
	}
	ctx := mgr.scheduler.stop()
	<-ctx.Done()
	mgr.logger.Infoln("audit plan manager stopped")
}

func (mgr *Manager) syncTask(auditPlanId uint) error {
	ap, exist, err := mgr.persist.GetActiveAuditPlanDetail(auditPlanId)
	if err != nil {
		return err
	}
	if !exist {
		return mgr.deleteAuditPlan(auditPlanId)
	} else {
		return mgr.startAuditPlan(ConvertModelToAuditPlanV2(ap))
	}
}

func (mgr *Manager) startAuditPlan(ap *AuditPlan) error {
	task, ok := mgr.tasks[ap.ID]
	if ok {
		err := task.Stop()
		if err != nil {
			return err
		}
	}

	task = NewTask(mgr.logger, ap)
	if err := task.Start(); err != nil {
		return err
	}
	mgr.tasks[ap.ID] = task

	return nil
}

func (mgr *Manager) deleteAuditPlan(id uint) error {
	if mgr.scheduler.hasJob(id) {
		err := mgr.scheduler.removeJob(mgr.logger, id)
		if err != nil {
			return err
		}
	}
	task, ok := mgr.tasks[id]
	if ok {
		err := task.Stop()
		if err != nil {
			return err
		}
		delete(mgr.tasks, id)
	}
	return nil
}

func (mgr *Manager) getTask(auditPlanId uint) (Task, error) {
	task, ok := mgr.tasks[auditPlanId]
	if !ok {
		return nil, errors.New(errors.DataNotExist, fmt.Errorf("task not found"))
	}
	return task, nil
}

// scheduler is not goroutine safe.
type scheduler struct {
	// cron is a AuditPlan scheduler.
	cron *cron.Cron

	// entryIDs maps audit plan name to it's job entry ID.
	entryIDs map[uint] /* audit plan id*/ cron.EntryID
}

func (s *scheduler) removeJob(entry *logrus.Entry, auditPlanId uint) error {
	entryID, ok := s.entryIDs[auditPlanId]
	if !ok {
		return ErrAuditPlanNotExist
	}

	s.cron.Remove(entryID)
	delete(s.entryIDs, auditPlanId)

	entry.WithFields(logrus.Fields{
		"id": auditPlanId,
	}).Infoln("stop audit scheduler")
	return nil
}

func (s *scheduler) addJob(entry *logrus.Entry, ap *AuditPlan, do func()) error {
	_, ok := s.entryIDs[ap.ID]
	if ok {
		return ErrAuditPlanExisted
	}

	entryID, err := s.cron.AddFunc(ap.CronExpression, do)
	if err != nil {
		return err
	}

	s.entryIDs[ap.ID] = entryID

	entry.WithFields(logrus.Fields{
		"id":              ap.ID,
		"name":            ap.Name,
		"cron_expression": ap.CronExpression,
	}).Infoln("start audit scheduler")
	return nil
}

func (s *scheduler) start() {
	s.cron.Start()
}

func (s *scheduler) stop() context.Context {
	return s.cron.Stop()
}

func (s *scheduler) hasJob(auditPlanId uint) bool {
	_, has := s.entryIDs[auditPlanId]
	return has
}
