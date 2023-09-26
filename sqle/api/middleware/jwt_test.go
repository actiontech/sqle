package middleware

import (
	"database/sql/driver"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	dmsCommonJwt "github.com/actiontech/dms/pkg/dms-common/api/jwt"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

func TestScannerVerifier(t *testing.T) {
	e := echo.New()

	apName := "test_audit_plan"
	projectName := "test_project"
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
		ctx.SetParamNames("audit_plan_name", "project_name")
		ctx.SetParamValues(apName, projectName)
		return ctx, res
	}

	{ // test audit plan name don't match the token
		token, err := dmsCommonJwt.GenJwtToken(dmsCommonJwt.WithUserName(testUser), dmsCommonJwt.WithExpiredTime(1*time.Hour), dmsCommonJwt.WithAuditPlanName(apName))
		assert.NoError(t, err)
		ctx, _ := newContextFunc(token, fmt.Sprintf("%s_modified", apName))
		err = mw(h)(ctx)
		assert.Contains(t, err.Error(), errAuditPlanMisMatch.Error())
	}

	{ // test unknown token
		token, err := dmsCommonJwt.GenJwtToken(dmsCommonJwt.WithUserName(testUser), dmsCommonJwt.WithExpiredTime(1*time.Hour))
		assert.NoError(t, err)
		ctx, _ := newContextFunc(token, apName)
		err = mw(h)(ctx)
		assert.Contains(t, err.Error(), "unknown token")
	}

	{ // test audit plan token incorrect
		token, err := dmsCommonJwt.GenJwtToken(dmsCommonJwt.WithUserName(testUser), dmsCommonJwt.WithExpiredTime(1*time.Hour), dmsCommonJwt.WithAuditPlanName(apName))
		assert.NoError(t, err)

		mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		model.InitMockStorage(mockDB)
		mock.ExpectQuery("SELECT `audit_plans`.* FROM `audit_plans` LEFT JOIN projects ON projects.id = audit_plans.project_id WHERE `audit_plans`.`deleted_at` IS NULL AND ((projects.name = ? AND audit_plans.name = ?))").
			WithArgs(projectName, apName).
			WillReturnRows(sqlmock.NewRows([]string{"name", "token"}).AddRow(driver.Value(testUser), "test-token"))
		mock.ExpectClose()

		ctx, _ := newContextFunc(token, apName)
		err = mw(h)(ctx)
		assert.Contains(t, err.Error(), errAuditPlanMisMatch.Error())

		mockDB.Close()
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	}

	{ // test audit plan not found
		token, err := dmsCommonJwt.GenJwtToken(dmsCommonJwt.WithUserName(testUser), dmsCommonJwt.WithExpiredTime(1*time.Hour), dmsCommonJwt.WithAuditPlanName(apName))
		assert.NoError(t, err)

		mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		model.InitMockStorage(mockDB)
		mock.ExpectQuery("SELECT `audit_plans`.* FROM `audit_plans` LEFT JOIN projects ON projects.id = audit_plans.project_id WHERE `audit_plans`.`deleted_at` IS NULL AND ((projects.name = ? AND audit_plans.name = ?))").
			WithArgs(projectName, apName).
			WillReturnError(gorm.ErrRecordNotFound)
		assert.NoError(t, err)
		model.InitMockStorage(mockDB)
		mock.ExpectClose()

		ctx, _ := newContextFunc(token, apName)
		err = mw(h)(ctx)
		assert.Contains(t, err.Error(), errAuditPlanMisMatch.Error())

		mockDB.Close()
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	}

	{ // test check success
		token, err := dmsCommonJwt.GenJwtToken(dmsCommonJwt.WithUserName(testUser), dmsCommonJwt.WithExpiredTime(1*time.Hour), dmsCommonJwt.WithAuditPlanName(apName))
		assert.NoError(t, err)

		mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		model.InitMockStorage(mockDB)
		mock.ExpectQuery("SELECT `audit_plans`.* FROM `audit_plans` LEFT JOIN projects ON projects.id = audit_plans.project_id WHERE `audit_plans`.`deleted_at` IS NULL AND ((projects.name = ? AND audit_plans.name = ?))").
			WithArgs(projectName, apName).
			WillReturnRows(sqlmock.NewRows([]string{"name", "token"}).AddRow(testUser, token))
		mock.ExpectClose()

		ctx, res := newContextFunc(token, apName)
		err = mw(h)(ctx)
		assert.NoError(t, err)
		assert.Contains(t, res.Body.String(), "hello, world")

		mockDB.Close()
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	}

	{ // test default auth scheme success
		token, err := dmsCommonJwt.GenJwtToken(dmsCommonJwt.WithUserName(testUser), dmsCommonJwt.WithExpiredTime(1*time.Hour), dmsCommonJwt.WithAuditPlanName(apName))
		assert.NoError(t, err)

		mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		model.InitMockStorage(mockDB)
		mock.ExpectQuery("SELECT `audit_plans`.* FROM `audit_plans` LEFT JOIN projects ON projects.id = audit_plans.project_id WHERE `audit_plans`.`deleted_at` IS NULL AND ((projects.name = ? AND audit_plans.name = ?))").
			WithArgs(projectName, apName).
			WillReturnRows(sqlmock.NewRows([]string{"name", "token"}).AddRow(testUser, token))
		mock.ExpectClose()

		tokenWithSchema := fmt.Sprintf("%s %s", middleware.DefaultJWTConfig.AuthScheme, token)
		ctx, res := newContextFunc(tokenWithSchema, apName)
		err = mw(h)(ctx)
		assert.NoError(t, err)
		assert.Contains(t, res.Body.String(), "hello, world")

		mockDB.Close()
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	}
}
