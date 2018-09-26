package os

import (
	"actiontech/ucommon/log"
	"errors"
	"fmt"
	"io/ioutil"
	"log/syslog"
	"net"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/tatsushid/go-fastping"
	"github.com/ungerik/go-dry"
	"path/filepath"
	"regexp"
)

func HasNic(nic string) (bool, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return false, err
	}
	for _, ele := range interfaces {
		if nic == ele.Name {
			return true, nil
		}
	}
	return false, nil
}

func getLocalIpNetAndNicBySip(sip string) (retIpNet *net.IPNet, retNic string, err error) {
	sipIp := net.ParseIP(sip)
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, ele := range interfaces {
		addrs, err := ele.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			if ipnet.Contains(sipIp) {
				return ipnet, ele.Name, nil
			}
		}
	}
	return nil, "", errors.New("the host is not in the same network segment with SIP, SIP couldn't bind")
}

func GetAllIps() ([]string, error) {
	interfaces, err := net.Interfaces()
	if nil != err {
		return nil, err
	}
	ret := make([]string, 0)
	for _, ele := range interfaces {
		addrs, err := ele.Addrs()
		if nil != err {
			continue
		}
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			ret = append(ret, ipnet.IP.String())
		}
	}
	return ret, nil
}

func GetLocalNicBySip(sip string) (ret string, err error) {
	_, ret, err = getLocalIpNetAndNicBySip(sip)
	return
}

func Pings(stage *log.Stage, ips []string, timeoutSeconds int) ([]string, []string) {
	pingable := []string{}
	retMutex := sync.Mutex{}

	p := fastping.NewPinger()
	p.MaxRTT = time.Duration(timeoutSeconds) * time.Second
	for _, ip := range ips {
		p.AddIP(ip)
	}
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		retMutex.Lock()
		defer retMutex.Unlock()
		pingable = append(pingable, addr.IP.String())

		if len(ips) == len(pingable) {
			go p.Stop()
		}
	}
	p.OnIdle = func() {}
	if err := p.Run(); nil != err {
		log.Key(stage, "pinger error: %v", err)
	}
	retMutex.Lock()
	defer retMutex.Unlock()
	notPingable := []string{}
	for _, ip := range ips {
		if !dry.StringInSlice(ip, pingable) {
			notPingable = append(notPingable, ip)
		}
	}
	return pingable, notPingable
}

func IsPortInUse(stage *log.Stage, port string) bool {
	_, retCode, err := Cmdf(stage, "exec 6<>/dev/tcp/127.0.0.1/"+port)
	defer func() {
		Cmdf(stage, "exec 6>&-")
		Cmdf(stage, "exec 6<&-")
	}()
	return nil == err && 0 == retCode
}

func IsBindSip(stage *log.Stage, ip string) bool {
	return HasLocalIp(stage.Go(), ip)
}

func HasLocalIp(stage *log.Stage, ip string) bool {
	if output, retCode, err := Cmdf(stage, "ip addr"); nil != err || 0 != retCode {
		return false
	} else {
		return strings.Contains(output, ip+"/")
	}
}

var rootDir string

//first call should be kept in OS main thread, otherwise maybe "permission denied"
func GetRootDir() string {
	if "" == rootDir {
		rootDir = path.Dir(GetExecDir())
	}
	return rootDir
}

func Uncompress(stage *log.Stage, from, to string, stripComponents int) error {
	if err := EnsureDir(stage, to, "", 0750); nil != err {
		return err
	}
	if ret, retCode, err := Cmdf(stage, `tar x -z -C %v --strip-components %v -f %v`, to, stripComponents, from); nil != err || 0 != retCode {
		return fmt.Errorf("tar failed, err=%v, retCode=%v, ret=%v", err, retCode, ret)
	}
	return nil
}

func LockFile(stage *log.Stage, file *os.File) error {
	return syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
}

func PidUsePort(stage *log.Stage, pid, port int) bool {
	_, retCode, err := Cmdf(stage, "lsof -i :%v -t | grep ^%v$", port, pid)
	return nil == err && 0 == retCode
}

func IsProcessExist(stage *log.Stage, pid int) bool {
	p, err := os.FindProcess(pid)
	if nil != err {
		return false
	}
	err = p.Signal(syscall.Signal(0))
	return nil == err || "operation not permitted" == err.Error()
}

func GetProcessPids(stage *log.Stage, runUser, processName string) ([]string, error) {
	if ret, retCode, err := Cmdf(stage, "pgrep -U%v -f \"%v$\"", runUser, processName); nil != err {
		return nil, fmt.Errorf("pgrep -U%v -f \"%v$\" got err=%v", runUser, processName, err)
	} else if 0 != retCode {
		return make([]string, 0), nil
	} else {
		return strings.Split(strings.TrimSpace(ret), "\n"), nil
	}
}

