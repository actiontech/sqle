package v1

import (
	base "github.com/actiontech/dms/pkg/dms-common/api/base/v1"

	"github.com/go-openapi/strfmt"
)

// swagger:parameters ListNamespaces
type ListNamespaceReq struct {
	// the maximum count of namespace to be returned
	// in:query
	// Required: true
	PageSize uint32 `query:"page_size" json:"page_size" validate:"required"`
	// the offset of namespaces to be returned, default is 0
	// in:query
	PageIndex uint32 `query:"page_index" json:"page_index"`
	// Multiple of ["name"], default is ["name"]
	// in:query
	OrderBy NamespaceOrderByField `query:"order_by" json:"order_by"`
	// filter the namespace name
	FilterByName string `query:"filter_by_name" json:"filter_by_name"`
	// filter the namespace UID
	FilterByUID string `query:"filter_by_uid" json:"filter_by_uid"`
}

// swagger:enum NamespaceOrderByField
type NamespaceOrderByField string

const (
	NamespaceOrderByName NamespaceOrderByField = "name"
)

// A dms namespace
type ListNamespace struct {
	// namespace uid
	NamespaceUid string `json:"uid"`
	// namespace name
	Name string `json:"name"`
	// namespace is archived
	Archived bool `json:"archived"`
	// namespace desc
	Desc string `json:"desc"`
	// create user
	CreateUser UidWithName `json:"create_user"`
	// create time
	CreateTime strfmt.DateTime `json:"create_time"`
}

// swagger:model ListNamespaceReply
type ListNamespaceReply struct {
	// List namespace reply
	Payload struct {
		Namespaces []*ListNamespace `json:"namespaces"`
		Total      int64            `json:"total"`
	} `json:"payload"`

	// Generic reply
	base.GenericResp
}
