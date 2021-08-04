package server

import (
	"context"
	_errors "errors"
	"fmt"
	"strings"
	"sync"

	"actiontech.cloud/sqle/sqle/sqle/driver"
	_ "actiontech.cloud/sqle/sqle/sqle/driver/mysql"
	"actiontech.cloud/sqle/sqle/sqle/errors"
	"actiontech.cloud/sqle/sqle/sqle/log"
	"actiontech.cloud/sqle/sqle/sqle/model"
	"actiontech.cloud/sqle/sqle/sqle/utils"

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
	queue chan *Action
}

func InitSqled(exit chan struct{}) {
	sqled = &Sqled{
		exit:        exit,
		currentTask: map[string]struct{}{},
		queue:       make(chan *Action, 1024),
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
func (s *Sqled) addTask(taskId string, typ int) (*Action, error) {
	var err error
	var d driver.Driver
	entry := log.NewEntry().WithField("task_id", taskId)
	action := &Action{
		typ:   typ,
		entry: entry,
		Done:  make(chan struct{}),
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

	task, exist, err := model.GetStorage().GetTaskDetailById(taskId)
	if err != nil {
		goto Error
	}
	if !exist {
		err = errors.New(errors.TaskNotExist, fmt.Errorf("task not exist"))
		goto Error
	}

	if err = action.validation(task); err != nil {
		goto Error
	}
	action.task = task

	// d will be closed in Sqled.do().
	if d, err = driver.NewDriver(entry, task.Instance, task.Schema); err != nil {
		goto Error
	}
	action.driver = d

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
	<-action.Done
	return action.task, action.Error
}

func (s *Sqled) Start() {
	go s.taskLoop()
	go s.cleanLoop()
}

// taskLoop is a task loop used to receive action from queue.
func (s *Sqled) taskLoop() {
	for {
		select {
		case <-s.exit:
			return
		case action := <-s.queue:
			go s.do(action)
		}
	}
}

func (s *Sqled) do(action *Action) error {
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
		action.Error = err
	}

	action.driver.Close()

	s.Lock()
	taskId := fmt.Sprintf("%d", action.task.ID)
	delete(s.currentTask, taskId)
	s.Unlock()

	select {
	case action.Done <- struct{}{}:
	default:
	}
	return err
}

const (
	ActionTypeAudit = iota + 1
	ActionTypeExecute
	ActionTypeRollback
)

// Action is an action for the task;
// when you want to execute a task, you can define an action whose type is rollback.
type Action struct {
	sync.Mutex

	// driver is interface which communicate with specify instance.
	driver driver.Driver

	task  *model.Task
	entry *logrus.Entry

	// typ is action type.
	typ   int
	Error error
	Done  chan struct{}
}

var (
	ErrActionExecuteOnExecutedTask       = _errors.New("task has been executed, can not do execute on it")
	ErrActionExecuteOnNonAuditedTask     = _errors.New("task has not been audited, can not do execute on it")
	ErrActionRollbackOnRollbackedTask    = _errors.New("task has been rollbacked, can not do rollback on it")
	ErrActionRollbackOnExecuteFailedTask = _errors.New("task has been executed failed, can not do rollback on it")
	ErrActionRollbackOnNonExecutedTask   = _errors.New("task has not been executed, can not do rollback on it")
)

// validation validate whether task can do action type(a.typ) or not.
func (a *Action) validation(task *model.Task) error {
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

func (a *Action) audit() error {
	st := model.GetStorage()

	task := a.task

	rules, err := st.GetRulesByInstanceId(fmt.Sprintf("%v", task.InstanceId))
	if err != nil {
		return err
	}
	var ptrRules []*model.Rule
	for i := range rules {
		ptrRules = append(ptrRules, &rules[i])
	}

	whitelist, _, err := st.GetSqlWhitelist(0, 0)
	if err != nil {
		return err
	}
	for _, executeSQL := range task.ExecuteSQLs {
		nodes, err := a.driver.Parse(executeSQL.Content)
		if err != nil {
			return err
		}

		if len(nodes) != 1 {
			return driver.ErrNodesCountExceedOne
		}

		sourceFP, err := nodes[0].Fingerprint()
		if err != nil {
			return err
		}

		var whitelistMatch bool
		for _, wl := range whitelist {
			if wl.MatchType == model.SQLWhitelistFPMatch {
				wlNodes, err := a.driver.Parse(wl.Value)
				if err != nil {
					return err
				}
				if len(wlNodes) != 1 {
					return driver.ErrNodesCountExceedOne
				}

				wlFP, err := wlNodes[0].Fingerprint()
				if err != nil {
					return err
				}

				if sourceFP == wlFP {
					whitelistMatch = true
				}
			} else {
				rawSQL := nodes[0].Text()
				if wl.CapitalizedValue == strings.ToUpper(rawSQL) {
					whitelistMatch = true
				}
			}
		}

		result := driver.NewInspectResults()
		if whitelistMatch {
			result.Add(model.RuleLevelNormal, "白名单")
		} else {
			result, err = a.driver.Audit(ptrRules, executeSQL.Content)
			if err != nil {
				return err
			}
		}

		executeSQL.AuditStatus = model.SQLAuditStatusFinished
		executeSQL.AuditLevel = result.Level()
		executeSQL.AuditResult = result.Message()
		executeSQL.AuditFingerprint = utils.Md5String(string(append([]byte(result.Message()), []byte(sourceFP)...)))

		a.entry.WithFields(logrus.Fields{
			"SQL":    executeSQL.Content,
			"level":  executeSQL.AuditLevel,
			"result": executeSQL.AuditResult}).Info("audit finished")
	}

	if task.SQLSource == model.TaskSQLSourceFromMyBatisXMLFile {
		a.entry.Warn("skip generate rollback SQLs")
	} else {
		var rollbackSQLs []*model.RollbackSQL
		for _, executeSQL := range task.ExecuteSQLs {
			rollbackSQL, reason, err := a.driver.GenRollbackSQL(executeSQL.Content)
			if err != nil {
				return err
			}
			result := driver.NewInspectResults()
			result.Add(executeSQL.AuditLevel, executeSQL.AuditResult)
			result.Add(model.RuleLevelNotice, reason)
			executeSQL.AuditLevel = result.Level()
			executeSQL.AuditResult = result.Message()

			rollbackSQLs = append(rollbackSQLs, &model.RollbackSQL{
				BaseSQL: model.BaseSQL{
					TaskId:  executeSQL.TaskId,
					Content: rollbackSQL,
				},
				ExecuteSQLId: executeSQL.ID,
			})
		}

		if err = st.UpdateRollbackSQLs(rollbackSQLs); err != nil {
			a.entry.Errorf("save rollback SQLs error:%v", err)
			return err
		}
	}

	if err = st.UpdateExecuteSQLs(task.ExecuteSQLs); err != nil {
		a.entry.Errorf("save SQLs error:%v", err)
		return err
	}

	var normalCount float64
	for _, executeSQL := range task.ExecuteSQLs {
		if executeSQL.AuditLevel == model.RuleLevelNormal {
			normalCount += 1
		}
	}

	var hasDDL bool
	var hasDML bool
	for _, executeSQL := range task.ExecuteSQLs {
		nodes, err := a.driver.Parse(executeSQL.Content)
		if err != nil {
			a.entry.Error(err.Error())
			continue
		}
		if len(nodes) != 1 {
			a.entry.Errorf(driver.ErrNodesCountExceedOne.Error())
			continue
		}

		switch nodes[0].Type() {
		case model.SQLTypeDDL:
			hasDDL = true
		case model.SQLTypeDML:
			hasDML = true
		}
	}

	var sqlType string
	if hasDML && hasDDL {
		sqlType = model.SQLTypeMulti
	} else if hasDDL {
		sqlType = model.SQLTypeDDL
	} else if hasDML {
		sqlType = model.SQLTypeDML
	}

	task.Status = model.TaskStatusAudited
	if err = st.UpdateTask(task, map[string]interface{}{
		"sql_type":  sqlType,
		"pass_rate": utils.Round(normalCount/float64(len(task.ExecuteSQLs)), 4),
		"status":    model.TaskStatusAudited,
	}); err != nil {
		a.entry.Errorf("update task error:%v", err)
		return err
	}
	return nil
}

func (a *Action) execute() (err error) {
	task := a.task

	a.entry.Info("start execution...")

	if err = model.GetStorage().UpdateTaskStatusById(task.ID, model.TaskStatusExecuting); nil != err {
		return
	}

	var txSQLs []*model.ExecuteSQL
	for i, executeSQL := range task.ExecuteSQLs {
		var nodes []driver.Node
		if nodes, err = a.driver.Parse(executeSQL.Content); err != nil {
			goto UpdateTask
		}

		switch nodes[0].Type() {
		case model.SQLTypeDML:
			txSQLs = append(txSQLs, executeSQL)

			if i == len(task.ExecuteSQLs)-1 {
				if err = a.execSQLs(txSQLs); err != nil {
					goto UpdateTask
				}
			}
		case model.SQLTypeDDL:
			if len(txSQLs) > 0 {
				if err = a.execSQLs(txSQLs); err != nil {
					goto UpdateTask
				}
				txSQLs = nil
			}
			if err = a.execSQL(executeSQL); err != nil {
				goto UpdateTask
			}

		default:
			err = fmt.Errorf("unknown SQL type %v", nodes[0].Type())
			goto UpdateTask
		}
	}

UpdateTask:
	taskStatus := model.TaskStatusExecuteSucceeded
	if err != nil {
		taskStatus = model.TaskStatusExecuteFailed
	}
	for _, sql := range task.ExecuteSQLs {
		if sql.ExecStatus == model.SQLExecuteStatusFailed {
			taskStatus = model.TaskStatusExecuteFailed
			break
		}
	}

	a.entry.WithField("task_status", taskStatus).Infof("execution is complated, err:%v", err)
	return model.GetStorage().UpdateTaskStatusById(task.ID, taskStatus)
}

// execSQL execute SQL and update SQL's executed status to storage.
func (a *Action) execSQL(executeSQL *model.ExecuteSQL) error {
	st := model.GetStorage()

	if err := st.UpdateExecuteSqlStatus(&executeSQL.BaseSQL, model.SQLExecuteStatusDoing, ""); err != nil {
		return err
	}

	_, err := a.driver.Exec(context.TODO(), executeSQL.Content)
	if err != nil {
		executeSQL.ExecStatus = model.SQLExecuteStatusFailed
		executeSQL.ExecResult = err.Error()
	} else {
		executeSQL.ExecStatus = model.SQLExecuteStatusSucceeded
		executeSQL.ExecResult = model.TaskExecResultOK
	}
	if err := st.Save(executeSQL); err != nil {
		return err
	}
	return nil
}

// execSQLs execute SQLs and update SQLs' executed status to storage.
func (a *Action) execSQLs(executeSQLs []*model.ExecuteSQL) error {
	st := model.GetStorage()

	if err := st.UpdateExecuteSQLStatusByTaskId(a.task, model.SQLExecuteStatusDoing); err != nil {
		return err
	}

	qs := make([]string, 0, len(executeSQLs))
	for _, executeSQL := range executeSQLs {
		qs = append(qs, executeSQL.Content)
	}

	results, txErr := a.driver.Tx(context.TODO(), qs...)
	for idx, executeSQL := range executeSQLs {
		if txErr != nil {
			executeSQL.ExecStatus = model.SQLExecuteStatusFailed
			executeSQL.ExecResult = txErr.Error()
			continue
		}
		rowAffects, _ := results[idx].RowsAffected()
		executeSQL.RowAffects = rowAffects
		executeSQL.ExecStatus = model.SQLExecuteStatusSucceeded
		executeSQL.ExecResult = model.TaskExecResultOK
	}

	return st.UpdateExecuteSQLs(executeSQLs)
}

func (a *Action) rollback() (err error) {
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

		nodes, err := a.driver.Parse(rollbackSQL.Content)
		if err != nil {
			return err
		}
		// todo: execute in transaction
		for _, node := range nodes {
			currentSQL := model.RollbackSQL{BaseSQL: model.BaseSQL{
				TaskId:  rollbackSQL.TaskId,
				Content: node.Text(),
			}, ExecuteSQLId: rollbackSQL.ExecuteSQLId}
			_, execErr := a.driver.Exec(context.TODO(), node.Text())
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
