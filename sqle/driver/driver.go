package driver

import (
	"context"
	"database/sql"
	"database/sql/driver"
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

type handler func(log *logrus.Entry, inst *model.Instance, schema string) Driver

func Register(name string, h handler, rs []*model.Rule) {
	_, exist := drivers[name]
	if exist {
		panic("duplicated driver name")
	}

	driversMu.Lock()
	drivers[name] = h
	driversMu.Unlock()

	rulesMu.Lock()
	for _, r := range rs {
		rules = append(rules, r)
	}
	rulesMu.Unlock()
}

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

type Driver interface {
	Close()
	Ping(ctx context.Context) error
	Exec(ctx context.Context, query string) (driver.Result, error)
	Tx(ctx context.Context, queries ...string) ([]driver.Result, error)
	Query(ctx context.Context, query string, args ...interface{}) ([]map[string]sql.NullString, error)

	Parse(sqlText string) ([]Node, error)
	Audit(rules []*model.Rule, sql string) (*AuditResult, error)
}

type Node interface {
	Text() string
	Type() string
	Fingerprint() (string, error)
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
	rs.results = append(rs.results, &auditResult{
		level:   level,
		message: fmt.Sprintf(message, args...),
	})
}
