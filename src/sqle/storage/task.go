package storage

import (
	"fmt"
)

type Task struct {
	Model
	Name     string
	Desc     string
	Schema   string
	ReqSql   string `json:"sql"`
	Action   int
	Progress int
	Inst     Instance `gorm:"foreignkey:InstId"`
	InstId   uint     `json:"-"`
	Sqls     []Sql    `json:"result" gorm:"foreignkey:TaskId"`
}

type Sql struct {
	Model
	TaskId      uint `json:"-"`
	Number      uint
	CommitSql   string `json:"commit_sql"`
	RollbackSql string `json:"rollback_sql"`
	Inspect     Action `gorm:"foreignkey:InspectId"`
	InspectId   uint   `json:"-"`
	Commit      Action `gorm:"foreignkey:CommitId"`
	CommitId    uint   `json:"-"`
	Rollback    Action `gorm:"foreignkey:RollbackId"`
	RollbackId  uint   `json:"-"`
}

type Action struct {
	Model
	Status bool
	Action string
	Result string
	Error  string
}

// task progress
const (
	TASK_PROGRESS_INIT = iota

	TASK_PROGRESS_INSPECT_START
	TASK_PROGRESS_INSPECT_END

	TASK_PROGRESS_COMMIT_START
	TASK_PROGRESS_COMMIT_END

	TASK_PROGRESS_ROLLACK_START
	TASK_PROGRESS_ROLLACK_END

	TASK_PROGRESS_ERROR
)

// task action
const (
	TASK_ACTION_INIT = iota
	TASK_ACTION_INSPECT
	TASK_ACTION_COMMIT
	TASK_ACTION_ROLLBACK
)

var ActionMap = map[int]string{
	TASK_ACTION_INSPECT:  "",
	TASK_ACTION_COMMIT:   "",
	TASK_ACTION_ROLLBACK: "",
}

func (s *Storage) GetTaskByName(id string) (*Task, error) {
	task := &Task{}
	err := s.db.Preload("Inst").Preload("Sqls").First(&task, id).Error
	return task, err
}

func (s *Storage) GetTasks() ([]Task, error) {
	tasks := []Task{}
	err := s.db.Preload("Inst").Preload("Sqls").Find(&tasks).Error
	return tasks, err
}

func (s *Storage) UpdateTaskSql(task *Task, sqls []*Sql) error {
	return s.db.Model(task).Association("Sqls").Replace(sqls).Error
}

func (s *Storage) UpdateTaskById(taskId string, attrs ...interface{}) error {
	return s.db.Table("tasks").Where("id = ?", taskId).Update(attrs...).Error
}

func (s *Storage) InspectTask(task *Task) error {
	//if task.Action == TASK_ACTION_INSPECT {
	//	return nil
	//}
	//if task.Action != TASK_ACTION_INIT {
	//	return fmt.Errorf("action exist: %s", ActionMap[task.Action])
	//}
	return s.UpdateTaskById(fmt.Sprintf("%v", task.ID), "action", TASK_ACTION_INSPECT)
}

func (s *Storage) CommitTask(task *Task) error {
	//if task.Action == TASK_ACTION_COMMIT {
	//	return nil
	//}
	//if task.Action != TASK_ACTION_INIT {
	//	return fmt.Errorf("action exist: %s", ActionMap[task.Action])
	//}
	//if task.Progress >= TASK_PROGRESS_COMMIT_START {
	//	return errors.New("has commit")
	//}
	return s.UpdateTaskById(fmt.Sprintf("%v", task.ID), "action", TASK_ACTION_COMMIT)
}

func (s *Storage) RollbackTask(task *Task) error {
	//if task.Action == TASK_ACTION_ROLLBACK {
	//	return nil
	//}
	//if task.Action != TASK_ACTION_INIT {
	//	return fmt.Errorf("action exist: %s", ActionMap[task.Action])
	//}
	//if task.Progress != TASK_PROGRESS_COMMIT_END {
	//	return errors.New("not commit")
	//}
	return s.UpdateTaskById(fmt.Sprintf("%v", task.ID), "action", TASK_ACTION_ROLLBACK)
}
