package v2

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	"github.com/labstack/echo/v4"
)

type InstanceForCreatingTask struct {
	InstanceName   string `json:"instance_name" form:"instance_name" valid:"required"`
	InstanceSchema string `json:"instance_schema" form:"instance_schema"`
}

type BatchCreateAndAuditTasksReqV1 struct {
	Instances []*InstanceForCreatingTask `json:"instances" valid:"dive,required"`
	Sql       string                     `json:"sql" form:"sql"`
}

type BatchCreateAndAuditTasksResV2 struct {
	controller.BaseRes
	Data []*v1.AuditTaskResV1 `json:"data"`
}

// BatchCreateAndAuditTasks
// @Summary 多数据源批量创建SQL审核任务并提交审核
// @Description create and audit tasks for each instance, you can upload sql content in three ways, any one can be used, but only one is effective.
// @Description 1. formData[sql]: sql content;
// @Description 2. file[input_sql_file]: it is a sql file;
// @Description 3. file[input_mybatis_xml_file]: it is mybatis xml file, sql will be parsed from it.
// @Accept mpfd
// @Produce json
// @Tags task
// @Id createAndAuditTasks
// @Security ApiKeyAuth
// @Param req body v2.BatchCreateAndAuditTasksReqV1 true "parameters for creating and audit tasks"
// @Param input_sql_file formData file false "input SQL file"
// @Param input_mybatis_xml_file formData file false "input mybatis XML file"
// @Success 200 {object} v2.BatchCreateAndAuditTasksResV2
// @router /v2/tasks/audits [post]
func BatchCreateAndAuditTasks(c echo.Context) error {
	return nil
}
