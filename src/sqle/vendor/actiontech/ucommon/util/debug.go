package util

import (
	"actiontech/ucommon/log"
	"fmt"
	"io/ioutil"
	"time"
)

var autoQaDebugCache string //performance

func HasAutoQaDebugCondf(cond string, args ...interface{}) (ok bool, matches []string) {
	if "" == autoQaDebugCache {
		if IsFileExist("./auto_qa.debug") {
			autoQaDebugCache = "true"
		} else {
			autoQaDebugCache = "false"
		}
	}
	if "true" == autoQaDebugCache && IsFileExist("./auto_qa.debug") {
		if a, err := ioutil.ReadFile("./auto_qa.debug"); nil == err {
			return RegexMatch(string(a), fmt.Sprintf(cond, args...))
		}
	}
	return false, nil
}

func DebugError(stage *log.Stage, cond string) error {
	ok, _ := HasAutoQaDebugCondf(cond)
	if ok {
		log.Key(stage, "[auto_qa.debug] "+cond)
		return fmt.Errorf("[auto_qa.debug] " + cond)
	}
	return nil
}

func DebugPause(cond string, args ...interface{}) {
	if ok, _ := HasAutoQaDebugCondf(cond, args...); !ok {
		return
	}
	log.Key(log.NewStage(), "DEBUG PAUSE %v", fmt.Sprintf(cond, args...))
	for {
		if ok, _ := HasAutoQaDebugCondf(cond, args...); ok {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	log.Key(log.NewStage(), "DEBUG CONT %v", fmt.Sprintf(cond, args...))
}

func DebugPanic(a interface{}) {
	ok, _ := HasAutoQaDebugCondf(fmt.Sprintf("%v", a))
	if ok {
		log.Key(log.NewStage(), "[auto_qa.debug] %v", a)
		panic(a)
	}
}
