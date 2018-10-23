package api

import (
	"fmt"
	"sqle/api/controller"

	"github.com/labstack/echo"
	"github.com/swaggo/echo-swagger"
	_ "sqle/docs"
)

// @title Sqle API Docs
// @version 1.0
// @description This is a sample server for dev.
// @BasePath /
func StartApi(port int, exitChan chan struct{}) {
	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.POST("/instances", controller.CreateInst)
	e.GET("/instances/:inst_id/connection", controller.PingInst)
	e.GET("/instances", controller.GetInsts)

	e.GET("/rule_templates", controller.GetAllTpl)
	e.POST("/rule_templates", controller.CreateTemplate)

	e.GET("/tasks", controller.GetTasks)
	e.POST("/tasks", controller.CreateTask)
	e.POST("/tasks/:task_id/inspection",controller.InspectTask)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", port)))
	close(exitChan)
}

func StartDocs(port int, exitChan chan struct{}) {
	e := echo.New()
	e.HideBanner = true
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Start(fmt.Sprintf(":%v", port))
	close(exitChan)
}
