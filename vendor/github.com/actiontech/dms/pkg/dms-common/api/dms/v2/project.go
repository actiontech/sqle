package v2

import (
	base "github.com/actiontech/dms/pkg/dms-common/api/base/v1"

	"github.com/go-openapi/strfmt"

	v1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
)

// swagger:parameters ListProjectsV2
type ListProjectReq struct {
	// the maximum count of Project to be returned
	// in:query
	// Required: true
	PageSize uint32 `query:"page_size" json:"page_size" validate:"required"`
	// the offset of Projects to be returned, default is 0
	// in:query
	PageIndex uint32 `query:"page_index" json:"page_index"`
	// Multiple of ["name"], default is ["name"]
	// in:query
	OrderBy v1.ProjectOrderByField `query:"order_by" json:"order_by"`
	// filter the Project name
	FilterByName string `query:"filter_by_name" json:"filter_by_name"`
	// filter the Project UID
	FilterByUID string `query:"filter_by_uid" json:"filter_by_uid"`
	// filter project by project id list, using in condition
	// in:query
	FilterByProjectUids []string `query:"filter_by_project_uids" json:"filter_by_project_uids"`
	// filter project by project priority
	// in:query
	FilterByProjectPriority v1.ProjectPriority `query:"filter_by_project_priority" json:"filter_by_project_priority"`
	// filter project by business tag
	// in:query
	FilterByBusinessTag string `query:"filter_by_business_tag" json:"filter_by_business_tag"`
	// filter the Project By Project description
	FilterByDesc string `query:"filter_by_desc" json:"filter_by_desc"`
}

// swagger:model ListProjectV2
type ListProject struct {
	// Project uid
	ProjectUid string `json:"uid"`
	// Project name
	Name string `json:"name"`
	// Project is archived
	Archived bool `json:"archived"`
	// Project desc
	Desc string `json:"desc"`
	// project business tag
	BusinessTag *BusinessTag `json:"business_tag"`
	// create user
	CreateUser v1.UidWithName `json:"create_user"`
	// create time
	CreateTime strfmt.DateTime `json:"create_time"`
	// project priority
	ProjectPriority v1.ProjectPriority `json:"project_priority" enums:"high,medium,low"`
}

// swagger:model BusinessTagCommon
type BusinessTag struct {
	UID string `json:"uid,omitempty"`
	// 业务标签最多50个字符
	Name string `json:"name" validate:"max=50"`
}

// swagger:model ListProjectReplyV2
type ListProjectReply struct {
	// List project reply
	Data  []*ListProject `json:"data"`
	Total int64          `json:"total_nums"`

	// Generic reply
	base.GenericResp
}
