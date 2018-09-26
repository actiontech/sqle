//+build linux

package ubootstrap

import (
	"actiontech/ucommon/log"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"
)

func setresgid(rgid int, egid int, sgid int) (err error) {
	return syscall.Setresgid(rgid, egid, sgid)
}

func setresuid(ruid int, euid int, suid int) (err error) {
	return syscall.Setresuid(ruid, euid, suid)
}

func RotateStdout() error {
	// check if stdout is a file
	stat, err := os.Stdout.Stat()
	if nil != err {
		return errors.New("get stdout stat fail: " + err.Error())
	}
	if 0 != (stat.Mode() & os.ModeType) {
		return nil
	}

	// rotate
	path, err := os.Readlink("/proc/self/fd/1")
	if nil != err {
		return errors.New("read stdout link fail: " + err.Error())
	}

	init := ""
	last := ""
	if idx := strings.LastIndexAny(path, `./`); -1 == idx || '/' == path[idx] {
		init = path
	} else {
		init = path[:idx]
		last = path[idx:]
	}

	if strings.HasSuffix(init, "_0") {
		fmt.Println("\n\n\n--------------------------------- ROTATE LOG ---------------------------------")
		fmt.Println("START TIME: ", time.Now().Format(log.LogTimeStamp))
		return nil
	}
	fmt.Println("START TIME: ", time.Now().Format(log.LogTimeStamp))

	os.Remove(init + "_2" + last)
	os.Rename(init+"_1"+last, init+"_2"+last)
	os.Rename(init+"_0"+last, init+"_1"+last)
	return os.Rename(path, init+"_0"+last)
}
