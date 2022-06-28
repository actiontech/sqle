//go:build enterprise
// +build enterprise

package driver

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func init() {
	SQLEGrpcServer = eeSQLEGrpcServer
}

var (
	eeSQLEGrpcServer = func(opts []grpc.ServerOption) *grpc.Server {
		opts = append(opts, grpc.ChainUnaryInterceptor(SQLEUnaryInterceptor))
		return grpc.NewServer(opts...)
	}
	// 默认无任何动作
	SQLEUnaryInterceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("metadata should not be empty")
		}
		if values := md.Get(ProjectKey); len(values) > 0 && values[0] == ProjectValue {
			return handler(ctx, req)
		}
		return nil, errors.New("the Enterprise Edition plugin only supports running on the Enterprise Edition SQLE")
	}
)
