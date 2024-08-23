package v1

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
)

type CreateBlacklistReqV1 struct {
	Type    string `json:"type" example:"sql" enums:"sql,fp_sql,ip,cidr,host,instance" valid:"required,oneof=sql fp_sql ip cidr host instance"`
	Desc    string `json:"desc" example:"used for rapid release"`
	Content string `json:"content" example:"select * from t1" valid:"required"`
}

// CreateBlacklist
// @Summary 添加黑名单
// @Description create a blacklist
// @Accept json
// @Id createBlacklistV1
// @Tags blacklist
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instance body v1.CreateBlacklistReqV1 true "add blacklist req"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/blacklist [post]
func CreateBlacklist(c echo.Context) error {
	req := new(CreateBlacklistReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	s := model.GetStorage()
	err = s.Save(&model.BlackListAuditPlanSQL{
		ProjectId:     model.ProjectUID(projectUid),
		FilterType:    model.BlacklistFilterType(req.Type),
		FilterContent: req.Content,
		Desc:          req.Desc,
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

// DeleteBlacklist
// @Description delete a blacklist
// @Id deleteBlackList
// @Tags blacklist
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param blacklist_id path string true "blacklist id"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/blacklist/{blacklist_id}/ [delete]
func DeleteBlacklist(c echo.Context) error {
	blacklistId := c.Param("blacklist_id")

	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	blacklist, exist, err := s.GetBlacklistByID(model.ProjectUID(projectUid), blacklistId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("blacklist is not exist")))
	}

	if err := s.Delete(blacklist); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type UpdateBlacklistReqV1 struct {
	Type    *string `json:"type" example:"sql" enums:"sql,fp_sql,ip,cidr,host,instance"`
	Desc    *string `json:"desc" example:"used for rapid release"`
	Content *string `json:"content" example:"select * from t1"`
}

// UpdateBlacklist
// @Summary 更新黑名单
// @Description update a blacklist
// @Accept json
// @Id updateBlacklistV1
// @Tags blacklist
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param blacklist_id path string true "blacklist id"
// @Param instance body v1.UpdateBlacklistReqV1 true "update blacklist req"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/blacklist/{blacklist_id}/ [patch]
func UpdateBlacklist(c echo.Context) error {
	req := new(UpdateBlacklistReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	blacklistId := c.Param("blacklist_id")
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	blacklist, exist, err := s.GetBlacklistByID(model.ProjectUID(projectUid), blacklistId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("blacklist is not exist")))
	}

	if req.Content != nil {
		blacklist.FilterContent = *req.Content
	}
	if req.Type != nil {
		blacklist.FilterType = model.BlacklistFilterType(*req.Type)
	}
	if req.Desc != nil {
		blacklist.Desc = *req.Desc
	}

	err = s.Save(blacklist)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type GetBlacklistReqV1 struct {
	FilterType         string `json:"filter_type" query:"filter_type" enums:"sql,fp_sql,ip,cidr,host,instance"  valid:"omitempty,oneof=sql fp_sql ip cidr host instance"`
	FuzzySearchContent string `json:"fuzzy_search_content" query:"fuzzy_search_content" valid:"omitempty"`
	PageIndex          uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize           uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetBlacklistResV1 struct {
	controller.BaseRes
	Data      []*BlacklistResV1 `json:"data"`
	TotalNums uint64            `json:"total_nums"`
}

type BlacklistResV1 struct {
	BlacklistID   uint       `json:"blacklist_id"`
	Content       string     `json:"content"`
	Desc          string     `json:"desc"`
	Type          string     `json:"type" enums:"sql,fp_sql,ip,cidr,host,instance"`
	MatchedCount  uint       `json:"matched_count"`
	LastMatchTime *time.Time `json:"last_match_time"`
}

// GetBlacklist
// @Summary 获取黑名单列表
// @Description get blacklist
// @Id getBlacklistV1
// @Tags blacklist
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param filter_type query string false "filter type" Enums(sql,fp_sql,ip,cidr,host,instance)
// @Param fuzzy_search_content query string false "fuzzy search content"
// @Param page_index query string true "page index"
// @Param page_size query string true "page size"
// @Success 200 {object} v1.GetBlacklistResV1
// @router /v1/projects/{project_name}/blacklist [get]
func GetBlacklist(c echo.Context) error {
	req := new(GetBlacklistReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	blacklistList, count, err := s.GetBlacklistList(model.ProjectUID(projectUid), model.BlacklistFilterType(req.FilterType), req.FuzzySearchContent, req.PageIndex, req.PageSize)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	res := make([]*BlacklistResV1, 0, len(blacklistList))
	for _, blacklist := range blacklistList {
		res = append(res, &BlacklistResV1{
			BlacklistID:   blacklist.ID,
			Content:       blacklist.FilterContent,
			Desc:          blacklist.Desc,
			Type:          string(blacklist.FilterType),
			MatchedCount:  blacklist.MatchedCount,
			LastMatchTime: blacklist.LastMatchTime,
		})
	}

	return c.JSON(http.StatusOK, &GetBlacklistResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      res,
		TotalNums: count,
	})
}
