package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowAvailableVersionResponse Response Object
type ShowAvailableVersionResponse struct {

	// 可选版本列表。
	AvailableVersions *[]string `json:"available_versions,omitempty"`
	HttpStatusCode    int       `json:"-"`
}

func (o ShowAvailableVersionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAvailableVersionResponse struct{}"
	}

	return strings.Join([]string{"ShowAvailableVersionResponse", string(data)}, " ")
}
