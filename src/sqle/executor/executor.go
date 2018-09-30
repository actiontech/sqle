package executor

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sqle/storage"
)


func OpenDbWithMeta(db *storage.Db) (*gorm.DB, error) {
	switch db.DbType {
	case storage.DB_TYPE_MYSQL:
		return gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			db.User, db.Password, db.Host, db.Port, ""))
	default:
		return nil, errors.New("db is not support")
	}
}

func Query(task *storage.Task, sql string) error {
	db, err := OpenDbWithMeta(&task.Db)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.DB().Exec(sql)
	return err
}

func Ping(database *storage.Db) error {
	db, err := OpenDbWithMeta(database)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.DB().Ping()
}

//func ShowDatabases(database *storage.Db) ([]string,error){
//	db, err := openDbWithMeta(database)
//	if err != nil {
//		return nil, err
//	}
//	defer db.Close()
//	return db.DB().QueryRow("show databases")
//}