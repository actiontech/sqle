package model

import "time"

type ProjectDetail struct {
	Name           string    `json:"name"`
	Desc           string    `json:"desc"`
	CreateUserName string    `json:"create_user_name"`
	CreateTime     time.Time `json:"create_time"`
}

var projectsQueryTpl = `SELECT
projects.name , projects.` + "`desc`" + `, users.login_name as create_user_name, projects.create_at as create_time
FROM projects
JOIN project_user on project_user.project_id = projects.id
JOIN users on users.id = project_user.user_id
WHERE users.deleted_at IS NULL
AND users.stat = 0

{{ if .filter_user_name }}
AND users.login_name = :filter_user_name
{{ end }}

UNION
SELECT
projects.name , projects.` + "`desc`" + `, users.login_name as create_user_name, projects.create_at as create_time
FROM projects
JOIN project_user_group on project_user_group.project_id = projects.id
JOIN user_group_users on project_user_group.user_group_id = user_group_users.user_group_id
JOIN users on users.id = user_group_user.user_id
WHERE users.deleted_at IS NULL
AND users.stat = 0

{{ if .filter_user_name }}
AND users.login_name = :filter_user_name
{{ end }}

{{ if .limit }}
LIMIT :limit OFFSET :offset
{{ end }}
{{- template "body" . -}}
`

var projectsCountTpl = `SELECT COUNT(1) FROM (SELECT DISTINCT projects.id
FROM projects
JOIN project_user on project_user.project_id = projects.id
JOIN users on users.id = project_user.user_id
WHERE users.deleted_at IS NULL
AND users.stat = 0

{{ if .filter_user_name }}
AND users.login_name = :filter_user_name
{{ end }}

UNION
SELECT
DISTINCT projects.id
FROM projects
JOIN project_user_group on project_user_group.project_id = projects.id
JOIN user_group_users on project_user_group.user_group_id = user_group_users.user_group_id
JOIN users on users.id = user_group_users.user_id
WHERE users.deleted_at IS NULL
AND users.stat = 0

{{ if .filter_user_name }}
AND users.login_name = :filter_user_name
{{ end }}
) as a
{{- template "body" . -}}
`

var projectsQueryBodyTpl = `
{{ define "body" . }}
{{ end }}
`

func (s *Storage) GetProjectsByReq(data map[string]interface{}) (
	result []*ProjectDetail, count uint64, err error) {

	if data["filter_user_name"] == DefaultAdminUser {
		data["filter_user_name"] = nil
	}

	err = s.getListResult(projectsQueryBodyTpl, projectsQueryTpl, data, &result)
	if err != nil {
		return result, 0, err
	}
	count, err = s.getCountResult(projectsQueryBodyTpl, projectsCountTpl, data)
	return result, count, err
}
