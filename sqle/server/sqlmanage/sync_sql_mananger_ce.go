//go:build !enterprise
// +build !enterprise

package sqlmanage

import (
	e "errors"
)

func (sa *SyncFromSqlAuditRecord) SyncSqlManager() error {
	return e.New("sql manage is enterprise version feature")
}

func (sa *SyncFromSqlAuditRecord) UpdateSqlManageRecord() error {
	return e.New("sql manage is enterprise version feature")
}
