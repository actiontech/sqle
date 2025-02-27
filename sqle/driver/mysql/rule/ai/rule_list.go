package ai

import (
	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
)

var sourceRuleHandlers []*rulepkg.SourceHandler

func init() {
	aiRuleHandlers := rulepkg.GenerateI18nRuleHandlers(plocale.Bundle, sourceRuleHandlers, driverV2.DriverTypeMySQL)
	rulepkg.AIRuleHandlers = append(rulepkg.AIRuleHandlers, aiRuleHandlers...)
	for k := range aiRuleHandlers {
		rulepkg.AIRuleHandlerMap[aiRuleHandlers[k].Rule.Name] = aiRuleHandlers[k]
		rulepkg.AllRules = append(rulepkg.AllRules, &aiRuleHandlers[k].Rule)
	}
}
