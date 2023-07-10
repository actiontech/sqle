//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/labstack/echo/v4"
)

func directGetSQLAnalysis(c echo.Context) error {
	return errors.New(errors.SQLAnalysisCommunityNotSupported, e.New("sql analysis community not supported"))
}
