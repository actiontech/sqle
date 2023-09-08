package v1

import (
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"

	"github.com/labstack/echo/v4"
)

type SqlAuditRecordReqV1 struct {
	InstanceName   string `json:"instance_name" form:"instance_name" example:"inst_1"`
	InstanceSchema string `json:"instance_schema" form:"instance_schema" example:"db1"`
	Sql            string `json:"sql" form:"sql" example:"alter table tb1 drop columns c1"`
}

type SqlAuditRecordResV1 struct {
	controller.BaseRes
	Data *SqlAuditRecordResData `json:"data"`
}

type SqlAuditRecordResData struct {
	Id   string          `json:"sql_audit_record_id"`
	Task *AuditTaskResV1 `json:"task"`
}

// CreateSQLAuditRecord
// @Summary 创建SQL审核记录
// @Id CreateSqlAuditRecordV1
// @Description create SQL audit record
// @Description 1. formData[sql]: sql content;
// @Description 2. file[input_sql_file]: it is a sql file;
// @Description 3. file[input_mybatis_xml_file]: it is mybatis xml file, sql will be parsed from it.
// @Description 4. file[input_zip_file]: it is ZIP file that sql will be parsed from xml or sql file inside it.
// @Accept mpfd
// @Produce json
// @Tags sql_audit_record
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instance_name formData string false "instance name"
// @Param instance_schema formData string false "schema of instance"
// @Param sql formData string false "sqls for audit"
// @Param input_sql_file formData file false "input SQL file"
// @Param input_mybatis_xml_file formData file false "input mybatis XML file"
// @Param input_zip_file formData file false "input ZIP file"
// @Success 200 {object} v1.SqlAuditRecordResV1
// @router /v1/projects/{project_name}/sql_audit_record [post]
func CreateSQLAuditRecord(c echo.Context) error {
	return nil
}

type UpdateSQLAuditRecordReqV1 struct {
	SQLAuditRecordId string    `json:"sql_audit_record_id" valid:"required"`
	Tags             *[]string `json:"tags"`
}

// UpdateSQLAuditRecordV1
// @Summary 更新SQL审核记录
// @Description update SQL audit record
// @Accept json
// @Id updateSQLAuditRecordV1
// @Tags sql_audit_record
// @Security ApiKeyAuth
// @Param param body v1.UpdateSQLAuditRecordReqV1 true "update SQL audit record"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/sql_audit_record/{sql_audit_record_id} [patch]
func UpdateSQLAuditRecordV1(c echo.Context) error {
	return nil
}

type GetSQLAuditRecordsReqV1 struct {
	FuzzySearchSQLAuditRecordId string `json:"fuzzy_search_sql_audit_record_id" query:"fuzzy_search_sql_audit_record_id"`
	FilterSQLAuditStatus        string `json:"filter_sql_audit_status" query:"filter_sql_audit_status" enums:"auditing,successfully,"`
	FilterTags                  string `json:"filter_tags" query:"filter_tags"`
}

type SQLAuditRecord struct {
	SQLAuditRecordId uint       `json:"sql_audit_record_id"`
	SQLAuditStatus   string     `json:"sql_audit_status"`
	Tags             []string   `json:"tags"`
	AuditScore       int32      `json:"audit_score"`
	AuditPassRate    float64    `json:"audit_pass_rate"`
	AuditStartedAt   *time.Time `json:"audit_started_at"`
	TaskId           uint       `json:"task_id"`
}

type GetSQLAuditRecordsResV1 struct {
	controller.BaseRes
	Data []SQLAuditRecord `json:"data"`
}

// GetSQLAuditRecordsV1
// @Summary 获取SQL审核记录列表
// @Description get sql audit records
// @Tags sql_audit_record
// @Id getSQLAuditRecordsV1
// @Security ApiKeyAuth
// @Param fuzzy_search_sql_audit_record_id query string false "fuzzy search sql audit record_id"
// @Param fuzzy_search_tags query string false "fuzzy search tags"
// @Param filter_sql_audit_status query string false "filter sql audit status" Enums(auditing,successfully)
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Param project_name path string true "project name"
// @Success 200 {object} v1.GetSQLAuditRecordsResV1
// @router /v1/projects/{project_name}/sql_audit_record [get]
func GetSQLAuditRecordsV1(c echo.Context) error {
	return nil
}

type GetSQLAuditRecordTagTipsResV1 struct {
	controller.BaseRes
	Tags []string `json:"data"`
}

// GetSQLAuditRecordTagTipsV1
// @Summary 获取SQL审核记录标签列表
// @Description get sql audit record tag tips
// @Tags sql_audit_record
// @Id GetSQLAuditRecordTagTipsV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetSQLAuditRecordTagTipsResV1
// @router /v1/projects/{project_name}/sql_audit_record/tag_tips [get]
func GetSQLAuditRecordTagTipsV1(c echo.Context) error {
	return nil
}
