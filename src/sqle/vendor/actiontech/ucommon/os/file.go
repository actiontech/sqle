package os

import (
	"actiontech/ucommon/conf"
	"actiontech/ucommon/log"
	user "actiontech/ucommon/user"
	"actiontech/ucommon/util"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

/*
	For example, path=/a/b/c
	1. If /a/b/c exists
		1.1 ensure c.perm>=perm
		1.2 ensure owner has access to /a,/a/b,/a/b/c
	2. If /a/b exists, /a/b/c not exists, mkdir /a/b/c
		2.1 ensure c.owner=owner
		2.2 ensure c.perm=perm
		2.3 ensure owner has access to /a,/a/b,/a/b/c
	3. If /a/b not exists, mkdir /a/b and /a/b/c
		3.1 ensure c.owner=owner and b.owner=owner
		3.2 ensure c.perm=perm and b.perm=perm
		3.3 ensure owner has access to /a,/a/b,/a/b/c
*/
func EnsureDir(stage *log.Stage, path, owner string /* empty = current user*/, perm os.FileMode) error {
	log.Detail(stage, "ensure dir %v ,owner(%v), perm(%v)", path, owner, perm)
	perm = perm | 0100

	var err error
	path, err = filepath.Abs(path)
	if nil != err {
		return err
	}

	stat, statErr := os.Stat(path)

	uid, gids, err := user.LookupUidGidByUser(owner)
	if nil != err {
		return err
	}

	if err := ensureDir(stage, path, owner, uid, gids[0], perm); nil != err {
		return err
	}

	if nil != statErr {
		if err := os.Chown(path, uid, gids[0]); nil != err {
			return err
		}
		if err := os.Chmod(path, perm); nil != err {
			return err
		}
	} else {
		if err := os.Chmod(path, stat.Mode()|perm); nil != err {
			return err
		}
	}
	return nil
}

func ensureDir(stage *log.Stage, path, owner string, uid, gid int, perm os.FileMode) error {
	var err error

	parent, _ := filepath.Split(path)
	if parent, err = filepath.Abs(parent); nil != err {
		return err
	}

	if "/" != parent {
		//ensure parent recursively
		if err := ensureDir(stage, parent, owner, uid, gid, perm); nil != err {
			return err
		}
	}

	//create current if need
	if stat, err := os.Stat(path); nil != err {
		if !os.IsNotExist(err) {
			return err
		}
		if err := os.Mkdir(path, perm); nil != err {
			return err
		}
		if err := os.Chown(path, uid, gid); nil != err {
			return err
		}
		if err := os.Chmod(path, perm); nil != err {
			return err
		}
	} else if !stat.IsDir() {
		return fmt.Errorf("%v is a file, expect dir", path)
	}

	return CheckPermission(stage, path, owner, 1)
}

func isGidsContains(gids []int, gid int) bool {
	for _, g := range gids {
		if g == gid {
			return true
		}
	}
	return false
}

func CheckPermission(stage *log.Stage, path, owner string /* empty = current user*/, permission int /* 0-7 */) error {
	uid, gids, err := user.LookupUidGidByUser(owner)
	if nil != err {
		return err
	}

	//check current access
	stat, err := os.Stat(path)
	if nil != err {
		return err
	}

	sysStat := stat.Sys().(*syscall.Stat_t)

	if 0 != uid {
		switch {
		case sysStat.Uid == uint32(uid):
			if os.FileMode(0100*permission) == (stat.Mode() & os.FileMode(0100*permission)) {
				return nil
			}
			return fmt.Errorf("no access to %v", path)
		case isGidsContains(gids, int(sysStat.Gid)):
			if os.FileMode(0010*permission) == (stat.Mode() & os.FileMode(0010*permission)) {
				return nil
			}
			return fmt.Errorf("no access to %v", path)
		default:
			if os.FileMode(0001*permission) != (stat.Mode() & os.FileMode(0001*permission)) {
				return fmt.Errorf("no access to %v", path)
			}
		}
	}
	return nil
}

func Remove(stage *log.Stage, path string) error {
	log.Detail(stage, "remove %v", path)

	if "" == strings.TrimRight(path, "/* ") {
		panic("cannot rm / path")
	}
	err := os.RemoveAll(path)
	if nil != err {
		log.Brief(stage, "remove %v error(%v)", path, err)
	}
	return err
}

func RecreateDir(stage *log.Stage, path, owner string /* empty = current user*/, perm os.FileMode) error {
	if err := Remove(stage, path); nil != err {
		return err
	}
	return EnsureDir(stage, path, owner, perm)
}

func ChownChgrp(stage *log.Stage, path, owner string) error {
	log.Detail(stage, "chown %v %v", owner, path)

	uid, gids, err := user.LookupUidGidByUser(owner)
	if nil != err {
		return err
	}
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if nil != err {
			return err
		}
		return os.Chown(path, uid, gids[0])
	})
}

func EnsureFile(stage *log.Stage, path, owner string, perm os.FileMode) error {
	log.Detail(stage, "ensure file %v ,owner(%v), perm(%v)", path, owner, perm)

	uid, gids, err := user.LookupUidGidByUser(owner)
	if nil != err {
		return err
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, perm)
	if nil != err {
		return err
	}
	if err := f.Chmod(perm); nil != err {
		return err
	}
	if err := f.Chown(uid, gids[0]); nil != err {
		return err
	}
	if err := f.Close(); nil != err {
		return err
	}
	return nil
}

var TempDirRoot string = "/tmp"

