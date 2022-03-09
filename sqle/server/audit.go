package server

import (
	"context"
	"fmt"
	"math"
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

	replenishTaskStatistics(task)

	return nil
}

func replenishTaskStatistics(task *model.Task) {
	var normalCount float64
	maxAuditLevel := driver.RuleLevelNormal
	for _, executeSQL := range task.ExecuteSQLs {
		if executeSQL.AuditLevel == string(driver.RuleLevelNormal) {
			normalCount += 1
		}
		if driver.RuleLevel(executeSQL.AuditLevel).More(maxAuditLevel) {
			maxAuditLevel = driver.RuleLevel(executeSQL.AuditLevel)
		}
	}
	task.PassRate = utils.Round(normalCount/float64(len(task.ExecuteSQLs)), 4)
	task.AuditLevel = string(maxAuditLevel)
	task.Score = scoreTask(task)

	task.Status = model.TaskStatusAudited
}

// Scoring rules from https://github.com/actiontech/sqle/issues/284
func scoreTask(task *model.Task) int32 {
	if len(task.ExecuteSQLs) == 0 {
		return 0
	}

	var (
		numberOfTask           float64
		numberOfLessThanError  float64
		numberOfLessThanWarn   float64
		numberOfLessThanNotice float64
		errorRate              float64
		warnRate               float64
		noticeRate             float64
		totalScore             float64
	)
	{ // ready to work
		numberOfTask = float64(len(task.ExecuteSQLs))

		for _, e := range task.ExecuteSQLs {
			switch driver.RuleLevel(e.AuditLevel) {
			case driver.RuleLevelError:
				numberOfLessThanError++
			case driver.RuleLevelWarn:
				numberOfLessThanWarn++
			case driver.RuleLevelNotice:
				numberOfLessThanNotice++
			}
		}

		numberOfLessThanNotice = numberOfLessThanNotice + numberOfLessThanWarn + numberOfLessThanError
		numberOfLessThanWarn = numberOfLessThanWarn + numberOfLessThanError

		errorRate = numberOfLessThanError / numberOfTask
		warnRate = numberOfLessThanWarn / numberOfTask
		noticeRate = numberOfLessThanNotice / numberOfTask
	}
	{ // calculate the total score
		// pass rate score
		totalScore = task.PassRate * 30
		// SQL occurrence probability below error level
		totalScore += (1 - errorRate) * 15
		// SQL occurrence probability below warn level
		totalScore += (1 - warnRate) * 10
		// SQL occurrence probability below notice level
		totalScore += (1 - noticeRate) * 5
		// SQL without error level
		if errorRate == 0 {
			totalScore += 15
		}
		// SQL without warn level
		if warnRate == 0 {
			totalScore += 10
		}
		// SQL without notice level
		if noticeRate == 0 {
			totalScore += 5
		}
		// the proportion of SQL with the level below error exceeds 90%
		if errorRate < 0.1 {
			totalScore += 5
		}
		// the proportion of SQL with the level below warn exceeds 90%
		if warnRate < 0.1 {
			totalScore += 3
		}
		// the proportion of SQL with the level below warn exceeds 90%
		if noticeRate < 0.1 {
			totalScore += 2
		}
	}

	return int32(math.Floor(totalScore))
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
