//go:build enterprise
// +build enterprise

package optimization

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"
	"sync"

	"github.com/actiontech/sqle/sqle/driver"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/model"
	opt "github.com/actiontech/sqle/sqle/server/optimization/rule"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/sirupsen/logrus"
)

// regIndexesRecommended to match IndexesRecommended in pawsql response
var regIndexesRecommended = regexp.MustCompile(`(?i)CREATE INDEX .+?\);`)

// OptimizationServer
type OptimizationOnlinePawSQLServer struct {
	OptimizationPawSQLServer
	Instance *model.Instance
	Schema   string
}

type OptimizationPawSQLServer struct {
	logger *logrus.Entry
}

type OptimizationWorkspace struct {
	DBName     string
	SchemaName string
	Host       string
	Port       string
	DbUser     string
	DbPassword string
	DDL        string
}

func (o OptimizationWorkspace) MD5() string {
	return utils.Md5String(o.DBName + o.SchemaName + o.Host + o.Port + o.DbUser + o.DbPassword + o.DDL)
}

var cacheOptimizationWorkspace sync.Map // map[string] /*MD5(dbName+schemaName++ Host + Port + DbUser + DbPassword + DDL)*/ string /*workspaceId*/

func NewOptimizationOnlinePawSQLServer(logger *logrus.Entry, instance *model.Instance, schema string) *OptimizationOnlinePawSQLServer {
	return &OptimizationOnlinePawSQLServer{
		OptimizationPawSQLServer: OptimizationPawSQLServer{
			logger: logger,
		},
		Instance: instance,
		Schema:   schema,
	}
}

// 在线模式：创建空间后创建分析任务
func (a *OptimizationOnlinePawSQLServer) Optimizate(ctx context.Context, OptimizationSQL string) (optimizationInfo *model.SQLOptimizationRecord, err error) {
	a.logger.Debugf("Optimization SQL %v", OptimizationSQL)
	optimizationWorkspace := OptimizationWorkspace{
		DBName:     a.Instance.Name,
		SchemaName: a.Schema,
		DDL:        OptimizationSQL,
	}

	workspaceId, ok := cacheOptimizationWorkspace.Load(optimizationWorkspace.MD5())
	if !ok {
		// 创建pawsql的空间,内存缓存
		workspaceId, err = a.createWorkspaceOnline(ctx, a.Instance, a.Schema)
		if err != nil {
			return nil, err
		}
		cacheOptimizationWorkspace.Store(optimizationWorkspace.MD5(), workspaceId)
	}

	// 创建分析任务
	id, err := a.createOptimization(ctx, workspaceId.(string), a.Instance, OptimizationSQL, "online")
	if err != nil {
		return nil, err
	}
	return a.getOptimizationInfo(ctx, id, a.Instance.DbType)
}

// 获取优化详情
func (a *OptimizationOnlinePawSQLServer) getOptimizationInfo(ctx context.Context, analysisId string, dbType string) (optimizationInfo *model.SQLOptimizationRecord, err error) {
	a.logger.Debugf("get Optimization info id : %v", analysisId)
	// 获取优化任务概览
	summary, err := a.getOptimizationSummary(ctx, analysisId)
	if err != nil {
		return optimizationInfo, err
	}
	a.logger.Debugf("get Optimization summary %v", summary)
	optimizationIndexInfo := trimKeyWord4Slice(summary.OptimizationIndexInfo)
	optimizationInfo = &model.SQLOptimizationRecord{
		NumberOfQuery:          summary.BasicSummary.NumberOfQuery,
		NumberOfSyntaxError:    summary.BasicSummary.NumberOfSyntaxError,
		NumberOfRewrite:        summary.BasicSummary.NumberOfRewrite,
		NumberOfRewrittenQuery: summary.BasicSummary.NumberOfRewrittenQuery,
		NumberOfIndex:          summary.BasicSummary.NumberOfIndex,
		NumberOfQueryIndex:     summary.BasicSummary.NumberOfQueryIndex,
		IndexRecommendations:   optimizationIndexInfo,
	}

	var performanceImprove float64
	for _, statementInfo := range summary.SummaryStatementInfo {

		// 获取优化任务详情
		detail, err := a.getOptimizationDetail(ctx, statementInfo.OptimizationStmtId)
		if err != nil {
			return optimizationInfo, err
		}

		a.logger.Debugf("get Optimization detail %v", detail)
		triggeredRule := make([]model.RewriteRule, 0)
		for _, v := range detail.RewrittenQuery {
			name, message := convertRuleNameAndMessage(ctx, v.RuleCode, v.RuleNameZh, dbType)
			triggeredRule = append(triggeredRule, model.RewriteRule{
				RuleName:            name,
				Message:             message,
				RewrittenQueriesStr: v.RewrittenQueriesStr,
				ViolatedQueriesStr:  v.ViolatedQueriesStr,
			})
		}

		// 触发多条重写规则的SQL是一致的，只需获取任意一条即可
		var optimizedSQL string
		if len(detail.RewrittenQuery) != 0 && detail.RewrittenQuery[0].RewrittenQueriesStr != "" {
			sqls := make([]string, 0)
			err = json.Unmarshal([]byte(detail.RewrittenQuery[0].RewrittenQueriesStr), &sqls)
			if err != nil {
				a.logger.Errorf("unmarshal rewriteQueriesStr error %v", err)
			}
			if len(sqls) > 0 {
				optimizedSQL = sqls[0]
			}
		}

		// 移除PAWSQL关键字
		indexRecommendeds := trimKeyWord4Slice(getIndexesRecommendedFromMD(detail.DetailMarkdownZh))
		contributingIndices := trimKeyWord4Slice([]string{statementInfo.ContributingIndices})
		detail.ValidationDetails.BeforePlan = trimKeyWord(detail.ValidationDetails.BeforePlan)
		detail.ValidationDetails.AfterPlan = trimKeyWord(detail.ValidationDetails.AfterPlan)

		// 详情和列表的性能提升值保持一致
		detail.ValidationDetails.PerformImprovePer = statementInfo.Performance
		optimizationInfo.OptimizationSQLs = append(optimizationInfo.OptimizationSQLs, &model.OptimizationSQL{
			OriginalSQL:              detail.StmtText,
			OptimizedSQL:             optimizedSQL,
			NumberOfRewrite:          statementInfo.NumberOfRewrite,
			NumberOfSyntaxError:      statementInfo.NumberOfSyntaxError,
			NumberOfIndex:            statementInfo.NumberOfIndex,
			NumberOfHitIndex:         statementInfo.NumberOfHitIndex,
			Performance:              statementInfo.Performance,
			ContributingIndices:      contributingIndices[0],
			TriggeredRules:           triggeredRule,
			IndexRecommendations:     indexRecommendeds,
			ExplainValidationDetails: model.ExplainValidationDetail{BeforeCost: detail.ValidationDetails.BeforeCost, AfterCost: detail.ValidationDetails.AfterCost, BeforePlan: detail.ValidationDetails.BeforePlan, AfterPlan: detail.ValidationDetails.AfterPlan, PerformImprovePer: detail.ValidationDetails.PerformImprovePer},
		})
		performanceImprove += statementInfo.Performance
	}
	// 概览接口返回的提升数据异常，自行计算
	optimizationInfo.PerformanceImprove = performanceImprove / float64(len(summary.SummaryStatementInfo))
	return
}

