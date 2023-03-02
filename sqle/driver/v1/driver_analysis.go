package driverV1

import (
	"context"

	"github.com/actiontech/sqle/sqle/driver/v1/proto"

	"google.golang.org/grpc/status"
)

// AnalysisDriver is a driver for SQL analysis and getting table metadata
type AnalysisDriver interface {
	ListTablesInSchema(ctx context.Context, conf *ListTablesInSchemaConf) (*ListTablesInSchemaResult, error)
	GetTableMetaByTableName(ctx context.Context, conf *GetTableMetaByTableNameConf) (*GetTableMetaByTableNameResult, error)
	GetTableMetaBySQL(ctx context.Context, conf *GetTableMetaBySQLConf) (*GetTableMetaBySQLResult, error)
	Explain(ctx context.Context, conf *ExplainConf) (*ExplainResult, error)
}

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
// the field Columns represents the column name of a table
// the field Rows represents the data of the table
// their relationship is as follows
/*
	| Columns[0]  | Columns[1]  | Columns[2]  |
	| Rows[0][0] | Rows[0][1] | Rows[0][2] |
	| Rows[1][0] | Rows[1][1] | Rows[1][2] |
*/
type AnalysisInfoInTableFormat struct {
	Columns []AnalysisInfoHead
	Rows    [][]string
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
				Columns: ColumnsInfoColumns,
				Rows:    ColumnsInfoRows,
			}},
		IndexesInfo: IndexesInfo{
			AnalysisInfoInTableFormat{
				Columns: IndexesInfoColumns,
				Rows:    IndexesInfoRows,
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
	if err != nil && status.Code(err) == grpcErrSQLIsNotSupported {
		return nil, ErrSQLIsNotSupported
	} else if err != nil {
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
					Columns: ColumnsInfoColumns,
					Rows:    ColumnsInfoRows,
				}},
			IndexesInfo: IndexesInfo{
				AnalysisInfoInTableFormat{
					Columns: IndexesInfoColumns,
					Rows:    IndexesInfoRows,
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
	if err != nil && status.Code(err) == grpcErrSQLIsNotSupported {
		return nil, ErrSQLIsNotSupported
	} else if err != nil {
		return nil, err
	}

	columns, rows := a.convertAnalysisInfoFromProtoToSqle(res.ClassicResult.AnalysisInfoInTableFormat)
	return &ExplainResult{
		ClassicResult: ExplainClassicResult{AnalysisInfoInTableFormat{
			Columns: columns,
			Rows:    rows,
		}},
	}, nil
}
