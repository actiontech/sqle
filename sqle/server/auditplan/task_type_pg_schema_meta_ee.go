//go:build enterprise
// +build enterprise

package auditplan

import (
	"context"
	"fmt"
	"time"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/sirupsen/logrus"
)

type PGSchemaMetaTaskV2 struct {
	BaseSchemaMetaTaskV2
}

func NewPGSchemaMetaTaskV2Fn() func() interface{} {
	return func() interface{} {
		return NewPGSchemaMetaTaskV2()
	}
}

func NewPGSchemaMetaTaskV2() *PGSchemaMetaTaskV2 {
	t := &PGSchemaMetaTaskV2{}
	t.BaseSchemaMetaTaskV2 = BaseSchemaMetaTaskV2{
		extract: t.extractSQL,
	}
	return t
}

func (at *PGSchemaMetaTaskV2) InstanceType() string {
	return InstanceTypePostgreSQL
}

func (at *PGSchemaMetaTaskV2) extractSQL(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) ([]*SchemaMetaSQL, error) {
	if ap.InstanceID == "" {
		return nil, fmt.Errorf("instance is not configured")
	}
	instance, _, err := dms.GetInstancesById(context.Background(), ap.InstanceID)
	if err != nil {
		return nil, fmt.Errorf("get instance fail, error: %v", err)
	}

	pluginMgr := driver.GetPluginManager()
	if !pluginMgr.IsOptionalModuleEnabled(instance.DbType, driverV2.OptionalModuleQuery) {
		return nil, fmt.Errorf("collect pg schema meta failed: %v", driver.NewErrPluginAPINotImplement(driverV2.OptionalModuleQuery))
	}
	plugin, err := pluginMgr.OpenPlugin(logger, instance.DbType, &driverV2.Config{DSN: &driverV2.DSN{
		Host:             instance.Host,
		Port:             instance.Port,
		User:             instance.User,
		Password:         instance.Password,
		AdditionalParams: instance.AdditionalParams,
	}})
	if err != nil {
		return nil, fmt.Errorf("connect to instance fail, error: %v", err)
	}
	defer pluginMgr.Stop()

	databases, err := plugin.Schemas(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get schemas from db2 failed, error: %v", err)
	}
	sqls := []*SchemaMetaSQL{}

	for _, db := range databases {
		// 获取所有的数据库及对应的schema
		schemas, err := at.GetAllUserSchemas(logger, plugin, db)
		if err != nil {
			return nil, fmt.Errorf("get databases=%s schemas fail: %s", db, err)
		}

		// 获取表和视图
		tablesAndViews, err := at.GetAllTablesAndViewsForPg(logger, plugin, db)
		if err != nil {
			return nil, fmt.Errorf("get all table and view fail, error: %s", err)
		}

		// 获取列信息
		columnsInfo, err := at.GetAllColumnsInfoForPg(logger, plugin, db)
		if err != nil {
			return nil, fmt.Errorf("get all columns information fail, error: %s", err)
		}

		// 获取约束信息
		constraints, err := at.GetAllConstraintsForPg(logger, plugin)
		if err != nil {
			return nil, fmt.Errorf("get all constraints fail, error: %s", err)
		}

		// 获取索引信息
		indexes, err := at.GetAllIndexesForPg(logger, plugin)
		if err != nil {
			return nil, fmt.Errorf("get all indexes fail, error: %s", err)
		}

		// 是否收集视图sql
		isCollectView := false
		if ap.Params.GetParam("collect_view").Bool() {
			isCollectView = true
		}

		var viewsSql []*PostgreSQLViewInfo
		if isCollectView {
			// 获取视图创建sql
			viewsSql, err = at.GetAllViewsSqlForPg(logger, plugin, db)
			if err != nil {
				return nil, fmt.Errorf("get all views sql fail, error: %s", err)
			}
		}

		for _, schema := range schemas {
			currentSchema := schema.schemaName
			for _, tableOrView := range tablesAndViews {
				if tableOrView.schemaName != currentSchema {
					continue
				}
				dataObjectType := ""
				tableOrViewName := tableOrView.tableName
				if tableOrView.tableType == "BASE TABLE" || tableOrView.tableType == "SYSTEM VIEW" { // 表
					dataObjectType = "table"
				} else if tableOrView.tableType == "VIEW" { // 视图
					dataObjectType = "view"
				}
				if dataObjectType == "table" {
					createDDL := at.createTableSqlForPg(currentSchema, tableOrViewName, columnsInfo, constraints, indexes)
					tableDDl := fmt.Sprintf("%s;%s", createDDL.createTableSql, createDDL.createIndexSql)

					sqls = append(sqls, &SchemaMetaSQL{
						SchemaName: db, // pg 的database 与MySQL schema一致
						MetaName:   tableOrViewName,
						MetaType:   dataObjectType,
						SQLContent: tableDDl,
					})

				} else if dataObjectType == "view" {
					if !isCollectView {
						continue
					}
					// 视图sql
					for _, view := range viewsSql {
						schemaName, tableName := view.schemaName, view.tableName
						if schemaName != schema.schemaName || tableName != tableOrViewName {
							continue
						}
						sqls = append(sqls, &SchemaMetaSQL{
							SchemaName: db,
							MetaName:   tableOrViewName,
							MetaType:   dataObjectType,
							SQLContent: view.viewSql,
						})
					}
				}
			}
		}
	}
	return sqls, nil
}

