package v1

import (
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

type GetProjectResV1 struct {
	controller.BaseRes
	Data      []*ProjectListItem `json:"data"`
	TotalNums uint64             `json:"total_nums"`
}

type ProjectListItem struct {
	Id             uint       `json:"id"`
	Name           string     `json:"name"`
	Desc           string     `json:"desc"`
	CreateUserName string     `json:"create_user_name"`
	CreateTime     *time.Time `json:"create_time"`
}

// GetProjectListV1
// @Summary 获取项目列表
// @Description get project list
// @Tags project
// @Id getProjectListV1
// @Security ApiKeyAuth
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page" default(50)
// @Success 200 {object} v1.GetProjectResV1
// @router /v1/projects [get]
func GetProjectListV1(c echo.Context) error {
	return nil
}

type GetProjectDetailResV1 struct {
	controller.BaseRes
	Data ProjectDetailItem `json:"data"`
}

type ProjectDetailItem struct {
	Id             uint       `json:"id"`
	Name           string     `json:"name"`
	Desc           string     `json:"desc"`
	CreateUserName string     `json:"create_user_name"`
	CreateTime     *time.Time `json:"create_time"`
}

// GetProjectDetailV1
// @Summary 获取项目详情
// @Description get project detail
// @Tags project
// @Id getProjectDetailV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Success 200 {object} v1.GetProjectDetailResV1
// @router /v1/projects/{project_name}/ [get]
func GetProjectDetailV1(c echo.Context) error {
	return nil
}

type CreateProjectReqV1 struct {
	Name string `json:"name" valid:"required"`
	Desc string `json:"desc"`
}

// CreateProjectV1
// @Summary 创建项目
// @Description create project
// @Accept json
// @Produce json
// @Tags project
// @Id createProjectV1
// @Security ApiKeyAuth
// @Param project body v1.CreateProjectReqV1 true "create project request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects [post]
func CreateProjectV1(c echo.Context) error {
	return nil
}

type UpdateProjectReqV1 struct {
	Desc *string `json:"desc"`
}

// UpdateProjectV1
// @Summary 更新项目
// @Description update project
// @Accept json
// @Produce json
// @Tags project
// @Id updateProjectV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param project body v1.UpdateProjectReqV1 true "create project request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/ [patch]
func UpdateProjectV1(c echo.Context) error {
	return nil
}

// DeleteProjectV1
// @Summary 删除项目
// @Description delete project
// @Id deleteProjectV1
// @Tags project
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/ [delete]
func DeleteProjectV1(c echo.Context) error {
	return nil
}

type GetProjectTipsResV1 struct {
	controller.BaseRes
	Data []ProjectTipResV1 `json:"data"`
}

type ProjectTipResV1 struct {
	Name string `json:"project_name"`
}

// GetProjectTipsV1
// @Summary 获取项目提示列表
// @Description get project tip list
// @Tags project
// @Id getProjectTipsV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetProjectTipsResV1
// @router /v1/project_tips [get]
func GetProjectTipsV1(c echo.Context) error {
	return nil
}
