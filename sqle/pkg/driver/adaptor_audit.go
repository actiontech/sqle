package driver

import (
	"context"
	"database/sql"
	"os"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/pkg/params"

	"github.com/hashicorp/go-hclog"
	goPlugin "github.com/hashicorp/go-plugin"
	"github.com/pkg/errors"
)

// AuditAdaptor is a wrapper for the sqle driver layer. It
// provides a more simpler interface for the database plugin.
type AuditAdaptor struct {
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

// NewAdaptor create a database plugin AuditAdaptor with dialector.
// NewAdaptor is actually NewAuditAdaptor, but the method name cannot be changed for historical reasons
func NewAdaptor(dt Dialector) *AuditAdaptor {
	return &AuditAdaptor{
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

func (a *AuditAdaptor) AddRule(r *driver.Rule, h rawSQLRuleHandler) {
	a.rules = append(a.rules, r)
	a.ruleToRawHandler[r.Name] = h
}

func (a *AuditAdaptor) AddAdditionalParams(p *params.Param) {
	a.additionalParams = append(a.additionalParams, p)
}

func (a *AuditAdaptor) AddRuleWithSQLParser(r *driver.Rule, h astSQLRuleHandler) {
	a.rules = append(a.rules, r)
	a.ruleToASTHandler[r.Name] = h
}

func (a *AuditAdaptor) Serve(opts ...AdaptorOption) {
	plugin := a.GeneratePlugin(opts...)
	a.l.Info("start serve plugin", "name", a.dt)
	p := driver.NewPlugin()
	p.AddPlugin(driver.PluginNameAuditDriver, driver.DefaultPluginVersion, plugin)
	p.Serve()
}

func (a *AuditAdaptor) GeneratePlugin(opts ...AdaptorOption) goPlugin.Plugin {
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

	r := &auditRegistererImpl{
		dt:               a.dt,
		rules:            a.rules,
		additionalParams: a.additionalParams,
	}

	newDriver := func(cfg *driver.Config) driver.Driver {
		if p, exist := pluginImpls[driver.PluginNameAuditDriver]; exist {
			return p
		}

		a.cfg = cfg

		di := &pluginImpl{a: a}

		if cfg.DSN == nil {
			pluginImpls[driver.PluginNameAuditDriver] = di
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
		pluginImpls[driver.PluginNameAuditDriver] = di
		return di
	}

	return driver.NewAuditDriverPlugin(r, newDriver)
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

var _ driver.Driver = (*pluginImpl)(nil)
var _ driver.Registerer = (*auditRegistererImpl)(nil)

type auditRegistererImpl struct {
	dt               Dialector
	rules            []*driver.Rule
	additionalParams params.Params
}

func (r *auditRegistererImpl) Name() string {
	return r.dt.String()
}

func (r *auditRegistererImpl) Rules() []*driver.Rule {
	return r.rules
}

func (r *auditRegistererImpl) AdditionalParams() params.Params {
	return r.additionalParams
}
