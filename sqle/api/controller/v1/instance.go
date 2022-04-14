package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/labstack/echo/v4"
)

var errInstanceNotExist = errors.New(errors.DataNotExist, fmt.Errorf("instance is not exist"))
var errInstanceNoAccess = errors.New(errors.DataNotExist, fmt.Errorf("instance is not exist or you can't access it"))
var errInstanceBind = errors.New(errors.DataExist, fmt.Errorf("an instance can only bind one rule template"))
var errWrongTimePeriod = errors.New(errors.DataInvalid, fmt.Errorf("wrong time period"))

type GetInstanceAdditionalMetasResV1 struct {
	controller.BaseRes
	Metas []*InstanceAdditionalMetaV1 `json:"data"`
}

type InstanceAdditionalMetaV1 struct {
	DBType string                          `json:"db_type"`
	Params []*InstanceAdditionalParamResV1 `json:"params"`
}

type InstanceAdditionalParamResV1 struct {
	Name        string `json:"name" example:"param name" form:"name"`
	Description string `json:"description" example:"参数项中文名" form:"description"`
	Type        string `json:"type" example:"int" form:"type"`
	Value       string `json:"value" example:"0" form:"value"`
}

// GetInstanceAdditionalMetas get instance additional metas
// @Summary 获取实例的额外属性列表
// @Description get instance additional metas
// @Id getInstanceAdditionalMetas
// @Tags instance
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetInstanceAdditionalMetasResV1
// @router /v1/instance_additional_metas [get]
func GetInstanceAdditionalMetas(c echo.Context) error {
	additionalParams := driver.AllAdditionalParams()
	res := &GetInstanceAdditionalMetasResV1{
		BaseRes: controller.NewBaseReq(nil),
		Metas:   []*InstanceAdditionalMetaV1{},
	}
	for name, params := range additionalParams {
		meta := &InstanceAdditionalMetaV1{
			DBType: name,
			Params: convertParamsToInstanceAdditionalParamRes(params),
		}

		res.Metas = append(res.Metas, meta)
	}
	return c.JSON(http.StatusOK, res)
}

func convertParamsToInstanceAdditionalParamRes(params params.Params) []*InstanceAdditionalParamResV1 {
	res := make([]*InstanceAdditionalParamResV1, len(params))
	for i, param := range params {
		res[i] = &InstanceAdditionalParamResV1{
			Name:        param.Key,
			Description: param.Desc,
			Type:        string(param.Type),
			Value:       param.Value,
		}
	}
	return res
}

type CreateInstanceReqV1 struct {
	Name                 string                          `json:"instance_name" form:"instance_name" example:"test" valid:"required,name"`
	DBType               string                          `json:"db_type" form:"db_type" example:"mysql"`
	User                 string                          `json:"db_user" form:"db_user" example:"root" valid:"required"`
	Host                 string                          `json:"db_host" form:"db_host" example:"10.10.10.10" valid:"required,ip_addr|uri|hostname|hostname_rfc1123"`
	Port                 string                          `json:"db_port" form:"db_port" example:"3306" valid:"required,port"`
	Password             string                          `json:"db_password" form:"db_password" example:"123456" valid:"required"`
	Desc                 string                          `json:"desc" example:"this is a test instance"`
	WorkflowTemplateName string                          `json:"workflow_template_name" form:"workflow_template_name"`
	MaintenanceTimes     []*MaintenanceTimeReqV1         `json:"maintenance_times" from:"maintenance_times"`
	RuleTemplates        []string                        `json:"rule_template_name_list" form:"rule_template_name_list"`
	Roles                []string                        `json:"role_name_list" form:"role_name_list"`
	AdditionalParams     []*InstanceAdditionalParamReqV1 `json:"additional_params" from:"additional_params"`
}

type InstanceAdditionalParamReqV1 struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type MaintenanceTimeReqV1 struct {
	MaintenanceStartTime *TimeReqV1 `json:"maintenance_start_time"`
	MaintenanceStopTime  *TimeReqV1 `json:"maintenance_stop_time"`
}

