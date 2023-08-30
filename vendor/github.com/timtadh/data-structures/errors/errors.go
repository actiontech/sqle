package errors

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

var SkipLogging map[string]bool

func init() {
	SkipLogging = make(map[string]bool, 10)
}

type Error struct {
	Errs  []error
	Stack []byte
}

func Errorf(format string, args ...interface{}) error {
	buf := make([]byte, 50000)
	n := runtime.Stack(buf, true)
	trace := make([]byte, n)
	copy(trace, buf)
	return &Error{
		Errs:  []error{fmt.Errorf(format, args...)},
		Stack: trace,
	}
}

func Logf(level, format string, args ...interface{}) {
	if SkipLogging[level] {
		return
	}
	pc, _, line, ok := runtime.Caller(1)
	if !ok {
		log.Printf(format, args)
		return
	}
	fn := runtime.FuncForPC(pc)
	msg := fmt.Sprintf(format, args...)
	log.Printf("%v (%v:%v): %v", level, fn.Name(), line, msg)
}

func (e *Error) Chain(err error) error {
	e.Errs = append(e.Errs, err)
	return e
}

func (e *Error) Error() string {
	if e == nil {
		return "Error <nil>"
	} else if len(e.Errs) == 0 {
		return fmt.Sprintf("%v\n%s", e.Errs, string(e.Stack))
	} else if len(e.Errs) == 1 {
		return fmt.Sprintf("%v\n%s", e.Errs[0], string(e.Stack))
	} else {
		errs := make([]string, 0, len(e.Errs))
		for _, err := range e.Errs {
			errs = append(errs, err.Error())
		}
		return fmt.Sprintf("{%v}\n%s", strings.Join(errs, ", "), string(e.Stack))
	}
}

func (e *Error) String() string {
	return e.Error()
}

type ErrorFmter func(a ...interface{}) error

func NotFound(a ...interface{}) error {
	// return fmt.Errorf("Key '%v' was not found.", a...)
	return Errorf("Key was not found.")
}

func NotFoundInBucket(a ...interface{}) error {
	return Errorf("Key, '%v', was not in bucket when expected.", a...)
}

func InvalidKey(a ...interface{}) error {
	return Errorf("Key, '%v', is invalid, %s", a...)
}

func TSTError(a ...interface{}) error {
	return Errorf("Internal TST error - "+a[0].(string), a[1:]...)
}

func NegativeSize(a ...interface{}) error {
	return Errorf("Negative size")
}

func BpTreeError(a ...interface{}) error {
	return Errorf("Internal B+ Tree error - "+a[0].(string), a[1:]...)
}

var Errors map[string]ErrorFmter = map[string]ErrorFmter{
	"not-found":           NotFound,
	"not-found-in-bucket": NotFoundInBucket,
	"invalid-key":         InvalidKey,
	"tst-error":           TSTError,
	"negative-size":       NegativeSize,
	"bptree-error":        BpTreeError,
}
