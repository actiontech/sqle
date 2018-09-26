// +build !cgo

package user

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
)

func LookupUidGidByUser(userName string) (uid int, gids []int, err error) {
	if "" != userName {
		ret, err := exec.Command("id", "-u", userName).CombinedOutput()
		if nil != err {
			if strings.Contains(strings.ToLower(string(ret)), "no such user") {
				return 0, []int{}, user.UnknownUserError(userName)
			}
			return 0, []int{}, fmt.Errorf("id -u %v error: %v", userName, err)
		}
		uid, _ = strconv.Atoi(strings.TrimSpace(string(ret)))

		ret, err = exec.Command("id", "-G", userName).Output()
		if nil != err {
			return 0, []int{}, fmt.Errorf("id -g %v error: %v", userName, err)
		}

		gids = []int{}
		for _, seg := range strings.Split(string(ret), " ") {
			seg = strings.TrimSpace(string(seg))
			if "" == seg {
				continue
			}
			gid, _ := strconv.Atoi(seg)
			gids = append(gids, gid)
		}
		return uid, gids, nil
	} else {
		uid = os.Getuid()
		gid := os.Getgid()
		return uid, []int{gid}, nil
	}

}
