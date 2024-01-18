package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"strings"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
)

var ErrRuleTemplateNotExist = errors.New(errors.DataNotExist, fmt.Errorf("rule template not exist"))

type CreateRuleTemplateReqV1 struct {
	Name     string      `json:"rule_template_name" valid:"required,name"`
	Desc     string      `json:"desc"`
	DBType   string      `json:"db_type" valid:"required"`
	RuleList []RuleReqV1 `json:"rule_list" form:"rule_list" valid:"required,dive,required"`
}

type RuleReqV1 struct {
	Name         string           `json:"name" form:"name" valid:"required" example:"ddl_check_index_count"`
	Level        string           `json:"level" form:"level" valid:"required" example:"error"`
	Params       []RuleParamReqV1 `json:"params" form:"params" valid:"dive,required"`
	IsCustomRule bool             `json:"is_custom_rule" form:"is_custom_rule"`
}

type RuleParamReqV1 struct {
	Key   string `json:"key" form:"key" valid:"required"`
	Value string `json:"value" form:"value" valid:"required"`
}

func checkAndGenerateRules(rulesReq []RuleReqV1, template *model.RuleTemplate) ([]model.RuleTemplateRule, error) {
	s := model.GetStorage()
	var rules map[string]model.Rule
	var err error

	ruleNames := make([]string, 0, len(rulesReq))
	for _, r := range rulesReq {
		ruleNames = append(ruleNames, r.Name)
	}
	rules, err = s.GetAndCheckRuleExist(ruleNames, template.DBType)
	if err != nil {
		return nil, err
	}

	templateRules := make([]model.RuleTemplateRule, 0, len(rulesReq))
	for _, r := range rulesReq {
		rule := rules[r.Name]
		params := rule.Params

		// check request params is equal rule params.
		if len(r.Params) != len(params) {
			reqParamsKey := make([]string, 0, len(r.Params))
			for _, p := range r.Params {
				reqParamsKey = append(reqParamsKey, p.Key)
			}
			paramsKey := make([]string, 0, len(params))
			for _, p := range params {
				paramsKey = append(paramsKey, p.Key)
			}
			return nil, fmt.Errorf("request rule \"%s'| params key is [%s], but need [%s]",
				r.Name, reqParamsKey, paramsKey)
		}
		for _, p := range r.Params {
			// set and valid param.
			err := params.SetParamValue(p.Key, p.Value)
			if err != nil {
				return nil, fmt.Errorf("set rule %s param error: %s", r.Name, err)
			}
		}
		templateRules = append(templateRules, model.NewRuleTemplateRule(template, &model.Rule{
			Name:   r.Name,
			Level:  r.Level,
			DBType: template.DBType,
			Params: params,
		}))
	}
	return templateRules, nil
}

func checkAndGenerateCustomRules(rulesReq []RuleReqV1, template *model.RuleTemplate) ([]model.RuleTemplateCustomRule, error) {
	s := model.GetStorage()
	var err error

	ruleNames := make([]string, 0, len(rulesReq))
	for _, r := range rulesReq {
		ruleNames = append(ruleNames, r.Name)
	}
	_, err = s.GetAndCheckCustomRuleExist(ruleNames)
	if err != nil {
		return nil, err
	}

	templateCustomRules := make([]model.RuleTemplateCustomRule, 0, len(rulesReq))
	for _, r := range rulesReq {
		templateCustomRules = append(templateCustomRules, model.NewRuleTemplateCustomRule(template, &model.CustomRule{
			RuleId: r.Name,
			Level:  r.Level,
			DBType: template.DBType,
		}))
	}
	return templateCustomRules, nil
}

func checkAndGenerateAllTypeRules(rulesReq []RuleReqV1, template *model.RuleTemplate) ([]model.RuleTemplateRule,
	[]model.RuleTemplateCustomRule, error) {

	var err error
	rules := []RuleReqV1{}
	customRules := []RuleReqV1{}
	for i := range rulesReq {
		if rulesReq[i].IsCustomRule {
			customRules = append(customRules, rulesReq[i])
		} else {
			rules = append(rules, rulesReq[i])
		}
	}

	templateRules := make([]model.RuleTemplateRule, 0, len(rules))
	templateCustomRules := make([]model.RuleTemplateCustomRule, 0, len(rules))
	if len(rules) > 0 {
		templateRules, err = checkAndGenerateRules(rules, template)
		if err != nil {
			return nil, nil, err
		}
	}
	if len(customRules) > 0 {
		templateCustomRules, err = checkAndGenerateCustomRules(customRules, template)
		if err != nil {
			return nil, nil, err
		}
	}
	return templateRules, templateCustomRules, nil
}

// @Summary 添加全局规则模板
// @Description create a global rule template
// @Id createRuleTemplateV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Accept json
// @Param req body v1.CreateRuleTemplateReqV1 true "add rule template request"
// @Success 200 {object} controller.BaseRes
// @router /v1/rule_templates [post]
func CreateRuleTemplate(c echo.Context) error {
	req := new(CreateRuleTemplateReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()
	exist, err := s.IsRuleTemplateExistFromAnyProject(model.ProjectIdForGlobalRuleTemplate, req.Name)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("rule template is exist")))
	}

	ruleTemplate := &model.RuleTemplate{
		ProjectId: model.ProjectIdForGlobalRuleTemplate,
		Name:      req.Name,
		Desc:      req.Desc,
		DBType:    req.DBType,
	}
	templateRules := []model.RuleTemplateRule{}
	templateCustomRules := []model.RuleTemplateCustomRule{}
	if len(req.RuleList) > 0 {
		templateRules, templateCustomRules, err = checkAndGenerateAllTypeRules(req.RuleList, ruleTemplate)
		if err != nil {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
		}
	}

	err = s.Save(ruleTemplate)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = s.UpdateRuleTemplateRules(ruleTemplate, templateRules...)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = s.UpdateRuleTemplateCustomRules(ruleTemplate, templateCustomRules...)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type UpdateRuleTemplateReqV1 struct {
	Desc     *string     `json:"desc"`
	RuleList []RuleReqV1 `json:"rule_list" form:"rule_list" valid:"dive,required"`
}

