package ubootstrap

import (
	"actiontech/ucommon/log"
	"actiontech/ucommon/os"
	user_ "actiontech/ucommon/user"
	"actiontech/ucommon/util"
	"fmt"
	"io/ioutil"
	os_ "os"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"syscall"
	"os/exec"
	"golang.org/x/net/trace"
	"net/http"
)

func init() {
	//bind init thread, otherwise GetCap will throw "no such process", since golang 1.7
	runtime.LockOSThread()

	if err := os_.Setenv("PATH", os_.Getenv("PATH")+":/sbin:/usr/sbin"); nil != err {
		panic(err.Error())
	}
	if err := os_.Setenv("LC_ALL", "en_US.UTF-8"); nil != err {
		panic(err.Error())
	}

	os.GetRootDir() //init root dir, keep it in OS main thread, otherwise maybe "permission denied"
}

func DefaultSocket() string {
	defaultSocket, err := filepath.Abs("./socket")
	if nil != err {
		panic("filepath.Abs panic")
	}
	return defaultSocket
}

func ListenKillSignal() chan os_.Signal {
	quitChan := make(chan os_.Signal, 1)
	signal.Notify(quitChan, os_.Interrupt, os_.Kill, syscall.SIGTERM, syscall.SIGUSR2 /*graceful-shutdown*/)
	return quitChan
}

func ChangeRunUser(runUser string, background bool) error {
	runUid, runGids, err := user_.LookupUidGidByUser(runUser)
	if nil != err {
		return err
	}
	
	if runUid != os_.Getuid() || runGids[0] != os_.Getgid() {
		if err := os.SetKeepCaps(); nil != err {
			return err
		}
		cmd := exec.Command(os_.Args[0], os_.Args[1:]...)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Credential: &syscall.Credential{
				Uid: uint32(runUid),
				Gid: uint32(runGids[0]),
			},
			Setsid: true,
		}
		cmd.Stdout = os_.Stdout
		cmd.Stderr = os_.Stderr
		if background {
			if err := cmd.Start(); nil != err {
				return err
			}
			fmt.Printf("Forked to process %d, make uid and cap work. Check syslog if failed.\n", cmd.Process.Pid)
			cmd.Process.Release()
		} else {
			if err := cmd.Run(); nil != err {
				if e2, ok := err.(*exec.ExitError); ok {
					if s, ok := e2.Sys().(syscall.WaitStatus); ok {
						os_.Exit(int(s.ExitStatus()))
					}
				}
				os_.Exit(1)
			}
		}
		os_.Exit(0)
	}

	{
		//mystery: if miss write() syscall, capset will throw EPERM.
		//totally make no sence
		fmt.Printf("")
	}

	capBeforeSetuid, err := os.GetCap()
	if nil != err {
		return fmt.Errorf("capget error: %v", err)
	}

	{
		//mystery: if miss write() syscall, capset will throw EPERM.
		//totally make no sence
		fmt.Printf("")
	}

	newCap := os.UserCapDataStruct{
		Effective:   capBeforeSetuid.Effective,
		Permitted:   capBeforeSetuid.Effective,
		Inheritable: capBeforeSetuid.Effective,
	}
	if err := os.SetCap(newCap); nil != err {
		return fmt.Errorf("capset error: %v", err)
	}

	return nil
}

func StartPid(pidFile string) error {
	stage := log.NewStage().Enter("start_pid")

	// lock file
	f, err := os_.OpenFile(pidFile, os_.O_CREATE, 0640)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		return fmt.Errorf("service has been starting...")
	}
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	util.DebugPause("wait another process")

	if bs, err := ioutil.ReadFile(pidFile); nil == err {
		pidString := strings.TrimSpace(string(bs))
		if pid, err := strconv.Atoi(pidString); nil == err {
			if os.IsProcessExist(stage, pid) && os.IsProcessMatchPidFile(pid, pidFile) && pid != os_.Getpid() {
				return fmt.Errorf("service is already running, pid is %v", pid)
			}
		}
	}
	if err := ioutil.WriteFile(pidFile, []byte(fmt.Sprintf("%v", os_.Getpid())), 0640); nil != err {
		return err
	}
	return nil
}

func StopPid(pidFile string) error {
	stage := log.NewStage().Enter("stop_pid")
	if pid, err := ioutil.ReadFile(pidFile); nil == err {
		pidString := strings.TrimSpace(string(pid))
		if pid, err := strconv.Atoi(pidString); nil == err {
			if pid == os_.Getpid() {
				return os.Remove(stage, pidFile)
			}
		}
	}
	return nil
}

func DumpLoop() {
	c := make(chan os_.Signal, 10)
	signal.Notify(c, syscall.Signal(0x15)) //0x15=SIGTTIN

	for {
		sig := <-c
		switch sig {
		case syscall.Signal(0x15):
			go func() {
				pprof.Lookup("goroutine").WriteTo(os_.Stdout, 1)
				if f, err := os_.OpenFile("dump", os_.O_WRONLY|os_.O_TRUNC|os_.O_CREATE, 0640); nil != err {
					fmt.Fprintf(os_.Stderr, "write dump error(%v)", err)
				} else {
					pprof.Lookup("goroutine").WriteTo(f, 1)
					f.Close()
				}
				{
					f, err := os_.OpenFile("heap_dump", os_.O_WRONLY|os_.O_TRUNC|os_.O_CREATE, 0640)
					if nil != err {
						fmt.Fprintf(os_.Stderr, "write heap_dump error(%v)", err)
					}
					pprof.WriteHeapProfile(f)
					f.Close()
				}
				return
			}()
		default:
		}
	}
}

func GrpcTraceHttpService(grpcTracePort int) {
	stage := log.NewStage().Enter("grpc_trace")
	if grpcTracePort > 0 {
		trace.AuthRequest = func(req *http.Request) (any, sensitive bool) {
			return true, true
		}
		log.Key(stage, "grpc trace listen on %v", grpcTracePort)
		if err := http.ListenAndServe(fmt.Sprintf(":%v", grpcTracePort), nil); nil != err {
			log.Key(stage, "grpc trace error: %v", err)
		}
	}
}

func SetLimitNOFILE(limit uint64) error {
	if err := setLimit(syscall.RLIMIT_NOFILE, limit); nil != err {
		return fmt.Errorf("set ulimit NOFILE %d error: %s\n", limit, err.Error())
	}
	return nil
}

func SetLimitNPROC(limit uint64) error {
	if err := setLimit(0x6 /* RLIMIT_NPROC */, limit); nil != err {
		return fmt.Errorf("set ulimit NPROC %d error: %s\n", limit, err.Error())
	}
	return nil
}

func setLimit(resource int, limit uint64) error {
	if limit <= 0 {
		return nil
	}
	var rlimit syscall.Rlimit
	err := syscall.Getrlimit(resource, &rlimit)
	if nil != err {
		return err
	}
	if limit == rlimit.Max {
		return nil
	}
	rlimit.Max = limit
	rlimit.Cur = limit
	return syscall.Setrlimit(resource, &rlimit)
}
