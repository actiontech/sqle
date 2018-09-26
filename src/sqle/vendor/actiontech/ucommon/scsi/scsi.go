package scsi

import (
	"actiontech/ucommon/log"
	"actiontech/ucommon/os"
	"actiontech/ucommon/util"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	maxCoolDownTimes  = 1 //for NHA, 2pc is not required (no brain-split), set maxCoolDownTimes=1 to avoid long wait in 2pc
	registerKeyReg    = regexp.MustCompile("\\s+([0-9a-fxA-FX]*)$")
	ErrFenced         = errors.New("Local node is fenced")
	FenceHexId        = "0x6841"
	ErrNeedPreempt    = errors.New("need preempt")
	reservationKeyReg = regexp.MustCompile("(?s)Key\\s*=\\s*(0x[0-9a-fA-F]+)\\s.*type\\s*[:=]\\s*(.*)")
)

const (
	MPATH_PERSIST = "mpathpersist"
	SG_PERSIST    = "sg_persist"
)

func GetPersistTool(stage *log.Stage, dev string) string {
	devLastPart := filepath.Base(dev)
	if os.IsFileExist("./force_use_sg_persist_" + devLastPart) {
		return SG_PERSIST
	}
	if os.IsFileExist("./force_use_mpath_persist_" + devLastPart) {
		return MPATH_PERSIST
	}
	_, err := os.Cmdf2(stage, "which multipath")
	if nil != err {
		return SG_PERSIST
	}
	ret, err := os.Cmdf2(stage, "SUDO multipath -l %s", dev)
	if nil != err || "" == ret {
		return SG_PERSIST
	} else {
		return MPATH_PERSIST
	}
}

func GetReservationKey(stage *log.Stage, dev string) (string, error) {
	key, _, err := getReservationInfo(stage, dev, GetPersistTool(stage, dev))
	return key, err
}

func getReservationInfo(stage *log.Stage, dev string, persistTool string) (string, string, error) {
	stage.Enter("Get_Reserve_Key")
	defer stage.Exit()

	command := ""
	switch persistTool {
	case SG_PERSIST:
		command = fmt.Sprintf("SUDO sg_persist -n -i -r -d %v", dev)
	case MPATH_PERSIST:
		command = fmt.Sprintf("SUDO mpathpersist -i -r -d %v", dev)
	}

	if ret, retCode, err := cmd(stage, command); nil != err || 0 != retCode {
		return "", "", fmt.Errorf("PR get reservation key failed")
	} else {
		matches := reservationKeyReg.FindStringSubmatch(ret)
		if len(matches) > 2 {
			rsvKey := matches[1]
			type_ := strings.TrimSpace(matches[2])
			log.Brief(stage, "reservation=%v type=%v", rsvKey, type_)
			return rsvKey, type_, nil
		}
		log.Brief(stage, "reservation= type=")
		return "", "", nil
	}

}

func Register(stage *log.Stage, dev string, hexId string) error {
	return register(stage, dev, hexId, GetPersistTool(stage, dev))
}

func register(stage *log.Stage, dev string, hexId string, persistTool string) error {
	stage.Enter("Scsi_register")
	defer stage.Exit()

	command := GetRegisterCmd(stage, hexId, dev)
	_, retCode, err := cmd(stage, command)
	if nil != err || 0 != retCode {
		return fmt.Errorf("PR register failed")
	}
	log.Key(stage, "registered dev(%v) key(%v)", dev, hexId)
	return nil
}

func Clear(stage *log.Stage, dev string, hexId string) error {
	return clear(stage, dev, hexId, GetPersistTool(stage, dev))
}

func clear(stage *log.Stage, dev string, hexId string, persistTool string) error {
	stage.Enter("Scsi_clear")
	defer stage.Exit()

	if err := register(stage, dev, hexId, persistTool); nil != err {
		return err
	}

	command := GetClearScsiCmd(stage, hexId, dev)
	_, retCode, err := cmd(stage, command)
	if nil != err || 0 != retCode {
		return fmt.Errorf("PR clear failed")
	}
	log.Key(stage, "cleared dev(%v) key(%v)", dev, hexId)
	return nil
}

