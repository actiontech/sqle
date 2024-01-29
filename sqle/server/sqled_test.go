package server

import (
	"context"
	_driver "database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/actiontech/sqle/sqle/driver"
	_ "github.com/actiontech/sqle/sqle/driver/mysql"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/assert"
)

func getAction(sqls []string, typ int, p driver.Plugin) *action {
	task := &model.Task{
		Model:      model.Model{ID: 1},
		SQLSource:  model.TaskSQLSourceFromMyBatisXMLFile,
		InstanceId: 1,
	}

	for _, sql := range sqls {
		task.ExecuteSQLs = append(task.ExecuteSQLs, &model.ExecuteSQL{
			BaseSQL: model.BaseSQL{Content: sql},
		})
	}

	entry := log.NewEntry().WithField("task_id", task.ID)
	return &action{
		task:   task,
		plugin: p,
		typ:    typ,
		entry:  entry,
		done:   make(chan struct{}),
	}
}

type mockDriver struct {
	parseError bool
}

func (d *mockDriver) Close(ctx context.Context) {}

func (d *mockDriver) Ping(ctx context.Context) error {
	return nil
}

func (d *mockDriver) KillProcess(ctx context.Context) error {
	return nil
}

func (d *mockDriver) Exec(ctx context.Context, query string) (_driver.Result, error) {
	return nil, nil
}

func (d *mockDriver) Tx(ctx context.Context, queries ...string) ([]_driver.Result, error) {
	return nil, nil
}

func (d *mockDriver) Schemas(ctx context.Context) ([]string, error) {
	return nil, nil
}

func (d *mockDriver) Parse(ctx context.Context, sqlText string) ([]driverV2.Node, error) {
	if d.parseError {
		return nil, errors.New("mock error: mockDriver.Parse")
	}

	return []driverV2.Node{{Text: sqlText}}, nil
}

func (d *mockDriver) Audit(ctx context.Context, sqls []string) ([]*driverV2.AuditResults, error) {
	return nil, nil
}

func (d *mockDriver) GenRollbackSQL(ctx context.Context, sql string) (string, string, error) {
	return "", "", nil
}

func (d *mockDriver) Explain(ctx context.Context, conf *driverV2.ExplainConf) (*driverV2.ExplainResult, error) {
	return nil, nil
}

func (d *mockDriver) GetTableMetaBySQL(ctx context.Context, conf *driver.GetTableMetaBySQLConf) (*driver.GetTableMetaBySQLResult, error) {
	return nil, nil
}

func (d *mockDriver) Query(ctx context.Context, sql string, conf *driverV2.QueryConf) (*driverV2.QueryResult, error) {
	return nil, nil
}

func (d *mockDriver) EstimateSQLAffectRows(ctx context.Context, sql string) (*driverV2.EstimatedAffectRows, error) {
	return nil, nil
}

