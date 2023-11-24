package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListSqlserverDatabasesRequest Request Object
type ListSqlserverDatabasesRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 分页页码，从1开始。
	Page int32 `json:"page"`

	// 每页数据条数。取值范围[1, 100]。
	Limit int32 `json:"limit"`

	// 数据库名。当指定该参数时，page和limit参数需要传入但不生效。
	DbName *string `json:"db-name,omitempty"`

	// 数据库恢复健康模式，取值：FULL  ：完整模式，SIMPLE  ：简单模式，BUlK_LOGGED ：大容量日志恢复模式（该参数仅用于SQL server引擎）
	RecoverModel *string `json:"recover_model,omitempty"`
}

func (o ListSqlserverDatabasesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSqlserverDatabasesRequest struct{}"
	}

	return strings.Join([]string{"ListSqlserverDatabasesRequest", string(data)}, " ")
}
