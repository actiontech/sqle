package server

import (
	"context"
	"fmt"
	"math"
	"runtime/debug"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"

	// "github.com/actiontech/sqle/sqle/driver/mysql/session"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// AfterAuditHook 审核完成后的钩子函数，由 EE 版本注册实现
// 用于记录代码规范规则触发统计等
var AfterAuditHook func(task *model.Task)

// AfterWorkflowCreateHook 工单创建后的钩子函数，由 EE 版本注册实现
// 用于记录提交工单时的规则触发统计
var AfterWorkflowCreateHook func(tasks []*model.Task, workflowID string)

func Audit(l *logrus.Entry, task *model.Task, projectId *model.ProjectUID, ruleTemplateName string) (err error) {
	return HookAudit(l, task, &EmptyAuditHook{}, projectId, ruleTemplateName)
}

func HookAudit(l *logrus.Entry, task *model.Task, hook AuditHook, projectId *model.ProjectUID, ruleTemplateName string) (err error) {
	st := model.GetStorage()
	if projectId == nil {
		return fmt.Errorf("HookAudit error because projectId is nil, taskId: %v", task.ID)
	}
	rules, customRules, err := st.GetAllRulesByTmpNameAndProjectIdInstanceDBType(ruleTemplateName, string(*projectId), task.Instance, task.DBType)
	if err != nil {
		return err
	}
	if task == nil {
		return fmt.Errorf("HookAudit error because task is nil, projectId: %v", string(*projectId))
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
	return hookAudit(l, task, plugin, hook, string(*projectId), customRules)
}

const AuditSchema = "AuditSchema"

func DirectAuditByInstance(l *logrus.Entry, sql, schemaName string, instance *model.Instance, ruleTemplateName string) (*model.Task, error) {
	st := model.GetStorage()
	rules, customRules, err := st.GetAllRulesByTmpNameAndProjectIdInstanceDBType(ruleTemplateName, instance.ProjectId, instance, instance.DbType)
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

	return task, audit(instance.ProjectId, l, task, plugin, customRules)
}

func AuditSQLByDBType(l *logrus.Entry, sql string, dbType string, projectId string, ruleTemplateName string) (*model.Task, error) {
	st := model.GetStorage()
	rules, customRules, err := st.GetAllRulesByTmpNameAndProjectIdInstanceDBType(ruleTemplateName, projectId, nil, dbType)
	if err != nil {
		return nil, err
	}
	plugin, err := newDriverManagerWithAudit(l, nil, "", dbType, rules)
	if err != nil {
		return nil, err
	}
	defer plugin.Close(context.TODO())

	return AuditSQLByDriver(projectId, l, sql, plugin, customRules)
}

func AuditSQLByRuleNames(l *logrus.Entry, sql string, dbType string, instance *model.Instance, schemaName string, ruleNames []string) (*model.Task, error) {
	st := model.GetStorage()
	rules, err := st.GetRulesByNamesAndDBType(ruleNames, dbType)
	if err != nil {
		return nil, err
	}
	plugin, err := newDriverManagerWithAudit(l, instance, schemaName, dbType, rules)
	if err != nil {
		return nil, err
	}
	defer plugin.Close(context.TODO())

	task, err := AuditSQLByDriver(instance.ProjectId, l, sql, plugin, nil)
	task.DBType = dbType
	task.Instance = instance
	task.InstanceId = instance.ID
	task.Schema = schemaName
	return task, err
}

func AuditSQLByDriver(projectId string, l *logrus.Entry, sql string, p driver.Plugin, customRules []*model.CustomRule) (*model.Task, error) {
	task, err := convertSQLsToTask(sql, p)
	if err != nil {
		return nil, err
	}
	return task, audit(projectId, l, task, p, customRules)
}

func convertSQLsToTask(sql string, p driver.Plugin) (*model.Task, error) {
	task := &model.Task{}
	executeSQLs, err := BuildExecuteSQLsFromSQL(context.TODO(), p, sql, BuildExecuteSQLsOptions{})
	if err != nil {
		return nil, err
	}
	task.ExecuteSQLs = executeSQLs
	return task, nil
}

type BuildExecuteSQLsOptions struct {
	StartNumber uint
	SourceFile  string
	StartLine   uint64
}

func BuildExecuteSQLsFromSQL(ctx context.Context, p driver.Plugin, sqlText string, opts BuildExecuteSQLsOptions) ([]*model.ExecuteSQL, error) {
	trimmedSQL := strings.TrimSpace(sqlText)
	if trimmedSQL == "" {
		return nil, nil
	}
	number := opts.StartNumber
	if number == 0 {
		number = 1
	}
	nodes, err := p.Parse(ctx, sqlText)
	if err != nil || len(nodes) == 0 {
		return []*model.ExecuteSQL{{
			BaseSQL: model.BaseSQL{
				Number:     number,
				Content:    trimmedSQL,
				SourceFile: opts.SourceFile,
				StartLine:  opts.StartLine,
			},
		}}, nil
	}

	executeSQLs := make([]*model.ExecuteSQL, 0, len(nodes))
	for i, node := range nodes {
		startLine := opts.StartLine
		if startLine == 0 {
			startLine = node.StartLine
		}
		executeSQLs = append(executeSQLs, &model.ExecuteSQL{
			BaseSQL: model.BaseSQL{
				Number:      number + uint(i),
				Content:     node.Text,
				SourceFile:  opts.SourceFile,
				StartLine:   startLine,
				SQLType:     node.Type,
				ExecBatchId: node.ExecBatchId,
			},
		})
	}
	return executeSQLs, nil
}

func audit(projectId string, l *logrus.Entry, task *model.Task, p driver.Plugin, customRules []*model.CustomRule) (err error) {
	return hookAudit(l, task, p, &EmptyAuditHook{}, projectId, customRules)
}

type AuditHook interface {
	BeforeAudit(sql *model.ExecuteSQL)
	AfterAudit(sql *model.ExecuteSQL)
}

type EmptyAuditHook struct{}

func (e *EmptyAuditHook) BeforeAudit(sql *model.ExecuteSQL) {}

func (e *EmptyAuditHook) AfterAudit(sql *model.ExecuteSQL) {}

func hookAudit(l *logrus.Entry, task *model.Task, p driver.Plugin, hook AuditHook, projectId string, customRules []*model.CustomRule) (err error) {
	defer func() {
		if errRecover := recover(); errRecover != nil {
			debug.PrintStack()
			// 为了将panic信息返回给调用者
			err = errors.New("An unknown error occurred, check std.log for details")
			l.Errorf("hookAudit panic: %v", errRecover)
		}
	}()

	st := model.GetStorage()
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
			appendManualConfirmWarn(executeSQL, executeSQL.Content, err)
			continue
		}
		var whitelistMatch bool
		var matchedWhitelistID uint
		for _, wl := range whitelist {
			if wl.MatchType == model.SQLWhitelistFPMatch {
				wlNode, err := parse(l, p, wl.Value)
				if err != nil {
					l.Errorf("parse whitelist sql error: %v,please check the accuracy of whitelist SQL: %s", err, wl.Value)
				}
				if node.Fingerprint == wlNode.Fingerprint {
					matchedWhitelistID = wl.ID
					whitelistMatch = true
				}
			} else {
				if wl.CapitalizedValue == strings.ToUpper(node.Text) {
					matchedWhitelistID = wl.ID
					whitelistMatch = true
				}
			}
		}
		if whitelistMatch {
			result := driverV2.NewAuditResults()
			result.Add(driverV2.RuleLevelNormal, "", plocale.Bundle.LocalizeAll(plocale.AuditResultMsgExcludedSQL))
			executeSQL.AuditStatus = model.SQLAuditStatusFinished
			executeSQL.AuditLevel = string(result.Level())
			executeSQL.AuditFingerprint = utils.Md5String(string(append([]byte(result.Message()), []byte(node.Fingerprint)...)))
			executeSQL.SqlFingerprint = node.Fingerprint
			appendExecuteSqlResults(executeSQL, result)
			if err := st.UpdateSqlWhitelistMatchedInfo(matchedWhitelistID, 1, time.Now()); err != nil {
				l.Errorf("update sql whitelist matched info error: %v", err)
			}
		} else {
			auditSqls = append(auditSqls, executeSQL)
			sqls = append(sqls, executeSQL.Content)
			nodes = append(nodes, node)
		}
	}
	if len(sqls) > 0 {
		for _, sql := range auditSqls {
			hook.BeforeAudit(sql)
		}

		results, err := p.Audit(context.TODO(), sqls)
		if err != nil {
			for i, sql := range auditSqls {
				result, singleErr := auditSingleSQL(l, p, task, sqls[i], customRules)
				hook.AfterAudit(sql)
				if singleErr != nil {
					appendManualConfirmWarn(sql, sqls[i], singleErr)
					continue
				}
				sql.AuditStatus = model.SQLAuditStatusFinished
				sql.AuditLevel = string(result.Level())
				sql.AuditFingerprint = utils.Md5String(string(append([]byte(result.Message()), []byte(nodes[i].Fingerprint)...)))
				sql.SqlFingerprint = nodes[i].Fingerprint
				appendExecuteSqlResults(sql, result)
			}
			ReplenishTaskStatistics(task)
			if AfterAuditHook != nil {
				go AfterAuditHook(task)
			}
			return nil
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
			sql.SqlFingerprint = nodes[i].Fingerprint
			appendExecuteSqlResults(sql, results[i])
		}
	}

	ReplenishTaskStatistics(task)

	// 审核完成后，记录代码规范规则触发统计（异步，不影响审核主流程）
	if AfterAuditHook != nil {
		go AfterAuditHook(task)
	}

	return nil
}

