package sqle

import (
	"actiontech/ucommon/log"
	"fmt"
	"sqle/executor"
	"sqle/inspector"
	"sqle/storage"
	"time"
)

type Sqled struct {
	stage *log.Stage
	// storage
	Storage *storage.Storage
}

func NewSqled(stage *log.Stage) *Sqled {
	return &Sqled{
		stage: stage,
	}
}

func (s *Sqled) Start(exitChan chan struct{}) {
	//go s.TaskLoop(exitChan)
}

func (s *Sqled) TaskLoop(exitChan chan struct{}) {
	t := time.Tick(5 * time.Second)
	for {
		select {
		case <-exitChan:
			return
		case <-t:
		}
		tasks, err := s.Storage.GetTasks()
		if err != nil {
			continue
		}
		for _, task := range tasks {
			currentTask := task
			switch currentTask.Action {
			case storage.TASK_ACTION_INSPECT:
				s.inspect(currentTask)
				s.Storage.UpdateTaskById(fmt.Sprintf("%v", task.ID),
					map[string]interface{}{"action": storage.TASK_ACTION_INIT, "progress": storage.TASK_PROGRESS_INSPECT_END})
			case storage.TASK_ACTION_COMMIT:
				s.commit(currentTask)
				s.Storage.UpdateTaskById(fmt.Sprintf("%v", task.ID),
					map[string]interface{}{"action": storage.TASK_ACTION_INIT, "progress": storage.TASK_PROGRESS_COMMIT_END})
			case storage.TASK_ACTION_ROLLBACK:
				s.rollback(currentTask)
				s.Storage.UpdateTaskById(fmt.Sprintf("%v", task.ID),
					map[string]interface{}{"action": storage.TASK_ACTION_INIT, "progress": storage.TASK_PROGRESS_ROLLACK_END})
			}
		}
	}
}

func (s *Sqled) inspect(task *storage.Task) error {
	sqls, err := inspector.Inspect(nil, task)
	if err != nil {
		return err
	}
	err = s.Storage.UpdateTaskSqls(task, sqls)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sqled) commit(task *storage.Task) error {
	for _, sql := range task.Sqls {
		if sql.CommitSql == "" {
			continue
		}
		// create rollback query
		rollbackQuery, err := inspector.CreateRollbackSql(task, sql.CommitSql)
		if err != nil {
			return err
		}
		fmt.Printf("rollback: %s\n", rollbackQuery)
		sql.RollbackSql = rollbackQuery
		err = executor.Exec(task, sql.CommitSql)
		if err != nil {
			sql.CommitResult = err.Error()
		}
		sql.CommitStatus = "1"
		fmt.Println(sql)
		err = s.Storage.Save(&sql)
	}
	return nil
}

func (s *Sqled) rollback(task *storage.Task) error {
	defer func() {

	}()
	for _, sql := range task.Sqls {
		if sql.RollbackSql == "" {
			continue
		}
		err := executor.Exec(task, sql.RollbackSql)
		if err != nil {
			sql.CommitResult = err.Error()
		}
		sql.RollbackStatus = "1"
		err = s.Storage.Save(&sql)
		if err != nil {
			return err
		}
	}
	return nil
}
