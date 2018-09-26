package dry

import (
	// "strings"
	// "fmt"
	"errors"
	"testing"
)

func Test_Error(t *testing.T) {
	err := AsError("TestError")
	if err == nil || err.Error() != "TestError" {
		t.Fail()
	}

	err = AsError(errors.New("TestError"))
	if err == nil || err.Error() != "TestError" {
		t.Fail()
	}
}
