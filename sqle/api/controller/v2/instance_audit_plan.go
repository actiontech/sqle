package v2

import (
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	"github.com/labstack/echo/v4"
)

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
// @Success 200 {object} v1.GetInstanceTipsResV1
// @router /v2/projects/{project_name}/instance_tips [get]
func GetInstanceTips(c echo.Context) error {
	return v1.GetInstanceTips(c)
}
