package v1

import (
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

type GetMessagePushConfigsListResV1 struct {
	controller.BaseRes
	Data []MessagePushConfigList `json:"data"`
}

type MessagePushConfigList struct {
	Type              string    `json:"type"`
	Status            string    `json:"status"`
	PushFrequencyType string    `json:"push_frequency_type"`
	PushFrequencyCron string    `json:"PushFrequency_cron"`
	PushUserType      string    `json:"push_user_Type"`
	PushUserList      []string  `json:"push_user_list"`
	LastPushTime      time.Time `json:"last_push_time"`
}

// GetMessagePushConfigList
// @Summary 获取消息推送配置列表
// @Description Get message push config list
// @Id GetMessagePushConfigList
// @Tags MessagePushConfig
// @Security ApiKeyAuth
// @Success 200 {object} GetMessagePushConfigsListResV1
// @Router /v1/project/{project_name}/message_push_configs [get]
func GetMessagePushConfigList(c echo.Context) error {
	return nil
}

type UpdateMessagePushConfigReqV1 struct {
	PushFrequencyType string   `json:"push_frequency_type"`
	PushFrequencyCron string   `json:"PushFrequency_cron"`
	PushUserType      string   `json:"push_user_Type"`
	PushUserList      []string `json:"push_user_list"`
	Status            string   `json:"status"`
}

// @Summary 更新消息推送配置
// @Description update message push config
// @Id UpdateMessagePushConfig
// @Tags message_push_config
// @Security ApiKeyAuth
// @Param message_push_config_id path string true "message push config id"
// @Param req body v1.UpdateMessagePushConfigReqV1 true "update message push config request"
// @Success 200 {object} controller.BaseRes
// @router /v1/project/{project_name}/message_push_configs/{message_push_config_id} [put]
func UpdateMessagePushConfig(c echo.Context) error {
	return nil
}
