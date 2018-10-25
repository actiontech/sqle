package sqle

import (
	"fmt"
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
	exit        chan struct{}
	currentTask map[string]struct{}
	queue       chan *Action
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
		err = fmt.Errorf("task not exist")
		goto Error
	}

	// valid
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
	go s.TaskLoop()
}

func (s *Sqled) TaskLoop() {
	for {
		select {
		case <-s.exit:
			return
		case action := <-s.queue:
			go s.Do(action)
		}
	}
}

func (s *Sqled) Do(action *Action) error {
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

	i := inspector.NewInspector(nil, task.Instance, task.CommitSqls, task.Schema)
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
	//for _, sql := range task.Sqls {
	//	if sql.CommitSql == "" {
	//		continue
	//	}
	//	// create rollback query
	//	rollbackQuery, err := inspector.CreateRollbackSql(task, sql.CommitSql)
	//	if err != nil {
	//		return err
	//	}
	//	fmt.Printf("rollback: %s\n", rollbackQuery)
	//	sql.RollbackSql = rollbackQuery
	//	err = executor.Exec(task, sql.CommitSql)
	//	//if err != nil {
	//	//	sql.CommitResult = err.Error()
	//	//}
	//	//sql.CommitStatus = "1"
	//	//fmt.Println(sql)
	//	//err = s.Storage.Save(&sql)
	//
	//}
	return nil
}

func (s *Sqled) rollback(task *model.Task) error {
	//defer func() {
	//
	//}()
	//for _, sql := range task.Sqls {
	//	if sql.RollbackSql == "" {
	//		continue
	//	}
	//	err := executor.Exec(task, sql.RollbackSql)
	//	//if err != nil {
	//	//	sql.CommitResult = err.Error()
	//	//}
	//	//sql.RollbackStatus = "1"
	//	//err = s.Storage.Save(&sql)
	//	if err != nil {
	//		return err
	//	}
	//}
	return nil
}