func IsPidProcess(stage *log.Stage, runUser, pidString, processName string) bool {
	if pids, err := GetProcessPids(stage, runUser, processName); nil == err {
		return dry.StringInSlice(pidString, pids)
	}
	return false
}

func ChkConfig(stage *log.Stage, initd string) error {
	var chkconfig string
	if IsSystemd() {
		chkconfig = SystemdEnableStr(initd)
	} else {
		chkconfig = ChkConfigStrSudo(initd)
	}
	_, err := Cmdf2(stage, "%v", chkconfig)
	if nil != err {
		return fmt.Errorf("chkconfig error: %v", err)
	}
	return nil
}

func ChkConfigStr(initd string) string {
	return fmt.Sprintf("egrep -o 'Ubuntu|Debian' /etc/issue && update-rc.d %v defaults 30 99 || chkconfig --add %v", initd, initd)
}

func ChkConfigStrSudo(initd string) string {
	return fmt.Sprintf("egrep -o 'Ubuntu|Debian' /etc/issue && SUDO update-rc.d %v defaults 30 99 || SUDO chkconfig --add %v", initd, initd)
}

func UnChkConfig(stage *log.Stage, initd string) error {
	var unchkconfig string
	if IsSystemd() {
		unchkconfig = SystemdDisableStr(initd)
	} else {
		unchkconfig = UnChkConfigStrSudo(initd)
	}
	_, err := Cmdf2(stage, "%v", unchkconfig)
	if nil != err {
		return fmt.Errorf("chkconfig error: %v", err)
	}
	return nil
}

func UnChkConfigStr(initd string) string {
	return fmt.Sprintf("egrep -o 'Ubuntu|Debian' /etc/issue && update-rc.d %v remove || chkconfig --del %v", initd, initd)
}

func UnChkConfigStrSudo(initd string) string {
	return fmt.Sprintf("egrep -o 'Ubuntu|Debian' /etc/issue && SUDO update-rc.d %v remove || SUDO chkconfig --del %v", initd, initd)
}

var isSystemd bool
var initIsSystemd sync.Once

func IsSystemd() bool {
	initIsSystemd.Do(func() {
		if data, err := ioutil.ReadFile("/proc/1/comm"); nil != err {
			// only linux 2.6.33 has this file, so if it is not exist, the system can not be systemd
			isSystemd = false
		} else {
			isSystemd = strings.Contains(string(data), "systemd")
		}
	})

	return isSystemd
}

func SystemdEnableStr(initd string) string {
	return fmt.Sprintf("SUDO systemctl daemon-reload && SUDO systemctl enable %v.service", initd)
}

func SystemdDisableStr(initd string) string {
	return fmt.Sprintf("SUDO systemctl disable %v.service", initd)
}

func Kill(pid int) error {
	return syscall.Kill(pid, syscall.SIGUSR1)
}

func CheckCommandExist(stage *log.Stage, cmd string) error {
	_, retCode, err := Cmdf(stage, "which %v", cmd)
	if nil != err {
		return err
	}
	if 0 != retCode {
		return fmt.Errorf("Cannot find \"%v\" command", cmd)
	}
	return nil
}

func CheckPingHostname(stage *log.Stage) error {
	_, retCode, err := Cmdf(stage, "ping -c 1 $(hostname)")
	if nil != err {
		return err
	}
	if 0 != retCode {
		return fmt.Errorf("Cannot ping hostname")
	}
	return nil
}

func CheckLibaio(stage *log.Stage) error {
	_, retCode, err := Cmdf(stage, "ldconfig -p | grep libaio")
	if nil != err {
		return err
	}
	if 0 != retCode {
		return fmt.Errorf("Cannot find libaio")
	}
	return nil
}

func CheckPerlModule(stage *log.Stage) error {
	_, retCode, err := Cmdf(stage, "perl -e 'use Data::Dumper'")
	if nil != err {
		return err
	}
	if 0 != retCode {
		return fmt.Errorf("Cannot find perl Data::Dumper")
	}
	_, retCode, err = Cmdf(stage, "perl -e 'use Time::HiRes'")
	if nil != err {
		return err
	}
	if 0 != retCode {
		return fmt.Errorf("Cannot find perl Time::HiRes")
	}
	_, retCode, err = Cmdf(stage, "perl -e 'use DBD::mysql'")
	if nil != err {
		return err
	}
	if 0 != retCode {
		return fmt.Errorf("Cannot find perl DBD::mysql")
	}

	_, retCode, err = Cmdf(stage, "perl -e 'use Digest::MD5'")
	if nil != err {
		return err
	}
	if 0 != retCode {
		return fmt.Errorf("Cannot find perl Digest::MD5")
	}
	return nil
}

