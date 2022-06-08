// Copyright 2016 TiKV Project Authors.
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
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/pingcap/errors"
	"github.com/pingcap/failpoint"
	"github.com/pingcap/kvproto/pkg/metapb"
	"github.com/pingcap/kvproto/pkg/pdpb"
	"github.com/pingcap/log"
	"github.com/tikv/pd/pkg/errs"
	"github.com/tikv/pd/pkg/grpcutil"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// Region contains information of a region's meta and its peers.
type Region struct {
	Meta         *metapb.Region
	Leader       *metapb.Peer
	DownPeers    []*metapb.Peer
	PendingPeers []*metapb.Peer
}

// Client is a PD (Placement Driver) client.
// It should not be used after calling Close().
type Client interface {
	// GetClusterID gets the cluster ID from PD.
	GetClusterID(ctx context.Context) uint64
	// GetAllMembers gets the members Info from PD
	GetAllMembers(ctx context.Context) ([]*pdpb.Member, error)
	// GetLeaderAddr returns current leader's address. It returns "" before
	// syncing leader from server.
	GetLeaderAddr() string
	// GetTS gets a timestamp from PD.
	GetTS(ctx context.Context) (int64, int64, error)
	// GetTSAsync gets a timestamp from PD, without block the caller.
	GetTSAsync(ctx context.Context) TSFuture
	// GetLocalTS gets a local timestamp from PD.
	GetLocalTS(ctx context.Context, dcLocation string) (int64, int64, error)
	// GetLocalTSAsync gets a local timestamp from PD, without block the caller.
	GetLocalTSAsync(ctx context.Context, dcLocation string) TSFuture
	// GetRegion gets a region and its leader Peer from PD by key.
	// The region may expire after split. Caller is responsible for caching and
	// taking care of region change.
	// Also it may return nil if PD finds no Region for the key temporarily,
	// client should retry later.
	GetRegion(ctx context.Context, key []byte) (*Region, error)
	// GetRegionFromMember gets a region from certain members.
	GetRegionFromMember(ctx context.Context, key []byte, memberURLs []string) (*Region, error)
	// GetPrevRegion gets the previous region and its leader Peer of the region where the key is located.
	GetPrevRegion(ctx context.Context, key []byte) (*Region, error)
	// GetRegionByID gets a region and its leader Peer from PD by id.
	GetRegionByID(ctx context.Context, regionID uint64) (*Region, error)
	// ScanRegion gets a list of regions, starts from the region that contains key.
	// Limit limits the maximum number of regions returned.
	// If a region has no leader, corresponding leader will be placed by a peer
	// with empty value (PeerID is 0).
	ScanRegions(ctx context.Context, key, endKey []byte, limit int) ([]*Region, error)
	// GetStore gets a store from PD by store id.
	// The store may expire later. Caller is responsible for caching and taking care
	// of store change.
	GetStore(ctx context.Context, storeID uint64) (*metapb.Store, error)
	// GetAllStores gets all stores from pd.
	// The store may expire later. Caller is responsible for caching and taking care
	// of store change.
	GetAllStores(ctx context.Context, opts ...GetStoreOption) ([]*metapb.Store, error)
	// Update GC safe point. TiKV will check it and do GC themselves if necessary.
	// If the given safePoint is less than the current one, it will not be updated.
	// Returns the new safePoint after updating.
	UpdateGCSafePoint(ctx context.Context, safePoint uint64) (uint64, error)

	// UpdateServiceGCSafePoint updates the safepoint for specific service and
	// returns the minimum safepoint across all services, this value is used to
	// determine the safepoint for multiple services, it does not trigger a GC
	// job. Use UpdateGCSafePoint to trigger the GC job if needed.
	UpdateServiceGCSafePoint(ctx context.Context, serviceID string, ttl int64, safePoint uint64) (uint64, error)
	// ScatterRegion scatters the specified region. Should use it for a batch of regions,
	// and the distribution of these regions will be dispersed.
	// NOTICE: This method is the old version of ScatterRegions, you should use the later one as your first choice.
	ScatterRegion(ctx context.Context, regionID uint64) error
	// ScatterRegions scatters the specified regions. Should use it for a batch of regions,
	// and the distribution of these regions will be dispersed.
	ScatterRegions(ctx context.Context, regionsID []uint64, opts ...RegionsOption) (*pdpb.ScatterRegionResponse, error)
	// SplitRegions split regions by given split keys
	SplitRegions(ctx context.Context, splitKeys [][]byte, opts ...RegionsOption) (*pdpb.SplitRegionsResponse, error)
	// GetOperator gets the status of operator of the specified region.
	GetOperator(ctx context.Context, regionID uint64) (*pdpb.GetOperatorResponse, error)
	// Close closes the client.
	Close()
}

// GetStoreOp represents available options when getting stores.
type GetStoreOp struct {
	excludeTombstone bool
}

// GetStoreOption configures GetStoreOp.
type GetStoreOption func(*GetStoreOp)

// WithExcludeTombstone excludes tombstone stores from the result.
func WithExcludeTombstone() GetStoreOption {
	return func(op *GetStoreOp) { op.excludeTombstone = true }
}

// RegionsOp represents available options when operate regions
type RegionsOp struct {
	group      string
	retryLimit uint64
}

// RegionsOption configures RegionsOp
type RegionsOption func(op *RegionsOp)

// WithGroup specify the group during Scatter/Split Regions
func WithGroup(group string) RegionsOption {
	return func(op *RegionsOp) { op.group = group }
}

// WithRetry specify the retry limit during Scatter/Split Regions
func WithRetry(retry uint64) RegionsOption {
	return func(op *RegionsOp) { op.retryLimit = retry }
}

type tsoRequest struct {
	start      time.Time
	clientCtx  context.Context
	requestCtx context.Context
	done       chan error
	physical   int64
	logical    int64
	dcLocation string
}

type tsoDispatcher struct {
	dispatcherCtx    context.Context
	dispatcherCancel context.CancelFunc
	tsoRequestCh     chan *tsoRequest
}

type lastTSO struct {
	physical int64
	logical  int64
}

