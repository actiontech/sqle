package v1

import (
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
)

type GetManagementPermissionsResV1 struct {
	controller.BaseRes
	Data []*ManagementPermissionResV1 `json:"data"`
}

type ManagementPermissionResV1 struct {
	Code uint   `json:"code"`
	Desc string `json:"desc"`
}

func generateManagementPermissionResV1(code uint) ManagementPermissionResV1 {
	return ManagementPermissionResV1{
		Code: code,
		Desc: model.GetManagementPermissionDesc(code),
	}
}

func generateManagementPermissionResV1s(code []uint) []*ManagementPermissionResV1 {
	m := []*ManagementPermissionResV1{}
	for _, u := range code {
		p := generateManagementPermissionResV1(u)
		m = append(m, &p)
	}
	return m
}

// GetManagementPermissions
// @Summary 获取平台管理权限列表
// @Description get platform management permissions
// @Id GetManagementPermissionsV1
// @Tags management_permission
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetManagementPermissionsResV1
// @Router /v1/management_permissions [get]
func GetManagementPermissions(c echo.Context) error {
	p := model.GetManagementPermission()
	data := []*ManagementPermissionResV1{}
	for u, s := range p {
		data = append(data, &ManagementPermissionResV1{
			Code: u,
			Desc: s,
		})
	}

	return c.JSON(http.StatusOK, GetManagementPermissionsResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    data,
	})
}
