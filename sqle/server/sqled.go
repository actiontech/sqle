package server

import (
	"fmt"
	"strings"
	"sync"

	"actiontech.cloud/sqle/sqle/sqle/driver"
	_ "actiontech.cloud/sqle/sqle/sqle/driver/mysql"
	"actiontech.cloud/sqle/sqle/sqle/errors"
	"actiontech.cloud/sqle/sqle/sqle/log"
	"actiontech.cloud/sqle/sqle/sqle/model"
	"actiontech.cloud/sqle/sqle/sqle/utils"
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

// Action is an action for the task;
// when you want to execute a task, you can define an action whose type is rollback.
type Action struct {
	sync.Mutex
	Task *model.Task
	// Typ is task type, include inspect, execute, rollback.
	Typ   int
	Error error
	Done  chan struct{}
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
	action := &Action{
		Typ:  typ,
		Done: make(chan struct{}),
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

	err = task.ValidAction(typ)
	if err != nil {
		goto Error
	}

	action.Task = task
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
	return action.Task, action.Error
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
	switch action.Typ {
	case model.TaskActionAudit:
		err = s.audit(action.Task)
	case model.TaskActionExecute:
		err = s.execute(action.Task)
	case model.TaskActionRollback:
		err = s.rollback(action.Task)
	}
	if err != nil {
		action.Error = err
	}
	s.Lock()
	taskId := fmt.Sprintf("%d", action.Task.ID)
	delete(s.currentTask, taskId)
	s.Unlock()

	select {
	case action.Done <- struct{}{}:
	default:
	}
	return err
}

func (s *Sqled) audit(task *model.Task) error {
	st := model.GetStorage()

	rules, err := st.GetRulesByInstanceId(fmt.Sprintf("%v", task.InstanceId))
	if err != nil {
		return err
	}
	var ptrRules []*model.Rule
	for i := range rules {
		ptrRules = append(ptrRules, &rules[i])
	}
	entry := log.NewEntry().WithField("task_id", task.ID)
	d, err := driver.NewDriver(entry, task.Instance, task.Schema)
	if err != nil {
		return err
	}
	defer d.Close()

	whitelist, _, err := st.GetSqlWhitelist(0, 0)
	if err != nil {
		return err
	}
	for _, executeSQL := range task.ExecuteSQLs {
		nodes, err := d.Parse(executeSQL.Content)
		if err != nil {
			return err
		}

		if len(nodes) != 1 {
			return fmt.Errorf("want get single node, but got %v", len(nodes))
		}

		sourceFP, err := nodes[0].Fingerprint()
		if err != nil {
			return err
		}

		var whitelistMatch bool
		for _, wl := range whitelist {
			if wl.MatchType == model.SQLWhitelistFPMatch {
				wlNodes, err := d.Parse(wl.Value)
				if err != nil {
					return err
				}
				if len(wlNodes) != 1 {
					return fmt.Errorf("want get single node, but got %v", len(nodes))
				}

				wlFP, err := wlNodes[0].Fingerprint()
				if err != nil {
					return err
				}

				if sourceFP == wlFP {
					whitelistMatch = true
				}
			} else {
				if wl.CapitalizedValue == strings.ToUpper(nodes[0].Text()) {
					whitelistMatch = true
				}
			}
		}

		result := driver.NewInspectResults()
		if whitelistMatch {
			result.Add(model.RuleLevelNormal, "白名单")
		} else {
			result, err = d.Audit(ptrRules, executeSQL.Content)
			if err != nil {
				return err
			}
		}

		executeSQL.AuditStatus = model.SQLAuditStatusFinished
		executeSQL.AuditLevel = result.Level()
		executeSQL.AuditResult = result.Message()
		executeSQL.AuditFingerprint = utils.Md5String(string(append([]byte(result.Message()), []byte(sourceFP)...)))
		entry.Infof("SQL=%s, level=%s, result=%s", executeSQL.Content, executeSQL.AuditLevel, executeSQL.AuditResult)
	}

	if task.SQLSource == model.TaskSQLSourceFromMyBatisXMLFile {
		entry.Warn("skip generate rollback SQLs")
	} else {
		var rollbackSQLs []*model.RollbackSQL
		for _, executeSQL := range task.ExecuteSQLs {
			rollbackSQL, reason, err := d.GenRollbackSQL(executeSQL.Content)
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
			entry.Errorf("save rollback SQLs error:%v", err)
			return err
		}
	}

	if err = st.UpdateExecuteSQLs(task.ExecuteSQLs); err != nil {
		entry.Errorf("save SQLs error:%v", err)
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
		nodes, err := d.Parse(executeSQL.Content)
		if err != nil {
			entry.Error(err.Error())
			continue
		}
		if len(nodes) != 1 {
			entry.Errorf("expected single SQL, but got %v", len(nodes))
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
		entry.Errorf("update task error:%v", err)
		return err
	}
	return nil
}

func (s *Sqled) execute(task *model.Task) error {
	if task.SQLType == model.SQLTypeDML {
		return s.executeDMLs(task)
	} else if task.SQLType == model.SQLTypeDDL {
		return s.executeDDLs(task)
	} else if task.SQLType == model.SQLTypeMulti {
		return errors.ErrSQLTypeConflict
	}
	return nil
}

func (s *Sqled) executeDDLs(task *model.Task) error {
	st := model.GetStorage()

	entry := log.NewEntry().WithField("task_id", task.ID)
	d, err := driver.NewDriver(entry, task.Instance, task.Schema)
	if err != nil {
		return err
	}
	defer d.Close()

	entry.Info("start execute")

	if err := st.UpdateTaskStatusById(task.ID, model.TaskStatusExecuting); nil != err {
		return err
	}

	for _, executeSQL := range task.ExecuteSQLs {
		if err := st.UpdateExecuteSqlStatus(&executeSQL.BaseSQL, model.SQLExecuteStatusDoing, ""); err != nil {
			return err
		}

		err := driver.Exec(d, &executeSQL.BaseSQL)
		err = st.Save(executeSQL)
		if err != nil {
			entry.Errorf("save SQL to storage error:%v", err)
			break
		}
	}

	taskStatus := model.TaskStatusExecuteSucceeded
	for _, sql := range task.ExecuteSQLs {
		if sql.ExecStatus == model.SQLExecuteStatusFailed {
			taskStatus = model.TaskStatusExecuteFailed
			break
		}
	}
	return st.UpdateTaskStatusById(task.ID, taskStatus)
}

func (s *Sqled) executeDMLs(task *model.Task) error {
	entry := log.NewEntry().WithField("task_id", task.ID)

	st := model.GetStorage()

	entry.Info("start execute")

	err := st.UpdateExecuteSQLStatusByTaskId(task, model.SQLExecuteStatusDoing)
	if err != nil {
		entry.Errorf("update execute SQL status to storage failed, error: %v", err)
		return err
	}

	if err := st.UpdateTaskStatusById(task.ID, model.TaskStatusExecuting); nil != err {
		return err
	}

	d, err := driver.NewDriver(entry, task.Instance, task.Schema)
	if err != nil {
		return err
	}
	defer d.Close()

	var baseSQLs []*model.BaseSQL
	for _, executeSQL := range task.ExecuteSQLs {
		baseSQLs = append(baseSQLs, &executeSQL.BaseSQL)
	}

	driver.Tx(d, baseSQLs)

	if err := st.UpdateExecuteSQLs(task.ExecuteSQLs); err != nil {
		entry.Errorf("save execute sql to storage failed, error: %v", err)
		if err := st.UpdateTaskStatusById(task.ID, model.TaskStatusExecuteFailed); nil != err {
			log.Logger().Errorf("update task exec_status failed: %v", err)
		}
		return err
	}

	taskStatus := model.TaskStatusExecuteSucceeded
	for _, commitSql := range task.ExecuteSQLs {
		if commitSql.ExecStatus == model.SQLExecuteStatusFailed {
			taskStatus = model.TaskStatusExecuteFailed
			break
		}
	}
	return st.UpdateTaskStatusById(task.ID, taskStatus)
}

func (s *Sqled) rollback(task *model.Task) error {
	entry := log.NewEntry().WithField("task_id", task.ID)
	entry.Info("start rollback SQL")

	d, err := driver.NewDriver(entry, task.Instance, task.Schema)
	if err != nil {
		return err
	}
	defer d.Close()

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

		nodes, err := d.Parse(rollbackSQL.Content)
		if err != nil {
			return err
		}
		// todo: execute in transaction
		for _, node := range nodes {
			currentSQL := model.RollbackSQL{BaseSQL: model.BaseSQL{
				TaskId:  rollbackSQL.TaskId,
				Content: node.Text(),
			}, ExecuteSQLId: rollbackSQL.ExecuteSQLId}

			execErr = driver.Exec(d, &currentSQL.BaseSQL)
			execErr = st.Save(currentSQL)
			if execErr != nil {
				break ExecSQLs
			}
		}
	}

	if execErr != nil {
		entry.Errorf("rollback SQL error:%v", execErr)
	} else {
		entry.Error("rollback SQL finished")
	}
	return execErr
}
