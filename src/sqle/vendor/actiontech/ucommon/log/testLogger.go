package log

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type testLogger struct {
	mutex            *sync.Mutex
	dir              string
	fileLimit        int
	totalLimit       int
	fileOwner        string
	levels           [5]level
	traceFilter      string
	configPath       string
	traceFilters     map[string]bool
	traceFiltersMutex *sync.RWMutex
	logContents      [][]string
}

func NewTestLogger() *testLogger {
	ret := testLogger{}
	ret.dir = ""
	ret.fileLimit = 100
	ret.totalLimit = 1000
	ret.fileOwner = "root"
	ret.configPath = ""
	ret.setLevel(user, true, "user")
	ret.setLevel(key, true, "key")
	ret.setLevel(brief, true, "brief")
	ret.setLevel(detail, false, "detail")
	ret.traceFiltersMutex = new(sync.RWMutex)
	ret.mutex = new(sync.Mutex)
	ret.logContents = make([][]string, 5)
	ret.setTraceFilters("ha")
	return &ret
}

func (logger *testLogger) setLevel(l int, enable bool, filePath string) {
	logger.levels[l] = level{}
	if enable {
		logger.levels[l].enable = 1
	} else {
		logger.levels[l].enable = 0
	}
	logger.levels[l].filePath = logger.dir + filePath
}

func (logger *testLogger) getLock() *sync.Mutex {
	return logger.mutex
}

func (logger *testLogger) printLog(level int, line string) {
	if !logger.getLevelAbility(level) {
		return
	}
	logger.logContents[level] = append(logger.logContents[level], line)
}

func (logger *testLogger) setFileLimit(fileLimit, totalLimit int) {
	logger.fileLimit = fileLimit
	logger.totalLimit = totalLimit
}

func (logger *testLogger) setOwner(currentUser string) {
	logger.fileOwner = currentUser
}

func (logger *testLogger) zipLoop() {
	go func() {
		fmt.Println("zipping at :", time.Now())
		time.Sleep(1 * time.Minute)
	}()
}
func (logger *testLogger) removeLoop() {
	go func() {
		fmt.Println("remove:", time.Now())
		time.Sleep(1 * time.Minute)
	}()

}

func (logger *testLogger) loadDynamicPropertiesLoop() {
	go func() {
		for {
			logger.loadDynamicProperties()
			time.Sleep(30 * time.Second)
		}
	}()
}

func (logger *testLogger) setEnableBrief(ability bool) {
	if ability {
		logger.levels[brief].enable = 1
	} else {
		logger.levels[brief].enable = 0
	}
}

func (logger *testLogger) setTraceFilters(traceFilter string) {
	logger.traceFiltersMutex.Lock()
	defer logger.traceFiltersMutex.Unlock()
	traceFilters := strings.Split(traceFilter, ";")
	logger.traceFilters = make(map[string]bool)
	for _, tf := range traceFilters {
		if len(tf) == 0 {
			continue
		}
		if tf[0] == '@' {
			logger.traceFilters[tf[1:]] = true
			continue
		}
		logger.traceFilters[tf] = false
	}
	if len(traceFilter) > 0 {
		logger.levels[detail].enable = 1
	} else {
		logger.levels[detail].enable = 0
	}
}
func (logger *testLogger) loadDynamicProperties() {
	fmt.Println("loadDynamicProperties...:", time.Now())
}

func (logger *testLogger) getTraceFilters() map[string]bool {
	return logger.traceFilters
}

func (logger *testLogger) getLevel(l int) level {
	return logger.levels[l]
}

func (logger *testLogger) getLevelAbility(level int) bool {
	return 1 == logger.levels[level].enable
}

func (logger *testLogger) setLevelAbility(level int, ability bool) {
	if ability {
		logger.levels[level].enable = 1
	} else {
		logger.levels[level].enable = 0

	}

}
