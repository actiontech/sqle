package v2

import (
	"fmt"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/driver"

	"github.com/labstack/echo/v4"
)

type GetDriversRes struct {
	controller.BaseRes
	Data []*DriverMeta `json:"data"`
}

type DriverMeta struct {
	Name        string `json:"driver_name"`
	DefaultPort uint   `json:"default_port"`
	LogoUrl     string `json:"logo_url"`
}

// GetDrivers get support Driver list.
// @Summary 获取当前 server 支持的审核类型
// @Description get drivers
// @Id getDriversV2
// @Tags configuration
// @Security ApiKeyAuth
// @Success 200 {object} v2.GetDriversRes
// @router /v2/configurations/drivers [get]
func GetDrivers(c echo.Context) error {

	metas := driver.GetPluginManager().AllDriverMetas()
	data := make([]*DriverMeta, len(metas))

	for i := range metas {
		meta := metas[i]
		data[i] = &DriverMeta{
			Name:        meta.PluginName,
			DefaultPort: uint(meta.DatabaseDefaultPort),
			LogoUrl:     fmt.Sprintf("/v1/static/instance_logo?instance_type=%s", meta.PluginName),
		}
	}

	return c.JSON(http.StatusOK, &GetDriversRes{
		BaseRes: controller.NewBaseReq(nil),
		Data:    data,
	})
}
