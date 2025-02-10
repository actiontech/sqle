//go:build enterprise
// +build enterprise

package v1

import (
	"fmt"
	"github.com/actiontech/sqle/sqle/pkg/im/coding"
	"net/http"
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/im"
	"github.com/actiontech/sqle/sqle/pkg/im/dingding"
	"github.com/actiontech/sqle/sqle/pkg/im/feishu"
	"github.com/actiontech/sqle/sqle/pkg/im/wechat"
	"github.com/labstack/echo/v4"
	larkContact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
)

const (
	ImTypeFeishu = "feishu"
	ImTypeWechat = "wechat"
)

func updateFeishuAuditConfigurationV1(c echo.Context) error {
	req := new(UpdateFeishuConfigurationReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()
	feishuCfg, _, err := s.GetImConfigByType(model.ImTypeFeishuAudit)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	{ // disable
		if req.IsFeishuNotificationEnabled != nil && !(*req.IsFeishuNotificationEnabled) {
			feishuCfg.IsEnable = false
			return controller.JSONBaseErrorReq(c, s.Save(feishuCfg))
		}
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

	// 如果是新创建的飞书审批配置，需要设置type
	feishuCfg.Type = model.ImTypeFeishuAudit

	if err := s.Save(feishuCfg); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	go im.CreateApprovalTemplate(model.ImTypeFeishuAudit)

	return controller.JSONBaseErrorReq(c, nil)
}

func getFeishuAuditConfigurationV1(c echo.Context) error {
	s := model.GetStorage()
	feishuCfg, exist, err := s.GetImConfigByType(model.ImTypeFeishuAudit)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return c.JSON(http.StatusOK, &GetFeishuAuditConfigurationResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data:    FeishuConfigurationV1{},
		})
	}

	return c.JSON(http.StatusOK, &GetFeishuAuditConfigurationResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: FeishuConfigurationV1{
			AppID:                       feishuCfg.AppKey,
			IsFeishuNotificationEnabled: feishuCfg.IsEnable,
		},
	})
}

