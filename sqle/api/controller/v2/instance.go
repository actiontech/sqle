package v2

import (
	"fmt"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/model"

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

type GetInstanceResV2 struct {
	controller.BaseRes
	Data InstanceResV2 `json:"data"`
}

// GetInstance get instance
// @Summary 获取实例信息
// @Description get instance db
// @Id getInstanceV2
// @Tags instance
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instance_name path string true "instance name"
// @Success 200 {object} v2.GetInstanceResV2
// @router /v2/projects/{project_name}/instances/{instance_name}/ [get]
func GetInstance(c echo.Context) error {
	instanceName := c.Param("instance_name")
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	instance, exist, err := dms.GetInstanceInProjectByName(c.Request().Context(), projectUid, instanceName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, v1.ErrInstanceNoAccess)
	}
	ruletemplate, err := model.GetStorage().GetRuleTemplateById(instance.RuleTemplateId)
	if err == nil {
		instance.RuleTemplates = []model.RuleTemplate{*ruletemplate}
	}
	return c.JSON(http.StatusOK, &GetInstanceResV2{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertInstanceToRes(instance),
	})
}

func convertInstanceToRes(instance *model.Instance) InstanceResV2 {
	instanceResV2 := InstanceResV2{
		Name:             instance.Name,
		Host:             instance.Host,
		Port:             instance.Port,
		User:             instance.User,
		Desc:             instance.Desc,
		DBType:           instance.DbType,
		MaintenanceTimes: v1.ConvertPeriodToMaintenanceTimeResV1(instance.MaintenancePeriod),
		AdditionalParams: []*v1.InstanceAdditionalParamResV1{},
		SQLQueryConfig: &v1.SQLQueryConfigResV1{
			MaxPreQueryRows:                  instance.SqlQueryConfig.MaxPreQueryRows,
			QueryTimeoutSecond:               instance.SqlQueryConfig.QueryTimeoutSecond,
			AuditEnabled:                     instance.SqlQueryConfig.AuditEnabled,
			AllowQueryWhenLessThanAuditLevel: instance.SqlQueryConfig.AllowQueryWhenLessThanAuditLevel,
		},
		Source: instance.Source,
	}

	if len(instance.RuleTemplates) != 0 {
		ruleTemplate := instance.RuleTemplates[0]
		instanceResV2.RuleTemplate = &RuleTemplateV2{
			Name:                 ruleTemplate.Name,
			IsGlobalRuleTemplate: ruleTemplate.ProjectId == model.ProjectIdForGlobalRuleTemplate,
		}
	}
	for _, param := range instance.AdditionalParams {
		instanceResV2.AdditionalParams = append(instanceResV2.AdditionalParams, &v1.InstanceAdditionalParamResV1{
			Name:        param.Key,
			Description: param.Desc,
			Type:        string(param.Type),
			Value:       fmt.Sprintf("%v", param.Value),
		})
	}
	return instanceResV2
}
