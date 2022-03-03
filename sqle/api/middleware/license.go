package middleware

import "github.com/labstack/echo/v4"

func LicenseAdapter() echo.MiddlewareFunc {
	//nolint:typecheck
	return licenseAdapter()
}
