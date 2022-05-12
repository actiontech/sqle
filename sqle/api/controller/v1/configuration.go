package v1

import (
	"fmt"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/notification"

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
	if req.EnableSMTPNotify != nil {
		// It is never possible to trigger an error here
		_ = smtpC.EnableSMTPNotify.Scan(*req.EnableSMTPNotify)
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
			EnableSMTPNotify: smtpC.EnableSMTPNotify.Bool,
			Host:             smtpC.Host,
			Port:             smtpC.Port,
			Username:         smtpC.Username,
		},
	})
}

type TestSMTPConfigurationReqV1 struct {
	RecipientAddr string `json:"recipient_addr" from:"recipient_addr" description:"消息接收者邮箱地址" valid:"required,email"`
}

type TestSMTPConfigurationResV1 struct {
	controller.BaseRes
	Data TestSMTPConfigurationResDataV1 `json:"data"`
}

type TestSMTPConfigurationResDataV1 struct {
	IsSMTPSendNormal bool   `json:"is_smtp_send_normal"`
	SendErrorMessage string `json:"send_error_message,omitempty"`
}

// TestSMTPConfigurationV1 used to test SMTP notifications
// @Summary 测试 邮箱 配置
// @Description test SMTP configuration
// @Accept json
// @Id testSMTPConfigurationV1
// @Tags configuration
// @Security ApiKeyAuth
// @Param req body v1.TestSMTPConfigurationReqV1 true "test SMTP configuration req"
// @Success 200 {object} v1.TestSMTPConfigurationResV1
// @router /v1/configurations/smtp/test [post]
func TestSMTPConfigurationV1(c echo.Context) error {
	req := new(TestSMTPConfigurationReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	addr := req.RecipientAddr
	notifier := &notification.EmailNotifier{}
	err := notifier.Notify(&notification.TestNotify{}, []*model.User{
		{
			Email: addr,
		},
	})
	if err != nil {
		return c.JSON(http.StatusOK, &TestSMTPConfigurationResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data: TestSMTPConfigurationResDataV1{
				IsSMTPSendNormal: false,
				SendErrorMessage: err.Error(),
			},
		})
	}
	return c.JSON(http.StatusOK, &TestSMTPConfigurationResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: TestSMTPConfigurationResDataV1{
			IsSMTPSendNormal: true,
			SendErrorMessage: "ok",
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
	EnableSSL           *bool   `json:"enable_ssl"`
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
	EnableSSL           bool   `json:"enable_ssl"`
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
			EnableSSL:           ldapC.EnableSSL,
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

		if req.EnableSSL != nil {
			ldapC.EnableSSL = *req.EnableSSL
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

type GetOauth2ConfigurationResV1 struct {
	controller.BaseRes
	Data GetOauth2ConfigurationResDataV1 `json:"data"`
}

type GetOauth2ConfigurationResDataV1 struct {
	EnableOauth2    bool     `json:"enable_oauth2"`
	ClientID        string   `json:"client_id"`
	ClientHost      string   `json:"client_host"`
	ServerAuthUrl   string   `json:"server_auth_url"`
	ServerTokenUrl  string   `json:"server_token_url"`
	ServerUserIdUrl string   `json:"server_user_id_url"`
	Scopes          []string `json:"scopes"`
	AccessTokenTag  string   `json:"access_token_tag"`
	UserIdTag       string   `json:"user_id_tag"`
	LoginTip        string   `json:"login_tip"`
}

// @Summary 获取 Oauth2 配置
// @Description get Oauth2 configuration
// @Id getOauth2ConfigurationV1
// @Tags configuration
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetOauth2ConfigurationResV1
// @router /v1/configurations/oauth2 [get]
func GetOauth2Configuration(c echo.Context) error {
	s := model.GetStorage()
	oauth2C, _, err := s.GetOauth2Configuration()
	return c.JSON(http.StatusOK, &GetOauth2ConfigurationResV1{
		BaseRes: controller.NewBaseReq(err),
		Data: GetOauth2ConfigurationResDataV1{
			EnableOauth2:    oauth2C.EnableOauth2,
			ClientID:        oauth2C.ClientID,
			ClientHost:      oauth2C.ClientHost,
			ServerAuthUrl:   oauth2C.ServerAuthUrl,
			ServerTokenUrl:  oauth2C.ServerTokenUrl,
			ServerUserIdUrl: oauth2C.ServerUserIdUrl,
			Scopes:          oauth2C.GetScopes(),
			AccessTokenTag:  oauth2C.AccessTokenTag,
			UserIdTag:       oauth2C.UserIdTag,
			LoginTip:        oauth2C.LoginTip,
		},
	})
}

type Oauth2ConfigurationReqV1 struct {
	EnableOauth2    *bool     `json:"enable_oauth2"`
	ClientID        *string   `json:"client_id"`
	ClientKey       *string   `json:"client_key"`
	ClientHost      *string   `json:"client_host"`
	ServerAuthUrl   *string   `json:"server_auth_url"`
	ServerTokenUrl  *string   `json:"server_token_url"`
	ServerUserIdUrl *string   `json:"server_user_id_url"`
	Scopes          *[]string `json:"scopes"`
	AccessTokenTag  *string   `json:"access_token_tag"`
	UserIdTag       *string   `json:"user_id_tag"`
	LoginTip        *string   `json:"login_tip"`
}

// @Summary 修改 Oauth2 配置
// @Description update Oauth2 configuration
// @Accept json
// @Id updateOauth2ConfigurationV1
// @Tags configuration
// @Security ApiKeyAuth
// @Param conf body v1.Oauth2ConfigurationReqV1 true "update Oauth2 configuration req"
// @Success 200 {object} controller.BaseRes
// @router /v1/configurations/oauth2 [patch]
func UpdateOauth2Configuration(c echo.Context) error {
	req := new(Oauth2ConfigurationReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()
	oauth2C, _, err := s.GetOauth2Configuration()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	{ // patch oauth2 config
		if req.EnableOauth2 != nil {
			oauth2C.EnableOauth2 = *req.EnableOauth2
		}
		if req.ClientID != nil {
			oauth2C.ClientID = *req.ClientID
		}
		if req.ClientKey != nil {
			oauth2C.ClientKey = *req.ClientKey
		}
		if req.ClientHost != nil {
			oauth2C.ClientHost = *req.ClientHost
		}
		if req.ServerAuthUrl != nil {
			oauth2C.ServerAuthUrl = *req.ServerAuthUrl
		}
		if req.ServerTokenUrl != nil {
			oauth2C.ServerTokenUrl = *req.ServerTokenUrl
		}
		if req.ServerUserIdUrl != nil {
			oauth2C.ServerUserIdUrl = *req.ServerUserIdUrl
		}
		if req.Scopes != nil {
			oauth2C.SetScopes(*req.Scopes)
		}
		if req.AccessTokenTag != nil {
			oauth2C.AccessTokenTag = *req.AccessTokenTag
		}
		if req.UserIdTag != nil {
			oauth2C.UserIdTag = *req.UserIdTag
		}
		if req.LoginTip != nil {
			oauth2C.LoginTip = *req.LoginTip
		}

	}

	return controller.JSONBaseErrorReq(c, s.Save(oauth2C))
}

type GetOauth2TipsResV1 struct {
	controller.BaseRes
	Data GetOauth2TipsResDataV1 `json:"data"`
}

type GetOauth2TipsResDataV1 struct {
	EnableOauth2 bool   `json:"enable_oauth2"`
	LoginTip     string `json:"login_tip"`
}

// @Summary 获取 Oauth2 基本信息
// @Description get Oauth2 tips
// @Id getOauth2Tips
// @Tags configuration
// @Success 200 {object} v1.GetOauth2TipsResV1
// @router /v1/configurations/oauth2/tips [get]
func GetOauth2Tips(c echo.Context) error {
	s := model.GetStorage()
	oauth2C, _, err := s.GetOauth2Configuration()
	return c.JSON(http.StatusOK, &GetOauth2TipsResV1{
		BaseRes: controller.NewBaseReq(err),
		Data: GetOauth2TipsResDataV1{
			EnableOauth2: oauth2C.EnableOauth2,
			LoginTip:     oauth2C.LoginTip,
		},
	})
}
