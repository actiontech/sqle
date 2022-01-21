package auditplan

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/assert"
)

func assertManager(t *testing.T, m *Manager, apCount int) {
	assert.Len(t, manager.scheduler.cron.Entries(), apCount)
	assert.Len(t, manager.scheduler.entryIDs, apCount)
}

func getGivenBeforeAddAP() (*model.AuditPlan, *model.User, *model.Instance, string) {
	adminUser := &model.User{
		Model: model.Model{
			ID: 1,
		},
		Name: model.DefaultAdminUser,
	}
	ap := &model.AuditPlan{
		Model:          model.Model{ID: 1},
		Name:           "test_audit_plan",
		CronExpression: "*/1 * * * *",
		DBType:         driver.DriverTypeMySQL,
		CreateUserID:   adminUser.ID,
	}
	inst := &model.Instance{
		Name:   "test_inst",
		DbType: driver.DriverTypePostgreSQL,
	}

	token := "mock token"
	gomonkey.ApplyMethod(reflect.TypeOf(&utils.JWT{}), "CreateToken", func(_ *utils.JWT, _ string, _ int64, _ ...utils.CustomClaimOption) (string, error) {
		return token, nil
	})
	return ap, adminUser, inst, token
}

func TestManager_AddStaticAuditPlan(t *testing.T) {
	mockDB, mockHandle, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()
	model.InitMockStorage(mockDB)

	storage := model.GetStorage()

	mockHandle.ExpectQuery("SELECT * FROM `audit_plans`  WHERE `audit_plans`.`deleted_at` IS NULL").
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"name", "cron_expression"}))
	exitCh := InitManager(storage)
	defer func() {
		exitCh <- struct{}{}
	}()

	ap, adminUser, _, token := getGivenBeforeAddAP()
	manager := GetManager()

	mockHandle.ExpectQuery("SELECT * FROM `users`  WHERE `users`.`deleted_at` IS NULL AND ((login_name = ?)) ORDER BY `users`.`id` ASC LIMIT 1").
		WithArgs(adminUser.Name).
		WillReturnRows(mockHandle.NewRows([]string{"id", "login_name"}).AddRow(adminUser.ID, adminUser.Name))
	mockHandle.ExpectBegin()
	mockHandle.ExpectExec("INSERT INTO `audit_plans` (`created_at`,`updated_at`,`deleted_at`,`name`,`cron_expression`,`db_type`,`token`,`instance_name`,`create_user_id`,`instance_database`,`type`,`params`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)").
		WithArgs(model.MockTime, model.MockTime, nil, ap.Name, ap.CronExpression, ap.DBType, token, "", ap.CreateUserID, "", "", nil).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockHandle.ExpectCommit()
	err = manager.AddStaticAuditPlan(ap.Name, ap.CronExpression, ap.DBType, adminUser.Name, "", nil)
	assert.NoError(t, err)
	assert.Len(t, manager.scheduler.cron.Entries(), 1)
	assert.Len(t, manager.scheduler.entryIDs, 1)

	err = manager.AddStaticAuditPlan(ap.Name, "", "", "", "", nil)
	assert.Equal(t, ErrAuditPlanExisted.Error(), err.Error())
	assertManager(t, manager, 1)

	err = mockHandle.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestManager_AddDynamicAuditPlan(t *testing.T) {
	mockDB, mockHandle, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()
	model.InitMockStorage(mockDB)

	storage := model.GetStorage()

	mockHandle.ExpectQuery("SELECT * FROM `audit_plans`  WHERE `audit_plans`.`deleted_at` IS NULL").
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"name", "cron_expression"}))
	exitCh := InitManager(storage)
	defer func() {
		exitCh <- struct{}{}
	}()

	ap, adminUser, inst, token := getGivenBeforeAddAP()
	database := "test_db"
	manager := GetManager()

	mockHandle.ExpectQuery("SELECT * FROM `users`  WHERE `users`.`deleted_at` IS NULL AND ((login_name = ?)) ORDER BY `users`.`id` ASC LIMIT 1").
		WithArgs(adminUser.Name).
		WillReturnRows(mockHandle.NewRows([]string{"id", "login_name"}).AddRow(adminUser.ID, adminUser.Name))
	mockHandle.ExpectQuery("SELECT * FROM `instances`  WHERE `instances`.`deleted_at` IS NULL AND ((name = ?)) ORDER BY `instances`.`id` ASC LIMIT 1").
		WithArgs(inst.Name).
		WillReturnRows(sqlmock.NewRows([]string{"db_type"}).AddRow(inst.DbType))
	mockHandle.ExpectBegin()
	mockHandle.ExpectExec("INSERT INTO `audit_plans` (`created_at`,`updated_at`,`deleted_at`,`name`,`cron_expression`,`db_type`,`token`,`instance_name`,`create_user_id`,`instance_database`,`type`,`params`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)").
		WithArgs(model.MockTime, model.MockTime, nil, ap.Name, ap.CronExpression, inst.DbType, token, inst.Name, ap.CreateUserID, database, "", nil).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockHandle.ExpectCommit()
	err = manager.AddDynamicAuditPlan(ap.Name, ap.CronExpression, inst.Name, database, adminUser.Name, "", nil)
	assert.NoError(t, err)
	assertManager(t, manager, 1)
}

