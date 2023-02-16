package driver

import (
	"context"

	"github.com/actiontech/sqle/sqle/driver/proto"

	goPlugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// queryDriverPlugin use for hide gRPC detail.
type analysisDriverGRPCServer struct {
	newDriver func(cfg *DSN) AnalysisDriver

	impl AnalysisDriver
}

func (a *analysisDriverGRPCServer) Init(c context.Context, req *proto.AnalysisDriverInitRequest) (*proto.Empty, error) {
	var dsn *DSN
	if req.GetDsn() != nil {
		dsn = &DSN{
			Host:             req.GetDsn().GetHost(),
			Port:             req.GetDsn().GetPort(),
			User:             req.GetDsn().GetUser(),
			Password:         req.GetDsn().GetPassword(),
			DatabaseName:     req.GetDsn().GetDatabase(),
			AdditionalParams: proto.ConvertProtoParamToParam(req.GetDsn().GetAdditionalParams()),
		}
	}
	a.impl = a.newDriver(dsn)
	return &proto.Empty{}, nil
}

func (a *analysisDriverGRPCServer) ListTablesInSchema(ctx context.Context, req *proto.ListTablesInSchemaRequest) (*proto.ListTablesInSchemaResponse, error) {
	schema := req.GetSchema()
	res, err := a.impl.ListTablesInSchema(ctx, &ListTablesInSchemaConf{Schema: schema})
	if err != nil {
		return &proto.ListTablesInSchemaResponse{}, err
	}

	tables := make([]*proto.Table, len(res.Tables))
	for i, table := range res.Tables {
		t := &proto.Table{Name: table.Name}
		tables[i] = t
	}
	resp := &proto.ListTablesInSchemaResponse{
		Tables: tables,
	}
	return resp, nil
}

func (a *analysisDriverGRPCServer) GetTableMetaByTableName(ctx context.Context, req *proto.GetTableMetaByTableNameRequest) (*proto.GetTableMetaByTableNameResponse, error) {
	conf := &GetTableMetaByTableNameConf{
		Schema: req.GetSchema(),
		Table:  req.GetTable(),
	}
	res, err := a.impl.GetTableMetaByTableName(ctx, conf)
	if err != nil {
		return &proto.GetTableMetaByTableNameResponse{}, err
	}

	columnInfoColumns, columnInfoRows := a.convertTableFormatInfoToProto(res.TableMeta.ColumnsInfo.AnalysisInfoInTableFormat)
	indexInfoColumns, indexInfoRows := a.convertTableFormatInfoToProto(res.TableMeta.IndexesInfo.AnalysisInfoInTableFormat)

	tableMeta := &proto.TableItem{
		Name:   res.TableMeta.Name,
		Schema: res.TableMeta.Schema,
		ColumnsInfo: &proto.ColumnsInfo{
			AnalysisInfoInTableFormat: &proto.AnalysisInfoInTableFormat{
				Columns: columnInfoColumns,
				Rows:    columnInfoRows,
			}},
		IndexesInfo: &proto.IndexesInfo{
			AnalysisInfoInTableFormat: &proto.AnalysisInfoInTableFormat{
				Columns: indexInfoColumns,
				Rows:    indexInfoRows,
			}},
		CreateTableSQL: res.TableMeta.CreateTableSQL,
	}
	resp := &proto.GetTableMetaByTableNameResponse{
		TableMeta: tableMeta,
	}
	return resp, nil
}

func (a *analysisDriverGRPCServer) convertTableFormatInfoToProto(analysisInfo AnalysisInfoInTableFormat) (columns []*proto.AnalysisInfoHead, rows []*proto.Row) {
	columns = make([]*proto.AnalysisInfoHead, len(analysisInfo.Columns))
	for i, c := range analysisInfo.Columns {
		columns[i] = &proto.AnalysisInfoHead{
			Name: c.Name,
			Desc: c.Desc,
		}
	}

	rows = make([]*proto.Row, len(analysisInfo.Rows))
	for i, r := range analysisInfo.Rows {
		rows[i] = &proto.Row{
			Items: r,
		}
	}
	return
}

func (a *analysisDriverGRPCServer) GetTableMetaBySQL(ctx context.Context, req *proto.GetTableMetaBySQLRequest) (*proto.GetTableMetaBySQLResponse, error) {
	conf := &GetTableMetaBySQLConf{
		Sql: req.Sql,
	}
	res, err := a.impl.GetTableMetaBySQL(ctx, conf)
	if err != nil && err == ErrSQLIsNotSupported {
		return &proto.GetTableMetaBySQLResponse{}, status.Error(grpcErrSQLIsNotSupported, err.Error())
	} else if err != nil {
		return &proto.GetTableMetaBySQLResponse{}, err
	}

	tableMetas := make([]*proto.TableMetaItemBySQL, len(res.TableMetas))
	for i, table := range res.TableMetas {
		columnInfoColumns, columnInfoRows := a.convertTableFormatInfoToProto(table.ColumnsInfo.AnalysisInfoInTableFormat)
		indexInfoColumns, indexInfoRows := a.convertTableFormatInfoToProto(table.IndexesInfo.AnalysisInfoInTableFormat)
		tableMetas[i] = &proto.TableMetaItemBySQL{
			Name:   table.Name,
			Schema: table.Schema,
			ColumnsInfo: &proto.ColumnsInfo{
				AnalysisInfoInTableFormat: &proto.AnalysisInfoInTableFormat{
					Columns: columnInfoColumns,
					Rows:    columnInfoRows,
				}},
			IndexesInfo: &proto.IndexesInfo{
				AnalysisInfoInTableFormat: &proto.AnalysisInfoInTableFormat{
					Columns: indexInfoColumns,
					Rows:    indexInfoRows,
				}},
			CreateTableSQL: table.CreateTableSQL,
			ErrMessage:     table.Message,
		}
	}

	resp := &proto.GetTableMetaBySQLResponse{
		TableMetas: tableMetas,
	}
	return resp, nil
}

func (a *analysisDriverGRPCServer) Explain(ctx context.Context, req *proto.ExplainRequest) (*proto.ExplainResponse, error) {
	conf := &ExplainConf{
		Sql: req.Sql,
	}
	res, err := a.impl.Explain(ctx, conf)
	if err != nil && err == ErrSQLIsNotSupported {
		return &proto.ExplainResponse{}, status.Error(grpcErrSQLIsNotSupported, err.Error())
	} else if err != nil {
		return &proto.ExplainResponse{}, err
	}

	columns, rows := a.convertTableFormatInfoToProto(res.ClassicResult.AnalysisInfoInTableFormat)
	classicRes := &proto.ExplainClassicResult{
		AnalysisInfoInTableFormat: &proto.AnalysisInfoInTableFormat{
			Columns: columns,
			Rows:    rows,
		}}

	resp := &proto.ExplainResponse{
		ClassicResult: classicRes,
	}
	return resp, nil
}

// analysisDriverPlugin implements goPlugin.GRPCPlugin
type analysisDriverPlugin struct {
	goPlugin.NetRPCUnsupportedPlugin

	Srv *analysisDriverGRPCServer
}

func NewAnalysisDriverPlugin(newDriver func(cfg *DSN) AnalysisDriver) *analysisDriverPlugin {
	return &analysisDriverPlugin{
		NetRPCUnsupportedPlugin: goPlugin.NetRPCUnsupportedPlugin{},
		Srv: &analysisDriverGRPCServer{
			newDriver: newDriver,
		},
	}
}

func (dp *analysisDriverPlugin) GRPCServer(broker *goPlugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterAnalysisDriverServer(s, dp.Srv)
	return nil
}

func (dp *analysisDriverPlugin) GRPCClient(ctx context.Context, broker *goPlugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return proto.NewAnalysisDriverClient(c), nil
}
