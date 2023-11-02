//go:build release
// +build release

package middleware

// import (
// 	"bytes"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"strings"

// 	"github.com/actiontech/sqle/sqle/api/controller"
// 	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
// 	"github.com/actiontech/sqle/sqle/errors"
// 	"github.com/actiontech/sqle/sqle/license"
// 	"github.com/actiontech/sqle/sqle/log"
// 	"github.com/actiontech/sqle/sqle/model"

// 	"github.com/labstack/echo/v4"
// )

// func licenseAdapter() echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			check := &LicenseChecker{ctx: c}
// 			err := check.Check()
// 			if err != nil {
// 				return errors.New(errors.ErrAccessDeniedError, fmt.Errorf("the operation is outside the scope of the license, %v", err))
// 			}
// 			return next(c)
// 		}
// 	}
// }

// type LicenseCheckerFn func() (skip bool, err error)

// type LicenseChecker struct {
// 	ctx     echo.Context
// 	license *license.License
// }

// func (c *LicenseChecker) getLicense() (*license.License, error) {
// 	if c.license != nil {
// 		return c.license, nil
// 	}
// 	s := model.GetStorage()
// 	l, exist, err := s.GetLicense()
// 	if err != nil {
// 		return nil, err
// 	}
// 	if !exist || l.Content == nil {
// 		return nil, fmt.Errorf("license is empty")
// 	}
// 	c.license = l.Content
// 	return c.license, nil
// }

// func (c *LicenseChecker) Check() error {
// 	var checks = []LicenseCheckerFn{
// 		c.allowGetMethod,
// 		c.allowApi,
// 		c.checkHardware,
// 		c.checkIsExpired,
// 		c.checkCreateUser,
// 		c.checkCreateInstance,
// 	}
// 	for _, check := range checks {
// 		skip, err := check()
// 		if err != nil {
// 			return err
// 		}
// 		if skip {
// 			return nil
// 		}
// 	}
// 	return nil
// }

// const (
// 	LimitTypeUser = "user"
// )

// var logger = log.NewEntry().WithField("action", "check-license")

// func (c *LicenseChecker) allowGetMethod() (bool, error) {
// 	if c.ctx.Request().Method == http.MethodGet {
// 		return true, nil
// 	}
// 	return false, nil
// }

// func (c *LicenseChecker) allowApi() (bool, error) {
// 	path := strings.TrimSuffix(c.ctx.Path(), "/")
// 	switch path {
// 	case "/v1/login",
// 		"/v2/login",
// 		"/v1/logout",
// 		"/v1/configurations/license",
// 		"/v1/configurations/license/check",
// 		"/v1/configurations/license/info":
// 		return true, nil
// 	default:
// 		return false, nil
// 	}
// }

// func (c *LicenseChecker) checkHardware() (bool, error) {
// 	l, err := c.getLicense()
// 	if err != nil {
// 		return true, err
// 	}
// 	s, err := license.CollectHardwareInfo()
// 	if err != nil {
// 		return false, err
// 	}
// 	err = l.CheckHardwareSignIsMatch(s)
// 	if err != nil {
// 		return false, err
// 	}
// 	return false, nil
// }

// func (c *LicenseChecker) checkIsExpired() (bool, error) {
// 	l, err := c.getLicense()
// 	if err != nil {
// 		return true, err
// 	}
// 	if err := l.CheckLicenseNotExpired(); err != nil {
// 		return true, err
// 	}
// 	return false, nil
// }

// func (c *LicenseChecker) checkCreateUser() (bool, error) {
// 	path := strings.TrimSuffix(c.ctx.Path(), "/")
// 	if c.ctx.Request().Method == http.MethodPost && path == "/v1/users" {
// 		s := model.GetStorage()
// 		count, err := s.GetAllUserCount()
// 		if err != nil {
// 			return true, err
// 		}
// 		l, err := c.getLicense()
// 		if err != nil {
// 			return true, err
// 		}
// 		err = l.CheckCanCreateUser(count)
// 		if err != nil {
// 			return true, err
// 		}
// 	}
// 	return false, nil
// }

// func (c *LicenseChecker) checkCreateInstance() (bool, error) {
// 	path := strings.TrimSuffix(c.ctx.Path(), "/")
// 	if !(c.ctx.Request().Method == http.MethodPost &&
// 		path == "/v1/projects/:project_name/instances" || path == "/v2/projects/:project_name/instances") {
// 		return false, nil
// 	}
// 	l, err := c.getLicense()
// 	if err != nil {
// 		return true, err
// 	}
// 	dbType, err := getDBTypeWithReq(c.ctx)
// 	if err != nil {
// 		return true, err
// 	}

// 	s := model.GetStorage()
// 	counts, err := s.GetAllInstanceCount()
// 	if err != nil {
// 		return true, err
// 	}

// 	usage := make(license.LimitOfEachType, len(counts))
// 	for _, v := range counts {
// 		usage[v.DBType] = license.LimitOfType{
// 			DBType: v.DBType,
// 			Count:  int(v.Count),
// 		}
// 	}
// 	err = l.CheckCanCreateInstance(dbType, usage)
// 	if err != nil {
// 		return true, err
// 	}
// 	return false, nil
// }

// // 复制结构体是为了防止循环依赖
// type CreateInstanceReqV1 struct {
// 	Name     string `json:"instance_name" form:"instance_name" example:"test" valid:"required,name"`
// 	DBType   string `json:"db_type" form:"db_type" example:"mysql"`
// 	User     string `json:"db_user" form:"db_user" example:"root" valid:"required"`
// 	Host     string `json:"db_host" form:"db_host" example:"10.10.10.10" valid:"required,ip_addr|uri|hostname|hostname_rfc1123"`
// 	Port     string `json:"db_port" form:"db_port" example:"3306" valid:"required,port"`
// 	Password string `json:"db_password" form:"db_password" example:"123456" valid:"required"`
// }

// func getDBTypeWithReq(c echo.Context) (dbType string, err error) {
// 	bodyBytes, _ := ioutil.ReadAll(c.Request().Body)
// 	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
// 	req := new(CreateInstanceReqV1)
// 	err = controller.BindAndValidateReq(c, req)
// 	if err != nil {
// 		return "", err
// 	}
// 	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

// 	dbType = fmt.Sprintf("%v", req.DBType)
// 	if dbType == "" {
// 		dbType = driverV2.DriverTypeMySQL
// 	}
// 	return dbType, err
// }
