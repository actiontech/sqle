package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SqlserverDatabaseForDetail 数据库信息。
type SqlserverDatabaseForDetail struct {

	// 数据库名称。
	Name string `json:"name"`

	// 数据库使用的字符集，例如Chinese_PRC_CI_AS等。
	CharacterSet string `json:"character_set"`

	// 数据库状态。取值如下:  Creating:表示创建中。 Running:表示使用中。 Deleting:表示删除中。 NotExists:表示不存在。
	State string `json:"state"`
}

func (o SqlserverDatabaseForDetail) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SqlserverDatabaseForDetail struct{}"
	}

	return strings.Join([]string{"SqlserverDatabaseForDetail", string(data)}, " ")
}
