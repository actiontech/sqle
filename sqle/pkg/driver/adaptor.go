package driver

import (
	"context"
	"database/sql"
	_driver "database/sql/driver"
	"os"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/hashicorp/go-hclog"
	goPlugin "github.com/hashicorp/go-plugin"
	"github.com/percona/go-mysql/query"
	"github.com/pkg/errors"
	"vitess.io/vitess/go/vt/sqlparser"
)

// Adaptor is a wrapper for the sqle driver layer. It
// provides a more simpler interface for the database plugin.
type Adaptor struct {
	l hclog.Logger

	cfg *driver.Config

	dt Dialector

	rules            []*driver.Rule
	ruleToRawHandler map[string] /*rule name*/ rawSQLRuleHandler
	ruleToASTHandler map[string] /*rule name*/ astSQLRuleHandler

	additionalParams params.Params

	ao *adaptorOptions
}

type adaptorOptions struct {
	sqlParser func(string) (interface{}, error)
}

type rawSQLRuleHandler func(ctx context.Context, rule *driver.Rule, rawSQL string) (string, error)
type astSQLRuleHandler func(ctx context.Context, rule *driver.Rule, astSQL interface{}) (string, error)

// NewAdaptor create a database plugin Adaptor with dialector.
func NewAdaptor(dt Dialector) *Adaptor {
	return &Adaptor{
		ao: &adaptorOptions{},

		dt: dt,
		l: hclog.New(&hclog.LoggerOptions{
			JSONFormat: true,
			Output:     os.Stderr,
			Level:      hclog.Trace,
		}),
		ruleToRawHandler: make(map[string]rawSQLRuleHandler),
		ruleToASTHandler: make(map[string]astSQLRuleHandler),
		additionalParams: params.Params{},
	}
}

func (a *Adaptor) AddRule(r *driver.Rule, h rawSQLRuleHandler) {
	a.rules = append(a.rules, r)
	a.ruleToRawHandler[r.Name] = h
}

func (a *Adaptor) AddAdditionalParams(p *params.Param) {
	a.additionalParams = append(a.additionalParams, p)
}

func (a *Adaptor) AddRuleWithSQLParser(r *driver.Rule, h astSQLRuleHandler) {
	a.rules = append(a.rules, r)
	a.ruleToASTHandler[r.Name] = h
}

func (a *Adaptor) Serve(opts ...AdaptorOption) {
	plugin := a.GeneratePlugin(opts...)
	a.l.Info("start serve plugin", "name", a.dt)
	p := driver.NewPlugin()
	p.AddPlugin(driver.PluginNameDriver, driver.DefaultPluginVersion, plugin)
	p.Serve()
}

func (a *Adaptor) GeneratePlugin(opts ...AdaptorOption) goPlugin.Plugin {
	defer func() {
		if err := recover(); err != nil {
			a.l.Error("panic", "err", err)
		}
	}()

	for _, opt := range opts {
		opt.apply(a.ao)
	}

	if len(a.rules) == 0 {
		a.l.Info("no rule in plugin adaptor", "name", a.dt)
	}

	if len(a.ruleToASTHandler) != 0 && a.ao.sqlParser == nil {
		panic("Add rule by AddRuleWithSQLParser(), but no SQL parser provided.")
	}

	r := &registererImpl{
		dt:               a.dt,
		rules:            a.rules,
		additionalParams: a.additionalParams,
	}

	newDriver := func(cfg *driver.Config) driver.Driver {
		a.cfg = cfg

		di := &driverImpl{a: a}

		if cfg.DSN == nil {
			return di
		}

		driverName, dsnDetail := a.dt.Dialect(cfg.DSN)
		db, err := sql.Open(driverName, dsnDetail)
		if err != nil {
			panic(errors.Wrap(err, "open database failed when new driver"))
		}
		conn, err := db.Conn(context.TODO())
		if err != nil {
			panic(errors.Wrap(err, "get database connection failed when new driver"))
		}
		if err := conn.PingContext(context.TODO()); err != nil {
			panic(errors.Wrap(err, "ping database connection failed when new driver"))
		}

		di.db = db
		di.conn = conn
		return di
	}

	return driver.NewDriverPlugin(r, newDriver)
}

// AdaptorOption store some custom options for the driver adaptor.
type AdaptorOption interface {
	apply(*adaptorOptions)
}

type optionFunc struct {
	f func(*adaptorOptions)
}

func newOptionFunc(f func(*adaptorOptions)) *optionFunc {
	return &optionFunc{
		f: f,
	}
}

func (this *optionFunc) apply(a *adaptorOptions) {
	this.f(a)
}

// WithSQLParser define custom SQL parser. If set, the adaptor
// will use it to parse the SQL. User can assert the SQL to correspond
// ast structure in ruleHandler.
func WithSQLParser(parser func(sql string) (ast interface{}, err error)) AdaptorOption {
	return newOptionFunc(func(a *adaptorOptions) {
		a.sqlParser = parser
	})
}

