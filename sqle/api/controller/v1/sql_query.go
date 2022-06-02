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
	SQLHistories []SQLHistoryItemResV1 `json:"sql_histories"`
}

type SQLHistoryItemResV1 struct {
	SQL string `json:"sql"`
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
// @router /v1/sql_query/history/{instance_name}/ [get]
func GetSQLQueryHistory(c echo.Context) error {
	return getSQLQueryHistory(c)
}

type GetSQLResultReqV1 struct {
	PageIndex uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetSQLResultResV1 struct {
	controller.BaseRes
	Data GetSQLResultResDataV1 `json:"data"`
}

type GetSQLResultResDataV1 struct {
	SQL         string                               `json:"sql"`
	StartLine   int                                  `json:"start_line"`
	EndLine     int                                  `json:"end_line"`
	CurrentPage int                                  `json:"current_page"`
	ExecuteTime int                                  `json:"execution_time"`
	Rows        []map[string] /* head name */ string `json:"rows"`
	Head        []SQLResultItemHeadResV1             `json:"head"`
}

type SQLResultItemHeadResV1 struct {
	FieldName string `json:"field_name"`
}

// GetSQLResult get sql query result
// @Summary 获取SQL查询结果
// @Description get sql query result
// @Id getSQLResult
// @Tags sql_query
// @Param query_id path string true "query sql id"
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetSQLResultResV1
// @router /v1/sql_query/results/{query_id}/ [get]
func GetSQLResult(c echo.Context) error {
	return getSQLResult(c)
}

type PrepareSQLQueryReqV1 struct {
	SQL            string `json:"sql" from:"sql"`
	InstanceSchema string `json:"instance_schema"`
}

type PrepareSQLQueryResV1 struct {
	controller.BaseRes
	Data PrepareSQLQueryResDataV1 `json:"data"`
}

type PrepareSQLQueryResDataV1 struct {
	QueryIds []PrepareSQLQueryResSQLV1 `json:"query_ids"`
}

type PrepareSQLQueryResSQLV1 struct {
	SQL     string `json:"sql"`
	QueryId string `json:"query_id"`
}

// PrepareSQLQuery prepare execute sql query
// @Summary 准备执行查询sql
// @Accept json
// @Description execute sql query
// @Id prepareSQLQuery
// @Tags sql_query
// @Param instance_name path string true "instance name"
// @Param req body v1.PrepareSQLQueryReqV1 true "exec sql"
// @Security ApiKeyAuth
// @Success 200 {object} v1.PrepareSQLQueryResV1
// @router /v1/sql_query/prepare/{instance_name}/ [post]
func PrepareSQLQuery(c echo.Context) error {
	return prepareSQLQuery(c)
}

type Table struct {
	Name string `json:"name"`
}

type ListTableBySchemaResV1 struct {
	controller.BaseRes
	Data []Table `json:"data"`
}

// ListTableBySchema list table by schema
// @Summary 获取数据库下的所有表
// @Description list table by schema
// @Id listTableBySchema
// @Tags sql_query
// @Param schema path string true "schema name"
// @Security ApiKeyAuth
// @Success 200 {object} v1.ListTableBySchemaResV1
// @router /v1/sql_query/tables/{schema}/ [get]
func ListTableBySchema(c echo.Context) error {
	return nil
}

type TableInfo struct {
	Items map[string]string `json:"items"`
}

type TableColumns struct {
	Rows []map[string]string      `json:"rows"`
	Head []SQLResultItemHeadResV1 `json:"head"`
}

type TableIndexes struct {
	Rows []map[string]string      `json:"rows"`
	Head []SQLResultItemHeadResV1 `json:"head"`
}

type TableMeta struct {
	Name    string       `json:"name"`
	Info    TableInfo    `json:"information"`
	Columns TableColumns `json:"columns"`
	Indexes TableIndexes `json:"indexes"`
}

type GetTableMetadataResV1 struct {
	controller.BaseRes
	Data TableMeta `json:"data"`
}

// GetTableMetadata get table metadata
// @Summary 获取表元数据
// @Description get table metadata
// @Id getTableMetadata
// @Tags sql_query
// @Param schema path string true "schema name"
// @Param table path string true "table name"
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetTableMetadataResV1
// @router /v1/sql_query/table/metadata/{schema}/{table}/ [get]
func GetTableMetadata(c echo.Context) error {
	return nil
}

type GetSQLExplainResV1 struct {
	controller.BaseRes
	Data SQLExplain `json:"data"`
}

// GetSQLExplain get SQL explain
// @Summary 获取SQL执行计划
// @Description get SQL explain
// @Id getSQLExplain
// @Tags sql_query
// @Param query_id path string true "query id"
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetSQLExplainResV1
// @router /v1/sql_query/explain/{query_id}/ [get]
func GetSQLExplain(c echo.Context) error {
	return nil
}
