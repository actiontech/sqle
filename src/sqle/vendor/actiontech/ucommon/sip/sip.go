package sip

import (
	"actiontech/ucommon/log"
	"actiontech/ucommon/os"
	"fmt"
	"regexp"
	"strings"
)

type Sip struct {
	Ip  string
	Dev string
}

var ipRegex = regexp.MustCompile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}")

//Sips could be "172.17.0.2,172.17.0.3/dev1,..."
func parseSips(sips string) ([]Sip, error) {
	sips = strings.TrimSpace(sips)
	ret := []Sip{}
	if "" == sips {
		return ret, nil
	}
	for _, a := range strings.Split(sips, ",") {
		sip := Sip{}
		if -1 != strings.Index(a, "/") {
			segs := strings.Split(a, "/")
			sip.Ip = segs[0]
			sip.Dev = segs[1]
		} else {
			sip.Ip = a
		}
		if !ipRegex.MatchString(sip.Ip) {
			return nil, fmt.Errorf("ip %v is invalid", sip.Ip)
		}
		ret = append(ret, sip)
	}
	return ret, nil
}

func IsBindSip(stage *log.Stage, sipsDesc string) (bool, error) {
	sips, err := parseSips(sipsDesc)
	if nil != err {
		return false, err
	}

	if 0 == len(sips) {
		return false, nil
	}

	for _, sip := range sips {
		if !os.IsBindSip(stage, sip.Ip) {
			return false, nil
		}
	}
	return true, nil
}

func BindSip(stage *log.Stage, sipsDesc string) error {
	sips, err := parseSips(sipsDesc)
	if nil != err {
		return err
	}

	for _, sip := range sips {
		nic := sip.Dev
		if "" == nic {
			a, err := os.GetLocalNicBySip(sip.Ip)
			if nil != err {
				return err
			}
			nic = a
		}
		os.Cmdf(stage, "SUDO ip addr add %v dev %v", sip.Ip, nic)
		os.Cmdf(stage, "SUDO arping -c 3 -A -I %v %v", nic, sip.Ip)
		if !os.IsBindSip(stage, sip.Ip) {
			return fmt.Errorf("bind sip %v failed", sip.Ip)
		}
	}
	return nil
}

func UnbindSip(stage *log.Stage, sipsDesc string) error {
	sips, err := parseSips(sipsDesc)
	if nil != err {
		return err
	}

	for _, sip := range sips {
		nic := sip.Dev
		if "" == nic {
			a, err := os.GetLocalNicBySip(sip.Ip)
			if nil != err {
				return err
			}
			nic = a
		}
		os.Cmdf(stage, "SUDO ip addr del %v/32 dev %v", sip.Ip, nic)
		if os.IsBindSip(stage, sip.Ip) {
			return fmt.Errorf("unbind sip %v failed", sip.Ip)
		}
	}

	return nil
}
