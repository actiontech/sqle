package server

import (
	"context"
	_errors "errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/go-sql-driver/mysql"

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
func (s *Sqled) addTask(taskId string, typ int) (*action, error) {
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

	task, exist, err := st.GetTaskDetailById(taskId)
	if err != nil {
		goto Error
	}
	if !exist {
		err = errors.New(errors.TaskNotExist, fmt.Errorf("task not exist"))
		goto Error
	}
	if task.InstanceId != 0 {
		instance, exist, err = dms.GetInstancesById(context.Background(), task.InstanceId)
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
	rules, customRules, err = st.GetAllRulesByTmpNameAndProjectIdInstanceDBType("", "", task.Instance, task.DBType)
	if err != nil {
		goto Error
	}
	p, err = newDriverManagerWithAudit(entry, task.Instance, task.Schema, task.DBType, rules)
	if err != nil {
		goto Error
	}
	action.plugin = p
	action.customRules = customRules
	action.rules = rules

	s.queue <- action

	return action, nil

Error:
	s.Lock()
	delete(s.currentTask, taskId)
	s.Unlock()
	return action, err
}

func (s *Sqled) AddTask(taskId string, typ int) error {
	_, err := s.addTask(taskId, typ)
	return err
}

func (s *Sqled) AddTaskWaitResult(taskId string, typ int) (*model.Task, error) {
	action, err := s.addTask(taskId, typ)
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

	plugin driver.Plugin

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

	err = audit(a.entry, a.task, a.plugin, a.customRules)
	if err != nil {
		return err
	}

	// skip generate if audit is static
	if a.task.SQLSource == model.TaskSQLSourceFromMyBatisXMLFile || a.task.SQLSource == model.TaskSQLSourceFromZipFile || a.task.SQLSource == model.TaskSQLSourceFromGitRepository || a.task.InstanceId == 0 {
		a.entry.Warn("skip generate rollback SQLs")
	} else if !driver.GetPluginManager().IsOptionalModuleEnabled(a.task.DBType, driverV2.OptionalModuleGenRollbackSQL) {
		a.entry.Infof("skip generate rollback SQLs, %v", driver.NewErrPluginAPINotImplement(driverV2.OptionalModuleGenRollbackSQL))
	} else {
		p, err := newDriverManagerWithAudit(a.entry, a.task.Instance, a.task.Schema, a.task.DBType, a.rules)
		if err != nil {
			return xerrors.Wrap(err, "new driver for generate rollback SQL")
		}
		defer p.Close(context.TODO())

		rollbackSQLs, err := genRollbackSQL(a.entry, a.task, p)
		if err != nil {
			return err
		}

		if err = st.UpdateRollbackSQLs(rollbackSQLs); err != nil {
			a.entry.Errorf("save rollback SQLs error:%v", err)
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
			taskStatus = model.TaskStatusExecuteFailed
		} else {
			taskStatus = model.TaskStatusExecuteSucceeded
		}
		// update task status by sql
		for _, sql := range task.ExecuteSQLs {
			if sql.ExecStatus == model.SQLExecuteStatusFailed ||
				sql.ExecStatus == model.SQLExecuteStatusTerminateSucc {
				taskStatus = model.TaskStatusExecuteFailed
				break
			}
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

		} else {
			a.terminatedSuccessfully() // NOTE: 如果中止成功，SQLs 状态已经被更新
		}
		taskStatus = model.TaskStatusExecuteFailed

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

	task := a.task

	// txSQLs keep adjacent DMLs, execute in one transaction.
	var txSQLs []*model.ExecuteSQL

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
		if a.hasTermination() && _errors.Is(mysql.ErrInvalidConn, execErr) {
			executeSQL.ExecStatus = model.SQLExecuteStatusTerminateSucc
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
		if txErr != nil {
			executeSQL.ExecStatus = model.SQLExecuteStatusFailed
			executeSQL.ExecResult = txErr.Error()
			if a.hasTermination() && _errors.Is(mysql.ErrInvalidConn, txErr) {
				executeSQL.ExecStatus = model.SQLExecuteStatusTerminateSucc
				if idx >= len(results) {
					continue
				}
				if results[idx] == nil {
					continue
				}
				rowAffects, _ := results[idx].RowsAffected()
				executeSQL.RowAffects = rowAffects
			}
			continue
		}
		rowAffects, _ := results[idx].RowsAffected()
		executeSQL.RowAffects = rowAffects
		executeSQL.ExecStatus = model.SQLExecuteStatusSucceeded
		executeSQL.ExecResult = model.TaskExecResultOK
	}
	if err := st.UpdateExecuteSQLs(executeSQLs); err != nil {
		return err
	}
	if txErr != nil {
		return txErr
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
			_, execErr := a.plugin.Exec(context.TODO(), node.Text)
			if execErr != nil {
				currentSQL.ExecStatus = model.SQLExecuteStatusFailed
				currentSQL.ExecResult = execErr.Error()
			} else {
				currentSQL.ExecStatus = model.SQLExecuteStatusSucceeded
				currentSQL.ExecResult = model.TaskExecResultOK
			}
			if execErr := st.Save(currentSQL); execErr != nil {
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
