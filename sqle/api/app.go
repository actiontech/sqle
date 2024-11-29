package api

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/actiontech/dms/pkg/dms-common/api/accesstoken"
	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	jwtPkg "github.com/actiontech/dms/pkg/dms-common/api/jwt"
	"github.com/actiontech/dms/pkg/dms-common/i18nPkg"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/locale"

	// "github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper"
	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	v2 "github.com/actiontech/sqle/sqle/api/controller/v2"
	v3 "github.com/actiontech/sqle/sqle/api/controller/v3"
	sqleMiddleware "github.com/actiontech/sqle/sqle/api/middleware"
	"github.com/actiontech/sqle/sqle/config"
	_ "github.com/actiontech/sqle/sqle/docs"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"

	"github.com/facebookgo/grace/gracenet"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	apiV1 = "v1"
	apiV2 = "v2"
	apiV3 = "v3"
)

type restApi struct {
	method        string
	path          string
	handlerFn     echo.HandlerFunc
	middleWareFns []echo.MiddlewareFunc
}

var restApis []restApi

func LoadRestApi(method string, path string, handlerFn echo.HandlerFunc, middleWareFns ...echo.MiddlewareFunc) {
	restApis = append(restApis, restApi{
		method:        method,
		path:          path,
		handlerFn:     handlerFn,
		middleWareFns: middleWareFns,
	})
}

func addCustomApis(e *echo.Group, apis []restApi) error {
	for _, api := range apis {
		switch api.method {
		case http.MethodGet:
			e.GET(api.path, api.handlerFn, api.middleWareFns...)
		case http.MethodPost:
			e.POST(api.path, api.handlerFn, api.middleWareFns...)
		case http.MethodPatch:
			e.PATCH(api.path, api.handlerFn, api.middleWareFns...)
		case http.MethodDelete:
			e.DELETE(api.path, api.handlerFn, api.middleWareFns...)
		case http.MethodPut:
			e.PUT(api.path, api.handlerFn, api.middleWareFns...)
		default:
			return fmt.Errorf("unsupported http method")
		}
	}
	return nil
}

