package v2

import (
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
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
	Number        uint           `json:"number"`
	ExecSQL       string         `json:"exec_sql"`
	SQLSourceFile string         `json:"sql_source_file"`
	SQLStartLine  uint64         `json:"sql_start_line"`
	AuditResult   []*AuditResult `json:"audit_result"`
	AuditLevel    string         `json:"audit_level"`
	AuditStatus   string         `json:"audit_status"`
	ExecResult    string         `json:"exec_result"`
	ExecStatus    string         `json:"exec_status"`
	RollbackSQL   string         `json:"rollback_sql,omitempty"`
	Description   string         `json:"description"`
	SQLType       string         `json:"sql_type"`
}

type AuditResult struct {
	Level    string `json:"level" example:"warn"`
	Message  string `json:"message" example:"避免使用不必要的内置函数md5()"`
	RuleName string `json:"rule_name"`
	DbType   string `json:"db_type"`
}

// @Summary 获取指定扫描任务的SQLs信息
// @Description get information of all SQLs belong to the specified audit task
// @Tags task
// @Id getAuditTaskSQLsV2
// @Security ApiKeyAuth
// @Param task_id path string true "task id"
// @Param filter_exec_status query string false "filter: exec status of task sql" Enums(initialized,doing,succeeded,failed,manually_executed,terminating,terminate_succeeded,terminate_failed)
// @Param filter_audit_status query string false "filter: audit status of task sql" Enums(initialized,doing,finished)
// @Param filter_audit_level query string false "filter: audit level of task sql" Enums(normal,notice,warn,error)
// @Param no_duplicate query boolean false "select unique (fingerprint and audit result) for task sql"
// @Param page_index query string true "page index"
// @Param page_size query string true "page size"
// @Success 200 {object} v2.GetAuditTaskSQLsResV2
// @router /v2/tasks/audits/{task_id}/sqls [get]
func GetTaskSQLs(c echo.Context) error {
	req := new(GetAuditTaskSQLsReqV2)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()
	taskId := c.Param("task_id")
	task, err := v1.GetTaskById(c.Request().Context(), taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = v1.CheckCurrentUserCanViewTask(c, task)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var offset uint32
	if req.PageIndex >= 1 {
		offset = req.PageSize * (req.PageIndex - 1)
	}
	data := map[string]interface{}{
		"task_id":             taskId,
		"filter_exec_status":  req.FilterExecStatus,
		"filter_audit_status": req.FilterAuditStatus,
		"filter_audit_level":  req.FilterAuditLevel,
		"no_duplicate":        req.NoDuplicate,
		"limit":               req.PageSize,
		"offset":              offset,
	}

	taskSQLs, count, err := s.GetTaskSQLsByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	taskSQLsRes := make([]*AuditTaskSQLResV2, 0, len(taskSQLs))
	for _, taskSQL := range taskSQLs {
		taskSQLRes := &AuditTaskSQLResV2{
			Number:        taskSQL.Number,
			Description:   taskSQL.Description,
			ExecSQL:       taskSQL.ExecSQL,
			SQLSourceFile: taskSQL.SQLSourceFile.String,
			SQLStartLine:  taskSQL.SQLStartLine,
			AuditLevel:    taskSQL.AuditLevel,
			AuditStatus:   taskSQL.AuditStatus,
			ExecResult:    taskSQL.ExecResult,
			ExecStatus:    taskSQL.ExecStatus,
			RollbackSQL:   taskSQL.RollbackSQL.String,
			SQLType:       taskSQL.SQLType.String,
		}
		for i := range taskSQL.AuditResults {
			ar := taskSQL.AuditResults[i]
			taskSQLRes.AuditResult = append(taskSQLRes.AuditResult, &AuditResult{
				Level:    ar.Level,
				Message:  ar.Message,
				RuleName: ar.RuleName,
				DbType:   task.DBType,
			})
		}

		taskSQLsRes = append(taskSQLsRes, taskSQLRes)
	}

	return c.JSON(http.StatusOK, &GetAuditTaskSQLsResV2{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      taskSQLsRes,
		TotalNums: count,
	})
}

type AffectRows struct {
	Count      int    `json:"count"`
	ErrMessage string `json:"err_message"`
}

type SQLExplain struct {
	SQL           string                   `json:"sql"`
	ErrMessage    string                   `json:"err_message"`
	ClassicResult *v1.ExplainClassicResult `json:"classic_result"`
}

type PerformanceStatistics struct {
	AffectRows *AffectRows `json:"affect_rows"`
}

type TableMetas struct {
	ErrMessage string          `json:"err_message"`
	Items      []*v1.TableMeta `json:"table_meta_items"`
}

type TaskAnalysisDataV2 struct {
	SQLExplain            *SQLExplain            `json:"sql_explain"`
	TableMetas            *TableMetas            `json:"table_metas"`
	PerformanceStatistics *PerformanceStatistics `json:"performance_statistics"`
}

type GetTaskAnalysisDataResV2 struct {
	controller.BaseRes
	Data *TaskAnalysisDataV2 `json:"data"`
}

// GetTaskAnalysisData get SQL explain and related table metadata for analysis
// @Summary 获取task相关的SQL执行计划和表元数据
// @Description get SQL explain and related table metadata for analysis
// @Id getTaskAnalysisDataV2
// @Tags task
// @Param task_id path string true "task id"
// @Param number path uint true "sql number"
// @Security ApiKeyAuth
// @Success 200 {object} v2.GetTaskAnalysisDataResV2
// @router /v2/tasks/audits/{task_id}/sqls/{number}/analysis [get]
func GetTaskAnalysisData(c echo.Context) error {
	return getTaskAnalysisData(c)
}
