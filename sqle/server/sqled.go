package server

import (
	"context"
	_errors "errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/utils"

	_ "github.com/actiontech/sqle/sqle/driver/mysql"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	xerrors "github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

var sqled *Sqled

func GetSqled() *Sqled {
	return sqled
}

// Sqled is an async task scheduling service.
// receive tasks from queue, the tasks include inspect, execute, rollback;
// and the task will only be executed once.
type Sqled struct {
	sync.Mutex
	// exit is Sqled service exit signal.
	exit chan struct{}
	// currentTask record the current task before execution,
	// and delete it after execution.
	currentTask map[string]struct{}
	// queue is a chan used to receive tasks.
	queue chan *action
}

func InitSqled(exit chan struct{}) {
	sqled = &Sqled{
		exit:        exit,
		currentTask: map[string]struct{}{},
		queue:       make(chan *action, 1024),
	}
	sqled.Start()
}

func (s *Sqled) HasTask(taskId string) bool {
	s.Lock()
	_, ok := s.currentTask[taskId]
	s.Unlock()
	return ok
}

// addTask receive taskId and action type, using taskId and typ to create an action;
// action will be validated, and sent to Sqled.queue.
func (s *Sqled) addTask(projectId string, taskId string, typ int, execSqlIds []uint) (*action, error) {
	var err error
	var p driver.Plugin
	var rules []*model.Rule
	var customRules []*model.CustomRule
	var instance *model.Instance
	st := model.GetStorage()
	// var drvMgr driver.DriverManager
	entry := log.NewEntry().WithField("task_id", taskId)
	action := &action{
		typ:   typ,
		entry: entry,
		done:  make(chan struct{}),
	}

	s.Lock()
	_, taskRunning := s.currentTask[taskId]
	if !taskRunning {
		s.currentTask[taskId] = struct{}{}
	}
	s.Unlock()
	if taskRunning {
		return action, errors.New(errors.TaskRunning, fmt.Errorf("task is running"))
	}

	task, exist, err := st.GetTaskDetailByIdWithExecSqlIds(taskId, execSqlIds)
	if err != nil {
		goto Error
	}
	if !exist {
		err = errors.New(errors.TaskNotExist, fmt.Errorf("task not exist"))
		goto Error
	}
	if task.InstanceId != 0 {
		instance, exist, err = dms.GetInstancesById(context.Background(), fmt.Sprintf("%d", task.InstanceId))
		if err != nil {
			goto Error
		}
		if !exist {
			err = errors.New(errors.DataNotExist, fmt.Errorf("instance not exist"))
			goto Error
		}

		task.Instance = instance
	}

	if err = action.validation(task); err != nil {
		goto Error
	}
	action.task = task

	// plugin will be closed by drvMgr in Sqled.do().
	rules, customRules, err = st.GetAllRulesByTmpNameAndProjectIdInstanceDBType(task.RuleTemplateName(), projectId, task.Instance, task.DBType)
	if err != nil {
		goto Error
	}

	p, err = newDriverManagerWithAudit(entry, task.Instance, task.Schema, task.DBType, modifyRulesWithBackupMaxRows(rules, task.DBType, task.BackupMaxRows))
	if err != nil {
		goto Error
	}
	action.plugin = p
	action.customRules = customRules
	action.rules = rules
	action.projectId = projectId

	s.queue <- action

	return action, nil

Error:
	s.Lock()
	delete(s.currentTask, taskId)
	s.Unlock()
	return action, err
}

func (s *Sqled) AddTask(projectId string, taskId string, typ int) error {
	_, err := s.addTask(projectId, taskId, typ, nil)
	return err
}

func (s *Sqled) AddTaskWaitResult(projectId string, taskId string, typ int) (*model.Task, error) {
	action, err := s.addTask(projectId, taskId, typ, nil)
	if err != nil {
		return nil, err
	}
	<-action.done
	return action.task, action.err
}

func (s *Sqled) AddTaskWaitResultWithSQLIds(projectId string, taskId string, execSqlIds []uint, typ int) (*model.Task, error) {
	action, err := s.addTask(projectId, taskId, typ, execSqlIds)
	if err != nil {
		return nil, err
	}
	<-action.done
	return action.task, action.err
}

func (s *Sqled) Start() {
	go s.taskLoop()
}

// taskLoop is a task loop used to receive action from queue.
func (s *Sqled) taskLoop() {
	for {
		select {
		case <-s.exit:
			return
		case action := <-s.queue:
			go func() {
				if err := s.do(action); err != nil {
					log.NewEntry().Error("sqled task loop do action failed, error:", err)
				}
			}()
		}
	}
}

func (s *Sqled) do(action *action) error {
	var err error
	switch action.typ {
	case ActionTypeAudit:
		err = action.audit()
	case ActionTypeExecute:
		err = action.execute()
	case ActionTypeRollback:
		err = action.rollback()
	}
	if err != nil {
		action.err = err
	}

	action.plugin.Close(context.TODO())

	s.Lock()
	taskId := fmt.Sprintf("%d", action.task.ID)
	delete(s.currentTask, taskId)
	s.Unlock()

	utils.TryClose(action.done)

	return err
}

const (
	ActionTypeAudit = iota + 1
	ActionTypeExecute
	ActionTypeRollback
)

// Action is an action for the task;
// when you want to execute a task, you can define an action whose type is rollback.
type action struct {
	sync.Mutex
	projectId string
	plugin    driver.Plugin

	task  *model.Task
	entry *logrus.Entry

	// typ is action type.
	typ  int
	err  error
	done chan struct{}

	terminateStatus int // 0:no terminate, 1,terminating, 2: terminate_succeeded, 3:terminate_failed

	customRules []*model.CustomRule
	rules       []*model.Rule
}

const (
	statusNoTermination = iota
	statusTerminating
	statusTerminateSucceeded
	statusTerminateFailed
)

func (a *action) hasTermination() bool {
	a.Lock()
	defer a.Unlock()
	return a.terminateStatus != statusNoTermination
}

func (a *action) terminate() {
	a.Lock()
	a.terminateStatus = statusTerminating
	a.Unlock()
}

func (a *action) terminatedSuccessfully() {
	a.Lock()
	a.terminateStatus = statusTerminateSucceeded
	a.Unlock()
}

func (a *action) terminatedFailed() {
	a.Lock()
	a.terminateStatus = statusTerminateFailed
	a.Unlock()
}

// isConnectionTerminatedError 判断执行错误是否由连接终止操作导致。
// 该函数通过匹配错误消息字符串来识别各种数据库驱动的连接终止错误。
// 支持的错误模式：
//   - MySQL: "invalid connection" (mysql.ErrInvalidConn 的 Error() 输出)
//   - PostgreSQL: "57P01" (SQLSTATE admin_shutdown，pg_terminate_backend 触发)
//   - PostgreSQL: "conn closed" (pgconn 连接已关闭状态)
//   - GaussDB / openGauss: "canceling statement due to user request"
//   - SQL Server: "connection is broken" (SqlException，连接被 KILL 后抛出)
//   - SQL Server: "Timeout expired" (KILL 后执行中的命令收到超时异常)
//   - Hive: "hive connection terminated:" (sqle-hive-plugin wrapHiveExecErr 在 Cancel 后加的稳定前缀)
//   - Hive: "Connection not open" (gohive 在 cursor.Cancel() 后 WaitForCompletion 常见的 thrift NOT_OPEN)
//   - Hive: "Context is done" / "Context was done before the query was executed" (gohive ctx 被取消)
//   - Hive: "operation in state" + "CANCELED" (HiveServer2 透出 OperationState=CANCELED)
func isConnectionTerminatedError(err error, dbType string) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	// MySQL: 连接被 KILL 命令终止后，mysql 驱动返回 ErrInvalidConn，
	// 其 Error() 为 "invalid connection"。错误可能经过 CodeError 或 fmt.Errorf 包装，
	// 但包装后的字符串仍包含 "invalid connection"。
	if dbType == driverV2.DriverTypeMySQL && strings.Contains(errMsg, "invalid connection") {
		return true
	}
	// PostgreSQL: 连接被 pg_terminate_backend 终止后，pgx 驱动返回 PgError，
	// 其 SQLSTATE 为 57P01 (admin_shutdown)。错误经过 gRPC 传输后以字符串形式保留。
	// 典型格式: "FATAL: terminating connection due to administrator command (SQLSTATE 57P01)"
	// PostgreSQL: 连接已被标记为关闭状态后，后续操作返回 "conn closed"。
	// 这是 pgconn 的 connLockError 错误，在连接被终止后尝试复用时出现。
	if dbType == driverV2.DriverTypePostgreSQL && (strings.Contains(errMsg, "57P01") || strings.Contains(errMsg, "conn closed")) {
		return true
	}
	// GaussDB / openGauss: 当 sqle-gaussdb-plugin 的 KillProcess 调用 pg_cancel_backend
	// 时，客户端会收到 "pq: canceling statement due to user request" (SQLSTATE 57014)。
	// 该字符串可能被 driver adaptor / gRPC 包装，但子串保留。
	// 仅在 dbType == GaussDB 时启用，避免误伤未来可能接入的纯 PostgreSQL 驱动的
	// 用户主动 cancel 场景（cancel != terminate，普通 PG 客户端不应被误判为已终止）。
	if dbType == driverV2.DriverTypeGaussDB &&
		strings.Contains(errMsg, "canceling statement due to user request") {
		return true
	}
	// SQL Server: 连接被 KILL 命令终止后，System.Data.SqlClient 抛出 SqlException，
	// 经过 C# gRPC 插件包装后，错误消息中包含 "connection is broken"。
	// SQL Server: 连接被 KILL 后，正在执行的命令可能收到超时异常而非连接断开异常。
	// 典型消息: "Timeout expired. The timeout period elapsed prior to completion of the operation or the server is not responding."
	// 经过 C# gRPC 插件的 "exec failed:" 前缀包装后仍可通过子串匹配识别。
	// SQL Server: KILL 完成后，会话进入 kill 状态，后续操作返回:
	// "Cannot continue the execution because the session is in the kill state."
	if dbType == driverV2.DriverTypeSQLServer &&
		(strings.Contains(errMsg, "connection is broken") ||
			strings.Contains(errMsg, "Timeout expired") ||
			strings.Contains(errMsg, "session is in the kill state")) {
		return true
	}
	// Hive: sqle-hive-plugin 在 KillProcess 已发出 Cancel 后给后续 exec/wait error
	// 统一加 "hive connection terminated:" 稳定前缀（见
	// sqle-hive-plugin driver/hive.go wrapHiveExecErr），sqled 端只用一条规则就能
	// 把这类 SQL 行打成 terminate_succ。即便底层 gohive/Thrift 错误消息后续升级
	// 变形，前缀稳定不变。
	if strings.Contains(errMsg, "hive connection terminated:") {
		return true
	}
	// Hive: cursor.Cancel() 后 cursor.WaitForCompletion / 后续 Execute 常见的
	// thrift "Connection not open" —— gohive 在 operationHandle 已被服务端关闭后
	// 再次复用 cursor 时透出。属于 Cancel 后的派生错误，应判为终止成功。
	if strings.Contains(errMsg, "Connection not open") {
		return true
	}
	// Hive: gohive 的 cursor.Execute / WaitForCompletion 在 ctx 被取消后会返回
	// "Context was done before the query was executed" 或裸 "Context is done"
	// （前者来自 gohive cursor.go 守卫；后者是上游 ctx.Err 包装）。两种形态都
	// 等价于「终止意图已收到 + 操作未真正跑完」。
	if strings.Contains(errMsg, "Context was done before the query was executed") {
		return true
	}
	if strings.Contains(errMsg, "Context is done") {
		return true
	}
	// Hive: HiveServer2 透出 OperationState 时，典型消息形如
	// "Invalid OperationHandle: OperationHandle [opType=EXECUTE_STATEMENT, ...]: operation in state CANCELED"。
	// 同时要求出现 "operation in state" 与 "CANCELED" 两个子串，避免把
	// "operation in state RUNNING/FINISHED" 等正常状态误判为终止。
	if strings.Contains(errMsg, "operation in state") && strings.Contains(errMsg, "CANCELED") {
		return true
	}

	return false
}

