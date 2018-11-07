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

	e.GET("/instances", controller.GetInsts)
	e.POST("/instances", controller.CreateInst)
	e.GET("/instances/:instance_id/", controller.GetInstance)
	e.DELETE("/instances/:instance_id/", controller.DeleteInstance)
	e.PUT("/instances/:instance_id/", controller.UpdateInstance)
	e.GET("/instances/:instance_id/connection", controller.PingInst)
	e.GET("/instances/:instance_id/schemas", controller.GetInstSchemas)

	e.GET("/rule_templates", controller.GetAllTpl)
	e.POST("/rule_templates", controller.CreateTemplate)
	e.GET("/rule_templates/:template_id/", controller.GetRuleTemplate)
	e.DELETE("/rule_templates/:template_id/", controller.DeleteRuleTemplate)
	e.PUT("/rule_templates/:template_id/", controller.UpdateRuleTemplate)

	e.GET("/rules", controller.GetRules)

	e.GET("/tasks", controller.GetTasks)
	e.POST("/tasks", controller.CreateTask)
	e.GET("/tasks/:task_id/", controller.GetTask)
	e.POST("/tasks/:task_id/", controller.DeleteTask)
	e.POST("/tasks/:task_id/inspection", controller.InspectTask)
	e.POST("/tasks/:task_id/commit", controller.CommitTask)
	e.POST("/tasks/:task_id/rollback", controller.RollbackTask)

	e.GET("/schemas", controller.GetAllSchemas)
	e.POST("/schemas/manual_update", controller.ManualUpdateAllSchemas)

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