package v1

import (
	"actiontech.cloud/universe/sqle/v4/sqle/api/controller"
	"actiontech.cloud/universe/sqle/v4/sqle/errors"
	"actiontech.cloud/universe/sqle/v4/sqle/model"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type CreateUserReqV1 struct {
	Name     string   `json:"user_name" form:"user_name" example:"test" valid:"required"`
	Password string   `json:"user_password" form:"user_name" example:"123456" valid:"required"`
	Email    string   `json:"email" form:"email" example:"test@email.com" valid:"email"`
	Roles    []string `json:"role_name_list" form:"role_name_list"`
}

// @Summary 创建用户
// @Description create user
// @Id createUserV1
// @Tags user
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param instance body v1.CreateUserReqV1 true "create user"
// @Success 200 {object} controller.BaseRes
// @router /v1/users [post]
func CreateUser(c echo.Context) error {
	req := new(CreateUserReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	s := model.GetStorage()
	_, exist, err := s.GetUserByName(req.Name)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DATA_EXIST, fmt.Errorf("user is exist")))
	}

	var roles []*model.Role
	if req.Roles != nil || len(req.Roles) > 0 {
		roles, err = s.GetAndCheckRoleExist(req.Roles)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}
	user := &model.User{
		Name:     req.Name,
		Password: req.Password,
		Email:    req.Email,
	}
	err = s.Save(user)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = s.UpdateUserRoles(user, roles...)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return controller.JSONBaseErrorReq(c, nil)
}

type UpdateUserReqV1 struct {
	Email *string  `json:"email"`
	Roles []string `json:"role_name_list" form:"role_name_list"`
}

// @Summary 更新用户信息
// @Description update user
// @Id updateUserV1
// @Tags user
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param user_name path string true "user name"
// @Param instance body v1.UpdateUserReqV1 true "update user"
// @Success 200 {object} controller.BaseRes
// @router /v1/users/{user_name}/ [patch]
func UpdateUser(c echo.Context) error {
	req := new(UpdateUserReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	userName := c.Param("user_name")
	s := model.GetStorage()
	user, exist, err := s.GetUserByName(userName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DATA_NOT_EXIST, fmt.Errorf("user is not exist")))
	}

	if req.Roles != nil || len(req.Roles) > 0 {
		roles, err := s.GetAndCheckRoleExist(req.Roles)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		err = s.UpdateUserRoles(user, roles...)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	if req.Email != nil {
		user.Email = *req.Email
		err = s.Save(user)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}
	return controller.JSONBaseErrorReq(c, nil)
}

// @Summary 删除用户
// @Description delete user
// @Id deleteUserV1
// @Tags user
// @Security ApiKeyAuth
// @Param user_name path string true "user name"
// @Success 200 {object} controller.BaseRes
// @router /v1/users/{user_name}/ [delete]
func DeleteUser(c echo.Context) error {
	userName := c.Param("user_name")
	s := model.GetStorage()
	user, exist, err := s.GetUserByName(userName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DATA_NOT_EXIST, fmt.Errorf("user is not exist")))
	}
	err = s.Delete(user)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return controller.JSONBaseErrorReq(c, nil)
}

type GetUserDetailResV1 struct {
	controller.BaseRes
	Data UserDetailResV1 `json:"data"`
}

type UserDetailResV1 struct {
	Name    string   `json:"user_name"`
	Email   string   `json:"email"`
	IsAdmin bool     `json:"is_admin"`
	Roles   []string `json:"role_name_list,omitempty"`
}

func convertUserToRes(user *model.User) UserDetailResV1 {
	userReq := UserDetailResV1{
		Name:    user.Name,
		Email:   user.Email,
		IsAdmin: user.Name == defaultAdminUser,
	}
	roleNames := make([]string, 0, len(user.Roles))
	for _, role := range user.Roles {
		roleNames = append(roleNames, role.Name)
	}
	userReq.Roles = roleNames
	return userReq
}

// @Summary 获取用户信息
// @Description get user info
// @Id getUserV1
// @Tags user
// @Security ApiKeyAuth
// @Param user_name path string true "user name"
// @Success 200 {object} v1.GetUserDetailResV1
// @router /v1/users/{user_name}/ [get]
func GetUser(c echo.Context) error {
	userName := c.Param("user_name")
	s := model.GetStorage()
	user, exist, err := s.GetUserDetailByName(userName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DATA_NOT_EXIST, fmt.Errorf("user is not exist")))
	}
	return c.JSON(http.StatusOK, &GetUserDetailResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertUserToRes(user),
	})
}

