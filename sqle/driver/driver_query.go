package driver

import (
	"context"
	"fmt"
	"sync"

	"github.com/actiontech/sqle/sqle/driver/proto"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"

	"github.com/sirupsen/logrus"
)

// SQLQueryDriver is a SQL rewrite and execute driver
type SQLQueryDriver interface {
	QueryPrepare(ctx context.Context, sql string, conf *QueryPrepareConf) (*QueryPrepareResult, error)
	Query(ctx context.Context, sql string, conf *QueryConf) (*QueryResult, error)
	Close(ctx context.Context)
}

type ErrorType string

const (
	ErrorTypeNotQuery = "not query"
	ErrorTypeNotError = "not error"
)

func init() {
	defaultPluginSet[DefaultPluginVersion][PluginNameQueryDriver] = &queryDriverPlugin{}
}

const (
	PluginNameQueryDriver = "query-driver"
)

type QueryPrepareConf struct {
	Limit  uint32
	Offset uint32
}

type QueryPrepareResult struct {
	NewSQL    string
	ErrorType ErrorType
	Error     string
}

type QueryConf struct {
	TimeOutSecond uint32
}

// The data location in Values should be consistent with that in Column
type QueryResult struct {
	Column params.Params
	Rows   []*QueryResultRow
}

type QueryResultRow struct {
	Values []*QueryResultValue
}

type QueryResultValue struct {
	Value string
}

// queryDriverPluginClient implement SQLQueryDriver. It use for hide gRPC detail, just like DriverGRPCServer.
type queryDriverPluginClient struct {
	plugin proto.QueryDriverClient

	// driverQuitCh produce a singal for telling caller that it's time to Client.Kill() plugin process.
	driverQuitCh chan struct{}
}

func (q *queryDriverPluginClient) Close(ctx context.Context) {
	close(q.driverQuitCh)
}

func (q *queryDriverPluginClient) QueryPrepare(ctx context.Context, sql string, conf *QueryPrepareConf) (*QueryPrepareResult, error) {
	req := &proto.QueryPrepareRequest{
		Sql: sql,
		Conf: &proto.QueryPrepareConf{
			Limit:  conf.Limit,
			Offset: conf.Offset,
		},
	}
	res, err := q.plugin.QueryPrepare(ctx, req)
	if err != nil {
		return nil, err
	}
	return &QueryPrepareResult{
		NewSQL:    res.GetNewSql(),
		ErrorType: ErrorType(res.GetErrorType()),
		Error:     res.GetError(),
	}, nil
}

func (q *queryDriverPluginClient) Query(ctx context.Context, sql string, conf *QueryConf) (*QueryResult, error) {
	req := &proto.QueryRequest{
		Sql: sql,
		Conf: &proto.QueryConf{
			TimeOutSecond: conf.TimeOutSecond,
		},
	}
	res, err := q.plugin.Query(ctx, req)
	if err != nil {
		return nil, err
	}
	result := &QueryResult{
		Column: params.Params{},
		Rows:   []*QueryResultRow{},
	}
	for _, p := range res.GetColumn() {
		result.Column = append(result.Column, &params.Param{
			Key:   p.GetKey(),
			Value: p.GetValue(),
			Desc:  p.GetDesc(),
			Type:  params.ParamType(p.GetType()),
		})
	}
	for _, row := range res.GetRows() {
		r := &QueryResultRow{
			Values: []*QueryResultValue{},
		}
		for _, value := range row.GetValues() {
			r.Values = append(r.Values, &QueryResultValue{
				Value: value.GetValue(),
			})
		}
		result.Rows = append(result.Rows, r)
	}
	return result, nil
}

var queryDriverMu = &sync.RWMutex{}
var queryDrivers = make(map[string]queryHandler)

// queryHandler is a template which SQLQueryDriver plugin should provide such function signature.
type queryHandler func(log *logrus.Entry, c *DSN) (SQLQueryDriver, error)

// NewSQLQueryDriver return a new instantiated SQLQueryDriver.
func NewSQLQueryDriver(log *logrus.Entry, dbType string, cfg *DSN) (SQLQueryDriver, error) {
	queryDriverMu.RLock()
	defer queryDriverMu.RUnlock()
	d, exist := queryDrivers[dbType]
	if !exist {
		return nil, fmt.Errorf("driver type %v is not supported", dbType)
	}
	return d(log, cfg)
}

// QueryDriverName = InstanceType
func GetQueryDriverNames() []string {
	queryDriverMu.RLock()
	defer queryDriverMu.RUnlock()
	names := []string{}
	for s := range queryDrivers {
		names = append(names, s)
	}
	return names
}

// RegisterSQLQueryDriver like sql.RegisterAuditDriver.
//
// RegisterSQLQueryDriver makes a database driver available by the provided driver name.
// SQLQueryDriver's initialize handler and audit rules register by RegisterSQLQueryDriver.
func RegisterSQLQueryDriver(name string, h queryHandler) {
	queryDriverMu.RLock()
	_, exist := queryDrivers[name]
	queryDriverMu.RUnlock()
	if exist {
		panic("duplicated driver name")
	}

	queryDriverMu.Lock()
	queryDrivers[name] = h
	queryDriverMu.Unlock()
}

func registerQueryPlugin(pluginName string, c PluginClient) error {
	closeCh := make(chan struct{})
	s, err := getQueryServerHandle(c, closeCh)
	if err != nil {
		return err
	}
	// The test target plugin implements the QueryDriver plugin
	_, err = s.Init(context.TODO(), &proto.InitRequest{})
	close(closeCh)
	if err != nil {
		return err
	}
	handler := func(log *logrus.Entry, config *DSN) (SQLQueryDriver, error) {
		pluginCloseCh := make(chan struct{})
		srv, err := getQueryServerHandle(c, pluginCloseCh)
		if err != nil {
			return nil, err
		}

		initRequest := &proto.InitRequest{
			Rules: []*proto.Rule{},
		}
		if config != nil {
			initRequest.Dsn = &proto.DSN{
				Host:             config.Host,
				Port:             config.Port,
				User:             config.User,
				Password:         config.Password,
				AdditionalParams: proto.ConvertParamToProtoParam(config.AdditionalParams),

				// database is to open.
				Database: config.DatabaseName,
			}
		}

		_, err = srv.Init(context.TODO(), initRequest)
		if err != nil {
			return nil, err
		}

		return &queryDriverPluginClient{srv, pluginCloseCh}, nil

	}

	RegisterSQLQueryDriver(pluginName, handler)
	log.Logger().WithFields(logrus.Fields{
		"plugin_name": pluginName,
		"plugin_type": PluginNameQueryDriver,
	}).Infoln("plugin inited")
	return nil
}

func getQueryServerHandle(client PluginClient, closeCh <-chan struct{}) (proto.QueryDriverClient, error) {
	gRPCClient, err := client.Client()
	if err != nil {
		return nil, err
	}
	go func() {
		<-closeCh
		client.Kill()
	}()

	rawI, err := gRPCClient.Dispense(PluginNameQueryDriver)
	if err != nil {
		return nil, err
	}
	// srv can only be proto.QueryDriverClient
	//nolint:forcetypeassert
	srv := rawI.(proto.QueryDriverClient)

	return srv, nil
}
