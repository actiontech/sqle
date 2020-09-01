package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"actiontech.cloud/universe/sqle/v3/sqle/api/server"
	"actiontech.cloud/universe/sqle/v3/sqle/errors"
	"actiontech.cloud/universe/sqle/v3/sqle/executor"
	"actiontech.cloud/universe/sqle/v3/sqle/log"
	"actiontech.cloud/universe/sqle/v3/sqle/model"
	"actiontech.cloud/universe/ucommon/v3/util"
	"github.com/labstack/echo/v4"
)

var nullString = ""

type CreateInstanceReq struct {
	Name *string `json:"name" example:"test" valid:"required"`
	// mysql, mycat, sqlserver
	DbType   *string `json:"db_type" example:"mysql" valid:"required,in(mysql|mycat|sqlserver)"`
	User     *string `json:"user" example:"root" valid:"required"`
	Host     *string `json:"host" example:"10.10.10.10" valid:"required,ipv4"`
	Port     *string `json:"port" example:"3306" valid:"required,range(1|65535)"`
	Password *string `json:"password" example:"123456" valid:"required"`
	Desc     *string `json:"desc" example:"this is a test instance" valid:"-"`
	// this a list for rule template name
	RuleTemplates []string `json:"rule_template_id_list" example:"1" valid:"-"`
	// mycat_config is required if db_type is "mycat"
	MycatConfig *model.MycatConfig `json:"mycat_config" valid:"-"`
}

type InstanceRes struct {
	BaseRes
	Data model.InstanceDetail `json:"data"`
}

// @Summary 添加实例
// @Description create a instance
// @Accept json
// @Param instance body controller.CreateInstanceReq true "add instance"
// @Success 200 {object} controller.InstanceRes
// @router /instances [post]
func CreateInst(c echo.Context) error {
	s := model.GetStorage()
	req := new(CreateInstanceReq)
	if err := c.Bind(req); err != nil {
		return c.JSON(200, NewBaseReq(err))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}

	_, exist, err := s.GetInstByName(*req.Name)
	if err != nil {
		return c.JSON(200, NewBaseReq(err))
	}
	if exist {
		err := fmt.Errorf("inst is exist")
		return c.JSON(200, NewBaseReq(errors.New(errors.INSTANCE_EXIST, err)))
	}

	instance := &model.Instance{
		Name:        *req.Name,
		DbType:      *req.DbType,
		User:        *req.User,
		Host:        *req.Host,
		Port:        *req.Port,
		Password:    *req.Password,
		Desc:        *req.Desc,
		MycatConfig: req.MycatConfig,
	}

	tpls := []model.RuleTemplate{}
	notExistTpls := []string{}
	for _, tplId := range req.RuleTemplates {
		t, exist, err := s.GetTemplateById(tplId)
		if err != nil {
			return c.JSON(200, NewBaseReq(err))
		}
		if !exist {
			notExistTpls = append(notExistTpls, tplId)
		}
		tpls = append(tpls, t)
	}

	if len(notExistTpls) > 0 {
		err := fmt.Errorf("rule_template id : %s not exist", strings.Join(notExistTpls, ", "))
		return c.JSON(200, NewBaseReq(errors.New(errors.RULE_TEMPLATE_NOT_EXIST, err)))
	}

	err = s.Save(instance)
	if err != nil {
		return c.JSON(200, NewBaseReq(err))
	}

	err = s.UpdateInstRuleTemplate(instance, tpls...)
	if err != nil {
		return c.JSON(200, NewBaseReq(err))
	}

	go server.GetSqled().UpdateAndGetInstanceStatus(log.NewEntry(), instance)

	return c.JSON(200, &InstanceRes{
		BaseRes: NewBaseReq(nil),
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
			BaseRes: NewBaseReq(err),
			Data:    false,
		})
	}
	if !exist {
		return c.JSON(200, PingInstRes{
			BaseRes: INSTANCE_NOT_EXIST_ERROR,
			Data:    false,
		})
	}
	return c.JSON(200, &InstanceRes{
		BaseRes: NewBaseReq(nil),
		Data:    instance.Detail(),
	})
}

