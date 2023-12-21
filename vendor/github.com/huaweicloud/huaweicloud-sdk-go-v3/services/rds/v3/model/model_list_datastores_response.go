package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListDatastoresResponse Response Object
type ListDatastoresResponse struct {

	// 数据库引擎信息。
	DataStores     *[]LDatastore `json:"dataStores,omitempty"`
	HttpStatusCode int           `json:"-"`
}

func (o ListDatastoresResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListDatastoresResponse struct{}"
	}

	return strings.Join([]string{"ListDatastoresResponse", string(data)}, " ")
}
