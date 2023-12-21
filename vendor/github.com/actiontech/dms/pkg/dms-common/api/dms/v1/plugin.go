package v1

import (
	"fmt"

	base "github.com/actiontech/dms/pkg/dms-common/api/base/v1"
)

// swagger:enum IPluginDBType
type IPluginDBType string

const (
	IPluginDBTypeDBTypeMySQL          IPluginDBType = "MySQL"
	IPluginDBTypeDBTypeOceanBaseMySQL IPluginDBType = "OceanBaseMySQL"
)

func ParseIPluginDBType(s string) (IPluginDBType, error) {
	switch s {
	case string(IPluginDBTypeDBTypeMySQL):
		return IPluginDBTypeDBTypeMySQL, nil
	case string(IPluginDBTypeDBTypeOceanBaseMySQL):
		return IPluginDBTypeDBTypeOceanBaseMySQL, nil
	default:
		return "", fmt.Errorf("invalid db type: %s", s)
	}
}

type IPluginDBService struct {
	Name                 string
	DBType               string
	Host                 string
	Port                 string
	User                 string
	Business             string
	SQLERuleTemplateName string
	SQLERuleTemplateId   string
	// TODO: more
}

type Plugin struct {
	// 插件名称
	Name string `json:"name" validate:"required"`
	// 添加数据源预检查接口地址, 如果为空表示没有检查, eg: http://127.0.0.1:7602/v1/auth/services/precheck/add
	AddDBServicePreCheckUrl string `json:"add_db_service_pre_check_url"`
	// 删除数据源预检查接口地址, 如果为空表示没有检查, eg: http://127.0.0.1:7602/v1/auth/services/precheck/del
	DelDBServicePreCheckUrl string `json:"del_db_service_pre_check_url"`
	// 删除用户预检查接口地址,如果为空表示没有检查, eg: http://127.0.0.1:7602/v1/auth/users/precheck/del
	DelUserPreCheckUrl string `json:"del_user_pre_check_url"`
	// 删除用户组预检查接口地址,如果为空表示没有检查, eg: http://127.0.0.1:7602/v1/auth/usergroups/precheck/del
	DelUserGroupPreCheckUrl string `json:"del_user_group_pre_check_url"`
	// 操作资源处理接口地址,如果为空表示没有检查, eg: http://127.0.0.1:7602/v1/auth/data_resource_operate/handle
	OperateDataResourceHandleUrl string `json:"operate_data_resource_handle_url"`
}

// swagger:parameters RegisterDMSPlugin
type RegisterDMSPluginReq struct {
	// Register dms plugin
	// in:body
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

// swagger:parameters AddDBServicePreCheck
type AddDBServicePreCheckReq struct {
	// Check if dms can add db service
	// in:body
	DBService *IPluginDBService `json:"db_service" validate:"required"`
}

func (u *AddDBServicePreCheckReq) String() string {
	if u == nil {
		return "AddDBServicePreCheckReq{nil}"
	}
	return fmt.Sprintf("AddDBServicePreCheckReq{Name:%s,DBType:%s Host:%s}", u.DBService.Name, u.DBService.DBType, u.DBService.Host)
}

// swagger:model AddDBServicePreCheckReply
type AddDBServicePreCheckReply struct {
	// Generic reply
	base.GenericResp
}

// swagger:parameters DelDBServicePreCheck
type DelDBServicePreCheckReq struct {
	// Check if dms can del db service
	// in:body
	DBServiceUid string `json:"db_service_uid" validate:"required"`
}

func (u *DelDBServicePreCheckReq) String() string {
	if u == nil {
		return "DelDBServicePreCheckReq{nil}"
	}
	return fmt.Sprintf("DelDBServicePreCheckReq{Uid:%s}", u.DBServiceUid)
}

// swagger:model DelDBServicePreCheckReply
type DelDBServicePreCheckReply struct {
	// Generic reply
	base.GenericResp
}

// swagger:parameters DelUserPreCheck
type DelUserPreCheckReq struct {
	// Check if dms can del db service
	// in:body
	UserUid string `json:"user_uid" validate:"required"`
}

func (u *DelUserPreCheckReq) String() string {
	if u == nil {
		return "DelUserPreCheckReq{nil}"
	}
	return fmt.Sprintf("DelUserPreCheckReq{Uid:%s}", u.UserUid)
}

// swagger:model DelUserPreCheckReply
type DelUserPreCheckReply struct {
	// Generic reply
	base.GenericResp
}

// swagger:parameters DelUserGroupPreCheck
type DelUserGroupPreCheckReq struct {
	// Check if dms can del db service
	// in:body
	UserGroupUid string `json:"user_group_uid" validate:"required"`
}

func (u *DelUserGroupPreCheckReq) String() string {
	if u == nil {
		return "DelUserGroupPreCheckReq{nil}"
	}
	return fmt.Sprintf("DelUserGroupPreCheckReq{Uid:%s}", u.UserGroupUid)
}

// swagger:model DelUserGroupPreCheckReply
type DelUserGroupPreCheckReply struct {
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
	OperationTimingAfter      OperationTimingType = "after"
)

// swagger:parameters OperateDataResourceHandle
type OperateDataResourceHandleReq struct {
	DataResourceUid  string              `json:"data_resource_uid"`
	DataResourceType DataResourceType    `json:"data_resource_type"`
	OperationType    OperationType       `json:"operation_type"`
	OperationTiming  OperationTimingType `json:"operation_timing"`
	// TODO ExtraParams  need extra params for pre check？
}

// swagger:model OperateDataResourceHandleReply
type OperateDataResourceHandleReply struct {
	// Generic reply
	base.GenericResp
}