// @Summary 获取实例信息
// @Description get instance db
// @Param instance_name path string true "Instance Name"
// @Success 200 {object} controller.InstanceRes
// @router /instances/{instance_name}/get_instance_by_name [get]
func GetInstanceByName(c echo.Context) error {
	s := model.GetStorage()
	instanceName := c.Param("instance_name")
	instance, exist, err := s.GetInstByName(instanceName)
	if err != nil {
		return c.JSON(200, PingInstRes{
			BaseRes: NewBaseReq(err),
			Data:    false,
		})
	}
	if !exist {
		return c.JSON(200, PingInstRes{
			BaseRes: INSTANCE_NOT_EXIST_ERROR,
			Data:    false,
		})
	}
	return c.JSON(200, &InstanceRes{
		BaseRes: NewBaseReq(nil),
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
			BaseRes: NewBaseReq(err),
			Data:    false,
		})
	}
	if !exist {
		return c.JSON(200, PingInstRes{
			BaseRes: INSTANCE_NOT_EXIST_ERROR,
			Data:    false,
		})
	}
	err = s.Delete(instance)
	if err != nil {
		return c.JSON(200, NewBaseReq(err))
	}

	server.GetSqled().DeleteInstanceStatus(instance)
	return c.JSON(200, NewBaseReq(nil))
}

// @Summary 更新实例
// @Description update instance db
// @Param instance_id path string true "Instance ID"
// @param instance body controller.CreateInstanceReq true "update instance"
// @Success 200 {object} controller.BaseRes
// @router /instances/{instance_id}/ [patch]
func UpdateInstance(c echo.Context) error {
	s := model.GetStorage()
	instanceId := c.Param("instance_id")
	instance, exist, err := s.GetInstById(instanceId)
	if err != nil {
		return c.JSON(200, NewBaseReq(err))
	}
	if !exist {
		return c.JSON(200, INSTANCE_NOT_EXIST_ERROR)
	}

	req := new(CreateInstanceReq)
	if err := c.Bind(req); err != nil {
		return c.JSON(200, NewBaseReq(err))
	}

	updateMap := map[string]string{}
	if req.Name != nil {
		if instance.Name != *req.Name {
			_, exist, err := s.GetInstByName(*req.Name)
			if err != nil {
				return c.JSON(200, NewBaseReq(err))
			}
			if exist {
				return c.JSON(200, NewBaseReq(errors.New(errors.INSTANCE_EXIST,
					fmt.Errorf("instance name is exist"))))
			}
		}
	}
	if req.Name != nil {
		updateMap["name"] = *req.Name
	}
	if req.Desc != nil {
		updateMap["desc"] = *req.Desc
	}
	if req.DbType != nil {
		updateMap["db_type"] = *req.DbType
	}
	if req.Host != nil {
		updateMap["host"] = *req.Host
	}
	if req.Port != nil {
		updateMap["port"] = *req.Port
	}
	if req.User != nil {
		updateMap["user"] = *req.User
	}
	if req.Password != nil {
		password, err := util.AesEncrypt(*req.Password)
		if err != nil {
			return c.JSON(200, NewBaseReq(err))
		}
		updateMap["password"] = password
	}

	if req.MycatConfig != nil {
		instance.MycatConfig = req.MycatConfig
		config, err := json.Marshal(req.MycatConfig)
		if err != nil {
			return c.JSON(200, NewBaseReq(err))
		}
		updateMap["mycat_config"] = string(config)
	}

	if req.RuleTemplates != nil {
		notExistTs := []string{}
		ruleTemplates := []model.RuleTemplate{}
		for _, tplName := range req.RuleTemplates {
			t, exist, err := s.GetTemplateByName(tplName)
			if err != nil {
				return c.JSON(200, NewBaseReq(err))
			}
			if !exist {
				notExistTs = append(notExistTs, tplName)
			}
			ruleTemplates = append(ruleTemplates, t)
		}

		if len(notExistTs) > 0 {
			err := fmt.Errorf("rule_template %s not exist", strings.Join(notExistTs, ", "))
			return c.JSON(200, NewBaseReq(errors.New(errors.RULE_TEMPLATE_NOT_EXIST, err)))
		}

		err = s.UpdateInstRuleTemplate(instance, ruleTemplates...)
		if err != nil {
			return c.JSON(200, NewBaseReq(err))
		}
	}

	err = s.UpdateInstanceById(instanceId, updateMap)
	if err != nil {
		return c.JSON(200, NewBaseReq(err))
	}

	go server.GetSqled().UpdateAndGetInstanceStatus(log.NewEntry(), instance)
	return c.JSON(200, NewBaseReq(nil))
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
	databases, err := s.GetInstances()
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	return c.JSON(http.StatusOK, &GetAllInstReq{
		BaseRes: NewBaseReq(nil),
		Data:    databases,
	})
}

type PingInstRes struct {
	BaseRes
	// true: ping success; false: ping failed
	Data bool `json:"data"`
}

