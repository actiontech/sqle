package v1

import (
	"fmt"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/labstack/echo/v4"
)

type CreateUserGroupReqV1 struct {
	Name  string   `json:"user_group_name" form:"user_group_name" example:"test" valid:"required,name"`
	Desc  string   `json:"user_group_desc" form:"user_group_desc" example:"this is a group"`
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
			return controller.JSONNewDataExistErr(c, "user group already exist")
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

	// user group
	ug := &model.UserGroup{
		Name: req.Name,
		Desc: req.Desc,
	}

	if err := s.SaveUserGroupAndAssociations(ug, users); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return controller.JSONBaseErrorReq(c, nil)
}

type GetUserGroupsReqV1 struct {
	FilterUserGroupName string `json:"filter_user_group_name" query:"filter_user_group_name"`
	PageIndex           uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize            uint32 `json:"page_size" query:"page_size" valid:"required"`
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
	req := new(GetUserGroupsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()

	var offset uint32
	if req.PageIndex >= 1 {
		offset = req.PageSize * (req.PageIndex - 1)
	}

	data := map[string]interface{}{
		"filter_user_group_name": req.FilterUserGroupName,
		"limit":                  req.PageSize,
		"offset":                 offset,
	}

	userGroups, count, err := s.GetUserGroupsByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	resData := make([]*UserGroupListItemResV1, len(userGroups))
	for i := range userGroups {
		userGroupItem := &UserGroupListItemResV1{
			Name:       userGroups[i].Name,
			Desc:       userGroups[i].Desc,
			IsDisabled: userGroups[i].IsDisabled(),
			Users:      userGroups[i].UserNames,
		}
		resData[i] = userGroupItem
	}

	return c.JSON(http.StatusOK, &GetUserGroupsResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      resData,
		TotalNums: count,
	})
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

	err = s.RemoveMemberGroupFromAllProjectByUserGroupID(ug.ID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = s.Delete(ug)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return controller.JSONBaseErrorReq(c, nil)
}

type PatchUserGroupReqV1 struct {
	Desc       *string   `json:"user_group_desc,omitempty" form:"user_group_desc" example:"this is a group"`
	Users      *[]string `json:"user_name_list,omitempty" form:"user_name_list"`
	IsDisabled *bool     `json:"is_disabled,omitempty" form:"is_disabled"`
}

// @Summary 更新用户组
// @Description update user group
// @Id updateUserGroupV1
// @Tags user_group
// @Security ApiKeyAuth
// @Param user_group_name path string true "user_group_name"
// @Param instance body v1.PatchUserGroupReqV1 true "update user group"
// @Success 200 {object} controller.BaseRes
// @router /v1/user_groups/{user_group_name}/ [patch]
func UpdateUserGroup(c echo.Context) (err error) {

	userGroupName := c.Param("user_group_name")

	req := new(PatchUserGroupReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	s := model.GetStorage()

	// check if user group already exist
	var ug *model.UserGroup
	{
		var isExist bool
		ug, isExist, err = s.GetUserGroupByName(userGroupName)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if !isExist {
			return controller.JSONNewDataNotExistErr(c, "user_group<%v> not exist", userGroupName)
		}
	}

	// update stat
	if req.IsDisabled != nil {
		if *req.IsDisabled {
			ug.SetStat(model.Disabled)
		} else {
			ug.SetStat(model.Enabled)
		}
	}

	// update desc
	if req.Desc != nil {
		ug.Desc = *req.Desc
	}

	// users
	var users []*model.User
	{
		if req.Users != nil {
			if len(*req.Users) > 0 {
				users, err = s.GetAndCheckUserExist(*req.Users)
				if err != nil {
					return controller.JSONBaseErrorReq(c, err)
				}
			} else {
				users = make([]*model.User, 0)
			}
		}
	}

	if err := s.SaveUserGroupAndAssociations(ug, users); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return controller.JSONBaseErrorReq(c, nil)
}

type UserGroupTipsReqV1 struct {
	FilterProject string `json:"filter_project"`
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
// @Param filter_project query string false "project name"
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetUserGroupTipsResV1
// @router /v1/user_group_tips [get]
func GetUserGroupTips(c echo.Context) error {
	s := model.GetStorage()

	projectName := c.Param("filter_project")

	userGroupNames, err := s.GetUserGroupTipByProject(projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	userGroupTipsRes := make([]*UserGroupTipListItem, len(userGroupNames))
	for i := range userGroupNames {
		userGroupTipsRes[i] = &UserGroupTipListItem{
			Name: userGroupNames[i].Name,
		}
	}

	return c.JSON(http.StatusOK, &GetUserGroupTipsResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    userGroupTipsRes,
	})
}

type CreateMemberGroupReqV1 struct {
	UserGroupName string          `json:"user_group_name" valid:"required"`
	Roles         []BindRoleReqV1 `json:"roles" valid:"required"`
}

// AddMemberGroup
// @Summary 添加成员组
// @Description add member group
// @Id addMemberGroupV1
// @Tags user_group
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param project_name path string true "project name"
// @Param data body v1.CreateMemberGroupReqV1 true "add member group"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/member_groups [post]
func AddMemberGroup(c echo.Context) error {
	req := new(CreateMemberGroupReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	projectName := c.Param("project_name")
	userName := controller.GetUserName(c)

	err := CheckIsProjectManager(userName, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	// 检查用户组是否已添加过
	isMember, err := s.CheckUserGroupIsMember(req.UserGroupName, projectName)
	if err != nil {
		return err
	}
	if isMember {
		return errors.New(errors.DataExist, fmt.Errorf("user group %v is in project %v", req.UserGroupName, projectName))
	}

	role := []model.BindRole{}
	instNames := []string{}
	roleNames := []string{}
	for _, r := range req.Roles {
		role = append(role, model.BindRole{
			RoleNames:    r.RoleNames,
			InstanceName: r.InstanceName,
		})
		instNames = append(instNames, r.InstanceName)
		roleNames = append(roleNames, r.RoleNames...)
	}

	// 检查实例是否存在
	exist, err := s.CheckInstancesExist(projectName, utils.RemoveDuplicate(instNames))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("prohibit binding non-existent instances")))
	}

	// 检查角色是否存在
	exist, err = s.CheckRolesExist(utils.RemoveDuplicate(roleNames))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("prohibit binding non-existent roles")))
	}

	return controller.JSONBaseErrorReq(c, s.AddMemberGroup(req.UserGroupName, projectName, role))
}

