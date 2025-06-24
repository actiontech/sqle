package constant

import (
	"fmt"

	dmsCommonV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
)

// internel build-in uid
const (
	UIDOfOpPermissionCreateProject          = "700001"
	UIDOfOpPermissionProjectAdmin           = "700002"
	UIDOfOpPermissionCreateWorkflow         = "700003"
	UIDOfOpPermissionAuditWorkflow          = "700004"
	UIDOfOpPermissionAuthDBServiceData      = "700005"
	UIDOfOpPermissionExecuteWorkflow        = "700006"
	UIDOfOpPermissionViewOthersWorkflow     = "700007"
	UIDOfOpPermissionViewOthersAuditPlan    = "700008"
	UIDOfOpPermissionSaveAuditPlan          = "700009"
	UIDOfOpPermissionSQLQuery               = "700010"
	UIDOfOpPermissionExportApprovalReject   = "700011"
	UIDOfOpPermissionExportCreate           = "700012"
	UIDOfOpPermissionCreateOptimization     = "700013"
	UIDOfOpPermissionViewOthersOptimization = "700014"
	UIDOfOpPermissionCreatePipeline         = "700015"
	// UIDOfOpPermissionGlobalView 可以查看全局资源,但是不能修改资源
	UIDOfOpPermissionGlobalView = "700016"
	// UIDOfOpPermissionGlobalManagement 可以操作和查看全局资源,但是权限级别低于admin,admin可以修改全局资源权限,全局资源权限不能修改admin
	// 拥有全局资源权限用户不能同级权限用户
	UIDOfOpPermissionGlobalManagement = "700017"
	UIDOfOrdinaryUser                 = "700018"
	UIDOfOpPermissionViewOperationRecord  = "700019"
	UIDOfOpPermissionViewExportTask = "700020"
	UIDOfPermissionViewQuickAuditRecord = "700021"
	UIDOfOpPermissionViewIDEAuditRecord = "700022"
	UIDOfOpPermissionViewOptimizationRecord = "700023"
	UIDOfOpPermissionViewVersionManage = "700024"
	UIDOfOpPermissionVersionManage = "700025"
	UIdOfOpPermissionViewPipeline  = "700026"
	UIdOfOpPermissionManageProjectDataSource  = "700028"
	UIdOfOpPermissionManageAuditRuleTemplate = "700029"
	UIdOfOpPermissionManageApprovalTemplate  = "700030"
	UIdOfOpPermissionManageMember = "700031"
	UIdOfOpPermissionPushRule = "700032"
	UIdOfOpPermissionMangeAuditSQLWhiteList = "700033"
	UIdOfOpPermissionManageSQLMangeWhiteList = "700034"
	UIdOfOpPermissionManageRoleMange  = "700035"
	UIdOfOpPermissionDesensitization = "700036"

	UIDOfDMSConfig = "700100"

	UIDOfUserAdmin = "700200"
	UIDOfUserSys   = "700201"

	UIDOfProjectDefault = "700300"

	UIDOfRoleProjectAdmin = "700400"
	UIDOfRoleDevEngineer  = "700403"
	UIDOfRoleDevManager   = "700404"
	UIDOfRoleOpsEngineer  = "700405"
)

func ConvertPermissionIdToType(opPermissionUid string) (apiOpPermissionTyp dmsCommonV1.OpPermissionType, err error) {
	switch opPermissionUid {
	case UIDOfOpPermissionCreateWorkflow:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionTypeCreateWorkflow
	case UIDOfOpPermissionAuditWorkflow:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionTypeAuditWorkflow
	case UIDOfOpPermissionAuthDBServiceData:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionTypeAuthDBServiceData
	case UIDOfOpPermissionProjectAdmin:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionTypeProjectAdmin
	case UIDOfOpPermissionCreateProject:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionTypeCreateProject
	case UIDOfOpPermissionExecuteWorkflow:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionTypeExecuteWorkflow
	case UIDOfOpPermissionViewOthersWorkflow:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionTypeViewOthersWorkflow
	case UIDOfOpPermissionSaveAuditPlan:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionTypeSaveAuditPlan
	case UIDOfOpPermissionViewOthersAuditPlan:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionTypeViewOtherAuditPlan
	case UIDOfOpPermissionSQLQuery:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionTypeSQLQuery
	case UIDOfOpPermissionExportApprovalReject:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionTypeAuditExportWorkflow
	case UIDOfOpPermissionExportCreate:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionTypeExportCreate
	case UIDOfOpPermissionCreatePipeline:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionTypeCreatePipeline
	case UIDOfOpPermissionViewOperationRecord:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionViewOperationRecord
	case UIDOfOpPermissionViewExportTask:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionViewExportTask
	case UIDOfPermissionViewQuickAuditRecord:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionViewQuickAuditRecord
	case UIDOfOpPermissionViewIDEAuditRecord:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionViewIDEAuditRecord
	case UIDOfOpPermissionViewOptimizationRecord:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionViewOptimizationRecord
	case UIDOfOpPermissionViewVersionManage:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionViewVersionManage
	case UIDOfOpPermissionVersionManage:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionVersionManage
	case UIdOfOpPermissionViewPipeline:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionViewPipeline
	case UIdOfOpPermissionManageProjectDataSource:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionManageProjectDataSource
	case UIdOfOpPermissionManageAuditRuleTemplate:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionManageAuditRuleTemplate
	case UIdOfOpPermissionManageApprovalTemplate:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionManageApprovalTemplate
	case UIdOfOpPermissionManageMember:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionManageMember
	case UIdOfOpPermissionPushRule:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionPushRule
	case UIdOfOpPermissionMangeAuditSQLWhiteList:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionMangeAuditSQLWhiteList
	case UIdOfOpPermissionManageSQLMangeWhiteList:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionManageSQLMangeWhiteList
	case UIdOfOpPermissionManageRoleMange:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionManageRoleMange
	case UIdOfOpPermissionDesensitization:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionDesensitization
	case UIDOfOrdinaryUser:
		apiOpPermissionTyp = dmsCommonV1.OpPermissionTypeNone
	default:
		return dmsCommonV1.OpPermissionTypeUnknown, fmt.Errorf("get user op permission type error: invalid op permission uid: %v", opPermissionUid)

	}
	return apiOpPermissionTyp, nil
}

