package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/common"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/notification"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/actiontech/sqle/sqle/server/auditplan"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/labstack/echo/v4"
	dry "github.com/ungerik/go-dry"
)

var tokenExpire = 365 * 24 * time.Hour

var (
	errAuditPlanNotExist         = errors.New(errors.DataNotExist, fmt.Errorf("audit plan is not exist"))
	errAuditPlanExisted          = errors.New(errors.DataNotExist, fmt.Errorf("audit plan existed"))
	errAuditPlanInstanceConflict = errors.New(errors.DataConflict, fmt.Errorf("instance_name can not be empty while instance_database is not empty"))
)

type GetAuditPlanMetasReqV1 struct {
	FilterInstanceType *string `json:"filter_instance_type" query:"filter_instance_type"`
}

type GetAuditPlanMetasResV1 struct {
	controller.BaseRes
	Data []AuditPlanMetaV1 `json:"data"`
}

type AuditPlanMetaV1 struct {
	Type         string                `json:"audit_plan_type"`
	Desc         string                `json:"audit_plan_type_desc"`
	InstanceType string                `json:"instance_type"`
	Params       []AuditPlanParamResV1 `json:"audit_plan_params,omitempty"`
}

type AuditPlanParamResV1 struct {
	Key   string `json:"key"`
	Desc  string `json:"desc"`
	Value string `json:"value"`
	Type  string `json:"type" enums:"string,int,bool"`
}

func convertAuditPlanMetaToRes(meta auditplan.Meta) AuditPlanMetaV1 {
	res := AuditPlanMetaV1{
		Type:         meta.Type,
		Desc:         meta.Desc,
		InstanceType: meta.InstanceType,
	}
	if meta.Params != nil && len(meta.Params) > 0 {
		paramsRes := make([]AuditPlanParamResV1, 0, len(meta.Params))
		for _, p := range meta.Params {
			paramRes := AuditPlanParamResV1{
				Key:   p.Key,
				Desc:  p.Desc,
				Type:  string(p.Type),
				Value: p.Value,
			}
			paramsRes = append(paramsRes, paramRes)
		}
		res.Params = paramsRes
	}
	return res
}

// @Summary 获取扫描任务元信息
// @Description get audit plan metas
// @Id getAuditPlanMetasV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param filter_instance_type query string false "filter instance type"
// @Success 200 {object} v1.GetAuditPlanMetasResV1
// @router /v1/audit_plan_metas [get]
func GetAuditPlanMetas(c echo.Context) error {
	req := new(GetAuditPlanMetasReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	var metas []AuditPlanMetaV1
	for _, meta := range auditplan.Metas {
		// filter instance type
		if req.FilterInstanceType == nil ||
			meta.InstanceType == auditplan.InstanceTypeAll ||
			meta.InstanceType == *req.FilterInstanceType {
			metas = append(metas, convertAuditPlanMetaToRes(meta))
		}
	}
	return c.JSON(http.StatusOK, &GetAuditPlanMetasResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    metas,
	})
}

type AuditPlanTypesV1 struct {
	Type string `json:"type"`
	Desc string `json:"desc"`
}

type GetAuditPlanTypesResV1 struct {
	controller.BaseRes
	Data []AuditPlanTypesV1 `json:"data"`
}

// GetAuditPlanTypes
// @Summary 获取扫描任务类型
// @Description get audit plan types
// @Id getAuditPlanTypesV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetAuditPlanTypesResV1
// @router /v1/audit_plan_types [get]
func GetAuditPlanTypes(c echo.Context) error {
	auditPlanTypesV1 := make([]AuditPlanTypesV1, 0, len(auditplan.Metas))
	for _, meta := range auditplan.Metas {
		auditPlanTypesV1 = append(auditPlanTypesV1, AuditPlanTypesV1{
			Type: meta.Type,
			Desc: meta.Desc,
		})
	}

	return c.JSON(http.StatusOK, &GetAuditPlanTypesResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    auditPlanTypesV1,
	})
}

type CreateAuditPlanReqV1 struct {
	Name             string                `json:"audit_plan_name" form:"audit_plan_name" example:"audit_plan_for_java_repo_1" valid:"required,name"`
	Cron             string                `json:"audit_plan_cron" form:"audit_plan_cron" example:"0 */2 * * *" valid:"required,cron"`
	InstanceType     string                `json:"audit_plan_instance_type" form:"audit_plan_instance_type" example:"mysql" valid:"required"`
	InstanceName     string                `json:"audit_plan_instance_name" form:"audit_plan_instance_name" example:"test_mysql"`
	InstanceDatabase string                `json:"audit_plan_instance_database" form:"audit_plan_instance_database" example:"app1"`
	RuleTemplateName string                `json:"rule_template_name" from:"rule_template_name" example:"default_MySQL"`
	Type             string                `json:"audit_plan_type" form:"audit_plan_type" example:"slow log"`
	Params           []AuditPlanParamReqV1 `json:"audit_plan_params" valid:"dive,required"`
}

type AuditPlanParamReqV1 struct {
	Key   string `json:"key" form:"key" valid:"required"`
	Value string `json:"value" form:"value" valid:"required"`
}

