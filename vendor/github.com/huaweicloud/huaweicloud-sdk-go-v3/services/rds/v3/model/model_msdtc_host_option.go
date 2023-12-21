package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type MsdtcHostOption struct {

	// 主机名称 hostname
	HostName string `json:"host_name"`

	// 主机ip
	Ip string `json:"ip"`
}

func (o MsdtcHostOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MsdtcHostOption struct{}"
	}

	return strings.Join([]string{"MsdtcHostOption", string(data)}, " ")
}
