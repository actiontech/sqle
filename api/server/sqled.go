package server

import (
	"fmt"
	"sqle/errors"
	"sqle/executor"
	"sqle/inspector"
	"sqle/model"
	"sync"
)

var sqled *Sqled

func GetSqled() *Sqled {
	return sqled
}

type Sqled struct {
	sync.Mutex
	exit            chan struct{}
	currentTask     map[string]struct{}
	queue           chan *Action
	instancesStatus map[uint]*InstanceStatus
}

type Action struct {
	sync.Mutex
	Task  *model.Task
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
		return action, errors.New(errors.TASK_RUNNING, fmt.Errorf("task is running"))
	}

	task, exist, err := model.GetStorage().GetTaskById(taskId)
	if err != nil {
		goto Error
	}
	if !exist {
		err = errors.New(errors.TASK_NOT_EXIST, fmt.Errorf("task not exist"))
		goto Error
	}

	// valid
	//err = task.ValidAction(typ)
	//if err != nil {
	//	goto Error
	//}

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
	go s.statusLoop()
}

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
	case model.TASK_ACTION_INSPECT:
		err = s.inspect(action.Task)
	case model.TASK_ACTION_COMMIT:
		err = s.commit(action.Task)
	case model.TASK_ACTION_ROLLBACK:
		err = s.rollback(action.Task)
	}
	if err != nil {
		action.Error = err
	}
	s.Lock()
	taskId := fmt.Sprintf(fmt.Sprintf("%d", action.Task.ID))
	delete(s.currentTask, taskId)
	s.Unlock()

	select {
	case action.Done <- struct{}{}:
	default:
	}
	return err
}

func (s *Sqled) inspect(task *model.Task) error {
	fmt.Printf("start inspect task %d\n", task.ID)
	st := model.GetStorage()

	rules, err := st.GetRulesByInstanceId(fmt.Sprintf("%v", task.InstanceId))
	if err != nil {
		return err
	}
	i := inspector.NewInspector(rules, task.Instance, task.CommitSqls, task.Schema)
	err = i.Advise()
	if err != nil {
		return err
	}
	for _, sql := range task.CommitSqls {
		if err := st.Save(&sql); err != nil {
			return err
		}
	}
	return st.UpdateNormalRate(task)
}

func (s *Sqled) commit(task *model.Task) error {
	fmt.Printf("start commit task %d\n", task.ID)
	st := model.GetStorage()

	rules, err := st.GetRulesByInstanceId(fmt.Sprintf("%v", task.InstanceId))
	if err != nil {
		return err
	}
	i := inspector.NewInspector(rules, task.Instance, task.CommitSqls, task.Schema)

	err = i.Prepare()
	if err != nil {
		return err
	}
	sqls, err := i.GenerateRollbackSql()
	if err != nil {
		return err
	}
	rollbackSql := []model.RollbackSql{}
	for _, sql := range sqls {
		rollbackSql = append(rollbackSql, model.RollbackSql{
			Sql: model.Sql{
				Content: sql,
			},
		})
	}
	err = st.UpdateRollbackSql(task, rollbackSql)
	if err != nil {
		return err
	}

	for i.NextCommitSql() {
		sql := i.GetCommitSql()
		st.UpdateCommitSqlStatus(sql, model.TASK_ACTION_DOING, "")
		err := i.Commit()
		st.UpdateCommitSqlStatus(sql, sql.ExecStatus, sql.ExecResult)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Sqled) rollback(task *model.Task) error {
	fmt.Printf("start rollback task %d\n", task.ID)
	st := model.GetStorage()

	// TODO: 1. using transaction for dml; 2. support mycat
	for _, sql := range task.RollbackSqls {
		if sql.Content == "" {
			continue
		}
		err := st.UpdateRollbackSqlStatus(sql, model.TASK_ACTION_DOING, "")
		if err != nil {
			return err
		}
		status := model.TASK_ACTION_DONE
		result := "ok"
		_, err = executor.Exec(task, sql.Content)
		if err != nil {
			status = model.TASK_ACTION_ERROR
			result = err.Error()
		}
		err = st.UpdateRollbackSqlStatus(sql, status, result)
		if err != nil {
			return err
		}
	}
	return nil
}
