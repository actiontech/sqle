package driver

import (
	"context"

	"github.com/actiontech/sqle/sqle/driver/proto"

	goPlugin "github.com/hashicorp/go-plugin"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const (
	PluginNameAuditDriver = "audit-driver"
)

// auditDriverPlugin use for hide gRPC detail.
type auditDriverGRPCServer struct {
	newDriver func(cfg *Config) Driver

	impl Driver

	// Registerer provide some plugin info to host process.
	r Registerer
}

func (d *auditDriverGRPCServer) Init(ctx context.Context, req *proto.InitRequest) (*proto.Empty, error) {
	var driverRules = make([]*Rule, 0, len(req.GetRules()))
	for _, rule := range req.GetRules() {
		driverRules = append(driverRules, convertRuleFromProtoToDriver(rule))
	}

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

	cfg, err := NewConfig(dsn, driverRules)
	if err != nil {
		return nil, errors.Wrap(err, "init config")
	}
	d.impl = d.newDriver(cfg)
	return &proto.Empty{}, nil
}

func (d *auditDriverGRPCServer) Close(ctx context.Context, req *proto.Empty) (*proto.Empty, error) {
	d.impl.Close(ctx)
	return &proto.Empty{}, nil
}

func (d *auditDriverGRPCServer) Ping(ctx context.Context, req *proto.Empty) (*proto.Empty, error) {
	return &proto.Empty{}, d.impl.Ping(ctx)
}

func (d *auditDriverGRPCServer) Exec(ctx context.Context, req *proto.ExecRequest) (*proto.ExecResponse, error) {
	result, err := d.impl.Exec(ctx, req.GetQuery())
	if err != nil {
		return &proto.ExecResponse{}, err
	}

	resp := &proto.ExecResponse{}
	lastInsertId, lastInsertIdErr := result.LastInsertId()
	resp.LastInsertId = lastInsertId
	if lastInsertIdErr != nil {
		resp.LastInsertIdError = lastInsertIdErr.Error()
	}
	rowsAffected, rowsAffectedErr := result.RowsAffected()
	resp.RowsAffected = rowsAffected
	if rowsAffectedErr != nil {
		resp.RowsAffectedError = rowsAffectedErr.Error()
	}
	return resp, nil
}

func (d *auditDriverGRPCServer) Tx(ctx context.Context, req *proto.TxRequest) (*proto.TxResponse, error) {
	results, err := d.impl.Tx(ctx, req.GetQueries()...)
	if err != nil {
		return &proto.TxResponse{}, err
	}

	txResp := &proto.TxResponse{}
	for _, result := range results {
		resp := &proto.ExecResponse{}

		lastInsertId, lastInsertIdErr := result.LastInsertId()
		resp.LastInsertId = lastInsertId
		if lastInsertIdErr != nil {
			resp.LastInsertIdError = lastInsertIdErr.Error()
		}
		rowsAffected, rowsAffectedErr := result.RowsAffected()
		resp.RowsAffected = rowsAffected
		if rowsAffectedErr != nil {
			resp.RowsAffectedError = rowsAffectedErr.Error()
		}

		txResp.Results = append(txResp.Results, resp)
	}
	return txResp, nil
}

func (d *auditDriverGRPCServer) Databases(ctx context.Context, req *proto.Empty) (*proto.DatabasesResponse, error) {
	databases, err := d.impl.Schemas(ctx)
	return &proto.DatabasesResponse{Databases: databases}, err
}

func (d *auditDriverGRPCServer) Parse(ctx context.Context, req *proto.ParseRequest) (*proto.ParseResponse, error) {
	nodes, err := d.impl.Parse(ctx, req.GetSqlText())
	if err != nil {
		return &proto.ParseResponse{}, err
	}

	resp := &proto.ParseResponse{}
	for _, node := range nodes {
		resp.Nodes = append(resp.Nodes, &proto.Node{
			Text:        node.Text,
			Type:        node.Type,
			Fingerprint: node.Fingerprint,
		})
	}
	return resp, nil
}

func (d *auditDriverGRPCServer) Audit(ctx context.Context, req *proto.AuditRequest) (*proto.AuditResponse, error) {
	auditResults, err := d.impl.Audit(ctx, req.GetSql())
	if err != nil {
		return &proto.AuditResponse{}, err
	}

	resp := &proto.AuditResponse{}
	for _, result := range auditResults.results {
		resp.Results = append(resp.Results, &proto.AuditResult{
			Level:   string(result.level),
			Message: result.message,
		})
	}
	return resp, nil
}

func (d *auditDriverGRPCServer) GenRollbackSQL(ctx context.Context, req *proto.GenRollbackSQLRequest) (*proto.GenRollbackSQLResponse, error) {
	rollbackSQL, reason, err := d.impl.GenRollbackSQL(ctx, req.GetSql())
	return &proto.GenRollbackSQLResponse{
		Sql:    rollbackSQL,
		Reason: reason,
	}, err
}

func (d *auditDriverGRPCServer) Metas(ctx context.Context, req *proto.Empty) (*proto.MetasResponse, error) {
	protoRules := make([]*proto.Rule, len(d.r.Rules()))

	for i, r := range d.r.Rules() {
		protoRules[i] = convertRuleFromDriverToProto(r)
	}

	return &proto.MetasResponse{
		Name:             d.r.Name(),
		Rules:            protoRules,
		AdditionalParams: proto.ConvertParamToProtoParam(d.r.AdditionalParams()),
	}, nil
}

// auditDriverPlugin implements goPlugin.GRPCPlugin
type auditDriverPlugin struct {
	goPlugin.NetRPCUnsupportedPlugin

	Srv *auditDriverGRPCServer
}

func NewAuditDriverPlugin(r Registerer, newDriver func(cfg *Config) Driver) *auditDriverPlugin {
	return &auditDriverPlugin{
		NetRPCUnsupportedPlugin: goPlugin.NetRPCUnsupportedPlugin{},
		Srv: &auditDriverGRPCServer{
			newDriver: newDriver,
			r:         r,
		},
	}
}

func (dp *auditDriverPlugin) GRPCServer(broker *goPlugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterDriverServer(s, dp.Srv)
	return nil
}

func (dp *auditDriverPlugin) GRPCClient(ctx context.Context, broker *goPlugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return proto.NewDriverClient(c), nil
}
