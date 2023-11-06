package common

import (
	"math"

	"google.golang.org/grpc"
)

// custom GRPC client options
var GRPCDialOptions = []grpc.DialOption{}

// custom GRPC server options
var GRPCServerOptions = []grpc.ServerOption{
	grpc.MaxRecvMsgSize(math.MaxInt32),
}

func NewGRPCServer(opts []grpc.ServerOption) *grpc.Server {
	return grpc.NewServer(append(opts, GRPCServerOptions...)...)
}
