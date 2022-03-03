//go:build !enterprise
// +build !enterprise

package middleware

import (
	"github.com/labstack/echo/v4"
)

func licenseAdapter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}
