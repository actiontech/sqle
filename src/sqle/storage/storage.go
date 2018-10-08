package storage

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

const (
	DB_TYPE_MYSQL = iota
	DB_TYPE_MYCAT
	DB_TYPE_SQLSERVER
)

type User struct {
	gorm.Model
	Name     string
	Password string
}

type Db struct {
	gorm.Model
	DbType   int
	Alias    string
	Host     string
	Port     string
	User     string
	Password string
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

func (t *Task) canInspect() bool {
	if t.Action != TASK_ACTION_INIT {
		return false
	}
	return true
}

func (t *Task) canCommit() bool {
	if t.Action != TASK_ACTION_INIT {
		return false
	}
	if !t.Approved {
		return false
	}
	if t.Progress >= TASK_PROGRESS_COMMIT_START {
		return false
	}
	return true
}

func (t *Task) canRollback() bool {
	if t.Action != TASK_ACTION_INIT {
		return false
	}
	if !t.Approved {
		return false
	}
	if t.Progress != TASK_PROGRESS_COMMIT_END {
		return false
	}
	return true
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

type InspectConfig struct {
	gorm.Model
	Code       int
	ConfigType int
	StmtType   int
	Variable   string
	Desc       string
	Level      int
	Disable    bool
}

func NewMysql(user, password, host, port, schema string) (*Storage, error) {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user, password, host, port, schema))
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	if err := createTable(db, &User{}); err != nil {
		return nil, err
	}
	if err := createTable(db, &Db{}); err != nil {
		return nil, err
	}
	if err := createTable(db, &Sql{}); err != nil {
		return nil, err
	}
	if err := createTable(db, &Task{}); err != nil {
		return nil, err
	}
	return &Storage{
		db: db,
	}, nil
}

func createTable(db *gorm.DB, model interface{}) error {
	hasTable := db.HasTable(model)
	if db.Error != nil {
		return db.Error
	}
	if !hasTable {
		return db.CreateTable(model).Error
	}
	return nil
}

type Storage struct {
	db *gorm.DB
}

func (s *Storage) Exist(model interface{}) (bool, error) {
	var count int
	err := s.db.Model(model).Where(model).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *Storage) Create(model interface{}) error {
	return s.db.Create(model).Error
}

func (s *Storage) Save(model interface{}) error {
	return s.db.Save(model).Error
}

func (s *Storage) GetUserById(id string) (*User, error) {
	user := &User{}
	err := s.db.First(user, id).Error
	return user, err
}

func (s *Storage) UpdateUser(user *User) error {
	return s.db.Save(user).Error
}

func (s *Storage) DelUser(user *User) error {
	return s.db.Delete(user).Error
}

func (s *Storage) GetUsers() ([]*User, error) {
	users := []*User{}
	err := s.db.Find(users).Error
	return users, err
}

func (s *Storage) GetDatabaseById(id string) (*Db, error) {
	database := &Db{}
	err := s.db.First(database).Error
	return database, err
}

func (s *Storage) UpdateDatabase(database *Db) error {
	return s.db.Save(database).Error
}

func (s *Storage) DelDatabase(database *Db) error {
	return s.db.Delete(database).Error
}

func (s *Storage) GetDatabases() ([]*Db, error) {
	databases := []*Db{}
	err := s.db.Find(databases).Error
	return databases, err
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
	return s.db.Model(task).UpdateColumns(Task{Action: TASK_ACTION_INSPECT}).Error
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
	return s.db.Model(task).UpdateColumns(Task{Action: TASK_ACTION_COMMIT}).Error
}
