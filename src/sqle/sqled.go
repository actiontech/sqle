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
	currentTask map[string]chan *Action
	queue       chan *Action
}

type Action struct {
	Task  *model.Task
	Typ   int
	Error error
}

func InitSqled(exit chan struct{}) {
	sqled = &Sqled{
		exit:        exit,
		currentTask: map[string]chan *Action{},
		queue:       make(chan *Action, 1024),
	}
	sqled.Start()
}

func (s *Sqled) AddTask(taskId string, typ int) (chan *Action, error) {
	s.Lock()
	if _, ok := s.currentTask[taskId]; ok {
		return nil, fmt.Errorf("action is exist")
	}
	done := make(chan *Action)
	s.currentTask[taskId] = done
	s.Unlock()

	task, exist, err := model.GetStorage().GetTaskById(taskId)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("task not exist")
	}
	action := &Action{
		Task: task,
		Typ:  typ,
	}
	s.queue <- action
	return done, nil
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
	done, ok := s.currentTask[taskId]
	delete(s.currentTask, taskId)
	s.Unlock()
	if ok {
		select {
		case done <- action:
		default:
		}
		close(done)
	}
	return err
}

func (s *Sqled) inspect(task *model.Task) error {
	st := model.GetStorage()
	err := st.UpdateProgress(task, model.TASK_PROGRESS_INSPECT_START)
	if err != nil {
		return err
	}
	i := inspector.NewInspector(nil, task.Instance, task.Schema, task.Sql)
	sqls, err := i.Inspect()
	if err != nil {
		return err
	}
	fmt.Println("1111")
	err = st.UpdateCommitSql(task, sqls)
	if err != nil {
		return err
	}
	fmt.Println("2222")
	return st.UpdateProgress(task, model.TASK_PROGRESS_INSPECT_END)
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
