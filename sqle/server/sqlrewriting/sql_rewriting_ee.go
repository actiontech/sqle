//go:build enterprise
// +build enterprise

package sqlrewriting

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	pkgHttp "github.com/actiontech/dms/pkg/dms-common/pkg/http"
	"github.com/actiontech/sqle/sqle/common"
	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/actiontech/sqle/sqle/server/fillsql"
	"github.com/sirupsen/logrus"
	"github.com/ungerik/go-dry"
	"golang.org/x/text/language"
)

func getUrl() string {
	return config.GetOptions().SqleOptions.SQLRewritingConfig.RewritingURL
}

// Suggestion 定义单个重写建议的结构
type Suggestion struct {
	// 重写建议的规则名称(现有规则的RuleName)
	RuleName string `json:"rule_name"`
	// 重写建议的规则ID(知识库新规则的ID)
	RuleID string `json:"rule_id"`
	// 重写建议的类型
	Type string `json:"type"`
	// 重写描述
	Description string `json:"desc,omitempty"`
	// 重写SQL
	RewrittenSql string `json:"rewritten_sql,omitempty"`
	// 需要执行的DDL DCL描述
	DDLDCLDesc string `json:"ddl_dcl_desc,omitempty"`
	// 需要执行的DDL DCL
	DDLDCL string `json:"ddl_dcl,omitempty"`
}

// 获取外部进程完成重写的响应结构
type GetRewriteResponse struct {
	// 重写前SQL的业务描述
	BusinessDesc string `json:"business_desc"`
	// 重写前SQL的执行逻辑描述
	LogicDesc string `json:"logic_desc"`
	// 重写建议列表
	Suggestions []Suggestion `json:"suggestions"`
	// 重写后的SQL的业务差异描述
	BusinessNonEquivalent string `json:"business_non_equivalent,omitempty"`
	// 重写后SQL的业务描述
	BusinessDescAfterOptimize string `json:"business_desc_after_optimize"`
	// 重写后的SQL执行逻辑描述
	LogicDescAfterOptimize string `json:"logic_desc_after_optimize"`
	// 当开启EnableProgressiveRewrite时，表示所有规则重写是否完成
	IsRewriteDone bool `json:"is_rewrite_done"`
}

type Rule struct {
	// 审核规则Id
	Id string `json:"id"`
	// 审核规则参数
	Params map[string]string `json:"params"`
	// 规则审核结果信息
	Msg string `json:"msg"`
}

// 调用外部进程完成重写的请求结构
type CallRewriteSQLRequest struct {
	DBType                   string `json:"db_type"`
	Rules                    []Rule `json:"rules"`
	OriginalSql              string `json:"original_sql"`
	ProgressiveRewrittenSQL  string `json:"progressive_rewritten_sql"`
	TableStructures          string `json:"table_structures"`
	Explain                  string `json:"explain"`
	EnableStructureType      bool   `json:"enable_structure_type"`      // 是否启用涉及数据库结构化的重写
	EnableProgressiveRewrite bool   `json:"enable_progressive_rewrite"` // 是否启用渐进式重写模式：每次只重写一条规则，重写后重新审核，继续处理剩余触发的规则
}

// TODO: 更多预检查规则补充
var rulePreCheck = []string{"dml_enable_explain_pre_check"}

type ruleIdConvert struct {
	Name   string `json:"Name"`
	CH     string `json:"CH"`
	RuleId string `json:"RuleId"`
}

// 已经使用新规则ID实现的数据库类型插件，不需要进行规则ID转换
func noNeedConvert(dbType string) bool {
	switch strings.ToLower(dbType) {
	case "tbase", "hana":
		return true
	}
	return false
}

func getRuleIdConvert(dbType string) ([]ruleIdConvert, error) {
	switch strings.ToLower(dbType) {
	case "mysql":
		return MySQLRuleIdConvert, nil
	case "postgresql":
		return PostgreSQLRuleIdConvert, nil
	case "oracle":
		return OracleRuleIdConvert, nil
	default:
		return nil, fmt.Errorf("unsupported db type: %s", dbType)
	}
}

func ConvertRuleNameToRuleId(dbType string, ruleName string) (string, error) {
	if ruleName == "" {
		return "", fmt.Errorf("rule name is empty")
	}
	if noNeedConvert(dbType) {
		return ruleName, nil
	}
	ruleIdConvert, err := getRuleIdConvert(dbType)
	if err != nil {
		return "", err
	}

	for _, r := range ruleIdConvert {
		if r.Name == ruleName {
			return r.RuleId, nil
		}
	}
	return "", fmt.Errorf("can't find rule id for rule name: %s", ruleName)
}

