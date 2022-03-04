//go:build release
// +build release

package middleware

import (
	e "errors"
	"sync"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/license"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
)

var ErrLicenseVerifyFailed = errors.New(errors.ErrAccessDeniedError, e.New("the operation is outside the scope of the license"))

var ErrInitLicenseFailed = errors.New(errors.ConnectStorageError, e.New("license information initialization failed, will retry initialization the next time the license is verified"))

var once = &sync.Once{}

func initLicense() error {
	s := model.GetStorage()
	l, _, err := s.GetLicense()
	if err != nil {
		return err
	}
	return license.InitChecker(l.Content)
}

func licenseAdapter() echo.MiddlewareFunc {
	initFailed := false
	once.Do(func() {
		if initLicense() != nil {
			once = &sync.Once{}
			initFailed = true
		}
	})

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			pass, err := license.Check(c)
			if err != nil {
				return ErrInitLicenseFailed
			}
			if !pass {
				return ErrLicenseVerifyFailed
			}
			if initFailed {
				return ErrInitLicenseFailed
			}
			return next(c)
		}
	}
}