// terminatedExecResult 把 driver/plugin 返回的原始 error 包装成对用户可读的
// 「因中止上线中断」文案，作为 ExecuteSQL.ExecResult 写入。命中
// isConnectionTerminatedError 时调用：上层 UI（工单 SQL 详情）就不会再展示
// `hive exec failed (sql=...): EOF` / `invalid connection` 这类对终端用户毫无
// 意义的裸 driver 错误，而是「因中止上线中断: <原文本>」——明确告诉用户
// 这条 SQL 是用户主动点了「中止上线」之后被取消的，并保留原始 error 便于
// 后续 sqled.log / dev 排查回溯。
const terminatedExecResultPrefix = "因中止上线中断"

func terminatedExecResult(err error) string {
	if err == nil {
		return terminatedExecResultPrefix
	}
	return fmt.Sprintf("%s: %s", terminatedExecResultPrefix, err.Error())
}

var (
	ErrActionExecuteOnExecutedTask       = _errors.New("task has been executed, can not do execute on it")
	ErrActionExecuteOnNonAuditedTask     = _errors.New("task has not been audited, can not do execute on it")
	ErrActionRollbackOnRollbackedTask    = _errors.New("task has been rollbacked, can not do rollback on it")
	ErrActionRollbackOnExecuteFailedTask = _errors.New("task has been executed failed, can not do rollback on it")
	ErrActionRollbackOnNonExecutedTask   = _errors.New("task has not been executed, can not do rollback on it")
)