type PostgreSQLSchema struct {
	schemaName string
}

type PostgreSQLTablesAndViews struct {
	schemaName string
	tableName  string
	tableType  string
}

type PostgreSQLViewInfo struct {
	schemaName string
	tableName  string
	viewSql    string
}

type PostgreSQLCreateTableSql struct {
	createTableSql string
	createIndexSql string
}

type PostgreSQLTableColumnInfo struct {
	schemaName string
	tableName  string
	columnSql  string
}

type PostgreSQLConstraint struct {
	schemaName           string
	tableName            string
	constraintDefinition string
}

type PostgreSQLIndex struct {
	schemaName      string
	tableName       string
	indexName       string
	indexDefinition string
}

func (at *PGSchemaMetaTaskV2) createTableSqlForPg(schema, tableOrViewName string, columnsInfo []*PostgreSQLTableColumnInfo, constraints []*PostgreSQLConstraint, indexes []*PostgreSQLIndex) *PostgreSQLCreateTableSql {
	tableDDl := fmt.Sprintf("CREATE TABLE %s.%s(", schema, tableOrViewName)
	// 列信息
	for _, columnInfo := range columnsInfo {
		schemaName, tableName := columnInfo.schemaName, columnInfo.tableName
		if schemaName != schema || tableName != tableOrViewName {
			continue
		}
		tableDDl += columnInfo.columnSql
	}
	// 约束信息
	for _, constraintInfo := range constraints {
		schemaName, tableName := constraintInfo.schemaName, constraintInfo.tableName
		if schemaName != schema || tableName != tableOrViewName {
			continue
		}
		tableDDl += fmt.Sprintf(",%s\n", constraintInfo.constraintDefinition)
	}
	tableDDl += ")"
	indexDDl := ""
	// 索引信息
	for _, indexInfo := range indexes {
		schemaName, tableName := indexInfo.schemaName, indexInfo.tableName
		if schemaName != schema || tableName != tableOrViewName {
			continue
		}
		indexDDl += fmt.Sprintf("%s;\n", indexInfo.indexDefinition)
	}
	return &PostgreSQLCreateTableSql{
		createTableSql: tableDDl,
		createIndexSql: indexDDl,
	}
}

func (at *PGSchemaMetaTaskV2) GetAllUserSchemas(logger *logrus.Entry, plugin driver.Plugin, database string) ([]*PostgreSQLSchema, error) {
	result := make([]*PostgreSQLSchema, 0)
	querySql := fmt.Sprintf(`
        SELECT schema_name FROM information_schema.schemata WHERE catalog_name = '%s'
		AND schema_name NOT LIKE 'pg_%%' AND schema_name != 'information_schema' ORDER BY schema_name`, database)
	res, err := at.GetResult(logger, plugin, querySql)
	if err != nil {
		return result, err
	}
	for _, value := range res {
		if len(value) == 0 {
			continue
		}
		result = append(result, &PostgreSQLSchema{value[0]})
	}
	if len(result) == 0 {
		return result, fmt.Errorf("database=%s has no schema", database)
	}
	return result, nil
}

func (at *PGSchemaMetaTaskV2) GetAllTablesAndViewsForPg(logger *logrus.Entry, plugin driver.Plugin, database string) ([]*PostgreSQLTablesAndViews, error) {
	querySql := fmt.Sprintf(`select table_schema, table_name, table_type from information_schema.tables 
		where table_catalog = '%s' and table_schema not like 'pg_%%' AND table_schema != 'information_schema'
		ORDER BY table_name`, database)
	result := make([]*PostgreSQLTablesAndViews, 0)
	ret, err := at.GetResult(logger, plugin, querySql)
	if err != nil {
		return result, err
	}
	for _, value := range ret {
		if len(value) < 3 {
			return result, fmt.Errorf("get tables and views error, column length is not three")
		}
		result = append(result, &PostgreSQLTablesAndViews{
			schemaName: value[0],
			tableName:  value[1],
			tableType:  value[2],
		})
	}
	return result, nil
}

