package v1

import (
	"context"
	"github.com/pkg/errors"
	"strconv"

	"github.com/actiontech/sqle/sqle/common"
	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/sirupsen/logrus"
)

type AnalysisResult struct {
	Cost             *float64
	ExplainResult    *driverV2.ExplainResult
	ExplainResultErr error

	TableMetaResult    *driver.GetTableMetaBySQLResult
	TableMetaResultErr error

	AffectRowsResult    *driverV2.EstimatedAffectRows
	AffectRowsResultErr error
}

func GetSQLAnalysisResult(l *logrus.Entry, instance *model.Instance, schema, sql string) (res *AnalysisResult, err error) {
	dsn, err := common.NewDSN(instance, schema)
	if err != nil {
		return nil, err
	}

	plugin, err := driver.GetPluginManager().
		OpenPlugin(l, instance.DbType, &driverV2.Config{DSN: dsn})
	if err != nil {
		return nil, err
	}
	defer plugin.Close(context.TODO())

	res = &AnalysisResult{}

	res.Cost, err = GetQueryCost(plugin, sql)
	res.ExplainResult, res.ExplainResultErr = Explain(instance.DbType, plugin, sql)
	res.TableMetaResult, res.TableMetaResultErr = GetTableMetas(instance.DbType, plugin, sql)
	res.AffectRowsResult, res.AffectRowsResultErr = GetRowsAffected(instance.DbType, plugin, sql)

	return res, nil
}

func GetQueryCost(plugin driver.Plugin, sql string) (cost *float64, err error) {
	explainJSONResult, err := plugin.ExplainJSONFormat(context.TODO(), &driverV2.ExplainConf{Sql: sql})
	if err != nil {
		return nil, err
	}
	sqlNodes, err := plugin.Parse(context.TODO(), sql)
	if err != nil {
		return nil, err
	}
	if len(sqlNodes) == 0 {
		return nil, errors.Errorf("failed to get query cost invalid sql: %v", sql)
	}
	// 不是SELECT语句直接忽略
	if sqlNodes[0].Type != driverV2.SQLTypeDQL {
		return nil, errors.Errorf("failed to get query cost because it is not DQL: %v", sql)
	}
	costFloat, err := strconv.ParseFloat(explainJSONResult.QueryBlock.CostInfo.QueryCost, 64)
	if err != nil {
		return nil, err
	}
	return &costFloat, nil
}

func Explain(dbType string, plugin driver.Plugin, sql string) (res *driverV2.ExplainResult, err error) {
	if !driver.GetPluginManager().
		IsOptionalModuleEnabled(dbType, driverV2.OptionalModuleExplain) {
		return nil, driver.NewErrPluginAPINotImplement(driverV2.OptionalModuleExplain)
	}

	return plugin.Explain(context.TODO(), &driverV2.ExplainConf{Sql: sql})
}

func GetTableMetas(dbType string, plugin driver.Plugin, sql string) (res *driver.GetTableMetaBySQLResult, err error) {
	if !driver.GetPluginManager().
		IsOptionalModuleEnabled(dbType, driverV2.OptionalModuleGetTableMeta) {
		return nil, driver.NewErrPluginAPINotImplement(driverV2.OptionalModuleGetTableMeta)
	}

	return plugin.GetTableMetaBySQL(context.TODO(), &driver.GetTableMetaBySQLConf{Sql: sql})
}

func GetRowsAffected(dbType string, plugin driver.Plugin, sql string) (res *driverV2.EstimatedAffectRows, err error) {
	if !driver.GetPluginManager().
		IsOptionalModuleEnabled(dbType, driverV2.OptionalModuleEstimateSQLAffectRows) {
		return nil, driver.NewErrPluginAPINotImplement(driverV2.OptionalModuleEstimateSQLAffectRows)
	}

	return plugin.EstimateSQLAffectRows(context.TODO(), sql)
}
