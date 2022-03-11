package driver

import (
	"context"
	"database/sql/driver"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/actiontech/sqle/sqle/driver/proto"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"
	goPlugin "github.com/hashicorp/go-plugin"
	"github.com/pingcap/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func convertRuleFromProtoToDriver(rule *proto.Rule) *Rule {
	var ps = make(params.Params, 0, len(rule.Params))
	for _, p := range rule.Params {
		ps = append(ps, &params.Param{
			Key:   p.Key,
			Value: p.Value,
			Desc:  p.Desc,
			Type:  params.ParamType(p.Type),
		})
	}
	return &Rule{
		Name:     rule.Name,
		Category: rule.Category,
		Desc:     rule.Desc,
		Level:    RuleLevel(rule.Level),
		Params:   ps,
	}
}

func convertRuleFromDriverToProto(rule *Rule) *proto.Rule {
	var params = make([]*proto.Param, 0, len(rule.Params))
	for _, p := range rule.Params {
		params = append(params, &proto.Param{
			Key:   p.Key,
			Value: p.Value,
			Desc:  p.Desc,
			Type:  string(p.Type),
		})
	}
	return &proto.Rule{
		Name:     rule.Name,
		Desc:     rule.Desc,
		Level:    string(rule.Level),
		Category: rule.Category,
		Params:   params,
	}
}

// InitPlugins init plugins at plugins directory. It should be called on host process.
func InitPlugins(pluginDir string) error {
	if pluginDir == "" {
		return nil
	}

	getServerHandle := func(path string, closeCh <-chan struct{}) (proto.DriverClient, error) {
		client := goPlugin.NewClient(&goPlugin.ClientConfig{
			HandshakeConfig: handshakeConfig,
			Plugins: goPlugin.PluginSet{
				filepath.Base(path): &driverPlugin{},
			},
			Cmd:              exec.Command(path),
			AllowedProtocols: []goPlugin.Protocol{goPlugin.ProtocolGRPC},
		})
		go func() {
			<-closeCh
			client.Kill()
		}()

		gRPCClient, err := client.Client()
		if err != nil {
			return nil, err
		}
		rawI, err := gRPCClient.Dispense(filepath.Base(path))
		if err != nil {
			return nil, err
		}
		// srv can only be proto.DriverClient
		//nolint:forcetypeassert
		srv := rawI.(proto.DriverClient)
		return srv, nil
	}

	var plugins []os.FileInfo
	if err := filepath.Walk(pluginDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "init plugin")
		}

		if info.IsDir() || info.Mode()&0111 == 0 {
			return nil
		}
		plugins = append(plugins, info)
		return nil
	}); err != nil {
		return err
	}

	for _, p := range plugins {
		binaryPath := filepath.Join(pluginDir, p.Name())

		closeCh := make(chan struct{})
		srv, err := getServerHandle(binaryPath, closeCh)
		if err != nil {
			return err
		}
		pluginMeta, err := srv.Metas(context.TODO(), &proto.Empty{})
		if err != nil {
			return err
		}
		close(closeCh)

		// driverRules get from plugin when plugin initialize.
		var driverRules = make([]*Rule, 0, len(pluginMeta.Rules))
		for _, rule := range pluginMeta.Rules {
			driverRules = append(driverRules, convertRuleFromProtoToDriver(rule))
		}

		handler := func(log *logrus.Entry, config *Config) (Driver, error) {
			pluginCloseCh := make(chan struct{})
			srv, err := getServerHandle(binaryPath, pluginCloseCh)
			if err != nil {
				return nil, err
			}

			// protoRules send to plugin for Audit.
			var protoRules []*proto.Rule
			for _, rule := range config.Rules {
				protoRules = append(protoRules, convertRuleFromDriverToProto(rule))
			}

			initRequest := &proto.InitRequest{
				Rules: protoRules,
			}
			if config.DSN != nil {
				initRequest.Dsn = &proto.DSN{
					Host:             config.DSN.Host,
					Port:             config.DSN.Port,
					User:             config.DSN.User,
					Password:         config.DSN.Password,
					AdditionalParams: proto.ConvertParamToProtoParam(config.DSN.AdditionalParams),

					// database is to open.
					Database: config.DSN.DatabaseName,
				}
			}

			_, err = srv.Init(context.TODO(), initRequest)
			if err != nil {
				return nil, err
			}
			return &driverPluginClient{srv, pluginCloseCh}, nil

		}

		Register(pluginMeta.Name, handler, driverRules, proto.ConvertProtoParamToParam(pluginMeta.GetAdditionalParams()))

		log.Logger().WithFields(logrus.Fields{
			"plugin_name": pluginMeta.Name,
		}).Infoln("plugin inited")
	}

	return nil
}

