//go:build enterprise
// +build enterprise

package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

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

type callbackRedirectData struct {
	UserExist   bool
	SqleToken   string
	Oauth2Token string
	Error       string
}

func (c callbackRedirectData) generateQuery(uri string) string {
	params := url.Values{}
	params.Set("user_exist", strconv.FormatBool(c.UserExist))
	if c.SqleToken != "" {
		params.Set("sqle_token", c.SqleToken)
	}
	if c.Oauth2Token != "" {
		params.Set("oauth2_token", c.Oauth2Token)
	}
	if c.Error != "" {
		params.Set("error", c.Error)
	}
	return fmt.Sprintf("%v/user/bind?%v", uri, params.Encode())
}

func oauth2Callback(c echo.Context) error {
	s := model.GetStorage()
	oauth2C, _, err := s.GetOauth2Configuration()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// TODO sqle https should also support
	uri := oauth2C.ClientHost
	data := callbackRedirectData{}

	// check callback request
	state := c.QueryParam("state")
	if state != oauthState {
		data.Error = fmt.Sprintf("invalid state: %v", state)
		return c.Redirect(http.StatusFound, data.generateQuery(uri))
	}
	code := c.QueryParam("code")
	if code == "" {
		data.Error = "code is nil"
		return c.Redirect(http.StatusFound, data.generateQuery(uri))
	}

	// get oauth2 token
	oauth2Token, err := generateOauth2Config(oauth2C).Exchange(context.Background(), code)
	if err != nil {
		data.Error = err.Error()
		return c.Redirect(http.StatusFound, data.generateQuery(uri))
	}
	data.Oauth2Token = oauth2Token.AccessToken

	//get user is exist
	userID, err := getOauth2UserID(oauth2C, oauth2Token.AccessToken)
	if err != nil {
		data.Error = err.Error()
		return c.Redirect(http.StatusFound, data.generateQuery(uri))
	}
	user, exist, err := s.GetUserByThirdPartyUserID(userID)
	if err != nil {
		data.Error = err.Error()
		return c.Redirect(http.StatusFound, data.generateQuery(uri))
	}
	data.UserExist = exist

	// the user has successfully logged in at the third party, and the token can be returned directly
	if exist {
		t, err := generateToken(user.Name)
		if err != nil {
			data.Error = err.Error()
			return c.Redirect(http.StatusFound, data.generateQuery(uri))
		}
		data.SqleToken = t
	}

	return c.Redirect(http.StatusFound, data.generateQuery(uri))
}

func getOauth2UserID(conf *model.Oauth2Configuration, token string) (userID string, err error) {
	uri := fmt.Sprintf("%v?%v=%v", conf.ServerUserIdUrl, conf.AccessTokenTag, token)
	resp, err := (&http.Client{}).Get(uri)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to get third-party user ID, unable to parse response")
	}
	data := map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", fmt.Errorf("failed to get third-party user ID, unrecognized response format")
	}
	user, ok := data[conf.UserIdTag]
	if !ok {
		return "", fmt.Errorf("not found third-party user ID")
	}
	return fmt.Sprintf("%v", user), nil
}

// prevent users from malicious password attempts
var errBeenBoundOrThePasswordIsWrong = errors.New(errors.DataExist, fmt.Errorf("the platform user has been bound or the password is wrong"))

func bindOauth2User(c echo.Context) error {
	req := new(BindOauth2UserReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	if req.Oauth2Token == "" {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf("oauth2 token can not empty")))
	}

	s := model.GetStorage()
	oauth2C, _, err := s.GetOauth2Configuration()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	oauth2UserID, err := getOauth2UserID(oauth2C, req.Oauth2Token)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, err))
	}

	// check third-party users have bound sqle user
	_, exist, err := s.GetUserByThirdPartyUserID(oauth2UserID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if exist {
		return controller.JSONBaseErrorReq(c, errBeenBoundOrThePasswordIsWrong)
	}

	user, exist, err := s.GetUserByName(req.UserName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// create user if not exist
	if !exist {
		user = &model.User{
			Name:                   req.UserName,
			Password:               req.Pwd,
			ThirdPartyUserID:       oauth2UserID,
			UserAuthenticationType: model.UserAuthenticationTypeOAUTH2,
		}
		err = s.Save(user)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	} else {

		// check password
		if user.Password != req.Pwd {
			return controller.JSONBaseErrorReq(c, errBeenBoundOrThePasswordIsWrong)
		}

		// check user login type
		if user.UserAuthenticationType != model.UserAuthenticationTypeOAUTH2 &&
			user.UserAuthenticationType != model.UserAuthenticationTypeSQLE &&
			user.UserAuthenticationType != "" {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("the user has bound other login methods")))
		}

		// check user bind third party users
		if user.ThirdPartyUserID != oauth2UserID && user.ThirdPartyUserID != "" {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("the user has bound other third-party user")))
		}

		// modify user login type
		if user.UserAuthenticationType != model.UserAuthenticationTypeOAUTH2 {
			user.ThirdPartyUserID = oauth2UserID
			user.UserAuthenticationType = model.UserAuthenticationTypeOAUTH2
			err := s.Save(user)
			if err != nil {
				return controller.JSONBaseErrorReq(c, err)
			}
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
