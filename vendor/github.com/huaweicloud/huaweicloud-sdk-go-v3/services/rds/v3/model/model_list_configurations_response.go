package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListConfigurationsResponse Response Object
type ListConfigurationsResponse struct {
	Configurations *[]ConfigurationSummary `json:"configurations,omitempty"`
	HttpStatusCode int                     `json:"-"`
}

func (o ListConfigurationsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListConfigurationsResponse struct{}"
	}

	return strings.Join([]string{"ListConfigurationsResponse", string(data)}, " ")
}
