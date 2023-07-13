//go:build enterprise
// +build enterprise

package v1

import (
	"fmt"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"

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

	queryFields := "rule_id, rule_name, db_type, `desc`, level, type"
	var err error
	var rules []*model.CustomRule
	if req.FilterDBType == "" && req.FilterRuleName == "" {
		rules, err = s.GetCustomRules(queryFields)
	} else {
		rules, err = s.GetCustomRulesByDBTypeAndFuzzyRuleName(queryFields, req.FilterDBType, req.FilterRuleName)
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

func createCustomRule(c echo.Context) error {
	s := model.GetStorage()
	req := new(CreateCustomRuleReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	ruleName := req.RuleName
	dBType := req.DBType

	_, exist, err := s.GetCustomRulesByRuleNameAndDBType(ruleName, dBType)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("rule name:[%v] is exist in %v", ruleName, dBType)))
	}

	uid, err := utils.GenUid()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	ruleId := fmt.Sprintf("rule_id_%v", uid)
	customRule := &model.CustomRule{
		RuleId:     ruleId,
		RuleName:   ruleName,
		DBType:     dBType,
		Desc:       req.Desc,
		Level:      req.Level,
		Typ:        req.Type,
		RuleScript: req.RuleScript,
	}
	err = s.Save(customRule)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

func updateCustomRule(c echo.Context) error {
	req := new(UpdateCustomRuleReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	ruleId := c.Param("rule_id")
	s := model.GetStorage()
	rule, exist, err := s.GetCustomRuleByRuleId(ruleId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("rule id:[%v] is not exist", ruleId)))
	}

	updateMap := map[string]interface{}{}
	if req.RuleName != nil {
		if *req.RuleName != rule.RuleName {
			_, exist, err = s.GetCustomRulesByRuleNameAndDBType(*req.RuleName, rule.DBType)
			if err != nil {
				return controller.JSONBaseErrorReq(c, err)
			}
			if exist {
				return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("rule name:[%v] is exist", *req.RuleName)))
			}
			updateMap["rule_name"] = req.RuleName
		}
	}
	if req.Desc != nil {
		updateMap["desc"] = *req.Desc
	}
	if req.Level != nil {
		updateMap["level"] = *req.Level
	}
	if req.Type != nil {
		if "" == *req.Type {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("type is not allowed to be empty")))
		}
		updateMap["type"] = *req.Type
	}
	if req.RuleScript != nil {
		updateMap["rule_script"] = *req.RuleScript
	}

	err = s.UpdateCustomRuleByRuleId(ruleId, updateMap)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

func getCustomRule(c echo.Context) error {
	ruleId := c.Param("rule_id")

	s := model.GetStorage()
	rule, exist, err := s.GetCustomRuleByRuleId(ruleId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("rule id:[%v] is not exist", ruleId)))
	}
	res := CustomRuleResV1{
		RuleId:     ruleId,
		DBType:     rule.DBType,
		RuleName:   rule.RuleName,
		Desc:       rule.Desc,
		Level:      rule.Level,
		Type:       rule.Typ,
		RuleScript: rule.RuleScript,
	}
	return c.JSON(http.StatusOK, &GetCustomRuleResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    res,
	})
}


func getRuleTypeByDBType(c echo.Context) error {
	dbType := c.Param("db_type")

	s := model.GetStorage()
	allRuleTypes, err := s.GetRuleTypeByDBType(dbType)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	existRuleTypeCount, err := s.GetCustomRuleTypeCountByDBType(dbType)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	dbTypeMap := make(map[string]uint, len(allRuleTypes))
	for i := range allRuleTypes {
		ruleType := allRuleTypes[i]
		count := uint(0)
		if typeCount, exist := existRuleTypeCount[ruleType]; exist {
			count = typeCount
		}
		dbTypeMap[ruleType] = count
	}

	ruleTypeV1s := []RuleTypeV1{}
	for k, v := range dbTypeMap {
		ruleTypeV1s = append(ruleTypeV1s, RuleTypeV1{
			RuleType:  k,
			RuleCount: v,
		})
	}

	return c.JSON(http.StatusOK, &GetRuleTypeByDBTypeResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    ruleTypeV1s,
	})
}
