//go:build enterprise
// +build enterprise

package auditplan

import (
	"context"
	"fmt"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/sirupsen/logrus"
)

type DB2SchemaMetaTaskV2 struct {
	BaseSchemaMetaTaskV2
}

func NewDB2SchemaMetaTaskV2Fn() func() interface{} {
	return func() interface{} {
		t := &DB2SchemaMetaTaskV2{}
		t.BaseSchemaMetaTaskV2 = BaseSchemaMetaTaskV2{
			extract: t.extractSQL,
		}
		return t
	}
}

func (at *DB2SchemaMetaTaskV2) InstanceType() string {
	return InstanceTypeDB2
}

func (at *DB2SchemaMetaTaskV2) getTablesFromSchema(ctx context.Context, plugin driver.Plugin, schema string) (tables []string, err error) {
	res, err := plugin.Query(ctx, fmt.Sprintf(`SELECT TABNAME FROM SYSCAT.TABLES WHERE TABSCHEMA = '%v' AND TYPE = 'T'`, schema), &driverV2.QueryConf{TimeOutSecond: 10})
	if err != nil {
		return nil, fmt.Errorf("query sql failed: %v", err)
	}
	if len(res.Column) != 1 || res.Column[0].Key != "TABNAME" {
		return nil, fmt.Errorf("parse query results failed")
	}
	for _, row := range res.Rows {
		for _, value := range row.Values {
			tables = append(tables, value.Value)
			continue
		}
	}
	return tables, nil
}

func (at *DB2SchemaMetaTaskV2) getViewsFromSchema(ctx context.Context, plugin driver.Plugin, schema string) (views []string, err error) {
	res, err := plugin.Query(ctx, fmt.Sprintf(`SELECT VIEWNAME FROM SYSCAT.VIEWS WHERE VIEWSCHEMA = '%v'`, schema), &driverV2.QueryConf{TimeOutSecond: 10})
	if err != nil {
		return nil, fmt.Errorf("query sql failed: %v", err)
	}
	if len(res.Column) != 1 || res.Column[0].Key != "VIEWNAME" {
		return nil, fmt.Errorf("parse query results failed")
	}
	for _, row := range res.Rows {
		for _, value := range row.Values {
			views = append(views, value.Value)
			continue
		}
	}
	return views, nil
}