// validation validate whether task can do action type(a.typ) or not.
func (a *action) validation(task *model.Task) error {
	switch a.typ {
	case ActionTypeAudit:
		// audit sql allowed at all times
		return nil
	case ActionTypeExecute:
		if task.HasDoingExecute() {
			return errors.New(errors.TaskActionDone, ErrActionExecuteOnExecutedTask)
		}
		if !task.HasDoingAudit() {
			return errors.New(errors.TaskActionInvalid, ErrActionExecuteOnNonAuditedTask)
		}
	case ActionTypeRollback:
		if task.HasDoingRollback() {
			return errors.New(errors.TaskActionDone, ErrActionRollbackOnRollbackedTask)
		}
		if task.IsExecuteFailed() {
			return errors.New(errors.TaskActionInvalid, ErrActionRollbackOnExecuteFailedTask)
		}
		if !task.HasDoingExecute() {
			return errors.New(errors.TaskActionInvalid, ErrActionRollbackOnNonExecutedTask)
		}
	}
	return nil
}

func (a *action) audit() (err error) {
	st := model.GetStorage()

	err = audit(a.projectId, a.entry, a.task, a.plugin, a.customRules)
	if err != nil {
		return err
	}
	backupService := BackupService{}
	if backupService.CheckCanTaskBackup(a.task) {
		backupTasks := make([]*model.BackupTask, 0, len(a.task.ExecuteSQLs))
		for _, sql := range a.task.ExecuteSQLs {
			backupTasks = append(backupTasks, initModelBackupTask(a.plugin, a.task, sql))
		}
		err = st.BatchCreateBackupTasks(backupTasks)
		if err != nil {
			return err
		}
	}

	if err = st.UpdateExecuteSQLs(a.task.ExecuteSQLs); err != nil {
		a.entry.Errorf("save SQLs error:%v", err)
		return err
	}

	if err = st.UpdateTask(a.task, map[string]interface{}{
		"pass_rate":   a.task.PassRate,
		"audit_level": a.task.AuditLevel,
		"status":      a.task.Status,
		"score":       a.task.Score,
	}); err != nil {
		a.entry.Errorf("update task error:%v", err)
		return err
	}
	return nil
}

