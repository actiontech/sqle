//go:build enterprise
// +build enterprise

package mysql

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/driver/mysql/differ"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/utils"
)

func addOptionModules(metas *driverV2.DriverMetas) {
	metas.EnabledOptionalModule = append(metas.EnabledOptionalModule, driverV2.OptionalBackup)
}

// func (*MysqlDriverImpl) QueryPrepare(ctx context.Context, sql string, conf *driver.QueryPrepareConf) (*driver.QueryPrepareResult, error) {
// 	return QueryPrepare(ctx, sql, conf)
// }

// func QueryPrepare(ctx context.Context, sql string, conf *driverV2.QueryPrepareConf) (*driverV2.QueryPrepareResult, error) {
// 	node, err := sqlparser.Parse(sql)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// check is query sql
// 	stmt, ok := node.(*sqlparser.Select)
// 	if !ok {
// 		return &driver.QueryPrepareResult{
// 			ErrorType: driver.ErrorTypeNotQuery,
// 			Error:     driver.ErrorTypeNotQuery,
// 		}, nil
// 	}

// 	// Generate new limit
// 	limit, offset := -1, 0
// 	if stmt.Limit != nil {
// 		if stmt.Limit.Rowcount != nil {
// 			limit, _ = strconv.Atoi(stmt.Limit.Rowcount.(*sqlparser.Literal).Val)
// 		}
// 		if stmt.Limit.Offset != nil {
// 			offset, _ = strconv.Atoi(stmt.Limit.Offset.(*sqlparser.Literal).Val)
// 		} else if limit != -1 {
// 			offset = 0
// 		}
// 	}
// 	appendLimit, appendOffset := -1, -1
// 	if conf != nil {
// 		appendLimit, appendOffset = int(conf.Limit), int(conf.Offset)
// 	}
// 	if appendLimit != -1 && appendOffset == -1 {
// 		appendLimit = 0
// 	}

// 	newLimit, newOffset := CalculateOffset(limit, offset, appendLimit, appendOffset)

// 	if newLimit != -1 {
// 		l := &sqlparser.Limit{
// 			Offset: &sqlparser.Literal{
// 				Type: sqlparser.IntVal,
// 				Val:  strconv.Itoa(newOffset),
// 			},
// 			Rowcount: &sqlparser.Literal{
// 				Type: sqlparser.IntVal,
// 				Val:  strconv.Itoa(newLimit),
// 			},
// 		}
// 		stmt.SetLimit(l)
// 	}

// 	// rewrite
// 	return &driver.QueryPrepareResult{
// 		NewSQL:    sqlparser.String(stmt),
// 		ErrorType: driver.ErrorTypeNotError,
// 	}, nil
// }

// // 1 means this item has no value or no limit
// func CalculateOffset(oldLimit, oldOffset, appendLimit, appendOffset int) (newLimit, newOffset int) {
// 	if checkIsInvalidCalculateOffset(oldLimit, oldOffset, appendLimit, appendOffset) {
// 		return oldLimit, oldOffset
// 	}
// 	return calculateOffset(oldLimit, oldOffset, appendLimit, appendOffset)
// }

// func calculateOffset(oldLimit, oldOffset, appendLimit, appendOffset int) (newLimit, newOffset int) {
// 	if oldLimit == -1 {
// 		return appendLimit, appendOffset
// 	}
// 	newOffset = oldOffset + appendOffset
// 	newLimit = appendLimit
// 	if newOffset+newLimit > oldLimit+oldOffset {
// 		newLimit = oldLimit - appendOffset
// 	}

// 	return newLimit, newOffset
// }

// func checkIsInvalidCalculateOffset(oldLimit, oldOffset, appendLimit, appendOffset int) bool {
// 	if appendLimit == -1 {
// 		return true
// 	}
// 	if oldLimit != -1 && appendOffset > oldLimit+oldOffset {
// 		return true
// 	}

// 	return false
// }

func (i *MysqlDriverImpl) Query(ctx context.Context, sql string, conf *driverV2.QueryConf) (*driverV2.QueryResult, error) {
	if i.IsOfflineAudit() {
		return nil, nil
	}
	conn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}
	// add timeout
	cancel := func() {}
	if conf != nil && conf.TimeOutSecond > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(conf.TimeOutSecond)*time.Second)
		defer cancel()
	}

	columns, rows, err := conn.Db.QueryWithContext(ctx, sql)
	if err != nil {
		return nil, err
	}

	// generate result
	res := &driverV2.QueryResult{
		Column: params.Params{},
		Rows:   []*driverV2.QueryResultRow{},
	}
	for _, column := range columns {
		res.Column = append(res.Column, &params.Param{
			Key:   column,
			Value: column,
		})
	}
	for _, row := range rows {
		r := &driverV2.QueryResultRow{
			Values: []*driverV2.QueryResultValue{},
		}
		for _, s := range row {
			r.Values = append(r.Values, &driverV2.QueryResultValue{
				Value: s.String,
			})
		}
		res.Rows = append(res.Rows, r)
	}
	return res, nil
}

