package optimization

import (
	"context"
	"errors"
	"fmt"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/actiontech/sqle/sqle/log"
)

type OptimizationServeror interface {
	Optimizate(ctx context.Context, optimizationSQL string) (optimizationId string, err error)
	GetOptimizationInfo(ctx context.Context, optimizationId string) (optimizationInfo model.SQLOptimizationRecord, err error)
}

// SQL优化任务入口
func Optimizate(ctx context.Context, user, projectId string, instanceName *string, schema *string, optimizationName, OptimizationSQL string) (optimizationId string, err error) {
	logger := log.NewEntry()
	// 参数校验
	if instanceName == nil || schema == nil {
		return "", errors.New("online optimizate sql with nil instance is not supported")
	}
	instance, exist, err := dms.GetInstanceInProjectByName(ctx, projectId, *instanceName)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", fmt.Errorf("instance %s not exist", *instanceName)
	}

	// SQL优化服务 newOptimizateServer
	server := NewOptimizationOnlinePawSQLServer(logger, instance, *schema)

	// 调用SQL优化
	id, err := server.Optimizate(ctx, OptimizationSQL)
	if err != nil {
		return "", err
	}
	logger.Debugf("optimization sql successful,optimization id: %s", id)
	// 保存优化SQL记录
	go func() {
		ctx := context.TODO()
		optimizationRecord, err := server.GetOptimizationInfo(ctx, id)
		if err != nil {
			logger.Error(err)
			return
		}
		optimizationRecord.Creator = user
		optimizationRecord.ProjectId = projectId
		optimizationRecord.InstanceName = *instanceName
		optimizationRecord.SchemaName = *schema
		optimizationRecord.DBType = instance.DbType
		optimizationRecord.OptimizationId = id
		optimizationRecord.OptimizationName = optimizationName

		// 保存优化任务详情
		s := model.GetStorage()
		err = s.Save(optimizationRecord)
		if err != nil {
			logger.Error(err)
			return
		}
	}()
	return id, nil
}
