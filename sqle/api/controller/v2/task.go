package v2

import (
	"github.com/labstack/echo/v4"

	"github.com/actiontech/sqle/sqle/api/controller"
)

type GetAuditTaskSQLsReqV2 struct {
	FilterExecStatus  string `json:"filter_exec_status" query:"filter_exec_status"`
	FilterAuditStatus string `json:"filter_audit_status" query:"filter_audit_status"`
	FilterAuditLevel  string `json:"filter_audit_level" query:"filter_audit_level"`
	NoDuplicate       bool   `json:"no_duplicate" query:"no_duplicate"`
	PageIndex         uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize          uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetAuditTaskSQLsResV2 struct {
	controller.BaseRes
	Data      []*AuditTaskSQLResV2 `json:"data"`
	TotalNums uint64               `json:"total_nums"`
}

type AuditTaskSQLResV2 struct {
	Number      uint           `json:"number"`
	ExecSQL     string         `json:"exec_sql"`
	AuditResult string         `json:"audit_result"`
	AuditLevel  string         `json:"audit_level"`
	AuditStatus string         `json:"audit_status"`
	ExecResult  []*AuditResult `json:"exec_result"`
	ExecStatus  string         `json:"exec_status"`
	RollbackSQL string         `json:"rollback_sql,omitempty"`
	Description string         `json:"description"`
}

type AuditResult struct {
	Level    string `json:"level"`
	Message  string `json:"message"`
	RuleName string `json:"rule_name"`
}

// @Summary 获取指定扫描任务的SQLs信息
// @Description get information of all SQLs belong to the specified audit task
// @Tags task
// @Id getAuditTaskSQLsV2
// @Security ApiKeyAuth
// @Param task_id path string true "task id"
// @Param filter_exec_status query string false "filter: exec status of task sql" Enums(initialized,doing,succeeded,failed,manually_executed)
// @Param filter_audit_status query string false "filter: audit status of task sql" Enums(initialized,doing,finished)
// @Param filter_audit_level query string false "filter: audit level of task sql" Enums(normal,notice,warn,error)
// @Param no_duplicate query boolean false "select unique (fingerprint and audit result) for task sql"
// @Param page_index query string true "page index"
// @Param page_size query string true "page size"
// @Success 200 {object} v2.GetAuditTaskSQLsResV2
// @router /v2/tasks/audits/{task_id}/sqls [get]
func GetTaskSQLs(c echo.Context) error {
	return controller.JSONNewNotImplementedErr(c)
}
