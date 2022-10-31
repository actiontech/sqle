package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"
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

// GetManagementPermissions
// @Summary 获取平台管理权限列表
// @Description get platform management permissions
// @Id GetManagementPermissionsV1
// @Tags management_permission
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetManagementPermissionsResV1
// @Router /v1/management_permissions [get]
func GetManagementPermissions(c echo.Context) error {
	return nil
}
