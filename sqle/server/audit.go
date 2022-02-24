package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func Audit(l *logrus.Entry, task *model.Task) (err error) {
	d, err := newDriverWithAudit(l, task.Instance, task.Schema, task.DBType)
	if err != nil {
		return err
	}
	defer d.Close(context.TODO())

	return audit(l, task, d)
}

func audit(l *logrus.Entry, task *model.Task, d driver.Driver) (err error) {
	st := model.GetStorage()

	whitelist, _, err := st.GetSqlWhitelist(0, 0)
	if err != nil {
		return err
	}
	for _, executeSQL := range task.ExecuteSQLs {
		// We always trust the ExecuteSQL.Content is single SQL.
		//
		// The audit() function has two producers for now:
		// 1. from API controller
		//		- the API controller should call Parse before audit.
		//      - If Parse() can not splits SQL to expected case, user can add SQL to whitelist for workaround.
		// 2. from audit plan
		//		- the audit plan may collect SQLs which plugins can not Parse.
		//      - In these case, we pass the raw SQL to plugins, it's ok.
		node, err := parse(l, d, executeSQL.Content)
		if err != nil {
			return err
		}
		var whitelistMatch bool
		for _, wl := range whitelist {
			if wl.MatchType == model.SQLWhitelistFPMatch {
				wlNode, err := parse(l, d, wl.Value)
				if err != nil {
					return err
				}
				if node.Fingerprint == wlNode.Fingerprint {
					whitelistMatch = true
				}
			} else {
				if wl.CapitalizedValue == strings.ToUpper(node.Text) {
					whitelistMatch = true
				}
			}
		}
		result := driver.NewInspectResults()
		if whitelistMatch {
			result.Add(driver.RuleLevelNormal, "白名单")
		} else {
			result, err = d.Audit(context.TODO(), executeSQL.Content)
			if err != nil {
				return err
			}
		}
		executeSQL.AuditStatus = model.SQLAuditStatusFinished
		executeSQL.AuditLevel = string(result.Level())
		executeSQL.AuditResult = result.Message()
		executeSQL.AuditFingerprint = utils.Md5String(string(append([]byte(result.Message()), []byte(node.Fingerprint)...)))

		l.WithFields(logrus.Fields{
			"SQL":    executeSQL.Content,
			"level":  executeSQL.AuditLevel,
			"result": executeSQL.AuditResult}).Info("audit finished")
	}
	return nil
}

func parse(l *logrus.Entry, d driver.Driver, sql string) (node driver.Node, err error) {
	nodes, err := d.Parse(context.TODO(), sql)
	if err != nil {
		return node, errors.Wrapf(err, "parse sql: %s", sql)
	}
	if len(nodes) == 0 {
		return node, fmt.Errorf("the node is empty after parse")
	}
	if len(nodes) > 1 {
		l.Errorf("the SQL is not single SQL: %s", sql)
	}
	return nodes[0], nil
}

func genRollbackSQL(l *logrus.Entry, task *model.Task, d driver.Driver) ([]*model.RollbackSQL, error) {
	rollbackSQLs := make([]*model.RollbackSQL, 0, len(task.ExecuteSQLs))
	for _, executeSQL := range task.ExecuteSQLs {
		rollbackSQL, reason, err := d.GenRollbackSQL(context.TODO(), executeSQL.Content)
		if err != nil {
			l.Errorf("gen rollback sql error, %v", err)
			return nil, err
		}
		result := driver.NewInspectResults()
		result.Add(driver.RuleLevel(executeSQL.AuditLevel), executeSQL.AuditResult)
		result.Add(driver.RuleLevelNotice, reason)
		executeSQL.AuditLevel = string(result.Level())
		executeSQL.AuditResult = result.Message()

		rollbackSQLs = append(rollbackSQLs, &model.RollbackSQL{
			BaseSQL: model.BaseSQL{
				TaskId:  executeSQL.TaskId,
				Content: rollbackSQL,
			},
			ExecuteSQLId: executeSQL.ID,
		})
	}
	return rollbackSQLs, nil
}
