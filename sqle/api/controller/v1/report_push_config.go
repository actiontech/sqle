package v1

import (
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

type GetReportPushConfigsListResV1 struct {
	controller.BaseRes
	Data []ReportPushConfigList `json:"data"`
}

type ReportPushConfigList struct {
	Type              string    `json:"type"`
	Status            string    `json:"status"`
	PushFrequencyType string    `json:"push_frequency_type"`
	PushFrequencyCron string    `json:"PushFrequency_cron"`
	PushUserType      string    `json:"push_user_Type"`
	PushUserList      []string  `json:"push_user_list"`
	LastPushTime      time.Time `json:"last_push_time"`
}

// GetReportPushConfigList
// @Summary 获取消息推送配置列表
// @Description Get report push config list
// @Id GetReportPushConfigList
// @Tags ReportPushConfig
// @Security ApiKeyAuth
// @Success 200 {object} GetReportPushConfigsListResV1
// @Router /v1/project/{project_name}/report_push_configs [get]
func GetReportPushConfigList(c echo.Context) error {
	return nil
}

type UpdateReportPushConfigReqV1 struct {
	PushFrequencyType string   `json:"push_frequency_type"`
	PushFrequencyCron string   `json:"PushFrequency_cron"`
	PushUserType      string   `json:"push_user_Type"`
	PushUserList      []string `json:"push_user_list"`
	Status            string   `json:"status"`
}

// @Summary 更新消息推送配置
// @Description update report push config
// @Id UpdateReportPushConfig
// @Tags report_push_config
// @Security ApiKeyAuth
// @Param report_push_config_id path string true "report push config id"
// @Param req body v1.UpdateReportPushConfigReqV1 true "update report push config request"
// @Success 200 {object} controller.BaseRes
// @router /v1/project/{project_name}/report_push_configs/{report_push_config_id} [put]
func UpdateReportPushConfig(c echo.Context) error {
	return nil
}
