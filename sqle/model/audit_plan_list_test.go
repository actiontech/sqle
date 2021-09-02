package model

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestStorage_GetAuditPlansByReq(t *testing.T) {
	// 1. test for common user
	tableAndRowOfSQL := `
	FROM audit_plans
	WHERE audit_plans.deleted_at IS NULL 
	AND audit_plans.create_user_id = ? 
	AND audit_plans.db_type = ?
	`
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	initMockStorage(mockDB)
	mock.ExpectPrepare(fmt.Sprintf(`SELECT audit_plans.name, audit_plans.cron_expression, audit_plans.db_type, audit_plans.token, audit_plans.instance_name, audit_plans.instance_database %v LIMIT ? OFFSET ?`, tableAndRowOfSQL)).
		ExpectQuery().WithArgs(1, "mysql", 100, 10).WillReturnRows(sqlmock.NewRows([]string{"name", "cron_expression", "db_type", "token", "instance_name", "instance_database"}).AddRow("audit_plan_1", "* */2 * * *", "mysql", "fake token", "inst_1", "db_1"))
	mock.ExpectPrepare(fmt.Sprintf(`SELECT COUNT(*) %v`, tableAndRowOfSQL)).
		ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{
		"COUNT(*)",
	}).AddRow("2"))
	nameFields := map[string]interface{}{
		"current_user_id":           1,
		"filter_audit_plan_db_type": "mysql",
		"limit":                     100,
		"offset":                    10}
	reslut, count, err := GetStorage().GetAuditPlansByReq(nameFields)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), count)
	assert.Len(t, reslut, 1)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
	mockDB.Close()

	// 2. test for admin user
	tableAndRowOfSQL1 := `
	FROM audit_plans
	WHERE audit_plans.deleted_at IS NULL
	`
	mockDB, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	initMockStorage(mockDB)
	mock.ExpectPrepare(fmt.Sprintf(`
	SELECT audit_plans.name, audit_plans.cron_expression, audit_plans.db_type, audit_plans.token, audit_plans.instance_name, audit_plans.instance_database
	%v
	LIMIT ? OFFSET ?`, tableAndRowOfSQL1)).
		ExpectQuery().WithArgs(100, 10).WillReturnRows(sqlmock.NewRows([]string{
		"name", "cron_expression", "db_type", "token", "instance_name", "instance_database",
	}).AddRow("audit_plan_1", "* */2 * * *", "mysql", "fake token", "inst_1", "db_1"))
	mock.ExpectPrepare(fmt.Sprintf(`SELECT COUNT(*) %v`, tableAndRowOfSQL1)).
		ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow("2"))
	nameFields = map[string]interface{}{
		"current_user_id":       1,
		"current_user_is_admin": true,
		"limit":                 100,
		"offset":                10}
	reslut, count, err = GetStorage().GetAuditPlansByReq(nameFields)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), count)
	assert.Len(t, reslut, 1)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
	mockDB.Close()
}