func GetClearScsiCmd(stage *log.Stage, hexId string, dev string) string {
	persistTool := GetPersistTool(stage, dev)

	command := ""
	switch persistTool {
	case SG_PERSIST:
		command = fmt.Sprintf("SUDO sg_persist -n -o -C -K %v -d %v", hexId, dev)
	case MPATH_PERSIST:
		command = fmt.Sprintf("SUDO mpathpersist -o -C -K %v -d %v", hexId, dev)
	}
	return command
}

func GetRegisterCmd(stage *log.Stage, hexId string, dev string) string {
	persistTool := GetPersistTool(stage, dev)

	command := " "
	switch persistTool {
	case SG_PERSIST:
		command = fmt.Sprintf("SUDO sg_persist -n -o -I -S %v -d %v", hexId, dev)
	case MPATH_PERSIST:
		command = fmt.Sprintf("SUDO mpathpersist -o -I -S %v -d %v", hexId, dev)
	}
	return command
}

func Reserve(stage *log.Stage, dev string, hexId string, level int, checkFenceHexId bool) error {
	return reserve(stage, dev, hexId, level, checkFenceHexId, GetPersistTool(stage, dev))
}

func reserve(stage *log.Stage, dev string, hexId string, level int, checkFenceHexId bool, persistTool string) error {
	stage.Enter("Scsi_reserve")
	defer stage.Exit()

	command := ""
	switch persistTool {
	case SG_PERSIST:
		command = fmt.Sprintf("SUDO sg_persist -n -o -R -S 0x0 -T %v -K %v -d %v", level, hexId, dev)
	case MPATH_PERSIST:
		command = fmt.Sprintf("SUDO mpathpersist -o -R -T %v -K %v -d %v", level, hexId, dev)
	}

	logErr := func(err error) error {
		if nil == err {
			log.Key(stage, "succeed")
		} else {
			log.KeyDilute1(stage, dev+err.Error(), "return error (%v)", err)
		}
		return err
	}
retry:
	if err := register(stage, dev, hexId, persistTool); nil != err {
		return logErr(err)
	}
	_, retCode, err := cmd(stage, command)
	if nil != err {
		return logErr(fmt.Errorf("PR reserve error"))
	}
	reservationKey, type_, err := getReservationInfo(stage, dev, persistTool)
	if nil != err {
		return logErr(err)
	}
	if strings.HasPrefix(type_, "obsolete") {
		return logErr(ErrNeedPreempt)
	}
	if FenceHexId != hexId && FenceHexId == reservationKey {
		clear(stage, dev, FenceHexId, persistTool)
		if checkFenceHexId {
			return logErr(ErrFenced)
		} else {
			checkFenceHexId = true
			goto retry
		}
	}
	if hexId != reservationKey {
		return logErr(fmt.Errorf("PR reserve failed"))
	}
	if hexId == reservationKey && 0 != retCode {
		return logErr(ErrNeedPreempt)
	}
	if 0 != retCode {
		return logErr(fmt.Errorf("PR reserve failed, retCode=%v", retCode))
	}
	log.Key(stage, "reserved dev(%v) key(%v)", dev, hexId)
	return log.KeyRet(stage, nil)
}

func Release(stage *log.Stage, dev string, hexId string, level int) error {
	return release(stage, dev, hexId, level, GetPersistTool(stage, dev))
}

func release(stage *log.Stage, dev string, hexId string, level int, persistTool string) error {
	stage.Enter("Scsi_release")
	defer stage.Exit()

	if err := register(stage, dev, hexId, persistTool); nil != err {
		return err
	}
	if err := util.DebugError(stage, "scsi release fail"); nil != err {
		return err
	}

	command := ""
	switch persistTool {
	case SG_PERSIST:
		command = fmt.Sprintf("SUDO sg_persist -n -o -L -S 0x0 -T %v -K %v -d %v", level, hexId, dev)
	case MPATH_PERSIST:
		command = fmt.Sprintf("SUDO mpathpersist -o -L -T %v -K %v -d %v", level, hexId, dev)
	}

	if _, retCode, err := cmd(stage, command); nil != err || 0 != retCode {
		return fmt.Errorf("PR release failed")
	}
	reservationKey, _, err := getReservationInfo(stage, dev, persistTool)
	if nil != err {
		return err
	}
	if reservationKey == hexId {
		return fmt.Errorf("PR release failed")
	}
	log.Key(stage, "released dev(%v) key(%v)", dev, hexId)
	return nil
}

