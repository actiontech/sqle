package v2

import (
	"github.com/actiontech/sqle/sqle/api/controller"

	"github.com/labstack/echo/v4"
)

type CreateRoleReqV2 struct {
	Name           string   `json:"role_name" form:"role_name" valid:"required,name"`
	Desc           string   `json:"role_desc" form:"role_desc"`
	Instances      []string `json:"instance_name_list" form:"instance_name_list"`
	OperationCodes []uint   `json:"operation_code_list" form:"operation_code_list"`
	Users          []string `json:"user_name_list" form:"user_name_list"`
	UserGroups     []string `json:"user_group_name_list" form:"user_group_name_list"`
}

// @Summary 创建角色
// @Description create role
// @Id createRoleV2
// @Tags role
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param instance body v2.CreateRoleReqV2 true "create role"
// @Success 200 {object} controller.BaseRes
// @router /v2/roles [post]
func CreateRole(c echo.Context) error {
	return controller.JSONNewNotImplementedErr(c)
}

type GetRolesReqV2 struct {
	FilterRoleName     string `json:"filter_role_name" query:"filter_role_name"`
	FilterUserName     string `json:"filter_user_name" query:"filter_user_name"`
	FilterInstanceName string `json:"filter_instance_name" query:"filter_instance_name"`
	PageIndex          uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize           uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetRolesResV2 struct {
	controller.BaseRes
	Data      []*RoleResV2 `json:"data"`
	TotalNums uint64       `json:"total_nums"`
}

type Operation struct {
	Code uint   `json:"operation_code"`
	Desc string `json:"operation_desc"`
}

type RoleResV2 struct {
	Name       string       `json:"role_name"`
	Desc       string       `json:"role_desc"`
	Users      []string     `json:"user_name_list,omitempty"`
	Instances  []string     `json:"instance_name_list,omitempty"`
	Operations []*Operation `json:"operation_list,omitempty"`
	UserGroups []string     `json:"user_group_name_list" form:"user_group_name_list"`
}

// @Summary 获取角色列表
// @Description get role list
// @Id getRoleListV2
// @Tags role
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param filter_role_name query string false "filter role name"
// @Param filter_user_name query string false "filter user name"
// @Param filter_instance_name query string false "filter instance name"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v2.GetRolesResV2
// @router /v2/roles [get]
func GetRoles(c echo.Context) error {
	return controller.JSONNewNotImplementedErr(c)
}

type UpdateRoleReqV2 struct {
	Desc           *string   `json:"role_desc" form:"role_desc"`
	Users          *[]string `json:"user_name_list,omitempty" form:"user_name_list"`
	Instances      []string  `json:"instance_name_list,omitempty" form:"instance_name_list"`
	OperationCodes []uint    `json:"operation_code_list" form:"operation_code_list"`
	UserGroups     *[]string `json:"user_group_name_list,omitempty" form:"user_group_name_list"`
}

// @Summary 更新角色信息
// @Description update role
// @Id updateRoleV2
// @Tags role
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param role_name path string true "role name"
// @Param instance body v2.UpdateRoleReqV2 true "update role request"
// @Success 200 {object} controller.BaseRes
// @router /v2/roles/{role_name}/ [patch]
func UpdateRole(c echo.Context) error {
	return controller.JSONNewNotImplementedErr(c)
}
