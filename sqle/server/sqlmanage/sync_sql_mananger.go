package sqlmanage

import (
	"github.com/actiontech/sqle/sqle/model"
)

type SyncSqlManager interface {
	SyncSqlManager() error
	UpdateSqlManageRecord() error
}

type SyncFromSqlAuditRecord struct {
	Task             *model.Task
	Source           string
	ProjectId        string
	SqlAuditRecordID string
}

func NewSyncFromSqlAudit(task *model.Task, source string, projectID string, sqlAuditID string) SyncSqlManager {
	return &SyncFromSqlAuditRecord{
		Task:             task,
		ProjectId:        projectID,
		SqlAuditRecordID: sqlAuditID,
		Source:           source,
	}
}
