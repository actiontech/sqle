package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
)

type GetDatabaseComparisonReqV1 struct {
	BaseDBObject       *DatabaseComparisonObject `json:"base_db_object" query:"base_db_object"`
	ComparisonDBObject *DatabaseComparisonObject `json:"comparison_db_object" query:"comparison_db_object"`
}

type DatabaseComparisonObject struct {
	InstanceId string  `json:"instance_id" query:"instance_id"`
	SchemaName *string `json:"schema_name,omitempty" query:"schema_name"`
}
type DatabaseComparisonResV1 struct {
	controller.BaseRes
	Data []*SchemaObject `json:"data"`
}

type SchemaObject struct {
	BaseSchemaName       string                `json:"base_schema_name"`
	ComparisonSchemaName string                `json:"comparison_schema_name"`
	ComparisonResult     string                `json:"comparison_result" enums:"same,inconsistent,base_not_exist,comparison_not_exist"`
	DatabaseDiffObjects  []*DatabaseDiffObject `json:"database_diff_objects,omitempty"`
	InconsistentNum      int                   `json:"inconsistent_num"`
}

type DatabaseDiffObject struct {
	InconsistentNum    int                 `json:"inconsistent_num"`
	ObjectType         string              `json:"object_type" enums:"TABLE,VIEW,PROCEDURE,TIGGER,EVENT,FUNCTION"`
	ObjectsDiffResults []*ObjectDiffResult `json:"objects_diff_result,omitempty"`
}

type ObjectDiffResult struct {
	ComparisonResult string `json:"comparison_result" enums:"same,inconsistent,base_not_exist,comparison_not_exist"`
	ObjectName       string `json:"object_name"`
}

// @Summary 执行数据库结构对比并获取结果
// @Description get database comparison
// @Id executeDatabaseComparisonV1
// @Tags database_comparison
// @Security ApiKeyAuth
// @Accept json
// @Param project_name path string true "project name"
// @Param database_comparison body v1.GetDatabaseComparisonReqV1 true "get database comparison request"
// @Success 200 {object} v1.DatabaseComparisonResV1
// @router /v1/projects/{project_name}/database_comparison/execute_comparison [post]
func ExecuteDatabaseComparison(c echo.Context) error {

	return getDatabaseComparison(c)
}

type DatabaseObject struct {
	ObjectName string `json:"object_name"`
	ObjectType string `json:"object_type" enums:"TABLE,VIEW,PROCEDURE,TIGGER,EVENT,FUNCTION"`
}

type GetComparisonStatementsReqV1 struct {
	DatabaseComparisonObject GetDatabaseComparisonReqV1 `json:"database_comparison_object"`
	DatabaseObject           *DatabaseObject            `json:"database_object"`
}

type DatabaseComparisonStatementsResV1 struct {
	controller.BaseRes
	Data *DatabaseComparisonStatements `json:"data"`
}

type DatabaseComparisonStatements struct {
	BaseSQL        *SQLStatement `json:"base_sql"`
	ComparisondSQL *SQLStatement `json:"comparison_sql"`
}

type SQLStatement struct {
	AuditError            string                       `json:"audit_error,omitempty"`
	SQLStatementWithAudit *SQLStatementWithAuditResult `json:"sql_statement_with_audit"`
}

type SQLStatementWithAuditResult struct {
	SQLStatement string            `json:"sql_statement"`
	AuditResults []*SQLAuditResult `json:"audit_results,omitempty"`
}

type SQLAuditResult struct {
	Level               string                    `json:"level" example:"warn"`
	Message             string                    `json:"message" example:"避免使用不必要的内置函数md5()"`
	RuleName            string                    `json:"rule_name"`
	DbType              string                    `json:"db_type"`
	ExecutionFailed     bool                      `json:"execution_failed" example:"false"`
	ErrorInfo           string                    `json:"error_info"`
	I18nAuditResultInfo model.I18nAuditResultInfo `json:"i18n_audit_result_info"`
}

// @Summary 获取对比语句
// @Description get database comparison detail
// @Id getComparisonStatementV1
// @Tags database_comparison
// @Security ApiKeyAuth
// @Accept json
// @Param project_name path string true "project name"
// @Param database_comparison_object body v1.GetComparisonStatementsReqV1 true "get database comparison statement request"
// @Success 200 {object} v1.DatabaseComparisonStatementsResV1
// @router /v1/projects/{project_name}/database_comparison/comparison_statements [post]
func GetComparisonStatement(c echo.Context) error {

	return getComparisonStatement(c)
}

type GenModifylSQLReqV1 struct {
	BaseInstanceId        string                  `json:"base_instance_id"`
	ComparisonInstanceId  string                  `json:"comparison_instance_id"`
	DatabaseSchemaObjects []*DatabaseSchemaObject `json:"database_schema_objects"`
}

type DatabaseSchemaObject struct {
	BaseSchemaName       string            `json:"base_schema_name"`
	ComparisonSchemaName string            `json:"comparison_schema_name"`
	DatabaseObjects      []*DatabaseObject `json:"database_objects"`
}

type GenModifySQLResV1 struct {
	controller.BaseRes
	Data []*DatabaseDiffModifySQL `json:"data"`
}

type DatabaseDiffModifySQL struct {
	SchemaName string                         `json:"schema_name"`
	ModifySQLs []*SQLStatementWithAuditResult `json:"modify_sqls"`
	AuditError string                         `json:"audit_error,omitempty"`
}

// @Summary 生成变更SQL
// @Description generate database diff modify sqls
// @Id genDatabaseDiffModifySQLsV1
// @Tags database_comparison
// @Security ApiKeyAuth
// @Accept json
// @Param project_name path string true "project name"
// @Param gen_modify_sql body v1.GenModifylSQLReqV1 true "generate database diff modify sqls request"
// @Success 200 {object} v1.GenModifySQLResV1
// @router /v1/projects/{project_name}/database_comparison/modify_sql_statements [post]
func GenDatabaseDiffModifySQLs(c echo.Context) error {

	return genDatabaseDiffModifySQLs(c)
}
