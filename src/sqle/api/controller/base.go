package controller

type BaseRes struct {
	Code    int    `json:"code" example:"0"`
	Message string `json:"message" example:"ok"`
}

func NewBaseReq(code int, message string) BaseRes {
	return BaseRes{
		Code:    code,
		Message: message,
	}
}

