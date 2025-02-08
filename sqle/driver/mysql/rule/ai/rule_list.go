package ai

import (
	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
)

var sourceRuleHandlers []*rulepkg.SourceHandler

func init() {
	aiRuleHandlers := rulepkg.GenerateI18nRuleHandlers(plocale.Bundle, sourceRuleHandlers)
	rulepkg.AIRuleHandlers = append(rulepkg.AIRuleHandlers, aiRuleHandlers...)
	for k := range aiRuleHandlers {
		rulepkg.AIRuleHandlerMap[aiRuleHandlers[k].Rule.Name] = aiRuleHandlers[k]
		rulepkg.AllRules = append(rulepkg.AllRules, &aiRuleHandlers[k].Rule)
	}
}
