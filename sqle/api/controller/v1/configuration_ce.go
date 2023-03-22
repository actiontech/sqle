//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"
	"fmt"
	"net/http"
	"os"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/config"
	"github.com/labstack/echo/v4"
)

var errCommunityEditionNotSupportCostumeLogo = e.New("costume logo is enterprise version feature")

const (
	// LogoUrl sqle static 服务接口的url
	LogoUrl = "/static/media/logo.410ecb70.png"
	Title   = "SQLE"

	// LogoPath sqle logo 的本地路径
	LogoPath = "./ui/static/media/logo.410ecb70.png"
)

func uploadLogo(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportCostumeLogo)
}

func getLogo(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportCostumeLogo)
}

func updatePersonaliseConfig(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportCostumeLogo)
}

func getSQLEInfo(c echo.Context) error {
	fileInfo, err := os.Stat(LogoPath)
	if err != nil {
		return controller.JSONBaseErrorReq(c, e.New("logo file not found"))
	}

	modifyTime := fileInfo.ModTime().Format("2006-01-02 15:04:05")
	logoUrl := fmt.Sprintf("%s?timestamp=%s", LogoUrl, modifyTime)

	return c.JSON(http.StatusOK, &GetSQLEInfoResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: GetSQLEInfoResDataV1{
			Version: config.Version,
			LogoUrl: logoUrl,
			Title:   Title,
		},
	})
}
