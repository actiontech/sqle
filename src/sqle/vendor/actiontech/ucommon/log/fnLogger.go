package log

import (
	"sync"
)

type FnLogger struct {
	mutex sync.Mutex
	LogFn func(string)
}

func (n *FnLogger) setOwner(currentOwner string)           {}
func (n *FnLogger) setFileLimit(fileLimit, totalLimit int) {}
func (n *FnLogger) printLog(level int, line string) {
	n.LogFn(line)
}
func (n *FnLogger) zipLoop()                            {}
func (n *FnLogger) removeLoop()                         {}
func (n *FnLogger) getLock() *sync.Mutex                { return &n.mutex }
func (n *FnLogger) loadDynamicPropertiesLoop()          {}
func (n *FnLogger) getTraceFilters() map[string]bool    { return nil }
func (n *FnLogger) getLevelAbility(l int) bool          { return true }
func (n *FnLogger) setLevelAbility(l int, ability bool) {}
func (n *FnLogger) setTraceFilters(traceFilters string) {}