func checkAndGenerateAuditPlanParams(auditPlanType, instanceType string, paramsReq []AuditPlanParamReqV1) (params.Params, error) {
	meta, err := auditplan.GetMeta(auditPlanType)
	if err != nil {
		return nil, err
	}
	if meta.InstanceType != auditplan.InstanceTypeAll && meta.InstanceType != instanceType {
		return nil, fmt.Errorf("audit plan type %s not found", auditPlanType)
	}
	// check request params is equal params.
	if len(paramsReq) != len(meta.Params) {
		reqParamsKey := make([]string, 0, len(paramsReq))
		for _, p := range paramsReq {
			reqParamsKey = append(reqParamsKey, p.Key)
		}
		paramsKey := make([]string, 0, len(meta.Params))
		for _, p := range meta.Params {
			paramsKey = append(paramsKey, p.Key)
		}
		return nil, fmt.Errorf("request params key is [%s], but need [%s]",
			strings.Join(reqParamsKey, ", "), strings.Join(paramsKey, ", "))
	}
	for _, p := range paramsReq {
		// set and valid param.
		err := meta.Params.SetParamValue(p.Key, p.Value)
		if err != nil {
			return nil, fmt.Errorf("set param error: %s", err)
		}
	}
	return meta.Params, nil
}

// @Summary 添加扫描任务
// @Description create audit plan
// @Id createAuditPlanV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Accept json
// @Param project_name path string true "project name"
// @Param audit_plan body v1.CreateAuditPlanReqV1 true "create audit plan"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/audit_plans [post]
func CreateAuditPlan(c echo.Context) error {
	s := model.GetStorage()

	req := new(CreateAuditPlanReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if !dry.StringInSlice(req.InstanceType, driver.AllDrivers()) {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DriverNotExist, &driver.DriverNotSupportedError{DriverTyp: req.InstanceType}))
	}

	if req.InstanceDatabase != "" && req.InstanceName == "" {
		return controller.JSONBaseErrorReq(c, errAuditPlanInstanceConflict)
	}

	// check project
	projectName := c.Param("project_name")
	project, exist, err := s.GetProjectByName(projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errProjectNotExist)
	}

	// check user
	currentUserName := controller.GetUserName(c)

	err = CheckIsProjectMember(currentUserName, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	user, exist, err := s.GetUserByName(currentUserName)
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("user is not exist")))
	} else if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
	}

	// check audit plan name
	_, exist, err = s.GetAuditPlanFromProjectByName(projectName, req.Name)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if exist {
		return controller.JSONBaseErrorReq(c, errAuditPlanExisted)
	}

	// check instance
	var instanceType string
	if req.InstanceName != "" {
		inst, exist, err := s.GetInstanceByNameAndProjectName(req.InstanceName, projectName)
		if !exist {
			return controller.JSONBaseErrorReq(c, ErrInstanceNotExist)
		} else if err != nil {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
		}
		// check instance database
		if req.InstanceDatabase != "" {
			drvMgr, err := common.NewDriverManagerWithoutAudit(log.NewEntry(), inst, "")
			if err != nil {
				return controller.JSONBaseErrorReq(c, err)
			}
			defer drvMgr.Close(context.TODO())

			d, err := drvMgr.GetAuditDriver()
			if err != nil {
				return controller.JSONBaseErrorReq(c, err)
			}

			schemas, err := d.Schemas(context.TODO())
			if err != nil {
				return controller.JSONBaseErrorReq(c, err)
			}
			if !dry.StringInSlice(req.InstanceDatabase, schemas) {
				return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("database %v is not exist in instance", req.InstanceDatabase)))
			}
		}
		instanceType = inst.DbType

		// check operation
		can, err := s.CheckUserCanCreateAuditPlan(user, req.InstanceName, instanceType) // todo: refactor permissions.
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if !can {
			return controller.JSONBaseErrorReq(c, errors.NewUserNotPermissionError(model.GetOperationCodeDesc(uint(model.OP_AUDIT_PLAN_SAVE))))
		}

	} else {
		instanceType = req.InstanceType
	}

	// todo: check and select rule template binding in project.
	// check rule template name
	if req.RuleTemplateName != "" {
		exist, err = s.IsRuleTemplateExist(req.RuleTemplateName)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if !exist {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("rule template does not exist")))
		}
	}
	ruleTemplateName, err := autoSelectRuleTemplate(req.RuleTemplateName, req.InstanceName, req.InstanceType)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// check params
	if req.Type == "" {
		req.Type = auditplan.TypeDefault
	}
	ps, err := checkAndGenerateAuditPlanParams(req.Type, instanceType, req.Params)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
	}

	// generate token
	j := utils.NewJWT(utils.JWTSecretKey)
	t, err := j.CreateToken(currentUserName, time.Now().Add(tokenExpire).Unix(),
		utils.WithAuditPlanName(req.Name))
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
	}

	ap := &model.AuditPlan{
		Name:             req.Name,
		CronExpression:   req.Cron,
		Type:             req.Type,
		Params:           ps,
		CreateUserID:     user.ID,
		Token:            t,
		DBType:           instanceType,
		RuleTemplateName: ruleTemplateName,
		InstanceName:     req.InstanceName,
		InstanceDatabase: req.InstanceDatabase,
		ProjectId:        project.ID,
	}
	err = s.Save(ap)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	manager := auditplan.GetManager()
	return controller.JSONBaseErrorReq(c, manager.SyncTask(ap.ID))
}

// customRuleTemplateName如果为空, 将返回instanceName绑定的规则模板, 如果customRuleTemplateName,和instanceName都为空, 将返回dbType对应默认模板, dbType不能为空, 函数不做参数校验
// 规则模板选择规则: 指定规则模板 -- > 数据源绑定的规则模板 -- > 数据库类型默认模板
func autoSelectRuleTemplate(customRuleTemplateName string, instanceName string, dbType string) (ruleTemplateName string, err error) {
	s := model.GetStorage()

	if customRuleTemplateName != "" {
		return customRuleTemplateName, nil
	}

	if instanceName != "" {
		ruleTemplate, exist, err := s.GetRuleTemplatesByInstanceName(instanceName)
		if err != nil {
			return "", err
		}
		if exist {
			return ruleTemplate.Name, nil
		}
	}

	return s.GetDefaultRuleTemplateName(dbType), nil

}

