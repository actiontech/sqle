package driver

import (
	"context"
	"os"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	hclog "github.com/hashicorp/go-hclog"
)

type rawSQLRuleHandler func(ctx context.Context, rule *driverV2.Rule, rawSQL string, nextSQL []string) (string, error)

type astSQLRuleHandler func(ctx context.Context, rule *driverV2.Rule, astSQL interface{}, nextSQL []string) (string, error)

type DriverBuilder struct {
	l    hclog.Logger
	Dt   Dialector
	Opts *adaptorOptions

	Meta             *driverV2.DriverMetas
	RuleToRawHandler map[string] /*rule name*/ rawSQLRuleHandler
	RuleToASTHandler map[string] /*rule name*/ astSQLRuleHandler
}

func NewDriverBuilder(dt Dialector, opts ...AdaptorOption) *DriverBuilder {
	b := &DriverBuilder{
		l: hclog.New(&hclog.LoggerOptions{
			JSONFormat: true,
			Output:     os.Stderr,
			Level:      hclog.Trace,
		}),
		Dt: dt,
		Meta: &driverV2.DriverMetas{
			PluginName: dt.String(),
		},
		RuleToRawHandler: make(map[string]rawSQLRuleHandler),
		RuleToASTHandler: make(map[string]astSQLRuleHandler),
	}
	for _, opt := range opts {
		opt.apply(b.Opts)
	}
	return b
}

func (b *DriverBuilder) DefaultServe() {
	b.SetEnableOptionalModule(driverV2.OptionalModuleQuery)
	b.Serve(NewDriverImpl)
}

func (b *DriverBuilder) Serve(fn func(*DriverBuilder, *driverV2.Config) (driverV2.Driver, error)) {
	defer func() {
		if err := recover(); err != nil {
			b.l.Error("panic", "err", err)
		}
	}()

	if len(b.Meta.Rules) == 0 {
		b.l.Info("no rule in plugin adaptor", "name", b.Dt)
	}

	if len(b.RuleToASTHandler) != 0 && b.Opts.sqlParser == nil {
		panic("Add rule by AddRuleWithSQLParser(), but no SQL parser provided.")
	}

	newDriver := func(cfg *driverV2.Config) (driverV2.Driver, error) {
		return fn(b, cfg)

	}
	driverV2.ServePlugin(*b.Meta, newDriver)
}

func (a *DriverBuilder) AddRule(r *driverV2.Rule, h rawSQLRuleHandler) {
	a.Meta.Rules = append(a.Meta.Rules, r)
	a.RuleToRawHandler[r.Name] = h
}

func (a *DriverBuilder) AddRuleWithSQLParser(r *driverV2.Rule, h astSQLRuleHandler) {
	a.Meta.Rules = append(a.Meta.Rules, r)
	a.RuleToASTHandler[r.Name] = h
}

func (a *DriverBuilder) AddDatabaseAdditionalParam(p *params.Param) {
	a.Meta.DatabaseAdditionalParams = append(a.Meta.DatabaseAdditionalParams, p)
}

func (b *DriverBuilder) SetEnableOptionalModule(modules ...driverV2.OptionalModule) {
	b.Meta.EnabledOptionalModule = modules
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

type adaptorOptions struct {
	sqlParser func(string) (interface{}, error)
}

// WithSQLParser define custom SQL parser. If set, the adaptor
// will use it to parse the SQL. User can assert the SQL to correspond
// ast structure in ruleHandler.
func WithSQLParser(parser func(sql string) (ast interface{}, err error)) AdaptorOption {
	return newOptionFunc(func(a *adaptorOptions) {
		a.sqlParser = parser
	})
}