const (
	defaultPDTimeout      = 3 * time.Second
	dialTimeout           = 3 * time.Second
	updateMemberTimeout   = time.Second // Use a shorter timeout to recover faster from network isolation.
	tsLoopDCCheckInterval = time.Minute
	maxMergeTSORequests   = 10000 // should be higher if client is sending requests in burst
	maxInitClusterRetries = 100
	retryInterval         = 1 * time.Second
	maxRetryTimes         = 5
)

// LeaderHealthCheckInterval might be chagned in the unit to shorten the testing time.
var LeaderHealthCheckInterval = time.Second

var (
	// errFailInitClusterID is returned when failed to load clusterID from all supplied PD addresses.
	errFailInitClusterID = errors.New("[pd] failed to get cluster id")
	// errClosing is returned when request is canceled when client is closing.
	errClosing = errors.New("[pd] closing")
	// errTSOLength is returned when the number of response timestamps is inconsistent with request.
	errTSOLength = errors.New("[pd] tso length in rpc response is incorrect")
)

type client struct {
	*baseClient
	// tsoDispatcher is used to dispatch different TSO requests to
	// the corresponding dc-location TSO channel.
	tsoDispatcher sync.Map // Same as map[string]chan *tsoRequest
	// dc-location -> deadline
	tsDeadline sync.Map // Same as map[string]chan deadline
	// dc-location -> *lastTSO
	lastTSMap sync.Map // Same as map[string]*lastTSO

	checkTSDeadlineCh chan struct{}

	leaderNetworkFailure int32
}

// NewClient creates a PD client.
func NewClient(pdAddrs []string, security SecurityOption, opts ...ClientOption) (Client, error) {
	return NewClientWithContext(context.Background(), pdAddrs, security, opts...)
}

// NewClientWithContext creates a PD client with context.
func NewClientWithContext(ctx context.Context, pdAddrs []string, security SecurityOption, opts ...ClientOption) (Client, error) {
	log.Info("[pd] create pd client with endpoints", zap.Strings("pd-address", pdAddrs))
	base, err := newBaseClient(ctx, addrsToUrls(pdAddrs), security, opts...)
	if err != nil {
		return nil, err
	}
	c := &client{
		baseClient:        base,
		checkTSDeadlineCh: make(chan struct{}),
	}

	c.updateTSODispatcher()
	c.wg.Add(3)
	go c.tsLoop()
	go c.tsCancelLoop()
	go c.leaderCheckLoop()

	return c, nil
}

func (c *client) updateTSODispatcher() {
	// Set up the new TSO dispatcher
	c.allocators.Range(func(dcLocationKey, _ interface{}) bool {
		dcLocation := dcLocationKey.(string)
		if !c.checkTSODispatcher(dcLocation) {
			log.Info("[pd] create tso dispatcher", zap.String("dc-location", dcLocation))
			c.createTSODispatcher(dcLocation)
			dispatcher, _ := c.tsoDispatcher.Load(dcLocation)
			dispatcherCtx := dispatcher.(*tsoDispatcher).dispatcherCtx
			tsoRequestCh := dispatcher.(*tsoDispatcher).tsoRequestCh
			// Each goroutine is responsible for handling the tso stream request for its dc-location.
			// The only case that will make the dispatcher goroutine exit
			// is that the loopCtx is done, otherwise there is no circumstance
			// this goroutine should exit.
			go c.handleDispatcher(dispatcherCtx, dcLocation, tsoRequestCh)
		}
		return true
	})
	// Clean up the unused TSO dispatcher
	c.tsoDispatcher.Range(func(dcLocationKey, _ interface{}) bool {
		dcLocation := dcLocationKey.(string)
		// Skip the Global TSO Allocator
		if dcLocation == globalDCLocation {
			return true
		}
		if dispatcher, exist := c.allocators.Load(dcLocation); !exist {
			log.Info("[pd] delete unused tso dispatcher", zap.String("dc-location", dcLocation))
			dispatcher.(*tsoDispatcher).dispatcherCancel()
			c.tsoDispatcher.Delete(dcLocation)
		}
		return true
	})
}

func (c *client) leaderCheckLoop() {
	defer c.wg.Done()

	leaderCheckLoopCtx, leaderCheckLoopCancel := context.WithCancel(c.ctx)
	defer leaderCheckLoopCancel()

	ticker := time.NewTicker(LeaderHealthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.checkLeaderHealth(leaderCheckLoopCtx)
		}
	}
}

func (c *client) checkLeaderHealth(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	if cc, ok := c.clientConns.Load(c.GetLeaderAddr()); ok {
		healthCli := healthpb.NewHealthClient(cc.(*grpc.ClientConn))
		resp, err := healthCli.Check(ctx, &healthpb.HealthCheckRequest{Service: ""})
		rpcErr, ok := status.FromError(err)
		failpoint.Inject("unreachableNetwork1", func() {
			resp = nil
			err = status.New(codes.Unavailable, "unavailable").Err()
		})
		if (ok && isNetworkError(rpcErr.Code())) || resp.GetStatus() != healthpb.HealthCheckResponse_SERVING {
			atomic.StoreInt32(&(c.leaderNetworkFailure), int32(1))
		} else {
			atomic.StoreInt32(&(c.leaderNetworkFailure), int32(0))
		}
	}
}

type deadline struct {
	timer  <-chan time.Time
	done   chan struct{}
	cancel context.CancelFunc
}

func (c *client) tsCancelLoop() {
	defer c.wg.Done()

	tsCancelLoopCtx, tsCancelLoopCancel := context.WithCancel(c.ctx)
	defer tsCancelLoopCancel()

	ticker := time.NewTicker(tsLoopDCCheckInterval)
	defer ticker.Stop()
	for {
		// Watch every dc-location's tsDeadlineCh
		c.allocators.Range(func(dcLocation, _ interface{}) bool {
			c.watchTSDeadline(tsCancelLoopCtx, dcLocation.(string))
			return true
		})
		select {
		case <-c.checkTSDeadlineCh:
			continue
		case <-ticker.C:
			continue
		case <-tsCancelLoopCtx.Done():
			return
		}
	}
}