// @Summary 删除扫描任务
// @Description delete audit plan
// @Id deleteAuditPlanV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param audit_plan_name path string true "audit plan name"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/audit_plans/{audit_plan_name}/ [delete]
func DeleteAuditPlan(c echo.Context) error {
	s := model.GetStorage()
	projectName := c.Param("project_name")
	userName := controller.GetUserName(c)
	err := CheckIsProjectMember(userName, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	apName := c.Param("audit_plan_name")
	err = CheckCurrentUserCanAccessAuditPlan(c, apName, 0) // todo: refactor permissions.
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	ap, exist, err := s.GetAuditPlanFromProjectByName(projectName, apName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errAuditPlanNotExist)
	}
	err = s.Delete(ap)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	manager := auditplan.GetManager()
	return controller.JSONBaseErrorReq(c, manager.SyncTask(ap.ID))
}

type UpdateAuditPlanReqV1 struct {
	Cron             *string               `json:"audit_plan_cron" form:"audit_plan_cron" example:"0 */2 * * *" valid:"omitempty,cron"`
	InstanceName     *string               `json:"audit_plan_instance_name" form:"audit_plan_instance_name" example:"test_mysql"`
	InstanceDatabase *string               `json:"audit_plan_instance_database" form:"audit_plan_instance_database" example:"app1"`
	RuleTemplateName *string               `json:"rule_template_name" form:"rule_template_name" example:"default_MySQL"`
	Params           []AuditPlanParamReqV1 `json:"audit_plan_params" valid:"dive,required"`
}

// @Summary 更新扫描任务
// @Description update audit plan
// @Id updateAuditPlanV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param audit_plan_name path string true "audit plan name"
// @param audit_plan body v1.UpdateAuditPlanReqV1 true "update audit plan"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/audit_plans/{audit_plan_name}/ [patch]
func UpdateAuditPlan(c echo.Context) error {
	req := new(UpdateAuditPlanReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	projectName := c.Param("project_name")
	userName := controller.GetUserName(c)
	err := CheckIsProjectMember(userName, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	apName := c.Param("audit_plan_name")

	err = CheckCurrentUserCanAccessAuditPlan(c, apName, 0) // todo: refactor permissions.
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	storage := model.GetStorage()
	ap, exist, err := storage.GetAuditPlanFromProjectByName(projectName, apName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errAuditPlanNotExist)
	}

	updateAttr := make(map[string]interface{})
	if req.Cron != nil {
		updateAttr["cron_expression"] = *req.Cron
	}
	if req.InstanceName != nil {
		updateAttr["instance_name"] = *req.InstanceName
	}
	if req.InstanceDatabase != nil {
		updateAttr["instance_database"] = *req.InstanceDatabase
	}
	if req.RuleTemplateName != nil {
		exist, err = storage.IsRuleTemplateExist(*req.RuleTemplateName)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if !exist {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("rule template does not exist")))
		}
		updateAttr["rule_template_name"] = *req.RuleTemplateName
	}
	if req.Params != nil {
		ps, err := checkAndGenerateAuditPlanParams(ap.Type, ap.DBType, req.Params)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		updateAttr["params"] = ps
	}

	err = storage.UpdateAuditPlanById(ap.ID, updateAttr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	manager := auditplan.GetManager()
	return controller.JSONBaseErrorReq(c, manager.SyncTask(ap.ID))
}