func TestManager_UpdateAuditPlan(t *testing.T) {
	mockDB, mockHandle, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()
	model.InitMockStorage(mockDB)

	storage := model.GetStorage()

	mockHandle.ExpectQuery("SELECT * FROM `audit_plans`  WHERE `audit_plans`.`deleted_at` IS NULL").
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"name", "cron_expression"}))
	exitCh := InitManager(storage)
	defer func() {
		exitCh <- struct{}{}
	}()

	ap, adminUser, _, token := getGivenBeforeAddAP()

	manager := GetManager()

	mockHandle.ExpectQuery("SELECT * FROM `users`  WHERE `users`.`deleted_at` IS NULL AND ((login_name = ?)) ORDER BY `users`.`id` ASC LIMIT 1").
		WithArgs(adminUser.Name).
		WillReturnRows(mockHandle.NewRows([]string{"id", "login_name"}).AddRow(adminUser.ID, adminUser.Name))
	mockHandle.ExpectBegin()
	mockHandle.ExpectExec("INSERT INTO `audit_plans` (`created_at`,`updated_at`,`deleted_at`,`name`,`cron_expression`,`db_type`,`token`,`instance_name`,`create_user_id`,`instance_database`,`type`,`params`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)").
		WithArgs(model.MockTime, model.MockTime, nil, ap.Name, ap.CronExpression, ap.DBType, token, "", ap.CreateUserID, "", "", nil).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockHandle.ExpectCommit()
	err = manager.AddStaticAuditPlan(ap.Name, ap.CronExpression, ap.DBType, adminUser.Name, "", nil)
	assert.NoError(t, err)
	assertManager(t, manager, 1)

	err = manager.UpdateAuditPlan("not_exist_ap_name", nil)
	assert.Equal(t, ErrAuditPlanNotExist.Error(), err.Error())

	updateAttr := map[string]interface{}{
		"cron_expression":   "*/2 * * * *",
		"instance_name":     "test_inst",
		"instance_database": "test_db",
	}

	mockHandle.ExpectBegin()
	mockHandle.ExpectExec("UPDATE `audit_plans` SET `cron_expression` = ?, `instance_database` = ?, `instance_name` = ?, `updated_at` = ? WHERE `audit_plans`.`deleted_at` IS NULL AND ((name = ?))").
		WithArgs(updateAttr["cron_expression"], updateAttr["instance_database"], updateAttr["instance_name"], model.MockTime, ap.Name).WillReturnResult(sqlmock.NewResult(1, 1))
	mockHandle.ExpectCommit()
	mockHandle.ExpectQuery("SELECT * FROM `audit_plans`  WHERE `audit_plans`.`deleted_at` IS NULL AND ((name = ?))").
		WithArgs(ap.Name).WillReturnRows(sqlmock.NewRows([]string{"name", "cron_expression"}).AddRow(ap.Name, updateAttr["cron_expression"]))
	err = manager.UpdateAuditPlan(ap.Name, updateAttr)
	assert.NoError(t, err)
	assertManager(t, manager, 1)
}

