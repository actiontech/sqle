package v1

import "github.com/labstack/echo/v4"

type CreateSyncInstanceTaskReqV1 struct {
	Source       string                         `json:"source" form:"source" validate:"required" example:"actiontech-dmp"`
	Version      string                         `json:"version" form:"version" validate:"required" example:"5.23.01.0"`
	URL          string                         `json:"url" form:"url" validate:"required" example:"http://10.186.62.56:10000"`
	DbType       string                         `json:"db_type" form:"db_type" validate:"required" example:"mysql"`
	RuleTemplate string                         `json:"rule_template" form:"rule_template" validate:"required" example:"default_mysql"`
	Cron         string                         `json:"cron" form:"cron" validate:"required" example:"0 0 * * *"`
	Params       []SyncTaskAdditionalParamReqV1 `json:"params" form:"params" validate:"dive"`
}

type SyncTaskAdditionalParamReqV1 struct {
	Key   string `json:"key" form:"key" valid:"required"`
	Value string `json:"value" form:"value" valid:"required"`
}

// CreateSyncInstanceTask create sync instance task
// @Summary 创建同步实例任务
// @Description create sync instance task
// @Id createSyncInstanceTaskV1
// @Tags sync_instance
// @Security ApiKeyAuth
// @Accept json
// @Param sync_task body v1.CreateSyncInstanceTaskReqV1 true "create sync instance task request"
// @Success 200 {object} controller.BaseRes
// @router /v1/task/sync_instance [post]
func CreateSyncInstanceTask(c echo.Context) error {
	return createSyncInstanceTask(c)
}

type UpdateSyncInstanceTaskReqV1 struct {
	Id           int                            `json:"id" form:"id" validate:"required" example:"1"`
	Version      *string                        `json:"version" form:"version" validate:"required" example:"5.23.01.0"`
	URL          *string                        `json:"url" form:"url" validate:"required" example:"http://10.186.62.56:10000"`
	RuleTemplate *string                        `json:"rule_template" form:"rule_template" validate:"required" example:"default_mysql"`
	Cron         *string                        `json:"cron" form:"cron" validate:"required" example:"0 0 * * *"`
	Params       []SyncTaskAdditionalParamReqV1 `json:"params" form:"params" validate:"dive"`
}

// UpdateSyncInstanceTask update sync instance task
// @Summary 更新同步实例任务
// @Description update sync instance task
// @Id updateSyncInstanceTaskV1
// @Tags sync_instance
// @Security ApiKeyAuth
// @param sync_task body v1.UpdateSyncInstanceTaskReqV1 true "update sync instance request"
// @Success 200 {object} controller.BaseRes
// @router /v1/task/sync_instance [patch]
func UpdateSyncInstanceTask(c echo.Context) error {
	return updateSyncInstanceTask(c)
}

type DeleteSyncInstanceTaskReqV1 struct {
	Id int `json:"id" form:"id" validate:"required" example:"1"`
}

// DeleteSyncInstanceTask delete sync instance task
// @Summary 删除同步实例任务
// @Description delete sync instance task
// @Id deleteSyncInstanceTaskV1
// @Tags sync_instance
// @Security ApiKeyAuth
// @param sync_task body v1.DeleteSyncInstanceTaskReqV1 true "delete sync instance request"
// @Success 200 {object} controller.BaseRes
// @router /v1/task/sync_instance [delete]
func DeleteSyncInstanceTask(c echo.Context) error {
	return deleteSyncInstanceTask(c)
}
