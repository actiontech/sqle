//go:build release
// +build release

package v1

import (
	e "errors"
	"fmt"
	"mime"
	"net/http"
	"strconv"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/license"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
)

var ErrNoLicenseRequired = errors.New(errors.ErrAccessDeniedError, e.New("sqle-ce no license required"))

const (
	HardwareInfoFileName = "collected.infos"
	LicenseFileParamKey  = "license_file"
)

func getLicense(c echo.Context) error {
	s := model.GetStorage()
	l, exist, err := s.GetLicense()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return c.JSON(http.StatusOK, GetLicenseResV1{
			BaseRes: controller.NewBaseReq(nil),
		})
	}

	permission, collectedInfosContent, err := license.DecodeLicense(l.Content)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, license.ErrInvalidLicense))
	}

	items := generateLicenseItems(permission, collectedInfosContent)

	return c.JSON(http.StatusOK, GetLicenseResV1{
		BaseRes: controller.NewBaseReq(nil),
		Content: l.Content,
		License: items,
	})

}

func getSQLELicenseInfo(c echo.Context) error {
	info, err := license.CollectHardwareInfo()
	if err != nil {
		return controller.JSONBaseErrorReq(c, license.ErrCollectLicenseInfo)
	}

	c.Response().Header().Set(echo.HeaderContentDisposition,
		mime.FormatMediaType("attachment", map[string]string{"filename": HardwareInfoFileName}))

	return c.Blob(http.StatusOK, echo.MIMETextPlain, []byte(info))
}

func setLicense(c echo.Context) error {
	file, exist, err := controller.ReadFileContent(c, LicenseFileParamKey)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, license.ErrLicenseEmpty))
	}

	{ // check license info
		_, collectedInfosContent, err := license.DecodeLicense(file)
		if err != nil {
			return controller.JSONBaseErrorReq(c, license.ErrInvalidLicense)
		}
		collected, err := license.CollectHardwareInfo()
		if err != nil {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataParseFail, license.ErrCollectLicenseInfo))
		}

		if collected != collectedInfosContent {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, license.ErrInvalidLicense))
		}
	}

	s := model.GetStorage()
	err = s.Delete(&model.License{})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = s.Save(&model.License{Content: file, WorkDurationHour: 0})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = license.UpdateLicense(file)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

func checkLicense(c echo.Context) error {
	file, exist, err := controller.ReadFileContent(c, LicenseFileParamKey)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, license.ErrLicenseEmpty))
	}

	permission, collectedInfosContent, err := license.DecodeLicense(file)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, license.ErrInvalidLicense))
	}
	collected, err := license.CollectHardwareInfo()
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataParseFail, license.ErrCollectLicenseInfo))
	}

	if collected != collectedInfosContent {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, license.ErrInvalidLicense))
	}

	items := generateLicenseItems(permission, collectedInfosContent)

	return c.JSON(http.StatusOK, GetLicenseResV1{
		BaseRes: controller.NewBaseReq(nil),
		Content: file,
		License: items,
	})

}

func generateLicenseItems(permission *license.LicensePermission, collectedInfosContent string) []LicenseItem {
	items := []LicenseItem{}

	for n, i := range permission.NumberOfInstanceOfEachType {
		items = append(items, LicenseItem{
			Description: fmt.Sprintf("[%v]类型实例数", n),
			Name:        n,
			Limit:       strconv.Itoa(i.Count),
		})
	}

	items = append(items, []LicenseItem{
		{
			Description: "用户数",
			Name:        "user",
			Limit:       strconv.Itoa(permission.UserCount),
		}, {
			Description: "机器信息",
			Name:        "info",
			Limit:       collectedInfosContent,
		}, {
			Description: "授权运行时长(天)",
			Name:        "work duration day",
			Limit:       strconv.Itoa(permission.WorkDurationDay),
		}, {
			Description: "SQLE版本",
			Name:        "version",
			Limit:       permission.Version,
		},
	}...)

	return items
}
