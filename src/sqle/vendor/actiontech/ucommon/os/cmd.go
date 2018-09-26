package os

import (
	"actiontech/ucommon/log"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

//Max cmd output is limited to 10MB
type HaCmdOutputBuffer struct {
	buf   bytes.Buffer
	mutex sync.Mutex
}

func (h *HaCmdOutputBuffer) Bytes() []byte {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	return h.buf.Bytes()
}

func (h *HaCmdOutputBuffer) Write(p []byte) (n int, err error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if h.buf.Len() < 1024*1024*10 {
		return h.buf.Write(p)
	}
	return len(p), nil
}

type HaCmd struct {
	cmd     *exec.Cmd
	deferFn func()
	str     string
	buf     HaCmdOutputBuffer
	mutex   sync.Mutex
}

func (h *HaCmd) Kill() error {
	var process *os.Process

	h.mutex.Lock()
	if nil != h.cmd && nil != h.cmd.Process {
		process = h.cmd.Process
	}
	h.mutex.Unlock()

	if nil == process {
		return fmt.Errorf("cmd.Process is nil, cannot kill")
	}
	err := KillWithAllDescendantPids(process.Pid)
	if nil == err {
		return nil
	} else {
		return fmt.Errorf("kill cmd tree error: %v", err)
	}
}

func (h *HaCmd) Pid() (int, error) {
	var process *os.Process

	h.mutex.Lock()
	if nil != h.cmd && nil != h.cmd.Process {
		process = h.cmd.Process
	}
	h.mutex.Unlock()

	if nil == process {
		return -1, fmt.Errorf("cmd.Process is nil, cannot kill")
	}
	return process.Pid, nil
}

func (h *HaCmd) Destroy() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.deferFn()
}

//retStr and retCode works when nil != err
func (h *HaCmd) Start(stage *log.Stage) (retStr string, retCode int, err error) {
	return h.start(stage)
}

//retStr and retCode works when nil != err
func (h *HaCmd) start(stage *log.Stage) (retStr string, retCode int, err error) {
	stage.Enter("Cmd_start")
	defer stage.Exit()

	h.mutex.Lock()
	defer h.mutex.Unlock()

	log.Detail(stage, "{%v}", h.str)

	if h.cmd.Stdout != nil {
		return h.handleError(stage, errors.New("exec: Stdout already set"))
	} else if h.cmd.Stderr != nil {
		return h.handleError(stage, errors.New("exec: Stderr already set"))
	}
	h.cmd.Stdout = &h.buf
	h.cmd.Stderr = &h.buf
	retStr, retCode, err = h.handleError(stage, h.cmd.Start())
	if nil != err || 0 != retCode {
		log.Write(stage).Brief("{%v},%v,%v", h.str, retCode, err).
			Detail("output=%v", retStr).Done()
	}
	return retStr, retCode, err
}

func (h *HaCmd) Wait(stage *log.Stage) (retStr string, retCode int, err error) {
	return h.wait(stage)
}

func (h *HaCmd) wait(stage *log.Stage) (retStr string, retCode int, err error) {
	stage.Enter("Cmd_wait")
	defer stage.Exit()

	retStr, retCode, err = h.handleError(stage, h.cmd.Wait())
	if nil == err && 0 == retCode {
		log.Write(stage).Detail("{%v},%v,%v,output=%v", h.str, retCode, err, retStr).Done()
	} else {
		log.Write(stage).Brief("{%v},%v,%v", h.str, retCode, err).
			Detail("output=%v", retStr).Done()
	}
	return retStr, retCode, err
}

func (h *HaCmd) Run(stage *log.Stage) (output string, retCode int, err error) {
	if output, retCode, err = h.start(stage); nil == err {
		output, retCode, err = h.wait(stage)
	}
	return
}

func (h *HaCmd) handleError(stage *log.Stage, e error) (retStr string, retCode int, err error) {
	retStr = string(h.buf.Bytes())
	retStr = strings.TrimSpace(retStr)
	err = e
	if nil != err {
		if e2, ok := err.(*exec.ExitError); ok {
			if s, ok := e2.Sys().(syscall.WaitStatus); ok {
				return retStr, int(s.ExitStatus()), nil
			}
		}
		return "", 0, err
	}
	return retStr, 0, nil
}

