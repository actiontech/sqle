package sqle

import (
	"fmt"
	"sqle/errors"
	"sqle/executor"
	"sqle/inspector"
	"sqle/model"
	"sync"
	"time"
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
	if _, ok := s.currentTask[taskId]; ok {
		return action, fmt.Errorf("action is exist")
	}
	s.currentTask[taskId] = struct{}{}
	s.Unlock()

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
	st := model.GetStorage()

	rules, err := st.GetRulesByInstanceId(fmt.Sprintf("%v", task.InstanceId))
	if err != nil {
		return err
	}
	i := inspector.NewInspector(model.GetRuleMapFromAllArray(rules), task.Instance, task.CommitSqls, task.Schema)
	sqlArray, err := i.Inspect()
	if err != nil {
		return err
	}
	for _, sql := range sqlArray {
		if err := st.Save(&sql); err != nil {
			return err
		}
	}
	return nil
}

func (s *Sqled) commit(task *model.Task) error {
	st := model.GetStorage()

	rules, err := st.GetRulesByInstanceId(fmt.Sprintf("%v", task.InstanceId))
	if err != nil {
		return err
	}
	i := inspector.NewInspector(model.GetRuleMapFromAllArray(rules), task.Instance, task.CommitSqls, task.Schema)
	sqls, err := i.GenerateRollbackSql()
	if err != nil {
		return err
	}
	rollbackSql := []model.RollbackSql{}
	for _, sql := range sqls {
		rollbackSql = append(rollbackSql, model.RollbackSql{
			Sql: sql,
		})
	}
	err = st.UpdateRollbackSql(task, rollbackSql)
	if err != nil {
		return err
	}
	// TODO: 1. using transaction for dml; 2. support mycat
	for _, sql := range task.CommitSqls {
		if sql.Sql == "" {
			continue
		}
		err := st.UpdateCommitSqlStatus(sql, model.TASK_ACTION_DOING, "")
		if err != nil {
			return err
		}
		status := model.TASK_ACTION_DONE
		result := "ok"
		err = executor.Exec(task, sql.Sql)
		if err != nil {
			status = model.TASK_ACTION_ERROR
			result = err.Error()
		}
		err = st.UpdateCommitSqlStatus(sql, status, result)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Sqled) rollback(task *model.Task) error {
	st := model.GetStorage()

	// TODO: 1. using transaction for dml; 2. support mycat
	for _, sql := range task.RollbackSqls {
		if sql.Sql == "" {
			continue
		}
		err := st.UpdateRollbackSqlStatus(sql, model.TASK_ACTION_DOING, "")
		if err != nil {
			return err
		}
		status := model.TASK_ACTION_DONE
		result := "ok"
		err = executor.Exec(task, sql.Sql)
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

type InstanceStatus struct {
	ID              uint     `json:"id"`
	Name            string   `json:"name"`
	Host            string   `json:"host"`
	Port            string   `json:"port"`
	IsConnectFailed bool     `json:"is_connect_failed"`
	Schemas         []string `json:"schema_list"`
}

func (s *Sqled) statusLoop() {
	tick := time.Tick(1 * time.Hour)
	s.UpdateAllInstanceStatus()
	for {
		select {
		case <-s.exit:
			return
		case <-tick:
			s.UpdateAllInstanceStatus()
		}
	}
}

func (s *Sqled) UpdateAllInstanceStatus() error {
	st := model.GetStorage()
	instances, err := st.GetInstances()
	if err != nil {
		return err
	}
	instancesStatus := map[uint]*InstanceStatus{}
	wait := sync.WaitGroup{}
	mutex := sync.Mutex{}
	for _, instance := range instances {
		wait.Add(1)
		currentInstance := instance
		go func() {
			status := &InstanceStatus{
				ID:   currentInstance.ID,
				Name: currentInstance.Name,
				Host: currentInstance.Host,
				Port: currentInstance.Port,
			}
			schemas, err := executor.ShowDatabases(currentInstance)
			if err != nil {
				status.IsConnectFailed = true
			} else {
				status.Schemas = schemas
			}
			mutex.Lock()
			instancesStatus[currentInstance.ID] = status
			mutex.Unlock()
			wait.Done()
		}()
	}
	wait.Wait()
	s.Lock()
	s.instancesStatus = instancesStatus
	s.Unlock()
	return nil
}

func (s *Sqled) UpdateAndGetInstanceStatus(instance model.Instance) (*InstanceStatus, error) {
	status := &InstanceStatus{
		ID:   instance.ID,
		Name: instance.Name,
		Host: instance.Host,
		Port: instance.Port,
	}
	schemas, err := executor.ShowDatabases(instance)
	if err != nil {
		status.IsConnectFailed = true
	} else {
		status.Schemas = schemas
	}
	s.Lock()
	s.instancesStatus[instance.ID] = status
	s.Unlock()
	return status, err
}

func (s *Sqled) GetAllInstanceStatus() []InstanceStatus {
	statusList := make([]InstanceStatus, 0, len(s.instancesStatus))
	s.Lock()
	for _, status := range s.instancesStatus {
		statusList = append(statusList, *status)
	}
	s.Unlock()
	return statusList
}

func (s *Sqled) DeleteInstanceStatus(instance model.Instance) {
	s.Lock()
	delete(s.instancesStatus, instance.ID)
	s.Unlock()
}
