package v1

import (
	"context"
	"fmt"
	"net/http"

	baseV1 "github.com/actiontech/dms/pkg/dms-common/api/base/v1"
	v1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/common"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/labstack/echo/v4"
)

var ErrInstanceNotExist = errors.New(errors.DataNotExist, fmt.Errorf("instance is not exist"))
var ErrInstanceNoAccess = errors.New(errors.DataNotExist, fmt.Errorf("instance is not exist or you can't access it"))
var errInstanceBind = errors.New(errors.DataExist, fmt.Errorf("an instance can only bind one rule template"))
var ErrWrongTimePeriod = errors.New(errors.DataInvalid, fmt.Errorf("wrong time period"))

type InstanceAdditionalParamResV1 struct {
	Name        string `json:"name" example:"param name" form:"name"`
	Description string `json:"description" example:"参数项中文名" form:"description"`
	Type        string `json:"type" example:"int" form:"type"`
	Value       string `json:"value" example:"0" form:"value"`
}

type InstanceAdditionalParamReqV1 struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type SQLQueryConfigResV1 struct {
	MaxPreQueryRows                  int    `json:"max_pre_query_rows"`
	QueryTimeoutSecond               int    `json:"query_timeout_second"`
	AuditEnabled                     bool   `json:"audit_enabled"`
	AllowQueryWhenLessThanAuditLevel string `json:"allow_query_when_less_than_audit_level"  enums:"normal,notice,warn,error"`
}

type MaintenanceTimeResV1 struct {
	MaintenanceStartTime *TimeResV1 `json:"maintenance_start_time"`
	MaintenanceStopTime  *TimeResV1 `json:"maintenance_stop_time"`
}

type TimeResV1 struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
}

func ConvertPeriodToMaintenanceTimeResV1(mt model.Periods) []*MaintenanceTimeResV1 {
	periods := make([]*MaintenanceTimeResV1, len(mt))
	for i, time := range mt {
		periods[i] = &MaintenanceTimeResV1{
			MaintenanceStartTime: &TimeResV1{
				Hour:   time.StartHour,
				Minute: time.StartMinute,
			},
			MaintenanceStopTime: &TimeResV1{
				Hour:   time.EndHour,
				Minute: time.EndMinute,
			},
		}
	}
	return periods
}

type GetInstanceConnectableResV1 struct {
	controller.BaseRes
	Data InstanceConnectableResV1 `json:"data"`
}

type InstanceConnectableResV1 struct {
	IsInstanceConnectable bool   `json:"is_instance_connectable"`
	ConnectErrorMessage   string `json:"connect_error_message,omitempty"`
}

func newInstanceConnectableResV1(err error) InstanceConnectableResV1 {
	if err == nil {
		return InstanceConnectableResV1{
			IsInstanceConnectable: true,
		}
	}
	return InstanceConnectableResV1{
		IsInstanceConnectable: false,
		ConnectErrorMessage:   err.Error(),
	}
}

// CheckInstanceIsConnectableByName test instance db connection
// @Summary 实例连通性测试（实例提交后）
// @Description test instance db connection
// @Id checkInstanceIsConnectableByNameV1
// @Tags instance
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instance_name path string true "instance name"
// @Success 200 {object} v1.GetInstanceConnectableResV1
// @router /v1/projects/{project_name}/instances/{instance_name}/connection [get]
func CheckInstanceIsConnectableByName(c echo.Context) error {
	instanceName := c.Param("instance_name")
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	instance, exist, err := dms.GetInstanceInProjectByName(c.Request().Context(), projectUid, instanceName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrInstanceNoAccess)
	}

	can, err := CheckCurrentUserCanAccessInstances(c.Request().Context(), projectUid, controller.GetUserID(c), []*model.Instance{instance})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !can {
		return controller.JSONBaseErrorReq(c, ErrInstanceNoAccess)
	}

	l := log.NewEntry()
	err = common.CheckInstanceIsConnectable(instance)
	if err != nil {
		l.Warnf("instance %s is not connectable, err: %s", instanceName, err)
	}

	return c.JSON(http.StatusOK, GetInstanceConnectableResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    newInstanceConnectableResV1(err),
	})
}

