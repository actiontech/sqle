package model

import (
	"encoding/json"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestStorage_GetAuditPlansByReq(t *testing.T) {
	// 1. test for common user
	tableAndRowOfSQL := `
	FROM
	audit_plans 
WHERE
	audit_plans.deleted_at IS NULL 
	AND ( audit_plans.create_user_id = ? ) 
	AND audit_plans.db_type = ?
	`
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	InitMockStorage(mockDB)
	mock.ExpectPrepare(fmt.Sprintf(`SELECT audit_plans.name, audit_plans.cron_expression, audit_plans.db_type, audit_plans.token, audit_plans.instance_name, audit_plans.instance_database, audit_plans.rule_template_name, audit_plans.type, audit_plans.params %v LIMIT ? OFFSET ?`, tableAndRowOfSQL)).
		ExpectQuery().WithArgs(1, "mysql", 100, 10).WillReturnRows(sqlmock.NewRows([]string{"name", "cron_expression", "db_type", "token", "instance_name", "instance_database", "rule_template_name", "type", "params"}).
		AddRow("audit_plan_1", "* */2 * * *", "mysql", "fake token", "inst_1", "template_1", "db_1", "", nil))
	mock.ExpectPrepare(fmt.Sprintf(`SELECT COUNT(*) %v`, tableAndRowOfSQL)).
		ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{
		"COUNT(*)",
	}).AddRow("2"))
	nameFields := map[string]interface{}{
		"current_user_id":           1,
		"filter_audit_plan_db_type": "mysql",
		"limit":                     100,
		"offset":                    10}
	result, count, err := GetStorage().GetAuditPlansByReq(nameFields)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), count)
	assert.Len(t, result, 1)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
	mockDB.Close()

	// 2. test for admin user
	tableAndRowOfSQL1 := `
	FROM
	audit_plans 
WHERE
	audit_plans.deleted_at IS NULL 
	`
	mockDB, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	InitMockStorage(mockDB)
	mock.ExpectPrepare(fmt.Sprintf(`
	SELECT audit_plans.name, audit_plans.cron_expression, audit_plans.db_type, audit_plans.token, audit_plans.instance_name, audit_plans.instance_database, audit_plans.rule_template_name, audit_plans.type, audit_plans.params
	%v
	LIMIT ? OFFSET ?`, tableAndRowOfSQL1)).
		ExpectQuery().WithArgs(100, 10).WillReturnRows(sqlmock.NewRows([]string{
		"name", "cron_expression", "db_type", "token", "instance_name", "instance_database", "rule_template_name", "type", "params",
	}).AddRow("audit_plan_1", "* */2 * * *", "mysql", "fake token", "inst_1", "template_1", "db_1", "", nil))
	mock.ExpectPrepare(fmt.Sprintf(`SELECT COUNT(*) %v`, tableAndRowOfSQL1)).
		ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow("2"))
	nameFields = map[string]interface{}{
		"current_user_id":       1,
		"current_user_is_admin": true,
		"limit":                 100,
		"offset":                10}
	result, count, err = GetStorage().GetAuditPlansByReq(nameFields)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), count)
	assert.Len(t, result, 1)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
	mockDB.Close()
}