type TimeReqV1 struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
}

func convertMaintenanceTimeReqV1ToPeriod(mt []*MaintenanceTimeReqV1) model.Periods {
	periods := make(model.Periods, len(mt))
	for i, time := range mt {
		periods[i] = &model.Period{
			StartHour:   time.MaintenanceStartTime.Hour,
			StartMinute: time.MaintenanceStartTime.Minute,
			EndHour:     time.MaintenanceStopTime.Hour,
			EndMinute:   time.MaintenanceStopTime.Minute,
		}
	}
	return periods
}

// CreateInstance create instance
// @Summary 添加实例
// @Description create a instance
// @Id createInstanceV1
// @Tags instance
// @Security ApiKeyAuth
// @Accept json
// @Param instance body v1.CreateInstanceReqV1 true "add instance"
// @Success 200 {object} controller.BaseRes
// @router /v1/instances [post]
func CreateInstance(c echo.Context) error {
	s := model.GetStorage()
	req := new(CreateInstanceReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	_, exist, err := s.GetInstanceByName(req.Name)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist, fmt.Errorf("instance is exist")))
	}

	if req.DBType == "" {
		req.DBType = driver.DriverTypeMySQL
	}

	maintenancePeriod := convertMaintenanceTimeReqV1ToPeriod(req.MaintenanceTimes)
	if !maintenancePeriod.SelfCheck() {
		return controller.JSONBaseErrorReq(c, errWrongTimePeriod)
	}

	additionalParams := driver.AllAdditionalParams()[req.DBType]
	for _, additionalParam := range req.AdditionalParams {
		err = additionalParams.SetParamValue(additionalParam.Name, additionalParam.Value)
		if err != nil {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, err))
		}
	}

	instance := &model.Instance{
		DbType:            req.DBType,
		Name:              req.Name,
		User:              req.User,
		Host:              req.Host,
		Port:              req.Port,
		Password:          req.Password,
		Desc:              req.Desc,
		AdditionalParams:  additionalParams,
		MaintenancePeriod: maintenancePeriod,
	}
	// set default workflow template
	if req.WorkflowTemplateName == "" {
		workflowTemplate, exist, err := s.GetWorkflowTemplateByName(model.DefaultWorkflowTemplate)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if exist {
			instance.WorkflowTemplateId = workflowTemplate.ID
		}
	} else {
		workflowTemplate, exist, err := s.GetWorkflowTemplateByName(req.WorkflowTemplateName)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if !exist {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("workflow template is not exist")))
		}
		instance.WorkflowTemplateId = workflowTemplate.ID
	}

	templates, err := s.GetAndCheckRuleTemplateExist(req.RuleTemplates)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	roles, err := s.GetAndCheckRoleExist(req.Roles)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if !CheckInstanceCanBindOneRuleTemplate(req.RuleTemplates) {
		return controller.JSONBaseErrorReq(c, errInstanceBind)
	}

	err = CheckInstanceAndRuleTemplateDbType(templates, instance)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = s.Save(instance)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = s.UpdateInstanceRoles(instance, roles...)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = s.UpdateInstanceRuleTemplates(instance, templates...)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

// 1. admin user have all access to all instance
// 2. non-admin user have access to instance which is bound to one of his roles
func checkCurrentUserCanAccessInstance(c echo.Context, instance *model.Instance) error {
	if controller.GetUserName(c) == model.DefaultAdminUser {
		return nil
	}
	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return err
	}
	s := model.GetStorage()
	access, err := s.UserCanAccessInstance(user, instance)
	if err != nil {
		return err
	}
	if !access {
		return errInstanceNoAccess
	}
	return nil
}

