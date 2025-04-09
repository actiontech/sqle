package v2

import (
	"context"
	"fmt"
	"net/http"

	dmsCommonV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	dmsCommonV2 "github.com/actiontech/dms/pkg/dms-common/api/dms/v2"
	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server"
	opt "github.com/actiontech/sqle/sqle/server/optimization/rule"
	"github.com/labstack/echo/v4"
)

type InstanceTipReqV2 struct {
	FilterDBType             string `json:"filter_db_type" query:"filter_db_type"`
	FilterByEnvironmentTag   string `json:"filter_by_environment_tag" query:"filter_by_environment_tag"`
	FilterWorkflowTemplateId uint32 `json:"filter_workflow_template_id" query:"filter_workflow_template_id"`
	FunctionalModule         string `json:"functional_module" query:"functional_module" enums:"create_audit_plan,create_workflow,sql_manage,create_optimization,create_pipeline" valid:"omitempty,oneof=create_audit_plan create_workflow sql_manage create_optimization create_pipeline"`
}


type GetInstanceTipsResV2 struct {
	controller.BaseRes
	Data []InstanceTipResV2 `json:"data"`
}

type InstanceTipResV2 struct {
	ID                      string   `json:"instance_id"`
	Name                    string   `json:"instance_name"`
	Type                    string   `json:"instance_type"`
	WorkflowTemplateId      uint32   `json:"workflow_template_id"`
	Host                    string   `json:"host"`
	Port                    string   `json:"port"`
	EnableBackup            bool     `json:"enable_backup"`
	SupportedBackupStrategy []string `json:"supported_backup_strategy"  enums:"none,manual,reverse_sql,original_row"`
	BackupMaxRows           uint64   `json:"backup_max_rows"`
}

// GetInstanceTips get instance tip list
// @Summary 获取实例提示列表
// @Description get instance tip list
// @Tags instance
// @Id getInstanceTipListV2
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param filter_db_type query string false "filter db type"
// @Param filter_by_environment_tag query string false "filter by environment tag"
// @Param filter_workflow_template_id query string false "filter workflow template id"
// @Param functional_module query string false "functional module" Enums(create_audit_plan,create_workflow,sql_manage,create_optimization,create_pipeline)
// @Success 200 {object} v2.GetInstanceTipsResV2
// @router /v2/projects/{project_name}/instance_tips [get]
func GetInstanceTips(c echo.Context) error {
	req := new(InstanceTipReqV2)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	projectUid, err := dms.GetProjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var operationType dmsCommonV1.OpPermissionType
	switch req.FunctionalModule {
	case v1.FunctionalModuleCreateAuditPlan:
		operationType = dmsCommonV1.OpPermissionTypeSaveAuditPlan
	case v1.FunctionalModuleCreateWorkflow:
		operationType = dmsCommonV1.OpPermissionTypeCreateWorkflow
	case v1.FunctionalModuleCreateOptimization:
		operationType = dmsCommonV1.OpPermissionTypeCreateOptimization
	case v1.FunctionalModuleCreatePipeline:
		operationType = dmsCommonV1.OpPermissionTypeCreatePipeline
	default:
	}
	dbServiceReq := &dmsCommonV2.ListDBServiceReq{
		FilterByEnvironmentTag: req.FilterByEnvironmentTag,
		ProjectUid:             projectUid,
		FilterByDBType:         req.FilterDBType,
	}

	instances, err := GetCanOperationInstances(c.Request().Context(), user, dbServiceReq, operationType)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	s := model.GetStorage()
	template, exist, err := s.GetWorkflowTemplateByProjectId(model.ProjectUID(projectUid))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("current project doesn't has workflow template"))
	}

	instanceTipsResV1 := make([]InstanceTipResV2, 0, len(instances))
	svc := server.BackupService{}
	for _, inst := range instances {
		if operationType == dmsCommonV1.OpPermissionTypeCreateOptimization && !opt.CanOptimizeDbType(inst.DbType) {
			continue
		}
		instanceTipRes := InstanceTipResV2{
			ID:                      inst.GetIDStr(),
			Name:                    inst.Name,
			Type:                    inst.DbType,
			Host:                    inst.Host,
			Port:                    inst.Port,
			WorkflowTemplateId:      uint32(template.ID),
			EnableBackup:            inst.EnableBackup,
			BackupMaxRows:           inst.BackupMaxRows,
			SupportedBackupStrategy: svc.SupportedBackupStrategy(inst.DbType),
		}
		instanceTipsResV1 = append(instanceTipsResV1, instanceTipRes)
	}
	return c.JSON(http.StatusOK, &GetInstanceTipsResV2{
		BaseRes: controller.NewBaseReq(nil),
		Data:    instanceTipsResV1,
	})
}

// 根据用户权限获取能访问/操作的实例列表
func GetCanOperationInstances(ctx context.Context, user *model.User, req *dmsCommonV2.ListDBServiceReq, operationType dmsCommonV1.OpPermissionType) ([]*model.Instance, error) {
	// 获取当前项目下指定数据库类型的全部实例
	instances, err := dms.GetInstancesInProjectByTypeAndBusiness(ctx, req.ProjectUid, req.FilterByDBType, req.FilterByEnvironmentTag)
	if err != nil {
		return nil, err
	}

	userOpPermissions, isAdmin, err := dmsobject.GetUserOpPermission(ctx, req.ProjectUid, user.GetIDStr(), controller.GetDMSServerAddress())
	if err != nil {
		return nil, err
	}

	if isAdmin || operationType == "" {
		return instances, nil
	}
	canOperationInstance := make([]*model.Instance, 0)
	for _, instance := range instances {
		if v1.CanOperationInstance(userOpPermissions, []dmsCommonV1.OpPermissionType{operationType}, instance) {
			canOperationInstance = append(canOperationInstance, instance)
		}
	}
	return canOperationInstance, nil
}
