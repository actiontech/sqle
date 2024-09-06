package notification

import (
	"context"
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/model"
)

type SQLmanageRecordNotification struct {
	SQLmanageRecords []*model.SQLManageRecord
	Config           SQLmanageRecordNotifyConfig
}

type SQLmanageRecordNotifyConfig struct {
	SQLEUrl     string
	ProjectName string
	StartTime   string
	EndTime     string
}

func NewSQLmanageRecordNotification(config SQLmanageRecordNotifyConfig, SQLmanageRecord []*model.SQLManageRecord) *SQLmanageRecordNotification {
	return &SQLmanageRecordNotification{
		SQLmanageRecords: SQLmanageRecord,
		Config:           config,
	}
}

func (a *SQLmanageRecordNotification) NotificationSubject() string {
	return "SQL管控记录"
}

func (a *SQLmanageRecordNotification) NotificationBody() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("记录时间周期：%v - %v\n", a.Config.StartTime, a.Config.EndTime))
	builder.WriteString(fmt.Sprintf(`所属项目  %v`, a.Config.ProjectName))

	if a.Config.SQLEUrl != "" && a.Config.ProjectName != "" {
		// TODO 当前只跳转到管控页面，无筛选条件
		// builder.WriteString(fmt.Sprintf("\n- SQL管控记录链接: %v/v2/projects/%v/sql_manages?filter_priority=hight&filter_last_audit_start_time_from=%v&filter_last_audit_start_time_to=%v&page_index=1&page_size=20&filter_status=unhandled", strings.TrimRight(a.Config.SQLEUrl, "/"), a.Config.ProjectName, a.Config.StartTime, a.Config.EndTime))
		projectId, err := dms.GetPorjectUIDByName(context.TODO(), a.Config.ProjectName)
		if err == nil {
			builder.WriteString(fmt.Sprintf("\n- SQL管控记录链接: %v/project/%v/sql-management", strings.TrimRight(a.Config.SQLEUrl, "/"), projectId))
		}

	}
	for _, sqlManagerRecord := range a.SQLmanageRecords {
		db := dms.GetInstancesByIdWithoutError(sqlManagerRecord.InstanceID)
		builder.WriteString(fmt.Sprintf(`
- SQLID	%v
- 所在数据源名称 %v
- 所属业务 %v
- SQL %v
- 触发规则级别 %v
- SQL审核建议 %v
================================`,
			sqlManagerRecord.SQLID,
			fmt.Sprintf("%v(%v:%v)", db.Name, db.Host, db.Port),
			db.Business,
			sqlManagerRecord.SqlText,
			sqlManagerRecord.AuditLevel,
			sqlManagerRecord.AuditResults,
		))
	}

	return builder.String()
}