type InstanceResV1 struct {
	Name                 string                          `json:"instance_name"`
	DBType               string                          `json:"db_type" example:"mysql"`
	Host                 string                          `json:"db_host" example:"10.10.10.10"`
	Port                 string                          `json:"db_port" example:"3306"`
	User                 string                          `json:"db_user" example:"root"`
	Desc                 string                          `json:"desc" example:"this is a instance"`
	WorkflowTemplateName string                          `json:"workflow_template_name,omitempty"`
	MaintenanceTimes     []*MaintenanceTimeResV1         `json:"maintenance_times" from:"maintenance_times"`
	RuleTemplates        []string                        `json:"rule_template_name_list,omitempty"`
	Roles                []string                        `json:"role_name_list,omitempty"`
	AdditionalParams     []*InstanceAdditionalParamResV1 `json:"additional_params"`
}

type MaintenanceTimeResV1 struct {
	MaintenanceStartTime *TimeResV1 `json:"maintenance_start_time"`
	MaintenanceStopTime  *TimeResV1 `json:"maintenance_stop_time"`
}

type TimeResV1 struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
}

func convertPeriodToMaintenanceTimeResV1(mt model.Periods) []*MaintenanceTimeResV1 {
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

type GetInstanceResV1 struct {
	controller.BaseRes
	Data InstanceResV1 `json:"data"`
}

func convertInstanceToRes(instance *model.Instance) InstanceResV1 {
	instanceResV1 := InstanceResV1{
		Name:             instance.Name,
		Host:             instance.Host,
		Port:             instance.Port,
		User:             instance.User,
		Desc:             instance.Desc,
		DBType:           instance.DbType,
		MaintenanceTimes: convertPeriodToMaintenanceTimeResV1(instance.MaintenancePeriod),
		AdditionalParams: []*InstanceAdditionalParamResV1{},
	}
	if instance.WorkflowTemplate != nil {
		instanceResV1.WorkflowTemplateName = instance.WorkflowTemplate.Name
	}
	if len(instance.RuleTemplates) > 0 {
		ruleTemplateNames := make([]string, 0, len(instance.RuleTemplates))
		for _, rt := range instance.RuleTemplates {
			ruleTemplateNames = append(ruleTemplateNames, rt.Name)
		}
		instanceResV1.RuleTemplates = ruleTemplateNames
	}
	if len(instance.Roles) > 0 {
		roleNames := make([]string, 0, len(instance.Roles))
		for _, r := range instance.Roles {
			roleNames = append(roleNames, r.Name)
		}
		instanceResV1.Roles = roleNames
	}
	for _, param := range instance.AdditionalParams {
		instanceResV1.AdditionalParams = append(instanceResV1.AdditionalParams, &InstanceAdditionalParamResV1{
			Name:        param.Key,
			Description: param.Desc,
			Type:        string(param.Type),
			Value:       fmt.Sprintf("%v", param.Value),
		})
	}
	return instanceResV1
}

// GetInstance get instance
// @Summary 获取实例信息
// @Description get instance db
// @Id getInstanceV1
// @Tags instance
// @Security ApiKeyAuth
// @Param instance_name path string true "instance name"
// @Success 200 {object} v1.GetInstanceResV1
// @router /v1/instances/{instance_name}/ [get]
func GetInstance(c echo.Context) error {
	s := model.GetStorage()
	instanceName := c.Param("instance_name")
	instance, exist, err := s.GetInstanceDetailByName(instanceName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errInstanceNoAccess)
	}

	err = checkCurrentUserCanAccessInstance(c, instance)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetInstanceResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertInstanceToRes(instance),
	})
}