type InstanceForCheckConnection struct {
	Name string `json:"name"`
}

type BatchCheckInstanceConnectionsReqV1 struct {
	Instances []InstanceForCheckConnection `json:"instances" valid:"dive,required"`
}

type BatchGetInstanceConnectionsResV1 struct {
	controller.BaseRes
	Data []InstanceConnectionResV1 `json:"data"`
}

type InstanceConnectionResV1 struct {
	InstanceName string `json:"instance_name"`
	InstanceConnectableResV1
}

// BatchCheckInstanceConnections test instance db connection
// @Summary 批量测试实例连通性（实例提交后）
// @Description batch test instance db connections
// @Id batchCheckInstanceIsConnectableByName
// @Tags instance
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instances body v1.BatchCheckInstanceConnectionsReqV1 true "instances"
// @Success 200 {object} v1.BatchGetInstanceConnectionsResV1
// @router /v1/projects/{project_name}/instances/connections [post]
func BatchCheckInstanceConnections(c echo.Context) error {
	req := new(BatchCheckInstanceConnectionsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	instanceNames := make([]string, 0, len(req.Instances))
	for _, instance := range req.Instances {
		instanceNames = append(instanceNames, instance.Name)
	}

	distinctInstNames := utils.RemoveDuplicate(instanceNames)

	instances, err := dms.GetInstancesInProjectByNames(c.Request().Context(), projectUid, distinctInstNames)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if len(distinctInstNames) != len(instances) {
		return controller.JSONBaseErrorReq(c, ErrInstanceNoAccess)
	}

	can, err := CheckCurrentUserCanAccessInstances(c.Request().Context(), projectUid, controller.GetUserID(c), instances)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !can {
		return controller.JSONBaseErrorReq(c, ErrInstanceNoAccess)
	}
	l := log.NewEntry()

	instanceConnectionResV1 := make([]InstanceConnectionResV1, len(instances))
	for i, instance := range instances {
		err := common.CheckInstanceIsConnectable(instance)
		if err != nil {
			l.Warnf("instance %s is not connectable, err: %s", instance.Name, err)
		}
		instanceConnectionResV1[i] = InstanceConnectionResV1{
			InstanceName:             instance.Name,
			InstanceConnectableResV1: newInstanceConnectableResV1(err),
		}
	}

	return c.JSON(http.StatusOK, BatchGetInstanceConnectionsResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    instanceConnectionResV1,
	})
}

type GetInstanceConnectableReqV1 struct {
	DBType           string                          `json:"db_type" form:"db_type" example:"mysql"`
	User             string                          `json:"user" form:"db_user" example:"root" valid:"required"`
	Host             string                          `json:"host" form:"db_host" example:"10.10.10.10" valid:"required,ip_addr|uri|hostname|hostname_rfc1123"`
	Port             string                          `json:"port" form:"db_port" example:"3306" valid:"required,port"`
	Password         string                          `json:"password" form:"db_password" example:"123456"`
	AdditionalParams []*InstanceAdditionalParamReqV1 `json:"additional_params" from:"additional_params"`
}

