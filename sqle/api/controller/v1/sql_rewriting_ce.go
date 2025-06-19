//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/labstack/echo/v4"
)

func getRewriteSQLData(c echo.Context) error {
	return errors.New(errors.EnterpriseEditionFeatures, e.New("sql rewriting is enterprise version feature"))
}

// getAsyncRewriteTaskStatus 获取异步重写任务状态
func getAsyncRewriteTaskStatus(c echo.Context) error {
	return errors.New(errors.EnterpriseEditionFeatures, e.New("sql rewriting is enterprise version feature"))
}
