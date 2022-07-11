//go:build release
// +build release

package license

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/moby/sys/mountinfo"
	"golang.org/x/sys/unix"
)

var (
	encoding = base64.NewEncoding("012345ghijklmnopq6789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefrstuvwxyz_~")
	// Location to perform UUID lookup
	uuidDirectory = "/dev/disk/by-uuid"
)

func CollectHardwareInfo() (string, error) {
	keys := make([]string, 0)

	bootDevUuid, err := getBootDevUuid()
	if nil != err {
		return "", err
	}

	keys = append(keys, bootDevUuid)

	netInterfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	var macs []string
	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if macAddr != "" {
			macs = append(macs, "HWaddr "+strings.ToUpper(macAddr))
		}
	}

	sort.Strings(macs)
	for _, mac := range macs {
		keys = append(keys, mac)
	}

	//encode
	return encoding.EncodeToString([]byte(strings.Join(keys, "|"))), nil
}

func getBootDevUuid() (string, error) {
	mounts, err := mountinfo.GetMounts(nil)
	if err != nil {
		return "", err
	}

	for _, mount := range mounts {
		if mount.Mountpoint == "/etc/hostname" {
			return "", nil // ignore docker
		}
	}

	for _, mount := range mounts {
		if mount.Mountpoint == "/" {
			deviceNumberFromMount, err := getDeviceNumber(mount.Source)
			if err != nil {
				return "", err
			}

			dirContents, err := ioutil.ReadDir(uuidDirectory)
			if err != nil {
				return "", err
			}

			for _, fileInfo := range dirContents {
				if fileInfo.Mode()&os.ModeSymlink != os.ModeSymlink {
					continue // ignore non-symlink
				}

				uuid := fileInfo.Name()
				uuidSymlinkPath := filepath.Join(uuidDirectory, uuid)
				deviceNumberFromUUID, err := getDeviceNumber(uuidSymlinkPath)
				if err != nil {
					return "", err
				}

				if deviceNumberFromMount == deviceNumberFromUUID {
					return uuid, nil
				}
			}
		}
	}

	return "", fmt.Errorf("\"/\" mount or uuid is empty")
}

// DeviceNumber represents a combined major:minor device number.
type DeviceNumber uint64

func (num DeviceNumber) String() string {
	return fmt.Sprintf("%d:%d", unix.Major(uint64(num)), unix.Minor(uint64(num)))
}

func getDeviceNumber(path string) (DeviceNumber, error) {
	var stat unix.Stat_t
	if err := unix.Stat(path, &stat); err != nil {
		return 0, err
	}
	return DeviceNumber(stat.Rdev), nil
}