func CheckInstanceIsConnectable(c echo.Context) error {
	req := new(v1.CheckDbConnectable)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	if req.DBType == "" {
		req.DBType = driverV2.DriverTypeMySQL
	}

	additionalParams := driver.GetPluginManager().AllAdditionalParams()[req.DBType]
	for _, additionalParam := range req.AdditionalParams {
		err := additionalParams.SetParamValue(additionalParam.Name, additionalParam.Value)
		if err != nil {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, err))
		}
	}

	instance := &model.Instance{
		DbType:           req.DBType,
		User:             req.User,
		Host:             req.Host,
		Port:             req.Port,
		Password:         req.Password,
		AdditionalParams: additionalParams,
	}

	l := log.NewEntry()

	err := common.CheckInstanceIsConnectable(instance)
	if err != nil {
		l.Warnf("check instance is connectable failed: %v", err)
		return c.JSON(http.StatusOK, baseV1.GenericResp{Code: http.StatusBadRequest, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, baseV1.GenericResp{Message: "OK"})
}

type GetInstanceSchemaResV1 struct {
	controller.BaseRes
	Data InstanceSchemaResV1 `json:"data"`
}

type InstanceSchemaResV1 struct {
	Schemas []string `json:"schema_name_list"`
}

// GetInstanceSchemas get instance schema list
// @Summary 实例 Schema 列表
// @Description instance schema list
// @Id getInstanceSchemasV1
// @Tags instance
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instance_name path string true "instance name"
// @Success 200 {object} v1.GetInstanceSchemaResV1
// @router /v1/projects/{project_name}/instances/{instance_name}/schemas [get]
func GetInstanceSchemas(c echo.Context) error {
	instanceName := c.Param("instance_name")
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	instance, exist, err := dms.GetInstanceInProjectByName(c.Request().Context(), projectUid, instanceName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrInstanceNoAccess)
	}

	can, err := CheckCurrentUserCanAccessInstances(c.Request().Context(), projectUid, controller.GetUserID(c), []*model.Instance{instance})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !can {
		return controller.JSONBaseErrorReq(c, ErrInstanceNoAccess)
	}

	plugin, err := common.NewDriverManagerWithoutAudit(log.NewEntry(), instance, "")
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	defer plugin.Close(context.TODO())

	schemas, err := plugin.Schemas(context.TODO())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &GetInstanceSchemaResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: InstanceSchemaResV1{
			Schemas: schemas,
		},
	})
}

const ( // InstanceTipReqV1.FunctionalModule Enums
	create_audit_plan = "create_audit_plan"
	create_workflow   = "create_workflow"
)

type InstanceTipReqV1 struct {
	FilterDBType             string `json:"filter_db_type" query:"filter_db_type"`
	FilterWorkflowTemplateId uint32 `json:"filter_workflow_template_id" query:"filter_workflow_template_id"`
	FunctionalModule         string `json:"functional_module" query:"functional_module" enums:"create_audit_plan,create_workflow,sql_manage" valid:"omitempty,oneof=create_audit_plan create_workflow sql_manage"`
}

type InstanceTipResV1 struct {
	ID                 string `json:"instance_id"`
	Name               string `json:"instance_name"`
	Type               string `json:"instance_type"`
	WorkflowTemplateId uint32 `json:"workflow_template_id"`
	Host               string `json:"host"`
	Port               string `json:"port"`
}

type GetInstanceTipsResV1 struct {
	controller.BaseRes
	Data []InstanceTipResV1 `json:"data"`
}

// GetInstanceTips get instance tip list
// @Summary 获取实例提示列表
// @Description get instance tip list
// @Tags instance
// @Id getInstanceTipListV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param filter_db_type query string false "filter db type"
// @Param filter_workflow_template_id query string false "filter workflow template id"
// @Param functional_module query string false "functional module" Enums(create_audit_plan,create_workflow,sql_manage)
// @Success 200 {object} v1.GetInstanceTipsResV1
// @router /v1/projects/{project_name}/instance_tips [get]
func GetInstanceTips(c echo.Context) error {
	req := new(InstanceTipReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var operationType v1.OpPermissionType
	switch req.FunctionalModule {
	case create_audit_plan:
		operationType = v1.OpPermissionTypeSaveAuditPlan
	case create_workflow:
		operationType = v1.OpPermissionTypeCreateWorkflow
	default:
	}

	instances, err := GetCanOperationInstances(c.Request().Context(), user, req.FilterDBType, projectUid, operationType)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	s := model.GetStorage()
	template, exist, err := s.GetWorkflowTemplateByProjectId(model.ProjectUID(projectUid))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("current project doesn't has workflow template"))
	}
	instanceTipsResV1 := make([]InstanceTipResV1, 0, len(instances))
	for _, inst := range instances {
		instanceTipRes := InstanceTipResV1{
			ID:                 inst.GetIDStr(),
			Name:               inst.Name,
			Type:               inst.DbType,
			Host:               inst.Host,
			Port:               inst.Port,
			WorkflowTemplateId: uint32(template.ID),
		}
		instanceTipsResV1 = append(instanceTipsResV1, instanceTipRes)
	}

	return c.JSON(http.StatusOK, &GetInstanceTipsResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    instanceTipsResV1,
	})
}

