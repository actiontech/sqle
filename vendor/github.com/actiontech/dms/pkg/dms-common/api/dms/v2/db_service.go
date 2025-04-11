package v2

import (
	base "github.com/actiontech/dms/pkg/dms-common/api/base/v1"
	v1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/go-openapi/strfmt"
)

// swagger:parameters ListDBServicesV2
type ListDBServiceReq struct {
	// the maximum count of db service to be returned
	// in:query
	// Required: true
	PageSize uint32 `query:"page_size" json:"page_size" validate:"required"`
	// the offset of users to be returned, default is 0
	// in:query
	PageIndex uint32 `query:"page_index" json:"page_index"`
	// Multiple of ["name"], default is ["name"]
	// in:query
	OrderBy v1.DBServiceOrderByField `query:"order_by" json:"order_by"`
	// the db service connection
	// enum: connect_success,connect_failed
	// in:query
	FilterLastConnectionTestStatus *string `query:"filter_last_connection_test_status" json:"filter_last_connection_test_status" validate:"omitempty,oneof=connect_success connect_failed"`
	// the db service host
	// in:query
	FilterByHost string `query:"filter_by_host" json:"filter_by_host"`
	// the db service uid
	// in:query
	FilterByUID string `query:"filter_by_uid" json:"filter_by_uid"`
	// the db service name
	// in:query
	FilterByName string `query:"filter_by_name" json:"filter_by_name"`
	// the db service port
	// in:query
	FilterByPort string `query:"filter_by_port" json:"filter_by_port"`
	// the db service db type
	// in:query
	FilterByDBType string `query:"filter_by_db_type" json:"filter_by_db_type"`
	// project id
	// in:path
	ProjectUid string `param:"project_uid" json:"project_uid"`
	// filter db services by db service id list using in condition
	// in:query
	FilterByDBServiceIds []string `query:"filter_by_db_service_ids" json:"filter_by_db_service_ids"`
	// filter db services by environment tag
	// in:query
	FilterByEnvironmentTagUID string `query:"filter_by_environment_tag_uid" json:"filter_by_environment_tag_uid"`
	// the db service fuzzy keyword,include host/port
	// in:query
	FuzzyKeyword string `query:"fuzzy_keyword" json:"fuzzy_keyword"`
	// is masking
	// in:query
	IsEnableMasking *bool `query:"is_enable_masking" json:"is_enable_masking"`
}

// swagger:model ListDBServiceReplyV2
type ListDBServiceReply struct {
	// List db service reply
	Data  []*ListDBService `json:"data"`
	Total int64            `json:"total_nums"`

	// Generic reply
	base.GenericResp
}

// swagger:model ListDBServiceV2
type ListDBService struct {
	// db service uid
	DBServiceUid string `json:"uid"`
	// db service name
	Name string `json:"name"`
	// db service DB type
	DBType string `json:"db_type"`
	// db service host
	Host string `json:"host"`
	// db service port
	Port string `json:"port"`
	// db service admin user
	User string `json:"user"`
	// db service admin encrypted password
	Password string `json:"password"`
	// DB Service environment tag
	EnvironmentTag *v1.EnvironmentTag `json:"environment_tag"`
	// DB Service maintenance time
	MaintenanceTimes []*v1.MaintenanceTime `json:"maintenance_times"`
	// DB desc
	Desc string `json:"desc"`
	// DB source
	Source string `json:"source"`
	// DB project uid
	ProjectUID string `json:"project_uid"`
	// sqle config
	SQLEConfig *v1.SQLEConfig `json:"sqle_config"`
	// DB Service Custom connection parameters
	AdditionalParams []*v1.AdditionalParam `json:"additional_params"`
	// is enable masking
	IsEnableMasking bool `json:"is_enable_masking"`
	// backup switch
	EnableBackup bool `json:"enable_backup"`
	// backup max rows
	BackupMaxRows uint64 `json:"backup_max_rows"`
	// audit plan types
	AuditPlanTypes []*v1.AuditPlanTypes `json:"audit_plan_types"`
	// instance audit plan id
	InstanceAuditPlanID uint `json:"instance_audit_plan_id,omitempty"`
	// DB connection test time
	LastConnectionTestTime strfmt.DateTime `json:"last_connection_test_time"`
	// DB connect test status
	LastConnectionTestStatus v1.LastConnectionTestStatus `json:"last_connection_test_status"`
	// DB connect test error message
	LastConnectionTestErrorMessage string `json:"last_connection_test_error_message,omitempty"`
}