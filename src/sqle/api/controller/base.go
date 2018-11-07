package controller

import (
	"fmt"
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
