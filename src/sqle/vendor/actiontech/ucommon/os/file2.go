//+build linux

package os

import (
	"fmt"
	"golang.org/x/sys/unix"
	os_ "os"
	"runtime"
	"syscall"
	user_ "actiontech/ucommon/user"
)

const (
	ACCESS_W = unix.W_OK
	ACCESS_R = unix.R_OK
	ACCESS_X = unix.X_OK
)

func CheckAccess(user string, path string, requiredAccess uint32) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	originalUid := os_.Getuid()
	originalGid := os_.Getgid()

	var uid int
	var gids []int

	var err error
	uid, gids, err = user_.LookupUidGidByUser(user)
	if nil != err {
		return err
	}


	var accessErr error

	for _, gid := range gids {
		{
			//mystery: if miss write() syscall, capset will throw EPERM.
			//totally make no sence
			fmt.Printf("")
		}

		originalCap, err := GetCap()
		if nil != err {
			return fmt.Errorf("GetCap error: %v", err)
		}

		//change user, and clear CAP
		{
			if err := SetKeepCaps(); nil != err {
				return fmt.Errorf("SetKeepCaps error: %v", err)
			}

			if err := syscall.Setresgid(gid, gid, gid); nil != err {
				return fmt.Errorf("Setresgid error: %v", err)
			}

			if err := syscall.Setresuid(uid, uid, uid); nil != err {
				return fmt.Errorf("Setresuid error: %v", err)
			}

			{
				//mystery: if miss write() syscall, capset will throw EPERM.
				//totally make no sence
				fmt.Printf("")
			}

			newCap := UserCapDataStruct{
				Effective:   1<<6 | 1<<7, /*CAP_SETGID,CAP_SETUID*/
				Permitted:   originalCap.Permitted,
				Inheritable: originalCap.Inheritable,
			}

			if err := SetCap(newCap); nil != err {
				return fmt.Errorf("SetCap error: %v", err)
			}
		}

		accessErr = unix.Access(path, requiredAccess)

		//change user back, and clear
		{
			if err := SetKeepCaps(); nil != err {
				return fmt.Errorf("SetKeepCaps error: %v", err)
			}

			if err := syscall.Setresgid(originalGid, originalGid, originalGid); nil != err {
				return fmt.Errorf("Setresgid as original error: %v", err)
			}

			if err := syscall.Setresuid(originalUid, originalUid, originalUid); nil != err {
				return fmt.Errorf("Setresuid as original error: %v", err)
			}

			{
				//mystery: if miss write() syscall, capset will throw EPERM.
				//totally make no sence
				fmt.Printf("")
			}

			if err := SetCap(originalCap); nil != err {
				return fmt.Errorf("SetCap as original error: %v", err)
			}
		}

		if nil == accessErr {
			return nil
		}
	}

	return accessErr
}
