package v1

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	v1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	dmsCommonJwt "github.com/actiontech/dms/pkg/dms-common/api/jwt"
	"github.com/actiontech/sqle/sqle/api/controller"
	scannerCmd "github.com/actiontech/sqle/sqle/cmd/scannerd/command"
	"github.com/actiontech/sqle/sqle/config"
	dms "github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/server/auditplan"

	"github.com/actiontech/sqle/sqle/utils"
	"github.com/labstack/echo/v4"
	dry "github.com/ungerik/go-dry"
)

type CreateInstanceAuditPlanReqV1 struct {
	InstanceId string `json:"instance_id" form:"instance_id"  valid:"required"`
	// 扫描类型
	AuditPlans []AuditPlan `json:"audit_plans" form:"audit_plans" valid:"required"`
}

type AuditPlan struct {
	RuleTemplateName        string                     `json:"rule_template_name" from:"rule_template_name" example:"default_MySQL"`
	Type                    string                     `json:"audit_plan_type" form:"audit_plan_type" example:"slow log"`
	Params                  []AuditPlanParamReqV1      `json:"audit_plan_params" valid:"dive,required"`
	HighPriorityConditions  []HighPriorityConditionReq `json:"high_priority_conditions" valid:"dive,required"`
	NeedMarkHighPrioritySQL bool                       `json:"need_mark_high_priority_sql"`
}
type HighPriorityConditionReq struct {
	Key      string `json:"key" form:"key" valid:"required"`
	Value    string `json:"value" form:"value" valid:"required"`
	Operator string `json:"operator" form:"operator" default:">" enums:">,=,<" valid:"oneof=> = <"`
}

type CreatInstanceAuditPlanResV1 struct {
	controller.BaseRes
	Data CreatInstanceAuditPlanRes `json:"data"`
}

type CreatInstanceAuditPlanRes struct {
	InstanceAuditPlanID string `json:"instance_audit_plan_id"`
}

func checkAndGenerateHighPriorityParams(auditPlanType, instanceType string, hpcParamsReq []HighPriorityConditionReq) (params.ParamsWithOperator, error) {
	meta, err := auditplan.GetMeta(auditPlanType)
	if err != nil {
		return nil, err
	}
	if meta.InstanceType != auditplan.InstanceTypeAll && meta.InstanceType != instanceType {
		return nil, fmt.Errorf("audit plan type %s not found", auditPlanType)
	}
	resetParams := make([]*params.ParamWithOperator, 0)
	for _, hpcParam := range hpcParamsReq {
		for _, p := range meta.HighPriorityParams {
			if p.Key != hpcParam.Key {
				continue
			}
			// set and valid param.
			p.Value = hpcParam.Value
			p.Operator.Value = params.OperatorValue(hpcParam.Operator)
			resetParams = append(resetParams, p)
			break
		}

	}
	return resetParams, nil
}

// @Summary 添加实例扫描任务
// @Description create instance audit plan
// @Id createInstanceAuditPlanV1
// @Tags instance_audit_plan
// @Security ApiKeyAuth
// @Accept json
// @Param project_name path string true "project name"
// @Param instannce_audit_plan body v1.CreateInstanceAuditPlanReqV1 true "create instance audit plan"
// @Success 200 {object} v1.CreatInstanceAuditPlanResV1
// @router /v1/projects/{project_name}/instance_audit_plans [post]
func CreateInstanceAuditPlan(c echo.Context) error {
	req := new(CreateInstanceAuditPlanReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// check instance
	inst, exist, err := dms.GetInstancesById(c.Request().Context(), req.InstanceId)
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrInstanceNotExist)
	} else if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
	}

	if !dry.StringInSlice(inst.DbType, driver.GetPluginManager().AllDrivers()) {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DriverNotExist, &driverV2.DriverNotSupportedError{DriverTyp: inst.DbType}))
	}

	// check instance audit plan exist
	_, exist, err = model.GetStorage().GetInstanceAuditPlanByInstanceID(int64(inst.ID))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if exist {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("current instance has audit plan"))
	}
	// check operation
	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	canCreateAuditPlan, err := CheckUserCanCreateAuditPlan(c.Request().Context(), projectUid, user, []*model.Instance{inst})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !canCreateAuditPlan {
		return controller.JSONBaseErrorReq(c, errors.NewUserNotPermissionError(model.GetOperationCodeDesc(c.Request().Context(), uint(model.OP_AUDIT_PLAN_SAVE))))
	}

	userId := controller.GetUserID(c)
	s := model.GetStorage()

	auditPlans := make([]*model.AuditPlanV2, 0)
	for _, auditPlan := range req.AuditPlans {
		if auditPlan.RuleTemplateName != "" {
			exist, err := s.IsRuleTemplateExist(auditPlan.RuleTemplateName, []string{projectUid, model.ProjectIdForGlobalRuleTemplate})
			if err != nil {
				return controller.JSONBaseErrorReq(c, err)
			}
			if !exist {
				return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("rule template does not exist")))
			}
		}
		// check rule template name
		ruleTemplateName, err := autoSelectRuleTemplate(c.Request().Context(), auditPlan.RuleTemplateName, inst.Name, inst.DbType, projectUid)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}

		// check params
		if auditPlan.Type == "" {
			auditPlan.Type = auditplan.TypeDefault
		}
		ps, err := checkAndGenerateAuditPlanParams(auditPlan.Type, inst.DbType, auditPlan.Params)
		if err != nil {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
		}

		hpc, err := checkAndGenerateHighPriorityParams(auditPlan.Type, inst.DbType, auditPlan.HighPriorityConditions)
		if err != nil {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
		}

		auditPlans = append(auditPlans, &model.AuditPlanV2{
			Type:                    auditPlan.Type,
			RuleTemplateName:        ruleTemplateName,
			Params:                  ps,
			HighPriorityParams:      hpc,
			NeedMarkHighPrioritySQL: auditPlan.NeedMarkHighPrioritySQL,
			ActiveStatus:            model.ActiveStatusNormal,
		})
	}

	ap := &model.InstanceAuditPlan{
		ProjectId:    model.ProjectUID(projectUid),
		InstanceID:   inst.ID,
		DBType:       inst.DbType,
		CreateUserID: userId,
		AuditPlans:   auditPlans,
		ActiveStatus: model.ActiveStatusNormal,
	}
	err = s.SaveInstanceAuditPlan(ap)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// generate token , 生成ID后根据ID生成token
	t, err := dmsCommonJwt.GenJwtToken(dmsCommonJwt.WithUserId(userId), dmsCommonJwt.WithExpiredTime(tokenExpire), dmsCommonJwt.WithAuditPlanName(utils.Md5(ap.GetIDStr())))
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
	}
	err = s.UpdateInstanceAuditPlanByID(ap.ID, map[string]interface{}{"token": t})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	resID := CreatInstanceAuditPlanRes{
		InstanceAuditPlanID: ap.GetIDStr(),
	}
	return c.JSON(http.StatusOK, &CreatInstanceAuditPlanResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    resID,
	})
}

