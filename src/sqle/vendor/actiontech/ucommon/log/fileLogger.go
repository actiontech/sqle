package log

import (
	"actiontech/ucommon/conf"
	commonUser "actiontech/ucommon/user"
	"fmt"
	"io/ioutil"
	"os"
	goUser "os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type fileLogger struct {
	mutex       *sync.Mutex
	dir         string
	fileLimit   int
	totalLimit  int
	fileOwner   string
	levels      [6]level
	traceFilter string
	configPath  string
}

type level struct {
	filePath             string
	enable               int32
	fileFd               *os.File
	lastCheckRemovalTime time.Time
	lastSyncTime         time.Time
}

const (
	user = iota
	key
	brief
	detail
	fileTimeStamp = "2006_01_02_15_04_05"
	LogTimeStamp  = "2006-01-02 15:04:05.000"
)

func InitFileLoggerWithHouseKeep(fileLimit, totalLimit int, currentUser string, enableDetail bool) {
	instanceMu.Lock()
	instance = newLogger("./logs", "./conf/log.config")
	instanceMu.Unlock()

	instance.setLevelAbility(detail, enableDetail)
	instance.loadDynamicPropertiesLoop()
	instance.setOwner(currentUser)
	instance.setFileLimit(fileLimit, totalLimit)
	instance.zipLoop()
	instance.removeLoop()
}

func InitFileLoggerWithoutHouseKeep(currentUser string, enableDetail bool) {
	instanceMu.Lock()
	instance = newLogger("./logs", "./conf/log.config")
	instanceMu.Unlock()

	instance.setLevelAbility(detail, enableDetail)
	instance.loadDynamicPropertiesLoop()
	instance.setOwner(currentUser)
}

func newLogger(dir, configPath string) Logger {
	ret := fileLogger{}
	ret.dir = dir
	ret.fileLimit = 100
	ret.totalLimit = 1000
	ret.traceFilter = ""
	ret.fileOwner = "root"
	ret.configPath = configPath
	ret.setLevel(user, 600, true, "/user.log")
	ret.setLevel(key, 10, true, "/key.log")
	ret.setLevel(brief, 10, true, "/brief.log")
	ret.setLevel(detail, 10, false, "/detail.log")
	ret.mutex = new(sync.Mutex)
	return &ret
}

func (logger *fileLogger) setLevel(l, repeatTime int, enable bool, filePath string) {
	logger.levels[l] = level{}
	logger.setLevelAbility(l, enable)
	logger.levels[l].filePath = logger.dir + filePath
	logger.levels[l].lastSyncTime = time.Now()
}

func (logger *fileLogger) getLock() *sync.Mutex {
	return logger.mutex
}

func (logger *fileLogger) newLogFile(level int) (*os.File, error) {
	if _, err := os.Stat(logger.dir); nil != err && os.IsNotExist(err) {
		if err := os.MkdirAll(logger.dir, 0750); nil != err {
			return nil, err
		}
	}

	file, err := openLogFile(logger.fileOwner, logger.levels[level].filePath)
	if nil != err {
		return nil, err
	}

	fi, err := file.Stat()
	if nil != err {
		file.Close()
		return nil, err
	}

	if "root" != logger.fileOwner {
		ownerId := fi.Sys().(*syscall.Stat_t).Uid
		uid, gids, err := commonUser.LookupUidGidByUser(logger.fileOwner)
		if nil != err {
			if _, ok := err.(goUser.UnknownUserError); ok { //when logger.fileOwner is not initialized, should ignore log file permission check
				return nil, nil
			}
			return nil, err
		}

		if uid != int(ownerId) {
			if err := file.Chown(uid, gids[0]); nil != err {
				return nil, err
			}
		}
	}
	return file, nil
}

const (
	DILUTE_LIMITS = 100000
)

