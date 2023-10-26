package server

import (
	"time"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/server/cluster"
	"github.com/sirupsen/logrus"
)

type ServerJob interface {
	Start()
	Stop()
}

var OnlyRunOnLeaderJobs = []func(entry *logrus.Entry) ServerJob{
	NewCleanJob,
	NewDingTalkJob,
	NewFeishuJob,
}

var RunOnAllJobs = []func(entry *logrus.Entry) ServerJob{
	NewWorkflowScheduleJob,
}

type ServerJobManager struct {
	clusterNode         cluster.Node
	onlyRunOnLeaderJobs []ServerJob
	runOnAllJobs        []ServerJob
	exitCh              chan struct{}
	doneCh              chan struct{}
	isLeader            bool
}

func NewServerJobManger(node cluster.Node) *ServerJobManager {
	return &ServerJobManager{
		onlyRunOnLeaderJobs: []ServerJob{},
		runOnAllJobs:        []ServerJob{},
		clusterNode:         node,
		exitCh:              make(chan struct{}),
		doneCh:              make(chan struct{}),
	}
}

func (s *ServerJobManager) Start() {
	entry := log.NewEntry().WithField("type", "server_job")
	entry.Infof("start job manager")

	defer s.startRunOnAllJob(entry)
	go func() {
		tick := time.NewTicker(5 * time.Second)
		defer tick.Stop()
		for {
			select {
			case <-tick.C:
				isLeader := s.clusterNode.IsLeader()
				if s.isLeader == isLeader {
					continue // leader not change. do nothing
				}
				s.isLeader = isLeader
				if isLeader {
					s.startOnlyRunOnLeaderJob(entry)
				} else {
					s.stopOnlyRunOnLeaderJob()
				}
			case <-s.exitCh:
				s.stopRunOnAllJob()
				s.stopOnlyRunOnLeaderJob()
				s.doneCh <- struct{}{}
				entry.Infof("stop job manager")
				return
			}
		}
	}()
}

func (s *ServerJobManager) Stop() {
	s.exitCh <- struct{}{}
	<-s.doneCh
}

func (s *ServerJobManager) startRunOnAllJob(entry *logrus.Entry) {
	for _, jobFn := range RunOnAllJobs {
		j := jobFn(entry)
		j.Start()
		s.runOnAllJobs = append(s.runOnAllJobs, j)
	}
}

func (s *ServerJobManager) stopRunOnAllJob() {
	for _, job := range s.runOnAllJobs {
		job.Stop()
	}
	s.runOnAllJobs = []ServerJob{}

}

func (s *ServerJobManager) startOnlyRunOnLeaderJob(entry *logrus.Entry) {
	for _, jobFn := range OnlyRunOnLeaderJobs {
		j := jobFn(entry)
		j.Start()
		s.onlyRunOnLeaderJobs = append(s.onlyRunOnLeaderJobs, j)
	}
}

func (s *ServerJobManager) stopOnlyRunOnLeaderJob() {
	for _, job := range s.onlyRunOnLeaderJobs {
		job.Stop()
	}
	s.onlyRunOnLeaderJobs = []ServerJob{}
}

type BaseJob struct {
	entry  *logrus.Entry
	exitCh chan struct{}
	doneCh chan struct{}

	internal time.Duration // the internal for do job.
	jobFn    func(entry *logrus.Entry)
}

func NewBaseJob(entry *logrus.Entry, internal time.Duration, jobFn func(entry *logrus.Entry)) *BaseJob {
	return &BaseJob{
		entry:    entry,
		exitCh:   make(chan struct{}),
		doneCh:   make(chan struct{}),
		internal: internal,
		jobFn:    jobFn,
	}
}

func (j *BaseJob) Start() {
	go func() {
		j.entry.Infof("start, internal is %s", j.internal)
		defer j.entry.Infof("stop")

		tick := time.NewTicker(j.internal)
		defer tick.Stop()
		for {
			select {
			case <-j.exitCh:
				j.doneCh <- struct{}{}
				return
			case <-tick.C:
				j.jobFn(j.entry)
			}
		}
	}()
}

func (j *BaseJob) Stop() {
	j.exitCh <- struct{}{}
	<-j.doneCh
}
