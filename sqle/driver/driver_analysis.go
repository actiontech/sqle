package driver

import (
	"context"
	"fmt"
	"sync"

	"github.com/actiontech/sqle/sqle/driver/proto"
	"github.com/actiontech/sqle/sqle/log"

	goPlugin "github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
)

// AnalysisDriver is a driver for SQL analysis and getting table metadata
type AnalysisDriver interface {
	ListTablesInSchema(ctx context.Context, conf *ListTablesInSchemaConf) (*ListTablesInSchemaResult, error)
	GetTableMetaByTableName(ctx context.Context, conf *GetTableMetaByTableNameConf) (*GetTableMetaByTableNameResult, error)
	GetTableMetaBySQL(ctx context.Context, conf *GetTableMetaBySQLConf) (*GetTableMetaBySQLResult, error)
	Explain(ctx context.Context, conf *ExplainConf) (*ExplainResult, error)
}

func init() {
	defaultPluginSet[DefaultPluginVersion][PluginNameAnalysisDriver] = &analysisDriverPlugin{}
}

const (
	PluginNameAnalysisDriver = "analysis-driver"
)

type ListTablesInSchemaConf struct {
	Schema string
}

type Table struct {
	Name string
}

type ListTablesInSchemaResult struct {
	Tables []Table
}

type AnalysisInfoHead struct {
	Name string
	Desc string
}

// AnalysisInfoInTableFormat
// the field Column represents the column name of a table
// the field Rows represents the data of the table
// their relationship is as follows
/*
	| Column[0]  | Column[1]  | Column[2]  |
	| Rows[0][0] | Rows[0][1] | Rows[0][2] |
	| Rows[1][0] | Rows[1][1] | Rows[1][2] |
*/
type AnalysisInfoInTableFormat struct {
	Column []AnalysisInfoHead
	Rows   [][]string
}

type GetTableMetaByTableNameConf struct {
	Schema string
	Table  string
}

type GetTableMetaBySQLConf struct {
	// this SQL should be a single SQL
	Sql string
}

type ColumnsInfo struct {
	AnalysisInfoInTableFormat
}

type IndexesInfo struct {
	AnalysisInfoInTableFormat
}

type TableMetaItem struct {
	Name           string
	Schema         string
	ColumnsInfo    ColumnsInfo
	IndexesInfo    IndexesInfo
	CreateTableSQL string
}

type GetTableMetaByTableNameResult struct {
	TableMeta TableMetaItem
}

type TableMetaItemBySQL struct {
	Name           string
	Schema         string
	ColumnsInfo    ColumnsInfo
	IndexesInfo    IndexesInfo
	CreateTableSQL string
	Message        string
}

type GetTableMetaBySQLResult struct {
	TableMetas []TableMetaItemBySQL
}

type ExplainConf struct {
	// this SQL should be a single SQL
	Sql string
}

type ExplainClassicResult struct {
	AnalysisInfoInTableFormat
}

type ExplainResult struct {
	ClassicResult ExplainClassicResult
}

var analysisDriverMu = &sync.RWMutex{}
var analysisDrivers = make(map[string]struct{})

// RegisterAnalysisDriver like sql.RegisterAuditDriver.
//
// RegisterAnalysisDriver makes a database driver available by the provided driver name.
// RegisterAnalysisDriver's initialize handler registered by RegisterAnalysisDriver.
func RegisterAnalysisDriver(name string) {
	analysisDriverMu.RLock()
	_, exist := analysisDrivers[name]
	analysisDriverMu.RUnlock()
	if exist {
		panic(fmt.Sprintf("duplicated driver name %v", name))
	}

	analysisDriverMu.Lock()
	analysisDrivers[name] = struct{}{}
	analysisDriverMu.Unlock()
}

// analysisDriverImpl implement AnalysisDriver. It use for hide gRPC detail, just like DriverGRPCServer.
type analysisDriverImpl struct {
	plugin proto.AnalysisDriverClient
}

func (a *analysisDriverImpl) ListTablesInSchema(ctx context.Context, conf *ListTablesInSchemaConf) (*ListTablesInSchemaResult, error) {
	req := &proto.ListTablesInSchemaRequest{
		Schema: conf.Schema,
	}
	res, err := a.plugin.ListTablesInSchema(ctx, req)
	if err != nil {
		return nil, err
	}

	tables := make([]Table, len(res.Tables))
	for i, t := range res.Tables {
		tables[i] = Table{Name: t.GetName()}
	}
	return &ListTablesInSchemaResult{
		Tables: tables,
	}, nil
}

