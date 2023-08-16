package v1

// GenericResp defines the return code and msg
// swagger:response GenericResp
type GenericResp struct {
	// code
	Code int `json:"code"`
	// message
	Msg string `json:"msg"`
}

func (r *GenericResp) SetCode(code int) {
	r.Code = code
}

func (r *GenericResp) SetMsg(msg string) {
	r.Msg = msg
}

type GenericResper interface {
	SetCode(code int)
	SetMsg(msg string)
}
