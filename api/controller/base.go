package controller

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/url"
	"sqle/errors"
)

var INSTANCE_NOT_EXIST_ERROR = NewBaseReq(errors.New(errors.INSTANCE_NOT_EXIST, fmt.Errorf("instance not exist")))
var INSTANCE_EXIST_ERROR = NewBaseReq(errors.New(errors.INSTANCE_EXIST, fmt.Errorf("inst is exist")))
var TASK_NOT_EXIST = NewBaseReq(errors.New(errors.TASK_NOT_EXIST, fmt.Errorf("task not exist")))

type BaseRes struct {
	Code    int    `json:"code" example:"0"`
	Message string `json:"message" example:"ok"`
}

func NewBaseReq(err error) BaseRes {
	res := BaseRes{}
	switch e := err.(type) {
	case *errors.CodeError:
		res.Code = e.Code()
		res.Message = e.Error()
	default:
		if err == nil {
			res.Code = 0
			res.Message = "ok"
		} else {
			res.Code = -1
			res.Message = e.Error()
		}
	}
	return res
}

func readFileToByte(c echo.Context, name string) (fileName string, data []byte, err error) {
	file, err := c.FormFile(name)
	if err != nil {
		err = errors.New(errors.READ_UPLOAD_FILE_ERROR, err)
		return
	}
	src, err := file.Open()
	if err != nil {
		err = errors.New(errors.READ_UPLOAD_FILE_ERROR, err)
		return
	}
	defer src.Close()
	data, err = ioutil.ReadAll(src)
	if err != nil {
		err = errors.New(errors.READ_UPLOAD_FILE_ERROR, err)
		return
	}
	return
}

type CustomValidator struct {
}

func (cv *CustomValidator) Validate(i interface{}) error {
	_, err := govalidator.ValidateStruct(i)
	return err
}

func unescapeParamString(params []*string) error {
	for i, p := range params {
		r, err := url.QueryUnescape(*p)
		if nil != err {
			return fmt.Errorf("unescape param [%v] failed: %v", params, err)
		}
		*params[i] = r
	}
	return nil
}
