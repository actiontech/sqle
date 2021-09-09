package reform

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Inspect returns suitable for logging representation of a query argument.
func Inspect(arg interface{}, addType bool) string {
	var s string
	v := reflect.ValueOf(arg)
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			s = "<nil>"
		} else {
			s = Inspect(v.Elem().Interface(), false)
		}

	case reflect.String:
		s = fmt.Sprintf("%#q", arg)

	default:
		s = fmt.Sprintf("%v", arg)
	}

	if addType {
		s += fmt.Sprintf(" (%T)", arg)
	}
	return s
}

// Logger is responsible to log queries before and after their execution.
type Logger interface {
	// Before logs query before execution.
	Before(query string, args []interface{})

	// After logs query after execution.
	After(query string, args []interface{}, d time.Duration, err error)
}

// Printf is a (fmt.Printf|log.Printf|testing.T.Logf)-like function.
type Printf func(format string, args ...interface{})

// PrintfLogger is a simple query logger.
type PrintfLogger struct {
	LogTypes bool
	printf   Printf
}

// NewPrintfLogger creates a new simple query logger for any Printf-like function.
func NewPrintfLogger(printf Printf) *PrintfLogger {
	return &PrintfLogger{false, printf}
}

// Before logs query before execution.
func (pl *PrintfLogger) Before(query string, args []interface{}) {
	// fast path
	if args == nil {
		pl.printf(">>> %s", query)
		return
	}

	ss := make([]string, len(args))
	for i, arg := range args {
		ss[i] = Inspect(arg, pl.LogTypes)
	}

	pl.printf(">>> %s [%s]", query, strings.Join(ss, ", "))
}

// After logs query after execution.
func (pl *PrintfLogger) After(query string, args []interface{}, d time.Duration, err error) {
	// fast path
	if args == nil {
		msg := fmt.Sprintf("%s %s", query, d)
		if err != nil {
			msg += ": " + err.Error()
		}
		pl.printf("<<< %s", msg)
		return
	}

	ss := make([]string, len(args))
	for i, arg := range args {
		ss[i] = Inspect(arg, pl.LogTypes)
	}

	msg := fmt.Sprintf("%s [%s] %s", query, strings.Join(ss, ", "), d)
	if err != nil {
		msg += ": " + err.Error()
	}
	pl.printf("<<< %s", msg)
}

// check interface
var _ Logger = (*PrintfLogger)(nil)
