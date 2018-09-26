// +build workaround_prof_self_perm_denied

package os

import (
	"path/filepath"
)

func GetExecDir() string {
	a, _ := filepath.Abs(".")
	return a
}