package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

type UpdateSqlBackupStrategyReq struct {
	Strategy string `json:"strategy" enum:"none,manual,reverse_sql,origin_row"`
}

// UpdateSqlBackupStrategy
// @Summary 更新单条SQL的备份策略
// @Description update back up strategy for one sql in workflow
// @Tags workflow
// @Accept json
// @Produce json
// @Id UpdateSqlBackupStrategyV1
// @Security ApiKeyAuth
// @Param task_id path string true "task id"
// @Param sql_id path string true "sql id"
// @Param strategy body v1.UpdateSqlBackupStrategyReq true "update back up strategy for one sql in workflow"
// @Success 200 {object} controller.BaseRes
// @router /v1/tasks/audits/{task_id}/sqls/{sql_id}/ [patch]
func UpdateSqlBackupStrategy(c echo.Context) error {
	return nil
}

type UpdateTaskBackupStrategyReq struct {
	Strategy string `json:"strategy" enum:"none,manual,reverse_sql,origin_row"`
}

// UpdateTaskBackupStrategy
// @Summary 更新工单中数据源对应所有SQL的备份策略
// @Description update back up strategy for all sqls in task
// @Tags workflow
// @Accept json
// @Produce json
// @Id UpdateTaskBackupStrategyV1
// @Security ApiKeyAuth
// @Param task_id path string true "task id"
// @Param strategy body v1.UpdateTaskBackupStrategyReq true "update back up strategy for sqls in workflow"
// @Success 200 {object} controller.BaseRes
// @router /v1/tasks/audits/{task_id}/ [patch]
func UpdateTaskBackupStrategy(c echo.Context) error {
	return nil
}

// @Summary 下载工单中的SQL备份
// @Description download SQL back up file for the audit task
// @Tags task
// @Id downloadBackupFileV1
// @Security ApiKeyAuth
// @Param workflow_id path string true "workflow id"
// @Param project_name path string true "project name"
// @Param task_id path string true "task id"
// @Success 200 file 1 "sql file"
// @router /v1/projects/{project_name}/workflows/{workflow_id}/tasks/{task_id}/backup_files/download [get]
func DownloadSqlBackupFile(c echo.Context) error {
	return nil
}

type BackupSqlListReq struct {
	FilterInstanceId string `json:"filter_instance_id" query:"filter_instance_id"`
	FilterExecStatus string `json:"filter_exec_status" query:"filter_exec_status"`
	PageIndex        uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize         uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type BackupSqlData struct {
	ExecOrder      uint     `json:"exec_order"`
	ExecSqlID      uint     `json:"exec_sql_id"`
	OriginSQL      string   `json:"origin_sql"`
	BackupSqls     []string `json:"backup_sqls"`
	BackupStrategy string   `json:"backup_strategy" enum:"none,manual,reverse_sql,origin_row"`
	InstanceName   string   `json:"instance_name"`
	InstanceId     string   `json:"instance_id "`
	ExecStatus     string   `json:"exec_status"`
	Description    string   `json:"description"`
}

type BackupSqlListRes struct {
	controller.BaseRes
	Data      []*BackupSqlData `json:"data"`
	TotalNums uint64           `json:"total_nums"`
}

// @Summary 获取工单下所有回滚SQL的列表
// @Description get backup sql list
// @Tags workflow
// @Id GetBackupSqlListV1
// @Security ApiKeyAuth
// @Param filter_exec_status query string false "filter: exec status of task sql" Enums(initialized,doing,succeeded,failed,manually_executed,terminating,terminate_succeeded,terminate_failed)
// @Param project_name path string true "project name"
// @Param workflow_id path string true "workflow id"
// @Param filter_instance_id query uint false "filter: instance id in workflow"
// @Param page_index query string true "page index"
// @Param page_size query string true "page size"
// @Success 200 {object} v1.BackupSqlListRes
// @router /v1/projects/{project_name}/workflows/{workflow_id}/backup_sqls [get]
func GetBackupSqlList(c echo.Context) error {
	return nil
}

type CreateRollbackWorkflowReq struct {
	Subject        string `json:"workflow_subject" form:"workflow_subject" valid:"required,name"`
	Desc           string `json:"desc" form:"desc"`
	SqlVersionID   *uint  `json:"sql_version_id" form:"sql_version_id"`
	TaskIds        []uint `json:"task_ids" form:"task_ids" valid:"required"`
	RollbackSqlIds []uint `json:"rollback_sql_ids" form:"rollback_sql_ids" valid:"required"`
}

type CreateRollbackWorkflowRes struct {
	controller.BaseRes
	Data *CreateRollbackWorkflowResData `json:"data"`
}

type CreateRollbackWorkflowResData struct {
	WorkflowID string `json:"workflow_id"`
}

// CreateRollbackWorkflow
// @Summary 创建回滚工单
// @Description create rollback workflow
// @Accept json
// @Produce json
// @Tags workflow
// @Id CreateRollbackWorkflow
// @Security ApiKeyAuth
// @Param instance body v1.CreateRollbackWorkflowReq true "create rollback workflow request"
// @Param project_name path string true "project name"
// @Param workflow_id path string true "origin workflow id to rollback"
// @Success 200 {object} CreateRollbackWorkflowRes
// @router /v1/projects/{project_name}/workflows/{workflow_id}/rollback [post]
func CreateRollbackWorkflow(c echo.Context) error {
	return nil
}
