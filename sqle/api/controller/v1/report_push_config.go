package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
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
	Id                string    `json:"report_push_config_id"`
	Type              string    `json:"type"`
	Enabled           bool      `json:"enabled"`
	TriggerType       string    `json:"trigger_type" enums:"immediately,timing"`
	PushFrequencyCron string    `json:"push_frequency_cron"`
	PushUserType      string    `json:"push_user_Type"  enums:"fixed,permission_match"`
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
			Id:                reportPushConfig.GetIDStr(),
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
	TriggerType       string   `json:"trigger_type" enums:"immediately,timing" valid:"oneof=immediately timing"`
	PushFrequencyCron string   `json:"push_frequency_cron"`
	PushUserType      string   `json:"push_user_Type"  enums:"fixed,permission_match" valid:"oneof=fixed permission_match"`
	PushUserList      []string `json:"push_user_list"`
	Enabled           bool     `json:"enabled"`
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
	req := new(UpdateReportPushConfigReqV1)
	err := controller.BindAndValidateReq(c, req)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	reportPushConfigId, err := strconv.Atoi(c.Param("report_push_config_id"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	s := model.GetStorage()
	config, exist, err := s.GetReportPushConfigById(uint(reportPushConfigId))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("report push configs %v not exist ,can't updatede", reportPushConfigId))
	}

	// 启停作为单独的行为，推送开关不一致，只做启停变更
	if req.Enabled != config.Enabled {
		config.Enabled = req.Enabled
	} else {
		if config.Type == model.TypeWorkflow {
			return controller.JSONBaseErrorReq(c, fmt.Errorf("report push configs %v update is not supported", config.Type))
		}
		config.TriggerType = req.TriggerType
		config.PushFrequencyCron = req.PushFrequencyCron
		config.PushUserType = req.PushUserType
		config.PushUserList = req.PushUserList
	}
	config.UpdateTime = time.Now()
	err = s.Save(config)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}
