package quota

import (
	"actiontech/ucommon/log"
	"actiontech/ucommon/os"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"syscall"
)

func GetAndSetQuotaLimitMBWithSubDir(stage *log.Stage, user, subDir string, limit int) error {
	qi, err := CheckAndGetQuotaInfo(stage)
	if nil != err {
		return err
	}
	err = qi.SetUserLimitMBWithSubDir(user, subDir, limit)
	if nil != err {
		return err
	}
	return nil
}

func CheckAndGetQuotaInfo(stage *log.Stage) (*QuotaInfo, error) {
	if err := CheckQuotaInstalled(stage); nil != err {
		return nil, err
	}
	return GetQuotaInfo(stage)
}

func CheckQuotaInstalled(stage *log.Stage) error {
	stage.Enter("check_quota_installed")
	defer stage.Exit()

	if _, err := os.Cmdf2(stage, "which quota"); nil != err {
		return errors.New("quota is not installed")
	}

	if _, err := os.Cmdf2(stage, "which repquota"); nil != err {
		return errors.New("repquota is not installed")
	}

	if _, err := os.Cmdf2(stage, "which quotaon"); nil != err {
		return errors.New("quotaon is not installed")
	}

	if _, err := os.Cmdf2(stage, "which setquota"); nil != err {
		return errors.New("setquota is not installed")
	}
	return nil
}

//GetQuotaInfo get quota information from local server
func GetQuotaInfo(stage *log.Stage) (*QuotaInfo, error) {
	qi := newQuotaInfo(stage)
	err := qi.getQuotaMountInfo()
	return qi, err
}

type QuotaUser struct {
	User        string
	BlockUsed   uint64
	BlockLimits uint64
	Grace       string
}

type QuotaDevice struct {
	FileSystem string
	MountDir   string
	Size       uint64
	Used       uint64
	Enable     bool
	Users      map[string]*QuotaUser
}

type QuotaInfo struct {
	Devices map[string]*QuotaDevice
	stage   *log.Stage
}

func newQuotaInfo(stage *log.Stage) *QuotaInfo {
	return &QuotaInfo{
		Devices: map[string]*QuotaDevice{},
		stage:   stage,
	}
}

func (qi *QuotaInfo) MatchQuotaDir(absDir string) (dir string, isMatch bool) {
	if "" == absDir {
		return "", false
	}
	for _, device := range qi.Devices {
		mountDir := device.MountDir
		// if devices has "/a/b", "/a", dir return "/a/b"
		if strings.HasPrefix(absDir, mountDir) && len(mountDir) > len(dir) {
			dir = mountDir
		}
	}
	if "" == dir {
		return "", false
	}
	return dir, true
}

func (qi *QuotaInfo) getQuotaMountInfo() error {
	qi.stage.Enter("get_quota_status")
	defer qi.stage.Exit()

	mounts, err := os.Cmdf2(qi.stage, "mount")
	if nil != err {
		return err
	}
	for _, line := range strings.Split(mounts, "\n") {
		if !strings.Contains(line, "usrquota") {
			continue
		}
		// mount format:
		// "<file system> on <mount point> type <type> <options>"
		//   1            2   3            4     5      6
		mount := strings.Split(line, " ")
		if 6 == len(mount) {
			mountDir := mount[2]
			if !strings.HasPrefix(mountDir, "/") {
				return errors.New(fmt.Sprintf(`mount dir: %s, the prefix must be "/"`, mountDir))
			}
			qi.Devices[mountDir] = &QuotaDevice{FileSystem: mount[0], MountDir: mountDir, Users: map[string]*QuotaUser{}}
		}
	}
	if 0 >= len(qi.Devices) {
		return errors.New("no quota disk")
	}
	return qi.getQuotaDeviceStatus()
}

func (qi *QuotaInfo) getQuotaDeviceStatus() error {
	if 0 >= len(qi.Devices) {
		return nil
	}
	for mountDir, device := range qi.Devices {
		fs := syscall.Statfs_t{}
		err := syscall.Statfs(mountDir, &fs)
		if nil != err {
			return err
		}
		//unit KB
		device.Size = fs.Blocks * uint64(fs.Bsize) / 1024
		device.Used = (fs.Blocks - fs.Bfree) * uint64(fs.Bsize) / 1024

		enable, err := checkUserQuotaEnable(qi.stage, mountDir)
		if nil != err {
			return err
		}
		device.Enable = enable
	}
	return qi.getQuotaUserStatus()
}

