package log

import (
	"fmt"
	"sync"
)

type StdLogger struct {
	mutex sync.Mutex
}

func (n *StdLogger) setOwner(currentOwner string)           {}
func (n *StdLogger) setFileLimit(fileLimit, totalLimit int) {}
func (n *StdLogger) printLog(level int, line string) {
	if 1 == level /* only key */ {
		fmt.Print(line + "\n")
	}
}
func (n *StdLogger) zipLoop()                            {}
func (n *StdLogger) removeLoop()                         {}
func (n *StdLogger) getLock() *sync.Mutex                { return &n.mutex }
func (n *StdLogger) loadDynamicPropertiesLoop()          {}
func (n *StdLogger) getTraceFilters() map[string]bool    { return nil }
func (n *StdLogger) getLevelAbility(l int) bool          { return true }
func (n *StdLogger) setLevelAbility(l int, ability bool) {}
func (n *StdLogger) setTraceFilters(traceFilters string) {}
