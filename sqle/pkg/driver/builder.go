package driver

import (
	"context"
	"os"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/i18nPkg"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
)

type RawSQLRuleHandler func(ctx context.Context, rule *driverV2.Rule, rawSQL string, nextSQL []string) (i18nPkg.I18nStr, error)

type AstSQLRuleHandler func(ctx context.Context, rule *driverV2.Rule, astSQL interface{}, nextSQL []string) (i18nPkg.I18nStr, error)

type AuditHandler struct {
	SqlParserFn      func(string) (interface{}, error)
	RuleToRawHandler map[string] /*rule name*/ RawSQLRuleHandler
	RuleToASTHandler map[string] /*rule name*/ AstSQLRuleHandler
}

func (a *AuditHandler) Audit(ctx context.Context, rule *driverV2.Rule, sql string, nextSQL []string) (*driverV2.AuditResult, error) {
	result := &driverV2.AuditResult{}
	message := i18nPkg.I18nStr{}
	var err error

	handler, ok := a.RuleToRawHandler[rule.Name]
	if ok {
		message, err = handler(ctx, rule, sql, nextSQL)
		if err != nil {
			return nil, errors.Wrapf(err, "audit SQL %s in driver adaptor", sql)
		}
	} else {
		handler, ok := a.RuleToASTHandler[rule.Name]
		if ok {
			var ast interface{}
			if a.SqlParserFn != nil {
				ast, err = a.SqlParserFn(sql)
				if err != nil {
					return nil, errors.Wrap(err, "parse sql")
				}
			}
			message, err = handler(ctx, rule, ast, nextSQL)
			if err != nil {
				return nil, errors.Wrapf(err, "audit SQL %s in driver adaptor", sql)
			}
		}
	}
	if len(message) != 0 {
		result.Level = rule.Level
		result.RuleName = rule.Name
	}
	for langTag, langMsg := range message {
		result.I18nAuditResultInfo[langTag] = driverV2.AuditResultInfo{
			Message: langMsg,
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
