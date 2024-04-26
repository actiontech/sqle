package v1

import (
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/server"

	"github.com/labstack/echo/v4"
)

type GetModuleStatusReqV1 struct {
	DbType     string `json:"db_type" query:"db_type"`
	ModuleName string `json:"module_name" query:"module_name"`
}

type GetModuleStatusResV1 struct {
	controller.BaseRes
	Data ModuleStatusRes `json:"data"`
}

type ModuleStatusRes struct {
	IsSupported bool `json:"is_supported"`
}

// @Summary 查询系统功能支持情况信息
// @Description get module status for modulealities in the system
// @Id getSystemModuleStatus
// @Tags system
// @Security ApiKeyAuth
// @Param db_type query string false "db type" Enums(MySQL,Oracle,TiDB,OceanBase For MySQL,PostgreSQL,DB2,SQL Server)
// @Param module_name query string false "module name" Enums(execute_sql_file_mode,sql_optimization)
// @Success 200 {object} v1.GetModuleStatusResV1
// @router /v1/system/module_status [get]
func GetSystemModuleStatus(c echo.Context) error {
	req := new(GetModuleStatusReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	checker, err := server.NewModuleStatusChecker(req.DbType, req.ModuleName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetModuleStatusResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: ModuleStatusRes{
			IsSupported: checker.CheckIsSupport(),
		},
	})
}
