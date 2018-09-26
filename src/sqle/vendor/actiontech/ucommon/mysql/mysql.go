package mysql

import (
	"actiontech/ucommon/log"
	"actiontech/ucommon/util"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type DB struct {
	*sql.DB
	sqlExecutionTimeout int
}

type DbMeta struct {
	User           string
	Password       string
	Socket         string
	Host           string
	Port           string
	ConnectTimeout int
	ExecSqlTimeout int
}

func pingDbOrTimeout(db *sql.DB, timeoutSeconds int) error {
	var pingErr error
	fin := <-util.Timeout(func() {
		pingErr = util.Ping(db)
	}, timeoutSeconds)
	if !fin {
		return errors.New("ping db timeout")
	}
	return pingErr
}

func OpenDb(stage *log.Stage, user, pass, socket, host, port string, connectTimeout, execSqlTimeout int) (db *DB, err error) {
	key := DbMeta{
		User:           user,
		Password:       pass,
		Socket:         socket,
		Host:           host,
		Port:           port,
		ConnectTimeout: connectTimeout,
		ExecSqlTimeout: execSqlTimeout,
	}
	return OpenDbWithDbMeta(stage, key)
}

/* workflow of OpenDbWithDbMeta, A B request for same db connection
              +                  +
            req A              req B
              |                  |
              |                  |
+------+------v---------+-----------------+---------------+
       |                |        |        |
       |                |        |      check
       |                |        |     old conn
       |             dial and    v        |
     check            check               v
    old conn         new conn             |
       |                |                 |
       |      or        |       or    dial and
       |                |              check
       |                |             new conn
       |                |                 |
       |                |                 |
+------v------+---------v--------+--------v----------------+
              |                  |
              A                  B
              v                  v
*/
var dbCache = make(map[DbMeta]*DB)
var DbCacheMutex = sync.Mutex{}
var healthCheckers = make(map[DbMeta]chan struct{})
var healthCheckersMutex = sync.Mutex{}
var once = sync.Once{}

func OpenDbWithDbMeta(stage *log.Stage, key DbMeta) (db *DB, err error) {
	once.Do(func() {
		go chCacheRotate()
	})

	ch := startCheckDbConn(stage, key)
	<-ch

	DbCacheMutex.Lock()
	if db, ok := dbCache[key]; ok {
		DbCacheMutex.Unlock()
		return db, nil
	}
	DbCacheMutex.Unlock()

	return nil, errors.New("no pingable connection and dial new connection fail")
}

func startCheckDbConn(stage *log.Stage, key DbMeta) <-chan struct{} {
	healthCheckersMutex.Lock()
	if finishCh, ok := healthCheckers[key]; ok {
		healthCheckersMutex.Unlock()
		return finishCh
	}

	finishCh := make(chan struct{})
	healthCheckers[key] = finishCh
	healthCheckersMutex.Unlock()

	go checkDbConn(stage, key, finishCh)
	return finishCh
}

func checkDbConn(stage *log.Stage, key DbMeta, finishCh chan struct{}) {
	defer func() {
		healthCheckersMutex.Lock()
		delete(healthCheckers, key)
		healthCheckersMutex.Unlock()
		close(finishCh)
	}()

	DbCacheMutex.Lock()
	db, ok := dbCache[key]

	// clean unclean cache, just like a connection use old password
	if !ok {
		for existKey := range dbCache {
			if key.Socket == existKey.Socket && key.Host == existKey.Host && key.Port == existKey.Port && key.User == existKey.User {
				dbCache[existKey].Close()
				log.Detail(stage, "db cache: %+v meta changed to ", existKey, key)
				delete(dbCache, existKey)
			}
		}
	}
	DbCacheMutex.Unlock()

	// check if cached connection pingable
	if ok {
		if err := pingDbOrTimeout(db.DB, 5 /* hard-coding */); nil != err {
			db.Close()
			log.Detail(stage, "db cache: %+v ping error: %v", key, err)

			DbCacheMutex.Lock()
			delete(dbCache, key)
			DbCacheMutex.Unlock()
		}
		return
	}

	// dial a new connection
	var err error
	if "" != key.Socket {
		db, err = OpenDbBySocketWithoutCache(stage, key.User, key.Password, key.Socket, key.ConnectTimeout, key.ExecSqlTimeout)
	} else {
		db, err = OpenDbWithoutCache(stage, key.User, key.Password, key.Host, key.Port, key.ConnectTimeout, key.ExecSqlTimeout)
	}

	if nil != err {
		log.KeyDilute2(stage, "dial_"+key.Host+":"+key.Port,
			"dial %s:%s with %s fail: %v", key.Host, key.Port, key.User, err)
		return
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)

	DbCacheMutex.Lock()
	dbCache[key] = db
	DbCacheMutex.Unlock()
}

// for hash map do not shrink after item deleted, https://github.com/golang/go/issues/20135
func chCacheRotate() {
	for range time.Tick(time.Hour) {
		newCheckers := make(map[DbMeta]chan struct{})

		healthCheckersMutex.Lock()
		for k, v := range healthCheckers {
			newCheckers[k] = v
		}
		healthCheckers = newCheckers
		healthCheckersMutex.Unlock()
	}
}

func OpenDbWithoutCache(stage *log.Stage, user, pass, host, port string, connectTimeout, execHaSqlTimeout int) (db *DB, err error) {
	return OpenDbWithoutCacheWithSchema(stage, user, pass, host, port, "", connectTimeout, execHaSqlTimeout)
}
func OpenDbWithoutCacheWithSchema(stage *log.Stage, user, pass, host, port, schema string, connectTimeout, execHaSqlTimeout int) (db *DB, err error) {
	log.Detail(stage, "host=%v; port=%v; schema=%v; connectTimeout=%v; execHaSqlTimeout=%v", host, port, schema, connectTimeout, execHaSqlTimeout)
	userPass := user
	if "" != pass {
		userPass += ":" + pass
	}
	return openDbByConnStrWithoutCache(stage, fmt.Sprintf("%s@(%s:%s)/%s?timeout=%ds&autocommit=1&multiStatements=true&loc=Local", userPass, host, port, schema, connectTimeout), host+":"+port, execHaSqlTimeout)
}
func OpenDbBySocketWithoutCache(stage *log.Stage, user, pass, socket string, connectTimeout, execHaSqlTimeout int) (db *DB, err error) {
	return OpenDbBySocketWithoutCacheWithSchema(stage, user, pass, socket, "", connectTimeout, execHaSqlTimeout)
}
func OpenDbBySocketWithoutCacheWithSchema(stage *log.Stage, user, pass, socket, schema string, connectTimeout, execHaSqlTimeout int) (db *DB, err error) {
	log.Detail(stage, "socket=%v; schema=%v; connectTimeout=%v; execHaSqlTimeout=%v", socket, schema, connectTimeout, execHaSqlTimeout)
	userPass := user
	if "" != pass {
		userPass += ":" + pass
	}
	return openDbByConnStrWithoutCache(stage, fmt.Sprintf("%s@unix(%s)/%s?timeout=%ds&autocommit=1&multiStatements=true&loc=Local", userPass, socket, schema, connectTimeout), socket, execHaSqlTimeout)
}
func openDbByConnStrWithoutCache(stage *log.Stage, connectCmd, flag string, execHaSqlTimeout int) (db *DB, err error) {
	stage.Enter("SqlOpen")
	defer stage.Exit()

	started := time.Now()
	defer func() {
		if delta := time.Now().Sub(started).Seconds(); delta > 5 {
			log.Brief(stage, "!!! open db %v exceed 5s (actual %vs)", flag, int(delta))
		}
	}()

	sqlDb, err := sql.Open("mysql", connectCmd)
	if nil != err {
		return nil, err
	}
	if err := pingDbOrTimeout(sqlDb, 5 /* hard-coding */); nil != err {
		sqlDb.Close()
		return nil, err
	}
	db = &DB{sqlDb, execHaSqlTimeout}
	return db, nil
}

func SqlQuery(stage *log.Stage, db *DB, s string, arguments ...interface{}) ([]map[string]string, error) {
	stage.Enter("SqlQuery")
	defer stage.Exit()
	result := []map[string]string{}
	mutex := sync.Mutex{} //TRY to fix mystery
	err := sqlQueryForEachRow(stage, db, s, []string{}, func(row map[string]string, hasNonEmptyField bool) bool {
		mutex.Lock()
		result = append(result, row)
		mutex.Unlock()
		return true
	}, arguments...)
	if nil != err {
		return nil, err
	}
	mutex.Lock()
	defer mutex.Unlock()
	return result, nil
}

func SqlExec(stage *log.Stage, db *DB, sql string, args ...interface{}) (err error) {
	stage.Enter("SqlExec")
	defer stage.Exit()

	mutex := sync.Mutex{} //TRY to fix mystery
	log.Detail(stage, "sql=%v; args=%v", sql, args)
	if !<-util.Timeout(func() {
		_, e := db.Exec(sql, args...)
		mutex.Lock()
		err = e
		mutex.Unlock()
	}, db.sqlExecutionTimeout) {
		log.Detail(stage, "sql=%v; args=%v; error (%v); timeout", sql, args, err)
		mutex.Lock()
		defer mutex.Unlock()
		return err
	} else {
		log.Detail(stage, "sql=%v; args=%v; error (%v)", sql, args, err)
		mutex.Lock()
		defer mutex.Unlock()
		return err
	}
}

func SqlQueryRow(stage *log.Stage, db *DB, s string, arguments ...interface{}) (result map[string]string, err error) {
	return SqlQueryRow2(stage, db, s, []string{}, arguments...)
}

func SqlQueryRow2(stage *log.Stage, db *DB, s string, columns []string, arguments ...interface{}) (map[string]string, error) {
	var result map[string]string
	mutex := sync.Mutex{} //TRY to fix mystery
	err := sqlQueryForEachRow(stage, db, s, columns, func(row map[string]string, hasNonEmptyField bool) bool {
		mutex.Lock()
		if !hasNonEmptyField {
			result = nil
		} else {
			result = row
		}
		mutex.Unlock()
		return false
	}, arguments...)
	if nil != err {
		return nil, err
	}
	mutex.Lock()
	defer mutex.Unlock()
	return result, nil
}

func sqlQueryForEachRow(stage *log.Stage, db *DB, s string, columns []string, iter func(row map[string]string, hasNonEmptyField bool) bool, arguments ...interface{}) (err error) {
	log.Detail(stage, "query=%v; args=%v", s, arguments)
	rowsChan := make(chan *sql.Rows, 1)
	errChan := make(chan error, 2)

	if !<-util.Timeout(func() {
		defer func() {
			if r := recover(); r != nil {
				errChan <- errors.New("panic in sql.query")
				return
			}
		}()

		rows, err := db.Query(s, arguments...)
		if nil != err {
			errChan <- err
			return
		}
		rowsChan <- rows
	}, db.sqlExecutionTimeout) {
		err := errors.New("sql-execution-timeout")
		log.Detail(stage, "query=%v; args=%v; error (%v); timeout", s, arguments, err)
		return err
	}

	var rows *sql.Rows
	select {
	case err := <-errChan:
		log.Detail(stage, "query=%v; args=%v; error (%v)", s, arguments, err)
		return err
	case rows = <-rowsChan:
	}

	if nil != rows {
		defer rows.Close()
		cols, _ := rows.Columns()
		buf := make([]interface{}, len(cols))
		data := make([]sql.NullString, len(cols))
		for i := range buf {
			buf[i] = &data[i]
		}

		for rows.Next() {
			rows.Scan(buf...)
			row := make(map[string]string)
			if len(columns) == 0 {
				for _, col := range cols {
					row[col] = ""
				}
			} else {
				for _, col := range columns {
					row[col] = ""
				}
			}
			hasNonEmptyField := false
			for k, col := range data {
				if _, found := row[cols[k]]; found {
					row[cols[k]] = col.String
					if hasNonEmptyField || len(col.String) > 0 {
						hasNonEmptyField = true
					}
				}
			}
			if !iter(row, hasNonEmptyField) {
				break
			}
		}
	}
	return
}
