package api

import (
	"actiontech.cloud/universe/sqle/v4/sqle/api/controller"
	"actiontech.cloud/universe/sqle/v4/sqle/api/controller/v1"
	"net/http"
	"strings"

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
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
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

	//e.POST("/instance/load_mycat_config", controller.UploadMycatConfig)
	//e.POST("/instance/connection", controller.PingInstance)
	//
	//e.GET("/instances", controller.GetInsts)
	//e.POST("/instances", controller.CreateInst)
	//e.GET("/instances/:instance_id/", controller.GetInstance)
	//e.GET("/instances/:instance_name/get_instance_by_name", controller.GetInstanceByName)
	//
	//e.DELETE("/instances/:instance_id/", controller.DeleteInstance)
	//e.PATCH("/instances/:instance_id/", controller.UpdateInstance)
	//e.GET("/instances/:instance_id/connection", controller.PingInstanceById)
	//e.GET("/instances/:instance_id/schemas", controller.GetInstSchemas)
	//
	//e.GET("/rule_templates", controller.GetAllTpl)
	//e.POST("/rule_templates", controller.CreateTemplate)
	//e.GET("/rule_templates/:template_id/", controller.GetRuleTemplate)
	//e.DELETE("/rule_templates/:template_id/", controller.DeleteRuleTemplate)
	//e.PATCH("/rule_templates/:template_id/", controller.UpdateRuleTemplate)
	//
	//e.GET("/rules", controller.GetRules)
	//e.PATCH("/rules", controller.UpdateRules)
	//e.GET("/tasks", controller.GetTasks)
	//e.POST("/tasks", controller.CreateTask)
	//e.GET("/tasks/:task_id/", controller.GetTask)
	//e.DELETE("/tasks/:task_id/", controller.DeleteTask)
	//e.POST("/tasks/:task_id/inspection", controller.InspectTask)
	//e.POST("/tasks/:task_id/commit", controller.CommitTask)
	//e.POST("/tasks/:task_id/rollback", controller.RollbackTask)
	//e.POST("/tasks/:task_id/upload_sql_file", controller.UploadSqlFile)
	//e.GET("/tasks/:task_id/uploaded_sqls", controller.GetUploadedSqls)
	//e.GET("/tasks/:task_id/execute_error_uploaded_sqls", controller.GetExecErrUploadedSqls)
	//
	//e.POST("/tasks/remove_by_task_ids", controller.DeleteTasks)
	//e.POST("/task/create_inspect", controller.CreateAndInspectTask)
	//
	//e.GET("/schemas", controller.GetAllSchemas)
	//e.POST("/schemas/manual_update", controller.ManualUpdateAllSchemas)
	//e.POST("/base/reload", controller.ReloadBaseInfo)
	//
	////SqlWhitelist
	//e.GET("/sql_whitelist/:sql_whitelist_id/", controller.GetSqlWhitelistItemById)
	//e.POST("/sql_whitelist", controller.CreateSqlWhitelistItem)
	//e.GET("/sql_whitelist", controller.GetSqlWhitelist)
	//e.PATCH("/sql_whitelist/:sql_white_id/", controller.UpdateSqlWhitelistItem)
	//e.DELETE("/sql_whitelist/:sql_white_id/", controller.RemoveSqlWhitelistItem)

	e.POST("/v1/login", v1.Login)

	v1Router := e.Group("/v1")
	v1Router.Use(JWTTokenAdapter(), middleware.JWT([]byte(v1.JWTSecret)))

	// v1 admin api, just admin user can access.
	{
		// user
		v1Router.GET("/test", v1.Test, AdminUserAllowed())
		v1Router.GET("/users", v1.GetUsers, AdminUserAllowed())
		v1Router.GET("/user_tips", v1.GetUserTips, AdminUserAllowed())
		v1Router.POST("/users", v1.CreateUser, AdminUserAllowed())
		v1Router.GET("/users/:user_name/", v1.GetUser, AdminUserAllowed())
		v1Router.PATCH("/users/:user_name/", v1.UpdateUser, AdminUserAllowed())
		v1Router.DELETE("/users/:user_name/", v1.DeleteUser, AdminUserAllowed())

		// role
		v1Router.GET("/roles", v1.GetRoles, AdminUserAllowed())
		v1Router.GET("/role_tips", v1.GetRoleTips, AdminUserAllowed())
		v1Router.POST("/roles", v1.CreateRole, AdminUserAllowed())
		v1Router.PATCH("/roles/:role_name/", v1.UpdateRole, AdminUserAllowed())
		v1Router.DELETE("/roles/:role_name/", v1.DeleteRole, AdminUserAllowed())

		// instance
		v1Router.POST("/instances", v1.CreateInstance, AdminUserAllowed())
		v1Router.DELETE("/instances/:instance_name/", v1.DeleteInstance, AdminUserAllowed())
		v1Router.PATCH("/instances/:instance_name/", v1.UpdateInstance, AdminUserAllowed())

		// rule template
		v1Router.POST("/rule_templates", v1.CreateRuleTemplate, AdminUserAllowed())
		v1Router.PATCH("/rule_templates/:rule_template_name", v1.UpdateRuleTemplate, AdminUserAllowed())
		v1Router.DELETE("/rule_templates/:rule_template_name", v1.DeleteRuleTemplate, AdminUserAllowed())
	}

	// user
	v1Router.GET("/user", v1.GetCurrentUser)

	// instance
	v1Router.GET("/instances", v1.GetInstances)
	v1Router.GET("/instances/:instance_name/", v1.GetInstance)
	v1Router.GET("/instances/:instance_name/connection", v1.CheckInstanceIsConnectableByName)
	v1Router.POST("/instance_connection", v1.CheckInstanceIsConnectable)
	v1Router.GET("/instances/:instance_name/schemas", v1.GetInstanceSchemas)
	v1Router.GET("/instance_tips", v1.GetInstanceTips)

	// rule template
	v1Router.GET("/rule_templates", v1.GetRuleTemplates)
	v1Router.GET("/rule_template_tips", v1.GetRuleTemplateTips)
	v1Router.GET("/rule_templates/:rule_template_name/", v1.GetRuleTemplate)

	//rule
	v1Router.GET("/rules", v1.GetRules)

	address := fmt.Sprintf(":%v", port)
	log.Logger().Infof("starting http server on %s", address)
	log.Logger().Fatal(e.Start(address))
	close(exitChan)
}

// JWTTokenAdapter is a `echo` middleware,　by rewriting the header, the jwt token support header
// "Authorization: {token}" and "Authorization: Bearer {token}".
func JWTTokenAdapter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get(echo.HeaderAuthorization)
			if auth != "" && !strings.HasPrefix(auth, middleware.DefaultJWTConfig.AuthScheme) {
				c.Request().Header.Set(echo.HeaderAuthorization,
					fmt.Sprintf("%s %s", middleware.DefaultJWTConfig.AuthScheme, auth))
			}
			return next(c)
		}
	}
}

// AdminUserAllowed is a `echo` middleware, only allow admin user to access next.
func AdminUserAllowed() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if controller.GetUserName(c) == "admin" {
				return next(c)
			}
			return echo.NewHTTPError(http.StatusForbidden)
		}
	}
}
