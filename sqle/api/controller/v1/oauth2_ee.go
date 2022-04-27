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
		RedirectURL:  fmt.Sprintf("http://%v/v1/oauth2/callback", conf.ClientHost),
		Scopes:       conf.GetScopes(),
		Endpoint: oauth2.Endpoint{
			AuthURL:  conf.ServerAuthUrl,
			TokenURL: conf.ServerTokenUrl,
		},
	}
}

type callbackRedirectData struct {
	UserExist    bool
	SqleToken    string
	Oauth2UserId string
	Error        string
}

func (c callbackRedirectData) generateQuery(uri string) string {
	params := url.Values{}
	params.Set("user_exist", strconv.FormatBool(c.UserExist))
	if c.SqleToken != "" {
		params.Set("sqle_token", c.SqleToken)
	}
	if c.Oauth2UserId != "" {
		params.Set("oauth2_user_id", c.Oauth2UserId)
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
	uri := fmt.Sprintf("http://%v", oauth2C.ClientHost)
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

	//get user id using oauth2 token
	userID, err := getOauth2UserID(oauth2C, oauth2Token)
	if err != nil {
		data.Error = err.Error()
		return c.Redirect(http.StatusFound, data.generateQuery(uri))
	}
	data.Oauth2UserId = userID

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

func getOauth2UserID(conf *model.Oauth2Configuration, token *oauth2.Token) (userID string, err error) {
	uri := fmt.Sprintf("%v?%v=%v", conf.ServerUserIdUrl, conf.AccessTokenTag, token.AccessToken)
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

func bindOauth2User(c echo.Context) error {
	req := new(BindOauth2UserReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	if req.Oauth2UserID == "" {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf("oauth2 user id can not empty")))
	}

	s := model.GetStorage()
	user, exist, err := s.GetUserByName(req.UserName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// create user if not exist
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

	// check password
	if user.Password != req.Pwd {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf("wrong password")))
	}

	// check user login type
	if user.UserAuthenticationType != model.UserAuthenticationTypeOAUTH2 &&
		user.UserAuthenticationType != model.UserAuthenticationTypeSQLE &&
		user.UserAuthenticationType != "" {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("the user has bound other login methods")))
	}

	// check user bind third party users
	if user.ThirdPartyUserID != req.Oauth2UserID && user.ThirdPartyUserID != "" {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("the user has bound other third-party user")))
	}

	// modify user login type
	if user.UserAuthenticationType != model.UserAuthenticationTypeOAUTH2 {
		user.ThirdPartyUserID = req.Oauth2UserID
		user.UserAuthenticationType = model.UserAuthenticationTypeOAUTH2
		err = s.UpdateUserAuthenticationTypeByName(user.Name, model.UserAuthenticationTypeOAUTH2)
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
