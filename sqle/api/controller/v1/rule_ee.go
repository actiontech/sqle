//go:build enterprise
// +build enterprise

package v1

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/labstack/echo/v4"
)

func convertCustomRuleToRes(rule *model.CustomRule) CustomRuleResV1 {
	ruleRes := CustomRuleResV1{
		RuleId:     rule.RuleId,
		Desc:       rule.Desc,
		Annotation: rule.Annotation,
		DBType:     rule.DBType,
		Level:      rule.Level,
		Type:       rule.Typ,
		RuleScript: rule.RuleScript,
		Categories: associateCategories(rule.Categories),
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

	queryFields := "rule_id, db_type, `desc`, annotation,level, type"
	var err error
	var rules []*model.CustomRule
	if req.FilterDBType == "" && req.FilterDesc == "" {
		rules, err = s.GetCustomRules(queryFields)
	} else {
		rules, err = s.GetCustomRulesByDBTypeAndFuzzyDesc(queryFields, req.FilterDBType, req.FilterDesc)
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
	_, exist, err := s.GetCustomRuleByRuleId(ruleId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("rule is not exist")))
	}

	err = s.DeleteCustomRule(ruleId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

func checkRuleScript(ruleScript string) error {
	_, err := regexp.Compile(ruleScript)
	return err
}

func createCustomRule(c echo.Context) error {
	s := model.GetStorage()
	req := new(CreateCustomRuleReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	desc := req.Desc
	dBType := req.DBType

	_, exist, err := s.GetCustomRulesByDescAndDBType(desc, dBType)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("desc:[%v] is exist in %v", desc, dBType)))
	}

	err = checkRuleScript(req.RuleScript)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	uid, err := utils.GenUid()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	ruleId := fmt.Sprintf("rule_id_%v", uid)
	customRule := &model.CustomRule{
		RuleId:     ruleId,
		DBType:     dBType,
		Desc:       desc,
		Annotation: req.Annotation,
		Level:      req.Level,
		Typ:        req.Type,
		RuleScript: req.RuleScript,
	}
	err = s.Save(customRule)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = createCustomRuleCategoryRels(ruleId, req.Tags)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

// createCustomRuleCategoryRels 创建自定义规则的分类关系
func createCustomRuleCategoryRels(customRuleId string, tags *[]string) error {
	if tags == nil {
		return nil
	}
	s := model.GetStorage()
	ruleCategories, err := s.GetAuditRuleCategoryByTagIn(*tags)
	if err != nil {
		return err
	}
	for _, category := range ruleCategories {
		customerCategoryRel := model.CustomRuleCategoryRel{CategoryId: category.ID, CustomRuleId: customRuleId}
		err = s.Save(customerCategoryRel)
		if err != nil {
			return err
		}
	}
	return err
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
	if req.Desc != nil {
		if *req.Desc != rule.Desc {
			if "" == *req.Desc {
				return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("desc is not allowed to be empty")))
			}
			_, exist, err = s.GetCustomRulesByDescAndDBType(*req.Desc, rule.DBType)
			if err != nil {
				return controller.JSONBaseErrorReq(c, err)
			}
			if exist {
				return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("rule name:[%v] is exist", *req.Desc)))
			}
			updateMap["desc"] = req.Desc
		}
	}
	if req.Annotation != nil {
		updateMap["annotation"] = *req.Annotation
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
		err := checkRuleScript(*req.RuleScript)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		updateMap["rule_script"] = *req.RuleScript
	}

	err = s.UpdateCustomRuleByRuleId(ruleId, updateMap)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = s.UpdateCustomRuleCategoriesByRuleId(ruleId, *req.Tags)
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
		Desc:       rule.Desc,
		Annotation: rule.Annotation,
		Level:      rule.Level,
		Type:       rule.Typ,
		RuleScript: rule.RuleScript,
		Categories: associateCategories(rule.Categories),
	}
	return c.JSON(http.StatusOK, &GetCustomRuleResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    res,
	})
}

func getRuleTypeByDBType(c echo.Context) error {
	dbType := c.Param("db_type")

	s := model.GetStorage()
	ctx := c.Request().Context()
	allRuleTypes, err := s.GetRuleTypeByDBType(ctx, dbType)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	existRuleTypeCount, err := s.GetCustomRuleTypeCountByDBType(dbType)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	dbTypeMap := make(map[string]uint, len(allRuleTypes))
	newRuleTypes := []model.CustomTypeCount{}
	for i := range allRuleTypes {
		ruleType := allRuleTypes[i]
		dbTypeMap[ruleType] = uint(0)
	}
	for _, ruleType := range existRuleTypeCount {
		_, exist := dbTypeMap[ruleType.Type]
		if exist {
			dbTypeMap[ruleType.Type] = ruleType.TypeCount
		} else {
			newRuleTypes = append(newRuleTypes, *ruleType)
		}
	}

	ruleTypeV1s := []RuleTypeV1{}
	for k, v := range dbTypeMap {
		ruleTypeV1s = append(ruleTypeV1s, RuleTypeV1{
			RuleType:  k,
			RuleCount: v,
		})
	}
	for _, ruleType := range newRuleTypes {
		ruleTypeV1s = append(ruleTypeV1s, RuleTypeV1{
			RuleType:         ruleType.Type,
			RuleCount:        ruleType.TypeCount,
			IsCustomRuleType: true,
		})
	}

	return c.JSON(http.StatusOK, &GetRuleTypeByDBTypeResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    ruleTypeV1s,
	})
}
