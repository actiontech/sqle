package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/driver"

	"github.com/actiontech/sqle/sqle/model"

	"github.com/actiontech/sqle/sqle/pkg/im"
	"github.com/actiontech/sqle/sqle/pkg/im/dingding"

	"github.com/labstack/echo/v4"
)

type GetDingTalkConfigurationResV1 struct {
	controller.BaseRes
	Data DingTalkConfigurationV1 `json:"data"`
}

type DingTalkConfigurationV1 struct {
	AppKey                 string `json:"app_key"`
	IsEnableDingTalkNotify bool   `json:"is_enable_ding_talk_notify"`
}

// GetDingTalkConfigurationV1
// @Summary 获取 dingTalk 配置
// @Description get dingTalk configuration
// @Id getDingTalkConfigurationV1
// @Tags configuration
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetDingTalkConfigurationResV1
// @router /v1/configurations/ding_talk [get]
func GetDingTalkConfigurationV1(c echo.Context) error {
	s := model.GetStorage()
	dingTalk, exist, err := s.GetImConfigByType(model.ImTypeDingTalk)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return c.JSON(http.StatusOK, &GetDingTalkConfigurationResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data:    DingTalkConfigurationV1{},
		})
	}

	return c.JSON(http.StatusOK, &GetDingTalkConfigurationResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: DingTalkConfigurationV1{
			AppKey:                 dingTalk.AppKey,
			IsEnableDingTalkNotify: dingTalk.IsEnable,
		},
	})
}

type UpdateDingTalkConfigurationReqV1 struct {
	AppKey                 *string `json:"app_key" form:"app_key"  validate:"required"`
	AppSecret              *string `json:"app_secret" form:"app_secret"  validate:"required"`
	IsEnableDingTalkNotify *bool   `json:"is_enable_ding_talk_notify" from:"is_enable_ding_talk_notify" validate:"required" description:"是否启用钉钉通知"`
}

// UpdateDingTalkConfigurationV1
// @Summary 添加或更新 DingTalk 配置
// @Description update DingTalk configuration
// @Accept json
// @Id updateDingTalkConfigurationV1
// @Tags configuration
// @Security ApiKeyAuth
// @Param instance body v1.UpdateDingTalkConfigurationReqV1 true "update DingTalk configuration req"
// @Success 200 {object} controller.BaseRes
// @router /v1/configurations/ding_talk [patch]
func UpdateDingTalkConfigurationV1(c echo.Context) error {
	req := new(UpdateDingTalkConfigurationReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	dingTalk, _, err := s.GetImConfigByType(model.ImTypeDingTalk)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	{ // disable
		if req.IsEnableDingTalkNotify != nil && !(*req.IsEnableDingTalkNotify) {
			dingTalk.IsEnable = false
			return controller.JSONBaseErrorReq(c, s.Save(dingTalk))
		}
	}

	if req.AppKey != nil {
		dingTalk.AppKey = *req.AppKey
	}
	if req.AppSecret != nil {
		dingTalk.AppSecret = *req.AppSecret
	}
	if req.IsEnableDingTalkNotify != nil {
		dingTalk.IsEnable = *req.IsEnableDingTalkNotify
	}

	dingTalk.Type = model.ImTypeDingTalk

	if err := s.Save(dingTalk); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	go im.CreateApprovalTemplate(model.ImTypeDingTalk)

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type TestDingTalkConfigResDataV1 struct {
	IsDingTalkSendNormal bool   `json:"is_ding_talk_send_normal"`
	SendErrorMessage     string `json:"send_error_message,omitempty"`
}

type TestDingTalkConfigResV1 struct {
	controller.BaseRes
	Data TestDingTalkConfigResDataV1 `json:"data"`
}

// TestDingTalkConfigV1
// @Summary 测试 DingTalk 配置
// @Description test DingTalk configuration
// @Accept json
// @Id testDingTalkConfigV1
// @Tags configuration
// @Security ApiKeyAuth
// @Success 200 {object} v1.TestDingTalkConfigResV1
// @router /v1/configurations/ding_talk/test [post]
func TestDingTalkConfigV1(c echo.Context) error {
	s := model.GetStorage()
	dingTalk, exist, err := s.GetImConfigByType(model.ImTypeDingTalk)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return c.JSON(http.StatusOK, &TestDingTalkConfigResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data: TestDingTalkConfigResDataV1{
				IsDingTalkSendNormal: false,
				SendErrorMessage:     "dingTalk config not exist",
			},
		})
	}

	_, err = dingding.GetToken(dingTalk.AppKey, dingTalk.AppSecret)
	if err != nil {
		return c.JSON(http.StatusOK, &TestDingTalkConfigResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data: TestDingTalkConfigResDataV1{
				IsDingTalkSendNormal: false,
				SendErrorMessage:     err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, &TestDingTalkConfigResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: TestDingTalkConfigResDataV1{
			IsDingTalkSendNormal: true,
		},
	})
}

type UpdateSystemVariablesReqV1 struct {
	WorkflowExpiredHours        *int    `json:"workflow_expired_hours" form:"workflow_expired_hours" example:"720"`
	Url                         *string `json:"url" form:"url" example:"http://10.186.61.32:8080" validate:"url"`
	OperationRecordExpiredHours *int    `json:"operation_record_expired_hours" form:"operation_record_expired_hours" example:"2160"`
}

// @Summary 修改系统变量
// @Description update system variables
// @Accept json
// @Id updateSystemVariablesV1
// @Tags configuration
// @Security ApiKeyAuth
// @Param instance body v1.UpdateSystemVariablesReqV1 true "update system variables request"
// @Success 200 {object} controller.BaseRes
// @router /v1/configurations/system_variables [patch]
func UpdateSystemVariables(c echo.Context) error {
	req := new(UpdateSystemVariablesReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	s := model.GetStorage()

	var systemVariables []model.SystemVariable
	if req.WorkflowExpiredHours != nil {
		systemVariables = append(systemVariables, model.SystemVariable{
			Key:   model.SystemVariableWorkflowExpiredHours,
			Value: fmt.Sprintf("%v", *req.WorkflowExpiredHours),
		})
	}

	if req.OperationRecordExpiredHours != nil {
		systemVariables = append(systemVariables, model.SystemVariable{
			Key:   model.SystemVariableOperationRecordExpiredHours,
			Value: strconv.Itoa(*req.OperationRecordExpiredHours),
		})
	}

	if req.Url != nil {
		systemVariables = append(systemVariables, model.SystemVariable{
			Key:   model.SystemVariableSqleUrl,
			Value: *req.Url,
		})
	}

	if err := s.PathSaveSystemVariables(systemVariables); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return controller.JSONBaseErrorReq(c, nil)
}

type GetSystemVariablesResV1 struct {
	controller.BaseRes
	Data SystemVariablesResV1 `json:"data"`
}

type SystemVariablesResV1 struct {
	WorkflowExpiredHours        int    `json:"workflow_expired_hours"`
	Url                         string `json:"url"`
	OperationRecordExpiredHours int    `json:"operation_record_expired_hours"`
}

// @Summary 获取系统变量
// @Description get system variables
// @Id getSystemVariablesV1
// @Tags configuration
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetSystemVariablesResV1
// @router /v1/configurations/system_variables [get]
func GetSystemVariables(c echo.Context) error {
	s := model.GetStorage()
	systemVariables, err := s.GetAllSystemVariables()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	expiredHours, err := strconv.Atoi(systemVariables[model.SystemVariableWorkflowExpiredHours].Value)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	operationRecordExpiredHours, err := strconv.Atoi(systemVariables[model.SystemVariableOperationRecordExpiredHours].Value)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetSystemVariablesResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: SystemVariablesResV1{
			WorkflowExpiredHours:        expiredHours,
			Url:                         systemVariables[model.SystemVariableSqleUrl].Value,
			OperationRecordExpiredHours: operationRecordExpiredHours,
		},
	})
}

