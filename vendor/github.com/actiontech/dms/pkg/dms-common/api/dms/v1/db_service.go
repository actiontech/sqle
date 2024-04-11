package v1

import (
	base "github.com/actiontech/dms/pkg/dms-common/api/base/v1"
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
	// the db service fuzzy keyword,include host/port
	// in:query
	FuzzyKeyword string `query:"fuzzy_keyword" json:"fuzzy_keyword"`
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
}

type SQLEConfig struct {
	// DB Service rule template name
	RuleTemplateName string `json:"rule_template_name"`
	// DB Service rule template id
	RuleTemplateID string `json:"rule_template_id"`
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
}

// swagger:model ListDBServiceReply
type ListDBServiceReply struct {
	// List db service reply
	Data  []*ListDBService `json:"data"`
	Total int64            `json:"total_nums"`

	// Generic reply
	base.GenericResp
}
