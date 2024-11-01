//go:build enterprise
// +build enterprise

package compare

import (
	"context"
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/common"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/sirupsen/logrus"
)

type Compared struct {
	BaseInstance     *model.Instance
	ComparedInstance *model.Instance
	ObjInfos         []*driverV2.DatabasCompareSchemaInfo
}

type SchemaObject struct {
	BaseSchemaName       string
	ComparisonSchemaName string
	ComparedResult       string
	DatabaseDiffObjects  []*DatabaseDiffObject
	InconsistentNum      int
}

type DatabaseDiffObject struct {
	InconsistentNum    int
	ObjectType         string
	ObjectsDiffResults []*ObjectDiffResult
}

type ObjectDiffResult struct {
	ComparedResult string
	ObjectName     string
}

const (
	DatabaseStructSame               string = "same"                 // 对比一致
	DatabaseStructInconsistent       string = "inconsistent"         // 对比有差异
	DatabaseStructBaseNotExist       string = "base_not_exist"       // 基准对象不存在
	DatabaseStructComparisonNotExist string = "comparison_not_exist" // 比对对象不存在
)

func (c *Compared) ExecDatabaseCompare(context context.Context, l *logrus.Entry) ([]*SchemaObject, error) {
	baseInfos := make([]*driverV2.DatabasSchemaInfo, len(c.ObjInfos))
	comparedInfos := make([]*driverV2.DatabasSchemaInfo, len(c.ObjInfos))
	if len(c.ObjInfos) > 0 {
		for i, objInfo := range c.ObjInfos {
			baseInfos[i] = &driverV2.DatabasSchemaInfo{
				SchemaName: objInfo.BaseSchemaName,
			}
			comparedInfos[i] = &driverV2.DatabasSchemaInfo{
				SchemaName: objInfo.ComparedSchemaName,
			}
		}
	}

	basePlugin, err := common.NewDriverManagerWithoutAudit(l, c.BaseInstance, "")
	if err != nil {
		return nil, err
	}
	defer basePlugin.Close(context)
	baseRes, err := basePlugin.GetDatabaseObjectDDL(context, baseInfos)
	if err != nil {
		return nil, err
	}
	comparedPlugin, err := common.NewDriverManagerWithoutAudit(l, c.ComparedInstance, "")
	if err != nil {
		return nil, err
	}
	defer comparedPlugin.Close(context)
	comparedRes, err := comparedPlugin.GetDatabaseObjectDDL(context, comparedInfos)
	if err != nil {
		return nil, err
	}
	schemaObjects := CompareDDL(baseRes, comparedRes, c.ObjInfos)
	return schemaObjects, err
}

// 根据对象名称对比基准数据源和对比数据源对象的差异
// 1. 当比对数据源中对象在基准数据源对象中不存在时，比对结果为基准对象不存在
// 2. 当基准数据源中对象在比对数据源对象中不存在时，比对结果为比对对象不存在
// 3. 当基准数据源中对象和比对数据源中对象DDL比对一致时，比对结果为比对一致
// 4. 当基准数据源中对象和比对数据源中对象DDL比对有差异时，比对结果为不一致
func CompareDDL(baseRes, comparedRes []*driverV2.DatabaseSchemaObjectResult, originCompareObjects []*driverV2.DatabasCompareSchemaInfo) []*SchemaObject {
	var schemaObjects []*SchemaObject

	baseMap := make(map[string]*driverV2.DatabaseSchemaObjectResult)
	comparedMap := make(map[string]*driverV2.DatabaseSchemaObjectResult)

	for _, base := range baseRes {
		baseMap[base.SchemaName] = base
	}

	for _, compared := range comparedRes {
		comparedMap[compared.SchemaName] = compared
	}

	if len(originCompareObjects) > 0 {
		// 对指定的schema进行对比
		for _, obj := range originCompareObjects {
			schemaObjects = append(schemaObjects, compareSchema(baseMap[obj.BaseSchemaName], comparedMap[obj.ComparedSchemaName]))
		}
	} else {
		// 对指定的数据源进行对比
		// 对基准数据源中的每个 schema 进行比较
		for schemaName := range baseMap {
			schemaObjects = append(schemaObjects, compareSchema(baseMap[schemaName], comparedMap[schemaName]))
		}
		// 对比对数据源中基准数据源没有的 schema 进行比较
		for schemaName := range comparedMap {
			if _, exists := baseMap[schemaName]; !exists {
				schemaObjects = append(schemaObjects, compareSchema(nil, comparedMap[schemaName]))
			}
		}

	}
	return schemaObjects
}

