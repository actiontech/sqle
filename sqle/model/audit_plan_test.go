package model

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestStorage_GetAuditPlans(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	initMockStorage(mockDB)
	mock.ExpectQuery("SELECT * FROM `audit_plans` WHERE `audit_plans`.`deleted_at` IS NULL").
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("audit_plan_1"))
	mock.ExpectClose()
	aps, err := GetStorage().GetAuditPlans()
	assert.NoError(t, err)
	assert.Len(t, aps, 1)
	mockDB.Close()
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestStorage_GetAuditPlanByName(t *testing.T) {
	// 1. test record exist
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	initMockStorage(mockDB)
	mock.ExpectQuery("SELECT * FROM `audit_plans` WHERE `audit_plans`.`deleted_at` IS NULL AND ((name = ?))").
		WithArgs("audit_plan_for_java_repo1").
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("audit_plan_1"))
	mock.ExpectClose()
	ap, exist, err := GetStorage().GetAuditPlanByName("audit_plan_for_java_repo1")
	assert.NoError(t, err)
	assert.True(t, exist)
	assert.Equal(t, "audit_plan_1", ap.Name)
	mockDB.Close()
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	// 2. test record not exist
	mockDB, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	initMockStorage(mockDB)
	mock.ExpectQuery("SELECT * FROM `audit_plans` WHERE `audit_plans`.`deleted_at` IS NULL AND ((name = ?))").
		WithArgs("audit_plan_for_java_repo1").
		WillReturnRows(sqlmock.NewRows([]string{"name"}))
	mock.ExpectClose()
	ap, exist, err = GetStorage().GetAuditPlanByName("audit_plan_for_java_repo1")
	assert.NoError(t, err)
	assert.False(t, exist)
	mockDB.Close()
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestStorage_GetAuditPlanSQLs(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	initMockStorage(mockDB)

	mockAuditPlanRow := AuditPlan{Model: Model{ID: 1}, Name: "audit_plan_for_java_repo1"}

	mock.ExpectQuery("SELECT * FROM `audit_plans` WHERE `audit_plans`.`deleted_at` IS NULL AND ((name = ?))").
		WithArgs("audit_plan_for_java_repo1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(mockAuditPlanRow.ID, mockAuditPlanRow.Name))
	mock.ExpectQuery("SELECT * FROM `audit_plan_sqls`  WHERE `audit_plan_sqls`.`deleted_at` IS NULL AND ((audit_plan_id = ?))").
		WithArgs(mockAuditPlanRow.ID).
		WillReturnRows(sqlmock.NewRows([]string{"fingerprint"}).AddRow("select * from t1 where id = ?").AddRow("select * from t2 where id = ?"))
	mock.ExpectClose()
	sqls, err := GetStorage().GetAuditPlanSQLs(mockAuditPlanRow.Name)
	assert.NoError(t, err)
	assert.Len(t, sqls, 2)
	mockDB.Close()
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	// 2. test update audit plan not exist
	mockDB, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	initMockStorage(mockDB)
	mock.ExpectQuery("SELECT * FROM `audit_plans` WHERE `audit_plans`.`deleted_at` IS NULL AND ((name = ?))").
		WithArgs("audit_plan_for_java_repo1").
		WillReturnRows(sqlmock.NewRows([]string{"name"}))
	mock.ExpectClose()
	_, err = GetStorage().GetAuditPlanSQLs("audit_plan_for_java_repo1")
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	mockDB.Close()
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// func TestStorage_UpdateAuditPlanByName(t *testing.T) {
// 	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
// 	assert.NoError(t, err)
// 	initMockStorage(mockDB)
// 	mockAuditPlanRow := AuditPlan{
// 		Model:            Model{ID: 1},
// 		Name:             "audit_plan_for_java_repo1",
// 		CronExpression:   "* * * * *",
// 		InstanceName:     "inst1",
// 		InstanceDatabase: "db_1"}

// 	updateTime := time.Now()
// 	mock.ExpectExec("UPDATE `audit_plans` SET `cron_expression` = ?, `instance_database` = ?, `instance_name` = ?, `updated_at` = ? WHERE `audit_plans`.`deleted_at` IS NULL AND ((name = ?))").
// 		WithArgs("* */2 * * *", "db_2", "inst2", updateTime, "audit_plan_for_java_repo1").
// 		WillReturnResult(sqlmock.NewResult(1, 1))
// 	mock.ExpectClose()
// 	updateAttrs := map[string]interface{}{
// 		"cron_expression":   "* */2 * * *",
// 		"updated_at":        updateTime,
// 		"instance_name":     "inst2",
// 		"instance_database": "db_2"}
// 	err = GetStorage().UpdateAuditPlanByName(mockAuditPlanRow.Name, updateAttrs)
// 	assert.NoError(t, err)
// 	mockDB.Close()
// 	err = mock.ExpectationsWereMet()
// 	assert.NoError(t, err)
// }
