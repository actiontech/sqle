package mysql

import (
	"context"
	"fmt"

	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/pkg/params"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func init() {
	driver.RegisterAnalysisDriver(driver.DriverTypeMySQL, newAnalysisDriverInspect)
}

func newAnalysisDriverInspect(log *logrus.Entry, dsn *driver.DSN) (driver.AnalysisDriver, error) {
	var inspect = &Inspect{}

	if dsn != nil {
		conn, err := executor.NewExecutor(log, dsn, dsn.DatabaseName)
		if err != nil {
			return nil, errors.Wrap(err, "new executor in inspect")
		}
		inspect.isConnected = true
		inspect.dbConn = conn
		inspect.inst = dsn

		ctx := session.NewContext(nil, session.WithExecutor(conn))
		ctx.SetCurrentSchema(dsn.DatabaseName)

		inspect.Ctx = ctx
	} else {
		ctx := session.NewContext(nil)
		inspect.Ctx = ctx
	}

	inspect.log = log
	inspect.result = driver.NewInspectResults()
	inspect.isOfflineAudit = dsn == nil

	inspect.cnf = &Config{
		DMLRollbackMaxRows: -1,
		DDLOSCMinSize:      -1,
		DDLGhostMinSize:    -1,
	}

	return inspect, nil
}

// ListTablesInSchema list tables in specified schema.
func (i *Inspect) ListTablesInSchema(ctx context.Context, conf *driver.ListTablesInSchemaConf) (*driver.ListTablesInSchemaResult, error) {
	conn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}
	defer conn.Db.Close()
	tables, err := conn.ShowSchemaTables(conf.Schema)
	if err != nil {
		return nil, err
	}

	resItems := make([]driver.Table, len(tables))
	for i, t := range tables {
		resItems[i].Name = t
	}
	return &driver.ListTablesInSchemaResult{Tables: resItems}, nil
}

// GetTableMetaByTableName get table's metadata by table name.
func (i *Inspect) GetTableMetaByTableName(ctx context.Context, conf *driver.GetTableMetaByTableNameConf) (*driver.GetTableMetaByTableNameResult, error) {
	return nil, nil
}

// GetTableMetaBySQL get table's metadata by SQL.
func (i *Inspect) GetTableMetaBySQL(ctx context.Context, conf *driver.GetTableMetaBySQLConf) (*driver.GetTableMetaBySQLResult, error) {
	return nil, nil
}

// Explain get explain result for SQL.
func (i *Inspect) Explain(ctx context.Context, conf *driver.ExplainConf) (*driver.ExplainResult, error) {
	// check sql
	// only support dml
	nodes, err := i.ParseSql(conf.Sql)
	if err != nil {
		return nil, err
	}
	switch nodes[0].(type) {
	case ast.DMLNode:
	default:
		return nil, fmt.Errorf("the sql is `%v`, but we only support DML", conf.Sql)
	}

	conn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}
	defer conn.Db.Close()

	columns, rows, err := conn.Db.QueryWithContext(context.TODO(), fmt.Sprintf("EXPLAIN %s", conf.Sql))
	if err != nil {
		return nil, err
	}

	resColumn := params.Params{}
	for _, column := range columns {
		resColumn = append(resColumn, &params.Param{
			Key:   column,
			Value: column,
		})
	}

	resRows := make([][]string, len(rows))
	for i, row := range rows {
		for _, s := range row {
			resRows[i] = append(resRows[i], s.String)
		}
	}
	res := driver.ExplainClassicResult{
		AnalysisInfoInTableFormat: driver.AnalysisInfoInTableFormat{
			Column: resColumn,
			Rows:   resRows,
		},
	}

	return &driver.ExplainResult{
		ClassicResult: res,
	}, nil
}
