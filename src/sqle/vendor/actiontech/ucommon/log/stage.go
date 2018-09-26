package log

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
)

type Stage struct {
	stages   []string
	stagei   int
	threadId string
	mutex    sync.Mutex
}

const (
	MAX_STAGE_COUNT = 100
)

func NewStage() *Stage {
	s := Stage{}
	s.stages = make([]string, MAX_STAGE_COUNT, MAX_STAGE_COUNT)
	s.stagei = -1
	s.threadId = genRandomThreadId()
	return &s
}

func genRandomThreadId() string {
	seq := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	l := len(seq)
	a := rand.Intn(l * l * l)
	return fmt.Sprintf("%c%c%c", seq[a%l], seq[(a/l)%l], seq[(a/l/l)%l])
}

func (s *Stage) Enter(desc string) *Stage {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.stagei >= MAX_STAGE_COUNT {
		panic("log.Stage exceed MAX_STAGE_COUNT")
	}
	s.stagei++
	s.stages[s.stagei] = desc
	return s
}

func (s *Stage) Exit() *Stage {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.stagei < 0 {
		panic("log.Stage exceed 0")
	}
	s.stagei--
	return s
}

func (s *Stage) Go() *Stage {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	ret := NewStage()
	copy(ret.stages, s.stages)
	ret.stagei = s.stagei
	return ret
}

func (s *Stage) GoNew() *Stage {
	return NewStage()
}

func (s *Stage) ToPrefix() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.stagei < 0 {
		return ""
	}
	return " " + "|" + s.threadId + "| <" + strings.Join(s.stages[0:s.stagei+1], ".") + ">"
}
func (s *Stage) String() string {
	return s.ToPrefix()
}
