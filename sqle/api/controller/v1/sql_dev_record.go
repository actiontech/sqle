package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

type GetSqlDEVRecordListReq struct {
	FuzzySearchSqlFingerprint *string `query:"fuzzy_search_sql_fingerprint" json:"fuzzy_search_sql_fingerprint,omitempty"`
	FilterInstanceName        *string `query:"filter_instance_name" json:"filter_instance_name,omitempty"`
	FilterCreator             *string `query:"filter_creator" json:"filter_creator,omitempty"`
	FilterSource              *string `query:"filter_source" json:"filter_source,omitempty"`
	FilterLastReceiveTimeFrom *string `query:"filter_last_receive_time_from" json:"filter_last_receive_time_from,omitempty"`
	FilterLastReceiveTimeTo   *string `query:"filter_last_receive_time_to" json:"filter_last_receive_time_to,omitempty"`
	FuzzySearchSchemaName     *string `query:"fuzzy_search_schema_name" json:"fuzzy_search_schema_name,omitempty"`
	SortField                 *string `query:"sort_field" json:"sort_field,omitempty" valid:"omitempty,oneof=first_appear_timestamp last_receive_timestamp fp_count" enums:"first_appear_timestamp,last_receive_timestamp,fp_count"`
	SortOrder                 *string `query:"sort_order" json:"sort_order,omitempty" valid:"omitempty,oneof=asc desc" enums:"asc,desc"`
	PageIndex                 uint32  `query:"page_index" valid:"required" json:"page_index"`
	PageSize                  uint32  `query:"page_size" valid:"required" json:"page_size"`
}

type GetSqlDEVRecordListResp struct {
	controller.BaseRes
	Data      []*SqlDEVRecord `json:"data"`
	TotalNums uint64          `json:"total_nums"`
}

type SqlDEVRecord struct {
	Id                   uint64         `json:"id"`
	SqlFingerprint       string         `json:"sql_fingerprint"`
	Sql                  string         `json:"sql"`
	Source               *RecordSource  `json:"source"`
	InstanceName         string         `json:"instance_name"`
	SchemaName           string         `json:"schema_name"`
	AuditResult          []*AuditResult `json:"audit_result"`
	FirstAppearTimeStamp string         `json:"first_appear_timestamp"`
	LastReceiveTimeStamp string         `json:"last_receive_timestamp"`
	FpCount              uint64         `json:"fp_count"`
	Creator              string         `json:"creator"` // create user name
}

type RecordSource struct {
	Name  string `json:"name" enums:"ide_plugin"`
	Value string `json:"value"`
}

// GetSqlDEVRecordList
// @Summary 获取开发sql记录
// @Description get sql dev record list
// @Tags SqlDEVRecord
// @Id GetSqlDEVRecordList
// @Security ApiKeyAuth
// @Param project_id path string true "project id"
// @Param fuzzy_search_sql_fingerprint query string false "fuzzy search sql fingerprint"
// @Param filter_instance_name query string false "instance name"
// @Param filter_source query string false "source" Enums(ide_plugin)
// @Param filter_creator query string false "creator name"
// @Param filter_last_receive_time_from query string false "last receive time from"
// @Param filter_last_receive_time_to query string false "last receive time to"
// @Param sort_field query string false "sort field" Enums(first_appear_timestamp,last_receive_timestamp,fp_count)
// @Param sort_order query string false "sort order" Enums(asc,desc)
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Success 200 {object} v1.GetSqlDEVRecordListResp
// @Router /v1/projects/{project_name}/sql_dev_records [get]
func GetSqlDEVRecordList(c echo.Context) error {
	return getSqlDEVRecordList(c)
}
