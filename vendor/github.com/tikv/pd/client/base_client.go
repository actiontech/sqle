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

package pd

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pingcap/errors"
	"github.com/pingcap/failpoint"
	"github.com/pingcap/kvproto/pkg/pdpb"
	"github.com/pingcap/log"
	"github.com/tikv/pd/pkg/errs"
	"github.com/tikv/pd/pkg/grpcutil"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// baseClient is a basic client for all other complex client.
type baseClient struct {
	urls      []string
	clusterID uint64
	// PD leader URL
	leader atomic.Value // Store as string
	// PD follower URLs
	followers atomic.Value // Store as []string
	// dc-location -> TSO allocator leader gRPC connection
	clientConns sync.Map // Store as map[string]*grpc.ClientConn
	// dc-location -> TSO allocator leader URL
	allocators sync.Map // Store as map[string]string

	checkLeaderCh        chan struct{}
	checkTSODispatcherCh chan struct{}

	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc

	security SecurityOption

	gRPCDialOptions  []grpc.DialOption
	timeout          time.Duration
	maxRetryTimes    int
	enableForwarding bool
}

// SecurityOption records options about tls
type SecurityOption struct {
	CAPath   string
	CertPath string
	KeyPath  string
}

// ClientOption configures client.
type ClientOption func(c *baseClient)

// WithGRPCDialOptions configures the client with gRPC dial options.
func WithGRPCDialOptions(opts ...grpc.DialOption) ClientOption {
	return func(c *baseClient) {
		c.gRPCDialOptions = append(c.gRPCDialOptions, opts...)
	}
}

// WithCustomTimeoutOption configures the client with timeout option.
func WithCustomTimeoutOption(timeout time.Duration) ClientOption {
	return func(c *baseClient) {
		c.timeout = timeout
	}
}

// WithForwardingOption configures the client with forwarding option.
func WithForwardingOption(enableForwarding bool) ClientOption {
	return func(c *baseClient) {
		c.enableForwarding = enableForwarding
	}
}

// WithMaxErrorRetry configures the client max retry times when connect meets error.
func WithMaxErrorRetry(count int) ClientOption {
	return func(c *baseClient) {
		c.maxRetryTimes = count
	}
}

// newBaseClient returns a new baseClient.
func newBaseClient(ctx context.Context, urls []string, security SecurityOption, opts ...ClientOption) (*baseClient, error) {
	ctx1, cancel := context.WithCancel(ctx)
	c := &baseClient{
		urls:                 urls,
		checkLeaderCh:        make(chan struct{}, 1),
		checkTSODispatcherCh: make(chan struct{}, 1),
		ctx:                  ctx1,
		cancel:               cancel,
		security:             security,
		timeout:              defaultPDTimeout,
		maxRetryTimes:        maxInitClusterRetries,
	}
	for _, opt := range opts {
		opt(c)
	}

	if err := c.initRetry(c.initClusterID); err != nil {
		c.cancel()
		return nil, err
	}
	if err := c.initRetry(c.updateMember); err != nil {
		c.cancel()
		return nil, err
	}
	log.Info("[pd] init cluster id", zap.Uint64("cluster-id", c.clusterID))

	c.wg.Add(1)
	go c.memberLoop()

	return c, nil
}

func (c *baseClient) initRetry(f func() error) error {
	var err error
	for i := 0; i < c.maxRetryTimes; i++ {
		if err = f(); err == nil {
			return nil
		}
		select {
		case <-c.ctx.Done():
			return err
		case <-time.After(time.Second):
		}
	}
	return errors.WithStack(err)
}

func (c *baseClient) memberLoop() {
	defer c.wg.Done()

	ctx, cancel := context.WithCancel(c.ctx)
	defer cancel()

	for {
		select {
		case <-c.checkLeaderCh:
		case <-time.After(time.Minute):
		case <-ctx.Done():
			return
		}
		failpoint.Inject("skipUpdateMember", func() {
			failpoint.Continue()
		})
		if err := c.updateMember(); err != nil {
			log.Error("[pd] failed updateMember", errs.ZapError(err))
		}
	}
}

// ScheduleCheckLeader is used to check leader.
func (c *baseClient) ScheduleCheckLeader() {
	select {
	case c.checkLeaderCh <- struct{}{}:
	default:
	}
}

func (c *baseClient) scheduleCheckTSODispatcher() {
	select {
	case c.checkTSODispatcherCh <- struct{}{}:
	default:
	}
}