func ConvertRuleIDToRuleName(dbType string, ruleId string) (string, error) {
	if ruleId == "" {
		return "", fmt.Errorf("rule id is empty")
	}
	if noNeedConvert(dbType) {
		return ruleId, nil
	}
	ruleIdConvert, err := getRuleIdConvert(dbType)
	if err != nil {
		return "", err
	}

	for _, r := range ruleIdConvert {
		if r.RuleId == ruleId {
			return r.Name, nil
		}
	}
	return "", fmt.Errorf("can't find rule name for rule id: %s", ruleId)
}

type SQLRewritingParams struct {
	Task                *model.Task             // 任务信息
	SQL                 *model.ExecuteSQL       // 需要重写的SQL
	TableStructures     []*driver.TableMeta     // 表结构
	Explain             *driverV2.ExplainResult // SQL 执行计划
	EnableStructureType bool                    // 是否启用涉及数据库结构化的重写
}

// ProgressiveRewriteSQL 启用渐进式重写模式：每次只重写一条规则，重写后重新审核，继续处理剩余触发的规则
func ProgressiveRewriteSQL(ctx context.Context, params *SQLRewritingParams) (*GetRewriteResponse, error) {
	l := log.NewEntry().WithField("get_rewrite_sql", "server")
	storage := model.GetStorage()
	lang := locale.Bundle.GetLangTagFromCtx(ctx)
	originalSQL := params.SQL.Content

	// 初始化返回结果
	ret := &GetRewriteResponse{}

	// 执行首次重写
	rewriteResult, err := performSQLRewriting(l, ctx, params, originalSQL, true /*开启渐进式重写*/)
	if err != nil {
		return nil, err
	}

	// 设置基础信息
	ret.BusinessDesc = rewriteResult.BusinessDesc
	ret.LogicDesc = rewriteResult.LogicDesc
	ret.LogicDescAfterOptimize = rewriteResult.LogicDescAfterOptimize
	ret.BusinessDescAfterOptimize = rewriteResult.BusinessDescAfterOptimize
	ret.BusinessNonEquivalent = rewriteResult.BusinessNonEquivalent

	// 如果已完成重写，直接返回结果
	if rewriteResult.IsRewriteDone {
		ret.Suggestions = rewriteResult.Suggestions
		return ret, nil
	}

	// 迭代重写直到完成
	for !rewriteResult.IsRewriteDone {

		// 没有重写完成，需要再次审核重写后的SQL，确认重写后的SQL是否触发了剩余的规则
		ruleNameNotRewritten := make([]string, 0)
		var sqlRewritten string
		var suggestionType string
		var ruleNameRewritten string
		for _, suggestion := range rewriteResult.Suggestions {
			if suggestion.RewrittenSql != "" {
				sqlRewritten = suggestion.RewrittenSql
				suggestionType = suggestion.Type
				rule, exist, err := storage.GetRule(suggestion.RuleName, params.Task.DBType)
				if err != nil {
					l.Errorf("get rule failed: %v", err)
					return nil, fmt.Errorf("get rule(%v) failed: %v", suggestion.RuleName, err)
				}
				if !exist {
					l.Errorf("rewritten rule not found: %s", suggestion.RuleName)
					return nil, fmt.Errorf("rule not found: %s", suggestion.RuleName)
				}
				ruleNameRewritten = rule.I18nRuleInfo.GetRuleInfoByLangTag(lang).Desc
				ret.Suggestions = append(ret.Suggestions, suggestion) // 保存已经重写的规则建议
			}
			if suggestion.RewrittenSql == "" {
				ruleNameNotRewritten = append(ruleNameNotRewritten, suggestion.RuleName) // 获取未重写的规则
			}
		}

		if sqlRewritten == "" {
			return nil, fmt.Errorf("rewritten sql is empty")
		}

		// 重新审核重写后的SQL
		l.Infof("audit rewritten sql: %v", sqlRewritten)
		task, err := server.AuditSQLByRuleNames(l, sqlRewritten, params.Task.DBType, params.Task.Instance, params.Task.Schema, ruleNameNotRewritten)
		if err != nil {
			l.Errorf("audit sqls failed: %v", err)
			return nil, fmt.Errorf("audit rewritten sql(%v) failed: %v", sqlRewritten, err)
		}
		if len(task.ExecuteSQLs) == 0 {
			return nil, fmt.Errorf("task.ExecuteSQLs is empty")
		}
		executeSQL := task.ExecuteSQLs[0]

		// 获取重写后的SQL不再触发的规则
		l.Infof("audit results: %v", executeSQL.AuditResults.String(ctx))
		ruleNamesResolved := findResolvedRules(ruleNameNotRewritten, executeSQL.AuditResults)
		l.Infof("rule name not rewritten: %v; rule name resolved: %v", ruleNameNotRewritten, ruleNamesResolved)

		for _, ruleName := range ruleNamesResolved {
			ret.Suggestions = append(ret.Suggestions, Suggestion{
				RuleName:     ruleName,
				RewrittenSql: sqlRewritten,
				Type:         suggestionType,
				Description:  fmt.Sprintf("根据规则“%v”重写后的SQL已不再触发本规则", ruleNameRewritten),
			})
		}

		// 再次获取重写后的SQL的执行计划
		explainResult := &driverV2.ExplainResult{}
		if res, err := ExplainTaskSQL(l, task); err != nil {
			l.Warnf("get explain failed: %v", err)
		} else {
			explainResult = res
		}

		params := &SQLRewritingParams{
			Task:                task,
			SQL:                 executeSQL,
			TableStructures:     params.TableStructures,
			Explain:             explainResult,
			EnableStructureType: params.EnableStructureType,
		}
		// 执行下一轮重写
		rewriteResult, err = performSQLRewriting(l, ctx, params, originalSQL, true /*开启渐进式重写*/)
		if err != nil {
			return nil, err
		}
		ret.BusinessDescAfterOptimize = rewriteResult.BusinessDescAfterOptimize
		ret.BusinessNonEquivalent = rewriteResult.BusinessNonEquivalent
	}
	ret.Suggestions = append(ret.Suggestions, rewriteResult.Suggestions...) // 添加最后一轮重写的建议(包括可能的一些无法重写的规则)
	return ret, nil
}