func TestAction_validation(t *testing.T) {
	actions := map[int]*action{
		ActionTypeAudit:    {typ: ActionTypeAudit},
		ActionTypeExecute:  {typ: ActionTypeExecute},
		ActionTypeRollback: {typ: ActionTypeRollback},
	}

	auditingTask := &model.Task{
		ExecuteSQLs: []*model.ExecuteSQL{
			{BaseSQL: model.BaseSQL{ExecStatus: model.SQLExecuteStatusInitialized}, AuditStatus: model.SQLAuditStatusDoing},
			{BaseSQL: model.BaseSQL{ExecStatus: model.SQLExecuteStatusInitialized}, AuditStatus: model.SQLAuditStatusInitialized},
		},
	}
	assert.Nil(t, actions[ActionTypeAudit].validation(auditingTask))

	executingTask := &model.Task{
		ExecuteSQLs: []*model.ExecuteSQL{
			{BaseSQL: model.BaseSQL{ExecStatus: model.SQLExecuteStatusInitialized}, AuditStatus: model.SQLAuditStatusFinished},
			{BaseSQL: model.BaseSQL{ExecStatus: model.SQLExecuteStatusDoing}, AuditStatus: model.SQLAuditStatusFinished},
		},
	}
	assert.EqualError(t, actions[ActionTypeExecute].validation(executingTask), ErrActionExecuteOnExecutedTask.Error())

	noAuditedTask := &model.Task{
		ExecuteSQLs: []*model.ExecuteSQL{
			{BaseSQL: model.BaseSQL{ExecStatus: model.SQLExecuteStatusInitialized}, AuditStatus: model.SQLAuditStatusInitialized},
			{BaseSQL: model.BaseSQL{ExecStatus: model.SQLExecuteStatusInitialized}, AuditStatus: model.SQLAuditStatusInitialized},
		},
	}
	assert.EqualError(t, actions[ActionTypeExecute].validation(noAuditedTask), ErrActionExecuteOnNonAuditedTask.Error())

	rollbackingTask := &model.Task{RollbackSQLs: []*model.RollbackSQL{
		{BaseSQL: model.BaseSQL{ExecStatus: model.SQLExecuteStatusDoing}},
		{BaseSQL: model.BaseSQL{ExecStatus: model.SQLExecuteStatusInitialized}},
	}}
	assert.EqualError(t, actions[ActionTypeRollback].validation(rollbackingTask), ErrActionRollbackOnRollbackedTask.Error())

	executedFailTask := &model.Task{
		ExecuteSQLs: []*model.ExecuteSQL{
			{BaseSQL: model.BaseSQL{ExecStatus: model.SQLExecuteStatusSucceeded}, AuditStatus: model.SQLAuditStatusFinished},
			{BaseSQL: model.BaseSQL{ExecStatus: model.SQLExecuteStatusFailed}, AuditStatus: model.SQLAuditStatusFinished},
		},
	}
	assert.EqualError(t, actions[ActionTypeRollback].validation(executedFailTask), ErrActionRollbackOnExecuteFailedTask.Error())

	noExecutedTask := &model.Task{ExecuteSQLs: []*model.ExecuteSQL{
		{BaseSQL: model.BaseSQL{ExecStatus: model.SQLExecuteStatusInitialized}, AuditStatus: model.SQLAuditStatusFinished},
		{BaseSQL: model.BaseSQL{ExecStatus: model.SQLExecuteStatusInitialized}, AuditStatus: model.SQLAuditStatusFinished},
	}}
	assert.EqualError(t, actions[ActionTypeRollback].validation(noExecutedTask), ErrActionRollbackOnNonExecutedTask.Error())
}