func Preempt(stage *log.Stage, dev string, hexId string, level int, checkFenceHexId bool) error {
	return preempt(stage, dev, hexId, level, checkFenceHexId, GetPersistTool(stage, dev))
}

func preempt(stage *log.Stage, dev string, hexId string, level int, checkFenceHexId bool, persistTool string) error {
	stage.Enter("Scsi_preempt")
	defer stage.Exit()

retry:
	oldReservationKey, type_, err := getReservationInfo(stage, dev, persistTool)
	if nil != err {
		return err
	}
	if "" == oldReservationKey {
		return reserve(stage, dev, hexId, level, checkFenceHexId, persistTool)
	}
	if strings.HasPrefix(type_, "obsolete") {
		//register by old key
		register(stage, dev, oldReservationKey, persistTool)
		//preempt by old key
		{
			command := ""
			switch persistTool {
			case SG_PERSIST:
				command = fmt.Sprintf("SUDO sg_persist -n -o -P -T %v -K %v -S %v -d %v", level, oldReservationKey, oldReservationKey, dev)
			case MPATH_PERSIST:
				command = fmt.Sprintf("SUDO mpathpersist -o -P -T %v -K %v -S %v -d %v", level, oldReservationKey, oldReservationKey, dev)
			}
			cmd(stage, command)
		}
		//log
		getReservationInfo(stage, dev, persistTool)
		//release old key
		{
			command := ""
			switch persistTool {
			case SG_PERSIST:
				command = fmt.Sprintf("SUDO sg_persist -n -o -L -S 0x0 -T %v -K %v -d %v", level, oldReservationKey, dev)
			case MPATH_PERSIST:
				command = fmt.Sprintf("SUDO mpathpersist -n -o -L -T %v -K %v -d %v", level, oldReservationKey, dev)
			}
			cmd(stage, command)
		}
		//unregister
		unregister(stage, dev, oldReservationKey, persistTool)
		_, type_, err := getReservationInfo(stage, dev, persistTool)
		if nil != err {
			return err
		}
		if strings.HasPrefix(type_, "obsolete") {
			return fmt.Errorf("fail to preempt obsolete key")
		}
		goto retry
	}
	if err := register(stage, dev, hexId, persistTool); nil != err {
		return err
	}
	if err := util.DebugError(stage, "scsi preempt fail"); nil != err {
		return err
	}
	{
		command := ""
		switch persistTool {
		case SG_PERSIST:
			command = fmt.Sprintf("SUDO sg_persist -n -o -P -T %v -K %v -S %v -d %v", level, hexId, oldReservationKey, dev)
		case MPATH_PERSIST:
			command = fmt.Sprintf("SUDO mpathpersist -o -P -T %v -K %v -S %v -d %v", level, hexId, oldReservationKey, dev)
		}
		if _, _, err := cmd(stage, command); nil != err {
			return fmt.Errorf("PR preempt error")
		}
	}
	reservationKey, type_, err := getReservationInfo(stage, dev, persistTool)
	if nil != err {
		return err
	}
	if reservationKey != hexId {
		if FenceHexId == reservationKey {
			clear(stage, dev, FenceHexId, persistTool)
			if checkFenceHexId {
				return ErrFenced
			} else {
				checkFenceHexId = true
				goto retry
			}
		}
		return fmt.Errorf("PR preempt failed")
	}

	log.Key(stage, "preempted dev(%v) key(%v)", dev, hexId)
	return nil
}

func ReserveOrTimeoutThenPreempt(stage *log.Stage, dev string, hexId string, level int, checkFenceHexId bool) (doneChan chan error, timeoutChan chan bool) {
	return reserveOrTimeoutThenPreempt(stage, dev, hexId, level, checkFenceHexId, GetPersistTool(stage, dev))
}