// @Summary 删除实例扫描任务
// @Description delete instance audit plan
// @Id deleteInstanceAuditPlanV1
// @Tags instance_audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instance_audit_plan_id path string true "instance audit plan id"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/instance_audit_plans/{instance_audit_plan_id}/ [delete]
func DeleteInstanceAuditPlan(c echo.Context) error {
	instanceAuditPlanID := c.Param("instance_audit_plan_id")
	projectUID, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	s := model.GetStorage()
	// check current user instance audit plan permission
	_, exist, err := GetInstanceAuditPlanIfCurrentUserCanAccess(c, projectUID, instanceAuditPlanID, v1.OpPermissionTypeSaveAuditPlan)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewInstanceAuditPlanNotExistErr())
	}

	err = s.DeleteInstanceAuditPlan(instanceAuditPlanID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return controller.JSONBaseErrorReq(c, nil)
}

type UpdateInstanceAuditPlanReqV1 struct {
	// 扫描类型
	AuditPlans []AuditPlan `json:"audit_plans" form:"audit_plans" valid:"required"`
}

// @Summary 更新实例扫描任务配置
// @Description update instance audit plan
// @Id updateInstanceAuditPlanV1
// @Tags instance_audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instance_audit_plan_id path string true "instance audit plan id"
// @param audit_plan body v1.UpdateInstanceAuditPlanReqV1 true "update instance audit plan"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/instance_audit_plans/{instance_audit_plan_id}/ [put]
func UpdateInstanceAuditPlan(c echo.Context) error {
	req := new(UpdateInstanceAuditPlanReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	instanceAuditPlanID := c.Param("instance_audit_plan_id")
	projectUID, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// check current user instance audit plan permission
	dbAuditPlans, exist, err := GetInstanceAuditPlanIfCurrentUserCanAccess(c, projectUID, instanceAuditPlanID, v1.OpPermissionTypeSaveAuditPlan)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewInstanceAuditPlanNotExistErr())
	}

	s := model.GetStorage()
	reqAuditPlansMap := make(map[string]AuditPlan)
	for _, auditPlan := range req.AuditPlans {
		reqAuditPlansMap[auditPlan.Type] = auditPlan
	}
	dbAuditPlansMap := make(map[string]*model.AuditPlanV2)
	// check db audit plans all are in the req audit plans
	for _, dbAuditPlan := range dbAuditPlans.AuditPlans {
		if _, ok := reqAuditPlansMap[dbAuditPlan.Type]; !ok {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid,
				fmt.Errorf("the audit plan is not allowed to be deleted at update")))
		}
		dbAuditPlansMap[dbAuditPlan.Type] = dbAuditPlan
	}

	resultAuditPlans := make([]*model.AuditPlanV2, 0)
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	inst := dms.GetInstancesByIdWithoutError(fmt.Sprintf("%d", dbAuditPlans.InstanceID))
	for _, auditPlan := range req.AuditPlans {
		if auditPlan.RuleTemplateName != "" {
			exist, err := s.IsRuleTemplateExist(auditPlan.RuleTemplateName, []string{projectUid, model.ProjectIdForGlobalRuleTemplate})
			if err != nil {
				return controller.JSONBaseErrorReq(c, err)
			}
			if !exist {
				return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("rule template does not exist")))
			}
		}
		// check rule template name
		ruleTemplateName, err := autoSelectRuleTemplate(c.Request().Context(), auditPlan.RuleTemplateName, inst.Name, dbAuditPlans.DBType, projectUid)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}

		// check params
		if auditPlan.Type == "" {
			auditPlan.Type = auditplan.TypeDefault
		}
		ps, err := checkAndGenerateAuditPlanParams(auditPlan.Type, dbAuditPlans.DBType, auditPlan.Params)
		if err != nil {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
		}
		hpc, err := checkAndGenerateHighPriorityParams(auditPlan.Type, dbAuditPlans.DBType, auditPlan.HighPriorityConditions)
		if err != nil {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
		}
		res := &model.AuditPlanV2{
			Type:                    auditPlan.Type,
			RuleTemplateName:        ruleTemplateName,
			Params:                  ps,
			HighPriorityParams:      hpc,
			NeedMarkHighPrioritySQL: auditPlan.NeedMarkHighPrioritySQL,
			InstanceAuditPlanID:     dbAuditPlans.ID,
		}

		// if the data exists in the database, update the data; if it does not exist, insert the data.
		if dbAuditPlan, ok := dbAuditPlansMap[auditPlan.Type]; ok {
			dbAuditPlan.RuleTemplateName = res.RuleTemplateName
			dbAuditPlan.Params = res.Params
			dbAuditPlan.HighPriorityParams = res.HighPriorityParams
			dbAuditPlan.NeedMarkHighPrioritySQL = res.NeedMarkHighPrioritySQL
			result := dbAuditPlan
			resultAuditPlans = append(resultAuditPlans, result)
		} else {
			if dbAuditPlans.ActiveStatus == model.ActiveStatusNormal {
				res.ActiveStatus = model.ActiveStatusNormal
			} else {
				res.ActiveStatus = model.ActiveStatusDisabled
			}
			resultAuditPlans = append(resultAuditPlans, res)
		}
	}

	err = s.BatchSaveAuditPlans(resultAuditPlans)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return controller.JSONBaseErrorReq(c, nil)
}

