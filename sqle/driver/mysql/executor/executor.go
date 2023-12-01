package executor

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

const DAIL_TIMEOUT = 5 * time.Second

type Db interface {
	Close()
	Ping() error
	Exec(query string) (driver.Result, error)
	Transact(qs ...string) ([]driver.Result, error)
	Query(query string, args ...interface{}) ([]map[string]sql.NullString, error)
	QueryWithContext(ctx context.Context, query string, args ...interface{}) (column []string, row [][]sql.NullString, err error)
	Logger() *logrus.Entry
	GetConnectionID() string
}

type BaseConn struct {
	log    *logrus.Entry
	host   string
	port   string
	user   string
	db     *sql.DB
	conn   *sql.Conn
	connID string
}

func newConn(entry *logrus.Entry, instance *driverV2.DSN, schema string) (*BaseConn, error) {
	var db *sql.DB
	var err error

	config := mysql.NewConfig()
	config.User = instance.User
	config.Passwd = instance.Password
	config.Addr = net.JoinHostPort(instance.Host, instance.Port)
	config.DBName = instance.DatabaseName
	config.ParseTime = true
	config.Loc = time.Local
	config.Timeout = DAIL_TIMEOUT
	config.Params = map[string]string{
		"charset": "utf8",
	}
	driver, err := mysql.NewConnector(config)
	if err != nil {
		entry.Error(err)
		return nil, errors.New(errors.ConnectRemoteDatabaseError, err)
	}
	db = sql.OpenDB(driver)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	entry.Infof("connecting to %s:%s with user(%s)", instance.Host, instance.Port, config.User)
	conn, err := db.Conn(context.Background())
	if err != nil {
		entry.Error(err)
		return nil, errors.New(errors.ConnectRemoteDatabaseError, err)
	}
	entry.Infof("connected to %s:%s", instance.Host, instance.Port)

	baseConn := &BaseConn{
		log:  entry,
		host: instance.Host,
		port: instance.Port,
		user: instance.User,
		db:   db,
		conn: conn,
	}
	baseConn.connID, err = baseConn.getConnectionID()
	if err != nil {
		entry.Errorf("get conn id failed, err: %v", err)
		// ignore the error to continue main process
	}
	return baseConn, nil
}

func (c *BaseConn) getConnectionID() (connID string, err error) {
	res, err := c.Query("SELECT connection_id() AS conn_id")
	if err != nil {
		return "", err
	}
	for i := range res {
		row := res[i]
		if row["conn_id"].String != "" {
			return row["conn_id"].String, nil
		}
	}

	return "", nil
}

func (c *BaseConn) Close() {
	c.conn.Close()
	c.db.Close()
}

func (c *BaseConn) GetConnectionID() string {
	return c.connID
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
	return errors.New(errors.ConnectRemoteDatabaseError, err)
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
	return result, errors.New(errors.ConnectRemoteDatabaseError, err)
}

