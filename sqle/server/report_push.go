package server

import (
	"time"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type ReportPushJob struct {
	BaseJob
	LastSyncTime time.Time
	Schedule     *cron.Cron
	ScheduleJob  map[string] /*reportPushConfigId*/ cron.EntryID /*entiryID*/
}

func NewReportPushJob(entry *logrus.Entry) ServerJob {
	entry = entry.WithField("job", "report_push")
	j := &ReportPushJob{Schedule: cron.New(), ScheduleJob: make(map[string]cron.EntryID)}
	j.Schedule.Start()
	j.BaseJob = *NewBaseJob(entry, 5*time.Second, j.SyncReportPushJob)
	return j
}

func (j *ReportPushJob) SyncReportPushJob(entry *logrus.Entry) {
	logger := log.NewEntry()
	s := model.GetStorage()

	configs, err := s.GetLastUpdateReportPushConfig(j.LastSyncTime)
	if err != nil {
		logger.Errorf("get last update report push config failed: %v", err)
		return
	}
	j.LastSyncTime = time.Now()

	for _, config := range configs {
		// 关闭/移除运行中任务
		if entryID, ok := j.ScheduleJob[config.GetIDStr()]; ok {
			j.Schedule.Remove(entryID)
			delete(j.ScheduleJob, config.GetIDStr())
		}

		// 重新启动任务
		if config.Enabled {
			pushTask, err := newPushJob(config)
			if err != nil {
				logger.Errorf("new push job failed: %v", err)
				continue
			}
			entryID, err := j.Schedule.AddJob(pushTask.Cron(), pushTask)
			if err != nil {
				logger.Errorf("add job failed: %v", err)
				continue
			}
			j.ScheduleJob[config.GetIDStr()] = entryID
		}
	}

}

// 执行推送任务
type PushJob interface {
	Cron() string
	Run()
}
