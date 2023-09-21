//go:build !enterprise
// +build !enterprise

package auditplan

import (
	e "errors"

	"github.com/actiontech/sqle/sqle/model"
)

func (sap *SyncFromAuditPlan) SyncSqlManager() error {
	return e.New("sql manage is enterprise version feature")
}

func (sa *SyncFromSqlAuditRecord) SyncSqlManager() error {
	return e.New("sql manage is enterprise version feature")
}

func SyncToSqlManage(sqls []*SQL, ap *model.AuditPlan) error {
	return e.New("sql manage is enterprise version feature")
}
