package v1

import (
	"fmt"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
)

// @Summary 删除角色
// @Description delete role
// @Id deleteRoleV1
// @Tags role
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param role_name path string true "role name"
// @Success 200 {object} controller.BaseRes
// @router /v1/roles/{role_name}/ [delete]
func DeleteRole(c echo.Context) error {
	roleName := c.Param("role_name")
	s := model.GetStorage()
	role, exist, err := s.GetRoleByName(roleName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("role is not exist")))
	}

	return controller.JSONBaseErrorReq(c,
		s.DeleteRoleAndAssociations(role))
}

type CreateRoleReqV1 struct {
	Name      string   `json:"role_name" form:"role_name" valid:"required,name"`
	Desc      string   `json:"role_desc" form:"role_desc"`
	Users     []string `json:"user_name_list" form:"user_name_list"`
	Instances []string `json:"instance_name_list" form:"instance_name_list"`
}

// @Summary 创建角色
// @Description create role
// @Deprecated
// @Id createRoleV1
// @Tags role
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param instance body v1.CreateRoleReqV1 true "create role"
// @Success 200 {object} controller.BaseRes
// @router /v1/roles [post]
func CreateRole(c echo.Context) error {
	req := new(CreateRoleReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	s := model.GetStorage()
	_, exist, err := s.GetRoleByName(req.Name)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("role is exist")))
	}

	var users []*model.User
	if req.Users != nil || len(req.Users) > 0 {
		users, err = s.GetAndCheckUserExist(req.Users)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	var instances []*model.Instance
	if req.Instances != nil || len(req.Instances) > 0 {
		instances, err = s.GetAndCheckInstanceExist(req.Instances)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}
	role := &model.Role{
		Name: req.Name,
		Desc: req.Desc,
	}
	err = s.Save(role)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = s.UpdateRoleUsers(role, users...)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = s.UpdateRoleInstances(role, instances...)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return controller.JSONBaseErrorReq(c, nil)
}

type UpdateRoleReqV1 struct {
	Desc      *string  `json:"role_desc" form:"role_desc"`
	Users     []string `json:"user_name_list" form:"user_name_list"`
	Instances []string `json:"instance_name_list" form:"instance_name_list"`
}

// @Summary 更新角色信息
// @Description update role
// @Deprecated
// @Id updateRoleV1
// @Tags role
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param role_name path string true "role name"
// @Param instance body v1.UpdateRoleReqV1 true "update role request"
// @Success 200 {object} controller.BaseRes
// @router /v1/roles/{role_name}/ [patch]
func UpdateRole(c echo.Context) error {
	req := new(UpdateRoleReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	roleName := c.Param("role_name")
	s := model.GetStorage()
	role, exist, err := s.GetRoleByName(roleName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("role is not exist")))
	}

	if req.Users != nil || len(req.Users) > 0 {
		users, err := s.GetAndCheckUserExist(req.Users)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		err = s.UpdateRoleUsers(role, users...)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	if req.Instances != nil || len(req.Instances) > 0 {
		instances, err := s.GetAndCheckInstanceExist(req.Instances)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		err = s.UpdateRoleInstances(role, instances...)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	if req.Desc != nil {
		role.Desc = *req.Desc
	}

	err = s.Save(role)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return controller.JSONBaseErrorReq(c, nil)
}

type GetRolesReqV1 struct {
	FilterRoleName     string `json:"filter_role_name" query:"filter_role_name"`
	FilterUserName     string `json:"filter_user_name" query:"filter_user_name"`
	FilterInstanceName string `json:"filter_instance_name" query:"filter_instance_name"`
	PageIndex          uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize           uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetRolesResV1 struct {
	controller.BaseRes
	Data      []RoleResV1 `json:"data"`
	TotalNums uint64      `json:"total_nums"`
}

type RoleResV1 struct {
	Name      string   `json:"role_name"`
	Desc      string   `json:"role_desc"`
	Users     []string `json:"user_name_list,omitempty"`
	Instances []string `json:"instance_name_list,omitempty"`
}

// @Summary 获取角色列表
// @Description get role list
// @Deprecated
// @Id getRoleListV1
// @Tags role
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param filter_role_name query string false "filter role name"
// @Param filter_user_name query string false "filter user name"
// @Param filter_instance_name query string false "filter instance name"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v1.GetRolesResV1
// @router /v1/roles [get]
func GetRoles(c echo.Context) error {
	req := new(GetRolesReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()

	var offset uint32
	if req.PageIndex >= 1 {
		offset = req.PageSize * (req.PageIndex - 1)
	}
	data := map[string]interface{}{
		"filter_role_name":     req.FilterRoleName,
		"filter_user_name":     req.FilterUserName,
		"filter_instance_name": req.FilterInstanceName,
		"limit":                req.PageSize,
		"offset":               offset,
	}

	roles, count, err := s.GetRolesByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	rolesRes := make([]RoleResV1, 0, len(roles))
	for _, role := range roles {
		roleRes := RoleResV1{
			Name:      role.Name,
			Desc:      role.Desc,
			Users:     role.UserNames,
			Instances: role.InstanceNames,
		}
		rolesRes = append(rolesRes, roleRes)
	}
	return c.JSON(http.StatusOK, &GetRolesResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      rolesRes,
		TotalNums: count,
	})
}

type RoleTipResV1 struct {
	Name string `json:"role_name"`
}

type GetRoleTipsResV1 struct {
	controller.BaseRes
	Data []RoleTipResV1 `json:"data"`
}

// @Summary 获取角色提示列表
// @Description get role tip list
// @Tags role
// @Id getRoleTipListV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetRoleTipsResV1
// @router /v1/role_tips [get]
func GetRoleTips(c echo.Context) error {
	s := model.GetStorage()
	roles, err := s.GetAllEnabledRoles()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	roleTipsRes := make([]RoleTipResV1, 0, len(roles))

	for _, role := range roles {
		roleTipRes := RoleTipResV1{
			Name: role.Name,
		}
		roleTipsRes = append(roleTipsRes, roleTipRes)
	}
	return c.JSON(http.StatusOK, &GetRoleTipsResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    roleTipsRes,
	})
}
