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

type CreateRoleReqV1 struct {
	Name           string `json:"role_name" form:"role_name" valid:"required,name"`
	Desc           string `json:"role_desc" form:"role_desc"`
	OperationCodes []uint `json:"operation_code_list" form:"operation_code_list"`
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
func CreateRole(c echo.Context) (err error) {

	req := new(CreateRoleReqV1)
	{
		if err := controller.BindAndValidateReq(c, req); err != nil {
			return err
		}
	}

	s := model.GetStorage()

	// check if role name already exists
	{
		_, exist, err := s.GetRoleByName(req.Name)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if exist {
			return controller.JSONNewDataExistErr(c, "role<%s> is exist", req.Name)
		}
	}

	// check operation codes
	{
		if len(req.OperationCodes) > 0 {
			if err := model.CheckIfOperationCodeValid(req.OperationCodes); err != nil {
				return controller.JSONBaseErrorReq(c, err)
			}
		}
	}

	newRole := &model.Role{
		Name: req.Name,
		Desc: req.Desc,
	}

	return controller.JSONBaseErrorReq(c,
		s.SaveRoleAndAssociations(newRole, req.OperationCodes),
	)
}

type GetRolesReqV1 struct {
	FilterRoleName string `json:"filter_role_name" query:"filter_role_name"`
	PageIndex      uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize       uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetRolesResV1 struct {
	controller.BaseRes
	Data      []*RoleResV1 `json:"data"`
	TotalNums uint64       `json:"total_nums"`
}

type Operation struct {
	Code uint   `json:"op_code"`
	Desc string `json:"op_desc"`
}

type RoleResV1 struct {
	Name       string       `json:"role_name"`
	Desc       string       `json:"role_desc"`
	Operations []*Operation `json:"operation_list,omitempty"`
	IsDisabled bool         `json:"is_disabled,omitempty"`
}

// @Summary 获取角色列表
// @Description get role list
// @Id getRoleListV1
// @Tags role
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param filter_role_name query string false "filter role name"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v1.GetRolesResV1
// @router /v1/roles [get]
func GetRoles(c echo.Context) error {
	req := new(GetRolesReqV1)
	{
		if err := controller.BindAndValidateReq(c, req); err != nil {
			return err
		}
	}

	s := model.GetStorage()

	var queryCondition map[string]interface{}
	{
		limit, offset := controller.GetLimitAndOffset(req.PageIndex, req.PageSize)
		queryCondition = map[string]interface{}{
			"filter_role_name": req.FilterRoleName,
			"limit":            limit,
			"offset":           offset,
		}
	}

	roles, count, err := s.GetRolesByReq(queryCondition)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	roleRes := make([]*RoleResV1, len(roles))
	for i := range roles {
		ops := make([]*Operation, len(roles[i].OperationsCodes))
		opCodes := roles[i].OperationsCodes.ForceConvertIntSlice()
		for i := range opCodes {
			ops[i] = &Operation{
				Code: opCodes[i],
				Desc: model.GetOperationCodeDesc(opCodes[i]),
			}
		}
		roleRes[i] = &RoleResV1{
			Name:       roles[i].Name,
			Desc:       roles[i].Desc,
			IsDisabled: roles[i].IsDisabled(),
			Operations: ops,
		}

	}

	return c.JSON(http.StatusOK, &GetRolesResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      roleRes,
		TotalNums: count,
	})
}

type UpdateRoleReqV1 struct {
	Desc           *string `json:"role_desc" form:"role_desc"`
	OperationCodes *[]uint `json:"operation_code_list,omitempty" form:"operation_code_list"`
	IsDisabled     *bool   `json:"is_disabled,omitempty"`
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
func UpdateRole(c echo.Context) (err error) {

	req := new(UpdateRoleReqV1)
	{
		if err := controller.BindAndValidateReq(c, req); err != nil {
			return err
		}
	}

	s := model.GetStorage()
	roleName := c.Param("role_name")

	// check if role name exists
	var role *model.Role
	{
		var isExist bool
		role, isExist, err = s.GetRoleByName(roleName)
		if err != nil {
			if err != nil {
				return controller.JSONBaseErrorReq(c, err)
			}
		}
		if !isExist {
			return controller.JSONNewDataNotExistErr(c,
				`role is not exist`)
		}
	}

	// update stat
	{
		if req.IsDisabled != nil {
			if *req.IsDisabled {
				role.Stat = model.Disabled
			} else {
				role.Stat = model.Enabled
			}

		}
	}

	// update desc
	if req.Desc != nil {
		role.Desc = *req.Desc
	}

	// check operation codes
	var opCodes []uint
	{
		if req.OperationCodes != nil {
			if len(*req.OperationCodes) > 0 {
				if err := model.CheckIfOperationCodeValid(*req.OperationCodes); err != nil {
					return controller.JSONBaseErrorReq(c, err)
				}
				opCodes = *req.OperationCodes
			} else {
				opCodes = make([]uint, 0)
			}
		}
	}

	return controller.JSONBaseErrorReq(c,
		s.SaveRoleAndAssociations(role, opCodes),
	)

}
