//go:build !enterprise
// +build !enterprise

package v2

import (
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	"github.com/labstack/echo/v4"
)

func getSqlManageList(c echo.Context) error {
	return v1.ErrCommunityEditionNotSupportSqlManage
}
