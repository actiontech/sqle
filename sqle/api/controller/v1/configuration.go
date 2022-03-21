package v1

import (
	"fmt"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
)

type UpdateSMTPConfigurationReqV1 struct {
	EnableSMTPNotify *bool   `json:"enable_smtp_notify" from:"enable_smtp_notify" description:"是否启用邮件通知"`
	Host             *string `json:"smtp_host" form:"smtp_host" example:"smtp.email.qq.com"`
	Port             *string `json:"smtp_port" form:"smtp_port" example:"465" valid:"omitempty,port"`
	Username         *string `json:"smtp_username" form:"smtp_username" example:"test@qq.com" valid:"omitempty,email"`
	Password         *string `json:"smtp_password" form:"smtp_password" example:"123"`
}

// @Summary 添加 SMTP 配置
// @Description update SMTP configuration
// @Accept json
// @Id updateSMTPConfigurationV1
// @Tags configuration
// @Security ApiKeyAuth
// @Param instance body v1.UpdateSMTPConfigurationReqV1 true "update SMTP configuration req"
// @Success 200 {object} controller.BaseRes
// @router /v1/configurations/smtp [patch]
func UpdateSMTPConfiguration(c echo.Context) error {
	req := new(UpdateSMTPConfigurationReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()
	smtpC, _, err := s.GetSMTPConfiguration()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if req.Host != nil {
		smtpC.Host = *req.Host
	}
	if req.Port != nil {
		smtpC.Port = *req.Port
	}
	if req.Username != nil {
		smtpC.Username = *req.Username
	}
	if req.Password != nil {
		smtpC.Password = *req.Password
	}

	if err := s.Save(smtpC); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return controller.JSONBaseErrorReq(c, nil)
}

type GetSMTPConfigurationResV1 struct {
	controller.BaseRes
	Data SMTPConfigurationResV1 `json:"data"`
}

type SMTPConfigurationResV1 struct {
	EnableSMTPNotify bool   `json:"enable_smtp_notify"`
	Host             string `json:"smtp_host"`
	Port             string `json:"smtp_port"`
	Username         string `json:"smtp_username"`
}

// @Summary 获取 SMTP 配置
// @Description get SMTP configuration
// @Id getSMTPConfigurationV1
// @Tags configuration
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetSMTPConfigurationResV1
// @router /v1/configurations/smtp [get]
func GetSMTPConfiguration(c echo.Context) error {
	s := model.GetStorage()
	smtpC, _, err := s.GetSMTPConfiguration()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, &GetSMTPConfigurationResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: SMTPConfigurationResV1{
			Host:     smtpC.Host,
			Port:     smtpC.Port,
			Username: smtpC.Username,
		},
	})
}

type TestWeChatConfigurationReqV1 struct {
	RecipientID string `json:"recipient_id" from:"recipient_id" description:"消息接收者企业微信ID"`
}

type TestWeChatConfigurationResV1 struct {
	controller.BaseRes
	Data TestWeChatConfigurationResDataV1 `json:"data"`
}

type TestWeChatConfigurationResDataV1 struct {
	IsWeChatSendNormal bool   `json:"is_wechat_send_normal"`
	SendErrorMessage   string `json:"send_error_message,omitempty"`
}

// TestWeChatConfigurationV1 used to test WeChat notifications
// @Summary 测试 企业微信 配置
// @Description test WeChat configuration
// @Accept json
// @Id testWeChatConfigurationV1
// @Tags configuration
// @Security ApiKeyAuth
// @Param instance body v1.TestWeChatConfigurationReqV1 true "test WeChat configuration req"
// @Success 200 {object} v1.TestWeChatConfigurationResV1
// @router /v1/configurations/wechat/test [post]
func TestWeChatConfigurationV1(c echo.Context) error {
	return testWeChatConfigurationV1(c)
}

type UpdateWeChatConfigurationReqV1 struct {
	EnableWeChatNotify *bool   `json:"enable_wechat_notify" from:"enable_wechat_notify" description:"是否启用微信通知"`
	CorpID             *string `json:"corp_id" from:"corp_id" description:"企业微信ID"`
	CorpSecret         *string `json:"corp_secret" from:"corp_secret" description:"企业微信ID对应密码"`
	AgentID            *int    `json:"agent_id" from:"agent_id" description:"企业微信应用ID"`
	SafeEnabled        *bool   `json:"safe_enabled" from:"safe_enabled" description:"是否对传输信息加密"`
	ProxyIP            *string `json:"proxy_ip" from:"proxy_ip" description:"企业微信代理服务器IP"`
}

// UpdateWeChatConfigurationV1 used to configure WeChat notifications
// @Summary 添加 企业微信 配置
// @Description update WeChat configuration
// @Accept json
// @Id updateWeChatConfigurationV1
// @Tags configuration
// @Security ApiKeyAuth
// @Param instance body v1.UpdateWeChatConfigurationReqV1 true "update WeChat configuration req"
// @Success 200 {object} controller.BaseRes
// @router /v1/configurations/wechat [patch]
func UpdateWeChatConfigurationV1(c echo.Context) error {
	return updateWeChatConfigurationV1(c)
}

type GetWeChatConfigurationResV1 struct {
	controller.BaseRes
	Data WeChatConfigurationResV1 `json:"data"`
}

type WeChatConfigurationResV1 struct {
	EnableWeChatNotify bool   `json:"enable_wechat_notify"`
	CorpID             string `json:"corp_id"`
	AgentID            int    `json:"agent_id"`
	SafeEnabled        bool   `json:"safe_enabled"`
	ProxyIP            string `json:"proxy_ip"`
}

