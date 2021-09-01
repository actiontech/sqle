package auditplan

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"actiontech.cloud/sqle/sqle/sqle/log"
	"actiontech.cloud/sqle/sqle/sqle/model"
	"actiontech.cloud/sqle/sqle/sqle/server"
	"actiontech.cloud/sqle/sqle/sqle/utils"

	"github.com/jinzhu/gorm"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var ErrAuditPlanNotExist = errors.New("audit plan not exist")
var ErrAuditPlanExisted = errors.New("audit plan existed")

var manager *Manager

func InitManager(s *model.Storage) chan struct{} {
	if manager == nil {
		manager = &Manager{
			scheduler: &scheduler{
				cron:     cron.New(),
				entryIDs: make(map[string]cron.EntryID),
			},
			persist: s,
			logger:  log.NewEntry(),
		}
		manager.start()
		exitCh := make(chan struct{})
		go func() {
			select {
			case <-exitCh:
				manager.stop()
			}
		}()
		return exitCh
	}
	return nil
}

func GetManager() *Manager {
	return manager
}

// Manager is the struct managing the persistent AuditPlans. It
// is *goroutine-safe*, since all exported methods are protected by a lock.
//
// All audit plan oprations except select should go through Manager.
type Manager struct {
	mu sync.Mutex

	scheduler *scheduler

	// persist is a database handle which store AuditPlan.
	persist *model.Storage

	logger *logrus.Entry
}

func (mgr *Manager) start() error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	mgr.scheduler.start()
	mgr.logger.Infoln("audit plan manager started")
	return mgr.loadAuditPlans()
}

func (mgr *Manager) stop() {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	ctx := mgr.scheduler.stop()
	<-ctx.Done()
	mgr.logger.Infoln("audit plan manager stopped")
}

func (mgr *Manager) AddStaticAuditPlan(name, cronExp, dbType, currentUserName string) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	ap := &model.AuditPlan{
		Name:           name,
		CronExpression: cronExp,
		DBType:         dbType,
	}
	return mgr.addAuditPlan(ap, currentUserName)
}

func (mgr *Manager) AddDynamicAuditPlan(name, cronExp, instanceName, instanceDatabase, currentUserName string) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	ap := &model.AuditPlan{
		Name:             name,
		CronExpression:   cronExp,
		InstanceName:     instanceName,
		InstanceDatabase: instanceDatabase,
	}
	return mgr.addAuditPlan(ap, currentUserName)
}

func (mgr *Manager) addAuditPlan(ap *model.AuditPlan, currentUserName string) error {
	user, exist, err := mgr.persist.GetUserByName(currentUserName)
	if !exist {
		return gorm.ErrRecordNotFound
	}
	if err != nil {
		return err
	}
	ap.CreateUserID = user.ID

	j := utils.NewJWT([]byte(utils.JWTSecret))
	t, err := j.CreateToken(currentUserName, time.Now().Add(time.Hour*24*365).Unix())
	if err != nil {
		return err
	}
	ap.Token = t

	if ap.InstanceName != "" {
		instance, exist, err := mgr.persist.GetInstanceByName(ap.InstanceName)
		if !exist {
			return gorm.ErrRecordNotFound
		}
		if err != nil {
			return err
		}
		ap.DBType = instance.DbType
	}

	err = mgr.persist.Save(ap)
	if err != nil {
		return err
	}
	return mgr.addAuditPlansToScheduler([]*model.AuditPlan{ap})
}

func (mgr *Manager) UpdateAuditPlan(name string, attrs map[string]interface{}) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if !mgr.scheduler.hasJob(name) {
		return ErrAuditPlanNotExist
	}

	err := mgr.persist.UpdateAuditPlanByName(name, attrs)
	if err != nil {
		return err
	}
	err = mgr.scheduler.removeJob(name)
	if err != nil {
		return err
	}
	ap, _, err := mgr.persist.GetAuditPlanByName(name)
	if err != nil {
		return err
	}
	return mgr.scheduler.addJob(ap, func() {
		mgr.runJob(ap)
	})
}

func (mgr *Manager) DeleteAuditPlan(name string) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if !mgr.scheduler.hasJob(name) {
		return ErrAuditPlanNotExist
	}

	ap, _, err := mgr.persist.GetAuditPlanByName(name)
	if err != nil {
		return err
	}
	err = mgr.persist.Delete(ap)
	if err != nil {
		return err
	}
	return mgr.scheduler.removeJob(name)
}

