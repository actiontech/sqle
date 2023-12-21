package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SetBackupPolicyResponse Response Object
type SetBackupPolicyResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o SetBackupPolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetBackupPolicyResponse struct{}"
	}

	return strings.Join([]string{"SetBackupPolicyResponse", string(data)}, " ")
}