// GetInstanceRules get instance all rule
// @Summary 获取实例应用的规则列表
// @Description get instance all rule
// @Id getInstanceRuleListV1
// @Tags instance
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param instance_name path string true "instance name"
// @Success 200 {object} v1.GetRulesResV1
// @router /v1/projects/{project_name}/instances/{instance_name}/rules [get]
func GetInstanceRules(c echo.Context) error {
	s := model.GetStorage()

	instanceName := c.Param("instance_name")
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	instance, exist, err := dms.GetInstanceInProjectByName(c.Request().Context(), projectUid, instanceName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("instance is not exist")))
	}

	can, err := CheckCurrentUserCanAccessInstances(c.Request().Context(), projectUid, controller.GetUserID(c), []*model.Instance{instance})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !can {
		return controller.JSONBaseErrorReq(c, ErrInstanceNoAccess)
	}

	rules, _, err := s.GetAllRulesByInstance(instance)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, &GetRulesResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertRulesToRes(rules),
	})
}

func CheckInstanceAndRuleTemplateDbType(ruleTemplates []*model.RuleTemplate, instances ...*model.Instance) error {
	if len(ruleTemplates) == 0 || len(instances) == 0 {
		return nil
	}

	dbType := ruleTemplates[0].DBType
	for _, rt := range ruleTemplates {
		if rt.DBType != dbType {
			return errors.New(errors.DataInvalid, fmt.Errorf("instance's and ruleTemplate's dbtype should be the same"))
		}
	}
	for _, inst := range instances {
		if inst.DbType != dbType {
			return errors.New(errors.DataInvalid, fmt.Errorf("instance's and ruleTemplate's dbtype should be the same"))
		}
	}
	return nil
}

type GetInstanceWorkflowTemplateResV1 struct {
	controller.BaseRes
	Data *WorkflowTemplateDetailResV1 `json:"data"`
}

type Table struct {
	Name string `json:"name"`
}

type ListTableBySchemaResV1 struct {
	controller.BaseRes
	Data []Table `json:"data"`
}

// ListTableBySchema list table by schema
// @Summary 获取数据库下的所有表
// @Description list table by schema
// @Id listTableBySchema
// @Tags instance
// @Param project_name path string true "project name"
// @Param instance_name path string true "instance name"
// @Param schema_name path string true "schema name"
// @Security ApiKeyAuth
// @Success 200 {object} v1.ListTableBySchemaResV1
// @router /v1/projects/{project_name}/instances/{instance_name}/schemas/{schema_name}/tables [get]
func ListTableBySchema(c echo.Context) error {
	return listTableBySchema(c)
}

type TableMetaItemHeadResV1 struct {
	FieldName string `json:"field_name"`
	Desc      string `json:"desc"`
}

type TableColumns struct {
	Rows []map[string]string      `json:"rows"`
	Head []TableMetaItemHeadResV1 `json:"head"`
}

type TableIndexes struct {
	Rows []map[string]string      `json:"rows"`
	Head []TableMetaItemHeadResV1 `json:"head"`
}

type InstanceTableMeta struct {
	Name           string       `json:"name"`
	Schema         string       `json:"schema"`
	Columns        TableColumns `json:"columns"`
	Indexes        TableIndexes `json:"indexes"`
	CreateTableSQL string       `json:"create_table_sql"`
}

type GetTableMetadataResV1 struct {
	controller.BaseRes
	Data InstanceTableMeta `json:"data"`
}

// GetTableMetadata get table metadata
// @Summary 获取表元数据
// @Description get table metadata
// @Id getTableMetadata
// @Tags instance
// @Param project_name path string true "project name"
// @Param instance_name path string true "instance name"
// @Param schema_name path string true "schema name"
// @Param table_name path string true "table name"
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetTableMetadataResV1
// @router /v1/projects/{project_name}/instances/{instance_name}/schemas/{schema_name}/tables/{table_name}/metadata [get]
func GetTableMetadata(c echo.Context) error {
	return getTableMetadata(c)
}
