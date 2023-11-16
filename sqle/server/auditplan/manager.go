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
	"github.com/actiontech/sqle/sqle/server"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var ErrAuditPlanNotExist = errors.New(errors.DataNotExist, fmt.Errorf("audit plan not exist"))
var ErrAuditPlanExisted = errors.New(errors.DataExist, fmt.Errorf("audit plan existed"))

func audit(auditPlanId uint, task Task) (*model.AuditPlanReportV2, error) {
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

	go func() {
		syncFromAuditPlan := NewSyncFromAuditPlan(auditPlanReport, auditResultResp.FilteredSqls, taskResp)
		if err := syncFromAuditPlan.SyncSqlManager(); err != nil {
			log.NewEntry().WithField("name", auditResultResp.AuditPlanID).Errorf("schedule to save sql manage failed, error: %v", err)
		}
	}()

	return auditPlanReport, notification.NotifyAuditPlan(auditPlanId, auditPlanReport)
}

func Audit(entry *logrus.Entry, ap *model.AuditPlan) (*model.AuditPlanReportV2, error) {
	task := NewTask(entry, ap)
	return audit(ap.ID, task)
}

func UploadSQLs(entry *logrus.Entry, ap *model.AuditPlan, sqls []*SQL, isPartialSync bool) error {
	go func() {
		err := SyncToSqlManage(sqls, ap)
		if err != nil {
			log.NewEntry().WithField("name", ap.Name).Errorf("schedule to save sql manage failed, error: %v", err)
		}
	}()

	task := NewTask(entry, ap)
	if isPartialSync {
		return task.PartialSyncSQLs(sqls)
	} else {
		return task.FullSyncSQLs(sqls)
	}
}

func GetSQLs(entry *logrus.Entry, ap *model.AuditPlan, args map[string]interface{}) ([]Head, []map[string] /* head name */ string, uint64, error) {
	args["audit_plan_id"] = ap.ID
	task := NewTask(entry, ap)
	return task.GetSQLs(args)
}

func init() {
	server.OnlyRunOnLeaderJobs = append(server.OnlyRunOnLeaderJobs, NewManager)
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
	mgr.scheduler.start()
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
		aps, err := dms.GetActiveAuditPlansWithInstance(mgr.persist.GetActiveAuditPlans)
		if err != nil {
			return err
		}
		mgr.isFullSyncDone = true
		for _, v := range aps {
			ap := v

			err := mgr.startAuditPlan(ap)
			if err != nil {
				mgr.logger.WithField("name", ap.Name).Errorf("start audit task failed, error: %v", err)
			}
		}
	}
	// 增量同步智能扫描任务，根据数据库记录的更新时间筛选，更新后将下次筛选的时间为上一次记录的最晚的更新时间。
	aps, err := mgr.persist.GetLatestAuditPlanRecords(*mgr.lastSyncTime)
	if err != nil {
		return err
	}

	for _, v := range aps {
		ap := v
		err := mgr.syncTask(ap.ID)
		if err != nil {
			mgr.logger.WithField("name", ap.Name).Errorf("sync audit task failed, error: %v", err)
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
	ap, exist, err := mgr.persist.GetActiveAuditPlanById(auditPlanId)
	if err != nil {
		return err
	}
	if !exist {
		return mgr.deleteAuditPlan(auditPlanId)
	} else {
		return mgr.startAuditPlan(ap)
	}
}

func (mgr *Manager) startAuditPlan(ap *model.AuditPlan) error {
	if mgr.scheduler.hasJob(ap.ID) {
		err := mgr.scheduler.removeJob(mgr.logger, ap.ID)
		if err != nil {
			return err
		}
	}
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

	return mgr.scheduler.addJob(mgr.logger, ap, func() {
		_, err := audit(ap.ID, task)
		if err != nil {
			mgr.logger.WithField("name", ap.Name).Errorf("schedule to audit task failed, error: %v", err)
		}
	})
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

func (s *scheduler) addJob(entry *logrus.Entry, ap *model.AuditPlan, do func()) error {
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
