package dry

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"runtime"
	"sync"
)

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
	buf := new(bytes.Buffer) // the returned data
	var lastFile string
	for i := 3; ; i++ {
		contin := fprintStackTraceLine(i, &lastFile, buf)
		if !contin {
			break
		}
	}
	return buf.String()
}

func StackTraceLine(skipFrames int) string {
	var buf bytes.Buffer
	var lastFile string
	fprintStackTraceLine(skipFrames, &lastFile, &buf)
	return buf.String()
}

func fprintStackTraceLine(i int, lastFile *string, buf *bytes.Buffer) bool {
	var lines [][]byte

	pc, file, line, ok := runtime.Caller(i)
	if !ok {
		return false
	}

	// Print this much at least.  If we can't find the source, it won't show.
	fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
	if file != *lastFile {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return true
		}
		lines = bytes.Split(data, []byte{'\n'})
		*lastFile = file
	}
	line-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	return true
}

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
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

func (self *DebugMutex) Lock() {
	fmt.Println("Mutex.Lock()\n" + StackTraceLine(3))
	self.m.Lock()
}

func (self *DebugMutex) Unlock() {
	fmt.Println("Mutex.Unlock()\n" + StackTraceLine(3))
	self.m.Unlock()
}

// DebugRWMutex wraps a sync.RWMutex and adds debug output
type DebugRWMutex struct {
	m sync.RWMutex
}

func (self *DebugRWMutex) RLock() {
	fmt.Println("RWMutex.RLock()\n" + StackTraceLine(3))
	self.m.RLock()
}

func (self *DebugRWMutex) RUnlock() {
	fmt.Println("RWMutex.RUnlock()\n" + StackTraceLine(3))
	self.m.RUnlock()
}

func (self *DebugRWMutex) Lock() {
	fmt.Println("RWMutex.Lock()\n" + StackTraceLine(3))
	self.m.Lock()
}

func (self *DebugRWMutex) Unlock() {
	fmt.Println("RWMutex.Unlock()\n" + StackTraceLine(3))
	self.m.Unlock()
}

func (self *DebugRWMutex) RLocker() sync.Locker {
	return self.m.RLocker()
}