type GetDriversResV1 struct {
	controller.BaseRes
	Data DriversResV1 `json:"data"`
}

type DriversResV1 struct {
	Drivers []string `json:"driver_name_list"`
}

// GetDrivers get support Driver list.
// @Summary 获取当前 server 支持的审核类型
// @Description get drivers
// @Id getDriversV1
// @Tags configuration
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetDriversResV1
// @router /v1/configurations/drivers [get]
func GetDrivers(c echo.Context) error {
	return c.JSON(http.StatusOK, &GetDriversResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    DriversResV1{Drivers: driver.GetPluginManager().AllDrivers()},
	})
}

type GetSQLEInfoResV1 struct {
	controller.BaseRes
	Data GetSQLEInfoResDataV1 `json:"data"`
}

type GetSQLEInfoResDataV1 struct {
	Version string `json:"version"`
	LogoUrl string `json:"logo_url"`
	Title   string `json:"title"`
}

// GetSQLEInfo get sqle basic info
// @Summary 获取 sqle 基本信息
// @Description get sqle basic info
// @Id getSQLEInfoV1
// @Tags global
// @Success 200 {object} v1.GetSQLEInfoResV1
// @router /v1/basic_info [get]
func GetSQLEInfo(c echo.Context) error {
	return getSQLEInfo(c)
}

type UploadLogoResV1 struct {
	controller.BaseRes
	Data UploadLogoResDataV1 `json:"data"`
}

type UploadLogoResDataV1 struct {
	LogoUrl string `json:"logo_url"`
}

// UploadLogo
// @Summary 上传Logo
// @Description upload logo
// @Id uploadLogo
// @Tags configuration
// @Accept mpfd
// @Security ApiKeyAuth
// @Param logo formData file true "logo file"
// @Success 200 {object} v1.UploadLogoResV1
// @router /v1/configurations/personalise/logo [post]
func UploadLogo(c echo.Context) error {
	return uploadLogo(c)
}

type GetLogoReqV1 struct {
	Timestamp string `query:"timestamp"`
}

// GetLogo
// @Summary 获取logo
// @Description get logo
// @Id getLogo
// @Tags configuration
// @Param timestamp query string false "timestamp"
// @Success 200 {file} file "get logo"
// @router /v1/static/logo [get]
func GetLogo(c echo.Context) error {
	return getLogo(c)
}

type PersonaliseReqV1 struct {
	Title *string `json:"title"`
}

// UpdatePersonaliseConfig
// @Summary 更新个性化设置
// @Description update personalise config
// @Id personalise
// @Tags configuration
// @Security ApiKeyAuth
// @Param conf body v1.PersonaliseReqV1 true "personalise req"
// @Success 200 {object} controller.BaseRes
// @router /v1/configurations/personalise [patch]
func UpdatePersonaliseConfig(c echo.Context) error {
	return updatePersonaliseConfig(c)
}
