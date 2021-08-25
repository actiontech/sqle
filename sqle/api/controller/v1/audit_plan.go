package v1

import (
	"actiontech.cloud/sqle/sqle/sqle/api/controller"

	"github.com/labstack/echo/v4"
)

type CreateAuditPlanReqV1 struct {
	Name             string `json:"audit_plan_name" form:"audit_plan_name" example:"audit_plan_for_java_repo_1" valid:"required,name"`
	Cron             string `json:"audit_plan_cron" form:"audit_plan_cron" example:"0 */2 * * *" valid:"required"`
	InstanceType     string `json:"audit_plan_instance_type" form:"audit_plan_instance_type" example:"mysql" valid:"required"`
	InstanceName     string `json:"audit_plan_instance_name" form:"audit_plan_instance_name" example:"test_mysql"`
	InstanceDatabase string `json:"audit_plan_instance_database" form:"audit_plan_instance_database" example:"app1"`
}

// @Summary 添加审核计划
// @Description create audit plan
// @Id createAuditPlanV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Accept json
// @Param audit_plan body v1.CreateAuditPlanReqV1 true "create audit plan"
// @Success 200 {object} controller.BaseRes
// @router /v1/audit_plans [post]
func CreateAuditPlan(c echo.Context) error { return nil }

// @Summary 删除审核计划
// @Description delete audit plan
// @Id deleteAuditPlanV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param audit_plan_name path string true "audit plan name"
// @Success 200 {object} controller.BaseRes
// @router /v1/audit_plans/{audit_plan_name}/ [delete]
func DeleteAuditPlan(c echo.Context) error { return nil }

type UpdateAuditPlanReqV1 struct {
	Cron             *string `json:"audit_plan_cron" form:"audit_plan_cron" example:"0 */2 * * *"`
	InstanceName     *string `json:"audit_plan_instance_name" form:"audit_plan_instance_name" example:"test_mysql"`
	InstanceDatabase *string `json:"audit_plan_instance_database" form:"audit_plan_instance_database" example:"app1"`
}

// @Summary 更新审核计划
// @Description update audit plan
// @Id updateAuditPlanV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param audit_plan_name path string true "audit plan name"
// @param audit_plan body v1.UpdateAuditPlanReqV1 true "update audit plan"
// @Success 200 {object} controller.BaseRes
// @router /v1/audit_plans/{audit_plan_name}/ [patch]
func UpdateAuditPlan(c echo.Context) error { return nil }

type GetAuditPlansReqV1 struct {
	FilterAuditPlanDBType string `json:"filter_audit_plan_db_type" query:"filter_audit_plan_db_type"`
	PageIndex             uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize              uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetAuditPlansResV1 struct {
	controller.BaseRes
	Data      []AuditPlanResV1 `json:"data"`
	TotalNums uint64           `json:"total_nums"`
}

type AuditPlanResV1 struct {
	Name             string `json:"audit_plan_name" example:"audit_for_java_app1"`
	Cron             string `json:"audit_plan_cron" example:"0 */2 * * *"`
	DBType           string `json:"audit_plan_db_type" example:"mysql"`
	Token            string `json:"audit_plan_token" example:"it's a JWT Token for scanner"`
	InstanceName     string `json:"audit_plan_instance_name" example:"test_mysql"`
	InstanceDatabase string `json:"audit_plan_instance_database" example:"app1"`
}

// @Summary 获取审核计划信息列表
// @Description get audit plan info list
// @Id getAuditPlansV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param filter_audit_plan_db_type query string false "filter audit plan db type"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v1.GetAuditPlansResV1
// @router /v1/audit_plans [get]
func GetAuditPlans(c echo.Context) error { return nil }

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
	Id        string `json:"audit_plan_report_id" example:"1"`
	Timestamp string `json:"audit_plan_report_timestamp" example:"RFC3339"`
}

// @Summary 获取指定审核计划的报告列表
// @Description get audit plan report list
// @Id getAuditPlanReportsV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param audit_plan_name path string true "audit plan name"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v1.GetAuditPlanReportsResV1
// @router /v1/audit_plans/{audit_plan_name}/reports [get]
func GetAuditPlanReports(c echo.Context) error { return nil }

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
	Fingerprint          string `json:"audit_plan_report_sql_fingerprint" example:"select * from t1 where id = ?"`
	LastReceiveText      string `json:"audit_plan_report_sql_last_receive_text" example:"select * from t1 where id = 1"`
	LastReceiveTimestamp string `json:"audit_plan_report_sql_last_receive_timestamp" example:"RFC3339"`
	AuditResult          string `json:"audit_plan_report_sql_audit_result" example:"same format as task audit result"`
}

