package v1

import (
	"fmt"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
)

type CreateRuleTemplateReqV1 struct {
	Name     string      `json:"rule_template_name" valid:"required,name"`
	Desc     string      `json:"desc"`
	DBType   string      `json:"db_type" valid:"required"`
	RuleList []RuleReqV1 `json:"rule_list" form:"rule_list" valid:"required,dive,required"`
}

type RuleReqV1 struct {
	Name   string           `json:"name" form:"name" valid:"required" example:"ddl_check_index_count"`
	Level  string           `json:"level" form:"level" valid:"required" example:"error"`
	Params []RuleParamReqV1 `json:"params" form:"params" valid:"dive,required"`
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

const projectIdForGlobalRuleTemplate = 0

// @Summary 添加全局规则模板
// @Description create a global rule template
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
	_, exist, err := s.GetRuleTemplateByProjectIdAndName(projectIdForGlobalRuleTemplate, req.Name)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("rule template is exist")))
	}

	ruleTemplate := &model.RuleTemplate{
		Name:   req.Name,
		Desc:   req.Desc,
		DBType: req.DBType,
	}
	templateRules := make([]model.RuleTemplateRule, 0, len(req.RuleList))
	if req.RuleList != nil || len(req.RuleList) > 0 {
		templateRules, err = checkAndGenerateRules(req.RuleList, ruleTemplate)
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
	template, exist, err := s.GetRuleTemplateByProjectIdAndName(projectIdForGlobalRuleTemplate, templateName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("rule template is not exist")))
	}

	templateRules := make([]model.RuleTemplateRule, 0, len(req.RuleList))
	if req.RuleList != nil || len(req.RuleList) > 0 {
		templateRules, err = checkAndGenerateRules(req.RuleList, template)
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
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type GetRuleTemplateResV1 struct {
	controller.BaseRes
	Data *RuleTemplateDetailResV1 `json:"data"`
}

type RuleTemplateDetailResV1 struct {
	Name      string                        `json:"rule_template_name"`
	Desc      string                        `json:"desc"`
	DBType    string                        `json:"db_type"`
	Instances []*GlobalRuleTemplateInstance `json:"instance_list,omitempty"`
	RuleList  []RuleResV1                   `json:"rule_list,omitempty"`
}

func convertRuleTemplateToRes(template *model.RuleTemplate, instNameToProjectName map[uint]model.ProjectAndInstance) *RuleTemplateDetailResV1 {
	instanceNames := make([]*GlobalRuleTemplateInstance, 0, len(template.Instances))
	for _, instance := range template.Instances {
		instanceNames = append(instanceNames, &GlobalRuleTemplateInstance{
			ProjectName:  instNameToProjectName[instance.ID].ProjectName,
			InstanceName: instance.Name,
		})
	}
	ruleList := make([]RuleResV1, 0, len(template.RuleList))
	for _, r := range template.RuleList {
		ruleList = append(ruleList, convertRuleToRes(r.GetRule()))
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
// @Success 200 {object} v1.GetRuleTemplateResV1
// @router /v1/rule_templates/{rule_template_name}/ [get]
func GetRuleTemplate(c echo.Context) error {
	s := model.GetStorage()
	templateName := c.Param("rule_template_name")
	template, exist, err := s.GetRuleTemplateDetailByNameAndProjectId(projectIdForGlobalRuleTemplate, templateName)
	if err != nil {
		return c.JSON(200, controller.NewBaseReq(err))
	}
	if !exist {
		return c.JSON(200, controller.NewBaseReq(errors.New(errors.DataNotExist,
			fmt.Errorf("rule template is not exist"))))
	}

	instanceIds := make([]uint, len(template.Instances))
	for i, inst := range template.Instances {
		instanceIds[i] = inst.ID
	}
	instanceNameToProjectName, err := s.GetProjectNamesByInstanceIds(instanceIds)
	if err != nil {
		return c.JSON(200, controller.NewBaseReq(err))
	}

	return c.JSON(http.StatusOK, &GetRuleTemplateResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertRuleTemplateToRes(template, instanceNameToProjectName),
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
	template, exist, err := s.GetRuleTemplateByProjectIdAndName(projectIdForGlobalRuleTemplate, templateName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("rule template is not exist")))
	}
	used, err := s.IsRuleTemplateBeingUsed(templateName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if used {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("rule template is being used")))
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
	Data      []RuleTemplateResV1 `json:"data"`
	TotalNums uint64              `json:"total_nums"`
}

type RuleTemplateResV1 struct {
	Name      string                        `json:"rule_template_name"`
	Desc      string                        `json:"desc"`
	DBType    string                        `json:"db_type"`
	Instances []*GlobalRuleTemplateInstance `json:"instance_list"`
}

type GlobalRuleTemplateInstance struct {
	ProjectName  string `json:"project_name"`
	InstanceName string `json:"instance_name"`
}

// @Summary 全局规则模板列表
// @Description get all global rule template
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

	ruleTemplatesReq := make([]RuleTemplateResV1, 0, len(ruleTemplates))
	for _, ruleTemplate := range ruleTemplates {
		ruleTemplateReq := RuleTemplateResV1{
			Name:   ruleTemplate.Name,
			Desc:   ruleTemplate.Desc,
			DBType: ruleTemplate.DBType,
			//Instances: ruleTemplate.InstanceNames,
		}
		ruleTemplatesReq = append(ruleTemplatesReq, ruleTemplateReq)
	}
	return c.JSON(http.StatusOK, &GetRuleTemplatesResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      ruleTemplatesReq,
		TotalNums: count,
	})
}

