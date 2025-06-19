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

// swagger:model GetMemberGroupReply
type GetMemberGroupReply struct {
	// List member reply
	Data *GetMemberGroup `json:"data"`

	// Generic reply
	base.GenericResp
}

type GetMemberGroup struct {
	Name string `json:"name"`
	// member group uid
	Uid string `json:"uid"`
	// member user
	Users []UidWithName `json:"users"`
	// Whether the member has project admin permission
	IsProjectAdmin bool `json:"is_project_admin"`
	// member op permission
	RoleWithOpRanges []ListMemberRoleWithOpRange `json:"role_with_op_ranges"`
}

type ListMemberRoleWithOpRange struct {
	// role uid
	RoleUID UidWithName `json:"role_uid" validate:"required"`
	// op permission range type, only support db service now
	OpRangeType OpRangeType `json:"op_range_type" validate:"required"`
	// op range uids
	RangeUIDs []UidWithName `json:"range_uids" validate:"required"`
}