// @Summary 更新实例扫描任务状态
// @Description stop/start instance audit plan
// @Id updateInstanceAuditPlanStatusV1
// @Tags instance_audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instance_audit_plan_id path string true "instance audit plan id"
// @param audit_plan body v1.UpdateInstanceAuditPlanStatusReqV1 true "update instance audit plan"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/instance_audit_plans/{instance_audit_plan_id}/ [patch]
func UpdateInstanceAuditPlanStatus(c echo.Context) error {
	req := new(UpdateInstanceAuditPlanStatusReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	s := model.GetStorage()
	instanceAuditPlanID := c.Param("instance_audit_plan_id")
	projectUID, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// check current user instance audit plan permission
	instanceAuditPlan, exist, err := GetInstanceAuditPlanIfCurrentUserCanAccess(c, projectUID, instanceAuditPlanID, v1.OpPermissionTypeSaveAuditPlan)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewInstanceAuditPlanNotExistErr())
	}

	instanceAuditPlan.ActiveStatus = req.Active
	if req.Active == model.ActiveStatusDisabled {
		for _, auditPlan := range instanceAuditPlan.AuditPlans {
			auditPlan.ActiveStatus = model.ActiveStatusDisabled
		}
		err = s.BatchSaveAuditPlans(instanceAuditPlan.AuditPlans)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}
	err = s.Save(instanceAuditPlan)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return controller.JSONBaseErrorReq(c, nil)
}

type AuditPlanTypeResBase struct {
	AuditPlanId       uint   `json:"audit_plan_id"`
	AuditPlanType     string `json:"type"`
	AuditPlanTypeDesc string `json:"desc"`
}

type GetInstanceAuditPlansReqV1 struct {
	FilterByBusiness      string `json:"filter_by_business" query:"filter_by_business"`
	FilterByDBType        string `json:"filter_by_db_type" query:"filter_by_db_type"`
	FilterByInstanceID    string `json:"filter_by_instance_id" query:"filter_by_instance_id"`
	FilterByAuditPlanType string `json:"filter_by_audit_plan_type" query:"filter_by_audit_plan_type"`
	FilterByActiveStatus  string `json:"filter_by_active_status" query:"filter_by_active_status"`
	FuzzySearch           string `json:"fuzzy_search" query:"fuzzy_search"`

	PageIndex uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetInstanceAuditPlansResV1 struct {
	controller.BaseRes
	Data      []InstanceAuditPlanResV1 `json:"data"`
	TotalNums uint64                   `json:"total_nums"`
}

type InstanceAuditPlanResV1 struct {
	InstanceAuditPlanId uint                   `json:"instance_audit_plan_id"`
	InstanceID          string                 `json:"instance_id"`
	InstanceName        string                 `json:"instance_name"`
	Business            string                 `json:"business"`
	InstanceType        string                 `json:"instance_type"`
	AuditPlanTypes      []AuditPlanTypeResBase `json:"audit_plan_types"`
	ActiveStatus        string                 `json:"active_status" enums:"normal,disabled"`
	// TODO 采集状态
	CreateTime string `json:"create_time"`
	Creator    string `json:"creator"`
}

// GetInstanceAuditPlans
// @Summary 获取实例扫描任务列表
// @Description get instance audit plan info list
// @Id getInstanceAuditPlansV1
// @Tags instance_audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param filter_by_business query string false "filter by business"
// @Param filter_by_db_type query string false "filter by db type"
// @Param filter_by_instance_id query string false "filter by instance id"
// @Param filter_by_audit_plan_type query string false "filter instance audit plan type"
// @Param filter_by_active_status query string false "filter instance audit plan active status"
// @Param fuzzy_search query string false "fuzzy search"
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Success 200 {object} v1.GetInstanceAuditPlansResV1
// @router /v1/projects/{project_name}/instance_audit_plans [get]
func GetInstanceAuditPlans(c echo.Context) error {
	s := model.GetStorage()

	req := new(GetInstanceAuditPlansReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"))
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
		"filter_by_business":                 req.FilterByBusiness,
		"filter_project_id":                  projectUid,
		"current_user_id":                    userId,
		"current_user_is_admin":              up.IsAdmin(),
		"filter_by_active_status":            req.FilterByActiveStatus,
		"limit":                              limit,
		"offset":                             offset,
	}
	if !up.IsAdmin() {
		accessinstanceId := up.GetInstancesByOP(v1.OpPermissionTypeViewOtherAuditPlan)
		if len(accessinstanceId) > 0 {
			data["accessible_instances_id"] = fmt.Sprintf("\"%s\"", strings.Join(accessinstanceId, "\",\""))
		}
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
				typeBase, err := ConvertAuditPlanTypeToResByID(c.Request().Context(), auditPlanId)
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
			Business:            inst.Business,
			InstanceType:        v.DBType,
			AuditPlanTypes:      typeBases,
			ActiveStatus:        v.ActiveStatus,
			CreateTime:          v.CreateTime,
			Creator:             dms.GetUserNameWithDelTag(v.CreateUserId),
		}
	}
	return c.JSON(http.StatusOK, &GetInstanceAuditPlansResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      resData,
		TotalNums: count,
	})
}