// ServePlugin start plugin process service. It should be called on plugin process.
func ServePlugin(r Registerer, newDriver func(cfg *Config) Driver) {
	name := r.Name()
	goPlugin.Serve(&goPlugin.ServeConfig{
		HandshakeConfig: handshakeConfig,

		Plugins: goPlugin.PluginSet{
			name: &driverPlugin{Srv: &driverGRPCServer{r: r, newDriver: newDriver}},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: goPlugin.DefaultGRPCServer,
	})
}

var handshakeConfig = goPlugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// driverPluginClient implement Driver. It use for hide gRPC detail, just like DriverGRPCServer.
type driverPluginClient struct {
	plugin proto.DriverClient

	// driverQuitCh produce a singal for telling caller that it's time to Client.Kill() plugin process.
	driverQuitCh chan struct{}
}

func (s *driverPluginClient) Close(ctx context.Context) {
	s.plugin.Close(ctx, &proto.Empty{})
	close(s.driverQuitCh)
}

func (s *driverPluginClient) Ping(ctx context.Context) error {
	_, err := s.plugin.Ping(ctx, &proto.Empty{})
	return err
}

type dbDriverResult struct {
	lastInsertId    int64
	lastInsertIdErr string
	rowsAffected    int64
	rowsAffectedErr string
}

func (s *dbDriverResult) LastInsertId() (int64, error) {
	if s.lastInsertIdErr != "" {
		return s.lastInsertId, fmt.Errorf(s.lastInsertIdErr)
	}
	return s.lastInsertId, nil
}

func (s *dbDriverResult) RowsAffected() (int64, error) {
	if s.rowsAffectedErr != "" {
		return s.rowsAffected, fmt.Errorf(s.rowsAffectedErr)
	}
	return s.rowsAffected, nil
}

func (s *driverPluginClient) Exec(ctx context.Context, query string) (driver.Result, error) {
	resp, err := s.plugin.Exec(ctx, &proto.ExecRequest{Query: query})
	if err != nil {
		return nil, err
	}
	return &dbDriverResult{
		lastInsertId:    resp.LastInsertId,
		lastInsertIdErr: resp.LastInsertIdError,
		rowsAffected:    resp.RowsAffected,
		rowsAffectedErr: resp.RowsAffectedError,
	}, nil
}

func (s *driverPluginClient) Tx(ctx context.Context, queries ...string) ([]driver.Result, error) {
	resp, err := s.plugin.Tx(ctx, &proto.TxRequest{Queries: queries})
	if err != nil {
		return nil, err
	}

	ret := make([]driver.Result, len(resp.Results))
	for i, result := range resp.Results {
		ret[i] = &dbDriverResult{
			lastInsertId:    result.LastInsertId,
			lastInsertIdErr: result.LastInsertIdError,
			rowsAffected:    result.RowsAffected,
			rowsAffectedErr: result.RowsAffectedError,
		}
	}
	return ret, nil
}

func (s *driverPluginClient) Schemas(ctx context.Context) ([]string, error) {
	resp, err := s.plugin.Databases(ctx, &proto.Empty{})
	if err != nil {
		return nil, err
	}
	return resp.Databases, nil
}

func (s *driverPluginClient) Parse(ctx context.Context, sqlText string) ([]Node, error) {
	resp, err := s.plugin.Parse(ctx, &proto.ParseRequest{SqlText: sqlText})
	if err != nil {
		return nil, err
	}

	nodes := make([]Node, len(resp.Nodes))
	for i, node := range resp.Nodes {
		nodes[i] = Node{
			Type:        node.Type,
			Text:        node.Text,
			Fingerprint: node.Fingerprint,
		}
	}
	return nodes, nil
}

func (s *driverPluginClient) Audit(ctx context.Context, sql string) (*AuditResult, error) {
	resp, err := s.plugin.Audit(ctx, &proto.AuditRequest{Sql: sql})
	if err != nil {
		return nil, err
	}

	ret := &AuditResult{}
	for _, result := range resp.Results {
		ret.results = append(ret.results, &auditResult{
			level:   RuleLevel(result.Level),
			message: result.Message,
		})
	}
	return ret, nil
}

func (s *driverPluginClient) GenRollbackSQL(ctx context.Context, sql string) (string, string, error) {
	resp, err := s.plugin.GenRollbackSQL(ctx, &proto.GenRollbackSQLRequest{Sql: sql})
	if err != nil {
		return "", "", err
	}

	return resp.Sql, resp.Reason, nil
}

// driverPlugin use for hide gRPC detail.
type driverGRPCServer struct {
	newDriver func(cfg *Config) Driver

	impl Driver

	// Registerer provide some plugin info to host process.
	r Registerer
}

func (d *driverGRPCServer) Init(ctx context.Context, req *proto.InitRequest) (*proto.Empty, error) {
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

func (d *driverGRPCServer) Close(ctx context.Context, req *proto.Empty) (*proto.Empty, error) {
	d.impl.Close(ctx)
	return &proto.Empty{}, nil
}

func (d *driverGRPCServer) Ping(ctx context.Context, req *proto.Empty) (*proto.Empty, error) {
	return &proto.Empty{}, d.impl.Ping(ctx)
}

func (d *driverGRPCServer) Exec(ctx context.Context, req *proto.ExecRequest) (*proto.ExecResponse, error) {
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

func (d *driverGRPCServer) Tx(ctx context.Context, req *proto.TxRequest) (*proto.TxResponse, error) {
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

func (d *driverGRPCServer) Databases(ctx context.Context, req *proto.Empty) (*proto.DatabasesResponse, error) {
	databases, err := d.impl.Schemas(ctx)
	return &proto.DatabasesResponse{Databases: databases}, err
}

func (d *driverGRPCServer) Parse(ctx context.Context, req *proto.ParseRequest) (*proto.ParseResponse, error) {
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

func (d *driverGRPCServer) Audit(ctx context.Context, req *proto.AuditRequest) (*proto.AuditResponse, error) {
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

func (d *driverGRPCServer) GenRollbackSQL(ctx context.Context, req *proto.GenRollbackSQLRequest) (*proto.GenRollbackSQLResponse, error) {
	rollbackSQL, reason, err := d.impl.GenRollbackSQL(ctx, req.GetSql())
	return &proto.GenRollbackSQLResponse{
		Sql:    rollbackSQL,
		Reason: reason,
	}, err
}

func (d *driverGRPCServer) Metas(ctx context.Context, req *proto.Empty) (*proto.MetasResponse, error) {
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

// driverPlugin implements goPlugin.GRPCPlugin
type driverPlugin struct {
	goPlugin.NetRPCUnsupportedPlugin

	Srv *driverGRPCServer
}

func (dp *driverPlugin) GRPCServer(broker *goPlugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterDriverServer(s, dp.Srv)
	return nil
}

func (dp *driverPlugin) GRPCClient(ctx context.Context, broker *goPlugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return proto.NewDriverClient(c), nil
}