func (c *BaseConn) Transact(qs ...string) ([]driver.Result, error) {
	var err error
	var tx *sql.Tx
	var results []driver.Result
	c.Logger().Infof("doing sql transact, host: %s, port: %s, user: %s", c.host, c.port, c.user)
	tx, err = c.conn.BeginTx(context.Background(), nil)
	if err != nil {
		return results, err
	}
	defer func() {
		if p := recover(); p != nil {
			c.Logger().Error("rollback sql transact")
			if err := tx.Rollback(); err != nil {
				c.Logger().Error("rollback sql transact failed, err:", err)
			}
			panic(p)
		}
		if err != nil {
			c.Logger().Error("rollback sql transact")
			if err := tx.Rollback(); err != nil {
				c.Logger().Error("rollback sql transact failed, err:", err)
			}
			return
		}
		err = tx.Commit()
		if err != nil {
			c.Logger().Error("transact commit failed")
		} else {
			c.Logger().Info("done sql transact")
		}
	}()
	for _, query := range qs {
		var txResult driver.Result
		txResult, err = tx.Exec(query)
		if err != nil {
			c.Logger().Errorf("exec sql failed, error: %s, query: %s", err, query)
			return results, err
		} else {
			results = append(results, txResult)
			c.Logger().Infof("exec sql success, query: %s", query)
		}
	}
	return results, nil
}
func (c *BaseConn) QueryWithContext(ctx context.Context, query string, args ...interface{}) (column []string, row [][]sql.NullString, err error) {
	rows, err := c.conn.QueryContext(ctx, query, args...)
	if err != nil {
		c.Logger().Errorf("query sql failed; host: %s, port: %s, user: %s, query: %s, error: %s\n",
			c.host, c.port, c.user, query, err.Error())
		return nil, nil, errors.New(errors.ConnectRemoteDatabaseError, err)
	} else {
		c.Logger().Infof("query sql success; host: %s, port: %s, user: %s, query: %s\n",
			c.host, c.port, c.user, query)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		// unknown error
		c.Logger().Error(err)
		return nil, nil, err
	}
	result := make([][]sql.NullString, 0)
	for rows.Next() {
		buf := make([]interface{}, len(columns))
		data := make([]sql.NullString, len(columns))
		for i := range buf {
			buf[i] = &data[i]
		}
		if err := rows.Scan(buf...); err != nil {
			c.Logger().Error(err)
			return nil, nil, err
		}
		value := make([]sql.NullString, len(columns))
		for i := 0; i < len(columns); i++ {
			value[i] = data[i]
		}
		result = append(result, value)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, err
	}
	return columns, result, nil
}

func (c *BaseConn) Query(query string, args ...interface{}) ([]map[string]sql.NullString, error) {
	columns, rows, err := c.QueryWithContext(context.TODO(), query, args...)
	if err != nil {
		return nil, err
	}
	result := make([]map[string]sql.NullString, len(rows))
	for j, row := range rows {
		value := make(map[string]sql.NullString)
		for i, s := range row {
			value[columns[i]] = s
		}
		result[j] = value
	}
	return result, nil
}

func (c *BaseConn) Logger() *logrus.Entry {
	return c.log
}

type Executor struct {
	Db                  Db
	lowerCaseTableNames bool
}

func (c *Executor) IsLowerCaseTableNames() bool {
	return c.lowerCaseTableNames
}

func (c *Executor) SetLowerCaseTableNames(lowerCaseTableNames bool) {
	c.lowerCaseTableNames = lowerCaseTableNames
}

func NewExecutor(entry *logrus.Entry, instance *driverV2.DSN, schema string) (*Executor, error) {
	var executor = &Executor{}
	var conn Db
	var err error
	conn, err = newConn(entry, instance, schema)
	if err != nil {
		return nil, err
	}
	executor.Db = conn
	return executor, nil
}

func Ping(entry *logrus.Entry, instance *driverV2.DSN) error {
	conn, err := NewExecutor(entry, instance, "")
	if err != nil {
		return err
	}
	defer conn.Db.Close()
	return conn.Db.Ping()
}

// When using keywords as table names, you need to pay attention to wrapping them in quotation marks
func (c *Executor) ShowCreateTable(schema, tableName string) (string, error) {
	query := fmt.Sprintf("show create table %s", tableName)
	if schema != "" {
		query = fmt.Sprintf("show create table %s.%s", schema, tableName)
	}
	result, err := c.Db.Query(query)
	if err != nil {
		return "", err
	}
	if len(result) != 1 {
		err := fmt.Errorf("show create table error, result is %v", result)
		c.Db.Logger().Error(err)
		return "", errors.New(errors.ConnectRemoteDatabaseError, err)
	}
	if query, ok := result[0]["Create Table"]; !ok {
		err := fmt.Errorf("show create table error, column \"Create Table\" not found")
		c.Db.Logger().Error(err)
		return "", errors.New(errors.ConnectRemoteDatabaseError, err)
	} else {
		return query.String, nil
	}
}