func compareSchema(base, compared *driverV2.DatabaseSchemaObjectResult) *SchemaObject {
	schemaObject := &SchemaObject{}
	// 对基准和对比数据源进行nil判断，为了兼顾两个实例直接对比，可能会出现schema不存在的情况
	if base != nil {
		schemaObject.BaseSchemaName = base.SchemaName
	} else {
		schemaObject.ComparedResult = DatabaseStructBaseNotExist
	}
	if compared != nil {
		schemaObject.ComparisonSchemaName = compared.SchemaName
	} else if base != nil {
		schemaObject.ComparedResult = DatabaseStructComparisonNotExist
	}

	diffObjectMap := make(map[string]*DatabaseDiffObject)

	baseObjMap := make(map[struct{ Name, Type string }]*driverV2.DatabaseObjectDDL)
	comparedObjMap := make(map[struct{ Name, Type string }]*driverV2.DatabaseObjectDDL)

	if base != nil && base.DatabaseObjectDDLs != nil {
		for _, obj := range base.DatabaseObjectDDLs {
			baseObjMap[struct{ Name, Type string }{obj.DatabaseObject.ObjectName, obj.DatabaseObject.ObjectType}] = obj
		}
	}

	if compared != nil && compared.DatabaseObjectDDLs != nil {
		for _, obj := range compared.DatabaseObjectDDLs {
			comparedObjMap[struct{ Name, Type string }{obj.DatabaseObject.ObjectName, obj.DatabaseObject.ObjectType}] = obj
		}
	}

	// 对基准数据源中的每个对象进行比较
	for objName, baseObj := range baseObjMap {
		if comparedObj, exists := comparedObjMap[objName]; exists {
			// 比对数据源中也存在相同的对象，则进行 DDL 对比
			diffObj := compareObjects(baseObj, comparedObj)
			updateDiffObjectMap(diffObjectMap, baseObj.DatabaseObject.ObjectType, diffObj)
		} else {
			// 比对数据源中不存在该对象
			updateDiffObjectMap(diffObjectMap, baseObj.DatabaseObject.ObjectType, &DatabaseDiffObject{
				InconsistentNum: 1,
				ObjectType:      baseObj.DatabaseObject.ObjectType,
				ObjectsDiffResults: []*ObjectDiffResult{
					{
						ComparedResult: DatabaseStructComparisonNotExist,
						ObjectName:     strings.Trim(objName.Name, "`"),
					},
				},
			})
		}
	}

	// 检查比对数据源中有但基准数据源中不存在的对象
	for objName, comparedObj := range comparedObjMap {
		if _, exists := baseObjMap[objName]; !exists {
			updateDiffObjectMap(diffObjectMap, comparedObj.DatabaseObject.ObjectType, &DatabaseDiffObject{
				InconsistentNum: 1,
				ObjectType:      comparedObj.DatabaseObject.ObjectType,
				ObjectsDiffResults: []*ObjectDiffResult{
					{
						ComparedResult: DatabaseStructBaseNotExist,
						ObjectName:     strings.Trim(objName.Name, "`"),
					},
				},
			})
		}
	}

	// 收集最终结果并计算不一致数量
	var diffObjects []*DatabaseDiffObject
	inconsistentCount := 0
	for _, diffObj := range diffObjectMap {
		diffObjects = append(diffObjects, diffObj)
		inconsistentCount += diffObj.InconsistentNum
	}

	// 填充对比结果和不一致的数量
	schemaObject.DatabaseDiffObjects = diffObjects
	schemaObject.InconsistentNum = inconsistentCount
	if schemaObject.ComparedResult == "" {
		if inconsistentCount == 0 {
			schemaObject.ComparedResult = DatabaseStructSame
		} else {
			schemaObject.ComparedResult = DatabaseStructInconsistent
		}
	}

	return schemaObject
}

