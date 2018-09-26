package grpc

import (
	"fmt"
	"sync"

	grpc_ "google.golang.org/grpc"
)

func init() {

}

var grpcConnsMu sync.RWMutex
var grpcConns map[string]*grpc_.ClientConn = make(map[string]*grpc_.ClientConn)

func GetRpcConnFromPool(ip, port string, connectableCheck func(conn *grpc_.ClientConn) error) (*grpc_.ClientConn, error) {
	grpcConnsMu.RLock()
	addr := ip + ":" + port
	conn := grpcConns[addr]
	grpcConnsMu.RUnlock()

	//health check, otherwise client will always get a closed connection
	if nil != conn {
		if err := connectableCheck(conn); nil != err {
			grpcConnsMu.Lock()
			if grpcConns[addr] == conn {
				grpcConns[addr] = nil
			}
			conn.Close()
			conn = nil
			grpcConnsMu.Unlock()
		}
	}

	if nil == conn {
		grpcConnsMu.Lock()
		conn = grpcConns[addr]
		if nil == conn {
			if c, err := Dial(addr); nil != err {
				grpcConnsMu.Unlock()
				return nil, fmt.Errorf("grpc connect to %v error: %v", addr, err)
			} else {
				grpcConns[addr] = c
				conn = c
			}
		}
		grpcConnsMu.Unlock()
	}
	return conn, nil
}