/*
示例：

	mysql> [透传语句]show databases where lower(`Database`) not in ('information_schema','performance_schema','mysql','sys', 'query_rewrite');

				↓包含透传语句时会多出info列
	+----------+------------------+
	| Database | info             |
	+----------+------------------+
	| system   | set_1700620716_1 |
	| test     | set_1700620716_1 |
	+----------+------------------+
*/
func (c *Executor) ShowDatabases(ignoreSysDatabase bool) ([]string, error) {
	var query string

	switch {
	case ignoreSysDatabase && c.IsLowerCaseTableNames():
		query = "show databases where lower(`Database`) not in ('information_schema','performance_schema','mysql','sys', 'query_rewrite')"
	case ignoreSysDatabase && !c.IsLowerCaseTableNames():
		query = "show databases where `Database` not in ('information_schema','performance_schema','mysql','sys', 'query_rewrite')"
	default:
		query = "show databases"
	}

	result, err := c.Db.Query(query)
	if err != nil {
		return nil, err
	}
	dbs := make([]string, len(result))
	for n, v := range result {
		if len(v) < 1 {
			err := fmt.Errorf("show databases error, result not match")
			c.Db.Logger().Error(err)
			return dbs, errors.New(errors.ConnectRemoteDatabaseError, err)
		}
		for key, value := range v {
			if key != "Database" {
				continue
			}
			dbs[n] = value.String
			break
		}
	}
	return dbs, nil
}

/*
示例：

	mysql> [透传语句]select TABLE_NAME from information_schema.tables where table_schema='test' and TABLE_TYPE in ('BASE TABLE','SYSTEM VIEW');
				  ↓包含透传语句时会多出info列
	+------------+------------------+
	| TABLE_NAME | info             |
	+------------+------------------+
	| test_table | set_1700620716_1 |
	+------------+------------------+

	mysql> [透传语句]select TABLE_NAME from information_schema.tables where lower(table_schema) = 'test' and TABLE_TYPE in ('BASE TABLE','SYSTEM VIEW');
				  ↓包含透传语句时会多出info列
	+------------+------------------+
	| TABLE_NAME | info             |
	+------------+------------------+
	| test_table | set_1700620716_1 |
	+------------+------------------+
*/
func (c *Executor) ShowSchemaTables(schema string) ([]string, error) {
	query := fmt.Sprintf(
		"select TABLE_NAME from information_schema.tables where table_schema='%s' and TABLE_TYPE in ('BASE TABLE','SYSTEM VIEW')", schema)

	if c.IsLowerCaseTableNames() {
		schema = strings.ToLower(schema)
		query = fmt.Sprintf(
			"select TABLE_NAME from information_schema.tables where lower(table_schema)='%s' and TABLE_TYPE in ('BASE TABLE','SYSTEM VIEW')", schema)

	}
	result, err := c.Db.Query(query)
	if err != nil {
		return nil, err
	}
	tables := make([]string, len(result))
	for n, v := range result {
		if len(v) < 1 {
			err := fmt.Errorf("show tables error, result not match")
			c.Db.Logger().Error(err)
			return tables, errors.New(errors.ConnectRemoteDatabaseError, err)
		}
		for key, table := range v {
			if key != "TABLE_NAME" {
				continue
			}
			tables[n] = table.String
			break
		}
	}
	return tables, nil
}

