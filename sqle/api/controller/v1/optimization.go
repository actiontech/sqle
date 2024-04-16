package v1

import (
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"

	"github.com/labstack/echo/v4"
)

type OptimizeSQLReq struct {
	DBType           string  `json:"db_type" form:"db_type" example:"MySQL"`
	SQLContent       string  `json:"sql_content" form:"sql_content" example:"select * from t1; select * from t2;" valid:"required"`
	OptimizationName string  `json:"optimization_name" form:"optimization_name" example:"optmz_2024031412091244" valid:"required"`
	InstanceName     *string `json:"instance_name" form:"instance_name" example:"instance1"`
	SchemaName       *string `json:"schema_name" form:"schema_name" example:"schema1"`
}

type OptimizeSQLRes struct {
	controller.BaseRes
	Data *OptimizeSQLResData `json:"data"`
}
type OptimizeSQLResData struct {
	OptimizationRecordId string `json:"sql_optimization_record_id"`
}

// @Summary 优化SQL
// @Description optimize sql
// @Description 1. formData[sql]: sql content;
// @Description 2. file[input_sql_file]: it is a sql file;
// @Description 3. file[input_mybatis_xml_file]: it is mybatis xml file, sql will be parsed from it.
// @Description 4. file[input_zip_file]: it is ZIP file that sql will be parsed from xml or sql file inside it.
// @Description 5. formData[git_http_url]:the url which scheme is http(s) and end with .git.
// @Description 6. formData[git_user_name]:The name of the user who owns the repository read access.
// @Description 7. formData[git_user_password]:The password corresponding to git_user_name.
// @Id OptimizeSQLReq
// @Tags sql_optimization
// @Security ApiKeyAuth
// @Param req body v1.OptimizeSQLReq true "sqls that should be optimization"
// @Param project_name path string true "project name"
// @Param instance_name formData string false "instance name"
// @Param schema_name formData string false "schema of instance"
// @Param db_type formData string false "db type of instance"
// @Param sql_content formData string false "sqls for audit"
// @Param optimization_name formData string true "optimization name"
// @Param input_sql_file formData file false "input SQL file"
// @Param input_mybatis_xml_file formData file false "input mybatis XML file"
// @Param input_zip_file formData file false "input ZIP file"
// @Param git_http_url formData string false "git repository url"
// @Param git_user_name formData string false "the name of user to clone the repository"
// @Param git_user_password formData string false "the password corresponding to git_user_name"
// @Success 200 {object} v1.OptimizeSQLRes
// @router /v1/projects/{project_name}/sql_optimization_records [post]
func SQLOptimizate(c echo.Context) error {
	return sqlOptimizate(c)
}

type GetOptimizationRecordReq struct {
	OptimizationId string `json:"optimization_id" query:"optimization_id" valid:"required"`
}

type GetOptimizationRecordRes struct {
	controller.BaseRes
	Data *OptimizationDetail `json:"data"`
}

type OptimizationDetail struct {
	OptimizationID       string              `json:"optimization_id"`
	OptimizationName     string              `json:"optimization_name"`
	InstanceNmae         string              `json:"instance_name"`
	DBType               string              `json:"db_type"`
	CreatedTime          time.Time           `json:"created_time"`
	CreatedUser          string              `json:"created_user"`
	Optimizationsummary  Optimizationsummary `json:"basic_summary"`
	IndexRecommendations []string            `json:"index_recommendations"`
}

type Optimizationsummary struct {
	NumberOfQuery          int     `json:"number_of_query"`
	NumberOfSyntaxError    int     `json:"number_of_syntax_error"`
	NumberOfRewrite        int     `json:"number_of_rewrite"`
	NumberOfRewrittenQuery int     `json:"number_of_rewritten_query"`
	NumberOfIndex          int     `json:"number_of_index"`
	NumberOfQueryIndex     int     `json:"number_of_query_index"`
	PerformanceGain        float64 `json:"performance_gain"`
}

// GetOptimizationRecord
// @Summary 获取SQL优化记录
// @Description get sql optimization record
// @Id GetOptimizationRecordReq
// @Tags sql_optimization
// @Param project_name path string true "project name"
// @Param optimization_record_id path string true "sql optimization record id"
// @Param sql query string false "sql"
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetOptimizationRecordRes
// @router /v1/projects/{project_name}/sql_optimization_records/{optimization_record_id}/ [get]
func GetOptimizationRecord(c echo.Context) error {
	return getOptimizationRecord(c)
}

