package executor

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sqle/model"
)

type Conn struct {
	*gorm.DB
}

func NewConn(dbType string, user, password, host, port, schema string) (*Conn, error) {
	var db *gorm.DB
	var err error
	switch dbType {
	case model.DB_TYPE_MYSQL:
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

func Ping(db model.Instance) error {
	conn, err := NewConn(db.DbType, db.User, db.Password, db.Host, db.Port, "")
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.ping()
}

func ShowDatabases(db model.Instance) ([]string, error) {
	conn, err := NewConn(db.DbType, db.User, db.Password, db.Host, db.Port, "")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return conn.ShowDatabases()
}

func OpenDbWithTask(task *model.Task) (*Conn, error) {
	db := task.Instance
	schema := task.Schema
	return NewConn(db.DbType, db.User, db.Password, db.Host, db.Port, schema)
}

func Exec(task *model.Task, sql string) error {
	conn, err := OpenDbWithTask(task)
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Exec(sql)
}

func (c *Conn) ping() error {
	return c.DB.DB().Ping()
}

func (c *Conn) Exec(query string) error {
	_, err := c.DB.DB().Exec(query)
	return err
}

func (c *Conn) Query(query string, args ...interface{}) ([]map[string]string, error) {
	rows, err := c.DB.DB().Query(query, args...)
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

func (c *Conn) ShowCreateTable(tableName string) (string, error) {
	result, err := c.Query(fmt.Sprintf("show create table %s", tableName))
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

func (c *Conn) ShowDatabases() ([]string, error) {
	result, err := c.Query("show databases")
	if err != nil {
		return nil, err
	}
	dbs := make([]string, len(result))
	for n, v := range result {
		dbs[n] = v["Database"]
	}
	return dbs, nil
}

func (c *Conn) ShowSchemaTables(schema string) ([]string, error) {
	result, err := c.Query("select table_name from information_schema.tables where table_schema = ?", schema)
	if err != nil {
		return nil, err
	}
	tables := make([]string, len(result))
	for n, v := range result {
		tables[n] = v["table_name"]
	}
	return tables, nil
}

type ExecutionPlanJson struct {
	QueryBlock struct {
		CostInfo struct {
			QueryCost string `json:"query_cost"`
		} `json:"cost_info"`
	} `json:"query_block"`
}

func (c *Conn) Explain(query string) (ExecutionPlanJson, error) {
	ep := ExecutionPlanJson{}
	result, err := c.Query(fmt.Sprintf("EXPLAIN FORMAT=\"json\" %s", query))
	if err != nil {
		return ep, err
	}
	if len(result) == 1 {
		json.Unmarshal([]byte(result[0]["EXPLAIN"]), &ep)
	}
	return ep, nil
}
