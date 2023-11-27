package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// TagWithKeyValue 键值对标签。
type TagWithKeyValue struct {

	// 标签键。长度为1-128个unicode字符。 可以包含任何语种字母、数字、空格和_.:=+-@， 但首尾不能含有空格，不能以sys开头。
	Key string `json:"key"`

	// 标签值。最大长度255个unicode字符。 可以为空字符串。可以包含任何语种字母、数字、空格和_.:=+-@，
	Value string `json:"value"`
}

func (o TagWithKeyValue) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TagWithKeyValue struct{}"
	}

	return strings.Join([]string{"TagWithKeyValue", string(data)}, " ")
}
