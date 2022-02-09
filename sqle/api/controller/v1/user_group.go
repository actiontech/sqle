package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
)

type CreateUserGroupReqV1 struct {
	Name  string   `json:"user_group_name" form:"user_group_name" example:"test" valid:"required,name"`
	Desc  string   `json:"user_group_desc" form:"user_group_desc" example:"this is a group"`
	Roles []string `json:"role_name_list" form:"role_name_list"`
	Users []string `json:"user_name_list" form:"user_name_list"`
}

// @Summary 创建用户组
// @Description create user group
// @Id CreateUserGroupV1
// @Tags user_group
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param instance body v1.CreateUserGroupReqV1 true "create user group"
// @Success 200 {object} controller.BaseRes
// @router /v1/user_groups [post]
func CreateUserGroup(c echo.Context) (err error) {

	req := new(CreateUserGroupReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	s := model.GetStorage()
	// check if user group already exist
	{
		_, isExist, err := s.GetUserGroupByName(req.Name)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if isExist {
			return controller.JSONNewDataExistErr(c, "user<%v> already exist", req.Name)
		}
	}

	// check users
	var users []*model.User
	{
		userNames := req.Users
		users, err = s.GetAndCheckUserExist(userNames)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	// check roles
	var roles []*model.Role
	{
		roleNames := req.Roles
		roles, err = s.GetAndCheckRoleExist(roleNames)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	// user group
	ug := &model.UserGroup{
		Name: req.Name,
		Desc: req.Desc,
	}

	if err := s.SaveUserGroupAndAssociations(ug, users, roles); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return controller.JSONBaseErrorReq(c, nil)
}

type GetUserGroupsResV1 struct {
	controller.BaseRes
	Data      []*UserGroupListItemResV1 `json:"data"`
	TotalNums uint64                    `json:"total_nums"`
}

type UserGroupListItemResV1 struct {
	Name       string   `json:"user_group_name"`
	Desc       string   `json:"user_group_desc"`
	IsDisabled bool     `json:"is_disabled,omitempty"`
	Users      []string `json:"user_name_list,omitempty"`
	Roles      []string `json:"role_name_list,omitempty"`
}

// @Summary 获取用户组列表
// @Description get user group info list
// @Id GetUserGroupsV1
// @Tags user_group
// @Id getUserGroupListV1
// @Security ApiKeyAuth
// @Param filter_user_group_name query string false "filter user group name"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Produce json
// @Success 200 {object} v1.GetUserGroupsResV1
// @router /v1/user_groups [get]
func GetUserGroups(c echo.Context) (err error) {
	// TODO: implementation
	return controller.JSONNewNotImplementedErr(c)
}

// @Summary 删除用户组
// @Description delete user group
// @Id deleteUserGroupV1
// @Tags user_group
// @Security ApiKeyAuth
// @Param user_group_name path string true "user_group_name"
// @Success 200 {object} controller.BaseRes
// @router /v1/user_groups/{user_group_name}/ [delete]
func DeleteUserGroup(c echo.Context) (err error) {

	userGroupName := c.Param("user_group_name")

	s := model.GetStorage()

	// check if user group exist
	var ug *model.UserGroup
	{
		var exist bool
		ug, exist, err = s.GetUserGroupByName(userGroupName)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if !exist {
			return controller.JSONBaseErrorReq(
				c, errors.NewDataNotExistErr("user group<%v> not exist", userGroupName))
		}
	}

	err = s.Delete(ug)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return controller.JSONBaseErrorReq(c, nil)
}

type PatchUserGroupReqV1 struct {
	Desc       string   `json:"user_group_desc" form:"user_group_desc" example:"this is a group"`
	Users      []string `json:"user_name_list" form:"user_name_list"`
	IsDisabled bool     `json:"is_disabled,omitempty" form:"is_disabled"`
	Roles      []string `json:"role_name_list" form:"role_name_list"`
}

// @Summary 更新用户组
// @Description update user group
// @Id updateUserGroupV1
// @Tags user_group
// @Security ApiKeyAuth
// @Param user_group_name path string true "user_group_name"
// @Success 200 {object} controller.BaseRes
// @router /v1/user_groups/{user_group_name}/ [patch]
func UpdateUserGroup(c echo.Context) error {
	// TODO: implementation
	return controller.JSONNewNotImplementedErr(c)
}

type UserGroupTipListItem struct {
	Name string `json:"user_group_name"`
}

type GetUserGroupTipsResV1 struct {
	controller.BaseRes
	Data []*UserGroupTipListItem `json:"data"`
}

// @Summary 获取用户组提示列表
// @Description get user group tip list
// @Tags user_group
// @Id getUserGroupTipListV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetUserGroupTipsResV1
// @router /v1/user_group_tips [get]
func GetUserGroupTips(c echo.Context) error {
	// TODO: implementation
	return controller.JSONNewNotImplementedErr(c)
}