func ReserveOrPreempt(stage *log.Stage, dev string, hexId string, level int, checkFenceHexId bool) (doneChan chan error) {
	dc, tc := reserveOrTimeoutThenPreempt(stage, dev, hexId, level, checkFenceHexId, GetPersistTool(stage, dev))
	close(tc)
	return dc
}

func reserveOrTimeoutThenPreempt(stage *log.Stage, dev string, hexId string, level int, checkFenceHexId bool, persistTool string) (doneChan chan error, timeoutChan chan bool) {
	stage.Enter("Scsi_reserve_or_preempt")
	defer stage.Exit()

	log.Key(stage, "Reserve or preempt dev(%v) key(%v)...", dev, hexId)

	doneChan = make(chan error, 1)
	timeoutChan = make(chan bool, 1)
	go func(stage *log.Stage) {
		for {
			err := reserve(stage, dev, hexId, level, checkFenceHexId, persistTool)
			if nil == err {
				doneChan <- nil
				return
			} else if ErrNeedPreempt == err {
				doneChan <- preempt(stage, dev, hexId, level, checkFenceHexId, persistTool)
				return
			} else if ErrFenced == err {
				doneChan <- err
				return
			}
			select {
			case <-timeoutChan:
				doneChan <- preempt(stage, dev, hexId, level, checkFenceHexId, persistTool)
				return
			case <-time.After(100 * time.Millisecond):
			}
		}
	}(stage.Go())
	return doneChan, timeoutChan
}