func (a *action) terminateExecution(ctx context.Context) error {
	if !driver.GetPluginManager().
		IsOptionalModuleEnabled(a.task.DBType, driverV2.OptionalModuleKillProcess) {
		return driver.NewErrPluginAPINotImplement(driverV2.OptionalModuleKillProcess)
	}
	return a.plugin.KillProcess(ctx)
}

func (a *action) execute() (err error) {
	st := model.GetStorage()
	task := a.task

	a.entry.Info("start execution...")

	attrs := map[string]interface{}{
		"status":        model.TaskStatusExecuting,
		"exec_start_at": time.Now(),
	}
	if err = st.UpdateTask(task, attrs); err != nil {
		return err
	}

	exeErrChan := make(chan error)
	terminateErrChan := make(chan error)

	{
		go func() { // execute
			exeErrChan <- a.execTask()
		}()

		go func() { // wait for kill signal
			for {
				select {
				case <-a.done:
					return
				default:
					if a.GetTaskStatus(st) == model.TaskStatusTerminating {
						a.terminate()
						ctx, cancel := context.WithTimeout(
							context.Background(), time.Minute*2)
						defer cancel()
						terminateErrChan <- a.terminateExecution(ctx)
						return
					}
				}
				time.Sleep(time.Millisecond * 500)
			}
		}()
	}

	// update task status
	taskStatus := model.TaskStatusExecuting

	select {
	case e := <-exeErrChan:
		err = e
		if e != nil {
			// 如果用户已经触发了上线终止 (a.terminate() 已被调用)，
			// 且执行错误能被识别为"连接终止 / 语句取消"类错误，则
			// 说明 SQL 实际是被 KillProcess 中断的，应当走 terminate_succeeded 收尾，
			// 而不是默认的 exec_failed。
			// 该路径覆盖：execTask 因 plugin Tx/Exec/ExecBatch 携带 cancel 错误返回，
			// 但 terminate goroutine 的 KillProcess RPC 尚未在 select 中胜出的常见竞态。
			if a.hasTermination() && isConnectionTerminatedError(e, a.task.DBType) {
				a.terminatedSuccessfully()
				taskStatus = model.TaskStatusTerminateSucc
			} else {
				taskStatus = model.TaskStatusExecuteFailed
			}
		} else {
			taskStatus = model.TaskStatusExecuteSucceeded
		}
		// update task status by sql
		// 验证task下所有的sql是否全部成功（工单中允许重新上线部分sql，所以需要验证全部sql是否成功）
		// 注意：SQLExecuteStatusTerminateSucc 也算"非成功"，但若 task 已被识别为 terminate_succeeded
		// (整体被用户中止)，则保留 terminate 收尾状态，不再回退到 exec_failed。
		failedSqls, queryErr := st.GetExecSqlsByTaskIdAndStatus(task.ID, []string{model.SQLExecuteStatusFailed, model.SQLExecuteStatusTerminateSucc})
		if queryErr != nil {
			return queryErr
		}
		if len(failedSqls) > 0 && taskStatus != model.TaskStatusTerminateSucc {
			taskStatus = model.TaskStatusExecuteFailed
		}

	case terminationErr := <-terminateErrChan:
		if terminationErr != nil {
			a.entry.Errorf("task(%v) termination failed, err: %v", task.ID, terminationErr)
			a.terminatedFailed()
			err = terminationErr

			{ //NOTE: 由于上线中止失败，需要更新 SQLs 状态
				for i := range task.ExecuteSQLs {
					sql := task.ExecuteSQLs[i]
					if sql.ExecStatus == model.SQLExecuteStatusDoing {
						sql.ExecStatus = model.SQLExecuteStatusTerminateFailed
						sql.ExecResult = fmt.Sprintf("%v", terminationErr)
					}
				}
				if err := st.UpdateExecuteSQLs(task.ExecuteSQLs); err != nil {
					return err
				}
			}
			taskStatus = model.TaskStatusTerminateFail
		} else {
			a.terminatedSuccessfully() // NOTE: 如果中止成功，SQLs 状态已经被更新
			taskStatus = model.TaskStatusTerminateSucc
		}

	}

	a.entry.WithField("task_status", taskStatus).
		Infof("execution is completed, err:%v", err)

	a.task.Status = taskStatus

	attrs = map[string]interface{}{
		"status":      taskStatus,
		"exec_end_at": time.Now(),
	}
	return st.UpdateTask(task, attrs)
}