func ReplenishTaskStatistics(task *model.Task) {
	if len(task.ExecuteSQLs) == 0 {
		task.PassRate = 1
		task.AuditLevel = string(driverV2.RuleLevelNull)
		task.Score = scoreTask(task)
		task.Status = model.TaskStatusAudited
		return
	}
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

func auditSingleSQL(l *logrus.Entry, p driver.Plugin, task *model.Task, sql string, customRules []*model.CustomRule) (*driverV2.AuditResults, error) {
	results, err := p.Audit(context.TODO(), []string{sql})
	if err != nil {
		return nil, err
	}
	if len(results) != 1 {
		return nil, fmt.Errorf("audit results [%d] does not match the number of SQL [1]", len(results))
	}
	CustomRuleAudit(l, task, []string{sql}, results, customRules)
	return results[0], nil
}

func appendManualConfirmWarn(executeSQL *model.ExecuteSQL, sql string, err error) {
	result := driverV2.NewAuditResults()
	result.AddResultWithError(driverV2.RuleLevelWarn, "", err.Error(), false, plocale.Bundle.LocalizeAll(plocale.UnsupportedSyntaxError))
	executeSQL.AuditStatus = model.SQLAuditStatusFinished
	executeSQL.AuditLevel = string(result.Level())
	executeSQL.AuditFingerprint = utils.Md5String(result.Message() + strings.TrimSpace(sql))
	executeSQL.SqlFingerprint = utils.Md5String(strings.TrimSpace(sql))
	appendExecuteSqlResults(executeSQL, result)
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

func appendExecuteSqlResults(executeSQL *model.ExecuteSQL, result *driverV2.AuditResults) {
	for i := range result.Results {
		executeSQL.AuditResults.Append(result.Results[i])
	}
}
