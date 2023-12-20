package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// TagDelWithKeyValue 键值对标签。
type TagDelWithKeyValue struct {

	// 标签键。最大长度127个unicode字符。 key不能为空，不能为空字符串。
	Key string `json:"key"`

	// 标签值。每个值最大长度255个unicode字符。 删除说明如下： - 如果“value”有值，按照“key”/“value”删除。 - 如果“value”没值，则按照“key”删除。
	Value *string `json:"value,omitempty"`
}

func (o TagDelWithKeyValue) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TagDelWithKeyValue struct{}"
	}

	return strings.Join([]string{"TagDelWithKeyValue", string(data)}, " ")
}
