//go:build enterprise
// +build enterprise

package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/common"
	dms "github.com/actiontech/sqle/sqle/dms"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/actiontech/sqle/sqle/server/compare"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func getDatabaseComparison(c echo.Context) error {
	req := new(GetDatabaseComparisonReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if req.BaseDBObject == nil || req.ComparisonDBObject == nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, fmt.Errorf("neither the base instance nor the comparison instance can be empty")))
	}
	if (req.BaseDBObject.SchemaName == nil) != (req.ComparisonDBObject.SchemaName == nil) {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, fmt.Errorf("the base instance and comparison instance must be consistent")))
	}
	logger := log.NewEntry().WithField("database comparison", "execute database comparison")
	// 获取基准数据源和比对数据源实例信息
	baseInst, err := getInstanceById(c, req.BaseDBObject.InstanceId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	comparedInst, err := getInstanceById(c, req.ComparisonDBObject.InstanceId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	schemaNames := make([]*driverV2.DatabasCompareSchemaInfo, 0, 1)
	if req.BaseDBObject.SchemaName != nil && req.ComparisonDBObject.SchemaName != nil {
		schemaNames = append(schemaNames, &driverV2.DatabasCompareSchemaInfo{
			BaseSchemaName:     *req.BaseDBObject.SchemaName,
			ComparedSchemaName: *req.ComparisonDBObject.SchemaName,
		})
	}
	compared := &compare.Compared{
		BaseInstance:     baseInst,
		ComparedInstance: comparedInst,
		ObjInfos:         schemaNames,
	}
	execCompareRes, err := compared.ExecDatabaseCompare(c.Request().Context(), logger)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
	}
	comparedRes := make([]*SchemaObject, len(execCompareRes))
	for i, compared := range execCompareRes {
		comparedRes[i] = &SchemaObject{
			BaseSchemaName:       compared.BaseSchemaName,
			ComparisonSchemaName: compared.ComparisonSchemaName,
			InconsistentNum:      compared.InconsistentNum,
			ComparisonResult:     compared.ComparedResult,
			DatabaseDiffObjects:  convertDatabaseDiffObject(compared.DatabaseDiffObjects),
		}
	}

	return c.JSON(http.StatusOK, &DatabaseComparisonResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    comparedRes,
	})
}

func convertDatabaseDiffObject(comparedDiffs []*compare.DatabaseDiffObject) []*DatabaseDiffObject {
	diffObjects := make([]*DatabaseDiffObject, len(comparedDiffs))
	for j, diffObject := range comparedDiffs {
		objects := make([]*ObjectDiffResult, len(diffObject.ObjectsDiffResults))
		for m, object := range diffObject.ObjectsDiffResults {
			objects[m] = &ObjectDiffResult{
				ObjectName:       object.ObjectName,
				ComparisonResult: object.ComparedResult,
			}
		}
		diffObjects[j] = &DatabaseDiffObject{
			ObjectType:         diffObject.ObjectType,
			InconsistentNum:    diffObject.InconsistentNum,
			ObjectsDiffResults: objects,
		}
	}
	return diffObjects
}

func getInstanceById(c echo.Context, instanceID string) (*model.Instance, error) {
	inst, exist, err := dms.GetInstancesById(c.Request().Context(), instanceID)
	if err != nil {
		return nil, errors.New(errors.DataConflict, err)
	}
	if !exist {
		return nil, ErrInstanceNotExist
	}
	return inst, nil
}

func getComparisonStatement(c echo.Context) error {
	req := new(GetComparisonStatementsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	projectID, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
	}
	baseObject := req.DatabaseComparisonObject.BaseDBObject
	comparisonObject := req.DatabaseComparisonObject.ComparisonDBObject
	databaseObject := req.DatabaseObject
	logger := log.NewEntry().WithField("get database object DDL", "comparison sql statement")
	var base *SQLStatement
	var comparison *SQLStatement
	if baseObject != nil {
		baseInst, err := getInstanceById(c, baseObject.InstanceId)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		base, err = getAuditedSQLStatement(c.Request().Context(), logger, baseInst, *baseObject.SchemaName, databaseObject.ObjectName, databaseObject.ObjectType, projectID)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}
	if comparisonObject != nil {
		comparedInst, err := getInstanceById(c, comparisonObject.InstanceId)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		comparison, err = getAuditedSQLStatement(c.Request().Context(), logger, comparedInst, *comparisonObject.SchemaName, databaseObject.ObjectName, databaseObject.ObjectType, projectID)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}
	return c.JSON(http.StatusOK, &DatabaseComparisonStatementsResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: &DatabaseComparisonStatements{
			BaseSQL:        base,
			ComparisondSQL: comparison,
		},
	})
}

