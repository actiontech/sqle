package driver

import (
	"context"
	"os"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
)

type RawSQLRuleHandler func(ctx context.Context, rule *driverV2.Rule, rawSQL string, nextSQL []string) (string, error)

type AstSQLRuleHandler func(ctx context.Context, rule *driverV2.Rule, astSQL interface{}, nextSQL []string) (string, error)

type AuditHandler struct {
	SqlParserFn      func(string) (interface{}, error)
	RuleToRawHandler map[string] /*rule name*/ RawSQLRuleHandler
	RuleToASTHandler map[string] /*rule name*/ AstSQLRuleHandler
}

func (a *AuditHandler) Audit(ctx context.Context, rules []*driverV2.Rule, sql string, nextSQL []string) (*driverV2.AuditResults, error) {
	result := driverV2.NewAuditResults()
	for _, rule := range rules {
		handler, ok := a.RuleToRawHandler[rule.Name]
		if ok {
			msg, err := handler(ctx, rule, sql, nextSQL)
			if err != nil {
				return nil, errors.Wrapf(err, "audit SQL %s in driver adaptor", sql)
			}
			result.Add(rule.Level, msg)
		} else {
			handler, ok := a.RuleToASTHandler[rule.Name]
			if ok {
				var err error
				var ast interface{}
				if a.SqlParserFn != nil {
					ast, err = a.SqlParserFn(sql)
					if err != nil {
						return nil, errors.Wrap(err, "parse sql")
					}
				}
				msg, err := handler(ctx, rule, ast, nextSQL)
				if err != nil {
					return nil, errors.Wrapf(err, "audit SQL %s in driver adaptor", sql)
				}
				result.Add(rule.Level, msg)
			}
		}
	}
	return result, nil
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

func (a *DriverBuilder) AddRule(r *driverV2.Rule, h RawSQLRuleHandler) {
	a.Meta.Rules = append(a.Meta.Rules, r)
	a.ah.RuleToRawHandler[r.Name] = h
}

func (a *DriverBuilder) AddRuleWithSQLParser(r *driverV2.Rule, h AstSQLRuleHandler) {
	a.Meta.Rules = append(a.Meta.Rules, r)
	a.ah.RuleToASTHandler[r.Name] = h
}

func (b *DriverBuilder) SetSQLParserFn(parser func(string) (interface{}, error)) {
	b.ah.SqlParserFn = parser
}

func (b *DriverBuilder) SetEnableOptionalModule(modules ...driverV2.OptionalModule) {
	b.Meta.EnabledOptionalModule = modules
}