func (at *PGSchemaMetaTaskV2) GetAllColumnsInfoForPg(logger *logrus.Entry, plugin driver.Plugin, database string) ([]*PostgreSQLTableColumnInfo, error) {
	columns := fmt.Sprintf(`select table_schema, table_name,
		string_agg(
			concat(
				column_name, ' ', 
				case 
					when lower(data_type) in ('char', 'varchar', 'character', 'character varying') then concat(data_type, '(', coalesce(character_maximum_length, 0), ')') 
					when lower(data_type) in ('numeric', 'decimal') then concat(data_type, '(', coalesce(numeric_precision, 0), ',', coalesce(numeric_scale, 0), ')') 
					when lower(data_type) in ('integer', 'smallint', 'bigint', 'text') then data_type
					else udt_name 
				end,
				case 
					when column_default != '' then concat(' default ', column_default) else '' 
				end,
				case 
					when is_nullable = 'no' then ' not null' else '' 
				end
			), ', ' order by ordinal_position
		) as columns_sql
	from information_schema.columns 
	where table_catalog = '%s' and table_schema not like 'pg_%%' and table_schema != 'information_schema' 
	group by table_schema, table_name`, database)
	result := make([]*PostgreSQLTableColumnInfo, 0)
	ret, err := at.GetResult(logger, plugin, columns)
	if err != nil {
		return result, err
	}
	for _, value := range ret {
		if len(value) < 3 {
			return result, fmt.Errorf("get column info error, column length is not three")
		}
		result = append(result, &PostgreSQLTableColumnInfo{
			schemaName: value[0],
			tableName:  value[1],
			columnSql:  value[2],
		})
	}
	return result, nil
}

func (at *PGSchemaMetaTaskV2) GetAllConstraintsForPg(logger *logrus.Entry, plugin driver.Plugin) ([]*PostgreSQLConstraint, error) {
	querySql := `select
		n.nspname as schema_name,
		c.relname as table_name,
		concat ( 'constraint ', r.conname, ' ', pg_catalog.pg_get_constraintdef ( r.oid, true ) ) as constraint_definition 
	from
		pg_catalog.pg_constraint r
		join pg_catalog.pg_class c on c.oid = r.conrelid
		join pg_catalog.pg_namespace n on n.oid = c.relnamespace 
	where
		n.nspname not like'pg_%' 
		and n.nspname != 'information_schema'`
	result := make([]*PostgreSQLConstraint, 0)
	ret, err := at.GetResult(logger, plugin, querySql)
	if err != nil {
		return result, err
	}
	for _, value := range ret {
		if len(value) < 3 {
			return result, fmt.Errorf("get constraint error, column length is not three")
		}
		result = append(result, &PostgreSQLConstraint{
			schemaName:           value[0],
			tableName:            value[1],
			constraintDefinition: value[2],
		})
	}
	return result, nil
}

func (at *PGSchemaMetaTaskV2) GetAllIndexesForPg(logger *logrus.Entry, plugin driver.Plugin) ([]*PostgreSQLIndex, error) {
	querySql := `select schemaname, tablename, indexname, indexdef from pg_indexes
		         where schemaname not like 'pg_%' and schemaname != 'information_schema'`
	result := make([]*PostgreSQLIndex, 0)
	ret, err := at.GetResult(logger, plugin, querySql)
	if err != nil {
		return result, err
	}
	for _, value := range ret {
		if len(value) < 4 {
			return result, fmt.Errorf("get index error, column length is not four")
		}
		result = append(result, &PostgreSQLIndex{
			schemaName:      value[0],
			tableName:       value[1],
			indexName:       value[2],
			indexDefinition: value[3],
		})
	}
	return result, nil
}

func (at *PGSchemaMetaTaskV2) GetAllViewsSqlForPg(logger *logrus.Entry, plugin driver.Plugin, database string) ([]*PostgreSQLViewInfo, error) {
	querySql := fmt.Sprintf(`select table_schema, table_name, 
		concat('create or replace view ', table_schema, '.', table_name, ' as ', view_definition) as create_view_statement 
	from information_schema.views 
	where table_catalog = '%s' 
		and table_schema not like 'pg_%%' 
		and table_schema != 'information_schema' 
	order by table_name`, database)
	result := make([]*PostgreSQLViewInfo, 0)
	ret, err := at.GetResult(logger, plugin, querySql)
	if err != nil {
		return result, err
	}
	for _, value := range ret {
		if len(value) < 3 {
			return result, fmt.Errorf("get view sql error, column length is not three")
		}
		result = append(result, &PostgreSQLViewInfo{
			schemaName: value[0],
			tableName:  value[1],
			viewSql:    value[2],
		})
	}
	return result, nil
}

func (at *PGSchemaMetaTaskV2) GetResult(logger *logrus.Entry, plugin driver.Plugin, sql string) ([][]string, error) {
	var ret [][]string
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	result, err := plugin.Query(ctx, sql, &driverV2.QueryConf{TimeOutSecond: 120})
	if err != nil {
		logger.Errorf("plugin.Query() failed:%s\n", err)
		return nil, err
	}
	rows := result.Rows
	for _, row := range rows {
		values := row.Values
		if len(values) == 0 {
			continue
		}
		var valueArr []string
		for _, value := range values {
			valueArr = append(valueArr, value.Value)
		}
		ret = append(ret, valueArr)
	}
	return ret, nil
}
