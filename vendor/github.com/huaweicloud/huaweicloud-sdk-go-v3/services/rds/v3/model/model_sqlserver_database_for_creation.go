package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SqlserverDatabaseForCreation 数据库信息。
type SqlserverDatabaseForCreation struct {

	// 数据库名称。 数据库名称长度可在1～64个字符之间，由字母、数字、中划线或下划线组成，不能包含其他特殊字符，且不能以RDS for SQL Server系统库开头或结尾。 RDS for SQL Server系统库包括master，msdb，model，tempdb，resource以及rdsadmin。
	Name string `json:"name"`
}

func (o SqlserverDatabaseForCreation) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SqlserverDatabaseForCreation struct{}"
	}

	return strings.Join([]string{"SqlserverDatabaseForCreation", string(data)}, " ")
}
