//go:build !enterprise
// +build !enterprise

package rule

import driverV2 "github.com/actiontech/sqle/sqle/driver/v2"

var defaultRuleKnowledgeMap = map[string]driverV2.RuleKnowledge{}
