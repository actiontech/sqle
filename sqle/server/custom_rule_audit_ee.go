//go:build enterprise
// +build enterprise

package server

import (
	"regexp"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/sirupsen/logrus"
)

func CustomRuleAudit(l *logrus.Entry, task *model.Task, sqls []string, results []*driverV2.AuditResults, customRules []*model.CustomRule) {
	if len(customRules) == 0 {
		return
	}
	
	if len(results) != len(sqls) {
		l.Errorf("audit results [%d] does not match the number of SQL [%d]", len(results), len(sqls))
		return
	}

	for i, sql := range sqls {
		for _, customRule := range customRules {
			ruleScript := customRule.RuleScript
			regex, err := regexp.Compile(ruleScript)
			if err != nil {
				l.Errorf("regexp compile failed:%v", err)
				continue
			}
			if regex.MatchString(sql) {
				res := driverV2.AuditResult{
					Message:  customRule.Desc,
					RuleName: customRule.RuleId,
					Level:    driverV2.RuleLevel(customRule.Level),
				}
				result := results[i]
				result.Results = append(result.Results, &res)
			}
		}
	}
}
