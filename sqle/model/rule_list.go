package model

import (
	"fmt"
	"strings"
)

type RuleTemplateDetail struct {
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	DBType string `json:"db_type"`
	// InstanceIds   RowList `json:"instance_ids"`
	// InstanceNames RowList `json:"instance_names"`
}

var ruleTemplatesQueryTpl = `SELECT rule_templates.name, rule_templates.desc, rule_templates.db_type
{{- template "body" . }}

{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var ruleTemplatesCountTpl = `SELECT COUNT(DISTINCT id)

{{- template "body" . -}}
`

var ruleTemplatesQueryBodyTpl = `
{{ define "body" }}
FROM rule_templates
WHERE
deleted_at IS NULL
AND project_id = :project_id
{{- end }}
`

func (s *Storage) GetRuleTemplatesByReq(data map[string]interface{}) (
	result []*RuleTemplateDetail, count uint64, err error) {

	err = s.getListResult(ruleTemplatesQueryBodyTpl, ruleTemplatesQueryTpl, data, &result)
	if err != nil {
		return result, 0, err
	}
	count, err = s.getCountResult(ruleTemplatesQueryBodyTpl, ruleTemplatesCountTpl, data)
	return result, count, err
}

func (s *Storage) GetRulesByReq(data map[string]interface{}) (
	result []*Rule, err error) {
	db := s.db.Preload("Categories")
	if data["filter_global_rule_template_name"] != "" {
		db = db.Joins("LEFT JOIN rule_template_rule ON rules.name = rule_template_rule.rule_name AND rules.db_type = rule_template_rule.db_type").
			Joins("LEFT JOIN rule_templates ON rule_template_rule.rule_template_id = rule_templates.id").
			Where("rule_templates.project_id = 0").
			Where("rule_templates.name = ?", data["filter_global_rule_template_name"].(string))
	}
	if data["filter_db_type"] != "" {
		db = db.Where("rules.db_type = ?", data["filter_db_type"])
	}
	if data["filter_rule_names"] != "" {
		if namesStr, yes := data["filter_rule_names"].(string); yes {
			db = db.Where("rules.name in (?)", strings.Split(namesStr, ","))
		}
	}
	if data["fuzzy_keyword_rule"] != "" {
		// todo i18n use json syntax to query?
		db = db.Where("rules.`i18n_rule_info` like ?", fmt.Sprintf("%%%s%%", data["fuzzy_keyword_rule"]))
	}
	if data["tags"] != "" {
		tags := strings.Split(data["tags"].(string), ",")
		db = db.Joins("LEFT JOIN audit_rule_category_rels on rules.name = audit_rule_category_rels.rule_name AND audit_rule_category_rels.rule_db_type = rules.db_type").
			Joins("LEFT JOIN audit_rule_categories on audit_rule_category_rels.category_id = audit_rule_categories.id").
			Where("audit_rule_categories.tag in (?)", tags).
			Group("rules.name, rules.db_type").
			Having("COUNT(DISTINCT audit_rule_categories.id) = ?", len(tags))
	}
	err = db.Find(&result).Error
	return result, err
}

func (s *Storage) GetCustomRulesByReq(data map[string]interface{}) (
	result []*CustomRule, err error) {
	db := s.db.Preload("Categories")
	if data["filter_global_rule_template_name"] != "" {
		db = db.Joins("LEFT JOIN rule_template_custom_rules ON custom_rules.rule_id = rule_template_custom_rules.rule_id").
			Joins("LEFT JOIN rule_templates ON rule_template_custom_rules.rule_template_id = rule_templates.id").
			Where("rule_templates.project_id = 0").
			Where("rule_templates.deleted_at is null").
			Where("rule_templates.name = ?", data["filter_global_rule_template_name"])
	}
	if data["filter_db_type"] != "" {
		db = db.Where("custom_rules.db_type = ?", data["filter_db_type"])
	}
	if data["filter_rule_names"] != "" {
		if namesStr, yes := data["filter_rule_names"].(string); yes {
			db = db.Where("custom_rules.rule_id in (?)", strings.Split(namesStr, ","))
		}
	}
	if data["fuzzy_keyword_rule"] != "" {
		db = db.Where("custom_rules.`desc` like ? OR custom_rules.annotation like ?", fmt.Sprintf("%%%s%%", data["fuzzy_keyword_rule"]), fmt.Sprintf("%%%s%%", data["fuzzy_keyword_rule"]))
	}
	if data["tags"] != "" {
		tags := strings.Split(data["tags"].(string), ",")
		db = db.Joins("LEFT JOIN custom_rule_category_rels on custom_rule_category_rels.custom_rule_id = custom_rules.rule_id").
			Joins("LEFT JOIN audit_rule_categories on custom_rule_category_rels.category_id = audit_rule_categories.id").
			Where("audit_rule_categories.tag in (?)", tags).
			Group("custom_rules.id").
			Having("COUNT(DISTINCT audit_rule_categories.id) = ?", len(tags))
	}
	err = db.Find(&result).Error
	return result, err
}