func (at *DB2SchemaMetaTaskV2) extractSQL(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) ([]*SchemaMetaSQL, error) {
	if ap.InstanceID == "" {
		return nil, fmt.Errorf("instance is not configured")
	}
	instance, exist, err := dms.GetInstancesById(context.Background(), ap.InstanceID)
	if err != nil {
		return nil, fmt.Errorf("get instance fail, error: %v", err)
	}
	if !exist {
		return nil, fmt.Errorf("instance: %v is not exist", ap.InstanceID)
	}
	pluginMgr := driver.GetPluginManager()
	if !pluginMgr.IsOptionalModuleEnabled(instance.DbType, driverV2.OptionalModuleQuery) {
		return nil, fmt.Errorf("collect DB2 schema meta failed: %v", driver.NewErrPluginAPINotImplement(driverV2.OptionalModuleQuery))
	}
	plugin, err := pluginMgr.OpenPlugin(logger, instance.DbType, &driverV2.Config{DSN: &driverV2.DSN{
		Host:             instance.Host,
		Port:             instance.Port,
		User:             instance.User,
		Password:         instance.Password,
		AdditionalParams: instance.AdditionalParams,
	}})
	if err != nil {
		return nil, fmt.Errorf("connect to instance fail, %v", err)
	}
	tempVariableName := fmt.Sprintf("%v.sqle_get_ddl_token", ap.ID)
	valIsCreated := false
	defer func() {
		if valIsCreated {
			_, err = plugin.Exec(context.Background(), fmt.Sprintf(`DROP VARIABLE %v`, tempVariableName))
			if err != nil {
				logger.Errorf("drop variable failed, error: %v", err)
			}
		}

		plugin.Close(context.Background())
	}()

	schemas, err := plugin.Schemas(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get schemas from db2 failed, error: %v", err)
	}
	sqls := []*SchemaMetaSQL{}

	for _, schema := range schemas {
		tables, err := at.getTablesFromSchema(context.Background(), plugin, schema)
		if err != nil {
			return nil, fmt.Errorf("get tables from schema [%v] failed, error: %v", schema, err)
		}

		var views []string
		if ap.Params.GetParam("collect_view").Bool() {
			views, err = at.getViewsFromSchema(context.Background(), plugin, schema)
			if err != nil {
				return nil, fmt.Errorf("get views from schema [%v] failed, error: %v", schema, err)
			}
		}

		_, err = plugin.Exec(context.Background(), fmt.Sprintf(`CREATE OR REPLACE VARIABLE %v integer`, tempVariableName))
		if err != nil {
			return nil, fmt.Errorf("create variable failed, error: %v", err)
		}
		valIsCreated = true
		for _, table := range tables {
			sql := fmt.Sprintf(`CALL SYSPROC.DB2LK_GENERATE_DDL('-t %v.%v -e',%v)`, schema, table, tempVariableName)
			_, err = plugin.Exec(context.Background(), sql)
			if err != nil {
				logger.Errorf("generate ddl failed, sql: %s, error: %v", sql, err)
				continue
			}
			result, err := plugin.Query(context.Background(), fmt.Sprintf(`
	SELECT VARCHAR(SQL_STMT,2000) AS CREATE_TABLE_DDL FROM SYSTOOLS.DB2LOOK_INFO WHERE OP_TOKEN = %v AND OBJ_TYPE = 'TABLE' ORDER BY OP_SEQUENCE ASC
	`, tempVariableName), &driverV2.QueryConf{TimeOutSecond: 10})
			if err != nil {
				logger.Errorf("get create table ddl for table [%v] failed, error: %v", table, err)
				continue
			}

			if len(result.Column) != 1 || result.Column[0].Key != "CREATE_TABLE_DDL" || len(result.Rows) != 1 {
				logger.Errorf("parse create table ddl records  for table [%v] failed", table)
				continue
			}
			createTableDDL := ""
			for _, value := range result.Rows[0].Values {
				createTableDDL = value.Value
			}
			if createTableDDL == "" {
				logger.Errorf("get empty create table ddl for table %v.%v", schema, table)
				continue
			}
			sqls = append(sqls, &SchemaMetaSQL{
				SchemaName: schema,
				MetaName:   table,
				MetaType:   "table",
				SQLContent: createTableDDL,
			})
		}

		for _, view := range views {
			sql := fmt.Sprintf(`call SYSPROC.DB2LK_GENERATE_DDL('-v %v.%v -e',%v)`, schema, view, tempVariableName)
			_, err = plugin.Exec(context.Background(), sql)
			if err != nil {
				logger.Errorf("generate ddl failed, sql: %s, error: %v", sql, err)
				continue
			}
			result, err := plugin.Query(context.Background(), fmt.Sprintf(`
	SELECT VARCHAR(SQL_STMT,2000) AS CREATE_VIEW_DDL FROM SYSTOOLS.DB2LOOK_INFO WHERE OP_TOKEN = %v AND OBJ_TYPE = 'VIEW' ORDER BY OP_SEQUENCE ASC
	`, tempVariableName), &driverV2.QueryConf{TimeOutSecond: 10})
			if err != nil {
				logger.Errorf("get create view ddl for view [%v] failed, error: %v", view, err)
				continue
			}

			if len(result.Column) != 1 || result.Column[0].Key != "CREATE_VIEW_DDL" || len(result.Rows) != 1 {
				logger.Errorf("parse create view ddl records  for view [%v] failed", view)
				continue
			}
			createViewDDL := ""
			for _, value := range result.Rows[0].Values {
				createViewDDL = value.Value
			}
			if createViewDDL == "" {
				logger.Errorf("get empty create view ddl for view %v.%v", schema, view)
				continue
			}
			sqls = append(sqls, &SchemaMetaSQL{
				SchemaName: schema,
				MetaName:   view,
				MetaType:   "view",
				SQLContent: createViewDDL,
			})
		}
	}
	return sqls, nil
}
