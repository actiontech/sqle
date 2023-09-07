//go:build !enterprise
// +build !enterprise

package server

import (
	"time"

	"github.com/sirupsen/logrus"
)

type FeishuJob struct {
	BaseJob
}

func NewFeishuJob(entry *logrus.Entry) ServerJob {
	f := new(FeishuJob)
	f.BaseJob = *NewBaseJob(entry, 60*time.Second, f.feishuRotation)
	return f
}

func (j *FeishuJob) feishuRotation(entry *logrus.Entry) {
	// do nothing
}
