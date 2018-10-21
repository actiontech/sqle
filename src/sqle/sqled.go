package sqle

import (
	"fmt"
	"sqle/storage"
	"sync"
	"time"
)

var sqled *Sqled

func SendAction(task *storage.Task, typ int, callback func(task *storage.Task) error, async bool) error {
	sqled.Lock()
	action := &Action{
		task:      task,
		actionTyp: typ,
	}
	if _, ok := sqled.listen[action.String()]; ok {
		return fmt.Errorf("exist")
	}
	end := make(chan struct{})
	sqled.listen[action.String()] = end
	sqled.queue <- action

	if async {
		return nil
	}
	<-end
	return nil
}

type Action struct {
	task      *storage.Task
	actionTyp int
	callback  func(task *storage.Task) error
}

func (a *Action) String() string {
	return fmt.Sprintf("%d:%d", a.task.ID, a.actionTyp)
}

type Sqled struct {
	sync.Mutex
	exit   chan struct{}
	listen map[string]chan struct{}
	queue  chan *Action
}

func InitSqled(exit chan struct{}) {
	sqled = &Sqled{
		exit:   exit,
		listen: map[string]chan struct{}{},
		queue:  make(chan *Action, 1024),
	}
	sqled.Start()
}

func (s *Sqled) Start() {
	go s.TaskLoop()
}

func (s *Sqled) TaskLoop() {

	t := time.Tick(5 * time.Second)
	for {
		select {
		case <-s.exit:
			return
		case <-t:
		}
		for action := range s.queue {
			currentAction := action
			switch currentAction.actionTyp {
			case storage.TASK_ACTION_INSPECT:
				go s.inspect(currentAction.task)
			case storage.TASK_ACTION_COMMIT:
				go s.commit(currentAction.task)
			case storage.TASK_ACTION_ROLLBACK:
				go s.rollback(currentAction.task)
			}
		}
		//tasks, err := s.Storage.GetTasks()
		//if err != nil {
		//	continue
		//}
		//for _, task := range tasks {
		//	currentTask := task
		//	switch currentTask.Action {
		//	case storage.TASK_ACTION_INSPECT:
		//		s.inspect(currentTask)
		//		s.Storage.UpdateTaskById(fmt.Sprintf("%v", task.ID),
		//			map[string]interface{}{"action": storage.TASK_ACTION_INIT, "progress": storage.TASK_PROGRESS_INSPECT_END})
		//	case storage.TASK_ACTION_COMMIT:
		//		s.commit(currentTask)
		//		s.Storage.UpdateTaskById(fmt.Sprintf("%v", task.ID),
		//			map[string]interface{}{"action": storage.TASK_ACTION_INIT, "progress": storage.TASK_PROGRESS_COMMIT_END})
		//	case storage.TASK_ACTION_ROLLBACK:
		//		s.rollback(currentTask)
		//		s.Storage.UpdateTaskById(fmt.Sprintf("%v", task.ID),
		//			map[string]interface{}{"action": storage.TASK_ACTION_INIT, "progress": storage.TASK_PROGRESS_ROLLACK_END})
		//	}
		//}
	}
}

func (s *Sqled) Do(action Action) error {
	var err error
	switch action.actionTyp {
	case storage.TASK_ACTION_INSPECT:
		err = s.inspect(action.task)
	case storage.TASK_ACTION_COMMIT:
		err = s.commit(action.task)
	case storage.TASK_ACTION_ROLLBACK:
		err = s.rollback(action.task)
	}
	if err != nil {
		return err
	}
	err = action.callback(action.task)
	//s.Lock()
	//_ := s.listen[action.String()]
	return err
}

func (s *Sqled) inspect(task *storage.Task) error {
	//_, err := inspector.Inspect(nil, task)
	//if err != nil {
	//	return err
	//}
	//err = s.Storage.UpdateTaskSqls(task, sqls)
	//if err != nil {
	//	return err
	//}
	//		s.Storage.UpdateTaskById(fmt.Sprintf("%v", task.ID),
	//			map[string]interface{}{"action": storage.TASK_ACTION_INIT, "progress": storage.TASK_PROGRESS_INSPECT_END})
	return nil
}

func (s *Sqled) commit(task *storage.Task) error {
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

func (s *Sqled) rollback(task *storage.Task) error {
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