func (mgr *Manager) TriggerAuditPlan(name string) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	ap, _, err := mgr.persist.GetAuditPlanByName(name)
	if err != nil {
		return err
	}
	mgr.runJob(ap)
	return nil
}

func (mgr *Manager) loadAuditPlans() error {
	aps, err := mgr.persist.GetAuditPlans()
	if err != nil {
		return err
	}
	return mgr.addAuditPlansToScheduler(aps)
}

func (mgr *Manager) runJob(ap *model.AuditPlan) {
	instance, _, err := mgr.persist.GetInstanceByName(ap.InstanceName)
	if err != nil {
		mgr.logger.WithField("name", ap.Name).Errorf("get instance error:%v\n", err)
		return
	}

	task := &model.Task{
		Schema:       ap.InstanceDatabase,
		InstanceId:   instance.ID,
		CreateUserId: ap.CreateUserID,
		SQLSource:    model.TaskSQLSourceFromAuditPlan,
		DBType:       ap.DBType,
	}
	{ // todo: extract common logic in CreateAndAuditTask
		sqls, err := mgr.persist.GetAuditPlanSQLs(ap.Name)
		if err != nil {
			mgr.logger.WithField("name", ap.Name).Errorf("get audit plan SQLs error:%v\n", err)
			return
		}
		for i, sql := range sqls {
			task.ExecuteSQLs = append(task.ExecuteSQLs, &model.ExecuteSQL{
				BaseSQL: model.BaseSQL{
					Number:  uint(i),
					Content: sql.LastSQLText,
				},
			})
		}
		err = mgr.persist.Save(task)
		if err != nil {
			mgr.logger.WithField("name", ap.Name).Errorf("save audit plan task error:%v\n", err)
			return
		}

		task, err = server.GetSqled().AddTaskWaitResult(fmt.Sprintf("%v", task.ID), server.ActionTypeAudit)
		if err != nil {
			mgr.logger.WithField("name", ap.Name).Errorf("audit task error:%v\n", err)
			return
		}
	}

	auditPlanReport := &model.AuditPlanReport{AuditPlanID: fmt.Sprintf("%v", ap.ID)}
	for _, executeSQL := range task.ExecuteSQLs {
		auditPlanReport.AuditPlanReportSQLs = append(auditPlanReport.AuditPlanReportSQLs, &model.AuditPlanReportSQL{
			AuditResult: executeSQL.AuditResult,
		})
	}

	err = mgr.persist.Save(auditPlanReport)
	if err != nil {
		mgr.logger.WithField("name", ap.Name).Errorf("save audit plan report error:%v\n", err)
		return
	}
}

func (mgr *Manager) addAuditPlansToScheduler(aps []*model.AuditPlan) error {
	for _, ap := range aps {
		mgr.scheduler.addJob(ap, func() {
			mgr.runJob(ap)
		})
		mgr.logger.WithFields(logrus.Fields{
			"name":            ap.Name,
			"cron_expression": ap.CronExpression}).Infoln("audit plan added")
	}
	return nil
}

// scheduler is not goroutine safe.
type scheduler struct {
	// cron is a AuditPlan scheduler.
	cron *cron.Cron

	// entryIDs maps audit plan name to it's job entry ID.
	entryIDs map[string]cron.EntryID
}

func (s *scheduler) removeJob(auditPlanName string) error {
	entryID, ok := s.entryIDs[auditPlanName]
	if !ok {
		return ErrAuditPlanNotExist
	}

	s.cron.Remove(entryID)
	delete(s.entryIDs, auditPlanName)
	return nil
}

func (s *scheduler) addJob(ap *model.AuditPlan, do func()) error {
	_, ok := s.entryIDs[ap.Name]
	if !ok {
		return ErrAuditPlanExisted
	}

	entryID, err := s.cron.AddFunc(ap.CronExpression, do)
	if err != nil {
		return err
	}
	s.entryIDs[ap.Name] = entryID
	return nil
}

func (s *scheduler) start() {
	s.cron.Start()
}

func (s *scheduler) stop() context.Context {
	return s.cron.Stop()
}

func (s *scheduler) hasJob(auditPlanName string) bool {
	_, has := s.entryIDs[auditPlanName]
	return has
}
