package params

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParam(t *testing.T) {
	ps := Params{
		&Param{
			Key:   "a",
			Value: "",
			Desc:  "a",
			Type:  ParamTypeString,
		},
		&Param{
			Key:   "b",
			Value: "",
			Desc:  "b",
			Type:  ParamTypeInt,
		},
		&Param{
			Key:   "c",
			Value: "",
			Desc:  "c",
			Type:  ParamTypeBool,
		},
	}
	// test set string
	err := ps.SetParamValue("a", "a")
	assert.NoError(t, err)
	assert.Equal(t, "a", ps.GetParam("a").String())

	err = ps.SetParamValue("a", "123")
	assert.NoError(t, err)
	assert.Equal(t, "123", ps.GetParam("a").String())

	err = ps.SetParamValue("a", "T")
	assert.NoError(t, err)
	assert.Equal(t, "T", ps.GetParam("a").String())

	err = ps.SetParamValue("a", "F")
	assert.NoError(t, err)
	assert.Equal(t, "F", ps.GetParam("a").String())

	// test set number
	err = ps.SetParamValue("b", "123")
	assert.NoError(t, err)
	assert.Equal(t, 123, ps.GetParam("b").Int())

	err = ps.SetParamValue("b", "b")
	assert.Error(t, err)
	assert.Equal(t, 123, ps.GetParam("b").Int()) // set value failed, value not change.

	err = ps.SetParamValue("b", "T")
	assert.Error(t, err)

	err = ps.SetParamValue("b", "F")
	assert.Error(t, err)

	// test set bool
	err = ps.SetParamValue("c", "c")
	assert.Error(t, err)
	assert.Equal(t, false, ps.GetParam("c").Bool())

	err = ps.SetParamValue("c", "1")
	assert.NoError(t, err)
	assert.Equal(t, true, ps.GetParam("c").Bool())

	err = ps.SetParamValue("c", "0")
	assert.NoError(t, err)
	assert.Equal(t, false, ps.GetParam("c").Bool())

	err = ps.SetParamValue("c", "T")
	assert.NoError(t, err)
	assert.Equal(t, true, ps.GetParam("c").Bool())

	err = ps.SetParamValue("c", "F")
	assert.NoError(t, err)
	assert.Equal(t, false, ps.GetParam("c").Bool())
}

func TestParams_ScanValue(t *testing.T) {
	ps := Params{
		&Param{
			Key:   "a",
			Value: "a",
			Desc:  "a",
			Type:  ParamTypeString,
		},
		&Param{
			Key:   "b",
			Value: "123",
			Desc:  "b",
			Type:  ParamTypeInt,
		},
		&Param{
			Key:   "c",
			Value: "T",
			Desc:  "c",
			Type:  ParamTypeBool,
		},
	}

	data, err := ps.Value()
	assert.NoError(t, err)

	var ps2 Params
	err = ps2.Scan(data)
	assert.NoError(t, err)
	assert.Equal(t, ps, ps2)

	ps3 := Params{
		&Param{
			Key:   "a",
			Value: "a",
			Desc:  "a",
			Type:  ParamTypeString,
		},
		&Param{
			Key:   "b",
			Value: "123",
			Desc:  "b",
			Type:  ParamTypeInt,
		},
		&Param{
			Key:   "c",
			Value: "T",
			Desc:  "d", // diff from ps
			Type:  ParamTypeBool,
		},
	}
	assert.NotEqual(t, ps3, ps2)

	var emptyPs Params
	data, err = emptyPs.Value()
	assert.NoError(t, err)

	var emptyPs2 Params
	err = emptyPs2.Scan(data)
	assert.NoError(t, err)
	assert.Equal(t, emptyPs, emptyPs2)

	data = []byte("this is test scan fail")
	var failPs Params
	err = failPs.Scan(data)
	assert.Error(t, err)
}
