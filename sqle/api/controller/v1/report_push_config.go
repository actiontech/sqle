package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	dms "github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
)

type GetReportPushConfigsListResV1 struct {
	controller.BaseRes
	Data []ReportPushConfigList `json:"data"`
}

type ReportPushConfigList struct {
	Type              string    `json:"type"`
	Enabled           bool      `json:"enabled"`
	TriggerType       string    `json:"trigger_type "`
	PushFrequencyCron string    `json:"push_frequency_cron"`
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
// @Param project_name path string true "project name"
// @Success 200 {object} GetReportPushConfigsListResV1
// @Router /v1/projects/{project_name}/report_push_configs [get]
func GetReportPushConfigList(c echo.Context) error {
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	reportPushConfigs, err := model.GetStorage().GetReportPushConfigListInProject(projectUid)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
	}

	ret := make([]ReportPushConfigList, 0, len(reportPushConfigs))
	for _, reportPushConfig := range reportPushConfigs {
		ret = append(ret, ReportPushConfigList{
			Type:              reportPushConfig.Type,
			Enabled:           reportPushConfig.Enabled,
			TriggerType:       reportPushConfig.TriggerType,
			PushFrequencyCron: reportPushConfig.PushFrequencyCron,
			PushUserType:      reportPushConfig.PushUserType,
			PushUserList:      reportPushConfig.PushUserList,
			LastPushTime:      reportPushConfig.LastPushTime,
		})
	}
	return c.JSON(http.StatusOK, GetReportPushConfigsListResV1{
		Data: ret,
	})
}

type UpdateReportPushConfigReqV1 struct {
	TriggerType       string   `json:"trigger_type "`
	PushFrequencyCron string   `json:"push_frequency_cron"`
	PushUserType      string   `json:"push_user_Type"`
	PushUserList      []string `json:"push_user_list"`
	Enabled           string   `json:"enabled"`
}

// @Summary 更新消息推送配置
// @Description update report push config
// @Id UpdateReportPushConfig
// @Tags ReportPushConfig
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param report_push_config_id path string true "report push config id"
// @Param req body v1.UpdateReportPushConfigReqV1 true "update report push config request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/report_push_configs/{report_push_config_id}/ [put]
func UpdateReportPushConfig(c echo.Context) error {
	return nil
}
