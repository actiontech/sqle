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
