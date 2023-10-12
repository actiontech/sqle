package v1

import (
	"context"
	"fmt"

	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
)

type CreateAuditWhitelistReqV1 struct {
	Value     string `json:"value" example:"create table" valid:"required"`
	MatchType string `json:"match_type" example:"exact_match" enums:"exact_match,fp_match" valid:"omitempty,oneof=exact_match fp_match"`
	Desc      string `json:"desc" example:"used for rapid release"`
}

// @Summary 添加SQL白名单
// @Description create a sql whitelist
// @Accept json
// @Id createAuditWhitelistV1
// @Tags audit_whitelist
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instance body v1.CreateAuditWhitelistReqV1 true "add sql whitelist req"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/audit_whitelist [post]
func CreateAuditWhitelist(c echo.Context) error {
	req := new(CreateAuditWhitelistReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"),true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	s := model.GetStorage()

	sqlWhitelist := &model.SqlWhitelist{
		ProjectId: model.ProjectUID(projectUid),
		Value:     req.Value,
		Desc:      req.Desc,
		MatchType: req.MatchType,
	}

	err = s.Save(sqlWhitelist)
	if err != nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type UpdateAuditWhitelistReqV1 struct {
	Value     *string `json:"value" example:"create table"`
	MatchType *string `json:"match_type" example:"exact_match" enums:"exact_match,fp_match"`
	Desc      *string `json:"desc" example:"used for rapid release"`
}

// @Summary 更新SQL白名单
// @Description update sql whitelist by id
// @Accept json
// @Id UpdateAuditWhitelistByIdV1
// @Tags audit_whitelist
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param audit_whitelist_id path string true "sql audit whitelist id"
// @Param instance body v1.UpdateAuditWhitelistReqV1 true "update sql whitelist req"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/audit_whitelist/{audit_whitelist_id}/ [patch]
func UpdateAuditWhitelistById(c echo.Context) error {
	req := new(UpdateAuditWhitelistReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()

	whitelistId := c.Param("audit_whitelist_id")
	sqlWhitelist, exist, err := s.GetSqlWhitelistByIdAndProjectUID(whitelistId, model.ProjectUID(projectUid))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("sql audit whitelist is not exist")))
	}

	// nothing to update
	if req.Value == nil && req.Desc == nil && req.MatchType == nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
	}

	if req.Value != nil {
		sqlWhitelist.Value = *req.Value
	}
	if req.MatchType != nil {
		sqlWhitelist.MatchType = *req.MatchType
	}
	if req.Desc != nil {
		sqlWhitelist.Desc = *req.Desc
	}

	err = s.Save(sqlWhitelist)
	if err != nil {
		return c.JSON(http.StatusOK, controller.NewBaseReq(err))
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

// @Summary 删除SQL白名单信息
// @Description remove sql white
// @Id deleteAuditWhitelistByIdV1
// @Tags audit_whitelist
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param audit_whitelist_id path string true "audit whitelist id"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/audit_whitelist/{audit_whitelist_id}/ [delete]
func DeleteAuditWhitelistById(c echo.Context) error {
	s := model.GetStorage()
	whitelistId := c.Param("audit_whitelist_id")
	// projectName := c.Param("project_name")

	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	sqlWhitelist, exist, err := s.GetSqlWhitelistByIdAndProjectUID(whitelistId, model.ProjectUID(projectUid))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("sql audit whitelist is not exist")))
	}
	err = s.Delete(sqlWhitelist)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type GetAuditWhitelistReqV1 struct {
	PageIndex uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetAuditWhitelistResV1 struct {
	controller.BaseRes
	Data      []*AuditWhitelistResV1 `json:"data"`
	TotalNums uint32                 `json:"total_nums"`
}

type AuditWhitelistResV1 struct {
	Id        uint   `json:"audit_whitelist_id"`
	Value     string `json:"value"`
	MatchType string `json:"match_type"`
	Desc      string `json:"desc"`
}

// @Summary 获取Sql审核白名单
// @Description get all whitelist
// @Id getAuditWhitelistV1
// @Tags audit_whitelist
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param page_index query string true "page index"
// @Param page_size query string true "page size"
// @Success 200 {object} v1.GetAuditWhitelistResV1
// @router /v1/projects/{project_name}/audit_whitelist [get]
func GetSqlWhitelist(c echo.Context) error {
	req := new(GetAuditWhitelistReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	sqlWhitelist, count, err := s.GetSqlWhitelistByProjectUID(req.PageIndex, req.PageSize, model.ProjectUID(projectUid))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	whitelistRes := make([]*AuditWhitelistResV1, 0, len(sqlWhitelist))
	for _, v := range sqlWhitelist {
		whitelistRes = append(whitelistRes, &AuditWhitelistResV1{
			Id:        v.ID,
			Value:     v.Value,
			Desc:      v.Desc,
			MatchType: v.MatchType,
		})
	}
	return c.JSON(http.StatusOK, &GetAuditWhitelistResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      whitelistRes,
		TotalNums: count,
	})
}
