package log

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

type Logger interface {
	setOwner(currentOwner string)
	setFileLimit(fileLimit, totalLimit int)
	printLog(level int, line string)
	zipLoop()
	removeLoop()
	getLock() *sync.Mutex
	loadDynamicPropertiesLoop()
	getLevelAbility(l int) bool
	setLevelAbility(level int, ability bool)
}

type Prefix interface {
	ToPrefix() string
}

const calldepth = 2 // for log call stack filename and line

var (
	instance   Logger
	instanceMu sync.RWMutex
)

func UserInfo(prefix Prefix, msg string, args ...interface{}) {
	Write(prefix).UserInfo(msg, args...).done(calldepth)
}

func UserWarn(prefix Prefix, msg string, args ...interface{}) {
	Write(prefix).UserWarn(msg, args...).done(calldepth)
}

func UserError(prefix Prefix, msg string, args ...interface{}) {
	Write(prefix).UserError(msg, args...).done(calldepth)
}

func Key(prefix Prefix, msg string, args ...interface{}) {
	Write(prefix).Key(msg, args...).done(calldepth)
}

func Brief(prefix Prefix, msg string, args ...interface{}) {
	Write(prefix).Brief(msg, args...).done(calldepth)
}

func Detail(prefix Prefix, msg string, args ...interface{}) {
	Write(prefix).Detail(msg, args...).done(calldepth)
}

func UserInfoDilute1(prefix Prefix, diluteKey, msg string, args ...interface{}) {
	WriteDilute(prefix, diluteKey, 10, 60).UserInfo(msg, args...).done(calldepth)
}

func UserWarnDilute1(prefix Prefix, diluteKey, msg string, args ...interface{}) {
	WriteDilute(prefix, diluteKey, 10, 60).UserWarn(msg, args...).done(calldepth)
}

func UserErrorDilute1(prefix Prefix, diluteKey, msg string, args ...interface{}) {
	WriteDilute(prefix, diluteKey, 10, 60).UserError(msg, args...).done(calldepth)
}

func KeyDilute1(prefix Prefix, diluteKey, msg string, args ...interface{}) {
	WriteDilute(prefix, diluteKey, 10, 60).Key(msg, args...).done(calldepth)
}

func BriefDilute1(prefix Prefix, diluteKey, msg string, args ...interface{}) {
	WriteDilute(prefix, diluteKey, 10, 60).Brief(msg, args...).done(calldepth)
}

func DetailDilute1(prefix Prefix, diluteKey, msg string, args ...interface{}) {
	WriteDilute(prefix, diluteKey, 10, 60).Detail(msg, args...).done(calldepth)
}

func UserInfoDilute2(prefix Prefix, diluteKey, msg string, args ...interface{}) {
	WriteDilute(prefix, diluteKey, 60, 300).UserInfo(msg, args...).done(calldepth)
}

func UserWarnDilute2(prefix Prefix, diluteKey, msg string, args ...interface{}) {
	WriteDilute(prefix, diluteKey, 60, 300).UserWarn(msg, args...).done(calldepth)
}

func UserErrorDilute2(prefix Prefix, diluteKey, msg string, args ...interface{}) {
	WriteDilute(prefix, diluteKey, 60, 300).UserError(msg, args...).done(calldepth)
}

func KeyDilute2(prefix Prefix, diluteKey, msg string, args ...interface{}) {
	WriteDilute(prefix, diluteKey, 60, 300).Key(msg, args...).done(calldepth)
}

func BriefDilute2(prefix Prefix, diluteKey, msg string, args ...interface{}) {
	WriteDilute(prefix, diluteKey, 60, 300).Brief(msg, args...).done(calldepth)
}

func DetailDilute2(prefix Prefix, diluteKey, msg string, args ...interface{}) {
	WriteDilute(prefix, diluteKey, 60, 300).Detail(msg, args...).done(calldepth)
}

func SetOwner(currentUser string) {
	instance.setOwner(currentUser)
}

func SetLogger(logger Logger) {
	instanceMu.Lock()
	defer instanceMu.Unlock()

	instance = logger
}

func shouldPrintLog(level int) bool {
	instanceMu.RLock()
	instSnapshot := instance
	instanceMu.RUnlock()

	for i := level; i <= detail; i++ {
		if instSnapshot.getLevelAbility(i) {
			return true
		}
	}
	return false
}

func zip(fileDir, fileName string) error {
	cmd := fmt.Sprintf("tar zcf %v/%v.tar.gz %v/%v", fileDir, fileName, fileDir, fileName)
	output, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if nil == err {
		os.Remove(fmt.Sprintf("%v/%v", fileDir, fileName))
		return nil
	} else {
		os.Remove(fmt.Sprintf("%v/%v.tar.gz", fileDir, fileName))
		return fmt.Errorf("%v, output=(%v)", err, string(output))
	}
}

func openLogFile(user, path string) (*os.File, error) {
	return os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0640)
}
