package controller

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"sqle/model"
)

type CreateTplReq struct {
	Name      string   `form:"name"`
	Desc      string   `form:"desc"`
	RulesName []string `json:"rule_name_list" form:"rule_name_list" example:"ddl_create_table_not_exist"`
}

// @Summary 添加规则模板
// @Description create a rule template
// @Accept json
// @Accept json
// @Param instance body controller.CreateTplReq true "add instance"
// @Success 200 {object} controller.BaseRes
// @router /rule_templates [post]
func CreateTemplate(c echo.Context) error {
	s := model.GetStorage()
	req := new(CreateTplReq)
	if err := c.Bind(req); err != nil {
		return err
	}

	fmt.Println(req)

	_, exist, err := s.GetTemplateByName(req.Name)
	if err != nil {
		return c.JSON(200, NewBaseReq(-1, err.Error()))
	}
	if exist {
		return c.JSON(200, NewBaseReq(-1, "template is exist"))
	}
	t := &model.RuleTemplate{
		Name: req.Name,
		Desc: req.Desc,
	}
	rules, err := s.GetAllRule()
	if err != nil {
		return c.JSON(200, NewBaseReq(-1, err.Error()))
	}
	ruleMap := model.GetRuleMapFromAllArray(rules)
	for _, name := range req.RulesName {
		fmt.Println("check: ", name)
		if rule, ok := ruleMap[name]; !ok {
			return c.JSON(200, NewBaseReq(-1, fmt.Sprintf("rule: %s is invalid", name)))
		} else {
			t.Rules = append(t.Rules, rule)
		}
	}

	err = s.Save(t)
	if err != nil {
		return c.JSON(200, NewBaseReq(-1, err.Error()))
	}

	return c.JSON(200, NewBaseReq(0, "ok"))
}

type GetAllTplRes struct {
	BaseRes
	Data []model.RuleTemplate `json:"data"`
}

// @Summary 规则模板列表
// @Description get all rule template
// @Success 200 {object} controller.GetAllTplRes
// @router /rule_templates [get]
func GetAllTpl(c echo.Context) error {
	s := model.GetStorage()
	if s == nil {
		c.String(500, "nil")
	}
	ts, err := s.GetAllTemplate()
	if err != nil {
		return c.String(500, err.Error())
	}
	return c.JSON(http.StatusOK, &GetAllTplRes{
		BaseRes: NewBaseReq(0, "ok"),
		Data:    ts,
	})
}

type GetAllRuleRes struct {
	BaseRes
	Data []model.Rule `json:"data"`
}

// @Summary 规则列表
// @Description get all rule template
// @Success 200 {object} controller.GetAllRuleRes
// @router /rules [get]
func GetRules(c echo.Context) error {
	s := model.GetStorage()
	rules, err := s.GetAllRule()
	if err != nil {
		return c.JSON(200, NewBaseReq(-1, err.Error()))
	}
	return c.JSON(200, &GetAllRuleRes{
		BaseRes: NewBaseReq(0, "ok"),
		Data:    rules,
	})
	return nil
}

type GetRuleTplRes struct {
	BaseRes
	Data model.RuleTemplate `json:"data"`
}

// @Summary 获取规则模板
// @Description get rule template
// @Param template_id path string true "Rule Template ID"
// @Success 200 {object} controller.GetAllRuleRes
// @router /rule_templates/{template_id}/ [get]
func GetRuleTemplate(c echo.Context) error {
	s := model.GetStorage()
	templateId := c.Param("template_id")
	template, exist, err := s.GetTemplateById(templateId)
	if err != nil {
		return c.JSON(200, NewBaseReq(-1, err.Error()))
	}
	if !exist {
		return c.JSON(200, NewBaseReq(-1, fmt.Sprintf("template id %v not exist", templateId)))
	}
	return c.JSON(http.StatusOK, &GetRuleTplRes{
		BaseRes: NewBaseReq(0, "ok"),
		Data:    template,
	})
	return c.JSON(200, template)
}

// @Summary 更新规则模板
// @Description update rule template
// @Param instance body controller.CreateTplReq true "update rule template"
// @Success 200 {object} controller.BaseRes
// @router /rule_templates/template_id/ [put]
func UpdateRuleTemplate(c echo.Context) error {
	return nil
}

// @Summary 删除规则模板
// @Description delete rule template
// @Success 200 {object} controller.BaseRes
// @router /rule_templates/template_id/ [delete]
func DeleteRuleTemplate(c echo.Context) error {
	return nil
}