type GetAuditPlansReqV1 struct {
	FilterAuditPlanDBType       string `json:"filter_audit_plan_db_type" query:"filter_audit_plan_db_type"`
	FuzzySearchAuditPlanName    string `json:"fuzzy_search_audit_plan_name" query:"fuzzy_search_audit_plan_name"`
	FilterAuditPlanType         string `json:"filter_audit_plan_type" query:"filter_audit_plan_type"`
	FilterAuditPlanInstanceName string `json:"filter_audit_plan_instance_name" query:"filter_audit_plan_instance_name"`
	PageIndex                   uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize                    uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetAuditPlansResV1 struct {
	controller.BaseRes
	Data      []AuditPlanResV1 `json:"data"`
	TotalNums uint64           `json:"total_nums"`
}

type AuditPlanResV1 struct {
	Name             string          `json:"audit_plan_name" example:"audit_for_java_app1"`
	Cron             string          `json:"audit_plan_cron" example:"0 */2 * * *"`
	DBType           string          `json:"audit_plan_db_type" example:"mysql"`
	Token            string          `json:"audit_plan_token" example:"it's a JWT Token for scanner"`
	InstanceName     string          `json:"audit_plan_instance_name" example:"test_mysql"`
	InstanceDatabase string          `json:"audit_plan_instance_database" example:"app1"`
	RuleTemplateName string          `json:"rule_template_name" example:"default_MySQL"`
	Meta             AuditPlanMetaV1 `json:"audit_plan_meta"`
}

// GetAuditPlans
// @Summary 获取扫描任务信息列表
// @Description get audit plan info list
// @Id getAuditPlansV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param filter_audit_plan_db_type query string false "filter audit plan db type"
// @Param fuzzy_search_audit_plan_name query string false "fuzzy search audit plan name"
// @Param filter_audit_plan_type query string false "filter audit plan type"
// @Param filter_audit_plan_instance_name query string false "filter audit plan instance name"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v1.GetAuditPlansResV1
// @router /v1/projects/{project_name}/audit_plans [get]
func GetAuditPlans(c echo.Context) error {
	s := model.GetStorage()

	req := new(GetAuditPlansReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	projectName := c.Param("project_name")
	userName := controller.GetUserName(c)
	err := CheckIsProjectMember(userName, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var offset uint32
	if req.PageIndex >= 1 {
		offset = req.PageSize * (req.PageIndex - 1)
	}

	currentUser, err := controller.GetCurrentUser(c)
	if err != nil {
		return err
	}
	instances, err := s.GetUserCanOpInstances(currentUser, []uint{model.OP_AUDIT_PLAN_VIEW_OTHERS}) // todo: refactor permissions.
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	names := []string{}
	for _, instance := range instances {
		names = append(names, instance.Name)
	}
	data := map[string]interface{}{
		"filter_audit_plan_db_type":       req.FilterAuditPlanDBType,
		"fuzzy_search_audit_plan_name":    req.FuzzySearchAuditPlanName,
		"filter_audit_plan_type":          req.FilterAuditPlanType,
		"filter_audit_plan_instance_name": req.FilterAuditPlanInstanceName,
		"current_user_name":               currentUser.Name,
		"current_user_is_admin":           model.DefaultAdminUser == currentUser.Name,
		"filter_project_name":             projectName,
		"limit":                           req.PageSize,
		"offset":                          offset,
	}
	if len(names) > 0 {
		data["accessible_instances_name"] = fmt.Sprintf("'%s'", strings.Join(names, "', '"))
	}
	auditPlans, count, err := s.GetAuditPlansByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	auditPlansResV1 := make([]AuditPlanResV1, len(auditPlans))
	for i, ap := range auditPlans {
		meta, err := auditplan.GetMeta(ap.Type.String)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		meta.Params = ap.Params
		auditPlansResV1[i] = AuditPlanResV1{
			Name:             ap.Name,
			Cron:             ap.Cron,
			DBType:           ap.DBType,
			InstanceName:     ap.InstanceName,
			InstanceDatabase: ap.InstanceDatabase,
			RuleTemplateName: ap.RuleTemplateName.String,
			Token:            ap.Token,
			Meta:             convertAuditPlanMetaToRes(meta),
		}
	}
	return c.JSON(http.StatusOK, &GetAuditPlansResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      auditPlansResV1,
		TotalNums: count,
	})
}

type GetAuditPlanResV1 struct {
	controller.BaseRes
	Data AuditPlanResV1 `json:"data"`
}

// @Summary 获取指定扫描任务
// @Description get audit plan
// @Id getAuditPlanV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param audit_plan_name path string true "audit plan name"
// @Success 200 {object} v1.GetAuditPlanResV1
// @router /v1/projects/{project_name}/audit_plans/{audit_plan_name}/ [get]
func GetAuditPlan(c echo.Context) error {
	projectName := c.Param("project_name")
	userName := controller.GetUserName(c)
	err := CheckIsProjectMember(userName, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	apName := c.Param("audit_plan_name")
	err = CheckCurrentUserCanAccessAuditPlan(c, apName, model.OP_AUDIT_PLAN_VIEW_OTHERS) // todo: refactor permissions.
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	storage := model.GetStorage()

	ap, exist, err := storage.GetAuditPlanFromProjectByName(projectName, apName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errAuditPlanNotExist)
	}
	meta, err := auditplan.GetMeta(ap.Type)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	meta.Params = ap.Params

	return c.JSON(http.StatusOK, &GetAuditPlanResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: AuditPlanResV1{
			Name:             ap.Name,
			Cron:             ap.CronExpression,
			DBType:           ap.DBType,
			InstanceName:     ap.InstanceName,
			InstanceDatabase: ap.InstanceDatabase,
			RuleTemplateName: ap.RuleTemplateName,
			Token:            ap.Token,
			Meta:             convertAuditPlanMetaToRes(meta),
		},
	})
}

type GetAuditPlanReportsReqV1 struct {
	PageIndex uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetAuditPlanReportsResV1 struct {
	controller.BaseRes
	Data      []AuditPlanReportResV1 `json:"data"`
	TotalNums uint64                 `json:"total_nums"`
}

type AuditPlanReportResV1 struct {
	Id         string  `json:"audit_plan_report_id" example:"1"`
	AuditLevel string  `json:"audit_level" enums:"normal,notice,warn,error,"`
	Score      int32   `json:"score"`
	PassRate   float64 `json:"pass_rate"`
	Timestamp  string  `json:"audit_plan_report_timestamp" example:"RFC3339"`
}

// @Summary 获取指定扫描任务的报告列表
// @Description get audit plan report list
// @Id getAuditPlanReportsV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param audit_plan_name path string true "audit plan name"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v1.GetAuditPlanReportsResV1
// @router /v1/projects/{project_name}/audit_plans/{audit_plan_name}/reports [get]
func GetAuditPlanReports(c echo.Context) error {
	s := model.GetStorage()

	req := new(GetAuditPlanReportsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	projectName := c.Param("project_name")
	userName := controller.GetUserName(c)
	err := CheckIsProjectMember(userName, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	apName := c.Param("audit_plan_name")
	err = CheckCurrentUserCanAccessAuditPlan(c, apName, model.OP_AUDIT_PLAN_VIEW_OTHERS) // todo: refactor permissions.
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var offset uint32
	if req.PageIndex >= 1 {
		offset = req.PageSize * (req.PageIndex - 1)
	}

	data := map[string]interface{}{
		"project_name":    projectName,
		"audit_plan_name": apName,
		"limit":           req.PageSize,
		"offset":          offset,
	}
	auditPlanReports, count, err := s.GetAuditPlanReportsByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	auditPlanReportsResV1 := make([]AuditPlanReportResV1, len(auditPlanReports))
	for i, auditPlanReport := range auditPlanReports {
		auditPlanReportsResV1[i] = AuditPlanReportResV1{
			Id:         auditPlanReport.ID,
			AuditLevel: auditPlanReport.AuditLevel.String,
			Score:      auditPlanReport.Score.Int32,
			PassRate:   auditPlanReport.PassRate.Float64,
			Timestamp:  auditPlanReport.CreateAt,
		}
	}
	return c.JSON(http.StatusOK, &GetAuditPlanReportsResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      auditPlanReportsResV1,
		TotalNums: count,
	})
}

type GetAuditPlanReportResV1 struct {
	controller.BaseRes
	Data AuditPlanReportResV1 `json:"data"`
}

// @Summary 获取指定扫描任务的SQL扫描记录统计信息
// @Description get audit plan report
// @Id getAuditPlanReportV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param audit_plan_name path string true "audit plan name"
// @Param audit_plan_report_id path string true "audit plan report id"
// @Success 200 {object} v1.GetAuditPlanReportResV1
// @router /v1/projects/{project_name}/audit_plans/{audit_plan_name}/reports/{audit_plan_report_id}/ [get]
func GetAuditPlanReport(c echo.Context) error {
	projectName := c.Param("project_name")
	userName := controller.GetUserName(c)
	err := CheckIsProjectMember(userName, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	apName := c.Param("audit_plan_name")
	err = CheckCurrentUserCanAccessAuditPlan(c, apName, model.OP_AUDIT_PLAN_VIEW_OTHERS) // todo: refactor permissions.
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	id := c.Param("audit_plan_report_id")
	reportID, err := strconv.Atoi(id)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf("parse audit plan report id failed: %v", err)))
	}
	s := model.GetStorage()
	report, exist, err := s.GetAuditPlanReportByID(uint(reportID))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("audit plan report not exist")))
	}

	return c.JSON(http.StatusOK, &GetAuditPlanReportResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: AuditPlanReportResV1{
			Id:         id,
			AuditLevel: report.AuditLevel,
			Score:      report.Score,
			PassRate:   report.PassRate,
			Timestamp:  report.CreatedAt.Format(time.RFC3339),
		},
	})
}

