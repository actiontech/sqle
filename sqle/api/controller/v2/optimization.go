package v2

import (
	"time"

	"github.com/actiontech/sqle/sqle/server/optimization/sql_flash"

	"github.com/actiontech/sqle/sqle/api/controller"

	"github.com/labstack/echo/v4"
)

type OptimizeSQLReq struct {
	OptimizationName string  `json:"optimization_name" form:"optimization_name" example:"optmz_2024031412091244" valid:"required"`
	SQLContent       string  `json:"sql_content" form:"sql_content" example:"select * from t1; select * from t2;"`
	InstanceName     *string `json:"instance_name" form:"instance_name" example:"instance1"`
	SchemaName       *string `json:"schema_name" form:"schema_name" example:"schema1"`
	ExplainInfo      string  `json:"explain_info" form:"explain_info"`
	Metadata         string  `json:"metadata" form:"metadata"`
}

type OptimizeSQLRes struct {
	controller.BaseRes
	Data *OptimizeSQLResData `json:"data"`
}

type OptimizeSQLResData struct {
	OptimizationRecordId string `json:"sql_optimization_record_id"`
}

// SQLOptimize
// @Summary 优化SQL
// @Description optimize sql
// @Id SQLOptimizeV2
// @Description 1. formData[sql]: sql content;
// @Description 2. file[input_sql_file]: it is a sql file;
// @Description 3. file[input_mybatis_xml_file]: it is mybatis xml file, sql will be parsed from it.
// @Description 4. file[input_zip_file]: it is ZIP file that sql will be parsed from xml or sql file inside it.
// @Description 5. formData[git_http_url]:the url which scheme is http(s) and end with .git.
// @Description 6. formData[git_user_name]:The name of the user who owns the repository read access.
// @Description 7. formData[git_user_password]:The password corresponding to git_user_name.
// @Accept mpfd
// @Produce json
// @Tags sql_optimization
// @Security ApiKeyAuth
// @Param req body v2.OptimizeSQLReq true "sqls that should be optimization"
// @Param project_name path string true "project name"
// @Param instance_name formData string false "instance name"
// @Param schema_name formData string false "schema of instance"
// @Param db_type formData string false "db type of instance"
// @Param optimization_name formData string true "optimization name"
// @Param sql_content formData string false "sqls for audit"
// @Param explain_info formData string false "explain info"
// @Param metadata formData string false "metadata"
// @Param input_sql_file formData file false "input SQL file"
// @Param input_zip_file formData file false "input ZIP file"
// @Param input_mybatis_xml_file formData file false "input mybatis XML file"
// @Param git_http_url formData string false "git repository url"
// @Param git_user_name formData string false "the name of user to clone the repository"
// @Param git_user_password formData string false "the password corresponding to git_user_name"
// @Success 200 {object} v2.OptimizeSQLRes
// @router /v2/projects/{project_name}/sql_optimization_records [post]
func SQLOptimize(c echo.Context) error {
	return sqlOptimize(c)
}

type GetOptimizationRecordsReq struct {
	FuzzySearch          string `json:"fuzzy_search" query:"fuzzy_search"`
	FilterInstanceName   string `json:"filter_instance_name" query:"filter_instance_name"`
	FilterCreateTimeFrom string `json:"filter_create_time_from" query:"filter_create_time_from"`
	FilterCreateTimeTo   string `json:"filter_create_time_to" query:"filter_create_time_to"`
	FilterStatus         string `json:"filter_status" query:"filter_status" enums:"optimizing,failed,finish"`
	PageIndex            uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize             uint32 `json:"page_size" query:"page_size" valid:"required"`
}

// SQL优化记录结构体
type OptimizationRecord struct {
	OptimizationID     string    `json:"optimization_id"` // 优化ID
	OptimizationName   string    `json:"optimization_name"`
	InstanceName       string    `json:"instance_name"`                           // 数据源
	DBType             string    `json:"db_type"`                                 // 数据库类型
	CreatedTime        time.Time `json:"created_time"`                            // 创建时间
	CreatedUser        string    `json:"created_user"`                            // 创建人
	Status             string    `json:"status" enums:"optimizing,failed,finish"` // 优化状态
	StatusDetail       string    `json:"status_detail"`                           // 优化状态详情
	PerformanceImprove float64   `json:"performance_improve"`                     // 优化提升性能
	NumberOfRule       int       `json:"number_of_rule"`                          // 优化规则数量
	NumberOfIndex      int       `json:"number_of_index"`                         // 优化索引数量
}

type GetOptimizationRecordsRes struct {
	controller.BaseRes
	Data      []OptimizationRecord `json:"data"`
	TotalNums uint64               `json:"total_nums"`
}

// GetOptimizationRecords
// @Summary 获取SQL优化记录列表
// @Description get sql optimization records
// @Tags sql_optimization
// @Id GetOptimizationRecordsV2
// @Security ApiKeyAuth
// @Param fuzzy_search query string false "fuzzy search for optimization_id or create_username"
// @Param filter_instance_name query string false "filter instance name"
// @Param filter_create_time_from query string false "filter create time from"
// @Param filter_create_time_to query string false "filter create time to"
// @Param filter_status query string false "filter status" Enums(optimizing,failed,finish)
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Param project_name path string true "project name"
// @Success 200 {object} v2.GetOptimizationRecordsRes
// @router /v2/projects/{project_name}/sql_optimization_records [get]
func GetOptimizationRecords(c echo.Context) error {
	return getOptimizationRecords(c)
}

type GetOptimizationDetailReq struct {
	OptimizationId string `json:"optimization_id" query:"optimization_id" valid:"required"`
}

type GetOptimizationDetailRes struct {
	controller.BaseRes
	Data *OptimizationSQLDetail `json:"data"`
}

type OptimizationSQLDetail struct {
	OptimizationId string `json:"optimization_id"`
	ID             uint   `json:"id"`
	Status         string `json:"status" enums:"optimizing,failed,finish"` // SQLe 维护的状态
	StatusDetail   string `json:"status_detail"`                           // SQLe 维护的状态详情

	// SQL Flash相关字段
	OriginSQL       string                    `json:"origin_sql"`        // 原始SQL
	Metadata        string                    `json:"metadata"`          // 数据库元数据信息
	TotalState      string                    `json:"total_state"`       // 总状态
	OriginQueryPlan *sql_flash.QueryPlan      `json:"origin_query_plan"` // 原始SQL查询计划
	OptimizeDetail  *sql_flash.OptimizeDetail `json:"optimize"`          // 优化详情
	TotalAnalysis   *sql_flash.TotalAnalysis  `json:"total_analysis"`    // 总体分析
	AdvisedIndex    *sql_flash.AdvisedIndex   `json:"advised_index"`     // 索引建议详情
}

// GetOptimizationSQLDetail
// @Summary 获取SQL优化语句详情
// @Description get sql optimization record
// @Id GetOptimizationSQLDetailV2
// @Tags sql_optimization
// @Param project_name path string true "project name"
// @Param optimization_record_id path string true "sql optimization record id"
// @Security ApiKeyAuth
// @Success 200 {object} v2.GetOptimizationDetailRes
// @router /v2/projects/{project_name}/sql_optimization_records/{optimization_record_id}/detail [get]
func GetOptimizationSQLDetail(c echo.Context) error {
	return getOptimizationSQL(c)
}