func TestStorage_GetAuditPlanSQLsByReq(t *testing.T) {
	tableAndRowOfSQL := `
	FROM audit_plan_sqls_v2 AS audit_plan_sqls
	JOIN audit_plans ON audit_plans.id = audit_plan_sqls.audit_plan_id
	WHERE audit_plan_sqls.deleted_at IS NULL
	AND  audit_plans.deleted_at IS NULL
	AND audit_plans.id = ?
	`
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()
	InitMockStorage(mockDB)
	mock.ExpectPrepare(fmt.Sprintf(`SELECT audit_plan_sqls.fingerprint, audit_plan_sqls.sql_content, audit_plan_sqls.schema, audit_plan_sqls.info %v order by audit_plan_sqls.id LIMIT ? OFFSET ?`, tableAndRowOfSQL)).
		ExpectQuery().
		WithArgs(1, 100, 10).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{"fingerprint", "sql_content", "schema", "info"}).
				AddRow(
					"select * from t1 where id = ?",
					"select * from t1 where id = 1",
					"schema",
					[]byte(`{"counter": 1, "last_receive_timestamp": "2021-09-01T13:46:13+08:00"}`),
				),
		)
	mock.ExpectPrepare(fmt.Sprintf(`SELECT COUNT(*) %v`, tableAndRowOfSQL)).
		ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow("2"))
	nameFields := map[string]interface{}{
		"audit_plan_id": 1,
		"limit":         100,
		"offset":        10}
	result, count, err := GetStorage().GetAuditPlanSQLsByReq(nameFields)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), count)
	assert.Len(t, result, 1)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestStorage_GetAuditPlanReportsByReq(t *testing.T) {
	tableAndRowOfSQL := `
	FROM
	audit_plan_reports_v2 AS reports
	JOIN audit_plans ON audit_plans.id = reports.audit_plan_id 
WHERE
	reports.deleted_at IS NULL 
	AND audit_plans.deleted_at IS NULL 
	AND audit_plans.name = ? 
	AND audit_plans.project_id = ? 
ORDER BY
	reports.created_at DESC ,
	reports.id DESC 
	`
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()
	InitMockStorage(mockDB)
	mock.ExpectPrepare(fmt.Sprintf(`SELECT reports.id, reports.score , reports.pass_rate, reports.audit_level, reports.created_at %v LIMIT ? OFFSET ?`, tableAndRowOfSQL)).
		ExpectQuery().WithArgs("audit_plan_for_jave_repo", "1", 100, 10).WillReturnRows(sqlmock.NewRows([]string{
		"id", "score", "pass_rate", "audit_level", "created_at"}).
		AddRow("1", 100, 1, "normal", "2021-09-01T13:46:13+08:00"))

	mock.ExpectPrepare(fmt.Sprintf(`SELECT COUNT(*) %v`, tableAndRowOfSQL)).
		ExpectQuery().WithArgs("audit_plan_for_jave_repo", "1").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow("2"))
	nameFields := map[string]interface{}{
		"audit_plan_name": "audit_plan_for_jave_repo",
		"project_id":      "1",
		"limit":           100,
		"offset":          10}
	result, count, err := GetStorage().GetAuditPlanReportsByReq(nameFields)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), count)
	assert.Len(t, result, 1)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestStorage_GetAuditPlanReportSQLsByReq(t *testing.T) {
	tableAndRowOfSQL := `
	FROM audit_plan_report_sqls_v2 AS report_sqls
	JOIN audit_plan_reports_v2 AS audit_plan_reports ON report_sqls.audit_plan_report_id = audit_plan_reports.id
	WHERE audit_plan_reports.deleted_at IS NULL
	AND report_sqls.deleted_at IS NULL
	AND report_sqls.audit_plan_report_id = ?
	AND audit_plan_reports.audit_plan_id = ?
	`

	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()
	InitMockStorage(mockDB)
	mockResult := []AuditResult{{Level: "error", Message: "FAKE AUDIT RESULT"}}
	mockResultBytes, _ := json.Marshal(mockResult)

	mock.ExpectPrepare(fmt.Sprintf(`SELECT report_sqls.sql, report_sqls.audit_results, report_sqls.number %v LIMIT ? OFFSET ?`, tableAndRowOfSQL)).
		ExpectQuery().WithArgs(1, 1, 100, 10).WillReturnRows(sqlmock.NewRows([]string{
		"sql", "audit_results", "number",
	}).AddRow("select * from t1 where id = 1", mockResultBytes, "1"))

	mock.ExpectPrepare(fmt.Sprintf(`SELECT COUNT(*) %v`, tableAndRowOfSQL)).
		ExpectQuery().WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow("2"))
	nameFields := map[string]interface{}{
		"audit_plan_id":        1,
		"audit_plan_report_id": 1,
		"limit":                100,
		"offset":               10}
	result, count, err := GetStorage().GetAuditPlanReportSQLsByReq(nameFields)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), count)
	assert.Len(t, result, 1)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
