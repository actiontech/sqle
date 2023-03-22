package v2

import (
	"fmt"
	"net/http"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/errors"

	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/labstack/echo/v4"
)

type GetInstancesReqV2 struct {
	FilterInstanceName     string `json:"filter_instance_name" query:"filter_instance_name"`
	FilterDBType           string `json:"filter_db_type" query:"filter_db_type"`
	FilterDBHost           string `json:"filter_db_host" query:"filter_db_host"`
	FilterDBPort           string `json:"filter_db_port" query:"filter_db_port"`
	FilterDBUser           string `json:"filter_db_user" query:"filter_db_user"`
	FilterRuleTemplateName string `json:"filter_rule_template_name" query:"filter_rule_template_name"`
	PageIndex              uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize               uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type RuleTemplateV2 struct {
	Name                 string `json:"name"`
	IsGlobalRuleTemplate bool   `json:"is_global_rule_template"`
}

type InstanceResV2 struct {
	Name             string                             `json:"instance_name"`
	DBType           string                             `json:"db_type" example:"mysql"`
	Host             string                             `json:"db_host" example:"10.10.10.10"`
	Port             string                             `json:"db_port" example:"3306"`
	User             string                             `json:"db_user" example:"root"`
	Desc             string                             `json:"desc" example:"this is a instance"`
	MaintenanceTimes []*v1.MaintenanceTimeResV1         `json:"maintenance_times" from:"maintenance_times"`
	RuleTemplate     *RuleTemplateV2                    `json:"rule_template,omitempty"`
	AdditionalParams []*v1.InstanceAdditionalParamResV1 `json:"additional_params"`
	SQLQueryConfig   *v1.SQLQueryConfigResV1            `json:"sql_query_config"`
	Source           string                             `json:"source" example:"SQLE"`
}

type GetInstancesResV2 struct {
	controller.BaseRes
	Data      []InstanceResV2 `json:"data"`
	TotalNums uint64          `json:"total_nums"`
}

// GetInstances get instances
// @Summary 获取实例信息列表
// @Description get instance info list
// @Id getInstanceListV2
// @Tags instance
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param filter_instance_name query string false "filter instance name"
// @Param filter_db_type query string false "filter db type"
// @Param filter_db_host query string false "filter db host"
// @Param filter_db_port query string false "filter db port"
// @Param filter_db_user query string false "filter db user"
// @Param filter_rule_template_name query string false "filter rule template name"
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Success 200 {object} GetInstancesResV2
// @router /v2/projects/{project_name}/instances [get]
func GetInstances(c echo.Context) error {
	req := new(GetInstancesReqV2)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	projectName := c.Param("project_name")
	err = v1.CheckIsProjectMember(user.Name, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var offset uint32
	if req.PageIndex >= 1 {
		offset = req.PageSize * (req.PageIndex - 1)
	}
	data := map[string]interface{}{
		"filter_instance_name":      req.FilterInstanceName,
		"filter_project_name":       projectName,
		"filter_db_host":            req.FilterDBHost,
		"filter_db_port":            req.FilterDBPort,
		"filter_db_user":            req.FilterDBUser,
		"filter_rule_template_name": req.FilterRuleTemplateName,
		"filter_db_type":            req.FilterDBType,
		"limit":                     req.PageSize,
		"offset":                    offset,
	}

	instances, count, err := s.GetInstancesByReq(data, user)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	templateNamesInProject, err := s.GetRuleTemplateNamesByProjectName(projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	instancesRes := []InstanceResV2{}
	for _, instance := range instances {
		var ruleTemplate *RuleTemplateV2
		if len(instance.RuleTemplateNames) >= 1 {
			ruleTemplate = &RuleTemplateV2{
				Name: instance.RuleTemplateNames[0],
			}

			if !utils.StringsContains(templateNamesInProject, ruleTemplate.Name) {
				ruleTemplate.IsGlobalRuleTemplate = true
			}
		}

		instanceReq := InstanceResV2{
			Name:             instance.Name,
			DBType:           instance.DbType,
			Host:             instance.Host,
			Port:             instance.Port,
			User:             instance.User,
			Desc:             instance.Desc,
			MaintenanceTimes: v1.ConvertPeriodToMaintenanceTimeResV1(instance.MaintenancePeriod),
			RuleTemplate:     ruleTemplate,
			SQLQueryConfig: &v1.SQLQueryConfigResV1{
				MaxPreQueryRows:                  instance.SqlQueryConfig.MaxPreQueryRows,
				QueryTimeoutSecond:               instance.SqlQueryConfig.QueryTimeoutSecond,
				AuditEnabled:                     instance.SqlQueryConfig.AuditEnabled,
				AllowQueryWhenLessThanAuditLevel: instance.SqlQueryConfig.AllowQueryWhenLessThanAuditLevel,
			},
			Source: instance.Source,
		}
		instancesRes = append(instancesRes, instanceReq)
	}
	return c.JSON(http.StatusOK, &GetInstancesResV2{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      instancesRes,
		TotalNums: count,
	})
}

type CreateInstanceReqV2 struct {
	Name             string                             `json:"instance_name" form:"instance_name" example:"test" valid:"required,name"`
	DBType           string                             `json:"db_type" form:"db_type" example:"mysql"`
	User             string                             `json:"db_user" form:"db_user" example:"root" valid:"required"`
	Host             string                             `json:"db_host" form:"db_host" example:"10.10.10.10" valid:"required,ip_addr|uri|hostname|hostname_rfc1123"`
	Port             string                             `json:"db_port" form:"db_port" example:"3306" valid:"required,port"`
	Password         string                             `json:"db_password" form:"db_password" example:"123456" valid:"required"`
	Desc             string                             `json:"desc" example:"this is a test instance"`
	SQLQueryConfig   *v1.SQLQueryConfigReqV1            `json:"sql_query_config" form:"sql_query_config"`
	MaintenanceTimes []*v1.MaintenanceTimeReqV1         `json:"maintenance_times" form:"maintenance_times"`
	RuleTemplateName string                             `json:"rule_template_name" form:"rule_template_name" valid:"required"`
	AdditionalParams []*v1.InstanceAdditionalParamReqV1 `json:"additional_params" form:"additional_params"`
}

// CreateInstance create instance
// @Summary 添加实例
// @Description create a instance
// @Id createInstanceV2
// @Tags instance
// @Security ApiKeyAuth
// @Accept json
// @Param project_name path string true "project name"
// @Param instance body v2.CreateInstanceReqV2 true "add instance"
// @Success 200 {object} controller.BaseRes
// @router /v2/projects/{project_name}/instances [post]
func CreateInstance(c echo.Context) error {
	s := model.GetStorage()
	req := new(CreateInstanceReqV2)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	projectName := c.Param("project_name")

	archived, err := s.IsProjectArchived(projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if archived {
		return controller.JSONBaseErrorReq(c, v1.ErrProjectArchived)
	}

	userName := controller.GetUserName(c)
	err = v1.CheckIsProjectManager(userName, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	_, exist, err := s.GetInstanceByNameAndProjectName(req.Name, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("instance is exist")))
	}

	if req.DBType == "" {
		req.DBType = driverV2.DriverTypeMySQL
	}

	maintenancePeriod := v1.ConvertMaintenanceTimeReqV1ToPeriod(req.MaintenanceTimes)
	if !maintenancePeriod.SelfCheck() {
		return controller.JSONBaseErrorReq(c, v1.ErrWrongTimePeriod)
	}

	additionalParams := driver.GetPluginManager().AllAdditionalParams()[req.DBType]
	for _, additionalParam := range req.AdditionalParams {
		err = additionalParams.SetParamValue(additionalParam.Name, additionalParam.Value)
		if err != nil {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, err))
		}
	}

	sqlQueryConfig := model.SqlQueryConfig{}
	if req.SQLQueryConfig != nil {
		sqlQueryConfig = model.SqlQueryConfig{
			MaxPreQueryRows:                  req.SQLQueryConfig.MaxPreQueryRows,
			QueryTimeoutSecond:               req.SQLQueryConfig.QueryTimeoutSecond,
			AuditEnabled:                     req.SQLQueryConfig.AuditEnabled,
			AllowQueryWhenLessThanAuditLevel: req.SQLQueryConfig.AllowQueryWhenLessThanAuditLevel,
		}
	}
	// default value
	if sqlQueryConfig.QueryTimeoutSecond == 0 {
		sqlQueryConfig.QueryTimeoutSecond = 10
	}
	// default value
	if sqlQueryConfig.MaxPreQueryRows == 0 {
		sqlQueryConfig.MaxPreQueryRows = 100
	}

	if sqlQueryConfig.AuditEnabled && sqlQueryConfig.AllowQueryWhenLessThanAuditLevel == "" {
		sqlQueryConfig.AllowQueryWhenLessThanAuditLevel = string(driverV2.RuleLevelError)
	}

	project, exist, err := s.GetProjectByName(projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrProjectNotExist(projectName))
	}

	instance := &model.Instance{
		DbType:             req.DBType,
		Name:               req.Name,
		User:               req.User,
		Host:               req.Host,
		Port:               req.Port,
		Password:           req.Password,
		Desc:               req.Desc,
		AdditionalParams:   additionalParams,
		MaintenancePeriod:  maintenancePeriod,
		SqlQueryConfig:     sqlQueryConfig,
		WorkflowTemplateId: project.WorkflowTemplateId,
		ProjectId:          project.ID,
		Source:             model.InstanceSourceSQLE,
	}

	var templates []*model.RuleTemplate
	template, exist, err := s.GetGlobalAndProjectRuleTemplateByNameAndProjectId(req.RuleTemplateName, project.ID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrRuleTemplateNotExist)
	}
	templates = append(templates, template)

	err = v1.CheckInstanceAndRuleTemplateDbType(templates, instance)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = s.Save(instance)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = s.UpdateInstanceRuleTemplates(instance, templates...)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}