type GetOptimizationRecordsReq struct {
	FuzzySearch          string `json:"fuzzy_search" query:"fuzzy_search"`
	FilterInstanceName   uint64 `json:"filter_instance_name" query:"filter_instance_name"`
	FilterCreateTimeFrom string `json:"filter_create_time_from" query:"filter_create_time_from"`
	FilterCreateTimeTo   string `json:"filter_create_time_to" query:"filter_create_time_from"`
	PageIndex            uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize             uint32 `json:"page_size" query:"page_size" valid:"required"`
}

// SQL优化记录结构体
type OptimizationRecord struct {
	OptimizationID   string    `json:"optimization_id"` // 优化ID
	OptimizationName string    `json:"optimization_name"`
	InstanceNmae     string    `json:"instance_name"`    // 数据源
	DBType           string    `json:"db_type"`          // 数据库类型
	PerformanceGain  float64   `json:"performance_gain"` // 优化提升性能
	CreatedTime      time.Time `json:"created_time"`     // 创建时间
	CreatedUser      string    `json:"created_user"`     // 创建人
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
// @Id getOptimizationRecords
// @Security ApiKeyAuth
// @Param fuzzy_search query string false "fuzzy search for optimization_id or create_username"
// @Param filter_instance_name query string false "filter instance name"
// @Param filter_create_time_from query string false "filter create time from"
// @Param filter_create_time_to query string false "filter create time to"
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Param project_name path string true "project name"
// @Success 200 {object} v1.GetOptimizationRecordsRes
// @router /v1/projects/{project_name}/sql_optimization_records [get]
func GetOptimizationRecords(c echo.Context) error {
	return getOptimizationRecords(c)
}

type GetOptimizationSQLReq struct {
	OptimizationId string `json:"optimization_id" query:"optimization_id" valid:"required"`
}

type GetOptimizationSQLRes struct {
	controller.BaseRes
	Data *OptimizationSQLDetail `json:"data"`
}

type OptimizationSQLDetail struct {
	OriginalSQL              string                  `json:"original_sql"`          // 原始SQL
	OptimizedSQL             string                  `json:"optimized_sql"`         // 优化后的SQL
	TriggeredRule            []RewriteRule           `json:"triggered_rule"`        // 触发的规则
	IndexRecommendations     []string                `json:"index_recommendations"` // 索引建议
	ExplainValidationDetails ExplainValidationDetail `json:"explain_validation_details"`
}
type ExplainValidationDetail struct {
	BeforeCost        float64 `json:"before_cost"`
	AfterCost         float64 `json:"after_cost"`
	BeforePlan        string  `json:"before_plan"`
	AfterPlan         string  `json:"after_plan"`
	PerformImprovePer float64 `json:"perform_improve_per"`
}

type RewriteRule struct {
	RuleCode            string `json:"rule_code"`
	RuleName            string `json:"rule_name"`
	Message             string `json:"message"`
	RewrittenQueriesStr string `json:"rewritten_queries_str"`
	ViolatedQueriesStr  string `json:"violated_queries_str"`
}

// GetOptimizationSQLDetail
// @Summary 获取SQL优化语句详情
// @Description get sql optimization record
// @Id GetOptimizationReq
// @Tags sql_optimization
// @Param project_name path string true "project name"
// @Param optimization_record_id path string true "sql optimization record id"
// @Param number path string true "optimization record sql  number"
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetOptimizationSQLRes
// @router /v1/projects/{project_name}/sql_optimization_records/{optimization_record_id}/sqls/{number}/ [get]
func GetOptimizationSQLDetail(c echo.Context) error {
	return getOptimizationSQL(c)
}

type GetOptimizationSQLsReq struct {
	PageIndex uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type OptimizationSQL struct {
	Number              uint64  `json:"number"`
	OriginalSQL         string  `json:"original_sql"`
	NumberOfRewrite     int     `json:"number_of_rewrite"`
	NumberOfSyntaxError int     `json:"number_of_syntax_error"`
	NumberOfIndex       int     `json:"number_of_index"`
	NumberOfHitIndex    int     `json:"number_of_hit_index"`
	Performance         float64 `json:"performance"`
	ContributingIndices string  `json:"contributing_indices"`
}

type GetOptimizationSQLsRes struct {
	controller.BaseRes
	Data      []OptimizationSQL `json:"data"`
	TotalNums uint64            `json:"total_nums"`
}

// GetOptimizationSQLs
// @Summary 获取SQL优化语句列表
// @Description get sql optimization sqls
// @Tags sql_optimization
// @Id getOptimizationSQLs
// @Security ApiKeyAuth
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Param project_name path string true "project name"
// @Param optimization_record_id path string true "optimization record id"
// @Success 200 {object} v1.GetOptimizationSQLsRes
// @router /v1/projects/{project_name}/sql_optimization_records/{optimization_record_id}/sqls [get]
func GetOptimizationSQLs(c echo.Context) error {
	return getOptimizationSQLs(c)
}
