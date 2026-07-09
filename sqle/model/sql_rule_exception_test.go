package model

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestStorage_CreateSQLRuleExceptionIfNotExist(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	mock.ExpectQuery("SELECT VERSION\\(\\)").WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("5.7"))
	InitMockStorage(mockDB)

	sqlRuleException := &SQLRuleException{
		ProjectId:      "700300",
		InstanceID:     1,
		SQLFingerprint: "select * from ?",
		RuleName:       "all_check_prepare_statement_placeholders",
		Reason:         "业务确认该 SQL 可例外",
	}
	mock.ExpectQuery("SELECT \\* FROM `sql_rule_exception` WHERE \\(project_id = \\? AND instance_id = \\? AND sql_fingerprint = \\? AND rule_name = \\?\\) AND `sql_rule_exception`.`deleted_at` IS NULL ORDER BY `sql_rule_exception`.`id` LIMIT 1").
		WithArgs("700300", uint64(1), "select * from ?", "all_check_prepare_statement_placeholders").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `sql_rule_exception`").WillReturnResult(sqlmock.NewResult(11, 1))
	mock.ExpectCommit()

	savedSQLRuleException, created, err := GetStorage().CreateSQLRuleExceptionIfNotExist(sqlRuleException)
	assert.NoError(t, err)
	if !assert.NotNil(t, savedSQLRuleException) {
		return
	}
	assert.True(t, created)
	assert.Equal(t, uint(11), savedSQLRuleException.ID)
	assert.NotEmpty(t, savedSQLRuleException.UniqueKey)

	mock.ExpectClose()
	assert.NoError(t, mockDB.Close())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStorage_CreateSQLRuleExceptionIfNotExistReturnsDataInvalidWhenFingerprintMissing(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	mock.ExpectQuery("SELECT VERSION\\(\\)").WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("5.7"))
	InitMockStorage(mockDB)

	savedSQLRuleException, created, err := GetStorage().CreateSQLRuleExceptionIfNotExist(&SQLRuleException{
		ProjectId:      "700300",
		InstanceID:     1,
		SQLFingerprint: "  ",
		RuleName:       "all_check_prepare_statement_placeholders",
		Reason:         "业务确认该 SQL 可例外",
	})

	assert.Nil(t, savedSQLRuleException)
	assert.False(t, created)
	assert.EqualError(t, err, SQLRuleExceptionMissingFingerprintMessage)
	codeErr, ok := err.(interface{ Code() int })
	if assert.True(t, ok) {
		assert.Equal(t, int(errors.DataInvalid), codeErr.Code())
	}

	mock.ExpectClose()
	assert.NoError(t, mockDB.Close())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStorage_CreateSQLRuleExceptionIfNotExistReturnsDataExist(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	mock.ExpectQuery("SELECT VERSION\\(\\)").WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("5.7"))
	InitMockStorage(mockDB)

	mock.ExpectQuery("SELECT \\* FROM `sql_rule_exception` WHERE \\(project_id = \\? AND instance_id = \\? AND sql_fingerprint = \\? AND rule_name = \\?\\) AND `sql_rule_exception`.`deleted_at` IS NULL ORDER BY `sql_rule_exception`.`id` LIMIT 1").
		WithArgs("700300", uint64(1), "select * from ?", "all_check_prepare_statement_placeholders").
		WillReturnRows(sqlmock.NewRows([]string{"id", "project_id", "instance_id", "sql_fingerprint", "rule_name", "reason"}).
			AddRow(11, "700300", 1, "select * from ?", "all_check_prepare_statement_placeholders", "业务确认该 SQL 可例外"))

	savedSQLRuleException, created, err := GetStorage().CreateSQLRuleExceptionIfNotExist(&SQLRuleException{
		ProjectId:      "700300",
		InstanceID:     1,
		SQLFingerprint: "select * from ?",
		RuleName:       "all_check_prepare_statement_placeholders",
		Reason:         "业务确认该 SQL 可例外",
	})
	assert.Error(t, err)
	if !assert.NotNil(t, savedSQLRuleException) {
		return
	}
	assert.False(t, created)
	assert.Equal(t, uint(11), savedSQLRuleException.ID)

	mock.ExpectClose()
	assert.NoError(t, mockDB.Close())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStorage_GetEffectiveSQLRuleExceptionsMatchesByRuleNameAndKeepsSnapshot(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	mock.ExpectQuery("SELECT VERSION\\(\\)").WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("5.7"))
	InitMockStorage(mockDB)

	mock.ExpectQuery("SELECT \\* FROM `sql_rule_exception` WHERE \\(project_id = \\? AND sql_fingerprint IN \\(\\?\\) AND rule_name IN \\(\\?\\)\\) AND `sql_rule_exception`.`deleted_at` IS NULL").
		WithArgs(ProjectUID("700300"), "SELECT * FROM ... WHERE ... IN (...)", "all_check_where_is_invalid").
		WillReturnRows(sqlmock.NewRows([]string{"id", "project_id", "instance_id", "sql_fingerprint", "rule_name", "rule_desc", "rule_level", "reason"}).
			AddRow(11, "700300", 2067529851245432800, "SELECT * FROM ... WHERE ... IN (...)", "all_check_where_is_invalid", "添加时规则描述快照", "error", "业务确认该 SQL 可例外"))

	sqlRuleExceptions, err := GetStorage().GetEffectiveSQLRuleExceptions("700300", 2067529851245432832, "SELECT * FROM ... WHERE ... IN (...)", []string{"all_check_where_is_invalid"})
	assert.NoError(t, err)
	if assert.Contains(t, sqlRuleExceptions, "all_check_where_is_invalid") {
		assert.Equal(t, "添加时规则描述快照", sqlRuleExceptions["all_check_where_is_invalid"].RuleDesc)
		assert.Equal(t, "error", sqlRuleExceptions["all_check_where_is_invalid"].RuleLevel)
	}

	mock.ExpectClose()
	assert.NoError(t, mockDB.Close())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStorage_GetSQLRuleExceptionsWithFilters(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	mock.ExpectQuery("SELECT VERSION\\(\\)").WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("5.7"))
	InitMockStorage(mockDB)

	instanceID := uint64(2067529851245432832)
	ruleName := "rule_a"
	createdBy := "admin"
	createdTimeFrom := "2026-06-18 00:00:00"
	createdTimeTo := "2026-06-19 00:00:00"
	sqlFingerprint := "select *"
	fuzzySearchValue := "business"
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `sql_rule_exception` LEFT JOIN instances ON instances.id = sql_rule_exception.instance_id WHERE sql_rule_exception.project_id = \\? AND sql_rule_exception.deleted_at IS NULL AND sql_rule_exception.instance_id = \\? AND sql_rule_exception.rule_name LIKE \\? AND sql_rule_exception.created_by = \\? AND sql_rule_exception.created_at > \\? AND sql_rule_exception.created_at < \\? AND sql_rule_exception.sql_fingerprint LIKE \\? AND \\(sql_rule_exception.sql_fingerprint LIKE \\? OR sql_rule_exception.rule_name LIKE \\? OR sql_rule_exception.rule_desc LIKE \\? OR sql_rule_exception.reason LIKE \\? OR sql_rule_exception.created_by LIKE \\?\\)").
		WithArgs("700300", instanceID, "%rule_a%", "admin", createdTimeFrom, createdTimeTo, "%select *%", "%business%", "%business%", "%business%", "%business%", "%business%").
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT sql_rule_exception.\\*, instances.name AS instance_name FROM `sql_rule_exception` LEFT JOIN instances ON instances.id = sql_rule_exception.instance_id WHERE sql_rule_exception.project_id = \\? AND sql_rule_exception.deleted_at IS NULL AND sql_rule_exception.instance_id = \\? AND sql_rule_exception.rule_name LIKE \\? AND sql_rule_exception.created_by = \\? AND sql_rule_exception.created_at > \\? AND sql_rule_exception.created_at < \\? AND sql_rule_exception.sql_fingerprint LIKE \\? AND \\(sql_rule_exception.sql_fingerprint LIKE \\? OR sql_rule_exception.rule_name LIKE \\? OR sql_rule_exception.rule_desc LIKE \\? OR sql_rule_exception.reason LIKE \\? OR sql_rule_exception.created_by LIKE \\?\\) AND `sql_rule_exception`.`deleted_at` IS NULL ORDER BY sql_rule_exception.id desc LIMIT 20").
		WithArgs("700300", instanceID, "%rule_a%", "admin", createdTimeFrom, createdTimeTo, "%select *%", "%business%", "%business%", "%business%", "%business%", "%business%").
		WillReturnRows(sqlmock.NewRows([]string{"id", "project_id", "instance_id", "sql_fingerprint", "rule_name", "reason", "created_by", "instance_name"}).
			AddRow(11, "700300", instanceID, "select * from ?", "rule_a", "business exception", "admin", "mysql_local_sqle"))

	sqlRuleExceptions, count, err := GetStorage().GetSQLRuleExceptions(1, 20, SQLRuleExceptionListFilter{
		ProjectID:        "700300",
		InstanceID:       &instanceID,
		RuleName:         &ruleName,
		CreatedBy:        &createdBy,
		CreatedTimeFrom:  &createdTimeFrom,
		CreatedTimeTo:    &createdTimeTo,
		SQLFingerprint:   &sqlFingerprint,
		FuzzySearchValue: &fuzzySearchValue,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
	if assert.Len(t, sqlRuleExceptions, 1) {
		assert.Equal(t, "mysql_local_sqle", sqlRuleExceptions[0].InstanceName)
	}

	mock.ExpectClose()
	assert.NoError(t, mockDB.Close())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStorage_DeleteSQLRuleExceptionReleasesUniqueKey(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	mock.ExpectQuery("SELECT VERSION\\(\\)").WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("5.7"))
	InitMockStorage(mockDB)

	sqlRuleException := &SQLRuleException{
		Model:          Model{ID: 11},
		ProjectId:      "700300",
		InstanceID:     1,
		SQLFingerprint: "select * from ?",
		RuleName:       "all_check_prepare_statement_placeholders",
		Reason:         "业务确认该 SQL 可例外",
		UniqueKey:      BuildSQLRuleExceptionUniqueKey("700300", 1, "select * from ?", "all_check_prepare_statement_placeholders"),
	}
	deletedUniqueKey := BuildDeletedSQLRuleExceptionUniqueKey(sqlRuleException.ID, sqlRuleException.UniqueKey)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE `+"`"+`sql_rule_exception`+"`"+` SET `+"`"+`unique_key`+"`"+`=\? WHERE id = \? AND deleted_at IS NULL`).
		WithArgs(deletedUniqueKey, sqlRuleException.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("UPDATE `sql_rule_exception` SET `deleted_at`=\\? WHERE `sql_rule_exception`.`id` = \\? AND `sql_rule_exception`.`deleted_at` IS NULL").
		WithArgs(sqlmock.AnyArg(), sqlRuleException.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	assert.NoError(t, GetStorage().DeleteSQLRuleException(sqlRuleException))

	mock.ExpectClose()
	assert.NoError(t, mockDB.Close())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStorage_GetSQLRuleExceptionByIDAndProjectIDNotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	mock.ExpectQuery("SELECT VERSION\\(\\)").WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("5.7"))
	InitMockStorage(mockDB)

	mock.ExpectQuery(`SELECT \* FROM `+"`"+`sql_rule_exception`+"`"+` WHERE \(id = \? AND project_id = \?\) AND `+"`"+`sql_rule_exception`+"`"+`.`+"`"+`deleted_at`+"`"+` IS NULL ORDER BY `+"`"+`sql_rule_exception`+"`"+`.`+"`"+`id`+"`"+` LIMIT 1`).
		WithArgs("11", "700300").
		WillReturnError(gorm.ErrRecordNotFound)

	sqlRuleException, exist, err := GetStorage().GetSQLRuleExceptionByIDAndProjectID("11", "700300")
	assert.NoError(t, err)
	assert.False(t, exist)
	assert.Nil(t, sqlRuleException)

	mock.ExpectClose()
	assert.NoError(t, mockDB.Close())
	assert.NoError(t, mock.ExpectationsWereMet())
}
