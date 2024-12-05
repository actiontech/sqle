//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/labstack/echo/v4"
)

func getTaskSQLRewrittenData(c echo.Context) error {
	return errors.New(errors.EnterpriseEditionFeatures, e.New("sql rewriting is enterprise version feature"))
}
