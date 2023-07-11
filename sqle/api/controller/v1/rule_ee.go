//go:build enterprise
// +build enterprise

package v1

import (
	"fmt"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
)

func convertCustomRuleToRes(rule *model.CustomRule) CustomRuleResV1 {
	ruleRes := CustomRuleResV1{
		RuleId:     rule.RuleId,
		DBType:     rule.DBType,
		RuleName:   rule.RuleName,
		Level:      rule.Level,
		Type:       rule.Typ,
		RuleScript: rule.RuleScript,
	}
	if rule.Params != nil && len(rule.Params) > 0 {
		paramsRes := make([]RuleParamResV1, 0, len(rule.Params))
		for _, p := range rule.Params {
			paramRes := RuleParamResV1{
				Key:   p.Key,
				Desc:  p.Desc,
				Type:  string(p.Type),
				Value: p.Value,
			}
			paramsRes = append(paramsRes, paramRes)
		}
		ruleRes.Params = paramsRes
	}
	return ruleRes
}

func convertCustomRulesToRes(rules []*model.CustomRule) []CustomRuleResV1 {
	rulesRes := make([]CustomRuleResV1, 0, len(rules))
	for _, rule := range rules {
		rulesRes = append(rulesRes, convertCustomRuleToRes(rule))
	}
	return rulesRes
}

func getCustomRules(c echo.Context) error {
	s := model.GetStorage()
	req := new(GetCustomRulesReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	queryFields := "rule_id, rule_name, db_type, `desc`, level, type, params"
	var err error
	var rules []*model.CustomRule
	if req.FilterDBType == "" && req.FilterRuleName == "" {
		rules, err = s.GetCustomRules(queryFields)
	} else {
		rules, err = s.GetCustomRulesByRuleNameAndDBType(queryFields, req.FilterDBType, req.FilterRuleName)
	}

	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetCustomRulesResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertCustomRulesToRes(rules),
	})
}

func deleteCustomRule(c echo.Context) error {
	ruleId := c.Param("rule_id")

	s := model.GetStorage()
	rule, exist, err := s.GetCustomRuleByRuleId(ruleId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("rule is not exist")))
	}

	err = s.Delete(rule)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}
