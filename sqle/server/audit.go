package server

import (
	"context"
	"fmt"
	"math"
	"runtime/debug"
	"strings"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func Audit(l *logrus.Entry, task *model.Task, projectId *model.ProjectUID, ruleTemplateName string) (err error) {
	return HookAudit(l, task, &EmptyAuditHook{}, projectId, ruleTemplateName)
}

func HookAudit(l *logrus.Entry, task *model.Task, hook AuditHook, projectId *model.ProjectUID, ruleTemplateName string) (err error) {
	st := model.GetStorage()
	rules, customRules, err := st.GetAllRulesByTmpNameAndProjectIdInstanceDBType(ruleTemplateName, string(*projectId), task.Instance, task.DBType)
	if err != nil {
		return err
	}
	plugin, err := newDriverManagerWithAudit(l, task.Instance, task.Schema, task.DBType, rules)
	if err != nil {
		return err
	}
	defer plugin.Close(context.TODO())

	// possible task is self build object, not model.Task{}
	if task.Instance == nil {
		task.Instance = &model.Instance{ProjectId: string(*projectId)}
	}

	return hookAudit(l, task, plugin, hook, customRules)
}

const AuditSchema = "AuditSchema"

func DirectAuditByInstance(l *logrus.Entry, sql, schemaName string, instance *model.Instance) (*model.Task, error) {
	st := model.GetStorage()
	rules, customRules, err := st.GetAllRulesByTmpNameAndProjectIdInstanceDBType("", "", instance, instance.DbType)
	if err != nil {
		return nil, err
	}
	plugin, err := newDriverManagerWithAudit(l, instance, schemaName, instance.DbType, rules)
	if err != nil {
		return nil, err
	}
	defer plugin.Close(context.TODO())

	task, err := convertSQLsToTask(sql, plugin)
	if err != nil {
		return nil, err
	}
	task.Instance = instance

	return task, audit(l, task, plugin, customRules)
}

func AuditSQLByDBType(l *logrus.Entry, sql string, dbType string) (*model.Task, error) {
	st := model.GetStorage()
	rules, customRules, err := st.GetAllRulesByTmpNameAndProjectIdInstanceDBType("", "", nil, dbType)
	if err != nil {
		return nil, err
	}
	plugin, err := newDriverManagerWithAudit(l, nil, "", dbType, rules)
	if err != nil {
		return nil, err
	}
	defer plugin.Close(context.TODO())

	return AuditSQLByDriver(l, sql, plugin, customRules)
}

func AuditSQLByDriver(l *logrus.Entry, sql string, p driver.Plugin, customRules []*model.CustomRule) (*model.Task, error) {
	task, err := convertSQLsToTask(sql, p)
	if err != nil {
		return nil, err
	}
	return task, audit(l, task, p, customRules)
}

func convertSQLsToTask(sql string, p driver.Plugin) (*model.Task, error) {
	task := &model.Task{}
	nodes, err := p.Parse(context.TODO(), sql)
	if err != nil {
		return nil, err
	}
	for n, node := range nodes {
		task.ExecuteSQLs = append(task.ExecuteSQLs, &model.ExecuteSQL{
			BaseSQL: model.BaseSQL{
				Number:  uint(n + 1),
				Content: node.Text,
			},
		})
	}
	return task, nil
}

func audit(l *logrus.Entry, task *model.Task, p driver.Plugin, customRules []*model.CustomRule) (err error) {
	return hookAudit(l, task, p, &EmptyAuditHook{}, customRules)
}

type AuditHook interface {
	BeforeAudit(sql *model.ExecuteSQL)
	AfterAudit(sql *model.ExecuteSQL)
}

type EmptyAuditHook struct{}

func (e *EmptyAuditHook) BeforeAudit(sql *model.ExecuteSQL) {}

func (e *EmptyAuditHook) AfterAudit(sql *model.ExecuteSQL) {}

func hookAudit(l *logrus.Entry, task *model.Task, p driver.Plugin, hook AuditHook, customRules []*model.CustomRule) (err error) {
	defer func() {
		if errRecover := recover(); errRecover != nil {
			debug.PrintStack()
			// 为了将panic信息返回给调用者
			err = errors.New("An unknown error occurred, check std.log for details")
			l.Errorf("hookAudit panic: %v", errRecover)
		}
	}()

	st := model.GetStorage()

	projectId := ""
	if task.Instance != nil {
		projectId = task.Instance.ProjectId
	}
	whitelist, err := st.GetSqlWhitelistByProjectId(projectId)
	if err != nil {
		return err
	}

	auditSqls := []*model.ExecuteSQL{}
	sqls := []string{}
	nodes := []driverV2.Node{}
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
		node, err := parse(l, p, strings.TrimSpace(executeSQL.Content))
		if err != nil {
			return err
		}
		var whitelistMatch bool
		for _, wl := range whitelist {
			if wl.MatchType == model.SQLWhitelistFPMatch {
				wlNode, err := parse(l, p, wl.Value)
				if err != nil {
					l.Errorf("parse whitelist sql error: %v,please check the accuracy of whitelist SQL: %s", err, wl.Value)
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
		if whitelistMatch {
			result := driverV2.NewAuditResults()
			result.Add(driverV2.RuleLevelNormal, "", "白名单")
			executeSQL.AuditStatus = model.SQLAuditStatusFinished
			executeSQL.AuditLevel = string(result.Level())
			executeSQL.AuditFingerprint = utils.Md5String(string(append([]byte(result.Message()), []byte(node.Fingerprint)...)))
			appendExecuteSqlResults(executeSQL, result)
		} else {
			auditSqls = append(auditSqls, executeSQL)
			sqls = append(sqls, executeSQL.Content)
			nodes = append(nodes, node)
		}
	}
	for _, sql := range auditSqls {
		hook.BeforeAudit(sql)
	}

	results, err := p.Audit(context.TODO(), sqls)
	if err != nil {
		return err
	}
	if len(results) != len(sqls) {
		return fmt.Errorf("audit results [%d] does not match the number of SQL [%d]", len(results), len(sqls))
	}
	CustomRuleAudit(l, task, sqls, results, customRules)
	for i, sql := range auditSqls {
		hook.AfterAudit(sql)
		sql.AuditStatus = model.SQLAuditStatusFinished
		sql.AuditLevel = string(results[i].Level())
		sql.AuditFingerprint = utils.Md5String(string(append([]byte(results[i].Message()), []byte(nodes[i].Fingerprint)...)))
		appendExecuteSqlResults(sql, results[i])
	}

	ReplenishTaskStatistics(task)
	return nil
}

func ReplenishTaskStatistics(task *model.Task) {
	var normalCount float64
	maxAuditLevel := driverV2.RuleLevelNull
	for _, executeSQL := range task.ExecuteSQLs {
		if driverV2.RuleLevelNormal.MoreOrEqual(driverV2.RuleLevel(executeSQL.AuditLevel)) {
			normalCount += 1
		}
		if driverV2.RuleLevel(executeSQL.AuditLevel).More(maxAuditLevel) {
			maxAuditLevel = driverV2.RuleLevel(executeSQL.AuditLevel)
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
			switch driverV2.RuleLevel(e.AuditLevel) {
			case driverV2.RuleLevelError:
				numberOfLessThanError++
			case driverV2.RuleLevelWarn:
				numberOfLessThanWarn++
			case driverV2.RuleLevelNotice:
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

func parse(l *logrus.Entry, p driver.Plugin, sql string) (node driverV2.Node, err error) {
	nodes, err := p.Parse(context.TODO(), sql)
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

func genRollbackSQL(l *logrus.Entry, task *model.Task, p driver.Plugin) ([]*model.RollbackSQL, error) {
	rollbackSQLs := make([]*model.RollbackSQL, 0, len(task.ExecuteSQLs))
	for _, executeSQL := range task.ExecuteSQLs {
		rollbackSQL, reason, err := p.GenRollbackSQL(context.TODO(), executeSQL.Content)
		if err != nil && session.IsParseShowCreateTableContentErr(err) {
			l.Errorf("gen rollback sql error, %v", err) // todo #1630 临时跳过创表语句解析错误
			return nil, nil
		} else if err != nil {
			l.Errorf("gen rollback sql error, %v", err)
			return nil, err
		}
		result := driverV2.NewAuditResults()
		for i := range executeSQL.AuditResults {
			ar := executeSQL.AuditResults[i]
			result.Add(driverV2.RuleLevel(ar.Level), ar.RuleName, ar.Message)
		}
		result.Add(driverV2.RuleLevelNotice, "", reason)

		executeSQL.AuditLevel = string(result.Level())
		appendExecuteSqlResults(executeSQL, result)

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

func appendExecuteSqlResults(executeSQL *model.ExecuteSQL, result *driverV2.AuditResults) {
	for i := range result.Results {
		ar := result.Results[i]
		executeSQL.AuditResults.Append(string(ar.Level), ar.RuleName, ar.Message)
	}
}
