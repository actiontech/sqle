package server

import (
	"context"
	_errors "errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	syncTask "github.com/actiontech/sqle/sqle/pkg/sync_Task"

	"github.com/actiontech/sqle/sqle/notification"

	imPkg "github.com/actiontech/sqle/sqle/pkg/im"
	"github.com/actiontech/sqle/sqle/pkg/im/dingding"

	"github.com/actiontech/sqle/sqle/driver"
	_ "github.com/actiontech/sqle/sqle/driver/mysql"
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
	var d driver.Driver
	var drvMgr driver.DriverManager
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

	// d will be closed by drvMgr in Sqled.do().
	drvMgr, err = newDriverManagerWithAudit(entry, task.Instance, task.Schema, task.DBType, nil, "")
	if err != nil {
		goto Error
	}
	if d, err = drvMgr.GetAuditDriver(); err != nil {
		goto Error
	}
	action.driver = d
	action.driverMgr = drvMgr

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
	go s.dingTalkLoop()
	go s.workflowScheduleLoop()
	go s.syncInstanceTaskLoop()
}

func (s *Sqled) syncInstanceTaskLoop() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	syncTask.EnableSyncInstanceTask(context.TODO())

	for {
		select {
		case <-s.exit:
			return
		case <-ticker.C:
			syncTask.ReloadSyncInstanceTask(context.Background(), "ticker reload")
		}
	}
}

func (s *Sqled) dingTalkLoop() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.exit:
			return
		case <-ticker.C:
			if err := s.dingTalkRotation(); err != nil {
				log.NewEntry().Error("dingTalkRotation failed, error:", err)
			}
		}
	}
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

	action.driverMgr.Close(context.TODO())

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

func (s *Sqled) dingTalkRotation() error {
	st := model.GetStorage()

	ims, err := st.GetAllIMConfig()
	if err != nil {
		log.NewEntry().Errorf("get all im config failed, error: %v", err)
	}

	for _, im := range ims {
		switch im.Type {
		case model.ImTypeDingTalk:
			d := &dingding.DingTalk{
				AppKey:    im.AppKey,
				AppSecret: im.AppSecret,
			}

			dingTalkInstances, err := st.GetDingTalkInstByStatus(model.ApproveStatusInitialized)
			if err != nil {
				log.NewEntry().Errorf("get ding talk status error: %v", err)
				continue
			}

			for _, dingTalkInstance := range dingTalkInstances {
				approval, err := d.GetApprovalDetail(dingTalkInstance.ApproveInstanceCode)
				if err != nil {
					log.NewEntry().Errorf("get ding talk approval detail error: %v", err)
					continue
				}

				switch *approval.Result {
				case model.ApproveStatusAgree:
					workflow, exist, err := st.GetWorkflowDetailById(strconv.Itoa(int(dingTalkInstance.WorkflowId)))
					if err != nil {
						log.NewEntry().Errorf("get workflow detail error: %v", err)
						continue
					}
					if !exist {
						log.NewEntry().Errorf("workflow not exist, id: %d", dingTalkInstance.WorkflowId)
						continue
					}

					nextStep := workflow.NextStep()

					userId := *approval.OperationRecords[1].UserId
					user, err := getUserByUserId(d, userId, workflow.CurrentStep().Assignees)
					if err != nil {
						log.NewEntry().Errorf("get user by user id error: %v", err)
						continue
					}

					if err := ApproveWorkflowProcess(workflow, user, st); err != nil {
						log.NewEntry().Errorf("approve workflow process error: %v", err)
						continue
					}

					dingTalkInstance.Status = model.ApproveStatusAgree
					if err := st.Save(&dingTalkInstance); err != nil {
						log.NewEntry().Errorf("save ding talk instance error: %v", err)
						continue
					}

					if nextStep.Template.Typ != model.WorkflowStepTypeSQLExecute {
						imPkg.CreateApprove(strconv.Itoa(int(workflow.ID)))
					}

				case model.ApproveStatusRefuse:
					workflow, exist, err := st.GetWorkflowDetailById(strconv.Itoa(int(dingTalkInstance.WorkflowId)))
					if err != nil {
						log.NewEntry().Errorf("get workflow detail error: %v", err)
						continue
					}
					if !exist {
						log.NewEntry().Errorf("workflow not exist, id: %d", dingTalkInstance.WorkflowId)
						continue
					}

					var reason string
					if approval.OperationRecords[1] != nil && approval.OperationRecords[1].Remark != nil {
						reason = *approval.OperationRecords[1].Remark
					} else {
						reason = "审批拒绝"
					}

					userId := *approval.OperationRecords[1].UserId
					user, err := getUserByUserId(d, userId, workflow.CurrentStep().Assignees)
					if err != nil {
						log.NewEntry().Errorf("get user by user id error: %v", err)
						continue
					}

					if err := RejectWorkflowProcess(workflow, reason, user, st); err != nil {
						log.NewEntry().Errorf("reject workflow process error: %v", err)
						continue
					}

					dingTalkInstance.Status = model.ApproveStatusRefuse
					if err := st.Save(&dingTalkInstance); err != nil {
						log.NewEntry().Errorf("save ding talk instance error: %v", err)
						continue
					}
				default:
					// ding talk rotation, no action
				}
			}
		}
	}

	return nil
}