type FullSyncAuditPlanSQLsReqV1 struct {
	SQLs []AuditPlanSQLReqV1 `json:"audit_plan_sql_list" form:"audit_plan_sql_list" valid:"dive"`
}

type AuditPlanSQLReqV1 struct {
	Fingerprint          string `json:"audit_plan_sql_fingerprint" form:"audit_plan_sql_fingerprint" example:"select * from t1 where id = ?"`
	Counter              string `json:"audit_plan_sql_counter" form:"audit_plan_sql_counter" example:"6" valid:"required"`
	LastReceiveText      string `json:"audit_plan_sql_last_receive_text" form:"audit_plan_sql_last_receive_text" example:"select * from t1 where id = 1"`
	LastReceiveTimestamp string `json:"audit_plan_sql_last_receive_timestamp" form:"audit_plan_sql_last_receive_timestamp" example:"RFC3339"`
	Schema               string `json:"audit_plan_sql_schema" from:"audit_plan_sql_schema" example:"db1"`
}

// @Summary 全量同步SQL到扫描任务
// @Description full sync audit plan SQLs
// @Id fullSyncAuditPlanSQLsV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param audit_plan_name path string true "audit plan name"
// @Param sqls body v1.FullSyncAuditPlanSQLsReqV1 true "full sync audit plan SQLs request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/audit_plans/{audit_plan_name}/sqls/full [post]
func FullSyncAuditPlanSQLs(c echo.Context) error {
	req := new(FullSyncAuditPlanSQLsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	projectName := c.Param("project_name")
	apName := c.Param("audit_plan_name")
	sqls, err := convertToModelAuditPlanSQL(c, projectName, apName, req.SQLs)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// TODO: just get ap id.
	s := model.GetStorage()
	ap, exist, err := s.GetAuditPlanFromProjectByName(projectName, apName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errAuditPlanNotExist)
	}

	manager := auditplan.GetManager()
	return controller.JSONBaseErrorReq(c, manager.UploadSQLs(ap.ID, sqls, false))
}

type PartialSyncAuditPlanSQLsReqV1 struct {
	SQLs []AuditPlanSQLReqV1 `json:"audit_plan_sql_list" form:"audit_plan_sql_list" valid:"dive"`
}

// @Summary 增量同步SQL到扫描任务
// @Description partial sync audit plan SQLs
// @Id partialSyncAuditPlanSQLsV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param audit_plan_name path string true "audit plan name"
// @Param sqls body v1.PartialSyncAuditPlanSQLsReqV1 true "partial sync audit plan SQLs request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/audit_plans/{audit_plan_name}/sqls/partial [post]
func PartialSyncAuditPlanSQLs(c echo.Context) error {
	req := new(PartialSyncAuditPlanSQLsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	projectName := c.Param("project_name")
	apName := c.Param("audit_plan_name")
	sqls, err := convertToModelAuditPlanSQL(c, projectName, apName, req.SQLs)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// TODO: just get ap id.
	s := model.GetStorage()
	ap, exist, err := s.GetAuditPlanFromProjectByName(projectName, apName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errAuditPlanNotExist)
	}

	manager := auditplan.GetManager()
	return controller.JSONBaseErrorReq(c, manager.UploadSQLs(ap.ID, sqls, true))
}