func (c *client) watchTSDeadline(ctx context.Context, dcLocation string) {
	if _, exist := c.tsDeadline.Load(dcLocation); !exist {
		tsDeadlineCh := make(chan deadline, 1)
		c.tsDeadline.Store(dcLocation, tsDeadlineCh)
		go func(dc string, tsDeadlineCh <-chan deadline) {
			for {
				select {
				case d := <-tsDeadlineCh:
					select {
					case <-d.timer:
						log.Error("tso request is canceled due to timeout", zap.String("dc-location", dc), errs.ZapError(errs.ErrClientGetTSOTimeout))
						d.cancel()
					case <-d.done:
					case <-ctx.Done():
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}(dcLocation, tsDeadlineCh)
	}
}

func (c *client) scheduleCheckTSDeadline() {
	select {
	case c.checkTSDeadlineCh <- struct{}{}:
	default:
	}
}

func (c *client) checkStreamTimeout(streamCtx context.Context, cancel context.CancelFunc, done chan struct{}) {
	select {
	case <-done:
		return
	case <-time.After(c.timeout):
		cancel()
	case <-streamCtx.Done():
	}
	<-done
}

func (c *client) GetAllMembers(ctx context.Context) ([]*pdpb.Member, error) {
	start := time.Now()
	defer func() { cmdDurationGetAllMembers.Observe(time.Since(start).Seconds()) }()

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	req := &pdpb.GetMembersRequest{Header: c.requestHeader()}
	ctx = grpcutil.BuildForwardContext(ctx, c.GetLeaderAddr())
	resp, err := c.getClient().GetMembers(ctx, req)
	cancel()
	if err != nil {
		cmdFailDurationGetAllMembers.Observe(time.Since(start).Seconds())
		c.ScheduleCheckLeader()
		return nil, errors.WithStack(err)
	}
	return resp.GetMembers(), nil
}

func (c *client) tsLoop() {
	defer c.wg.Done()

	loopCtx, loopCancel := context.WithCancel(c.ctx)
	defer loopCancel()

	ticker := time.NewTicker(tsLoopDCCheckInterval)
	defer ticker.Stop()
	for {
		c.updateTSODispatcher()
		select {
		case <-ticker.C:
		case <-c.checkTSODispatcherCh:
		case <-loopCtx.Done():
			return
		}
	}
}

func (c *client) createTsoStream(ctx context.Context, cancel context.CancelFunc, client pdpb.PDClient) (pdpb.PD_TsoClient, error) {
	done := make(chan struct{})
	// TODO: we need to handle a conner case that this goroutine is timeout while the stream is successfully created.
	go c.checkStreamTimeout(ctx, cancel, done)
	stream, err := client.Tso(ctx)
	done <- struct{}{}
	return stream, err
}

func (c *client) checkAllocator(dispatcherCtx context.Context, forwardCancel context.CancelFunc, dc, forwardedHostTrim, addrTrim, url string, streamCh streamCh, changedCh chan bool) {
	defer func() {
		// cancel the forward stream
		forwardCancel()
		close(streamCh)
		close(changedCh)
		requestForwarded.WithLabelValues(forwardedHostTrim, addrTrim).Set(0)
	}()
	cc, u := c.getAllocatorClientConnByDCLocation(dc)
	healthCli := healthpb.NewHealthClient(cc)
	for {
		// the pd/allocator leader change, we need to re-establish the stream
		if u != url {
			log.Info("the leader of the allocator leader is changed", zap.String("dc", dc), zap.String("origin", url), zap.String("new", u))
			return
		}
		healthCtx, healthCancel := context.WithTimeout(dispatcherCtx, c.timeout)
		resp, err := healthCli.Check(healthCtx, &healthpb.HealthCheckRequest{Service: ""})
		failpoint.Inject("unreachableNetwork", func() {
			resp.Status = healthpb.HealthCheckResponse_UNKNOWN
		})
		if err == nil && resp.GetStatus() == healthpb.HealthCheckResponse_SERVING {
			// create a stream of the original allocator
			cctx, cancel := context.WithCancel(dispatcherCtx)
			stream, err := c.createTsoStream(cctx, cancel, pdpb.NewPDClient(cc))
			if err == nil && stream != nil {
				log.Info("recover the original tso stream since the network has become normal", zap.String("dc", dc), zap.String("url", url))
				streamCh <- struct {
					cancel context.CancelFunc
					stream pdpb.PD_TsoClient
				}{
					cancel: cancel,
					stream: stream,
				}
				<-changedCh
				healthCancel()
				return
			}
		}
		healthCancel()
		select {
		case <-dispatcherCtx.Done():
			return
		case <-time.After(time.Second):
			// To ensure we can get the latest allocator leader
			// and once the leader is changed, we can exit this function.
			_, u = c.getAllocatorClientConnByDCLocation(dc)
		}
	}
}

func (c *client) checkTSODispatcher(dcLocation string) bool {
	dispatcher, ok := c.tsoDispatcher.Load(dcLocation)
	if !ok || dispatcher == nil {
		return false
	}
	return true
}

func (c *client) createTSODispatcher(dcLocation string) {
	dispatcherCtx, dispatcherCancel := context.WithCancel(c.ctx)
	dispatcher := &tsoDispatcher{
		dispatcherCtx:    dispatcherCtx,
		dispatcherCancel: dispatcherCancel,
		tsoRequestCh:     make(chan *tsoRequest, maxMergeTSORequests),
	}
	c.tsoDispatcher.Store(dcLocation, dispatcher)
}

type streamCh chan struct {
	cancel context.CancelFunc
	stream pdpb.PD_TsoClient
}

func (c *client) handleDispatcher(dispatcherCtx context.Context, dc string, tsoDispatcher chan *tsoRequest) {
	var (
		err           error
		cancel        context.CancelFunc
		stream        pdpb.PD_TsoClient
		opts          []opentracing.StartSpanOption
		requests      = make([]*tsoRequest, maxMergeTSORequests+1)
		needUpdate    = false
		streamCh      streamCh
		changedCh     chan bool
		connectionCtx connectionContext
	)
	defer func() {
		log.Info("[pd] exit tso dispatcher", zap.String("dc-location", dc))
		if cancel != nil {
			cancel()
		}
	}()
	for {
		// If the tso stream for the corresponding dc-location has not been created yet or needs to be re-created,
		// we will try to create the stream first.
		if stream == nil {
			if needUpdate {
				err = c.updateMember()
				if err != nil {
					select {
					case <-dispatcherCtx.Done():
						return
					default:
					}
					continue
				}
				needUpdate = false
			}

			connectionCtx, err = c.tryConnect(dispatcherCtx, dc)
			stream, cancel, streamCh, changedCh = connectionCtx.stream, connectionCtx.cancel, connectionCtx.streamCh, connectionCtx.changeCh

			if err != nil || stream == nil {
				select {
				case <-dispatcherCtx.Done():
					return
				default:
				}
				log.Error("[pd] create tso stream error", zap.String("dc-location", dc), errs.ZapError(errs.ErrClientCreateTSOStream, err))
				c.ScheduleCheckLeader()
				c.revokeTSORequest(errors.WithStack(err), tsoDispatcher)
				select {
				case <-time.After(time.Second):
				case <-dispatcherCtx.Done():
					return
				}
				continue
			}
		}
		select {
		case first := <-tsoDispatcher:
			pendingPlus1 := len(tsoDispatcher) + 1
			requests[0] = first
			for i := 1; i < pendingPlus1; i++ {
				requests[i] = <-tsoDispatcher
			}
			done := make(chan struct{})
			dl := deadline{
				timer:  time.After(c.timeout),
				done:   done,
				cancel: cancel,
			}
			tsDeadlineCh, ok := c.tsDeadline.Load(dc)
			for !ok || tsDeadlineCh == nil {
				c.scheduleCheckTSDeadline()
				time.Sleep(time.Millisecond * 100)
				tsDeadlineCh, ok = c.tsDeadline.Load(dc)
			}
			select {
			case tsDeadlineCh.(chan deadline) <- dl:
			case <-dispatcherCtx.Done():
				return
			}
			opts = extractSpanReference(requests[:pendingPlus1], opts[:0])
			select {
			case s, ok := <-streamCh:
				if ok {
					stream = s.stream
					changedCh <- true
					cancel = s.cancel
					log.Info("tso stream has changed back to pd or tso allocator leader")
				}
			default:
			}
			err = c.processTSORequests(stream, dc, requests[:pendingPlus1], opts)
			close(done)
		case <-dispatcherCtx.Done():
			return
		}
		// If error happens during tso stream handling, reset stream and run the next trial.
		if err != nil {
			select {
			case <-dispatcherCtx.Done():
				return
			default:
			}
			log.Error("[pd] getTS error", zap.String("dc-location", dc), errs.ZapError(errs.ErrClientGetTSO, err))
			c.ScheduleCheckLeader()
			cancel()
			stream = nil
			if IsLeaderChange(err) {
				needUpdate = true
			}
		}
	}
}

type connectionContext struct {
	stream   pdpb.PD_TsoClient
	cancel   context.CancelFunc
	streamCh streamCh
	changeCh chan bool
}

func (c *client) tryConnect(dispatcherCtx context.Context, dc string) (connectionContext, error) {
	var (
		url           string
		networkErrNum uint64
		err           error
		cc            *grpc.ClientConn
		stream        pdpb.PD_TsoClient
	)
	// retry several times before falling back to the follower when the network problem happens
	for i := 0; i < maxRetryTimes; i++ {
		cc, url = c.getAllocatorClientConnByDCLocation(dc)
		cctx, cancel := context.WithCancel(dispatcherCtx)
		stream, err = c.createTsoStream(cctx, cancel, pdpb.NewPDClient(cc))
		failpoint.Inject("unreachableNetwork", func() {
			stream = nil
			err = status.New(codes.Unavailable, "unavailable").Err()
		})
		if stream != nil && err == nil {
			return connectionContext{stream, cancel, nil, nil}, nil
		}

		if err != nil && c.enableForwarding {
			// The reason we need to judge if the error code is equal to "Canceled" here is that
			// when we create a stream we use a goroutine to manually control the timeout of the connection.
			// There is no need to wait for the transport layer timeout which can reduce the time of unavailability.
			// But it conflicts with the retry mechanism since we use the error code to decide if it is caused by network error.
			// And actually the `Canceled` error can be regarded as a kind of network error in some way.
			if rpcErr, ok := status.FromError(err); ok && (isNetworkError(rpcErr.Code()) || rpcErr.Code() == codes.Canceled) {
				networkErrNum++
			}
		}

		cancel()
		select {
		case <-dispatcherCtx.Done():
			return connectionContext{}, err
		case <-time.After(retryInterval):
		}
	}

	if networkErrNum == maxRetryTimes {
		// encounter the network error
		followerClient, addr := c.followerClient()
		if followerClient != nil {
			log.Info("fall back to use follower to forward tso stream", zap.String("dc", dc), zap.String("addr", addr))
			forwardedHost, ok := c.getAllocatorLeaderAddrByDCLocation(dc)
			if !ok {
				return connectionContext{}, errors.Errorf("cannot find the allocator leader in %s", dc)
			}

			// create the follower stream
			cctx, cancel := context.WithCancel(dispatcherCtx)
			cctx = grpcutil.BuildForwardContext(cctx, forwardedHost)
			stream, err1 := c.createTsoStream(cctx, cancel, followerClient)
			if err1 == nil {
				streamCh := make(streamCh)
				changedCh := make(chan bool)
				forwardedHostTrim := trimHTTPPrefix(forwardedHost)
				addrTrim := trimHTTPPrefix(addr)
				// the goroutine is used to check the network and change back to the original stream
				go c.checkAllocator(dispatcherCtx, cancel, dc, forwardedHostTrim, addrTrim, url, streamCh, changedCh)
				requestForwarded.WithLabelValues(forwardedHostTrim, addrTrim).Set(1)
				return connectionContext{stream, cancel, streamCh, changedCh}, nil
			}
			cancel()
			return connectionContext{}, err1
		}
	}
	return connectionContext{}, err
}

func extractSpanReference(requests []*tsoRequest, opts []opentracing.StartSpanOption) []opentracing.StartSpanOption {
	for _, req := range requests {
		if span := opentracing.SpanFromContext(req.requestCtx); span != nil {
			opts = append(opts, opentracing.ChildOf(span.Context()))
		}
	}
	return opts
}

func (c *client) processTSORequests(stream pdpb.PD_TsoClient, dcLocation string, requests []*tsoRequest, opts []opentracing.StartSpanOption) error {
	if len(opts) > 0 {
		span := opentracing.StartSpan("pdclient.processTSORequests", opts...)
		defer span.Finish()
	}
	count := int64(len(requests))
	start := time.Now()
	req := &pdpb.TsoRequest{
		Header:     c.requestHeader(),
		Count:      uint32(count),
		DcLocation: dcLocation,
	}

	if err := stream.Send(req); err != nil {
		err = errors.WithStack(err)
		c.finishTSORequest(requests, 0, 0, 0, err)
		return err
	}
	resp, err := stream.Recv()
	if err != nil {
		err = errors.WithStack(err)
		c.finishTSORequest(requests, 0, 0, 0, err)
		return err
	}
	requestDurationTSO.Observe(time.Since(start).Seconds())
	tsoBatchSize.Observe(float64(count))

	if resp.GetCount() != uint32(count) {
		err = errors.WithStack(errTSOLength)
		c.finishTSORequest(requests, 0, 0, 0, err)
		return err
	}

	physical, logical, suffixBits := resp.GetTimestamp().GetPhysical(), resp.GetTimestamp().GetLogical(), resp.GetTimestamp().GetSuffixBits()
	// logical is the highest ts here, we need to do the subtracting before we finish each TSO request.
	firstLogical := addLogical(logical, -count+1, suffixBits)
	c.compareAndSwapTS(dcLocation, physical, firstLogical, suffixBits, count)
	c.finishTSORequest(requests, physical, firstLogical, suffixBits, nil)
	return nil
}

// Because of the suffix, we need to shift the count before we add it to the logical part.
func addLogical(logical, count int64, suffixBits uint32) int64 {
	return logical + count<<suffixBits
}

func (c *client) compareAndSwapTS(dcLocation string, physical, firstLogical int64, suffixBits uint32, count int64) {
	highestLogical := addLogical(firstLogical, count-1, suffixBits)
	lastTSOInterface, loaded := c.lastTSMap.LoadOrStore(dcLocation, &lastTSO{
		physical: physical,
		// Save the highest logical part here
		logical: highestLogical,
	})
	if !loaded {
		return
	}
	lastTSOPointer := lastTSOInterface.(*lastTSO)
	lastPhysical := lastTSOPointer.physical
	lastLogical := lastTSOPointer.logical
	// The TSO we get is a range like [a, b), so we save the last TSO's b as the lastLogical and compare the new TSO's a with it.
	if tsLessEqual(physical, firstLogical, lastPhysical, lastLogical) {
		panic(errors.Errorf("%s timestamp fallback, newly acquired ts (%d, %d) is less or equal to last one (%d, %d)",
			dcLocation, physical, firstLogical, lastPhysical, lastLogical))
	}
	lastTSOPointer.physical = physical
	// Same as above, we save the highest logical part here.
	lastTSOPointer.logical = highestLogical
}

func tsLessEqual(physical, logical, thatPhysical, thatLogical int64) bool {
	if physical == thatPhysical {
		return logical <= thatLogical
	}
	return physical < thatPhysical
}

func (c *client) finishTSORequest(requests []*tsoRequest, physical, firstLogical int64, suffixBits uint32, err error) {
	for i := 0; i < len(requests); i++ {
		if span := opentracing.SpanFromContext(requests[i].requestCtx); span != nil {
			span.Finish()
		}
		requests[i].physical, requests[i].logical = physical, addLogical(firstLogical, int64(i), suffixBits)
		requests[i].done <- err
	}
}

func (c *client) revokeTSORequest(err error, tsoDispatcher chan *tsoRequest) {
	for i := 0; i < len(tsoDispatcher); i++ {
		req := <-tsoDispatcher
		req.done <- err
	}
}

func (c *client) Close() {
	c.cancel()
	c.wg.Wait()

	c.tsoDispatcher.Range(func(_, dispatcher interface{}) bool {
		if dispatcher != nil {
			c.revokeTSORequest(errors.WithStack(errClosing), dispatcher.(*tsoDispatcher).tsoRequestCh)
			dispatcher.(*tsoDispatcher).dispatcherCancel()
		}
		return true
	})

	c.clientConns.Range(func(_, cc interface{}) bool {
		if err := cc.(*grpc.ClientConn).Close(); err != nil {
			log.Error("[pd] failed to close gRPC clientConn", errs.ZapError(errs.ErrCloseGRPCConn, err))
		}
		return true
	})
}

// leaderClient gets the client of current PD leader.
func (c *client) leaderClient() pdpb.PDClient {
	if cc, ok := c.clientConns.Load(c.GetLeaderAddr()); ok {
		return pdpb.NewPDClient(cc.(*grpc.ClientConn))
	}
	return nil
}

// followerClient gets the client of current PD follower.
func (c *client) followerClient() (pdpb.PDClient, string) {
	addrs := c.GetFollowerAddr()
	var cc *grpc.ClientConn
	var err error
	if len(addrs) < 1 {
		return nil, ""
	}

	for i := 0; i < len(addrs); i++ {
		addr := addrs[rand.Intn(len(addrs))]
		if cc, err = c.getOrCreateGRPCConn(addr); err != nil {
			continue
		}
		healthCtx, healthCancel := context.WithTimeout(c.ctx, c.timeout)
		resp, err := healthpb.NewHealthClient(cc).Check(healthCtx, &healthpb.HealthCheckRequest{Service: ""})
		if err == nil && resp.GetStatus() == healthpb.HealthCheckResponse_SERVING {
			healthCancel()
			return pdpb.NewPDClient(cc), addr
		}
		healthCancel()
	}
	return nil, ""
}

func (c *client) getClient() pdpb.PDClient {
	if c.enableForwarding && atomic.LoadInt32(&c.leaderNetworkFailure) == 1 {
		followerClient, addr := c.followerClient()
		if followerClient != nil {
			log.Debug("use follower client", zap.String("addr", addr))
			return followerClient
		}
	}
	return c.leaderClient()
}

var tsoReqPool = sync.Pool{
	New: func() interface{} {
		return &tsoRequest{
			done:     make(chan error, 1),
			physical: 0,
			logical:  0,
		}
	},
}

func (c *client) GetTSAsync(ctx context.Context) TSFuture {
	return c.GetLocalTSAsync(ctx, globalDCLocation)
}

func (c *client) GetLocalTSAsync(ctx context.Context, dcLocation string) TSFuture {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan("GetLocalTSAsync", opentracing.ChildOf(span.Context()))
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	req := tsoReqPool.Get().(*tsoRequest)
	req.requestCtx = ctx
	req.clientCtx = c.ctx
	req.start = time.Now()
	req.dcLocation = dcLocation
	if err := c.dispatchRequest(dcLocation, req); err != nil {
		// Wait for a while and try again
		time.Sleep(50 * time.Millisecond)
		if err = c.dispatchRequest(dcLocation, req); err != nil {
			req.done <- err
		}
	}
	return req
}

func (c *client) dispatchRequest(dcLocation string, request *tsoRequest) error {
	dispatcher, ok := c.tsoDispatcher.Load(dcLocation)
	if !ok {
		err := errs.ErrClientGetTSO.FastGenByArgs(fmt.Sprintf("unknown dc-location %s to the client", dcLocation))
		log.Error("[pd] dispatch tso request error", zap.String("dc-location", dcLocation), errs.ZapError(err))
		c.ScheduleCheckLeader()
		return err
	}
	dispatcher.(*tsoDispatcher).tsoRequestCh <- request
	return nil
}

// TSFuture is a future which promises to return a TSO.
type TSFuture interface {
	// Wait gets the physical and logical time, it would block caller if data is not available yet.
	Wait() (int64, int64, error)
}

func (req *tsoRequest) Wait() (physical int64, logical int64, err error) {
	// If tso command duration is observed very high, the reason could be it
	// takes too long for Wait() be called.
	start := time.Now()
	cmdDurationTSOAsyncWait.Observe(start.Sub(req.start).Seconds())
	select {
	case err = <-req.done:
		err = errors.WithStack(err)
		defer tsoReqPool.Put(req)
		if err != nil {
			cmdFailDurationTSO.Observe(time.Since(req.start).Seconds())
			return 0, 0, err
		}
		physical, logical = req.physical, req.logical
		now := time.Now()
		cmdDurationWait.Observe(now.Sub(start).Seconds())
		cmdDurationTSO.Observe(now.Sub(req.start).Seconds())
		return
	case <-req.requestCtx.Done():
		return 0, 0, errors.WithStack(req.requestCtx.Err())
	case <-req.clientCtx.Done():
		return 0, 0, errors.WithStack(req.clientCtx.Err())
	}
}

func (c *client) GetTS(ctx context.Context) (physical int64, logical int64, err error) {
	resp := c.GetTSAsync(ctx)
	return resp.Wait()
}

func (c *client) GetLocalTS(ctx context.Context, dcLocation string) (physical int64, logical int64, err error) {
	resp := c.GetLocalTSAsync(ctx, dcLocation)
	return resp.Wait()
}

func handleRegionResponse(res *pdpb.GetRegionResponse) *Region {
	if res.Region == nil {
		return nil
	}

	r := &Region{
		Meta:         res.Region,
		Leader:       res.Leader,
		PendingPeers: res.PendingPeers,
	}
	for _, s := range res.DownPeers {
		r.DownPeers = append(r.DownPeers, s.Peer)
	}
	return r
}

func (c *client) GetRegion(ctx context.Context, key []byte) (*Region, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan("pdclient.GetRegion", opentracing.ChildOf(span.Context()))
		defer span.Finish()
	}
	start := time.Now()
	defer func() { cmdDurationGetRegion.Observe(time.Since(start).Seconds()) }()

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	req := &pdpb.GetRegionRequest{
		Header:    c.requestHeader(),
		RegionKey: key,
	}
	ctx = grpcutil.BuildForwardContext(ctx, c.GetLeaderAddr())
	resp, err := c.getClient().GetRegion(ctx, req)
	cancel()

	if err != nil {
		cmdFailDurationGetRegion.Observe(time.Since(start).Seconds())
		c.ScheduleCheckLeader()
		return nil, errors.WithStack(err)
	}
	return handleRegionResponse(resp), nil
}

func isNetworkError(code codes.Code) bool {
	return code == codes.Unavailable || code == codes.DeadlineExceeded
}

func (c *client) GetRegionFromMember(ctx context.Context, key []byte, memberURLs []string) (*Region, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan("pdclient.GetRegionFromMember", opentracing.ChildOf(span.Context()))
		defer span.Finish()
	}
	start := time.Now()
	defer func() { cmdDurationGetRegion.Observe(time.Since(start).Seconds()) }()

	var resp *pdpb.GetRegionResponse
	for _, url := range memberURLs {
		conn, err := c.getOrCreateGRPCConn(url)
		if err != nil {
			log.Error("[pd] can't get grpc connection", zap.String("member-URL", url), errs.ZapError(err))
			continue
		}
		cc := pdpb.NewPDClient(conn)
		resp, err = cc.GetRegion(ctx, &pdpb.GetRegionRequest{
			Header:    c.requestHeader(),
			RegionKey: key,
		})
		if err != nil {
			log.Error("[pd] can't get region info", zap.String("member-URL", url), errs.ZapError(err))
			continue
		}
		if resp != nil {
			break
		}
	}

	if resp == nil {
		cmdFailDurationGetRegion.Observe(time.Since(start).Seconds())
		c.ScheduleCheckLeader()
		errorMsg := fmt.Sprintf("[pd] can't get region info from member URLs: %+v", memberURLs)
		return nil, errors.WithStack(errors.New(errorMsg))
	}
	return handleRegionResponse(resp), nil
}

func (c *client) GetPrevRegion(ctx context.Context, key []byte) (*Region, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan("pdclient.GetPrevRegion", opentracing.ChildOf(span.Context()))
		defer span.Finish()
	}
	start := time.Now()
	defer func() { cmdDurationGetPrevRegion.Observe(time.Since(start).Seconds()) }()

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	req := &pdpb.GetRegionRequest{
		Header:    c.requestHeader(),
		RegionKey: key,
	}
	ctx = grpcutil.BuildForwardContext(ctx, c.GetLeaderAddr())
	resp, err := c.getClient().GetPrevRegion(ctx, req)
	cancel()

	if err != nil {
		cmdFailDurationGetPrevRegion.Observe(time.Since(start).Seconds())
		c.ScheduleCheckLeader()
		return nil, errors.WithStack(err)
	}
	return handleRegionResponse(resp), nil
}

