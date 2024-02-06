//go:build !enterprise
// +build !enterprise

package server

import (
	"time"

	"github.com/sirupsen/logrus"
)

type DingTalkJob struct {
	BaseJob
}

func NewDingTalkJob(entry *logrus.Entry) ServerJob {
	d := new(DingTalkJob)
	d.BaseJob = *NewBaseJob(entry, 60*time.Second, d.dingTalkRotation)
	return d
}

func (j *DingTalkJob) dingTalkRotation(entry *logrus.Entry) {
	// do nothing
}