func convertToModelAuditPlanSQL(c echo.Context, projectName, apName string, reqSQLs []AuditPlanSQLReqV1) ([]*auditplan.SQL, error) {
	s := model.GetStorage()

	ap, exist, err := s.GetAuditPlanFromProjectByName(projectName, apName)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errAuditPlanNotExist
	}

	var drvMgr driver.DriverManager
	var drv driver.Driver
	// lazy load driver
	initDriver := func() error {
		if drv == nil {
			drvMgr, err = common.NewDriverManagerWithoutCfg(log.NewEntry(), ap.DBType)
			if err != nil {
				return err
			}
			drv, err = drvMgr.GetAuditDriver()
			if err != nil {
				return err
			}
		}
		return nil
	}
	defer func() {
		if drvMgr != nil {
			drvMgr.Close(context.TODO())
		}
	}()

	sqls := make([]*auditplan.SQL, len(reqSQLs))
	for i, reqSQL := range reqSQLs {
		fp := reqSQL.Fingerprint
		// the caller may be written in a different language, such as (Java, Bash, Python), so the fingerprint is
		// generated in different ways. In order to maintain th same fingerprint generation logic, we provide a way to
		// generate it by sqle, if the request fingerprint is empty.
		if fp == "" {
			err := initDriver()
			if err != nil {
				return nil, err
			}
			nodes, err := drv.Parse(context.TODO(), reqSQL.LastReceiveText)
			if err != nil {
				return nil, err
			}
			if len(nodes) > 0 {
				fp = nodes[0].Fingerprint
			} else {
				fp = reqSQL.LastReceiveText
			}
		}
		counter, err := strconv.ParseUint(reqSQL.Counter, 10, 64)
		if err != nil {
			return nil, err
		}
		info := map[string]interface{}{
			"counter":                counter,
			"last_receive_timestamp": reqSQL.LastReceiveTimestamp,
			server.AuditSchema:       reqSQL.Schema,
		}
		sqls[i] = &auditplan.SQL{
			Fingerprint: fp,
			SQLContent:  reqSQL.LastReceiveText,
			Info:        info,
		}
	}
	return sqls, nil
}

type TriggerAuditPlanResV1 struct {
	controller.BaseRes
	Data AuditPlanReportResV1 `json:"data"`
}

// @Summary 触发扫描任务
// @Description trigger audit plan
// @Id triggerAuditPlanV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param audit_plan_name path string true "audit plan name"
// @Success 200 {object} v1.TriggerAuditPlanResV1
// @router /v1/projects/{project_name}/audit_plans/{audit_plan_name}/trigger [post]
func TriggerAuditPlan(c echo.Context) error {
	projectName := c.Param("project_name")
	userName := controller.GetUserName(c)
	err := CheckIsProjectMember(userName, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	apName := c.Param("audit_plan_name")
	err = CheckCurrentUserCanAccessAuditPlan(c, apName, 0) // todo: refactor permissions.
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// TODO: just get ap id.
	s := model.GetStorage()
	ap, exist, err := s.GetAuditPlanFromProjectByName(projectName, apName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errAuditPlanNotExist)
	}

	manager := auditplan.GetManager()
	report, err := manager.Audit(ap.ID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, &TriggerAuditPlanResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: AuditPlanReportResV1{
			Id:         fmt.Sprintf("%v", report.ID),
			AuditLevel: report.AuditLevel,
			Score:      report.Score,
			PassRate:   report.PassRate,
			Timestamp:  report.CreatedAt.Format(time.RFC3339),
		},
	})
}

func CheckCurrentUserCanAccessAuditPlan(c echo.Context, apName string, opCode int) error {
	storage := model.GetStorage()

	ap, exist, err := storage.GetAuditPlanByName(apName)
	if err != nil {
		return err
	}
	if !exist {
		return errAuditPlanNotExist
	}

	if controller.GetUserName(c) == model.DefaultAdminUser {
		return nil
	}

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return err
	}

	if ap.CreateUserID == user.ID {
		return nil
	}
	if opCode > 0 {
		instances, err := storage.GetUserCanOpInstances(user, []uint{uint(opCode)})
		if err != nil {
			return controller.JSONBaseErrorReq(c, errors.NewUserNotPermissionError(model.GetOperationCodeDesc(uint(opCode))))
		}
		for _, instance := range instances {
			if ap.InstanceName == instance.Name {
				return nil
			}
		}
	}
	return nil
}

type UpdateAuditPlanNotifyConfigReqV1 struct {
	NotifyInterval      *int    `json:"notify_interval" default:"10"`
	NotifyLevel         *string `json:"notify_level" default:"warn" enums:"normal,notice,warn,error" valid:"oneof=normal notice warn error"`
	EnableEmailNotify   *bool   `json:"enable_email_notify"`
	EnableWebHookNotify *bool   `json:"enable_web_hook_notify"`
	WebHookURL          *string `json:"web_hook_url"`
	WebHookTemplate     *string `json:"web_hook_template"`
}

