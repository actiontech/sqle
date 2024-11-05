package v1

import "github.com/labstack/echo/v4"

type UpdateSqlBackupStrategyReq struct {
	Strategy  string `json:"strategy" enum:"none,manual,reverse_sql,origin_row"`
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
	Strategy  string `json:"strategy" enum:"none,manual,reverse_sql,origin_row"`
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
