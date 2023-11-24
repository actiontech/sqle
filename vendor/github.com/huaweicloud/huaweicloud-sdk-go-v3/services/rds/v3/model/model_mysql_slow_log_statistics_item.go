package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type MysqlSlowLogStatisticsItem struct {

	// 执行次数。
	Count *string `json:"count,omitempty"`

	// 执行时间。
	Time *string `json:"time,omitempty"`

	// 等待锁时间。mysql支持
	LockTime *string `json:"lock_time,omitempty"`

	// 结果行数量。mysql支持
	RowsSent *int64 `json:"rows_sent,omitempty"`

	// 扫描的行数量。mysql支持
	RowsExamined *int64 `json:"rows_examined,omitempty"`

	// 所属数据库。
	Database *string `json:"database,omitempty"`

	// 帐号。
	Users *string `json:"users,omitempty"`

	// 执行语法。
	QuerySample *string `json:"query_sample,omitempty"`

	// IP地址。
	ClientIp *string `json:"client_ip,omitempty"`

	// 语句类型。
	Type *string `json:"type,omitempty"`
}

func (o MysqlSlowLogStatisticsItem) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlSlowLogStatisticsItem struct{}"
	}

	return strings.Join([]string{"MysqlSlowLogStatisticsItem", string(data)}, " ")
}