// @Summary 获取指定审核计划的SQL审核详情
// @Description get audit plan report SQLs
// @Id getAuditPlanReportSQLsV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param audit_plan_name path string true "audit plan name"
// @Param audit_plan_report_id path string true "audit plan report id"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v1.GetAuditPlanReportSQLsResV1
// @router /v1/audit_plans/{audit_plan_name}/report/{audit_plan_report_id}/ [get]
func GetAuditPlanReportSQLs(c echo.Context) error { return nil }

type FullSyncAuditPlanSQLsReqV1 struct {
	SQLs []AuditPlanSQLReqV1 `json:"audit_plan_sql_list" form:"audit_plan_sql_list"`
}

type AuditPlanSQLReqV1 struct {
	Fingerprint          string `json:"audit_plan_sql_fingerprint" form:"audit_plan_sql_fingerprint" example:"select * from t1 where id = ?" valid:"required"`
	Counter              string `json:"audit_plan_sql_counter" form:"audit_plan_sql_counter" example:"6" valid:"required"`
	LastReceiveText      string `json:"audit_plan_sql_last_receive_text" form:"audit_plan_sql_last_receive_text" example:"select * from t1 where id = 1"`
	LastReceiveTimestamp string `json:"audit_plan_sql_last_receive_timestamp" form:"audit_plan_sql_last_receive_timestamp" example:"RFC3339"`
}

// @Summary 全量同步SQL到审核计划
// @Description full sync audit plan SQLs
// @Id fullSyncAuditPlanSQLsV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param audit_plan_name path string true "audit plan name"
// @Param sqls body v1.FullSyncAuditPlanSQLsReqV1 true "full sync audit plan SQLs request"
// @Success 200 {object} controller.BaseRes
// @router /v1/audit_plans/{audit_plan_name}/sqls/full [post]
func FullSyncAuditPlanSQLs(c echo.Context) error { return nil }

type PartialSyncAuditPlanSQLsReqV1 struct {
	SQLs []AuditPlanSQLReqV1 `json:"audit_plan_sql_list" form:"audit_plan_sql_list"`
}

// @Summary 增量同步SQL到审核计划
// @Description partial sync audit plan SQLs
// @Id partialSyncAuditPlanSQLsV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param audit_plan_name path string true "audit plan name"
// @Param sqls body v1.PartialSyncAuditPlanSQLsReqV1 true "partial sync audit plan SQLs request"
// @Success 200 {object} controller.BaseRes
// @router /v1/audit_plans/{audit_plan_name}/sqls/partial [post]
func PartialSyncAuditPlanSQLs(c echo.Context) error { return nil }

type GetAuditPlanSQLsReqV1 struct {
	PageIndex uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetAuditPlanSQLsResV1 struct {
	controller.BaseRes
	Data      []AuditPlanSQLResV1 `json:"data"`
	TotalNums uint64              `json:"total_nums"`
}

type AuditPlanSQLResV1 struct {
	Fingerprint          string `json:"audit_plan_sql_fingerprint" example:"select * from t1 where id = ?"`
	Counter              string `json:"audit_plan_sql_counter" example:"6"`
	LastReceiveText      string `json:"audit_plan_sql_last_receive_text" example:"select * from t1 where id = 1"`
	LastReceiveTimestamp string `json:"audit_plan_sql_last_receive_timestamp" example:"RFC3339"`
}

// @Summary 获取指定审核计划的SQLs信息(不包括审核结果)
// @Description get audit plan SQLs
// @Id getAuditPlanSQLsV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param audit_plan_name path string true "audit plan name"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v1.GetAuditPlanSQLsResV1
// @router /v1/audit_plans/{audit_plan_name}/sqls [get]
func GetAuditPlanSQLs(c echo.Context) error { return nil }

// @Summary 触发审核计划
// @Description trigger audit plan
// @Id triggerAuditPlanV1
// @Tags audit_plan
// @Security ApiKeyAuth
// @Param audit_plan_name path string true "audit plan name"
// @Success 200 {object} controller.BaseRes
// @router /v1/audit_plans/{audit_plan_name}/trigger [post]
func TriggerAuditPlan(c echo.Context) error { return nil }
