//go:build enterprise
// +build enterprise

package auditplan

import (
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/pkg/params"
)

type TDMySQLSchemaMetaTaskV2 struct {
	MySQLSchemaMetaTaskV2
}

func NewTDMySQLSchemaMetaTaskV2Fn() func() interface{} {
	return func() interface{} {
		return NewTDMySQLSchemaMetaTaskV2()
	}
}

func NewTDMySQLSchemaMetaTaskV2() *TDMySQLSchemaMetaTaskV2 {
	t := &TDMySQLSchemaMetaTaskV2{}
	t.MySQLSchemaMetaTaskV2 = *NewMySQLSchemaMetaTaskV2()
	return t
}

func (at *TDMySQLSchemaMetaTaskV2) InstanceType() string {
	return InstanceTypeTDSQL
}

func (at *TDMySQLSchemaMetaTaskV2) Params(instanceId ...string) params.Params {
	return []*params.Param{
		{
			Key:      paramKeyCollectIntervalMinute,
			Value:    "60",
			Type:     params.ParamTypeInt,
			I18nDesc: locale.ShouldLocalizeAll(locale.ParamCollectIntervalMinuteOracle),
		},
	}
}
