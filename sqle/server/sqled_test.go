package server

import (
	"testing"

	_ "actiontech.cloud/sqle/sqle/sqle/driver/mysql"
	"actiontech.cloud/sqle/sqle/sqle/model"
	"github.com/stretchr/testify/assert"
)

func TestAction_validation(t *testing.T) {
	actions := map[int]*Action{
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
