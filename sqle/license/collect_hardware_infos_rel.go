//go:build release
// +build release

package license

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"sort"
	"strings"
)

var (
	dockerDevReg = regexp.MustCompile("([^\\s]+) on /etc/hostname type")
	devReg       = regexp.MustCompile("(/[^\\s]+) on / type")
	blkidReg     = regexp.MustCompile("UUID=\"([^ ]+)\"")
	macsReg      = regexp.MustCompile("link/ether ([^ ]+)")
	encoding     = base64.NewEncoding("012345ghijklmnopq6789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefrstuvwxyz_~")
)

func CollectHardwareInfo() (string, error) {
	keys := make([]string, 0)

	bootDevUuid, err := getBootDevUuid()
	if nil != err {
		return "", err
	}

	keys = append(keys, bootDevUuid)

	//ifconfig macs
	output, err := cmd("ip addr")
	if nil != err {
		return "", err
	}
	ipaddrs := macsReg.FindAllStringSubmatch(output, -1)
	macs := make([]string, 0, len(ipaddrs))
	for _, ipaddr := range ipaddrs {
		if len(ipaddr) < 1 {
			continue
		}
		// be compatible with early version license
		macs = append(macs, "HWaddr "+strings.ToUpper(ipaddr[1]))
	}
	sort.Strings(macs)
	for _, mac := range macs {
		keys = append(keys, mac)
	}

	//encode
	return encoding.EncodeToString([]byte(strings.Join(keys, "|"))), nil
}

func getBootDevUuid() (string, error) {
	bootDevUuid := ""
	output, err := cmd("mount -l")
	if nil != err {
		return "", err
	}
	matches := dockerDevReg.FindStringSubmatch(output)
	if nil != matches {
		return " ", nil // ignore docker
	}
	matches = devReg.FindStringSubmatch(output)
	if nil == matches {
		return "", fmt.Errorf("show \"/\" mount got empty")
	}

	bootDev := matches[1]
	output, err = cmd(fmt.Sprintf("blkid %v", bootDev))
	if nil != err {
		return "", err
	}
	matches = blkidReg.FindStringSubmatch(output)
	if nil == matches {
		return "", fmt.Errorf("read root DEV uuid got empty")
	}
	bootDevUuid = matches[1]
	return bootDevUuid, nil
}

func cmd(str string) (string, error) {
	cmd := exec.Command("bash", "--noprofile", "--norc", "-c", str)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	defer stdout.Close()
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}
	defer stderr.Close()

	if err = cmd.Start(); err != nil {
		return "", err
	}

	errBytes, err := ioutil.ReadAll(stderr)
	if err != nil {
		return "", err
	}
	if len(errBytes) > 0 {
		return "", fmt.Errorf((string(errBytes)))
	}

	opBytes, err := ioutil.ReadAll(stdout)
	return string(opBytes), err
}
