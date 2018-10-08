package executor

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sqle/storage"
)

func openDb(dbType int, user, password, host, port, schema string) (*gorm.DB, error) {
	switch dbType {
	case storage.DB_TYPE_MYSQL:
		return gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			user, password, host, port, schema))
	default:
		return nil, errors.New("db is not support")
	}
}

func Ping(db *storage.Db) error {
	conn, err := openDb(db.DbType, db.User, db.Password, db.Host, db.Port, "")
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.DB().Ping()
}

func OpenDbWithTask(task *storage.Task) (*gorm.DB, error) {
	db := task.Db
	schema := task.Schema
	return openDb(db.DbType, db.User, db.Password, db.Host, db.Port, schema)
}

func Exec(task *storage.Task, sql string) error {
	conn, err := OpenDbWithTask(task)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.DB().Exec(sql)
	return err
}