type UpdateMemberGroupReqV1 struct {
	Roles *[]BindRoleReqV1 `json:"roles"`
}

// UpdateMemberGroup
// @Summary 修改成员组
// @Description update member group
// @Id updateMemberGroupV1
// @Tags user_group
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param project_name path string true "project name"
// @Param user_group_name path string true "user group name"
// @Param data body v1.UpdateMemberGroupReqV1 true "update member_group"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/member_groups/{user_group_name}/ [patch]
func UpdateMemberGroup(c echo.Context) error {
	req := new(UpdateMemberGroupReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	projectName := c.Param("project_name")
	groupName := c.Param("user_group_name")
	currentUser := controller.GetUserName(c)

	err := CheckIsProjectManager(currentUser, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	isMember, err := s.CheckUserGroupIsMember(groupName, projectName)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New(errors.DataNotExist, fmt.Errorf("user group %v is not in project %v", groupName, projectName))
	}

	// 更新角色
	role := []model.BindRole{}
	if req.Roles != nil {
		for _, r := range *req.Roles {
			role = append(role, model.BindRole{
				RoleNames:    r.RoleNames,
				InstanceName: r.InstanceName,
			})
		}
	}
	if len(role) > 0 {
		return controller.JSONBaseErrorReq(c, s.UpdateUserGroupRoles(groupName, projectName, role))
	}

	return controller.JSONBaseErrorReq(c, nil)
}

// DeleteMemberGroup
// @Summary 删除成员组
// @Description delete member group
// @Id deleteMemberGroupV1
// @Tags user_group
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param user_group_name path string true "user group name"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/member_groups/{user_group_name}/ [delete]
func DeleteMemberGroup(c echo.Context) error {
	projectName := c.Param("project_name")
	groupName := c.Param("user_group_name")
	currentUser := controller.GetUserName(c)

	err := CheckIsProjectManager(currentUser, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()

	return controller.JSONBaseErrorReq(c, s.RemoveMemberGroup(groupName, projectName))

}

type GetMemberGroupReqV1 struct {
	FilterInstanceName  string `json:"filter_instance_name" query:"filter_instance_name"`
	FilterUserGroupName string `json:"filter_user_group_name" query:"filter_user_group_name"`
	PageIndex           uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize            uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetMemberGroupsRespV1 struct {
	controller.BaseRes
	Data      []GetMemberGroupRespDataV1 `json:"data"`
	TotalNums uint64                     `json:"total_nums"`
}

type GetMemberGroupRespDataV1 struct {
	UserGroupName string          `json:"user_group_name"`
	Roles         []BindRoleReqV1 `json:"roles"`
}

// GetMemberGroups
// @Summary 获取成员组列表
// @Description get member groups
// @Id getMemberGroupsV1
// @Tags user_group
// @Security ApiKeyAuth
// @Param filter_user_group_name query string false "filter user group name"
// @Param filter_instance_name query string false "filter instance name"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Param project_name path string true "project name"
// @Success 200 {object} v1.GetMemberGroupsRespV1
// @router /v1/projects/{project_name}/member_groups [get]
func GetMemberGroups(c echo.Context) error {
	req := new(GetMemberGroupReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	projectName := c.Param("project_name")
	currentUser := controller.GetUserName(c)

	err := CheckIsProjectMember(currentUser, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// 获取成员组列表
	limit, offset := controller.GetLimitAndOffset(req.PageIndex, req.PageSize)
	filter := model.GetMemberGroupFilter{
		Limit:             &limit,
		Offset:            &offset,
		FilterProjectName: &projectName,
	}
	if req.FilterInstanceName != "" {
		filter.FilterInstanceName = &req.FilterInstanceName
	}
	if req.FilterUserGroupName != "" {
		filter.FilterUserGroupName = &req.FilterUserGroupName
	}

	s := model.GetStorage()
	groups, err := s.GetMemberGroups(filter)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	total, err := s.GetMemberGroupCount(filter)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// 获取角色信息
	groupNames := []string{}
	for _, group := range groups {
		groupNames = append(groupNames, group.Name)
	}

	bindRole, err := s.GetBindRolesByMemberGroupNames(groupNames, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// 生成响应
	data := []GetMemberGroupRespDataV1{}
	for _, group := range groups {
		data = append(data, GetMemberGroupRespDataV1{
			UserGroupName: group.Name,
			Roles:         convertBindRoleToBindRoleReqV1(bindRole[group.Name]),
		})
	}

	return c.JSON(http.StatusOK, GetMemberGroupsRespV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      data,
		TotalNums: total,
	})
}

type GetMemberGroupRespV1 struct {
	controller.BaseRes
	Data GetMemberGroupRespDataV1 `json:"data"`
}

// GetMemberGroup
// @Summary 获取成员组信息
// @Description get member group
// @Id getMemberGroupV1
// @Tags user_group
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param user_group_name path string true "user group name"
// @Success 200 {object} v1.GetMemberGroupRespV1
// @router /v1/projects/{project_name}/member_groups/{user_group_name}/ [get]
func GetMemberGroup(c echo.Context) error {
	projectName := c.Param("project_name")
	groupName := c.Param("user_group_name")
	currentUser := controller.GetUserName(c)

	err := CheckIsProjectMember(currentUser, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	userGroup, err := s.GetMemberGroupByGroupName(projectName, groupName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	bindRole, err := s.GetBindRolesByMemberGroupNames([]string{userGroup.Name}, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, GetMemberGroupRespV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: GetMemberGroupRespDataV1{
			UserGroupName: userGroup.Name,
			Roles:         convertBindRoleToBindRoleReqV1(bindRole[userGroup.Name]),
		},
	})
}
