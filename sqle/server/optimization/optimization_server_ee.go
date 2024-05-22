//go:build enterprise
// +build enterprise

package optimization

import (
	"context"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"
)

type OptimizationServeror interface {
	Optimizate(ctx context.Context, optimizationSQL string) (optimizationInfo model.SQLOptimizationRecord, err error)
}

type OptimizateStatus string

// 优化状态常量
const (
	OptimizateStatusOptimizating OptimizateStatus = "optimizing"
	OptimizateStatusFinish       OptimizateStatus = "finish"
	OptimizateStatusFailed       OptimizateStatus = "failed"
)

func (dets OptimizateStatus) String() string {
	return string(dets)
}

// SQL优化任务入口
func Optimizate(ctx context.Context, user, projectId string, instance *model.Instance, schema *string, optimizationName, OptimizationSQL string) (optimizationId string, err error) {
	logger := log.NewEntry()

	id, err := utils.GenUid()
	if err != nil {
		return "", err
	}
	// 保存优化SQL记录基础信息
	optimizationRecord := new(model.SQLOptimizationRecord)
	optimizationRecord.Creator = user
	optimizationRecord.ProjectId = projectId
	optimizationRecord.InstanceId = instance.ID
	optimizationRecord.InstanceName = instance.Name
	optimizationRecord.SchemaName = *schema
	optimizationRecord.DBType = instance.DbType
	optimizationRecord.OptimizationId = id
	optimizationRecord.OptimizationName = optimizationName
	optimizationRecord.Status = OptimizateStatusOptimizating.String()

	err = model.GetStorage().Save(optimizationRecord)
	if err != nil {
		logger.Error(err)
		return
	}

	// 保存优化任务详情
	go func() {
		store := model.GetStorage()
		optimizationRecord, err := store.GetOptimizationRecordId(id)
		if err != nil {
			logger.Error(err)
			return
		}

		// SQL优化服务 newOptimizateServer
		server := NewOptimizationOnlinePawSQLServer(logger, instance, *schema)
		// 调用SQL优化
		optimizationInfo, err := server.Optimizate(context.TODO(), OptimizationSQL)
		if err != nil {
			optimizationRecord.Status = OptimizateStatusFailed.String()
			optimizationRecord.OptimizationSQLs = append(optimizationRecord.OptimizationSQLs, &model.OptimizationSQL{
				OptimizationId: id,
				OriginalSQL:    OptimizationSQL,
			})
			logger.Error(err)
		} else {
			optimizationRecord.PerformanceImprove = optimizationInfo.PerformanceImprove
			optimizationRecord.NumberOfQuery = optimizationInfo.NumberOfQuery
			optimizationRecord.NumberOfSyntaxError = optimizationInfo.NumberOfSyntaxError
			optimizationRecord.NumberOfRewrite = optimizationInfo.NumberOfRewrite
			optimizationRecord.NumberOfRewrittenQuery = optimizationInfo.NumberOfRewrittenQuery
			optimizationRecord.NumberOfIndex = optimizationInfo.NumberOfIndex
			optimizationRecord.NumberOfQueryIndex = optimizationInfo.NumberOfQueryIndex
			optimizationRecord.IndexRecommendations = optimizationInfo.IndexRecommendations
			for _, optSQL := range optimizationInfo.OptimizationSQLs {
				optSQL.OptimizationId = id
			}
			optimizationRecord.OptimizationSQLs = optimizationInfo.OptimizationSQLs
			optimizationRecord.Status = OptimizateStatusFinish.String()
		}
		err = store.Save(optimizationRecord)
		if err != nil {
			logger.Error(err)
			return
		}

	}()
	return id, nil
}
