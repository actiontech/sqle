package pgrpc

/*
	pgrpc is passive-grpc, which means server dial client, and then client invoke server's procedure.
	It's different from grpc, in which client dial server, and client invoke server's procedure.
*/

import (
	"actiontech/ucommon/log"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

/*
	serverd dials client, and provide connections to grpc.Listen
*/

type tServerd struct {
	mu         sync.Mutex
	connCounts map[ /*ip:port*/ string]int

	// addClient(1) && removeClient(1) && addClient(1), should have 2 version, connCounts is isolated in each version
	connsVersion map[ /*ip:port*/ string]int64

	newConns map[ /*ip:port*/ string](chan net.Conn)
	stopping bool
	dialQuit chan struct{}
	closed   chan struct{}
}

func newServerd() *tServerd {
	return &tServerd{
		connCounts:   make(map[string]int),
		connsVersion: make(map[string]int64),
		newConns:     make(map[string](chan net.Conn)),
		dialQuit:     make(chan struct{}, 0),
		closed:       make(chan struct{}, 0),
	}
}

func (s *tServerd) start(localKey string) {
	localKeyBs := make([]byte, UADDR_LEN)
	copy(localKeyBs, []byte(localKey))

	go s.dial(localKeyBs)
}

func (s *tServerd) stop() {
	s.mu.Lock()
	if s.stopping {
		s.mu.Unlock()
		return
	}
	s.stopping = true
	s.mu.Unlock()

	<-s.dialQuit

	s.mu.Lock()
	for _, chs := range s.newConns {
	LOOP:
		for {
			select {
			case conn := <-chs:
				go conn.Close()
			default:
				break LOOP
			}
		}
	}
	s.mu.Unlock()

	close(s.closed)
}

func (s *tServerd) dial(localKeyBs []byte) {
	defer func() {
		close(s.dialQuit)
	}()

	stage := log.NewStage().Enter("pgrpc_serverd_dial")

	for {
		s.mu.Lock()
		if s.stopping {
			s.mu.Unlock()
			return
		}

		for k, c := range s.connCounts {
			if c >= MAX_CONNS_PER_UADDR {
				continue
			}

			conn, err := net.Dial("tcp", k)
			if nil != err {
				log.Detail(stage, "dial error: %v", err)
				continue
			}
			{
				tcpConn, _ := conn.(*net.TCPConn)
				tcpConn.SetLinger(1)

				/*
					pgrpc still need tcp-layer keepalive
					tcp.Dial is called by pgrpc.Server, and if pgrpc.Client didn't use this connection in pgrpc-layer,
					the connection has no pgrpc-layer keepalive
				 */
				tcpConn.SetKeepAlive(true)
				tcpConn.SetKeepAlivePeriod(1 * time.Second)
			}

			if n, err := conn.Write(localKeyBs); nil != err || n < UADDR_LEN {
				go conn.Close()
				log.Detail(stage, "write key error: %v", err)
				continue
			}
			log.Detail(stage, "new connection to %v", k)

			s.connCounts[k]++
			s.newConns[k] <- newConnWithCloseSignal(conn, s, k, s.connsVersion[k])
		}
		s.mu.Unlock()
		time.Sleep(100 * time.Millisecond)
	}
}

func (s *tServerd) onConnClose(k string, version int64) {
	s.mu.Lock()
	if s.connsVersion[k] == version && s.connCounts[k] > 0 {
		s.connCounts[k]--
	}
	s.mu.Unlock()
}

type reconnectError string

func (r reconnectError) Error() string {
	return string(r)
}

func (s *tServerd) addClient(ipPort string) error {
	s.mu.Lock()
	if nil != s.newConns[ipPort] {
		s.mu.Unlock()
		return reconnectError(fmt.Sprintf("pgrpc error: serverd already has client %v", ipPort))
	}
	s.newConns[ipPort] = make(chan net.Conn, MAX_CONNS_PER_UADDR)
	s.connsVersion[ipPort] = time.Now().UnixNano()
	s.connCounts[ipPort] = 0
	s.mu.Unlock()
	return nil
}

func (s *tServerd) removeClient(ipPort string) error {
	s.mu.Lock()
	if nil == s.newConns[ipPort] {
		s.mu.Unlock()
		return nil
	}
	delete(s.connCounts, ipPort)
	chs := s.newConns[ipPort]
	delete(s.newConns, ipPort)
	s.mu.Unlock()

	for {
		select {
		case conn := <-chs:
			go conn.Close()
		default:
			return nil
		}
	}
}

//net.Listener impl
func (s *tServerd) Accept() (net.Conn, error) {
	for {
		s.mu.Lock()
		if s.stopping {
			s.mu.Unlock()
			return nil, fmt.Errorf("pgrpc: listener closed")
		}

		for _, ch := range s.newConns {
			select {
			case conn := <-ch:
				s.mu.Unlock()
				return conn, nil
			default:
			}
		}
		s.mu.Unlock()

		time.Sleep(100 * time.Millisecond)
	}
}

//net.Listener impl
func (s *tServerd) Close() error {
	s.stop()
	return nil
}

//net.Listener impl
func (s *tServerd) Addr() net.Addr {
	a, _ := net.ResolveIPAddr("tcp", "0.0.0.0:0") //mock addr
	return a
}

type connWithCloseSignal struct {
	conn net.Conn
	key  string

	closeChan  chan struct{}
	firstClose sync.Once
	serverd    *tServerd

	firstRead    int32 //atomic
	connsVersion int64
}

func newConnWithCloseSignal(conn net.Conn, serverd *tServerd, key string, connsVersion int64) *connWithCloseSignal {
	c := &connWithCloseSignal{
		conn:         conn,
		serverd:      serverd,
		key:          key,
		closeChan:    make(chan struct{}, 0),
		connsVersion: connsVersion,
	}
	go func() {
		select {
		case <-c.closeChan:
		case <-serverd.closed:
			//when listener is closed, should close pgrpc connections which's waiting for first read
			//otherwise, these connections, are not considered as grpc-live connection, and will block grpc.GracefulStop()
			if 0 == atomic.LoadInt32(&c.firstRead) {
				c.Close()
			}
		}
	}()
	return c
}

func (c *connWithCloseSignal) Read(b []byte) (n int, err error) {
	n, err = c.conn.Read(b)
	atomic.StoreInt32(&c.firstRead, 1)
	return n, err
}

func (c *connWithCloseSignal) Write(b []byte) (n int, err error) {
	return c.conn.Write(b)
}

func (c *connWithCloseSignal) Close() error {
	c.firstClose.Do(func() {
		c.serverd.onConnClose(c.key, c.connsVersion)
		close(c.closeChan)
	})
	return c.conn.Close()
}

func (c *connWithCloseSignal) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *connWithCloseSignal) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *connWithCloseSignal) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *connWithCloseSignal) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *connWithCloseSignal) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

