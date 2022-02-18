package errors

import "strings"

type multiErrors []error

func (e multiErrors) Error() string {
	var r strings.Builder
	r.WriteString("multi err: ")
	for _, err := range e {
		r.WriteString(err.Error())
		r.WriteString(" ")
	}
	return r.String()
}

func Combine(maybeError ...error) error {
	var errs multiErrors
	for _, err := range maybeError {
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}
