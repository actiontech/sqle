package dms

import (
	"context"
	"fmt"
	"strconv"
	"time"

	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	dmsCommonAes "github.com/actiontech/dms/pkg/dms-common/pkg/aes"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
)

func getInstances(ctx context.Context, req dmsV1.ListDBServiceReq) ([]*model.Instance, error) {
	var ret = make([]*model.Instance, 0)

	var limit, pageIndex uint32 = 20, 1

	for ; ; pageIndex++ {
		req.PageIndex = pageIndex
		req.PageSize = limit

		dbServices, _, err := func() ([]*dmsV1.ListDBService, int64, error) {
			newCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			return dmsobject.ListDbServices(newCtx, GetDMSServerAddress(), req)
		}()

		if err != nil {
			return nil, fmt.Errorf("get instances from dms error: %v", err)
		}

		for _, item := range dbServices {
			if item.SQLEConfig == nil || item.SQLEConfig.RuleTemplateID == "" {
				continue
			}

			instance, err := convertInstance(item)
			if err != nil {
				return nil, fmt.Errorf("convert instance error: %v", err)
			}

			ret = append(ret, instance)
		}

		if len(dbServices) < int(limit) {
			break
		}
	}

	return ret, nil
}

func getInstance(ctx context.Context, req dmsV1.ListDBServiceReq) (*model.Instance, bool, error) {
	newCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	dbServices, total, err := dmsobject.ListDbServices(newCtx, GetDMSServerAddress(), req)

	if err != nil {
		return nil, false, fmt.Errorf("get instances from dms error: %v", err)
	}

	if total == 0 {
		return nil, false, nil
	}

	instance, err := convertInstance(dbServices[0])
	if err != nil {
		return nil, false, fmt.Errorf("convert instance error: %v", err)
	}

	return instance, true, nil
}

func convertInstance(instance *dmsV1.ListDBService) (*model.Instance, error) {
	instanceId, err := strconv.ParseInt(instance.DBServiceUid, 0, 64)
	if err != nil {
		return nil, err
	}

	ruleTemplateId, err := strconv.ParseInt(instance.SQLEConfig.RuleTemplateID, 0, 64)
	if err != nil {
		return nil, err
	}

	var maintenancePeriod = make(model.Periods, 0)
	for _, item := range instance.MaintenanceTimes {
		maintenancePeriod = append(maintenancePeriod, &model.Period{
			StartHour:   item.MaintenanceStartTime.Hour,
			StartMinute: item.MaintenanceStartTime.Minute,
			EndHour:     item.MaintenanceStopTime.Hour,
			EndMinute:   item.MaintenanceStopTime.Minute,
		})
	}

	decryptPassword, err := dmsCommonAes.AesDecrypt(instance.Password)
	if err != nil {
		return nil, err
	}

	additionalParams := make(params.Params, 0, len(instance.AdditionalParams))
	for _, item := range instance.AdditionalParams {
		additionalParams = append(additionalParams, &params.Param{
			Key:   item.Name,
			Value: item.Value,
			Desc:  item.Description,
			Type:  params.ParamType(item.Type),
		})
	}

	sqlQueryConfig := model.SqlQueryConfig{}
	if instance.SQLEConfig != nil {
		sqlQueryConfig = model.SqlQueryConfig{
			MaxPreQueryRows:                  instance.SQLEConfig.SQLQueryConfig.MaxPreQueryRows,
			QueryTimeoutSecond:               instance.SQLEConfig.SQLQueryConfig.QueryTimeoutSecond,
			AuditEnabled:                     instance.SQLEConfig.SQLQueryConfig.AuditEnabled,
			AllowQueryWhenLessThanAuditLevel: string(instance.SQLEConfig.SQLQueryConfig.AllowQueryWhenLessThanAuditLevel),
		}
	}

	return &model.Instance{
		ID:                uint64(instanceId),
		Name:              instance.Name,
		DbType:            instance.DBType,
		RuleTemplateId:    uint64(ruleTemplateId),
		RuleTemplateName:  instance.SQLEConfig.RuleTemplateName,
		ProjectId:         instance.ProjectUID,
		MaintenancePeriod: maintenancePeriod,
		Host:              instance.Host,
		Port:              instance.Port,
		User:              instance.User,
		Password:          decryptPassword,
		Desc:              instance.Desc,
		AdditionalParams:  additionalParams,
		SqlQueryConfig:    sqlQueryConfig,
	}, nil
}

func GetInstancesInProject(ctx context.Context, projectUid string) ([]*model.Instance, error) {
	return getInstances(ctx, dmsV1.ListDBServiceReq{
		ProjectUid: projectUid,
	})
}

