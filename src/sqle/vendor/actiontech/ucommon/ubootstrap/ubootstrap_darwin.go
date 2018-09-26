//+build darwin

package ubootstrap

import (
	"fmt"
)

func setresgid(rgid int, egid int, sgid int) (err error) {
	return fmt.Errorf("setresgid not supported in darwin")
}

func setresuid(ruid int, euid int, suid int) (err error) {
	return fmt.Errorf("setresuid not supported in darwin")
}

func RotateStdout() error {
	fmt.Errorf("rotate stdout is not support on macOS")
	return nil
}
