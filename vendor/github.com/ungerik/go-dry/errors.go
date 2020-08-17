package dry

import (
	"fmt"
	"strings"
)

// PanicIfErr panics with a stack trace if any
// of the passed args is a non nil error
func PanicIfErr(args ...interface{}) {
	for _, v := range args {
		if err, _ := v.(error); err != nil {
			panic(fmt.Errorf("Panicking because of error: %s\nAt:\n%s\n", err, StackTrace(2)))
		}
	}
}

// GetError returns the last argument that is of type error,
// panics if none of the passed args is of type error.
// Note that GetError(nil) will panic because nil is not of type error but interface{}
func GetError(args ...interface{}) error {
	for i := len(args) - 1; i >= 0; i-- {
		arg := args[i]
		if arg != nil {
			if err, ok := arg.(error); ok {
				return err
			}
		}
	}
	panic("no argument of type error")
}

// AsError returns r as error, converting it when necessary
func AsError(r interface{}) error {
	if r == nil {
		return nil
	}
	if err, ok := r.(error); ok {
		return err
	}
	return fmt.Errorf("%v", r)
}

// FirstError returns the first non nil error, or nil
func FirstError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// LastError returns the last non nil error, or nil
func LastError(errs ...error) error {
	for i := len(errs) - 1; i >= 0; i-- {
		err := errs[i]
		if err != nil {
			return err
		}
	}
	return nil
}

// AsErrorList checks if err is already an ErrorList
// and returns it if this is the case.
// Else an ErrorList with err as element is created.
// Useful if a function potentially returns an ErrorList as error
// and you want to avoid creating nested ErrorLists.
func AsErrorList(err error) ErrorList {
	if list, ok := err.(ErrorList); ok {
		return list
	}
	return ErrorList{err}
}

/*
ErrorList holds a slice of errors.

Usage example:

	func maybeError() (int, error) {
		return
	}

	func main() {
		e := NewErrorList(maybeError())
		e.Collect(maybeError())
		e.Collect(maybeError())

		if e.Err() != nil {
			fmt.Println("Some calls of maybeError() returned errors:", e)
		} else {
			fmt.Println("No call of maybeError() returned an error")
		}
	}
*/
type ErrorList []error

// NewErrorList returns an ErrorList where Collect has been called for args.
// The returned list will be nil if there was no non nil error in args.
// Note that alle methods of ErrorList can be called with a nil ErrorList.
func NewErrorList(args ...interface{}) (list ErrorList) {
	list.Collect(args...)
	return list
}

// Error calls fmt.Println for of every error in the list
// and returns the concernated text.
// Can be called for a nil ErrorList.
func (list ErrorList) Error() string {
	if len(list) == 0 {
		return "Empty ErrorList"
	}
	var b strings.Builder
	for _, err := range list {
		fmt.Fprintln(&b, err)
	}
	return b.String()
}

// Err returns the list if it is not empty,
// or nil if it is empty.
// Can be called for a nil ErrorList.
func (list ErrorList) Err() error {
	if len(list) == 0 {
		return nil
	}
	return list
}

// First returns the first error in the list or nil.
// Can be called for a nil ErrorList.
func (list ErrorList) First() error {
	if len(list) == 0 {
		return nil
	}
	return list[0]
}

// Last returns the last error in the list or nil.
// Can be called for a nil ErrorList.
func (list ErrorList) Last() error {
	if len(list) == 0 {
		return nil
	}
	return list[len(list)-1]
}

// Collect adds any non nil errors in args to the list.
func (list *ErrorList) Collect(args ...interface{}) {
	for _, a := range args {
		if err, _ := a.(error); err != nil {
			*list = append(*list, err)
		}
	}
}
