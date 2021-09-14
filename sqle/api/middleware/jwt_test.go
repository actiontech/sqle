package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"actiontech.cloud/sqle/sqle/sqle/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAuditPlanVerifyAdapter(t *testing.T) {
	e := echo.New()

	jwt := utils.NewJWT([]byte(utils.JWTSecret))
	apName := "test_audit_plan"

	h := func(c echo.Context) error {
		return c.HTML(http.StatusOK, "hello, world")
	}

	mw := AuditPlanVerifyAdapter()

	{
		token, err := jwt.CreateToken("test_user", time.Now().Add(1*time.Hour).Unix(), utils.WithAuditPlanName(apName))
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/:audit_plan_name/", nil)
		req.Header.Set(echo.HeaderAuthorization, token)
		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)
		ctx.SetParamNames("audit_plan_name")
		ctx.SetParamValues("test_audit_plan_modified")

		err = mw(h)(ctx)
		assert.Contains(t, err.Error(), errAuditPlanMisMatch.Error())
	}

	{
		token, err := jwt.CreateToken("test_user", time.Now().Add(1*time.Hour).Unix())
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/:audit_plan_name/", nil)
		req.Header.Set(echo.HeaderAuthorization, token)
		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)

		err = mw(h)(ctx)
		assert.Contains(t, err.Error(), "unknown token")
	}

	{
		token, err := jwt.CreateToken("test_user", time.Now().Add(1*time.Hour).Unix(), utils.WithAuditPlanName(apName))
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/:audit_plan_name/", nil)
		req.Header.Set(echo.HeaderAuthorization, token)
		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)
		ctx.SetPath("/:audit_plan_name/")
		ctx.SetParamNames("audit_plan_name")
		ctx.SetParamValues(apName)

		err = mw(h)(ctx)
		assert.NoError(t, err)
		assert.Contains(t, string(res.Body.Bytes()), "hello, world")
	}
}
