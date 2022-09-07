package v2

import (
	"github.com/labstack/echo/v4"
)

type CreateWorkflowReqV2 struct {
	Subject      string   `json:"workflow_subject" form:"workflow_subject" valid:"required,name"`
	Desc         string   `json:"desc" form:"desc"`
	TaskIds      []string `json:"task_ids" form:"task_ids" valid:"required"`
}

// CreateWorkflowV2
// @Summary 创建工单
// @Description create workflow
// @Accept json
// @Produce json
// @Tags workflow
// @Id createWorkflowV1
// @Security ApiKeyAuth
// @Param instance body v2.CreateWorkflowReqV2 true "create workflow request"
// @Success 200 {object} controller.BaseRes
// @router /v2/workflows [post]
func CreateWorkflowV2(c echo.Context) error {
	return nil
}
