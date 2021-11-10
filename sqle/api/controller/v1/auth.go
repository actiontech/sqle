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

	s := model.GetStorage()
	user, exist, err := s.GetUserByName(req.UserName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist || !(req.UserName == user.Name && req.Password == user.Password) {
		return controller.JSONBaseErrorReq(c, errors.New(errors.LoginAuthFail,
			fmt.Errorf("password is wrong or user does not exist")))
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


type LoginChecker interface {
	login(userName, password string) (err error)
}

// ldapLoginV3 version 3 ldap login verification logic.
type ldapLoginV3 struct {
	config *model.LDAPConfiguration
}

func newLdapLoginV3(configuration *model.LDAPConfiguration) *ldapLoginV3 {
	return &ldapLoginV3{config: configuration}
}

func (l ldapLoginV3) login(userName, password string) (err error) {
	email, err := l.loginToLdap(userName, password)
	if err != nil {
		return err
	}
	return l.autoRegisterUser(userName, password, email)
}

func (l ldapLoginV3) loginToLdap(userName, password string) (email string, err error) {
	ldapC, _, err := model.GetStorage().GetLDAPConfiguration()
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("ldap://%s:%s", ldapC.Host, ldapC.Port)
	conn, err := ldap.DialURL(url)
	if err != nil {
		return "", fmt.Errorf("get ldap server connect failed: %v", err)
	}
	defer conn.Close()

	if err = conn.Bind(ldapC.ConnectDn, ldapC.ConnectPassword); err != nil {
		return "", fmt.Errorf("bind ldap manager user failed: %v", err)
	}
	searchRequest := ldap.NewSearchRequest(
		ldapC.BaseDn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(%s=%s)", ldapC.UserNameRdnKey, userName),
		[]string{},
		nil,
	)
	result, err := conn.Search(searchRequest)
	if err != nil {
		return "", fmt.Errorf("search user on ldap server failed: %v", err)
	}
	if len(result.Entries) != 1 {
		return "", fmt.Errorf("search user on ldap ,result size(%v) not unique", len(result.Entries))
	}
	userDn := result.Entries[0].DN
	if err = conn.Bind(userDn, password); err != nil {
		return "", fmt.Errorf("ldap login failed, username and password do not match")
	}

	return result.Entries[0].GetAttributeValue(ldapC.UserEmailRdnKey), nil
}

func (l ldapLoginV3) autoRegisterUser(userName, password, email string) (err error) {
	s := model.GetStorage()
	_, exist, err := s.GetUserByName(userName)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	user := &model.User{
		Name:     userName,
		Password: password,
		Email:    email,
	}
	return s.Save(user)
}

// sqleLogin sqle login verification logic
type sqleLogin struct {
}

func newSqleLogin() *sqleLogin {
	return &sqleLogin{}
}

func (s sqleLogin) login(userName, password string) (err error) {
	storage := model.GetStorage()
	user, exist, err := storage.GetUserByName(userName)
	if err != nil {
		return err
	}
	if !exist || !(userName == user.Name && password == user.Password) {
		return fmt.Errorf("password is wrong or user does not exist")
	}
	return nil
}
