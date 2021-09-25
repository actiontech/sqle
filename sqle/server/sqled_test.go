package server

import (
	"context"
	_driver "database/sql/driver"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/actiontech/sqle/sqle/driver"
	_ "github.com/actiontech/sqle/sqle/driver/mysql"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/stretchr/testify/assert"
)

func getAction(typ int) *action {
	task := &model.Task{
		Model:     model.Model{ID: 1},
		SQLSource: model.TaskSQLSourceFromMyBatisXMLFile,
		ExecuteSQLs: []*model.ExecuteSQL{
			{
				BaseSQL: model.BaseSQL{Content: "select * from t1"},
			},
		},
	}

	entry := log.NewEntry().WithField("task_id", task.ID)
	return &action{
		task:   task,
		driver: &mockDriver{},
		typ:    typ,
		entry:  entry,
		done:   make(chan struct{}),
	}
}

type mockDriver struct {
}

func (d *mockDriver) Close(ctx context.Context) {
	return
}

func (d *mockDriver) Ping(ctx context.Context) error {
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

func (d *mockDriver) Parse(ctx context.Context, sqlText string) ([]driver.Node, error) {
	return []driver.Node{{Text: sqlText}}, nil
}

func (d *mockDriver) Audit(ctx context.Context, rules []*model.Rule, sql string) (*driver.AuditResult, error) {
	return nil, nil
}

func (d *mockDriver) GenRollbackSQL(ctx context.Context, sql string) (string, string, error) {
	return "", "", nil
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
	act := getAction(ActionTypeAudit)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rule_templates`")).
		WillReturnRows(sqlmock.NewRows([]string{}))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `sql_whitelist`")).
		WillReturnRows(sqlmock.NewRows([]string{"value", "match_type"}).AddRow(whitelist.Value, whitelist.MatchType))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `sql_whitelist`")).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow("1"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `execute_sql_detail`")).
		WithArgs(model.MockTime, model.MockTime, nil, 0, 0, act.task.ExecuteSQLs[0].Content, "", 0, "", 0, 0, "", model.SQLAuditStatusFinished, "[normal]白名单", "2882fdbb7d5bcda7b49ea0803493467e", "normal").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `tasks`")).
		WithArgs(float64(1), model.TaskStatusAudited, act.task.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = act.audit()
	assert.NoError(t, err)
	assert.Equal(t, model.TaskStatusAudited, act.task.Status)
	assert.Equal(t, float64(1), act.task.PassRate)
}
