package v1

import (
	"context"
	"fmt"
	"net/http"

	"actiontech.cloud/sqle/sqle/sqle/api/controller"
	"actiontech.cloud/sqle/sqle/sqle/driver"
	"actiontech.cloud/sqle/sqle/sqle/errors"
	"actiontech.cloud/sqle/sqle/sqle/log"
	"actiontech.cloud/sqle/sqle/sqle/model"
	"actiontech.cloud/universe/ucommon/v4/util"

	"github.com/labstack/echo/v4"
)

var instanceNotExistError = errors.New(errors.DataNotExist, fmt.Errorf("instance is not exist"))
var instanceNoAccessError = errors.New(errors.DataNotExist, fmt.Errorf("instance is not exist or you can't access it"))

type CreateInstanceReqV1 struct {
	Name                 string   `json:"instance_name" form:"instance_name" example:"test" valid:"required,name"`
	DBType               string   `json:"db_type" form:"db_type" example:"mysql"`
	User                 string   `json:"db_user" form:"db_user" example:"root" valid:"required"`
	Host                 string   `json:"db_host" form:"db_host" example:"10.10.10.10" valid:"required,ipv4"`
	Port                 string   `json:"db_port" form:"db_port" example:"3306" valid:"required,port"`
	Password             string   `json:"db_password" form:"db_password" example:"123456" valid:"required"`
	Desc                 string   `json:"desc" example:"this is a test instance"`
	WorkflowTemplateName string   `json:"workflow_template_name" form:"workflow_template_name"`
	RuleTemplates        []string `json:"rule_template_name_list" form:"rule_template_name_list"`
	Roles                []string `json:"role_name_list" form:"role_name_list"`
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
		req.DBType = model.DBTypeMySQL
	}
	instance := &model.Instance{
		DbType:   req.DBType,
		Name:     req.Name,
		User:     req.User,
		Host:     req.Host,
		Port:     req.Port,
		Password: req.Password,
		Desc:     req.Desc,
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
	err = s.CheckInstanceBindCount(req.RuleTemplates)
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
		return instanceNoAccessError
	}
	return nil
}

type InstanceResV1 struct {
	Name                 string   `json:"instance_name"`
	DBType               string   `json:"db_type" example:"mysql"`
	Host                 string   `json:"db_host" example:"10.10.10.10"`
	Port                 string   `json:"db_port" example:"3306"`
	User                 string   `json:"db_user" example:"root"`
	Desc                 string   `json:"desc" example:"this is a instance"`
	WorkflowTemplateName string   `json:"workflow_template_name,omitempty"`
	RuleTemplates        []string `json:"rule_template_name_list,omitempty"`
	Roles                []string `json:"role_name_list,omitempty"`
}

type GetInstanceResV1 struct {
	controller.BaseRes
	Data InstanceResV1 `json:"data"`
}

func convertInstanceToRes(instance *model.Instance) InstanceResV1 {
	instanceResV1 := InstanceResV1{
		Name:   instance.Name,
		Host:   instance.Host,
		Port:   instance.Port,
		User:   instance.User,
		Desc:   instance.Desc,
		DBType: instance.DbType,
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
		return controller.JSONBaseErrorReq(c, instanceNoAccessError)
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
		return controller.JSONBaseErrorReq(c, instanceNotExistError)
	}
	err = s.Delete(instance)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

type UpdateInstanceReqV1 struct {
	DBType               *string  `json:"db_type" form:"db_type" example:"mysql"`
	User                 *string  `json:"db_user" form:"db_user" example:"root"`
	Host                 *string  `json:"db_host" form:"db_host" example:"10.10.10.10" valid:"omitempty,ipv4"`
	Port                 *string  `json:"db_port" form:"db_port" example:"3306" valid:"omitempty,port"`
	Password             *string  `json:"db_password" form:"db_password" example:"123456"`
	Desc                 *string  `json:"desc" example:"this is a test instance"`
	WorkflowTemplateName *string  `json:"workflow_template_name" form:"workflow_template_name"`
	RuleTemplates        []string `json:"rule_template_name_list" form:"rule_template_name_list"`
	Roles                []string `json:"role_name_list" form:"role_name_list"`
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

	s := model.GetStorage()
	instanceName := c.Param("instance_name")
	instance, exist, err := s.GetInstanceByName(instanceName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, instanceNotExistError)
	}

	err = s.CheckInstanceBindCount(req.RuleTemplates)
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
	if req.Password != nil {
		password, err := util.AesEncrypt(*req.Password)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		updateMap["db_password"] = password
	}
	if req.DBType != nil {
		updateMap["db_type"] = *req.DBType
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

	instances, count, err := s.GetInstancesByReq(data)
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

func checkInstanceIsConnectable(c echo.Context, instance *model.Instance) error {
	d, err := driver.NewDriver(log.NewEntry(), instance, "")
	if err != nil {
		return err
	}
	defer d.Close(context.TODO())
	if err := d.Ping(context.TODO()); err != nil {
		return c.JSON(http.StatusOK, GetInstanceConnectableResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data: InstanceConnectableResV1{
				IsInstanceConnectable: false,
				ConnectErrorMessage:   err.Error(),
			},
		})
	} else {
		return c.JSON(http.StatusOK, GetInstanceConnectableResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data: InstanceConnectableResV1{
				IsInstanceConnectable: true,
			},
		})
	}
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
		return controller.JSONBaseErrorReq(c, instanceNoAccessError)
	}
	err = checkCurrentUserCanAccessInstance(c, instance)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return checkInstanceIsConnectable(c, instance)
}

type GetInstanceConnectableReqV1 struct {
	DBType   string `json:"db_type" form:"db_type" example:"mysql"`
	User     string `json:"user" form:"db_user" example:"root" valid:"required"`
	Host     string `json:"host" form:"db_host" example:"10.10.10.10" valid:"required,ipv4"`
	Port     string `json:"port" form:"db_port" example:"3306" valid:"required,port"`
	Password string `json:"password" form:"db_password" example:"123456"`
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
		req.DBType = model.DBTypeMySQL
	}
	instance := &model.Instance{
		DbType:   req.DBType,
		User:     req.User,
		Host:     req.Host,
		Port:     req.Port,
		Password: req.Password,
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
		return controller.JSONBaseErrorReq(c, instanceNoAccessError)
	}
	err = checkCurrentUserCanAccessInstance(c, instance)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	d, err := driver.NewDriver(log.NewEntry(), instance, "")
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
// @Success 200 {object} v1.GetInstanceTipsResV1
// @router /v1/instance_tips [get]
func GetInstanceTips(c echo.Context) error {
	s := model.GetStorage()
	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	roles, err := s.GetUserInstanceTip(user)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	instanceTipsResV1 := make([]InstanceTipResV1, 0, len(roles))

	for _, role := range roles {
		instanceTipRes := InstanceTipResV1{
			Name: role.Name,
			Type: role.DbType,
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
	rules, err := s.GetRulesByInstanceId(fmt.Sprintf("%d", instance.ID))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, &GetRulesResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertRulesToRes(rules),
	})
}