func (c *client) GetRegionByID(ctx context.Context, regionID uint64) (*Region, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan("pdclient.GetRegionByID", opentracing.ChildOf(span.Context()))
		defer span.Finish()
	}
	start := time.Now()
	defer func() { cmdDurationGetRegionByID.Observe(time.Since(start).Seconds()) }()

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	req := &pdpb.GetRegionByIDRequest{
		Header:   c.requestHeader(),
		RegionId: regionID,
	}
	ctx = grpcutil.BuildForwardContext(ctx, c.GetLeaderAddr())
	resp, err := c.getClient().GetRegionByID(ctx, req)
	cancel()

	if err != nil {
		cmdFailedDurationGetRegionByID.Observe(time.Since(start).Seconds())
		c.ScheduleCheckLeader()
		return nil, errors.WithStack(err)
	}
	return handleRegionResponse(resp), nil
}

func (c *client) ScanRegions(ctx context.Context, key, endKey []byte, limit int) ([]*Region, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan("pdclient.ScanRegions", opentracing.ChildOf(span.Context()))
		defer span.Finish()
	}
	start := time.Now()
	defer cmdDurationScanRegions.Observe(time.Since(start).Seconds())

	var cancel context.CancelFunc
	scanCtx := ctx
	if _, ok := ctx.Deadline(); !ok {
		scanCtx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}
	req := &pdpb.ScanRegionsRequest{
		Header:   c.requestHeader(),
		StartKey: key,
		EndKey:   endKey,
		Limit:    int32(limit),
	}
	scanCtx = grpcutil.BuildForwardContext(scanCtx, c.GetLeaderAddr())
	resp, err := c.getClient().ScanRegions(scanCtx, req)

	if err != nil {
		cmdFailedDurationScanRegions.Observe(time.Since(start).Seconds())
		c.ScheduleCheckLeader()
		return nil, errors.WithStack(err)
	}

	return handleRegionsResponse(resp), nil
}

