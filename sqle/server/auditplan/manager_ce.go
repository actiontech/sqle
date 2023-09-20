//go:build !enterprise
// +build !enterprise

package auditplan

import (
	e "errors"
)

func (sap *SyncFromAuditPlan) SyncSqlManager() error {
	return e.New("sql manage is enterprise version feature")
}

func (sa *SyncFromSqlAuditRecord) SyncSqlManager() error {
	return e.New("sql manage is enterprise version feature")
}
