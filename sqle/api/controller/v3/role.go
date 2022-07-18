package v3

import (
	"fmt"
	"net/http"

	"github.com/actiontech/sqle/sqle/utils"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
)

type GetRolesReqV3 struct {
	FilterRoleName     string `json:"filter_role_name" query:"filter_role_name"`
	FilterUserName     string `json:"filter_user_name" query:"filter_user_name"`
	FilterInstanceName string `json:"filter_instance_name" query:"filter_instance_name"`
	PageIndex          uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize           uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type Operation struct {
	Code uint   `json:"op_code"`
	Desc string `json:"op_desc"`
}

type Instance struct {
	Name   string `json:"name"`
	DbType string `json:"db_type"`
}

type RoleResV3 struct {
	Name       string       `json:"role_name"`
	Desc       string       `json:"role_desc"`
	Users      []string     `json:"user_name_list,omitempty"`
	Instances  []Instance   `json:"instance_list,omitempty"`
	Operations []*Operation `json:"operation_list,omitempty"`
	UserGroups []string     `json:"user_group_name_list,omitempty" form:"user_group_name_list"`
	IsDisabled bool         `json:"is_disabled,omitempty"`
}

type GetRolesResV3 struct {
	controller.BaseRes
	Data      []*RoleResV3 `json:"data"`
	TotalNums uint64       `json:"total_nums"`
}

// GetRolesV3
// @Summary 获取角色列表
// @Description get role list
// @Id getRoleListV3
// @Tags role
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param filter_role_name query string false "filter role name"
// @Param filter_user_name query string false "filter user name"
// @Param filter_instance_name query string false "filter instance name"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v3.GetRolesResV3
// @router /v3/roles [get]
func GetRolesV3(c echo.Context) error {
	req := new(GetRolesReqV3)
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
			"filter_role_name":     req.FilterRoleName,
			"filter_user_name":     req.FilterUserName,
			"filter_instance_name": req.FilterInstanceName,
			"limit":                limit,
			"offset":               offset,
		}
	}

	roles, count, err := s.GetRolesByReq(queryCondition)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	roleToInstances, err := buildInstanceList(roles, s)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	roleRes := make([]*RoleResV3, len(roles))
	for i := range roles {
		ops := make([]*Operation, len(roles[i].OperationsCodes))
		opCodes := roles[i].OperationsCodes.ForceConvertIntSlice()
		for i := range opCodes {
			ops[i] = &Operation{
				Code: opCodes[i],
				Desc: model.GetOperationCodeDesc(opCodes[i]),
			}
		}

		roleRes[i] = &RoleResV3{
			Name:       roles[i].Name,
			Desc:       roles[i].Desc,
			Instances:  roleToInstances[roles[i].Name],
			UserGroups: roles[i].UserGroupNames,
			Users:      roles[i].UserNames,
			IsDisabled: roles[i].IsDisabled(),
			Operations: ops,
		}

	}

	return c.JSON(http.StatusOK, &GetRolesResV3{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      roleRes,
		TotalNums: count,
	})
}

func buildInstanceList(roles []*model.RoleDetail, s *model.Storage) (roleToInstances map[string][]Instance, err error) {
	var allInstanceNames []string
	for _, role := range roles {
		allInstanceNames = append(allInstanceNames, role.InstanceNames...)
	}
	allInstanceNames = utils.RemoveDuplicate(allInstanceNames)
	instancesFromDB, err := s.GetInstancesByNames(allInstanceNames)
	if err != nil {
		return nil, fmt.Errorf("get instances by name failed: %v", err)
	}

	roleToInstances = make(map[string][]Instance, len(roles))
	for _, role := range roles {
		instances := make([]Instance, len(role.InstanceNames))
		for i, instName := range role.InstanceNames {
			instances[i] = Instance{
				Name:   instName,
				DbType: getDbTypeFromInstancesByInstanceName(instName, instancesFromDB),
			}
		}
		roleToInstances[role.Name] = instances
	}

	return roleToInstances, nil
}

func getDbTypeFromInstancesByInstanceName(name string, instances []*model.Instance) string {
	for _, inst := range instances {
		if inst.Name == name {
			return inst.DbType
		}
	}
	return ""
}
