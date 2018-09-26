package pgrpc

/*
	pgrpc is passive-grpc, which means server dial client, and then client invoke server's procedure.
	It's different from grpc, in which client dial server, and client invoke server's procedure.
*/

import (
	"actiontech/ucommon/log"
	"actiontech/ucommon/secure"
	"bytes"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	grace "github.com/facebookgo/grace/gracenet"
	"golang.org/x/net/context"
	grpc_ "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

/*
	clientd keep connections from server, and provide connections to grpc.Dial
*/
type tClientd struct {
	listener    net.Listener
	conns       map[string](chan net.Conn)
	addrMapping map[string]string
	mu          sync.Mutex
	acceptQuit  chan bool
}

func newClientd(l net.Listener) *tClientd {
	return &tClientd{
		listener:    l,
		conns:       make(map[string](chan net.Conn)),
		acceptQuit:  make(chan bool, 1),
		addrMapping: map[string]string{},
	}
}

type pgrpcTransportCredentials struct {
	inner   credentials.TransportCredentials
	clientd *tClientd
}

func (p *pgrpcTransportCredentials) ClientHandshake(ctx context.Context, addr string, rawConn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	p.clientd.mu.Lock()
	realAddr := p.clientd.addrMapping[addr]
	p.clientd.mu.Unlock()
	return p.inner.ClientHandshake(ctx, realAddr, rawConn)
}

func (p *pgrpcTransportCredentials) ServerHandshake(rawConn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	return p.inner.ServerHandshake(rawConn)
}

func (p *pgrpcTransportCredentials) Info() credentials.ProtocolInfo {
	return p.inner.Info()
}

func (p *pgrpcTransportCredentials) Clone() credentials.TransportCredentials {
	copy := p.inner.Clone()
	return &pgrpcTransportCredentials{
		inner:   copy,
		clientd: p.clientd,
	}
}

func (p *pgrpcTransportCredentials) OverrideServerName(a string) error {
	return p.inner.OverrideServerName(a)
}

func (c *tClientd) buildTransportCredentials() (*pgrpcTransportCredentials, error) {
	if !secure.IsSecurityEnabled() {
		return nil, fmt.Errorf("secure is disabled")
	}
	creds, err := secure.GetClientTLSCredentials()
	if err != nil {
		return nil, err
	}
	ret := &pgrpcTransportCredentials{
		inner:   creds,
		clientd: c,
	}
	return ret, nil
}

func (c *tClientd) start() {
	go c.accept()
}

func (c *tClientd) stop() {
	c.listener.Close()
	<-c.acceptQuit

	c.mu.Lock()
	for _, v := range c.conns {
		if nil == v {
			continue
		}
		select {
		case conn := <-v:
			conn.Close()
		default:
		}
	}
	c.mu.Unlock()
}

func (c *tClientd) accept() {
	defer func() {
		c.acceptQuit <- true
	}()

	stage := log.NewStage().Enter("pgrpc_clientd_accept")
	defer stage.Exit()

	for {
		conn, err := c.listener.Accept()
		if nil != err {
			if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
				log.Detail(stage, "accept temporary error: %v", err)
				continue
			}
			log.Detail(stage, "accept error: %v", err)
			return
		}

		keyBs := make([]byte, UADDR_LEN)

		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		if _, err := io.ReadFull(conn, keyBs); nil != err {
			conn.Close()
			log.Detail(stage, "accept read key error: %v", err)
			continue
		}
		conn.SetReadDeadline(time.Time{})

		key := string(bytes.Trim(keyBs, "\x00"))
		addr := conn.RemoteAddr().String()
		log.Detail(stage, "accept new connection from %v %v", key, addr)

		c.mu.Lock()
		if nil == c.conns[key] {
			c.conns[key] = make(chan net.Conn, MAX_CONNS_PER_UADDR)
		}
		c.addrMapping[key] = addr
		connPool := c.conns[key]
		c.mu.Unlock()

	RETRY:
		select {
		case connPool <- conn:
		default:
			select {
			case oldConn := <-connPool:
				oldConn.Close()
				goto RETRY
			default:
				conn.Close()
			}
		}
	}
}

func (c *tClientd) getConn(addr string, timeout time.Duration) (net.Conn, error) {
	stage := log.NewStage().Enter("pgrpc_dialer")
	c.mu.Lock()
	connPool := c.conns[addr]
	c.mu.Unlock()

	if nil == connPool {
		log.Detail(stage, "connect to %v error: no connection exists", addr)
		return nil, fmt.Errorf("pgrpc error: no connection exists")
	}

	select {
	case c := <-connPool:
		return c, nil
	case <-time.After(timeout):
		log.Detail(stage, "connect to %v error: wait pgrpc connection timeout", addr)
		return nil, fmt.Errorf("pgrpc error: wait pgrpc connection timeout")
	}
}

var clientd *tClientd
var clientdMu sync.Mutex

func StartClientd(graceNet *grace.Net, port int) error {
	l, err := graceNet.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if nil != err {
		return err
	}
	clientdMu.Lock()
	if nil != clientd {
		clientdMu.Unlock()
		return fmt.Errorf("pgrpc error: clientd already exists")
	}
	clientd = newClientd(l)
	c := clientd
	clientdMu.Unlock()

	go c.start()
	return nil
}

func StopClientd() {
	clientdMu.Lock()
	c := clientd
	clientd = nil
	clientdMu.Unlock()

	if nil != c {
		c.stop()
	}
}

func PgrpcDialer(addr string, timeout time.Duration) (net.Conn, error) {
	clientdMu.Lock()
	c := clientd
	clientdMu.Unlock()

	if nil == c {
		return nil, fmt.Errorf("pgrpc is stopping")
	}
	return c.getConn(addr, timeout)
}

func PgrpcTransportCredentials() (credentials.TransportCredentials, error) {
	clientdMu.Lock()
	c := clientd
	clientdMu.Unlock()

	if nil == c {
		return nil, fmt.Errorf("pgrpc is stopping")
	}
	return c.buildTransportCredentials()
}

func Dial(addr string) (*grpc_.ClientConn, error) {
	dialer := grpc_.WithDialer(PgrpcDialer)

	kpParam := keepalive.ClientParameters{
		Time:                time.Second,
		Timeout:             3 * time.Second,
		PermitWithoutStream: true,
	}

	if secure.IsSecurityEnabled() {
		creds, err := PgrpcTransportCredentials()
		if err != nil {
			return nil, err
		}
		return grpc_.Dial(addr, grpc_.WithTimeout(1*time.Second), dialer, grpc_.WithTransportCredentials(creds),
			grpc_.WithKeepaliveParams(kpParam))
	} else {
		return grpc_.Dial(addr, grpc_.WithTimeout(1*time.Second), dialer, grpc_.WithInsecure(),
			grpc_.WithKeepaliveParams(kpParam))
	}
}