func (qi *QuotaInfo) getQuotaUserStatus() error {
	for k, v := range qi.Devices {
		cmd := fmt.Sprintf(`SUDO repquota %s |grep "\- "|awk -F " " '{print $1,$3,$5}'`, k)
		report, err := os.Cmdf2(qi.stage, cmd)
		if nil != err {
			return err
		}
		for _, line := range strings.Split(report, "\n") {
			info := strings.Split(line, " ")
			if 3 != len(info) {
				return errors.New(fmt.Sprintf("fail to exec \"repquota\" status.line: %s", line))
			}
			used, err := strconv.ParseUint(info[1], 10, 64)
			if nil != err {
				return errors.New(fmt.Sprintf("fail to exec \"repquota\" status.line: %s", line))
			}
			size, err := strconv.ParseUint(info[2], 10, 64)
			if nil != err {
				return errors.New(fmt.Sprintf("fail to exec \"repquota\" status.line: %s", line))
			}
			v.Users[info[0]] = &QuotaUser{User: info[0], BlockUsed: used, BlockLimits: size}
		}
	}
	return nil
}

func (qi *QuotaInfo) CancelUserLimitMBWithSubDir(user, subDir string) error {
	return qi.SetUserLimitMBWithSubDir(user, subDir, 0)
}

func (qi *QuotaInfo) SetUserLimitMBWithSubDir(user, subDir string, limit int) error {
	mountDir, ok := qi.MatchQuotaDir(subDir)
	if !ok {
		return errors.New(fmt.Sprintf("\"%s\" isn't quota file system", subDir))
	}
	return qi.SetUserLimitMB(user, mountDir, limit)
}

//SetQuotaLimit set quota limit from local server, set limit numbers in human friendly units (MB).
func (qi *QuotaInfo) SetUserLimitMB(user, dir string, limit int) error {
	return qi.setUserLimit(user, dir, limit*1024)
}

func (qi *QuotaInfo) setUserLimit(user, dir string, limit int) error {
	err := qi.filterUserLimit(user, dir, limit)
	if nil != err {
		return err
	}

	if !qi.Devices[dir].Enable {
		err := setUserQuotaEnable(qi.stage, dir)
		if nil != err {
			return err
		}
	}
	cmd := fmt.Sprintf("SUDO setquota -u %s %d %d 0 0 %s", user, limit, limit, dir)
	_, err = os.Cmdf2(qi.stage, cmd)
	if nil != err {
		return err
	}
	return nil
}

func (qi *QuotaInfo) filterUserLimit(user, dir string, limit int) error {
	device, ok := qi.Devices[dir]
	if !ok {
		return errors.New(fmt.Sprintf("mount dir \"%s\" isn't quota file system", dir))
	}
	if 0 == limit {
		return nil
	}
	if 0 > limit {
		return errors.New(fmt.Sprintf("limit size must be greater than 0, current is %d", limit))
	}
	if limit > int(device.Size) {
		return errors.New(fmt.Sprintf("file system %s total size: %d, limit size: %d", dir, device.Size, limit))
	}
	userInfo, ok := device.Users[user]
	if !ok {
		return nil
	}
	if limit < int(userInfo.BlockUsed) {
		return errors.New(fmt.Sprintf("user:%s used size:%d, limit size:%d", user, userInfo.BlockUsed, limit))
	}
	return nil
}

func checkUserQuotaEnable(stage *log.Stage, fileSystem string) (bool, error) {
	cmd := fmt.Sprintf(`SUDO quotaon -up %s`, fileSystem)
	_, code, err := os.Cmdf(stage, cmd)
	if nil != err {
		return false, errors.New(fmt.Sprintf("message: %s, cmd: %s", err, cmd))
	}
	if 1 == code {
		return true, nil
	}
	return false, nil
}

func setUserQuotaEnable(stage *log.Stage, fileSystem string) error {
	cmd := fmt.Sprintf(`SUDO quotaon -u %s`, fileSystem)
	_, err := os.Cmdf2(stage, cmd)
	return err
}
