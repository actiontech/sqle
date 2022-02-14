package server

import (
	"context"
	"strings"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"
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
		nodes, err := d.Parse(context.TODO(), executeSQL.Content)
		if err != nil {
			return err
		}
		if len(nodes) != 1 {
			return driver.ErrNodesCountExceedOne
		}
		var whitelistMatch bool
		for _, wl := range whitelist {
			if wl.MatchType == model.SQLWhitelistFPMatch {
				wlNodes, err := d.Parse(context.TODO(), wl.Value)
				if err != nil {
					return err
				}
				if len(wlNodes) != 1 {
					return driver.ErrNodesCountExceedOne
				}

				if nodes[0].Fingerprint == wlNodes[0].Fingerprint {
					whitelistMatch = true
				}
			} else {
				if wl.CapitalizedValue == strings.ToUpper(nodes[0].Text) {
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
		executeSQL.AuditFingerprint = utils.Md5String(string(append([]byte(result.Message()), []byte(nodes[0].Fingerprint)...)))

		l.WithFields(logrus.Fields{
			"SQL":    executeSQL.Content,
			"level":  executeSQL.AuditLevel,
			"result": executeSQL.AuditResult}).Info("audit finished")
	}
	return nil
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