func ApproveWorkflowProcess(workflow *model.Workflow, user *model.User, s *model.Storage) error {
	currentStep := workflow.CurrentStep()

	if workflow.Record.Status == model.WorkflowStatusWaitForExecution {
		return errors.New(errors.DataInvalid,
			fmt.Errorf("workflow has been approved, you should to execute it"))
	}

	currentStep.State = model.WorkflowStepStateApprove
	now := time.Now()
	currentStep.OperateAt = &now
	currentStep.OperationUserId = user.ID
	nextStep := workflow.NextStep()
	workflow.Record.CurrentWorkflowStepId = nextStep.ID
	if nextStep.Template.Typ == model.WorkflowStepTypeSQLExecute {
		workflow.Record.Status = model.WorkflowStatusWaitForExecution
	}

	err := s.UpdateWorkflowStatus(workflow, currentStep, nil)
	if err != nil {
		return fmt.Errorf("update workflow status failed, %v", err)
	}

	go notification.NotifyWorkflow(strconv.Itoa(int(workflow.ID)), notification.WorkflowNotifyTypeApprove)

	return nil
}

func RejectWorkflowProcess(workflow *model.Workflow, reason string, user *model.User, s *model.Storage) error {
	currentStep := workflow.CurrentStep()
	currentStep.State = model.WorkflowStepStateReject
	currentStep.Reason = reason
	now := time.Now()
	currentStep.OperateAt = &now
	currentStep.OperationUserId = user.ID

	workflow.Record.Status = model.WorkflowStatusReject
	workflow.Record.CurrentWorkflowStepId = 0

	if err := s.UpdateWorkflowStatus(workflow, currentStep, nil); err != nil {
		return fmt.Errorf("update workflow status failed, %v", err)
	}

	go notification.NotifyWorkflow(fmt.Sprintf("%v", workflow.ID), notification.WorkflowNotifyTypeReject)

	return nil
}

func getUserByUserId(d *dingding.DingTalk, userId string, assignees []*model.User) (*model.User, error) {
	phone, err := d.GetMobileByUserID(userId)
	if err != nil {
		return nil, fmt.Errorf("get user mobile error: %v", err)
	}

	for _, assignee := range assignees {
		if assignee.Phone == phone {
			return assignee, nil
		}
	}

	return nil, fmt.Errorf("user not found, phone: %s", phone)
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
	driver    driver.Driver
	driverMgr driver.DriverManager

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
		drvMgr, err := newDriverManagerWithAudit(a.entry, a.task.Instance, a.task.Schema, a.task.DBType, nil, "")
		if err != nil {
			return xerrors.Wrap(err, "new driver for generate rollback SQL")
		}
		defer drvMgr.Close(context.TODO())

		d, err := drvMgr.GetAuditDriver()
		if err != nil {
			return err
		}

		rollbackSQLs, err := genRollbackSQL(a.entry, a.task, d)
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

		default:
			if len(txSQLs) > 0 {
				if err = a.execSQLs(txSQLs); err != nil {
					break outerLoop
				}
				txSQLs = nil
			}
			if err = a.execSQL(executeSQL); err != nil {
				break outerLoop
			}
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
	task.Status = taskStatus

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

	_, execErr := a.driver.Exec(context.TODO(), executeSQL.Content)
	if execErr != nil {
		executeSQL.ExecStatus = model.SQLExecuteStatusFailed
		executeSQL.ExecResult = execErr.Error()
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

func newDriverManagerWithAudit(l *logrus.Entry, inst *model.Instance, database string, dbType string, projectId *uint, ruleTemplateName string) (driver.DriverManager, error) {
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

	// 填充规则
	{
		if ruleTemplateName != "" {
			if projectId == nil {
				return nil, xerrors.New("project id is needed when rule template name is given")
			}
			modelRules, err = st.GetRulesFromRuleTemplateByName([]uint{*projectId, model.ProjectIdForGlobalRuleTemplate}, ruleTemplateName)
		} else {
			if inst != nil {
				modelRules, err = st.GetRulesByInstanceId(fmt.Sprintf("%v", inst.ID))
			} else {
				templateName := st.GetDefaultRuleTemplateName(dbType)
				// 默认规则模板从全局模板里拿
				modelRules, err = st.GetRulesFromRuleTemplateByName([]uint{model.ProjectIdForGlobalRuleTemplate}, templateName)
			}
		}
		if err != nil {
			return nil, xerrors.Errorf("get rules error: %v", err)
		}
	}

	// 填充dsn
	if inst != nil {
		dsn = &driver.DSN{
			Host:             inst.Host,
			Port:             inst.Port,
			User:             inst.User,
			Password:         inst.Password,
			AdditionalParams: inst.AdditionalParams,

			DatabaseName: database,
		}
	}

	rules := make([]*driver.Rule, len(modelRules))
	for i, rule := range modelRules {
		rules[i] = model.ConvertRuleToDriverRule(rule)
	}

	cfg, err := driver.NewConfig(dsn, rules)
	if err != nil {
		return nil, xerrors.Wrap(err, "new driver with audit")
	}

	return driver.NewDriverManger(l, dbType, cfg)
}