func getAuditedSQLStatement(c context.Context, l *logrus.Entry, inst *model.Instance, schemaName, objectName, objectType, projectID string) (*SQLStatement, error) {

	ddl, err := compare.GetDatabaseObjectDDL(c, l, inst, schemaName, objectName, objectType)
	if err != nil {
		return nil, errors.New(errors.DataConflict, err)
	}

	if ddl == nil {
		return nil, nil
	}

	projectUID := model.ProjectUID(projectID)

	task := &model.Task{
		Instance: inst,
		Schema:   schemaName,
		DBType:   inst.DbType,
		ExecuteSQLs: []*model.ExecuteSQL{
			{
				BaseSQL: model.BaseSQL{
					Content: ddl.ObjectDDL,
				},
			},
		},
	}
	var sqlStatement *SQLStatement
	err = server.Audit(l, task, &projectUID, inst.RuleTemplateName)
	if err != nil {
		sqlStatement = &SQLStatement{
			AuditError: err.Error(),
			SQLStatementWithAudit: &SQLStatementWithAuditResult{
				SQLStatement: ddl.ObjectDDL,
			},
		}
		l.Errorf("auditing sql ddl statement error: %v", err)
	} else {
		sqlStatement = &SQLStatement{
			SQLStatementWithAudit: &SQLStatementWithAuditResult{
				AuditResults: convertTaskAuditResults(c, task.ExecuteSQLs[0].AuditResults, task.DBType),
				SQLStatement: ddl.ObjectDDL,
			},
		}
	}

	return sqlStatement, nil
}

func convertTaskAuditResults(c context.Context, taskAuditResults model.AuditResults, dbType string) []*SQLAuditResult {
	auditResults := make([]*SQLAuditResult, len(taskAuditResults))
	for i, auditResult := range taskAuditResults {
		auditResults[i] = &SQLAuditResult{
			Level:               auditResult.Level,
			Message:             auditResult.GetAuditMsgByLangTag(locale.Bundle.GetLangTagFromCtx(c)),
			RuleName:            auditResult.RuleName,
			DbType:              dbType,
			I18nAuditResultInfo: auditResult.I18nAuditResultInfo,
		}
	}
	return auditResults
}

func genDatabaseDiffModifySQLs(c echo.Context) error {
	req := new(GenModifylSQLReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	projectId, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
	}
	// 获取基准数据源和比对数据源实例信息
	baseInst, err := getInstanceById(c, req.BaseInstanceId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	comparedInst, err := getInstanceById(c, req.ComparisonInstanceId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	schemaObjs := make([]*driverV2.DatabasCompareSchemaInfo, len(req.DatabaseSchemaObjects))
	for i, schema := range req.DatabaseSchemaObjects {
		schemaObjs[i] = &driverV2.DatabasCompareSchemaInfo{
			BaseSchemaName:     schema.BaseSchemaName,
			ComparedSchemaName: schema.ComparisonSchemaName,
			DatabaseObjects:    convertDatabaseSchemaObjectToDriver(schema.DatabaseObjects),
		}
	}
	compared := &compare.Compared{
		BaseInstance:     baseInst,
		ComparedInstance: comparedInst,
		ObjInfos:         schemaObjs,
	}

	projectUID := model.ProjectUID(projectId)
	logger := log.NewEntry().WithField("database comparsion", "generate modify sql")
	diffSQLs, err := compared.GetDatabaseDiffModifySQLs(c.Request().Context(), logger)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	plugin, err := common.NewDriverManagerWithoutAudit(logger, comparedInst, "")
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	defer plugin.Close(c.Request().Context())
	comparisonDBSchemas, err := plugin.Schemas(c.Request().Context())
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// 允许在对比实例上生成schema创建语句，schema在该实例上不存在审核时不在task中指定schema
	comparisonDBSchemasMap := make(map[string]struct{}, len(comparisonDBSchemas))
	for _, schema := range comparisonDBSchemas {
		comparisonDBSchemasMap[schema] = struct{}{}
	}
	diffModifySQLs := make([]*DatabaseDiffModifySQL, len(diffSQLs))
	for i, diffSQL := range diffSQLs {
		task := &model.Task{Instance: comparedInst, DBType: comparedInst.DbType}
		if _, exist := comparisonDBSchemasMap[diffSQL.SchemaName]; exist {
			task.Schema = diffSQL.SchemaName
		}
		modifySQLs := make([]*SQLStatementWithAuditResult, len(diffSQL.ModifySQLs))
		for j, sql := range diffSQL.ModifySQLs {
			task.ExecuteSQLs = append(task.ExecuteSQLs, &model.ExecuteSQL{
				BaseSQL: model.BaseSQL{
					Number:  uint(j),
					Content: sql,
				},
			})
			modifySQLs[j] = &SQLStatementWithAuditResult{
				SQLStatement: sql,
			}
		}
		var sqlStatement *DatabaseDiffModifySQL
		err := server.Audit(logger, task, &projectUID, comparedInst.RuleTemplateName)
		if err != nil {
			sqlStatement = &DatabaseDiffModifySQL{
				AuditError: err.Error(),
				SchemaName: diffSQL.SchemaName,
				ModifySQLs: modifySQLs,
			}
			logger.Errorf("auditing modify sql statement error: %v", err)
		} else {
			for j, sql := range task.ExecuteSQLs {
				modifySQLs[j].AuditResults = convertTaskAuditResults(c.Request().Context(), sql.AuditResults, comparedInst.DbType)
			}
			sqlStatement = &DatabaseDiffModifySQL{
				SchemaName: diffSQL.SchemaName,
				ModifySQLs: modifySQLs,
			}
		}
		diffModifySQLs[i] = sqlStatement

	}
	return c.JSON(http.StatusOK, &GenModifySQLResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    diffModifySQLs,
	})
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
