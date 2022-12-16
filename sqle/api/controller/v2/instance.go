package v2

import (
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
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
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
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

			if !isTemplateExistsInProject(ruleTemplate.Name, templateNamesInProject) {
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
		}
		instancesRes = append(instancesRes, instanceReq)
	}
	return c.JSON(http.StatusOK, &GetInstancesResV2{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      instancesRes,
		TotalNums: count,
	})
}

func isTemplateExistsInProject(templateName string, templateNamesInProject []string) bool {
	return utils.IsStringExistsInStrings(templateName, templateNamesInProject)
}