// @Summary 更新全局规则模板
// @Description update global rule template
// @Id updateRuleTemplateV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Param rule_template_name path string true "rule template name"
// @Param req body v1.UpdateRuleTemplateReqV1 true "update rule template request"
// @Success 200 {object} controller.BaseRes
// @router /v1/rule_templates/{rule_template_name}/ [patch]
func UpdateRuleTemplate(c echo.Context) error {
	templateName := c.Param("rule_template_name")
	req := new(UpdateRuleTemplateReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()
	template, exist, err := s.GetRuleTemplateByProjectIdAndName(model.ProjectIdForGlobalRuleTemplate, templateName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("rule template is not exist")))
	}

	templateRules := []model.RuleTemplateRule{}
	templateCustomRules := []model.RuleTemplateCustomRule{}
	if len(req.RuleList) > 0 {
		templateRules, templateCustomRules, err = checkAndGenerateAllTypeRules(req.RuleList, template)
		if err != nil {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
		}
	}

	if req.Desc != nil {
		template.Desc = *req.Desc
		err = s.Save(&template)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}
	if req.RuleList != nil {
		err = s.UpdateRuleTemplateRules(template, templateRules...)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		err = s.UpdateRuleTemplateCustomRules(template, templateCustomRules...)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type GetRuleTemplateReqV1 struct {
	FuzzyKeywordRule string `json:"fuzzy_keyword_rule" query:"fuzzy_keyword_rule"`
}

type GetRuleTemplateResV1 struct {
	controller.BaseRes
	Data *RuleTemplateDetailResV1 `json:"data"`
}

type RuleTemplateDetailResV1 struct {
	Name     string      `json:"rule_template_name"`
	Desc     string      `json:"desc"`
	DBType   string      `json:"db_type"`
	RuleList []RuleResV1 `json:"rule_list,omitempty"`
}

func convertRuleTemplateToRes(template *model.RuleTemplate) *RuleTemplateDetailResV1 {

	ruleList := make([]RuleResV1, 0, len(template.RuleList))
	for _, r := range template.RuleList {
		if r.Rule == nil {
			continue
		}
		ruleList = append(ruleList, convertRuleToRes(r.GetRule()))
	}
	for _, r := range template.CustomRuleList {
		if r.CustomRule == nil {
			continue
		}
		ruleList = append(ruleList, convertCustomRuleToRuleResV1(r.GetRule()))
	}
	return &RuleTemplateDetailResV1{
		Name:     template.Name,
		Desc:     template.Desc,
		DBType:   template.DBType,
		RuleList: ruleList,
	}
}

// @Summary 获取全局规则模板信息
// @Description get global rule template
// @Id getRuleTemplateV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Param rule_template_name path string true "rule template name"
// @Param fuzzy_keyword_rule query string false "fuzzy rule,keyword for desc and annotation"
// @Success 200 {object} v1.GetRuleTemplateResV1
// @router /v1/rule_templates/{rule_template_name}/ [get]
func GetRuleTemplate(c echo.Context) error {
	s := model.GetStorage()
	templateName := c.Param("rule_template_name")
	req := new(GetRuleTemplateReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	template, exist, err := s.GetRuleTemplateDetailByNameAndProjectIds([]string{model.ProjectIdForGlobalRuleTemplate}, templateName, req.FuzzyKeywordRule)
	if err != nil {
		return c.JSON(200, controller.NewBaseReq(err))
	}
	if !exist {
		return c.JSON(200, controller.NewBaseReq(errors.New(errors.DataNotExist,
			fmt.Errorf("rule template is not exist"))))
	}

	return c.JSON(http.StatusOK, &GetRuleTemplateResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertRuleTemplateToRes(template),
	})
}

// @Summary 删除全局规则模板
// @Description delete global rule template
// @Id deleteRuleTemplateV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Param rule_template_name path string true "rule template name"
// @Success 200 {object} controller.BaseRes
// @router /v1/rule_templates/{rule_template_name}/ [delete]
func DeleteRuleTemplate(c echo.Context) error {
	s := model.GetStorage()
	templateName := c.Param("rule_template_name")
	template, exist, err := s.GetRuleTemplateByProjectIdAndName(model.ProjectIdForGlobalRuleTemplate, templateName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("rule template is not exist")))
	}

	// check audit plans
	{
		auditPlanNames, err := s.GetAuditPlanNamesByRuleTemplate(templateName)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}

		if len(auditPlanNames) > 0 {
			err = errors.NewDataInvalidErr("rule_templates[%v] is still in use, related audit_plan[%v]",
				templateName, strings.Join(auditPlanNames, ", "))
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	// check instance
	{
		instanceNames, err := dms.GetInstancesNameByRuleTemplateName(c.Request().Context(), templateName)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if len(instanceNames) > 0 {
			err = errors.NewDataInvalidErr("rule_templates[%v] is still in use, related instances[%v]",
				templateName, strings.Join(instanceNames, ", "))
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	err = s.Delete(template)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type GetRuleTemplatesReqV1 struct {
	PageIndex uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetRuleTemplatesResV1 struct {
	controller.BaseRes
	Data      []RuleTemplateResV1 `json:"data"`
	TotalNums uint64              `json:"total_nums"`
}

type RuleTemplateResV1 struct {
	Name   string `json:"rule_template_name"`
	Desc   string `json:"desc"`
	DBType string `json:"db_type"`
}

// @Summary 全局规则模板列表
// @Description get all global rule template
// @Id getRuleTemplateListV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Success 200 {object} v1.GetRuleTemplatesResV1
// @router /v1/rule_templates [get]
func GetRuleTemplates(c echo.Context) error {
	req := new(GetRuleTemplatesReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	limit, offset := controller.GetLimitAndOffset(req.PageIndex, req.PageSize)

	s := model.GetStorage()
	ruleTemplates, count, err := getRuleTemplatesByReq(s, limit, offset, model.ProjectIdForGlobalRuleTemplate)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetRuleTemplatesResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      convertRuleTemplatesToRes(ruleTemplates),
		TotalNums: count,
	})
}

func getRuleTemplatesByReq(s *model.Storage, limit, offset uint32, projectId string) (ruleTemplates []*model.RuleTemplateDetail, count uint64, err error) {
	data := map[string]interface{}{
		"limit":      limit,
		"offset":     offset,
		"project_id": projectId,
	}
	ruleTemplates, count, err = s.GetRuleTemplatesByReq(data)
	if err != nil {
		return nil, 0, err
	}
	return
}

func convertRuleTemplatesToRes(ruleTemplates []*model.RuleTemplateDetail) []RuleTemplateResV1 {

	ruleTemplatesReq := make([]RuleTemplateResV1, 0, len(ruleTemplates))
	for _, ruleTemplate := range ruleTemplates {
		ruleTemplateReq := RuleTemplateResV1{
			Name:   ruleTemplate.Name,
			Desc:   ruleTemplate.Desc,
			DBType: ruleTemplate.DBType,
		}
		ruleTemplatesReq = append(ruleTemplatesReq, ruleTemplateReq)
	}
	return ruleTemplatesReq
}

type GetRulesReqV1 struct {
	FilterDBType                 string `json:"filter_db_type" query:"filter_db_type"`
	FilterGlobalRuleTemplateName string `json:"filter_global_rule_template_name" query:"filter_global_rule_template_name"`
	FilterRuleNames              string `json:"filter_rule_names" query:"filter_rule_names"`
	FuzzyKeywordRule             string `json:"fuzzy_keyword_rule" query:"fuzzy_keyword_rule"`
}

type GetRulesResV1 struct {
	controller.BaseRes
	Data []RuleResV1 `json:"data"`
}

type RuleResV1 struct {
	Name         string           `json:"rule_name"`
	Desc         string           `json:"desc"`
	Annotation   string           `json:"annotation" example:"避免多次 table rebuild 带来的消耗、以及对线上业务的影响"`
	Level        string           `json:"level" example:"error" enums:"normal,notice,warn,error"`
	Typ          string           `json:"type" example:"全局配置" `
	DBType       string           `json:"db_type" example:"mysql"`
	Params       []RuleParamResV1 `json:"params,omitempty"`
	IsCustomRule bool             `json:"is_custom_rule"`
}

type RuleParamResV1 struct {
	Key   string `json:"key" form:"key"`
	Value string `json:"value" form:"value"`
	Desc  string `json:"desc" form:"desc"`
	Type  string `json:"type" form:"type" enums:"string,int,bool"`
}

func convertRuleToRes(rule *model.Rule) RuleResV1 {
	ruleRes := RuleResV1{
		Name:       rule.Name,
		Desc:       rule.Desc,
		Annotation: rule.Annotation,
		Level:      rule.Level,
		Typ:        rule.Typ,
		DBType:     rule.DBType,
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

func convertCustomRuleToRuleResV1(rule *model.CustomRule) RuleResV1 {
	ruleRes := RuleResV1{
		Name:         rule.RuleId,
		Desc:         rule.Desc,
		Annotation:   rule.Annotation,
		Level:        rule.Level,
		Typ:          rule.Typ,
		DBType:       rule.DBType,
		IsCustomRule: true,
	}
	return ruleRes
}

func convertRulesToRes(rules interface{}) []RuleResV1 {
	rulesRes := []RuleResV1{}
	switch ruleSlice := rules.(type) {
	case []*model.Rule:
		for _, rule := range ruleSlice {
			rulesRes = append(rulesRes, convertRuleToRes(rule))
		}
	case []*model.CustomRule:
		for _, rule := range ruleSlice {
			rulesRes = append(rulesRes, convertCustomRuleToRuleResV1(rule))
		}
	}
	return rulesRes
}

// @Summary 规则列表
// @Description get all rule template
// @Id getRuleListV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Param filter_db_type query string false "filter db type"
// @Param fuzzy_keyword_rule query string false "fuzzy rule,keyword for desc and annotation"
// @Param filter_global_rule_template_name query string false "filter global rule template name"
// @Param filter_rule_names query string false "filter rule name list"
// @Success 200 {object} v1.GetRulesResV1
// @router /v1/rules [get]
func GetRules(c echo.Context) error {
	req := new(GetRulesReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()
	var rules []*model.Rule
	var customRules []*model.CustomRule
	var err error
	rules, err = s.GetRulesByReq(map[string]interface{}{
		"filter_global_rule_template_name": req.FilterGlobalRuleTemplateName,
		"filter_db_type":                   req.FilterDBType,
		"filter_rule_names":                req.FilterRuleNames,
		"fuzzy_keyword_rule":               req.FuzzyKeywordRule,
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	customRules, err = s.GetCustomRulesByReq(map[string]interface{}{
		"filter_global_rule_template_name": req.FilterGlobalRuleTemplateName,
		"filter_db_type":                   req.FilterDBType,
		"filter_rule_names":                req.FilterRuleNames,
		"fuzzy_keyword_rule":               req.FuzzyKeywordRule,
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	ruleRes := convertRulesToRes(rules)
	customRuleRes := convertRulesToRes(customRules)
	ruleRes = append(ruleRes, customRuleRes...)
	return c.JSON(http.StatusOK, &GetRulesResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    ruleRes,
	})
}

type RuleTemplateTipReqV1 struct {
	FilterDBType string `json:"filter_db_type" query:"filter_db_type"`
}

type RuleTemplateTipResV1 struct {
	ID     string `json:"rule_template_id"`
	Name   string `json:"rule_template_name"`
	DBType string `json:"db_type"`
}

type GetRuleTemplateTipsResV1 struct {
	controller.BaseRes
	Data []RuleTemplateTipResV1 `json:"data"`
}

// @Summary 获取全局规则模板提示
// @Description get global rule template tips
// @Id getRuleTemplateTipsV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Param filter_db_type query string false "filter db type"
// @Success 200 {object} v1.GetRuleTemplateTipsResV1
// @router /v1/rule_template_tips [get]
func GetRuleTemplateTips(c echo.Context) error {
	req := new(RuleTemplateTipReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	return getRuleTemplateTips(c, model.ProjectIdForGlobalRuleTemplate, req.FilterDBType)
}

func getRuleTemplateTips(c echo.Context, projectId string, filterDBType string) error {
	s := model.GetStorage()
	ruleTemplates, err := s.GetRuleTemplateTips(projectId, filterDBType)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	ruleTemplateTipsRes := make([]RuleTemplateTipResV1, 0, len(ruleTemplates))
	for _, roleTemplate := range ruleTemplates {
		ruleTemplateTipRes := RuleTemplateTipResV1{
			ID:     roleTemplate.GetIDStr(),
			Name:   roleTemplate.Name,
			DBType: roleTemplate.DBType,
		}
		ruleTemplateTipsRes = append(ruleTemplateTipsRes, ruleTemplateTipRes)
	}
	return c.JSON(http.StatusOK, &GetRuleTemplateTipsResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    ruleTemplateTipsRes,
	})
}

type CloneRuleTemplateReqV1 struct {
	Name string `json:"new_rule_template_name" valid:"required"`
	Desc string `json:"desc"`
}

// @Summary 克隆全局规则模板
// @Description clone a rule template
// @Id CloneRuleTemplateV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Accept json
// @Param rule_template_name path string true "rule template name"
// @Param req body v1.CloneRuleTemplateReqV1 true "clone rule template request"
// @Success 200 {object} controller.BaseRes
// @router /v1/rule_templates/{rule_template_name}/clone [post]
func CloneRuleTemplate(c echo.Context) error {
	req := new(CloneRuleTemplateReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()
	exist, err := s.IsRuleTemplateExistFromAnyProject(model.ProjectIdForGlobalRuleTemplate, req.Name)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("rule template is exist")))
	}

	sourceTplName := c.Param("rule_template_name")
	sourceTpl, exist, err := s.GetRuleTemplateDetailByNameAndProjectIds([]string{model.ProjectIdForGlobalRuleTemplate}, sourceTplName, "")
	if err != nil {
		return c.JSON(200, controller.NewBaseReq(err))
	}
	if !exist {
		return c.JSON(200, controller.NewBaseReq(errors.New(errors.DataNotExist,
			fmt.Errorf("source rule template %s is not exist", sourceTplName))))
	}

	ruleTemplate := &model.RuleTemplate{
		ProjectId: model.ProjectIdForGlobalRuleTemplate,
		Name:      req.Name,
		Desc:      req.Desc,
		DBType:    sourceTpl.DBType,
	}
	err = s.Save(ruleTemplate)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = s.CloneRuleTemplateRules(sourceTpl, ruleTemplate)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = s.CloneRuleTemplateCustomRules(sourceTpl, ruleTemplate)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

func CheckRuleTemplateCanBeBindEachInstance(s *model.Storage, tplName string, instances []*model.Instance) error {
	for _, inst := range instances {
		currentBindTemplates, err := s.GetRuleTemplatesByInstance(inst)
		if err != nil {
			return err
		}
		if len(currentBindTemplates) > 1 {
			return errInstanceBind
		}
		if len(currentBindTemplates) == 1 && currentBindTemplates[0].Name != tplName {
			return errInstanceBind
		}
	}
	return nil
}

type CreateProjectRuleTemplateReqV1 struct {
	Name     string      `json:"rule_template_name" valid:"required,name"`
	Desc     string      `json:"desc"`
	DBType   string      `json:"db_type" valid:"required"`
	RuleList []RuleReqV1 `json:"rule_list" form:"rule_list" valid:"required,dive,required"`
}

// CreateProjectRuleTemplate
// @Summary 添加项目规则模板
// @Description create a rule template in project
// @Id createProjectRuleTemplateV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Accept json
// @Param project_name path string true "project name"
// @Param req body v1.CreateProjectRuleTemplateReqV1 true "add rule template request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/rule_templates [post]
func CreateProjectRuleTemplate(c echo.Context) error {
	req := new(CreateProjectRuleTemplateReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	_, exist, err := s.GetGlobalAndProjectRuleTemplateByNameAndProjectId(req.Name, projectUid)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("rule template is exist")))
	}

	ruleTemplate := &model.RuleTemplate{
		ProjectId: model.ProjectUID(projectUid),
		Name:      req.Name,
		Desc:      req.Desc,
		DBType:    req.DBType,
	}
	templateRules := []model.RuleTemplateRule{}
	templateCustomRules := []model.RuleTemplateCustomRule{}
	if len(req.RuleList) > 0 {
		templateRules, templateCustomRules, err = checkAndGenerateAllTypeRules(req.RuleList, ruleTemplate)
		if err != nil {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
		}
	}

	// var instances []*model.Instance
	// if len(req.Instances) > 0 {
	// 	instances, err = s.GetAndCheckInstanceExist(req.Instances, projectName)
	// 	if err != nil {
	// 		return controller.JSONBaseErrorReq(c, err)
	// 	}
	// }

	// err = CheckRuleTemplateCanBeBindEachInstance(s, req.Name, instances)
	// if err != nil {
	// 	return controller.JSONBaseErrorReq(c, err)
	// }

	// err = CheckInstanceAndRuleTemplateDbType([]*model.RuleTemplate{ruleTemplate}, instances...)
	// if err != nil {
	// 	return controller.JSONBaseErrorReq(c, err)
	// }

	err = s.Save(ruleTemplate)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = s.UpdateRuleTemplateRules(ruleTemplate, templateRules...)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = s.UpdateRuleTemplateCustomRules(ruleTemplate, templateCustomRules...)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// TODO SQLE会移除instance参数
	// err = s.UpdateRuleTemplateInstances(ruleTemplate, instances...)
	// if err != nil {
	// 	return controller.JSONBaseErrorReq(c, err)
	// }
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type UpdateProjectRuleTemplateReqV1 struct {
	Desc     *string     `json:"desc"`
	RuleList []RuleReqV1 `json:"rule_list" form:"rule_list" valid:"dive,required"`
}

// UpdateProjectRuleTemplate
// @Summary 更新项目规则模板
// @Description update rule template in project
// @Id updateProjectRuleTemplateV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param rule_template_name path string true "rule template name"
// @Param req body v1.UpdateProjectRuleTemplateReqV1 true "update rule template request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/rule_templates/{rule_template_name}/ [patch]
func UpdateProjectRuleTemplate(c echo.Context) error {
	templateName := c.Param("rule_template_name")
	req := new(UpdateProjectRuleTemplateReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	s := model.GetStorage()
	template, exist, err := s.GetGlobalAndProjectRuleTemplateByNameAndProjectId(templateName, projectUid)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("rule template is not exist")))
	}

	if template.ProjectId == model.ProjectIdForGlobalRuleTemplate {
		return controller.JSONBaseErrorReq(c, errors.New(errors.UserNotPermission, fmt.Errorf("you cannot update a global template from this api")))
	}

	templateRules := []model.RuleTemplateRule{}
	templateCustomRules := []model.RuleTemplateCustomRule{}
	if len(req.RuleList) > 0 {
		templateRules, templateCustomRules, err = checkAndGenerateAllTypeRules(req.RuleList, template)
		if err != nil {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
		}
	}

	if req.Desc != nil {
		template.Desc = *req.Desc
		err = s.Save(&template)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}
	if req.RuleList != nil {
		err = s.UpdateRuleTemplateRules(template, templateRules...)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		err = s.UpdateRuleTemplateCustomRules(template, templateCustomRules...)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type GetProjectRuleTemplateResV1 struct {
	controller.BaseRes
	Data *RuleProjectTemplateDetailResV1 `json:"data"`
}

type RuleProjectTemplateDetailResV1 struct {
	Name     string      `json:"rule_template_name"`
	Desc     string      `json:"desc"`
	DBType   string      `json:"db_type"`
	RuleList []RuleResV1 `json:"rule_list,omitempty"`
}

type ProjectRuleTemplateInstance struct {
	Name string `json:"name"`
}

// GetProjectRuleTemplate
// @Summary 获取项目规则模板信息
// @Description get rule template detail in project
// @Id getProjectRuleTemplateV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param rule_template_name path string true "rule template name"
// @Param fuzzy_keyword_rule query string false "fuzzy rule,keyword for desc and annotation"
// @Success 200 {object} v1.GetProjectRuleTemplateResV1
// @router /v1/projects/{project_name}/rule_templates/{rule_template_name}/ [get]
func GetProjectRuleTemplate(c echo.Context) error {
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	s := model.GetStorage()
	templateName := c.Param("rule_template_name")
	req := new(GetRuleTemplateReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	template, exist, err := s.GetRuleTemplateDetailByNameAndProjectIds([]string{projectUid}, templateName, req.FuzzyKeywordRule)
	if err != nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
	}
	if !exist {
		return c.JSON(http.StatusOK, controller.NewBaseReq(errors.New(errors.DataNotExist,
			fmt.Errorf("rule template is not exist"))))
	}

	return c.JSON(http.StatusOK, &GetProjectRuleTemplateResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertProjectRuleTemplateToRes(template),
	})
}

func convertProjectRuleTemplateToRes(template *model.RuleTemplate) *RuleProjectTemplateDetailResV1 {
	ruleList := make([]RuleResV1, 0, len(template.RuleList))
	for _, r := range template.RuleList {
		if r.Rule == nil {
			continue
		}
		ruleList = append(ruleList, convertRuleToRes(r.GetRule()))
	}
	for _, r := range template.CustomRuleList {
		if r.CustomRule == nil {
			continue
		}
		ruleList = append(ruleList, convertCustomRuleToRuleResV1(r.GetRule()))
	}
	return &RuleProjectTemplateDetailResV1{
		Name:     template.Name,
		Desc:     template.Desc,
		DBType:   template.DBType,
		RuleList: ruleList,
	}
}

// DeleteProjectRuleTemplate
// @Summary 删除项目规则模板
// @Description delete rule template in project
// @Id deleteProjectRuleTemplateV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param rule_template_name path string true "rule template name"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/rule_templates/{rule_template_name}/ [delete]
func DeleteProjectRuleTemplate(c echo.Context) error {
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()

	templateName := c.Param("rule_template_name")
	template, exist, err := s.GetGlobalAndProjectRuleTemplateByNameAndProjectId(templateName, projectUid)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("rule template is not exist")))
	}

	if template.ProjectId == model.ProjectIdForGlobalRuleTemplate {
		return controller.JSONBaseErrorReq(c, errors.New(errors.UserNotPermission, fmt.Errorf("you cannot delete a global template from this api")))
	}

	// check audit plans
	{
		auditPlanNames, err := s.GetAuditPlanNamesByRuleTemplateAndProject(templateName, projectUid)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}

		if len(auditPlanNames) > 0 {
			err = errors.NewDataInvalidErr("rule_templates[%v] is still in use, related audit_plan[%v]",
				templateName, strings.Join(auditPlanNames, ", "))
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	// check instance
	{
		instanceNames, err := dms.GetInstancesNameInProjectByRuleTemplateName(c.Request().Context(), projectUid, templateName)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if len(instanceNames) > 0 {
			err = errors.NewDataInvalidErr("rule_templates[%v] is still in use, related instances[%v]",
				templateName, strings.Join(instanceNames, ", "))
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	err = s.Delete(template)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type GetProjectRuleTemplatesResV1 struct {
	controller.BaseRes
	Data      []ProjectRuleTemplateResV1 `json:"data"`
	TotalNums uint64                     `json:"total_nums"`
}

type ProjectRuleTemplateResV1 struct {
	Name   string `json:"rule_template_name"`
	Desc   string `json:"desc"`
	DBType string `json:"db_type"`
}

// GetProjectRuleTemplates
// @Summary 项目规则模板列表
// @Description get all rule template in a project
// @Id getProjectRuleTemplateListV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Success 200 {object} v1.GetProjectRuleTemplatesResV1
// @router /v1/projects/{project_name}/rule_templates [get]
func GetProjectRuleTemplates(c echo.Context) error {
	req := new(GetRuleTemplatesReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	s := model.GetStorage()
	limit, offset := controller.GetLimitAndOffset(req.PageIndex, req.PageSize)
	ruleTemplates, count, err := getRuleTemplatesByReq(s, limit, offset, projectUid)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetProjectRuleTemplatesResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      convertProjectRuleTemplatesToRes(ruleTemplates),
		TotalNums: count,
	})
}

func convertProjectRuleTemplatesToRes(ruleTemplates []*model.RuleTemplateDetail) []ProjectRuleTemplateResV1 {
	ruleTemplatesRes := make([]ProjectRuleTemplateResV1, len(ruleTemplates))
	for i, t := range ruleTemplates {
		// instances := make([]*ProjectRuleTemplateInstance, len(t.InstanceNames))
		// for j, instName := range t.InstanceNames {
		// 	instances[j] = &ProjectRuleTemplateInstance{Name: instName}
		// }
		ruleTemplatesRes[i] = ProjectRuleTemplateResV1{
			Name:   t.Name,
			Desc:   t.Desc,
			DBType: t.DBType,
			// Instances: instances,
		}
	}
	return ruleTemplatesRes
}

type CloneProjectRuleTemplateReqV1 struct {
	Name string `json:"new_rule_template_name" valid:"required"`
	Desc string `json:"desc"`
}

// CloneProjectRuleTemplate
// @Summary 克隆项目规则模板
// @Description clone a rule template in project
// @Id cloneProjectRuleTemplateV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Accept json
// @Param project_name path string true "project name"
// @Param rule_template_name path string true "rule template name"
// @Param req body v1.CloneProjectRuleTemplateReqV1 true "clone rule template request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/rule_templates/{rule_template_name}/clone [post]
func CloneProjectRuleTemplate(c echo.Context) error {
	req := new(CloneProjectRuleTemplateReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	exist, err := s.IsRuleTemplateExistFromAnyProject(model.ProjectUID(projectUid), req.Name)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("rule template is exist")))
	}

	sourceTplName := c.Param("rule_template_name")
	sourceTpl, exist, err := s.GetRuleTemplateDetailByNameAndProjectIds([]string{projectUid}, sourceTplName, "")
	if err != nil {
		return c.JSON(200, controller.NewBaseReq(err))
	}
	if !exist {
		return c.JSON(200, controller.NewBaseReq(errors.New(errors.DataNotExist,
			fmt.Errorf("source rule template %s is not exist", sourceTplName))))
	}

	// var instances []*model.Instance
	// if len(req.Instances) > 0 {
	// 	instances, err = s.GetAndCheckInstanceExist(req.Instances, projectName)
	// 	if err != nil {
	// 		return controller.JSONBaseErrorReq(c, err)
	// 	}
	// }

	// err = CheckRuleTemplateCanBeBindEachInstance(s, req.Name, instances)
	// if err != nil {
	// 	return controller.JSONBaseErrorReq(c, err)
	// }

	// err = CheckInstanceAndRuleTemplateDbType([]*model.RuleTemplate{sourceTpl}, instances...)
	// if err != nil {
	// 	return controller.JSONBaseErrorReq(c, err)
	// }

	ruleTemplate := &model.RuleTemplate{
		ProjectId: model.ProjectUID(projectUid),
		Name:      req.Name,
		Desc:      req.Desc,
		DBType:    sourceTpl.DBType,
	}
	err = s.Save(ruleTemplate)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = s.CloneRuleTemplateRules(sourceTpl, ruleTemplate)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = s.CloneRuleTemplateCustomRules(sourceTpl, ruleTemplate)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

// GetProjectRuleTemplateTips
// @Summary 获取项目规则模板提示
// @Description get rule template tips in project
// @Id getProjectRuleTemplateTipsV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param filter_db_type query string false "filter db type"
// @Success 200 {object} v1.GetRuleTemplateTipsResV1
// @router /v1/projects/{project_name}/rule_template_tips [get]
func GetProjectRuleTemplateTips(c echo.Context) error {
	req := new(RuleTemplateTipReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return getRuleTemplateTips(c, projectUid, req.FilterDBType)
}

type ParseProjectRuleTemplateFileResV1 struct {
	controller.BaseRes
	Data ParseProjectRuleTemplateFileResDataV1 `json:"data"`
}

type ParseProjectRuleTemplateFileResDataV1 struct {
	Name     string      `json:"name"`
	Desc     string      `json:"desc"`
	DBType   string      `json:"db_type"`
	RuleList []RuleResV1 `json:"rule_list"`
}

// ParseProjectRuleTemplateFile parse rule template
// @Summary 解析规则模板文件
// @Description parse rule template
// @Id importProjectRuleTemplateV1
// @Tags rule_template
// @Accept mpfd
// @Security ApiKeyAuth
// @Param rule_template_file formData file true "SQLE rule template file"
// @Success 200 {object} v1.ParseProjectRuleTemplateFileResV1
// @router /v1/rule_templates/parse [post]
func ParseProjectRuleTemplateFile(c echo.Context) error {
	// 读取+解析文件
	_, file, exist, err := controller.ReadFile(c, "rule_template_file")
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("the file has not been uploaded or the key bound to the file is not 'rule_template_file'")))
	}

	resp := ParseProjectRuleTemplateFileResDataV1{}
	err = json.Unmarshal([]byte(file), &resp)
	if err != nil {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("the file format is incorrect. Please check the uploaded file, error: %v", err))
	}

	return c.JSON(http.StatusOK, ParseProjectRuleTemplateFileResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    resp,
	})
}

// ExportRuleTemplateFile
// @Summary 导出全局规则模板
// @Description export rule template
// @Id exportRuleTemplateV1
// @Tags rule_template
// @Param rule_template_name path string true "rule template name"
// @Security ApiKeyAuth
// @Success 200 file 1 "sqle rule template file"
// @router /v1/rule_templates/{rule_template_name}/export [get]
func ExportRuleTemplateFile(c echo.Context) error {
	templateName := c.Param("rule_template_name")
	return exportRuleTemplateFile(c, model.ProjectIdForGlobalRuleTemplate, templateName)
}

// ExportProjectRuleTemplateFile
// @Summary 导出项目规则模板
// @Description export rule template in a project
// @Id exportProjectRuleTemplateV1
// @Tags rule_template
// @Param project_name path string true "project name"
// @Param rule_template_name path string true "rule template name"
// @Security ApiKeyAuth
// @Success 200 file 1 "sqle rule template file"
// @router /v1/projects/{project_name}/rule_templates/{rule_template_name}/export [get]
func ExportProjectRuleTemplateFile(c echo.Context) error {
	// 权限检查
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	// 获取模板内容
	templateName := c.Param("rule_template_name")
	return exportRuleTemplateFile(c, projectUid, templateName)
}

func exportRuleTemplateFile(c echo.Context, projectID string, ruleTemplateName string) error {
	// 获取规则模板详情
	s := model.GetStorage()
	template, exist, err := s.GetRuleTemplateDetailByNameAndProjectIds([]string{projectID}, ruleTemplateName, "")
	if err != nil {
		return c.JSON(200, controller.NewBaseReq(err))
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrRuleTemplateNotExist)
	}

	// 补充缺失的信息(规则说明等描述信息)
	ruleNames := []string{}
	for _, rule := range template.RuleList {
		ruleNames = append(ruleNames, rule.RuleName)
	}

	rules, err := model.GetStorage().GetRulesByNamesAndDBType(ruleNames, template.DBType)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	ruleCache := map[string] /*rule name*/ model.Rule{}
	for _, rule := range rules {
		ruleCache[rule.Name] = rule
	}

	resp := ParseProjectRuleTemplateFileResDataV1{
		Name:     template.Name,
		Desc:     template.Desc,
		DBType:   template.DBType,
		RuleList: []RuleResV1{},
	}
	for _, rule := range template.RuleList {
		r := RuleResV1{
			Name:       rule.RuleName,
			Desc:       ruleCache[rule.RuleName].Desc,
			Annotation: ruleCache[rule.RuleName].Annotation,
			Level:      rule.RuleLevel,
			Typ:        ruleCache[rule.RuleName].Typ,
			DBType:     rule.RuleDBType,
			Params:     []RuleParamResV1{},
		}

		for _, param := range rule.RuleParams {
			r.Params = append(r.Params, RuleParamResV1{
				Key:   param.Key,
				Value: param.Value,
				Desc:  param.Desc,
				Type:  string(param.Type),
			})
		}

		resp.RuleList = append(resp.RuleList, r)
	}

	bt, err := json.Marshal(resp)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	buff := &bytes.Buffer{}
	buff.Write(bt)
	c.Response().Header().Set(echo.HeaderContentDisposition,
		mime.FormatMediaType("attachment", map[string]string{"filename": fmt.Sprintf("RuleTemplate-%v.json", ruleTemplateName)}))
	return c.Blob(http.StatusOK, "text/plain;charset=utf-8", buff.Bytes())
}

