package params

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/actiontech/dms/pkg/dms-common/i18nPkg"
	dmsCommonAes "github.com/actiontech/dms/pkg/dms-common/pkg/aes"
	"golang.org/x/text/language"
)

type Params []*Param

type ParamType string

const (
	ParamTypeString   ParamType = "string"
	ParamTypeInt      ParamType = "int"
	ParamTypeBool     ParamType = "bool"
	ParamTypeFloat64  ParamType = "float64"
	ParamTypePassword ParamType = "password"
)

type Param struct {
	Key      string          `json:"key"`
	Value    string          `json:"value"`
	Desc     string          `json:"desc"` // Deprecated: use I18nDesc instead
	I18nDesc i18nPkg.I18nStr `json:"i18n_desc"`
	Type     ParamType       `json:"type"`
	Enums    []EnumsValue    `json:"enums"`
}

type EnumsValue struct {
	Value    string          `json:"value"`
	Desc     string          `json:"desc"` // Deprecated: use I18nDesc instead
	I18nDesc i18nPkg.I18nStr `json:"i18n_desc"`
}

func (e EnumsValue) GetDesc(lang language.Tag) string {
	if e.Desc != "" {
		e.I18nDesc.SetStrInLang(i18nPkg.DefaultLang, e.Desc)
	}
	return e.I18nDesc.GetStrInLang(lang)
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

func (r *Param) GetDesc(lang language.Tag) string {
	if r == nil {
		return ""
	}
	if r.Desc != "" {
		r.I18nDesc.SetStrInLang(i18nPkg.DefaultLang, r.Desc)
	}
	return r.I18nDesc.GetStrInLang(lang)
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

	for _, p := range result {
		if p.Type == ParamTypePassword {
			p.Value, err = dmsCommonAes.AesDecrypt(p.Value)
			if err != nil {
				return fmt.Errorf("param %s value decrypt err: %v", p.Key, err)
			}
		}
	}

	*r = result
	return err
}

// Value impl sql.driver.Valuer interface
func (r Params) Value() (driver.Value, error) {
	if len(r) == 0 {
		return nil, nil
	}

	params := make([]Param, 0, len(r))

	for _, p := range r {
		param := Param{
			Key:      p.Key,
			Value:    p.Value,
			Desc:     p.Desc,
			I18nDesc: p.I18nDesc,
			Type:     p.Type,
		}

		if param.Type == ParamTypePassword {
			val, err := dmsCommonAes.AesEncrypt(p.Value)
			if err != nil {
				return nil, fmt.Errorf("param %s value encrypt err: %v", p.Key, err)
			}

			param.Value = val
		}

		params = append(params, param)
	}

	v, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json value: %v", v)
	}
	return v, err
}

func (r *Params) Copy() Params {
	ps := make(Params, 0, len(*r))
	for _, p := range *r {
		ps = append(ps, &Param{
			Key:      p.Key,
			Value:    p.Value,
			Desc:     p.Desc,
			I18nDesc: p.I18nDesc.Copy(),
			Type:     p.Type,
		})
	}
	return ps
}

type ParamsWithOperator []*ParamWithOperator
type ParamWithOperator struct {
	Param
	Operator Operator `json:"operator"`
}

// Scan impl sql.Scanner interface
func (r *ParamsWithOperator) Scan(value interface{}) error {
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
	result := ParamsWithOperator{}
	err := json.Unmarshal(bytes, &result)

	for _, p := range result {
		if p.Type == ParamTypePassword {
			p.Value, err = dmsCommonAes.AesDecrypt(p.Value)
			if err != nil {
				return fmt.Errorf("param %s value decrypt err: %v", p.Key, err)
			}
		}
	}

	*r = result
	return err
}

// Value impl sql.driver.Valuer interface
func (r ParamsWithOperator) Value() (driver.Value, error) {
	if len(r) == 0 {
		return nil, nil
	}

	params := make([]ParamWithOperator, 0, len(r))

	for _, p := range r {
		param := ParamWithOperator{
			Param: Param{
				Key:      p.Key,
				Value:    p.Value,
				Desc:     p.Desc,
				I18nDesc: p.I18nDesc,
				Type:     p.Type,
			},
			Operator: Operator{
				Value:      p.Operator.Value,
				EnumsValue: p.Operator.EnumsValue,
			},
		}

		if param.Type == ParamTypePassword {
			val, err := dmsCommonAes.AesEncrypt(p.Value)
			if err != nil {
				return nil, fmt.Errorf("param %s value encrypt err: %v", p.Key, err)
			}

			param.Value = val
		}

		params = append(params, param)
	}

	v, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json value: %v", v)
	}
	return v, err
}