func ConvertAuditPlanTypeToResByID(ctx context.Context, id string) (AuditPlanTypeResBase, error) {
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
				AuditPlanType:     auditPlan.Type,
				AuditPlanTypeDesc: locale.ShouldLocalizeMsg(ctx, meta.Desc),
				AuditPlanId:       auditPlan.ID,
			}, nil
		}
	}
	return AuditPlanTypeResBase{}, nil
}

func ConvertAuditPlanTypeToRes(ctx context.Context, id uint, auditPlanType string) AuditPlanTypeResBase {
	for _, meta := range auditplan.Metas {
		if meta.Type == auditPlanType {
			return AuditPlanTypeResBase{
				AuditPlanType:     auditPlanType,
				AuditPlanTypeDesc: locale.ShouldLocalizeMsg(ctx, meta.Desc),
				AuditPlanId:       id,
			}
		}
	}
	return AuditPlanTypeResBase{}
}

type GetInstanceAuditPlanDetailResV1 struct {
	controller.BaseRes
	Data InstanceAuditPlanDetailResV1 `json:"data"`
}

type InstanceAuditPlanDetailResV1 struct {
	Business     string `json:"business"     example:"test"`
	InstanceType string `json:"instance_type" example:"mysql" `
	InstanceName string `json:"instance_name" example:"test_mysql"`
	InstanceID   string `json:"instance_id" example:"instance_id"`
	// 扫描类型
	AuditPlans []AuditPlanRes `json:"audit_plans"`
}

type AuditPlanRes struct {
	RuleTemplateName        string                  `json:"rule_template_name" from:"rule_template_name" example:"default_MySQL"`
	Type                    AuditPlanTypeResBase    `json:"audit_plan_type" form:"audit_plan_type"`
	Params                  []AuditPlanParamResV1   `json:"audit_plan_params" valid:"dive,required"`
	NeedMarkHighPrioritySQL bool                    `json:"need_mark_high_priority_sql"`
	HighPriorityConditions  []HighPriorityCondition `json:"high_priority_conditions"`
}