func handleRegionsResponse(resp *pdpb.ScanRegionsResponse) []*Region {
	var regions []*Region
	if len(resp.GetRegions()) == 0 {
		// Make it compatible with old server.
		metas, leaders := resp.GetRegionMetas(), resp.GetLeaders()
		for i := range metas {
			r := &Region{Meta: metas[i]}
			if i < len(leaders) {
				r.Leader = leaders[i]
			}
			regions = append(regions, r)
		}
	} else {
		for _, r := range resp.GetRegions() {
			region := &Region{
				Meta:         r.Region,
				Leader:       r.Leader,
				PendingPeers: r.PendingPeers,
			}
			for _, p := range r.DownPeers {
				region.DownPeers = append(region.DownPeers, p.Peer)
			}
			regions = append(regions, region)
		}
	}
	return regions
}

func (c *client) GetStore(ctx context.Context, storeID uint64) (*metapb.Store, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan("pdclient.GetStore", opentracing.ChildOf(span.Context()))
		defer span.Finish()
	}
	start := time.Now()
	defer func() { cmdDurationGetStore.Observe(time.Since(start).Seconds()) }()

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	req := &pdpb.GetStoreRequest{
		Header:  c.requestHeader(),
		StoreId: storeID,
	}
	ctx = grpcutil.BuildForwardContext(ctx, c.GetLeaderAddr())
	resp, err := c.getClient().GetStore(ctx, req)
	cancel()

	if err != nil {
		cmdFailedDurationGetStore.Observe(time.Since(start).Seconds())
		c.ScheduleCheckLeader()
		return nil, errors.WithStack(err)
	}
	return handleStoreResponse(resp)
}