// DeleteInstance delete instance
// @Summary 删除实例
// @Description delete instance db
// @Id deleteInstanceV1
// @Tags instance
// @Security ApiKeyAuth
// @Param instance_name path string true "instance name"
// @Success 200 {object} controller.BaseRes
// @router /v1/instances/{instance_name}/ [delete]
func DeleteInstance(c echo.Context) error {
	s := model.GetStorage()
	instanceName := c.Param("instance_name")
	instance, exist, err := s.GetInstanceByName(instanceName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errInstanceNotExist)
	}

	tasks, err := s.GetTaskByInstanceId(instance.ID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	taskIds := make([]uint, 0, len(tasks))
	for _, task := range tasks {
		taskIds = append(taskIds, task.ID)
	}
	isRunning, err := s.TaskWorkflowIsRunning(taskIds)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if isRunning {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataExist,
			fmt.Errorf("%s can't be deleted,cause on_process workflow exist", instanceName)))
	}

	err = s.Delete(instance)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type UpdateInstanceReqV1 struct {
	DBType               *string                         `json:"db_type" form:"db_type" example:"mysql"`
	User                 *string                         `json:"db_user" form:"db_user" example:"root"`
	Host                 *string                         `json:"db_host" form:"db_host" example:"10.10.10.10" valid:"omitempty,ip_addr|uri|hostname|hostname_rfc1123"`
	Port                 *string                         `json:"db_port" form:"db_port" example:"3306" valid:"omitempty,port"`
	Password             *string                         `json:"db_password" form:"db_password" example:"123456"`
	Desc                 *string                         `json:"desc" example:"this is a test instance"`
	WorkflowTemplateName *string                         `json:"workflow_template_name" form:"workflow_template_name"`
	MaintenanceTimes     []*MaintenanceTimeReqV1         `json:"maintenance_times" from:"maintenance_times"`
	RuleTemplates        []string                        `json:"rule_template_name_list" form:"rule_template_name_list"`
	Roles                []string                        `json:"role_name_list" form:"role_name_list"`
	AdditionalParams     []*InstanceAdditionalParamReqV1 `json:"additional_params" from:"additional_params"`
}

