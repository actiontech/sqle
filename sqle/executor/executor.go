package executor

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"actiontech.cloud/sqle/sqle/sqle/errors"
	"actiontech.cloud/sqle/sqle/sqle/model"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

const DAIL_TIMEOUT = 5 * time.Second

type Db interface {
	Close()
	Ping() error
	Exec(query string) (driver.Result, error)
	Transact(qs ...string) ([]driver.Result, map[int]string, error)
	ExecDDL(query, schema, table string) error
	Query(query string, args ...interface{}) ([]map[string]sql.NullString, error)
	Logger() *logrus.Entry
}

type BaseConn struct {
	log  *logrus.Entry
	host string
	port string
	user string
	db   *sql.DB
	conn *sql.Conn
}

func newConn(entry *logrus.Entry, instance *model.Instance, schema string) (*BaseConn, error) {
	var db *sql.DB
	var err error
	switch instance.DbType {
	case model.DB_TYPE_MYSQL, model.DB_TYPE_MYCAT:
		db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=%s&charset=utf8&parseTime=True&loc=Local",
			instance.User, instance.Password, instance.Host, instance.Port, schema, DAIL_TIMEOUT))
	case model.DB_TYPE_SQLSERVER:
		query := url.Values{}
		query.Add("database", schema)
		query.Add("dial timeout", fmt.Sprintf("%.0f", DAIL_TIMEOUT.Seconds()))
		source := &url.URL{
			Scheme:   "sqlserver",
			User:     url.UserPassword(instance.User, instance.Password),
			Host:     fmt.Sprintf("%s:%s", instance.Host, instance.Port),
			RawQuery: query.Encode(),
		}
		db, err = sql.Open("mssql", source.String())

	default:
		err := fmt.Errorf("db type is not support")
		entry.Error(err)
		return nil, errors.New(errors.CONNECT_REMOTE_DB_ERROR, err)
	}
	if err != nil {
		entry.Error(err)
		return nil, errors.New(errors.CONNECT_REMOTE_DB_ERROR, err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	entry.Infof("connecting to %s %s:%s", instance.DbType, instance.Host, instance.Port)
	conn, err := db.Conn(context.Background())
	if err != nil {
		entry.Error(err)
		return nil, errors.New(errors.CONNECT_REMOTE_DB_ERROR, err)
	}
	entry.Infof("connected to %s %s:%s", instance.DbType, instance.Host, instance.Port)
	return &BaseConn{
		log:  entry,
		host: instance.Host,
		port: instance.Port,
		user: instance.User,
		db:   db,
		conn: conn,
	}, nil
}

func (c *BaseConn) Close() {
	c.conn.Close()
	c.db.Close()
}

func (c *BaseConn) Ping() error {
	c.Logger().Infof("ping %s:%s", c.host, c.port)
	ctx, cancel := context.WithTimeout(context.Background(), DAIL_TIMEOUT)
	defer cancel()
	err := c.conn.PingContext(ctx)
	if err != nil {
		c.Logger().Infof("ping %s:%s failed, %s", c.host, c.port, err)
	} else {
		c.Logger().Infof("ping %s:%s success", c.host, c.port)
	}
	return errors.New(errors.CONNECT_REMOTE_DB_ERROR, err)
}

func (c *BaseConn) Exec(query string) (driver.Result, error) {
	result, err := c.conn.ExecContext(context.Background(), query)
	if err != nil {
		c.Logger().Errorf("exec sql failed; host: %s, port: %s, user: %s, query: %s, error: %s",
			c.host, c.port, c.user, query, err.Error())
	} else {
		c.Logger().Infof("exec sql success; host: %s, port: %s, user: %s, query: %s",
			c.host, c.port, c.user, query)
	}
	return result, errors.New(errors.CONNECT_REMOTE_DB_ERROR, err)
}

func (c *BaseConn) Transact(qs ...string) ([]driver.Result, map[int] /*sql index*/ string /*exec result*/, error) {
	var err error
	var tx *sql.Tx
	var results []driver.Result
	qsExecResultMap := make(map[int] /*sql index*/ string /*exec result*/)
	c.Logger().Infof("doing sql transact, host: %s, port: %s, user: %s", c.host, c.port, c.user)
	tx, err = c.conn.BeginTx(context.Background(), nil)
	if err != nil {
		return results, qsExecResultMap, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			c.Logger().Error("rollback sql transact")
			panic(p)
		}
		if err != nil {
			tx.Rollback()
			c.Logger().Error("rollback sql transact")
			return
		}
		err = tx.Commit()
		if err != nil {
			c.Logger().Error("transact commit failed")
		} else {
			c.Logger().Info("done sql transact")
		}
	}()
	for index, query := range qs {
		var txResult driver.Result
		txResult, err = tx.Exec(query)
		if err != nil {
			qsExecResultMap[index] = err.Error()
			c.Logger().Errorf("exec sql failed, error: %s, query: %s", err, query)
			return results, qsExecResultMap, nil
		} else {
			qsExecResultMap[index] = "ok"
			results = append(results, txResult)
			c.Logger().Infof("exec sql success, query: %s", query)
		}
	}
	return results, qsExecResultMap, nil
}