// 由基准数据源(driver中数据源)和比对数据源(入参)所指定的schema对象信息生成变更sql，修改比对数据库对象的差异
// 生成变更sql大致逻辑是：
// 1、获取基准和比对schema对象的解析结果
// 2、比对schema对象解析结果的差异
// 3、根据差异生成变更sql
func (i *MysqlDriverImpl) GetDatabaseDiffModifySQL(ctx context.Context, calibratedDSN *driverV2.DSN, objInfos []*driverV2.DatabasCompareSchemaInfo) ([]*driverV2.DatabaseDiffModifySQLResult, error) {
	baseConn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}
	logEntry := log.NewEntry().WithField("generate sql", "generate sql based on database differences")
	compareConn, err := executor.NewExecutor(logEntry, &driverV2.DSN{
		Host:             calibratedDSN.Host,
		Port:             calibratedDSN.Port,
		User:             calibratedDSN.User,
		Password:         calibratedDSN.Password,
		AdditionalParams: calibratedDSN.AdditionalParams,
	}, "")
	if err != nil {
		return nil, err
	}
	defer compareConn.Db.Close()
	compareSchemaInfos := make([]*driverV2.DatabasSchemaInfo, len(objInfos))
	baseSchemaInfos := make([]*driverV2.DatabasSchemaInfo, len(objInfos))
	for i, objInfo := range objInfos {
		compareSchemaInfos[i] = &driverV2.DatabasSchemaInfo{
			SchemaName:      objInfo.ComparedSchemaName,
			DatabaseObjects: objInfo.DatabaseObjects,
		}
		baseSchemaInfos[i] = &driverV2.DatabasSchemaInfo{
			SchemaName:      objInfo.BaseSchemaName,
			DatabaseObjects: objInfo.DatabaseObjects,
		}
	}

	compareSchemas, err := differ.Schemas(true, compareConn, compareSchemaInfos)
	if err != nil {
		return nil, err
	}
	baseSchemas, err := differ.Schemas(true, baseConn, baseSchemaInfos)
	if err != nil {
		return nil, err
	}
	compareSchemaMap := make(map[string]*differ.Schema)
	baseSchemaMap := make(map[string]*differ.Schema)

	for _, compare := range compareSchemas {
		compareSchemaMap[compare.Name] = compare
	}
	for _, base := range baseSchemas {
		baseSchemaMap[base.Name] = base
	}

	cpAllSchemas, err := compareConn.ShowDatabases(true)
	if err != nil {
		return nil, err
	}
	cpAllSchemaMap := make(map[string]struct{}, len(cpAllSchemas))
	for _, schemaName := range cpAllSchemas {
		cpAllSchemaMap[schemaName] = struct{}{}
	}

	mods := differ.StatementModifiers{
		AllowUnsafe:            true, // 允许生成不安全语句，如不检查是否有内容的情况下生成drop table、drop database语句等
		StrictIndexOrder:       true, // 强制保持索引顺序一致
		StrictForeignKeyNaming: true, // 强制保持外键定义一致
	}
	// LockClause：
	// 该字段在生成的 ALTER TABLE 语句中用于指定锁定的行为。通过设置 LOCK=[value] 来决定当修改表时是否锁定表以及锁定的方式。
	// 例如，LOCK=none 表示在 ALTER TABLE 操作时不锁定任何内容。通常用于需要在不阻塞表读写的情况下修改表结构的场景。
	// AlgorithmClause：
	// 该字段用于生成 ALTER TABLE 语句时指定修改表的算法。ALGORITHM=[value] 用于选择更改表结构时使用的算法。
	// 例如，ALGORITHM=inplace 表示不需要拷贝整个表，而是在不重建整个表的情况下就地修改表结构。这种算法更高效，尤其适用于大表的变更。
	// 此处属于优化内容，且mysql各个版本支持度不唯一，此处可以不使用该属性
	// results, err := compareConn.Db.Query("SELECT VERSION()")
	// if err != nil {
	// 	return nil, err
	// }
	// flavor := differ.ParseFlavor(fmt.Sprintf("%s %s", "mysql:", results[0]["VERSION()"].String))
	// if !flavor.IsMySQL(5, 5) {
	// 	mods.LockClause, mods.AlgorithmClause = "none", "inplace"
	// }
	dbDiffSQLs := make([]*driverV2.DatabaseDiffModifySQLResult, 0)
	for _, objInfo := range objInfos {
		compareSchema := compareSchemaMap[objInfo.ComparedSchemaName]
		baseSchema := baseSchemaMap[objInfo.BaseSchemaName]
		// compare作为需要校对的schema，base为基准schema，另一种写法为diff.NewSchemaDiff(compareSchema, baseSchema)
		diff := compareSchema.Diff(baseSchema)
		objDiffs := diff.ObjectDiffs()
		diffSqls := make([]string, 0, len(objDiffs))

		// 当数据库存在于校准数据源时，在生成变更sql前，首先use database
		if _, ok := cpAllSchemaMap[objInfo.ComparedSchemaName]; ok {
			diffSqls = append(diffSqls, fmt.Sprintf("USE %s;\n", objInfo.ComparedSchemaName))
		}

		for _, objDiff := range objDiffs {
			stmt, err := objDiff.Statement(mods)
			if err != nil {
				return nil, err
			}
			if stmt != "" {
				diffSqls = append(diffSqls, fmt.Sprintf("%s;\n", stmt))
			}
			// 当sql语句为create database时，需要补充在其后use database;
			if objDiff.ObjectKey().Type == differ.ObjectTypeDatabase && objDiff.DiffType() == differ.DiffTypeCreate {
				diffSqls = append(diffSqls, fmt.Sprintf("USE %s;\n", objDiff.ObjectKey().Name))
			}

		}

		dbDiffSQLs = append(dbDiffSQLs, &driverV2.DatabaseDiffModifySQLResult{
			SchemaName: objInfo.ComparedSchemaName, // 此处使用校准数据源的schema，因为变更sql需要在校准数据源上执行
			ModifySQLs: diffSqls,
		})
	}

	return dbDiffSQLs, nil
}

