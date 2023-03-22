//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/actiontech/sqle/sqle/errors"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/config"
	"github.com/labstack/echo/v4"
)

var errCommunityEditionNotSupportCostumeLogo = e.New("costume logo is enterprise version feature")

const (
	// LogoUrlBase sqle static 服务接口的url前缀
	LogoUrlBase = "/static/media"
	Title       = "SQLE"

	// LogoDir sqle logo 的本地目录
	LogoDir = "./ui/static/media"
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

func getLogoFileInfo() (fs.FileInfo, error) {
	fileInfos, err := ioutil.ReadDir(LogoDir)
	if err != nil {
		return nil, e.New("read logo dir failed")
	}

	var hasLogoFile bool
	var logoFileInfo fs.FileInfo
	for _, fileInfo := range fileInfos {
		if strings.HasPrefix(fileInfo.Name(), "logo.") {
			hasLogoFile = true
			logoFileInfo = fileInfo
			break
		}
	}
	if !hasLogoFile {
		return nil, e.New("no logo file")
	}

	return logoFileInfo, nil
}
