//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/labstack/echo/v4"
)

var errCommunityEditionNotSupportCustomRule = errors.New(errors.CustomRuleEditionNotSupported, e.New("custom rule community not supported"))


func getCustomRules(c echo.Context) error {
	return errCommunityEditionNotSupportCustomRule
}

func deleteCustomRule(c echo.Context) error {
	return errCommunityEditionNotSupportCustomRule
}