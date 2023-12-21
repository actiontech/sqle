package v1

import (
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/dms"

	"github.com/labstack/echo/v4"
)

type UserTipsReqV1 struct {
	FilterProject string `json:"filter_project" query:"filter_project" valid:"required"`
}

type UserTipResV1 struct {
	UserID string `json:"user_id"`
	Name   string `json:"user_name"`
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
// @Param filter_project query string true "project name"
// @Success 200 {object} v1.GetUserTipsResV1
// @router /v1/user_tips [get]
func GetUserTips(c echo.Context) error {
	req := new(UserTipsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), req.FilterProject)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	users, err := dms.ListProjectUserTips(c.Request().Context(), projectUid)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	userTipsRes := make([]UserTipResV1, 0, len(users))
	for _, user := range users {
		userTipRes := UserTipResV1{
			Name:   user.Name,
			UserID: user.GetIDStr(),
		}
		userTipsRes = append(userTipsRes, userTipRes)
	}
	return c.JSON(http.StatusOK, &GetUserTipsResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    userTipsRes,
	})
}
