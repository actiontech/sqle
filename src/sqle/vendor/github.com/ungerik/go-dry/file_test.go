package dry

import "testing"

func Test_FileGetString(t *testing.T) {
	_, err := FileGetString("invalid_file")
	if err == nil {
		t.Fail()
	}

	str, err := FileGetString("testfile.txt")
	if err != nil {
		t.Error(err)
	}
	if str != "Hello World!" {
		t.Fail()
	}

	str, err = FileGetString("https://raw.github.com/ungerik/go-dry/master/testfile.txt")
	if err != nil {
		t.Error(err)
	}
	if str != "Hello World!" {
		t.Fail()
	}
}

func Test_FileIsDir(t *testing.T) {
	if FileIsDir("testfile.txt") {
		t.Fail()
	}
	if !FileIsDir(".") {
		t.Fail()
	}
}
