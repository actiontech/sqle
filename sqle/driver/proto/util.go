package proto

import (
	"github.com/actiontech/sqle/sqle/pkg/params"
)

func ParamToProtoParam(p []*params.Param) []*Param {
	pp := make([]*Param, len(p))
	for _, param := range p {
		pp = append(pp, &Param{
			Key:   param.Key,
			Value: param.Value,
			Desc:  param.Desc,
			Type:  string(param.Type),
		})
	}
	return pp
}

func ProtoParamToParam(p []*Param) []*params.Param {
	pp := make([]*params.Param, len(p))
	for _, param := range p {
		pp = append(pp, &params.Param{
			Key:   param.Key,
			Value: param.Value,
			Desc:  param.Desc,
			Type:  params.ParamType(param.Type),
		})
	}
	return pp
}
