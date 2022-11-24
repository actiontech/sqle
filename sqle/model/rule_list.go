package model

import "github.com/actiontech/sqle/sqle/errors"

type RuleTemplateDetail struct {
	Name          string  `json:"name"`
	Desc          string  `json:"desc"`
	DBType        string  `json:"db_type"`
	InstanceIds   RowList `json:"instance_ids"`
	InstanceNames RowList `json:"instance_names"`
}

var ruleTemplatesQueryTpl = `SELECT rt1.name, rt1.desc, rt1.db_type,
GROUP_CONCAT(DISTINCT COALESCE(instances.id,'')) AS instance_ids,
GROUP_CONCAT(DISTINCT COALESCE(instances.name,'')) AS instance_names
FROM rule_templates AS rt1
LEFT JOIN instance_rule_template ON rt1.id = instance_rule_template.rule_template_id
LEFT JOIN instances ON instance_rule_template.instance_id = instances.id AND instances.deleted_at IS NULL
WHERE (SELECT COUNT(DISTINCT(rule_templates.id))

{{- template "body" . }}
AND rt1.id = rule_templates.id
) > 0
GROUP BY rt1.id
{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var ruleTemplatesCountTpl = `SELECT COUNT(DISTINCT rule_templates.id)

{{- template "body" . -}}
`

var ruleTemplatesQueryBodyTpl = `
{{ define "body" }}
FROM rule_templates
LEFT JOIN instance_rule_template ON rule_templates.id = instance_rule_template.rule_template_id
LEFT JOIN instances ON instance_rule_template.instance_id = instances.id AND instances.deleted_at IS NULL
WHERE
rule_templates.deleted_at IS NULL
AND rule_templates.project_id = :project_id

{{- if .filter_instance_name }}
AND instances.name = :filter_instance_name
{{- end }}
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

func (s *Storage) GetRuleTemplateTotalByProjectName(projectName string) (uint64, error) {
	var count uint64
	err := s.db.
		Table("rule_templates").
		Joins("LEFT JOIN projects ON rule_templates.project_id = projects.id").
		Where("projects.name = ?", projectName).
		Where("rule_templates.deleted_at IS NULL").
		Count(&count).
		Error
	return count, errors.ConnectStorageErrWrapper(err)
}
