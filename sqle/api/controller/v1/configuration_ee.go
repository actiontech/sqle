//go:build enterprise
// +build enterprise

package v1

import (
	e "errors"
	"fmt"
	"net/http"
	"time"

	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/errors"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
)

const (
	Title = "SQLE企业版"

	// LogoUrl 用户配置的 logo url 接口
	LogoUrl = "/v1/static/logo"
)

func updatePersonaliseConfig(c echo.Context) error {
	req := new(PersonaliseReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	personaliseConfig, _, err := s.GetPersonaliseConfig()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if req.Title != nil {
		personaliseConfig.Title = *req.Title
	}

	return controller.JSONBaseErrorReq(c, s.Save(&personaliseConfig))
}

func uploadLogo(c echo.Context) error {
	return nil
}

func getLogo(c echo.Context) error {
	return nil
}

func getSQLEInfo(c echo.Context) error {
	baseInfo, err := getDefaultBaseInfo()
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("failed to get default base info: %w", err)))
	}

	s := model.GetStorage()
	personaliseConfig, exist, err := s.GetPersonaliseConfig()
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, fmt.Errorf("failed to get personalise config: %w", err)))
	}
	if !exist {
		return c.JSON(http.StatusOK, &GetSQLEInfoResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data:    baseInfo,
		})
	}

	if personaliseConfig.Logo != nil && personaliseConfig.LogoUpdateTime != nil {
		baseInfo.LogoUrl = getLogoUrl(personaliseConfig.LogoUpdateTime)
	}

	if personaliseConfig.Title != "" {
		baseInfo.Title = personaliseConfig.Title
	}

	return c.JSON(http.StatusOK, &GetSQLEInfoResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    baseInfo,
	})
}

func getLogoUrl(updateTime *time.Time) string {
	return fmt.Sprintf("%s?timestamp=%d", LogoUrl, updateTime.Unix())
}

func getDefaultBaseInfo() (GetSQLEInfoResDataV1, error) {
	logoUrl, err := GetDefaultLogoUrl()
	if err != nil {
		return GetSQLEInfoResDataV1{}, e.New("failed to get default logo url")
	}

	return GetSQLEInfoResDataV1{
		Version: config.Version,
		LogoUrl: logoUrl,
		Title:   Title,
	}, nil
}
