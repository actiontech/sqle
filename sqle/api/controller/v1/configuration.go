package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/service"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/notification"
	"github.com/actiontech/sqle/sqle/pkg/im"
	"github.com/actiontech/sqle/sqle/pkg/im/dingding"
	"github.com/actiontech/sqle/sqle/pkg/im/feishu"

	"github.com/labstack/echo/v4"
)

type UpdateSMTPConfigurationReqV1 struct {
	EnableSMTPNotify *bool   `json:"enable_smtp_notify" from:"enable_smtp_notify" description:"是否启用邮件通知"`
	Host             *string `json:"smtp_host" form:"smtp_host" example:"smtp.email.qq.com"`
	Port             *string `json:"smtp_port" form:"smtp_port" example:"465" valid:"omitempty,port"`
	Username         *string `json:"smtp_username" form:"smtp_username" example:"test@qq.com" valid:"omitempty,email"`
	Password         *string `json:"smtp_password" form:"smtp_password" example:"123"`
	IsSkipVerify     *bool   `json:"is_skip_verify" form:"is_skip_verify" description:"是否跳过安全认证"`
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
	if req.IsSkipVerify != nil {
		smtpC.IsSkipVerify = *req.IsSkipVerify
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
	IsSkipVerify     bool   `json:"is_skip_verify"`
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
			IsSkipVerify:     smtpC.IsSkipVerify,
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

	s := model.GetStorage()
	smtpC, exist, err := s.GetSMTPConfiguration()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return c.JSON(http.StatusOK, &TestSMTPConfigurationResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data:    newSmtpConnectableResV1(false, "SMTP is not configured"),
		})
	}

	if !smtpC.EnableSMTPNotify.Bool {
		return c.JSON(http.StatusOK, &TestSMTPConfigurationResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data:    newSmtpConnectableResV1(false, "SMTP notice is not enabled"),
		})
	}

	addr := req.RecipientAddr
	notifier := &notification.EmailNotifier{}
	err = notifier.Notify(&notification.TestNotify{}, []*model.User{
		{
			Email: addr,
		},
	})
	if err != nil {
		return c.JSON(http.StatusOK, &TestSMTPConfigurationResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data:    newSmtpConnectableResV1(false, err.Error()),
		})
	}

	return c.JSON(http.StatusOK, &TestSMTPConfigurationResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    newSmtpConnectableResV1(true, "ok"),
	})
}

func newSmtpConnectableResV1(isSendNormal bool, errMsg string) TestSMTPConfigurationResDataV1 {
	return TestSMTPConfigurationResDataV1{
		IsSMTPSendNormal: isSendNormal,
		SendErrorMessage: errMsg,
	}
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

func testWeChatConfigurationV1(c echo.Context) error {
	req := new(TestWeChatConfigurationReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	testID := req.RecipientID
	notifier := &notification.WeChatNotifier{}
	err := notifier.Notify(&notification.TestNotify{}, []*model.User{
		{
			Name:     testID,
			WeChatID: testID,
		},
	})
	if err != nil {
		return c.JSON(http.StatusOK, &TestWeChatConfigurationResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data: TestWeChatConfigurationResDataV1{
				IsWeChatSendNormal: false,
				SendErrorMessage:   err.Error(),
			},
		})
	}
	return c.JSON(http.StatusOK, &TestWeChatConfigurationResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: TestWeChatConfigurationResDataV1{
			IsWeChatSendNormal: true,
			SendErrorMessage:   "ok",
		},
	})
}

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

type GetFeishuConfigurationResV1 struct {
	controller.BaseRes
	Data FeishuConfigurationV1 `json:"data"`
}

type FeishuConfigurationV1 struct {
	AppID                       string `json:"app_id"`
	IsFeishuNotificationEnabled bool   `json:"is_feishu_notification_enabled"`
}