func (a *action) GetTaskStatus(st *model.Storage) string {
	taskStatus, err := st.GetTaskStatusByID(strconv.Itoa(int(a.task.ID)))
	if err != nil {
		a.entry.Error(err.Error())
		return ""
	}
	return taskStatus
}

func (a *action) execTask() (err error) {
	svc := BackupService{}
	if svc.CheckCanTaskBackup(a.task) {
		err = a.backupAndExecSql()
		if err != nil {
			return err
		}
		return nil
	}
	switch a.task.ExecMode {
	case model.ExecModeSqlFile:
		// check plugin can exec batch sqls
		execFileModeChecker, err := NewModuleStatusChecker(a.task.DBType, executeSqlFileMode)
		if err != nil {
			return err
		}
		if !execFileModeChecker.CheckIsSupport() {
			return fmt.Errorf("plugin %v does not support execute sql file", a.task.DBType)
		}
		err = a.execSqlFileMode()
		if err != nil {
			return err
		}
	default:
		err = a.execSqlSqlMode()
		if err != nil {
			return err
		}
	}
	return nil

}

/*
backupAndExecSql() 备份与执行SQL：

	按照顺序，先根据一条SQL备份，再执行该SQL。备份过程中涉及连库查询和保存数据。
*/
func (a *action) backupAndExecSql() error {
	for _, executeSQL := range a.task.ExecuteSQLs {
		backupMgr, err := getBackupManager(a.plugin, executeSQL, a.task.DBType, a.task.BackupMaxRows)
		if err != nil {
			return fmt.Errorf("in backupAndExecSql when getBackupManager, err %w , task: %v", err, a.task.ID)
		}
		if err = backupMgr.Backup(); err != nil {
			return fmt.Errorf("in backupAndExecSql when backupMgr Backup, err %w, backup manager: %v, task: %v", err, backupMgr, a.task.ID)
		}
		if err := a.execSQL(executeSQL); err != nil {
			return fmt.Errorf("in backupAndExecSql when execSQL %v, err %w, backup manager: %v, task: %v", executeSQL, err, backupMgr, a.task.ID)
		}
	}
	return nil
}

