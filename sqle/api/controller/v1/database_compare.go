package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
)

type DatabaseCompareResV1 struct {
	controller.BaseRes
	Data []*SchemaObject `json:"data"`
}

type SchemaObject struct {
	SchemaName          string                `json:"schema_name"`
	ComparedResult      string                `json:"compared_result"`
	DatabaseDiffObjects []*DatabaseDiffObject `json:"database_diff_objects,omitempty"`
	InconsistentNum     int                   `json:"inconsistent_num"`
}

type DatabaseDiffObject struct {
	InconsistentNum    int                 `json:"inconsistent_num"`
	ObjectType         string              `json:"object_type" enums:"TABLE,VIEW,PROCEDURE,TIGGER,EVENT,FUNCTION"`
	ObjectsDiffResults []*ObjectDiffResult `json:"objects_diff_result,omitempty"`
}

type ObjectDiffResult struct {
	ComparedResult string `json:"compared_result"`
	ObjectName     string `json:"object_name"`
}

type GetDatabaseCompareReqV1 struct {
	BaseDBObject     *DatabaseCompareObject `json:"base_db_object" query:"base_db_object"`
	ComparedDBObject *DatabaseCompareObject `json:"compared_db_object" query:"compared_db_object"`
}

type DatabaseCompareObject struct {
	InstanceId string  `json:"instance_id" query:"instance_id"`
	SchemaName *string `json:"schema_name,omitempty" query:"instance_id"`
}

// @Summary 获取对比结果
// @Description get database compare
// @Id getDatabaseCompareV1
// @Tags database_compare
// @Security ApiKeyAuth
// @Accept json
// @Param project_name path string true "project name"
// @Param database_compared body v1.GetDatabaseCompareReqV1 true "get database compare request"
// @Success 200 {object} v1.DatabaseCompareResV1
// @router /v1/projects/{project_name}/database_compare [post]
func GetDatabaseCompare(c echo.Context) error {

	return getDatabaseCompare(c)
}

type DatabaseObject struct {
	ObjectName string `json:"object_name"`
	ObjectType string `json:"object_type" enums:"TABLE,VIEW,PROCEDURE,TIGGER,EVENT,FUNCTION"`
}

type GetComparedStatementsReqV1 struct {
	DatabaseComparedObject GetDatabaseCompareReqV1 `json:"database_compared_object"`
	DatabaseObject         *DatabaseObject         `json:"database_object"`
}

type DatabaseCompareStatementsResV1 struct {
	BaseSQL     *SQLStatementWithAuditResult `json:"base_sql"`
	ComparedSQL *SQLStatementWithAuditResult `json:"compared_sql"`
}

type SQLStatementWithAuditResult struct {
	SQLStatement string             `json:"sql_statement"`
	AuditResults model.AuditResults `json:"audit_results"`
	AuditLevel   string             `json:"audit_level"`
}

// @Summary 获取对比语句
// @Description gen database compare detail
// @Id getCompareStatementV1
// @Tags database_compare
// @Security ApiKeyAuth
// @Accept json
// @Param project_name path string true "project name"
// @Param database_compared_object body v1.GetComparedStatementsReqV1 true "get database compared statement request"
// @Success 200 {object} v1.DatabaseCompareStatementsResV1
// @router /v1/projects/{project_name}/database_compare/statements [post]
func GetCompareStatement(c echo.Context) error {

	return getCompareStatement(c)
}

type GenAlterantlSQLReqV1 struct {
	BaseInstanceId        string                  `json:"base_instance_id"`
	ComparedInstanceId    string                  `json:"compared_instance_id"`
	DatabaseSchemaObjects []*DatabaseSchemaObject `json:"database_schema_objects"`
}

type DatabaseSchemaObject struct {
	BaseSchemaName     string            `json:"base_schema_name"`
	ComparedSchemaName string            `json:"compared_schema_name"`
	DatabaseObjects    []*DatabaseObject `json:"database_objects"`
}

type GenAlterantlSQLResV1 struct {
	controller.BaseRes
	Data []*DatabaseDiffModifySQL `json:"data"`
}

type DatabaseDiffModifySQL struct {
	SchemaName string                         `json:"schema_name"`
	ModifySQLs []*SQLStatementWithAuditResult `json:"modify_sqls"`
}

// @Summary 生成变更SQL
// @Description gen database diff modify sqls
// @Id genDatabaseDiffModifySQLsV1
// @Tags database_compare
// @Security ApiKeyAuth
// @Accept json
// @Param project_name path string true "project name"
// @Param gen_alterantl_sql body v1.GenAlterantlSQLReqV1 true "gen database diff modify sqls request"
// @Success 200 {object} v1.GenAlterantlSQLResV1
// @router /v1/projects/{project_name}/database_compare/alterant_sql [post]
func GenDatabaseDiffModifySQLs(c echo.Context) error {

	return genDatabaseDiffModifySQLs(c)
}