type CustomRuleResV1 struct {
	RuleId     string `json:"rule_id"`
	Desc       string `json:"desc" example:"this is test rule"`
	Annotation string `json:"annotation" example:"this is test rule"`
	DBType     string `json:"db_type" example:"MySQL"`
	Level      string `json:"level" example:"notice" enums:"normal,notice,warn,error"`
	Type       string `json:"type" example:"DDL规则"`
	RuleScript string `json:"rule_script,omitempty"`
}

type GetCustomRulesResV1 struct {
	controller.BaseRes
	Data []CustomRuleResV1 `json:"data"`
}

type GetCustomRulesReqV1 struct {
	FilterDBType string `json:"filter_db_type" query:"filter_db_type"`
	FilterDesc   string `json:"filter_desc" query:"filter_desc"`
}

// @Summary 自定义规则列表
// @Description get all custom rule template
// @Id getCustomRulesV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Param filter_db_type query string false "filter db type"
// @Param filter_desc query string false "filter desc"
// @Success 200 {object} v1.GetCustomRulesResV1
// @router /v1/custom_rules [get]
func GetCustomRules(c echo.Context) error {
	return getCustomRules(c)
}

// @Summary 删除自定义规则
// @Description delete custom rule
// @Id deleteCustomRuleV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Param rule_id path string true "rule id"
// @Success 200 {object} controller.BaseRes
// @router /v1/custom_rules/{rule_id} [delete]
func DeleteCustomRule(c echo.Context) error {
	return deleteCustomRule(c)
}