// GetClusterID returns the ClusterID.
func (c *baseClient) GetClusterID(context.Context) uint64 {
	return c.clusterID
}

// GetLeaderAddr returns the leader address.
func (c *baseClient) GetLeaderAddr() string {
	leaderAddr := c.leader.Load()
	if leaderAddr == nil {
		return ""
	}
	return leaderAddr.(string)
}

// GetLeaderAddr returns the follower address.
func (c *baseClient) GetFollowerAddr() []string {
	followerAddrs := c.followers.Load()
	if followerAddrs == nil {
		return []string{}
	}
	return followerAddrs.([]string)
}

// GetURLs returns the URLs.
// For testing use. It should only be called when the client is closed.
func (c *baseClient) GetURLs() []string {
	return c.urls
}

func (c *baseClient) GetAllocatorLeaderURLs() map[string]string {
	allocatorLeader := make(map[string]string)
	c.allocators.Range(func(dcLocation, url interface{}) bool {
		allocatorLeader[dcLocation.(string)] = url.(string)
		return true
	})
	return allocatorLeader
}

func (c *baseClient) getAllocatorLeaderAddrByDCLocation(dcLocation string) (string, bool) {
	url, exist := c.allocators.Load(dcLocation)
	if !exist {
		return "", false
	}
	return url.(string), true
}

func (c *baseClient) getAllocatorClientConnByDCLocation(dcLocation string) (*grpc.ClientConn, string) {
	url, ok := c.allocators.Load(dcLocation)
	if !ok {
		panic(fmt.Sprintf("the allocator leader in %s should exist", dcLocation))
	}
	cc, ok := c.clientConns.Load(url)
	if !ok {
		panic(fmt.Sprintf("the client connection of %s in %s should exist", url, dcLocation))
	}
	return cc.(*grpc.ClientConn), url.(string)
}

const globalDCLocation = "global"

func (c *baseClient) gcAllocatorLeaderAddr(curAllocatorMap map[string]*pdpb.Member) {
	// Clean up the old TSO allocators
	c.allocators.Range(func(dcLocationKey, _ interface{}) bool {
		dcLocation := dcLocationKey.(string)
		// Skip the Global TSO Allocator
		if dcLocation == globalDCLocation {
			return true
		}
		if _, exist := curAllocatorMap[dcLocation]; !exist {
			log.Info("[pd] delete unused tso allocator", zap.String("dc-location", dcLocation))
			c.allocators.Delete(dcLocation)
		}
		return true
	})
}

func (c *baseClient) initClusterID() error {
	ctx, cancel := context.WithCancel(c.ctx)
	defer cancel()
	for _, u := range c.urls {
		members, err := c.getMembers(ctx, u, c.timeout)
		if err != nil || members.GetHeader() == nil {
			log.Warn("[pd] failed to get cluster id", zap.String("url", u), errs.ZapError(err))
			continue
		}
		c.clusterID = members.GetHeader().GetClusterId()
		return nil
	}
	return errors.WithStack(errFailInitClusterID)
}

func (c *baseClient) updateMember() error {
	for _, u := range c.urls {
		members, err := c.getMembers(c.ctx, u, updateMemberTimeout)

		var errTSO error
		if err == nil {
			if members.GetLeader() == nil || len(members.GetLeader().GetClientUrls()) == 0 {
				err = errs.ErrClientGetLeader.FastGenByArgs("leader address don't exist")
			}
			// Still need to update TsoAllocatorLeaders, even if there is no PD leader
			errTSO = c.switchTSOAllocatorLeader(members.GetTsoAllocatorLeaders())
		}

		// Failed to get PD leader
		if err != nil {
			log.Info("[pd] cannot update member from this address",
				zap.String("address", u),
				errs.ZapError(err))
			select {
			case <-c.ctx.Done():
				return errors.WithStack(err)
			default:
				continue
			}
		}

		c.updateURLs(members.GetMembers())
		c.updateFollowers(members.GetMembers(), members.GetLeader())
		if err := c.switchLeader(members.GetLeader().GetClientUrls()); err != nil {
			return err
		}
		c.scheduleCheckTSODispatcher()

		// If `switchLeader` succeeds but `switchTSOAllocatorLeader` has an error,
		// the error of `switchTSOAllocatorLeader` will be returned.
		return errTSO
	}
	return errs.ErrClientGetLeader.FastGenByArgs(c.urls)
}

