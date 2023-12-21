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
	// project id
	// Required: true
	// in:path
	ProjectUid string `param:"project_uid" json:"project_uid" validate:"required"`
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
	// is member project admin, admin has all permissions
	IsAdmin bool `json:"is_admin"`
	// member op permissions
	MemberOpPermissionList []OpPermissionItem `json:"member_op_permission_list"`
}

// swagger:model ListMembersForInternalReply
type ListMembersForInternalReply struct {
	// List member reply
	Data  []*ListMembersForInternalItem `json:"data"`
	Total int64                         `json:"total_nums"`

	// Generic reply
	base.GenericResp
}
