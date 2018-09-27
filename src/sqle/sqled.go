package sqle

import (
	"actiontech/ucommon/log"
	"sqle/executor"
	"sqle/inspector"
	"sqle/storage"
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

func (s *Sqled) Inspect(task *storage.Task) error {
	sqls, err := inspector.Inspect(task)
	if err != nil {
		return err
	}
	return s.Storage.UpdateTaskSqls(task, sqls)
}

func (s *Sqled) Commit(task *storage.Task) error {
	for _, sql := range task.Sqls {
		err := executor.Query(task, sql.CommitSql)
		if err != nil {
			sql.CommitResult = err.Error()
			s.Storage.Save(sql)
			return nil
		}
	}
	return nil
}