func GetInstancesInProjectByType(ctx context.Context, projectUid, dbType string) ([]*model.Instance, error) {
	return getInstances(ctx, dmsV1.ListDBServiceReq{
		ProjectUid:     projectUid,
		FilterByDBType: dbType,
	})
}

func GetInstancesNameInProjectByRuleTemplateName(ctx context.Context, projectUid, ruleTemplateName string) ([]string, error) {
	instances, err := getInstances(ctx, dmsV1.ListDBServiceReq{
		ProjectUid: projectUid,
	})

	if err != nil {
		return nil, err
	}

	ret := make([]string, 0)
	for _, instance := range instances {
		if instance.RuleTemplateName == ruleTemplateName {
			ret = append(ret, instance.Name)
		}
	}

	return ret, nil
}

func GetInstancesNameByRuleTemplateName(ctx context.Context, ruleTemplateName string) ([]string, error) {
	instances, err := getInstances(ctx, dmsV1.ListDBServiceReq{})

	if err != nil {
		return nil, err
	}

	ret := make([]string, 0)
	for _, instance := range instances {
		if instance.RuleTemplateName == ruleTemplateName {
			ret = append(ret, instance.Name)
		}
	}

	return ret, nil
}

func GetInstanceInProjectByName(ctx context.Context, projectUid, name string) (*model.Instance, bool, error) {
	if len(projectUid) == 0 || len(name) == 0 {
		return nil, false, nil
	}

	return getInstance(ctx, dmsV1.ListDBServiceReq{
		PageSize:     1,
		FilterByName: name,
		ProjectUid:   projectUid,
	})
}

func GetInstancesInProjectByNames(ctx context.Context, projectUid string, names []string) (instances []*model.Instance, err error) {
	for _, name := range names {
		instance, isExist, err := getInstance(ctx, dmsV1.ListDBServiceReq{
			PageSize:     1,
			FilterByName: name,
			ProjectUid:   projectUid,
		})

		if err != nil {
			return nil, err
		}

		if isExist {
			instances = append(instances, instance)
		}
	}

	return instances, err
}

func GetInstanceNamesInProjectByIds(ctx context.Context, projectUid string, instanceIds []string) ([]string, error) {
	ret := make([]string, 0)
	for _, instanceId := range instanceIds {
		instance, exist, err := getInstance(ctx, dmsV1.ListDBServiceReq{
			PageSize:    1,
			FilterByUID: instanceId,
			ProjectUid:  projectUid,
		})

		if err != nil {
			return nil, err
		}

		if exist {
			ret = append(ret, instance.Name)
		}
	}

	return ret, nil
}

func GetInstanceNamesInProject(ctx context.Context, projectUid string) ([]string, error) {
	ret := make([]string, 0)

	instances, err := getInstances(ctx, dmsV1.ListDBServiceReq{
		PageSize:   1,
		ProjectUid: projectUid,
	})

	if err != nil {
		return nil, err
	}

	for _, instance := range instances {
		ret = append(ret, instance.Name)
	}

	return ret, nil
}

func GetInstancesById(ctx context.Context, instanceId uint64) (*model.Instance, bool, error) {
	if instanceId == 0 {
		return nil, false, nil
	}

	return getInstance(ctx, dmsV1.ListDBServiceReq{
		PageSize:    1,
		FilterByUID: strconv.FormatUint(instanceId, 10),
	})
}

func GetInstancesByIds(ctx context.Context, instanceIds []uint64) ([]*model.Instance, error) {
	ret := make([]*model.Instance, 0)
	for _, instanceId := range instanceIds {
		instance, exist, err := getInstance(ctx, dmsV1.ListDBServiceReq{
			PageSize:    1,
			FilterByUID: strconv.FormatUint(instanceId, 10),
		})

		if err != nil {
			return nil, err
		}

		if exist {
			ret = append(ret, instance)
		}
	}

	return ret, nil
}

func GetInstanceIdNameMapByIds(ctx context.Context, instanceIds []uint64) (map[uint64]string, error) {
	// todo: remove duplicate instance id
	ret := make(map[uint64]string)
	for _, instanceId := range instanceIds {
		instance, exist, err := getInstance(ctx, dmsV1.ListDBServiceReq{
			PageSize:    1,
			FilterByUID: strconv.FormatUint(instanceId, 10),
		})

		if err != nil {
			return nil, err
		}

		if exist {
			ret[instance.ID] = instance.Name
		}
	}

	return ret, nil
}