// UpdateInstance update instance
// @Summary 更新实例
// @Description update instance
// @Id updateInstanceV1
// @Tags instance
// @Security ApiKeyAuth
// @Param instance_name path string true "instance name"
// @param instance body v1.UpdateInstanceReqV1 true "update instance request"
// @Success 200 {object} controller.BaseRes
// @router /v1/instances/{instance_name}/ [patch]
func UpdateInstance(c echo.Context) error {
	req := new(UpdateInstanceReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	maintenancePeriod := convertMaintenanceTimeReqV1ToPeriod(req.MaintenanceTimes)
	if !maintenancePeriod.SelfCheck() {
		return controller.JSONBaseErrorReq(c, errWrongTimePeriod)
	}

	s := model.GetStorage()
	instanceName := c.Param("instance_name")
	instance, exist, err := s.GetInstanceByName(instanceName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errInstanceNotExist)
	}

	if !CheckInstanceCanBindOneRuleTemplate(req.RuleTemplates) {
		return controller.JSONBaseErrorReq(c, errInstanceBind)
	}

	ruleTemplates, err := s.GetRuleTemplatesByNames(req.RuleTemplates)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = CheckInstanceAndRuleTemplateDbType(ruleTemplates, instance)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	updateMap := map[string]interface{}{}
	if req.Desc != nil {
		updateMap["desc"] = *req.Desc
	}
	if req.Host != nil {
		updateMap["db_host"] = *req.Host
	}
	if req.Port != nil {
		updateMap["db_port"] = *req.Port
	}
	if req.User != nil {
		updateMap["db_user"] = *req.User
	}

	if req.MaintenanceTimes != nil {
		updateMap["maintenance_period"] = maintenancePeriod
	}

	if req.Password != nil {
		password, err := utils.AesEncrypt(*req.Password)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		updateMap["db_password"] = password
	}

	if req.WorkflowTemplateName != nil {
		// Workflow template name empty is unbound instance workflow template.
		if *req.WorkflowTemplateName == "" {
			updateMap["workflow_template_id"] = 0
		} else {
			workflowTemplate, exist, err := s.GetWorkflowTemplateByName(*req.WorkflowTemplateName)
			if err != nil {
				return controller.JSONBaseErrorReq(c, err)
			}
			if !exist {
				return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
					fmt.Errorf("workflow template is not exist")))
			}
			updateMap["workflow_template_id"] = workflowTemplate.ID
		}
	}

	if req.RuleTemplates != nil {
		ruleTemplates, err := s.GetAndCheckRuleTemplateExist(req.RuleTemplates)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		err = s.UpdateInstanceRuleTemplates(instance, ruleTemplates...)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}
	if req.Roles != nil {
		roles, err := s.GetAndCheckRoleExist(req.Roles)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		err = s.UpdateInstanceRoles(instance, roles...)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	if req.AdditionalParams != nil {
		additionalParams := driver.AllAdditionalParams()[instance.DbType]
		for _, additionalParam := range req.AdditionalParams {
			err = additionalParams.SetParamValue(additionalParam.Name, additionalParam.Value)
			if err != nil {
				return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, err))
			}
		}
		updateMap["additional_params"] = additionalParams
	}

	err = s.UpdateInstanceById(instance.ID, updateMap)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type GetInstancesReqV1 struct {
	FilterInstanceName         string `json:"filter_instance_name" query:"filter_instance_name"`
	FilterDBType               string `json:"filter_db_type" query:"filter_db_type"`
	FilterDBHost               string `json:"filter_db_host" query:"filter_db_host"`
	FilterDBPort               string `json:"filter_db_port" query:"filter_db_port"`
	FilterDBUser               string `json:"filter_db_user" query:"filter_db_user"`
	FilterWorkflowTemplateName string `json:"filter_workflow_template_name" query:"filter_workflow_template_name"`
	FilterRuleTemplateName     string `json:"filter_rule_template_name" query:"filter_rule_template_name"`
	FilterRoleName             string `json:"filter_role_name" query:"filter_role_name"`
	PageIndex                  uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize                   uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetInstancesResV1 struct {
	controller.BaseRes
	Data      []InstanceResV1 `json:"data"`
	TotalNums uint64          `json:"total_nums"`
}

// GetInstances get instances
// @Summary 获取实例信息列表
// @Description get instance info list
// @Id getInstanceListV1
// @Tags instance
// @Security ApiKeyAuth
// @Param filter_instance_name query string false "filter instance name"
// @Param filter_db_type query string false "filter db type"
// @Param filter_db_host query string false "filter db host"
// @Param filter_db_port query string false "filter db port"
// @Param filter_db_user query string false "filter db user"
// @Param filter_workflow_template_name query string false "filter workflow rule template name"
// @Param filter_rule_template_name query string false "filter rule template name"
// @Param filter_role_name query string false "filter role name"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Success 200 {object} v1.GetInstancesResV1
// @router /v1/instances [get]
func GetInstances(c echo.Context) error {
	req := new(GetInstancesReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	s := model.GetStorage()

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var offset uint32
	if req.PageIndex >= 1 {
		offset = req.PageSize * (req.PageIndex - 1)
	}
	data := map[string]interface{}{
		"filter_instance_name":          req.FilterInstanceName,
		"filter_db_host":                req.FilterDBHost,
		"filter_db_port":                req.FilterDBPort,
		"filter_db_user":                req.FilterDBUser,
		"filter_workflow_template_name": req.FilterWorkflowTemplateName,
		"filter_rule_template_name":     req.FilterRuleTemplateName,
		"filter_role_name":              req.FilterRoleName,
		"filter_db_type":                req.FilterDBType,
		"current_user_id":               user.ID,
		"check_user_can_access":         user.Name != model.DefaultAdminUser,
		"limit":                         req.PageSize,
		"offset":                        offset,
	}

	instances, count, err := s.GetInstancesByReq(data, user)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	instancesReq := []InstanceResV1{}
	for _, instance := range instances {
		instanceReq := InstanceResV1{
			Name:                 instance.Name,
			Desc:                 instance.Desc,
			Host:                 instance.Host,
			Port:                 instance.Port,
			User:                 instance.User,
			WorkflowTemplateName: instance.WorkflowTemplateName.String,
			Roles:                instance.RoleNames,
			MaintenanceTimes:     convertPeriodToMaintenanceTimeResV1(instance.MaintenancePeriod),
			RuleTemplates:        instance.RuleTemplateNames,
		}
		instancesReq = append(instancesReq, instanceReq)
	}
	return c.JSON(http.StatusOK, &GetInstancesResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      instancesReq,
		TotalNums: count,
	})
}

type GetInstanceConnectableResV1 struct {
	controller.BaseRes
	Data InstanceConnectableResV1 `json:"data"`
}

type InstanceConnectableResV1 struct {
	IsInstanceConnectable bool   `json:"is_instance_connectable"`
	ConnectErrorMessage   string `json:"connect_error_message,omitempty"`
}

func newGetInstanceConnectableResV1(err error) GetInstanceConnectableResV1 {
	if err == nil {
		return GetInstanceConnectableResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data: InstanceConnectableResV1{
				IsInstanceConnectable: true,
			},
		}
	}
	return GetInstanceConnectableResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: InstanceConnectableResV1{
			IsInstanceConnectable: false,
			ConnectErrorMessage:   err.Error(),
		},
	}
}

