package executor

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sqle/storage"
)

const CONNECT_TIMEOUT = 5

//func query(meta *storage.TaskMeta, sql string, execTimeout int) ([]map[string]string, error) {
//	stage := log.NewStage().Enter("mysql_query")
//	db, err := mysql.OpenDbWithoutCacheWithSchema(stage,
//		meta.Database.User, meta.Database.Password, meta.Database.Host, meta.Database.Password,
//		"", CONNECT_TIMEOUT, execTimeout)
//	if nil != err {
//		return nil, err
//	}
//	defer db.Close()
//
//	return mysql.SqlQuery(stage, db, sql)
//}
//
//func commitTask(meta *storage.TaskMeta) error {
//	for _, sqlMeta := range meta.SqlMetas {
//		_, err := query(meta, sqlMeta.ExecSql, 5)
//		if nil != err {
//			sqlMeta.CommitResult = err.Error()
//			return err
//		}
//	}
//	return nil
//}

func openDbWithMeta(db *storage.Db) (*gorm.DB, error) {
	switch db.DbType {
	case storage.DB_TYPE_MYSQL:
		return gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			db.User, db.Password, db.Host, db.Port, ""))
	default:
		return nil, errors.New("db is not support")
	}
}

func Query(task *storage.Task, sql string) error {
	db, err := openDbWithMeta(&task.Db)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.DB().Exec(sql)
	return err
}

func Ping(task *storage.Task) error {
	db, err := openDbWithMeta(&task.Db)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.DB().Ping()
}
