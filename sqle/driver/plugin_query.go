package driver

import (
	"context"
	"github.com/actiontech/sqle/sqle/driver/proto"
	goPlugin "github.com/hashicorp/go-plugin"

	"google.golang.org/grpc"
)

// queryDriverPlugin use for hide gRPC detail.
type queryDriverGRPCServer struct {
	newDriver func(cfg *DSN) SQLQueryDriver

	impl SQLQueryDriver
}

func (q *queryDriverGRPCServer) Init(c context.Context, req *proto.InitRequest) (*proto.Empty, error) {
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
	q.impl = q.newDriver(dsn)
	return &proto.Empty{}, nil
}

func (q *queryDriverGRPCServer) QueryPrepare(ctx context.Context, req *proto.QueryPrepareRequest) (*proto.QueryPrepareResponse, error) {
	conf := &QueryPrepareConf{
		Limit:  req.GetConf().GetLimit(),
		Offset: req.GetConf().GetOffset(),
	}
	res, err := q.impl.QueryPrepare(ctx, req.GetSql(), conf)
	if err != nil {
		return &proto.QueryPrepareResponse{}, err
	}

	resp := &proto.QueryPrepareResponse{
		NewSql:    res.NewSQL,
		ErrorType: string(res.ErrorType),
		Error:     res.Error,
	}
	return resp, nil
}

func (q *queryDriverGRPCServer) Query(ctx context.Context, req *proto.QueryRequest) (*proto.QueryResponse, error) {
	conf := &QueryConf{
		TimeOutSecond: req.GetConf().GetTimeOutSecond(),
	}
	res, err := q.impl.Query(ctx, req.GetSql(), conf)
	if err != nil {
		return &proto.QueryResponse{}, err
	}
	resp := &proto.QueryResponse{
		Column: []*proto.Param{},
		Rows:   []*proto.QueryResultRow{},
	}
	for _, param := range res.Column {
		resp.Column = append(resp.Column, &proto.Param{
			Key:   param.Key,
			Value: param.Value,
			Desc:  param.Desc,
			Type:  string(param.Type),
		})
	}
	for _, row := range res.Rows {
		rows := &proto.QueryResultRow{
			Values: []*proto.QueryResultValue{},
		}
		for _, value := range row.Values {
			rows.Values = append(rows.Values, &proto.QueryResultValue{
				Value: value.Value,
			})
		}
		resp.Rows = append(resp.Rows, rows)
	}
	return resp, nil
}

// queryDriverPlugin implements goPlugin.GRPCPlugin
type queryDriverPlugin struct {
	goPlugin.NetRPCUnsupportedPlugin

	Srv *queryDriverGRPCServer
}

func NewQueryDriverPlugin(newDriver func(cfg *DSN) SQLQueryDriver) *queryDriverPlugin {
	return &queryDriverPlugin{
		NetRPCUnsupportedPlugin: goPlugin.NetRPCUnsupportedPlugin{},
		Srv: &queryDriverGRPCServer{
			newDriver: newDriver,
		},
	}
}

func (dp *queryDriverPlugin) GRPCServer(broker *goPlugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterQueryDriverServer(s, dp.Srv)
	return nil
}

func (dp *queryDriverPlugin) GRPCClient(ctx context.Context, broker *goPlugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return proto.NewQueryDriverClient(c), nil
}