// 比较两个对象的 DDL，如果相同则返回对比一致，否则返回对比有差异
func compareObjects(baseObj, comparedObj *driverV2.DatabaseObjectDDL) *DatabaseDiffObject {
	diffObject := &DatabaseDiffObject{
		ObjectType: baseObj.DatabaseObject.ObjectType,
	}
	var comparedResult string
	if baseObj.ObjectDDL == comparedObj.ObjectDDL {
		// DDL相同
		comparedResult = DatabaseStructSame
	} else {
		// DDL不同，增加不一致计数
		comparedResult = DatabaseStructInconsistent
		diffObject.InconsistentNum++
	}

	diffObject.ObjectsDiffResults = []*ObjectDiffResult{
		{
			ComparedResult: comparedResult,
			ObjectName:     strings.Trim(baseObj.DatabaseObject.ObjectName, "`"),
		},
	}

	return diffObject
}

// 将比较结果更新到 diffObjectMap 中，并根据对象类型进行聚合
func updateDiffObjectMap(diffObjectMap map[string]*DatabaseDiffObject, objectType string, diffObj *DatabaseDiffObject) {
	if existingDiffObj, exists := diffObjectMap[objectType]; exists {
		// 如果该对象类型已存在于 map 中，累加不一致数量，并追加对比结果
		existingDiffObj.InconsistentNum += diffObj.InconsistentNum
		existingDiffObj.ObjectsDiffResults = append(existingDiffObj.ObjectsDiffResults, diffObj.ObjectsDiffResults...)
	} else {
		// 如果不存在，则直接将当前对比结果加入 map
		diffObjectMap[objectType] = diffObj
	}
}

func (c *Compared) GetDatabaseDiffModifySQLs(context context.Context, l *logrus.Entry) ([]*driverV2.DatabaseDiffModifySQLResult, error) {
	basePlugin, err := common.NewDriverManagerWithoutAudit(l, c.BaseInstance, "")
	if err != nil {
		return nil, err
	}
	defer basePlugin.Close(context)
	calibratedDSN := &driverV2.DSN{
		Host:             c.ComparedInstance.Host,
		Port:             c.ComparedInstance.Port,
		User:             c.ComparedInstance.User,
		Password:         c.ComparedInstance.Password,
		AdditionalParams: c.ComparedInstance.AdditionalParams,
	}
	modifySQLs, err := basePlugin.GetDatabaseDiffModifySQL(context, calibratedDSN, c.ObjInfos)
	if err != nil {
		return nil, err
	}
	dbDiffSQLs := make([]*driverV2.DatabaseDiffModifySQLResult, len(modifySQLs))
	for i, schemaDiff := range modifySQLs {
		dbDiffSQLs[i] = &driverV2.DatabaseDiffModifySQLResult{
			SchemaName: schemaDiff.SchemaName,
			ModifySQLs: schemaDiff.ModifySQLs,
		}
	}
	return dbDiffSQLs, nil
}

func GetDatabaseObjectDDL(c context.Context, l *logrus.Entry, instance *model.Instance, schemaName, objectName, objectType string) (*driverV2.DatabaseObjectDDL, error) {
	p, err := common.NewDriverManagerWithoutAudit(l, instance, "")
	if err != nil {
		return nil, err
	}
	defer p.Close(c)
	baseInfos := make([]*driverV2.DatabasSchemaInfo, 0, 1)
	dbObj := &driverV2.DatabaseObject{
		ObjectName: objectName,
		ObjectType: objectType,
	}
	baseInfos = append(baseInfos, &driverV2.DatabasSchemaInfo{
		SchemaName:      schemaName,
		DatabaseObjects: []*driverV2.DatabaseObject{dbObj},
	})
	objDDLRes, err := p.GetDatabaseObjectDDL(c, baseInfos)
	if err != nil {
		return nil, err
	}
	// 因为该方法只获取一个数据库对象的ddl语句，所以得到的结果也应该只有1条
	if !(len(objDDLRes) == 1 && len(objDDLRes[0].DatabaseObjectDDLs) == 1) {
		return nil, fmt.Errorf("the number of ddl statements to get the database is not as expected")
	}
	return &driverV2.DatabaseObjectDDL{
		DatabaseObject: &driverV2.DatabaseObject{
			ObjectName: objDDLRes[0].DatabaseObjectDDLs[0].DatabaseObject.ObjectName,
			ObjectType: objDDLRes[0].DatabaseObjectDDLs[0].DatabaseObject.ObjectType,
		},
		ObjectDDL: objDDLRes[0].DatabaseObjectDDLs[0].ObjectDDL,
	}, nil
}
