package executor

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"sqle/errors"
	"sqle/model"
	"strconv"
)

type Db interface {
	Close()
	Ping() error
	Exec(query string) (driver.Result, error)
	ExecDDL(query, schema, table string) error
	Query(query string, args ...interface{}) ([]map[string]sql.NullString, error)
	Logger() *logrus.Entry
}

type BaseConn struct {
	log  *logrus.Entry
	host string
	port string
	user string
	*gorm.DB
}

func newConn(entry *logrus.Entry, instance *model.Instance, schema string) (*BaseConn, error) {
	var db *gorm.DB
	var err error
	switch instance.DbType {
	case model.DB_TYPE_MYSQL, model.DB_TYPE_MYCAT:
		db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			instance.User, instance.Password, instance.Host, instance.Port, schema))
	case model.DB_TYPE_SQLSERVER:
		db, err = gorm.Open("mssql", fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
			instance.User, instance.Password, instance.Host, instance.Port, schema))

	default:
		err := fmt.Errorf("db type is not support")
		entry.Error(err)
		return nil, errors.New(errors.CONNECT_REMOTE_DB_ERROR, err)
	}
	if err != nil {
		err = fmt.Errorf("connect to %s:%s failed, %s", instance.Host, instance.Port, err)
		entry.Error(err)
		return nil, errors.New(errors.CONNECT_REMOTE_DB_ERROR, err)
	}
	return &BaseConn{
		log:  entry,
		host: instance.Host,
		port: instance.Port,
		user: instance.User,
		DB:   db,
	}, nil
}

func (c *BaseConn) Close() {
	c.DB.Close()
}

func (c *BaseConn) Ping() error {
	return errors.New(errors.CONNECT_REMOTE_DB_ERROR, c.DB.DB().Ping())
}

func (c *BaseConn) Exec(query string) (driver.Result, error) {
	result, err := c.DB.DB().Exec(query)
	if err != nil {
		c.Logger().Errorf("exec sql failed; host: %s, port: %s, user: %s, query: %s, error: %s\n",
			c.host, c.port, c.user, query, err.Error())
	} else {
		c.Logger().Infof("exec sql success; host: %s, port: %s, user: %s, query: %s\n",
			c.host, c.port, c.user, query)
	}
	return result, errors.New(errors.CONNECT_REMOTE_DB_ERROR, err)
}

func (c *BaseConn) ExecDDL(query, schema, table string) error {
	_, err := c.Exec(query)
	return err
}

func (c *BaseConn) Query(query string, args ...interface{}) ([]map[string]sql.NullString, error) {
	rows, err := c.DB.DB().Query(query, args...)
	if err != nil {
		c.Logger().Errorf("query sql failed; host: %s, port: %s, user: %s, query: %s, error: %s\n",
			c.host, c.port, c.user, query, err.Error())
		return nil, errors.New(errors.CONNECT_REMOTE_DB_ERROR, err)
	} else {
		c.Logger().Infof("query sql success; host: %s, port: %s, user: %s, query: %s\n",
			c.host, c.port, c.user, query)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		// unknown error
		c.Logger().Error(err)
		return nil, err
	}
	result := make([]map[string]sql.NullString, 0)
	for rows.Next() {
		buf := make([]interface{}, len(columns))
		data := make([]sql.NullString, len(columns))
		for i := range buf {
			buf[i] = &data[i]
		}
		if err := rows.Scan(buf...); err != nil {
			c.Logger().Error(err)
			return nil, err
		}
		value := make(map[string]sql.NullString, len(columns))
		for i := 0; i < len(columns); i++ {
			k := columns[i]
			v := data[i]
			value[k] = v
		}
		result = append(result, value)
	}
	return result, nil
}

func (c *BaseConn) Logger() *logrus.Entry {
	return c.log
}

type Executor struct {
	Db Db
}

func NewExecutor(entry *logrus.Entry, instance *model.Instance, schema string) (*Executor, error) {
	var executor = &Executor{}
	var conn Db
	var err error
	switch instance.DbType {
	case model.DB_TYPE_MYCAT:
		conn, err = newMycatConn(entry, instance, schema)
	default:
		conn, err = newConn(entry, instance, schema)
	}
	if err != nil {
		return nil, err
	}
	executor.Db = conn
	return executor, nil
}

func Ping(entry *logrus.Entry, instance *model.Instance) error {
	conn, err := NewExecutor(entry, instance, "")
	if err != nil {
		return err
	}
	defer conn.Db.Close()
	return conn.Db.Ping()
}

func ShowDatabases(entry *logrus.Entry, instance *model.Instance) ([]string, error) {
	conn, err := NewExecutor(entry, instance, "")
	if err != nil {
		return nil, err
	}
	defer conn.Db.Close()
	return conn.ShowDatabases()
}

