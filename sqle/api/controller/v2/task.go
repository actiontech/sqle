package v2

import (
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server"

	"github.com/labstack/echo/v4"
)

type GetAuditTaskSQLsReqV2 struct {
	FilterExecStatus  string `json:"filter_exec_status" query:"filter_exec_status"`
	FilterAuditStatus string `json:"filter_audit_status" query:"filter_audit_status"`
	FilterAuditLevel  string `json:"filter_audit_level" query:"filter_audit_level"`
	FilterAuditFileId uint   `json:"filter_audit_file_id" query:"filter_audit_file_id"`
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
	ExecSqlID                   uint                          `json:"exec_sql_id"`
	Number                      uint                          `json:"number"`
	ExecSQL                     string                        `json:"exec_sql"`
	SQLSourceFile               string                        `json:"sql_source_file"`
	SQLStartLine                uint64                        `json:"sql_start_line"`
	AuditResult                 []*AuditResult                `json:"audit_result"`
	AuditLevel                  string                        `json:"audit_level"`
	AuditStatus                 string                        `json:"audit_status"`
	ExecResult                  string                        `json:"exec_result"`
	ExecStatus                  string                        `json:"exec_status"`
	RollbackSQLs                []string                      `json:"rollback_sqls,omitempty"`
	Description                 string                        `json:"description"`
	SQLType                     string                        `json:"sql_type"`
	BackupStrategy              string                        `json:"backup_strategy" enums:"none,manual,reverse_sql,original_row"`
	BackupStrategyTip           string                        `json:"backup_strategy_tip"`
	BackupStatus                string                        `json:"backup_status" enums:"waiting_for_execution,executing,failed,succeed"`
	BackupResult                string                        `json:"backup_result"`
	AssociatedRollbackWorkflows []*AssociatedRollbackWorkflow `json:"associated_rollback_workflows"`
}

type AssociatedRollbackWorkflow struct {
	WorkflowID   string `json:"workflow_id"`
	WorkflowName string `json:"workflow_name"`
	Status       string `json:"status" enums:"wait_for_audit,wait_for_execution,rejected,canceled,exec_failed,executing,finished"`
}

type AuditResult struct {
	Level               string                    `json:"level" example:"warn"`
	Message             string                    `json:"message" example:"避免使用不必要的内置函数md5()"`
	RuleName            string                    `json:"rule_name"`
	DbType              string                    `json:"db_type"`
	I18nAuditResultInfo model.I18nAuditResultInfo `json:"i18n_audit_result_info"`
}

// @Summary 获取指定扫描任务的SQLs信息
// @Description get information of all SQLs belong to the specified audit task
// @Tags task
// @Id getAuditTaskSQLsV2
// @Security ApiKeyAuth
// @Param task_id path string true "task id"
// @Param filter_exec_status query string false "filter: exec status of task sql" Enums(initialized,doing,succeeded,failed,manually_executed,terminating,terminate_succeeded,terminate_failed,execute_rollback)
// @Param filter_audit_status query string false "filter: audit status of task sql" Enums(initialized,doing,finished)
// @Param filter_audit_level query string false "filter: audit level of task sql" Enums(normal,notice,warn,error)
// @Param no_duplicate query boolean false "select unique (fingerprint and audit result) for task sql"
// @Param filter_audit_file_id query uint false "filter: audit file id of task"
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

	limit, offset := controller.GetLimitAndOffset(req.PageIndex, req.PageSize)
	data := map[string]interface{}{
		"task_id":              taskId,
		"filter_exec_status":   req.FilterExecStatus,
		"filter_audit_status":  req.FilterAuditStatus,
		"filter_audit_level":   req.FilterAuditLevel,
		"filter_audit_file_id": req.FilterAuditFileId,
		"no_duplicate":         req.NoDuplicate,
		"limit":                limit,
		"offset":               offset,
	}

	taskSQLs, count, err := s.GetTaskSQLsByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	relations, err := s.GetExecuteSqlRollbackWorkflowRelationByTaskId(task.ID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	associatedRollbackWorkflowsMap := make(map[uint][]*AssociatedRollbackWorkflow)
	workflowIdMap := make(map[string]struct{})
	for _, relation := range relations {
		associatedRollbackWorkflowsMap[relation.ExecuteSqlId] = append(associatedRollbackWorkflowsMap[relation.ExecuteSqlId], &AssociatedRollbackWorkflow{
			WorkflowID:   relation.RollbackWorkflowId,
			WorkflowName: relation.RollbackWorkflowSubject,
			Status:       relation.RollbackWorkflowStatus,
		})
		workflowIdMap[relation.RollbackWorkflowId] = struct{}{}
	}
	taskSQLsRes := make([]*AuditTaskSQLResV2, 0, len(taskSQLs))
	backupService := server.BackupService{}
	rollbackSqlMap, err := backupService.GetRollbackSqlsMap(task.ID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	backupTaskMap, err := backupService.GetBackupTasksMap(task.ID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	for _, taskSQL := range taskSQLs {
		taskSQLRes := &AuditTaskSQLResV2{
			ExecSqlID:                   taskSQL.Id,
			Number:                      taskSQL.Number,
			Description:                 taskSQL.Description,
			ExecSQL:                     taskSQL.ExecSQL,
			SQLSourceFile:               taskSQL.SQLSourceFile.String,
			SQLStartLine:                taskSQL.SQLStartLine,
			AuditLevel:                  taskSQL.AuditLevel,
			AuditStatus:                 taskSQL.AuditStatus,
			ExecResult:                  taskSQL.ExecResult,
			ExecStatus:                  taskSQL.ExecStatus,
			RollbackSQLs:                rollbackSqlMap[taskSQL.Id],
			SQLType:                     taskSQL.SQLType.String,
			BackupStrategy:              backupTaskMap.GetBackupStrategy(taskSQL.Id),
			BackupStrategyTip:           backupTaskMap.GetBackupStrategyTip(taskSQL.Id),
			AssociatedRollbackWorkflows: associatedRollbackWorkflowsMap[taskSQL.Id],
		}
		for i := range taskSQL.AuditResults {
			ar := taskSQL.AuditResults[i]
			taskSQLRes.AuditResult = append(taskSQLRes.AuditResult, &AuditResult{
				Level:               ar.Level,
				Message:             ar.GetAuditMsgByLangTag(locale.Bundle.GetLangTagFromCtx(c.Request().Context())),
				RuleName:            ar.RuleName,
				DbType:              task.DBType,
				I18nAuditResultInfo: ar.I18nAuditResultInfo,
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
