package controller

import (
	"github.com/labstack/echo"
	"net/http"
	"sqle/storage"
)

type CreateTplReq struct {
	Name  string   `form:"name"`
	Desc  string   `form:"desc"`
	Rules []string `form:"rules" example:"check_name_length"`
}

// @Summary 添加规则模板
// @Description create a rule template
// @Accept json
// @Accept json
// @Param instance body controller.CreateTplReq true "add instance"
// @Success 200 {object} controller.BaseRes
//// @router /rule_templates [post]
func CreateTemplate(c echo.Context) error {
	s := storage.GetStorage()
	req := new(CreateTplReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	//fmt.Println(req)

	_, exist, err := s.GetTemplateByName(req.Name)
	if err != nil {
		return c.JSON(200, NewBaseReq(-1, err.Error()))
	}
	if exist {
		return c.JSON(200, NewBaseReq(-1, err.Error()))
	}
	t := &storage.RuleTemplate{
		Name: req.Name,
		Desc: req.Desc,
	}
	err = s.Save(t)
	if err != nil {
		return c.JSON(200, NewBaseReq(-1, err.Error()))
	}
	//rules := make([]storage.Rule, 0, len(req.Rules))
	//for _, rule := range req.Rules {
	//	//r := storage.Rule{
	//	//	TemplateId: t.ID,
	//	//	Code:       rule.Code,
	//	//	Value:      rule.Value,
	//	//	Level:      rule.Level,
	//	//}
	//	//rules = append(rules, r)
	//}
	//if len(rules) > 0 {
	//	err := s.UpdateRules(t, rules...)
	//	if err != nil {
	//		return c.JSON(200, NewBaseReq(-1, err.Error()))
	//	}
	//}
	return c.JSON(200, NewBaseReq(0, "ok"))
}

type GetAllTplRes struct {
	BaseRes
	Data []storage.RuleTemplate `json:"data"`
}

// @Summary 规则模板列表
// @Description get all rule template
// @Success 200 {object} controller.GetAllTplRes
// @router /rule_templates [get]
func GetAllTpl(c echo.Context) error {
	s := storage.GetStorage()
	if s == nil {
		c.String(500, "nil")
	}
	ts, err := s.GetAllTpl()
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
	Data []storage.Rule `json:"data"`
}

// @Summary 规则列表
// @Description get all rule template
// @Success 200 {object} controller.GetAllRuleRes
// @router /rules [get]
func GetRules(c echo.Context) error {
	return nil
}

// @Summary 获取规则模板
// @Description get rule template
// @Success 200 {object} controller.GetAllRuleRes
// @router /rule_template/template_id/ [get]
func GetRuleTemplate(c echo.Context) error {
	return nil
}

// @Summary 更新规则模板
// @Description update rule template
// @Param instance body controller.CreateTplReq true "update rule template"
// @Success 200 {object} controller.BaseRes
// @router /rule_template/template_id/ [put]
func UpdateRuleTemplate(c echo.Context) error {
	return nil
}

// @Summary 删除规则模板
// @Description delete rule template
// @Success 200 {object} controller.BaseRes
// @router /rule_template/template_id/ [delete]
func DeleteRuleTemplate(c echo.Context) error {
	return nil
}
