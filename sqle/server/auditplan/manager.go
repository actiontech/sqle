package auditplan

import (
	"context"
	"fmt"
	"sync"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/notification"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var ErrAuditPlanNotExist = errors.New(errors.DataNotExist, fmt.Errorf("audit plan not exist"))
var ErrAuditPlanExisted = errors.New(errors.DataExist, fmt.Errorf("audit plan existed"))

var manager *Manager

func InitManager(s *model.Storage) chan struct{} {
	manager = &Manager{
		scheduler: &scheduler{
			cron:     cron.New(),
			entryIDs: make(map[uint]cron.EntryID),
		},
		persist: s,
		logger:  log.NewEntry().WithField("type", "audit_plan"),
		tasks:   map[uint]Task{},
	}

	err := manager.start()
	if err != nil {
		panic(err)
	}

	exitCh := make(chan struct{})

	go func() {
		<-exitCh
		manager.stop()
	}()

	return exitCh
}

func GetManager() *Manager {
	return manager
}

// Manager is the struct managing the persistent AuditPlans. It
// is *goroutine-safe*, since all exported methods are protected by a lock.
//
// All audit plan operations except select should go through Manager.
type Manager struct {
	mu sync.Mutex

	scheduler *scheduler

	// persist is a database handle which store AuditPlan.
	persist *model.Storage

	logger *logrus.Entry

	tasks map[uint] /* audit plan id*/ Task
}

func (mgr *Manager) start() error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	mgr.scheduler.start()
	mgr.logger.Infoln("audit plan manager started")

	aps, err := mgr.persist.GetAuditPlans()
	if err != nil {
		return err
	}
	for _, v := range aps {
		ap := v
		err := mgr.startAuditPlan(ap)
		if err != nil {
			mgr.logger.WithField("name", ap.Name).Errorf("start audit task failed, error: %v", err)
		}
	}
	return nil
}

func (mgr *Manager) stop() {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

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

func (mgr *Manager) SyncTask(auditPlanId uint) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	ap, exist, err := mgr.persist.GetAuditPlanById(auditPlanId)
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
		_, err := mgr.Audit(ap.ID)
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

func (mgr *Manager) Audit(auditPlanId uint) (*model.AuditPlanReportV2, error) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	task, err := mgr.getTask(auditPlanId)
	if err != nil {
		return nil, err
	}
	report, err := task.Audit()
	if err != nil {
		return nil, err
	}
	return report, notification.NotifyAuditPlan(auditPlanId, report)
}

func (mgr *Manager) UploadSQLs(auditPlanId uint, sqls []*SQL, isPartialSync bool) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	task, err := mgr.getTask(auditPlanId)
	if err != nil {
		return err
	}
	if isPartialSync {
		return task.PartialSyncSQLs(sqls)
	} else {
		return task.FullSyncSQLs(sqls)
	}
}

func (mgr *Manager) GetSQLs(auditPlanId uint, args map[string]interface{}) ([]Head, []map[string] /* head name */ string, uint64, error) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	args["audit_plan_id"] = auditPlanId

	task, err := mgr.getTask(auditPlanId)
	if err != nil {
		return nil, nil, 0, err
	}
	return task.GetSQLs(args)
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