func ConvertPermissionTypeToId(opPermissionType dmsCommonV1.OpPermissionType) (permissionId string, err error) {
	switch opPermissionType {
	case dmsCommonV1.OpPermissionTypeCreateWorkflow:
		permissionId = UIDOfOpPermissionCreateWorkflow
	case dmsCommonV1.OpPermissionTypeAuditWorkflow:
		permissionId = UIDOfOpPermissionAuditWorkflow
	case dmsCommonV1.OpPermissionTypeAuthDBServiceData:
		permissionId = UIDOfOpPermissionAuthDBServiceData
	case dmsCommonV1.OpPermissionTypeProjectAdmin:
		permissionId = UIDOfOpPermissionProjectAdmin
	case dmsCommonV1.OpPermissionTypeCreateProject:
		permissionId = UIDOfOpPermissionCreateProject
	case dmsCommonV1.OpPermissionTypeGlobalManagement:
		permissionId = UIDOfOpPermissionGlobalManagement
	case dmsCommonV1.OpPermissionTypeGlobalView:
		permissionId = UIDOfOpPermissionGlobalView
	case dmsCommonV1.OpPermissionTypeExecuteWorkflow:
		permissionId = UIDOfOpPermissionExecuteWorkflow
	case dmsCommonV1.OpPermissionTypeViewOthersWorkflow:
		permissionId = UIDOfOpPermissionViewOthersWorkflow
	case dmsCommonV1.OpPermissionTypeSaveAuditPlan:
		permissionId = UIDOfOpPermissionSaveAuditPlan
	case dmsCommonV1.OpPermissionTypeViewOtherAuditPlan:
		permissionId = UIDOfOpPermissionViewOthersAuditPlan
	case dmsCommonV1.OpPermissionTypeSQLQuery:
		permissionId = UIDOfOpPermissionSQLQuery
	case dmsCommonV1.OpPermissionTypeAuditExportWorkflow:
		permissionId = UIDOfOpPermissionExportApprovalReject
	case dmsCommonV1.OpPermissionTypeExportCreate:
		permissionId = UIDOfOpPermissionExportCreate
	case dmsCommonV1.OpPermissionTypeCreatePipeline:
		permissionId = UIDOfOpPermissionCreatePipeline
	case dmsCommonV1.OpPermissionViewOperationRecord:
		permissionId = UIDOfOpPermissionViewOperationRecord
	case dmsCommonV1.OpPermissionViewExportTask:
		permissionId = UIDOfOpPermissionViewExportTask
	case dmsCommonV1.OpPermissionViewQuickAuditRecord:
		permissionId = UIDOfPermissionViewQuickAuditRecord
	case dmsCommonV1.OpPermissionViewIDEAuditRecord:
		permissionId = UIDOfOpPermissionViewIDEAuditRecord
	case dmsCommonV1.OpPermissionViewOptimizationRecord:
		permissionId = UIDOfOpPermissionViewOptimizationRecord
	case dmsCommonV1.OpPermissionViewVersionManage:
		permissionId = UIDOfOpPermissionViewVersionManage
	case dmsCommonV1.OpPermissionVersionManage:
		permissionId = UIDOfOpPermissionVersionManage
	case dmsCommonV1.OpPermissionViewPipeline:
		permissionId = UIdOfOpPermissionViewPipeline
	case dmsCommonV1.OpPermissionManageProjectDataSource:
		permissionId = UIdOfOpPermissionManageProjectDataSource
	case dmsCommonV1.OpPermissionManageAuditRuleTemplate:
		permissionId = UIdOfOpPermissionManageAuditRuleTemplate
	case dmsCommonV1.OpPermissionManageApprovalTemplate:
		permissionId = UIdOfOpPermissionManageApprovalTemplate
	case dmsCommonV1.OpPermissionManageMember:
		permissionId = UIdOfOpPermissionManageMember
	case dmsCommonV1.OpPermissionPushRule:
		permissionId = UIdOfOpPermissionPushRule
	case dmsCommonV1.OpPermissionMangeAuditSQLWhiteList:
		permissionId = UIdOfOpPermissionMangeAuditSQLWhiteList
	case dmsCommonV1.OpPermissionManageSQLMangeWhiteList:
		permissionId = UIdOfOpPermissionManageSQLMangeWhiteList
	case dmsCommonV1.OpPermissionManageRoleMange:
		permissionId = UIdOfOpPermissionManageRoleMange
	case dmsCommonV1.OpPermissionDesensitization:
		permissionId = UIdOfOpPermissionDesensitization
	case dmsCommonV1.OpPermissionTypeNone:
		permissionId = UIDOfOrdinaryUser
	default:
		return "", fmt.Errorf("get user op permission id error: invalid op permission type: %v", opPermissionType)
	}

	return permissionId, nil
}