func cmd(stage *log.Stage, c string) (ret string, retCode int, err error) {
	//sg_persist may get "PR in: unit attention", need retry
	for i := 0; i < 3; i++ {
		ret, retCode, err = os.Cmdf(stage, c)
		if nil != err {
			return
		}
		if 0 == retCode {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	return
}

func parsingRegKeyExp(exp string) []string {
	keys := make([]string, 0)
	for _, line := range strings.Split(exp, "\n") {
		if !registerKeyReg.MatchString(line) {
			continue
		}
		matches := registerKeyReg.FindStringSubmatch(line)
		keys = append(keys, matches[1])
	}
	return keys
}

func GetRegsiterKeys(stage *log.Stage, dev string) ([]string, error) {
	return getRegsiterKeys(stage, dev, GetPersistTool(stage, dev))
}

func getRegsiterKeys(stage *log.Stage, dev string, persistTool string) ([]string, error) {
	stage.Enter("Get_Regsiter_Key")
	defer stage.Exit()

	command := ""
	switch persistTool {
	case SG_PERSIST:
		command = fmt.Sprintf("SUDO sg_persist -n -i -k -d %v", dev)
	case MPATH_PERSIST:
		command = fmt.Sprintf("SUDO mpathpersist -i -k -d %v", dev)
	}

	if ret, retCode, err := cmd(stage, command); nil != err || 0 != retCode {
		return nil, fmt.Errorf("PR get regsiter key failed")
	} else {
		ret := parsingRegKeyExp(ret)
		log.Brief(stage, "get register keys(%v) of dev(%v)", ret, dev)
		return ret, nil
	}
}

func Unregister(stage *log.Stage, dev string, hexId string) error {
	return unregister(stage, dev, hexId, GetPersistTool(stage, dev))
}

func unregister(stage *log.Stage, dev string, hexId string, persistTool string) error {
	stage.Enter("Unregister")
	defer stage.Exit()

	if err := util.DebugError(stage, "scsi unregister fail"); nil != err {
		return err
	}

	command := ""
	switch persistTool {
	case SG_PERSIST:
		command = fmt.Sprintf("SUDO sg_persist -n -o -G -S 0x0 -K %v -d %v", hexId, dev)
	case MPATH_PERSIST:
		command = fmt.Sprintf("SUDO mpathpersist -o -G -S 0x0 -K %v -d %v", hexId, dev)
	}

	if _, retCode, err := cmd(stage, command); nil != err || 0 != retCode {
		return fmt.Errorf("PR unregister failed")
	}
	registerKeys, err := getRegsiterKeys(stage, dev, persistTool)
	if nil != err {
		return err
	}
	for _, registerKey := range registerKeys {
		if registerKey == hexId {
			return fmt.Errorf("PR unregister failed")
		}
	}

	log.Key(stage, "unregisted dev(%v) key(%v)", dev, hexId)
	return nil
}

//2pc
func PreemptPrepare(stage *log.Stage, localDev, remoteDev, localHexId, remoteHexId string, level int) error {
	stage.Enter("2pc_prepare")
	defer stage.Exit()

	localPersistTool := GetPersistTool(stage, localDev)
	remotePersistTool := GetPersistTool(stage, localDev)

	coolDownTimes := 0
START:
	log.Key(stage, "scsi 2pc preparing, round (%v)...", coolDownTimes)
	coolDownTimes++
	//treat remote node is crash after max cool down times -> force fence
	if coolDownTimes > maxCoolDownTimes {
		log.Brief(stage, "check local dev register key beyond max times,force commit")
		return nil
	}
	//register remote dev
	//need unregister before every return
	if err := util.DebugError(stage, "scsi 2pc prepare register fail"); nil != err {
		return err
	}

	if err := register(stage, remoteDev, localHexId, remotePersistTool); nil != err {
		log.Brief(stage, "register remote dev fail")
		return err
	}
	log.Brief(stage, "register remote dev succeed")

	util.DebugPause("pause after scsi 2pc register")

	// check local dev regsiter key
	registerKeys, err := getRegsiterKeys(stage, localDev, localPersistTool)
	if nil != err {
		log.Brief(stage, "get remote dev register key fail, unregister")
		unregister(stage, remoteDev, localHexId, remotePersistTool)
		return err
	}
	for _, registerKey := range registerKeys {
		if remoteHexId == registerKey {
			//check local dev reservation key
			err := checkReservationKey(stage, localHexId, localDev, localPersistTool)
			if nil != err {
				log.Brief(stage, "check local dev reservation key fail, unregister")
				unregister(stage, remoteDev, localHexId, remotePersistTool)
				return err
			}

			//cool down
			unregister(stage, remoteDev, localHexId, remotePersistTool)
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
			log.Brief(stage, "check local dev register key fail, enter next round")
			goto START
		}
	}
	log.Brief(stage, "check local dev register key succeed")

	//check local dev reservation key
	err = checkReservationKey(stage, localHexId, localDev, localPersistTool)
	if nil != err {
		log.Brief(stage, "check local dev reservation key fail, unregister")
		unregister(stage, remoteDev, localHexId, remotePersistTool)
		return err
	}
	log.Brief(stage, "check local dev reservation key succeed")
	return nil
}

func PreemptCommit(stage *log.Stage, dev string, hexId string, level int, checkFenceHexId bool) error {
	stage.Enter("2pc_commit")
	defer stage.Exit()
	util.DebugPause("pause before scsi 2pc commit")
	return log.KeyRet(stage, preempt(stage, dev, hexId, level, checkFenceHexId, GetPersistTool(stage, dev)))
}

func checkReservationKey(stage *log.Stage, hexId, dev string, persistTool string) error {
	rsvKey, _, err := getReservationInfo(stage, dev, persistTool)
	if nil != err {
		return err
	}
	if hexId != rsvKey {
		return ErrFenced
	}
	return nil
}

type DevNotWritableError struct {
	mountPoint string
}

func (d *DevNotWritableError) Error() string {
	return fmt.Sprintf("scsi write point %v cannot be write", d.mountPoint)
}

func CheckDevWritable(stage *log.Stage, dev, hexId, mountPoint string) error {
	key, _, err := getReservationInfo(stage, dev, GetPersistTool(stage, dev))
	if nil != err {
		return err
	}
	if key != hexId {
		return fmt.Errorf("scsi reservation key changed to %v", key)
	}

	_, err = ioutil.ReadDir(mountPoint)
	if nil != err {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(mountPoint, "U_TEST_SCSI_WRITABLE"), []byte("This file is created by Actiontech Universe"), 0700); nil != err {
		if strings.Contains(err.Error(), "no space left") {
			//ignore disk-space-full error
			return nil
		}
		return &DevNotWritableError{mountPoint}
	}

	return nil
}
