package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ChangeFailoverStrategyResponse Response Object
type ChangeFailoverStrategyResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o ChangeFailoverStrategyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeFailoverStrategyResponse struct{}"
	}

	return strings.Join([]string{"ChangeFailoverStrategyResponse", string(data)}, " ")
}
