package errors

import (
	_errors "errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCombine(t *testing.T) {
	err1 := _errors.New("err1")
	err2 := _errors.New("err2")
	err := Combine(err1, err2)
	assert.Equal(t, "multi err: err1 err2 ", err.Error())

	err = Combine(nil)
	assert.Nil(t, err)
}