func ClearScreen() {
	clearCmd := exec.Command("clear")
	clearCmd.Stdout = os.Stdout
	clearCmd.Run()
}

func GetAllDescendantPids(pid int) []int {
	pids := []int{pid}
	i := 0
	for i < len(pids) {
		output, _ := exec.Command("pgrep", "-P", fmt.Sprintf("%v", pids[i])).CombinedOutput()
		for _, line := range strings.Split(string(output), "\n") {
			line = strings.TrimSpace(line)
			pid, err := strconv.Atoi(line)
			if nil != err {
				continue
			}
			if 1 != pid && pid != pids[i] {
				pids = append(pids, pid)
			}
		}
		i++
	}
	return pids
}

func KillWithAllDescendantPids(pid int) error {
	for _, pid := range GetAllDescendantPids(pid) {
		syscall.Kill(pid, syscall.SIGKILL)
	}
	return nil
}

func CheckRootDirExistLoop(rootDir string) {
	for {
		if IsFileExist(rootDir) {
			time.Sleep(30 * time.Second)
			continue
		}
		errMsg := fmt.Sprintf("work dir %v gone, suicide", rootDir)
		if errSyslog, err := syslog.New(syslog.LOG_ERR, os.Args[0]); nil == err {
			errSyslog.Err(errMsg)
			errSyslog.Close()
		}
		fmt.Println(errMsg)
		os.Exit(1)
	}
}

func ErrExit(e error) {
	if errSyslog, err := syslog.New(syslog.LOG_ERR, os.Args[0]); nil == err {
		errSyslog.Err(e.Error())
		errSyslog.Close()
	}
	fmt.Fprintf(os.Stderr, "ERROR: %v\n", e)
	os.Exit(1)
}

func HaltIfShutdown(stage *log.Stage) {
	if IsFileExist("/tmp/.universe_is_reboot") {
		log.UserInfo(stage, "System halt")
		os.Exit(99)
	}
}

func IsProcessMatchPidFile(pid int, pidFile string) bool {
	if cmdline, err := ioutil.ReadFile(fmt.Sprintf("/proc/%v/cmdline", pid)); nil == err {
		segs := strings.Split(string(cmdline), string(byte(0)))
		if len(segs) > 0 && pidFile == fmt.Sprintf("%v.pid", filepath.Base(segs[0])) {
			return true
		}

		/*
			uelasticsearch cmdline[0] is "java", which is not match previous rule
			here, we check every arg, find if any arg ends with "uelasticsearch.pid", which means java take pidfile as arguments
		 */
		for _, seg := range segs {
			if strings.HasSuffix(seg, pidFile) {
				return true
			}
		}
	}
	return false
}

var currentOsId string
var initCurrentOsId sync.Once

func GetOsId() string {
	initCurrentOsId.Do(func() {
		if bs, err := ioutil.ReadFile("/etc/os-release"); nil != err {
			currentOsId = "UNKNOWN-OS"
		} else {
			matches := regexp.MustCompile(`(?m:^ID="?(.+?)"?$)`).FindStringSubmatch(string(bs))
			if len(matches) != 2 {
				currentOsId = "UNKNOWN-OS"
			} else {
				currentOsId = strings.ToLower(matches[1])
			}
		}
		log.Key(log.NewStage().Enter("init_os_id"), "current OS is "+currentOsId)
	})
	return currentOsId
}

var currentSuseRelease string
var initCurrentSuseRelease sync.Once

func GetSlesVersion() string {
	initCurrentSuseRelease.Do(func() {
		if bs, err := ioutil.ReadFile("/etc/SuSE-release"); nil != err {
			currentSuseRelease = "No-SLES"
		} else {
			currentSuseRelease = string(bs)
		}
		log.Key(log.NewStage().Enter("init_suse_version"), "current suse version is "+currentSuseRelease)
	})
	return currentSuseRelease
}

func IsSles11() bool {
	suseRelease := GetSlesVersion()
	return strings.Contains(suseRelease, "SUSE Linux Enterprise Server 11")
}

func IsSles112() bool {
	suseRelease := GetSlesVersion()
	return strings.Contains(suseRelease, "VERSION = 11\nPATCHLEVEL = 2")
}

func IsSles114() bool {
	suseRelease := GetSlesVersion()
	return strings.Contains(suseRelease, "VERSION = 11\nPATCHLEVEL = 4")
}

func IsSles() bool {
	return "sles" == GetOsId()
}

var rcRoot string
var initRcRoot sync.Once

func RcRoot() string {
	initRcRoot.Do(func() {
		rcRoot = "/etc"
		if !IsFileExist(rcRoot + "/rc0.d") {
			rcRoot = rcRoot + "/rc.d"
		}
		log.Key(log.NewStage().Enter("init_rc_root"), "current rc root is "+rcRoot)
	})
	return rcRoot
}
