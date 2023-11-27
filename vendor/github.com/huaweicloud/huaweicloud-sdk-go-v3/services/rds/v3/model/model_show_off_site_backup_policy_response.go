package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowOffSiteBackupPolicyResponse Response Object
type ShowOffSiteBackupPolicyResponse struct {

	// 备份策略对象，包括备份类型、备份保留天数、目标区域ID和目标project ID。
	PolicyPara     *[]GetOffSiteBackupPolicy `json:"policy_para,omitempty"`
	HttpStatusCode int                       `json:"-"`
}

func (o ShowOffSiteBackupPolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowOffSiteBackupPolicyResponse struct{}"
	}

	return strings.Join([]string{"ShowOffSiteBackupPolicyResponse", string(data)}, " ")
}
