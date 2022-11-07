package model

import "time"

type ProjectDetail struct {
	Name           string    `json:"name"`
	Desc           string    `json:"desc"`
	CreateUserName string    `json:"create_user_name"`
	CreateTime     time.Time `json:"create_time"`
}

var projectsQueryTpl = `SELECT
projects.name , projects.` + "`desc`" + `, users.login_name as create_user_name, projects.created_at as create_time

{{- template "body" . -}}

{{ if .limit }}
LIMIT :limit OFFSET :offset
{{ end }}
`

var projectsCountTpl = `SELECT COUNT(DISTINCT projects.id)

{{- template "body" . -}}
`

var projectsQueryBodyTpl = `
{{ define "body" }}

FROM projects
JOIN project_user on project_user.project_id = projects.id
JOIN users on users.id = project_user.user_id
JOIN project_user_group on project_user_group.project_id = projects.id
JOIN user_group_users on project_user_group.user_group_id = user_group_users.user_group_id
RIGHT JOIN users as u on u.id = user_group_users.user_id
WHERE users.deleted_at IS NULL
AND u.deleted_at IS NULL
AND users.stat = 0
AND u.stat = 0

{{ if .filter_user_name }}
AND
( 
	users.login_name = :filter_user_name
OR 
	u.login_name = :filter_user_name
)
{{ end }}


{{ end }}
`

func (s *Storage) GetProjectsByReq(data map[string]interface{}) (
	result []*ProjectDetail, count uint64, err error) {

	if data["filter_user_name"] == DefaultAdminUser {
		delete(data, "filter_user_name")
	}

	err = s.getListResult(projectsQueryBodyTpl, projectsQueryTpl, data, &result)
	if err != nil {
		return result, 0, err
	}
	count, err = s.getCountResult(projectsQueryBodyTpl, projectsCountTpl, data)
	return result, count, err
}
