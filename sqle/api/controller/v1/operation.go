package v1

import (
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
)

type GetOperationsResV1 struct {
	controller.BaseRes
	Data []*OperationListItemResV1 `json:"data"`
}

type OperationListItemResV1 struct {
	Code uint   `json:"op_code"`
	Desc string `json:"op_desc"`
}

// @Summary 获取权限动作列表
// @Description get permission operations
// @Id GetOperationsV1
// @Tags operation
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetOperationsResV1
// @Router /v1/operations [get]
func GetOperations(c echo.Context) error {

	opCodes := model.GetOperationCodeList()

	respData := make([]*OperationListItemResV1, len(opCodes))

	for i := range opCodes {
		respData[i] = &OperationListItemResV1{
			Code: uint(opCodes[i]),
			Desc: model.GetOperationCodeDesc(opCodes[i]),
		}
	}

	return c.JSON(http.StatusOK, &GetOperationsResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    respData,
	})
}
