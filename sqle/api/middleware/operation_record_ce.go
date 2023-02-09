//go:build !enterprise
// +build !enterprise

package middleware

import (
	"github.com/labstack/echo/v4"
)

func OperationLogRecord() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			return next(c)
		}
	}
}