// @title Sqle API Docs
// @version 1.0
// @description This is a sample server for dev.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /sqle
func StartApi(net *gracenet.Net, exitChan chan struct{}, config *config.SqleOptions, swaggerYaml []byte) {
	defer close(exitChan)

	e := echo.New()
	output := log.NewRotateFile(config.Service.LogPath, "/api.log", config.Service.LogMaxSizeMB /*MB*/, config.Service.LogMaxBackupNumber)
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

	e.GET("/swagger_file", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct {
			Content []byte `json:"content"`
		}{
			Content: swaggerYaml,
		})
	})

	v1Router := e.Group(apiV1)
	v1Router.Use(sqleMiddleware.JWTTokenAdapter(), sqleMiddleware.JWTWithConfig(dmsV1.JwtSigningKey), sqleMiddleware.VerifyUserIsDisabled(), locale.Bundle.EchoMiddlewareByCustomFunc(dms.GetCurrentUserLanguage, i18nPkg.GetLangByAcceptLanguage), sqleMiddleware.OperationLogRecord(), accesstoken.CheckLatestAccessToken(controller.GetDMSServerAddress(), jwtPkg.GetTokenDetailFromContextWithOldJwt))
	v2Router := e.Group(apiV2)
	v2Router.Use(sqleMiddleware.JWTTokenAdapter(), sqleMiddleware.JWTWithConfig(dmsV1.JwtSigningKey), sqleMiddleware.VerifyUserIsDisabled(), locale.Bundle.EchoMiddlewareByCustomFunc(dms.GetCurrentUserLanguage, i18nPkg.GetLangByAcceptLanguage), sqleMiddleware.OperationLogRecord(), accesstoken.CheckLatestAccessToken(controller.GetDMSServerAddress(), jwtPkg.GetTokenDetailFromContextWithOldJwt))
	v3Router := e.Group(apiV3)
	v3Router.Use(sqleMiddleware.JWTTokenAdapter(), sqleMiddleware.JWTWithConfig(dmsV1.JwtSigningKey), sqleMiddleware.VerifyUserIsDisabled(), locale.Bundle.EchoMiddlewareByCustomFunc(dms.GetCurrentUserLanguage, i18nPkg.GetLangByAcceptLanguage), sqleMiddleware.OperationLogRecord(), accesstoken.CheckLatestAccessToken(controller.GetDMSServerAddress(), jwtPkg.GetTokenDetailFromContextWithOldJwt))
	// v1 admin api, just admin user can access.
	{
		// rule template
		v1Router.POST("/rule_templates", v1.CreateRuleTemplate, sqleMiddleware.OpGlobalAllowed())
		v1Router.POST("/rule_templates/:rule_template_name/clone", v1.CloneRuleTemplate, sqleMiddleware.OpGlobalAllowed())
		v1Router.PATCH("/rule_templates/:rule_template_name/", v1.UpdateRuleTemplate, sqleMiddleware.OpGlobalAllowed())
		v1Router.DELETE("/rule_templates/:rule_template_name/", v1.DeleteRuleTemplate, sqleMiddleware.OpGlobalAllowed())
		v1Router.GET("/rule_templates/:rule_template_name/export", v1.ExportRuleTemplateFile, sqleMiddleware.ViewGlobalAllowed())
		v1Router.DELETE("/custom_rules/:rule_id", v1.DeleteCustomRule, sqleMiddleware.OpGlobalAllowed())
		v1Router.POST("/custom_rules", v1.CreateCustomRule, sqleMiddleware.OpGlobalAllowed())
		v1Router.PATCH("/custom_rules/:rule_id", v1.UpdateCustomRule, sqleMiddleware.OpGlobalAllowed())
		v1Router.PATCH("/rule_knowledge/db_types/:db_type/rules/:rule_name/", v1.UpdateRuleKnowledgeV1, sqleMiddleware.OpGlobalAllowed())
		v1Router.PATCH("/rule_knowledge/db_types/:db_type/custom_rules/:rule_name/", v1.UpdateCustomRuleKnowledgeV1, sqleMiddleware.OpGlobalAllowed())
		// configurations
		v1Router.GET("/configurations/ding_talk", v1.GetDingTalkConfigurationV1, sqleMiddleware.ViewGlobalAllowed())
		v1Router.PATCH("/configurations/ding_talk", v1.UpdateDingTalkConfigurationV1, sqleMiddleware.OpGlobalAllowed())
		v1Router.POST("/configurations/ding_talk/test", v1.TestDingTalkConfigV1, sqleMiddleware.OpGlobalAllowed())
		v1Router.GET("/configurations/system_variables", v1.GetSystemVariables, sqleMiddleware.ViewGlobalAllowed())
		v1Router.PATCH("/configurations/system_variables", v1.UpdateSystemVariables, sqleMiddleware.OpGlobalAllowed())
		v1Router.GET("/configurations/license", v1.GetLicense, sqleMiddleware.ViewGlobalAllowed())
		v1Router.POST("/configurations/license", v1.SetLicense, sqleMiddleware.OpGlobalAllowed())
		v1Router.GET("/configurations/license/info", v1.GetSQLELicenseInfo, sqleMiddleware.ViewGlobalAllowed())
		v1Router.POST("/configurations/license/check", v1.CheckLicense, sqleMiddleware.OpGlobalAllowed())
		v1Router.PATCH("/configurations/feishu_audit", v1.UpdateFeishuAuditConfigurationV1, sqleMiddleware.OpGlobalAllowed())
		v1Router.GET("/configurations/feishu_audit", v1.GetFeishuAuditConfigurationV1, sqleMiddleware.ViewGlobalAllowed())
		v1Router.POST("/configurations/feishu_audit/test", v1.TestFeishuAuditConfigV1, sqleMiddleware.OpGlobalAllowed())
		v1Router.PATCH("/configurations/wechat_audit", v1.UpdateWechatAuditConfigurationV1, sqleMiddleware.OpGlobalAllowed())
		v1Router.GET("/configurations/wechat_audit", v1.GetWechatAuditConfigurationV1, sqleMiddleware.ViewGlobalAllowed())
		v1Router.POST("/configurations/wechat_audit/test", v1.TestWechatAuditConfigV1, sqleMiddleware.OpGlobalAllowed())

		// statistic
		v1Router.GET("/statistic/instances/type_percent", v1.GetInstancesTypePercentV1, sqleMiddleware.ViewGlobalAllowed())
		v1Router.GET("/statistic/instances/sql_average_execution_time", v1.GetSqlAverageExecutionTimeV1, sqleMiddleware.ViewGlobalAllowed())
		v1Router.GET("/statistic/license/usage", v1.GetLicenseUsageV1, sqleMiddleware.ViewGlobalAllowed())
		v1Router.GET("/statistic/workflows/rejected_percent_group_by_creator", v1.GetWorkflowRejectedPercentGroupByCreatorV1, sqleMiddleware.ViewGlobalAllowed())
		v1Router.GET("/statistic/workflows/counts", v1.GetWorkflowCountsV1, sqleMiddleware.ViewGlobalAllowed())
		v1Router.GET("/statistic/workflows/duration_of_waiting_for_audit", v1.GetWorkflowDurationOfWaitingForAuditV1, sqleMiddleware.ViewGlobalAllowed())
		v1Router.GET("/statistic/workflows/audit_pass_percent", v1.GetWorkflowAuditPassPercentV1, sqleMiddleware.ViewGlobalAllowed())
		v1Router.GET("/statistic/workflows/each_day_counts", v1.GetWorkflowCreatedCountsEachDayV1, sqleMiddleware.ViewGlobalAllowed())
		v1Router.GET("/statistic/workflows/status_count", v1.GetWorkflowStatusCountV1, sqleMiddleware.ViewGlobalAllowed())
		v1Router.GET("/statistic/workflows/instance_type_percent", v1.GetWorkflowPercentCountedByInstanceTypeV1, sqleMiddleware.ViewGlobalAllowed())

		// operation record
		v1Router.GET("/operation_records/operation_type_names", v1.GetOperationTypeNameList, sqleMiddleware.ViewGlobalAllowed())
		v1Router.GET("/operation_records/operation_actions", v1.GetOperationActionList, sqleMiddleware.ViewGlobalAllowed())
		v1Router.GET("/operation_records", v1.GetOperationRecordListV1, sqleMiddleware.ViewGlobalAllowed())
		v1Router.GET("/operation_records/exports", v1.GetExportOperationRecordListV1, sqleMiddleware.ViewGlobalAllowed())

		// 企业公告
		v1Router.PATCH("/company_notice", v1.UpdateCompanyNotice, sqleMiddleware.OpGlobalAllowed())

		// 内部调用
		v1Router.POST("/data_resource/handle", v1.OperateDataResourceHandle, sqleMiddleware.OpGlobalAllowed())
		v1Router.POST(fmt.Sprintf("%s/connection", dmsV1.InternalDBServiceRouterGroup), v1.CheckInstanceIsConnectable, sqleMiddleware.OpGlobalAllowed())
	}

	// project admin and global manage router
	v1OpProjectRouter := v1Router.Group("/projects", sqleMiddleware.OpProjectAllowed())
	{
		// audit whitelist
		v1OpProjectRouter.POST("/:project_name/audit_whitelist", v1.CreateAuditWhitelist)
		v1OpProjectRouter.PATCH("/:project_name/audit_whitelist/:audit_whitelist_id/", v1.UpdateAuditWhitelistById)
		v1OpProjectRouter.DELETE("/:project_name/audit_whitelist/:audit_whitelist_id/", v1.DeleteAuditWhitelistById)

		// blacklist
		v1OpProjectRouter.POST("/:project_name/blacklist", v1.CreateBlacklist)
		v1OpProjectRouter.DELETE("/:project_name/blacklist/:blacklist_id/", v1.DeleteBlacklist)
		v1OpProjectRouter.PATCH("/:project_name/blacklist/:blacklist_id/", v1.UpdateBlacklist)

		// rule template
		v1OpProjectRouter.POST("/:project_name/rule_templates", v1.CreateProjectRuleTemplate)
		v1OpProjectRouter.PATCH("/:project_name/rule_templates/:rule_template_name/", v1.UpdateProjectRuleTemplate)
		v1OpProjectRouter.DELETE("/:project_name/rule_templates/:rule_template_name/", v1.DeleteProjectRuleTemplate)
		v1OpProjectRouter.POST("/:project_name/rule_templates/:rule_template_name/clone", v1.CloneProjectRuleTemplate)

		// workflow template
		v1OpProjectRouter.PATCH("/:project_name/workflow_template", v1.UpdateWorkflowTemplate)

		// report push
		v1OpProjectRouter.PUT("/:project_name/report_push_configs/:report_push_config_id/", v1.UpdateReportPushConfig)

		// sql version
		v1OpProjectRouter.POST("/:project_name/sql_versions", v1.CreateSqlVersion)
		v1OpProjectRouter.PATCH("/:project_name/sql_versions/:sql_version_id/", v1.UpdateSqlVersion)
		v1OpProjectRouter.DELETE("/:project_name/sql_versions/:sql_version_id/", v1.DeleteSqlVersion)
		v1OpProjectRouter.POST("/:project_name/sql_versions/:sql_version_id/lock", v1.LockSqlVersion)
	}

	// project admin and global view router
	v1ViewProjectRouter := v1Router.Group("/projects", sqleMiddleware.ViewProjectAllowed())
	{
		v1ViewProjectRouter.GET("/:project_name/blacklist", v1.GetBlacklist)
	}

	// project member router
	v1ProjectOpRouter := v1Router.Group("/projects", sqleMiddleware.ProjectMemberOpAllowed())
	{
		// instance
		v1ProjectOpRouter.POST("/:project_name/instances/connections", v1.BatchCheckInstanceConnections)
		// workflow
		v1ProjectOpRouter.POST("/:project_name/workflows", DeprecatedBy(apiV2))

		v1ProjectOpRouter.POST("/:project_name/workflows/:workflow_name/steps/:workflow_step_id/approve", DeprecatedBy(apiV2))
		v1ProjectOpRouter.POST("/:project_name/workflows/:workflow_name/steps/:workflow_step_id/reject", DeprecatedBy(apiV2))
		v1ProjectOpRouter.POST("/:project_name/workflows/:workflow_name/cancel", DeprecatedBy(apiV2))
		v1ProjectOpRouter.POST("/:project_name/workflows/cancel", DeprecatedBy(apiV2))
		v1ProjectOpRouter.POST("/:project_name/workflows/complete", DeprecatedBy(apiV2))
		v1ProjectOpRouter.POST("/:project_name/workflows/:workflow_name/tasks/:task_id/execute", DeprecatedBy(apiV2))
		v1ProjectOpRouter.POST("/:project_name/workflows/:workflow_id/tasks/:task_id/terminate", v1.TerminateSingleTaskByWorkflowV1)
		v1ProjectOpRouter.POST("/:project_name/workflows/:workflow_name/tasks/execute", DeprecatedBy(apiV2))
		v1ProjectOpRouter.POST("/:project_name/workflows/:workflow_id/tasks/terminate", v1.TerminateMultipleTaskByWorkflowV1)
		v1ProjectOpRouter.PUT("/:project_name/workflows/:workflow_name/tasks/:task_id/schedule", DeprecatedBy(apiV2))
		v1ProjectOpRouter.PATCH("/:project_name/workflows/:workflow_name/", DeprecatedBy(apiV2))
		v1ProjectOpRouter.POST("/:project_name/workflows/:workflow_id/tasks/:task_id/order_file", v1.UpdateSqlFileOrderByWorkflowV1)
		v1ProjectOpRouter.GET("/:project_name/workflows/:workflow_id/backup_sqls", v1.GetBackupSqlList)
		v1ProjectOpRouter.POST("/:project_name/workflows/:workflow_id/create_rollback_workflow", v1.CreateRollbackWorkflow)

		// sql version
		v1ProjectOpRouter.POST("/:project_name/sql_versions/:sql_version_id/batch_release_workflows", v1.BatchReleaseWorkflows)
		v1ProjectOpRouter.POST("/:project_name/sql_versions/:sql_version_id/batch_execute_workflows", v1.BatchExecuteWorkflows)
		v1ProjectOpRouter.POST("/:project_name/sql_versions/:sql_version_id/sql_version_stages/:sql_version_stage_id/associate_workflows", v1.BatchAssociateWorkflowsWithVersion)

		// audit plan; 智能扫描任务
		v1ProjectOpRouter.POST("/:project_name/audit_plans", v1.CreateAuditPlan)
		v1ProjectOpRouter.DELETE("/:project_name/audit_plans/:audit_plan_name/", v1.DeleteAuditPlan)
		v1ProjectOpRouter.PATCH("/:project_name/audit_plans/:audit_plan_name/", v1.UpdateAuditPlan)

		v1ProjectOpRouter.POST("/:project_name/audit_plans/:audit_plan_name/trigger", v1.TriggerAuditPlan)
		v1ProjectOpRouter.PATCH("/:project_name/audit_plans/:audit_plan_name/notify_config", v1.UpdateAuditPlanNotifyConfig)

		// scanner token auth
		v1ProjectOpRouter.POST("/:project_name/audit_plans/:audit_plan_name/sqls/full", v1.FullSyncAuditPlanSQLs, sqleMiddleware.ScannerVerifier())
		v1ProjectOpRouter.POST("/:project_name/audit_plans/:audit_plan_name/sqls/partial", v1.PartialSyncAuditPlanSQLs, sqleMiddleware.ScannerVerifier())

		// instance audti plan 实例智能扫描任务
		v1ProjectOpRouter.POST("/:project_name/instance_audit_plans", v1.CreateInstanceAuditPlan)
		v1ProjectOpRouter.DELETE("/:project_name/instance_audit_plans/:instance_audit_plan_id/", v1.DeleteInstanceAuditPlan)
		v1ProjectOpRouter.PUT("/:project_name/instance_audit_plans/:instance_audit_plan_id/", v1.UpdateInstanceAuditPlan)
		v1ProjectOpRouter.PATCH("/:project_name/instance_audit_plans/:instance_audit_plan_id/", v1.UpdateInstanceAuditPlanStatus)

		// audit plan; 智能扫描任务
		v1ProjectOpRouter.DELETE("/:project_name/instance_audit_plans/:instance_audit_plan_id/audit_plans/:audit_plan_id/", v1.DeleteAuditPlanById)
		v1ProjectOpRouter.PATCH("/:project_name/instance_audit_plans/:instance_audit_plan_id/audit_plans/:audit_plan_id/", v1.UpdateAuditPlanStatus)

		v1ProjectOpRouter.POST("/:project_name/instance_audit_plans/:instance_audit_plan_id/audit_plans/:audit_plan_id/sql_data", v1.GetInstanceAuditPlanSQLData)
		v1ProjectOpRouter.POST("/:project_name/instance_audit_plans/:instance_audit_plan_id/audit_plans/:audit_plan_id/sql_export", v1.GetInstanceAuditPlanSQLExport)
		v1ProjectOpRouter.POST("/:project_name/instance_audit_plans/:instance_audit_plan_id/audit_plans/:audit_plan_id/audit", v1.AuditPlanTriggerSqlAudit)

		// sql manager
		v1ProjectOpRouter.PATCH("/:project_name/sql_manages/batch", v1.BatchUpdateSqlManage)

		// sql audit record
		v1ProjectOpRouter.POST("/:project_name/sql_audit_records", v1.CreateSQLAuditRecord)

		v1ProjectOpRouter.PATCH("/:project_name/sql_audit_records/:sql_audit_record_id/", v1.UpdateSQLAuditRecordV1)

		// sql optimization
		v1ProjectOpRouter.POST("/:project_name/sql_optimization_records", v1.SQLOptimizate)

		// task
		v1ProjectOpRouter.POST("/:project_name/tasks/audits", v1.CreateAndAuditTask)

		// pipeline
		v1ProjectOpRouter.POST("/:project_name/pipelines", v1.CreatePipeline)
		v1ProjectOpRouter.DELETE("/:project_name/pipelines/:pipeline_id/", v1.DeletePipeline)
		v1ProjectOpRouter.PATCH("/:project_name/pipelines/:pipeline_id/", v1.UpdatePipeline)

		// database_compare
		v1ProjectOpRouter.POST("/:project_name/database_comparison/execute_comparison", v1.ExecuteDatabaseComparison)
		v1ProjectOpRouter.POST("/:project_name/database_comparison/comparison_statements", v1.GetComparisonStatement)
		v1ProjectOpRouter.POST("/:project_name/database_comparison/modify_sql_statements", v1.GenDatabaseDiffModifySQLs)
	}

	// project member router
	v1ProjectViewRouter := v1Router.Group("/projects", sqleMiddleware.ProjectMemberViewAllowed())
	{
		// statistic
		v1ProjectViewRouter.GET("/:project_name/statistics", v1.GetProjectStatisticsV1)
		v1ProjectViewRouter.GET("/:project_name/statistics", v1.GetProjectStatisticsV1)
		v1ProjectViewRouter.GET("/:project_name/statistic/workflow_status", v1.StatisticWorkflowStatusV1)
		v1ProjectViewRouter.GET("/:project_name/statistic/risk_workflow", v1.StatisticRiskWorkflowV1)
		v1ProjectViewRouter.GET("/:project_name/statistic/audit_plans", v1.StatisticAuditPlanV1)
		v1ProjectViewRouter.GET("/:project_name/statistic/risk_audit_plans", v1.GetRiskAuditPlanV1)
		v1ProjectViewRouter.GET("/:project_name/statistic/role_user", v1.GetRoleUserCountV1)
		v1ProjectViewRouter.GET("/:project_name/statistic/project_score", v1.GetProjectScoreV1)
		v1ProjectViewRouter.GET("/:project_name/statistic/instance_health", v1.GetInstanceHealthV1)
		v1ProjectViewRouter.GET("/:project_name/statistic/audited_sqls", v1.StatisticsAuditedSQLV1)
		v1ProjectViewRouter.GET("/:project_name/statistic/optimization_record_overview", v1.GetOptimizationRecordOverview)
		v1ProjectViewRouter.GET("/:project_name/statistic/optimization_performance_improve_overview", v1.GetDBPerformanceImproveOverview)
		v1ProjectViewRouter.GET("/:project_name/audit_whitelist", v1.GetSqlWhitelist)

		v1ProjectViewRouter.GET("/:project_name/instances/:instance_name/connection", v1.CheckInstanceIsConnectableByName)
		v1ProjectViewRouter.GET("/:project_name/instances/:instance_name/schemas", v1.GetInstanceSchemas)
		v1ProjectViewRouter.GET("/:project_name/instance_tips", v1.GetInstanceTips)
		v1ProjectViewRouter.GET("/:project_name/instances/:instance_name/rules", v1.GetInstanceRules)
		v1ProjectViewRouter.GET("/:project_name/instances/:instance_name/schemas/:schema_name/tables", v1.ListTableBySchema)
		v1ProjectViewRouter.GET("/:project_name/instances/:instance_name/schemas/:schema_name/tables/:table_name/metadata", v1.GetTableMetadata)

		// rule template
		v1ProjectViewRouter.GET("/:project_name/rule_templates/:rule_template_name/", v1.GetProjectRuleTemplate)
		v1ProjectViewRouter.GET("/:project_name/rule_templates", v1.GetProjectRuleTemplates)
		v1ProjectViewRouter.GET("/:project_name/rule_template_tips", v1.GetProjectRuleTemplateTips)
		v1ProjectViewRouter.GET("/:project_name/rule_templates/:rule_template_name/export", v1.ExportProjectRuleTemplateFile)

		// workflow template
		v1ProjectViewRouter.GET("/:project_name/workflow_template", v1.GetWorkflowTemplate)
		v1ProjectViewRouter.GET("/:project_name/workflows/:workflow_name/", DeprecatedBy(apiV2))
		v1ProjectViewRouter.GET("/:project_name/workflows", v1.GetWorkflowsV1)
		v1ProjectViewRouter.GET("/:project_name/workflows/:workflow_name/tasks", DeprecatedBy(apiV2))
		v1ProjectViewRouter.GET("/:project_name/workflows/exports", v1.ExportWorkflowV1)
		v1ProjectViewRouter.GET("/:project_name/workflows/:workflow_id/tasks/:task_id/attachment", v1.GetWorkflowTaskAuditFile)
		v1ProjectViewRouter.GET("/:project_name/sql_versions", v1.GetSqlVersionList)
		v1ProjectViewRouter.GET("/:project_name/sql_versions/:sql_version_id/", v1.GetSqlVersionDetail)
		v1ProjectViewRouter.GET("/:project_name/sql_versions/:sql_version_id/sql_version_stages/:sql_version_stage_id/dependencies", v1.GetDependenciesBetweenStageInstance)
		v1ProjectViewRouter.GET("/:project_name/sql_versions/:sql_version_id/sql_version_stages/:sql_version_stage_id/associate_workflows", v1.GetWorkflowsThatCanBeAssociatedToVersion)
		v1ProjectViewRouter.GET("/:project_name/audit_plans", v1.GetAuditPlans)

		v1ProjectViewRouter.GET("/:project_name/audit_plans/:audit_plan_name/", v1.GetAuditPlan)
		v1ProjectViewRouter.GET("/:project_name/audit_plans/:audit_plan_name/reports", v1.GetAuditPlanReports)
		v1ProjectViewRouter.GET("/:project_name/audit_plans/:audit_plan_name/reports/:audit_plan_report_id/", v1.GetAuditPlanReport)
		v1ProjectViewRouter.GET("/:project_name/audit_plans/:audit_plan_name/sqls", v1.GetAuditPlanSQLs)

		v1ProjectViewRouter.GET("/:project_name/audit_plans/:audit_plan_name/notify_config", v1.GetAuditPlanNotifyConfig)
		v1ProjectViewRouter.GET("/:project_name/audit_plans/:audit_plan_name/notify_config/test", v1.TestAuditPlanNotifyConfig)
		v1ProjectViewRouter.GET("/:project_name/audit_plans/:audit_plan_name/reports/:audit_plan_report_id/sqls/:number/analysis", v1.GetAuditPlanAnalysisData)
		v1ProjectViewRouter.GET("/:project_name/audit_plans/:audit_plan_name/reports/:audit_plan_report_id/sqls", v1.GetAuditPlanReportSQLsV1)
		v1ProjectViewRouter.GET("/:project_name/audit_plans/:audit_plan_name/reports/:audit_plan_report_id/export", v1.ExportAuditPlanReportV1)

		v1ProjectViewRouter.GET("/:project_name/instance_audit_plans", v1.GetInstanceAuditPlans)
		v1ProjectViewRouter.GET("/:project_name/instance_audit_plans/:instance_audit_plan_id", v1.GetInstanceAuditPlanDetail)
		v1ProjectViewRouter.GET("/:project_name/instance_audit_plans/:instance_audit_plan_id/audit_plans", v1.GetInstanceAuditPlanOverview)

		v1ProjectViewRouter.GET("/:project_name/sql_manages", v1.GetSqlManageList)
		v1ProjectViewRouter.GET("/:project_name/sql_manages/exports", v1.ExportSqlManagesV1)
		v1ProjectViewRouter.GET("/:project_name/sql_manages/rule_tips", v1.GetSqlManageRuleTips)
		v1ProjectViewRouter.GET("/:project_name/sql_manages/:sql_manage_id/sql_analysis", v1.GetSqlManageSqlAnalysisV1)

		// sql dev records
		v1ProjectViewRouter.GET("/:project_name/sql_dev_records", v1.GetSqlDEVRecordList)

		v1ProjectViewRouter.GET("/:project_name/sql_audit_records", v1.GetSQLAuditRecordsV1)
		v1ProjectViewRouter.GET("/:project_name/sql_audit_records/:sql_audit_record_id/", v1.GetSQLAuditRecordV1)

		v1ProjectViewRouter.GET("/:project_name/sql_audit_records/tag_tips", v1.GetSQLAuditRecordTagTipsV1)

		v1ProjectViewRouter.GET("/:project_name/sql_optimization_records", v1.GetOptimizationRecords)
		v1ProjectViewRouter.GET("/:project_name/sql_optimization_records/:optimization_record_id/", v1.GetOptimizationRecord)
		v1ProjectViewRouter.GET("/:project_name/sql_optimization_records/:optimization_record_id/sqls", v1.GetOptimizationSQLs)
		v1ProjectViewRouter.GET("/:project_name/sql_optimization_records/:optimization_record_id/sqls/:number/", v1.GetOptimizationSQLDetail)

		v1ProjectViewRouter.GET("/:project_name/report_push_configs", v1.GetReportPushConfigList)

		v1ProjectViewRouter.GET("/:project_name/pipelines", v1.GetPipelines)
		v1ProjectViewRouter.GET("/:project_name/pipelines/:pipeline_id/", v1.GetPipelineDetail)

		v1ProjectViewRouter.GET("/:project_name/instance_audit_plans/:instance_audit_plan_id/audit_plans/:audit_plan_id/sqls", v1.GetInstanceAuditPlanSQLs) // 弃用
		v1ProjectViewRouter.GET("/:project_name/instance_audit_plans/:instance_audit_plan_id/audit_plans/:audit_plan_id/sql_meta", v1.GetInstanceAuditPlanSQLMeta)
		v1ProjectViewRouter.GET("/:project_name/instance_audit_plans/:instance_audit_plan_id/sqls/:id/analysis", v1.GetAuditPlanSqlAnalysisData)

		v1ProjectViewRouter.GET("/:project_name/sql_versions", v1.GetSqlVersionList)
		v1ProjectViewRouter.GET("/:project_name/sql_versions/:sql_version_id/", v1.GetSqlVersionDetail)
		v1ProjectViewRouter.GET("/:project_name/sql_versions/:sql_version_id/sql_version_stages/:sql_version_stage_id/dependencies", v1.GetDependenciesBetweenStageInstance)
		v1ProjectViewRouter.GET("/:project_name/sql_versions/:sql_version_id/sql_version_stages/:sql_version_stage_id/associate_workflows", v1.GetWorkflowsThatCanBeAssociatedToVersion)
	}

	// project member router
	v2ProjectOpRouter := v2Router.Group("/projects", sqleMiddleware.ProjectMemberOpAllowed())
	{
		// workflow
		v2ProjectOpRouter.POST("/:project_name/workflows", v2.CreateWorkflowV2)
		v2ProjectOpRouter.POST("/:project_name/workflows/:workflow_id/steps/:workflow_step_id/approve", v2.ApproveWorkflowV2)
		v2ProjectOpRouter.POST("/:project_name/workflows/:workflow_id/steps/:workflow_step_id/reject", v2.RejectWorkflowV2)
		v2ProjectOpRouter.POST("/:project_name/workflows/:workflow_id/cancel", v2.CancelWorkflowV2)
		v2ProjectOpRouter.POST("/:project_name/workflows/cancel", v2.BatchCancelWorkflowsV2)
		v2ProjectOpRouter.POST("/:project_name/workflows/complete", v2.BatchCompleteWorkflowsV2)
		v2ProjectOpRouter.POST("/:project_name/workflows/:workflow_id/tasks/:task_id/execute", v2.ExecuteOneTaskOnWorkflowV2)
		v2ProjectOpRouter.POST("/:project_name/workflows/:workflow_id/tasks/execute", v2.ExecuteTasksOnWorkflowV2)
		v2ProjectOpRouter.PUT("/:project_name/workflows/:workflow_id/tasks/:task_id/schedule", v2.UpdateWorkflowScheduleV2)
		v2ProjectOpRouter.PATCH("/:project_name/workflows/:workflow_id/", v2.UpdateWorkflowV2)

		// scanner token auth
		v2ProjectOpRouter.POST("/:project_name/audit_plans/:audit_plan_name/sqls/full", v2.FullSyncAuditPlanSQLs, sqleMiddleware.ScannerVerifier())
		v2ProjectOpRouter.POST("/:project_name/audit_plans/:audit_plan_name/sqls/partial", v2.PartialSyncAuditPlanSQLs, sqleMiddleware.ScannerVerifier())
		v2ProjectOpRouter.POST("/:project_name/audit_plans/:audit_plan_id/sqls/upload", v2.UploadInstanceAuditPlanSQLs, sqleMiddleware.ScannerVerifier())
	}

	v3ProjectOpRouter := v3Router.Group("/projects", sqleMiddleware.ProjectMemberOpAllowed())
	{
		// workflow
		v3ProjectOpRouter.POST("/:project_name/workflows/complete", v3.BatchCompleteWorkflowsV3)

	}

	v2ProjectViewRouter := v2Router.Group("/projects", sqleMiddleware.ProjectMemberViewAllowed())
	{
		v2ProjectViewRouter.GET("/:project_name/workflows/:workflow_id/", v2.GetWorkflowV2)
		v2ProjectViewRouter.GET("/:project_name/workflows/:workflow_id/tasks", v2.GetSummaryOfWorkflowTasksV2)
		// instance
		v2ProjectViewRouter.GET("/:project_name/instances/:instance_name/", v2.GetInstance)
		// audit plan; 智能扫描任务
		v2ProjectViewRouter.GET("/:project_name/audit_plans", v2.GetAuditPlans)
		v2ProjectViewRouter.GET("/:project_name/audit_plans/:audit_plan_name/reports/:audit_plan_report_id/sqls/:number/analysis", v2.GetAuditPlanAnalysisData)
		v2ProjectViewRouter.GET("/:project_name/audit_plans/:audit_plan_name/reports/:audit_plan_report_id/sqls", v2.GetAuditPlanReportSQLs)

		// sql managers
		v2ProjectViewRouter.GET("/:project_name/sql_manages", v2.GetSqlManageList)
	}

	{
		v1Router.GET("/user_tips", v1.GetUserTips)

		// 全局 rule template
		v1Router.GET("/rule_templates", v1.GetRuleTemplates)
		v1Router.GET("/rule_template_tips", v1.GetRuleTemplateTips)
		v1Router.GET("/rule_templates/:rule_template_name/", v1.GetRuleTemplate)

		v1Router.POST("/rule_templates/parse", v1.ParseProjectRuleTemplateFile)
		v1Router.GET("/import_rule_template", v1.GetImportRuleTemplateFile)

		// 全局 workflow
		v1Router.GET("/rule_knowledge/db_types/:db_type/rules/:rule_name/", v1.GetRuleKnowledge)
		v1Router.GET("/rule_knowledge/db_types/:db_type/custom_rules/:rule_name/", v1.GetCustomRuleKnowledge)
		v1Router.GET("/workflows/statistic_of_instances", v1.GetWorkflowStatisticOfInstances)

		//rule
		v1Router.GET("/rules", v1.GetRules)
		v1Router.GET("/custom_rules", v1.GetCustomRules)
		v1Router.GET("/custom_rules/:rule_id", v1.GetCustomRule)
		v1Router.GET("/custom_rules/:db_type/rule_types", v1.GetRuleTypeByDBType)

		// task
		v1Router.GET("/tasks/audits/:task_id/", v1.GetTask)
		v1Router.GET("/tasks/audits/:task_id/sqls", v1.GetTaskSQLs)
		v1Router.PATCH("/tasks/audits/:task_id/sqls/:sql_id/backup_strategy", v1.UpdateSqlBackupStrategy)
		v1Router.PATCH("/tasks/audits/:task_id/backup_strategy", v1.UpdateTaskBackupStrategy)
		v2Router.GET("/tasks/audits/:task_id/sqls", v2.GetTaskSQLs)
		v2Router.GET("/tasks/audits/:task_id/files", v2.GetAuditFileList)
		v2Router.GET("/tasks/audits/:task_id/files/:file_id/", v2.GetAuditFileExecStatistic)
		v1Router.GET("/tasks/audits/:task_id/sql_report", v1.DownloadTaskSQLReportFile)
		v1Router.GET("/tasks/audits/:task_id/sql_file", v1.DownloadTaskSQLFile)
		v1Router.GET("/tasks/audits/:task_id/audit_file", v1.DownloadAuditFile)
		v1Router.GET("/tasks/audits/:task_id/sql_content", v1.GetAuditTaskSQLContent)
		v1Router.PATCH("/tasks/audits/:task_id/sqls/:number", v1.UpdateAuditTaskSQLs)
		v1Router.GET("/tasks/audits/:task_id/sqls/:number/analysis", v1.GetTaskAnalysisData)
		v2Router.GET("/tasks/audits/:task_id/sqls/:number/analysis", v2.GetTaskAnalysisData)
		v1Router.POST("/projects/:project_name/task_groups", v1.CreateAuditTasksGroupV1)
		v1Router.POST("/task_groups/audit", v1.AuditTaskGroupV1)
		v1Router.GET("/tasks/file_order_methods", v1.GetSqlFileOrderMethodV1)

		// dashboard
		v1Router.GET("/dashboard", v1.Dashboard)
		// 全局 sql manage
		v1Router.GET("/dashboard/sql_manages", v1.GetGlobalSqlManageList)
		v1Router.GET("/dashboard/sql_manages/statistics", v1.GetGlobalSqlManageStatistics)
		v1Router.GET("/dashboard/workflows", v1.GetGlobalWorkflowsV1)
		v1Router.GET("/dashboard/workflows/statistics", v1.GetGlobalWorkflowStatistics)

		// configurations
		v1Router.GET("/configurations/drivers", v1.GetDrivers)
		v2Router.GET("/configurations/drivers", v2.GetDrivers)
		v1Router.GET("/configurations/workflows/schedule/default_option", v1.GetScheduledTaskDefaultOptionV1)

		// audit plan
		v1Router.GET("/audit_plan_metas", v1.GetAuditPlanMetas)
		v1Router.GET("/audit_plan_types", v1.GetAuditPlanTypes)

		// sql audit
		v1Router.POST("/sql_audit", v1.DirectAudit)
		v2Router.POST("/sql_audit", v2.DirectAudit)
		v1Router.POST("/audit_files", v1.DirectAuditFiles)
		v2Router.POST("/audit_files", v2.DirectAuditFiles)
		v1Router.GET("/sql_analysis", v1.DirectGetSQLAnalysis)
		// 企业公告
		v1Router.GET("/company_notice", v1.GetCompanyNotice)
		// 系统功能开关
		v1Router.GET("/system/module_status", v1.GetSystemModuleStatus)
		v1Router.GET("/system/module_red_dots", v1.GetSystemModuleRedDots)
	}

	// enterprise customized apis
	err := addCustomApis(v1Router, restApis)
	if err != nil {
		log.Logger().Fatalf("failed to register custom api, %v", err)
		return
	}

	address := fmt.Sprintf(":%v", config.APIServiceOpts.Port)
	log.Logger().Infof("starting http server on %s", address)

	// start http server
	l, err := net.Listen("tcp4", address)
	if err != nil {
		log.Logger().Fatal(err)
		return
	}

	if config.APIServiceOpts.EnableHttps {
		// Usually, it is easier to create an tls server using echo#StartTLS;
		// but I need create a graceful listener.
		if config.APIServiceOpts.CertFilePath == "" || config.APIServiceOpts.KeyFilePath == "" {
			log.Logger().Fatal("invalid tls configuration")
			return
		}
		tlsConfig := new(tls.Config)
		tlsConfig.Certificates = make([]tls.Certificate, 1)
		tlsConfig.Certificates[0], err = tls.LoadX509KeyPair(config.APIServiceOpts.CertFilePath, config.APIServiceOpts.KeyFilePath)
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