type CreateCustomRuleReqV1 struct {
	DBType     string `json:"db_type" form:"db_type" example:"MySQL" valid:"required"`
	Desc       string `json:"desc" form:"desc" example:"this is test rule" valid:"required"`
	Annotation string `json:"annotation" form:"annotation" example:"this is test rule"`
	Level      string `json:"level" form:"level" example:"notice" valid:"required" enums:"normal,notice,warn,error"`
	Type       string `json:"type" form:"type" example:"DDL规则" valid:"required"`
	RuleScript string `json:"rule_script" form:"rule_script" valid:"required"`
}

// @Summary 添加自定义规则
// @Description create custom rule
// @Id createCustomRuleV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Param instance body v1.CreateCustomRuleReqV1 true "add custom rule"
// @Success 200 {object} controller.BaseRes
// @router /v1/custom_rules [post]
func CreateCustomRule(c echo.Context) error {
	return createCustomRule(c)
}

type UpdateCustomRuleReqV1 struct {
	Desc       *string `json:"desc" form:"desc" example:"this is test rule"`
	Annotation *string `json:"annotation" form:"annotation" example:"this is test rule"`
	Level      *string `json:"level" form:"level" example:"notice" enums:"normal,notice,warn,error"`
	Type       *string `json:"type" form:"type" example:"DDL规则"`
	RuleScript *string `json:"rule_script" form:"rule_script"`
}

