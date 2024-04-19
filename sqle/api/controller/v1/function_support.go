package v1

import (
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/server"

	"github.com/labstack/echo/v4"
)

type GetFunctionSupportReqV1 struct {
	DbType       string `json:"db_type" query:"db_type"`
	FunctionName string `json:"function_name" query:"function_name"`
}

type GetFunctionSupportResV1 struct {
	controller.BaseRes
	Support bool `json:"support"`
}

// @Summary 查询系统功能支持情况信息
// @Description get support for functionalities in the system
// @Id getFunctionSupport
// @Tags system
// @Security ApiKeyAuth
// @Param db_type query string false "db type" Enums(MySQL,Oracle,TiDB,OceanBase For MySQL,PostgreSQL,DB2,SQL Server)
// @Param function_name query string false "function name" Enums(execute_sql_file_mode)
// @Success 200 {object} v1.GetFunctionSupportResV1
// @router /v1/function_support [get]
func GetFunctionSupport(c echo.Context) error {
	req := new(GetFunctionSupportReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	checker, err := server.NewFunctionSupportChecker(req.DbType, req.FunctionName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetFunctionSupportResV1{
		BaseRes: controller.NewBaseReq(nil),
		Support: checker.CheckIsSupport(),
	})
}
