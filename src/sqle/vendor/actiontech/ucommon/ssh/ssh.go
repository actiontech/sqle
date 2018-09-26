package ssh

import (
	"actiontech/ucommon/log"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	goOs "os"
	"path/filepath"
	"strings"
)

type Ssh struct {
	user       string
	sshClient  *ssh.Client
	sftpClient *sftp.Client
}

func NewSsh(ip, port, user, passwd string) (*Ssh, error) {
	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", ip, port), &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(passwd),
		},
	})
	if nil != err {
		return nil, err
	}
	sftpClient, err := sftp.NewClient(sshClient)
	if nil != err {
		sshClient.Close()
		return nil, err
	}
	ret := Ssh{
		sshClient:  sshClient,
		sftpClient: sftpClient,
		user:       user,
	}
	return &ret, nil
}

func (s *Ssh) PushFile(srcFile, targetFile string) error {
	src, err := goOs.Open(srcFile)
	if nil != err {
		return err
	}
	defer src.Close()

	target, err := s.sftpClient.OpenFile(targetFile, goOs.O_WRONLY|goOs.O_CREATE|goOs.O_TRUNC)
	if nil != err {
		return err
	}
	defer target.Close()

	if _, err := io.Copy(target, src); nil != err {
		return err
	}
	return nil
}

func (s *Ssh) PushContent(srcContent, targetFile string) error {
	target, err := s.sftpClient.OpenFile(targetFile, goOs.O_WRONLY|goOs.O_CREATE|goOs.O_TRUNC)
	if nil != err {
		return err
	}
	defer target.Close()

	if _, err := io.Copy(target, strings.NewReader(srcContent)); nil != err {
		return err
	}
	return nil
}

func (s *Ssh) RmFile(targetFile string) error {
	return s.sftpClient.Remove(targetFile)
}

func (s *Ssh) Cmdf(cmd string, args ...interface{}) (string, int, error) {
	stage := log.NewStage().Enter("ssh_cmdf")

	cmd = fmt.Sprintf(cmd, args...)
	sudo := ""
	suroot := "bash "
	if "root" != s.user {
		sudo = "sudo -S "
		suroot = "sudo -S su root "
	}
	cmd = strings.Replace(cmd, "SUAS(", "SU $(SUDO stat -c %U ", -1)
	cmd = strings.Replace(cmd, "SUDO ", sudo, -1)
	cmd = strings.Replace(cmd, "SU ", sudo+"su ", -1)
	cmd = strings.Replace(cmd, "SUROOT ", suroot, -1)

	log.Detail(stage, cmd)

	session, err := s.sshClient.NewSession()
	if nil != err {
		return "", 0, err
	}
	defer session.Close()
	if err := session.Setenv("LC_ALL", "en_US.UTF-8"); nil != err {
		return "", 0, err
	}
	retStr, err := session.CombinedOutput(cmd)
	if nil != err {
		if e2, ok := err.(*ssh.ExitError); ok {
			return string(retStr), e2.ExitStatus(), nil
		}
		return string(retStr), 0, err
	}
	return string(retStr), 0, nil
}

func (s *Ssh) Cmdf2(cmd string, args ...interface{}) (string, error) {
	ret, retCode, err := s.Cmdf(cmd, args...)
	if nil == err && 0 != retCode {
		return ret, fmt.Errorf("%v", ret)
	}
	return ret, err
}

func (s *Ssh) Destroy() {
	s.sftpClient.Close()
	s.sshClient.Close()
}

func (s *Ssh) EnsureNewDirAndGroupAccess(dir string, userGroup string) (bool, error) {
	dir = filepath.Clean(dir)
	if "/" == dir {
		return false, nil
	}
	_, retCode, err := s.Cmdf("SUDO ls " + dir)
	if nil != err {
		return false, err
	}
	if 0 == retCode {
		return false, fmt.Errorf("Dir %v already exist", dir)
	}
	return s.ensureGroupAccess(dir, userGroup)
}

func (s *Ssh) ensureGroupAccess(dir, userGroup string) (bool, error) {
	if _, retCode, err := s.Cmdf("SUDO ls " + dir); nil != err {
		return false, err
	} else if 0 == retCode {
		return false, nil
	}
	if _, err := s.ensureGroupAccess(filepath.Dir(dir), userGroup); nil != err {
		return false, err
	}
	if _, err := s.Cmdf2("SUDO mkdir %s && SUDO chgrp %s %s && SUDO chmod 750 %s", dir, userGroup, dir, dir); nil != err {
		return false, err
	}
	return true, nil
}
