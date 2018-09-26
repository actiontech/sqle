package storage

import (
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

type Task struct {
	gorm.Model
	User       User `gorm:"foreignkey:UserId"`
	UserId     int
	Approver   User `gorm:"foreignkey:ApproverId"`
	ApproverId int
	Db         Db `gorm:"foreignkey:DbId"`
	DbId       int
	ReqSql     string
	Sqls       []Sql `gorm:"foreignkey:TaskId"`
}

type Sql struct {
	gorm.Model
	TaskId         string
	CommitSql      string
	RollbackSql    string
	InspectResult  string
	CommitResult   string
	RollbackResult string
}

func NewMysql(user, password, host, port, schema string) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user, password, host, port, schema))
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	_, err = db.DB().Exec(fmt.Sprintf("create database if not exists %s", schema))
	if err != nil {
		fmt.Println("create database error")
		return nil, err
	}
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
	return db, nil
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
