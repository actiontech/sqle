package api

import (
	"actiontech.cloud/universe/sqle/v4/sqle/api/controller"

	"fmt"

	_ "actiontech.cloud/universe/sqle/v4/sqle/docs"
	"actiontech.cloud/universe/sqle/v4/sqle/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Sqle API Docs
// @version 1.0
// @description This is a sample server for dev.
// @BasePath /
func StartApi(port int, exitChan chan struct{}, logPath string) {
	e := echo.New()
	output := log.NewRotateFile(logPath, "/api.log", 1024 /*1GB*/)
	defer output.Close()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: output,
	}))
	e.HideBanner = true
	e.HidePort = true
	e.Validator = &controller.CustomValidator{}

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.POST("/instance/load_mycat_config", controller.UploadMycatConfig)
	e.POST("/instance/connection", controller.PingInstance)

	e.GET("/instances", controller.GetInsts)
	e.POST("/instances", controller.CreateInst)
	e.GET("/instances/:instance_id/", controller.GetInstance)
	e.GET("/instances/:instance_name/get_instance_by_name", controller.GetInstanceByName)

	e.DELETE("/instances/:instance_id/", controller.DeleteInstance)
	e.PATCH("/instances/:instance_id/", controller.UpdateInstance)
	e.GET("/instances/:instance_id/connection", controller.PingInstanceById)
	e.GET("/instances/:instance_id/schemas", controller.GetInstSchemas)

	e.GET("/rule_templates", controller.GetAllTpl)
	e.POST("/rule_templates", controller.CreateTemplate)
	e.GET("/rule_templates/:template_id/", controller.GetRuleTemplate)
	e.DELETE("/rule_templates/:template_id/", controller.DeleteRuleTemplate)
	e.PATCH("/rule_templates/:template_id/", controller.UpdateRuleTemplate)

	e.GET("/rules", controller.GetRules)
	e.PATCH("/rules", controller.UpdateRules)
	e.GET("/tasks", controller.GetTasks)
	e.POST("/tasks", controller.CreateTask)
	e.GET("/tasks/:task_id/", controller.GetTask)
	e.DELETE("/tasks/:task_id/", controller.DeleteTask)
	e.POST("/tasks/:task_id/inspection", controller.InspectTask)
	e.POST("/tasks/:task_id/commit", controller.CommitTask)
	e.POST("/tasks/:task_id/rollback", controller.RollbackTask)
	e.POST("/tasks/:task_id/upload_sql_file", controller.UploadSqlFile)
	e.GET("/tasks/:task_id/uploaded_sqls", controller.GetUploadedSqls)
	e.GET("/tasks/:task_id/execute_error_uploaded_sqls", controller.GetExecErrUploadedSqls)

	e.POST("/tasks/remove_by_task_ids", controller.DeleteTasks)
	e.POST("/task/create_inspect", controller.CreateAndInspectTask)

	e.GET("/schemas", controller.GetAllSchemas)
	e.POST("/schemas/manual_update", controller.ManualUpdateAllSchemas)
	e.POST("/base/reload", controller.ReloadBaseInfo)

	//SqlWhitelist
	e.GET("/sql_whitelist/:sql_whitelist_id/", controller.GetSqlWhitelistItem)
	e.POST("/sql_whitelist", controller.CreateSqlWhitelistItem)
	e.GET("/sql_whitelist", controller.GetAllWhitelist)
	e.PATCH("/sql_whitelist/:sql_white_id/", controller.UpdateSqlWhitelistItem)
	e.DELETE("/sql_whitelist/:sql_white_id/", controller.RemoveSqlWhitelistItem)
	address := fmt.Sprintf(":%v", port)
	log.Logger().Infof("starting http server on %s", address)
	log.Logger().Fatal(e.Start(address))
	close(exitChan)
}
