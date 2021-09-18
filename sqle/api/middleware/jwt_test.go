package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/actiontech/sqle/sqle/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

func TestScannerVerifier(t *testing.T) {
	e := echo.New()

	jwt := utils.NewJWT([]byte(utils.JWTSecret))
	apName := "test_audit_plan"
	testUser := "test_user"

	h := func(c echo.Context) error {
		return c.HTML(http.StatusOK, "hello, world")
	}

	mw := ScannerVerifier()

	newContextFunc := func(token, apName string) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodGet, "/:audit_plan_name/", nil)
		req.Header.Set(echo.HeaderAuthorization, token)
		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)
		ctx.SetParamNames("audit_plan_name")
		ctx.SetParamValues(apName)
		return ctx, res
	}

	{
		token, err := jwt.CreateToken(testUser, time.Now().Add(1*time.Hour).Unix(), utils.WithAuditPlanName(apName))
		assert.NoError(t, err)
		ctx, _ := newContextFunc(token, fmt.Sprintf("%s_modified", apName))
		err = mw(h)(ctx)
		assert.Contains(t, err.Error(), errAuditPlanMisMatch.Error())
	}

	{
		token, err := jwt.CreateToken(testUser, time.Now().Add(1*time.Hour).Unix())
		assert.NoError(t, err)
		ctx, _ := newContextFunc(token, apName)
		err = mw(h)(ctx)
		assert.Contains(t, err.Error(), "unknown token")
	}

	{
		token, err := jwt.CreateToken(testUser, time.Now().Add(1*time.Hour).Unix(), utils.WithAuditPlanName(apName))
		assert.NoError(t, err)
		ctx, res := newContextFunc(token, apName)
		err = mw(h)(ctx)
		assert.NoError(t, err)
		assert.Contains(t, string(res.Body.Bytes()), "hello, world")
	}

	{
		token, err := jwt.CreateToken(testUser, time.Now().Add(1*time.Hour).Unix(), utils.WithAuditPlanName(apName))
		tokenWithSchema := fmt.Sprintf("%s %s", middleware.DefaultJWTConfig.AuthScheme, token)
		assert.NoError(t, err)
		ctx, res := newContextFunc(tokenWithSchema, apName)
		err = mw(h)(ctx)
		assert.NoError(t, err)
		assert.Contains(t, string(res.Body.Bytes()), "hello, world")
	}
}
