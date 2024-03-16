package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DataIpRequest struct {

	// 内网ip
	NewIp string `json:"new_ip"`
}

func (o DataIpRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DataIpRequest struct{}"
	}

	return strings.Join([]string{"DataIpRequest", string(data)}, " ")
}
