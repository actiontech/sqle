package test

import (
	"actiontech/ucommon/log"
	"actiontech/ucommon/os"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var compiled map[string]bool = make(map[string]bool)

func DockerTest(t *testing.T, dockerImage string) bool {
	return innerDockerTest(t, dockerImage, "root", "")
}

func DockerTestWithUser(t *testing.T, dockerImage string, runUser string) bool {
	return innerDockerTest(t, dockerImage, runUser, "")
}

func DockerTestWithUserWithCap(t *testing.T, dockerImage string, runUser string, caps string) bool {
	return innerDockerTest(t, dockerImage, runUser, caps)
}

func InDocker() bool {
	return os.IsFileExist("/tmp/THIS_IS_IN_CONTAINER")
}

/*
	When set spec user, also set umask 0077
*/
func innerDockerTest(t *testing.T, dockerImage string, runUser string, caps string) bool {
	stage := log.NewStage().Enter("in_docker_test")
	if os.IsFileExist("/tmp/THIS_IS_IN_CONTAINER") {
		return true
	}
	callerFn := ""
	{
		pc, _, _, ok := runtime.Caller(2)
		details := runtime.FuncForPC(pc)
		if !ok || nil == details {
			t.Fatalf("unable to recognize caller: %v", details.Name())
		}
		arr := strings.Split(details.Name(), ".")
		callerFn = arr[len(arr)-1]
	}

	goPath, _ := filepath.Abs(".")
	relPath := ""
	{
		for filepath.Base(goPath) != "src" {
			relPath = filepath.Base(goPath) + "/" + relPath
			goPath = filepath.Dir(goPath)
		}
		relPath = filepath.Base(goPath) + "/" + relPath
		relPath = filepath.Clean(relPath)
		goPath = filepath.Dir(goPath)
	}

	compiledFileName := ""
	{
		cwd, _ := filepath.Abs(".")
		compiledFileName = filepath.Base(cwd) + ".test"
	}

	if !compiled[compiledFileName] {
		_, err := os.Cmdf2(stage, "GOPATH=%v GOOS=linux go test -c", goPath)
		if nil != err {
			t.Fatalf("compile error: %v", err)
		}
		compiled[compiledFileName] = true
	}

	inShellCmd := fmt.Sprintf("touch /tmp/THIS_IS_IN_CONTAINER; mkdir /tmp/test_root; cp /testing/%v/%v /tmp/test_root/%v", relPath, compiledFileName, compiledFileName)
	if "" != caps {
		inShellCmd = inShellCmd + fmt.Sprintf(";sudo setcap %v=+eip /tmp/test_root/%v", caps, compiledFileName)
	}
	inShellCmd = fmt.Sprintf("%v; cd /tmp/test_root; PATH=/usr/local/go/bin/:$PATH ./%v -test.run \"^%v\\$\" -test.v", inShellCmd, compiledFileName, callerFn)
	shellCmd := fmt.Sprintf("sh -c '%s'", inShellCmd)
	if runUser != "root" {
		//group 999 is for virtualbox mapped folder
		shellCmd = fmt.Sprintf("sh -c \"(groupadd -g 999 test_group) && (useradd -g test_group %s) && (echo '%s ALL=(ALL) NOPASSWD: ALL' >> /etc/sudoers) && su %s -c 'umask 0077; %s'\"", runUser, runUser, runUser, inShellCmd)
	}
	cmd := fmt.Sprintf("%s run --privileged -d -t -v %s:/testing -v /usr/local/go:/usr/local/go %s", getDockerCmd(), goPath, dockerImage)
	out, err := os.Cmdf2(stage, "%v", cmd)
	if nil != err {
		t.Fatalf("\n---INNER LOG\n%v\n---END INNER LOG", out)
		return false
	}
	segs := strings.Split(strings.TrimSpace(out), "\n")
	dockerId := segs[len(segs)-1]
	defer os.Cmdf2(stage, "%s rm -f %s", getDockerCmd(), dockerId)
	if out, err = os.Cmdf2(stage, "%s exec -i %s %s", getDockerCmd(), dockerId, shellCmd); nil != err {
		t.Fatalf("\n---INNER LOG\n%v\n---END INNER LOG", out)
	} else {
		t.Logf("\n---INNER LOG\n%v\n---END INNER LOG", out)
	}

	return false
}
