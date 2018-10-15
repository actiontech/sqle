package storage

import (
	"github.com/jinzhu/gorm"
	"fmt"
	"errors"
)

type Task struct {
	gorm.Model
	// meta
	User       User `gorm:"foreignkey:UserId"`
	UserId     int
	Db         Db `gorm:"foreignkey:DbId"`
	DbId       int
	Approver   User `gorm:"foreignkey:ApproverId"`
	ApproverId int
	Schema     string
	ReqSql     string

	// status
	Approved bool
	Action   int
	Progress int
	Sqls     []Sql `gorm:"foreignkey:TaskId"`
}

type Sql struct {
	gorm.Model
	TaskId         int
	CommitSql      string
	RollbackSql    string
	InspectError   string
	InspectWarn    string
	CommitStatus   string
	CommitResult   string
	RollbackStatus string
	RollbackResult string
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

func (s *Storage) GetTaskById(id string) (*Task, error) {
	task := &Task{}
	err := s.db.Preload("User").Preload("Approver").Preload("Db").Preload("Sqls").First(&task, id).Error
	return task, err
}

func (s *Storage) GetTasks() ([]*Task, error) {
	tasks := []*Task{}
	err := s.db.Preload("User").Preload("Approver").Preload("Db").Preload("Sqls").Find(&tasks).Error
	return tasks, err
}

func (s *Storage) UpdateTaskSqls(task *Task, sqls []*Sql) error {
	return s.db.Model(task).Association("Sqls").Replace(sqls).Error
}

func (s *Storage) InspectTask(taskId string) error {
	task, err := s.GetTaskById(taskId)
	if err != nil {
		return err
	}
	if task.Action == TASK_ACTION_INSPECT {
		return nil
	}
	if task.Action != TASK_ACTION_INIT {
		return fmt.Errorf("action exist: %s", ActionMap[task.Action])
	}
	return s.db.Model(task).Update("action", TASK_ACTION_INSPECT).Error
}

func (s *Storage) CommitTask(taskId string) error {
	task, err := s.GetTaskById(taskId)
	if err != nil {
		return err
	}
	if task.Action == TASK_ACTION_COMMIT {
		return nil
	}
	if task.Action != TASK_ACTION_INIT {
		return fmt.Errorf("action exist: %s", ActionMap[task.Action])
	}
	if !task.Approved {
		return errors.New("not approve")
	}
	if task.Progress >= TASK_PROGRESS_COMMIT_START {
		return errors.New("has commit")
	}
	return s.db.Model(task).Update("action", TASK_ACTION_COMMIT).Error
}

func (s *Storage) RollbackTask(taskId string) error {
	task, err := s.GetTaskById(taskId)
	if err != nil {
		return err
	}
	if task.Action == TASK_ACTION_ROLLBACK {
		return nil
	}
	if task.Action != TASK_ACTION_INIT {
		return fmt.Errorf("action exist: %s", ActionMap[task.Action])
	}
	if !task.Approved {
		return errors.New("not approve")
	}
	if task.Progress != TASK_PROGRESS_COMMIT_END {
		return errors.New("not commit")
	}
	return s.db.Model(task).Update("action", TASK_ACTION_ROLLBACK).Error
}
