package os

import (
	"actiontech/ucommon/log"
	"actiontech/ucommon/util"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

func Mount(stage *log.Stage, dev string, point string, readonly bool, mysqlRunUser string) error {
	if err := EnsureDir(stage, point, mysqlRunUser, 0750); nil != err {
		return err
	}
	if readonly {
		if ret, retCode, err := Cmdf(stage, "SUDO mount -o ro %s %s", dev, point); nil != err || (0 != retCode && !(32 == retCode && strings.Contains(ret, "already mounted"))) {
			return fmt.Errorf("mount %v(%v) failed", dev, point)
		}
	} else {
		if ret, retCode, err := Cmdf(stage, "SUDO mount %s %s", dev, point); nil != err || (0 != retCode && !(32 == retCode && strings.Contains(ret, "already mounted"))) {
			return fmt.Errorf("mount %v(%v) failed", dev, point)
		}
	}
	return nil
}

var umountNotMountedRegexp = regexp.MustCompile("not .*mount")

func Umount(stage *log.Stage, point string) error {

	if _, err := os.Stat(point); nil != err && os.IsNotExist(err) {
		return nil
	}

	//3 try
	for i := 0; i < 3; i++ {
		if ret, retCode, err := Cmdf(stage, "SUDO umount %s", point); nil != err {
			return fmt.Errorf("umount %v failed, error(%v)", point, err)
		} else if (0 == retCode) || (0 != retCode && umountNotMountedRegexp.MatchString(ret)) {
			return nil
		}
		Cmdf(stage, "SUDO lsof +D %s", point) //for log
		time.Sleep(100 * time.Millisecond)
	}

	util.DebugPause("pause after umount")

	Cmdf(stage, "SUDO lsof +D %s", point) //for log

	// //2 try
	// for i := 0; i < 2; i++ {
	// 	//fuser
	// 	if ret, _, err := Cmdf(stage, "SUDO fuser -m %s", point); nil != err /* ignore retCode since retCode=1 means no pid found */ {
	// 		return fmt.Errorf("umount %v failed, fuser error(%v)", point, err)
	// 	} else {
	// 		pids := util.SplitAndTrimSpace(ret, " ")
	// 		if dry.StringInSlice("1", pids) {
	// 			//root device
	// 			return nil
	// 		}
	// 		currentPid := os.Getpid()
	// 		if dry.StringInSlice(fmt.Sprintf("%v", currentPid), pids) {
	// 			return fmt.Errorf("umount %v failed, current process(%v) is using it", point, currentPid)
	// 		}
	// 		Cmdf(stage, "SUDO kill %v", ret)
	// 	}

	// 	//umount again
	// 	if ret, retCode, err := Cmdf(stage, "SUDO umount %s", point); nil != err {
	// 		return fmt.Errorf("umount %v failed, error(%v)", point, err)
	// 	} else if (0 == retCode) || (0 != retCode && umountNotMountedRegexp.MatchString(ret)) {
	// 		return nil
	// 	}
	// }
	return fmt.Errorf("umount %v failed", point)
}