func (c *BaseConn) ExecDDL(query, schema, table string) error {
	_, err := c.Exec(query)
	return err
}

func (c *BaseConn) Query(query string, args ...interface{}) ([]map[string]sql.NullString, error) {
	rows, err := c.conn.QueryContext(context.Background(), query, args...)
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
	dbType string
	Db     Db
}

func NewExecutor(entry *logrus.Entry, instance *model.Instance, schema string) (*Executor, error) {
	var executor = &Executor{
		dbType: instance.DbType,
	}
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
	return conn.ShowDatabases(true)
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

func (c *Executor) ShowDatabases(ignoreSysDatabase bool) ([]string, error) {
	var query string
	if c.dbType == model.DB_TYPE_SQLSERVER {
		if ignoreSysDatabase {
			query = "select name from sys.databases where name not in ('master','tempdb','model','msdb','distribution')"
		} else {
			query = "select name from sys.databases"
		}
	} else {
		if ignoreSysDatabase {
			query = "show databases where `Database` not in ('information_schema','performance_schema','mysql','sys')"
		} else {
			query = "show databases"
		}
	}
	result, err := c.Db.Query(query)
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

type ExplainRecord struct {
	Id           string `json:"id"`
	SelectType   string `json:"select_type"`
	Table        string `json:"table"`
	Partitions   string `json:"partitions"`
	Type         string `json:"type"`
	PossibleKeys string `json:"possible_keys"`
	Key          string `json:"key"`
	KeyLen       string `json:"key_len"`
	Ref          string `json:"ref"`
	Rows         int64  `json:"rows"`
	Filtered     string `json:"filtered"`
	Extra        string `json:"extra"`
}

// https://dev.mysql.com/doc/refman/5.7/en/explain-output.html#explain_rows
const (
	ExplainRecordExtraUsingFilesort  = "Using filesort"
	ExplainRecordExtraUsingTemporary = "Using temporary"

	ExplainRecordAccessTypeAll = "ALL"
)

func (c *Executor) Explain(query string) ([]*ExplainRecord, error) {
	records, err := c.Db.Query(fmt.Sprintf("EXPLAIN %s", query))
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no explain record for sql %v", query)
	}

	var ret []*ExplainRecord
	for _, record := range records {
		rows, _ := strconv.ParseInt(record["rows"].String, 10, 64)
		ret = append(ret, &ExplainRecord{
			Id:           record["id"].String,
			SelectType:   record["select_type"].String,
			Table:        record["table"].String,
			Partitions:   record["partitions"].String,
			Type:         record["type"].String,
			PossibleKeys: record["possible_keys"].String,
			Key:          record["key"].String,
			KeyLen:       record["key_len"].String,
			Ref:          record["ref"].String,
			Rows:         rows,
			Filtered:     record["filtered"].String,
			Extra:        record["Extra"].String,
		})
	}
	return ret, nil
}

func (c *Executor) ShowMasterStatus() ([]map[string]sql.NullString, error) {
	if c.dbType == model.DB_TYPE_SQLSERVER {
		return []map[string]sql.NullString{}, nil
	}
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
	if c.dbType == model.DB_TYPE_SQLSERVER {
		return "", 0, nil
	}
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
	if c.dbType == model.DB_TYPE_SQLSERVER {
		return 0, nil
	}
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
func (c *Executor) ShowDefaultConfiguration(sql, column string) (string, error) {
	if c.dbType != model.DB_TYPE_MYSQL {
		return "", nil
	}
	result, err := c.Db.Query(sql)
	if err != nil {
		return "", err
	}
	// table not found, rows = 0
	if len(result) == 0 {
		return "", nil
	}
	ret, ok := result[0][column]
	if !ok {
		return "", nil
	}
	return ret.String, nil
}
