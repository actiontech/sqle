package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"actiontech.cloud/universe/sqle/v4/sqle/utils"

	"actiontech.cloud/universe/sqle/v4/sqle/model"
	"github.com/labstack/echo/v4"
)

type SqlWhitelistItemRes struct {
	BaseRes
	Data []model.SqlWhitelist `json:"data"`
}

type GetSqlWhitelistRes struct {
	BaseRes
	Data      []model.SqlWhitelist `json:"data"`
	TotalNums uint32               `json:"total_nums"`
}

// @Summary 获取指定SQL白名单信息
// @Description get sql whitelist item
// @Param sql_whitelist_id path string true "sql whitelist item ID"
// @Success 200 {object} controller.SqlWhitelistItemRes
// @router /sql_whitelist/{sql_whitelist_id}/ [get]
func GetSqlWhitelistItemById(c echo.Context) error {
	s := model.GetStorage()
	sqlWhiteId := c.Param("sql_whitelist_id")
	sqlWhitelistItem, exist, err := s.GetSqlWhitelistItemById(sqlWhiteId)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	if !exist {
		return c.JSON(http.StatusOK, NewBaseReq(fmt.Errorf("sql whitelist is not exist")))
	}
	return c.JSON(http.StatusOK, &SqlWhitelistItemRes{
		BaseRes: NewBaseReq(nil),
		Data:    []model.SqlWhitelist{*sqlWhitelistItem},
	})
}

type CreateSqlWhitelistItemReq struct {
	Value *string `json:"value" example:"create table" valid:"required"`
	Desc  *string `json:"desc" example:"used for rapid release" valid:"-"`
}

// @Summary 添加SQL白名单
// @Description create a sql whitelist item
// @Accept json
// @Param instance body controller.CreateSqlWhitelistItemReq true "add sql whitelist item"
// @Success 200 {object} controller.SqlWhitelistItemRes
// @router /sql_whitelist [post]
func CreateSqlWhitelistItem(c echo.Context) error {
	s := model.GetStorage()
	req := new(CreateSqlWhitelistItemReq)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}

	sqlWhitelistItem := &model.SqlWhitelist{
		Value:         *req.Value,
		Desc:          *req.Desc,
		MessageDigest: utils.Md5String(strings.ToUpper(*req.Value)),
	}
	err := s.Save(sqlWhitelistItem)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	sqlWhitelistItem.PutSqlWhitelistMD5()
	return c.JSON(http.StatusOK, &SqlWhitelistItemRes{
		BaseRes: NewBaseReq(nil),
		Data:    []model.SqlWhitelist{*sqlWhitelistItem},
	})
}

// @Summary 更新SQL白名单
// @Description update a sql whitelist item
// @Accept json
// @Param sql_whitelist_id path string true "sql whitelist item ID"
// @Param instance body controller.CreateSqlWhitelistItemReq true "update sql whitelist item"
// @Success 200 {object} controller.SqlWhitelistItemRes
// @router /sql_whitelist/{sql_whitelist_id}/ [patch]
func UpdateSqlWhitelistItem(c echo.Context) error {
	s := model.GetStorage()
	sqlWhiteId := c.Param("sql_whitelist_id")
	sqlWhitelistItem, exist, err := s.GetSqlWhitelistItemById(sqlWhiteId)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	if !exist {
		return c.JSON(http.StatusOK, NewBaseReq(fmt.Errorf("sql whitelist is not exist")))
	}
	req := new(CreateSqlWhitelistItemReq)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	sqlWhitelistItem.Value = *req.Value
	sqlWhitelistItem.Desc = *req.Desc
	sqlWhitelistItem.MessageDigest = utils.Md5String(strings.ToUpper(*req.Value))
	err = s.Save(sqlWhitelistItem)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	sqlWhitelistItem.PutSqlWhitelistMD5()
	return c.JSON(http.StatusOK, &SqlWhitelistItemRes{
		BaseRes: NewBaseReq(nil),
		Data:    []model.SqlWhitelist{*sqlWhitelistItem},
	})
}

// @Summary 获取Sql审核白名单
// @Description get all whitelist
// @Param page_index query string false "page index"
// @Param page_size query string false "page size"
// @Success 200 {object} controller.GetSqlWhitelistRes
// @router /sql_whitelist [get]
func GetSqlWhitelist(c echo.Context) error {
	s := model.GetStorage()
	index, err := url.QueryUnescape(c.QueryParam("page_index"))
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	pageIndex, err := FormatStringToInt(index)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	size, err := url.QueryUnescape(c.QueryParam("page_size"))
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	pageSize, err := FormatStringToInt(size)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	sqlWhitelist, count, err := s.GetSqlWhitelist(pageIndex, pageSize)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}

	return c.JSON(http.StatusOK, &GetSqlWhitelistRes{
		BaseRes:   NewBaseReq(nil),
		Data:      sqlWhitelist,
		TotalNums: count,
	})
}

// @Summary 删除SQL白名单信息
// @Description remove sql white
// @Param sql_whitelist_id path string true "sql whitelist item ID"
// @Success 200 {object} controller.SqlWhitelistItemRes
// @router /sql_whitelist/{sql_whitelist_id}/ [delete]
func RemoveSqlWhitelistItem(c echo.Context) error {
	s := model.GetStorage()
	sqlWhiteId := c.Param("sql_whitelist_id")
	sqlWhitelistItem, exist, err := s.GetSqlWhitelistItemById(sqlWhiteId)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	if !exist {
		return c.JSON(http.StatusOK, NewBaseReq(fmt.Errorf("sql whitelist is not exist")))
	}
	err = s.Delete(sqlWhitelistItem)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	sqlWhitelistItem.RemoveSqlWhitelistMD5()
	return c.JSON(http.StatusOK, NewBaseReq(nil))
}
