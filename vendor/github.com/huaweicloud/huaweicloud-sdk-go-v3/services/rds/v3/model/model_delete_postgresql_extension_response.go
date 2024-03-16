package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeletePostgresqlExtensionResponse Response Object
type DeletePostgresqlExtensionResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeletePostgresqlExtensionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeletePostgresqlExtensionResponse struct{}"
	}

	return strings.Join([]string{"DeletePostgresqlExtensionResponse", string(data)}, " ")
}