// @Summary 获取实例扫描任务详情
// @Description get instance audit plan detail
// @Id getInstanceAuditPlanDetailV1
// @Tags instance_audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instance_audit_plan_id path string true "instance audit plan id"
// @Success 200 {object} v1.GetInstanceAuditPlanDetailResV1
// @router /v1/projects/{project_name}/instance_audit_plans/{instance_audit_plan_id} [get]
func GetInstanceAuditPlanDetail(c echo.Context) error {
	instanceAuditPlanID := c.Param("instance_audit_plan_id")
	projectUID, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// check current user instance audit plan permission
	detail, exist, err := GetInstanceAuditPlanIfCurrentUserCanAccess(c, projectUID, instanceAuditPlanID, v1.OpPermissionTypeViewOtherAuditPlan)
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
	resData := InstanceAuditPlanDetailResV1{
		Business:     inst.Business,
		InstanceType: detail.DBType,
		InstanceName: inst.Name,
		InstanceID:   inst.GetIDStr(),
		AuditPlans:   auditPlans,
	}
	return c.JSON(http.StatusOK, &GetInstanceAuditPlanDetailResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    resData,
	})
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
			paramsRes := make([]AuditPlanParamResV1, 0, len(meta.Params()))
			for _, p := range meta.Params() {
				val := p.Value
				if p.Type == params.ParamTypePassword {
					val = ""
				}
				paramRes := AuditPlanParamResV1{
					Key:   p.Key,
					Desc:  p.GetDesc(locale.GetLangTagFromCtx(ctx)),
					Type:  string(p.Type),
					Value: val,
				}
				paramsRes = append(paramsRes, paramRes)
			}
			resAuditPlan.Params = paramsRes
		}

		if v.HighPriorityParams != nil && len(v.HighPriorityParams) > 0 {
			hppParamsRes := make([]HighPriorityCondition, len(v.HighPriorityParams))
			for i, hpp := range v.HighPriorityParams {
				for _, metaHpp := range meta.HighPriorityParams {
					if metaHpp.Key != hpp.Key {
						continue
					}
					highParamRes := HighPriorityCondition{
						Key:   metaHpp.Key,
						Desc:  metaHpp.GetDesc(locale.GetLangTagFromCtx(ctx)),
						Value: hpp.Value,
						Type:  string(metaHpp.Type),
						Operator: Operator{
							Value:      string(hpp.Operator.Value),
							EnumsValue: metaHpp.Operator.EnumsValue,
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

type GetInstanceAuditPlanOverviewResV1 struct {
	controller.BaseRes
	Data []InstanceAuditPlanInfo `json:"data"`
}

type AuditPlanRuleTemplate struct {
	Name                 string `json:"name"`
	IsGlobalRuleTemplate bool   `json:"is_global_rule_template"`
}

type InstanceAuditPlanInfo struct {
	ID                 uint                   `json:"id"`
	Type               AuditPlanTypeResBase   `json:"audit_plan_type"`
	DBType             string                 `json:"audit_plan_db_type" example:"mysql"`
	InstanceName       string                 `json:"audit_plan_instance_name" example:"test_mysql"`
	ExecCmd            string                 `json:"exec_cmd" example:"./scanner xxx"`
	RuleTemplate       *AuditPlanRuleTemplate `json:"audit_plan_rule_template,omitempty" `
	TotalSQLNums       int64                  `json:"total_sql_nums"`
	UnsolvedSQLNums    int64                  `json:"unsolved_sql_nums"`
	LastCollectionTime *time.Time             `json:"last_collection_time"`
	ActiveStatus       string                 `json:"active_status" enums:"normal,disabled"`
}

// @Summary 获取实例扫描任务概览
// @Description get audit plan overview
// @Id getInstanceAuditPlanOverviewV1
// @Tags instance_audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instance_audit_plan_id path string true "instance audit plan id"
// @Success 200 {object} v1.GetInstanceAuditPlanOverviewResV1
// @router /v1/projects/{project_name}/instance_audit_plans/{instance_audit_plan_id}/audit_plans [get]
func GetInstanceAuditPlanOverview(c echo.Context) error {
	instanceAuditPlanID := c.Param("instance_audit_plan_id")
	projectName := c.Param("project_name")
	projectUID, err := dms.GetPorjectUIDByName(c.Request().Context(), projectName, true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// check current user instance audit plan permission
	detail, exist, err := GetInstanceAuditPlanIfCurrentUserCanAccess(c, projectUID, instanceAuditPlanID, v1.OpPermissionTypeViewOtherAuditPlan)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewInstanceAuditPlanNotExistErr())
	}
	s := model.GetStorage()

	inst := dms.GetInstancesByIdWithoutError(fmt.Sprintf("%d", detail.InstanceID))
	resAuditPlans := make([]InstanceAuditPlanInfo, 0, len(detail.AuditPlans))
	for _, v := range detail.AuditPlans {
		execCmd := GetAuditPlanExecCmd(projectName, detail, v)

		totalSQLNums, err := s.GetAuditPlanTotalSQL(v.ID)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}

		unsolvedSQLNums, err := getAuditPlanUnsolvedSQLCount(v.ID)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		template, exist, err := s.GetGlobalAndProjectRuleTemplateByNameAndProjectId(v.RuleTemplateName, projectUID)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		ruleTemplate := &AuditPlanRuleTemplate{}
		if exist {
			ruleTemplate.Name = v.RuleTemplateName
			ruleTemplate.IsGlobalRuleTemplate = (template.ProjectId == model.ProjectIdForGlobalRuleTemplate)
		}

		typeBase := ConvertAuditPlanTypeToRes(c.Request().Context(), v.ID, v.Type)
		resAuditPlan := InstanceAuditPlanInfo{
			ID:              v.ID,
			Type:            typeBase,
			DBType:          detail.DBType,
			InstanceName:    inst.Name,
			ExecCmd:         execCmd,
			RuleTemplate:    ruleTemplate,
			TotalSQLNums:    totalSQLNums,
			UnsolvedSQLNums: unsolvedSQLNums,
			ActiveStatus:    v.ActiveStatus,
		}
		if v.AuditPlanTaskInfo != nil {
			resAuditPlan.LastCollectionTime = v.AuditPlanTaskInfo.LastCollectionTime
		}
		resAuditPlans = append(resAuditPlans, resAuditPlan)
	}

	return c.JSON(http.StatusOK, &GetInstanceAuditPlanOverviewResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    resAuditPlans,
	})
}

func GetAuditPlanExecCmd(projectName string, iap *model.InstanceAuditPlan, ap *model.AuditPlanV2) string {
	logger := log.NewEntry().WithField("get audit plan exec cmd", fmt.Sprintf("inst id:%d,audit plan type : %s", iap.InstanceID, ap.Type))
	_, ok := auditplan.GetSupportedScannerAuditPlanType()[ap.Type]
	if !ok {
		return ""
	}

	address := config.GetOptions().SqleOptions.DMSServerAddress
	parsedURL, err := url.Parse(address)
	if err != nil {
		logger.Info("parse server address failed ", err)
		return ""
	}
	ip, port, err := net.SplitHostPort(parsedURL.Host)
	if err != nil {
		logger.Info("split server host failed ", err)
		return ""
	}

	scannerd, err := scannerCmd.GetScannerdCmd(ap.Type)
	if err != nil {
		logger.Infof("get scannerd %s failed %s", ap.Type, err)
		return ""
	}
	cmd, err := scannerd.GenCommand("./scannerd", map[string]string{
		scannerCmd.FlagHost:        ip,
		scannerCmd.FlagPort:        port,
		scannerCmd.FlagProject:     projectName,
		scannerCmd.FlagToken:       iap.Token,
		scannerCmd.FlagAuditPlanID: fmt.Sprint(ap.ID),
	})
	if err != nil {
		logger.Infof("generate scannerd %s command failed %s", ap.Type, err)
		return ""
	}
	return cmd
}

type UpdateInstanceAuditPlanStatusReqV1 struct {
	// 任务状态
	Active string `json:"active" form:"active" enums:"normal,disabled"`
}

// @Summary 删除扫描任务
// @Description delete audit plan by type
// @Id deleteAuditPlanByTypeV1
// @Tags instance_audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instance_audit_plan_id path string true "instance audit plan id"
// @Param audit_plan_id path string true "audit plan id"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/instance_audit_plans/{instance_audit_plan_id}/audit_plans/{audit_plan_id}/ [delete]
func DeleteAuditPlanById(c echo.Context) error {
	instanceAuditPlanID := c.Param("instance_audit_plan_id")
	projectUID, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// check current user instance audit plan permission
	_, exist, err := GetInstanceAuditPlanIfCurrentUserCanAccess(c, projectUID, instanceAuditPlanID, v1.OpPermissionTypeSaveAuditPlan)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewInstanceAuditPlanNotExistErr())
	}
	s := model.GetStorage()

	audit_plan_id, err := strconv.Atoi(c.Param("audit_plan_id"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.NewInstanceAuditPlanNotExistErr())
	}
	err = s.DeleteAuditPlan(audit_plan_id)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return controller.JSONBaseErrorReq(c, nil)
}

