package v1

import (
	base "github.com/actiontech/dms/pkg/dms-common/api/base/v1"

	"github.com/go-openapi/strfmt"
)

// swagger:parameters ListProjects
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
	OrderBy ProjectOrderByField `query:"order_by" json:"order_by"`
	// filter the Project name
	FilterByName string `query:"filter_by_name" json:"filter_by_name"`
	// filter the Project UID
	FilterByUID string `query:"filter_by_uid" json:"filter_by_uid"`
	// filter project by project id list, using in condition
	// in:query
	FilterByProjectUids []string `query:"filter_by_project_uids" json:"filter_by_project_uids"`
	// filter project by project priority
	// in:query
	FilterByProjectPriority ProjectPriority `query:"filter_by_project_priority" json:"filter_by_project_priority"`
	// filter the Project By Project description
	FilterByDesc string `query:"filter_by_desc" json:"filter_by_desc"`
}

// swagger:enum ProjectOrderByField
type ProjectOrderByField string

const (
	ProjectOrderByName ProjectOrderByField = "name"
)

// swagger:enum ProjectPriority
type ProjectPriority string

const (
	ProjectPriorityHigh    ProjectPriority = "high"
	ProjectPriorityMedium  ProjectPriority = "medium"
	ProjectPriorityLow     ProjectPriority = "low"
	ProjectPriorityUnknown ProjectPriority = "unknown" // 当数据库中数据存在问题时，返回该状态
)

func ToPriorityNum(priority ProjectPriority) uint8 {
	switch priority {
	case ProjectPriorityHigh:
		return 30
	case ProjectPriorityMedium:
		return 20
	case ProjectPriorityLow:
		return 10
	default:
		return 20 // 默认优先级为中
	}
}

func ToPriority(priority uint8) ProjectPriority {
	switch priority {
	case 10:
		return ProjectPriorityLow
	case 20:
		return ProjectPriorityMedium
	case 30:
		return ProjectPriorityHigh
	default:
		return ProjectPriorityUnknown
	}
}

// swagger:model ListProjectV1
type ListProject struct {
	// Project uid
	ProjectUid string `json:"uid"`
	// Project name
	Name string `json:"name"`
	// Project is archived
	Archived bool `json:"archived"`
	// Project desc
	Desc string `json:"desc"`
	// is fixed business
	IsFixedBusiness bool `json:"is_fixed_business"`
	// Project business
	Business []Business `json:"business"`
	// create user
	CreateUser UidWithName `json:"create_user"`
	// create time
	CreateTime strfmt.DateTime `json:"create_time"`
	// project priority
	ProjectPriority ProjectPriority `json:"project_priority" enums:"high,medium,low"`
}

type Business struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	IsUsed bool   `json:"is_used"`
}

// swagger:model ListProjectReply
type ListProjectReply struct {
	// List project reply
	Data  []*ListProject `json:"data"`
	Total int64          `json:"total_nums"`

	// Generic reply
	base.GenericResp
}