func TestStorage_GetAuditPlanSQLsByReq(t *testing.T) {
	tableAndRowOfSQL := `
	FROM audit_plan_sqls
	JOIN audit_plans ON audit_plans.id = audit_plan_sqls.audit_plan_id
	WHERE audit_plan_sqls.deleted_at IS NULL
	AND audit_plans.name = ?
	`
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()
	initMockStorage(mockDB)
	mock.ExpectPrepare(fmt.Sprintf(`SELECT audit_plan_sqls.fingerprint, audit_plan_sqls.counter, audit_plan_sqls.last_sql, audit_plan_sqls.last_receive_timestamp %v LIMIT ? OFFSET ?`, tableAndRowOfSQL)).
		ExpectQuery().WithArgs("audit_plan_for_jave_repo", 100, 10).WillReturnRows(sqlmock.NewRows([]string{
		"fingerprint", "counter", "last_sql", "last_receive_timestamp",
	}).AddRow("select * from t1 where id = ?", "3", "select * from t1 where id = 1", "2021-09-01T13:46:13+08:00"))
	mock.ExpectPrepare(fmt.Sprintf(`SELECT COUNT(*) %v`, tableAndRowOfSQL)).
		ExpectQuery().WithArgs("audit_plan_for_jave_repo").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow("2"))
	nameFields := map[string]interface{}{
		"audit_plan_name": "audit_plan_for_jave_repo",
		"limit":           100,
		"offset":          10}
	reslut, count, err := GetStorage().GetAuditPlanSQLsByReq(nameFields)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), count)
	assert.Len(t, reslut, 1)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestStorage_GetAuditPlanReportsByReq(t *testing.T) {
	tableAndRowOfSQL := `
	FROM audit_plan_reports
	JOIN audit_plans ON audit_plans.id = audit_plan_reports.audit_plan_id
	WHERE audit_plan_reports.deleted_at IS NULL
	AND audit_plans.name = ?
	`
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()
	initMockStorage(mockDB)
	mock.ExpectPrepare(fmt.Sprintf(`SELECT audit_plan_reports.id, audit_plan_reports.created_at %v LIMIT ? OFFSET ?`, tableAndRowOfSQL)).
		ExpectQuery().WithArgs("audit_plan_for_jave_repo", 100, 10).WillReturnRows(sqlmock.NewRows([]string{
		"id", "created_at"}).AddRow("1", "2021-09-01T13:46:13+08:00"))

	mock.ExpectPrepare(fmt.Sprintf(`SELECT COUNT(*) %v`, tableAndRowOfSQL)).
		ExpectQuery().WithArgs("audit_plan_for_jave_repo").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow("2"))
	nameFields := map[string]interface{}{
		"audit_plan_name": "audit_plan_for_jave_repo",
		"limit":           100,
		"offset":          10}
	reslut, count, err := GetStorage().GetAuditPlanReportsByReq(nameFields)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), count)
	assert.Len(t, reslut, 1)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestStorage_GetAuditPlanReportSQLsByReq(t *testing.T) {
	tableAndRowOfSQL := `
	FROM audit_plan_report_sqls
	JOIN audit_plan_reports ON audit_plan_report_sqls.audit_plan_report_id = audit_plan_reports.id
	JOIN audit_plans ON audit_plan_reports.audit_plan_id = audit_plans.id
	JOIN audit_plan_sqls ON audit_plan_sqls.id = audit_plan_report_sqls.audit_plan_sql_id
	WHERE audit_plan_report_sqls.deleted_at IS NULL
	AND audit_plans.name = ?
	AND audit_plan_reports.id = ?
	`

	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()
	initMockStorage(mockDB)
	mock.ExpectPrepare(fmt.Sprintf(`SELECT audit_plan_report_sqls.audit_result, audit_plan_sqls.fingerprint, audit_plan_sqls.last_sql, audit_plan_sqls.last_receive_timestamp %v LIMIT ? OFFSET ?`, tableAndRowOfSQL)).
		ExpectQuery().WithArgs("audit_plan_for_jave_repo", 1, 100, 10).WillReturnRows(sqlmock.NewRows([]string{
		"audit_result", "fingerprint", "last_sql", "last_receive_timestamp",
	}).AddRow("FAKE AUDIT RESULT", "select * from t1 where id = ?", "select * from t1 where id = 1", "2021-09-01T13:46:13+08:00"))

	mock.ExpectPrepare(fmt.Sprintf(`SELECT COUNT(*) %v`, tableAndRowOfSQL)).
		ExpectQuery().WithArgs("audit_plan_for_jave_repo", 1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow("2"))
	nameFields := map[string]interface{}{
		"audit_plan_name":      "audit_plan_for_jave_repo",
		"audit_plan_report_id": 1,
		"limit":                100,
		"offset":               10}
	reslut, count, err := GetStorage().GetAuditPlanReportSQLsByReq(nameFields)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), count)
	assert.Len(t, reslut, 1)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
