//go:build enterprise
// +build enterprise

package sync

import (
	"context"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/sync/dmp"
	"github.com/robfig/cron/v3"
)

var ExitCronChan chan string

func init() {
	ExitCronChan = make(chan string)
}

const ActiontechDmp = "actiontech-dmp"

type SyncInstance interface {
	Sync(context.Context) func()
}

func ReloadInstance(ctx context.Context, reloadReason string) {
	// 退出当前运行cron任务
	ExitCronChan <- reloadReason
	go EnableInstanceSync(ctx)
}

func EnableInstanceSync(ctx context.Context) {
	newLog := log.NewEntry()

	c := cron.New()

	s := model.GetStorage()
	syncTasks, err := s.GetAllSyncTasks()
	if err != nil {
		newLog.Errorf("get all sync tasks error: %v", err)
	}

	for _, syncTask := range syncTasks {
		var syncFunc func()

		switch syncTask.Source {
		case ActiontechDmp:
			dmpSync := dmp.NewDmpSync(newLog, syncTask.URL, syncTask.Version, syncTask.DbType, syncTask.RuleTemplate.Name)
			syncFunc = dmpSync.Sync(ctx)
		default:
			continue
		}

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
