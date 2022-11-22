//go:build enterprise
// +build enterprise

package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/common"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
)

func listTableBySchema(c echo.Context) error {
	s := model.GetStorage()
	instanceName := c.Param("instance_name")
	schemaName := c.Param("schema_name")
	projectName := c.Param("project_name")

	instance, exist, err := s.GetInstanceByNameAndProjectName(instanceName, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("instance is not exist")))
	}

	can, err := checkCurrentUserCanAccessInstance(c, instance)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !can {
		return controller.JSONBaseErrorReq(c, errInstanceNoAccess)
	}

	dsn, err := common.NewDSN(instance, "")
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	drvMgr, err := driver.NewDriverManger(log.NewEntry(), instance.DbType, &driver.Config{DSN: dsn})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	defer drvMgr.Close(context.TODO())

	analysisDriver, err := drvMgr.GetAnalysisDriver()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	tablesInSchema, err := analysisDriver.ListTablesInSchema(context.TODO(), &driver.ListTablesInSchemaConf{Schema: schemaName})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &ListTableBySchemaResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertTablesToRes(tablesInSchema),
	})
}

func convertTablesToRes(listTablesInSchemaResult *driver.ListTablesInSchemaResult) []Table {
	result := make([]Table, 0, len(listTablesInSchemaResult.Tables))
	for _, table := range listTablesInSchemaResult.Tables {
		result = append(result, Table{Name: table.Name})
	}

	return result
}

func getTableMetadata(c echo.Context) error {
	s := model.GetStorage()
	instanceName := c.Param("instance_name")
	schemaName := c.Param("schema_name")
	tableName := c.Param("table_name")
	projectName := c.Param("project_name")

	instance, exist, err := s.GetInstanceDetailByNameAndProjectName(instanceName, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrInstanceNotExist)
	}

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if user.Name != model.DefaultAdminUser {
		exist, err := s.CheckUserHasOpToInstances(user, []*model.Instance{instance}, []uint{model.OP_SQL_QUERY_QUERY})
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if !exist {
			return controller.JSONBaseErrorReq(c, errInstanceNoAccess)
		}
	}

	dsn, err := common.NewDSN(instance, schemaName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	drvMgr, err := driver.NewDriverManger(log.NewEntry(), instance.DbType, &driver.Config{DSN: dsn})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	defer drvMgr.Close(context.TODO())

	analysisDriver, err := drvMgr.GetAnalysisDriver()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	tableMetadata, err := analysisDriver.GetTableMetaByTableName(context.TODO(), &driver.GetTableMetaByTableNameConf{
		Schema: schemaName,
		Table:  tableName,
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetTableMetadataResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertTableMetadataToRes(tableMetadata),
	})
}

func convertTableMetadataToRes(metaByTableNameResult *driver.GetTableMetaByTableNameResult) InstanceTableMeta {
	columnsInfo := metaByTableNameResult.TableMeta.ColumnsInfo
	indexesInfo := metaByTableNameResult.TableMeta.IndexesInfo

	tableMeta := InstanceTableMeta{
		Columns: TableColumns{
			Rows: make([]map[string]string, len(columnsInfo.Rows)),
			Head: make([]TableMetaItemHeadResV1, len(columnsInfo.Columns)),
		},
		Indexes: TableIndexes{
			Rows: make([]map[string]string, len(indexesInfo.Rows)),
			Head: make([]TableMetaItemHeadResV1, len(indexesInfo.Columns)),
		},
	}

	for i, column := range columnsInfo.Columns {
		tableMeta.Columns.Head[i].FieldName = column.Name
		tableMeta.Columns.Head[i].Desc = column.Desc
	}

	for i, rows := range columnsInfo.Rows {
		tableMeta.Columns.Rows[i] = make(map[string]string)
		for k, row := range rows {
			columnName := columnsInfo.Columns[k].Name
			tableMeta.Columns.Rows[i][columnName] = row
		}
	}

	for i, column := range indexesInfo.Columns {
		tableMeta.Indexes.Head[i].FieldName = column.Name
		tableMeta.Indexes.Head[i].Desc = column.Desc
	}

	for i, rows := range indexesInfo.Rows {
		tableMeta.Indexes.Rows[i] = make(map[string]string)
		for k, row := range rows {
			columnName := indexesInfo.Columns[k].Name
			tableMeta.Indexes.Rows[i][columnName] = row
		}
	}

	tableMeta.CreateTableSQL = metaByTableNameResult.TableMeta.CreateTableSQL
	tableMeta.Name = metaByTableNameResult.TableMeta.Name
	tableMeta.Schema = metaByTableNameResult.TableMeta.Schema

	return tableMeta
}
