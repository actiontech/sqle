package driver

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"actiontech.cloud/sqle/sqle/sqle/model"

	"github.com/sirupsen/logrus"
)

var (
	drivers   = make(map[string]handler)
	driversMu sync.RWMutex

	rules   []*model.Rule
	rulesMu sync.RWMutex
)

// handler is a template which Driver plugin should provide such function signature.
type handler func(log *logrus.Entry, inst *model.Instance, schema string) Driver

// Register like sql.Register.
//
// Register makes a database driver available by the provided driver name.
// Driver's initialize handler and audit rules register by Register.
func Register(name string, h handler, rs []*model.Rule) {
	_, exist := drivers[name]
	if exist {
		panic("duplicated driver name")
	}

	driversMu.Lock()
	drivers[name] = h
	driversMu.Unlock()

	rulesMu.Lock()
	rules = append(rules, rs...)
	rulesMu.Unlock()
}

// NewDriver return a new instantiated Driver.
func NewDriver(log *logrus.Entry, inst *model.Instance, schema string) (Driver, error) {
	driversMu.RLock()
	defer driversMu.RUnlock()

	d, exist := drivers[inst.DbType]
	if !exist {
		return nil, fmt.Errorf("driver type %v is not supported", inst.DbType)
	}

	return d(log, inst, schema), nil
}

func AllRules() []*model.Rule {
	rulesMu.RLock()
	defer rulesMu.RUnlock()
	return rules
}

func AllDrivers() []string {
	rulesMu.RLock()
	defer rulesMu.RUnlock()

	driverNames := make([]string, 0, len(drivers))
	for n := range drivers {
		driverNames = append(driverNames, n)
	}
	return driverNames
}

var ErrNodesCountExceedOne = errors.New("after parse, nodes count exceed one")

// Driver is a interface that must be implemented by a database.
//
// Driver is responsible for two primary things:
// 1. privode handle to communicate with database
// 2. audit SQL with rules
type Driver interface {
	Close()
	Ping(ctx context.Context) error
	Exec(ctx context.Context, query string) (driver.Result, error)
	Tx(ctx context.Context, queries ...string) ([]driver.Result, error)
	Query(ctx context.Context, query string, args ...interface{}) ([]map[string]sql.NullString, error)

	// Schemas export all supported schemas.
	//
	// For example, performance_schema/performance_schema... which in MySQL is not allowed for auditing.
	Schemas(ctx context.Context) ([]string, error)

	// Parse parse sqlText to Node array.
	//
	// sqlText may be single SQL or batch SQLs.
	Parse(sqlText string) ([]Node, error)

	// Audit sql with rules. sql is single SQL text.
	//
	// Multi Audit call may be in one context.
	// For example:
	//		driver, _ := NewDriver(..., ..., ...)
	// 		driver.Audit(..., "CREATE TABLE t1(id int)")
	// 		driver.Audit(..., "SELECT * FROM t1 WHERE id = 1")
	//      ...
	// driver should keep SQL context during it's lifecycle.
	Audit(rules []*model.Rule, sql string) (*AuditResult, error)

	// GenRollbackSQL generate sql's rollback SQL.
	GenRollbackSQL(sql string) (string, string, error)
}

// Node is a interface which unify SQL ast tree. It produce by Driver.Parse.
type Node interface {
	// Text get the raw SQL text of Node.
	Text() string

	// Type return type of SQL, such as DML/DDL/DCL.
	Type() string

	// Fingerprint generate fingerprint of Node's raw SQL.
	//
	// For example:
	// 		driver, _ := NewDriver(..., ..., ...)
	// 		nodes, _ := driver.Parse("select * from t1 where id = 1")
	//		f, _ := nodes[0].Fingerprint() // f == SELECT * FROM `t1` WHERE id = ?
	Fingerprint() (string, error)
}

func Tx(d Driver, baseSQLs []*model.BaseSQL) error {
	var retErr error
	var results []driver.Result
	qs := make([]string, 0, len(baseSQLs))

	for _, baseSQL := range baseSQLs {
		qs = append(qs, baseSQL.Content)
	}

	// todo(@wy): missing binlog fields of BaseSQL
	defer func() {
		for idx, baseSQL := range baseSQLs {
			if retErr != nil {
				baseSQL.ExecStatus = model.SQLExecuteStatusFailed
				baseSQL.ExecResult = retErr.Error()
				continue
			}
			rowAffects, _ := results[idx].RowsAffected()
			baseSQL.RowAffects = rowAffects
			baseSQL.ExecStatus = model.SQLExecuteStatusSucceeded
			baseSQL.ExecResult = "ok"
		}
	}()

	results, err := d.Tx(context.TODO(), qs...)
	if err != nil {
		retErr = err
	} else if len(results) != len(qs) {
		retErr = fmt.Errorf("number of transaction result does not match number of SQLs")
	}
	return retErr
}

func Exec(d Driver, baseSQL *model.BaseSQL) error {
	_, err := d.Exec(context.TODO(), baseSQL.Content)
	if err != nil {
		baseSQL.ExecStatus = model.SQLExecuteStatusFailed
		baseSQL.ExecResult = err.Error()
	} else {
		baseSQL.ExecStatus = model.SQLExecuteStatusSucceeded
		baseSQL.ExecResult = "ok"
	}
	return err
}

type AuditResult struct {
	results []*auditResult
}

type auditResult struct {
	level   string
	message string
}

func NewInspectResults() *AuditResult {
	return &AuditResult{
		results: []*auditResult{},
	}
}

// Level find highest Level in result
func (rs *AuditResult) Level() string {
	level := model.RuleLevelNormal
	for _, result := range rs.results {
		if model.RuleLevelMap[level] < model.RuleLevelMap[result.level] {
			level = result.level
		}
	}
	return level
}

func (rs *AuditResult) Message() string {
	messages := make([]string, len(rs.results))
	for n, result := range rs.results {
		var message string
		match, _ := regexp.MatchString(fmt.Sprintf(`^\[%s|%s|%s|%s|%s\]`,
			model.RuleLevelError, model.RuleLevelWarn, model.RuleLevelNotice, model.RuleLevelNormal, "osc"),
			result.message)
		if match {
			message = result.message
		} else {
			message = fmt.Sprintf("[%s]%s", result.level, result.message)
		}
		messages[n] = message
	}
	return strings.Join(messages, "\n")
}

func (rs *AuditResult) Add(level, message string, args ...interface{}) {
	if level == "" || message == "" {
		return
	}

	rs.results = append(rs.results, &auditResult{
		level:   level,
		message: fmt.Sprintf(message, args...),
	})
}