type UpdateAuditPlanStatusReqV1 struct {
	// 任务状态
	Active string `json:"active" form:"active" enums:"normal,disabled"`
}

// @Summary 更新扫描任务状态
// @Description stop/start audit plan
// @Id updateAuditPlanStatusV1
// @Tags instance_audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instance_audit_plan_id path string true "instance audit plan id"
// @Param audit_plan_id path string true "audit plan id"
// @param audit_plan body v1.UpdateAuditPlanStatusReqV1 true "update audit plan status"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/instance_audit_plans/{instance_audit_plan_id}/audit_plans/{audit_plan_id}/ [patch]
func UpdateAuditPlanStatus(c echo.Context) error {
	req := new(UpdateAuditPlanStatusReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	instanceAuditPlanID := c.Param("instance_audit_plan_id")
	projectUID, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// check current user instance audit plan permission
	deatil, exist, err := GetInstanceAuditPlanIfCurrentUserCanAccess(c, projectUID, instanceAuditPlanID, v1.OpPermissionTypeSaveAuditPlan)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewInstanceAuditPlanNotExistErr())
	}
	if deatil.ActiveStatus != model.ActiveStatusNormal && req.Active == model.ActiveStatusNormal {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf("instance audit plan active status not normal")))
	}
	audit_plan_id, err := strconv.Atoi(c.Param("audit_plan_id"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.NewInstanceAuditPlanNotExistErr())
	}
	s := model.GetStorage()
	auditPlan, exist, err := s.GetAuditPlanByID(audit_plan_id)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewAuditPlanNotExistErr())
	}
	auditPlan.ActiveStatus = req.Active
	err = s.Save(auditPlan)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return controller.JSONBaseErrorReq(c, nil)
}

// @Summary 获取指定扫描任务的SQLs信息
// @Description get audit plan SQLs
// @Id getInstanceAuditPlanSQLsV1
// @Tags instance_audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instance_audit_plan_id path string true "instance audit plan id"
// @Param audit_plan_id path string true "audit plan id"
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Success 200 {object} v1.GetAuditPlanSQLsResV1
// @router /v1/projects/{project_name}/instance_audit_plans/{instance_audit_plan_id}/audit_plans/{audit_plan_id}/sqls [get]
func GetInstanceAuditPlanSQLs(c echo.Context) error {
	req := new(GetAuditPlanSQLsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	instanceAuditPlanID := c.Param("instance_audit_plan_id")
	projectUID, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// check current user instance audit plan permission
	_, exist, err := GetInstanceAuditPlanIfCurrentUserCanAccess(c, projectUID, instanceAuditPlanID, v1.OpPermissionTypeViewOtherAuditPlan)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewInstanceAuditPlanNotExistErr())
	}
	apID, err := strconv.Atoi(c.Param("audit_plan_id"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	ap, err := s.GetAuditPlanDetailByID(uint(apID))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	limit, offset := controller.GetLimitAndOffset(req.PageIndex, req.PageSize)

	data := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}

	head, rows, count, err := auditplan.GetSQLs(log.NewEntry(), auditplan.ConvertModelToAuditPlanV2(ap), data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// convert head and rows to audit planSQL response
	res := AuditPlanSQLResV1{
		Rows: rows,
	}
	for _, v := range head {
		res.Head = append(res.Head, AuditPlanSQLHeadV1{
			Name: v.Name,
			Desc: locale.ShouldLocalizeMsg(c.Request().Context(), v.Desc),
			Type: v.Type,
		})
	}
	return c.JSON(http.StatusOK, &GetAuditPlanSQLsResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      res,
		TotalNums: count,
	})
}

type GetAuditPlanSQLMetaResV1 struct {
	controller.BaseRes
	Data AuditPlanSQLMetaResV1 `json:"data"`
}

type AuditPlanSQLMetaResV1 struct {
	Head           []AuditPlanSQLHeadV1 `json:"head"`
	FilterMetaList []FilterMeta         `json:"filter_meta_list"`
}

type FilterMeta struct {
	Name            string      `json:"filter_name"`
	Desc            string      `json:"desc"`
	FilterInputType string      `json:"filter_input_type" enums:"int,string,date_time"`
	FilterOpType    string      `json:"filter_op_type" enums:"equal,between"`
	FilterTips      []FilterTip `json:"filter_tip_list"`
}

type FilterTip struct {
	Value string `json:"value"`
	Desc  string `json:"desc"`
	Group string `json:"group"`
}

type Filter struct {
	Name                  string             `json:"filter_name"`
	FilterComparisonValue string             `json:"filter_compare_value"`
	FilterBetweenValue    FilterBetweenValue `json:"filter_between_value"`
}

type FilterBetweenValue struct {
	From string
	To   string
}

