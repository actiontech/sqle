package v1

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/locale"
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
	// This parameter is deprecated
	FilterBusiness         *string `query:"filter_business" json:"filter_business,omitempty"`
	FilterByEnvironmentTag *string `query:"filter_by_environment_tag" json:"filter_by_environment_tag,omitempty"`
	FilterPriority         *string `query:"filter_priority" json:"filter_priority,omitempty" enums:"high,low"`
	FuzzySearchEndpoint    *string `query:"fuzzy_search_endpoint" json:"fuzzy_search_endpoint,omitempty"`
	FuzzySearchSchemaName  *string `query:"fuzzy_search_schema_name" json:"fuzzy_search_schema_name,omitempty"`
	SortField              *string `query:"sort_field" json:"sort_field,omitempty" valid:"omitempty,oneof=first_appear_timestamp last_receive_timestamp fp_count" enums:"first_appear_timestamp,last_receive_timestamp,fp_count"`
	SortOrder              *string `query:"sort_order" json:"sort_order,omitempty" valid:"omitempty,oneof=asc desc" enums:"asc,desc"`
	PageIndex              uint32  `query:"page_index" valid:"required" json:"page_index"`
	PageSize               uint32  `query:"page_size" valid:"required" json:"page_size"`
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
	Status          string         `json:"status" enums:"unhandled,solved,ignored,manual_audited,sent"`
	Remark          string         `json:"remark"`
	Endpoint        string         `json:"endpoint"`
}

// AuditResult 用于SQL全生命周期展示的AuditResult
type AuditResult struct {
	Level           string `json:"level" example:"warn"`
	Message         string `json:"message" example:"避免使用不必要的内置函数md5()"`
	RuleName        string `json:"rule_name"`
	ErrorInfo       string `json:"error_info"`
	ExecutionFailed bool   `json:"execution_failed"`
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
	return nil
}

type BatchUpdateSqlManageReq struct {
	SqlManageIdList []*uint64 `json:"sql_manage_id_list"`
	Status          *string   `json:"status" enums:"solved,ignored,manual_audited"`
	Priority        *string   `json:"priority" enums:",high"`
	Assignees       []string  `json:"assignees"`
	Remark          *string   `json:"remark"`
}

type SqlManageCodingReq struct {
	SqlManageIdList   []*uint64       `json:"sql_manage_id_list"`
	Priority          *CodingPriority `json:"priority" enums:"LOW,MEDIUM,HIGH,EMERGENCY"`
	CodingProjectName *string         `json:"coding_project_name"`
	Type              *CodingType     `json:"type" enums:"DEFECT,MISSION,REQUIREMENT,EPIC,SUB_TASK"`
}

type CodingType string

const (
	CodingTypeDefect      CodingType = "DEFECT"
	CodingTypeMission     CodingType = "MISSION"
	CodingTypeRequirement CodingType = "REQUIREMENT"
	CodingTypeEpic        CodingType = "EPIC"
	CodingTypeSubTask     CodingType = "SUB_TASK"
)

type CodingPriority string

const (
	CodingPriorityLow       CodingPriority = "LOW"
	CodingPriorityMedium                   = "MEDIUM"
	CodingPriorityHigh                     = "HIGH"
	CodingPriorityEmergency                = "EMERGENCY"
)

