package controller

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateName(t *testing.T) {
	assert.Equal(t, true, validateName("姓名name"))
	assert.Equal(t, true, validateName("姓名name_"))
	assert.Equal(t, true, validateName("姓名name-"))
	assert.Equal(t, true, validateName("姓名name-1"))
	assert.Equal(t, true, validateName("name-1"))
	assert.Equal(t, true, validateName("姓名name_1"))
	assert.Equal(t, true, validateName("name_1"))

	len120 := `test_name_length_120_000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000`
	assert.Equal(t, 120, len(len120))
	assert.Equal(t, true, validateName(len120))

	assert.Equal(t, false, validateName(""))
	assert.Equal(t, false, validateName("1name姓名"))
	assert.Equal(t, false, validateName("_姓名name"))
	assert.Equal(t, false, validateName("-姓名name"))
	assert.Equal(t, false, validateName("姓名name*"))
	assert.Equal(t, false, validateName("name*"))
	assert.Equal(t, false, validateName("name姓名*name"))
	assert.Equal(t, false, validateName("name*name"))
	assert.Equal(t, false, validateName("*name姓名"))

	len121 := `test_name_length_121_0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000`
	assert.Equal(t, 121, len(len121))
	assert.Equal(t, false, validateName(len121))
	assert.Equal(t, false, validateName(len121+"*"))
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
	assert.Equal(t, "tSingleError.name must match regexp `^[a-zA-Z一-龥][a-zA-Z0-9一-龥_-]{0,119}$`",
		Validate(&tSingleError{Name: "_name"}).Error())

	type tMultiError struct {
		Name string `json:"name" valid:"name"`
		Port string `json:"port" valid:"port"`
	}
	assert.Equal(t, "tMultiError.name must match regexp `^[a-zA-Z一-龥][a-zA-Z0-9一-龥_-]{0,119}$`; "+
		"tMultiError.port is invalid port",
		Validate(&tMultiError{Name: "_name", Port: "0"}).Error())
}