// @Summary 更新自定义规则
// @Description update custom rule
// @Id updateCustomRuleV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Param rule_id path string true "rule id"
// @Param instance body v1.UpdateCustomRuleReqV1 true "update custom rule"
// @Success 200 {object} controller.BaseRes
// @router /v1/custom_rules/{rule_id} [patch]
func UpdateCustomRule(c echo.Context) error {
	return updateCustomRule(c)
}

type GetCustomRuleResV1 struct {
	controller.BaseRes
	Data CustomRuleResV1 `json:"data"`
}

// @Summary 获取自定义规则
// @Description get custom rule by rule_id
// @Id getCustomRuleV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Param rule_id path string true "rule id"
// @Success 200 {object} v1.GetCustomRuleResV1
// @router /v1/custom_rules/{rule_id} [get]
func GetCustomRule(c echo.Context) error {
	return getCustomRule(c)
}

type RuleTypeV1 struct {
	RuleType         string `json:"rule_type"`
	RuleCount        uint   `json:"rule_count"`
	IsCustomRuleType bool   `json:"is_custom_rule_type"`
}

type GetRuleTypeByDBTypeResV1 struct {
	controller.BaseRes
	Data []RuleTypeV1 `json:"data"`
}

// @Summary 获取规则分类
// @Description get rule type by db type
// @Id getRuleTypeByDBTypeV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Param db_type query string true "db type"
// @Success 200 {object} v1.GetRuleTypeByDBTypeResV1
// @router /v1/custom_rules/{db_type}/rule_types [get]
func GetRuleTypeByDBType(c echo.Context) error {
	return getRuleTypeByDBType(c)
}

