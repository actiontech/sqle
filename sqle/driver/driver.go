package driver

import (
	"context"
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	// drivers store instantiate handlers for MySQL or gRPC plugin.
	drivers   = make(map[string]handler)
	driversMu sync.RWMutex

	// rules store audit rules for each driver.
	rules   map[string][]*Rule
	rulesMu sync.RWMutex

	// additionalParams store driver additional params
	additionalParams   map[string]params.Params
	additionalParamsMu sync.RWMutex
)

const (
	SQLTypeDML = "dml"
	SQLTypeDDL = "ddl"
)

const (
	DriverTypeMySQL      = "mysql"
	DriverTypePostgreSQL = "PostgreSQL"
)

// DSN provide necessary information to connect to database.
type DSN struct {
	Host             string
	Port             string
	User             string
	Password         string
	AdditionalParams params.Params

	// DatabaseName is the default database to connect.
	DatabaseName string
}

type RuleLevel string

const (
	RuleLevelNormal RuleLevel = "normal"
	RuleLevelNotice RuleLevel = "notice"
	RuleLevelWarn   RuleLevel = "warn"
	RuleLevelError  RuleLevel = "error"
)

var ruleLevelMap = map[RuleLevel]int{
	RuleLevelNormal: 0,
	RuleLevelNotice: 1,
	RuleLevelWarn:   2,
	RuleLevelError:  3,
}

func (r RuleLevel) LessOrEqual(l RuleLevel) bool {
	return ruleLevelMap[r] <= ruleLevelMap[l]
}

func (r RuleLevel) More(l RuleLevel) bool {
	return ruleLevelMap[r] > ruleLevelMap[l]
}

type Rule struct {
	Name string
	Desc string

	// Category is the category of the rule. Such as "Naming Conventions"...
	// Rules will be displayed on the SQLE rule list page by category.
	Category string
	Level    RuleLevel
	Params   params.Params
}

//func (r *Rule) GetValueInt(defaultRule *Rule) int64 {
//	value := r.getValue(DefaultSingleParamKeyName, defaultRule)
//	i, err := strconv.ParseInt(value, 10, 64)
//	if err != nil {
//		return 0
//	}
//	return i
//}
//
//func (r *Rule) GetSingleValue() string {
//	value, _ := r.Params.GetParamValue(DefaultSingleParamKeyName)
//	return value
//}
//
//func (r *Rule) GetSingleValueInt() int {
//	value := r.GetSingleValue()
//	i, err := strconv.Atoi(value)
//	if err != nil {
//		return 0
//	}
//	return i
//}

//func (r *Rule) getValue(key string, defaultRule *Rule) string {
//	var value string
//	var exist bool
//	value, exist = r.Params.GetParamValue(key)
//	if !exist {
//		value, _ = defaultRule.Params.GetParamValue(key)
//	}
//	return value
//}

// Config define the configuration for driver.
type Config struct {
	DSN   *DSN
	Rules []*Rule
}

// NewConfig return a config for driver.
//
// 1. dsn is nil, rules is not nil. Use drive to do Offline Audit.
// 2. dsn is not nil, rule is nil. Use drive to communicate with database only.
// 3. dsn is not nil, rule is not nil. Most common usecase.
func NewConfig(dsn *DSN, rules []*Rule) (*Config, error) {
	if dsn == nil && rules == nil {
		fmt.Println("dsn is nil, and rules is nil, nothing can be done by driver")
	}

	return &Config{
		DSN:   dsn,
		Rules: rules,
	}, nil
}

// handler is a template which Driver plugin should provide such function signature.
type handler func(log *logrus.Entry, c *Config) (Driver, error)

// Register like sql.Register.
//
// Register makes a database driver available by the provided driver name.
// Driver's initialize handler and audit rules register by Register.
func Register(name string, h handler, rs []*Rule, ap params.Params) {
	_, exist := drivers[name]
	if exist {
		panic("duplicated driver name")
	}

	driversMu.Lock()
	drivers[name] = h
	driversMu.Unlock()

	rulesMu.Lock()
	if rules == nil {
		rules = make(map[string][]*Rule)
	}
	rules[name] = rs
	rulesMu.Unlock()

	additionalParamsMu.Lock()
	if additionalParams == nil {
		additionalParams = make(map[string]params.Params)
	}
	additionalParams[name] = ap
	additionalParamsMu.Unlock()
}

type DriverNotSupportedError struct {
	DriverTyp string
}

