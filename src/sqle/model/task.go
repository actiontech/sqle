package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

// task progress
//const (
//	TASK_PROGRESS_INIT = ""
//
//	TASK_PROGRESS_INSPECT_START = "inspecting"
//	TASK_PROGRESS_INSPECT_END   = "inspected"
//
//	TASK_PROGRESS_COMMIT_START = "committing"
//	TASK_PROGRESS_COMMIT_END   = "committed"
//
//	TASK_PROGRESS_ROLLACK_START = "rolling back"
//	TASK_PROGRESS_ROLLACK_END   = "rolled back"
//
//	TASK_PROGRESS_ERROR = "error"
//)

// task action
const (
	TASK_ACTION_INSPECT = iota
	TASK_ACTION_COMMIT
	TASK_ACTION_ROLLBACK
)

const (
	TASK_ACTION_INIT  = ""
	TASK_ACTION_DOING = "doing"
	TASK_ACTION_DONE  = "done"
	TASK_ACTION_ERROR = "error"
)

var ActionMap = map[int]string{
	TASK_ACTION_INSPECT:  "",
	TASK_ACTION_COMMIT:   "",
	TASK_ACTION_ROLLBACK: "",
}

type CommitSql struct {
	Model
	TaskId        uint   `json:"-"`
	Number        int    `json:"number"`
	Sql           string `json:"sql"`
	InspectStatus string `json:"inspect_status"`
	InspectResult string `json:"inspect_result"`
	InspectLevel  string `json:"inspect_level"`
	ExecStatus    string `json:"exec_status"`
	ExecResult    string `json:"exec_result"`
}

func (s CommitSql) TableName() string {
	return "commit_sql_detail"
}

type RollbackSql struct {
	Model
	TaskId     uint   `json:"-"`
	Number     uint   `json:"number"`
	Sql        string `json:"sql"`
	ExecStatus string `json:"exec_status"`
	ExecResult string `json:"exec_result"`
}

func (s RollbackSql) TableName() string {
	return "rollback_sql_detail"
}

type Task struct {
	Model
	Name         string         `json:"name" example:"REQ201812578"`
	Desc         string         `json:"desc" example:"this is a task"`
	Schema       string         `json:"schema" example:"db1"`
	Instance     Instance       `json:"-" gorm:"foreignkey:InstanceId"`
	InstanceId   uint           `json:"instance_id"`
	Sql          string         `json:"sql"`
	CommitSqls   []*CommitSql   `json:"-" gorm:"foreignkey:TaskId"`
	RollbackSqls []*RollbackSql `json:"-" gorm:"foreignkey:TaskId"`
}

type TaskDetail struct {
	Task
	Instance     Instance      `json:"instance"`
	InstanceId   uint          `json:"-"`
	CommitSqls   []*CommitSql   `json:"commit_sql_list"`
	RollbackSqls []*RollbackSql `json:"rollback_sql_list"`
}

func (t *Task) Detail() TaskDetail {
	data := TaskDetail{
		Task:         *t,
		InstanceId:   t.InstanceId,
		Instance:     t.Instance,
		CommitSqls:   t.CommitSqls,
		RollbackSqls: t.RollbackSqls,
	}
	if t.RollbackSqls == nil {
		data.RollbackSqls = []*RollbackSql{}
	}
	if t.CommitSqls == nil {
		data.CommitSqls = []*CommitSql{}
	}
	return data
}

func (t *Task) ValidAction(typ int) error {
	// inspect sql allowed at all times
	if typ == TASK_ACTION_INSPECT {
		return nil
	}
	// commit is only allowed to commit once
	if typ == TASK_ACTION_COMMIT {
		if t.CommitSqls != nil {
			for _, commitSql := range t.CommitSqls {
				if commitSql.ExecStatus != "" {
					return fmt.Errorf("task has committed")
				}
			}
		}
	}
	// rollback is only allowed to commit once
	// and when commit success
	if typ == TASK_ACTION_ROLLBACK {
		if t.RollbackSqls != nil {
			for _, rollbackSql := range t.RollbackSqls {
				if rollbackSql.ExecStatus != "" {
					return fmt.Errorf("task has rolled back")
				}
			}
			for _, commitSql := range t.CommitSqls {
				if commitSql.ExecStatus != TASK_ACTION_DONE {
					return fmt.Errorf("task has committed")
				}
			}
		}
	}
	return nil
}

func (s *Storage) GetTaskById(taskId string) (*Task, bool, error) {
	task := &Task{}
	err := s.db.Preload("Instance").Preload("CommitSqls").Preload("RollbackSqls").First(&task, taskId).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	return task, true, nil
}

func (s *Storage) GetTasks() ([]Task, error) {
	tasks := []Task{}
	err := s.db.Find(&tasks).Error
	return tasks, err
}

func (s *Storage) UpdateTaskById(taskId string, attrs ...interface{}) error {
	return s.db.Table("tasks").Where("id = ?", taskId).Update(attrs...).Error
}

func (s *Storage) UpdateCommitSql(task *Task, commitSql []CommitSql) error {
	return s.db.Model(task).Association("CommitSqls").Replace(commitSql).Error
}

func (s *Storage) UpdateRollbackSql(task *Task, rollbackSql []*RollbackSql) error {
	return s.db.Model(task).Association("RollbackSqls").Replace(rollbackSql).Error
}

func (s *Storage) UpdateProgress(task *Task, progress string) error {
	return s.UpdateTaskById(fmt.Sprintf("%v", task.ID), map[string]interface{}{"progress": progress})
}
