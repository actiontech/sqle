package v1

import (
	"net/http"
	"strconv"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/driver"

	"github.com/actiontech/sqle/sqle/model"

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
	return getDingTalkConfigurationV1(c)
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
	return updateDingTalkConfigurationV1(c)
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
	return testDingTalkConfigV1(c)
}

type UpdateSystemVariablesReqV1 struct {
	Url                         *string `json:"url" form:"url" example:"http://10.186.61.32:8080" validate:"url"`
	OperationRecordExpiredHours *int    `json:"operation_record_expired_hours" form:"operation_record_expired_hours" example:"2160"`
	CbOperationLogsExpiredHours *int    `json:"cb_operation_logs_expired_hours" form:"cb_operation_logs_expired_hours" example:"2160"`
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

	if req.OperationRecordExpiredHours != nil {
		systemVariables = append(systemVariables, model.SystemVariable{
			Key:   model.SystemVariableOperationRecordExpiredHours,
			Value: strconv.Itoa(*req.OperationRecordExpiredHours),
		})
	}

	if req.CbOperationLogsExpiredHours != nil {
		systemVariables = append(systemVariables, model.SystemVariable{
			Key:   model.SystemVariableCbOperationLogsExpiredHours,
			Value: strconv.Itoa(*req.CbOperationLogsExpiredHours),
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
	Url                         string `json:"url"`
	OperationRecordExpiredHours int    `json:"operation_record_expired_hours"`
	CbOperationLogsExpiredHours int    `json:"cb_operation_logs_expired_hours"`
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
	operationRecordExpiredHours, err := strconv.Atoi(systemVariables[model.SystemVariableOperationRecordExpiredHours].Value)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	cbOperationLogsExpiredHours, err := strconv.Atoi(systemVariables[model.SystemVariableCbOperationLogsExpiredHours].Value)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetSystemVariablesResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: SystemVariablesResV1{
			Url:                         systemVariables[model.SystemVariableSqleUrl].Value,
			OperationRecordExpiredHours: operationRecordExpiredHours,
			CbOperationLogsExpiredHours: cbOperationLogsExpiredHours,
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

type GetFeishuAuditConfigurationResV1 struct {
	controller.BaseRes
	Data FeishuConfigurationV1 `json:"data"`
}

type FeishuConfigurationV1 struct {
	AppID                       string `json:"app_id"`
	IsFeishuNotificationEnabled bool   `json:"is_feishu_notification_enabled"`
}

// GetFeishuAuditConfigurationV1
// @Summary 获取飞书审核配置
// @Description get feishu audit configuration
// @Id getFeishuAuditConfigurationV1
// @Tags configuration
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetFeishuAuditConfigurationResV1
// @router /v1/configurations/feishu_audit [get]
func GetFeishuAuditConfigurationV1(c echo.Context) error {
	return getFeishuAuditConfigurationV1(c)
}

type UpdateFeishuConfigurationReqV1 struct {
	AppID                       *string `json:"app_id" form:"app_id" validate:"required" description:"飞书应用ID"`
	AppSecret                   *string `json:"app_secret" form:"app_secret" validate:"required" description:"飞书应用Secret"`
	IsFeishuNotificationEnabled *bool   `json:"is_feishu_notification_enabled" from:"is_feishu_notification_enabled" validate:"required" description:"是否启用飞书推送"`
}

// UpdateFeishuAuditConfigurationV1
// @Summary 添加或更新飞书配置
// @Description update feishu audit configuration
// @Accept json
// @Id updateFeishuAuditConfigurationV1
// @Tags configuration
// @Security ApiKeyAuth
// @Param param body v1.UpdateFeishuConfigurationReqV1 true "update feishu audit configuration req"
// @Success 200 {object} controller.BaseRes
// @router /v1/configurations/feishu_audit [patch]
func UpdateFeishuAuditConfigurationV1(c echo.Context) error {
	return updateFeishuAuditConfigurationV1(c)
}

type TestFeishuConfigurationReqV1 struct {
	AccountType string `json:"account_type" form:"account_type" enums:"email,phone" valid:"required"`
	Account     string `json:"account" form:"account" valid:"required" description:"绑定了飞书的手机号或邮箱"`
}

type TestFeishuConfigResDataV1 struct {
	IsMessageSentNormally bool   `json:"is_message_sent_normally"`
	ErrorMessage          string `json:"error_message,omitempty"`
}

type TestFeishuConfigResV1 struct {
	controller.BaseRes
	Data TestFeishuConfigResDataV1 `json:"data"`
}

const (
	FeishuAccountTypeEmail = "email"
	FeishuAccountTypePhone = "phone"
)

// TestFeishuAuditConfigV1
// @Summary 测试飞书审批配置
// @Description test feishu audit configuration
// @Accept json
// @Id testFeishuAuditConfigV1
// @Tags configuration
// @Security ApiKeyAuth
// @Param req body v1.TestFeishuConfigurationReqV1 true "test feishu configuration req"
// @Success 200 {object} v1.TestFeishuConfigResV1
// @router /v1/configurations/feishu_audit/test [post]
func TestFeishuAuditConfigV1(c echo.Context) error {
	return testFeishuAuditConfigV1(c)
}

type GetWechatAuditConfigurationResV1 struct {
	controller.BaseRes
	Data WechatConfigurationV1 `json:"data"`
}

type WechatConfigurationV1 struct {
	CorpID string `json:"corp_id"`

	IsWechatNotificationEnabled bool `json:"is_wechat_notification_enabled"`
}

type GetCodingAuditConfigurationResV1 struct {
	controller.BaseRes
	Data CodingConfigurationV1 `json:"data"`
}

type CodingConfigurationV1 struct {
	CodingUrl string `json:"coding_url"`

	IsCodingNotificationEnabled bool `json:"is_coding_notification_enabled"`
}

// GetWechatAuditConfigurationV1
// @Summary 获取微信审核配置
// @Description get wechat audit configuration
// @Id getWechatAuditConfigurationV1
// @Tags configuration
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetWechatAuditConfigurationResV1
// @router /v1/configurations/wechat_audit [get]
func GetWechatAuditConfigurationV1(c echo.Context) error {
	return getWechatAuditConfigurationV1(c)
}

// GetCodingAuditConfigurationV1
// @Summary 获取Coding审核配置
// @Description get coding audit configuration
// @Id getCodingAuditConfigurationV1
// @Tags configuration
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetCodingAuditConfigurationResV1
// @router /v1/configurations/coding_audit [get]
func GetCodingAuditConfigurationV1(c echo.Context) error {
	return getCodingAuditConfigurationV1(c)
}

type UpdateWechatConfigurationReqV1 struct {
	CorpID                      *string `json:"corp_id" from:"corp_id" description:"微信企业号ID"`
	CorpSecret                  *string `json:"corp_secret" from:"corp_secret" description:"企业微信ID对应密码"`
	IsWechatNotificationEnabled *bool   `json:"is_wechat_notification_enabled" from:"is_wechat_notification_enabled" validate:"required" description:"是否启用微信对接流程"`
}

// UpdateWechatAuditConfigurationV1
// @Summary 添加或更新微信配置
// @Description update wechat audit configuration
// @Accept json
// @Id updateWechatAuditConfigurationV1
// @Tags configuration
// @Security ApiKeyAuth
// @Param param body v1.UpdateWechatConfigurationReqV1 true "update wechat audit configuration req"
// @Success 200 {object} controller.BaseRes
// @router /v1/configurations/wechat_audit [patch]
func UpdateWechatAuditConfigurationV1(c echo.Context) error {
	return updateWechatAuditConfigurationV1(c)
}

type UpdateCodingAuditConfigurationReqV1 struct {
	CodingUrl                   *string `json:"coding_url" from:"coding_url" description:"Coding平台的地址"`
	Token                       *string `json:"token" from:"token" description:"访问令牌"`
	IsCodingNotificationEnabled *bool   `json:"is_coding_notification_enabled" from:"is_coding_notification_enabled" description:"是否启用Coding对接流程"`
}

// UpdateCodingAuditConfigurationV1
// @Summary 添加或更新Coding配置
// @Description update coding audit configuration
// @Accept json
// @Id UpdateCodingAuditConfigurationV1
// @Tags configuration
// @Security ApiKeyAuth
// @Param param body v1.UpdateCodingAuditConfigurationReqV1 true "update coding audit configuration req"
// @Success 200 {object} controller.BaseRes
// @router /v1/configurations/coding_audit [patch]
func UpdateCodingAuditConfigurationV1(c echo.Context) error {
	return updateCodingAuditConfigurationV1(c)
}

type TestWechatConfigResDataV1 struct {
	IsMessageSentNormally bool   `json:"is_message_sent_normally"`
	ErrorMessage          string `json:"error_message,omitempty"`
}

type TestWechatConfigResV1 struct {
	controller.BaseRes
	Data TestWechatConfigResDataV1 `json:"data"`
}

type TestCodingConfigResV1 struct {
	controller.BaseRes
	Data TestCodingConfigResDataV1 `json:"data"`
}

type TestCodingConfigResDataV1 struct {
	IsMessageSentNormally bool   `json:"is_message_sent_normally"`
	ErrorMessage          string `json:"error_message,omitempty"`
}

type TestWechatConfigurationReqV1 struct {
	WechatId string `json:"wechat_id" form:"wechat_id" valid:"required" description:"用户个人企业微信ID"`
}

// TestWechatAuditConfigV1
// @Summary 测试微信审批配置
// @Description test wechat audit configuration
// @Accept json
// @Id testWechatAuditConfigV1
// @Tags configuration
// @Security ApiKeyAuth
// @Param req body v1.TestWechatConfigurationReqV1 true "test wechat configuration req"
// @Success 200 {object} v1.TestWechatConfigResV1
// @router /v1/configurations/wechat_audit/test [post]
func TestWechatAuditConfigV1(c echo.Context) error {
	return testWechatAuditConfigV1(c)
}

type TestCodingConfigurationReqV1 struct {
	ProjectName string `json:"project_name" form:"project_name" valid:"required" description:"项目名称"`
}

// TestCodingAuditConfigV1
// @Summary 测试Coding审批配置
// @Description test coding audit configuration
// @Accept json
// @Id testCodingAuditConfigV1
// @Tags configuration
// @Param req body v1.TestCodingConfigurationReqV1 true "test coding configuration req"
// @Security ApiKeyAuth
// @Success 200 {object} v1.TestCodingConfigResV1
// @router /v1/configurations/coding_audit/test [post]
func TestCodingAuditConfigV1(c echo.Context) error {
	return testCodingAuditConfigV1(c)
}

type ScheduleTaskDefaultOption struct {
	DefaultSelector string `json:"default_selector" enums:"wechat,feishu"`
}

type ScheduledTaskDefaultOptionV1Rsp struct {
	controller.BaseRes
	Data ScheduleTaskDefaultOption `json:"data"`
}

// GetScheduledTaskDefaultOptionV1
// @Summary 获取工单定时上线二次确认默认选项
// @Description get scheduled task default option
// @Tags workflow
// @Id getScheduledTaskDefaultOptionV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.ScheduledTaskDefaultOptionV1Rsp
// @Router /v1/configurations/workflows/schedule/default_option [get]
func GetScheduledTaskDefaultOptionV1(c echo.Context) error {
	return getScheduledTaskDefaultOptionV1(c)
}