func TestManager_DeleteAuditPlan(t *testing.T) {
	mockDB, mockHandle, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()
	model.InitMockStorage(mockDB)

	storage := model.GetStorage()

	mockHandle.ExpectQuery("SELECT * FROM `audit_plans`  WHERE `audit_plans`.`deleted_at` IS NULL").
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"name", "cron_expression"}))
	exitCh := InitManager(storage)
	defer func() {
		exitCh <- struct{}{}
	}()

	ap, adminUser, _, token := getGivenBeforeAddAP()

	manager := GetManager()

	err = manager.DeleteAuditPlan("not_exist_ap_name")
	assert.Equal(t, ErrAuditPlanNotExist.Error(), err.Error())

	mockHandle.ExpectQuery("SELECT * FROM `users`  WHERE `users`.`deleted_at` IS NULL AND ((login_name = ?)) ORDER BY `users`.`id` ASC LIMIT 1").
		WithArgs(adminUser.Name).
		WillReturnRows(mockHandle.NewRows([]string{"id", "login_name"}).AddRow(adminUser.ID, adminUser.Name))
	mockHandle.ExpectBegin()
	mockHandle.ExpectExec("INSERT INTO `audit_plans` (`created_at`,`updated_at`,`deleted_at`,`name`,`cron_expression`,`db_type`,`token`,`instance_name`,`create_user_id`,`instance_database`,`type`,`params`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)").
		WithArgs(model.MockTime, model.MockTime, nil, ap.Name, ap.CronExpression, ap.DBType, token, "", ap.CreateUserID, "", "", nil).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockHandle.ExpectCommit()
	err = manager.AddStaticAuditPlan(ap.Name, ap.CronExpression, ap.DBType, adminUser.Name, "", nil)
	assert.NoError(t, err)
	assertManager(t, manager, 1)

	mockHandle.ExpectQuery("SELECT * FROM `audit_plans`  WHERE `audit_plans`.`deleted_at` IS NULL AND ((name = ?))").
		WithArgs(ap.Name).WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(ap.Name))
	mockHandle.ExpectBegin()
	mockHandle.ExpectExec("UPDATE `audit_plans` SET `deleted_at`=? WHERE `audit_plans`.`deleted_at` IS NULL").
		WithArgs(model.MockTime).WillReturnResult(sqlmock.NewResult(1, 1))
	mockHandle.ExpectCommit()
	err = manager.DeleteAuditPlan(ap.Name)
	assert.NoError(t, err)
	assertManager(t, manager, 0)
}

func TestInitManager(t *testing.T) {
	mockDB, mockHandle, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()
	model.InitMockStorage(mockDB)

	storage := model.GetStorage()

	mockHandle.ExpectQuery("SELECT * FROM `audit_plans`  WHERE `audit_plans`.`deleted_at` IS NULL").
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"name", "cron_expression"}).AddRow("test_ap_1", "*/1 * * * *").AddRow("test_ap_2", "*/2 * * * *"))
	InitManager(storage)

	manager := GetManager()
	assertManager(t, manager, 2)

	mockHandle.ExpectQuery("SELECT * FROM `audit_plans`  WHERE `audit_plans`.`deleted_at` IS NULL").
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"name", "cron_expression"}).AddRow("test_ap_1", "*/1 * * * *"))
	InitManager(storage)
	manager = GetManager()
	assertManager(t, manager, 1)
}

func TestManager_runJob(t *testing.T) {
	mockDB, mockHandle, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	model.InitMockStorage(mockDB)

	storage := model.GetStorage()

	mockHandle.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `audit_plans`")).
		WillReturnRows(sqlmock.NewRows([]string{"name", "cron_expression"}))
	exitCh := InitManager(storage)
	err = mockHandle.ExpectationsWereMet()
	assert.NoError(t, err)
	defer func() {
		exitCh <- struct{}{}
	}()

	// no SQL in audit plan, runJob should skip audit.
	m := GetManager()
	ap, _, _, _ := getGivenBeforeAddAP()
	mockHandle.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `audit_plans`")).
		WithArgs(ap.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(ap.ID))
	mockHandle.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `audit_plan_sqls`")).
		WithArgs(ap.ID).
		WillReturnRows(sqlmock.NewRows([]string{}))
	report := m.runJob(ap)
	err = mockHandle.ExpectationsWereMet()
	assert.NoError(t, err)
	assert.Nil(t, report)
}