// performSQLRewriting 处理 SQL 重写的核心逻辑
func performSQLRewriting(l *logrus.Entry, ctx context.Context, params *SQLRewritingParams, originalSQL string, enableProgressiveRewrite bool) (*GetRewriteResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params is nil")
	}
	if params.SQL == nil {
		return nil, fmt.Errorf("sql is nil")
	}

	dbType := params.Task.DBType
	s := model.GetStorage()

	var rules []Rule
	var ruleNamesCanNotRewrite []string
	// 记录 ruleID 到 ruleName 的映射关系，用于处理多个 ruleName 对应同一个 ruleID 的情况
	ruleIDToNames := make(map[string][]string)
	// 用于去重的 ruleID 集合
	uniqueRuleIDs := make(map[string]struct{})
	// 定义返回结果
	reply := &GetRewriteResponse{}

	for _, ar := range params.SQL.AuditResults {
		if ar.RuleName == "" {
			continue
		}
		// 不对预检查规则进行重写
		if dry.StringInSlice(ar.RuleName, rulePreCheck) {
			ruleNamesCanNotRewrite = append(ruleNamesCanNotRewrite, ar.RuleName)
			continue
		}

		ruleID, err := ConvertRuleNameToRuleId(dbType, ar.RuleName)
		if err != nil {
			l.Errorf("can't convert rule name(%v) to rule id: %v", ar.RuleName, err)
			ruleNamesCanNotRewrite = append(ruleNamesCanNotRewrite, ar.RuleName)
			continue
		}

		// 记录 ruleID 到 ruleName 的映射
		ruleIDToNames[ruleID] = append(ruleIDToNames[ruleID], ar.RuleName)

		// 如果该 ruleID 已经存在，则跳过
		if _, exists := uniqueRuleIDs[ruleID]; exists {
			continue
		}
		uniqueRuleIDs[ruleID] = struct{}{}

		r, exist, err := s.GetRule(ar.RuleName, dbType)
		if err != nil {
			return nil, fmt.Errorf("get rule failed: %v", err)
		}
		if !exist {
			return nil, fmt.Errorf("rule not found: %s", ar.RuleName)
		}
		ruleParams := map[string]string{}
		for _, p := range r.Params {
			ruleParams[p.Key] = p.Value
		}

		rules = append(rules, Rule{
			Id:     ruleID,
			Params: ruleParams,
			Msg:    ar.GetAuditMsgByLangTag(language.Chinese),
		})

	}

	if len(rules) == 0 {
		// 没有规则需要重写
		l.Infof("no rules need to rewrite")
		reply.IsRewriteDone = true
	} else {
		// 定义要发送的参数
		req := &CallRewriteSQLRequest{
			DBType:                  dbType,
			Rules:                   rules,
			OriginalSql:             originalSQL,
			ProgressiveRewrittenSQL: params.SQL.Content,
			EnableStructureType:     params.EnableStructureType,
		}

		if enableProgressiveRewrite {
			req.EnableProgressiveRewrite = true
		}

		if s, err := json.Marshal(params.Explain); err != nil {
			return nil, fmt.Errorf("marshal explain failed: %v", err)
		} else {
			req.Explain = string(s)
		}

		for _, table := range params.TableStructures {
			s, err := json.Marshal(table)
			if err != nil {
				return nil, fmt.Errorf("marshal table structure failed: %v", err)
			}
			req.TableStructures += string(s) + "\n\n"
		}

		// 定义 API 端点
		apiURL := getUrl()

		// 先测试API端点的连通性
		connectivityTimeout := 5 * time.Second // 设置连通性测试的超时时间
		if err := checkAPIConnectivity(apiURL, connectivityTimeout); err != nil {
			return nil, fmt.Errorf("请检查合规重写配置及联通性: API connectivity check failed: %v.", err)
		}

		// 发送 HTTP POST 请求
		var callRewriteSQLTimeout int64 = 600 // 设置重写请求的超时时间
		if err := pkgHttp.POST(pkgHttp.SetTimeoutValueContext(ctx, callRewriteSQLTimeout), apiURL, nil, req, reply); err != nil {
			return nil, fmt.Errorf("failed to call %v: %v", apiURL, err)
		}

		// 处理返回结果，将同一个 ruleID 对应的多个 ruleName 都添加到建议中(重写功能使用了新的规则ID，不需要现有的规则名称，这里填充回现有的规则名称只是为了页面展示)
		var expandedSuggestions []Suggestion
		for _, suggestion := range reply.Suggestions {
			ruleNames := ruleIDToNames[suggestion.RuleID]
			for _, ruleName := range ruleNames {
				newSuggestion := suggestion
				newSuggestion.RuleName = ruleName
				expandedSuggestions = append(expandedSuggestions, newSuggestion)
			}
		}
		reply.Suggestions = expandedSuggestions
	}

	// 对于无法重写的规则，需要将其加入到Suggestions中
	for _, rn := range ruleNamesCanNotRewrite {
		reply.Suggestions = append(reply.Suggestions, Suggestion{
			RuleName:    rn,
			Type:        "other",
			Description: "需要人工处理",
		})
	}
	return reply, nil

}

