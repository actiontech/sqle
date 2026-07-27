package server

import (
	"strings"

	"github.com/actiontech/sqle/sqle/model"
)

func ensureNonEmptyFailReason(reason string) string {
	if strings.TrimSpace(reason) == "" {
		return model.OnlineFailReasonFallback
	}
	return reason
}

// markSQLsAfterNotExecuted 将 afterSQL 之后尚未执行的 SQL 标为 not_executed。
func (a *action) markSQLsAfterNotExecuted(afterSQL *model.ExecuteSQL, stage string) {
	if a.task == nil {
		return
	}
	seen := afterSQL == nil
	for _, sql := range a.task.ExecuteSQLs {
		if !seen {
			if afterSQL != nil && sql.ID == afterSQL.ID {
				seen = true
			}
			continue
		}
		if sql.ExecStatus == model.SQLExecuteStatusInitialized || sql.ExecStatus == "" {
			sql.ExecStatus = model.SQLExecuteStatusNotExecuted
			sql.ExecResult = model.SQLNotExecutedReason
			sql.FailStage = stage
		}
	}
}

func (a *action) applyTaskExecFailSummary(stage, reason string, failedSQL *model.ExecuteSQL) {
	if a.task == nil {
		return
	}
	reason = ensureNonEmptyFailReason(reason)
	a.task.ExecFailStage = stage
	a.task.ExecFailReason = reason
	// 归因失败条数：只计 exec_status=failed，不计 execute_rollback / not_executed
	count := 0
	for _, sql := range a.task.ExecuteSQLs {
		if sql.ExecStatus == model.SQLExecuteStatusFailed {
			count++
		}
	}
	if count == 0 {
		count = 1
	}
	a.task.ExecFailSQLCount = count
	if failedSQL != nil {
		a.task.ExecFailSQLNumber = failedSQL.Number
		a.task.ExecFailSQLID = failedSQL.ID
	}
}

// persistOnlineFailure 将失败 SQL 与后续未执行 SQL 双写落库，并填充任务级摘要（内存；由 execute 收尾 UpdateTask 落库）。
func (a *action) persistOnlineFailure(failedSQL *model.ExecuteSQL, stage, reason string) error {
	st := model.GetStorage()
	reason = ensureNonEmptyFailReason(reason)
	if failedSQL != nil {
		failedSQL.ExecStatus = model.SQLExecuteStatusFailed
		failedSQL.ExecResult = reason
		failedSQL.FailStage = stage
	}
	a.markSQLsAfterNotExecuted(failedSQL, stage)
	a.applyTaskExecFailSummary(stage, reason, failedSQL)
	if a.task == nil || len(a.task.ExecuteSQLs) == 0 {
		return nil
	}
	return st.UpdateExecuteSQLs(a.task.ExecuteSQLs)
}

// PersistTaskDatasourceConnectFailure 连接预检失败时落库（任务 exec_failed + SQL 明细）。
func PersistTaskDatasourceConnectFailure(task *model.Task, connectErr error) error {
	if task == nil {
		return nil
	}
	st := model.GetStorage()
	reason := ensureNonEmptyFailReason("")
	if connectErr != nil {
		reason = ensureNonEmptyFailReason(connectErr.Error())
	}
	stage := model.OnlineFailStageDatasourceConnect

	if len(task.ExecuteSQLs) == 0 {
		sqls, err := st.GetExecuteSQLsByTaskID(task.ID)
		if err != nil {
			return err
		}
		task.ExecuteSQLs = sqls
	}

	var failedSQL *model.ExecuteSQL
	for i, sql := range task.ExecuteSQLs {
		if i == 0 {
			sql.ExecStatus = model.SQLExecuteStatusFailed
			sql.ExecResult = reason
			sql.FailStage = stage
			failedSQL = sql
			continue
		}
		if sql.ExecStatus == model.SQLExecuteStatusInitialized || sql.ExecStatus == "" {
			sql.ExecStatus = model.SQLExecuteStatusNotExecuted
			sql.ExecResult = model.SQLNotExecutedReason
			sql.FailStage = stage
		}
	}
	if len(task.ExecuteSQLs) > 0 {
		if err := st.UpdateExecuteSQLs(task.ExecuteSQLs); err != nil {
			return err
		}
	}

	task.Status = model.TaskStatusExecuteFailed
	task.ExecFailStage = stage
	task.ExecFailReason = reason
	task.ExecFailSQLCount = 1
	attrs := map[string]interface{}{
		"status":              model.TaskStatusExecuteFailed,
		"exec_fail_stage":     stage,
		"exec_fail_reason":    reason,
		"exec_fail_sql_count": 1,
	}
	if failedSQL != nil {
		task.ExecFailSQLNumber = failedSQL.Number
		task.ExecFailSQLID = failedSQL.ID
		attrs["exec_fail_sql_number"] = failedSQL.Number
		attrs["exec_fail_sql_id"] = failedSQL.ID
	}
	return st.UpdateTask(task, attrs)
}
