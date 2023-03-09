package common

import (
	"google.golang.org/grpc"
)

// custom GRPC client options
var GRPCDialOptions = []grpc.DialOption{}

// custom GRPC server options
var GRPCServerOptions = []grpc.ServerOption{}

func NewGRPCServer(opts []grpc.ServerOption) *grpc.Server {
	return grpc.NewServer(append(opts, GRPCServerOptions...)...)
}
