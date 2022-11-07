package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
)

type GetProjectReqV1 struct {
	PageIndex uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetProjectResV1 struct {
	controller.BaseRes
	Data      []*ProjectListItem `json:"data"`
	TotalNums uint64             `json:"total_nums"`
}

type ProjectListItem struct {
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
	req := new(GetProjectReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	limit, offset := controller.GetLimitAndOffset(req.PageIndex, req.PageSize)

	user := controller.GetUserName(c)

	mp := map[string]interface{}{
		"limit":            limit,
		"offset":           offset,
		"filter_user_name": user,
	}

	s := model.GetStorage()
	projects, total, err := s.GetProjectsByReq(mp)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	resp := []*ProjectListItem{}
	for _, project := range projects {
		resp = append(resp, &ProjectListItem{
			Name:           project.Name,
			Desc:           project.Desc,
			CreateUserName: project.CreateUserName,
			CreateTime:     &project.CreateTime,
		})
	}

	return c.JSON(http.StatusOK, &GetProjectResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      resp,
		TotalNums: total,
	})
}

type GetProjectDetailResV1 struct {
	controller.BaseRes
	Data ProjectDetailItem `json:"data"`
}

type ProjectDetailItem struct {
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
	return createProjectV1(c)
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
	req := new(UpdateProjectReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	projectIDStr := c.Param("filter_project")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("project id should be uint but not"))
	}

	s := model.GetStorage()
	sure, err := s.CheckUserCanUpdateProject(uint(projectID), user.ID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !sure {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("you can not modify this project"))
	}

	attr := map[string]interface{}{}
	if req.Desc != nil {
		attr["desc"] = *req.Desc
	}

	return controller.JSONBaseErrorReq(c, s.UpdateProjectInfoByID(uint(projectID), attr))
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
	return deleteProjectV1(c)
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
