package auditplan

import (
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"github.com/actiontech/sqle/sqle/log"

	"github.com/actiontech/sqle/sqle/model"
	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	mockDB, mockHandle, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()
	model.InitMockStorage(mockDB)
	storage := model.GetStorage()

	j := NewManager(log.NewEntry())
	m, _ := j.(*Manager)
	m.persist = storage

	nextTime := m.lastSyncTime.Add(5 * time.Second)

	// test init
	mockHandle.ExpectQuery("SELECT * FROM `audit_plans` WHERE `audit_plans`.`deleted_at` IS NULL AND ((project_status = ?))").
		WithArgs("active").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "cron_expression"}).AddRow(1, "test_ap_1", "default", "*/1 * * * *"))

	mockHandle.ExpectQuery("SELECT id, updated_at FROM `audit_plans` WHERE (updated_at > ?) ORDER BY updated_at").
		WithArgs(m.lastSyncTime).
		WillReturnRows(sqlmock.NewRows([]string{"id", "updated_at"}).AddRow(2, nextTime))

	mockHandle.ExpectQuery("SELECT * FROM `audit_plans` WHERE `audit_plans`.`deleted_at` IS NULL AND ((project_status = ?) AND (id = ?)) ").
		WithArgs("active", 2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "cron_expression"}).AddRow(2, "test_ap_2", "default", "*/2 * * * *"))

	assert.NoError(t, m.sync())

	assert.Len(t, m.scheduler.cron.Entries(), 2)
	assert.Len(t, m.tasks, 2)
	task, err := m.getTask(1)
	assert.NoError(t, err)
	dt, ok := task.(*DefaultTask)
	assert.Equal(t, ok, true)
	assert.Equal(t, dt.ap.CronExpression, "*/1 * * * *")

	task, err = m.getTask(2)
	assert.NoError(t, err)
	dt, ok = task.(*DefaultTask)
	assert.Equal(t, ok, true)
	assert.Equal(t, dt.ap.CronExpression, "*/2 * * * *")

	assert.Equal(t, &nextTime, m.lastSyncTime)
	assert.Equal(t, true, m.isFullSyncDone)
	// return

	// test add task

	nextTimeMore := nextTime.Add(6 * time.Second)
	mockHandle.ExpectQuery("SELECT id, updated_at FROM `audit_plans` WHERE (updated_at > ?) ORDER BY updated_at").
		WithArgs(nextTime).
		WillReturnRows(sqlmock.NewRows([]string{"id", "updated_at"}).AddRow(3, nextTimeMore))

	mockHandle.ExpectQuery("SELECT * FROM `audit_plans` WHERE `audit_plans`.`deleted_at` IS NULL AND ((project_status = ?) AND (id = ?))").
		WithArgs("active", 3).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "cron_expression"}).AddRow(3, "test_ap_3", "default", "*/3 * * * *"))

	m.sync()
	assert.Len(t, m.scheduler.cron.Entries(), 3)
	assert.Len(t, m.tasks, 3)
	task, err = m.getTask(3)
	assert.NoError(t, err)
	dt, ok = task.(*DefaultTask)
	assert.Equal(t, ok, true)
	assert.Equal(t, dt.ap.CronExpression, "*/3 * * * *")

	assert.Equal(t, &nextTimeMore, m.lastSyncTime)
	assert.Equal(t, true, m.isFullSyncDone)

	// test delete task
	nextTimeMore2 := nextTimeMore.Add(6 * time.Second)
	mockHandle.ExpectQuery("SELECT id, updated_at FROM `audit_plans` WHERE (updated_at > ?) ORDER BY updated_at").
		WithArgs(nextTimeMore).
		WillReturnRows(sqlmock.NewRows([]string{"id", "updated_at"}).AddRow(3, nextTimeMore2))

	mockHandle.ExpectQuery("SELECT * FROM `audit_plans` WHERE `audit_plans`.`deleted_at` IS NULL AND ((project_status = ?) AND (id = ?))").
		WithArgs("active", 3).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "cron_expression"}))

	m.sync()
	assert.Len(t, m.scheduler.cron.Entries(), 2)
	assert.Len(t, m.tasks, 2)
	_, err = m.getTask(3)
	assert.Error(t, err)

	assert.Equal(t, &nextTimeMore2, m.lastSyncTime)
	assert.Equal(t, true, m.isFullSyncDone)

	// test no task change
	mockHandle.ExpectQuery("SELECT id, updated_at FROM `audit_plans` WHERE (updated_at > ?) ORDER BY updated_at").
		WithArgs(nextTimeMore2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "updated_at"}))

	m.sync()
	assert.Len(t, m.scheduler.cron.Entries(), 2)
	assert.Len(t, m.tasks, 2)

	assert.Equal(t, &nextTimeMore2, m.lastSyncTime) // next time do not change.
}
