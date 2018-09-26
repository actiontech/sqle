// +build cgo

package user

import (
	"os"
	"os/user"
	"strconv"
)

func LookupUidGidByUser(username string) (uid int, gids []int, err error) {
	if "" != username {
		u, err := user.Lookup(username)
		if nil != err {
			return 0, []int{}, err
		}
		uid, _ = strconv.Atoi(u.Uid)

		// On different systems, gid's sorting is not the same.
		// In order to get a consistent gid, we set the primary gid in the first place.
		pGid, _ := strconv.Atoi(u.Gid)
		gids = []int{pGid}
		{
			groupIds, err := u.GroupIds()
			if nil != err {
				return 0, []int{}, err
			}
			for _, gidStr := range groupIds {
				gid, _ := strconv.Atoi(gidStr)
				if pGid == gid {
					continue
				}
				gids = append(gids, gid)
			}
		}
	} else {
		uid = os.Getuid()
		pGid := os.Getgid()
		gids = []int{pGid}
		if groups, err := os.Getgroups(); nil == err {
			for _, gid := range groups {
				if pGid == gid {
					continue
				}
				gids = append(gids, gid)
			}

		}
	}

	return uid, gids, nil
}
