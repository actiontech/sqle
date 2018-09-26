package log

import (
	"sync"
)

type NoLogger struct {
	mutex sync.Mutex
}

func init() {
	instance = &NoLogger{}
}

func (n *NoLogger) setOwner(currentOwner string)                                                                              {}
func (n *NoLogger) setFileLimit(fileLimit, totalLimit int)                                                                    {}
func (n *NoLogger) printLog(level int, line string) {}
func (n *NoLogger) zipLoop()                                                                                                  {}
func (n *NoLogger) removeLoop()                                                                                               {}
func (n *NoLogger) getLock() *sync.Mutex                                                                                      { return &n.mutex }
func (n *NoLogger) loadDynamicPropertiesLoop()                                                                                {}
func (n *NoLogger) getLevelAbility(l int) bool                                                                                { return false }
func (n *NoLogger) setLevelAbility(l int, ability bool)                                                                       {}