/*
示例：

	当使用透传语句时，会多出一列info
	mysql> [透传语句]select TABLE_NAME from information_schema.tables where table_schema='test' and TABLE_TYPE='VIEW';
				  ↓包含透传语句时会多出info列
	+------------+------------------+
	| TABLE_NAME | info             |
	+------------+------------------+
	| test_table | set_1700620716_1 |
	+------------+------------------+

	mysql> [透传语句]select TABLE_NAME from information_schema.tables where lower(table_schema) = 'test' and TABLE_TYPE='VIEW';
				  ↓包含透传语句时会多出info列
	+------------+------------------+
	| TABLE_NAME | info             |
	+------------+------------------+
	| test_table | set_1700620716_1 |
	+------------+------------------+
*/
func (c *Executor) ShowSchemaViews(schema string) ([]string, error) {
	query := fmt.Sprintf(
		"select TABLE_NAME from information_schema.tables where table_schema='%s' and TABLE_TYPE='VIEW'", schema)

	if c.IsLowerCaseTableNames() {
		schema = strings.ToLower(schema)

		query = fmt.Sprintf(
			"select TABLE_NAME from information_schema.tables where lower(table_schema)='%s' and TABLE_TYPE='VIEW'", schema)
	}

	result, err := c.Db.Query(query)
	if err != nil {
		return nil, err
	}
	tables := make([]string, len(result))
	for n, v := range result {
		if len(v) < 1 {
			err := fmt.Errorf("show views error, result not match")
			c.Db.Logger().Error(err)
			return tables, errors.New(errors.ConnectRemoteDatabaseError, err)
		}
		for key, table := range v {
			if key != "TABLE_NAME" {
				continue
			}
			tables[n] = table.String
			break
		}
	}
	return tables, nil
}

// When using keywords as view names, you need to pay attention to wrapping them in quotation marks
func (c *Executor) ShowCreateView(tableName string) (string, error) {
	result, err := c.Db.Query(fmt.Sprintf("show create view %s", tableName))
	if err != nil {
		return "", err
	}
	if len(result) != 1 {
		err := fmt.Errorf("show create view error, result is %v", result)
		c.Db.Logger().Error(err)
		return "", errors.New(errors.ConnectRemoteDatabaseError, err)
	}
	if query, ok := result[0]["Create View"]; !ok {
		err := fmt.Errorf("show create view error, column \"Create View\" not found")
		c.Db.Logger().Error(err)
		return "", errors.New(errors.ConnectRemoteDatabaseError, err)
	} else {
		return query.String, nil
	}
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
	ExplainRecordExtraUsingFilesort         = "Using filesort"
	ExplainRecordExtraUsingTemporary        = "Using temporary"
	ExplainRecordExtraUsingIndexForSkipScan = "Using index for skip scan"
	ExplainRecordExtraUsingWhere            = "Using where"

	ExplainRecordAccessTypeAll   = "ALL"
	ExplainRecordAccessTypeIndex = "index"

	ExplainRecordPrimaryKey = "PRIMARY"
)

func (c *Executor) Explain(query string) (columns []string, rows [][]sql.NullString, err error) {
	columns, rows, err = c.Db.QueryWithContext(context.TODO(), fmt.Sprintf("EXPLAIN %s", query))
	if err != nil {
		return nil, nil, err
	}

	if len(rows) == 0 {
		return nil, nil, fmt.Errorf("no explain record for sql %v", query)
	}

	return columns, rows, nil
}

func (c *Executor) GetExplainRecord(query string) ([]*ExplainRecord, error) {
	columns, rows, err := c.Explain(query)
	if err != nil {
		return nil, err
	}

	records := make([]map[string]sql.NullString, len(rows))
	for j, row := range rows {
		value := make(map[string]sql.NullString)
		for i, s := range row {
			value[columns[i]] = s
		}
		records[j] = value
	}

	ret := make([]*ExplainRecord, len(records))
	for i, record := range records {
		rows, _ := strconv.ParseInt(record["rows"].String, 10, 64)
		ret[i] = &ExplainRecord{
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
		}
	}
	return ret, nil
}