func handleStoreResponse(resp *pdpb.GetStoreResponse) (*metapb.Store, error) {
	store := resp.GetStore()
	if store == nil {
		return nil, errors.New("[pd] store field in rpc response not set")
	}
	if store.GetState() == metapb.StoreState_Tombstone {
		return nil, nil
	}
	return store, nil
}

func (c *client) GetAllStores(ctx context.Context, opts ...GetStoreOption) ([]*metapb.Store, error) {
	// Applies options
	options := &GetStoreOp{}
	for _, opt := range opts {
		opt(options)
	}

	if span := opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan("pdclient.GetAllStores", opentracing.ChildOf(span.Context()))
		defer span.Finish()
	}
	start := time.Now()
	defer func() { cmdDurationGetAllStores.Observe(time.Since(start).Seconds()) }()

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	req := &pdpb.GetAllStoresRequest{
		Header:                 c.requestHeader(),
		ExcludeTombstoneStores: options.excludeTombstone,
	}
	ctx = grpcutil.BuildForwardContext(ctx, c.GetLeaderAddr())
	resp, err := c.getClient().GetAllStores(ctx, req)
	cancel()

	if err != nil {
		cmdFailedDurationGetAllStores.Observe(time.Since(start).Seconds())
		c.ScheduleCheckLeader()
		return nil, errors.WithStack(err)
	}
	return resp.GetStores(), nil
}

