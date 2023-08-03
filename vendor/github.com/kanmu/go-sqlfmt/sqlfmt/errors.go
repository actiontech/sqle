package sqlfmt

import (
	"fmt"
)

// FormatError is an error that occurred while sqlfmt.Process
type FormatError struct {
	msg string
}

func (e *FormatError) Error() string {
	return fmt.Sprint(e.msg)
}
