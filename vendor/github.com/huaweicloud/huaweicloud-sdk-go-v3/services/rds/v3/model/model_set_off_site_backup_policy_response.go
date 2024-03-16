package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SetOffSiteBackupPolicyResponse Response Object
type SetOffSiteBackupPolicyResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o SetOffSiteBackupPolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetOffSiteBackupPolicyResponse struct{}"
	}

	return strings.Join([]string{"SetOffSiteBackupPolicyResponse", string(data)}, " ")
}
