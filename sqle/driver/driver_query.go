package driver

import (
	"context"
	"fmt"
	"sync"

	"github.com/actiontech/sqle/sqle/driver/proto"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"

	goPlugin "github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
)

// SQLQueryDriver is a SQL rewrite and execute driver
type SQLQueryDriver interface {
	QueryPrepare(ctx context.Context, sql string, conf *QueryPrepareConf) (*QueryPrepareResult, error)
	Query(ctx context.Context, sql string, conf *QueryConf) (*QueryResult, error)
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
}

func (q *queryDriverImpl) QueryPrepare(ctx context.Context, sql string, conf *QueryPrepareConf) (*QueryPrepareResult, error) {
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
var queryDrivers = make(map[string]struct{})

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
func RegisterSQLQueryDriver(name string) {
	queryDriverMu.RLock()
	_, exist := queryDrivers[name]
	queryDriverMu.RUnlock()
	if exist {
		panic(fmt.Sprintf("duplicated driver name %v", name))
	}

	queryDriverMu.Lock()
	queryDrivers[name] = struct{}{}
	queryDriverMu.Unlock()
}

func registerQueryDriver(pluginName string, gRPCClient goPlugin.ClientProtocol) error {
	rawI, err := gRPCClient.Dispense(PluginNameQueryDriver)
	if err != nil {
		return err
	}
	// srv can only be proto.QueryDriverClient
	//nolint:forcetypeassert
	s := rawI.(proto.QueryDriverClient)

	// The test target plugin implements the QueryDriver plugin
	_, err = s.Init(context.TODO(), &proto.InitRequest{})
	if err != nil {
		return err
	}

	RegisterSQLQueryDriver(pluginName)
	log.Logger().WithFields(logrus.Fields{
		"plugin_name": pluginName,
		"plugin_type": PluginNameQueryDriver,
	}).Infoln("plugin inited")
	return nil
}
