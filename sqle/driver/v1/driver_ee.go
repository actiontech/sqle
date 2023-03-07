//go:build enterprise
// +build enterprise

package driver

import (
	"context"

	"google.golang.org/grpc"
)

func init() {
	SQLEGRPCDialOptions = append(SQLEGRPCDialOptions, grpc.WithPerRPCCredentials(new(sqleCredential)))
}

const (
	ProjectKey   = "project"
	ProjectValue = "action-sqle-ee-asdf" // asdf用于防止value被蒙到
)

// GrpcMetadata 用于存放希望发送给插件的信息
var GrpcMetadata = map[string]string{
	ProjectKey: ProjectValue,
}

type sqleCredential struct {
}

func (sqleCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return GrpcMetadata, nil
}

func (sqleCredential) RequireTransportSecurity() bool {
	return false
}