func Cmdf(stage *log.Stage, base string, args ...interface{}) (output string, retCode int, err error) {
	stage.Enter("Cmd")
	defer stage.Exit()

	return innerCmdfTimeout(stage, base, 0, args...)
}

func Cmdf2(stage *log.Stage, base string, args ...interface{}) (output string, err error) {
	stage.Enter("Cmd2")
	defer stage.Exit()

	ret, retCode, err := innerCmdfTimeout(stage, base, 0, args...)
	if nil != err {
		return ret, err
	}
	if 0 != retCode {
		return ret, fmt.Errorf("retCode=%v (%v)", retCode, ret)
	}
	return ret, nil
}

func CmdfWithQuitChan(stage *log.Stage, base string, quitChan chan bool, args ...interface{}) (output string, retCode int, err error) {
	stage.Enter("Cmd_with_quit_chan")
	defer stage.Exit()

	return innerCmdfWithQuitChan(stage, base, quitChan, args...)
}

func innerCmdfWithQuitChan(stage *log.Stage, base string, quitChan chan bool, args ...interface{}) (string, int, error) {
	str := fmt.Sprintf(base, args...)

	cmd, err := NewHaCmd(stage, str)
	defer cmd.Destroy()
	if nil != err {
		return "", 1, err
	}
	if nil == quitChan {
		return cmd.Run(stage)
	}
	okChan := make(chan bool, 1)

	var output string
	var retCode int

	go func(stage *log.Stage) {
		output, retCode, err = cmd.Run(stage)
		okChan <- true
	}(stage.Go())
	select {
	case <-okChan:
		return output, retCode, err
	case <-quitChan:
		if err := cmd.Kill(); nil == err {
			log.Brief(stage, "{%v} timeout, killed", str)
		} else {
			log.Brief(stage, "{%v} timeout, kill error (%v)", str, err)
		}
		return "", 1, fmt.Errorf("timeout")
	}
}

func CmdfTimeout(stage *log.Stage, base string, timeoutSeconds int, args ...interface{}) (output string, retCode int, err error) {
	stage.Enter(fmt.Sprintf("Cmd_timeout_%v", timeoutSeconds))
	defer stage.Exit()

	return innerCmdfTimeout(stage, base, timeoutSeconds, args...)
}

func innerCmdfTimeout(stage *log.Stage, base string, timeoutSeconds int, args ...interface{}) (output string, retCode int, err error) {
	if 0 == timeoutSeconds {
		return innerCmdfWithQuitChan(stage, base, nil, args...)
	}
	deferChan := make(chan bool, 1)
	quitChan := make(chan bool, 1)
	go func() {
		select {
		case <-time.After(time.Duration(timeoutSeconds) * time.Second):
			quitChan <- true
		case <-deferChan:
			quitChan <- true
		}
		close(quitChan)
	}()
	defer func() {
		deferChan <- true
		close(deferChan)
	}()
	return innerCmdfWithQuitChan(stage, base, quitChan, args...)
}

func NewHaCmd(stage *log.Stage, str string) (a *HaCmd, err error) {
	a = &HaCmd{}
	a.deferFn = func() {
	}
	sudo := ""
	suroot := "bash "
	if 0 != os.Getuid() {
		sudo = "sudo -S "
		suroot = "sudo -S su -s $(which bash) root "
	}
	str = strings.Replace(str, "SUAS(", "SU $(SUDO stat -c %U ", -1)
	str = strings.Replace(str, "SUDO ", sudo, -1)
	str = strings.Replace(str, "SU ", sudo+"su -s $(which bash) ", -1)
	str = strings.Replace(str, "SUROOT ", suroot, -1)
	a.str = str
	a.cmd = exec.Command("bash", "--noprofile", "--norc", "-c", fmt.Sprintf("%v", str))
	return a, nil
}
