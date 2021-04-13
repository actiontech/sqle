package controller

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateName(t *testing.T) {
	assert.Equal(t, true, validateName("name"))
	assert.Equal(t, true, validateName("name_"))
	assert.Equal(t, true, validateName("name-"))
	assert.Equal(t, true, validateName("name-1"))
	assert.Equal(t, true, validateName("name_1"))

	len60 := `test_name_length_60_0000000000000000000000000000000000000000`
	assert.Equal(t, 60, len(len60))
	assert.Equal(t, true, validateName(len60))

	assert.Equal(t, false, validateName(""))
	assert.Equal(t, false, validateName("1name"))
	assert.Equal(t, false, validateName("_name"))
	assert.Equal(t, false, validateName("-name"))
	assert.Equal(t, false, validateName("name*"))
	assert.Equal(t, false, validateName("name*name"))
	assert.Equal(t, false, validateName("*name"))

	len61 := `test_name_length_61_00000000000000000000000000000000000000000`
	assert.Equal(t, 61, len(len61))
	assert.Equal(t, false, validateName(len61))
	assert.Equal(t, false, validateName(len61+"*"))
}

func TestValidatePort(t *testing.T) {
	assert.Equal(t, true, validatePort("1"))
	assert.Equal(t, true, validatePort("3306"))
	assert.Equal(t, true, validatePort("65535"))

	assert.Equal(t, false, validatePort("0"))
	assert.Equal(t, false, validatePort("65536"))
	assert.Equal(t, false, validatePort(""))
	assert.Equal(t, false, validatePort("port"))
	assert.Equal(t, false, validatePort("_"))
}

func TestCustomValidateErrorMessage(t *testing.T) {
	type tSingleError struct {
		Name string `json:"name" valid:"name"`
	}
	assert.Equal(t, "tSingleError.name must match regexp `^[a-zA-Z][a-zA-Z0-9\\_\\-]{0,59}$`",
		Validate(&tSingleError{Name: "_name"}).Error())

	type tMultiError struct {
		Name string `json:"name" valid:"name"`
		Port string `json:"port" valid:"port"`
	}
	assert.Equal(t, "tMultiError.name must match regexp `^[a-zA-Z][a-zA-Z0-9\\_\\-]{0,59}$`; "+
		"tMultiError.port is invalid port",
		Validate(&tMultiError{Name: "_name", Port: "0"}).Error())
}