func GetInstanceInProjectById(ctx context.Context, projectUid string, instanceId uint64) (*model.Instance, bool, error) {
	if len(projectUid) == 0 || instanceId == 0 {
		return nil, false, nil
	}

	return getInstance(ctx, dmsV1.ListDBServiceReq{
		PageSize:    1,
		FilterByUID: strconv.FormatUint(instanceId, 10),
		ProjectUid:  projectUid,
	})
}

func GetInstancesInProjectByIds(ctx context.Context, projectUid string, instanceIds []uint64) ([]*model.Instance, error) {
	ret := make([]*model.Instance, 0)
	for _, instanceId := range instanceIds {
		instance, exist, err := getInstance(ctx, dmsV1.ListDBServiceReq{
			PageSize:    1,
			FilterByUID: strconv.FormatUint(instanceId, 10),
			ProjectUid:  projectUid,
		})

		if err != nil {
			return nil, err
		}

		if exist {
			ret = append(ret, instance)
		}
	}

	return ret, nil
}

type InstanceTypeCount struct {
	DBType string `json:"db_type"`
	Count  int64  `json:"count"`
}

func GetInstanceCountGroupType(ctx context.Context) ([]InstanceTypeCount, error) {
	instances, err := getInstances(ctx, dmsV1.ListDBServiceReq{})

	if err != nil {
		return nil, err
	}

	var typeCountMap = map[string]int64{}

	for _, instance := range instances {
		typeCountMap[instance.DbType]++
	}

	ret := make([]InstanceTypeCount, 0, len(typeCountMap))
	for dbType, count := range typeCountMap {
		ret = append(ret, InstanceTypeCount{
			DBType: dbType,
			Count:  count,
		})
	}

	return ret, nil
}

func GetWorkflowDetailByWorkflowId(projectId, workflowId string, fn func(projectId, workflowId string) (*model.Workflow, bool, error)) (*model.Workflow, error) {
	workflow, exist, err := fn(projectId, workflowId)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New(errors.DataNotExist, fmt.Errorf("workflow is not exist or you can't access it"))
	}

	instanceIds := make([]uint64, 0, len(workflow.Record.InstanceRecords))
	for _, item := range workflow.Record.InstanceRecords {
		instanceIds = append(instanceIds, item.InstanceId)
	}

	if len(instanceIds) == 0 {
		return workflow, nil
	}

	instances, err := GetInstancesInProjectByIds(context.Background(), string(workflow.ProjectId), instanceIds)
	if err != nil {
		return nil, err
	}
	instanceMap := map[uint64]*model.Instance{}
	for _, instance := range instances {
		instanceMap[instance.ID] = instance
	}
	for i, item := range workflow.Record.InstanceRecords {
		if instance, ok := instanceMap[item.InstanceId]; ok {
			workflow.Record.InstanceRecords[i].Instance = instance
		}
	}

	return workflow, nil
}

func GetAuditPlanWithInstanceFromProjectByName(projectId, name string, fn func(projectId, name string) (*model.AuditPlan, bool, error)) (*model.AuditPlan, bool, error) {
	auditPlan, exist, err := fn(projectId, name)
	if err != nil {
		return nil, false, err
	}
	if !exist {
		return nil, false, nil
	}

	instance, exists, err := GetInstanceInProjectByName(context.Background(), projectId, auditPlan.InstanceName)
	if err != nil {
		return nil, false, err
	}
	if exists {
		auditPlan.Instance = instance
	}
	return auditPlan, true, nil
}

func GetActiveAuditPlansWithInstance(fn func() ([]*model.AuditPlan, error)) ([]*model.AuditPlan, error) {
	auditPlans, err := fn()
	if err != nil {
		return nil, err
	}

	for i, item := range auditPlans {
		// todo dms不支持跨项目查询实例，所以单个查询
		instance, exists, err := GetInstanceInProjectByName(context.Background(), string(item.ProjectId), item.Name)
		if err != nil {
			continue
		}
		if exists {
			auditPlans[i].Instance = instance
		}
	}
	return auditPlans, nil
}

func GetAuditPlanWithInstanceById(id uint, fn func(id uint) (*model.AuditPlan, bool, error)) (*model.AuditPlan, bool, error) {
	auditPlan, exist, err := fn(id)
	if err != nil {
		return nil, false, err
	}
	if !exist {
		return nil, false, nil
	}

	instance, exists, err := GetInstanceInProjectByName(context.Background(), string(auditPlan.ProjectId), auditPlan.InstanceName)
	if err != nil {
		return nil, false, err
	}
	if exists {
		auditPlan.Instance = instance
	}
	return auditPlan, true, nil
}
