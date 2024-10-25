package v1

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/api/controller"
	dms "github.com/actiontech/sqle/sqle/dms"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/server/compare"
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
	req := new(GetDatabaseCompareReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	baseInst, exist, err := dms.GetInstancesById(c.Request().Context(), req.BaseDBObject.InstanceId)
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrInstanceNotExist)
	} else if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
	}
	comparedInst, exist, err := dms.GetInstancesById(c.Request().Context(), req.ComparedDBObject.InstanceId)
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrInstanceNotExist)
	} else if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
	}
	compared := &compare.Compared{
		BaseInstance:     baseInst,
		ComparedInstance: comparedInst,
	}
	compared.ExecDatabaseCompare(req.BaseDBObject.SchemaName, req.ComparedDBObject.SchemaName)
	return nil
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
	SQLStatement string         `json:"sql_statement"`
	AuditResult  []*AuditResult `json:"audit_result"`
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

	return nil
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
	Data []DatabaseDiffModifySQL `json:"data"`
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
// @Param gen_alterantl_sql body v1.GenAlterantlSQLReqV1 true "gen alterantl sql request"
// @Success 200 {object} v1.GenAlterantlSQLResV1
// @router /v1/projects/{project_name}/database_compare/alterant_sql [post]
func GenDatabaseDiffModifySQLs(c echo.Context) error {
	req := new(GenAlterantlSQLReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	baseInst, exist, err := dms.GetInstancesById(c.Request().Context(), req.BaseInstanceId)
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrInstanceNotExist)
	} else if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
	}
	comparedInst, exist, err := dms.GetInstancesById(c.Request().Context(), req.ComparedInstanceId)
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrInstanceNotExist)
	} else if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
	}
	compared := &compare.Compared{
		BaseInstance:     baseInst,
		ComparedInstance: comparedInst,
	}
	// baseSchemas := make([]*driverV2.DatabasSchemaInfo, len(req.BaseDatabaseObjects))
	// for i, baseSchema := range req.BaseDatabaseObjects {
	// 	baseSchemas[i] = &driverV2.DatabasSchemaInfo{
	// 		ScheamName:      baseSchema.SchemaName,
	// 		DatabaseObjects: convertDatabaseSchemaObjectToDriver(baseSchema.DatabaseObjects),
	// 	}
	// }
	schemaObjs := make([]*driverV2.DatabasCompareSchemaInfo, len(req.DatabaseSchemaObjects))
	for i, schema := range req.DatabaseSchemaObjects {
		schemaObjs[i] = &driverV2.DatabasCompareSchemaInfo{
			BaseScheamName:     schema.BaseSchemaName,
			ComparedScheamName: schema.ComparedSchemaName,
			DatabaseObjects:    convertDatabaseSchemaObjectToDriver(schema.DatabaseObjects),
		}
	}

	diffsqls, err := compared.GetDatabaseDiffModifySQLs(schemaObjs)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	fmt.Printf(diffsqls[0].ModifySQLs[0])
	return nil
}

func convertDatabaseSchemaObjectToDriver(Objs []*DatabaseObject) []*driverV2.DatabaseObject {
	baseObjects := make([]*driverV2.DatabaseObject, len(Objs))
	for i, baseObject := range Objs {
		baseObjects[i] = &driverV2.DatabaseObject{
			ObjectName: baseObject.ObjectName,
			ObjectType: baseObject.ObjectType,
		}
	}
	return baseObjects
}