func (logger *fileLogger) printLog(level int, line string) {
	if !logger.getLevelAbility(level) {
		return
	}

	tsNow := time.Now()

	//check file existence
	if nil != logger.levels[level].fileFd &&
		time.Now().After(logger.levels[level].lastCheckRemovalTime.Add(10*time.Second)) {
		logger.levels[level].lastCheckRemovalTime = time.Now()
		if _, err := os.Stat(logger.levels[level].filePath); nil != err {
			logger.levels[level].fileFd.Close()
			logger.levels[level].fileFd = nil
		}
	}

	if nil == logger.levels[level].fileFd {
		if fd, err := logger.newLogFile(level); nil != err {
			fmt.Fprintf(os.Stderr, "[%v][LOG][ERROR] %v\n", tsNow.Format(LogTimeStamp), err)
			return
		} else if nil == fd {
			//ignore some error
			return
		} else {
			logger.levels[level].fileFd = fd
		}
	}

	file := logger.levels[level].fileFd
	if _, err := file.WriteString(line + "\n"); nil != err {
		file.Close()
		logger.levels[level].fileFd = nil
		return
	}
	now := time.Now()
	if logger.levels[level].lastSyncTime.Add(1 * time.Second).Before(now) {
		file.Sync()
		logger.levels[level].lastSyncTime = now
	}
	fi, err := file.Stat()
	if nil != err {
		file.Close()
		logger.levels[level].fileFd = nil
		return
	}
	size := fi.Size()
	if size > int64(logger.fileLimit)*1024*1024 {
		file.Close()
		logger.levels[level].fileFd = nil

		path := logger.levels[level].filePath
		fileDir := filepath.Dir(path)
		fileName := filepath.Base(path)
		t := tsNow.Format(fileTimeStamp)
		os.Rename(path, filepath.Join(fileDir, t+"_"+fileName))
	}
}

func (logger *fileLogger) setFileLimit(fileLimit, totalLimit int) {
	logger.mutex.Lock()
	logger.fileLimit = fileLimit
	logger.totalLimit = totalLimit
	logger.mutex.Unlock()
}

func (logger *fileLogger) setOwner(currentUser string) {
	logger.fileOwner = currentUser
}

func (logger *fileLogger) zipLoop() {
	go func() {
		reg := regexp.MustCompile("^\\d{4}_(\\d{2}_){5}.*\\.log$")
		for {
			files, err := ioutil.ReadDir(logger.dir)
			if nil != err {
				continue
			}
			isFound := false
			for _, file := range files {
				if !file.IsDir() && reg.Match([]byte(file.Name())) {
					var zipErr error
					isFound = true
					for i := 0; i < 3; i++ {
						zipErr = zip(logger.dir, file.Name())
						if nil == zipErr {
							break
						}
					}
					if nil != zipErr {
						fmt.Fprintf(os.Stderr, "[%v][LOG][ERROR] %v\n", time.Now().Format(LogTimeStamp), zipErr)
					}
					os.Remove(fmt.Sprintf("%v/%v", logger.dir, file.Name()))
					break
				}
			}
			if isFound {
				continue
			}

			time.Sleep(1 * time.Minute)
		}
	}()
}
func (logger *fileLogger) removeLoop() {
	go func() {
		rmExp := "^\\d{4}_(\\d{2}_){5}.*\\.log.tar.gz$"
		rmReg := regexp.MustCompile(rmExp)
		zipExp := "^\\d{4}_(\\d{2}_){5}.*\\.log$"
		zipReg := regexp.MustCompile(zipExp)
		for {
		loop:
			toRemoves := []string{}
			toZips := []string{}
			var size int64
			files, err := ioutil.ReadDir(logger.dir)
			if nil != err {
				continue
			}
			//init size, toRemoves, toZips
			for _, file := range files {
				if !file.IsDir() {
					rmMatched := rmReg.Match([]byte(file.Name()))
					zipMatched := zipReg.Match([]byte(file.Name()))
					if rmMatched || zipMatched {
						size += file.Size()
					}
					if file.Name() == "user.log" || file.Name() == "detail.log" || file.Name() == "key.log" || file.Name() == "brief.log" {
						size += file.Size()
					}
					if rmMatched {
						toRemoves = append(toRemoves, file.Name())
					}
					if zipMatched {
						toZips = append(toZips, file.Name())
					}
				}

			}
			//remove first file in toRemoves-toZip, and loop
			logger.mutex.Lock()
			totalLimit := logger.totalLimit
			logger.mutex.Unlock()
			if size >= int64(totalLimit)*1024*1024 {
				for _, toRemove := range toRemoves {
					isRemove := true
					for _, toZip := range toZips {
						if strings.HasPrefix(toRemove, toZip) { // traverse toZips, if toRemove and toZip have same prefix means this log is on zipping, should not remove this log.
							isRemove = false
							break
						}
					}
					if !isRemove {
						continue
					}
					removeFilePath := fmt.Sprintf("%v/%v", logger.dir, toRemove)
					fmt.Fprintf(os.Stderr, "[%v][LOG][REMOVE]%v\n", time.Now().Format(LogTimeStamp), removeFilePath)
					err := os.Remove(removeFilePath)
					if nil != err {
						fmt.Fprintf(os.Stderr, "[%v][LOG][ERROR] remove file err:%v\n", time.Now().Format(LogTimeStamp), err)
					}

					goto loop
				}
			}
			time.Sleep(1 * time.Minute)
		}
	}()

}