var _ driver.Driver = (*driverImpl)(nil)
var _ driver.Registerer = (*registererImpl)(nil)

type registererImpl struct {
	dt               Dialector
	rules            []*driver.Rule
	additionalParams params.Params
}

func (r *registererImpl) Name() string {
	return r.dt.String()
}

func (r *registererImpl) Rules() []*driver.Rule {
	return r.rules
}

func (r *registererImpl) AdditionalParams() params.Params {
	return r.additionalParams
}

type driverImpl struct {
	a    *Adaptor
	db   *sql.DB
	conn *sql.Conn
}

func (d *driverImpl) Close(ctx context.Context) {
	if err := d.conn.Close(); err != nil {
		d.a.l.Error("failed to close connection in driver adaptor", "err", err)
	}
	if err := d.db.Close(); err != nil {
		d.a.l.Error("failed to close database in driver adaptor", "err", err)
	}
}

func (d *driverImpl) Ping(ctx context.Context) error {
	if err := d.conn.PingContext(ctx); err != nil {
		return errors.Wrap(err, "ping in driver adaptor")
	}
	return nil
}

func (d *driverImpl) Exec(ctx context.Context, sql string) (_driver.Result, error) {
	res, err := d.conn.ExecContext(ctx, sql)
	if err != nil {
		return nil, errors.Wrap(err, "exec sql in driver adaptor")
	}
	return res, nil
}

func (d *driverImpl) Tx(ctx context.Context, sqls ...string) ([]_driver.Result, error) {
	var (
		err error
		tx  *sql.Tx
	)

	tx, err = d.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "begin tx in driver adaptor")
	}

	defer func() {
		if err != nil {
			if err = tx.Rollback(); err != nil {
				err = errors.Wrap(err, "rollback tx in driver adaptor")
				return
			}
		} else {
			if err = tx.Commit(); err != nil {
				err = errors.Wrap(err, "commit tx in driver adaptor")
				return
			}
		}
	}()

	results := make([]_driver.Result, 0, len(sqls))
	for _, sql := range sqls {
		result, e := tx.ExecContext(ctx, sql)
		if e != nil {
			err = errors.Wrap(e, "exec sql in driver adaptor")
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func (d *driverImpl) Schemas(ctx context.Context) ([]string, error) {
	rows, err := d.conn.QueryContext(ctx, d.a.dt.ShowDatabaseSQL())
	if err != nil {
		return nil, errors.Wrap(err, "query database in driver adaptor")
	}
	defer rows.Close()

	var schemas []string
	for rows.Next() {
		var schema string
		if err := rows.Scan(&schema); err != nil {
			return nil, errors.Wrap(err, "scan database in driver adaptor")
		}
		schemas = append(schemas, schema)
	}

	if rows.Err() != nil {
		return nil, errors.Wrap(rows.Err(), "scan database in driver adaptor")
	}

	return schemas, nil
}

func (d *driverImpl) Parse(ctx context.Context, sql string) ([]driver.Node, error) {
	sqls, err := sqlparser.SplitStatementToPieces(sql)
	if err != nil {
		return nil, errors.Wrap(err, "split sql")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "split sql %s error", sql)
	}

	nodes := make([]driver.Node, 0, len(sqls))
	for _, sql := range sqls {
		n := driver.Node{
			Text:        sql,
			Type:        classifySQL(sql),
			Fingerprint: query.Fingerprint(sql),
		}
		nodes = append(nodes, n)
	}
	return nodes, nil
}

func classifySQL(sql string) (sqlType string) {
	if utils.HasPrefix(sql, "update", false) ||
		utils.HasPrefix(sql, "insert", false) ||
		utils.HasPrefix(sql, "delete", false) {
		return driver.SQLTypeDML
	}

	return driver.SQLTypeDDL
}

func (d *driverImpl) Audit(ctx context.Context, sql string) (*driver.AuditResult, error) {
	var err error
	var ast interface{}
	if d.a.ao.sqlParser != nil {
		ast, err = d.a.ao.sqlParser(sql)
		if err != nil {
			return nil, errors.Wrap(err, "parse sql")
		}
	}

	result := driver.NewInspectResults()
	for _, rule := range d.a.cfg.Rules {
		handler, ok := d.a.ruleToRawHandler[rule.Name]
		if ok {
			msg, err := handler(ctx, rule, sql)
			if err != nil {
				return nil, errors.Wrapf(err, "audit SQL %s in driver adaptor", sql)
			}
			result.Add(rule.Level, msg)
		} else {
			handler, ok := d.a.ruleToASTHandler[rule.Name]
			if ok {
				msg, err := handler(ctx, rule, ast)
				if err != nil {
					return nil, errors.Wrapf(err, "audit SQL %s in driver adaptor", sql)
				}
				result.Add(rule.Level, msg)
			}
		}
	}

	return result, nil
}

func (d *driverImpl) GenRollbackSQL(ctx context.Context, sql string) (string, string, error) {
	return "", "", nil
}