func (c *client) UpdateGCSafePoint(ctx context.Context, safePoint uint64) (uint64, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan("pdclient.UpdateGCSafePoint", opentracing.ChildOf(span.Context()))
		defer span.Finish()
	}
	start := time.Now()
	defer func() { cmdDurationUpdateGCSafePoint.Observe(time.Since(start).Seconds()) }()

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	req := &pdpb.UpdateGCSafePointRequest{
		Header:    c.requestHeader(),
		SafePoint: safePoint,
	}
	ctx = grpcutil.BuildForwardContext(ctx, c.GetLeaderAddr())
	resp, err := c.getClient().UpdateGCSafePoint(ctx, req)
	cancel()

	if err != nil {
		cmdFailedDurationUpdateGCSafePoint.Observe(time.Since(start).Seconds())
		c.ScheduleCheckLeader()
		return 0, errors.WithStack(err)
	}
	return resp.GetNewSafePoint(), nil
}

// UpdateServiceGCSafePoint updates the safepoint for specific service and
// returns the minimum safepoint across all services, this value is used to
// determine the safepoint for multiple services, it does not trigger a GC
// job. Use UpdateGCSafePoint to trigger the GC job if needed.
func (c *client) UpdateServiceGCSafePoint(ctx context.Context, serviceID string, ttl int64, safePoint uint64) (uint64, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan("pdclient.UpdateServiceGCSafePoint", opentracing.ChildOf(span.Context()))
		defer span.Finish()
	}

	start := time.Now()
	defer func() { cmdDurationUpdateServiceGCSafePoint.Observe(time.Since(start).Seconds()) }()

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	req := &pdpb.UpdateServiceGCSafePointRequest{
		Header:    c.requestHeader(),
		ServiceId: []byte(serviceID),
		TTL:       ttl,
		SafePoint: safePoint,
	}
	ctx = grpcutil.BuildForwardContext(ctx, c.GetLeaderAddr())
	resp, err := c.getClient().UpdateServiceGCSafePoint(ctx, req)
	cancel()

	if err != nil {
		cmdFailedDurationUpdateServiceGCSafePoint.Observe(time.Since(start).Seconds())
		c.ScheduleCheckLeader()
		return 0, errors.WithStack(err)
	}
	return resp.GetMinSafePoint(), nil
}

