package sqlQuery

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/sirupsen/logrus"
)

var ErrSqlQueryAuditLevelIsNotAllowed = errors.New(errors.DataExist, fmt.Errorf("the audit level is not allowed to perform sql query"))

func Audit(sqls []string, instance *model.Instance) error {
	task := &model.Task{
		DBType:   instance.DbType,
		Instance: instance,
	}
	for i, sql := range sqls {
		task.ExecuteSQLs = append(task.ExecuteSQLs, &model.ExecuteSQL{
			BaseSQL: model.BaseSQL{
				Number:  uint(i),
				Content: sql,
			},
		})
	}

	logger := log.NewEntry().WithField("type", "sql_query")
	err := server.Audit(logger, task)
	if err != nil {
		return err
	}

	allowQueryWhenLessThanAuditLevel := driver.RuleLevel(instance.SqlQueryConfig.AllowQueryWhenLessThanAuditLevel)
	if allowQueryWhenLessThanAuditLevel.LessOrEqual(driver.RuleLevel(task.AuditLevel)) {
		auditResults := make(map[string]struct{})
		for _, sql := range task.ExecuteSQLs {
			if allowQueryWhenLessThanAuditLevel.LessOrEqual(driver.RuleLevel(sql.AuditLevel)) {
				auditResults[sql.AuditResult] = struct{}{}
				logger.WithFields(logrus.Fields{
					"sql":                                    sql.Content,
					"audit_result":                           sql.AuditResult,
					"audit_level":                            sql.AuditLevel,
					"allow_query_when_less_than_audit_level": instance.SqlQueryConfig.AllowQueryWhenLessThanAuditLevel,
				}).Errorln(ErrSqlQueryAuditLevelIsNotAllowed.Error())
			}
		}

		var auditResultsMsg string
		for res, _ := range auditResults {
			auditResultsMsg = fmt.Sprintf("%v\n%v", auditResultsMsg, res)
		}

		return fmt.Errorf("%v: %v", ErrSqlQueryAuditLevelIsNotAllowed, auditResultsMsg)
	}

	return nil
}
