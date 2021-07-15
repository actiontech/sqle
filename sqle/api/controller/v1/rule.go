package v1

import (
	"actiontech.cloud/sqle/sqle/sqle/api/controller"
	"actiontech.cloud/sqle/sqle/sqle/errors"
	"actiontech.cloud/sqle/sqle/sqle/model"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CreateRuleTemplateReqV1 struct {
	Name      string   `json:"rule_template_name" valid:"required,name"`
	Desc      string   `json:"desc"`
	Rules     []string `json:"rule_name_list" example:"ddl_check_index_count"`
	Instances []string `json:"instance_name_list"`
}

// @Summary 添加规则模板
// @Description create a rule template
// @Id createRuleTemplateV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Accept json
// @Param instance body v1.CreateRuleTemplateReqV1 true "add rule template request"
// @Success 200 {object} controller.BaseRes
// @router /v1/rule_templates [post]
func CreateRuleTemplate(c echo.Context) error {
	req := new(CreateRuleTemplateReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()
	_, exist, err := s.GetRuleTemplateByName(req.Name)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("rule template is exist")))
	}
	var rules []model.Rule
	if req.Rules != nil || len(req.Rules) > 0 {
		rules, err = s.GetAndCheckRuleExist(req.Rules)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	var instances []*model.Instance
	if req.Instances != nil || len(req.Instances) > 0 {
		instances, err = s.GetAndCheckInstanceExist(req.Instances)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	ruleTemplate := &model.RuleTemplate{
		Name: req.Name,
		Desc: req.Desc,
	}
	err = s.Save(ruleTemplate)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = s.UpdateRuleTemplateRules(ruleTemplate, rules...)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = s.UpdateRuleTemplateInstances(ruleTemplate, instances...)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type UpdateRuleTemplateReqV1 struct {
	Desc      *string  `json:"desc"`
	Rules     []string `json:"rule_name_list" example:"ddl_check_index_count"`
	Instances []string `json:"instance_name_list" example:"mysql-xxx"`
}

// @Summary 更新规则模板
// @Description update rule template
// @Id updateRuleTemplateV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Param rule_template_name path string true "rule template name"
// @Param instance body v1.UpdateRuleTemplateReqV1 true "update rule template request"
// @Success 200 {object} controller.BaseRes
// @router /v1/rule_templates/{rule_template_name}/ [patch]
func UpdateRuleTemplate(c echo.Context) error {
	templateName := c.Param("rule_template_name")
	req := new(UpdateRuleTemplateReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()
	template, exist, err := s.GetRuleTemplateByName(templateName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("rule template is not exist")))
	}
	var rules []model.Rule
	if req.Rules != nil || len(req.Rules) > 0 {
		rules, err = s.GetAndCheckRuleExist(req.Rules)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	var instances []*model.Instance
	if req.Instances != nil || len(req.Instances) > 0 {
		instances, err = s.GetAndCheckInstanceExist(req.Instances)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	if req.Desc != nil {
		template.Desc = *req.Desc
		err = s.Save(&template)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}
	if req.Rules != nil {
		err = s.UpdateRuleTemplateRules(template, rules...)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	if req.Instances != nil {
		err = s.UpdateRuleTemplateInstances(template, instances...)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type GetRuleTemplateResV1 struct {
	controller.BaseRes
	Data *RuleTemplateDetailResV1 `json:"data"`
}

type RuleTemplateDetailResV1 struct {
	Name      string   `json:"rule_template_name"`
	Desc      string   `json:"desc"`
	Rules     []string `json:"rule_name_list,omitempty"`
	Instances []string `json:"instance_name_list,omitempty"`
}

func convertRuleTemplateToRes(template *model.RuleTemplate) *RuleTemplateDetailResV1 {
	ruleNames := make([]string, 0, len(template.Rules))
	for _, rule := range template.Rules {
		ruleNames = append(ruleNames, rule.Name)
	}
	instanceNames := make([]string, 0, len(template.Instances))
	for _, instance := range template.Instances {
		instanceNames = append(instanceNames, instance.Name)
	}
	return &RuleTemplateDetailResV1{
		Name:      template.Name,
		Desc:      template.Desc,
		Rules:     ruleNames,
		Instances: instanceNames,
	}
}

// @Summary 获取规则模板信息
// @Description get rule template
// @Id getRuleTemplateV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Param rule_template_name path string true "rule template name"
// @Success 200 {object} v1.GetRuleTemplateResV1
// @router /v1/rule_templates/{rule_template_name}/ [get]
func GetRuleTemplate(c echo.Context) error {
	s := model.GetStorage()
	templateName := c.Param("rule_template_name")
	template, exist, err := s.GetRuleTemplateDetailByName(templateName)
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

// @Summary 删除规则模板
// @Description delete rule template
// @Id deleteRuleTemplateV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Param rule_template_name path string true "rule template name"
// @Success 200 {object} controller.BaseRes
// @router /v1/rule_templates/{rule_template_name}/ [delete]
func DeleteRuleTemplate(c echo.Context) error {
	s := model.GetStorage()
	templateName := c.Param("rule_template_name")
	template, exist, err := s.GetRuleTemplateByName(templateName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("rule template is not exist")))
	}
	err = s.Delete(template)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type GetRuleTemplatesReqV1 struct {
	FilterInstanceName string `json:"filter_instance_name" query:"filter_instance_name"`
	PageIndex          uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize           uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetRuleTemplatesResV1 struct {
	controller.BaseRes
	Data      []RuleTemplateDetailResV1 `json:"data"`
	TotalNums uint64                    `json:"total_nums"`
}

// @Summary 规则模板列表
// @Description get all rule template
// @Id getRuleTemplateListV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Param filter_instance_name query string false "filter instance name"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v1.GetRuleTemplatesResV1
// @router /v1/rule_templates [get]
func GetRuleTemplates(c echo.Context) error {
	req := new(GetRuleTemplatesReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	var offset uint32
	if req.PageIndex >= 1 {
		offset = req.PageSize * (req.PageIndex - 1)
	}
	data := map[string]interface{}{
		"filter_instance_name": req.FilterInstanceName,
		"limit":                req.PageSize,
		"offset":               offset,
	}
	s := model.GetStorage()
	ruleTemplates, count, err := s.GetRuleTemplatesByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	ruleTemplatesReq := make([]RuleTemplateDetailResV1, 0, len(ruleTemplates))
	for _, ruleTemplate := range ruleTemplates {
		ruleTemplateReq := RuleTemplateDetailResV1{
			Name:      ruleTemplate.Name,
			Desc:      ruleTemplate.Desc,
			Rules:     ruleTemplate.RuleNames,
			Instances: ruleTemplate.InstanceNames,
		}
		ruleTemplatesReq = append(ruleTemplatesReq, ruleTemplateReq)
	}
	return c.JSON(http.StatusOK, &GetRuleTemplatesResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      ruleTemplatesReq,
		TotalNums: count,
	})
}

type GetRulesResV1 struct {
	controller.BaseRes
	Data []RuleResV1 `json:"data"`
}

type RuleResV1 struct {
	Name   string `json:"rule_name"`
	Desc   string `json:"desc"`
	Value  string `json:"value"`
	Level  string `json:"level" example:"error" enums:"normal,notice,warn,error"`
	Typ    string `json:"type" example:"全局配置" `
	DBType string `json:"db_type" example:"mysql"`
}

func convertRulesToRes(rules []model.Rule) []RuleResV1 {
	rulesRes := make([]RuleResV1, 0, len(rules))
	for _, rule := range rules {
		rulesRes = append(rulesRes, RuleResV1{
			Name:  rule.Name,
			Desc:  rule.Desc,
			Value: rule.Value,
			Level: rule.Level,
			Typ:   rule.Typ,
		})
	}
	return rulesRes
}

// @Summary 规则列表
// @Description get all rule template
// @Id getRuleListV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetRulesResV1
// @router /v1/rules [get]
func GetRules(c echo.Context) error {
	s := model.GetStorage()
	rules, err := s.GetAllRule()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, &GetRulesResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertRulesToRes(rules),
	})
}

type RuleTemplateTipResV1 struct {
	Name string `json:"rule_template_name"`
}

type GetRuleTemplateTipsResV1 struct {
	controller.BaseRes
	Data []RuleTemplateTipResV1 `json:"data"`
}

// @Summary 获取规则模板提示
// @Description get rule template tips
// @Id getRuleTemplateTipsV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetRuleTemplateTipsResV1
// @router /v1/rule_template_tips [get]
func GetRuleTemplateTips(c echo.Context) error {
	s := model.GetStorage()
	ruleTemplates, err := s.GetRuleTemplateTips()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	ruleTemplateTipsRes := make([]RuleTemplateTipResV1, 0, len(ruleTemplates))

	for _, roleTemplate := range ruleTemplates {
		ruleTemplateTipRes := RuleTemplateTipResV1{
			Name: roleTemplate.Name,
		}
		ruleTemplateTipsRes = append(ruleTemplateTipsRes, ruleTemplateTipRes)
	}
	return c.JSON(http.StatusOK, &GetRuleTemplateTipsResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    ruleTemplateTipsRes,
	})
}
