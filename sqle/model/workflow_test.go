package model

import (
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestStorage_UpdateWorkflowSchedule(t *testing.T) {
	// test set schedule time
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	InitMockStorage(mockDB)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `workflow_records` SET `schedule_user_id` = ?, `scheduled_at` = ?, `updated_at` = ? WHERE `workflow_records`.`deleted_at` IS NULL AND ((id = ?))").
		WithArgs(1, AnyTime{}, AnyTime{}, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	record := &WorkflowRecord{
		Model: Model{ID: 1},
	}
	workflow := &Workflow{Record: record}
	scheduleTime := time.Date(2021, 12, 1, 12, 00, 00, 00, time.Local)
	err = GetStorage().UpdateWorkflowSchedule(workflow, 1, &scheduleTime)
	assert.NoError(t, err)
	mockDB.Close()
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	// test del schedule time, set to null
	mockDB, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	InitMockStorage(mockDB)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `workflow_records` SET `schedule_user_id` = ?, `scheduled_at` = ?, `updated_at` = ? WHERE `workflow_records`.`deleted_at` IS NULL AND ((id = ?))").
		WithArgs(2, nil, AnyTime{}, 2).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	record = &WorkflowRecord{
		Model: Model{ID: 2},
	}
	workflow = &Workflow{Record: record}
	err = GetStorage().UpdateWorkflowSchedule(workflow, 2, nil)
	if err != nil {
		fmt.Printf("err: [%v]", err)
	}
	assert.NoError(t, err)
	mockDB.Close()
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
