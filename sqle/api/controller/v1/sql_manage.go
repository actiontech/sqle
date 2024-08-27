package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

type GetSqlManageListReq struct {
	FuzzySearchSqlFingerprint    *string `query:"fuzzy_search_sql_fingerprint" json:"fuzzy_search_sql_fingerprint,omitempty"`
	FilterAssignee               *string `query:"filter_assignee" json:"filter_assignee,omitempty"`
	FilterInstanceID             *string `query:"filter_instance_id" json:"filter_instance_id,omitempty"`
	FilterSource                 *string `query:"filter_source" json:"filter_source,omitempty"`
	FilterAuditLevel             *string `query:"filter_audit_level" json:"filter_audit_level,omitempty"`
	FilterLastAuditStartTimeFrom *string `query:"filter_last_audit_start_time_from" json:"filter_last_audit_start_time_from,omitempty"`
	FilterLastAuditStartTimeTo   *string `query:"filter_last_audit_start_time_to" json:"filter_last_audit_start_time_to,omitempty"`
	FilterStatus                 *string `query:"filter_status" json:"filter_status,omitempty"`
	FilterDbType                 *string `query:"filter_db_type" json:"filter_db_type,omitempty"`
	FilterRuleName               *string `query:"filter_rule_name" json:"filter_rule_name,omitempty"`
	FilterBusiness               *string `query:"filter_business" json:"filter_business,omitempty"`
	FilterPriority               *string `query:"filter_priority" json:"filter_priority,omitempty" enums:"high,low"`
	FuzzySearchEndpoint          *string `query:"fuzzy_search_endpoint" json:"fuzzy_search_endpoint,omitempty"`
	FuzzySearchSchemaName        *string `query:"fuzzy_search_schema_name" json:"fuzzy_search_schema_name,omitempty"`
	SortField                    *string `query:"sort_field" json:"sort_field,omitempty" valid:"omitempty,oneof=first_appear_timestamp last_receive_timestamp fp_count" enums:"first_appear_timestamp,last_receive_timestamp,fp_count"`
	SortOrder                    *string `query:"sort_order" json:"sort_order,omitempty" valid:"omitempty,oneof=asc desc" enums:"asc,desc"`
	PageIndex                    uint32  `query:"page_index" valid:"required" json:"page_index"`
	PageSize                     uint32  `query:"page_size" valid:"required" json:"page_size"`
}

type GetSqlManageListResp struct {
	controller.BaseRes
	Data                  []*SqlManage `json:"data"`
	SqlManageTotalNum     uint64       `json:"sql_manage_total_num"`
	SqlManageBadNum       uint64       `json:"sql_manage_bad_num"`
	SqlManageOptimizedNum uint64       `json:"sql_manage_optimized_num"`
}

type SqlManage struct {
	Id              uint64         `json:"id"`
	SqlFingerprint  string         `json:"sql_fingerprint"`
	Sql             string         `json:"sql"`
	Source          *Source        `json:"source"`
	InstanceName    string         `json:"instance_name"`
	SchemaName      string         `json:"schema_name"`
	AuditResult     []*AuditResult `json:"audit_result"`
	FirstAppearTime string         `json:"first_appear_time"`
	LastAppearTime  string         `json:"last_appear_time"`
	AppearNum       uint64         `json:"appear_num"`
	Assignees       []string       `json:"assignees"`
	Status          string         `json:"status" enums:"unhandled,solved,ignored,manual_audited"`
	Remark          string         `json:"remark"`
	Endpoint        string         `json:"endpoint"`
}

// AuditResult 用于SQL全生命周期展示的AuditResult
type AuditResult struct {
	Level    string `json:"level" example:"warn"`
	Message  string `json:"message" example:"避免使用不必要的内置函数md5()"`
	RuleName string `json:"rule_name"`
}

type Source struct {
	SqlSourceType string   `json:"sql_source_type"`
	SqlSourceDesc string   `json:"sql_source_desc"`
	SqlSourceIDs  []string `json:"sql_source_ids"`
}

