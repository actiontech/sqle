package v1

import (
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

type CreateSqlVersionReqV1 struct {
	Version         string                  `json:"version" form:"version" valid:"required" example:"2.23"`
	Desc            string                  `json:"desc" form:"desc"`
	SqlVersionStage []CreateSqlVersionStage `json:"create_sql_version_stage" valid:"dive,required"`
}

type CreateSqlVersionStage struct {
	Name                    string                    `json:"name" form:"name" valid:"required" example:"生产"`
	StageSequence           int                       `json:"stage_sequence" form:"stage_sequence" valid:"required"`
	CreateStagesInstanceDep []CreateStagesInstanceDep `json:"create_stages_instance_dep"`
}

type CreateStagesInstanceDep struct {
	StageInstanceID     string `json:"stage_instance_id"`
	NextStageInstanceID string `json:"next_stage_instance_id"`
}

// @Summary 创建SQL版本记录
// @Description create sql version
// @Id createSqlVersionV1
// @Tags sql_version
// @Security ApiKeyAuth
// @Accept json
// @Param project_name path string true "project name"
// @Param sql_version body v1.CreateSqlVersionReqV1 true "create sql version request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/sql_versions [post]
func CreateSqlVersion(c echo.Context) error {

	/**
		1、save version
		2、save stage
		3、遍历stage，stage的id作为dependency的stage id
		4、实现方法getNextSatgeBySequence，获取next stage id（SqlVersionStage表StageSequence + 1的id作为NextStageID）
		5、save dependenc
	**/

	return createSqlVersion(c)
}

type GetSqlVersionListReqV1 struct {
	FilterByCreatedAtFrom *string `json:"filter_by_created_at_from,omitempty" query:"filter_by_created_at_from"`
	FilterByCreatedAtTo   *string `json:"filter_by_created_at_to,omitempty" query:"filter_by_created_at_to"`
	FilterByLockTimeFrom  *string `json:"filter_by_lock_time_from,omitempty" query:"filter_by_lock_time_from"`
	FilterByLockTimeTo    *string `json:"filter_by_lock_time_to,omitempty" query:"filter_by_lock_time_to"`
	FilterByVersionStatus *string `json:"filter_by_version_status,omitempty" query:"filter_by_version_status"`
	FuzzySearch           *string `json:"fuzzy_search,omitempty" query:"fuzzy_search"`

	PageIndex uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetSqlVersionListResV1 struct {
	controller.BaseRes
	Data      []*SqlVersionResV1 `json:"data"`
	TotalNums uint64             `json:"total_nums"`
}

type SqlVersionResV1 struct {
	VersionID             uint       `json:"version_id"`
	Version               string     `json:"version"`
	Desc                  string     `json:"desc"`
	Status                string     `json:"status" enums:"is_being_released,locked"`
	LockTime              *time.Time `json:"lock_time"`
	CreatedAt             *time.Time `json:"created_at"`
	HasAssociatedWorkflow bool       `json:"has_associated_workflow"`
}

// @Summary 获取SQL版本列表
// @Description sql version list
// @Id getSqlVersionListV1
// @Tags sql_version
// @Security ApiKeyAuth
// @Accept json
// @Param project_name path string true "project name"
// @Param filter_by_created_at_from query string false "filter by created at from"
// @Param filter_by_created_at_to query string false "filter by created at to"
// @Param filter_by_lock_time_from query string false "filter by lock time from"
// @Param filter_by_lock_time_to query string false "filter by lock time to"
// @Param filter_by_version_status query string false "filter by version status"
// @Param fuzzy_search query string false "fuzzy search"
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Success 200 {object} v1.GetSqlVersionListResV1
// @router /v1/projects/{project_name}/sql_versions [get]
func GetSqlVersionList(c echo.Context) error {

	// 获取列表同时返回版本是否有关联工单
	return getSqlVersionList(c)
}

type GetSqlVersionDetailResV1 struct {
	controller.BaseRes
	Data []*SqlVersionDetailResV1 `json:"data"`
}

type SqlVersionDetailResV1 struct {
	StageID         uint                          `json:"stage_id"`
	StageName       string                        `json:"stage_name"`
	StageSequence   int                           `json:"stage_sequence"`
	WorkflowDetails []*WorkflowDetailWithInstance `json:"workflow_details"`
}

type WorkflowDetailWithInstance struct {
	Name              string              `json:"workflow_name"`
	WorkflowId        string              `json:"workflow_id"`
	Desc              string              `json:"desc,omitempty"`
	Status            string              `json:"status" enums:"wait_for_audit,wait_for_execution,rejected,canceled,exec_failed,executing,finished"`
	WorkflowInstances []*WorkflowInstance `json:"workflow_instances"`
}
type WorkflowInstance struct {
	InstanceID     string `json:"instances_id"`
	InstanceName   string `json:"instances_name"`
	InstanceSchema string `json:"instance_schema"`
}

// @Summary 获取SQL版本详情
// @Description get sql version detail
// @Id getSqlVersionDetailV1
// @Tags sql_version
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param sql_version_id path string true "sql version id"
// @Success 200 {object} v1.GetSqlVersionDetailResV1
// @router /v1/projects/{project_name}/sql_versions/{sql_version_id}/ [get]
func GetSqlVersionDetail(c echo.Context) error {

	/**
		1、getStageBySqlVersionID
		2、获取工单信息dms.GetWorkflowDetailByWorkflowId(projectUid, workflowID, s.GetWorkflowDetailWithoutInstancesByWorkflowID)
		3、遍历workflow获取WorkflowDetail遍历workflow.record.instancerecords获取WorkflowInstance
	**/
	return getSqlVersionDetail(c)
}

type UpdateSqlVersionReqV1 struct {
	Version         *string                  `json:"version" form:"version" example:"2.23"`
	Desc            *string                  `json:"desc" form:"desc"`
	SqlVersionStage *[]UpdateSqlVersionStage `json:"update_sql_version_stage"`
}

type UpdateSqlVersionStage struct {
	StageID                 *uint                      `json:"stage_id" form:"stage_id"`
	Name                    *string                    `json:"name" form:"name" valid:"required" example:"生产"`
	StageSequence           *int                       `json:"stage_sequence" form:"stage_sequence" valid:"required"`
	CreateStagesInstanceDep *[]UpdateStagesInstanceDep `json:"update_stages_instance_dep"`
}

type UpdateStagesInstanceDep struct {
	StageInstanceID     string `json:"stage_instance_id"`
	NextStageInstanceID string `json:"next_stage_instance_id"`
}

// @Summary 更新SQL版本信息
// @Description update sql version
// @Id updateSqlVersionV1
// @Tags sql_version
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param sql_version_id path string true "sql version id"
// @Param sql_version body v1.UpdateSqlVersionReqV1 false "update sql version request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/sql_versions/{sql_version_id}/ [patch]
func UpdateSqlVersion(c echo.Context) error {
	/**
		1、getStageBySqlVersionID
		2、save sql version
		3、如果要更新stage，传入完整的stage，覆盖式更新stage及dependency
	  **/
	return updateSqlVersion(c)
}

type LockSqlVersionReqV1 struct {
	IsLocked bool `json:"is_locked"`
}

// @Summary 锁定SQL版本
// @Description lock sql version
// @Id lockSqlVersionV1
// @Tags sql_version
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param sql_version_id path string true "sql version id"
// @Param sql_version body v1.LockSqlVersionReqV1 true "lock sql version request"
// @Success 200 {object}  controller.BaseRes
// @router /v1/projects/{project_name}/sql_versions/{sql_version_id}/lock [post]
func LockSqlVersion(c echo.Context) error {

	return lockSqlVersion(c)
}

// @Summary 删除SQL版本
// @Description delete sql version
// @Id deleteSqlVersionV1
// @Tags sql_version
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param sql_version_id path string true "sql version id"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/sql_versions/{sql_version_id}/ [delete]
func DeleteSqlVersion(c echo.Context) error {
	// delete sql version
	// delete stage
	// delete dependency
	// delete workflow release stage
	return deleteSqlVersion(c)
}

type GetDepBetweenStageInstanceResV1 struct {
	controller.BaseRes
	Data []*DepBetweenStageInstance `json:"data"`
}

type DepBetweenStageInstance struct {
	StageInstanceID       string `json:"stage_instance_id"`
	StageInstanceName     string `json:"stage_instance_name"`
	NextStageInstanceID   string `json:"next_stage_instance_id"`
	NextStageInstanceName string `json:"next_stage_instance_name"`
}

// @Summary 获取当前阶段与下一阶段实例的依赖信息
// @Description get dependencies between stage instance
// @Tags sql_version
// @Id getDependenciesBetweenStageInstanceV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param sql_version_id path string true "sql version id"
// @Param sql_version_stage_id path string true "sql version stage id"
// @Success 200 {object} v1.GetDepBetweenStageInstanceResV1
// @router /v1/projects/{project_name}/sql_versions/{sql_version_id}/sql_version_stages/{sql_version_stage_id}/dependencies [get]
func GetDependenciesBetweenStageInstance(c echo.Context) error {
	/**
		select * from SqlVersionStagesDependency where SqlVersionStageID = sql_version_stage_id
	 **/
	return getDependenciesBetweenStageInstance(c)
}

type BatchReleaseWorkflowReqV1 struct {
	ReleaseWorkflows []ReleaseWorkflows `json:"release_workflows" form:"release_workflows" valid:"dive,required"`
}

type ReleaseWorkflows struct {
	WorkFlowID             string                  `json:"workflow_id" form:"workflow_id" valid:"required"`
	TargetReleaseInstances []TargetReleaseInstance `json:"target_release_instances" valid:"dive,required"`
}

type TargetReleaseInstance struct {
	InstanceID           string `json:"instance_id" form:"instance_id" valid:"required"`
	InstanceSchema       string `json:"instance_schema" form:"instance_schema"`
	TargetInstanceID     string `json:"target_instance_id" form:"target_instance_id" valid:"required"`
	TargetInstanceSchema string `json:"target_instance_schema" form:"target_instance_schema"`
}

// @Summary 批量发布工单（在版本的下一阶段创建工单）
// @Description batch release workflow
// @Id batchReleaseWorkflowsV1
// @Tags sql_version
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param sql_version_id path string true "sql version id"
// @Param data body v1.BatchReleaseWorkflowReqV1 true "batch release workflow request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/sql_versions/{sql_version_id}/batch_release_workflows [post]
func BatchReleaseWorkflows(c echo.Context) error {
	/**
		1、遍历获取工单信息，dms.GetWorkflowDetailByWorkflowId(projectUid, workflowID, s.GetWorkflowDetailWithoutInstancesByWorkflowID)
		2、遍历workflow.task[]，获取task信息
		3、找到task的instance id，对应到target instance id
		4、dms.GetInstancesById(c.Request().Context(), req.InstanceId)
		5、参考func CreateAndAuditTask(c echo.Context) error 创建审核，如果原task的sql source是sql_file，
			参考func GetWorkflowTaskAuditFile(c echo.Context) error获取原始文件
		6、提交工单工单（忽略工单审批流程模板的审核等级限制）参考func CreateWorkflowV2(c echo.Context) error
	**/
	return batchReleaseWorkflows(c)
}

type BatchExecuteTasksOnWorkflowReqV1 struct {
	WorkflowIDs []string `json:"workflow_ids" valid:"required"`
}

// @Summary 工单批量上线
// @Description batch execute tasks on workflow
// @Tags sql_version
// @Id batchExecuteTasksOnWorkflowV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param sql_version_id path string true "sql version id"
// @Param data body v1.BatchExecuteTasksOnWorkflowReqV1 true "batch execute tasks on workflow request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/sql_versions/{sql_version_id}/batch_execute_workflows [post]
func BatchExecuteTasksOnWorkflow(c echo.Context) error {
	/**
		1、遍历workflow id，获取workflow信息
		2、参考func ExecuteTasksOnWorkflowV2(c echo.Context) error 执行上线
	**/
	return batchExecuteTasksOnWorkflow(c)
}

type RetryExecWorkflowReqV1 struct {
	WorkflowID string `json:"workflow_ids" valid:"required"`
	TaskIds    []uint `json:"task_ids" form:"task_ids" valid:"required"`
}

// @Summary 工单重试（上线失败修改sql重试上线）
// @Description reject exec failed workflow
// @Tags sql_version
// @Id retryExecWorkflowV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param sql_version_id path string true "sql version id"
// @Param data body v1.RetryExecWorkflowReqV1 true "retry execute workflow request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/sql_versions/{sql_version_id}/retry_workflow [post]
func RetryExecWorkflow(c echo.Context) error {
	/**
		暂不考虑重复上线问题，用户在修改sql时自行回滚或删除上线成功的sql。
		发起人有权限重试并提交修改的sql，
		1、参考驳回后修改sql重新审核提交工单逻辑
		2、不过滤上线失败的工单
		3、调整工单上线逻辑，上线失败current_workflow_step_id不要置为0（考虑影响）
		4、workflow record history调整，可以在工单详情中查看到上线失败的history
	**/
	return retryExecWorkflow(c)
}

type BatchAssociateWorkflowsWithVersionReqV1 struct {
	StageAndWorkflows []StageAndWorkflows `json:"stage_Workflows" form:"stage_Workflows" valid:"dive,required"`
}
type StageAndWorkflows struct {
	SqlVersionStageID *string `json:"sql_version_stage_id"`
	WorkflowID        *string `json:"workflow_id" valid:"required"`
}

// @Summary 批量关联工单到版本
// @Description batch associate workflows with version
// @Tags sql_version
// @Id batchAssociateWorkflowsWithVersionV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param sql_version_id path string true "sql version id"
// @Param data body v1.BatchAssociateWorkflowsWithVersionReqV1 true "batch associate workflows with version request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/sql_versions/{sql_version_id}/associate_workflows [post]
func BatchAssociateWorkflowsWithVersion(c echo.Context) error {

	// 与第一个阶段进行关联
	return batchAssociateWorkflowsWithVersion(c)
}

type GetWorkflowsThatCanBeAssociatedToVersionResV1 struct {
	controller.BaseRes
	Data []*AssociateWorkflows `json:"data"`
}
type AssociateWorkflows struct {
	WorkflowID   string `json:"workflow_id"`
	WorkflowName string `json:"workflow_name"`
	WorkflowDesc string `json:"desc"`
}

// @Summary 获取可与版本关联的工单
// @Description get workflows that can be associated to version
// @Tags sql_version
// @Id GetWorkflowsThatCanBeAssociatedToVersionV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param sql_version_id path string true "sql version id"
// @Param sql_version_stage_id path string true "sql version stage id"
// @Success 200 {object} v1.GetWorkflowsThatCanBeAssociatedToVersionResV1
// @router /v1/projects/{project_name}/sql_versions/{sql_version_id}/sql_version_stages/{sql_version_stage_id}/associate_workflows [get]
func GetWorkflowsThatCanBeAssociatedToVersion(c echo.Context) error {

	return getWorkflowsThatCanBeAssociatedToVersion(c)
}
