//go:build enterprise
// +build enterprise

package server

import (
	"fmt"
	"regexp"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/model"
	xerrors "github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

func CustomRuleAudit(l *logrus.Entry, task *model.Task, sqls []string, results []*driverV2.AuditResults, projectId *uint, ruleTemplateName string) {
	st := model.GetStorage()

	var customRules []*model.CustomRule
	var err error
	if ruleTemplateName != "" {
		if projectId == nil {
			err = xerrors.New("project id is needed when rule template name is given")
		} else {
			customRules, err = st.GetCustomRulesFromRuleTemplateByName([]uint{*projectId, model.ProjectIdForGlobalRuleTemplate}, ruleTemplateName)
		}
	} else {
		if task.Instance != nil {
			customRules, err = st.GetCustomRulesByInstanceId(fmt.Sprintf("%v", task.Instance.ID))
		} else {
			templateName := st.GetDefaultRuleTemplateName(task.DBType)
			// 默认规则模板从全局模板里拿
			customRules, err = st.GetCustomRulesFromRuleTemplateByName([]uint{model.ProjectIdForGlobalRuleTemplate}, templateName)
		}
	}
	if err != nil {
		l.Errorf("get rules error: %v", err)
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
				result := results[i] // TODO: 判断越界
				result.Results = append(result.Results, &res)
			}
		}
	}
}