// GetFeishuConfigurationV1
// @Summary 获取飞书配置
// @Description get feishu configuration
// @Id getFeishuConfigurationV1
// @Tags configuration
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetFeishuConfigurationResV1
// @router /v1/configurations/feishu [get]
func GetFeishuConfigurationV1(c echo.Context) error {
	s := model.GetStorage()
	feishuCfg, exist, err := s.GetImConfigByType(model.ImTypeFeishu)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return c.JSON(http.StatusOK, &GetFeishuConfigurationResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data:    FeishuConfigurationV1{},
		})
	}

	return c.JSON(http.StatusOK, &GetFeishuConfigurationResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: FeishuConfigurationV1{
			AppID:                       feishuCfg.AppKey,
			IsFeishuNotificationEnabled: feishuCfg.IsEnable,
		},
	})
}

type UpdateFeishuConfigurationReqV1 struct {
	AppID                       *string `json:"app_id" form:"app_id"`
	AppSecret                   *string `json:"app_secret" form:"app_secret" `
	IsFeishuNotificationEnabled *bool   `json:"is_feishu_notification_enabled" from:"is_feishu_notification_enabled" description:"是否启用飞书推送"`
}

// UpdateFeishuConfigurationV1
// @Summary 添加或更新飞书配置
// @Description update feishu configuration
// @Accept json
// @Id updateFeishuConfigurationV1
// @Tags configuration
// @Security ApiKeyAuth
// @Param param body v1.UpdateFeishuConfigurationReqV1 true "update feishu configuration req"
// @Success 200 {object} controller.BaseRes
// @router /v1/configurations/feishu [patch]
func UpdateFeishuConfigurationV1(c echo.Context) error {
	req := new(UpdateFeishuConfigurationReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()
	feishuCfg, _, err := s.GetImConfigByType(model.ImTypeFeishu)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if req.AppID != nil {
		feishuCfg.AppKey = *req.AppID
	}
	if req.AppSecret != nil {
		feishuCfg.AppSecret = *req.AppSecret
	}
	if req.IsFeishuNotificationEnabled != nil {
		feishuCfg.IsEnable = *req.IsFeishuNotificationEnabled
	}
	feishuCfg.Type = model.ImTypeFeishu

	if err := s.Save(feishuCfg); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return controller.JSONBaseErrorReq(c, nil)
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

// TestFeishuConfigV1
// @Summary 测试飞书配置
// @Description test feishu configuration
// @Accept json
// @Id testFeishuConfigV1
// @Tags configuration
// @Security ApiKeyAuth
// @Param req body v1.TestFeishuConfigurationReqV1 true "test feishu configuration req"
// @Success 200 {object} v1.TestFeishuConfigResV1
// @router /v1/configurations/feishu/test [post]
func TestFeishuConfigV1(c echo.Context) error {
	req := new(TestFeishuConfigurationReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var email, phone []string
	switch req.AccountType {
	case FeishuAccountTypeEmail:
		err := controller.Validate(struct {
			Email string `valid:"email"`
		}{req.Account})
		if err != nil {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, err))
		}
		email = append(email, req.Account)
	case FeishuAccountTypePhone:
		phone = append(phone, req.Account)
	default:
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf("unknown account type: %v", req.AccountType)))
	}

	s := model.GetStorage()
	feishuCfg, exist, err := s.GetImConfigByType(model.ImTypeFeishu)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return c.JSON(http.StatusOK, &TestFeishuConfigResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data: TestFeishuConfigResDataV1{
				IsMessageSentNormally: false,
				ErrorMessage:          "feishu configuration doesn't exist",
			},
		})
	}

	client := feishu.NewFeishuClient(feishuCfg.AppKey, feishuCfg.AppSecret)
	feishuUsers, err := client.GetUsersByEmailOrMobileWithLimitation(email, phone)
	if err != nil {
		return c.JSON(http.StatusOK, &TestFeishuConfigResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data: TestFeishuConfigResDataV1{
				IsMessageSentNormally: false,
				ErrorMessage:          fmt.Sprintf("get user_ids failed: %v", err),
			},
		})
	}

	if len(feishuUsers) == 0 {
		return c.JSON(http.StatusOK, &TestFeishuConfigResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data: TestFeishuConfigResDataV1{
				IsMessageSentNormally: false,
				ErrorMessage:          "can not find matched feishu user",
			},
		})
	}

	n := &notification.TestNotify{}
	content, err := notification.BuildFeishuMessageBody(n)
	if err != nil {
		return c.JSON(http.StatusOK, &TestFeishuConfigResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data: TestFeishuConfigResDataV1{
				IsMessageSentNormally: false,
				ErrorMessage:          fmt.Sprintf("convert content failed: %v", err),
			},
		})
	}
	for uid := range feishuUsers {
		if err = client.SendMessage(feishu.FeishuReceiverIdTypeUserId, uid, feishu.FeishuSendMessageMsgTypePost, content); err != nil {
			return c.JSON(http.StatusOK, &TestFeishuConfigResV1{
				BaseRes: controller.NewBaseReq(nil),
				Data: TestFeishuConfigResDataV1{
					IsMessageSentNormally: false,
					ErrorMessage:          err.Error(),
				},
			})
		}
	}

	return c.JSON(http.StatusOK, &TestFeishuConfigResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: TestFeishuConfigResDataV1{
			IsMessageSentNormally: true,
		},
	})

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