// @Summary 获取当前用户信息
// @Description get current user info
// @Id getCurrentUserV1
// @Tags user
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetUserDetailResV1
// @router /v1/user [get]
func GetCurrentUser(c echo.Context) error {
	userName := controller.GetUserName(c)
	s := model.GetStorage()
	user, exist, err := s.GetUserDetailByName(userName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DATA_NOT_EXIST, fmt.Errorf("user is not exist")))
	}
	return c.JSON(http.StatusOK, &GetUserDetailResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertUserToRes(user),
	})
}

type GetUsersReqV1 struct {
	FilterUserName string `json:"filter_user_name" query:"filter_user_name"`
	FilterRoleName string `json:"filter_role_name" query:"filter_role_name"`
	PageIndex      uint32 `json:"page_index" query:"page_index" valid:"required,int"`
	PageSize       uint32 `json:"page_size" query:"page_size" valid:"required,int"`
}

type GetUsersResV1 struct {
	controller.BaseRes
	Data      []UserResV1 `json:"data"`
	TotalNums uint64      `json:"total_nums"`
}

type UserResV1 struct {
	Name  string   `json:"user_name"`
	Email string   `json:"email"`
	Roles []string `json:"role_name_list,omitempty"`
}

// @Summary 获取用户信息列表
// @Description get user info list
// @Tags user
// @Id getUserListV1
// @Security ApiKeyAuth
// @Param filter_user_name query string false "filter user name"
// @Param filter_role_name query string false "filter role name"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v1.GetUsersResV1
// @router /v1/users [get]
func GetUsers(c echo.Context) error {
	req := new(GetUsersReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()

	var offset uint32
	if req.PageIndex >= 1 {
		offset = req.PageSize * (req.PageIndex - 1)
	}
	data := map[string]interface{}{
		"filter_user_name": req.FilterUserName,
		"filter_role_name": req.FilterRoleName,
		"limit":            req.PageSize,
		"offset":           offset,
	}

	users, count, err := s.GetUsersByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	usersReq := []UserResV1{}
	for _, user := range users {
		userReq := UserResV1{
			Name:  user.Name,
			Email: user.Email,
		}
		if user.RoleNames != "" {
			userReq.Roles = strings.Split(user.RoleNames, ",")
		}
		usersReq = append(usersReq, userReq)
	}
	return c.JSON(http.StatusOK, &GetUsersResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      usersReq,
		TotalNums: count,
	})
}

type UserTipResV1 struct {
	Name string `json:"user_name"`
}

type GetUserTipsResV1 struct {
	controller.BaseRes
	Data []UserTipResV1 `json:"data"`
}

// @Summary 获取用户提示列表
// @Description get user tip list
// @Tags user
// @Id getUserTipListV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetUserTipsResV1
// @router /v1/user_tips [get]
func GetUserTips(c echo.Context) error {
	s := model.GetStorage()
	users, err := s.GetAllUserTip()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	userTipsRes := make([]UserTipResV1, 0, len(users))

	for _, user := range users {
		userTipRes := UserTipResV1{
			Name: user.Name,
		}
		userTipsRes = append(userTipsRes, userTipRes)
	}
	return c.JSON(http.StatusOK, &GetUserTipsResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    userTipsRes,
	})
}

type CreateRoleReqV1 struct {
	Name      string   `json:"role_name" form:"role_name" valid:"required"`
	Desc      string   `json:"role_desc" form:"role_desc"`
	Users     []string `json:"user_name_list" form:"user_name_list"`
	Instances []string `json:"instance_name_list" form:"instance_name_list"`
}

// @Summary 创建角色
// @Description create role
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
		return controller.JSONBaseErrorReq(c, errors.New(errors.DATA_EXIST, fmt.Errorf("role is exist")))
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
		return controller.JSONBaseErrorReq(c, errors.New(errors.DATA_NOT_EXIST, fmt.Errorf("role is not exist")))
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
		return controller.JSONBaseErrorReq(c, errors.New(errors.DATA_NOT_EXIST, fmt.Errorf("role is not exist")))
	}
	err = s.Delete(role)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return controller.JSONBaseErrorReq(c, nil)
}

type GetRolesReqV1 struct {
	FilterRoleName     string `json:"filter_role_name" query:"filter_role_name"`
	FilterUserName     string `json:"filter_user_name" query:"filter_user_name"`
	FilterInstanceName string `json:"filter_instance_name" query:"filter_instance_name"`
	PageIndex          uint32 `json:"page_index" query:"page_index" valid:"required,int"`
	PageSize           uint32 `json:"page_size" query:"page_size" valid:"required,int"`
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
			Name: role.Name,
			Desc: role.Desc,
		}
		if role.UserNames != "" {
			roleRes.Users = strings.Split(role.UserNames, ",")
		}
		if role.InstanceNames != "" {
			roleRes.Instances = strings.Split(role.InstanceNames, ",")
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
	roles, err := s.GetAllRoleTip()
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
