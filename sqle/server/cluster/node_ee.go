package cluster

import (
	"time"

	"github.com/actiontech/sqle/sqle/model"
)

func init() {
	DefaultNode = NewBaseOnModelCluster()
}

type BaseOnModelClusterNode struct {
	ServerId string
	exitCh   chan struct{}
	doneCh   chan struct{}
}

func NewBaseOnModelCluster() Node {
	return &BaseOnModelClusterNode{
		exitCh: make(chan struct{}),
		doneCh: make(chan struct{}),
	}
}

func (c *BaseOnModelClusterNode) IsLeader() bool {
	s := model.GetStorage()
	id, err := s.GetClusterLeader()
	if err != nil {
		return false
	}
	return c.ServerId == id
}

func (c *BaseOnModelClusterNode) Join(serverId string) {
	c.ServerId = serverId
	s := model.GetStorage()
	s.AttemptClusterLeadership(c.ServerId)
	go func() {
		tick := time.NewTicker(time.Second * 5)
		defer tick.Stop()
		for {
			select {
			case <-tick.C:
				s.AttemptClusterLeadership(c.ServerId)
			case <-c.exitCh:
				c.doneCh <- struct{}{}
				return
			}
		}
	}()
}

func (c *BaseOnModelClusterNode) Leave() {
	c.exitCh <- struct{}{}
	<-c.doneCh
}
