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
