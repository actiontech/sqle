//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/labstack/echo/v4"
)

var ErrCommunityEditionNotSupportSqlManage = errors.New(errors.EnterpriseEditionFeatures, e.New("sql manage is enterprise version feature"))

func getSqlManageList(c echo.Context) error {
	return ErrCommunityEditionNotSupportSqlManage
}

func batchUpdateSqlManage(c echo.Context) error {
	return ErrCommunityEditionNotSupportSqlManage
}

func exportSqlManagesV1(c echo.Context) error {
	return ErrCommunityEditionNotSupportSqlManage
}

func getSqlManageRuleTips(c echo.Context) error {
	return ErrCommunityEditionNotSupportSqlManage
}

func getSqlManageSqlAnalysisV1(c echo.Context) error {
	return ErrCommunityEditionNotSupportSqlManage
}
