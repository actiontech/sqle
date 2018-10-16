package sqle

import (
	"actiontech/ucommon/log"
	"sqle/executor"
	"sqle/inspector"
	"sqle/storage"
	"time"
)

var sqled *Sqled

func GetSqled() *Sqled {
	return sqled
}

type Sqled struct {
	stage *log.Stage
	// storage
	Storage *storage.Storage
}

func InitSqled(stage *log.Stage, s *storage.Storage) {
	sqled = &Sqled{
		stage:   stage,
		Storage: s,
	}
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
				go s.inspect(currentTask)
			case storage.TASK_ACTION_COMMIT:
				go s.commit(currentTask)
			case storage.TASK_ACTION_ROLLBACK:
				go s.rollback(currentTask)
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
	return s.Storage.Update(task, map[string]interface{}{"action": storage.TASK_ACTION_INIT, "progress": storage.TASK_PROGRESS_INSPECT_END})
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
		sql.RollbackSql = rollbackQuery
		err = executor.Exec(task, sql.CommitSql)
		if err != nil {
			sql.CommitResult = err.Error()
		}
		sql.CommitStatus = "1"
		s.Storage.Save(sql)
		if err != nil {
			return err
		}
	}
	return s.Storage.Update(task, map[string]interface{}{"action": storage.TASK_ACTION_INIT, "progress": storage.TASK_PROGRESS_COMMIT_END})
}

func (s *Sqled) rollback(task *storage.Task) error {
	for _, sql := range task.Sqls {
		if sql.RollbackSql == "" {
			continue
		}
		err := executor.Exec(task, sql.RollbackSql)
		if err != nil {
			sql.CommitResult = err.Error()
		}
		sql.RollbackStatus = "1"
		s.Storage.Save(sql)
		if err != nil {
			return err
		}
	}
	return s.Storage.Update(task, map[string]interface{}{"action": storage.TASK_ACTION_INIT, "progress": storage.TASK_PROGRESS_ROLLACK_END})
}
