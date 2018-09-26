// +build !workaround_prof_self_perm_denied

package os

import (
	"runtime"
	"path/filepath"
	"os"
	"path"
)

func GetExecDir() string {
	if "darwin" == runtime.GOOS {
		a, _ := filepath.Abs(".")
		return a
	}
	p, err := os.Readlink("/proc/self/exe")
	if nil != err {
		panic("getRootDir got err:" + err.Error())
	}
	return path.Dir(p)
}