func (c *client) ScatterRegion(ctx context.Context, regionID uint64) error {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan("pdclient.ScatterRegion", opentracing.ChildOf(span.Context()))
		defer span.Finish()
	}
	return c.scatterRegionsWithGroup(ctx, regionID, "")
}

func (c *client) scatterRegionsWithGroup(ctx context.Context, regionID uint64, group string) error {
	start := time.Now()
	defer func() { cmdDurationScatterRegion.Observe(time.Since(start).Seconds()) }()

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	req := &pdpb.ScatterRegionRequest{
		Header:   c.requestHeader(),
		RegionId: regionID,
		Group:    group,
	}
	ctx = grpcutil.BuildForwardContext(ctx, c.GetLeaderAddr())
	resp, err := c.getClient().ScatterRegion(ctx, req)
	cancel()
	if err != nil {
		return err
	}
	if resp.Header.GetError() != nil {
		return errors.Errorf("scatter region %d failed: %s", regionID, resp.Header.GetError().String())
	}
	return nil
}

func (c *client) ScatterRegions(ctx context.Context, regionsID []uint64, opts ...RegionsOption) (*pdpb.ScatterRegionResponse, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan("pdclient.ScatterRegions", opentracing.ChildOf(span.Context()))
		defer span.Finish()
	}
	return c.scatterRegionsWithOptions(ctx, regionsID, opts...)
}

func (c *client) GetOperator(ctx context.Context, regionID uint64) (*pdpb.GetOperatorResponse, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan("pdclient.GetOperator", opentracing.ChildOf(span.Context()))
		defer span.Finish()
	}
	start := time.Now()
	defer func() { cmdDurationGetOperator.Observe(time.Since(start).Seconds()) }()

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	req := &pdpb.GetOperatorRequest{
		Header:   c.requestHeader(),
		RegionId: regionID,
	}
	ctx = grpcutil.BuildForwardContext(ctx, c.GetLeaderAddr())
	return c.getClient().GetOperator(ctx, req)
}

// SplitRegions split regions by given split keys
func (c *client) SplitRegions(ctx context.Context, splitKeys [][]byte, opts ...RegionsOption) (*pdpb.SplitRegionsResponse, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan("pdclient.SplitRegions", opentracing.ChildOf(span.Context()))
		defer span.Finish()
	}
	start := time.Now()
	defer func() { cmdDurationSplitRegions.Observe(time.Since(start).Seconds()) }()
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	options := &RegionsOp{}
	for _, opt := range opts {
		opt(options)
	}
	req := &pdpb.SplitRegionsRequest{
		Header:     c.requestHeader(),
		SplitKeys:  splitKeys,
		RetryLimit: options.retryLimit,
	}
	ctx = grpcutil.BuildForwardContext(ctx, c.GetLeaderAddr())
	return c.getClient().SplitRegions(ctx, req)
}

func (c *client) requestHeader() *pdpb.RequestHeader {
	return &pdpb.RequestHeader{
		ClusterId: c.clusterID,
	}
}

func (c *client) scatterRegionsWithOptions(ctx context.Context, regionsID []uint64, opts ...RegionsOption) (*pdpb.ScatterRegionResponse, error) {
	start := time.Now()
	defer func() { cmdDurationScatterRegions.Observe(time.Since(start).Seconds()) }()
	options := &RegionsOp{}
	for _, opt := range opts {
		opt(options)
	}
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	req := &pdpb.ScatterRegionRequest{
		Header:     c.requestHeader(),
		Group:      options.group,
		RegionsId:  regionsID,
		RetryLimit: options.retryLimit,
	}

	ctx = grpcutil.BuildForwardContext(ctx, c.GetLeaderAddr())
	resp, err := c.getClient().ScatterRegion(ctx, req)
	cancel()

	if err != nil {
		return nil, err
	}
	if resp.Header.GetError() != nil {
		return nil, errors.Errorf("scatter regions %v failed: %s", regionsID, resp.Header.GetError().String())
	}
	return resp, nil
}

func addrsToUrls(addrs []string) []string {
	// Add default schema "http://" to addrs.
	urls := make([]string, 0, len(addrs))
	for _, addr := range addrs {
		if strings.Contains(addr, "://") {
			urls = append(urls, addr)
		} else {
			urls = append(urls, "http://"+addr)
		}
	}
	return urls
}

// IsLeaderChange will determine whether there is a leader change.
func IsLeaderChange(err error) bool {
	errMsg := err.Error()
	return strings.Contains(errMsg, errs.NotLeaderErr) || strings.Contains(errMsg, errs.MismatchLeaderErr)
}

func trimHTTPPrefix(str string) string {
	str = strings.TrimPrefix(str, "http://")
	str = strings.TrimPrefix(str, "https://")
	return str
}