func (logger *fileLogger) loadDynamicPropertiesLoop() {
	logger.loadDynamicProperties() //sync at fileLogger init
	go func() {
		for {
			time.Sleep(30 * time.Second)
			logger.loadDynamicProperties()
		}
	}()
}

func (logger *fileLogger) loadDynamicProperties() {
	file, err := conf.ReadConfigFile(logger.configPath)
	if nil != err {
		//ignore
		return
	}

	enableDetail, err := file.GetBool("logger", "enable_detail")
	if nil != err {
		if !strings.HasPrefix(err.Error(), "option 'enable_detail' not found in section 'logger'") {
			fmt.Fprintf(os.Stderr, "[%v][LOG][ERROR] load log config \"enable_detail\" err :%v\n", time.Now().Format(LogTimeStamp), err)
		}
	} else {
		logger.setLevelAbility(detail, enableDetail)
	}

	// load filelimit and totallimit, the default filelimit is 100M, totallimit is 1G.
	filelimit, err := file.GetInt("logger", "file_limit")
	if nil != err {
		if !strings.HasPrefix(err.Error(), "option 'file_limit' not found in section 'logger'") {
			fmt.Fprintf(os.Stderr, "[%v][LOG][ERROR] load log config \"file_limit\" err :%v\n", time.Now().Format(LogTimeStamp), err)
		}
		filelimit = logger.fileLimit
	}
	totallimit, err := file.GetInt("logger", "total_limit")
	if nil != err {
		if !strings.HasPrefix(err.Error(), "option 'total_limit' not found in section 'logger'") {
			fmt.Fprintf(os.Stderr, "[%v][LOG][ERROR] load log config \"total_limit\" err :%v\n", time.Now().Format(LogTimeStamp), err)
		}
		totallimit = logger.totalLimit
	}
	logger.setFileLimit(filelimit, totallimit)
}

func getStack(skip int) string {
	depth := 0
	pcs := make([]uintptr, 100)
	runtime.Callers(skip, pcs)
	for i, pc := range pcs {
		if pc == 0 {
			depth = i - 1
			break
		}
	}

	ret := ""
	level := 0
	for i := depth; i >= 0; i-- {
		f := runtime.FuncForPC(pcs[i])
		fileName, line := f.FileLine(pcs[i])
		fileName = fileName + ""
		ret = ret + fmt.Sprintf("|-%v:%v :%v\n", fileName, f.Name(), line)
		level++
	}
	return ret
}

func (logger *fileLogger) getLevelAbility(level int) bool {
	return 1 == atomic.LoadInt32(&logger.levels[level].enable)
}

func (logger *fileLogger) setLevelAbility(level int, ability bool) {
	if ability {
		atomic.StoreInt32(&logger.levels[level].enable, 1)
	} else {
		atomic.StoreInt32(&logger.levels[level].enable, 0)
	}
}