func OpenDbWithTask(entry *logrus.Entry, task *model.Task) (*Executor, error) {
	return NewExecutor(entry, task.Instance, task.Schema)
}

func Exec(entry *logrus.Entry, task *model.Task, sql string) (driver.Result, error) {
	conn, err := OpenDbWithTask(entry, task)
	if err != nil {
		return nil, err
	}
	defer conn.Db.Close()
	return conn.Db.Exec(sql)
}

func (c *Executor) ShowCreateTable(tableName string) (string, error) {
	result, err := c.Db.Query(fmt.Sprintf("show create table %s", tableName))
	if err != nil {
		return "", err
	}
	if len(result) != 1 {
		err := fmt.Errorf("show create table error, result is %v", result)
		c.Db.Logger().Error(err)
		return "", errors.New(errors.CONNECT_REMOTE_DB_ERROR, err)
	}
	if query, ok := result[0]["Create Table"]; !ok {
		err := fmt.Errorf("show create table error, column \"Create Table\" not found")
		c.Db.Logger().Error(err)
		return "", errors.New(errors.CONNECT_REMOTE_DB_ERROR, err)
	} else {
		return query.String, nil
	}
}

func (c *Executor) ShowDatabases() ([]string, error) {
	result, err := c.Db.Query("show databases")
	if err != nil {
		return nil, err
	}
	dbs := make([]string, len(result))
	for n, v := range result {
		if len(v) != 1 {
			err := fmt.Errorf("show databases error, result not match")
			c.Db.Logger().Error(err)
			return dbs, errors.New(errors.CONNECT_REMOTE_DB_ERROR, err)
		}
		for _, db := range v {
			dbs[n] = db.String
			break
		}
	}
	return dbs, nil
}

func (c *Executor) ShowSchemaTables(schema string) ([]string, error) {
	result, err := c.Db.Query(fmt.Sprintf("show tables from %s", schema))
	if err != nil {
		return nil, err
	}
	tables := make([]string, len(result))
	for n, v := range result {
		if len(v) != 1 {
			err := fmt.Errorf("show tables error, result not match")
			c.Db.Logger().Error(err)
			return tables, errors.New(errors.CONNECT_REMOTE_DB_ERROR, err)
		}
		for _, table := range v {
			tables[n] = table.String
			break
		}
	}
	return tables, nil
}

type ExecutionPlanJson struct {
	QueryBlock struct {
		CostInfo struct {
			QueryCost string `json:"query_cost"`
		} `json:"cost_info"`
		TABLE struct {
			Rows int `json:"rows_examined_per_scan"`
		}
	} `json:"query_block"`
}

func (c *Executor) Explain(query string) (ExecutionPlanJson, error) {
	ep := ExecutionPlanJson{}
	result, err := c.Db.Query(fmt.Sprintf("EXPLAIN FORMAT=\"json\" %s", query))
	if err != nil {
		return ep, err
	}
	if len(result) == 1 {
		json.Unmarshal([]byte(result[0]["EXPLAIN"].String), &ep)
	}
	return ep, nil
}

func (c *Executor) ShowMasterStatus() ([]map[string]sql.NullString, error) {
	result, err := c.Db.Query(fmt.Sprintf("show master status"))
	if err != nil {
		return nil, err
	}
	// result may be empty
	if len(result) != 1 && len(result) != 0 {
		err := fmt.Errorf("show master status error, result is %v", result)
		c.Db.Logger().Error(err)
		return nil, errors.New(errors.CONNECT_REMOTE_DB_ERROR, err)
	}
	return result, nil
}

func (c *Executor) FetchMasterBinlogPos() (string, int64, error) {
	result, err := c.ShowMasterStatus()
	if err != nil {
		return "", 0, err
	}
	if len(result) == 0 {
		return "", 0, nil
	}
	file := result[0]["File"].String
	pos, err := strconv.ParseInt(result[0]["Position"].String, 10, 64)
	if err != nil {
		c.Db.Logger().Error(err)
		return "", 0, err
	}
	return file, pos, nil
}

func (c *Executor) ShowTableSizeMB(schema, table string) (float64, error) {
	sql := fmt.Sprintf(`select (DATA_LENGTH + INDEX_LENGTH)/1024/1024 as Size from information_schema.tables 
where table_schema = '%s' and table_name = '%s'`, schema, table)
	result, err := c.Db.Query(sql)
	if err != nil {
		return 0, err
	}
	// table not found, rows = 0
	if len(result) == 0 {
		return 0, nil
	}
	sizeStr := result[0]["Size"].String
	if sizeStr == "" {
		return 0, nil
	}
	size, err := strconv.ParseFloat(sizeStr, 10)
	if err != nil {
		c.Db.Logger().Error(err)
		return 0, errors.New(errors.CONNECT_REMOTE_DB_ERROR, err)
	}
	return size, nil
}
