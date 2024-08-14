//go:build enterprise
// +build enterprise

package server

import (
	"fmt"
	"time"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/notification"
	"github.com/robfig/cron/v3"
)

const (
	SQLManagePushMessageType = "sql_manage"
	SQLReportPushMessageType = "sql_report"
)

func newPushJob(p *model.ReportPushConfig) (PushJob, error) {
	switch p.Type {
	case SQLManagePushMessageType:
		pushJobSQLMessage := SQLManageRecordPushJob{
			Config: p,
		}
		return pushJobSQLMessage, nil
	default:
		return nil, fmt.Errorf("type %v not supported", p.Type)
	}
}

// SQL 管控推送任务
type SQLManageRecordPushJob struct {
	Config *model.ReportPushConfig
}

func (p SQLManageRecordPushJob) Cron() string {
	return p.Config.PushFrequencyCron
}
func (p SQLManageRecordPushJob) Run() {
	logger := log.NewEntry()
	s := model.GetStorage()

	sqls, err := s.GetLastHightLevelSQLs(p.Config.ProjectId, p.Config.LastPushTime)
	if err != nil {
		logger.Error(err)
		return
	}
	if len(sqls) == 0 {
		return
	}
	url, err := s.GetSqleUrl()
	if err != nil {
		logger.Error(err)
		return
	}

	project, err := dms.GetProjectByID(p.Config.ProjectId)
	if err != nil {
		logger.Error(err)
		return
	}
	notify := notification.NewSQLmanageRecordNotification(notification.SQLmanageRecordNotifyConfig{
		SQLEUrl:     url,
		ProjectName: project.Name,
		StartTime:   p.Config.LastPushTime.Format(time.RFC3339),
		EndTime:     time.Now().Format(time.RFC3339),
	}, sqls)
	err = notification.Notify(notify, p.Config.PushUserList)
	if err != nil {
		logger.Error(err)
		return
	}
	p.Config.LastPushTime = time.Now()
	err = s.Save(p.Config)
	if err != nil {
		logger.Error(err)
		return
	}
}

// SQL 管控报告
type SQLReportPushJob struct {
	scheduler *cron.Cron
	config    model.ReportPushConfig
}

func (p SQLReportPushJob) Cron() string {
	return p.config.PushFrequencyCron
}

func (p SQLReportPushJob) Run() {
}