// todo : 该接口已废弃，后续会删除
// GetSqlManageList
// @Deprecated
// @Summary 获取管控sql列表
// @Description get sql manage list
// @Tags SqlManage
// @Id GetSqlManageList
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param fuzzy_search_sql_fingerprint query string false "fuzzy search sql fingerprint"
// @Param filter_assignee query string false "assignee"
// @Param filter_instance_name query string false "instance name"
// @Param filter_source query string false "source" Enums(audit_plan,sql_audit_record)
// @Param filter_audit_level query string false "audit level" Enums(normal,notice,warn,error)
// @Param filter_last_audit_start_time_from query string false "last audit start time from"
// @Param filter_last_audit_start_time_to query string false "last audit start time to"
// @Param filter_status query string false "status" Enums(unhandled,solved,ignored,manual_audited)
// @Param filter_rule_name query string false "rule name"
// @Param filter_db_type query string false "db type"
// @Param fuzzy_search_endpoint query string false "fuzzy search endpoint"
// @Param fuzzy_search_schema_name query string false "fuzzy search schema name"
// @Param sort_field query string false "sort field" Enums(first_appear_timestamp,last_receive_timestamp,fp_count)
// @Param sort_order query string false "sort order" Enums(asc,desc)
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Success 200 {object} v1.GetSqlManageListResp
// @Router /v1/projects/{project_name}/sql_manages [get]
func GetSqlManageList(c echo.Context) error {
	return getSqlManageList(c)
}

type BatchUpdateSqlManageReq struct {
	SqlManageIdList []*uint64 `json:"sql_manage_id_list"`
	Status          *string   `json:"status" enums:"solved,ignored,manual_audited"`
	Priority        *string   `json:"priority" enums:",high"`
	Assignees       []string  `json:"assignees"`
	Remark          *string   `json:"remark"`
}

// BatchUpdateSqlManage batch update sql manage
// @Summary 批量更新SQL管控
// @Description batch update sql manage
// @Tags SqlManage
// @Id BatchUpdateSqlManage
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param BatchUpdateSqlManageReq body BatchUpdateSqlManageReq true "batch update sql manage request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/sql_manages/batch [PATCH]
func BatchUpdateSqlManage(c echo.Context) error {
	return batchUpdateSqlManage(c)
}

type ExportSqlManagesReq struct {
	FuzzySearchSqlFingerprint    *string `query:"fuzzy_search_sql_fingerprint" json:"fuzzy_search_sql_fingerprint,omitempty"`
	FilterAssignee               *string `query:"filter_assignee" json:"filter_assignee,omitempty"`
	FilterBusiness               *string `query:"filter_business" json:"filter_business,omitempty"`
	FilterInstanceID             *string `query:"filter_instance_id" json:"filter_instance_id,omitempty"`
	FilterSource                 *string `query:"filter_source" json:"filter_source,omitempty"`
	FilterAuditLevel             *string `query:"filter_audit_level" json:"filter_audit_level,omitempty"`
	FilterLastAuditStartTimeFrom *string `query:"filter_last_audit_start_time_from" json:"filter_last_audit_start_time_from,omitempty"`
	FilterLastAuditStartTimeTo   *string `query:"filter_last_audit_start_time_to" json:"filter_last_audit_start_time_to,omitempty"`
	FilterStatus                 *string `query:"filter_status" json:"filter_status,omitempty"`
	FilterDbType                 *string `query:"filter_db_type" json:"filter_db_type,omitempty"`
	FilterRuleName               *string `query:"filter_rule_name" json:"filter_rule_name,omitempty"`
	FilterPriority               *string `query:"filter_priority" json:"filter_priority,omitempty" enums:"high,low"`
	FuzzySearchEndpoint          *string `query:"fuzzy_search_endpoint" json:"fuzzy_search_endpoint,omitempty"`
	FuzzySearchSchemaName        *string `query:"fuzzy_search_schema_name" json:"fuzzy_search_schema_name,omitempty"`
	SortField                    *string `query:"sort_field" json:"sort_field,omitempty" valid:"omitempty,oneof=first_appear_timestamp last_receive_timestamp fp_count" enums:"first_appear_timestamp,last_receive_timestamp,fp_count"`
	SortOrder                    *string `query:"sort_order" json:"sort_order,omitempty" valid:"omitempty,oneof=asc desc" enums:"asc,desc"`
}