func (c *Executor) ShowMasterStatus() ([]map[string]sql.NullString, error) {
	result, err := c.Db.Query("show master status")
	if err != nil {
		return nil, err
	}
	// result may be empty
	if len(result) != 1 && len(result) != 0 {
		err := fmt.Errorf("show master status error, result is %v", result)
		c.Db.Logger().Error(err)
		return nil, errors.New(errors.ConnectRemoteDatabaseError, err)
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

	if c.IsLowerCaseTableNames() {
		schema = strings.ToLower(schema)
		table = strings.ToLower(table)

		sql = fmt.Sprintf(`select (DATA_LENGTH + INDEX_LENGTH)/1024/1024 as Size from information_schema.tables 
where lower(table_schema) = '%s' and lower(table_name) = '%s'`, schema, table)
	}

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
	size, err := strconv.ParseFloat(sizeStr, 64)
	if err != nil {
		c.Db.Logger().Error(err)
		return 0, errors.New(errors.ConnectRemoteDatabaseError, err)
	}
	return size, nil
}
func (c *Executor) ShowDefaultConfiguration(sql, column string) (string, error) {
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

type TableColumnsInfo struct {
	ColumnName       string
	ColumnType       string
	CharacterSetName string
	IsNullable       string
	ColumnKey        string
	ColumnDefault    string
	Extra            string
	ColumnComment    string
}

func (c *Executor) GetTableColumnsInfo(schema, tableName string) ([]*TableColumnsInfo, error) {
	query := "SELECT COLUMN_NAME, COLUMN_TYPE, CHARACTER_SET_NAME, IS_NULLABLE, COLUMN_KEY, COLUMN_DEFAULT, EXTRA, COLUMN_COMMENT FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=? AND TABLE_NAME=?"

	if c.IsLowerCaseTableNames() {
		schema = strings.ToLower(schema)
		tableName = strings.ToLower(tableName)
		query = "SELECT COLUMN_NAME, COLUMN_TYPE, CHARACTER_SET_NAME, IS_NULLABLE, COLUMN_KEY, COLUMN_DEFAULT, EXTRA, COLUMN_COMMENT FROM INFORMATION_SCHEMA.COLUMNS WHERE lower(TABLE_SCHEMA)=? AND lower(TABLE_NAME)=?"
	}

	records, err := c.Db.Query(query, schema, tableName)
	if err != nil {
		return nil, err
	}

	ret := make([]*TableColumnsInfo, len(records))
	for i, record := range records {
		ret[i] = &TableColumnsInfo{
			ColumnName:       record["COLUMN_NAME"].String,
			ColumnType:       record["COLUMN_TYPE"].String,
			CharacterSetName: record["CHARACTER_SET_NAME"].String,
			IsNullable:       record["IS_NULLABLE"].String,
			ColumnKey:        record["COLUMN_KEY"].String,
			ColumnDefault:    record["COLUMN_DEFAULT"].String,
			Extra:            record["EXTRA"].String,
			ColumnComment:    record["COLUMN_COMMENT"].String,
		}
	}

	return ret, nil
}

type TableIndexesInfo struct {
	ColumnName  string
	KeyName     string
	NonUnique   string
	SeqInIndex  string
	Cardinality string
	Null        string
	IndexType   string
	Comment     string
}

// When using keywords as view names, you need to pay attention to wrapping them in quotation marks
func (c *Executor) GetTableIndexesInfo(schema, tableName string) ([]*TableIndexesInfo, error) {
	records, err := c.Db.Query(fmt.Sprintf("SHOW INDEX FROM %s.%s", schema, tableName))
	if err != nil {
		return nil, err
	}

	ret := make([]*TableIndexesInfo, len(records))
	for i, record := range records {
		ret[i] = &TableIndexesInfo{
			ColumnName:  record["Column_name"].String,
			KeyName:     record["Key_name"].String,
			NonUnique:   record["Non_unique"].String,
			SeqInIndex:  record["Seq_in_index"].String,
			Cardinality: record["Cardinality"].String,
			Null:        record["Null"].String,
			IndexType:   record["Index_type"].String,
			Comment:     record["Comment"].String,
		}
	}
	return ret, nil
}
