package proto

import (
	"github.com/actiontech/dms/pkg/dms-common/i18nPkg"
	"github.com/actiontech/sqle/sqle/pkg/params"
)

func ConvertParamToProtoParam(p params.Params) []*Param {
	pp := make([]*Param, len(p))
	for i, v := range p {
		if v == nil {
			continue
		}
		pp[i] = &Param{
			Key:   v.Key,
			Value: v.Value,
			Desc:  v.GetDesc(i18nPkg.DefaultLang),
			Type:  string(v.Type),
		}
	}
	return pp
}

func ConvertProtoParamToParam(p []*Param) params.Params {
	pp := make(params.Params, len(p))
	for i, v := range p {
		if v == nil {
			continue
		}
		pp[i] = &params.Param{
			Key:   v.Key,
			Value: v.Value,
			Desc:  v.Desc,
			Type:  params.ParamType(v.Type),
		}
	}
	return pp
}
