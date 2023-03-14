package driver

import (
	"context"
	"os"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	hclog "github.com/hashicorp/go-hclog"
)

type RawSQLRuleHandler interface {
	GetRule() *driverV2.Rule
	BeforeAudit(interface{}) (interface{}, error)
	Audit(ctx context.Context, rawSQL string, nextSQL []string) (string, error)
}

type AstSQLRuleHandler interface {
	GetRule() *driverV2.Rule
	BeforeAudit(interface{}) (interface{}, error)
	Audit(ctx context.Context, astSQL interface{}, nextSQL []string) (string, error)
}

type AuditHandler struct {
	SqlParserFn      func(string) (interface{}, error)
	RuleToRawHandler map[string] /*rule name*/ RawSQLRuleHandler
	RuleToASTHandler map[string] /*rule name*/ AstSQLRuleHandler
}

type DriverBuilder struct {
	l hclog.Logger

	Dt Dialector
	ah *AuditHandler

	Meta *driverV2.DriverMetas
}

func NewDriverBuilder(dt Dialector) *DriverBuilder {
	b := &DriverBuilder{
		l: hclog.New(&hclog.LoggerOptions{
			JSONFormat: true,
			Output:     os.Stderr,
			Level:      hclog.Trace,
		}),
		Dt: dt,
		Meta: &driverV2.DriverMetas{
			PluginName:               dt.String(),
			DatabaseAdditionalParams: dt.DatabaseAdditionalParam(),
		},
		ah: &AuditHandler{
			RuleToRawHandler: make(map[string]RawSQLRuleHandler),
			RuleToASTHandler: make(map[string]AstSQLRuleHandler),
		},
	}
	return b
}

func (b *DriverBuilder) DefaultServe() {
	b.SetEnableOptionalModule(driverV2.OptionalModuleQuery)
	b.Serve(NewDriverImpl)
}

func (b *DriverBuilder) Serve(fn func(hclog.Logger, Dialector, *AuditHandler, *driverV2.Config) (driverV2.Driver, error)) {
	defer func() {
		if err := recover(); err != nil {
			b.l.Error("panic", "err", err)
		}
	}()

	if len(b.Meta.Rules) == 0 {
		b.l.Info("no rule in plugin adaptor", "name", b.Dt)
	}

	if len(b.ah.RuleToASTHandler) != 0 && b.ah.SqlParserFn == nil {
		panic("Add rule by AddRuleWithSQLParser(), but no SQL parser provided.")
	}

	newDriver := func(cfg *driverV2.Config) (driverV2.Driver, error) {
		return fn(b.l, b.Dt, b.ah, cfg)

	}
	driverV2.ServePlugin(*b.Meta, newDriver)
}

func (a *DriverBuilder) AddRule(h RawSQLRuleHandler) {
	a.Meta.Rules = append(a.Meta.Rules, h.GetRule())
	a.ah.RuleToRawHandler[h.GetRule().Name] = h
}

func (a *DriverBuilder) AddRuleWithSQLParser(h AstSQLRuleHandler) {
	a.Meta.Rules = append(a.Meta.Rules, h.GetRule())
	a.ah.RuleToASTHandler[h.GetRule().Name] = h
}

func (b *DriverBuilder) SetSQLParserFn(parser func(string) (interface{}, error)) {
	b.ah.SqlParserFn = parser
}

func (b *DriverBuilder) SetEnableOptionalModule(modules ...driverV2.OptionalModule) {
	b.Meta.EnabledOptionalModule = modules
}
