//go:build enterprise
// +build enterprise

package v1

import (
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/notification"

	"github.com/labstack/echo/v4"
)

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

func getSQLQueryConfiguration(c echo.Context) error {
	return c.JSON(http.StatusOK, GetSQLQueryConfigurationResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: GetSQLQueryConfigurationResDataV1{
			EnableSQLQuery:  false,
			SQLQueryRootURI: "/sql_query",
		},
	})
}