func checkInstanceIsConnectable(c echo.Context, instance *model.Instance) error {
	d, err := newDriverWithoutAudit(log.NewEntry(), instance, "")
	if err != nil {
		return c.JSON(http.StatusOK, newGetInstanceConnectableResV1(err))
	}
	defer d.Close(context.TODO())
	if err := d.Ping(context.TODO()); err != nil {
		return c.JSON(http.StatusOK, newGetInstanceConnectableResV1(err))
	}

	return c.JSON(http.StatusOK, newGetInstanceConnectableResV1(nil))
}

// CheckInstanceIsConnectableByName test instance db connection
// @Summary 实例连通性测试（实例提交后）
// @Description test instance db connection
// @Id checkInstanceIsConnectableByNameV1
// @Tags instance
// @Security ApiKeyAuth
// @Param instance_name path string true "instance name"
// @Success 200 {object} v1.GetInstanceConnectableResV1
// @router /v1/instances/{instance_name}/connection [get]
func CheckInstanceIsConnectableByName(c echo.Context) error {
	s := model.GetStorage()
	instanceName := c.Param("instance_name")
	instance, exist, err := s.GetInstanceByName(instanceName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errInstanceNoAccess)
	}
	err = checkCurrentUserCanAccessInstance(c, instance)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return checkInstanceIsConnectable(c, instance)
}

type GetInstanceConnectableReqV1 struct {
	DBType           string                          `json:"db_type" form:"db_type" example:"mysql"`
	User             string                          `json:"user" form:"db_user" example:"root" valid:"required"`
	Host             string                          `json:"host" form:"db_host" example:"10.10.10.10" valid:"required,ip_addr|uri|hostname|hostname_rfc1123"`
	Port             string                          `json:"port" form:"db_port" example:"3306" valid:"required,port"`
	Password         string                          `json:"password" form:"db_password" example:"123456"`
	AdditionalParams []*InstanceAdditionalParamReqV1 `json:"additional_params" from:"additional_params"`
}

