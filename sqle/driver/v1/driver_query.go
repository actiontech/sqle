package v1

import (
	"context"

	"github.com/actiontech/sqle/sqle/driver/v1/proto"
	"github.com/actiontech/sqle/sqle/pkg/params"
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

// queryDriverImpl implement SQLQueryDriver. It use for hide gRPC detail, just like DriverGRPCServer.
type queryDriverImpl struct {
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

func (q *queryDriverImpl) Query(ctx context.Context, sql string, conf *QueryConf) (*QueryResult, error) {
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
