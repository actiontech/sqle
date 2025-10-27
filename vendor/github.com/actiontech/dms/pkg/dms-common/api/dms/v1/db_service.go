package v1

import (
	base "github.com/actiontech/dms/pkg/dms-common/api/base/v1"
	"github.com/go-openapi/strfmt"
)

type CheckDbConnectable struct {
	// DB Service type
	// Required: true
	// example: MySQL
	DBType string `json:"db_type"  example:"mysql" validate:"required"`
	// DB Service admin user
	// Required: true
	// example: root
	User string `json:"user"  example:"root" valid:"required"`
	// DB Service host
	// Required: true
	// example: 127.0.0.1
	Host string `json:"host"  example:"10.10.10.10" valid:"required,ip_addr|uri|hostname|hostname_rfc1123"`
	// DB Service port
	// Required: true
	// example: 3306
	Port string `json:"port"  example:"3306" valid:"required,port"`
	// DB Service admin password
	// Required: true
	// example: 123456
	Password string `json:"password"  example:"123456"`
	// DB Service Custom connection parameters
	// Required: false
	AdditionalParams []*AdditionalParam `json:"additional_params" from:"additional_params"`
}

type AdditionalParam struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description" example:"参数项中文名" form:"description"`
	Type        string `json:"type" example:"int" form:"type"`
}

// swagger:parameters ListDBServices
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
	OrderBy DBServiceOrderByField `query:"order_by" json:"order_by"`
	// the db service business name
	// in:query
	FilterByBusiness string `query:"filter_by_business" json:"filter_by_business"`
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
	// the db service fuzzy keyword,include host/port
	// in:query
	FuzzyKeyword string `query:"fuzzy_keyword" json:"fuzzy_keyword"`
	// is masking
	// in:query
	IsEnableMasking *bool `query:"is_enable_masking" json:"is_enable_masking"`
}

// swagger:enum DBServiceOrderByField
type DBServiceOrderByField string

const (
	DBServiceOrderByName DBServiceOrderByField = "name"
)

type MaintenanceTime struct {
	MaintenanceStartTime *Time `json:"maintenance_start_time"`
	MaintenanceStopTime  *Time `json:"maintenance_stop_time"`
}

type Time struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
}

type AuditPlanTypes struct {
	AuditPlanId       uint   `json:"audit_plan_id"`
	AuditPlanType     string `json:"type"`
	AuditPlanTypeDesc string `json:"desc"`
}

// swagger:enum LastConnectionTestStatus
type LastConnectionTestStatus string

const (
	LastConnectionTestStatusSuccess LastConnectionTestStatus = "connect_success"
	LastConnectionTestStatusFailed  LastConnectionTestStatus = "connect_failed"
)

// A dms db Service
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
	// TODO This parameter is deprecated and will be removed soon.
	// the db service business name
	Business string `json:"business"`
	// DB Service maintenance time
	MaintenanceTimes []*MaintenanceTime `json:"maintenance_times"`
	// DB desc
	Desc string `json:"desc"`
	// DB source
	Source string `json:"source"`
	// DB project uid
	ProjectUID string `json:"project_uid"`
	// sqle config
	SQLEConfig *SQLEConfig `json:"sqle_config"`
	// DB Service Custom connection parameters
	AdditionalParams []*AdditionalParam `json:"additional_params"`
	// is enable masking
	IsEnableMasking bool `json:"is_enable_masking"`
	// backup switch
	EnableBackup bool `json:"enable_backup"`
	// backup max rows
	BackupMaxRows uint64 `json:"backup_max_rows"`
	// audit plan types
	AuditPlanTypes []*AuditPlanTypes `json:"audit_plan_types"`
	// instance audit plan id
	InstanceAuditPlanID uint `json:"instance_audit_plan_id,omitempty"`
	// DB connection test time
	LastConnectionTestTime strfmt.DateTime `json:"last_connection_test_time"`
	// DB connect test status
	LastConnectionTestStatus LastConnectionTestStatus `json:"last_connection_test_status"`
	// DB connect test error message
	LastConnectionTestErrorMessage string `json:"last_connection_test_error_message,omitempty"`
}

// swagger:model
type EnvironmentTag struct {
	UID string `json:"uid,omitempty"`
	// 环境属性标签最多50个字符
	Name string `json:"name" validate:"max=50"`
}

type SQLEConfig struct {
	// DB Service audit enabled
	AuditEnabled bool `json:"audit_enabled" example:"false"`
	// DB Service rule template name
	RuleTemplateName string `json:"rule_template_name"`
	// DB Service rule template id
	RuleTemplateID string `json:"rule_template_id"`
	// DB Service data export rule template name
	DataExportRuleTemplateName string `json:"data_export_rule_template_name"`
	// DB Service data export rule template id
	DataExportRuleTemplateID string `json:"data_export_rule_template_id"`
	// DB Service SQL query config
	SQLQueryConfig *SQLQueryConfig `json:"sql_query_config"`
}

// swagger:enum SQLAllowQueryAuditLevel
type SQLAllowQueryAuditLevel string

const (
	AuditLevelNormal SQLAllowQueryAuditLevel = "normal"
	AuditLevelNotice SQLAllowQueryAuditLevel = "notice"
	AuditLevelWarn   SQLAllowQueryAuditLevel = "warn"
	AuditLevelError  SQLAllowQueryAuditLevel = "error"
)

type SQLQueryConfig struct {
	MaxPreQueryRows                  int                     `json:"max_pre_query_rows" example:"100"`
	QueryTimeoutSecond               int                     `json:"query_timeout_second" example:"10"`
	AuditEnabled                     bool                    `json:"audit_enabled" example:"false"`
	AllowQueryWhenLessThanAuditLevel SQLAllowQueryAuditLevel `json:"allow_query_when_less_than_audit_level" enums:"normal,notice,warn,error" valid:"omitempty,oneof=normal notice warn error " example:"error"`
	RuleTemplateName                 string                  `json:"rule_template_name"`
	RuleTemplateID                   string                  `json:"rule_template_id"`
}

// swagger:model ListDBServiceReply
type ListDBServiceReply struct {
	// List db service reply
	Data  []*ListDBService `json:"data"`
	Total int64            `json:"total_nums"`

	// Generic reply
	base.GenericResp
}

type DBServiceUidWithNameInfo struct {
	DBServiceUid  string
	DBServiceName string
}