func checkAPIConnectivity(apiURL string, timeout time.Duration) error {
	parsedURL, err := url.Parse(apiURL)
	if err != nil {
		return fmt.Errorf("parse rewriting_url failed: %v, please check config.yml", err)
	}

	// 仅支持http和https协议
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("unsupported URL scheme: %s", parsedURL.Scheme)
	}

	host := parsedURL.Host
	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		return fmt.Errorf("cannot connect to API endpoint %s: %v", apiURL, err)
	}
	defer conn.Close()
	return nil
}

func ExplainTaskSQL(l *logrus.Entry, task *model.Task) (res *driverV2.ExplainResult, err error) {
	if task.Instance == nil {
		return nil, fmt.Errorf("task.Instance is nil")
	}

	if len(task.ExecuteSQLs) == 0 {
		return nil, fmt.Errorf("task.ExecuteSQLs is empty")
	}
	taskSql := task.ExecuteSQLs[0]
	instance := task.Instance
	sqlContent, err := fillsql.FillingSQLWithParamMarker(taskSql.Content, task)
	if err != nil {
		l.Errorf("fill param marker sql failed: %v", err)
		sqlContent = taskSql.Content
	}

	dsn, err := common.NewDSN(instance, task.Schema)
	if err != nil {
		return nil, err
	}

	plugin, err := driver.GetPluginManager().
		OpenPlugin(l, instance.DbType, &driverV2.Config{DSN: dsn})
	if err != nil {
		return nil, err
	}
	defer plugin.Close(context.TODO())

	if !driver.GetPluginManager().
		IsOptionalModuleEnabled(instance.DbType, driverV2.OptionalModuleExplain) {
		return nil, driver.NewErrPluginAPINotImplement(driverV2.OptionalModuleExplain)
	}

	return plugin.Explain(context.TODO(), &driverV2.ExplainConf{Sql: sqlContent})
}

// findResolvedRules 找出重写后不再触发的规则名称
// params: originalRules 原始触发的规则名称列表
// params: newAuditResults 重写后的SQL审核结果
// return 已解决（不再触发）的规则名称列表
func findResolvedRules(originalRules []string, newAuditResults model.AuditResults) []string {
	// 使用 map 存储当前仍触发的规则，便于快速查找
	stillTriggered := make(map[string]bool)
	for _, result := range newAuditResults {
		stillTriggered[result.RuleName] = true
	}

	// 找出原始规则中不再出现在新审核结果中的规则名称
	resolved := make([]string, 0, len(originalRules))
	for _, ruleName := range originalRules {
		if !stillTriggered[ruleName] {
			resolved = append(resolved, ruleName)
		}
	}

	return resolved
}