type DBType string

func (d DBType) String() string {
	return string(d)
}

func ParseDBType(s string) (DBType, error) {
	switch s {
	case "MySQL":
		return DBTypeMySQL, nil
	case "TDSQL For InnoDB":
		return DBTypeTDSQLForInnoDB, nil
	case "TiDB":
		return DBTypeTiDB, nil
	case "PostgreSQL":
		return DBTypePostgreSQL, nil
	case "Oracle":
		return DBTypeOracle, nil
	case "DB2":
		return DBTypeDB2, nil
	case "SQL Server":
		return DBTypeSQLServer, nil
	case "OceanBase For MySQL":
		return DBTypeOceanBaseMySQL, nil
	case "GoldenDB":
		return DBTypeGoldenDB, nil
	case "TBase":
		return DBTypeTBase, nil
	default:
		return "", fmt.Errorf("invalid db type: %s", s)
	}
}

const (
	DBTypeMySQL          DBType = "MySQL"
	DBTypePostgreSQL     DBType = "PostgreSQL"
	DBTypeTiDB           DBType = "TiDB"
	DBTypeSQLServer      DBType = "SQL Server"
	DBTypeOracle         DBType = "Oracle"
	DBTypeDB2            DBType = "DB2"
	DBTypeOceanBaseMySQL DBType = "OceanBase For MySQL"
	DBTypeTDSQLForInnoDB DBType = "TDSQL For InnoDB"
	DBTypeGoldenDB       DBType = "GoldenDB"
	DBTypeTBase          DBType = "TBase"
)

var SupportedDataExportDBTypes = map[DBType]struct{}{
	DBTypeMySQL:      {},
	DBTypePostgreSQL: {},
	DBTypeOracle:     {},
	DBTypeSQLServer:  {},
}

type FilterCondition struct {
	// Filter For Preload Table
	Table         string
	KeywordSearch bool
	Field         string
	Operator      FilterOperator
	Value         interface{}
}

type FilterOperator string

const (
	FilterOperatorEqual              FilterOperator = "="
	FilterOperatorIsNull             FilterOperator = "isNull"
	FilterOperatorNotEqual           FilterOperator = "<>"
	FilterOperatorContains           FilterOperator = "like"
	FilterOperatorNotContains        FilterOperator = "not like"
	FilterOperatorGreaterThanOrEqual FilterOperator = ">="
	FilterOperatorLessThanOrEqual    FilterOperator = "<="
	FilterOperatorIn                 FilterOperator = "in"
)

type DBServiceSourceName string

const (
	DBServiceSourceNameDMP           DBServiceSourceName = "Actiontech DMP"
	DBServiceSourceNameDMS           DBServiceSourceName = "Actiontech DMS"
	DBServiceSourceNameSQLE          DBServiceSourceName = "SQLE"
	DBServiceSourceNameExpandService DBServiceSourceName = "Expand Service"
)

func ParseDBServiceSource(s string) (DBServiceSourceName, error) {
	switch s {
	case string(DBServiceSourceNameDMP):
		return DBServiceSourceNameDMP, nil
	case string(DBServiceSourceNameDMS):
		return DBServiceSourceNameDMS, nil
	case string(DBServiceSourceNameSQLE):
		return DBServiceSourceNameSQLE, nil
	case string(DBServiceSourceNameExpandService):
		return DBServiceSourceNameExpandService, nil
	default:
		return "", fmt.Errorf("invalid data object source name: %s", s)
	}
}

const (
	DMSToken        = "dms-token"
	DMSRefreshToken = "dms-refresh-token"
)
