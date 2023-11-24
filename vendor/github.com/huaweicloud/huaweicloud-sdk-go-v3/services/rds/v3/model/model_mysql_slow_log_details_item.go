package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type MysqlSlowLogDetailsItem struct {

	// 执行次数。
	Count *string `json:"count,omitempty"`

	// 执行时间。
	Time *string `json:"time,omitempty"`

	// 等待锁时间。mysql支持
	LockTime *string `json:"lock_time,omitempty"`

	// 结果行数量。mysql支持
	RowsSent *string `json:"rows_sent,omitempty"`

	// 扫描的行数量。mysql支持
	RowsExamined *string `json:"rows_examined,omitempty"`

	// 所属数据库。
	Database *string `json:"database,omitempty"`

	// 帐号。
	Users *string `json:"users,omitempty"`

	// 执行语法。慢日志默认脱敏显示，如需明文显示，请联系客服人员添加白名单。
	QuerySample *string `json:"query_sample,omitempty"`

	// 语句类型。
	Type *string `json:"type,omitempty"`

	// 发生时间，UTC时间。
	StartTime *string `json:"start_time,omitempty"`

	// IP地址。
	ClientIp *string `json:"client_ip,omitempty"`

	// 日志单行序列号。
	LineNum *string `json:"line_num,omitempty"`
}

func (o MysqlSlowLogDetailsItem) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlSlowLogDetailsItem struct{}"
	}

	return strings.Join([]string{"MysqlSlowLogDetailsItem", string(data)}, " ")
}
