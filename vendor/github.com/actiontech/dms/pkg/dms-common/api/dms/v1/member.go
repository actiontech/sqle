package v1

import base "github.com/actiontech/dms/pkg/dms-common/api/base/v1"

// swagger:parameters ListMembersForInternal
type ListMembersForInternalReq struct {
	// the maximum count of member to be returned
	// in:query
	// Required: true
	PageSize uint32 `query:"page_size" json:"page_size" validate:"required"`
	// the offset of members to be returned, default is 0
	// in:query
	PageIndex uint32 `query:"page_index" json:"page_index"`
	// the member namespace uid
	// in:query
	NamespaceUid string `query:"namespace_uid" json:"namespace_uid" validate:"required"`
}

// swagger:enum MemberForInternalOrderByField
type MemberForInternalOrderByField string

const (
	MemberForInternalOrderByUserUid MemberForInternalOrderByField = "user_uid"
)

// A dms member for internal
type ListMembersForInternalItem struct {
	// member user
	User UidWithName `json:"user"`
	// is member namespace admin, admin has all permissions
	IsAdmin bool `json:"is_admin"`
	// member op permissions
	MemberOpPermissionList []OpPermissionItem `json:"member_op_permission_list"`
}

// swagger:model ListMembersForInternalReply
type ListMembersForInternalReply struct {
	// List member reply
	Payload struct {
		Members []*ListMembersForInternalItem `json:"members"`
		Total   int64                         `json:"total"`
	} `json:"payload"`

	// Generic reply
	base.GenericResp
}