func (a *action) execSqlSqlMode() error {
	// txSQLs keep adjacent DMLs, execute in one transaction.
	task := a.task
	var txSQLs []*model.ExecuteSQL
	var err error
	for i := range task.ExecuteSQLs {
		executeSQL := task.ExecuteSQLs[i]
		var nodes []driverV2.Node
		if nodes, err = a.plugin.Parse(context.TODO(), executeSQL.Content); err != nil {
			return err
		}

		switch nodes[0].Type {
		case driverV2.SQLTypeDML, driverV2.SQLTypeDQL:
			txSQLs = append(txSQLs, executeSQL)
			if i == len(task.ExecuteSQLs)-1 {
				if err = a.execSQLs(txSQLs); err != nil {
					return err
				}
			}

		default:
			if len(txSQLs) > 0 {
				if err = a.execSQLs(txSQLs); err != nil {
					return err
				}
				txSQLs = nil
			}
			if err = a.execSQL(executeSQL); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *action) execSqlFileMode() error {
	files, err := getFilesSortByExecOrder(a.task.GetIDStr())
	if err != nil {
		return err
	}
	sqlsInFile := groupSqlsByFile(a.task.ExecuteSQLs)
	// execute sqls in the order of files
	for _, file := range files {
		sqls, ok := sqlsInFile[file.FileName]
		if !ok {
			continue
		}
		err = a.executeSqlsGroupByBatchId(sqls)
		if err != nil {
			if err == ErrExecuteFileFailed {
				return nil
			}
			return err
		}
	}
	return nil
}

func getFilesSortByExecOrder(taskId string) ([]*model.AuditFile, error) {
	st := model.GetStorage()
	files, err := st.GetFileByTaskId(taskId)
	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ExecOrder < files[j].ExecOrder
	})
	return files, nil
}

func groupSqlsByFile(executeSQLs []*model.ExecuteSQL) map[string][]*model.ExecuteSQL {
	fileSqlMap := make(map[string][]*model.ExecuteSQL)
	for _, executeSQL := range executeSQLs {
		fileSqlMap[executeSQL.SourceFile] = append(fileSqlMap[executeSQL.SourceFile], executeSQL)
	}
	return fileSqlMap
}

// 存在SQL文件执行失败，不再执行其他SQL文件
var ErrExecuteFileFailed error = fmt.Errorf("execute file failed, please stop execute other file")

func (a *action) executeSqlsGroupByBatchId(sqls []*model.ExecuteSQL) error {
	sqlBatch := make([]*model.ExecuteSQL, 0)
	for idx, sql := range sqls {
		sqlBatch = append(sqlBatch, sql)
		if idx < len(sqls)-1 {
			// not the last sql
			if sql.ExecBatchId != sqls[idx+1].ExecBatchId {
				// when batch id is changed, execute sql batch
				if err := a.executeSQLBatch(sqlBatch); err != nil {
					return err
				}
				if err := checkBatchExecuteStatus(sqlBatch); err != nil {
					return err
				}
				// clear sql batch
				sqlBatch = make([]*model.ExecuteSQL, 0)
			}
		} else {
			// when encount the last sql in this file execute sql batch
			if err := a.executeSQLBatch(sqlBatch); err != nil {
				return err
			}
			if err := checkBatchExecuteStatus(sqlBatch); err != nil {
				return err
			}
		}

	}
	return nil
}