// ExportSqlManagesV1
// @Summary 导出SQL管控
// @Description export sql manage
// @Id exportSqlManageV1
// @Tags SqlManage
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param fuzzy_search_sql_fingerprint query string false "fuzzy search sql fingerprint"
// @Param filter_assignee query string false "assignee"
// @Param filter_business query string false "business"
// @Param filter_priority query string false "priority" Enums(high,low)
// @Param filter_instance_id query string false "instance id"
// @Param filter_source query string false "source" Enums(audit_plan,sql_audit_record)
// @Param filter_audit_level query string false "audit level" Enums(normal,notice,warn,error)
// @Param filter_last_audit_start_time_from query string false "last audit start time from"
// @Param filter_last_audit_start_time_to query string false "last audit start time to"
// @Param filter_status query string false "status" Enums(unhandled,solved,ignored,manual_audited)
// @Param filter_db_type query string false "db type"
// @Param filter_rule_name query string false "rule name"
// @Param fuzzy_search_endpoint query string false "fuzzy search endpoint"
// @Param fuzzy_search_schema_name query string false "fuzzy search schema name"
// @Param sort_field query string false "sort field" Enums(first_appear_timestamp,last_receive_timestamp,fp_count)
// @Param sort_order query string false "sort order" Enums(asc,desc)
// @Success 200 {file} file "export sql manage"
// @Router /v1/projects/{project_name}/sql_manages/exports [get]
func ExportSqlManagesV1(c echo.Context) error {
	return exportSqlManagesV1(c)
}

type RuleRespV1 struct {
	RuleName string `json:"rule_name"`
	Desc     string `json:"desc"`
}

type RuleTips struct {
	DbType string       `json:"db_type"`
	Rule   []RuleRespV1 `json:"rule"`
}

type GetSqlManageRuleTipsResp struct {
	controller.BaseRes
	Data []RuleTips `json:"data"`
}

// GetSqlManageRuleTips
// @Summary 获取管控规则tips
// @Description get sql manage rule tips
// @Id GetSqlManageRuleTips
// @Tags SqlManage
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Success 200 {object} v1.GetSqlManageRuleTipsResp
// @Router /v1/projects/{project_name}/sql_manages/rule_tips [get]
func GetSqlManageRuleTips(c echo.Context) error {
	return getSqlManageRuleTips(c)
}

type AffectRows struct {
	Count      int    `json:"count"`
	ErrMessage string `json:"err_message"`
}

type PerformanceStatistics struct {
	AffectRows *AffectRows `json:"affect_rows"`
}

type TableMetas struct {
	ErrMessage string       `json:"err_message"`
	Items      []*TableMeta `json:"table_meta_items"`
}

type SqlAnalysis struct {
	SQLExplain            *SQLExplain            `json:"sql_explain"`
	TableMetas            *TableMetas            `json:"table_metas"`
	PerformanceStatistics *PerformanceStatistics `json:"performance_statistics"`
}

type GetSqlManageSqlAnalysisResp struct {
	controller.BaseRes
	// V1版本不能引用V2版本的结构体,所以只能复制一份
	Data *SqlAnalysis `json:"data"`
}

// GetSqlManageSqlAnalysisV1
// @Summary 获取SQL管控SQL分析
// @Description get sql manage analysis
// @Id GetSqlManageSqlAnalysisV1
// @Tags SqlManage
// @Param project_name path string true "project name"
// @Param sql_manage_id path string true "sql manage id"
// @Security ApiKeyAuth
// @Success 200 {object} GetSqlManageSqlAnalysisResp
// @Router /v1/projects/{project_name}/sql_manages/{sql_manage_id}/sql_analysis [get]
func GetSqlManageSqlAnalysisV1(c echo.Context) error {
	return getSqlManageSqlAnalysisV1(c)
}

