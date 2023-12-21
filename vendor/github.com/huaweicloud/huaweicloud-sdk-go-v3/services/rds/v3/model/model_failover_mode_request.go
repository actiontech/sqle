package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type FailoverModeRequest struct {

	// 同步模式，各引擎可选择方式具体如下： MySQL： - async：异步。 - semisync：半同步。
	Mode string `json:"mode"`
}

func (o FailoverModeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "FailoverModeRequest struct{}"
	}

	return strings.Join([]string{"FailoverModeRequest", string(data)}, " ")
}