func (e *DriverNotSupportedError) Error() string {
	return fmt.Sprintf("driver type %v is not supported", e.DriverTyp)
}

// NewDriver return a new instantiated Driver.
func NewDriver(log *logrus.Entry, dbType string, cfg *Config) (Driver, error) {
	driversMu.RLock()
	defer driversMu.RUnlock()

	d, exist := drivers[dbType]
	if !exist {
		return nil, fmt.Errorf("driver type %v is not supported", dbType)
	}

	return d(log, cfg)
}

func AllRules() map[string][]*Rule {
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

func AllAdditionalParams() map[string] /*driver name*/ params.Params {
	additionalParamsMu.RLock()
	defer additionalParamsMu.RUnlock()

	newParams := map[string]params.Params{}
	for k, v := range additionalParams {
		newParams[k] = v.Copy()
	}
	return newParams
}

var ErrNodesCountExceedOne = errors.New("after parse, nodes count exceed one")

// Driver is a interface that must be implemented by a database.
//
// It's implementation maybe on the same process or over gRPC(by go-plugin).
//
// Driver is responsible for two primary things:
// 1. provides handle to communicate with database
// 2. audit SQL with rules
type Driver interface {
	Close(ctx context.Context)
	Ping(ctx context.Context) error
	Exec(ctx context.Context, query string) (driver.Result, error)
	Tx(ctx context.Context, queries ...string) ([]driver.Result, error)

	// Schemas export all supported schemas.
	//
	// For example, performance_schema/performance_schema... which in MySQL is not allowed for auditing.
	Schemas(ctx context.Context) ([]string, error)

	// Parse parse sqlText to Node array.
	//
	// sqlText may be single SQL or batch SQLs.
	Parse(ctx context.Context, sqlText string) ([]Node, error)

	// Audit sql with rules. sql is single SQL text.
	//
	// Multi Audit call may be in one context.
	// For example:
	//		driver, _ := NewDriver(..., ..., ...)
	// 		driver.Audit(..., "CREATE TABLE t1(id int)")
	// 		driver.Audit(..., "SELECT * FROM t1 WHERE id = 1")
	//      ...
	// driver should keep SQL context during it's lifecycle.
	Audit(ctx context.Context, sql string) (*AuditResult, error)

	// GenRollbackSQL generate sql's rollback SQL.
	GenRollbackSQL(ctx context.Context, sql string) (string, string, error)
}

// Registerer is the interface that all SQLe plugins must support.
type Registerer interface {
	// Name returns plugin name.
	Name() string

	// Rules returns all rules that plugin supported.
	Rules() []*Rule

	// AdditionalParams returns all additional params that plugin supported.
	AdditionalParams() params.Params
}

// Node is a interface which unify SQL ast tree. It produce by Driver.Parse.
type Node struct {
	// Text is the raw SQL text of Node.
	Text string

	// Type is type of SQL, such as DML/DDL/DCL.
	Type string

	// Fingerprint is fingerprint of Node's raw SQL.
	Fingerprint string
}

// // DSN like https://github.com/go-sql-driver/mysql/blob/master/dsn.go. type Config struct
// type DSN struct {
// 	Type string

// 	Host   string
// 	Port   string
// 	User   string
// 	Pass   string
// 	DBName string
// }

type AuditResult struct {
	results []*auditResult
}

type auditResult struct {
	level   RuleLevel
	message string
}

func NewInspectResults() *AuditResult {
	return &AuditResult{
		results: []*auditResult{},
	}
}

// Level find highest Level in result
func (rs *AuditResult) Level() RuleLevel {
	level := RuleLevelNormal
	for _, curr := range rs.results {
		if ruleLevelMap[curr.level] > ruleLevelMap[level] {
			level = curr.level
		}
	}
	return level
}

func (rs *AuditResult) Message() string {
	messages := make([]string, len(rs.results))
	for n, result := range rs.results {
		var message string
		match, _ := regexp.MatchString(fmt.Sprintf(`^\[%s|%s|%s|%s|%s\]`,
			RuleLevelError, RuleLevelWarn, RuleLevelNotice, RuleLevelNormal, "osc"),
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

func (rs *AuditResult) Add(level RuleLevel, message string, args ...interface{}) {
	if level == "" || message == "" {
		return
	}

	rs.results = append(rs.results, &auditResult{
		level:   level,
		message: fmt.Sprintf(message, args...),
	})
}

func (rs *AuditResult) HasResult() bool {
	return len(rs.results) != 0
}
