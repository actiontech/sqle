//go:build !enterprise
// +build !enterprise

package server

import (
	"time"

	"github.com/sirupsen/logrus"
)

type WechatJob struct {
	BaseJob
}

func NewWechatJob(entry *logrus.Entry) ServerJob {
	w := new(WechatJob)
	w.BaseJob = *NewBaseJob(entry, 60*time.Second, w.wechatRotation)
	return w
}

// 当前企微轮询只为二次确认工单定时上线功能
// https://github.com/actiontech/sqle-ee/issues/1441
func (w *WechatJob) wechatRotation(entry *logrus.Entry) {
	// nothing
}