var serverd *tServerd
var serverdMu sync.Mutex

func IsServerdStarted() bool {
	serverdMu.Lock()
	s := serverd
	serverdMu.Unlock()
	return nil != s
}

func StartServerd(localKey string) error {
	serverdMu.Lock()
	if nil != serverd {
		serverdMu.Unlock()
		return fmt.Errorf("pgrpc error: serverd already exists")
	}
	serverd = newServerd()
	s := serverd
	serverdMu.Unlock()

	s.start(localKey)
	return nil
}

func StopServerd() {
	serverdMu.Lock()
	s := serverd
	serverd = nil
	serverdMu.Unlock()

	if nil != s {
		s.stop()
	}
}

func AddClient(ipPort string) error {
	serverdMu.Lock()
	s := serverd
	serverdMu.Unlock()

	if nil == s {
		return fmt.Errorf("pgrpc error: no serverd")
	}
	if err := s.addClient(ipPort); nil != err {
		if _, ok := err.(reconnectError); !ok {
			return err
		}
	}
	return nil
}

func RemoveClient(ipPort string) error {
	serverdMu.Lock()
	s := serverd
	serverdMu.Unlock()

	if nil == s {
		return fmt.Errorf("pgrpc error: no serverd")
	}
	return s.removeClient(ipPort)
}

func PgrpcListener() (net.Listener, error) {
	serverdMu.Lock()
	s := serverd
	serverdMu.Unlock()

	if nil == s {
		return nil, fmt.Errorf("pgrpc error: no serverd")
	}

	return s, nil
}
