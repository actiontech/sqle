//go:build enterprise
// +build enterprise

package server

import (
	"regexp"
	"strconv"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/sirupsen/logrus"
)

func CustomRuleAudit(l *logrus.Entry, DbType string, sqls []string, results []*driverV2.AuditResults) {
	st := model.GetStorage()
	customRules, exist, err := st.GetCustomRuleByDBTypeAndScriptType(DbType, "regular")
	if err != nil {
		l.Errorf("query custom rules failed:%v", err)
	}
	if exist {
		return
	}
	for i, sql := range sqls {
		for _, customRule := range customRules {
			ruleScript := customRule.RuleScript
			regex, err := regexp.Compile(ruleScript)
			if err != nil {
				// 一些特殊字符会报错，需要将他们转为原始字符。例如:[\u4e00-\u9fa5]
				ruleScript, err = strconv.Unquote(`"` + ruleScript + `"`)
				if err != nil {
					l.Errorf("str unquote failed:%v", err)
					continue
				}
				regex, err = regexp.Compile(ruleScript)
				if err != nil {
					l.Errorf("regexp compile failed:%v", err)
					continue
				}
			}
			if regex.MatchString(sql) {
				res := driverV2.AuditResult{
					Message:  customRule.Desc,
					RuleName: customRule.RuleId,
					Level:    driverV2.RuleLevel(customRule.Level),
				}
				result := results[i]  // TODO: 判断越界
				result.Results = append(result.Results, &res)
			}
		}
	}
}