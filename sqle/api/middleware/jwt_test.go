package middleware

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	dmsCommon "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	dmsCommonJwt "github.com/actiontech/dms/pkg/dms-common/api/jwt"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

func mockServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		response := dmsCommon.ListProjectReply{
			Data: []*dmsCommon.ListProject{{
				ProjectUid: "700300",
				Name:       "default",
				Archived:   false,
			}},
			Total: 1,
		}
		res, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("response err %v", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		_, err = w.Write(res)
		if err != nil {
			log.Printf("response err %v ", err)
		}
	}
	return httptest.NewServer(http.HandlerFunc(f))
}

func TestScannerVerifier(t *testing.T) {
	server := mockServer()
	defer server.Close()
	controller.InitDMSServerAddress(server.URL)

	e := echo.New()

	apName := "test_audit_plan"
	projectUID := "700300"
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
		ctx.SetParamValues(apName, projectUID)
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
		mock.ExpectQuery("SELECT * FROM `audit_plans` WHERE `audit_plans`.`deleted_at` IS NULL AND ((project_id = ? AND name = ?))").
			WithArgs(projectUID, apName).
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
		mock.ExpectQuery("SELECT * FROM `audit_plans` WHERE `audit_plans`.`deleted_at` IS NULL AND ((project_id = ? AND name = ?))").
			WithArgs(projectUID, apName).
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
		mock.ExpectQuery("SELECT * FROM `audit_plans` WHERE `audit_plans`.`deleted_at` IS NULL AND ((project_id = ? AND name = ?))").
			WithArgs(projectUID, apName).
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
		mock.ExpectQuery("SELECT * FROM `audit_plans` WHERE `audit_plans`.`deleted_at` IS NULL AND ((project_id = ? AND name = ?))").
			WithArgs(projectUID, apName).
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

func TestScannerVerifierIssue1758(t *testing.T) {
	server := mockServer()
	defer server.Close()
	controller.InitDMSServerAddress(server.URL)
	e := echo.New()

	jwt := utils.NewJWT(utils.JWTSecretKey)
	apName120 := "test_name_length_120_000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	projectUid := "700300"
	userName := "admin"
	assert.Equal(t, 120, len(apName120))
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
		ctx.SetParamValues(apName, projectUid)
		return ctx, res
	}
	{ // test check success
		token, err := dmsCommonJwt.GenJwtToken(dmsCommonJwt.WithUserName(utils.Md5(userName)), dmsCommonJwt.WithExpiredTime(1*time.Hour), dmsCommonJwt.WithAuditPlanName(utils.Md5(apName120)))
		assert.NoError(t, err)

		mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		model.InitMockStorage(mockDB)
		mock.ExpectQuery("SELECT * FROM `audit_plans` WHERE `audit_plans`.`deleted_at` IS NULL AND ((project_id = ? AND name = ?))").
			WithArgs(projectUid, apName120).
			WillReturnRows(sqlmock.NewRows([]string{"name", "token"}).AddRow(userName, token))
		mock.ExpectClose()

		ctx, res := newContextFunc(token, apName120) //这里模拟上下文不需要哈希
		err = mw(h)(ctx)
		assert.NoError(t, err)
		assert.Contains(t, res.Body.String(), "hello, world")

		mockDB.Close()
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	}
	{ // test audit plan name don't match the token
		token, err := jwt.CreateToken(utils.Md5(userName), time.Now().Add(1*time.Hour).Unix(), utils.WithAuditPlanName(utils.Md5(apName120)))
		assert.NoError(t, err)
		ctx, _ := newContextFunc(token, fmt.Sprintf("%s_modified", apName120))
		err = mw(h)(ctx)
		assert.Contains(t, err.Error(), errAuditPlanMisMatch.Error())
	}
	{ // test unknown token
		token, err := jwt.CreateToken(utils.Md5(userName), time.Now().Add(1*time.Hour).Unix())
		assert.NoError(t, err)
		ctx, _ := newContextFunc(token, apName120)
		err = mw(h)(ctx)
		assert.Contains(t, err.Error(), "unknown token")
	}
	{ // test old token
		token, err := jwt.CreateToken(userName, time.Now().Add(1*time.Hour).Unix(), utils.WithAuditPlanName(apName120))
		assert.NoError(t, err)
		mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		model.InitMockStorage(mockDB)
		mock.ExpectQuery("SELECT * FROM `audit_plans` WHERE `audit_plans`.`deleted_at` IS NULL AND ((project_id = ? AND name = ?))").
			WithArgs(projectUid, apName120).
			WillReturnRows(sqlmock.NewRows([]string{"name", "token"}).AddRow(userName, token))
		mock.ExpectClose()

		ctx, _ := newContextFunc(token, apName120)
		err = mw(h)(ctx)
		assert.NoError(t, err)

		mockDB.Close()
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	}
}
