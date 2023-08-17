package v1

import (
	"crypto/tls"
	_errors "errors"
	"fmt"
	"net/http"
	"time"

	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper"
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
func LoginV1(c echo.Context) error {
	req := new(UserLoginReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	t, err := Login(c, req.UserName, req.Password)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetUserLoginResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: UserLoginResV1{
			Token: t, // this token won't be used any more
		},
	})
}

func Login(c echo.Context, userName, password string) (token string, err error) {
	loginChecker, err := GetLoginCheckerByUserName(userName)
	if err != nil {
		return "", errors.New(errors.LoginAuthFail, err)
	}
	err = loginChecker.login(password)
	if err != nil {
		return "", errors.New(errors.LoginAuthFail, err)
	}

	token, err = generateToken(userName)
	if err != nil {
		return "", errors.New(http.StatusInternalServerError, err)
	}

	SetCookie(c, token)

	return
}

func SetCookie(c echo.Context, token string) {
	c.SetCookie(&http.Cookie{
		Name:    "sqle-token",
		Value:   token,
		Expires: time.Now().Add(tokenLifeTime),
		Path:    "/",
	})
}

const tokenLifeTime = time.Hour * 24

func generateToken(userName string) (string, error) {
	j := utils.NewJWT(utils.JWTSecretKey)
	return j.CreateToken(userName, time.Now().Add(tokenLifeTime).Unix())
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

	checkerType := loginCheckerTypeUnknown
	exist := false
	{ // get login checker type
		var u *model.User = nil
		var l *model.LDAPConfiguration = nil
		if userExist {
			u = user
		}
		if ldapExist {
			l = ldapC
		}
		checkerType, exist = getLoginCheckerType(u, l)
	}

	// match login method
	switch checkerType {
	case loginCheckerTypeLDAP:
		if !exist {
			return newLdapLoginV3WhenUserNotExist(ldapC, userName), nil
		}
		return newLdapLoginV3WhenUserExist(ldapC, user), nil
	case loginCheckerTypeSQLE:
		return newSqleLogin(user), nil
	default:
		return nil, fmt.Errorf("the user does not exist or the password is wrong")
	}
}

type checkerType int

const (
	loginCheckerTypeUnknown checkerType = iota
	loginCheckerTypeSQLE
	loginCheckerTypeLDAP
)

// determine whether the login conditions are met according to the order of login priority
func getLoginCheckerType(user *model.User, ldapC *model.LDAPConfiguration) (checkerType checkerType, userExist bool) {

	// ldap login condition
	if ldapC != nil && ldapC.Enable {
		if user != nil && user.UserAuthenticationType == model.UserAuthenticationTypeLDAP {
			return loginCheckerTypeLDAP, true
		}
		if user == nil {
			return loginCheckerTypeLDAP, false
		}
	}

	// sqle login condition, oauth 2 and other login types of users can also log in through the account and password
	if user != nil && (user.UserAuthenticationType != model.UserAuthenticationTypeLDAP) {
		return loginCheckerTypeSQLE, true
	}

	// no alternative login method
	return loginCheckerTypeUnknown, user != nil
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

func (l *ldapLoginV3) login(password string) (err error) {
	err = l.loginToLdap(password)
	if err != nil {
		return err
	}
	return l.autoRegisterUser(password)
}

var errLdapLoginFailed = _errors.New("ldap login failed, username and password do not match")

const ldapServerErrorFormat = "search user on ldap server failed: %v"

func (l *ldapLoginV3) loginToLdap(password string) (err error) {
	ldapC, _, err := model.GetStorage().GetLDAPConfiguration()
	if err != nil {
		return err
	}

	var conn *ldap.Conn
	if l.config.EnableSSL {
		url := fmt.Sprintf("ldaps://%s:%s", ldapC.Host, ldapC.Port)
		conn, err = ldap.DialURL(url, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	} else {
		url := fmt.Sprintf("ldap://%s:%s", ldapC.Host, ldapC.Port)
		conn, err = ldap.DialURL(url)
	}
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
		return fmt.Errorf(ldapServerErrorFormat, err)
	}
	if len(result.Entries) == 0 {
		return errLdapLoginFailed
	}
	if len(result.Entries) != 1 {
		return fmt.Errorf(ldapServerErrorFormat, "the queried user is not unique, please check whether the relevant configuration is correct")
	}
	userDn := result.Entries[0].DN
	if err = conn.Bind(userDn, password); err != nil {
		return errLdapLoginFailed
	}
	l.email = result.Entries[0].GetAttributeValue(ldapC.UserEmailRdnKey)
	return nil
}

func (l *ldapLoginV3) autoRegisterUser(pwd string) (err error) {
	if l.userExist {
		return l.updateUser(pwd)
	}
	return l.registerUser(pwd)
}

func (l *ldapLoginV3) registerUser(pwd string) (err error) {
	user := &model.User{
		Name:                   l.user.Name,
		Password:               pwd,
		Email:                  l.email,
		UserAuthenticationType: model.UserAuthenticationTypeLDAP,
	}
	return model.GetStorage().Save(user)
}

func (l *ldapLoginV3) updateUser(pwd string) (err error) {
	if l.user.Password == pwd {
		return nil
	}
	return model.GetStorage().UpdatePassword(l.user, pwd)
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

func (s *sqleLogin) login(password string) (err error) {
	if password != s.user.Password {
		return fmt.Errorf("password is wrong or user does not exist")
	}
	return nil
}

// LogoutV1
// @Summary 用户登出
// @Description user logout
// @Tags user
// @Id logoutV1
// @Success 200 {object} controller.BaseRes
// @router /v1/logout [post]
func LogoutV1(c echo.Context) error {
	return c.JSON(http.StatusOK, controller.NewBaseReq(logout(c)))
}

func logout(c echo.Context) error {
	cookie, err := c.Cookie("sqle-token")
	if err != nil {
		return err
	}
	cookie.MaxAge = -1 // MaxAge<0 means delete cookie now
	cookie.Path = "/"
	c.SetCookie(cookie)
	cloudbeaver_wrapper.UnbindCBSessionIdBySqleToken(cookie.Value)
	return nil
}
