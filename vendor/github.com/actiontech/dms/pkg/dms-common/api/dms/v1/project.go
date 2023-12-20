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
}

// swagger:enum ProjectOrderByField
type ProjectOrderByField string

const (
	ProjectOrderByName ProjectOrderByField = "name"
)

// A dms Project
type ListProject struct {
	// Project uid
	ProjectUid string `json:"uid"`
	// Project name
	Name string `json:"name"`
	// Project is archived
	Archived bool `json:"archived"`
	// Project desc
	Desc string `json:"desc"`
	// create user
	CreateUser UidWithName `json:"create_user"`
	// create time
	CreateTime strfmt.DateTime `json:"create_time"`
}

// swagger:model ListProjectReply
type ListProjectReply struct {
	// List project reply
	Data  []*ListProject `json:"data"`
	Total int64          `json:"total_nums"`

	// Generic reply
	base.GenericResp
}
