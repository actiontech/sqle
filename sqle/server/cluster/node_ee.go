//go:build enterprise
// +build enterprise

package cluster

import (
	"time"

	"github.com/actiontech/sqle/sqle/license"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/sirupsen/logrus"
)

func init() {
	DefaultNode = NewBaseOnModelCluster()
}

type BaseOnModelClusterNode struct {
	entry    *logrus.Entry
	ServerId string
	exitCh   chan struct{}
	doneCh   chan struct{}
}

func NewBaseOnModelCluster() Node {
	return &BaseOnModelClusterNode{
		entry:  log.NewEntry().WithField("type", "cluster"),
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
	h, err := license.CollectHardwareInfo()
	if err != nil {
		c.entry.Errorf("collect hardware info failed, error: %v", err)
	}
	err = s.RegisterClusterNode(serverId, h)
	if err != nil {
		c.entry.Errorf("register cluster node info failed, error: %v", err)
	}

	err = s.AttemptClusterLeadership(c.ServerId)
	if err != nil {
		c.entry.Error("attempt cluster leader ship failed, error: %v", err)
	}
	go func() {
		tick := time.NewTicker(time.Second * 5)
		defer tick.Stop()
		for {
			select {
			case <-tick.C:
				err := s.AttemptClusterLeadership(c.ServerId)
				if err != nil {
					c.entry.Error("attempt cluster leader ship failed, error: %v", err)
				}
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
