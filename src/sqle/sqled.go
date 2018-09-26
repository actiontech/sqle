package sqle

import (
	"actiontech/ucommon/log"
	"github.com/jinzhu/gorm"
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
	Db *gorm.DB
}

func InitSqled(stage *log.Stage, db *gorm.DB) {
	sqled = &Sqled{
		stage: stage,
		Db:    db,
	}
}

func (s *Sqled) Exist(model interface{}) (bool, error) {
	var count int
	err := s.Db.Model(model).Where(model).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

//func (s *Sqled) GetTaskById(taskId string) (storage.Task, error){
//	s.Db.
//}

func (s *Sqled) Inspect(task *storage.Task) error {
	sqls, err := inspector.Inspect(task)
	if err != nil {
		return err
	}

	return s.Db.Model(task).Association("Sqls").Replace(sqls).Error
}

func (s *Sqled) ExecSql(task *storage.Task) error {
	for _, sql := range task.Sqls {
		err := executor.Query(task, sql.CommitSql)
		if err != nil {
			sql.CommitResult = err.Error()
			s.Db.Save(sql)
			return nil
		}
	}
	return nil
}