// @Summary 实例连通性测试（实例提交后）
// @Description test instance db connection
// @Param instance_id path string true "Instance ID"
// @Success 200 {object} controller.PingInstRes
// @router /instances/{instance_id}/connection [get]
func PingInstanceById(c echo.Context) error {
	s := model.GetStorage()
	instId := c.Param("instance_id")
	inst, exist, err := s.GetInstById(instId)
	if err != nil {
		return c.JSON(200, PingInstRes{
			BaseRes: NewBaseReq(err),
			Data:    false,
		})
	}
	if !exist {
		return c.JSON(200, PingInstRes{
			BaseRes: INSTANCE_NOT_EXIST_ERROR,
			Data:    false,
		})
	}
	if err := executor.Ping(log.NewEntry(), inst); err != nil {
		return c.JSON(200, PingInstRes{
			BaseRes: NewBaseReq(errors.New(errors.STATUS_OK, err)),
			Data:    false,
		})
	}
	return c.JSON(200, PingInstRes{
		BaseRes: NewBaseReq(nil),
		Data:    true,
	})
}

type PingInstanceReq struct {
	// mysql, mycat, sqlserver
	DbType   string `json:"db_type" example:"mysql"`
	User     string `json:"user" example:"root"`
	Host     string `json:"host" example:"10.10.10.10"`
	Port     string `json:"port" example:"3306"`
	Password string `json:"password" example:"123456"`
	// mycat_config is required if db_type is "mycat"
	MycatConfig *model.MycatConfig `json:"mycat_config"`
}

// @Summary 实例连通性测试（实例提交前）
// @Description test instance db connection 注：可直接提交创建实例接口的body，该接口的json 内容是创建实例的 json 的子集
// @Accept json
// @Param instance body controller.PingInstanceReq true "instance info"
// @Success 200 {object} controller.PingInstRes
// @router /instance/connection [post]
func PingInstance(c echo.Context) error {
	req := new(PingInstanceReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	instance := &model.Instance{
		DbType:      req.DbType,
		User:        req.User,
		Host:        req.Host,
		Port:        req.Port,
		Password:    req.Password,
		MycatConfig: req.MycatConfig,
	}

	if err := executor.Ping(log.NewEntry(), instance); err != nil {
		return c.JSON(200, PingInstRes{
			BaseRes: NewBaseReq(errors.New(errors.STATUS_OK, err)),
			Data:    false,
		})
	}
	return c.JSON(200, PingInstRes{
		BaseRes: NewBaseReq(nil),
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
		return c.JSON(200, NewBaseReq(err))
	}
	if !exist {
		return c.JSON(200, INSTANCE_NOT_EXIST_ERROR)
	}
	status, err := server.GetSqled().UpdateAndGetInstanceStatus(log.NewEntry(), instance)
	if err != nil {
		return c.JSON(200, NewBaseReq(err))
	}
	return c.JSON(200, &GetSchemaRes{
		BaseRes: NewBaseReq(nil),
		Data:    status.Schemas,
	})
}

type GetAllSchemasRes struct {
	BaseRes
	Data []server.InstanceStatus `json:"data"`
}

// @Summary 所有实例的 Schema 列表
// @Description all schema list
// @Success 200 {object} controller.GetAllSchemasRes
// @router /schemas [get]
func GetAllSchemas(c echo.Context) error {
	return c.JSON(200, &GetAllSchemasRes{
		BaseRes: NewBaseReq(nil),
		Data:    server.GetSqled().GetAllInstanceStatus(),
	})
}

// @Summary 手动刷新 Schema 列表
// @Description update schema list
// @Success 200 {object} controller.BaseRes
// @router /schemas/manual_update [post]
func ManualUpdateAllSchemas(c echo.Context) error {
	go server.GetSqled().UpdateAllInstanceStatus(log.NewEntry())
	return c.JSON(200, NewBaseReq(nil))
}

type LoadMycatConfigRes struct {
	BaseRes
	Data CreateInstanceReq `json:"data"`
}

// @Summary 解析 mycat 配置
// @Description unmarshal mycat config file
// @Accept mpfd
// @Param server_xml formData file true "server.xml"
// @Param schema_xml formData file true "schema.xml"
// @Param rule_xml formData file true "rule.xml"
// @Success 200 {object} controller.LoadMycatConfigRes
// @router /instance/load_mycat_config [post]
func UploadMycatConfig(c echo.Context) error {
	_, server, err := readFileToByte(c, "server_xml")
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	_, schemas, err := readFileToByte(c, "schema_xml")
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	_, rules, err := readFileToByte(c, "rule_xml")
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	instance, err := model.LoadMycatServerFromXML(server, schemas, rules)
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	return c.JSON(http.StatusOK, LoadMycatConfigRes{
		BaseRes: NewBaseReq(nil),
		Data: CreateInstanceReq{
			Name:          &nullString,
			Desc:          &nullString,
			Host:          &nullString,
			Port:          &nullString,
			Password:      &nullString,
			DbType:        &instance.DbType,
			User:          &instance.User,
			RuleTemplates: []string{},
			MycatConfig:   instance.MycatConfig,
		},
	})
}
