//go:build enterprise
// +build enterprise

package v2

import (
	"github.com/actiontech/sqle/sqle/errors"

	"github.com/labstack/echo/v4"
)

func getTaskAnalysisData(c echo.Context) error {
	return errors.NewNotImplementedError("no impl yet")
}
