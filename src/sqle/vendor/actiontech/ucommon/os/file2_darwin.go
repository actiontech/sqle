//+build darwin

package os

import "golang.org/x/sys/unix"

const (
	ACCESS_W = unix.W_OK
	ACCESS_R = unix.R_OK
	ACCESS_X = unix.X_OK
)

func CheckAccess(user string, path string, requiredAccess uint32) error {
	return nil
}
