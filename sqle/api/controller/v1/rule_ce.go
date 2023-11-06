//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/labstack/echo/v4"
)

var errCommunityEditionNotSupportCustomRule = errors.New(errors.CustomRuleEditionNotSupported, e.New("custom rule community not supported"))
var errCommunityEditionNotSupportRuleKnowledge = errors.New(errors.CustomRuleEditionNotSupported, e.New("community do not support rule knowledge"))

func getCustomRules(c echo.Context) error {
	return errCommunityEditionNotSupportCustomRule
}

func deleteCustomRule(c echo.Context) error {
	return errCommunityEditionNotSupportCustomRule
}

func createCustomRule(c echo.Context) error {
	return errCommunityEditionNotSupportCustomRule
}

func updateCustomRule(c echo.Context) error {
	return errCommunityEditionNotSupportCustomRule
}

func getCustomRule(c echo.Context) error {
	return errCommunityEditionNotSupportCustomRule
}

func getRuleTypeByDBType(c echo.Context) error {
	return errCommunityEditionNotSupportCustomRule
}

func getRuleKnowledge(c echo.Context) error {
	return errCommunityEditionNotSupportRuleKnowledge
}

func updateRuleKnowledge(c echo.Context) error {
	return errCommunityEditionNotSupportRuleKnowledge
}

func getCustomRuleKnowledge(c echo.Context) error {
	return errCommunityEditionNotSupportRuleKnowledge
}

func updateCustomRuleKnowledge(c echo.Context) error {
	return errCommunityEditionNotSupportRuleKnowledge
}
