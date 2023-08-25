package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	dmsRegister "github.com/actiontech/dms/pkg/dms-common/register"

	// "github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper"
	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	v2 "github.com/actiontech/sqle/sqle/api/controller/v2"
	sqleMiddleware "github.com/actiontech/sqle/sqle/api/middleware"
	"github.com/actiontech/sqle/sqle/config"
	_ "github.com/actiontech/sqle/sqle/docs"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"

	"github.com/facebookgo/grace/gracenet"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

const (
	apiV1 = "v1"
	apiV2 = "v2"
)

// @title Sqle API Docs
// @version 1.0
// @description This is a sample server for dev.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /
func StartApi(net *gracenet.Net, exitChan chan struct{}, config config.SqleConfig) {
	defer close(exitChan)

	e := echo.New()
	output := log.NewRotateFile(config.LogPath, "/api.log", config.LogMaxSizeMB /*MB*/, config.LogMaxBackupNumber)
	defer output.Close()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: output,
	}))
	e.HideBanner = true
	e.HidePort = true

	// custom handler http error
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if _, ok := err.(*errors.CodeError); ok {
			if err = controller.JSONBaseErrorReq(c, err); err != nil {
				log.NewEntry().Error("send json error response failed, error:", err)
			}
		} else {
			e.DefaultHTTPErrorHandler(err, c)
		}
	}

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// e.POST("/v1/login", v1.LoginV1)
	// e.POST("/v2/login", v2.LoginV2)

	// the operation of obtaining the basic information of the platform should be for all users, not the users who log in to the platform
	e.GET("/v1/basic_info", v1.GetSQLEInfo)
	e.GET("/v1/static/logo", v1.GetLogo)

	// oauth2 interface does not require login authentication
	// e.GET("/v1/configurations/oauth2/tips", v1.GetOauth2Tips)
	// e.GET("/v1/oauth2/link", v1.Oauth2Link)
	// e.GET("/v1/oauth2/callback", v1.Oauth2Callback)
	// e.POST("/v1/oauth2/user/bind", v1.BindOauth2User)

	v1Router := e.Group(apiV1)
	v1Router.Use(sqleMiddleware.JWTTokenAdapter(), sqleMiddleware.JWTWithConfig(dmsV1.JwtSigningKey), sqleMiddleware.VerifyUserIsDisabled(), sqleMiddleware.OperationLogRecord())
	v2Router := e.Group(apiV2)
	v2Router.Use(sqleMiddleware.JWTTokenAdapter(), sqleMiddleware.JWTWithConfig(dmsV1.JwtSigningKey), sqleMiddleware.VerifyUserIsDisabled(), sqleMiddleware.OperationLogRecord())

	// v1 admin api, just admin user can access.
	{
		// user
		// v1Router.GET("/users", v1.GetUsers, AdminUserAllowed())
		// v1Router.POST("/users", v1.CreateUser, AdminUserAllowed())
		// v1Router.GET("/users/:user_name/", v1.GetUser, AdminUserAllowed())
		// v1Router.PATCH("/users/:user_name/", v1.UpdateUser, AdminUserAllowed())
		// v1Router.DELETE("/users/:user_name/", v1.DeleteUser, AdminUserAllowed())
		// v1Router.PATCH("/users/:user_name/password", v1.UpdateOtherUserPassword, AdminUserAllowed())

		// user_group
		// v1Router.POST("/user_groups", v1.CreateUserGroup, AdminUserAllowed())
		// v1Router.GET("/user_groups", v1.GetUserGroups, AdminUserAllowed())
		// v1Router.DELETE("/user_groups/:user_group_name/", v1.DeleteUserGroup, AdminUserAllowed())
		// v1Router.PATCH("/user_groups/:user_group_name/", v1.UpdateUserGroup, AdminUserAllowed())

		// role
		// v1Router.GET("/roles", v1.GetRoles, AdminUserAllowed())
		// v1Router.POST("/roles", v1.CreateRole, AdminUserAllowed())
		// v1Router.PATCH("/roles/:role_name/", v1.UpdateRole, AdminUserAllowed())
		// v1Router.DELETE("/roles/:role_name/", v1.DeleteRole, AdminUserAllowed())

		// rule template
		v1Router.POST("/rule_templates", v1.CreateRuleTemplate, sqleMiddleware.AdminUserAllowed())
		v1Router.POST("/rule_templates/:rule_template_name/clone", v1.CloneRuleTemplate, sqleMiddleware.AdminUserAllowed())
		v1Router.PATCH("/rule_templates/:rule_template_name/", v1.UpdateRuleTemplate, sqleMiddleware.AdminUserAllowed())
		v1Router.DELETE("/rule_templates/:rule_template_name/", v1.DeleteRuleTemplate, sqleMiddleware.AdminUserAllowed())
		v1Router.GET("/rule_templates/:rule_template_name/export", v1.ExportRuleTemplateFile, sqleMiddleware.AdminUserAllowed())
		v1Router.DELETE("/custom_rules/:rule_id", v1.DeleteCustomRule, sqleMiddleware.AdminUserAllowed())
		v1Router.POST("/custom_rules", v1.CreateCustomRule, sqleMiddleware.AdminUserAllowed())
		v1Router.PATCH("/custom_rules/:rule_id", v1.UpdateCustomRule, sqleMiddleware.AdminUserAllowed())

		// configurations
		// v1Router.GET("/configurations/ldap", v1.GetLDAPConfiguration, AdminUserAllowed())
		// v1Router.PATCH("/configurations/ldap", v1.UpdateLDAPConfiguration, AdminUserAllowed())
		// v1Router.GET("/configurations/smtp", v1.GetSMTPConfiguration, sqleMiddleware.AdminUserAllowed())
		// v1Router.POST("/configurations/smtp/test", v1.TestSMTPConfigurationV1, sqleMiddleware.AdminUserAllowed())
		// v1Router.PATCH("/configurations/smtp", v1.UpdateSMTPConfiguration, sqleMiddleware.AdminUserAllowed())
		// v1Router.GET("/configurations/wechat", v1.GetWeChatConfiguration, sqleMiddleware.AdminUserAllowed())
		// v1Router.PATCH("/configurations/wechat", v1.UpdateWeChatConfigurationV1, sqleMiddleware.AdminUserAllowed())
		// v1Router.POST("/configurations/wechat/test", v1.TestWeChatConfigurationV1, sqleMiddleware.AdminUserAllowed())
		v1Router.GET("/configurations/ding_talk", v1.GetDingTalkConfigurationV1, sqleMiddleware.AdminUserAllowed())
		v1Router.PATCH("/configurations/ding_talk", v1.UpdateDingTalkConfigurationV1, sqleMiddleware.AdminUserAllowed())
		v1Router.POST("/configurations/ding_talk/test", v1.TestDingTalkConfigV1, sqleMiddleware.AdminUserAllowed())
		// v1Router.GET("/configurations/feishu", v1.GetFeishuConfigurationV1, sqleMiddleware.AdminUserAllowed())
		// v1Router.PATCH("/configurations/feishu", v1.UpdateFeishuConfigurationV1, sqleMiddleware.AdminUserAllowed())
		// v1Router.POST("/configurations/feishu/test", v1.TestFeishuConfigV1, sqleMiddleware.AdminUserAllowed())
		v1Router.GET("/configurations/system_variables", v1.GetSystemVariables, sqleMiddleware.AdminUserAllowed())
		v1Router.PATCH("/configurations/system_variables", v1.UpdateSystemVariables, sqleMiddleware.AdminUserAllowed())
		v1Router.GET("/configurations/license", v1.GetLicense, sqleMiddleware.AdminUserAllowed())
		v1Router.POST("/configurations/license", v1.SetLicense, sqleMiddleware.AdminUserAllowed())
		v1Router.GET("/configurations/license/info", v1.GetSQLELicenseInfo, sqleMiddleware.AdminUserAllowed())
		v1Router.POST("/configurations/license/check", v1.CheckLicense, sqleMiddleware.AdminUserAllowed())
		// v1Router.GET("/configurations/oauth2", v1.GetOauth2Configuration, AdminUserAllowed())
		// v1Router.PATCH("/configurations/oauth2", v1.UpdateOauth2Configuration, AdminUserAllowed())
		v1Router.POST("/configurations/personalise/logo", v1.UploadLogo, sqleMiddleware.AdminUserAllowed())
		v1Router.PATCH("/configurations/personalise", v1.UpdatePersonaliseConfig, sqleMiddleware.AdminUserAllowed())
		// v1Router.PATCH("/configurations/webhook", v1.UpdateWorkflowWebHookConfig, sqleMiddleware.AdminUserAllowed())
		// v1Router.GET("/configurations/webhook", v1.GetWorkflowWebHookConfig, sqleMiddleware.AdminUserAllowed())
		// v1Router.POST("/configurations/webhook/test", v1.TestWorkflowWebHookConfig, sqleMiddleware.AdminUserAllowed())

		// statistic
		// v1Router.GET("/statistic/instances/type_percent", v1.GetInstancesTypePercentV1, AdminUserAllowed())
		// v1Router.GET("/statistic/instances/sql_average_execution_time", v1.GetSqlAverageExecutionTimeV1, AdminUserAllowed())
		// v1Router.GET("/statistic/instances/sql_execution_fail_percent", v1.GetSqlExecutionFailPercentV1, AdminUserAllowed())
		// v1Router.GET("/statistic/license/usage", v1.GetLicenseUsageV1, AdminUserAllowed())
		v1Router.GET("/statistic/workflows/rejected_percent_group_by_creator", v1.GetWorkflowRejectedPercentGroupByCreatorV1, sqleMiddleware.AdminUserAllowed())
		//v1Router.GET("/statistic/workflows/rejected_percent_group_by_instance", v1.GetWorkflowRejectedPercentGroupByInstanceV1, AdminUserAllowed())
		v1Router.GET("/statistic/workflows/counts", v1.GetWorkflowCountsV1, sqleMiddleware.AdminUserAllowed())
		v1Router.GET("/statistic/workflows/duration_of_waiting_for_audit", v1.GetWorkflowDurationOfWaitingForAuditV1, sqleMiddleware.AdminUserAllowed())
		//v1Router.GET("/statistic/workflows/duration_of_waiting_for_execution", v1.GetWorkflowDurationOfWaitingForExecutionV1, AdminUserAllowed())
		//v1Router.GET("/statistic/workflows/pass_percent", v1.GetWorkflowPassPercentV1, AdminUserAllowed())
		v1Router.GET("/statistic/workflows/audit_pass_percent", v1.GetWorkflowAuditPassPercentV1, sqleMiddleware.AdminUserAllowed())
		v1Router.GET("/statistic/workflows/each_day_counts", v1.GetWorkflowCreatedCountsEachDayV1, sqleMiddleware.AdminUserAllowed())
		v1Router.GET("/statistic/workflows/status_count", v1.GetWorkflowStatusCountV1, sqleMiddleware.AdminUserAllowed())
		v1Router.GET("/statistic/workflows/instance_type_percent", v1.GetWorkflowPercentCountedByInstanceTypeV1, sqleMiddleware.AdminUserAllowed())
		v1Router.GET("/static/instance_logo", v1.GetInstanceTypeLogo)

		// sync instance
		// v1Router.POST("/sync_instances", v1.CreateSyncInstanceTask, AdminUserAllowed())
		// v1Router.GET("/sync_instances", v1.GetSyncInstanceTaskList, AdminUserAllowed())
		// v1Router.GET("/sync_instances/:task_id/", v1.GetSyncInstanceTask, AdminUserAllowed())
		// v1Router.PATCH("/sync_instances/:task_id/", v1.UpdateSyncInstanceTask, AdminUserAllowed())
		// v1Router.GET("/sync_instances/source_tips", v1.GetSyncTaskSourceTips, AdminUserAllowed())
		// v1Router.DELETE("/sync_instances/:task_id/", v1.DeleteSyncInstanceTask, AdminUserAllowed())
		// v1Router.POST("/sync_instances/:task_id/trigger", v1.TriggerSyncInstance, AdminUserAllowed())

		// operation record
		v1Router.GET("/operation_records/operation_type_names", v1.GetOperationTypeNameList, sqleMiddleware.AdminUserAllowed())
		v1Router.GET("/operation_records/operation_actions", v1.GetOperationActionList, sqleMiddleware.AdminUserAllowed())
		v1Router.GET("/operation_records", v1.GetOperationRecordListV1, sqleMiddleware.AdminUserAllowed())
		v1Router.GET("/operation_records/exports", v1.GetExportOperationRecordListV1, sqleMiddleware.AdminUserAllowed())

		// other
		// v1Router.GET("/management_permissions", v1.GetManagementPermissions, AdminUserAllowed())

		// statistic
		v1Router.GET("/projects/:project_name/statistic/audited_sqls", v1.StatisticsAuditedSQLV1)
		// 数据资源更新处理 内部调用
		v1Router.POST("/data_resource/handle", v1.OperateDataResourceHandle, sqleMiddleware.AdminUserAllowed())
	}

	// auth
	// v1Router.POST("/logout", v1.LogoutV1)

	// project admin router
	v1ProjectAdminRouter := v1Router.Group("/projects", sqleMiddleware.ProjectAdminUserAllowed())
	{
		// audit whitelist
		v1ProjectAdminRouter.POST("/:project_name/audit_whitelist", v1.CreateAuditWhitelist)
		v1ProjectAdminRouter.PATCH("/:project_name/audit_whitelist/:audit_whitelist_id/", v1.UpdateAuditWhitelistById)
		v1ProjectAdminRouter.DELETE("/:project_name/audit_whitelist/:audit_whitelist_id/", v1.DeleteAuditWhitelistById)
		// rule template
		v1ProjectAdminRouter.POST("/:project_name/rule_templates", v1.CreateProjectRuleTemplate)
		v1ProjectAdminRouter.PATCH("/:project_name/rule_templates/:rule_template_name/", v1.UpdateProjectRuleTemplate)
		v1ProjectAdminRouter.DELETE("/:project_name/rule_templates/:rule_template_name/", v1.DeleteProjectRuleTemplate)
		v1ProjectAdminRouter.POST("/:project_name/rule_templates/:rule_template_name/clone", v1.CloneProjectRuleTemplate)
		// workflow template
		v1ProjectAdminRouter.PATCH("/:project_name/workflow_template", v1.UpdateWorkflowTemplate)
	}

	// project member router
	v1ProjectRouter := v1Router.Group("/projects", sqleMiddleware.ProjectMemberAllowed())
	{
		// statistic
		v1ProjectRouter.GET("/:project_name/statistics", v1.GetProjectStatisticsV1)
		v1ProjectRouter.GET("/:project_name/statistics", v1.GetProjectStatisticsV1)
		v1ProjectRouter.GET("/:project_name/statistic/workflow_status", v1.StatisticWorkflowStatusV1)
		v1ProjectRouter.GET("/:project_name/statistic/risk_workflow", v1.StatisticRiskWorkflowV1)
		v1ProjectRouter.GET("/:project_name/statistic/audit_plans", v1.StatisticAuditPlanV1)
		v1ProjectRouter.GET("/:project_name/statistic/risk_audit_plans", v1.GetRiskAuditPlanV1)
		v1ProjectRouter.GET("/:project_name/statistic/role_user", v1.GetRoleUserCountV1)
		v1ProjectRouter.GET("/:project_name/statistic/project_score", v1.GetProjectScoreV1)
		v1ProjectRouter.GET("/:project_name/statistic/instance_health", v1.GetInstanceHealthV1)
		// audit whitelist
		v1ProjectRouter.GET("/:project_name/audit_whitelist", v1.GetSqlWhitelist)
		// instance
		// v1ProjectRouter.GET("/:project_name/instances/:instance_name/", v1.GetInstance)
		v1ProjectRouter.GET("/:project_name/instances/:instance_name/connection", v1.CheckInstanceIsConnectableByName)
		v1ProjectRouter.POST("/:project_name/instances/connections", v1.BatchCheckInstanceConnections)
		v1ProjectRouter.GET("/:project_name/instances/:instance_name/schemas", v1.GetInstanceSchemas)
		v1ProjectRouter.GET("/:project_name/instance_tips", v1.GetInstanceTips)
		v1ProjectRouter.GET("/:project_name/instances/:instance_name/rules", v1.GetInstanceRules)
		v1ProjectRouter.GET("/:project_name/instances/:instance_name/schemas/:schema_name/tables", v1.ListTableBySchema)
		v1ProjectRouter.GET("/:project_name/instances/:instance_name/schemas/:schema_name/tables/:table_name/metadata", v1.GetTableMetadata)
		// rule template
		v1ProjectRouter.GET("/:project_name/rule_templates/:rule_template_name/", v1.GetProjectRuleTemplate)
		v1ProjectRouter.GET("/:project_name/rule_templates", v1.GetProjectRuleTemplates)
		v1ProjectRouter.GET("/:project_name/rule_template_tips", v1.GetProjectRuleTemplateTips)
		v1ProjectRouter.GET("/:project_name/rule_templates/:rule_template_name/export", v1.ExportProjectRuleTemplateFile)
		// workflow template
		v1ProjectRouter.GET("/:project_name/workflow_template", v1.GetWorkflowTemplate)
		// workflow
		v1ProjectRouter.POST("/:project_name/workflows", DeprecatedBy(apiV2))
		v1ProjectRouter.GET("/:project_name/workflows/:workflow_name/", DeprecatedBy(apiV2))
		v1ProjectRouter.GET("/:project_name/workflows", v1.GetWorkflowsV1)
		v1ProjectRouter.POST("/:project_name/workflows/:workflow_name/steps/:workflow_step_id/approve", DeprecatedBy(apiV2))
		v1ProjectRouter.POST("/:project_name/workflows/:workflow_name/steps/:workflow_step_id/reject", DeprecatedBy(apiV2))
		v1ProjectRouter.POST("/:project_name/workflows/:workflow_name/cancel", DeprecatedBy(apiV2))
		v1ProjectRouter.POST("/:project_name/workflows/cancel", DeprecatedBy(apiV2))
		v1ProjectRouter.POST("/:project_name/workflows/complete", DeprecatedBy(apiV2))
		v1ProjectRouter.POST("/:project_name/workflows/:workflow_name/tasks/:task_id/execute", DeprecatedBy(apiV2))
		v1ProjectRouter.POST("/:project_name/workflows/:workflow_id/tasks/:task_id/terminate", v1.TerminateSingleTaskByWorkflowV1)
		v1ProjectRouter.GET("/:project_name/workflows/:workflow_name/tasks", DeprecatedBy(apiV2))
		v1ProjectRouter.POST("/:project_name/workflows/:workflow_name/tasks/execute", DeprecatedBy(apiV2))
		v1ProjectRouter.POST("/:project_name/workflows/:workflow_id/tasks/terminate", v1.TerminateMultipleTaskByWorkflowV1)
		v1ProjectRouter.PUT("/:project_name/workflows/:workflow_name/tasks/:task_id/schedule", DeprecatedBy(apiV2))
		v1ProjectRouter.PATCH("/:project_name/workflows/:workflow_name/", DeprecatedBy(apiV2))
		v1ProjectRouter.GET("/:project_name/workflows/exports", v1.ExportWorkflowV1)
		// audit plan; 智能扫描任务
		v1ProjectRouter.POST("/:project_name/audit_plans", v1.CreateAuditPlan)
		v1ProjectRouter.GET("/:project_name/audit_plans", v1.GetAuditPlans)
		v1ProjectRouter.DELETE("/:project_name/audit_plans/:audit_plan_name/", v1.DeleteAuditPlan)
		v1ProjectRouter.PATCH("/:project_name/audit_plans/:audit_plan_name/", v1.UpdateAuditPlan)
		v1ProjectRouter.GET("/:project_name/audit_plans/:audit_plan_name/", v1.GetAuditPlan)
		v1ProjectRouter.GET("/:project_name/audit_plans/:audit_plan_name/reports", v1.GetAuditPlanReports)
		v1ProjectRouter.GET("/:project_name/audit_plans/:audit_plan_name/reports/:audit_plan_report_id/", v1.GetAuditPlanReport)
		v1ProjectRouter.GET("/:project_name/audit_plans/:audit_plan_name/sqls", v1.GetAuditPlanSQLs)
		v1ProjectRouter.POST("/:project_name/audit_plans/:audit_plan_name/trigger", v1.TriggerAuditPlan)
		v1ProjectRouter.PATCH("/:project_name/audit_plans/:audit_plan_name/notify_config", v1.UpdateAuditPlanNotifyConfig)
		v1ProjectRouter.GET("/:project_name/audit_plans/:audit_plan_name/notify_config", v1.GetAuditPlanNotifyConfig)
		v1ProjectRouter.GET("/:project_name/audit_plans/:audit_plan_name/notify_config/test", v1.TestAuditPlanNotifyConfig)
		v1ProjectRouter.GET("/:project_name/audit_plans/:audit_plan_name/reports/:audit_plan_report_id/sqls/:number/analysis", v1.GetAuditPlanAnalysisData)
		v1ProjectRouter.GET("/:project_name/audit_plans/:audit_plan_name/reports/:audit_plan_report_id/sqls", v1.GetAuditPlanReportSQLsV1)
	}

	// project member router
	v2ProjectRouter := v2Router.Group("/projects", sqleMiddleware.ProjectMemberAllowed())
	{
		// workflow
		v2ProjectRouter.POST("/:project_name/workflows", v2.CreateWorkflowV2)
		v2ProjectRouter.GET("/:project_name/workflows/:workflow_id/", v2.GetWorkflowV2)
		v2ProjectRouter.POST("/:project_name/workflows/:workflow_id/steps/:workflow_step_id/approve", v2.ApproveWorkflowV2)
		v2ProjectRouter.POST("/:project_name/workflows/:workflow_id/steps/:workflow_step_id/reject", v2.RejectWorkflowV2)
		v2ProjectRouter.POST("/:project_name/workflows/:workflow_id/cancel", v2.CancelWorkflowV2)
		v2ProjectRouter.POST("/:project_name/workflows/cancel", v2.BatchCancelWorkflowsV2)
		v2ProjectRouter.POST("/:project_name/workflows/complete", v2.BatchCompleteWorkflowsV2)
		v2ProjectRouter.POST("/:project_name/workflows/:workflow_id/tasks/:task_id/execute", v2.ExecuteOneTaskOnWorkflowV2)
		v2ProjectRouter.GET("/:project_name/workflows/:workflow_id/tasks", v2.GetSummaryOfWorkflowTasksV2)
		v2ProjectRouter.POST("/:project_name/workflows/:workflow_id/tasks/execute", v2.ExecuteTasksOnWorkflowV2)
		v2ProjectRouter.PUT("/:project_name/workflows/:workflow_id/tasks/:task_id/schedule", v2.UpdateWorkflowScheduleV2)
		v2ProjectRouter.PATCH("/:project_name/workflows/:workflow_id/", v2.UpdateWorkflowV2)
		// audit plan; 智能扫描任务
		v2ProjectRouter.GET("/:project_name/audit_plans", v2.GetAuditPlans)
		v2ProjectRouter.GET("/:project_name/audit_plans/:audit_plan_name/reports/:audit_plan_report_id/sqls/:number/analysis", v2.GetAuditPlanAnalysisData)
		v2ProjectRouter.GET("/:project_name/audit_plans/:audit_plan_name/reports/:audit_plan_report_id/sqls", v2.GetAuditPlanReportSQLs)
	}

	// scanner token auth
	v1Router.POST("/projects/:project_name/audit_plans/:audit_plan_name/sqls/full", v1.FullSyncAuditPlanSQLs, sqleMiddleware.ScannerVerifier())
	v1Router.POST("/projects/:project_name/audit_plans/:audit_plan_name/sqls/partial", v1.PartialSyncAuditPlanSQLs, sqleMiddleware.ScannerVerifier())

	// project
	// v1Router.PATCH("/projects/:project_name/", v1.UpdateProjectV1)
	// v1Router.DELETE("/projects/:project_name/", v1.DeleteProjectV1)
	// v1Router.POST("/projects", v1.CreateProjectV1)
	// v1Router.POST("/projects/:project_name/archive", v1.ArchiveProjectV1)
	// v1Router.POST("/projects/:project_name/unarchive", v1.UnarchiveProjectV1)
	// v1Router.GET("/projects", v1.GetProjectListV1)

	// role
	// v1Router.GET("/role_tips", v1.GetRoleTips)

	// user
	// v1Router.GET("/user", v1.GetCurrentUser)
	// v1Router.PATCH("/user", v1.UpdateCurrentUser)
	v1Router.GET("/user_tips", v1.GetUserTips)
	// v1Router.PUT("/user/password", v1.UpdateCurrentUserPassword)
	// v1Router.POST("/projects/:project_name/members", v1.AddMember)
	// v1Router.PATCH("/projects/:project_name/members/:user_name/", v1.UpdateMember)
	// v1Router.DELETE("/projects/:project_name/members/:user_name/", v1.DeleteMember)
	// v1Router.GET("/projects/:project_name/members", v1.GetMembers)
	// v1Router.GET("/projects/:project_name/members/:user_name/", v1.GetMember)
	// v1Router.GET("/projects/:project_name/member_tips", v1.GetMemberTips)

	// user group
	// v1Router.POST("/projects/:project_name/member_groups", v1.AddMemberGroup)
	// v1Router.PATCH("/projects/:project_name/member_groups/:user_group_name/", v1.UpdateMemberGroup)
	// v1Router.DELETE("/projects/:project_name/member_groups/:user_group_name/", v1.DeleteMemberGroup)
	// v1Router.GET("/projects/:project_name/member_groups", v1.GetMemberGroups)
	// v1Router.GET("/projects/:project_name/member_groups/:user_group_name/", v1.GetMemberGroup)
	// v1Router.GET("/user_group_tips", v1.GetUserGroupTips)

	// operations
	// v1Router.GET("/operations", v1.GetOperations)

	// instance
	// v1Router.GET("/projects/:project_name/instances", v1.GetInstances)
	// v2Router.GET("/projects/:project_name/instances", v2.GetInstances)

	v2Router.GET("/projects/:project_name/instances/:instance_name/", v2.GetInstance)

	v1Router.POST("/instance_connection", v1.CheckInstanceIsConnectable) // permission

	// v1Router.POST("/projects/:project_name/instances", DeprecatedBy(apiV2))
	// v2Router.POST("/projects/:project_name/instances", v2.CreateInstance)
	// v1Router.GET("/instance_additional_metas", v1.GetInstanceAdditionalMetas)
	// v1Router.DELETE("/projects/:project_name/instances/:instance_name/", v1.DeleteInstance)
	// v1Router.PATCH("/projects/:project_name/instances/:instance_name/", v1.UpdateInstance)

	// 全局 rule template
	v1Router.GET("/rule_templates", v1.GetRuleTemplates)
	v1Router.GET("/rule_template_tips", v1.GetRuleTemplateTips)
	v1Router.GET("/rule_templates/:rule_template_name/", v1.GetRuleTemplate)

	v1Router.POST("/rule_templates/parse", v1.ParseProjectRuleTemplateFile)

	// 全局 workflow
	v1Router.GET("/workflows", v1.GetGlobalWorkflowsV1)

	//rule
	v1Router.GET("/rules", v1.GetRules)
	v1Router.GET("/custom_rules", v1.GetCustomRules)
	v1Router.GET("/custom_rules/:rule_id", v1.GetCustomRule)
	v1Router.GET("/custom_rules/:db_type/rule_types", v1.GetRuleTypeByDBType)

	// task
	v1Router.POST("/projects/:project_name/tasks/audits", v1.CreateAndAuditTask)
	v1Router.GET("/tasks/audits/:task_id/", v1.GetTask)
	v1Router.GET("/tasks/audits/:task_id/sqls", v1.GetTaskSQLs)
	v2Router.GET("/tasks/audits/:task_id/sqls", v2.GetTaskSQLs)
	v1Router.GET("/tasks/audits/:task_id/sql_report", v1.DownloadTaskSQLReportFile)
	v1Router.GET("/tasks/audits/:task_id/sql_file", v1.DownloadTaskSQLFile)
	v1Router.GET("/tasks/audits/:task_id/sql_content", v1.GetAuditTaskSQLContent)
	v1Router.PATCH("/tasks/audits/:task_id/sqls/:number", v1.UpdateAuditTaskSQLs)
	v1Router.GET("/tasks/audits/:task_id/sqls/:number/analysis", v1.GetTaskAnalysisData)
	v2Router.GET("/tasks/audits/:task_id/sqls/:number/analysis", v2.GetTaskAnalysisData)
	v1Router.POST("/projects/:project_name/task_groups", v1.CreateAuditTasksGroupV1)
	v1Router.POST("/task_groups/audit", v1.AuditTaskGroupV1)

	// dashboard
	v1Router.GET("/dashboard", v1.Dashboard)
	v1Router.GET("/dashboard/project_tips", v1.DashboardProjectTipsV1)

	// configurations
	v1Router.GET("/configurations/drivers", v1.GetDrivers)
	v2Router.GET("/configurations/drivers", v2.GetDrivers)
	// v1Router.GET("/configurations/sql_query", v1.GetSQLQueryConfiguration)

	// audit plan
	v1Router.GET("/audit_plan_metas", v1.GetAuditPlanMetas)
	v1Router.GET("/audit_plan_types", v1.GetAuditPlanTypes)

	// sql audit
	v1Router.POST("/sql_audit", v1.DirectAudit)
	v2Router.POST("/sql_audit", v2.DirectAudit)

	v1Router.GET("/sql_analysis", v1.DirectGetSQLAnalysis)

	// UI
	e.File("/", "ui/index.html")
	e.Static("/static", "ui/static")
	e.File("/favicon.png", "ui/favicon.png")
	e.GET("/*", func(c echo.Context) error {
		return c.File("ui/index.html")
	})

	address := fmt.Sprintf(":%v", config.SqleServerPort)
	log.Logger().Infof("starting http server on %s", address)

	// start http server
	l, err := net.Listen("tcp4", address)
	if err != nil {
		log.Logger().Fatal(err)
		return
	}
	if config.EnableHttps {
		// Usually, it is easier to create an tls server using echo#StartTLS;
		// but I need create a graceful listener.
		if config.CertFilePath == "" || config.KeyFilePath == "" {
			log.Logger().Fatal("invalid tls configuration")
			return
		}
		tlsConfig := new(tls.Config)
		tlsConfig.Certificates = make([]tls.Certificate, 1)
		tlsConfig.Certificates[0], err = tls.LoadX509KeyPair(config.CertFilePath, config.KeyFilePath)
		if err != nil {
			log.Logger().Fatal("load x509 key pair failed, error:", err)
			return
		}
		e.TLSServer.TLSConfig = tlsConfig
		e.TLSListener = tls.NewListener(l, tlsConfig)

		log.Logger().Fatal(e.StartServer(e.TLSServer))
	} else {
		e.Listener = l
		log.Logger().Fatal(e.Start(""))
	}
}

