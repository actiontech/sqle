package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

type GetCompanyNoticeResp struct {
	controller.BaseRes
	Data CompanyNotice `json:"data"`
}

type CompanyNotice struct {
	NoticeStr string `json:"notice_str"`
}

// GetCompanyNotice
// @Summary 获取企业公告
// @Description get company notice info
// @Id getCompanyNotice
// @Tags companyNotice
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetCompanyNoticeResp
// @Router /v1/company_notice [get]
func GetCompanyNotice(c echo.Context) error {
	return getCompanyNotice(c)
}

type UpdateCompanyNoticeReq struct {
	NoticeStr *string `json:"notice_str" valid:"omitempty"`
}

// UpdateCompanyNotice
// @Summary 更新企业公告
// @Description update company notice info
// @Id updateCompanyNotice
// @Tags companyNotice
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param companyNotice body v1.UpdateCompanyNoticeReq true "company notice"
// @Success 200 {object} controller.BaseRes
// @Router /v1/company_notice [patch]
func UpdateCompanyNotice(c echo.Context) error {
	return updateCompanyNotice(c)
}