// CheckInstanceIsConnectable test instance db connection
// @Summary 实例连通性测试（实例提交前）
// @Description test instance db connection 注：可直接提交创建实例接口的body，该接口的json 内容是创建实例的 json 的子集
// @Accept json
// @Id checkInstanceIsConnectableV1
// @Tags instance
// @Security ApiKeyAuth
// @Param instance body v1.GetInstanceConnectableReqV1 true "instance info"
// @Success 200 {object} v1.GetInstanceConnectableResV1
// @router /v1/instance_connection [post]
func CheckInstanceIsConnectable(c echo.Context) error {
	req := new(GetInstanceConnectableReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	if req.DBType == "" {
		req.DBType = driver.DriverTypeMySQL
	}

	additionalParams := driver.AllAdditionalParams()[req.DBType]
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
	return checkInstanceIsConnectable(c, instance)
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
// @Param instance_name path string true "instance name"
// @Success 200 {object} v1.GetInstanceSchemaResV1
// @router /v1/instances/{instance_name}/schemas [get]
func GetInstanceSchemas(c echo.Context) error {
	s := model.GetStorage()
	instanceName := c.Param("instance_name")
	instance, exist, err := s.GetInstanceByName(instanceName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errInstanceNoAccess)
	}
	err = checkCurrentUserCanAccessInstance(c, instance)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	d, err := newDriverWithoutAudit(log.NewEntry(), instance, "")
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	defer d.Close(context.TODO())
	schemas, err := d.Schemas(context.TODO())
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
)

type InstanceTipReqV1 struct {
	FilterDBType     string `json:"filter_db_type" query:"filter_db_type"`
	FunctionalModule string `json:"functional_module" query:"functional_module" enums:"create_audit_plan" valid:"omitempty,oneof=create_audit_plan"`
}

type InstanceTipResV1 struct {
	Name string `json:"instance_name"`
	Type string `json:"instance_type"`
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
// @Param filter_db_type query string false "filter db type"
// @Param functional_module query string false "functional module" Enums(create_audit_plan)
// @Success 200 {object} v1.GetInstanceTipsResV1
// @router /v1/instance_tips [get]
func GetInstanceTips(c echo.Context) error {
	req := new(InstanceTipReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	s := model.GetStorage()
	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	var instances []*model.Instance
	switch req.FunctionalModule {
	case create_audit_plan:
		instances, err = s.GetInstanceTipsByUserAndOperation(user, req.FilterDBType, model.OP_AUDIT_PLAN_SAVE)
	default:
		instances, err = s.GetInstanceTipsByUser(user, req.FilterDBType)
	}
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	instanceTipsResV1 := make([]InstanceTipResV1, 0, len(instances))

	for _, inst := range instances {
		instanceTipRes := InstanceTipResV1{
			Name: inst.Name,
			Type: inst.DbType,
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
// @Param instance_name path string true "instance name"
// @Success 200 {object} v1.GetRulesResV1
// @router /v1/instances/{instance_name}/rules [get]
func GetInstanceRules(c echo.Context) error {
	s := model.GetStorage()
	instanceName := c.Param("instance_name")
	instance, exist, err := s.GetInstanceByName(instanceName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("instance is not exist")))
	}

	err = checkCurrentUserCanAccessInstance(c, instance)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	rules, err := s.GetRulesByInstanceId(fmt.Sprintf("%d", instance.ID))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, &GetRulesResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertRulesToRes(rules),
	})
}

func CheckInstanceCanBindOneRuleTemplate(ruleTemplates []string) bool {
	return len(ruleTemplates) <= 1
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

// GetInstanceWorkflowTemplate get instance workflow template
// @Summary 获取实例应用的工作流程模板
// @Description get instance workflow template
// @Id getInstanceWorkflowTemplateV1
// @Tags instance
// @Security ApiKeyAuth
// @Param instance_name path string true "instance name"
// @Success 200 {object} v1.GetInstanceWorkflowTemplateResV1
// @router /v1/instances/{instance_name}/workflow_template [get]
func GetInstanceWorkflowTemplate(c echo.Context) error {
	s := model.GetStorage()
	instanceName := c.Param("instance_name")
	instance, exist, err := s.GetInstanceByName(instanceName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("instance is not exist")))
	}

	err = checkCurrentUserCanAccessInstance(c, instance)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	template, exist, err := s.GetWorkflowTemplateById(instance.WorkflowTemplateId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("the instance is not bound workflow template")))
	}

	res, err := getWorkflowTemplateDetailByTemplate(template)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, &GetInstanceWorkflowTemplateResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    res,
	})
}
