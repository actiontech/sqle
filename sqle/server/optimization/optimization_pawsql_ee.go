//go:build enterprise
// +build enterprise

package optimization

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/sirupsen/logrus"
)

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
	id, err := a.createOptimization(ctx, workspaceId.(string), a.Instance.DbType, OptimizationSQL, "online")
	if err != nil {
		return nil, err
	}
	return a.getOptimizationInfo(ctx, id)
}

// 获取优化详情
func (a *OptimizationOnlinePawSQLServer) getOptimizationInfo(ctx context.Context, analysisId string) (optimizationInfo *model.SQLOptimizationRecord, err error) {
	a.logger.Debugf("get Optimization info id : %v", analysisId)
	// 获取优化任务概览
	summary, err := a.getOptimizationSummary(ctx, analysisId)
	if err != nil {
		return optimizationInfo, err
	}
	a.logger.Debugf("get Optimization summary %v", summary)
	optimizationIndexInfo := trimKeyWord(summary.OptimizationIndexInfo)
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
			triggeredRule = append(triggeredRule, model.RewriteRule{
				RuleName:            v.RuleCode,
				Message:             fmt.Sprintf("%s \n %s", v.RuleNameZh, v.RuleNameEn),
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
			optimizedSQL = sqls[0]
		}
		// 索引数据移除pawsql关键字
		indexRecommendeds := trimKeyWord(detail.IndexRecommended)
		contributingIndices := trimKeyWord([]string{statementInfo.ContributingIndices})

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
	optimizationInfo.PerformanceImprove = performanceImprove
	return
}
func trimKeyWord(slices []string) (ret []string) {
	for _, s := range slices {
		ret = append(ret, strings.ReplaceAll(s, "PAWSQL", "OPTIMIZATION"))
	}
	return
}