func checkBatchExecuteStatus(sqlBatch []*model.ExecuteSQL) error {
	for _, sql := range sqlBatch {
		if sql.ExecStatus == model.SQLExecuteStatusFailed {
			return ErrExecuteFileFailed
		}
	}
	return nil
}

// executeSQLBatch executes a batch of SQLs and updates their status.
func (a *action) executeSQLBatch(executeSQLs []*model.ExecuteSQL) error {
	st := model.GetStorage()
	// update status befor execute
	for _, executeSQL := range executeSQLs {
		executeSQL.ExecStatus = model.SQLExecuteStatusDoing
	}
	if err := st.UpdateExecuteSQLs(executeSQLs); err != nil {
		return err
	}

	sqls := make([]string, 0, len(executeSQLs))
	for _, sql := range executeSQLs {
		sqls = append(sqls, sql.Content)
	}

	results, execErr := a.plugin.ExecBatch(context.TODO(), sqls...)
	if execErr != nil {
		for idx, executeSQL := range executeSQLs {
			executeSQL.ExecStatus = model.SQLExecuteStatusFailed
			executeSQL.ExecResult = execErr.Error()
			if a.hasTermination() && isConnectionTerminatedError(execErr, a.task.DBType) {
				executeSQL.ExecStatus = model.SQLExecuteStatusTerminateSucc
				executeSQL.ExecResult = terminatedExecResult(execErr)
				if idx >= len(results) || results[idx] == nil {
					continue
				}
				rowAffects, _ := results[idx].RowsAffected()
				executeSQL.RowAffects = rowAffects
			}
		}
	} else {
		for idx, executeSQL := range executeSQLs {
			rowAffects, _ := results[idx].RowsAffected()
			executeSQL.RowAffects = rowAffects
			executeSQL.ExecStatus = model.SQLExecuteStatusSucceeded
			executeSQL.ExecResult = model.TaskExecResultOK
		}
	}

	err := st.BatchSaveExecuteSqls(executeSQLs)
	if err != nil {
		return err
	}
	return nil
}

// execSQL execute SQL and update SQL's executed status to storage.
func (a *action) execSQL(executeSQL *model.ExecuteSQL) error {
	st := model.GetStorage()

	if err := st.UpdateExecuteSqlStatus(&executeSQL.BaseSQL, model.SQLExecuteStatusDoing, ""); err != nil {
		return err
	}

	_, execErr := a.plugin.Exec(context.TODO(), executeSQL.Content)
	if execErr != nil {
		executeSQL.ExecStatus = model.SQLExecuteStatusFailed
		executeSQL.ExecResult = execErr.Error()
		if a.hasTermination() && isConnectionTerminatedError(execErr, a.task.DBType) {
			executeSQL.ExecStatus = model.SQLExecuteStatusTerminateSucc
			executeSQL.ExecResult = terminatedExecResult(execErr)
		}
	} else {
		executeSQL.ExecStatus = model.SQLExecuteStatusSucceeded
		executeSQL.ExecResult = model.TaskExecResultOK
	}
	if err := st.Save(executeSQL); err != nil {
		return err
	}
	if execErr != nil {
		return execErr
	}
	return nil
}

