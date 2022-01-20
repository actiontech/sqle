package auditplan

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/jinzhu/gorm"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var tokenExpire = 365 * 24 * time.Hour

var ErrAuditPlanNotExist = errors.New("audit plan not exist")

var ErrAuditPlanExisted = errors.New("audit plan existed")

var manager *Manager

func InitManager(s *model.Storage) chan struct{} {
	manager = &Manager{
		scheduler: &scheduler{
			cron:     cron.New(),
			entryIDs: make(map[string]cron.EntryID),
		},
		persist: s,
		logger:  log.NewEntry(),
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

func (mgr *Manager) AddStaticAuditPlan(name, cronExp, dbType, currentUserName,
	auditPlanType string, ps params.Params) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	ap := &model.AuditPlan{
		Name:           name,
		CronExpression: cronExp,
		DBType:         dbType,
		Type:           auditPlanType,
		Params:         ps,
	}

	return mgr.addAuditPlan(ap, currentUserName)
}

func (mgr *Manager) AddDynamicAuditPlan(name, cronExp, instanceName, instanceDatabase, currentUserName,
	auditPlanType string, ps params.Params) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	ap := &model.AuditPlan{
		Name:             name,
		CronExpression:   cronExp,
		InstanceName:     instanceName,
		InstanceDatabase: instanceDatabase,
		Type:             auditPlanType,
		Params:           ps,
	}

	return mgr.addAuditPlan(ap, currentUserName)
}

func (mgr *Manager) addAuditPlan(ap *model.AuditPlan, currentUserName string) error {
	if mgr.scheduler.hasJob(ap.Name) {
		return ErrAuditPlanExisted
	}

	user, exist, err := mgr.persist.GetUserByName(currentUserName)
	if !exist {
		return gorm.ErrRecordNotFound
	} else if err != nil {
		return err
	}

	j := utils.NewJWT([]byte(utils.JWTSecret))

	t, err := j.CreateToken(currentUserName, time.Now().Add(tokenExpire).Unix(),
		utils.WithAuditPlanName(ap.Name))
	if err != nil {
		return err
	}

	if ap.InstanceName != "" {
		instance, exist, err := mgr.persist.GetInstanceByName(ap.InstanceName)
		if !exist {
			return gorm.ErrRecordNotFound
		} else if err != nil {
			return err
		}

		ap.DBType = instance.DbType
	}

	ap.Token = t
	ap.CreateUserID = user.ID

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

var errNoSQLInAuditPlan = errors.New("there is no SQLs in audit plan")

func (mgr *Manager) TriggerAuditPlan(name string) (*model.AuditPlanReport, error) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	ap, _, err := mgr.persist.GetAuditPlanByName(name)
	if err != nil {
		return nil, err
	}

	report := mgr.runJob(ap)
	if report == nil {
		return nil, errNoSQLInAuditPlan
	}

	return report, nil
}

func (mgr *Manager) loadAuditPlans() error {
	aps, err := mgr.persist.GetAuditPlans()
	if err != nil {
		return err
	}

	return mgr.addAuditPlansToScheduler(aps)
}

// TODO: runJob is a async task, it's report should send by channel.
func (mgr *Manager) runJob(ap *model.AuditPlan) *model.AuditPlanReport {
	task := &model.Task{
		Schema:       ap.InstanceDatabase,
		CreateUserId: ap.CreateUserID,
		SQLSource:    model.TaskSQLSourceFromAuditPlan,
		DBType:       ap.DBType,
	}

	// todo: extract common logic in CreateAndAuditTask
	auditPlanSQLs, err := mgr.persist.GetAuditPlanSQLs(ap.Name)
	if err != nil {
		mgr.logger.WithField("name", ap.Name).Errorf("get audit plan SQLs error:%v\n", err)
		return nil
	}

	if len(auditPlanSQLs) == 0 {
		mgr.logger.WithField("name", ap.Name).Warnf("skip audit, %v", errNoSQLInAuditPlan)
		return nil
	}

	for i, sql := range auditPlanSQLs {
		task.ExecuteSQLs = append(task.ExecuteSQLs, &model.ExecuteSQL{
			BaseSQL: model.BaseSQL{
				Number:  uint(i),
				Content: sql.LastSQL,
			},
		})
	}

	instance, _, err := mgr.persist.GetInstanceByName(ap.InstanceName)
	if err != nil {
		mgr.logger.WithField("name", ap.Name).Errorf("get instance error:%v\n", err)
		return nil
	}

	task.InstanceId = instance.ID

	err = mgr.persist.Save(task)
	if err != nil {
		mgr.logger.WithField("name", ap.Name).Errorf("save audit plan task error:%v\n", err)
		return nil
	}

	task, err = server.GetSqled().AddTaskWaitResult(fmt.Sprintf("%v", task.ID), server.ActionTypeAudit)
	if err != nil {
		mgr.logger.WithField("name", ap.Name).Errorf("audit task error:%v\n", err)
		return nil
	}

	auditPlanReport := &model.AuditPlanReport{AuditPlanID: ap.ID}
	for i, executeSQL := range task.ExecuteSQLs {
		auditPlanReport.AuditPlanReportSQLs = append(auditPlanReport.AuditPlanReportSQLs, &model.AuditPlanReportSQL{
			AuditPlanSQLID: auditPlanSQLs[i].ID,
			AuditResult:    executeSQL.AuditResult,
		})
	}

	err = mgr.persist.Save(auditPlanReport)
	if err != nil {
		mgr.logger.WithField("name", ap.Name).Errorf("save audit plan report error:%v\n", err)
		return nil
	}

	return auditPlanReport
}

func (mgr *Manager) addAuditPlansToScheduler(aps []*model.AuditPlan) error {
	for _, v := range aps {
		ap := v

		err := mgr.scheduler.addJob(ap, func() {
			mgr.runJob(ap)
		})
		if err != nil {
			return err
		}

		mgr.logger.WithFields(logrus.Fields{
			"name":            ap.Name,
			"cron_expression": ap.CronExpression,
		}).Infoln("audit plan added")
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
	if ok {
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
