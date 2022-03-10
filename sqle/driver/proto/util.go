package proto

import (
	"github.com/actiontech/sqle/sqle/pkg/params"
)

func ParamToProtoParam(p []*params.Param) []*Param {
	pp := make([]*Param, len(p))
	for i, v := range p {
		pp[i] = &Param{
			Key:   v.Key,
			Value: v.Value,
			Desc:  v.Desc,
			Type:  string(v.Type),
		}
	}
	return pp
}

func ProtoParamToParam(p []*Param) []*params.Param {
	pp := make([]*params.Param, len(p))
	for i, v := range p {
		pp[i] = &params.Param{
			Key:   v.Key,
			Value: v.Value,
			Desc:  v.Desc,
			Type:  params.ParamType(v.Type),
		}
	}
	return pp
}
