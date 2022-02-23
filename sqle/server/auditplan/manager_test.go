package auditplan

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	mockDB, mockHandle, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()
	model.InitMockStorage(mockDB)
	storage := model.GetStorage()

	// test init
	mockHandle.ExpectQuery("SELECT * FROM `audit_plans`  WHERE `audit_plans`.`deleted_at` IS NULL").
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"name", "type", "cron_expression"}).AddRow("test_ap_1", "default", "*/1 * * * *"))

	InitManager(storage)

	assert.Len(t, manager.scheduler.cron.Entries(), 1)
	assert.Len(t, manager.tasks, 1)
	task, err := manager.getTask("test_ap_1")
	assert.NoError(t, err)
	dt, ok := task.(*DefaultTask)
	assert.Equal(t, ok, true)
	assert.Equal(t, dt.ap.CronExpression, "*/1 * * * *")

	// test add task
	mockHandle.ExpectQuery("SELECT * FROM `audit_plans`  WHERE `audit_plans`.`deleted_at` IS NULL AND ((name = ?))").
		WithArgs("test_ap_2").
		WillReturnRows(sqlmock.NewRows([]string{"name", "type", "cron_expression"}).AddRow("test_ap_2", "default", "*/1 * * * *"))

	manager.SyncTask("test_ap_2")
	assert.Len(t, manager.scheduler.cron.Entries(), 2)
	assert.Len(t, manager.tasks, 2)
	task, err = manager.getTask("test_ap_2")
	assert.NoError(t, err)
	dt, ok = task.(*DefaultTask)
	assert.Equal(t, ok, true)
	assert.Equal(t, dt.ap.CronExpression, "*/1 * * * *")

	// test delete task
	mockHandle.ExpectQuery("SELECT * FROM `audit_plans`  WHERE `audit_plans`.`deleted_at` IS NULL AND ((name = ?))").
		WithArgs("test_ap_2").
		WillReturnRows(sqlmock.NewRows([]string{"name", "type", "cron_expression"}))

	manager.SyncTask("test_ap_2")
	assert.Len(t, manager.scheduler.cron.Entries(), 1)
	assert.Len(t, manager.tasks, 1)
	_, err = manager.getTask("test_ap_1")
	assert.NoError(t, err)

	// test update task
	mockHandle.ExpectQuery("SELECT * FROM `audit_plans`  WHERE `audit_plans`.`deleted_at` IS NULL AND ((name = ?))").
		WithArgs("test_ap_1").
		WillReturnRows(sqlmock.NewRows([]string{"name", "type", "cron_expression"}).AddRow("test_ap_1", "default", "*/2 * * * *"))

	manager.SyncTask("test_ap_1")
	assert.Len(t, manager.scheduler.cron.Entries(), 1)
	assert.Len(t, manager.tasks, 1)
	task, err = manager.getTask("test_ap_1")
	assert.NoError(t, err)
	dt, ok = task.(*DefaultTask)
	assert.Equal(t, ok, true)
	assert.Equal(t, dt.ap.CronExpression, "*/2 * * * *")
}
