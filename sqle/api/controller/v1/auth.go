package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/go-ldap/ldap/v3"
	"github.com/labstack/echo/v4"
)

type UserLoginReqV1 struct {
	UserName string `json:"username" form:"username" example:"test" valid:"required"`
	Password string `json:"password" form:"password" example:"123456" valid:"required"`
}

type GetUserLoginResV1 struct {
	controller.BaseRes
	Data UserLoginResV1 `json:"data"`
}

type UserLoginResV1 struct {
	Token string `json:"token" example:"this is a jwt token string"`
}

// @Summary 用户登录
// @Description user login
// @Tags user
// @Id loginV1
// @Param user body v1.UserLoginReqV1 true "user login request"
// @Success 200 {object} v1.GetUserLoginResV1
// @router /v1/login [post]
func Login(c echo.Context) error {
	req := new(UserLoginReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	loginChecker, err := GetLoginCheckerByUserName(req.UserName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.LoginAuthFail, err))
	}
	err = loginChecker.login(req.Password)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.LoginAuthFail, err))
	}

	j := utils.NewJWT([]byte(utils.JWTSecret))
	t, err := j.CreateToken(req.UserName, time.Now().Add(time.Hour*24).Unix())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, &GetUserLoginResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: UserLoginResV1{
			Token: t,
		},
	})
}

// GetLoginCheckerByUserName get login checker by user name and init login checker
func GetLoginCheckerByUserName(userName string) (LoginChecker, error) {
	// get user metadata and config
	s := model.GetStorage()
	user, userExist, err := s.GetUserByName(userName)
	if err != nil {
		return nil, err
	}
	ldapC, ldapExist, err := s.GetLDAPConfiguration()
	if err != nil {
		return nil, err
	}

	state := loginCheckerTypeUnknown
	{ // fix user login method
		var u *model.User = nil
		var l *model.LDAPConfiguration = nil
		if userExist {
			u = user
		}
		if ldapExist {
			l = ldapC
		}
		state = getLoginCheckerType(u, l)
	}

	// match login method
	switch state {
	case loginCheckerTypeLDAP:
		return newLdapLoginV3WhenUserExist(ldapC, user), nil
	case loginCheckerTypeLDAPUserNotExist:
		return newLdapLoginV3WhenUserNotExist(ldapC, userName), nil
	case loginCheckerTypeSQLE:
		return newSqleLogin(user), nil
	default:
		return nil, fmt.Errorf("the user does not exist or the password is wrong")
	}
}

type userState int

const (
	loginCheckerTypeUnknown userState = iota
	loginCheckerTypeSQLE
	loginCheckerTypeLDAP
	loginCheckerTypeLDAPUserNotExist
)

// determine whether the login conditions are met according to the order of login priority
func getLoginCheckerType(user *model.User, ldapC *model.LDAPConfiguration) userState {

	// ldap login condition
	if ldapC != nil && ldapC.Enable {
		if user != nil && user.UserAuthenticationType == model.UserAuthenticationTypeLDAP {
			return loginCheckerTypeLDAP
		}
		if user == nil {
			return loginCheckerTypeLDAPUserNotExist
		}
	}

	// sqle login condition
	if user != nil && (user.UserAuthenticationType == model.UserAuthenticationTypeSQLE || user.UserAuthenticationType == "") {
		return loginCheckerTypeSQLE
	}

	// no alternative login method
	return loginCheckerTypeUnknown
}

type LoginChecker interface {
	login(password string) (err error)
}

type baseLoginChecker struct {
	user *model.User
}

// ldapLoginV3 version 3 ldap login verification logic.
type ldapLoginV3 struct {
	baseLoginChecker
	config    *model.LDAPConfiguration
	email     string
	userExist bool
}

func newLdapLoginV3WhenUserExist(configuration *model.LDAPConfiguration, user *model.User) *ldapLoginV3 {
	return &ldapLoginV3{
		config:    configuration,
		userExist: true,
		baseLoginChecker: baseLoginChecker{
			user: user,
		},
	}
}

func newLdapLoginV3WhenUserNotExist(configuration *model.LDAPConfiguration, userName string) *ldapLoginV3 {
	return &ldapLoginV3{
		config:    configuration,
		userExist: false,
		baseLoginChecker: baseLoginChecker{
			user: &model.User{
				Name: userName,
			},
		},
	}
}

func (l ldapLoginV3) login(password string) (err error) {
	err = l.loginToLdap(password)
	if err != nil {
		return err
	}
	return l.autoRegisterUser()
}

func (l ldapLoginV3) loginToLdap(password string) (err error) {
	ldapC, _, err := model.GetStorage().GetLDAPConfiguration()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("ldap://%s:%s", ldapC.Host, ldapC.Port)
	conn, err := ldap.DialURL(url)
	if err != nil {
		return fmt.Errorf("get ldap server connect failed: %v", err)
	}
	defer conn.Close()

	if err = conn.Bind(ldapC.ConnectDn, ldapC.ConnectPassword); err != nil {
		return fmt.Errorf("bind ldap manager user failed: %v", err)
	}
	searchRequest := ldap.NewSearchRequest(
		ldapC.BaseDn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(%s=%s)", ldapC.UserNameRdnKey, l.user.Name),
		[]string{},
		nil,
	)
	result, err := conn.Search(searchRequest)
	if err != nil {
		return fmt.Errorf("search user on ldap server failed: %v", err)
	}
	if len(result.Entries) != 1 {
		return fmt.Errorf("search user on ldap ,result size(%v) not unique", len(result.Entries))
	}
	userDn := result.Entries[0].DN
	if err = conn.Bind(userDn, password); err != nil {
		return fmt.Errorf("ldap login failed, username and password do not match")
	}
	l.email = result.Entries[0].GetAttributeValue(ldapC.UserEmailRdnKey)
	return nil
}

func (l ldapLoginV3) autoRegisterUser() (err error) {
	if l.userExist {
		return nil
	}
	user := &model.User{
		Name:                   l.user.Name,
		Password:               "this password will not be used",
		Email:                  l.email,
		UserAuthenticationType: model.UserAuthenticationTypeLDAP,
	}
	return model.GetStorage().Save(user)
}

// sqleLogin sqle login verification logic
type sqleLogin struct {
	baseLoginChecker
}

func newSqleLogin(user *model.User) *sqleLogin {
	return &sqleLogin{
		baseLoginChecker: baseLoginChecker{
			user: user,
		},
	}
}

func (s sqleLogin) login(password string) (err error) {
	if password != s.user.Password {
		return fmt.Errorf("password is wrong or user does not exist")
	}
	return nil
}
