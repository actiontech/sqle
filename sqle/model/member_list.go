package model

import (
	"fmt"
	"github.com/actiontech/sqle/sqle/errors"
)

type MemberDetail struct {
	UserName  string `json:"user_name"`
	IsManager bool   `json:"is_manager"`
}

var membersQueryTpl = `
SELECT users.login_name AS user_name , members.is_manager AS is_manager

{{ template "body" . }}

{{ if .limit }}
LIMIT :limit OFFSET :offset
{{ end }}
`

var membersCountTpl = `
SELECT COUNT(DISTINCT members.id)
{{ template "body" . }}
`

var membersQueryBodyTpl = `
{{ define "body" }}

FROM members
JOIN users ON users.id = members.user_id
JOIN projects ON projects.id = members.projects_id
JOIN instances ON instances.id = members.instance_id
WHERE users.stat = 0
AND users.deleted_at
AND projects.name = :filter_project_name

{{ if .filter_user_name }}
AND users.login_name = :filter_user_name
{{ end }}

{{ if .filter_instance_name }}
AND instances.name = :filter_instance_name
{{ end }}

{{ end }}
`

func (s *Storage) GetMembersByReq(data map[string]interface{}) (
	result []*MemberDetail, count uint64, err error) {

	if data["filter_project_name"] == nil {
		return nil, 0, errors.New(errors.DataInvalid, fmt.Errorf("project name must be exist"))
	}

	err = s.getListResult(membersQueryBodyTpl, membersQueryTpl, data, &result)
	if err != nil {
		return result, 0, err
	}
	count, err = s.getCountResult(membersQueryBodyTpl, membersCountTpl, data)
	return result, count, err
}