// @Summary 更新扫描任务通知设置
// @Description update audit plan notify config
// @Id updateAuditPlanNotifyConfigV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param audit_plan_name path string true "audit plan name"
// @param config body v1.UpdateAuditPlanNotifyConfigReqV1 true "update audit plan notify config"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/audit_plans/{audit_plan_name}/notify_config [patch]
func UpdateAuditPlanNotifyConfig(c echo.Context) error {
	req := new(UpdateAuditPlanNotifyConfigReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	projectName := c.Param("project_name")
	userName := controller.GetUserName(c)
	err := CheckIsProjectMember(userName, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	apName := c.Param("audit_plan_name")

	err = CheckCurrentUserCanAccessAuditPlan(c, apName, 0) // todo: refactor permissions.
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	updateAttr := make(map[string]interface{})
	if req.EnableWebHookNotify != nil {
		updateAttr["enable_web_hook_notify"] = *req.EnableWebHookNotify
	}
	if req.EnableEmailNotify != nil {
		updateAttr["enable_email_notify"] = *req.EnableEmailNotify
	}
	if req.NotifyInterval != nil {
		updateAttr["notify_interval"] = *req.NotifyInterval
	}
	if req.NotifyLevel != nil {
		updateAttr["notify_level"] = *req.NotifyLevel
	}
	if req.WebHookURL != nil {
		updateAttr["web_hook_url"] = *req.WebHookURL
	}
	if req.WebHookTemplate != nil {
		updateAttr["web_hook_template"] = *req.WebHookTemplate
	}

	// TODO: just get ap id.
	s := model.GetStorage()
	ap, exist, err := s.GetAuditPlanFromProjectByName(projectName, apName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errAuditPlanNotExist)
	}

	storage := model.GetStorage()
	err = storage.UpdateAuditPlanById(ap.ID, updateAttr)
	return controller.JSONBaseErrorReq(c, err)
}

type GetAuditPlanNotifyConfigResV1 struct {
	controller.BaseRes
	Data GetAuditPlanNotifyConfigResDataV1 `json:"data"`
}

type GetAuditPlanNotifyConfigResDataV1 struct {
	NotifyInterval      int    `json:"notify_interval"`
	NotifyLevel         string `json:"notify_level"`
	EnableEmailNotify   bool   `json:"enable_email_notify"`
	EnableWebHookNotify bool   `json:"enable_web_hook_notify"`
	WebHookURL          string `json:"web_hook_url"`
	WebHookTemplate     string `json:"web_hook_template"`
}

// @Summary 获取扫描任务消息推送设置
// @Description get audit plan notify config
// @Id getAuditPlanNotifyConfigV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param audit_plan_name path string true "audit plan name"
// @Success 200 {object} v1.GetAuditPlanNotifyConfigResV1
// @router /v1/projects/{project_name}/audit_plans/{audit_plan_name}/notify_config [get]
func GetAuditPlanNotifyConfig(c echo.Context) error {
	projectName := c.Param("project_name")
	userName := controller.GetUserName(c)
	err := CheckIsProjectMember(userName, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	apName := c.Param("audit_plan_name")
	err = CheckCurrentUserCanAccessAuditPlan(c, apName, model.OP_AUDIT_PLAN_VIEW_OTHERS) // todo: refactor permissions.
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	ap, _, err := s.GetAuditPlanFromProjectByName(projectName, apName)

	return c.JSON(http.StatusOK, GetAuditPlanNotifyConfigResV1{
		BaseRes: controller.NewBaseReq(err),
		Data: GetAuditPlanNotifyConfigResDataV1{
			NotifyInterval:      ap.NotifyInterval,
			NotifyLevel:         ap.NotifyLevel,
			EnableEmailNotify:   ap.EnableEmailNotify,
			EnableWebHookNotify: ap.EnableWebHookNotify,
			WebHookURL:          ap.WebHookURL,
			WebHookTemplate:     ap.WebHookTemplate,
		},
	})
}

type TestAuditPlanNotifyConfigResV1 struct {
	controller.BaseRes
	Data TestAuditPlanNotifyConfigResDataV1 `json:"data"`
}

type TestAuditPlanNotifyConfigResDataV1 struct {
	IsNotifySendNormal bool   `json:"is_notify_send_normal"`
	SendErrorMessage   string `json:"send_error_message,omitempty"`
}

// @Summary 测试扫描任务消息推送
// @Description Test audit task message push
// @Id testAuditPlanNotifyConfigV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param audit_plan_name path string true "audit plan name"
// @Success 200 {object} v1.TestAuditPlanNotifyConfigResV1
// @router /v1/projects/{project_name}/audit_plans/{audit_plan_name}/notify_config/test [get]
func TestAuditPlanNotifyConfig(c echo.Context) error {
	projectName := c.Param("project_name")
	userName := controller.GetUserName(c)
	err := CheckIsProjectMember(userName, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	apName := c.Param("audit_plan_name")
	s := model.GetStorage()

	err = CheckCurrentUserCanAccessAuditPlan(c, apName, 0) // todo: refactor permissions.
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	ap, _, err := s.GetAuditPlanFromProjectByName(projectName, apName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	user, exist, err := s.GetUserByID(ap.CreateUserID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return c.JSON(http.StatusOK, TestAuditPlanNotifyConfigResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data: TestAuditPlanNotifyConfigResDataV1{
				IsNotifySendNormal: false,
				SendErrorMessage:   "audit plan create user not exist",
			},
		})
	}
	ap.CreateUser = user

	err = notification.GetAuditPlanNotifier().Send(&notification.TestNotify{}, ap)
	if err != nil {
		return c.JSON(http.StatusOK, TestAuditPlanNotifyConfigResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data: TestAuditPlanNotifyConfigResDataV1{
				IsNotifySendNormal: false,
				SendErrorMessage:   err.Error(),
			},
		})
	}
	return c.JSON(http.StatusOK, TestAuditPlanNotifyConfigResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: TestAuditPlanNotifyConfigResDataV1{
			IsNotifySendNormal: true,
			SendErrorMessage:   "success",
		},
	})
}

type TableMeta struct {
	Name           string       `json:"name"`
	Schema         string       `json:"schema"`
	Columns        TableColumns `json:"columns"`
	Indexes        TableIndexes `json:"indexes"`
	CreateTableSQL string       `json:"create_table_sql"`
	Message        string       `json:"message"`
}

type GetSQLAnalysisDataResItemV1 struct {
	SQLExplain SQLExplain  `json:"sql_explain"`
	TableMetas []TableMeta `json:"table_metas"`
}

type GetAuditPlanAnalysisDataResV1 struct {
	controller.BaseRes
	Data GetSQLAnalysisDataResItemV1 `json:"data"`
}

// GetAuditPlanAnalysisData get SQL explain and related table metadata for analysis
// @Summary 获取task相关的SQL执行计划和表元数据
// @Description get SQL explain and related table metadata for analysis
// @Id getTaskAnalysisData
// @Tags audit_plan
// @Param project_name path string true "project name"
// @Param audit_plan_report_id path string true "audit plan report id"
// @Param number path string true "sql number"
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetAuditPlanAnalysisDataResV1
// @router /v1/projects/{project_name}/audit_plans/reports/{audit_plan_report_id}/sqls/{number}/analysis [get]
func GetAuditPlanAnalysisData(c echo.Context) error {
	return getAuditPlanAnalysisData(c)
}

type GetAuditPlanSQLsReqV1 struct {
	PageIndex uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetAuditPlanSQLsResV1 struct {
	controller.BaseRes
	Data      AuditPlanSQLResV1 `json:"data"`
	TotalNums uint64            `json:"total_nums"`
}

type AuditPlanSQLResV1 struct {
	Head []AuditPlanSQLHeadV1                 `json:"head"`
	Rows []map[string] /* head name */ string `json:"rows"`
}

type AuditPlanSQLHeadV1 struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
	Type string `json:"type,omitempty" enums:"sql"`
}

// @Summary 获取指定扫描任务的SQLs信息(不包括扫描结果)
// @Description get audit plan SQLs
// @Id getAuditPlanSQLsV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param audit_plan_name path string true "audit plan name"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v1.GetAuditPlanSQLsResV1
// @router /v1/projects/{project_name}/audit_plans/{audit_plan_name}/sqls [get]
func GetAuditPlanSQLs(c echo.Context) error {
	req := new(GetAuditPlanSQLsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	projectName := c.Param("project_name")
	userName := controller.GetUserName(c)
	err := CheckIsProjectMember(userName, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	apName := c.Param("audit_plan_name")
	err = CheckCurrentUserCanAccessAuditPlan(c, apName, model.OP_AUDIT_PLAN_VIEW_OTHERS) // todo: refactor permissions.
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var offset uint32
	if req.PageIndex >= 1 {
		offset = req.PageSize * (req.PageIndex - 1)
	}

	data := map[string]interface{}{
		"limit":  req.PageSize,
		"offset": offset,
	}

	// TODO: just get ap id.
	s := model.GetStorage()
	ap, exist, err := s.GetAuditPlanFromProjectByName(projectName, apName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errAuditPlanNotExist)
	}

	manager := auditplan.GetManager()
	head, rows, count, err := manager.GetSQLs(ap.ID, data)
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
			Desc: v.Desc,
			Type: v.Type,
		})
	}
	return c.JSON(http.StatusOK, &GetAuditPlanSQLsResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      res,
		TotalNums: count,
	})
}

