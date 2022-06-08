// Copyright 2019 TiKV Project Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package grpcutil

import (
	"context"
	"crypto/tls"
	"net/url"

	"github.com/pingcap/log"
	"github.com/tikv/pd/pkg/errs"
	"go.etcd.io/etcd/pkg/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

// ForwardMetadataKey is used to record the forwarded host of PD.
const ForwardMetadataKey = "pd-forwarded-host"

// TLSConfig is the configuration for supporting tls.
type TLSConfig struct {
	// CAPath is the path of file that contains list of trusted SSL CAs. if set, following four settings shouldn't be empty
	CAPath string `toml:"cacert-path" json:"cacert-path"`
	// CertPath is the path of file that contains X509 certificate in PEM format.
	CertPath string `toml:"cert-path" json:"cert-path"`
	// KeyPath is the path of file that contains X509 key in PEM format.
	KeyPath string `toml:"key-path" json:"key-path"`
	// CertAllowedCN is a CN which must be provided by a client
	CertAllowedCN []string `toml:"cert-allowed-cn" json:"cert-allowed-cn"`
}

// ToTLSConfig generates tls config.
func (s TLSConfig) ToTLSConfig() (*tls.Config, error) {
	if len(s.CertPath) == 0 && len(s.KeyPath) == 0 {
		return nil, nil
	}
	allowedCN, err := s.GetOneAllowedCN()
	if err != nil {
		return nil, err
	}

	tlsInfo := transport.TLSInfo{
		CertFile:      s.CertPath,
		KeyFile:       s.KeyPath,
		TrustedCAFile: s.CAPath,
		AllowedCN:     allowedCN,
	}

	tlsConfig, err := tlsInfo.ClientConfig()
	if err != nil {
		return nil, errs.ErrEtcdTLSConfig.Wrap(err).GenWithStackByCause()
	}
	return tlsConfig, nil
}

// GetOneAllowedCN only gets the first one CN.
func (s TLSConfig) GetOneAllowedCN() (string, error) {
	switch len(s.CertAllowedCN) {
	case 1:
		return s.CertAllowedCN[0], nil
	case 0:
		return "", nil
	default:
		return "", errs.ErrSecurityConfig.FastGenByArgs("only supports one CN")
	}
}

// GetClientConn returns a gRPC client connection.
// creates a client connection to the given target. By default, it's
// a non-blocking dial (the function won't wait for connections to be
// established, and connecting happens in the background). To make it a blocking
// dial, use WithBlock() dial option.
//
// In the non-blocking case, the ctx does not act against the connection. It
// only controls the setup steps.
//
// In the blocking case, ctx can be used to cancel or expire the pending
// connection. Once this function returns, the cancellation and expiration of
// ctx will be noop. Users should call ClientConn.Close to terminate all the
// pending operations after this function returns.
func GetClientConn(ctx context.Context, addr string, tlsCfg *tls.Config, do ...grpc.DialOption) (*grpc.ClientConn, error) {
	opt := grpc.WithInsecure()
	if tlsCfg != nil {
		creds := credentials.NewTLS(tlsCfg)
		opt = grpc.WithTransportCredentials(creds)
	}
	u, err := url.Parse(addr)
	if err != nil {
		return nil, errs.ErrURLParse.Wrap(err).GenWithStackByCause()
	}
	cc, err := grpc.DialContext(ctx, u.Host, append(do, opt)...)
	if err != nil {
		return nil, errs.ErrGRPCDial.Wrap(err).GenWithStackByCause()
	}
	return cc, nil
}

// BuildForwardContext creates a context with receiver metadata information.
// It is used in client side.
func BuildForwardContext(ctx context.Context, addr string) context.Context {
	md := metadata.Pairs(ForwardMetadataKey, addr)
	return metadata.NewOutgoingContext(ctx, md)
}

// ResetForwardContext is going to reset the forwarded host in metadata.
func ResetForwardContext(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Error("failed to get forwarding metadata")
	}
	md.Set(ForwardMetadataKey, "")
	return metadata.NewOutgoingContext(ctx, md)
}
