package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"

	"github.com/labstack/echo/v4"
)

type GetSQLQueryHistoryReqV1 struct {
	FilterFuzzySearch string `json:"filter_fuzzy_search" query:"filter_fuzzy_search"`
	PageIndex         uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize          uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetSQLQueryHistoryResV1 struct {
	controller.BaseRes
	Data GetSQLQueryHistoryResDataV1 `json:"data"`
}

type GetSQLQueryHistoryResDataV1 struct {
	SQLHistories []string `json:"sql_histories"`
}

// GetSQLQueryHistory get current user sql query history
// @Summary 获取当前用户历史查询SQL
// @Description get sql query history
// @Id getSQLQueryHistory
// @Tags sql_query
// @Param instance_name path string true "instance name"
// @Param filter_fuzzy_search query string false "fuzzy search filter"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetSQLQueryHistoryResV1
// @router /v1/sql_query/{instance_name}/history [get]
func GetSQLQueryHistory(c echo.Context) error {
	return nil
}

type ExecSQLQueryReqV1 struct {
	PageIndex uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required"`
	SQL       string `json:"sql" from:"sql"`
}

type ExecSQLQueryResV1 struct {
	controller.BaseRes
	Data ExecSQLQueryResDataV1 `json:"data"`
}

type ExecSQLQueryResDataV1 struct {
	StartLine     int                           `json:"start_line"`
	EndLine       int                           `json:"end_line"`
	CurrentPage   int                           `json:"current_page"`
	ExecuteTime   int                           `json:"execution_time"`
	ExecuteResult []ExecSQLQueryResResultItemV1 `json:"execute_result"`
}

type ExecSQLQueryResResultItemV1 struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ExecSQLQuery execute sql query
// @Summary 执行查询sql
// @Accept json
// @Description execute sql query
// @Id execSQLQuery
// @Tags sql_query
// @Param instance_name path string true "instance name"
// @Param exec body v1.ExecSQLQueryReqV1 true "exec sql"
// @Security ApiKeyAuth
// @Success 200 {object} v1.ExecSQLQueryResV1
// @router /v1/sql_query/{instance_name}/execute [post]
func ExecSQLQuery(c echo.Context) error {
	return nil
}
