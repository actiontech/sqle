package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DbsInstanceHostInfoResult struct {

	// host  id
	Id *string `json:"id,omitempty"`

	// host地址
	Host *string `json:"host,omitempty"`

	// host 名称
	HostName *string `json:"host_name,omitempty"`
}

func (o DbsInstanceHostInfoResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DbsInstanceHostInfoResult struct{}"
	}

	return strings.Join([]string{"DbsInstanceHostInfoResult", string(data)}, " ")
}