type RuleInfo struct {
	Desc       string `json:"desc" example:"this is test rule"`
	Annotation string `json:"annotation" example:"this is test rule"`
}

type RuleKnowledgeResV1 struct {
	Rule             RuleInfo `json:"rule"`
	KnowledgeContent string   `json:"knowledge_content"`
}

type GetRuleKnowledgeResV1 struct {
	controller.BaseRes
	Data RuleKnowledgeResV1 `json:"data"`
}

// GetRuleKnowledge
// @Summary 查看规则知识库
// @Description get rule knowledge
// @Id getRuleKnowledgeV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Param rule_name path string true "rule name"
// @Param db_type path string true "db type of rule"
// @Success 200 {object} v1.GetRuleKnowledgeResV1
// @router /v1/rule_knowledge/db_types/{db_type}/rules/{rule_name}/ [get]
func GetRuleKnowledge(c echo.Context) error {
	return getRuleKnowledge(c)
}

type UpdateRuleKnowledgeReq struct {
	KnowledgeContent *string `json:"knowledge_content" form:"knowledge_content"`
}

// UpdateRuleKnowledgeV1
// @Summary 更新规则知识库
// @Description update rule knowledge
// @Id updateRuleKnowledge
// @Tags rule_template
// @Security ApiKeyAuth
// @Param rule_name path string true "rule name"
// @Param db_type path string true "db type of rule"
// @Param req body v1.UpdateRuleKnowledgeReq true "update rule knowledge"
// @Success 200 {object} controller.BaseRes
// @router /v1/rule_knowledge/db_types/{db_type}/rules/{rule_name}/ [patch]
func UpdateRuleKnowledgeV1(c echo.Context) error {
	return updateRuleKnowledge(c)
}

// GetCustomRuleKnowledge
// @Summary 查看自定义规则知识库
// @Description get custom rule knowledge
// @Id getCustomRuleKnowledgeV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Param rule_name path string true "rule name"
// @Param db_type path string true "db type of rule"
// @Success 200 {object} v1.GetRuleKnowledgeResV1
// @router /v1/rule_knowledge/db_types/{db_type}/custom_rules/{rule_name}/ [get]
func GetCustomRuleKnowledge(c echo.Context) error {
	return getCustomRuleKnowledge(c)
}

// UpdateCustomRuleKnowledgeV1
// @Summary 更新自定义规则知识库
// @Description update custom rule knowledge
// @Id updateCustomRuleKnowledge
// @Tags rule_template
// @Security ApiKeyAuth
// @Param rule_name path string true "rule name"
// @Param db_type path string true "db type of rule"
// @Param req body v1.UpdateRuleKnowledgeReq true "update rule knowledge"
// @Success 200 {object} controller.BaseRes
// @router /v1/rule_knowledge/db_types/{db_type}/custom_rules/{rule_name}/ [patch]
func UpdateCustomRuleKnowledgeV1(c echo.Context) error {
	return updateCustomRuleKnowledge(c)
}