func getIndexesRecommendedFromMD(md string) []string {
	start := strings.Index(md, "推荐的索引")
	if start < 0 {
		return []string{}
	}
	then := strings.Index(md[start:], "\n#")
	if then < 0 {
		return []string{}
	}
	createIndexes := regIndexesRecommended.FindAllString(md[start:start+then], -1)
	return createIndexes
}

func trimKeyWord(s string) string {
	return strings.ReplaceAll(s, "PAWSQL", "OPTIMIZATION")
}

func trimKeyWord4Slice(slice []string) (ret []string) {
	ret = make([]string, len(slice))
	for k := range slice {
		ret[k] = trimKeyWord(slice[k])
	}
	return
}

func getOptimizationReqRules(instance *model.Instance) ([]*CreateOptimizationRules, error) {
	store := model.GetStorage()
	templateRules, err := store.GetAllOptimizationRulesByInstance(instance)
	if err != nil {
		return nil, err
	}
	reqRules := convertRulesToOptimizationReqRules(templateRules, instance.DbType)
	return reqRules, nil
}

func convertRulesToOptimizationReqRules(templateRules []*model.Rule, dbType string) []*CreateOptimizationRules {
	allRules := opt.OptimizationRuleMap[dbType]
	reqRules := []*CreateOptimizationRules{}
	for _, templateRule := range templateRules {
		var ruleCode string
		for _, opRule := range allRules {
			if templateRule.Name == opRule.Rule.Name {
				ruleCode = opRule.RuleCode
				break
			}
		}
		// 构建SQL优化任务使用的规则，当前所有重写规则的阈值都为模板规则参数的first_key
		rule := &CreateOptimizationRules{
			RuleCode:  ruleCode,
			Rewrite:   true,
			Threshold: templateRule.Params.GetParam(rulepkg.DefaultSingleParamKeyName).String(),
		}
		reqRules = append(reqRules, rule)
	}
	return reqRules
}

func convertRuleNameAndMessage(ctx context.Context, ruleCode string, ruleMessage string, dbType string) (string, string) {
	lang := locale.Bundle.GetLangTagFromCtx(ctx)
	name, message := ruleCode, ruleMessage
	// 获取对应的重写规则
	rule, exist := opt.GetOptimizationRuleByRuleCode(ruleCode, dbType)
	if exist && rule != nil {
		// 用重写规则的名称和描述
		name, message = rule.Name, rule.I18nRuleInfo.GetRuleInfoByLangTag(lang).Desc
		dm := driver.GetPluginManager().GetDriverMetasOfPlugin(dbType)
		if dm != nil {
			for _, v := range dm.Rules {
				if v.Name == name {
					// 获取复用规则的Desc（保持与规则模板中的规则名一致）
					message = v.I18nRuleInfo.GetRuleInfoByLangTag(lang).Desc
					break
				}
			}
		}
	}
	return name, message
}
