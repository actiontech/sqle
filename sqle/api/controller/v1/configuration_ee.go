//go:build enterprise
// +build enterprise

package v1

import (
	e "errors"
	"fmt"
	"net/http"
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
)

const (
	Title = "SQLE企业版"

	// LogoUrl 用户配置的 logo url 接口
	LogoUrl = "/v1/static/logo"

	// LogoFileKey logo 文件key
	LogoFileKey = "logo"

	// MaxByteSizeOfLogo logo 最大字节数, 100KB
	MaxByteSizeOfLogo = 1024 * 100
)

var (
	logoUrl = func(time time.Time) string {
		return fmt.Sprintf("%s?timestamp=%d", LogoUrl, time.Unix())
	}

	// logoCache logo 缓存
	logoCache = map[string] /*logo update unix time*/ []byte{} /*logo*/
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
	logo, exist, err := controller.ReadFileContent(c, LogoFileKey)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf("failed to read logo file: %w", err)))
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, e.New("logo file not exist")))
	}

	if isLogoMoreThanMaxSize([]byte(logo)) {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf("logo file size is too large, large than max byte %d", MaxByteSizeOfLogo)))
	}

	s := model.GetStorage()
	logoConfig, _, err := s.GetLogoConfigWithoutLogoImage()
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, fmt.Errorf("failed to get logo config: %w", err)))
	}

	logoConfig.Logo = []byte(logo)

	if err := s.Save(&logoConfig); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, UploadLogoResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: UploadLogoResDataV1{
			LogoUrl: logoUrl(logoConfig.UpdatedAt),
		},
	})
}

func isLogoMoreThanMaxSize(logo []byte) bool {
	if len(logo) > MaxByteSizeOfLogo {
		return true
	}
	return false
}

func getLogo(c echo.Context) error {
	req := new(GetLogoReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if logo, ok := logoCache[req.Timestamp]; ok {
		return c.Blob(http.StatusOK, "image/png", logo)
	}

	s := model.GetStorage()
	logoConfig, exist, err := s.GetLogoConfig()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, e.New("logoConfig not exist")))
	}

	logoCache[req.Timestamp] = logoConfig.Logo

	return c.Blob(http.StatusOK, "image/png", logoConfig.Logo)
}

func getSQLEInfo(c echo.Context) error {
	s := model.GetStorage()
	personaliseConfig, _, err := s.GetPersonaliseConfig()
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, fmt.Errorf("failed to get personalise config: %w", err)))
	}

	logo, _, err := s.GetLogoConfigWithoutLogoImage()
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, fmt.Errorf("failed to get logo config: %w", err)))
	}

	baseInfo, err := getDefaultBaseInfo()
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("failed to get default base info: %w", err)))
	}

	if !logo.UpdatedAt.Equal(time.Time{}) {
		baseInfo.LogoUrl = logoUrl(logo.UpdatedAt)
	}

	if personaliseConfig.Title != "" {
		baseInfo.Title = personaliseConfig.Title
	}

	return c.JSON(http.StatusOK, &GetSQLEInfoResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    baseInfo,
	})
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

func updateFeishuAuditConfigurationV1(c echo.Context) error {
	req := new(UpdateFeishuConfigurationReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()
	feishuCfg, _, err := s.GetImConfigByType(model.ImTypeFeishuApproval)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	{ // disable
		if req.IsFeishuNotificationEnabled != nil && !(*req.IsFeishuNotificationEnabled) {
			feishuCfg.IsEnable = false
			return controller.JSONBaseErrorReq(c, s.Save(feishuCfg))
		}
	}

	if req.AppID != nil {
		feishuCfg.AppKey = *req.AppID
	}
	if req.AppSecret != nil {
		feishuCfg.AppSecret = *req.AppSecret
	}
	if req.IsFeishuNotificationEnabled != nil {
		feishuCfg.IsEnable = *req.IsFeishuNotificationEnabled
	}

	feishuCfg.Type = model.ImTypeFeishuApproval

	if err := s.Save(feishuCfg); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return controller.JSONBaseErrorReq(c, nil)
}

func getFeishuAuditConfigurationV1(c echo.Context) error {
	s := model.GetStorage()
	feishuCfg, exist, err := s.GetImConfigByType(model.ImTypeFeishuApproval)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return c.JSON(http.StatusOK, &GetFeishuConfigurationResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data:    FeishuConfigurationV1{},
		})
	}

	return c.JSON(http.StatusOK, &GetFeishuConfigurationResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: FeishuConfigurationV1{
			AppID:                       feishuCfg.AppKey,
			IsFeishuNotificationEnabled: feishuCfg.IsEnable,
		},
	})
}
