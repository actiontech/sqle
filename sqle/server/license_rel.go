//go:build release
// +build release

package server

import (
	"time"

	"github.com/actiontech/sqle/sqle/model"
	"github.com/sirupsen/logrus"
)

func init() {
	OnlyRunOnLeaderJobs = append(OnlyRunOnLeaderJobs, NewLicenseJob)
}

type LicenseJob struct {
	BaseJob
}

func NewLicenseJob(entry *logrus.Entry) ServerJob {
	entry = entry.WithField("job", "license")
	j := &LicenseJob{}
	j.BaseJob = *NewBaseJob(entry, 1*time.Hour, j.UpdateLicense)
	return j
}

func (j *LicenseJob) UpdateLicense(entry *logrus.Entry) {
	s := model.GetStorage()
	l, exist, err := s.GetLicense()
	if err != nil {
		entry.Errorf("fail to get license, error: %v", err)
		return
	}
	if !exist || l.Content == nil {
		return
	}
	l.Content.WorkDurationHour += 1
	l.WorkDurationHour = l.WorkDurationHour + 1
	err = s.Save(l)
	if err != nil {
		entry.Errorf("fail to update license, error: %v", err)
	}
}
