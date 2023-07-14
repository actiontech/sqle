//go:build !enterprise
// +build !enterprise

package server

import (
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"

	"github.com/sirupsen/logrus"
)

func CustomRuleAudit(l *logrus.Entry, DbType string, sqls []string, results []*driverV2.AuditResults) {}
