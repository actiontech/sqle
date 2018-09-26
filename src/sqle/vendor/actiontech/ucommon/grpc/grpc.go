package grpc

import (
	"actiontech/ucommon/secure"
	grpc_ "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"strings"
	"time"
)

func DialSocket(socket string) (*grpc_.ClientConn, error) {
	dialer := grpc_.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
		return net.DialTimeout("unix", addr, timeout)
	})
	return grpcDial(socket, grpc_.WithTimeout(1*time.Second), dialer)
}

func DialTcp(socket string) (*grpc_.ClientConn, error) {
	dialer := grpc_.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
		conn, err := net.DialTimeout("tcp4", addr, timeout)
		if nil != err {
			return nil, err
		}
		tcpConn, _ := conn.(*net.TCPConn)
		tcpConn.SetLinger(1)
		return tcpConn, nil
	})
	return grpcDial(socket, grpc_.WithTimeout(1*time.Second), dialer)
}

func Dial(target string) (*grpc_.ClientConn, error) {
	if strings.HasPrefix(target, "/") {
		return DialSocket(target)
	} else {
		return DialTcp(target)
	}
}

func grpcDial(target string, timeout grpc_.DialOption, dialer grpc_.DialOption) (*grpc_.ClientConn, error) {
	kpParam := keepalive.ClientParameters{
		Time: time.Second,
		Timeout: 3 * time.Second,
		PermitWithoutStream: true,
	}

	if secure.IsSecurityEnabled() {
		creds, err := secure.GetClientTLSCredentials()
		if err != nil {
			return nil, err
		}
		return grpc_.Dial(target, timeout, dialer, grpc_.WithTransportCredentials(creds),
			grpc_.WithKeepaliveParams(kpParam))
	} else {
		return grpc_.Dial(target, timeout, dialer, grpc_.WithInsecure(),
			grpc_.WithKeepaliveParams(kpParam))
	}
}