func (a *analysisDriverImpl) GetTableMetaByTableName(ctx context.Context, conf *GetTableMetaByTableNameConf) (*GetTableMetaByTableNameResult, error) {
	req := &proto.GetTableMetaByTableNameRequest{
		Schema: conf.Schema,
		Table:  conf.Table,
	}
	res, err := a.plugin.GetTableMetaByTableName(ctx, req)
	if err != nil {
		return nil, err
	}

	resTableMeta := res.GetTableMeta()

	ColumnsInfoColumns, ColumnsInfoRows := a.convertAnalysisInfoFromProtoToSqle(resTableMeta.ColumnsInfo.AnalysisInfoInTableFormat)
	IndexesInfoColumns, IndexesInfoRows := a.convertAnalysisInfoFromProtoToSqle(resTableMeta.IndexesInfo.AnalysisInfoInTableFormat)

	tableMeta := TableMetaItem{
		Name:   resTableMeta.GetName(),
		Schema: resTableMeta.GetSchema(),
		ColumnsInfo: ColumnsInfo{
			AnalysisInfoInTableFormat{
				Column: ColumnsInfoColumns,
				Rows:   ColumnsInfoRows,
			}},
		IndexesInfo: IndexesInfo{
			AnalysisInfoInTableFormat{
				Column: IndexesInfoColumns,
				Rows:   IndexesInfoRows,
			},
		},
		CreateTableSQL: resTableMeta.GetCreateTableSQL(),
	}
	return &GetTableMetaByTableNameResult{
		TableMeta: tableMeta,
	}, nil
}

func (a *analysisDriverImpl) convertAnalysisInfoFromProtoToSqle(protoInfo *proto.AnalysisInfoInTableFormat) (columns []AnalysisInfoHead, rows [][]string) {
	columns = make([]AnalysisInfoHead, len(protoInfo.Columns))
	for i, c := range protoInfo.Columns {
		columns[i] = AnalysisInfoHead{
			Name: c.GetName(),
			Desc: c.GetDesc(),
		}
	}

	rows = make([][]string, len(protoInfo.Rows))
	for i, r := range protoInfo.Rows {
		rows[i] = r.GetItems()
	}

	return
}

func (a *analysisDriverImpl) GetTableMetaBySQL(ctx context.Context, conf *GetTableMetaBySQLConf) (*GetTableMetaBySQLResult, error) {
	req := &proto.GetTableMetaBySQLRequest{
		Sql: conf.Sql,
	}
	res, err := a.plugin.GetTableMetaBySQL(ctx, req)
	if err != nil {
		return nil, err
	}

	tableMetas := make([]TableMetaItemBySQL, len(res.GetTableMetas()))
	for i, resTableMeta := range res.GetTableMetas() {
		ColumnsInfoColumns, ColumnsInfoRows := a.convertAnalysisInfoFromProtoToSqle(resTableMeta.ColumnsInfo.AnalysisInfoInTableFormat)
		IndexesInfoColumns, IndexesInfoRows := a.convertAnalysisInfoFromProtoToSqle(resTableMeta.IndexesInfo.AnalysisInfoInTableFormat)

		tableMeta := TableMetaItemBySQL{
			Name:   resTableMeta.GetName(),
			Schema: resTableMeta.GetSchema(),
			ColumnsInfo: ColumnsInfo{
				AnalysisInfoInTableFormat{
					Column: ColumnsInfoColumns,
					Rows:   ColumnsInfoRows,
				}},
			IndexesInfo: IndexesInfo{
				AnalysisInfoInTableFormat{
					Column: IndexesInfoColumns,
					Rows:   IndexesInfoRows,
				},
			},
			CreateTableSQL: resTableMeta.GetCreateTableSQL(),
			Message:        resTableMeta.GetErrMessage(),
		}
		tableMetas[i] = tableMeta
	}

	return &GetTableMetaBySQLResult{
		TableMetas: tableMetas,
	}, nil
}

func (a *analysisDriverImpl) Explain(ctx context.Context, conf *ExplainConf) (*ExplainResult, error) {
	req := &proto.ExplainRequest{
		Sql: conf.Sql,
	}
	res, err := a.plugin.Explain(ctx, req)
	if err != nil {
		return nil, err
	}

	columns, rows := a.convertAnalysisInfoFromProtoToSqle(res.ClassicResult.AnalysisInfoInTableFormat)
	return &ExplainResult{
		ClassicResult: ExplainClassicResult{AnalysisInfoInTableFormat{
			Column: columns,
			Rows:   rows,
		}},
	}, nil
}

func registerAnalysisDriver(pluginName string, gRPCClient goPlugin.ClientProtocol) error {
	rawI, err := gRPCClient.Dispense(PluginNameAnalysisDriver)
	if err != nil {
		return err
	}
	//nolint:forcetypeassert
	s := rawI.(proto.AnalysisDriverClient)

	// The test target plugin implements the AnalysisDriver plugin
	_, err = s.Init(context.TODO(), &proto.AnalysisDriverInitRequest{})
	if err != nil {
		return err
	}

	RegisterAnalysisDriver(pluginName)
	log.Logger().WithFields(logrus.Fields{
		"plugin_name": pluginName,
		"plugin_type": PluginNameAnalysisDriver,
	}).Infoln("plugin inited")
	return nil
}
