package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SetBinlogClearPolicyResponse Response Object
type SetBinlogClearPolicyResponse struct {

	// 操作结果。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o SetBinlogClearPolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetBinlogClearPolicyResponse struct{}"
	}

	return strings.Join([]string{"SetBinlogClearPolicyResponse", string(data)}, " ")
}