// execSQLs execute SQLs and update SQLs' executed status to storage.
func (a *action) execSQLs(executeSQLs []*model.ExecuteSQL) error {
	st := model.GetStorage()

	for _, executeSQL := range executeSQLs {
		executeSQL.ExecStatus = model.SQLExecuteStatusDoing
	}
	if err := st.UpdateExecuteSQLs(executeSQLs); err != nil {
		return err
	}

	qs := make([]string, 0, len(executeSQLs))
	for _, executeSQL := range executeSQLs {
		qs = append(qs, executeSQL.Content)
	}

	results, txErr := a.plugin.Tx(context.TODO(), qs...)
	for idx, executeSQL := range executeSQLs {
		if results != nil && idx < len(results.ExecResult) {
			rowAffects, _ := results.ExecResult[idx].RowsAffected()
			executeSQL.RowAffects = rowAffects
		}
		if txErr != nil {
			executeSQL.ExecStatus = model.SQLExecuteStatusFailed
			executeSQL.ExecResult = txErr.Error()
			if a.hasTermination() && isConnectionTerminatedError(txErr, a.task.DBType) {
				executeSQL.ExecStatus = model.SQLExecuteStatusTerminateSucc
				executeSQL.ExecResult = terminatedExecResult(txErr)
			}
			continue
		}
		if results != nil && results.ExecErr != nil {
			if results.ExecErr.ErrSqlIndex == uint32(idx) {
				executeSQL.ExecStatus = model.SQLExecuteStatusFailed
				executeSQL.ExecResult = results.ExecErr.SqlExecErrMsg
				if a.hasTermination() && isConnectionTerminatedError(fmt.Errorf("%s", results.ExecErr.SqlExecErrMsg), a.task.DBType) {
					executeSQL.ExecStatus = model.SQLExecuteStatusTerminateSucc
					executeSQL.ExecResult = terminatedExecResult(fmt.Errorf(results.ExecErr.SqlExecErrMsg))
				}
			} else {
				executeSQL.ExecStatus = model.SQLExecuteStatusFailed
				executeSQL.ExecResult = model.TaskExecResultRollback
			}
			continue
		}
		executeSQL.ExecStatus = model.SQLExecuteStatusSucceeded
		executeSQL.ExecResult = model.TaskExecResultOK
	}
	if err := st.UpdateExecuteSQLs(executeSQLs); err != nil {
		return err
	}
	if txErr != nil {
		return txErr
	}
	if results != nil && results.ExecErr != nil {
		return fmt.Errorf("sql execute err msg: %v", results.ExecErr.SqlExecErrMsg)
	}
	return nil
}

func (a *action) rollback() (err error) {
	task := a.task
	a.entry.Info("start rollback SQL")

	var execErr error
	st := model.GetStorage()
ExecSQLs:
	for _, rollbackSQL := range task.RollbackSQLs {
		if rollbackSQL.Content == "" {
			continue
		}
		if err = st.UpdateRollbackSqlStatus(&rollbackSQL.BaseSQL, model.SQLExecuteStatusDoing, ""); err != nil {
			return err
		}

		nodes, err := a.plugin.Parse(context.TODO(), rollbackSQL.Content)
		if err != nil {
			return err
		}
		// todo: execute in transaction
		for _, node := range nodes {
			currentSQL := model.RollbackSQL{BaseSQL: model.BaseSQL{
				TaskId:  rollbackSQL.TaskId,
				Content: node.Text,
			}, ExecuteSQLId: rollbackSQL.ExecuteSQLId}
			_, execErr = a.plugin.Exec(context.TODO(), node.Text)
			if execErr != nil {
				currentSQL.ExecStatus = model.SQLExecuteStatusFailed
				currentSQL.ExecResult = execErr.Error()
			} else {
				currentSQL.ExecStatus = model.SQLExecuteStatusSucceeded
				currentSQL.ExecResult = model.TaskExecResultOK
			}
			if execErr = st.Save(currentSQL); execErr != nil {
				break ExecSQLs
			}
		}
	}

	if execErr != nil {
		a.entry.Errorf("rollback SQL error:%v", execErr)
	} else {
		a.entry.Error("rollback SQL finished")
	}
	return execErr
}

func newDriverManagerWithAudit(l *logrus.Entry, inst *model.Instance, database string, dbType string, modelRules []*model.Rule) (driver.Plugin, error) {
	if inst == nil && dbType == "" {
		return nil, xerrors.Errorf("instance is nil and dbType is nil")
	}

	if dbType == "" {
		dbType = inst.DbType
	}

	var dsn *driverV2.DSN

	// 填充dsn
	if inst != nil {
		dsn = &driverV2.DSN{
			Host:             inst.Host,
			Port:             inst.Port,
			User:             inst.User,
			Password:         inst.Password,
			AdditionalParams: inst.AdditionalParams,

			DatabaseName: database,
		}
	}

	rules := make([]*driverV2.Rule, len(modelRules))
	for i, rule := range modelRules {
		rules[i] = model.ConvertRuleToDriverRule(rule)
	}

	cfg := &driverV2.Config{
		DSN:   dsn,
		Rules: rules,
	}
	return driver.GetPluginManager().OpenPlugin(l, dbType, cfg)
}
