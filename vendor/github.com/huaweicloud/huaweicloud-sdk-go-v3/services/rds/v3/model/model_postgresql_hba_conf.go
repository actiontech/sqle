package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PostgresqlHbaConf struct {

	// 连接类型，枚举，host、hostssl、hostnossl
	Type string `json:"type"`

	// 数据库名，除template0，template1的数据库名，多个以逗号隔开
	Database string `json:"database"`

	// 用户名，all，除内置用户（rdsAdmin, rdsMetric, rdsBackup, rdsRepl, rdsProxy）以外，多个以逗号隔开
	User string `json:"user"`

	// 客户端IP地址。0.0.0.0/0表示允许用户从任意IP地址访问数据库
	Address string `json:"address"`

	// 掩码，默认为空字符串
	Mask *string `json:"mask,omitempty"`

	// 认证方式。枚举：reject、md5、scram-sha-256
	Method string `json:"method"`

	// 优先级，表示配置的先后
	Priority int32 `json:"priority"`
}

func (o PostgresqlHbaConf) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostgresqlHbaConf struct{}"
	}

	return strings.Join([]string{"PostgresqlHbaConf", string(data)}, " ")
}