func (i *MysqlDriverImpl) GetDatabaseObjectDDL(ctx context.Context, objInfos []*driverV2.DatabasSchemaInfo) ([]*driverV2.DatabaseSchemaObjectResult, error) {
	conn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}
	// 复制切片以避免修改原始数据
	databasSchemaInfo := make([]*driverV2.DatabasSchemaInfo, len(objInfos))
	copy(databasSchemaInfo, objInfos)
	// 未指定schema时默认获取全量
	if len(databasSchemaInfo) == 0 {
		schemas, err := conn.ShowDatabases(true)
		if err != nil {
			return nil, err
		}
		for _, schema := range schemas {
			schemaObjs, err := getAllSchemaObjects(conn, schema)
			if err != nil {
				return nil, err
			}
			databasSchemaInfo = append(databasSchemaInfo, &driverV2.DatabasSchemaInfo{
				SchemaName:      schema,
				DatabaseObjects: schemaObjs,
			})

		}
	} else {
		// 未指定schema对象（表、视图、存储过程等）时默认获取全量
		for _, objInfo := range databasSchemaInfo {
			if len(objInfo.DatabaseObjects) == 0 {
				schemaObjs, err := getAllSchemaObjects(conn, objInfo.SchemaName)
				if err != nil {
					return nil, err
				}
				objInfo.DatabaseObjects = schemaObjs
			}
		}
	}

	res := make([]*driverV2.DatabaseSchemaObjectResult, 0, len(databasSchemaInfo))
	for _, objInfo := range databasSchemaInfo {
		err = conn.UseSchema(objInfo.SchemaName)
		if err != nil {
			return nil, fmt.Errorf("use schema fail, error: %v", err)
		}
		schemaDDL, err := conn.ShowCreateSchema(utils.SupplementalQuotationMarks(objInfo.SchemaName))
		if err != nil {
			return nil, err
		}
		dbDDLs := make([]*driverV2.DatabaseObjectDDL, 0, len(objInfo.DatabaseObjects))
		for _, obj := range objInfo.DatabaseObjects {
			objetDDL := ""
			if obj.ObjectType == driverV2.ObjectType_TABLE {
				tableDDL, err := conn.ShowCreateTable(objInfo.SchemaName, obj.ObjectName)
				engine, err := conn.GetTableEngine(objInfo.SchemaName, strings.Trim(obj.ObjectName, "`"))
				if err != nil {
					return nil, err
				}
				objetDDL = sanitizeTableDDL(tableDDL, engine)
			}
			if obj.ObjectType == driverV2.ObjectType_PROCEDURE {
				objetDDL, err = conn.ShowCreateProcedure(obj.ObjectName)
				if err != nil {
					return nil, err
				}
			}
			if obj.ObjectType == driverV2.ObjectType_FUNCTION {
				objetDDL, err = conn.ShowCreateFunction(obj.ObjectName)
				if err != nil {
					return nil, err
				}
			}

			// TODO 数据库结构对比目前视图、事件、触发器还未支持，若后续支持在此处补充获取这些对象的ddl（代码须经验证）
			// if obj.ObjectType == driverV2.ObjectType_VIEW {
			// 	objetDDL, err = conn.ShowCreateView(obj.ObjectName)
			// 	if err != nil {
			// 		return nil, err
			// 	}
			// }

			dbDDLs = append(dbDDLs, &driverV2.DatabaseObjectDDL{
				ObjectDDL: objetDDL,
				DatabaseObject: &driverV2.DatabaseObject{
					ObjectName: obj.ObjectName,
					ObjectType: obj.ObjectType,
				},
			})
		}
		res = append(res, &driverV2.DatabaseSchemaObjectResult{
			SchemaName:         objInfo.SchemaName,
			SchemaDDL:          schemaDDL,
			DatabaseObjectDDLs: dbDDLs,
		})
	}

	return res, nil
}

