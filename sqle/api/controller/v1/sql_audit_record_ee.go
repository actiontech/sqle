//go:build enterprise
// +build enterprise

package v1

import (
	"context"
	"fmt"

	"github.com/actiontech/sqle/sqle/common"
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

	plugin, err := common.NewDriverManagerWithoutCfg(log.NewEntry(), record.Task.DBType)
	if err != nil {
		return fmt.Errorf("open plugin failed: %v", err)
	}
	defer plugin.Close(context.TODO())

	sqlFpMap := make(map[string]string)
	for _, sql := range record.Task.ExecuteSQLs {
		node, err := plugin.Parse(context.TODO(), sql.Content)
		if err != nil {
			return fmt.Errorf("parse sqls failed: %v", err)
		}
		sqlFpMap[sql.Content] = node[0].Fingerprint
	}
	newSyncFromSqlAudit := sqlmanage.NewSyncFromSqlAudit(record.Task, sqlFpMap, record.ProjectId, record.AuditRecordId)
	if len(tags) > 0 {
		if err := newSyncFromSqlAudit.SyncSqlManager(model.SQLManageSourceSqlAuditRecord); err != nil {
			return fmt.Errorf("sync sql manager failed, error: %v", err)
		}
	} else if len(tags) == 0 {

		if err := newSyncFromSqlAudit.UpdateSqlManageRecord(record.AuditRecordId, model.SQLManageSourceSqlAuditRecord); err != nil {
			return fmt.Errorf("sync sql manager failed, error: %v", err)
		}

		s := model.GetStorage()
		err := s.UpdateSqlManage(record.ID)
		if err != nil {
			return fmt.Errorf("update sql manage failed: %v", err)
		}
	}

	return nil
}
