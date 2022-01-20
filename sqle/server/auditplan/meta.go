package auditplan

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/pkg/params"
)

type Meta struct {
	Type         string        `json:"audit_plan_type"`
	Desc         string        `json:"audit_plan_type_desc"`
	InstanceType string        `json:"instance_type"`
	Params       params.Params `json:"audit_plan_params,omitempty"`
}

const (
	TypeMySQLSlowLog = "mysql_slow_log"
	TypeMySQLMybatis = "mysql_mybatis"
	TypeDefault      = "default"
)

const (
	InstanceTypeAll   = ""
	InstanceTypeMySQL = "mysql"
)

var Metas = []Meta{
	{
		Type:         TypeDefault,
		Desc:         "自定义",
		InstanceType: InstanceTypeAll,
	},
	{
		Type:         TypeMySQLSlowLog,
		Desc:         "慢日志",
		InstanceType: InstanceTypeMySQL,
	},
	{
		Type:         TypeMySQLMybatis,
		Desc:         "Mybatis 扫描",
		InstanceType: InstanceTypeMySQL,
	},
}

var MetaMap = map[string]Meta{}

func init() {
	for _, meta := range Metas {
		MetaMap[meta.Type] = meta
	}
}

func GetMeta(typ string) (Meta, error) {
	if typ == "" {
		typ = TypeDefault
	}
	meta, ok := MetaMap[typ]
	if !ok {
		return Meta{}, fmt.Errorf("audit plan type %s not found", typ)
	}
	return Meta{
		Type:         meta.Type,
		Desc:         meta.Desc,
		InstanceType: meta.InstanceType,
		Params:       meta.Params.Copy(),
	}, nil
}