func (codingPriority CodingPriority) Weight() string {
	weight := "-1"
	switch codingPriority {
	case CodingPriorityLow:
		weight = "0"
	case CodingPriorityMedium:
		weight = "1"
	case CodingPriorityHigh:
		weight = "2"
	case CodingPriorityEmergency:
		weight = "3"
	default:
		weight = "-1"
	}
	return weight
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
	FuzzySearchSqlFingerprint *string `query:"fuzzy_search_sql_fingerprint" json:"fuzzy_search_sql_fingerprint,omitempty"`
	FilterAssignee            *string `query:"filter_assignee" json:"filter_assignee,omitempty"`
	// This parameter is deprecated
	FilterBusiness               *string `query:"filter_business" json:"filter_business,omitempty"`
	FilterByEnvironmentTag       *string `query:"filter_by_environment_tag" json:"filter_by_environment_tag,omitempty"`
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

// @Deprecated
// ExportSqlManagesV1
// @Summary 导出SQL管控
// @Description export sql manage
// @Id exportSqlManageV1
// @Tags SqlManage
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param fuzzy_search_sql_fingerprint query string false "fuzzy search sql fingerprint"
// @Param filter_assignee query string false "assignee"
// @Param filter_business query string false "filter by business" // This parameter is deprecated
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
	return nil
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

type SqlAnalysisChart struct {
	Points  *[]ChartPoint `json:"points"`
	XInfo   *string       `json:"x_info"`
	YInfo   *string       `json:"y_info"`
	Message string        `json:"message"`
}

type ChartPoint struct {
	X     *string             `json:"x"`
	Y     *float64            `json:"y"`
	Infos []map[string]string `json:"info"`
}

type GetSqlManageSqlAnalysisResp struct {
	controller.BaseRes
	// V1版本不能引用V2版本的结构体,所以只能复制一份
	Data *SqlAnalysis `json:"data"`
}

type PostSqlManageCodingResp struct {
	controller.BaseRes
	Data *CodingResp `json:"data"`
}

type CodingResp struct {
	Message string `json:"message"`
	Code    string `json:"code"`
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

type SqlManageAnalysisChartReq struct {
	LatestPointEnabled bool    `query:"latest_point_enabled" json:"latest_point_enabled"`
	StartTime          *string `query:"start_time" json:"start_time"`
	EndTime            *string `query:"end_time" json:"end_time"`
	MetricName         *string `query:"metric_name" json:"metric_name"`
}

type SqlManageAnalysisChartResp struct {
	controller.BaseRes
	Data *SqlAnalysisChart `json:"data"`
}

// GetSqlManageSqlAnalysisChartV1
// @Summary 获取SQL管控SQL执行计划Cost趋势图表
// @Description get sql manage analysis
// @Id GetSqlManageSqlAnalysisChartV1
// @Tags SqlManage
// @Param project_name path string true "project name"
// @Param sql_manage_id path string true "sql manage id"
// @Param latest_point_enabled query bool true "latest_point_enabled"
// @Param start_time query string true "start time"
// @Param end_time query string true "end time"
// @Param metric_name query string true "metric_name"
// @Security ApiKeyAuth
// @Success 200 {object} SqlManageAnalysisChartResp
// @Router /v1/projects/{project_name}/sql_manages/{sql_manage_id}/sql_analysis_chart [get]
func GetSqlManageSqlAnalysisChartV1(c echo.Context) error {
	return getSqlManageSqlAnalysisChartV1(c)
}

// SendSqlManage
// @Summary 推送SQL管控结果到外部系统
// @Description get sql manage analysis
// @Id SendSqlManage
// @Tags SqlManage
// @Param project_name path string true "project name"
// @Param SqlManageCodingReq body SqlManageCodingReq true "batch update sql manage request"
// @Security ApiKeyAuth
// @Success 200 {object} PostSqlManageCodingResp
// @Router /v1/projects/{project_name}/sql_manages/send [post]
func SendSqlManage(c echo.Context) error {
	return sendSqlManage(c)
}

func convertSQLAnalysisResultToRes(ctx context.Context, res *AnalysisResult, rawSQL string) *SqlAnalysis {
	data := &SqlAnalysis{}

	// explain
	{
		data.SQLExplain = &SQLExplain{SQL: rawSQL}
		data.SQLExplain.Cost = *res.Cost
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
				classicResult.Head[i].Desc = col.I18nDesc.GetStrInLang(locale.Bundle.GetLangTagFromCtx(ctx))
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
					tableMetaColumnData.Head[j].Desc = col.I18nDesc.GetStrInLang(locale.Bundle.GetLangTagFromCtx(ctx))
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
					tableMetaIndexData.Head[j].Desc = tableMetaIndexInfo.Columns[j].I18nDesc.GetStrInLang(locale.Bundle.GetLangTagFromCtx(ctx))
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

type GetGlobalSqlManageListReq struct {
	FilterProjectUid      *string                `query:"filter_project_uid" json:"filter_project_uid,omitempty"`
	FilterInstanceId      *string                `query:"filter_instance_id" json:"filter_instance_id,omitempty"`
	FilterProjectPriority *dmsV1.ProjectPriority `query:"filter_project_priority" json:"filter_project_priority,omitempty" enums:"high,medium,low"`
	PageIndex             uint32                 `query:"page_index" valid:"required" json:"page_index"`
	PageSize              uint32                 `query:"page_size" valid:"required" json:"page_size"`
}

type GetGlobalSqlManageListResp struct {
	controller.BaseRes
	Data      []*GlobalSqlManage `json:"data"`
	TotalNums uint64             `json:"total_nums"`
}

type GlobalSqlManage struct {
	Id                   uint64                `json:"id"`
	Sql                  string                `json:"sql"`
	Source               *Source               `json:"source"`
	AuditResult          []*AuditResult        `json:"audit_result"`
	ProjectName          string                `json:"project_name"`
	ProjectUid           string                `json:"project_uid"`
	InstanceName         string                `json:"instance_name"`
	InstanceId           string                `json:"instance_id"`
	Status               string                `json:"status" enums:"unhandled,solved,ignored,manual_audited,sent"`
	ProjectPriority      dmsV1.ProjectPriority `json:"project_priority" enums:"high,medium,low"`
	FirstAppearTimeStamp string                `json:"first_appear_timestamp"`
	ProblemDescriptions  []string              `json:"problem_descriptions"` // 根据来源信息拼接
}

// GetGlobalSqlManageList
// @Summary 获取全局管控sql列表
// @Description get global sql manage list
// @Tags SqlManage
// @Id GetGlobalSqlManageList
// @Security ApiKeyAuth
// @Param filter_project_uid query string false "project uid"
// @Param filter_instance_id query string false "instance id"
// @Param filter_project_priority query string false "project priority" Enums(high,medium,low)
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Success 200 {object} v1.GetGlobalSqlManageListResp
// @Router /v1/dashboard/sql_manages [get]
func GetGlobalSqlManageList(c echo.Context) error {
	return getGlobalSqlManageList(c)
}

type GetGlobalSqlManageStatisticsReq struct {
	FilterProjectUid      *string                `query:"filter_project_uid" json:"filter_project_uid,omitempty"`
	FilterInstanceId      *string                `query:"filter_instance_id" json:"filter_instance_id,omitempty"`
	FilterProjectPriority *dmsV1.ProjectPriority `query:"filter_project_priority" json:"filter_project_priority,omitempty" enums:"high,medium,low"`
}

type GetGlobalSqlManageStatisticsResp struct {
	controller.BaseRes
	TotalNums uint64 `json:"total_nums"`
}

// GetGlobalSqlManageStatistics
// @Summary 获取全局管控sql统计信息
// @Description get global sql manage statistics
// @Tags SqlManage
// @Id GetGlobalSqlManageStatistics
// @Security ApiKeyAuth
// @Param filter_project_uid query string false "project uid"
// @Param filter_instance_id query string false "instance id"
// @Param filter_project_priority query string false "project priority" Enums(high,medium,low)
// @Success 200 {object} v1.GetGlobalSqlManageStatisticsResp
// @Router /v1/dashboard/sql_manages/statistics  [get]
func GetGlobalSqlManageStatistics(c echo.Context) error {
	return getGlobalSqlManageStatistics(c)
}

type GetAbnormalAuditPlanInstancesResp struct {
	controller.BaseRes
	Data []*AbnormalAuditPlanInstance `json:"data"`
}

type AbnormalAuditPlanInstance struct {
	InstanceName        string `json:"instance_name" example:"MySQL"`
	InstanceAuditPlanID uint   `json:"instance_audit_plan_id"`
}

// GetAbnormalInstanceAuditPlans get the instance of audit plan execution abnormal
// @Summary 获取执行异常的扫描任务实例
// @Description get the instance of audit plan execution abnormal
// @Id getAbnormalInstanceAuditPlansV1
// @Tags SqlManage
// @Param project_name path string true "project name"
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetAbnormalAuditPlanInstancesResp
// @Router /v1/projects/{project_name}/sql_manages/abnormal_audit_plan_instance [get]
func GetAbnormalInstanceAuditPlans(c echo.Context) error {
	return getAbnormalInstanceAuditPlans(c)
}

type GetSqlManageSqlPerformanceInsightsResp struct {
	controller.BaseRes
	Data *SqlManageSqlPerformanceInsights `json:"data"`
}

type SqlManageSqlPerformanceInsights struct {
	XInfo   *string `json:"x_info"`
	YInfo   *string `json:"y_info"`
	Message string  `json:"message"`
	Lines   *[]Line `json:"lines"`
}

type Line struct {
	LineName string        `json:"line_name"`
	Points   *[]ChartPoint `json:"points"`
}

// 定义SQL性能洞察指标类型
type MetricName string

const (
	MetricNameComprehensiveTrend MetricName = "comprehensive_trend"  // 数据源综合趋势
	MetricNameSlowSQLTrend       MetricName = "slow_sql_trend"       // 慢SQL趋势
	MetricNameTopSQLTrend        MetricName = "top_sql_trend"        // TopSQL趋势
	MetricNameActiveSessionTrend MetricName = "active_session_trend" // 活跃会话数趋势
)

// GetSqlManageSqlPerformanceInsights
// @Summary 获取SQL管控SQL性能洞察图表数据
// @Description get sql manage sql performance insights
// @Id GetSqlManageSqlPerformanceInsights
// @Tags SqlManage
// @Param project_name path string true "project name"
// @Param metric_name query string true "metric name" Enums(comprehensive_trend,slow_sql_trend,top_sql_trend,active_session_trend)
// @Param start_time query string true "start time"
// @Param end_time query string true "end time"
// @Param instance_name query string true "instance name"
// @Security ApiKeyAuth
// @Success 200 {object} GetSqlManageSqlPerformanceInsightsResp
// @Router /v1/projects/{project_name}/sql_performance_insights [get]
func GetSqlManageSqlPerformanceInsights(c echo.Context) error {
	// 获取指标类型参数
	metricNameStr := c.QueryParam("metric_name")
	metricName := MetricName(metricNameStr)

	// 构建时间点数据 - 所有图表使用相同的时间范围
	timePoints := []string{
		"2025-05-07T00:00:00+08:00", "2025-05-07T03:00:00+08:00", "2025-05-07T06:00:00+08:00", "2025-05-07T09:00:00+08:00",
		"2025-05-07T12:00:00+08:00", "2025-05-07T15:00:00+08:00", "2025-05-07T18:00:00+08:00", "2025-05-07T21:00:00+08:00",
		"2025-05-08T00:00:00+08:00",
	}

	var lines []Line
	var xInfo, yInfo string

	switch metricName {
	case MetricNameComprehensiveTrend:
		// 数据源综合性能趋势 - 四条线
		xInfo = "时间"
		yInfo = "指标值"

		// CPU使用率数据
		cpuPoints := make([]ChartPoint, len(timePoints))
		cpuValues := []float64{50.0, 35.0, 70.0, 45.0, 30.0, 45.0, 65.0, 90.0, 75.0}
		for i, t := range timePoints {
			x := t
			y := cpuValues[i]
			cpuPoints[i] = ChartPoint{
				X: &x,
				Y: &y,
				Infos: []map[string]string{
					{"值": fmt.Sprintf("%.1f%%", y)},
				},
			}
		}

		// 磁盘I/O数据
		diskIOPoints := make([]ChartPoint, len(timePoints))
		diskIOValues := []float64{40.0, 45.0, 35.0, 50.0, 60.0, 45.0, 55.0, 65.0, 55.0}
		for i, t := range timePoints {
			x := t
			y := diskIOValues[i]
			diskIOPoints[i] = ChartPoint{
				X: &x,
				Y: &y,
				Infos: []map[string]string{
					{"值": fmt.Sprintf("%.1f MB/s", y)},
				},
			}
		}

		// 连接数数据
		connectionPoints := make([]ChartPoint, len(timePoints))
		connectionValues := []float64{60.0, 50.0, 45.0, 40.0, 50.0, 45.0, 40.0, 35.0, 45.0}
		for i, t := range timePoints {
			x := t
			y := connectionValues[i]
			connectionPoints[i] = ChartPoint{
				X: &x,
				Y: &y,
				Infos: []map[string]string{
					{"值": fmt.Sprintf("%.0f", y)},
				},
			}
		}

		// 网络流量数据
		networkPoints := make([]ChartPoint, len(timePoints))
		networkValues := []float64{55.0, 65.0, 60.0, 70.0, 75.0, 65.0, 80.0, 70.0, 85.0}
		for i, t := range timePoints {
			x := t
			y := networkValues[i]
			networkPoints[i] = ChartPoint{
				X: &x,
				Y: &y,
				Infos: []map[string]string{
					{"值": fmt.Sprintf("%.1f MB/s", y)},
				},
			}
		}

		// 构建四条线
		lines = []Line{
			{
				LineName: "CPU使用率",
				Points:   &cpuPoints,
			},
			{
				LineName: "磁盘I/O",
				Points:   &diskIOPoints,
			},
			{
				LineName: "连接数",
				Points:   &connectionPoints,
			},
			{
				LineName: "网络流量",
				Points:   &networkPoints,
			},
		}

	case MetricNameSlowSQLTrend:
		// 慢SQL趋势 - 一条线
		xInfo = "时间"
		yInfo = "慢SQL数量"

		// 慢SQL数量数据
		slowSQLPoints := make([]ChartPoint, len(timePoints))
		slowSQLValues := []float64{5.0, 8.0, 6.0, 10.0, 15.0, 10.0, 20.0, 50.0, 35.0}
		for i, t := range timePoints {
			x := t
			y := slowSQLValues[i]
			slowSQLPoints[i] = ChartPoint{
				X: &x,
				Y: &y,
				Infos: []map[string]string{
					{"数量": fmt.Sprintf("%.0f", y)},
				},
			}
		}

		lines = []Line{
			{
				LineName: "慢SQL数量",
				Points:   &slowSQLPoints,
			},
		}

	case MetricNameTopSQLTrend:
		// TopSQL执行趋势 - 两条线
		xInfo = "时间"
		yInfo = "执行指标"

		// 执行次数数据
		execCountPoints := make([]ChartPoint, len(timePoints))
		execCountValues := []float64{30.0, 35.0, 32.0, 40.0, 45.0, 55.0, 65.0, 75.0, 70.0}
		for i, t := range timePoints {
			x := t
			y := execCountValues[i]
			execCountPoints[i] = ChartPoint{
				X: &x,
				Y: &y,
				Infos: []map[string]string{
					{"执行次数": fmt.Sprintf("%.0f", y)},
				},
			}
		}

		// 执行时间数据
		execTimePoints := make([]ChartPoint, len(timePoints))
		execTimeValues := []float64{25.0, 30.0, 25.0, 35.0, 45.0, 60.0, 70.0, 85.0, 75.0}
		for i, t := range timePoints {
			x := t
			y := execTimeValues[i]
			execTimePoints[i] = ChartPoint{
				X: &x,
				Y: &y,
				Infos: []map[string]string{
					{"执行时间": fmt.Sprintf("%.0f ms", y)},
				},
			}
		}

		lines = []Line{
			{
				LineName: "执行次数",
				Points:   &execCountPoints,
			},
			{
				LineName: "执行时间",
				Points:   &execTimePoints,
			},
		}

	case MetricNameActiveSessionTrend:
		// 活跃会话数趋势 - 一条线
		xInfo = "时间"
		yInfo = "会话数"

		// 活跃会话数据
		sessionPoints := make([]ChartPoint, len(timePoints))
		sessionValues := []float64{40.0, 50.0, 40.0, 50.0, 45.0, 60.0, 80.0, 150.0, 120.0}
		for i, t := range timePoints {
			x := t
			y := sessionValues[i]
			sessionPoints[i] = ChartPoint{
				X: &x,
				Y: &y,
				Infos: []map[string]string{
					{"会话数": fmt.Sprintf("%.0f", y)},
				},
			}
		}

		lines = []Line{
			{
				LineName: "活跃会话数",
				Points:   &sessionPoints,
			},
		}

	default:
		// 默认返回空数据
		xInfo = "时间"
		yInfo = "指标值"
		lines = []Line{}
	}

	// 返回结果
	return c.JSON(http.StatusOK, GetSqlManageSqlPerformanceInsightsResp{
		BaseRes: controller.BaseRes{
			Code:    0,
			Message: "success",
		},
		Data: &SqlManageSqlPerformanceInsights{
			XInfo:   &xInfo,
			YInfo:   &yInfo,
			Message: "",
			Lines:   &lines,
		},
	})
}

type GetSqlManageSqlPerformanceInsightsRelatedSQLResp struct {
	controller.BaseRes
	Data      []*RelatedSQLInfo `json:"data"`
	TotalNums uint32            `json:"total_nums"`
}

type SqlSourceTypeEnum string

const (
	SqlSourceTypeOrder     SqlSourceTypeEnum = "order"
	SqlSourceTypeSqlManage SqlSourceTypeEnum = "sql_manage"
)

type RelatedSQLInfo struct {
	SqlFingerprint     string                   `json:"sql_fingerprint"`
	Source             SqlSourceTypeEnum        `json:"source" enums:"order,sql_manage"`
	ExecuteStartTime   string                   `json:"execute_start_time"`
	ExecuteEndTime     string                   `json:"execute_end_time"`
	ExecuteTime        float64                  `json:"execute_time"`                   // 执行时间(s)
	LockWaitTime       float64                  `json:"lock_wait_time"`                 // 锁等待时间(s)
	ExecutionCostTrend *SqlAnalysisScatterChart `json:"execution_cost_trend,omitempty"` // SQL执行代价趋势图表
}

// 散点图结构体，专用于SQL执行代价的散点图表示
type SqlAnalysisScatterChart struct {
	Points  *[]ScatterPoint `json:"points"`
	XInfo   *string         `json:"x_info"`
	YInfo   *string         `json:"y_info"`
	Message string          `json:"message"`
}

type ScatterPoint struct {
	Time            *string             `json:"time"`
	Cost            *float64            `json:"cost"`
	SQL             *string             `json:"sql"`
	Id              uint64              `json:"id"`
	IsInTransaction bool                `json:"is_in_transaction"`
	Infos           []map[string]string `json:"info"`
}

type GetSqlManageSqlPerformanceInsightsRelatedSQLReq struct {
	InstanceName string             `query:"instance_name" json:"instance_name" valid:"required"`
	StartTime    string             `query:"start_time" json:"start_time" valid:"required"`
	EndTime      string             `query:"end_time" json:"end_time" valid:"required"`
	FilterSource *SqlSourceTypeEnum `query:"filter_source" json:"filter_source,omitempty" enums:"order,sql_manage"`
	OrderBy      *string            `query:"order_by" json:"order_by,omitempty"`
	IsAsc        *bool              `query:"is_asc" json:"is_asc,omitempty"`
	PageIndex    uint32             `query:"page_index" valid:"required" json:"page_index"`
	PageSize     uint32             `query:"page_size" valid:"required" json:"page_size"`
}

// fixme: 这个接口的设计上由于产品对于SQL指纹浮动的图表的展示情况还有疑问。所以对应的设计应该是缺失了部分。
// 1、由于后续还需要查询具体sql的管理事务等。所以可能需要设计一个 id 之类的字段以便后续查询。但是现在list里的内容是sql指纹，不是具体的sql，所以这个ID要放在哪还需要等产品的结论。
// 2、同样的，对于关联事物功能。目前产品的设计里，关联事物打开的是具体的SQL，而不是SQL指纹。但是又不是所有SQL都在事务里。
// 2.1、所以问题一为，点击按钮的时候要打开的是哪一条具体的SQL。
// 2.2、需要设计字段显示的告诉前端这条SQL是否在一个事务当中。以便前端展示 "关联事务" 按钮是否可用。
// 3、这个table还有一个跳转到SQL分析的功能。但是现有的SQL分析也是针对单条SQL。而这个地方也是SQL指纹。
// 3.1、而且看上去现有的SQL是基于生成了分析结果。然后这个结果有个类似于 sql_manage_id 的玩意来获取的。但是从这里跳过去就没这个id了。需要前后端讨论一下具体实现。
// GetSqlManageSqlPerformanceInsightsRelatedSQL
// @Summary 获取sql洞察 时间选区 的关联SQL
// @Description Get related SQL for the selected time range in SQL performance insights
// @Id GetSqlManageSqlPerformanceInsightsRelatedSQL
// @Tags SqlManage
// @Param project_name path string true "project name"
// @Param instance_name query string true "instance name"
// @Param start_time query string true "start time"
// @Param end_time query string true "end time"
// @Param filter_source query string false "filter by SQL source" Enums(order,sql_manage)
// @Param order_by query string false "order by field"
// @Param is_asc query bool false "is ascending order"
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Security ApiKeyAuth
// @Success 200 {object} GetSqlManageSqlPerformanceInsightsRelatedSQLResp
// @Router /v1/projects/{project_name}/sql_performance_insights/related_sql [get]
func GetSqlManageSqlPerformanceInsightsRelatedSQL(c echo.Context) error {
	// 解析请求参数
	req := &GetSqlManageSqlPerformanceInsightsRelatedSQLReq{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, controller.BaseRes{
			Code:    1,
			Message: err.Error(),
		})
	}

	// 创建SQL执行代价趋势示例数据
	createExecutionCostTrend := func(startValue, endValue float64) *SqlAnalysisScatterChart {
		timePoints := []string{
			"2023-05-01T00:00:00+08:00", "2023-05-01T00:00:00+08:00", "2023-05-02T00:00:00+08:00", "2023-05-03T00:00:00+08:00",
			"2023-05-04T00:00:00+08:00", "2023-05-05T00:00:00+08:00", "2023-05-05T00:00:00+08:00", "2023-05-06T00:00:00+08:00",
			"2023-05-07T00:00:00+08:00",
		}

		// 线性变化的代价值
		step := (endValue - startValue) / float64(len(timePoints)-1)
		points := make([]ScatterPoint, len(timePoints))

		sql1 := "SELECT * FROM users WHERE user_id = 100"
		sql2 := "SELECT * FROM orders WHERE order_id = 5001"

		for i, t := range timePoints {
			x := t
			y := startValue + step*float64(i)
			sql := sql1
			if i%2 == 0 {
				sql = sql2
			}
			points[i] = ScatterPoint{
				Time:            &x,
				Cost:            &y,
				Id:              uint64(i),
				IsInTransaction: i%2 == 0,
				SQL:             &sql,
				Infos: []map[string]string{
					{"代价": fmt.Sprintf("%.2f", y)},
				},
			}
		}

		xInfo := "时间"
		yInfo := "执行代价"

		return &SqlAnalysisScatterChart{
			Points:  &points,
			XInfo:   &xInfo,
			YInfo:   &yInfo,
			Message: "",
		}
	}

	// 扩展示例数据到10条
	mockData := []*RelatedSQLInfo{
		{
			SqlFingerprint:     "SELECT * FROM users WHERE user_id = 100",
			Source:             SqlSourceTypeSqlManage,
			ExecuteStartTime:   "2023-05-07T10:15:00+08:00",
			ExecuteEndTime:     "2023-05-07T10:15:01+08:00",
			ExecuteTime:        1.2,
			LockWaitTime:       0.1,
			ExecutionCostTrend: createExecutionCostTrend(10.5, 25.8),
		},
		{
			SqlFingerprint:     "UPDATE orders SET status = 'completed' WHERE order_id = 5001",
			Source:             SqlSourceTypeOrder,
			ExecuteStartTime:   "2023-05-07T10:20:00+08:00",
			ExecuteEndTime:     "2023-05-07T10:20:02+08:00",
			ExecuteTime:        2.0,
			LockWaitTime:       0.5,
			ExecutionCostTrend: createExecutionCostTrend(15.3, 12.1),
		},
		{
			SqlFingerprint:   "SELECT * FROM products WHERE category_id = 3 ORDER BY price DESC LIMIT 10",
			Source:           SqlSourceTypeSqlManage,
			ExecuteStartTime: "2023-05-07T10:25:00+08:00",
			ExecuteEndTime:   "2023-05-07T10:25:00+08:00",
			ExecuteTime:      0.5,
			LockWaitTime:     0.0,
			// 不是所有记录都有趋势数据
		},
		{
			SqlFingerprint:     "INSERT INTO user_logs (user_id, action, timestamp) VALUES (102, 'login', NOW())",
			Source:             SqlSourceTypeOrder,
			ExecuteStartTime:   "2023-05-07T10:30:00+08:00",
			ExecuteEndTime:     "2023-05-07T10:30:01+08:00",
			ExecuteTime:        0.8,
			LockWaitTime:       0.2,
			ExecutionCostTrend: createExecutionCostTrend(5.2, 18.9),
		},
		{
			SqlFingerprint:     "SELECT u.name, o.order_date FROM users u JOIN orders o ON u.id = o.user_id WHERE o.status = 'pending'",
			Source:             SqlSourceTypeSqlManage,
			ExecuteStartTime:   "2023-05-07T10:35:00+08:00",
			ExecuteEndTime:     "2023-05-07T10:35:03+08:00",
			ExecuteTime:        3.2,
			LockWaitTime:       0.3,
			ExecutionCostTrend: createExecutionCostTrend(30.5, 45.2),
		},
		{
			SqlFingerprint:     "DELETE FROM cart_items WHERE created_at < DATE_SUB(NOW(), INTERVAL 30 DAY)",
			Source:             SqlSourceTypeOrder,
			ExecuteStartTime:   "2023-05-07T10:40:00+08:00",
			ExecuteEndTime:     "2023-05-07T10:40:05+08:00",
			ExecuteTime:        4.5,
			LockWaitTime:       1.2,
			ExecutionCostTrend: createExecutionCostTrend(22.7, 18.1),
		},
		{
			SqlFingerprint:   "UPDATE inventory SET quantity = quantity - 1 WHERE product_id = 555",
			Source:           SqlSourceTypeSqlManage,
			ExecuteStartTime: "2023-05-07T10:45:00+08:00",
			ExecuteEndTime:   "2023-05-07T10:45:01+08:00",
			ExecuteTime:      0.9,
			LockWaitTime:     0.4,
			// 不是所有记录都有趋势数据
		},
		{
			SqlFingerprint:     "SELECT COUNT(*) FROM customer_support_tickets WHERE status != 'closed' GROUP BY priority",
			Source:             SqlSourceTypeOrder,
			ExecuteStartTime:   "2023-05-07T10:50:00+08:00",
			ExecuteEndTime:     "2023-05-07T10:50:02+08:00",
			ExecuteTime:        1.8,
			LockWaitTime:       0.1,
			ExecutionCostTrend: createExecutionCostTrend(8.3, 27.5),
		},
		{
			SqlFingerprint:     "SELECT p.name, SUM(oi.quantity) FROM products p JOIN order_items oi ON p.id = oi.product_id GROUP BY p.id ORDER BY SUM(oi.quantity) DESC LIMIT 5",
			Source:             SqlSourceTypeSqlManage,
			ExecuteStartTime:   "2023-05-07T10:55:00+08:00",
			ExecuteEndTime:     "2023-05-07T10:55:04+08:00",
			ExecuteTime:        3.8,
			LockWaitTime:       0.2,
			ExecutionCostTrend: createExecutionCostTrend(40.1, 35.6),
		},
		{
			SqlFingerprint:     "ALTER TABLE user_sessions ADD COLUMN last_activity TIMESTAMP",
			Source:             SqlSourceTypeOrder,
			ExecuteStartTime:   "2023-05-07T11:00:00+08:00",
			ExecuteEndTime:     "2023-05-07T11:00:10+08:00",
			ExecuteTime:        10.0,
			LockWaitTime:       1.5,
			ExecutionCostTrend: createExecutionCostTrend(18.0, 60.0),
		},
	}

	// 应用筛选条件
	var filteredData []*RelatedSQLInfo
	if req.FilterSource != nil {
		for _, item := range mockData {
			if item.Source == *req.FilterSource {
				filteredData = append(filteredData, item)
			}
		}
	} else {
		filteredData = mockData
	}

	// 应用排序
	if req.OrderBy != nil {
		// 实现排序逻辑
		sort.Slice(filteredData, func(i, j int) bool {
			isDescending := req.IsAsc != nil && !*req.IsAsc

			// 根据OrderBy字段进行不同的排序
			switch *req.OrderBy {
			case "execute_start_time":
				if isDescending {
					return filteredData[i].ExecuteStartTime > filteredData[j].ExecuteStartTime
				}
				return filteredData[i].ExecuteStartTime < filteredData[j].ExecuteStartTime
			case "execute_time":
				if isDescending {
					return filteredData[i].ExecuteTime > filteredData[j].ExecuteTime
				}
				return filteredData[i].ExecuteTime < filteredData[j].ExecuteTime
			case "lock_wait_time":
				if isDescending {
					return filteredData[i].LockWaitTime > filteredData[j].LockWaitTime
				}
				return filteredData[i].LockWaitTime < filteredData[j].LockWaitTime
			default:
				// 默认按ExecuteStartTime排序
				if isDescending {
					return filteredData[i].ExecuteStartTime > filteredData[j].ExecuteStartTime
				}
				return filteredData[i].ExecuteStartTime < filteredData[j].ExecuteStartTime
			}
		})
	}

	// 应用分页
	totalNums := uint32(33)

	// 计算分页范围
	startIndex := (req.PageIndex - 1) * req.PageSize
	endIndex := startIndex + req.PageSize

	// 边界检查
	if startIndex >= uint32(len(filteredData)) {
		startIndex = 0
		endIndex = 0
	}

	if endIndex > uint32(len(filteredData)) {
		endIndex = uint32(len(filteredData))
	}

	// 获取当前页的数据
	var pageData []*RelatedSQLInfo
	if endIndex > startIndex {
		pageData = filteredData[startIndex:endIndex]
	} else {
		pageData = []*RelatedSQLInfo{}
	}

	return c.JSON(http.StatusOK, GetSqlManageSqlPerformanceInsightsRelatedSQLResp{
		BaseRes: controller.BaseRes{
			Code:    0,
			Message: "success",
		},
		Data:      pageData,
		TotalNums: totalNums,
	})
}

type LockType string

const (
	LockTypeShared                   LockType = "SHARED"                     // 共享锁
	LockTypeExclusive                LockType = "EXCLUSIVE"                  // 排他锁
	LockTypeIntentionShared          LockType = "INTENTION_SHARED"           // 意向共享锁
	LockTypeIntentionExclusive       LockType = "INTENTION_EXCLUSIVE"        // 意向排他锁
	LockTypeSharedIntentionExclusive LockType = "SHARED_INTENTION_EXCLUSIVE" // 共享意向排他锁
	LockTypeRowLock                  LockType = "ROW_LOCK"                   // 行锁
	LockTypeTableLock                LockType = "TABLE_LOCK"                 // 表锁
	LockTypeMetadataLock             LockType = "METADATA_LOCK"              // 元数据锁
)

type TransactionState string

const (
	TransactionStateRunning   TransactionState = "RUNNING"   // 执行中
	TransactionStateCompleted TransactionState = "COMPLETED" // 已完成
)

type TransactionInfo struct {
	TransactionId        string           `json:"transaction_id"`
	LockType             LockType         `json:"lock_type" enums:"SHARED,EXCLUSIVE,INTENTION_SHARED,INTENTION_EXCLUSIVE,SHARED_INTENTION_EXCLUSIVE,ROW_LOCK,TABLE_LOCK,METADATA_LOCK"`
	TransactionStartTime string           `json:"transaction_start_time"`
	TransactionEndTime   string           `json:"transaction_end_time"`
	TransactionDuration  float64          `json:"transaction_duration"`
	TransactionState     TransactionState `json:"transaction_state" enums:"RUNNING,COMPLETED"`
}

type TransactionTimelineItem struct {
	StartTime   string `json:"start_time"`
	Description string `json:"description"`
}

type TransactionTimeline struct {
	Timeline         []*TransactionTimelineItem `json:"timeline"`
	CurrentStepIndex int                        `json:"current_step_index"`
}

// fixme： 这里同样有SQL分析功能。所以有跟上面3.1相同的问题。
type TransactionSQL struct {
	SQL             string   `json:"sql"`
	ExecuteDuration float64  `json:"execute_duration"`
	LockType        LockType `json:"lock_type" enums:"SHARED,EXCLUSIVE,INTENTION_SHARED,INTENTION_EXCLUSIVE,SHARED_INTENTION_EXCLUSIVE,ROW_LOCK,TABLE_LOCK,METADATA_LOCK"`
}

type TransactionLockInfo struct {
	LockType      LockType `json:"lock_type" enums:"SHARED,EXCLUSIVE,INTENTION_SHARED,INTENTION_EXCLUSIVE,SHARED_INTENTION_EXCLUSIVE,ROW_LOCK,TABLE_LOCK,METADATA_LOCK"`
	TableName     string   `json:"table_name"`
	CreateLockSQL string   `json:"create_lock_sql"`
}

type RelatedTransactionInfo struct {
	TransactionInfo     *TransactionInfo       `json:"transaction_info"`
	TransactionTimeline *TransactionTimeline   `json:"transaction_timeline"`
	TransactionSQLList  []*TransactionSQL      `json:"related_sql_info"`
	TransactionLockInfo []*TransactionLockInfo `json:"transaction_lock_info"`
}

type GetSqlManageSqlPerformanceInsightsRelatedTransactionResp struct {
	controller.BaseRes
	Data *RelatedTransactionInfo `json:"data"`
}

// GetSqlManageSqlPerformanceInsightsRelatedTransaction
// @Summary 获取sql洞察 相关SQL中具体一条SQL 的关联事务
// @Description Get related transaction for the selected SQL in SQL performance insights
// @Id GetSqlManageSqlPerformanceInsightsRelatedTransaction
// @Tags SqlManage
// @Param project_name path string true "project name"
// @Param instance_name query string true "instance name"
// @Param sql_id query string true "sql id"
// @Security ApiKeyAuth
// @Success 200 {object} GetSqlManageSqlPerformanceInsightsRelatedTransactionResp
// @Router /v1/projects/{project_name}/sql_performance_insights/related_sql/related_transaction [get]
func GetSqlManageSqlPerformanceInsightsRelatedTransaction(c echo.Context) error {
	// 创建模拟的事务信息
	transactionInfo := &TransactionInfo{
		TransactionId:        "TRX-123456789",
		LockType:             LockTypeExclusive,
		TransactionStartTime: "2023-05-07T10:45:00+08:00",
		TransactionEndTime:   "2023-05-07T10:45:30+08:00",
		TransactionDuration:  30.0,
		TransactionState:     TransactionStateCompleted,
	}

	// 创建模拟的事务时间线
	transactionTimeline := &TransactionTimeline{
		Timeline: []*TransactionTimelineItem{
			{
				StartTime:   "2023-05-07T10:45:00+08:00",
				Description: "事务开始",
			},
			{
				StartTime:   "2023-05-07T10:45:05+08:00",
				Description: "执行第一条SQL语句",
			},
			{
				StartTime:   "2023-05-07T10:45:10+08:00",
				Description: "获取表锁",
			},
			{
				StartTime:   "2023-05-07T10:45:15+08:00",
				Description: "执行第二条SQL语句",
			},
			{
				StartTime:   "2023-05-07T10:45:25+08:00",
				Description: "提交事务",
			},
			{
				StartTime:   "2023-05-07T10:45:30+08:00",
				Description: "事务完成",
			},
		},
		CurrentStepIndex: 5,
	}

	// 创建模拟的SQL列表
	transactionSQLList := []*TransactionSQL{
		{
			SQL:             "BEGIN TRANSACTION",
			ExecuteDuration: 0.1,
			LockType:        LockTypeIntentionExclusive,
		},
		{
			SQL:             "UPDATE inventory SET quantity = quantity - 1 WHERE product_id = 555",
			ExecuteDuration: 5.2,
			LockType:        LockTypeRowLock,
		},
		{
			SQL:             "UPDATE order_items SET status = 'shipped' WHERE order_id = 12345",
			ExecuteDuration: 10.3,
			LockType:        LockTypeTableLock,
		},
		{
			SQL:             "COMMIT",
			ExecuteDuration: 0.2,
			LockType:        LockTypeIntentionExclusive,
		},
	}

	// 创建模拟的锁信息
	transactionLockInfo := []*TransactionLockInfo{
		{
			LockType:      LockTypeTableLock,
			TableName:     "order_items",
			CreateLockSQL: "UPDATE order_items SET status = 'shipped' WHERE order_id = 12345",
		},
		{
			LockType:      LockTypeIntentionShared,
			TableName:     "inventory",
			CreateLockSQL: "SELECT * FROM inventory WHERE product_id = 555 FOR UPDATE",
		},
		{
			LockType:      LockTypeIntentionExclusive,
			TableName:     "products",
			CreateLockSQL: "SELECT * FROM products WHERE category_id = 3 ORDER BY price DESC LIMIT 10",
		},
	}

	// 组合所有信息到RelatedTransactionInfo
	relatedTransactionInfo := &RelatedTransactionInfo{
		TransactionInfo:     transactionInfo,
		TransactionTimeline: transactionTimeline,
		TransactionSQLList:  transactionSQLList,
		TransactionLockInfo: transactionLockInfo,
	}

	// 返回响应
	return c.JSON(http.StatusOK, GetSqlManageSqlPerformanceInsightsRelatedTransactionResp{
		BaseRes: controller.BaseRes{
			Code:    0,
			Message: "success",
		},
		Data: relatedTransactionInfo,
	})
}
