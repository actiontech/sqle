package controller
//
//import (
//	"fmt"
//	"net/http"
//	"strings"
//
//	"actiontech.cloud/universe/sqle/v4/sqle/errors"
//	"actiontech.cloud/universe/sqle/v4/sqle/model"
//	"github.com/labstack/echo/v4"
//)
//
//type CreateTplReq struct {
//	Name         *string  `json:"name" valid:"required"`
//	Desc         *string  `json:"desc" valid:"-"`
//	RulesName    []string `json:"rule_name_list" example:"ddl_check_index_count" valid:"-"`
//	InstanceName []string `json:"instance_name_list" example:"mysql-xxx" valid:"-"`
//}
//
//// @Summary 添加规则模板
//// @Description create a rule template
//// @Accept json
//// @Param instance body controller.CreateTplReq true "add instance"
//// @Success 200 {object} controller.BaseRes
//// @router /rule_templates [post]
//func CreateTemplate(c echo.Context) error {
//	s := model.GetStorage()
//	req := new(CreateTplReq)
//	if err := c.Bind(req); err != nil {
//		return c.JSON(http.StatusOK, NewBaseReq(err))
//	}
//	if err := c.Validate(req); err != nil {
//		return c.JSON(http.StatusOK, NewBaseReq(err))
//	}
//
//	_, exist, err := s.GetTemplateByName(*req.Name)
//	if err != nil {
//		return c.JSON(200, NewBaseReq(err))
//	}
//	if exist {
//		return c.JSON(200, NewBaseReq(errors.New(errors.RULE_TEMPLATE_EXIST,
//			fmt.Errorf("template is exist"))))
//	}
//	t := &model.RuleTemplate{}
//	if req.Name != nil {
//		t.Name = *req.Name
//	}
//	if req.Desc != nil {
//		t.Desc = *req.Desc
//	}
//
//	rules, err := s.GetAllRule()
//	if err != nil {
//		return c.JSON(200, NewBaseReq(err))
//	}
//	ruleMap := model.GetRuleMapFromAllArray(rules)
//	for _, name := range req.RulesName {
//		if rule, ok := ruleMap[name]; !ok {
//			return c.JSON(200, NewBaseReq(errors.New(errors.RULE_NOT_EXIST,
//				fmt.Errorf("rule: %s is invalid", name))))
//		} else {
//			t.Rules = append(t.Rules, rule)
//		}
//	}
//
//	notExistInsts := []string{}
//	for _, instName := range req.InstanceName {
//		i, exist, err := s.GetInstByName(instName)
//		if err != nil {
//			return c.JSON(200, NewBaseReq(err))
//		}
//		if !exist {
//			notExistInsts = append(notExistInsts, instName)
//		}
//		t.Instances = append(t.Instances, *i)
//	}
//
//	if len(notExistInsts) > 0 {
//		err := fmt.Errorf("instance name : %s not exist", strings.Join(notExistInsts, ", "))
//		return c.JSON(200, NewBaseReq(errors.New(errors.RULE_TEMPLATE_NOT_EXIST, err)))
//	}
//
//	err = s.Save(t)
//	if err != nil {
//		return c.JSON(200, NewBaseReq(err))
//	}
//
//	return c.JSON(200, NewBaseReq(nil))
//}
//
//type GetRuleTplRes struct {
//	BaseRes
//	Data model.RuleTemplateDetail `json:"data"`
//}
//
//// @Summary 获取规则模板
//// @Description get rule template
//// @Param template_id path string true "Rule Template ID"
//// @Success 200 {object} controller.GetRuleTplRes
//// @router /rule_templates/{template_id}/ [get]
//func GetRuleTemplate(c echo.Context) error {
//	s := model.GetStorage()
//	templateId := c.Param("template_id")
//	template, exist, err := s.GetTemplateById(templateId)
//	if err != nil {
//		return c.JSON(200, NewBaseReq(err))
//	}
//	if !exist {
//		return c.JSON(200, NewBaseReq(errors.New(errors.RULE_TEMPLATE_NOT_EXIST,
//			fmt.Errorf("rule template is not exist"))))
//	}
//
//	return c.JSON(http.StatusOK, &GetRuleTplRes{
//		BaseRes: NewBaseReq(nil),
//		Data:    template.Detail(),
//	})
//}
//
//// @Summary 删除规则模板
//// @Description delete rule template
//// @Param template_id path string true "Template ID"
//// @Success 200 {object} controller.BaseRes
//// @router /rule_templates/{template_id}/ [delete]
//func DeleteRuleTemplate(c echo.Context) error {
//	s := model.GetStorage()
//	templateId := c.Param("template_id")
//	template, exist, err := s.GetTemplateById(templateId)
//	if err != nil {
//		return c.JSON(200, NewBaseReq(err))
//	}
//	if !exist {
//		return c.JSON(200, NewBaseReq(errors.New(errors.RULE_TEMPLATE_NOT_EXIST,
//			fmt.Errorf("rule template is not exist"))))
//	}
//	instances, err := s.GetInstancesNameByTemplate(template)
//	if err != nil {
//		return c.JSON(200, NewBaseReq(err))
//	}
//	if len(instances) > 0 {
//		return c.JSON(200, NewBaseReq(errors.New(errors.RULE_TEMPLATE_IS_USED,
//			fmt.Errorf(fmt.Sprintf("rule template is used by instance %s",
//				strings.Join(instances, ","))))))
//	}
//
//	err = s.Delete(template)
//	if err != nil {
//		return c.JSON(200, NewBaseReq(err))
//	}
//	return c.JSON(200, NewBaseReq(nil))
//}
//
//// @Summary 更新规则模板
//// @Description update rule template
//// @Param template_id path string true "Template ID"
//// @Param instance body controller.CreateTplReq true "update rule template"
//// @Success 200 {object} controller.BaseRes
//// @router /rule_templates/{template_id}/ [patch]
//func UpdateRuleTemplate(c echo.Context) error {
//	s := model.GetStorage()
//	templateId := c.Param("template_id")
//	req := new(CreateTplReq)
//	if err := c.Bind(req); err != nil {
//		return c.JSON(200, NewBaseReq(err))
//	}
//
//	template, exist, err := s.GetTemplateById(templateId)
//	if err != nil {
//		return c.JSON(200, NewBaseReq(err))
//	}
//	if !exist {
//		return c.JSON(200, NewBaseReq(errors.New(errors.RULE_TEMPLATE_NOT_EXIST,
//			fmt.Errorf("rule template is not exist"))))
//	}
//
//	if req.Name != nil && template.Name != *req.Name {
//		_, exist, err := s.GetTemplateByName(*req.Name)
//		if err != nil {
//			return c.JSON(200, NewBaseReq(err))
//		}
//		if exist {
//			return c.JSON(200, NewBaseReq(errors.New(errors.RULE_TEMPLATE_EXIST,
//				fmt.Errorf("template name %s is exist", *req.Name))))
//		}
//	}
//
//	if req.Name != nil {
//		template.Name = *req.Name
//	}
//	if req.Desc != nil {
//		template.Desc = *req.Desc
//	}
//
//	templateRules := []model.Rule{}
//	if req.RulesName != nil {
//		template.Rules = nil
//		rules, err := s.GetAllRule()
//		if err != nil {
//			return c.JSON(200, NewBaseReq(err))
//		}
//		ruleMap := model.GetRuleMapFromAllArray(rules)
//		for _, name := range req.RulesName {
//			if rule, ok := ruleMap[name]; !ok {
//				return c.JSON(200, NewBaseReq(errors.New(errors.RULE_NOT_EXIST,
//					fmt.Errorf("rule: %s is invalid", name))))
//			} else {
//				templateRules = append(templateRules, rule)
//			}
//		}
//
//	}
//
//	instances := []model.Instance{}
//	if req.InstanceName != nil {
//		template.Instances = nil
//		notExistInsts := []string{}
//		for _, instName := range req.InstanceName {
//			t, exist, err := s.GetInstByName(instName)
//			if err != nil {
//				return c.JSON(200, NewBaseReq(err))
//			}
//			if !exist {
//				notExistInsts = append(notExistInsts, instName)
//			}
//			instances = append(instances, *t)
//		}
//
//		if len(notExistInsts) > 0 {
//			err := fmt.Errorf("instance name : %s not exist", strings.Join(notExistInsts, ", "))
//			return c.JSON(200, NewBaseReq(errors.New(errors.RULE_TEMPLATE_NOT_EXIST, err)))
//		}
//	}
//
//	err = s.Save(&template)
//	if err != nil {
//		return c.JSON(200, NewBaseReq(err))
//	}
//	if req.RulesName != nil {
//		err = s.UpdateTemplateRules(&template, templateRules...)
//		if err != nil {
//			return c.JSON(200, NewBaseReq(err))
//		}
//	}
//
//	if req.InstanceName != nil {
//		err = s.UpdateTemplateInst(&template, instances...)
//		if err != nil {
//			return c.JSON(200, NewBaseReq(err))
//		}
//	}
//
//	return c.JSON(200, NewBaseReq(nil))
//}
//
//type GetAllTplRes struct {
//	BaseRes
//	Data []model.RuleTemplate `json:"data"`
//}
//
//// @Summary 规则模板列表
//// @Description get all rule template
//// @Success 200 {object} controller.GetAllTplRes
//// @router /rule_templates [get]
//func GetAllTpl(c echo.Context) error {
//	s := model.GetStorage()
//	ts, err := s.GetAllTemplate()
//	if err != nil {
//		return c.JSON(200, NewBaseReq(err))
//	}
//	return c.JSON(http.StatusOK, &GetAllTplRes{
//		BaseRes: NewBaseReq(nil),
//		Data:    ts,
//	})
//}
//
//type GetAllRuleRes struct {
//	BaseRes
//	Data []model.Rule `json:"data"`
//}
//
//// @Summary 规则列表
//// @Description get all rule template
//// @Success 200 {object} controller.GetAllRuleRes
//// @router /rules [get]
//func GetRules(c echo.Context) error {
//	s := model.GetStorage()
//	rules, err := s.GetAllRule()
//	if err != nil {
//		return c.JSON(200, NewBaseReq(err))
//	}
//	return c.JSON(200, &GetAllRuleRes{
//		BaseRes: NewBaseReq(nil),
//		Data:    rules,
//	})
//}
//
//type UpdateRuleReq struct {
//	Name  string `json:"name"`
//	Value string `json:"value"`
//}
//
//type UpdateAllRuleReq struct {
//	Configs []UpdateRuleReq `json:"rule_list"`
//}
//
//// @Summary 修改配置
//// @Description update rules
//// @Accept json
//// @Produce json
//// @Param instance body controller.UpdateAllRuleReq true "update rule"
//// @Success 200 {object} controller.BaseRes
//// @router /rules [patch]
//func UpdateRules(c echo.Context) error {
//	s := model.GetStorage()
//	reqs := new(UpdateAllRuleReq)
//	if err := c.Bind(reqs); err != nil {
//		return c.JSON(http.StatusOK, NewBaseReq(err))
//	}
//	rules, err := s.GetAllRule()
//	if err != nil {
//		return c.JSON(http.StatusOK, NewBaseReq(err))
//	}
//	ruleMap := model.GetRuleMapFromAllArray(rules)
//	for _, req := range reqs.Configs {
//		if _, ok := ruleMap[req.Name]; ok {
//			err := s.UpdateRuleValueByName(req.Name, req.Value)
//			if err != nil {
//				return c.JSON(http.StatusOK, NewBaseReq(err))
//			}
//		}
//	}
//	return c.JSON(200, NewBaseReq(nil))
//}
