package ai

import (
	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
)

var sourceRuleHandlers []*rulepkg.SourceHandler

func init() {
	aiRuleHandlers := rulepkg.GenerateI18nRuleHandlers(plocale.Bundle, sourceRuleHandlers)
	rulepkg.AIRuleHandlers = append(rulepkg.AIRuleHandlers, aiRuleHandlers...)
	for _, handler := range aiRuleHandlers {
		rulepkg.AIRuleHandlerMap[handler.Rule.Name] = handler
	}
}
