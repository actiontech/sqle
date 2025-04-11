package v1

import (
	"fmt"

	base "github.com/actiontech/dms/pkg/dms-common/api/base/v1"
	"github.com/actiontech/dms/pkg/params"
)

type IPluginDBService struct {
	Name   string
	DBType string
	Host   string
	Port   string
	User   string
	// Business             string
	EnvironmentTag       *EnvironmentTag
	SQLERuleTemplateName string
	SQLERuleTemplateId   string
	AdditionalParams     params.Params `json:"additional_params" from:"additional_params"`
}

type IPluginProject struct {
	// Project name
	Name string `json:"name"`
	// Project is archived
	Archived bool `json:"archived"`
	// Project desc
	Desc string `json:"desc"`
}

type Plugin struct {
	// 插件名称
	Name string `json:"name" validate:"required"`
	// 操作资源处理接口地址,如果为空表示没有检查, eg: http://127.0.0.1:7602/v1/auth/data_resource_operate/handle
	// 该地址目的是统一调用其他服务 数据资源变更前后校验/更新数据的 接口
	// eg: 删除数据源前：
	// 需要sqle服务中实现接口逻辑，判断该数据源上已经没有进行中的工单
	OperateDataResourceHandleUrl string `json:"operate_data_resource_handle_url"`
	GetDatabaseDriverOptionsUrl  string `json:"get_database_driver_options_url"`
	GetDatabaseDriverLogosUrl    string `json:"get_database_driver_logos_url"`
}

// swagger:model
type RegisterDMSPluginReq struct {
	Plugin *Plugin `json:"plugin" validate:"required"`
}

func (u *RegisterDMSPluginReq) String() string {
	if u == nil {
		return "RegisterDMSPluginReq{nil}"
	}
	return fmt.Sprintf("RegisterDMSPluginReq{Name:%s}", u.Plugin.Name)
}

// swagger:model RegisterDMSPluginReply
type RegisterDMSPluginReply struct {
	// Generic reply
	base.GenericResp
}

// swagger:enum DataResourceType
type DataResourceType string

const (
	DataResourceTypeDBService DataResourceType = "db_service"
	DataResourceTypeProject   DataResourceType = "project"
	DataResourceTypeUser      DataResourceType = "user"
	DataResourceTypeUserGroup DataResourceType = "user_group"
)

// swagger:enum OperationType
type OperationType string

const (
	OperationTypeCreate OperationType = "create"
	OperationTypeUpdate OperationType = "update"
	OperationTypeDelete OperationType = "delete"
)

// swagger:enum OperationTimingType
type OperationTimingType string

const (
	OperationTimingTypeBefore OperationTimingType = "before"
	OperationTimingTypeAfter  OperationTimingType = "after"
)

// swagger:parameters OperateDataResourceHandle
type OperateDataResourceHandleReq struct {
	DataResourceUid  string              `json:"data_resource_uid"`
	DataResourceType DataResourceType    `json:"data_resource_type"`
	OperationType    OperationType       `json:"operation_type"`
	OperationTiming  OperationTimingType `json:"operation_timing"`
	ExtraParams      string              `json:"extra_params"`
}

// swagger:model OperateDataResourceHandleReply
type OperateDataResourceHandleReply struct {
	// Generic reply
	base.GenericResp
}
