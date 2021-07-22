package driver

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
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
	Audit(rules []*model.Rule, baseSQLs []*model.BaseSQL, isSkip func(node Node) bool) ([]*model.ExecuteSQL, []*model.RollbackSQL, error)
}

type Node interface {
	Text() string
	Type() string
	Fingerprint() (string, error)
}