// 去除DDL中不需要关注的token
// 在建表的DDL语句中，表选项AUTO_INCREMENT表示自增起点为
// 当表中存有数据时，已有记录数量 + 1，这个值根据表中记录数量所变化，获取ddl后续处理时一般不关心自增起点
func sanitizeTableDDL(ddl string, tableEngine string) string {
	re := regexp.MustCompile(`AUTO_INCREMENT\s*=\s*\d+\s`)
	result := re.ReplaceAllString(ddl, "")
	if tableEngine == "InnoDB" {
		// 删除不影响 InnoDB 的创建选项
		result = differ.NormalizeCreateOptions(result)
	}
	return result
}

func getAllSchemaObjects(conn *executor.Executor, schemaName string) ([]*driverV2.DatabaseObject, error) {
	dbObjs := make([]*driverV2.DatabaseObject, 0)
	tables, err := conn.ShowSchemaTables(schemaName)
	if err != nil {
		return nil, err
	}
	for _, tableName := range tables {
		dbObjs = append(dbObjs, &driverV2.DatabaseObject{
			ObjectName: utils.SupplementalQuotationMarks(tableName),
			ObjectType: driverV2.ObjectType_TABLE,
		})
	}

	procedures, err := conn.ShowSchemaProcedures(schemaName)
	if err != nil {
		return nil, err
	}
	for _, procedureName := range procedures {
		dbObjs = append(dbObjs, &driverV2.DatabaseObject{
			ObjectName: utils.SupplementalQuotationMarks(procedureName),
			ObjectType: driverV2.ObjectType_PROCEDURE,
		})
	}

	functions, err := conn.ShowSchemaFunctions(schemaName)
	if err != nil {
		return nil, err
	}
	for _, functionName := range functions {
		dbObjs = append(dbObjs, &driverV2.DatabaseObject{
			ObjectName: utils.SupplementalQuotationMarks(functionName),
			ObjectType: driverV2.ObjectType_FUNCTION,
		})
	}

	// TODO 数据库结构对比目前视图、事件、触发器还未支持，若后续支持放开此处代码注释以获取这些对象的名称（代码须经验证）
	// views, err := conn.ShowSchemaViews(schemaName)
	// if err != nil {
	// 	return nil, err
	// }
	// for _, viewName := range views {
	// 	dbObjs = append(dbObjs, &driverV2.DatabaseObject{
	// 		ObjectName: utils.SupplementalQuotationMarks(viewName),
	// 		ObjectType: driverV2.ObjectType_VIEW,
	// 	})
	// }
	// 	triggers, err := conn.ShowSchemaTriggers(schemaName)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	for _, triggerName := range triggers {
	// 		dbObjs = append(dbObjs, &driverV2.DatabaseObject{
	// 			ObjectName: utils.SupplementalQuotationMarks(triggerName),
	// 			ObjectType: driverV2.ObjectType_TRIGGER,
	// 		})
	// 	}
	//
	// 	events, err := conn.ShowSchemaEvents(schemaName)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	for _, eventName := range events {
	// 		dbObjs = append(dbObjs, &driverV2.DatabaseObject{
	// 			ObjectName: utils.SupplementalQuotationMarks(eventName),
	// 			ObjectType: driverV2.ObjectType_EVENT,
	// 		})
	// 	}

	return dbObjs, nil
}
