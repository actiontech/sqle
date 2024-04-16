package optimization

import (
	"context"
	"fmt"
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
	UPdateTime string
	DDL        string
}

func (o OptimizationWorkspace) MD5() string {
	return utils.Md5String(o.DBName + o.SchemaName + o.DDL)
}

var cacheOptimizationWorkspace sync.Map // map[string] /*MD5(dbName+schemaName+DDL)*/ string /*workspaceId*/

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
func (a *OptimizationOnlinePawSQLServer) Optimizate(ctx context.Context, OptimizationSQL string) (optimizationId string, err error) {
	a.logger.Debugf("Optimization SQL %v", OptimizationSQL)
	optimizationWorkspace := OptimizationWorkspace{
		DBName:     a.Instance.Name,
		SchemaName: a.Schema,
		DDL:        OptimizationSQL,
	}

	workspaceId, ok := cacheOptimizationWorkspace.Load(optimizationWorkspace.MD5())
	if !ok {
		// 创建pawsql的空间,内存缓存
		workspaceId, err = a.CreateWorkspaceOnline(ctx, a.Instance, a.Schema)
		if err != nil {
			return "", err
		}
		cacheOptimizationWorkspace.Store(optimizationWorkspace.MD5(), workspaceId)
	}

	// 创建分析任务
	return a.CreateOptimization(ctx, workspaceId.(string), a.Instance.DbType, OptimizationSQL, "online")
}

// 获取优化详情
func (a *OptimizationOnlinePawSQLServer) GetOptimizationInfo(ctx context.Context, optimizationId string) (optimizationInfo *model.SQLOptimizationRecord, err error) {
	a.logger.Debugf("get Optimization detail %v", optimizationId)
	// 获取优化任务概览
	summary, err := a.GetOptimizationSummary(ctx, optimizationId)
	if err != nil {
		return optimizationInfo, err
	}
	a.logger.Debugf("get Optimization summary %v", summary)
	optimizationInfo = &model.SQLOptimizationRecord{
		PerformanceImprove:     summary.BasicSummary.PerformanceImprove,
		NumberOfQuery:          summary.BasicSummary.NumberOfQuery,
		NumberOfSyntaxError:    summary.BasicSummary.NumberOfSyntaxError,
		NumberOfRewrite:        summary.BasicSummary.NumberOfRewrite,
		NumberOfRewrittenQuery: summary.BasicSummary.NumberOfRewrittenQuery,
		NumberOfIndex:          summary.BasicSummary.NumberOfIndex,
		NumberOfQueryIndex:     summary.BasicSummary.NumberOfQueryIndex,
		IndexRecommendations:   summary.OptimizationIndexInfo,
	}
	for _, statementInfo := range summary.SummaryStatementInfo {

		// 获取优化任务详情
		detail, err := a.GetOptimizationDetail(ctx, statementInfo.OptimizationStmtId)
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

		optimizationInfo.OptimizationSQLs = append(optimizationInfo.OptimizationSQLs, &model.OptimizationSQL{
			OptimizationId:           optimizationId,
			OriginalSQL:              statementInfo.StmtText,
			OptimizedSQL:             detail.StmtText,
			NumberOfRewrite:          statementInfo.NumberOfRewrite,
			NumberOfSyntaxError:      statementInfo.NumberOfSyntaxError,
			NumberOfIndex:            statementInfo.NumberOfIndex,
			NumberOfHitIndex:         statementInfo.NumberOfHitIndex,
			Performance:              statementInfo.Performance,
			ContributingIndices:      statementInfo.ContributingIndices,
			TriggeredRules:           triggeredRule,
			IndexRecommendations:     detail.IndexRecommended,
			ExplainValidationDetails: model.ExplainValidationDetail{BeforeCost: detail.ValidationDetails.BeforeCost, AfterCost: detail.ValidationDetails.AfterCost, BeforePlan: detail.ValidationDetails.BeforePlan, AfterPlan: detail.ValidationDetails.AfterPlan, PerformImprovePer: detail.ValidationDetails.PerformImprovePer},
		})

	}
	return
}
