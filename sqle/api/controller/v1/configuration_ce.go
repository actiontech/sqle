//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"
	"fmt"
	"net/http"

	"github.com/actiontech/sqle/sqle/errors"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/config"
	"github.com/labstack/echo/v4"
)

var (
	errCommunityEditionNotSupportCostumeLogo             = errors.New(errors.EnterpriseEditionFeatures, e.New("costume logo is enterprise version feature"))
	errCommunityEditionNotSupportUpdatePersonaliseConfig = errors.New(errors.EnterpriseEditionFeatures, e.New("update personalise config is enterprise version feature"))
)

const (
	Title = "SQLE"
)

func uploadLogo(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportCostumeLogo)
}

func getLogo(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportCostumeLogo)
}

func updatePersonaliseConfig(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportUpdatePersonaliseConfig)
}

func getSQLEInfo(c echo.Context) error {
	fileInfo, err := getLogoFileInfo()
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, e.New("no logo file")))
	}

	modifyTime := fileInfo.ModTime().Unix()
	logoUrl := fmt.Sprintf("%s/%s?timestamp=%d", LogoUrlBase, fileInfo.Name(), modifyTime)

	return c.JSON(http.StatusOK, &GetSQLEInfoResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: GetSQLEInfoResDataV1{
			Version: config.Version,
			LogoUrl: logoUrl,
			Title:   Title,
		},
	})
}
