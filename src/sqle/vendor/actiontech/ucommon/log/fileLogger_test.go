package log

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestRemoveLoop(t *testing.T) {
	dir := "/tmp/testLog/"
	instanceMu.Lock()
	instance = newLogger(dir, "/tmp/conf/log.conf")
	instanceMu.Unlock()
	instance.setFileLimit(1, 3)
	instance.removeLoop()

	data := make([]byte, 1024*1024)
	tarFile1 := "2017_12_26_05_29_35_detail.log.tar.gz"
	tarFile2 := "2017_12_26_05_29_36_detail.log.tar.gz"
	std1 := "std_01.log"
	std2 := "std_02.log"
	std3 := "std_03.log"
	std4 := "std_04.log"
	writeFile(dir+tarFile1, data)
	writeFile(dir+tarFile2, data)
	writeFile(dir+std1, data)
	writeFile(dir+std2, data)
	writeFile(dir+std3, data)
	writeFile(dir+std4, data)

	// (removeLoop should be set to 1s.)
	time.Sleep(2 * time.Second)
	if !isExistFile(tarFile1) || !isExistFile(tarFile2) {
		t.Errorf("totalLimit is:%vM, but total size is:%vM\n", 3, 2)
	}

	t.Logf("before SetTotalLimit, totalLimit is:%vM, current size is:%vM\n", 3, 2)
	instance.setFileLimit(2, 2)
	t.Logf("after SetTotalLimit, totalLimit is:%vM, current size is:%vM\n", 2, 2)

	time.Sleep(2 * time.Second)
	if isExistFile(tarFile1) || !isExistFile(tarFile2) {
		t.Errorf("totalLimit is:%vM, but total size is:%vM\n", 2, 2)
	}
}

func writeFile(fileName string, data []byte) {
	f, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	f.Write(data)
	f.Sync()
	f.Close()
}

func isExistFile(fileName string) bool {
	var exist bool
	files, _ := ioutil.ReadDir("/tmp/testLog")
	for _, file := range files {
		if !file.IsDir() && file.Name() == fileName {
			exist = true
		}
	}
	return exist
}