type GetRulesReqV1 struct {
	FilterDBType string `json:"filter_db_type" query:"filter_db_type"`
}

type GetRulesResV1 struct {
	controller.BaseRes
	Data []RuleResV1 `json:"data"`
}

type RuleResV1 struct {
	Name   string           `json:"rule_name"`
	Desc   string           `json:"desc"`
	Level  string           `json:"level" example:"error" enums:"normal,notice,warn,error"`
	Typ    string           `json:"type" example:"全局配置" `
	DBType string           `json:"db_type" example:"mysql"`
	Params []RuleParamResV1 `json:"params,omitempty"`
}

type RuleParamResV1 struct {
	Key   string `json:"key" form:"key"`
	Value string `json:"value" form:"value"`
	Desc  string `json:"desc" form:"desc"`
	Type  string `json:"type" form:"type" enums:"string,int,bool"`
}

func convertRuleToRes(rule *model.Rule) RuleResV1 {
	ruleRes := RuleResV1{
		Name:   rule.Name,
		Desc:   rule.Desc,
		Level:  rule.Level,
		Typ:    rule.Typ,
		DBType: rule.DBType,
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

func convertRulesToRes(rules []*model.Rule) []RuleResV1 {
	rulesRes := make([]RuleResV1, 0, len(rules))
	for _, rule := range rules {
		rulesRes = append(rulesRes, convertRuleToRes(rule))
	}
	return rulesRes
}

// @Summary 规则列表
// @Description get all rule template
// @Id getRuleListV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Param filter_db_type query string false "filter db type"
// @Success 200 {object} v1.GetRulesResV1
// @router /v1/rules [get]
func GetRules(c echo.Context) error {
	req := new(GetRulesReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()
	var rules []*model.Rule
	var err error
	if req.FilterDBType != "" {
		rules, err = s.GetAllRuleByDBType(req.FilterDBType)
	} else {
		rules, err = s.GetAllRule()
	}
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, &GetRulesResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertRulesToRes(rules),
	})
}

type RuleTemplateTipReqV1 struct {
	FilterDBType string `json:"filter_db_type" query:"filter_db_type"`
}

type RuleTemplateTipResV1 struct {
	Name   string `json:"rule_template_name"`
	DBType string `json:"db_type"`
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
// @Param filter_db_type query string false "filter db type"
// @Success 200 {object} v1.GetRuleTemplateTipsResV1
// @router /v1/rule_template_tips [get]
func GetRuleTemplateTips(c echo.Context) error {
	req := new(RuleTemplateTipReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	s := model.GetStorage()
	ruleTemplates, err := s.GetRuleTemplateTips(req.FilterDBType)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	ruleTemplateTipsRes := make([]RuleTemplateTipResV1, 0, len(ruleTemplates))
	for _, roleTemplate := range ruleTemplates {
		ruleTemplateTipRes := RuleTemplateTipResV1{
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
// @Param instance body v1.CloneRuleTemplateReqV1 true "clone rule template request"
// @Success 200 {object} controller.BaseRes
// @router /v1/rule_templates/{rule_template_name}/clone [post]
func CloneRuleTemplate(c echo.Context) error {
	req := new(CloneRuleTemplateReqV1)
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

	sourceTplName := c.Param("rule_template_name")
	sourceTpl, exist, err := s.GetRuleTemplateDetailByName(sourceTplName)
	if err != nil {
		return c.JSON(200, controller.NewBaseReq(err))
	}
	if !exist {
		return c.JSON(200, controller.NewBaseReq(errors.New(errors.DataNotExist,
			fmt.Errorf("source rule template %s is not exist", sourceTplName))))
	}

	var instances []*model.Instance
	if req.Instances != nil || len(req.Instances) > 0 {
		instances, err = s.GetAndCheckInstanceExist(req.Instances)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	err = CheckRuleTemplateCanBeBindEachInstance(s, req.Name, instances)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = CheckInstanceAndRuleTemplateDbType([]*model.RuleTemplate{sourceTpl}, instances...)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	ruleTemplate := &model.RuleTemplate{
		Name:   req.Name,
		Desc:   req.Desc,
		DBType: sourceTpl.DBType,
	}
	err = s.Save(ruleTemplate)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = s.CloneRuleTemplateRules(sourceTpl, ruleTemplate)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = s.UpdateRuleTemplateInstances(ruleTemplate, instances...)
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
	Name      string      `json:"rule_template_name" valid:"required,name"`
	Desc      string      `json:"desc"`
	DBType    string      `json:"db_type" valid:"required"`
	Instances []string    `json:"instance_name_list"`
	RuleList  []RuleReqV1 `json:"rule_list" form:"rule_list" valid:"required,dive,required"`
}

// CreateProjectRuleTemplate
// @Summary 添加项目规则模板
// @Description create a rule template in project
// @Id createProjectRuleTemplateV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Accept json
// @Param project_name path string true "project name"
// @Param instance body v1.CreateProjectRuleTemplateReqV1 true "add rule template request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/rule_templates [post]
func CreateProjectRuleTemplate(c echo.Context) error {
	return nil
}

type UpdateProjectRuleTemplateReqV1 struct {
	Desc      *string     `json:"desc"`
	Instances []string    `json:"instance_name_list" example:"mysql-xxx"`
	RuleList  []RuleReqV1 `json:"rule_list" form:"rule_list" valid:"dive,required"`
}

// UpdateProjectRuleTemplate
// @Summary 更新项目规则模板
// @Description update rule template in project
// @Id updateProjectRuleTemplateV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param rule_template_name path string true "rule template name"
// @Param instance body v1.UpdateProjectRuleTemplateReqV1 true "update rule template request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/rule_templates/{rule_template_name}/ [patch]
func UpdateProjectRuleTemplate(c echo.Context) error {
	return nil
}

type GetProjectRuleTemplateResV1 struct {
	controller.BaseRes
	Data *RuleProjectTemplateDetailResV1 `json:"data"`
}

type RuleProjectTemplateDetailResV1 struct {
	Name      string                         `json:"rule_template_name"`
	Desc      string                         `json:"desc"`
	DBType    string                         `json:"db_type"`
	Instances []*ProjectRuleTemplateInstance `json:"instance_list,omitempty"`
	RuleList  []RuleResV1                    `json:"rule_list,omitempty"`
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
// @Success 200 {object} v1.GetProjectRuleTemplateResV1
// @router /v1/projects/{project_name}/rule_templates/{rule_template_name}/ [get]
func GetProjectRuleTemplate(c echo.Context) error {
	return nil
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
	return nil
}

type GetProjectRuleTemplatesResV1 struct {
	controller.BaseRes
	Data      []ProjectRuleTemplateResV1 `json:"data"`
	TotalNums uint64                     `json:"total_nums"`
}

type ProjectRuleTemplateResV1 struct {
	Name      string                         `json:"rule_template_name"`
	Desc      string                         `json:"desc"`
	DBType    string                         `json:"db_type"`
	Instances []*ProjectRuleTemplateInstance `json:"instance_list"`
}

// GetProjectRuleTemplates
// @Summary 项目规则模板列表
// @Description get all rule template in a project
// @Id getProjectRuleTemplateListV1
// @Tags rule_template
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v1.GetProjectRuleTemplatesResV1
// @router /v1/projects/{project_name}/rule_templates [get]
func GetProjectRuleTemplates(c echo.Context) error {
	return nil
}

type CloneProjectRuleTemplateReqV1 struct {
	Name      string   `json:"new_rule_template_name" valid:"required"`
	Desc      string   `json:"desc"`
	Instances []string `json:"instance_name_list"`
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
// @Param instance body v1.CloneProjectRuleTemplateReqV1 true "clone rule template request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/rule_templates/{rule_template_name}/clone [post]
func CloneProjectRuleTemplate(c echo.Context) error {
	return nil
}
