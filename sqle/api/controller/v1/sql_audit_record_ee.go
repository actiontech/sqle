//go:build enterprise
// +build enterprise

package v1

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server/sqlmanage"
)

func syncSqlManage(record *model.SQLAuditRecord, tags []string) (err error) {
	// 以下两种情况不同步:
	// 1.更换标签
	// 2.在已存在标签的基础上增加新标签
	if record.Tags != nil && len(tags) > 0 {
		log.NewEntry().WithField("sync_sql_audit_record", record.AuditRecordId).Info("audit record tags is not empty,do not need to sync")
		return
	}

	if record.Task == nil || record.Task.ExecuteSQLs == nil {
		return fmt.Errorf("sql audit record task is nil")
	}

	newSyncFromSqlAudit := sqlmanage.NewSyncFromSqlAudit(record.Task, model.SQLManageSourceSqlAuditRecord, record.ProjectId, record.AuditRecordId)
	if len(tags) > 0 {
		if err := newSyncFromSqlAudit.SyncSqlManager(); err != nil {
			return fmt.Errorf("sync sql manager failed, error: %v", err)
		}
	} else if len(tags) == 0 {
		if err := newSyncFromSqlAudit.UpdateSqlManageRecord(); err != nil {
			return fmt.Errorf("sync sql manager failed, error: %v", err)
		}
	}

	return nil
}