// GetWeChatConfiguration used to get wechat configure
// @Summary 获取 企业微信 配置
// @Description get WeChat configuration
// @Id getWeChatConfigurationV1
// @Tags configuration
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetWeChatConfigurationResV1
// @router /v1/configurations/wechat [get]
func GetWeChatConfiguration(c echo.Context) error {
	return getWeChatConfiguration(c)
}

type GetLDAPConfigurationResV1 struct {
	controller.BaseRes
	Data LDAPConfigurationResV1 `json:"data"`
}

type LDAPConfigurationReqV1 struct {
	EnableLdap          *bool   `json:"enable_ldap"`
	LdapServerHost      *string `json:"ldap_server_host"`
	LdapServerPort      *string `json:"ldap_server_port"`
	LdapConnectDn       *string `json:"ldap_connect_dn"`
	LdapConnectPwd      *string `json:"ldap_connect_pwd"`
	LdapSearchBaseDn    *string `json:"ldap_search_base_dn"`
	LdapUserNameRdnKey  *string `json:"ldap_user_name_rdn_key"`
	LdapUserEmailRdnKey *string `json:"ldap_user_email_rdn_key"`
}

type LDAPConfigurationResV1 struct {
	EnableLdap          bool   `json:"enable_ldap"`
	LdapServerHost      string `json:"ldap_server_host"`
	LdapServerPort      string `json:"ldap_server_port"`
	LdapConnectDn       string `json:"ldap_connect_dn"`
	LdapSearchBaseDn    string `json:"ldap_search_base_dn"`
	LdapUserNameRdnKey  string `json:"ldap_user_name_rdn_key"`
	LdapUserEmailRdnKey string `json:"ldap_user_email_rdn_key"`
}

// @Summary 获取 LDAP 配置
// @Description get LDAP configuration
// @Id getLDAPConfigurationV1
// @Tags configuration
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetLDAPConfigurationResV1
// @router /v1/configurations/ldap [get]
func GetLDAPConfiguration(c echo.Context) error {
	s := model.GetStorage()
	ldapC, _, err := s.GetLDAPConfiguration()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, &GetLDAPConfigurationResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: LDAPConfigurationResV1{
			EnableLdap:          ldapC.Enable,
			LdapServerHost:      ldapC.Host,
			LdapServerPort:      ldapC.Port,
			LdapConnectDn:       ldapC.ConnectDn,
			LdapSearchBaseDn:    ldapC.BaseDn,
			LdapUserNameRdnKey:  ldapC.UserNameRdnKey,
			LdapUserEmailRdnKey: ldapC.UserEmailRdnKey,
		},
	})
}

// @Summary 添加 LDAP 配置
// @Description update LDAP configuration
// @Accept json
// @Id updateLDAPConfigurationV1
// @Tags configuration
// @Security ApiKeyAuth
// @Param instance body v1.LDAPConfigurationReqV1 true "update LDAP configuration req"
// @Success 200 {object} controller.BaseRes
// @router /v1/configurations/ldap [patch]
func UpdateLDAPConfiguration(c echo.Context) error {
	req := new(LDAPConfigurationReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()
	ldapC, _, err := s.GetLDAPConfiguration()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	{ // patch ldap config

		if req.EnableLdap != nil {
			ldapC.Enable = *req.EnableLdap
		}

		if req.LdapServerHost != nil {
			ldapC.Host = *req.LdapServerHost
		}

		if req.LdapServerPort != nil {
			ldapC.Port = *req.LdapServerPort
		}

		if req.LdapConnectDn != nil {
			ldapC.ConnectDn = *req.LdapConnectDn
		}

		if req.LdapConnectPwd != nil {
			ldapC.ConnectPassword = *req.LdapConnectPwd
		}

		if req.LdapSearchBaseDn != nil {
			ldapC.BaseDn = *req.LdapSearchBaseDn
		}

		if req.LdapUserNameRdnKey != nil {
			ldapC.UserNameRdnKey = *req.LdapUserNameRdnKey
		}

		if req.LdapUserEmailRdnKey != nil {
			ldapC.UserEmailRdnKey = *req.LdapUserEmailRdnKey
		}

	}
	if err := s.Save(ldapC); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return controller.JSONBaseErrorReq(c, nil)
}

type UpdateSystemVariablesReqV1 struct {
	WorkflowExpiredHours *int `json:"workflow_expired_hours" form:"workflow_expired_hours" example:"720"`
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

	if req.WorkflowExpiredHours != nil {
		sv := &model.SystemVariable{
			Key:   model.SystemVariableWorkflowExpiredHours,
			Value: fmt.Sprintf("%v", *req.WorkflowExpiredHours)}

		if err := s.Save(sv); err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}
	return controller.JSONBaseErrorReq(c, nil)
}

type GetSystemVariablesResV1 struct {
	controller.BaseRes
	Data SystemVariablesResV1 `json:"data"`
}

type SystemVariablesResV1 struct {
	WorkflowExpiredHours int `json:"workflow_expired_hours"`
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
	wfExpiredHours, err := s.GetWorkflowExpiredHoursOrDefault()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetSystemVariablesResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: SystemVariablesResV1{
			WorkflowExpiredHours: int(wfExpiredHours),
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
		Data:    DriversResV1{Drivers: driver.AllDrivers()},
	})
}

type GetSQLEInfoResV1 struct {
	controller.BaseRes
	Version string `json:"version"`
}

// GetSQLEInfo get sqle basic info
// @Summary 获取 sqle 基本信息
// @Description get sqle basic info
// @Id getSQLEInfoV1
// @Tags global
// @Success 200 {object} v1.GetSQLEInfoResV1
// @router /v1/basic_info [get]
func GetSQLEInfo(c echo.Context) error {
	return c.JSON(http.StatusOK, &GetSQLEInfoResV1{
		BaseRes: controller.NewBaseReq(nil),
		Version: config.Version,
	})
}
