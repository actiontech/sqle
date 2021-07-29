package model

import (
	"bytes"
	"text/template"
)

type RoleDetail struct {
	Id            int
	Name          string `json:"name"`
	Desc          string
	UserNames     RowList `json:"user_names"`
	InstanceNames RowList `json:"instance_names"`
}

var rolesQueryTpl = `SELECT roles.id, roles.name, roles.desc,
GROUP_CONCAT(DISTINCT COALESCE(users.login_name,'')) AS user_names,
GROUP_CONCAT(DISTINCT COALESCE(instances.name,'')) AS instance_names
FROM roles
LEFT JOIN user_role ON roles.id = user_role.role_id
LEFT JOIN users ON user_role.user_id = users.id AND users.deleted_at IS NULL
LEFT JOIN instance_role ON roles.id = instance_role.role_id
LEFT JOIN instances ON instance_role.instance_id = instances.id AND instances.deleted_at IS NULL
WHERE
roles.id in (SELECT DISTINCT(roles.id)

{{- template "body" . -}}
)
GROUP BY roles.id
{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var rolesCountTpl = `SELECT COUNT(DISTINCT roles.id)

{{- template "body" . -}}
`

var rolesQueryBodyTpl = `
{{ define "body" }}
FROM roles 
LEFT JOIN user_role ON roles.id = user_role.role_id
LEFT JOIN users ON user_role.user_id = users.id AND users.deleted_at IS NULL
LEFT JOIN instance_role ON roles.id = instance_role.role_id
LEFT JOIN instances ON instance_role.instance_id = instances.id AND instances.deleted_at IS NULL
WHERE
roles.deleted_at IS NULL

{{- if .filter_role_name }}
AND roles.name = :filter_role_name
{{- end }}

{{- if .filter_user_name }}
AND users.login_name = :filter_user_name
{{- end }}

{{- if .filter_instance_name }}
AND instances.name = :filter_instance_name
{{- end }}
{{- end }}
`

func (s *Storage) GetRolesByReq(data map[string]interface{}) (
	result []*RoleDetail, count uint64, err error) {

	err = s.getListResult(rolesQueryBodyTpl, rolesQueryTpl, data, &result)
	if err != nil {
		return result, 0, err
	}
	count, err = s.getCountResult(rolesQueryBodyTpl, rolesCountTpl, data)
	return result, count, err
}

func getSelectQuery(bodyTpl, queryTpl string, data interface{}) (string, error) {
	var buff bytes.Buffer
	tpl, err := template.New("getQuery").Parse(bodyTpl)
	if err != nil {
		return "", err
	}
	_, err = tpl.Parse(queryTpl)
	if err != nil {
		return "", err
	}
	err = tpl.Execute(&buff, data)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}

func getCountQuery(bodyTpl, countTpl string, data interface{}) (string, error) {
	var buff bytes.Buffer
	tpl, err := template.New("getCount").Parse(bodyTpl)
	if err != nil {
		return "", err
	}
	_, err = tpl.Parse(countTpl)
	if err != nil {
		return "", err
	}
	err = tpl.Execute(&buff, data)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}

func (s *Storage) getListResult(bodyTpl, queryTpl string, data map[string]interface{},
	result interface{}) error {
	selectQuery, err := getSelectQuery(bodyTpl, queryTpl, data)
	if err != nil {
		return err
	}

	sqlxDb := GetSqlxDb()
	nstmtTasksQuery, err := sqlxDb.PrepareNamed(selectQuery)
	if err != nil {
		return err
	}
	err = nstmtTasksQuery.Select(result, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) getCountResult(bodyTpl, countTpl string, data map[string]interface{}) (uint64, error) {
	sqlxDb := GetSqlxDb()
	countQuery, err := getCountQuery(bodyTpl, countTpl, data)
	if err != nil {
		return 0, err
	}
	nstmtCountQuery, err := sqlxDb.PrepareNamed(countQuery)
	if err != nil {
		return 0, err
	}
	var count uint64
	err = nstmtCountQuery.Get(&count, data)
	if err != nil {
		return 0, err
	}
	return count, nil
}
