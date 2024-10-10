package notification

import (
	"context"
	"fmt"
	"strings"

	"github.com/actiontech/dms/pkg/dms-common/i18nPkg"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/locale"
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

func (a *SQLmanageRecordNotification) NotificationSubject() i18nPkg.I18nStr {
	return locale.Bundle.LocalizeAll(locale.NotifyManageRecordSubject)
}

func (a *SQLmanageRecordNotification) NotificationBody() i18nPkg.I18nStr {
	timeInBody := locale.Bundle.LocalizeAllWithArgs(locale.NotifyManageRecordBodyTime, a.Config.StartTime, a.Config.EndTime)
	projInBody := locale.Bundle.LocalizeAllWithArgs(locale.NotifyManageRecordBodyProj, a.Config.ProjectName)

	var linkInBody i18nPkg.I18nStr
	if a.Config.SQLEUrl != "" && a.Config.ProjectName != "" {
		// TODO 当前只跳转到管控页面，无筛选条件
		// builder.WriteString(fmt.Sprintf("\n- SQL管控记录链接: %v/v2/projects/%v/sql_manages?filter_priority=hight&filter_last_audit_start_time_from=%v&filter_last_audit_start_time_to=%v&page_index=1&page_size=20&filter_status=unhandled", strings.TrimRight(a.Config.SQLEUrl, "/"), a.Config.ProjectName, a.Config.StartTime, a.Config.EndTime))
		projectId, err := dms.GetPorjectUIDByName(context.TODO(), a.Config.ProjectName)
		if err == nil {
			link := fmt.Sprintf("%v/project/%v/sql-management", strings.TrimRight(a.Config.SQLEUrl, "/"), projectId)
			linkInBody = locale.Bundle.LocalizeAllWithArgs(locale.NotifyManageRecordBodyLink, link)
		}

	}

	bodyStr := make([]i18nPkg.I18nStr, 0, len(a.SQLmanageRecords)+3)
	bodyStr = append(bodyStr, timeInBody, projInBody, linkInBody)

	for _, sqlManagerRecord := range a.SQLmanageRecords {
		auditResults := make(i18nPkg.I18nStr, len(locale.Bundle.LanguageTags()))
		for _, langTag := range locale.Bundle.LanguageTags() {
			auditResults.SetStrInLang(langTag, sqlManagerRecord.AuditResults.GetAuditJsonStrByLangTag(langTag))
		}

		db := dms.GetInstancesByIdWithoutError(sqlManagerRecord.InstanceID)
		recordI18n := locale.Bundle.LocalizeAllWithArgs(locale.NotifyManageRecordBodyRecord,
			sqlManagerRecord.SQLID,
			fmt.Sprintf("%v(%v:%v)", db.Name, db.Host, db.Port),
			db.Business,
			sqlManagerRecord.SqlText,
			sqlManagerRecord.AuditLevel,
			auditResults,
		)
		bodyStr = append(bodyStr, recordI18n)
	}

	return locale.Bundle.JoinI18nStr(bodyStr, "\n")
}