// @Summary 获取指定扫描任务的SQL列表元信息
// @Description get audit plan SQL meta
// @Id getInstanceAuditPlanSQLMetaV1
// @Tags instance_audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instance_audit_plan_id path string true "instance audit plan id"
// @Param audit_plan_id path string true "audit plan id"
// @Success 200 {object} v1.GetAuditPlanSQLMetaResV1
// @router /v1/projects/{project_name}/instance_audit_plans/{instance_audit_plan_id}/audit_plans/{audit_plan_id}/sql_meta [get]
func GetInstanceAuditPlanSQLMeta(c echo.Context) error {
	instanceAuditPlanID := c.Param("instance_audit_plan_id")
	projectUID, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// check current user instance audit plan permission
	_, exist, err := GetInstanceAuditPlanIfCurrentUserCanAccess(c, projectUID, instanceAuditPlanID, v1.OpPermissionTypeViewOtherAuditPlan)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewInstanceAuditPlanNotExistErr())
	}
	ctx := c.Request().Context()
	s := model.GetStorage()
	auditPlanId, err := strconv.Atoi(c.Param("audit_plan_id"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.NewAuditPlanNotExistErr())
	}
	apDetail, err := s.GetAuditPlanDetailByID(uint(auditPlanId))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	ap := auditplan.ConvertModelToAuditPlanV2(apDetail)
	head, err := auditplan.GetSQLHead(ap, s)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	filter, err := auditplan.GetSQLFilterMeta(ctx, ap, s)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	data := AuditPlanSQLMetaResV1{}
	for _, v := range head {
		data.Head = append(data.Head, AuditPlanSQLHeadV1{
			Name:     v.Name,
			Desc:     locale.ShouldLocalizeMsg(ctx, v.Desc),
			Type:     v.Type,
			Sortable: v.Sortable,
		})
	}
	for _, v := range filter {
		data.FilterMetaList = append(data.FilterMetaList, FilterMeta{
			Name:            v.Name,
			Desc:            locale.ShouldLocalizeMsg(ctx, v.Desc),
			FilterInputType: string(v.FilterInputType),
			FilterOpType:    string(v.FilterOpType),
			FilterTips:      ConvertFilterTipsToRes(v.FilterTips),
		})
	}

	return c.JSON(http.StatusOK, &GetAuditPlanSQLMetaResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    data,
	})
}

type GetAuditPlanSQLDataReqV1 struct {
	PageIndex uint32   `json:"page_index" valid:"required"`
	PageSize  uint32   `json:"page_size" valid:"required"`
	OrderBy   string   `json:"order_by"`
	IsAsc     bool     `json:"is_asc"`
	Filters   []Filter `json:"filter_list"`
}

type GetAuditPlanSQLDataResV1 struct {
	controller.BaseRes
	Data      AuditPlanSQLDataResV1 `json:"data"`
	TotalNums uint64                `json:"total_nums"`
}

type AuditPlanSQLDataResV1 struct {
	Rows []map[string] /* head name */ string `json:"rows"`
}

// @Summary 获取指定扫描任务的SQL列表
// @Description get audit plan SQLs
// @Id getInstanceAuditPlanSQLDataV1
// @Tags instance_audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instance_audit_plan_id path string true "instance audit plan id"
// @Param audit_plan_id path string true "audit plan id"
// @param audit_plan_sql_request body v1.GetAuditPlanSQLDataReqV1 true "audit plan sql data request"
// @Success 200 {object} v1.GetAuditPlanSQLDataResV1
// @router /v1/projects/{project_name}/instance_audit_plans/{instance_audit_plan_id}/audit_plans/{audit_plan_id}/sql_data [post]
func GetInstanceAuditPlanSQLData(c echo.Context) error {
	req := new(GetAuditPlanSQLDataReqV1)
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if err = json.Unmarshal(body, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if err = controller.Validate(req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	instanceAuditPlanID := c.Param("instance_audit_plan_id")
	projectUID, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// check current user instance audit plan permission
	_, exist, err := GetInstanceAuditPlanIfCurrentUserCanAccess(c, projectUID, instanceAuditPlanID, v1.OpPermissionTypeViewOtherAuditPlan)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewInstanceAuditPlanNotExistErr())
	}
	s := model.GetStorage()
	auditPlanId, err := strconv.Atoi(c.Param("audit_plan_id"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.NewAuditPlanNotExistErr())
	}
	apDetail, err := s.GetAuditPlanDetailByID(uint(auditPlanId))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	ap := auditplan.ConvertModelToAuditPlanV2(apDetail)

	data, count, err := auditplan.GetSQLData(ap, s, ConvertReqToAuditPlanFilter(req.Filters), req.OrderBy, req.IsAsc, int(req.PageSize), int((req.PageIndex-1)*req.PageSize))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetAuditPlanSQLDataResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      AuditPlanSQLDataResV1{Rows: data},
		TotalNums: count,
	})
}

type GetAuditPlanSQLExportReqV1 struct {
	OrderBy string   `json:"order_by"`
	IsAsc   bool     `json:"is_asc"`
	Filters []Filter `json:"filter_list"`
}