type GetAuditPlanReportSQLsReqV1 struct {
	PageIndex uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetAuditPlanReportSQLsResV1 struct {
	controller.BaseRes
	Data      []AuditPlanReportSQLResV1 `json:"data"`
	TotalNums uint64                    `json:"total_nums"`
}

type AuditPlanReportSQLResV1 struct {
	SQL         string `json:"audit_plan_report_sql" example:"select * from t1 where id = 1"`
	AuditResult string `json:"audit_plan_report_sql_audit_result" example:"same format as task audit result"`
	Number      uint   `json:"number" example:"1"`
}

// GetAuditPlanReportSQLsV1 is to fix the irregular uri used by GetAuditPlanReportSQLs
// issue: https://github.com/actiontech/sqle/issues/429
// @Summary 获取指定扫描任务的SQL扫描详情
// @Description get audit plan report SQLs
// @Id getAuditPlanReportsSQLsV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param audit_plan_name path string true "audit plan name"
// @Param audit_plan_report_id path string true "audit plan report id"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v1.GetAuditPlanReportSQLsResV1
// @router /v1/projects/{project_name}/audit_plans/{audit_plan_name}/reports/{audit_plan_report_id}/sqls [get]
func GetAuditPlanReportSQLsV1(c echo.Context) error {
	s := model.GetStorage()

	req := new(GetAuditPlanReportSQLsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	projectName := c.Param("project_name")
	userName := controller.GetUserName(c)
	err := CheckIsProjectMember(userName, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	apName := c.Param("audit_plan_name")
	err = CheckCurrentUserCanAccessAuditPlan(c, apName, model.OP_AUDIT_PLAN_VIEW_OTHERS) // todo: refactor permissions.
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var offset uint32
	if req.PageIndex >= 1 {
		offset = req.PageSize * (req.PageIndex - 1)
	}

	data := map[string]interface{}{
		"audit_plan_report_id": c.Param("audit_plan_report_id"),
		"limit":                req.PageSize,
		"offset":               offset,
	}
	auditPlanReportSQLs, count, err := s.GetAuditPlanReportSQLsByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	auditPlanReportSQLsResV1 := make([]AuditPlanReportSQLResV1, len(auditPlanReportSQLs))
	for i, auditPlanReportSQL := range auditPlanReportSQLs {
		auditPlanReportSQLsResV1[i] = AuditPlanReportSQLResV1{
			SQL:         auditPlanReportSQL.SQL,
			AuditResult: auditPlanReportSQL.AuditResult,
			Number:      auditPlanReportSQL.Number,
		}
	}
	return c.JSON(http.StatusOK, &GetAuditPlanReportSQLsResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      auditPlanReportSQLsResV1,
		TotalNums: count,
	})
}
