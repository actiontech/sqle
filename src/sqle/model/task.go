package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type CommitSql struct {
	Model
	TaskId        uint   `json:"-"`
	Number        uint   `json:"number"`
	Sql           string `json:"sql"`
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
	Name         string        `json:"name" example:"REQ201812578"`
	Desc         string        `json:"desc" example:"this is a task"`
	Schema       string        `json:"schema" example:"db1"`
	Instance     Instance      `json:"_" gorm:"foreignkey:InstanceId"`
	InstanceId   uint          `json:"instance_id"`
	Sql          string        `json:"sql"`
	Progress     string        `json:"progress"`
	CommitSqls   []CommitSql   `json:"-" gorm:"foreignkey:TaskId"`
	RollbackSqls []RollbackSql `json:"-" gorm:"foreignkey:TaskId"`
}

type TaskDetail struct {
	Task
	Instance     Instance      `json:"instance"`
	InstanceId   uint          `json:"-"`
	CommitSqls   []CommitSql   `json:"commit_sql_list"`
	RollbackSqls []RollbackSql `json:"rollback_sql_list"`
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
		data.RollbackSqls = []RollbackSql{}
	}
	if t.CommitSqls == nil {
		data.CommitSqls = []CommitSql{}
	}
	return data
}

//type Sql struct {
//	Model
//	Sql          string        `json:"sql"`
//	Progress     string        `json:"progress"`
//	CommitSqls   []CommitSql   `json:"-" gorm:"foreignkey:SqlId"`
//	RollbackSqls []RollbackSql `json:"-" gorm:"foreignkey:SqlId"`
//}
//
//func (s Sql) TableName() string {
//	return "sqls"
//}

// task progress
const (
	TASK_PROGRESS_INIT = ""

	TASK_PROGRESS_INSPECT_START = "inspecting"
	TASK_PROGRESS_INSPECT_END   = "inspected"

	TASK_PROGRESS_COMMIT_START = "committing"
	TASK_PROGRESS_COMMIT_END   = "committed"

	TASK_PROGRESS_ROLLACK_START = "rolling back"
	TASK_PROGRESS_ROLLACK_END   = "rolled back"

	TASK_PROGRESS_ERROR = "error"
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

//func (s *Storage) GetSqlById(SqlId string) (*Sql, bool, error) {
//	sql := &Sql{}
//	err := s.db.Preload("CommitSqls").Preload("RollbackSqls").First(&sql, SqlId).Error
//	if err == gorm.ErrRecordNotFound {
//		return nil, false, nil
//	}
//	return sql, true, err
//}

func (s *Storage) UpdateTaskById(taskId string, attrs ...interface{}) error {
	return s.db.Table("tasks").Where("id = ?", taskId).Update(attrs...).Error
}

//func (s *Storage) UpdateSqlsById(sqlsId uint, attrs ...interface{}) error {
//	return s.db.Table(Sql{}.TableName()).Where("id = ?", sqlsId).Update(attrs...).Error
//}

func (s *Storage) UpdateCommitSql(task *Task, commitSql []*CommitSql) error {
	return s.db.Model(task).Association("CommitSqls").Replace(commitSql).Error
}

func (s *Storage) UpdateRollbackSql(task *Task, rollbackSql []*RollbackSql) error {
	return s.db.Model(task).Association("RollbackSqls").Replace(rollbackSql).Error
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

func (s *Storage) UpdateProgress(task *Task, progress string) error {
	return s.UpdateTaskById(fmt.Sprintf("%v", task.ID), map[string]interface{}{"progress": progress})
}
