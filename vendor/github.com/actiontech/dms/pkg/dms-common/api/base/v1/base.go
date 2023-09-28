package v1

// GenericResp defines the return code and msg
// swagger:response GenericResp
type GenericResp struct {
	// code
	Code int `json:"code"`
	// message
	Message string `json:"message"`
}

func (r *GenericResp) SetCode(code int) {
	r.Code = code
}

func (r *GenericResp) SetMsg(msg string) {
	r.Message = msg
}

type GenericResper interface {
	SetCode(code int)
	SetMsg(msg string)
}