func (r *ParamsWithOperator) GetParam(key string) *ParamWithOperator {
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

type Operator struct {
	Value      OperatorValue `json:"boolean_operator_value"`
	EnumsValue []EnumsValue  `json:"boolean_operator_enums_value"`
}

type OperatorValue string

const (
	LessThanOperator             OperatorValue = "<"
	GreaterThanOperator          OperatorValue = ">"
	LessThanOrEqualToOperator    OperatorValue = "<="
	GreaterThanOrEqualToOperator OperatorValue = ">="
	EqualToOperator              OperatorValue = "="
	NotEqualToOperator           OperatorValue = "<>"
	InOperator                   OperatorValue = "IN"
	IsOperator                   OperatorValue = "IS"
	ContainsOperator             OperatorValue = "CONTAINS"
)

func (r *ParamsWithOperator) CompareParamValue(key string, inputValue string) (bool, error) {
	paramNotFoundErrMsg := "param %s not found"
	if r == nil {
		return false, fmt.Errorf(paramNotFoundErrMsg, key)
	}

	param := r.GetParam(key)
	if param == nil {
		return false, fmt.Errorf(paramNotFoundErrMsg, key)
	}

	// Perform comparison based on the type of the parameter
	switch param.Type {
	case ParamTypeInt:
		paramValue, err := strconv.Atoi(param.Value)
		if err != nil {
			return false, fmt.Errorf("failed to convert param value to int: %v", err)
		}
		inputIntValue, err := strconv.Atoi(inputValue)
		if err != nil {
			return false, fmt.Errorf("failed to convert input value to int: %v", err)
		}
		return compareInt(paramValue, inputIntValue, param.Operator.Value), nil

	case ParamTypeFloat64:
		paramValue, err := strconv.ParseFloat(param.Value, 64)
		if err != nil {
			return false, fmt.Errorf("failed to convert param value to float64: %v", err)
		}
		inputFloatValue, err := strconv.ParseFloat(inputValue, 64)
		if err != nil {
			return false, fmt.Errorf("failed to convert input value to float64: %v", err)
		}
		return compareFloat64(paramValue, inputFloatValue, param.Operator.Value), nil

	case ParamTypeString:
		return compareString(param.Value, inputValue, param.Operator.Value), nil

	default:
		return false, fmt.Errorf("unsupported ParamType: %s", param.Type)
	}
}

// Helper functions to perform comparison based on Operator
func compareInt(paramValue, inputValue int, operator OperatorValue) bool {
	switch operator {
	case LessThanOperator:
		return inputValue < paramValue
	case GreaterThanOperator:
		return inputValue > paramValue
	case LessThanOrEqualToOperator:
		return inputValue <= paramValue
	case GreaterThanOrEqualToOperator:
		return inputValue >= paramValue
	case EqualToOperator:
		return inputValue == paramValue
	case NotEqualToOperator:
		return inputValue != paramValue
	default:
		return false
	}
}

func compareFloat64(paramValue, inputValue float64, operator OperatorValue) bool {
	switch operator {
	case LessThanOperator:
		return inputValue < paramValue
	case GreaterThanOperator:
		return inputValue > paramValue
	case LessThanOrEqualToOperator:
		return inputValue <= paramValue
	case GreaterThanOrEqualToOperator:
		return inputValue >= paramValue
	case EqualToOperator:
		return inputValue == paramValue
	case NotEqualToOperator:
		return inputValue != paramValue
	default:
		return false
	}
}

func compareString(paramValue, inputValue string, operator OperatorValue) bool {
	switch operator {
	case LessThanOperator:
		return inputValue < paramValue
	case GreaterThanOperator:
		return inputValue > paramValue
	case LessThanOrEqualToOperator:
		return inputValue <= paramValue
	case GreaterThanOrEqualToOperator:
		return inputValue >= paramValue
	case EqualToOperator:
		return inputValue == paramValue
	case NotEqualToOperator:
		return inputValue != paramValue
	case ContainsOperator:
		return contains(paramValue, inputValue)
	default:
		return false
	}
}

func contains(paramValue, inputValue string) bool {
	return paramValue == inputValue
}
