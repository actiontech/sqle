package server

import (
	"context"
	_errors "errors"
	"fmt"
	"sync"
	"math"
	"time"

	"github.com/actiontech/sqle/sqle/driver"
	_ "github.com/actiontech/sqle/sqle/driver/mysql"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"
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
	var d driver.Driver
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
	if d, err = newDriverWithAudit(entry, task.Instance, task.Schema, task.DBType); err != nil {
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
	<-action.done
	return action.task, action.err
}

func (s *Sqled) Start() {
	go s.taskLoop()
	go s.cleanLoop()
	go s.workflowScheduleLoop()
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

	action.driver.Close(context.TODO())

	s.Lock()
	taskId := fmt.Sprintf("%d", action.task.ID)
	delete(s.currentTask, taskId)
	s.Unlock()

	select {
	case action.done <- struct{}{}:
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
type action struct {
	sync.Mutex

	// driver is interface which communicate with specify instance.
	driver driver.Driver

	task  *model.Task
	entry *logrus.Entry

	// typ is action type.
	typ  int
	err  error
	done chan struct{}
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

	err = audit(a.entry, a.task, a.driver)
	if err != nil {
		return err
	}

	// skip generate if audit is static
	if a.task.SQLSource == model.TaskSQLSourceFromMyBatisXMLFile || a.task.InstanceId == 0 {
		a.entry.Warn("skip generate rollback SQLs")
	} else {
		d, err := newDriverWithAudit(a.entry, a.task.Instance, a.task.Schema, a.task.DBType)
		if err != nil {
			return xerrors.Wrap(err, "new driver for generate rollback SQL")
		}
		defer d.Close(context.TODO())

		rollbackSQLs ,err := genRollbackSQL(a.entry, a.task, d)
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

	var normalCount float64
	maxAuditLevel := driver.RuleLevelNormal
	for _, executeSQL := range a.task.ExecuteSQLs {
		if executeSQL.AuditLevel == string(driver.RuleLevelNormal) {
			normalCount += 1
		}
		if driver.RuleLevel(executeSQL.AuditLevel).More(maxAuditLevel) {
			maxAuditLevel = driver.RuleLevel(executeSQL.AuditLevel)
		}
	}
	a.task.PassRate = utils.Round(normalCount/float64(len(a.task.ExecuteSQLs)), 4)
	a.task.AuditLevel = string(maxAuditLevel)
	a.task.Score = scoreTask(a.task)

	a.task.Status = model.TaskStatusAudited
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

// Scoring rules from https://github.com/actiontech/sqle/issues/284
func scoreTask(task *model.Task) int32 {
	var (
		numberOfTask           float64
		numberOfLessThanError  float64
		numberOfLessThanWarn   float64
		numberOfLessThanNotice float64
		errorRate              float64
		warnRate               float64
		noticeRate             float64
		totalScore             float64
	)
	{ // ready to work
		numberOfTask = float64(len(task.ExecuteSQLs))

		for _, e := range task.ExecuteSQLs {
			switch driver.RuleLevel(e.AuditLevel) {
			case driver.RuleLevelError:
				numberOfLessThanError++
			case driver.RuleLevelWarn:
				numberOfLessThanWarn++
			case driver.RuleLevelNotice:
				numberOfLessThanNotice++
			}
		}

		numberOfLessThanNotice = numberOfLessThanNotice + numberOfLessThanWarn + numberOfLessThanError
		numberOfLessThanWarn = numberOfLessThanWarn + numberOfLessThanError

		errorRate = numberOfLessThanError / numberOfTask
		warnRate = numberOfLessThanWarn / numberOfTask
		noticeRate = numberOfLessThanNotice / numberOfTask
	}
	{ // calculate the total score
		// pass rate score
		totalScore = task.PassRate * 30
		// SQL occurrence probability below error level
		totalScore += (1 - errorRate) * 15
		// SQL occurrence probability below warn level
		totalScore += (1 - warnRate) * 10
		// SQL occurrence probability below notice level
		totalScore += (1 - noticeRate) * 5
		// SQL without error level
		if errorRate == 0 {
			totalScore += 15
		}
		// SQL without warn level
		if warnRate == 0 {
			totalScore += 10
		}
		// SQL without notice level
		if noticeRate == 0 {
			totalScore += 5
		}
		// the proportion of SQL with the level below error exceeds 90%
		if errorRate < 0.1 {
			totalScore += 5
		}
		// the proportion of SQL with the level below warn exceeds 90%
		if warnRate < 0.1 {
			totalScore += 3
		}
		// the proportion of SQL with the level below warn exceeds 90%
		if noticeRate < 0.1 {
			totalScore += 2
		}
	}

	return int32(math.Floor(totalScore))
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

	// txSQLs keep adjacent DMLs, execute in one transaction.
	var txSQLs []*model.ExecuteSQL

outerLoop:
	for i, executeSQL := range task.ExecuteSQLs {
		var nodes []driver.Node
		if nodes, err = a.driver.Parse(context.TODO(), executeSQL.Content); err != nil {
			break outerLoop
		}

		switch nodes[0].Type {
		case driver.SQLTypeDML:
			txSQLs = append(txSQLs, executeSQL)

			if i == len(task.ExecuteSQLs)-1 {
				if err = a.execSQLs(txSQLs); err != nil {
					break outerLoop
				}
			}

		case driver.SQLTypeDDL:
			if len(txSQLs) > 0 {
				if err = a.execSQLs(txSQLs); err != nil {
					break outerLoop
				}
				txSQLs = nil
			}
			if err = a.execSQL(executeSQL); err != nil {
				break outerLoop
			}

		default:
			err = fmt.Errorf("unknown SQL type %v", nodes[0].Type)
			break outerLoop
		}
	}

	taskStatus := model.TaskStatusExecuteSucceeded

	if err != nil {
		taskStatus = model.TaskStatusExecuteFailed
	} else {
		for _, sql := range task.ExecuteSQLs {
			if sql.ExecStatus == model.SQLExecuteStatusFailed {
				taskStatus = model.TaskStatusExecuteFailed
				break
			}
		}
	}

	a.entry.WithField("task_status", taskStatus).Infof("execution is completed, err:%v", err)

	attrs = map[string]interface{}{
		"status":      taskStatus,
		"exec_end_at": time.Now(),
	}
	return st.UpdateTask(task, attrs)
}

// execSQL execute SQL and update SQL's executed status to storage.
func (a *action) execSQL(executeSQL *model.ExecuteSQL) error {
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

		nodes, err := a.driver.Parse(context.TODO(), rollbackSQL.Content)
		if err != nil {
			return err
		}
		// todo: execute in transaction
		for _, node := range nodes {
			currentSQL := model.RollbackSQL{BaseSQL: model.BaseSQL{
				TaskId:  rollbackSQL.TaskId,
				Content: node.Text,
			}, ExecuteSQLId: rollbackSQL.ExecuteSQLId}
			_, execErr := a.driver.Exec(context.TODO(), node.Text)
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

func newDriverWithAudit(l *logrus.Entry, inst *model.Instance, database string, dbType string) (driver.Driver, error) {
	if inst == nil && dbType == "" {
		return nil, xerrors.Errorf("instance is nil and dbType is nil")
	}

	if dbType == "" {
		dbType = inst.DbType
	}

	st := model.GetStorage()

	var err error
	var dsn *driver.DSN
	var modelRules []*model.Rule

	if inst == nil {
		templateName := st.GetDefaultRuleTemplateName(dbType)
		modelRules, err = st.GetRulesFromRuleTemplateByName(templateName)
	} else {
		dsn = &driver.DSN{
			Host:     inst.Host,
			Port:     inst.Port,
			User:     inst.User,
			Password: inst.Password,

			DatabaseName: database,
		}

		modelRules, err = st.GetRulesByInstanceId(fmt.Sprintf("%v", inst.ID))
	}

	if err != nil {
		return nil, xerrors.Errorf("get rules error: %v", err)
	}

	rules := make([]*driver.Rule, len(modelRules))
	for i, rule := range modelRules {
		rules[i] = model.ConvertRuleToDriverRule(rule)
	}

	cfg, err := driver.NewConfig(dsn, rules)
	if err != nil {
		return nil, xerrors.Wrap(err, "new driver with audit")
	}

	return driver.NewDriver(l, dbType, cfg)
}
