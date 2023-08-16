package model

// import (
// 	"fmt"
// 	"time"

// 	"github.com/actiontech/sqle/sqle/errors"
// )

// type ProjectDetail struct {
// 	Name           string    `json:"name"`
// 	Desc           string    `json:"desc"`
// 	CreateUserName string    `json:"create_user_name"`
// 	CreateTime     time.Time `json:"create_time"`
// 	Status         string    `json:"status"`
// }

// var projectsQueryTpl = `SELECT
// DISTINCT projects.name , projects.` + "`desc`" + `, cu.login_name as create_user_name, projects.created_at as create_time, projects.status as status

// {{- template "body" . -}}

// {{ if .limit }}
// LIMIT :limit OFFSET :offset
// {{ end }}
// `

// var projectsCountTpl = `SELECT COUNT(DISTINCT projects.id)

// {{- template "body" . -}}
// `

// var projectsQueryBodyTpl = `
// {{ define "body" }}

// FROM projects
// LEFT JOIN project_user ON project_user.project_id = projects.id
// LEFT JOIN users ON users.id = project_user.user_id
// LEFT JOIN project_user_group ON project_user_group.project_id = projects.id
// LEFT JOIN user_group_users ON project_user_group.user_group_id = user_group_users.user_group_id
// LEFT JOIN users AS u ON u.id = user_group_users.user_id
// LEFT JOIN users AS cu ON cu.id = projects.create_user_id
// WHERE
// projects.deleted_at IS NULL

// {{ if .filter_user_name }}
// AND
// ((
// 	users.login_name = :filter_user_name
// 	AND users.deleted_at IS NULL
// 	AND users.stat = 0
// )
// OR (
// 	u.login_name = :filter_user_name
// 	AND u.deleted_at IS NULL
// 	AND u.stat = 0
// ))
// {{ end }}

// {{ end }}
// `

// func (s *Storage) GetProjectsByReq(data map[string]interface{}) (
// 	result []*ProjectDetail, count uint64, err error) {

// 	if data["filter_user_name"] == DefaultAdminUser {
// 		delete(data, "filter_user_name")
// 	}

// 	err = s.getListResult(projectsQueryBodyTpl, projectsQueryTpl, data, &result)
// 	if err != nil {
// 		return result, 0, err
// 	}
// 	count, err = s.getCountResult(projectsQueryBodyTpl, projectsCountTpl, data)
// 	return result, count, err
// }

// type MemberDetail struct {
// 	UserName  string `json:"user_name"`
// 	IsManager bool   `json:"is_manager"`
// }

// var membersQueryTpl = `
// SELECT DISTINCT users.login_name AS user_name

// {{ template "body" . }}

// {{ if .limit }}
// LIMIT :limit OFFSET :offset
// {{ end }}

// `

// var membersCountTpl = `
// SELECT COUNT(DISTINCT project_user.user_id)
// {{ template "body" . }}
// `

// var membersQueryBodyTpl = `
// {{ define "body" }}

// FROM project_user
// LEFT JOIN users ON users.id = project_user.user_id
// LEFT JOIN projects ON projects.id = project_user.project_id

// {{ if .filter_instance_name }}
// LEFT JOIN instances ON instances.project_id = projects.id
// {{ end }}

// {{  if .filter_instance_name  }}
// JOIN project_member_roles ON project_member_roles.user_id = users.id AND project_member_roles.instance_id = instances.id
// {{ end }}

// WHERE users.stat = 0
// AND users.deleted_at IS NULL
// AND projects.deleted_at IS NULL

// {{ if .filter_instance_name }}
// AND instances.deleted_at IS NULL
// {{ end }}

// AND projects.name = :filter_project_name

// {{ if .filter_user_name }}
// AND users.login_name = :filter_user_name
// {{ end }}

// {{ if .filter_instance_name }}
// AND instances.name = :filter_instance_name
// {{ end }}

// {{ end }}
// `

// func (s *Storage) GetMembersByReq(data map[string]interface{}) (
// 	result []*MemberDetail, count uint64, err error) {

// 	if data["filter_project_name"] == nil {
// 		return nil, 0, errors.New(errors.DataInvalid, fmt.Errorf("project name must be exist"))
// 	}

// 	members := []*struct {
// 		UserName string `json:"user_name"`
// 	}{}
// 	err = s.getListResult(membersQueryBodyTpl, membersQueryTpl, data, &members)
// 	if err != nil {
// 		return result, 0, err
// 	}

// 	project, _, err := s.GetProjectByName(fmt.Sprintf("%v", data["filter_project_name"]))
// 	if err != nil {
// 		return result, 0, err
// 	}

// 	for _, member := range members {
// 		isManager := false
// 		for _, manager := range project.Managers {
// 			if manager.Name == member.UserName {
// 				isManager = true
// 				break
// 			}
// 		}
// 		result = append(result, &MemberDetail{
// 			UserName:  member.UserName,
// 			IsManager: isManager,
// 		})
// 	}

// 	count, err = s.getCountResult(membersQueryBodyTpl, membersCountTpl, data)
// 	return result, count, err
// }
