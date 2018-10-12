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
			switch task.Action {
			case storage.TASK_ACTION_INSPECT:
				s.Inspect(task)

			case storage.TASK_ACTION_COMMIT:
				s.Commit(task)
			}
		}
	}
}

func (s *Sqled) Inspect(task *storage.Task) error {
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

func (s *Sqled) Commit(task *storage.Task) error {
	for _, sql := range task.Sqls {
		// create rollback query
		rollbackQuery, err := inspector.CreateRollbackSql(task, sql.CommitSql)
		if err != nil {
			return err
		}
		sql.RollbackSql = rollbackQuery
		err = executor.Exec(task, sql.CommitSql)
		if err != nil {
			sql.CommitResult = err.Error()
			s.Storage.Save(sql)
			return nil
		}
	}
	return s.Storage.Update(task, map[string]interface{}{"action": storage.TASK_ACTION_INIT, "progress": storage.TASK_PROGRESS_COMMIT_END})
}