// DeprecatedBy is a controller used to mark deprecated and used to replace the original controller.
func DeprecatedBy(version string) func(echo.Context) error {
	return func(ctx echo.Context) error {
		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf(
			"the API has been deprecated, please using the %s version", version))
	}
}

func RegisterAsDMSTarget(config config.SqleConfig) error {
	controller.InitDMSServerAddress(config.DMSServerAddress)
	ctx := context.Background()
	// 向DMS注册反向代理
	if err := dmsRegister.RegisterDMSProxyTarget(ctx, controller.GetDMSServerAddress(), "sqle-api", fmt.Sprintf("http://%v:%v", config.SqleServerHost, config.SqleServerPort) /* TODO https的处理*/, []string{"/sqle/v"}); nil != err {
		return fmt.Errorf("failed to register dms proxy target: %v", err)
	}
	// 注册校验接口
	if err := dmsRegister.RegisterDMSPlugin(ctx, controller.GetDMSServerAddress(), &dmsV1.Plugin{
		Name:                         "sqle",
		OperateDataResourceHandleUrl: fmt.Sprintf("http://%s:%d/%s/%s", config.SqleServerHost, config.SqleServerPort, apiV1, "data_resource/handle"),
	}); err != nil {
		return fmt.Errorf("failed to register  dms plugin for operation data source handle")
	}

	return nil
}