func (c *baseClient) getMembers(ctx context.Context, url string, timeout time.Duration) (*pdpb.GetMembersResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	cc, err := c.getOrCreateGRPCConn(url)
	if err != nil {
		return nil, err
	}
	members, err := pdpb.NewPDClient(cc).GetMembers(ctx, &pdpb.GetMembersRequest{})
	if err != nil {
		attachErr := errors.Errorf("error:%s target:%s status:%s", err, cc.Target(), cc.GetState().String())
		return nil, errs.ErrClientGetMember.Wrap(attachErr).GenWithStackByCause()
	}
	return members, nil
}

func (c *baseClient) updateURLs(members []*pdpb.Member) {
	urls := make([]string, 0, len(members))
	for _, m := range members {
		urls = append(urls, m.GetClientUrls()...)
	}

	sort.Strings(urls)
	// the url list is same.
	if reflect.DeepEqual(c.urls, urls) {
		return
	}

	log.Info("[pd] update member urls", zap.Strings("old-urls", c.urls), zap.Strings("new-urls", urls))
	c.urls = urls
}

func (c *baseClient) switchLeader(addrs []string) error {
	// FIXME: How to safely compare leader urls? For now, only allows one client url.
	addr := addrs[0]
	oldLeader := c.GetLeaderAddr()
	if addr == oldLeader {
		return nil
	}

	if _, err := c.getOrCreateGRPCConn(addr); err != nil {
		log.Warn("[pd] failed to connect leader", zap.String("leader", addr), errs.ZapError(err))
		return err
	}
	// Set PD leader and Global TSO Allocator (which is also the PD leader)
	c.leader.Store(addr)
	c.allocators.Store(globalDCLocation, addr)
	log.Info("[pd] switch leader", zap.String("new-leader", addr), zap.String("old-leader", oldLeader))
	return nil
}

func (c *baseClient) updateFollowers(members []*pdpb.Member, leader *pdpb.Member) {
	var addrs []string
	for _, member := range members {
		if member.GetMemberId() != leader.GetMemberId() {
			if len(member.GetClientUrls()) > 0 {
				addrs = append(addrs, member.GetClientUrls()...)
			}
		}
	}
	c.followers.Store(addrs)
}

func (c *baseClient) switchTSOAllocatorLeader(allocatorMap map[string]*pdpb.Member) error {
	if len(allocatorMap) == 0 {
		return nil
	}
	// Switch to the new one
	for dcLocation, member := range allocatorMap {
		if len(member.GetClientUrls()) == 0 {
			continue
		}
		addr := member.GetClientUrls()[0]
		oldAddr, exist := c.getAllocatorLeaderAddrByDCLocation(dcLocation)
		if exist && addr == oldAddr {
			continue
		}
		if _, err := c.getOrCreateGRPCConn(addr); err != nil {
			log.Warn("[pd] failed to connect dc tso allocator leader",
				zap.String("dc-location", dcLocation),
				zap.String("leader", addr),
				errs.ZapError(err))
			return err
		}
		c.allocators.Store(dcLocation, addr)
		log.Info("[pd] switch dc tso allocator leader",
			zap.String("dc-location", dcLocation),
			zap.String("new-leader", addr),
			zap.String("old-leader", oldAddr))
	}
	// Garbage collection of the old TSO allocator leaders
	c.gcAllocatorLeaderAddr(allocatorMap)
	return nil
}

func (c *baseClient) getOrCreateGRPCConn(addr string) (*grpc.ClientConn, error) {
	conn, ok := c.clientConns.Load(addr)
	if ok {
		return conn.(*grpc.ClientConn), nil
	}
	tlsCfg, err := grpcutil.TLSConfig{
		CAPath:   c.security.CAPath,
		CertPath: c.security.CertPath,
		KeyPath:  c.security.KeyPath,
	}.ToTLSConfig()
	if err != nil {
		return nil, err
	}
	dCtx, cancel := context.WithTimeout(c.ctx, dialTimeout)
	defer cancel()
	cc, err := grpcutil.GetClientConn(dCtx, addr, tlsCfg, c.gRPCDialOptions...)
	if err != nil {
		return nil, err
	}
	if old, ok := c.clientConns.Load(addr); ok {
		cc.Close()
		log.Debug("use old connection", zap.String("target", cc.Target()), zap.String("state", cc.GetState().String()))
		return old.(*grpc.ClientConn), nil
	}
	c.clientConns.Store(addr, cc)
	return cc, nil
}
