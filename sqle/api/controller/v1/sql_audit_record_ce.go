//go:build !enterprise
// +build !enterprise

package v1

import (
	"github.com/actiontech/sqle/sqle/model"
)

func syncSqlManage(record *model.SQLAuditRecord, tags []string) (err error) {
	return nil
}
