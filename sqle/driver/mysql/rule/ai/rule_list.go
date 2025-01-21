package ai

import (
	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
)

var sourceRuleHandlers []*rulepkg.SourceHandler

func init() {
	aiRuleHandlers := rulepkg.GenerateI18nRuleHandlers(plocale.Bundle, sourceRuleHandlers)
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, aiRuleHandlers...)
}
