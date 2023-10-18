package dms

import (
	"context"
	"fmt"
	"strconv"
	"time"

	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	dmsCommonAes "github.com/actiontech/dms/pkg/dms-common/pkg/aes"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/jinzhu/gorm"
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

	var ruleTemplates []model.RuleTemplate
	if instance.SQLEConfig.RuleTemplateID != "" {
		ruleTemplates = []model.RuleTemplate{{ProjectId: model.ProjectUID(instance.ProjectUID), Name: instance.SQLEConfig.RuleTemplateName}}
	}

	return &model.Instance{
		ID:                uint64(instanceId),
		Name:              instance.Name,
		DbType:            string(instance.DBType),
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
		RuleTemplates:     ruleTemplates,
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

func GetInstancesByIds(ctx context.Context, instanceIds []uint64) ([]*model.Instance, error) {
	if len(instanceIds) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

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

func GetInstanceInProjectById(ctx context.Context, projectUid string, instanceId uint64) (*model.Instance, bool, error) {
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
