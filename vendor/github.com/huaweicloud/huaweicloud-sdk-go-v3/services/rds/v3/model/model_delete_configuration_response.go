package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteConfigurationResponse Response Object
type DeleteConfigurationResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteConfigurationResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteConfigurationResponse struct{}"
	}

	return strings.Join([]string{"DeleteConfigurationResponse", string(data)}, " ")
}
