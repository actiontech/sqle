package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateConfigurationResponse Response Object
type UpdateConfigurationResponse struct {
	Configuration  *UpdateConfigurationRspConfiguration `json:"configuration,omitempty"`
	HttpStatusCode int                                  `json:"-"`
}

func (o UpdateConfigurationResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateConfigurationResponse struct{}"
	}

	return strings.Join([]string{"UpdateConfigurationResponse", string(data)}, " ")
}