func testFeishuAuditConfigV1(c echo.Context) error {
	req := new(TestFeishuConfigurationReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var email, phone string
	switch req.AccountType {
	case FeishuAccountTypeEmail:
		err := controller.Validate(struct {
			Email string `valid:"email"`
		}{req.Account})
		if err != nil {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, err))
		}
		email = req.Account
	case FeishuAccountTypePhone:
		phone = req.Account
	default:
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf("unknown account type: %v", req.AccountType)))
	}

	s := model.GetStorage()
	feishuCfg, exist, err := s.GetImConfigByType(model.ImTypeFeishuAudit)
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
	feishuUsers, err := client.GetUsersByEmailOrMobileWithLimitation([]string{email}, []string{phone}, larkContact.UserIdTypeOpenId)
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
				ErrorMessage:          "no user found in feishu",
			},
		})
	}

	for uid := range feishuUsers {
		_, err := client.CreateApprovalInstance(c.Request().Context(), feishuCfg.ProcessCode, locale.Bundle.LocalizeMsgByCtx(c.Request().Context(), locale.ConfigTestAudit), uid,
			[]string{uid}, "", "", locale.Bundle.LocalizeMsgByCtx(c.Request().Context(), locale.ConfigFeishuTestContent), "", "")
		if err != nil {
			return c.JSON(http.StatusOK, &TestFeishuConfigResV1{
				BaseRes: controller.NewBaseReq(nil),
				Data: TestFeishuConfigResDataV1{
					IsMessageSentNormally: false,
					ErrorMessage:          fmt.Sprintf("send approval failed: %v", err),
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

func getDingTalkConfigurationV1(c echo.Context) error {
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

func updateDingTalkConfigurationV1(c echo.Context) error {
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

func testDingTalkConfigV1(c echo.Context) error {
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

func getWechatAuditConfigurationV1(c echo.Context) error {
	s := model.GetStorage()
	wechat, exist, err := s.GetImConfigByType(model.ImTypeWechatAudit)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return c.JSON(http.StatusOK, &GetWechatAuditConfigurationResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data:    WechatConfigurationV1{},
		})
	}

	return c.JSON(http.StatusOK, &GetWechatAuditConfigurationResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: WechatConfigurationV1{
			CorpID:                      wechat.AppKey,
			IsWechatNotificationEnabled: wechat.IsEnable,
		},
	})
}

func updateWechatAuditConfigurationV1(c echo.Context) error {
	req := new(UpdateWechatConfigurationReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	wechat, _, err := s.GetImConfigByType(model.ImTypeWechatAudit)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if req.IsWechatNotificationEnabled != nil && !(*req.IsWechatNotificationEnabled) {
		wechat.IsEnable = false
		return controller.JSONBaseErrorReq(c, s.Save(wechat))
	}

	var isChanged bool
	if req.CorpID != nil {
		if wechat.AppKey != *req.CorpID {
			wechat.AppKey = *req.CorpID
			isChanged = true
		}
	}
	if req.CorpSecret != nil {
		if wechat.AppSecret != *req.CorpSecret {
			wechat.AppSecret = *req.CorpSecret
			isChanged = true
		}
	}
	if req.IsWechatNotificationEnabled != nil {
		if wechat.IsEnable != *req.IsWechatNotificationEnabled {
			wechat.IsEnable = *req.IsWechatNotificationEnabled
			isChanged = true
		}
	}

	wechat.Type = model.ImTypeWechatAudit

	if isChanged {
		if err := s.Save(wechat); err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	go im.CreateApprovalTemplate(model.ImTypeWechatAudit)

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

func testWechatAuditConfigV1(c echo.Context) error {
	req := new(TestWechatConfigurationReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	s := model.GetStorage()
	wechatCfg, exist, err := s.GetImConfigByType(model.ImTypeWechatAudit)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return c.JSON(http.StatusOK, &TestWechatConfigResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data: TestWechatConfigResDataV1{
				IsMessageSentNormally: false,
				ErrorMessage:          "wechat configuration doesn't exist",
			},
		})
	}

	client := wechat.NewWechatClient(wechatCfg.AppKey, wechatCfg.AppSecret)
	_, err = client.CreateApprovalInstance(c.Request().Context(), wechatCfg.ProcessCode, "", req.WechatId, []string{req.WechatId},
		"", "", locale.Bundle.LocalizeMsgByCtx(c.Request().Context(), locale.ConfigFeishuTestContent), nil)

	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &TestWechatConfigResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: TestWechatConfigResDataV1{
			IsMessageSentNormally: true,
		},
	})
}

func getCodingConfigurationV1(c echo.Context) error {
	s := model.GetStorage()
	codingConfig, exist, err := s.GetImConfigByType(model.ImTypeCoding)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return c.JSON(http.StatusOK, &GetCodingConfigurationResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data:    CodingConfigurationV1{},
		})
	}

	return c.JSON(http.StatusOK, &GetCodingConfigurationResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: CodingConfigurationV1{
			CodingUrl:       codingConfig.AppKey,
			IsCodingEnabled: codingConfig.IsEnable,
		},
	})
}

func updateCodingConfigurationV1(c echo.Context) error {
	req := new(UpdateCodingConfigurationReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	s := model.GetStorage()
	codingCfg, _, err := s.GetImConfigByType(model.ImTypeCoding)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if req.IsCodingEnabled != nil && !(*req.IsCodingEnabled) {
		codingCfg.IsEnable = false
		return controller.JSONBaseErrorReq(c, s.Save(codingCfg))
	}
	var isChanged bool
	if req.CodingUrl != nil {
		if codingCfg.AppKey != *req.CodingUrl {
			codingCfg.AppKey = *req.CodingUrl
			isChanged = true
		}
	}
	if req.Token != nil {
		if codingCfg.AppSecret != *req.Token {
			codingCfg.AppSecret = *req.Token
			isChanged = true
		}
	}
	if req.IsCodingEnabled != nil {
		if codingCfg.IsEnable != *req.IsCodingEnabled {
			codingCfg.IsEnable = *req.IsCodingEnabled
			isChanged = true
		}
	}
	codingCfg.Type = model.ImTypeCoding
	if isChanged {
		if err := s.Save(codingCfg); err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

func testCodingAuditConfigV1(c echo.Context) error {
	req := new(TestCodingConfigurationReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()
	codingCfg, exist, err := s.GetImConfigByType(model.ImTypeCoding)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return c.JSON(http.StatusOK, &TestCodingConfigResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data: TestCodingConfigResDataV1{
				IsMessageSentNormally: false,
				ErrorMessage:          "coding configuration doesn't exist",
			},
		})
	}
	issueName := fmt.Sprintf("SQLE_TEST_%s", time.Now().Format("20060102150405"))
	description := locale.Bundle.LocalizeMsgByCtx(c.Request().Context(), locale.ConfigCoding)
	_, err = coding.NewCodingClient(codingCfg.AppKey, codingCfg.AppSecret).CreateIssue(coding.CreateIssueRequestBody{Name: issueName, Priority: CodingPriorityLow.Weight(), ProjectName: req.CodingProjectName, Type: string(CodingTypeMission), Description: description, WorkingHours: 24})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, &TestCodingConfigResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: TestCodingConfigResDataV1{
			IsMessageSentNormally: true,
		},
	})
}

func getScheduledTaskDefaultOptionV1(c echo.Context) error {
	s := model.GetStorage()
	wir, exist, err := s.GetLastNeedNotifyScheduledRecord()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return c.JSON(http.StatusOK, &ScheduledTaskDefaultOptionV1Rsp{
			BaseRes: controller.NewBaseReq(nil),
			Data:    ScheduleTaskDefaultOption{},
		})
	}

	fr, err := s.GetFeishuRecordsByTaskIds([]uint{wir.TaskId})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if len(fr) > 0 {
		return c.JSON(http.StatusOK, &ScheduledTaskDefaultOptionV1Rsp{
			BaseRes: controller.NewBaseReq(nil),
			Data: ScheduleTaskDefaultOption{
				DefaultSelector: ImTypeFeishu,
			},
		})
	}

	wr, err := s.GetWechatRecordsByTaskIds([]uint{wir.TaskId})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if len(wr) > 0 {
		return c.JSON(http.StatusOK, &ScheduledTaskDefaultOptionV1Rsp{
			BaseRes: controller.NewBaseReq(nil),
			Data: ScheduleTaskDefaultOption{
				DefaultSelector: ImTypeWechat,
			},
		})
	}

	log.NewEntry().Error("failed to retrieve default option: Unable to locate the 'im' type. This issue may be attributed to the data truncation in the ScheduledRecord table. ")
	return c.JSON(http.StatusOK, &ScheduledTaskDefaultOptionV1Rsp{
		BaseRes: controller.NewBaseReq(nil),
		Data:    ScheduleTaskDefaultOption{},
	})
}
