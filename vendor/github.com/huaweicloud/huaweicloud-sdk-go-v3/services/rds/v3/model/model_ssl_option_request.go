package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SslOptionRequest struct {

	// - true, 打开ssl开关。 - false, 关闭ssl开关。
	SslOption bool `json:"ssl_option"`
}

func (o SslOptionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SslOptionRequest struct{}"
	}

	return strings.Join([]string{"SslOptionRequest", string(data)}, " ")
}