func Test_action_audit_UpdateTask(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	model.InitMockStorage(mockDB)

	whitelist := model.SqlWhitelist{
		Value:     "select * from t1",
		MatchType: model.SQLWhitelistExactMatch,
	}
	act := getAction([]string{"select * from t1"}, ActionTypeAudit, &mockDriver{})

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `sql_whitelist` WHERE `sql_whitelist`.`deleted_at` IS NULL AND ((sql_whitelist.project_id = ?))")).
		WithArgs("").
		WillReturnRows(sqlmock.NewRows([]string{"value", "match_type"}).AddRow(whitelist.Value, whitelist.MatchType))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `execute_sql_detail`")).
		WithArgs(model.MockTime, model.MockTime, nil, 0, 0, act.task.ExecuteSQLs[0].Content, "", "", 0, "", 0, 0, "", "", "", 0, "", model.SQLAuditStatusFinished, `[{"level":"normal","message":"白名单","rule_name":""}]`, "2882fdbb7d5bcda7b49ea0803493467e", "normal").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `tasks`")).
		WithArgs(driverV2.RuleLevelNormal, float64(1), 100, model.TaskStatusAudited, act.task.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = act.audit()
	assert.NoError(t, err)
	assert.Equal(t, model.TaskStatusAudited, act.task.Status)
	assert.Equal(t, float64(1), act.task.PassRate)
}

func Test_action_execute(t *testing.T) {
	driver.GetPluginManager().Start("", nil)
	mockUpdateTaskStatus := func(t *testing.T) {
		gomonkey.ApplyMethod(reflect.TypeOf(&model.Storage{}), "UpdateTask", func(_ *model.Storage, _ *model.Task, attr ...interface{}) error {
			a, ok := attr[0].(map[string]interface{})
			if !ok {
				assert.Error(t, fmt.Errorf("updateTask args type expect is map[string]interface{}"))
				return nil
			}
			status, ok := a["status"].(string)
			if !ok {
				assert.Error(t, fmt.Errorf("updateTask args attr[\"status\"] type expect is string"))
				return nil
			}
			if status == model.TaskStatusExecuting {
				return nil
			}

			assert.Equal(t, model.TaskStatusExecuteFailed, status)
			return nil
		})

		gomonkey.ApplyMethod(reflect.TypeOf(&model.Storage{}), "GetRulesFromRuleTemplateByName", func(_ *model.Storage, _ []string, _ string) ([]*model.Rule, []*model.CustomRule, error) {
			return nil, nil, nil
		})
	}

	newDriver := func() (driver.Plugin, error) {
		rules, _, err := model.GetStorage().GetAllRulesByTmpNameAndProjectIdInstanceDBType("", "", nil, driverV2.DriverTypeMySQL)
		if err != nil {
			return nil, err
		}
		p, err := newDriverManagerWithAudit(log.NewEntry(), nil, "", driverV2.DriverTypeMySQL, rules)
		if err != nil {
			return nil, err
		}
		defer p.Close(context.TODO())

		// d, err := drvMgr.GetAuditDriver()
		// if err != nil {
		// 	return nil, err
		// }
		return p, nil
	}

	tests := []struct {
		name    string
		setUp   func(t *testing.T) (driver.Plugin, error)
		sqls    []string
		wantErr bool
	}{
		{
			name: "Given: one SQL;Parse error, then Update Task Status to Failed",
			setUp: func(t *testing.T) (driver.Plugin, error) {
				mockUpdateTaskStatus(t)
				return &mockDriver{parseError: true}, nil
			},
			sqls:    []string{"select * from t1"},
			wantErr: false,
		},

		{
			name: "Given: one SQL;execSQLs error, then Update Task Status to Failed",
			setUp: func(t *testing.T) (driver.Plugin, error) {
				mockUpdateTaskStatus(t)

				gomonkey.ApplyMethod(reflect.TypeOf(&model.Storage{}), "UpdateExecuteSQLs", func(_ *model.Storage, _ []*model.ExecuteSQL) error {
					return errors.New("mock error: Storage.UpdateExecuteSQLs")
				})

				return newDriver()
			},
			sqls:    []string{"select * from t1"},
			wantErr: false,
		},

		{
			name: "Given: one SQL;execSQL error, then Update Task Status to Failed",
			setUp: func(t *testing.T) (driver.Plugin, error) {
				mockUpdateTaskStatus(t)

				gomonkey.ApplyMethod(reflect.TypeOf(&model.Storage{}), "UpdateExecuteSqlStatus", func(_ *model.Storage, _ *model.BaseSQL, _ string, _ string) error {
					return errors.New("mock error: Storage.UpdateExecuteSqlStatus")
				})

				return newDriver()
			},
			sqls:    []string{"create table t1(id int)"},
			wantErr: false,
		},

		{
			name: "Given: two SQLs;execSQLs error, then Update Task Status to Failed",
			setUp: func(t *testing.T) (driver.Plugin, error) {
				mockUpdateTaskStatus(t)

				gomonkey.ApplyMethod(reflect.TypeOf(&model.Storage{}), "UpdateExecuteSQLs", func(_ *model.Storage, _ []*model.ExecuteSQL) error {
					return errors.New("mock error: Storage.UpdateExecuteSQLs")
				})

				return newDriver()
			},
			sqls:    []string{"select * from t1", "create table t1(id int)"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := tt.setUp(t)
			assert.NoError(t, err)

			a := getAction(tt.sqls, ActionTypeExecute, d)
			if err := a.execute(); (err != nil) != tt.wantErr {
				t.Errorf("action.execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestScoreTask(t *testing.T) {
	task := &model.Task{
		PassRate: 0.5,
		ExecuteSQLs: []*model.ExecuteSQL{
			{
				AuditLevel: "warn",
			}, {
				AuditLevel: "normal",
			}, {
				AuditLevel: "warn",
			}, {
				AuditLevel: "normal",
			}, {
				AuditLevel: "notice",
			}, {
				AuditLevel: "normal",
			}, {
				AuditLevel: "error",
			}, {
				AuditLevel: "normal",
			}, {
				AuditLevel: "normal",
			}, {
				AuditLevel: "normal",
			}, {
				AuditLevel: "normal",
			}, {
				AuditLevel: "normal",
			}, {
				AuditLevel: "normal",
			},
		},
	}
	score := scoreTask(task)

	assert.Equal(t, int32(45), score)
}
