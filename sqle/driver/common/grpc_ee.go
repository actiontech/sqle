//go:build enterprise || trial
// +build enterprise trial

package common

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func init() {
	GRPCDialOptions = append(GRPCDialOptions, grpc.WithPerRPCCredentials(new(EECheckCredential)))
	GRPCServerOptions = append(GRPCServerOptions, grpc.ChainUnaryInterceptor(EECheckInterceptor))
}

const (
	ProjectKey   = "project"
	ProjectValue = "action-sqle-ee-asdf" // asdf用于防止value被蒙到
)

// GrpcMetadata 用于存放希望发送给插件的信息
var GrpcMetadata = map[string]string{
	ProjectKey: ProjectValue,
}

type EECheckCredential struct{}

func (EECheckCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return GrpcMetadata, nil
}

func (EECheckCredential) RequireTransportSecurity() bool {
	return false
}

var EECheckInterceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	err = errors.New("the Enterprise Edition plugin only supports running on the Enterprise Edition SQLE")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, err
	}

	if values := md.Get(ProjectKey); len(values) > 0 && values[0] == ProjectValue {
		return handler(ctx, req)
	}
	return nil, err
}
