package driver

import (
	"context"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
)

type AstSQLRuleHandlerImpl struct {
	Rule    *driverV2.Rule
	Message string
}

func (r *AstSQLRuleHandlerImpl) GetRule() *driverV2.Rule {
	return r.Rule
}

func (r *AstSQLRuleHandlerImpl) BeforeAudit(interface{}) (interface{}, error) {
	return nil, nil
}

func (r *AstSQLRuleHandlerImpl) Audit(ctx context.Context, astSQL interface{}, nextSQL []string) (string, error) {
	return "", nil
}

type RawSQLRuleHandlerImpl struct {
	Rule    *driverV2.Rule
	Message string
}

func (r *RawSQLRuleHandlerImpl) GetRule() *driverV2.Rule {
	return r.Rule
}

func (r *RawSQLRuleHandlerImpl) BeforeAudit(interface{}) (interface{}, error) {
	return nil, nil
}

func (r *RawSQLRuleHandlerImpl) Audit(ctx context.Context, rawSQL string, nextSQL []string) (string, error) {
	return "", nil
}
