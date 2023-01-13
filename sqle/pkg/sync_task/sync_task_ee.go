//go:build enterprise
// +build enterprise

package sync_task

import (
	"context"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/sync_task/dmp"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var ExitCronChan chan string

func init() {
	ExitCronChan = make(chan string)
}

const SyncTaskActiontechDmp = "actiontech-dmp"

type SyncInstanceTask interface {
	GetSyncInstanceTaskFunc(context.Context) func()
}

func ReloadSyncInstanceTask(ctx context.Context, reloadReason string) {
	// 退出当前运行cron任务
	ExitCronChan <- reloadReason
	go EnableSyncInstanceTask(ctx)
}

func EnableSyncInstanceTask(ctx context.Context) {
	newLog := log.NewEntry()

	c := cron.New()

	s := model.GetStorage()
	syncTasks, err := s.GetAllSyncInstanceTasks()
	if err != nil {
		newLog.Errorf("get all sync tasks error: %v", err)
	}

	for _, syncTask := range syncTasks {
		var syncFunc func()

		syncInstance := NewSyncInstanceTask(newLog, syncTask.URL, syncTask.Version, syncTask.DbType, syncTask.RuleTemplate.Name)
		syncFunc = syncInstance.GetSyncInstanceTaskFunc(ctx)

		_, err := c.AddFunc(syncTask.SyncInstanceInterval, syncFunc)
		if err != nil {
			newLog.Errorf("add cron task error: %v", err)
		}
	}

	c.Start()

	exitReason := <-ExitCronChan

	c.Stop()

	newLog.Infof("exit cron task, reason: %s", exitReason)
}

func NewSyncInstanceTask(log *logrus.Entry, url, dmpVersion, dbType, ruleTemplateName string) SyncInstanceTask {
	switch dbType {
	case SyncTaskActiontechDmp:
		return dmp.NewDmpSync(log, url, dmpVersion, dbType, ruleTemplateName)
	}
	return nil
}
