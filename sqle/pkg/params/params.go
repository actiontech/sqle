package params

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

type Params []*Param

type ParamType string

const (
	ParamTypeString  ParamType = "string"
	ParamTypeInt     ParamType = "int"
	ParamTypeBool    ParamType = "bool"
	ParamTypeFloat64 ParamType = "float64"
)

type Param struct {
	Key   string    `json:"key"`
	Value string    `json:"value"`
	Desc  string    `json:"desc"`
	Type  ParamType `json:"type"`
}

func (r *Params) SetParamValue(key, value string) error {
	paramNotFoundErrMsg := "param %s not found"
	if r == nil {
		return fmt.Errorf(paramNotFoundErrMsg, key)
	}
	for _, p := range *r {
		var err error
		if p.Key == key {
			switch p.Type {
			case ParamTypeBool:
				_, err = strconv.ParseBool(value)
			case ParamTypeInt:
				_, err = strconv.Atoi(value)
			default:
			}
			if err != nil {
				return fmt.Errorf("param %s value don't match \"%s\"", key, p.Type)
			}
			p.Value = value
			return nil
		}
	}
	return fmt.Errorf(paramNotFoundErrMsg, key)
}

func (r *Params) GetParam(key string) *Param {
	if r == nil {
		return nil
	}
	for _, p := range *r {
		if p.Key == key {
			return p
		}
	}
	return nil
}

func (r *Param) String() string {
	if r == nil {
		return ""
	}
	return r.Value
}

func (r *Param) Int() int {
	if r == nil {
		return 0
	}
	i, err := strconv.Atoi(r.Value)
	if err != nil {
		return 0
	}
	return i
}

func (r *Param) Float64() float64 {
	if r == nil {
		return 0
	}

	i, err := strconv.ParseFloat(r.Value, 64)
	if err != nil {
		return 0
	}
	return i
}

func (r *Param) Bool() bool {
	if r == nil {
		return false
	}
	b, err := strconv.ParseBool(r.Value)
	if err != nil {
		return false
	}
	return b
}

// Scan impl sql.Scanner interface
func (r *Params) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal json value: %v", value)
	}
	if len(bytes) == 0 {
		return nil
	}
	result := Params{}
	err := json.Unmarshal(bytes, &result)
	*r = result
	return err
}

// Value impl sql.driver.Valuer interface
func (r Params) Value() (driver.Value, error) {
	if len(r) == 0 {
		return nil, nil
	}
	v, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json value: %v", v)
	}
	return v, err
}

func (r *Params) Copy() Params {
	ps := make(Params, 0, len(*r))
	for _, p := range *r {
		ps = append(ps, &Param{
			Key:   p.Key,
			Value: p.Value,
			Desc:  p.Desc,
			Type:  p.Type,
		})
	}
	return ps
}
