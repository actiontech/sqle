package v2

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	dmsCommonV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	dmsCommonV2 "github.com/actiontech/dms/pkg/dms-common/api/dms/v2"
	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/actiontech/sqle/sqle/server/auditplan"
	opt "github.com/actiontech/sqle/sqle/server/optimization/rule"
	"github.com/labstack/echo/v4"
)

type InstanceTipReqV2 struct {
	FilterDBType             string `json:"filter_db_type" query:"filter_db_type"`
	FilterByEnvironmentTag   string `json:"filter_by_environment_tag" query:"filter_by_environment_tag"`
	FilterWorkflowTemplateId uint32 `json:"filter_workflow_template_id" query:"filter_workflow_template_id"`
	FunctionalModule         string `json:"functional_module" query:"functional_module" enums:"create_audit_plan,create_workflow,sql_manage,create_optimization,create_pipeline,create_version,view_sql_insight" valid:"omitempty,oneof=create_audit_plan create_workflow sql_manage create_optimization create_pipeline create_version view_sql_insight"`
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
	EnvironmentTagName      string   `json:"environment_tag_name"`
	EnvironmentTagUID       string   `json:"environment_tag_uid"`
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
// @Param functional_module query string false "functional module" Enums(create_audit_plan,create_workflow,sql_manage,create_optimization,create_pipeline,create_version,view_sql_insight)
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
	case v1.FunctionalModuleCreateVersion:
		operationType = dmsCommonV1.OpPermissionVersionManage
	case v1.FunctionalModuleViewSQLInsight:
		operationType = dmsCommonV1.OpPermissionTypeViewSQLInsight
	default:
	}
	dbServiceReq := &dmsCommonV2.ListDBServiceReq{
		FilterByEnvironmentTagUID: req.FilterByEnvironmentTag,
		ProjectUid:                projectUid,
		FilterByDBType:            req.FilterDBType,
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
			EnvironmentTagName:      inst.EnvironmentTagName,
			EnvironmentTagUID:       inst.EnvironmentTagUID,
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
	instances, err := dms.GetInstancesInProjectByTypeAndEnvironmentTag(ctx, req.ProjectUid, req.FilterByDBType, req.FilterByEnvironmentTagUID)
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

type GetInstanceAuditPlansReq struct {
	FilterByEnvironmentTag string `json:"filter_by_environment_tag" query:"filter_by_environment_tag"`
	FilterByDBType         string `json:"filter_by_db_type" query:"filter_by_db_type"`
	FilterByInstanceID     string `json:"filter_by_instance_id" query:"filter_by_instance_id"`
	FilterByAuditPlanType  string `json:"filter_by_audit_plan_type" query:"filter_by_audit_plan_type"`
	FilterByActiveStatus   string `json:"filter_by_active_status" query:"filter_by_active_status"`
	FuzzySearch            string `json:"fuzzy_search" query:"fuzzy_search"`

	PageIndex uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetInstanceAuditPlansRes struct {
	controller.BaseRes
	Data      []InstanceAuditPlanResV1 `json:"data"`
	TotalNums uint64                   `json:"total_nums"`
}

type InstanceAuditPlanResV1 struct {
	InstanceAuditPlanId uint                   `json:"instance_audit_plan_id"`
	InstanceID          string                 `json:"instance_id"`
	InstanceName        string                 `json:"instance_name"`
	Environment         string                 `json:"environment"`
	InstanceType        string                 `json:"instance_type"`
	AuditPlanTypes      []AuditPlanTypeResBase `json:"audit_plan_types"`
	ActiveStatus        string                 `json:"active_status" enums:"normal,disabled"`
	// TODO 采集状态
	CreateTime string `json:"create_time"`
	Creator    string `json:"creator"`
}

type AuditPlanTypeResBase struct {
	AuditPlanId          uint   `json:"audit_plan_id"`
	AuditPlanType        string `json:"type"`
	AuditPlanTypeDesc    string `json:"desc"`
	Token                string `json:"token"`
	ActiveStatus         string `json:"active_status" enums:"normal,disabled"`
	LastCollectionStatus string `json:"last_collection_status" enums:"normal,abnormal"`
}

// GetInstanceAuditPlans
// @Summary 获取实例扫描任务列表
// @Description get instance audit plan info list
// @Id getInstanceAuditPlansV2
// @Tags instance_audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param filter_by_environment_tag query string false "filter by environment tag"
// @Param filter_by_db_type query string false "filter by db type"
// @Param filter_by_instance_id query string false "filter by instance id"
// @Param filter_by_audit_plan_type query string false "filter instance audit plan type"
// @Param filter_by_active_status query string false "filter instance audit plan active status"
// @Param fuzzy_search query string false "fuzzy search"
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Success 200 {object} v2.GetInstanceAuditPlansRes
// @router /v2/projects/{project_name}/instance_audit_plans [get]
func GetInstanceAuditPlans(c echo.Context) error {
	s := model.GetStorage()

	req := new(GetInstanceAuditPlansReq)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	projectUid, err := dms.GetProjectUIDByName(c.Request().Context(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	limit, offset := controller.GetLimitAndOffset(req.PageIndex, req.PageSize)

	userId := controller.GetUserID(c)

	up, err := dms.NewUserPermission(userId, projectUid)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	data := map[string]interface{}{
		"filter_instance_audit_plan_db_type": req.FilterByDBType,
		"filter_audit_plan_type":             req.FilterByAuditPlanType,
		"filter_audit_plan_instance_id":      req.FilterByInstanceID,
		"filter_project_id":                  projectUid,
		"current_user_id":                    userId,
		"current_user_is_admin":              up.IsAdmin(),
		"filter_by_active_status":            req.FilterByActiveStatus,
		"limit":                              limit,
		"offset":                             offset,
	}
	if !up.CanViewProject() {
		// 如果有配置SQL管控权限，那么可以查看自己创建的或者该权限对应数据源的
		accessibleInstanceId := up.GetInstancesByOP(dmsCommonV1.OpPermissionTypeSaveAuditPlan)
		if len(accessibleInstanceId) > 0 {
			data["accessible_instances_id"] = fmt.Sprintf("\"%s\"", strings.Join(accessibleInstanceId, "\",\""))
		}
	}

	if req.FilterByEnvironmentTag != "" {
		instances, err := dms.GetInstancesInProjectByTypeAndEnvironmentTag(c.Request().Context(), projectUid, "", req.FilterByEnvironmentTag)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		instIds := make([]string, len(instances))
		for i, v := range instances {
			instIds[i] = v.GetIDStr()
		}

		data["filter_instance_ids_by_env"] = fmt.Sprintf("\"%s\"", strings.Join(instIds, "\",\""))
	}

	instanceAuditPlans, count, err := s.GetInstanceAuditPlansByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	resData := make([]InstanceAuditPlanResV1, len(instanceAuditPlans))
	for i, v := range instanceAuditPlans {
		auditPlanIds := strings.Split(v.AuditPlanIds.String, ",")
		typeBases := make([]AuditPlanTypeResBase, 0, len(auditPlanIds))
		for _, auditPlanId := range auditPlanIds {
			if auditPlanId != "" {
				typeBase, err := ConvertAuditPlanTypeToResByID(c.Request().Context(), auditPlanId, v.Token)
				if err != nil {
					return controller.JSONBaseErrorReq(c, err)
				}
				typeBases = append(typeBases, typeBase)

			}
		}
		inst := dms.GetInstancesByIdWithoutError(v.InstanceID)
		resData[i] = InstanceAuditPlanResV1{
			InstanceAuditPlanId: v.Id,
			InstanceID:          strconv.FormatUint(inst.ID, 10),
			InstanceName:        inst.Name,
			Environment:         inst.EnvironmentTagName,
			InstanceType:        v.DBType,
			AuditPlanTypes:      typeBases,
			ActiveStatus:        v.ActiveStatus,
			CreateTime:          v.CreateTime,
			Creator:             dms.GetUserNameWithDelTag(v.CreateUserId),
		}
	}
	return c.JSON(http.StatusOK, &GetInstanceAuditPlansRes{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      resData,
		TotalNums: count,
	})
}

func ConvertAuditPlanTypeToResByID(ctx context.Context, id string, token string) (AuditPlanTypeResBase, error) {
	auditPlanID, err := strconv.Atoi(id)
	if err != nil {
		return AuditPlanTypeResBase{}, err
	}
	s := model.GetStorage()
	auditPlan, exist, err := s.GetAuditPlanByID(auditPlanID)
	if err != nil {
		return AuditPlanTypeResBase{}, err
	}
	if !exist {
		return AuditPlanTypeResBase{}, nil
	}
	for _, meta := range auditplan.Metas {
		if meta.Type == auditPlan.Type {
			return AuditPlanTypeResBase{
				AuditPlanType:        auditPlan.Type,
				AuditPlanTypeDesc:    locale.Bundle.LocalizeMsgByCtx(ctx, meta.Desc),
				AuditPlanId:          auditPlan.ID,
				Token:                token,
				ActiveStatus:         auditPlan.ActiveStatus,
				LastCollectionStatus: auditPlan.AuditPlanTaskInfo.LastCollectionStatus,
			}, nil
		}
	}
	return AuditPlanTypeResBase{}, nil
}

type GetInstanceAuditPlanDetailRes struct {
	controller.BaseRes
	Data InstanceAuditPlanDetailRes `json:"data"`
}

type InstanceAuditPlanDetailRes struct {
	Environment  string `json:"environment" example:"prod"`
	InstanceType string `json:"instance_type" example:"mysql" `
	InstanceName string `json:"instance_name" example:"test_mysql"`
	InstanceID   string `json:"instance_id" example:"instance_id"`
	// 扫描类型
	AuditPlans []AuditPlanRes `json:"audit_plans"`
}

type AuditPlanRes struct {
	RuleTemplateName        string                          `json:"rule_template_name" from:"rule_template_name" example:"default_MySQL"`
	Type                    AuditPlanTypeResBase            `json:"audit_plan_type" form:"audit_plan_type"`
	Params                  []v1.AuditPlanParamResV1        `json:"audit_plan_params" valid:"dive,required"`
	NeedMarkHighPrioritySQL bool                            `json:"need_mark_high_priority_sql"`
	HighPriorityConditions  []v1.HighPriorityConditionResV1 `json:"high_priority_conditions"`
}

type HighPriorityCondition struct {
	Key         string              `json:"key"`
	Desc        string              `json:"desc"`
	Value       string              `json:"value"`
	Type        string              `json:"type" enums:"string,int,bool,password"`
	EnumsValues []params.EnumsValue `json:"enums_value"`
	Operator    v1.Operator         `json:"operator"`
}

// @Summary 获取实例扫描任务详情
// @Description get instance audit plan detail
// @Id getInstanceAuditPlanDetailV2
// @Tags instance_audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instance_audit_plan_id path string true "instance audit plan id"
// @Success 200 {object} v2.GetInstanceAuditPlanDetailRes
// @router /v2/projects/{project_name}/instance_audit_plans/{instance_audit_plan_id} [get]
func GetInstanceAuditPlanDetail(c echo.Context) error {
	instanceAuditPlanID := c.Param("instance_audit_plan_id")
	projectUID, err := dms.GetProjectUIDByName(c.Request().Context(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// check current user instance audit plan permission
	detail, exist, err := v1.GetInstanceAuditPlanIfCurrentUserCanView(c, projectUID, instanceAuditPlanID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewInstanceAuditPlanNotExistErr())
	}
	auditPlans, err := ConvertAuditPlansToRes(c.Request().Context(), detail.AuditPlans)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	inst := dms.GetInstancesByIdWithoutError(fmt.Sprintf("%d", detail.InstanceID))
	resData := InstanceAuditPlanDetailRes{
		InstanceType: detail.DBType,
		InstanceName: inst.Name,
		Environment:  inst.EnvironmentTagUID,
		InstanceID:   inst.GetIDStr(),
		AuditPlans:   auditPlans,
	}
	return c.JSON(http.StatusOK, &GetInstanceAuditPlanDetailRes{
		BaseRes: controller.NewBaseReq(nil),
		Data:    resData,
	})
}

func ConvertAuditPlanTypeToRes(ctx context.Context, id uint, auditPlanType string) AuditPlanTypeResBase {
	for _, meta := range auditplan.Metas {
		if meta.Type == auditPlanType {
			return AuditPlanTypeResBase{
				AuditPlanType:     auditPlanType,
				AuditPlanTypeDesc: locale.Bundle.LocalizeMsgByCtx(ctx, meta.Desc),
				AuditPlanId:       id,
			}
		}
	}
	return AuditPlanTypeResBase{}
}

func ConvertAuditPlansToRes(ctx context.Context, auditPlans []*model.AuditPlanV2) ([]AuditPlanRes, error) {
	resAuditPlans := make([]AuditPlanRes, 0, len(auditPlans))
	for _, v := range auditPlans {
		typeBase := ConvertAuditPlanTypeToRes(ctx, v.ID, v.Type)
		resAuditPlan := AuditPlanRes{
			RuleTemplateName: v.RuleTemplateName,
			Type:             typeBase,
		}
		meta, err := auditplan.GetMeta(v.Type)
		if err != nil {
			return nil, err
		}
		meta.Params = func(instanceId ...string) params.Params { return v.Params }
		if meta.Params != nil && len(meta.Params()) > 0 {
			paramsRes := make([]v1.AuditPlanParamResV1, 0, len(meta.Params()))
			for _, p := range meta.Params() {
				val := p.Value
				if p.Type == params.ParamTypePassword {
					val = ""
				}
				paramRes := v1.AuditPlanParamResV1{
					Key:   p.Key,
					Desc:  p.GetDesc(locale.Bundle.GetLangTagFromCtx(ctx)),
					Type:  string(p.Type),
					Value: val,
				}
				paramsRes = append(paramsRes, paramRes)
			}
			resAuditPlan.Params = paramsRes
		}

		if v.HighPriorityParams != nil && len(v.HighPriorityParams) > 0 {
			hppParamsRes := make([]v1.HighPriorityConditionResV1, len(v.HighPriorityParams))
			for i, hpp := range v.HighPriorityParams {
				for _, metaHpp := range meta.HighPriorityParams {
					if metaHpp.Key != hpp.Key {
						continue
					}
					highParamRes := v1.HighPriorityConditionResV1{
						Key:   metaHpp.Key,
						Desc:  metaHpp.GetDesc(locale.Bundle.GetLangTagFromCtx(ctx)),
						Value: hpp.Value,
						Type:  string(metaHpp.Type),
						Operator: v1.OperatorResV1{
							Value:      string(hpp.Operator.Value),
							EnumsValue: v1.ConvertEnumsValuesToRes(ctx, metaHpp.Operator.EnumsValue),
						},
					}
					hppParamsRes[i] = highParamRes
					break
				}
			}
			resAuditPlan.HighPriorityConditions = hppParamsRes
			resAuditPlan.NeedMarkHighPrioritySQL = v.NeedMarkHighPrioritySQL
		}

		resAuditPlans = append(resAuditPlans, resAuditPlan)
	}
	return resAuditPlans, nil
}