func updateWeChatConfigurationV1(c echo.Context) error {
	req := new(UpdateWeChatConfigurationReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()
	wechatC, _, err := s.GetWeChatConfiguration()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if req.CorpID != nil {
		wechatC.CorpID = *req.CorpID
	}
	if req.CorpSecret != nil {
		wechatC.CorpSecret = *req.CorpSecret
	}
	if req.AgentID != nil {
		wechatC.AgentID = *req.AgentID
	}
	if req.ProxyIP != nil {
		wechatC.ProxyIP = *req.ProxyIP
	}
	if req.EnableWeChatNotify != nil {
		wechatC.EnableWeChatNotify = *req.EnableWeChatNotify
	}
	if req.SafeEnabled != nil {
		wechatC.SafeEnabled = *req.SafeEnabled
	}

	if err := s.Save(wechatC); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return controller.JSONBaseErrorReq(c, nil)
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

func getWeChatConfiguration(c echo.Context) error {
	s := model.GetStorage()
	wechatC, _, err := s.GetWeChatConfiguration()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, &GetWeChatConfigurationResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: WeChatConfigurationResV1{
			EnableWeChatNotify: wechatC.EnableWeChatNotify,
			CorpID:             wechatC.CorpID,
			AgentID:            wechatC.AgentID,
			SafeEnabled:        wechatC.SafeEnabled,
			ProxyIP:            wechatC.ProxyIP,
		},
	})
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
	WorkflowExpiredHours *int    `json:"workflow_expired_hours" form:"workflow_expired_hours" example:"720"`
	Url                  *string `json:"url" form:"url" example:"http://10.186.61.32:8080" validate:"url"`
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
	WorkflowExpiredHours int    `json:"workflow_expired_hours"`
	Url                  string `json:"url"`
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

	return c.JSON(http.StatusOK, &GetSystemVariablesResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: SystemVariablesResV1{
			WorkflowExpiredHours: expiredHours,
			Url:                  systemVariables[model.SystemVariableSqleUrl].Value,
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

type GetSQLQueryConfigurationResV1 struct {
	controller.BaseRes
	Data GetSQLQueryConfigurationResDataV1 `json:"data"`
}

type GetSQLQueryConfigurationResDataV1 struct {
	EnableSQLQuery  bool   `json:"enable_sql_query"`
	SQLQueryRootURI string `json:"sql_query_root_uri"`
}

// @Summary 获取SQL查询配置信息
// @Description get sqle query configuration
// @Id getSQLQueryConfiguration
// @Tags configuration
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetSQLQueryConfigurationResV1
// @router /v1/configurations/sql_query [get]
func GetSQLQueryConfiguration(c echo.Context) error {
	return c.JSON(http.StatusOK, GetSQLQueryConfigurationResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: GetSQLQueryConfigurationResDataV1{
			EnableSQLQuery:  service.IsCloudBeaverConfigured(),
			SQLQueryRootURI: service.CbRootUri,
		},
	})
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

type WebHookConfigV1 struct {
	Enable               *bool   `json:"enable" description:"是否启用"`
	MaxRetryTimes        *int    `json:"max_retry_times" description:"最大重试次数"`
	RetryIntervalSeconds *int    `json:"retry_interval_seconds" description:"请求重试间隔"`
	Token                *string `json:"token" description:"token 令牌"`
	URL                  *string `json:"url" description:"回调API URL"`
}

// UpdateWorkflowWebHookConfig
// @Summary 更新工单 WebHook 配置
// @Description update webhook config
// @Id updateGlobalWebHookConfig
// @Tags configuration
// @Security ApiKeyAuth
// @Param request body v1.WebHookConfigV1 true "update webhook config"
// @Success 200 {object} controller.BaseRes
// @Router /v1/configurations/webhook [patch]
func UpdateWorkflowWebHookConfig(c echo.Context) error {
	req := new(WebHookConfigV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()
	cfg, _, err := s.GetWorkflowWebHookConfig()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if req.Enable != nil {
		cfg.Enable = *req.Enable
	}
	if req.MaxRetryTimes != nil {
		if *req.MaxRetryTimes < 0 || *req.MaxRetryTimes > 5 {
			err = errors.NewDataInvalidErr(
				"ouf of range[0-5] for max_retry_times[%v]", *req.MaxRetryTimes)
			return controller.JSONBaseErrorReq(c, err)
		}
		cfg.MaxRetryTimes = *req.MaxRetryTimes
	}
	if req.RetryIntervalSeconds != nil {
		if *req.RetryIntervalSeconds < 1 || *req.RetryIntervalSeconds > 5 {
			err = errors.NewDataInvalidErr(
				"out of range[1-5] for retry_interval_seconds[%v]", *req.RetryIntervalSeconds)
			return controller.JSONBaseErrorReq(c, err)
		}
		cfg.RetryIntervalSeconds = *req.RetryIntervalSeconds
	}
	if req.Token != nil {
		cfg.Token = *req.Token
	}
	if req.URL != nil {
		if !strings.HasPrefix(*req.URL, "http://") {
			err = errors.NewDataInvalidErr("url must start with 'http://'")
			return controller.JSONBaseErrorReq(c, err)
		}
		cfg.URL = *req.URL
	}
	return controller.JSONBaseErrorReq(c, s.Save(cfg))
}

type GetWorkflowWebHookConfigResV1 struct {
	controller.BaseRes
	Data WebHookConfigV1 `json:"data"`
}

// GetWorkflowWebHookConfig
// @Summary 获取全局工单 WebHook 配置
// @Description get workflow webhook config
// @Id getGlobalWorkflowWebHookConfig
// @Tags configuration
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetWorkflowWebHookConfigResV1
// @Router /v1/configurations/webhook [get]
func GetWorkflowWebHookConfig(c echo.Context) error {
	s := model.GetStorage()
	cfg, _, err := s.GetWorkflowWebHookConfig()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, &GetWorkflowWebHookConfigResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: WebHookConfigV1{
			Enable:               &cfg.Enable,
			MaxRetryTimes:        &cfg.MaxRetryTimes,
			RetryIntervalSeconds: &cfg.RetryIntervalSeconds,
			Token:                &cfg.Token,
			URL:                  &cfg.URL,
		},
	})
}

type TestWorkflowWebHookConfigResDataV1 struct {
	SendErrorMessage string `json:"send_error_message,omitempty"`
}

type TestWorkflowWebHookConfigResV1 struct {
	controller.BaseRes
	Data TestWorkflowWebHookConfigResDataV1 `json:"data"`
}

// TestWorkflowWebHookConfig
// @Summary 测试全局工单 WebHook 配置
// @Description test workflow webhook config
// @Id testGlobalWorkflowWebHookConfig
// @Tags configuration
// @Security ApiKeyAuth
// @Success 200 {object} v1.TestWorkflowWebHookConfigResV1
// @Router /v1/configurations/webhook/test [post]
func TestWorkflowWebHookConfig(c echo.Context) error {
	data := &TestWorkflowWebHookConfigResDataV1{}
	err := notification.TestWorkflowConfig()
	if err != nil {
		data.SendErrorMessage = err.Error()
	}
	return c.JSON(http.StatusOK, &TestWorkflowWebHookConfigResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    *data,
	})
}