// @Summary 导出指定扫描任务的 SQL CSV 列表
// @Description export audit plan SQL report as CSV
// @Id getInstanceAuditPlanSQLExportV1
// @Tags instance_audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instance_audit_plan_id path string true "instance audit plan id"
// @Param audit_plan_id path string true "audit plan id"
// @param audit_plan_sql_request body v1.GetAuditPlanSQLExportReqV1 true "audit plan sql export request"
// @Success 200 {file} file "export audit plan sql report"
// @router /v1/projects/{project_name}/instance_audit_plans/{instance_audit_plan_id}/audit_plans/{audit_plan_id}/sql_export [post]
func GetInstanceAuditPlanSQLExport(c echo.Context) error {
	req := new(GetAuditPlanSQLExportReqV1)
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if err = json.Unmarshal(body, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if err = controller.Validate(req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	instanceAuditPlanID := c.Param("instance_audit_plan_id")
	projectUID, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// check current user instance audit plan permission
	_, exist, err := GetInstanceAuditPlanIfCurrentUserCanAccess(c, projectUID, instanceAuditPlanID, v1.OpPermissionTypeViewOtherAuditPlan)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewInstanceAuditPlanNotExistErr())
	}
	s := model.GetStorage()
	auditPlanId, err := strconv.Atoi(c.Param("audit_plan_id"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.NewAuditPlanNotExistErr())
	}
	apDetail, err := s.GetAuditPlanDetailByID(uint(auditPlanId))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	ap := auditplan.ConvertModelToAuditPlanV2(apDetail)
	head, err := auditplan.GetSQLHead(ap, s)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	rows, _, err := auditplan.GetSQLData(ap, s, ConvertReqToAuditPlanFilter(req.Filters), req.OrderBy, req.IsAsc, 0, 0)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	buff := new(bytes.Buffer)
	buff.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	csvWriter := csv.NewWriter(buff)
	toWrite := make([]string, len(head))
	for col, h := range head {
		toWrite[col] = locale.ShouldLocalizeMsg(c.Request().Context(), h.Desc)
	}
	if err = csvWriter.Write(toWrite); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	for _, rowMap := range rows {
		for col, h := range head {
			toWrite[col] = utils.TruncateAndMarkForExcelCell(rowMap[h.Name])
		}
		if err = csvWriter.Write(toWrite); err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}
	csvWriter.Flush()

	fileName := fmt.Sprintf("sql_export_%s_%s.csv", ap.Type, time.Now().Format("20060102150405"))
	c.Response().Header().Set(echo.HeaderContentDisposition, mime.FormatMediaType("attachment", map[string]string{"filename": fileName}))

	return c.Blob(http.StatusOK, "text/csv", buff.Bytes())
}

func ConvertFilterTipsToRes(fts []auditplan.FilterTip) []FilterTip {
	resAuditPlans := make([]FilterTip, 0, len(fts))
	for _, v := range fts {
		resAuditPlans = append(resAuditPlans, FilterTip{
			Value: v.Value,
			Desc:  v.Desc,
			Group: v.Group,
		})
	}
	return resAuditPlans
}

func ConvertReqToAuditPlanFilter(fs []Filter) []auditplan.Filter {
	filters := make([]auditplan.Filter, 0, len(fs))
	for _, v := range fs {
		filters = append(filters, auditplan.Filter{
			Name:                  v.Name,
			FilterComparisonValue: v.FilterComparisonValue,
			FilterBetweenValue: auditplan.FilterBetweenValue{
				From: v.FilterBetweenValue.From,
				To:   v.FilterBetweenValue.To,
			},
		})
	}
	return filters
}

// GetAuditPlanSqlAnalysisData get SQL explain and related table metadata for analysis
// @Summary 获取扫描任务相关的SQL执行计划和表元数据
// @Description get SQL explain and related table metadata for analysis
// @Id getAuditPlanSqlAnalysisDataV1
// @Tags instance_audit_plan
// @Param project_name path string true "project name"
// @Param instance_audit_plan_id path string true "instance audit plan id"
// @Param id path string true "audit plan sql id"
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetSqlManageSqlAnalysisResp
// @router /v1/projects/{project_name}/instance_audit_plans/{instance_audit_plan_id}/sqls/{id}/analysis [get]
func GetAuditPlanSqlAnalysisData(c echo.Context) error {
	insAuditPlanID := c.Param("instance_audit_plan_id")
	sqlManageRecordId := c.Param("id")
	projectUID, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	detail, exist, err := GetInstanceAuditPlanIfCurrentUserCanAccess(c, projectUID, insAuditPlanID, v1.OpPermissionTypeViewOtherAuditPlan)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewInstanceAuditPlanNotExistErr())
	}
	instance, exist, err := dms.GetInstancesById(c.Request().Context(), strconv.FormatUint(detail.InstanceID, 10))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewInstanceNoExistErr())
	}
	s := model.GetStorage()
	originSQL, exist, err := s.GetManageSQLById(sqlManageRecordId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, err)
	}

	res, err := GetSQLAnalysisResult(log.NewEntry(), instance, originSQL.SchemaName, originSQL.SqlText)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetSqlManageSqlAnalysisResp{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertSQLAnalysisResultToRes(res, originSQL.SqlText),
	})
}

// @Summary 扫描任务触发sql审核
// @Description audit plan trigger sql audit
// @Id auditPlanTriggerSqlAuditV1
// @Tags instance_audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instance_audit_plan_id path string true "instance audit plan id"
// @Param audit_plan_id path string true "audit plan id"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/instance_audit_plans/{instance_audit_plan_id}/audit_plans/{audit_plan_id}/audit [post]
func AuditPlanTriggerSqlAudit(c echo.Context) error {
	insAuditPlanID := c.Param("instance_audit_plan_id")
	auditPlanID, err := strconv.Atoi(c.Param("audit_plan_id"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf("parse audit plan report id failed: %v", err)))
	}
	projectUID, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	_, exist, err := GetInstanceAuditPlanIfCurrentUserCanAccess(c, projectUID, insAuditPlanID, v1.OpPermissionTypeSaveAuditPlan)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewInstanceAuditPlanNotExistErr())
	}
	s := model.GetStorage()
	auditPlanSqls, err := s.GetManagerSQLListByAuditPlanId(uint(auditPlanID))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	auditedSqls, err := auditplan.BatchAuditSQLs(auditPlanSqls, false)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = s.Save(auditedSqls)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return controller.JSONBaseErrorReq(c, nil)
}
