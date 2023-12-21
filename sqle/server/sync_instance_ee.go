//go:build enterprise
// +build enterprise

package server

// import (
// 	"context"
// 	"sync"
// 	"time"

// 	"github.com/actiontech/sqle/sqle/model"
// 	syncTask "github.com/actiontech/sqle/sqle/pkg/sync_task"
// 	"github.com/robfig/cron/v3"

// 	"github.com/sirupsen/logrus"
// )

// func init() {
// 	OnlyRunOnLeaderJobs = append(OnlyRunOnLeaderJobs, NewSyncInstanceJob)
// }

// type SyncInstanceJob struct {
// 	BaseJob
// 	cron           *cron.Cron
// 	once           sync.Once
// 	lastReloadTime time.Time
// }

// func NewSyncInstanceJob(entry *logrus.Entry) ServerJob {
// 	entry = entry.WithField("job", "sync_instance")
// 	j := &SyncInstanceJob{
// 		once:           sync.Once{},
// 		lastReloadTime: time.Now(),
// 	}
// 	j.BaseJob = *NewBaseJob(entry, 10*time.Second, j.ReloadSyncInstanceTask)
// 	return j
// }

// func (j *SyncInstanceJob) Stop() {
// 	j.BaseJob.Stop()
// 	j.stopSyncInstanceTask()
// }

// func (j *SyncInstanceJob) ReloadSyncInstanceTask(entry *logrus.Entry) {
// 	j.once.Do(func() {
// 		entry.Infof("start load task")
// 		j.startSyncInstanceTask(entry)
// 		entry.Infof("end load task")
// 		return
// 	})
// 	s := model.GetStorage()
// 	task, exist, err := s.GetLatestSyncInstanceTask(j.lastReloadTime)
// 	if err != nil {
// 		return
// 	}
// 	if exist {
// 		entry.Infof("start reload task")
// 		j.stopSyncInstanceTask()
// 		j.startSyncInstanceTask(entry)
// 		j.lastReloadTime = task.UpdatedAt
// 		entry.Infof("end reload task")
// 	}
// }

// func (j *SyncInstanceJob) startSyncInstanceTask(entry *logrus.Entry) {
// 	j.cron = cron.New()

// 	s := model.GetStorage()

// 	tasks, err := s.GetAllSyncInstanceTasks()
// 	if err != nil {
// 		entry.Errorf("get all sync tasks error: %v", err)
// 	}
// 	for _, task := range tasks {
// 		var syncFunc func()

// 		syncInstance := syncTask.NewSyncInstanceTask(entry, task.ID, task.Source, task.URL, task.Version, task.DbType, task.RuleTemplate.Name)
// 		ctx := context.TODO()
// 		syncFunc = syncInstance.GetSyncInstanceTaskFunc(ctx)

// 		_, err := j.cron.AddFunc(task.SyncInstanceInterval, syncFunc)
// 		if err != nil {
// 			entry.Errorf("add cron task error: %v", err)
// 		}
// 	}

// 	j.cron.Start()
// }

// func (j *SyncInstanceJob) stopSyncInstanceTask() {
// 	if j.cron != nil {
// 		j.cron.Stop()
// 		j.cron = nil
// 	}
// }
