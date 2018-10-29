package controller

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"sqle/executor"
	"sqle/model"
	"strings"
)

type CreateInstanceReq struct {
	Name string `form:"name" example:"test"`
	// mysql, mycat, sqlserver
	DbType   string `json:"db_type" form:"type" example:"mysql"`
	User     string `form:"user" example:"root"`
	Host     string `form:"host" example:"10.10.10.10"`
	Port     string `form:"port" example:"3306"`
	Password string `form:"password" example:"123456"`
	Desc     string `form:"desc" example:"this is a test instance"`
	// this a list for rule template name
	RuleTemplates []string `json:"rule_template_name_list" form:"rule_template_name_list" example:"template_1"`
}

type InstanceRes struct {
	BaseRes
	Data model.InstanceDetail `json:"data"`
}

// @Summary 添加实例
// @Description create a instance
// @Description create a instance
// @Accept json
// @Accept json
// @Param instance body controller.CreateInstanceReq true "add instance"
// @Success 200 {object} controller.InstanceRes
// @router /instances [post]
func CreateInst(c echo.Context) error {
	s := model.GetStorage()
	req := new(CreateInstanceReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	_, exist, err := s.GetInstByName(req.Name)
	if err != nil {
		return c.JSON(200, NewBaseReq(-1, err.Error()))
	}
	if exist {
		return c.JSON(200, NewBaseReq(-1, "inst is exist"))
	}

	instance := model.Instance{
		Name:     req.Name,
		DbType:   req.DbType,
		User:     req.User,
		Host:     req.Host,
		Port:     req.Port,
		Password: req.Password,
		Desc:     req.Desc,
	}

	notExistTs := []string{}
	for _, tplName := range req.RuleTemplates {
		t, exist, err := s.GetTemplateByName(tplName)
		if err != nil {
			return c.JSON(200, NewBaseReq(-1, err.Error()))
		}
		if !exist {
			notExistTs = append(notExistTs, tplName)
		}
		instance.RuleTemplates = append(instance.RuleTemplates, t)
	}

	if len(notExistTs) > 0 {
		return c.JSON(200, NewBaseReq(-1, fmt.Sprintf("rule_template %s not exist", strings.Join(notExistTs, ", "))))
	}

	err = s.Save(&instance)
	if err != nil {
		return c.JSON(200, NewBaseReq(-1, err.Error()))
	}
	return c.JSON(200, &InstanceRes{
		BaseRes: NewBaseReq(0, "ok"),
		Data:    instance.Detail(),
	})
}

// @Summary 获取实例信息
// @Description get instance db
// @Param instance_id path string true "Instance ID"
// @Success 200 {object} controller.InstanceRes
// @router /instances/{instance_id}/ [get]
func GetInstance(c echo.Context) error {
	s := model.GetStorage()
	instanceId := c.Param("instance_id")
	instance, exist, err := s.GetInstById(instanceId)
	if err != nil {
		return c.JSON(200, PingInstRes{
			BaseRes: NewBaseReq(-1, err.Error()),
			Data:    false,
		})
	}
	if !exist {
		return c.JSON(200, PingInstRes{
			BaseRes: NewBaseReq(-1, "inst not exist"),
			Data:    false,
		})
	}
	return c.JSON(200, &InstanceRes{
		BaseRes: NewBaseReq(0, "ok"),
		Data:    instance.Detail(),
	})
}

// @Summary 删除实例
// @Description delete instance db
// @Param instance_id path string true "Instance ID"
// @Success 200 {object} controller.BaseRes
// @router /instances/{instance_id}/ [delete]
func DeleteInstance(c echo.Context) error {
	s := model.GetStorage()
	instanceId := c.Param("instance_id")
	instance, exist, err := s.GetInstById(instanceId)
	if err != nil {
		return c.JSON(200, PingInstRes{
			BaseRes: NewBaseReq(-1, err.Error()),
			Data:    false,
		})
	}
	if !exist {
		return c.JSON(200, PingInstRes{
			BaseRes: NewBaseReq(-1, "inst not exist"),
			Data:    false,
		})
	}
	err = s.Delete(instance)
	if err != nil {
		return c.JSON(200, NewBaseReq(-1, err.Error()))
	}
	return c.JSON(200, NewBaseReq(0, "ok"))
}

// @Summary 更新实例
// @Description update instance db
// @Param instance_id path string true "Instance ID"
// @param instance body controller.CreateInstanceReq true "update instance"
// @Success 200 {object} controller.BaseRes
// @router /instances/{instance_id}/ [put]
func UpdateInstance(c echo.Context) error {
	s := model.GetStorage()
	instanceId := c.Param("instance_id")
	req := new(CreateInstanceReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	instance, exist, err := s.GetInstById(instanceId)
	if err != nil {
		return c.JSON(200, NewBaseReq(-1, err.Error()))
	}
	if !exist {
		return c.JSON(200, NewBaseReq(-1, "instance not exist"))
	}
	if instance.Name != req.Name {
		_, exist, err := s.GetInstByName(req.Name)
		if err != nil {
			return c.JSON(200, NewBaseReq(-1, err.Error()))
		}
		if exist {
			return c.JSON(200, NewBaseReq(-1, "instance name is exist"))
		}
	}

	instance.Name = req.Name
	instance.Desc = req.Desc
	instance.DbType = req.DbType
	instance.Host = req.Host
	instance.Port = req.Port
	instance.User = req.User
	instance.Password = req.Password
	instance.RuleTemplates = nil

	notExistTs := []string{}
	ruleTemplates := []model.RuleTemplate{}
	for _, tplName := range req.RuleTemplates {
		t, exist, err := s.GetTemplateByName(tplName)
		if err != nil {
			return c.JSON(200, NewBaseReq(-1, err.Error()))
		}
		if !exist {
			notExistTs = append(notExistTs, tplName)
		}
		ruleTemplates = append(ruleTemplates, t)
	}

	if len(notExistTs) > 0 {
		return c.JSON(200, NewBaseReq(-1, fmt.Sprintf("rule_template %s not exist", strings.Join(notExistTs, ", "))))
	}

	err = s.Save(&instance)
	if err != nil {
		return c.JSON(200, NewBaseReq(-1, err.Error()))
	}

	err = s.UpdateInstRuleTemplate(&instance, ruleTemplates...)
	if err != nil {
		return c.JSON(200, NewBaseReq(-1, err.Error()))
	}

	return c.JSON(200, NewBaseReq(0, "ok"))
}

type GetAllInstReq struct {
	BaseRes
	Data []model.Instance `json:"data"`
}

// @Summary 实例列表
// @Description get all instances
// @Success 200 {object} controller.GetAllInstReq
// @router /instances [get]
func GetInsts(c echo.Context) error {
	s := model.GetStorage()
	if s == nil {
		c.String(500, "nil")
	}
	databases, err := s.GetInstances()
	if err != nil {
		return c.String(500, err.Error())
	}
	return c.JSON(http.StatusOK, &GetAllInstReq{
		BaseRes: NewBaseReq(0, "ok"),
		Data:    databases,
	})
}

type PingInstRes struct {
	BaseRes
	// true: ping success; false: ping failed
	Data bool `json:"data"`
}

// @Summary 实例连通性测试
// @Description test instance db connection
// @Param instance_id path string true "Instance ID"
// @Success 200 {object} controller.PingInstRes
// @router /instances/{instance_id}/connection [get]
func PingInst(c echo.Context) error {
	s := model.GetStorage()
	instId := c.Param("instance_id")
	inst, exist, err := s.GetInstById(instId)
	if err != nil {
		return c.JSON(200, PingInstRes{
			BaseRes: NewBaseReq(-1, err.Error()),
			Data:    false,
		})
	}
	if !exist {
		return c.JSON(200, PingInstRes{
			BaseRes: NewBaseReq(-1, "inst not exist"),
			Data:    false,
		})
	}
	if err := executor.Ping(inst); err != nil {
		return c.JSON(200, PingInstRes{
			BaseRes: NewBaseReq(-1, err.Error()),
			Data:    false,
		})
	}
	return c.JSON(200, PingInstRes{
		BaseRes: NewBaseReq(0, ""),
		Data:    true,
	})
}

type GetSchemaRes struct {
	BaseRes
	Data []string `json:"data" example:"db1"`
}

// @Summary 实例 Schema 列表
// @Description instance schema list
// @Param instance_id path string true "Instance ID"
// @Success 200 {object} controller.GetSchemaRes
// @router /instances/{instance_id}/schemas [get]
func GetInstSchemas(c echo.Context) error {
	s := model.GetStorage()
	instanceId := c.Param("instance_id")
	instance, exist, err := s.GetInstById(instanceId)
	if err != nil {
		return c.JSON(200, NewBaseReq(-1, err.Error()))
	}
	if !exist {
		return c.JSON(200, NewBaseReq(-1, "instance not exist"))
	}
	schemas, err := executor.ShowDatabases(instance)
	if err != nil {
		return c.JSON(200, NewBaseReq(-1, err.Error()))
	}
	return c.JSON(200, &GetSchemaRes{
		BaseRes: NewBaseReq(0, "ok"),
		Data:    schemas,
	})
}

type GetAllSchemasRes struct {
	BaseRes
	Data []GetAllSchemasData `json:"data"`
}

type GetAllSchemasData struct {
	Name    string
	Host    string
	Port    string
	Schemas []string
}

// @Summary 所有实例的 Schema 列表
// @Description all schema list
// @Success 200 {object} controller.GetAllSchemasRes
// @router /schemas/ [get]
func GetAllSchemas(c echo.Context) error {
	return nil
}