func SetTempDirRoot(path string) error {
	if a, err := filepath.Abs(path); nil != err {
		return err
	} else {
		TempDirRoot = a
	}
	return nil
}

func TempDir(stage *log.Stage, prefix string) (string, error) {
	if err := EnsureDir(stage, TempDirRoot, "", 0750); nil != err {
		return "", err
	}
	return ioutil.TempDir(TempDirRoot, prefix)
}

func Move(stage *log.Stage, from, to string, owner string, perm os.FileMode) error {
	log.Detail(stage, "move %v %v, owner(%v), perm(%v)", from, to, owner, perm)
	if err := os.Rename(from, to); nil != err {
		return err
	}
	if err := ChownChgrp(stage, to, owner); nil != err {
		return err
	}
	if err := os.Chmod(to, perm); nil != err {
		return err
	}
	return nil
}

func MoveOrCopyFile(stage *log.Stage, from, to string, owner string, perm os.FileMode) error {
	if err := Move(stage, from, to, owner, perm); nil != err {
		if strings.Contains(err.Error(), "invalid cross-device link") {
			return CopyFile(stage, from, to, owner, perm)
		}
		return err
	}
	return nil
}

//for text file busy error while copy file
func SafeMoveOrCopyFile(stage *log.Stage, from, to string, owner string, perm os.FileMode) error {
	if err := Move(stage, from, to, owner, perm); nil != err {
		if strings.Contains(err.Error(), "invalid cross-device link") {
			if !IsFileExist(to) {
				return CopyFile(stage, from, to, owner, perm)
			}

			if err := CopyFile(stage, from, to+".safe", owner, perm); nil != err {
				os.Remove(to + ".safe")
				return err
			}
			return Move(stage, to+".safe", to, owner, perm)
		}
		return err
	}
	return nil
}

func MoveOrCopyDir(stage *log.Stage, from, to string, owner string, perm os.FileMode) error {
	if err := Move(stage, from, to, owner, perm); nil != err {
		if strings.Contains(err.Error(), "invalid cross-device link") {
			return CopyDir(stage, from, to, owner, perm)
		}
		return err
	}
	return nil
}

func WriteFile(stage *log.Stage, path string, content string, owner string, perm os.FileMode) error {
	if err := EnsureFile(stage, path, owner, perm); nil != err {
		return err
	}
	return ioutil.WriteFile(path, []byte(content), perm)
}

func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	return nil == err
}

//in case disk full
func SafeWriteFile(stage *log.Stage, filename string, data string, owner string, perm os.FileMode) error {
	if !IsFileExist(filename) {
		if err := WriteFile(stage, filename, data, owner, perm); nil != err {
			return err
		}
		return nil
	}
	if err := WriteFile(stage, filename+".safe", data, owner, perm); nil != err {
		os.Remove(filename + ".safe")
		return err
	}
	return Move(stage, filename+".safe", filename, owner, perm)
}

func CopyFile(stage *log.Stage, srcPath, targetPath string, owner string, perm os.FileMode) error {
	log.Detail(stage, "copy %v %v, owner(%v), perm(%v)", srcPath, targetPath, owner, perm)

	src, err := os.OpenFile(srcPath, os.O_RDONLY, 0750)
	if nil != err {
		return err
	}
	defer src.Close()
	target, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if nil != err {
		return err
	}
	defer target.Close()
	if _, err := io.Copy(target, src); nil != err {
		return err
	}
	if err := ChownChgrp(stage, targetPath, owner); nil != err {
		return err
	}
	return os.Chmod(targetPath, perm)
}

func CopyDir(stage *log.Stage, srcPath, targetPath string, owner string, perm os.FileMode) error {
	log.Detail(stage, "copy %v %v, owner(%v), perm(%v)", srcPath, targetPath, owner, perm)

	ret, retCode, err := Cmdf(stage, "cp -R %v %v", srcPath, targetPath)
	if nil != err {
		return err
	}
	if 0 != retCode {
		return fmt.Errorf("cp %v to %v error(%v)", srcPath, targetPath, ret)
	}
	if err := ChownChgrp(stage, targetPath, owner); nil != err {
		return err
	}
	return os.Chmod(targetPath, perm)
}

//mainly for disk full
func SafeWriteConfigFile(stage *log.Stage, filename string, f *conf.ConfigFile) error {
	if !util.IsFileExist(filename) {
		return f.WriteConfigFile(filename, 0640, "", []string{})
	}

	if err := f.WriteConfigFile(filename+".safe", 0640, "", []string{}); nil != err {
		os.Remove(filename + ".safe")
		return err
	}
	return Move(stage, filename+".safe", filename, "", 0640)
}

func GetNumFromFile(file string) (uint64, error) {
	data, err := ioutil.ReadFile(file)
	if nil != err {
		return 0, err
	}
	dataStr := strings.TrimSpace(string(data))
	limit, err := strconv.ParseUint(dataStr, 10, 64)
	if nil != err {
		return 0, err
	}
	return limit, nil
}

func Mkdir(stage *log.Stage, path string, owner string, perm os.FileMode) error {
	uid, gids, err := user.LookupUidGidByUser(owner)
	if nil != err {
		return err
	}

	if err := os.Mkdir(path, perm); nil != err {
		return err
	}

	if err := os.Chmod(path, perm); nil != err {
		return err
	}

	if err := os.Chown(path, uid, gids[0]); nil != err {
		return err
	}

	log.Brief(stage, "mkdir %v, owner=%v, perm=%v", path, owner, perm)
	return nil
}

func GetFileContent(path string) (string, error) {
	if !IsFileExist(path) {
		return "", os.ErrNotExist
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil

}
