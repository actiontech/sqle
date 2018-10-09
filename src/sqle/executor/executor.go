package executor

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sqle/storage"
)

type Conn struct {
	*gorm.DB
}

func NewConn(dbType int, user, password, host, port, schema string) (*Conn, error) {
	var db *gorm.DB
	var err error
	switch dbType {
	case storage.DB_TYPE_MYSQL:
		db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			user, password, host, port, schema))
	default:
		return nil, errors.New("db is not support")
	}
	if err != nil {
		return nil, err
	}
	return &Conn{db}, nil
}

func Ping(db *storage.Db) error {
	conn, err := NewConn(db.DbType, db.User, db.Password, db.Host, db.Port, "")
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.ping()
}

func ShowDatabase(db *storage.Db) ([]string, error) {
	conn, err := NewConn(db.DbType, db.User, db.Password, db.Host, db.Port, "")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	result, err := conn.query("show databases")
	if err != nil {
		return nil, err
	}
	dbs := make([]string, 0, len(result))
	for _, v := range result {
		dbName := v["Database"]
		dbs = append(dbs, fmt.Sprintf("%v", dbName))
	}
	return dbs, nil
}

func OpenDbWithTask(task *storage.Task) (*Conn, error) {
	db := task.Db
	schema := task.Schema
	return NewConn(db.DbType, db.User, db.Password, db.Host, db.Port, schema)
}

func Exec(task *storage.Task, sql string) error {
	conn, err := OpenDbWithTask(task)
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.exec(sql)
}

func (c *Conn) ping() error {
	return c.DB.DB().Ping()
}

func (c *Conn) exec(query string) error {
	_, err := c.DB.DB().Exec(query)
	return err
}

func (c *Conn) query(query string) ([]map[string]string, error) {
	rows, err := c.DB.DB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	result := make([]map[string]string, 0)
	for rows.Next() {
		buf := make([]interface{}, len(columns))
		data := make([]sql.NullString, len(columns))
		for i := range buf {
			buf[i] = &data[i]
		}
		if err := rows.Scan(buf...); err != nil {
			return nil, err
		}
		value := make(map[string]string, len(columns))
		for i := 0; i < len(columns); i++ {
			k := columns[i]
			v := data[i].String
			value[k] = v
		}
		result = append(result, value)
	}
	return result, nil
}

func (c *Conn) ShowCreateDatabase(tableName string) (string, error) {
	result, err := c.query(fmt.Sprintf("show create table %s", tableName))
	if err != nil {
		return "", err
	}
	if len(result) != 1 {
		return "", fmt.Errorf("show create table error, result is %v", result)
	}
	if query, ok := result[0]["Create Table"]; !ok {
		return "", fmt.Errorf("show create table error, column \"Create Table\" not found")
	} else {
		return query, nil
	}
}
