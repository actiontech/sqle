//go:build enterprise
// +build enterprise

package sync_task

// import (
// 	"context"

// 	"github.com/actiontech/sqle/sqle/model"
// 	"github.com/sirupsen/logrus"
// )

// type SyncInstanceTask interface {
// 	GetSyncInstanceTaskFunc(context.Context) func()
// }

// func NewSyncInstanceTask(log *logrus.Entry, id uint, source, url, dmpVersion, dbType, ruleTemplateName string) SyncInstanceTask {
// 	switch source {
// 	case model.SyncTaskSourceActiontechDmp:
// 		return NewDmpSync(log, id, url, dmpVersion, dbType, ruleTemplateName)
// 	}
// 	return nil
// }
