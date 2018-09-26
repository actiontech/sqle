package user

import (
	"os"
)

func IsCurrentUser(user string) (bool, error) {
	uid, _, err := LookupUidGidByUser(user)
	if nil != err {
		return false, err
	}
	currentUid := os.Getuid()
	return uid == currentUid, nil
}