func convertSQLAnalysisResultToRes(res *AnalysisResult, rawSQL string) *SqlAnalysis {
	data := &SqlAnalysis{}

	// explain
	{
		data.SQLExplain = &SQLExplain{SQL: rawSQL}
		if res.ExplainResultErr != nil {
			data.SQLExplain.Message = res.ExplainResultErr.Error()
		} else {
			classicResult := ExplainClassicResult{
				Head: make([]TableMetaItemHeadResV1, len(res.ExplainResult.ClassicResult.Columns)),
				Rows: make([]map[string]string, len(res.ExplainResult.ClassicResult.Rows)),
			}

			// head
			for i := range res.ExplainResult.ClassicResult.Columns {
				col := res.ExplainResult.ClassicResult.Columns[i]
				classicResult.Head[i].FieldName = col.Name
				classicResult.Head[i].Desc = col.Desc
			}

			// rows
			for i := range res.ExplainResult.ClassicResult.Rows {
				row := res.ExplainResult.ClassicResult.Rows[i]
				classicResult.Rows[i] = make(map[string]string, len(row))
				for k := range row {
					colName := res.ExplainResult.ClassicResult.Columns[k].Name
					classicResult.Rows[i][colName] = row[k]
				}
			}
			data.SQLExplain.ClassicResult = classicResult
		}
	}

	// table_metas
	{
		data.TableMetas = &TableMetas{}
		if res.TableMetaResultErr != nil {
			data.TableMetas.ErrMessage = res.TableMetaResultErr.Error()
		} else {
			tableMetaItemsData := make([]*TableMeta, len(res.TableMetaResult.TableMetas))
			for i := range res.TableMetaResult.TableMetas {
				tableMeta := res.TableMetaResult.TableMetas[i]
				tableMetaColumnsInfo := tableMeta.ColumnsInfo
				tableMetaIndexInfo := tableMeta.IndexesInfo
				tableMetaItemsData[i] = &TableMeta{}
				tableMetaItemsData[i].Columns = TableColumns{
					Rows: make([]map[string]string, len(tableMetaColumnsInfo.Rows)),
					Head: make([]TableMetaItemHeadResV1, len(tableMetaColumnsInfo.Columns)),
				}

				tableMetaItemsData[i].Indexes = TableIndexes{
					Rows: make([]map[string]string, len(tableMetaIndexInfo.Rows)),
					Head: make([]TableMetaItemHeadResV1, len(tableMetaIndexInfo.Columns)),
				}

				tableMetaColumnData := tableMetaItemsData[i].Columns
				for j := range tableMetaColumnsInfo.Columns {
					col := tableMetaColumnsInfo.Columns[j]
					tableMetaColumnData.Head[j].FieldName = col.Name
					tableMetaColumnData.Head[j].Desc = col.Desc
				}

				for j := range tableMetaColumnsInfo.Rows {
					tableMetaColumnData.Rows[j] = make(map[string]string, len(tableMetaColumnsInfo.Rows[j]))
					for k := range tableMetaColumnsInfo.Rows[j] {
						colName := tableMetaColumnsInfo.Columns[k].Name
						tableMetaColumnData.Rows[j][colName] = tableMetaColumnsInfo.Rows[j][k]
					}
				}

				tableMetaIndexData := tableMetaItemsData[i].Indexes
				for j := range tableMetaIndexInfo.Columns {
					tableMetaIndexData.Head[j].FieldName = tableMetaIndexInfo.Columns[j].Name
					tableMetaIndexData.Head[j].Desc = tableMetaIndexInfo.Columns[j].Desc
				}

				for j := range tableMetaIndexInfo.Rows {
					tableMetaIndexData.Rows[j] = make(map[string]string, len(tableMetaIndexInfo.Rows[j]))
					for k := range tableMetaIndexInfo.Rows[j] {
						colName := tableMetaIndexInfo.Columns[k].Name
						tableMetaIndexData.Rows[j][colName] = tableMetaIndexInfo.Rows[j][k]
					}
				}

				tableMetaItemsData[i].Name = tableMeta.Name
				tableMetaItemsData[i].Schema = tableMeta.Schema
				tableMetaItemsData[i].CreateTableSQL = tableMeta.CreateTableSQL
				tableMetaItemsData[i].Message = tableMeta.Message
			}
			data.TableMetas.Items = tableMetaItemsData
		}
	}

	// performance_statistics
	{
		data.PerformanceStatistics = &PerformanceStatistics{}

		// affect_rows
		data.PerformanceStatistics.AffectRows = &AffectRows{}
		if res.AffectRowsResultErr != nil {
			data.PerformanceStatistics.AffectRows.ErrMessage = res.AffectRowsResultErr.Error()
		} else {
			data.PerformanceStatistics.AffectRows.ErrMessage = res.AffectRowsResult.ErrMessage
			data.PerformanceStatistics.AffectRows.Count = int(res.AffectRowsResult.Count)
		}

	}

	return data
}
