//go:build enterprise
// +build enterprise

package v1

import (
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
)

func getCompanyNotice(c echo.Context) error {
	s := model.GetStorage()
	notice, err := s.GetCompanyNotice()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetCompanyNoticeResp{
		BaseRes: controller.NewBaseReq(nil),
		Data: CompanyNotice{
			NoticeStr: notice.NoticeStr,
		},
	})
}

func updateCompanyNotice(c echo.Context) error {
	req := new(UpdateCompanyNoticeReq)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if req.NoticeStr == nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
	}

	s := model.GetStorage()
	notice, err := s.GetCompanyNotice()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	notice.NoticeStr = *req.NoticeStr

	return controller.JSONBaseErrorReq(c, s.Save(&notice))
}
