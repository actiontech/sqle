package controller

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"sqle/executor"
	"sqle/storage"
)

type CreateInstReq struct {
	Name     string `form:"name" example:"test"`
	DbType   int    `form:"type" example:"1"`
	User     string `form:"user" example:"root"`
	Host     string `form:"host" example:"10.10.10.10"`
	Port     string `form:"port" example:"3306"`
	Password string `form:"password" example:"123456"`
	Desc     string `form:"desc" example:"this is a test instance"`
}

// @Title createInstance
// @Description create a instance
// @Accept json
// @Accept json
// @Param instance body controller.CreateInstReq true "add instance"
// @Success 200 {object} controller.BaseReq
//// @router /instances [post]
func CreateInst(c echo.Context) error {
	s := storage.GetStorage()
	req := new(CreateInstReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	database := &storage.Db{
		DbType:   req.DbType,
		User:     req.User,
		Host:     req.Host,
		Port:     req.Port,
		Password: req.Password,
		Alias:    req.Desc,
	}

	exist, err := s.Exist(database)
	if err != nil {
		return err
	}
	if exist {
		return c.JSON(200, NewBaseReq(-1, "inst is exist"))
	}
	err = s.Save(database)
	if err != nil {
		return err
	}
	return c.JSON(200, NewBaseReq(0, "ok"))
}

//func (c *BaseController) PingDatabase() {
//	database := &storage.Db{}
//	database.User = c.GetString("user")
//	database.DbType, _ = c.GetInt("db_type", 0)
//	database.Host = c.GetString("host")
//	database.Port = c.GetString("port")
//	database.Password = c.GetString("password")
//	err := executor.Ping(database)
//	if err != nil {
//		c.CustomAbort(500, err.Error())
//	}
//	c.Ctx.WriteString("ok")
//	return
//}
//
//func (c *BaseController) GetDatabaseSchemas() {
//	dbId := c.Ctx.Input.Param(":dbId")
//	db, err := c.storage.GetDatabaseById(dbId)
//	if err != nil {
//		c.CustomAbort(500, err.Error())
//	}
//	schemas, err := executor.ShowDatabase(db)
//	if err != nil {
//		c.CustomAbort(500, err.Error())
//	}
//	c.serveJson(schemas)
//}
//
//func (c *BaseController) DatabaseList() {
//	databases, err := c.storage.GetDatabases()
//	if err != nil {
//		c.CustomAbort(500, err.Error())
//	}
//	c.serveJson(databases)
//}

type GetAllInstReq struct {
	BaseReq
	Data []storage.Db `json:"data"`
}

// @Title getInstanceList
// @Description get all instances
// @Success 200 {object} controller.GetAllInstReq
// @router /instances [get]
func GetInsts(c echo.Context) error {
	s := storage.GetStorage()
	if s == nil {
		c.String(500, "nil")
	}
	databases, err := s.GetDatabases()
	if err != nil {
		return c.String(500, err.Error())
	}
	return c.JSON(http.StatusOK, &GetAllInstReq{
		BaseReq: NewBaseReq(0, "ok"),
		Data:    databases,
	})
}

type PingInstReq struct {
	BaseReq
	Status bool `json:"status"`
}

// @Title getInstanceList
// @Description get all instances
// @Param inst_id path string true "Instance ID"
// @Success 200 {object} controller.PingInstReq
// @router /instances/{inst_id}/connection [post]
func PingInst(c echo.Context) error {
	s := storage.GetStorage()
	//req := new(CreateInstReq)
	//if err := c.Bind(req); err != nil {
	//	return err
	//}
	instId := c.Param("inst_id")
	fmt.Println("inst id: ", instId)
	inst, exist, err := s.GetDatabaseById(instId)
	if err != nil {
		return c.JSON(200, PingInstReq{
			BaseReq: NewBaseReq(-1, err.Error()),
			Status:  false,
		})
	}
	if !exist {
		return c.JSON(200, PingInstReq{
			BaseReq: NewBaseReq(0, "inst not exist"),
			Status:  false,
		})
	}
	fmt.Println(inst)
	if err := executor.Ping(inst); err != nil {
		return c.JSON(200, PingInstReq{
			BaseReq: NewBaseReq(0, err.Error()),
			Status:  false,
		})
	}
	return c.JSON(200, PingInstReq{
		BaseReq: NewBaseReq(0, ""),
		Status:  true,
	})
}
