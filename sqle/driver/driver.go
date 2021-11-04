package driver

import (
	"context"
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/actiontech/sqle/sqle/model"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	// drivers store instantiate handlers for MySQL or gRPC plugin.
	drivers   = make(map[string]handler)
	driversMu sync.RWMutex

	rules   []*model.Rule
	rulesMu sync.RWMutex
)

type Config struct {
	IsOfflineAudit bool
	Schema         string
	Inst           *model.Instance
	Rules          []*model.Rule
}

func NewConfig(inst *model.Instance, rules []*model.Rule, schema string, isOfflineAudit bool) *Config {
	return &Config{
		Inst:           inst,
		Rules:          rules,
		Schema:         schema,
		IsOfflineAudit: isOfflineAudit,
	}
}

// handler is a template which Driver plugin should provide such function signature.
type handler func(log *logrus.Entry, config *Config) (Driver, error)

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

type ErrDriverNotSupported struct {
	DriverTyp string
}

func (e *ErrDriverNotSupported) Error() string {
	return fmt.Sprintf("driver type %v is not supported", e.DriverTyp)
}

// NewDriver return a new instantiated Driver.
func NewDriver(log *logrus.Entry, inst *model.Instance, isOfflineAudit bool, dbType, schema string) (Driver, error) {
	driversMu.RLock()
	defer driversMu.RUnlock()

	d, exist := drivers[dbType]
	if !exist {
		return nil, fmt.Errorf("driver type %v is not supported", inst.DbType)
	}

	st := model.GetStorage()

	var err error
	var rules []*model.Rule

	if isOfflineAudit {
		// use default_{db_type}'s rules if audit is offline
		// refer: model.utils.CreateDefaultTemplate
		// TODO: add function to generate default rule template name
		templateName := fmt.Sprintf("default_%v", dbType)
		rules, err = st.GetRulesFromRuleTemplateByName(templateName)
	} else {
		rules, err = st.GetRulesByInstanceId(fmt.Sprintf("%v", inst.ID))
	}

	if err != nil {
		return nil, errors.Wrap(err, "get rules for audit")
	}

	return d(log, NewConfig(inst, rules, schema, isOfflineAudit))
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
// It's implementation maybe on the same process or over gRPC(by go-plugin).
//
// Driver is responsible for two primary things:
// 1. privode handle to communicate with database
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
	Rules() []*model.Rule
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
