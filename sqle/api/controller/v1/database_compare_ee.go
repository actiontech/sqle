//go:build enterprise
// +build enterprise

package v1

import (
	"fmt"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	dms "github.com/actiontech/sqle/sqle/dms"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server/compare"
	"github.com/labstack/echo/v4"
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
	execCompareRes, err := compared.ExecDatabaseCompare()
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

	return c.JSON(http.StatusOK, &DatabaseComparisonStatementsResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    nil,
	})
}

func genDatabaseDiffModifySQLs(c echo.Context) error {

	return c.JSON(http.StatusOK, &GenModifySQLResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    nil,
	})
}
