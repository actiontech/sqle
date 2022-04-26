//go:build enterprise
// +build enterprise

package v1

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

var errOauth2NotEnable = errors.New(errors.ErrAccessDeniedError, fmt.Errorf("oauth2 not enable"))
var oauthState = "sqle-ee-action"

func oauth2Link(c echo.Context) error {
	entry := log.NewEntry()

	s := model.GetStorage()
	oauth2C, _, err := s.GetOauth2Configuration()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !oauth2C.EnableOauth2 {
		return controller.JSONBaseErrorReq(c, errOauth2NotEnable)
	}

	uri := generateOauth2Config(oauth2C).AuthCodeURL(oauthState)
	_, err = url.ParseRequestURI(uri)
	if err != nil {
		entry.Errorf("parse oauth2 link failed: %v", err)
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataParseFail, fmt.Errorf("parse oauth2 link failed, please check oauth2 config.")))
	}

	return c.Redirect(http.StatusFound, uri)
}

func generateOauth2Config(conf *model.Oauth2Configuration) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientKey,
		RedirectURL:  fmt.Sprintf("%v/v1/oauth2/callback", conf.ClientHost),
		Scopes:       conf.GetScopes(),
		Endpoint: oauth2.Endpoint{
			AuthURL:  conf.ServerAuthUrl,
			TokenURL: conf.ServerTokenUrl,
		},
	}
}

func oauth2Callback(c echo.Context) error {
	return nil
}

func bindOauth2User(c echo.Context) error {
	req := new(BindOauth2UserReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	if req.Oauth2UserID == "" {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf("oauth2 user id can not empty")))
	}

	s := model.GetStorage()
	user, exist, err := s.GetUserByThirdPartyUserID(req.Oauth2UserID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		user = &model.User{
			Name:                   req.UserName,
			Password:               req.Pwd,
			ThirdPartyUserID:       req.Oauth2UserID,
			UserAuthenticationType: model.UserAuthenticationTypeOAUTH2,
		}
		err = s.Save(user)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	// check user login type
	if user.UserAuthenticationType != model.UserAuthenticationTypeOAUTH2 &&
		user.UserAuthenticationType != model.UserAuthenticationTypeSQLE &&
		user.UserAuthenticationType != "" {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("the user has bound other login methods")))
	}

	// modify user login type
	if user.UserAuthenticationType != model.UserAuthenticationTypeOAUTH2 {
		user.ThirdPartyUserID = req.Oauth2UserID
		user.UserAuthenticationType = model.UserAuthenticationTypeOAUTH2
		err = s.Update(user)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	t, err := generateToken(req.UserName)
	return c.JSON(http.StatusOK, BindOauth2UserResV1{
		BaseRes: controller.NewBaseReq(err),
		Data: BindOauth2UserResDataV1{
			Token: t,
		},
	})
}
