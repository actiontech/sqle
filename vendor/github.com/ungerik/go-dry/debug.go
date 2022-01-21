package dry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
	"sync"
)

// PrettyPrintAsJSON marshalles input as indented JSON
// and calles fmt.Println with the result.
// If indent arguments are given, they are joined into
// a string and used as JSON line indent.
// If no indet argument is given, two spaces will be used
// to indent JSON lines.
// A byte slice as input will be marshalled as json.RawMessage.
func PrettyPrintAsJSON(input interface{}, indent ...string) error {
	var indentStr string
	if len(indent) == 0 {
		indentStr = "  "
	} else {
		indentStr = strings.Join(indent, "")
	}
	if b, ok := input.([]byte); ok {
		input = json.RawMessage(b)
	}
	data, err := json.MarshalIndent(input, "", indentStr)
	if err != nil {
		return err
	}
	_, err = fmt.Println(string(data))
	return err
}

// Nop is a dummy function that can be called in source files where
// other debug functions are constantly added and removed.
// That way import "github.com/ungerik/go-quick" won't cause an error when
// no other debug function is currently used.
// Arbitrary objects can be passed as arguments to avoid "declared and not used"
// error messages when commenting code out and in.
// The result is a nil interface{} dummy value.
func Nop(dummiesIn ...interface{}) (dummyOut interface{}) {
	return nil
}

func StackTrace(skipFrames int) string {
	var b strings.Builder
	var lastFile string
	for i := 3; ; i++ {
		contin := fprintStackTraceLine(i, &lastFile, &b)
		if !contin {
			break
		}
	}
	return b.String()
}

func StackTraceLine(skipFrames int) string {
	var b strings.Builder
	var lastFile string
	fprintStackTraceLine(skipFrames, &lastFile, &b)
	return b.String()
}

func fprintStackTraceLine(i int, lastFile *string, b *strings.Builder) bool {
	var lines [][]byte

	pc, file, line, ok := runtime.Caller(i)
	if !ok {
		return false
	}

	// Print this much at least.  If we can't find the source, it won't show.
	fmt.Fprintf(b, "%s:%d (0x%x)\n", file, line, pc)
	if file != *lastFile {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return true
		}
		lines = bytes.Split(data, []byte{'\n'})
		*lastFile = file
	}
	line-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	fmt.Fprintf(b, "\t%s: %s\n", function(pc), source(lines, line))
	return true
}

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	sep       = []byte("/")
)

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.Trim(lines[n], " \t")
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod

	if period := bytes.LastIndex(name, sep); period >= 0 {
		name = name[period+1:]
	}

	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

// DebugMutex wraps a sync.Mutex and adds debug output
type DebugMutex struct {
	m sync.Mutex
}

func (d *DebugMutex) Lock() {
	fmt.Println("Mutex.Lock()\n" + StackTraceLine(3))
	d.m.Lock()
}

func (d *DebugMutex) Unlock() {
	fmt.Println("Mutex.Unlock()\n" + StackTraceLine(3))
	d.m.Unlock()
}

// DebugRWMutex wraps a sync.RWMutex and adds debug output
type DebugRWMutex struct {
	m sync.RWMutex
}

func (d *DebugRWMutex) RLock() {
	fmt.Println("RWMutex.RLock()\n" + StackTraceLine(3))
	d.m.RLock()
}

func (d *DebugRWMutex) RUnlock() {
	fmt.Println("RWMutex.RUnlock()\n" + StackTraceLine(3))
	d.m.RUnlock()
}

func (d *DebugRWMutex) Lock() {
	fmt.Println("RWMutex.Lock()\n" + StackTraceLine(3))
	d.m.Lock()
}

func (d *DebugRWMutex) Unlock() {
	fmt.Println("RWMutex.Unlock()\n" + StackTraceLine(3))
	d.m.Unlock()
}

func (d *DebugRWMutex) RLocker() sync.Locker {
	return d.m.RLocker()
}
