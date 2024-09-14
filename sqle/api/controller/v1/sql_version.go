package v1

import (
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

type SqlVersionReqV1 struct {
	VersionNumber   string            `json:"version_number" form:"version_number" valid:"required" example:"2.23"`
	Desc            string            `json:"desc" form:"desc"`
	SqlVersionStage []SqlVersionStage `json:"sql_version_stage" valid:"dive,required"`
}

type SqlVersionStage struct {
	Name             string                       `json:"name" form:"name" valid:"required" example:"生产"`
	StageSequence    int                          `json:"stage_sequence" form:"stage_sequence" valid:"required"`
	StagesDependency []SqlVersionStagesDependency `json:"stages_dependency" valid:"dive,required"`
}

type SqlVersionStagesDependency struct {
	StageInstanceID     uint `json:"source_instance_id"`
	NextStageInstanceID uint `json:"target_instance_id"`
}

// @Summary 创建sql版本记录
// @Description create sql version
// @Id createSqlVersionV1
// @Tags sql_version
// @Security ApiKeyAuth
// @Accept json
// @Param project_name path string true "project name"
// @Param sql_version body v1.SqlVersionReqV1 true "create sql version request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/sql_version [post]
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
	VersionNumber         string     `json:"version_number"`
	Desc                  string     `json:"desc"`
	Status                string     `json:"status"`
	LockTime              *time.Time `json:"lock_time"`
	CreatedAt             *time.Time `json:"created_at"`
	HasAssociatedWorkflow bool       `json:"has_associated_workflow"`
}

// @Summary 获取sql版本列表
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
// @Param fuzzy_search query string false "fuzzy search"
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page"
// @Success 200 {object} v1.GetSqlVersionListResV1
// @router /v1/projects/{project_name}/sql_version [get]
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

// @Summary 获取sql版本详情
// @Description get sql version detail
// @Id getSqlVersionDetailV1
// @Tags sql_version
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param sql_version_id path string true "sql version id"
// @Success 200 {object} v1.GetSqlVersionDetailResV1
// @router /v1/projects/{project_name}/sql_version/{sql_version_id}/ [get]
func GetSqlVersionDetail(c echo.Context) error {

	/**
		1、getStageBySqlVersionID
		2、获取工单信息dms.GetWorkflowDetailByWorkflowId(projectUid, workflowID, s.GetWorkflowDetailWithoutInstancesByWorkflowID)
		3、遍历workflow获取WorkflowDetail遍历workflow.record.instancerecords获取WorkflowInstance
	**/
	return getSqlVersionDetail(c)
}

// @Summary 更新sql版本信息
// @Description update sql version
// @Id updateSqlVersionV1
// @Tags sql_version
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param sql_version_id path string true "sql version id"
// @Param sql_version body v1.SqlVersionReqV1 true "update sql version request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/sql_version/{sql_version_id}/ [put]
func UpdateSqlVersion(c echo.Context) error {
	/**
		没有工单：
		1、getStageBySqlVersionID
		2、save sql version
		3、遍历db stage，嵌套遍历req satge，根据stage顺号对比更新、删除、添加stage
		有工单：
		save sql version
	  **/
	return updateSqlVersion(c)
}

type LockSqlVersionReqV1 struct {
	IsLocked bool `json:"is_locked"`
}

// @Summary 锁定sql版本
// @Description lock sql version
// @Id lockSqlVersionV1
// @Tags sql_version
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param sql_version_id path string true "sql version id"
// @Param sql_version body v1.LockSqlVersionReqV1 true "lock sql version request"
// @Success 200 {object}  controller.BaseRes
// @router /v1/projects/{project_name}/sql_version/{sql_version_id}/lock [patch]
func LockSqlVersion(c echo.Context) error {

	return lockSqlVersion(c)
}

// @Summary 删除sql版本
// @Description delete sql version
// @Id deleteSqlVersionV1
// @Tags sql_version
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param sql_version_id path string true "sql version id"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/sql_version/{sql_version_id}/ [delete]
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
	StageInstanceID     string `json:"stage_instance_id"`
	NextStageInstanceID string `json:"next_stage_instance_id"`
}

// @Summary 获取当前阶段与下一阶段实例的依赖信息
// @Description get dependencies between stage instance
// @Tags sql_version
// @Id getInstanceTipListV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param sql_version_id path string true "sql version id"
// @Param sql_version_stage_id path string true "sql version stage id"
// @Success 200 {object} v1.GetInstanceTipsResV1
// @router /v1/projects/{project_name}/sql_version/{sql_version_id}/sql_version_stage/{sql_version_stage_id}/dependencies [get]
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

// @Summary 批量发布工单
// @Description batch release workflow
// @Id batchReleaseWorkflowV1
// @Tags sql_version
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param sql_version_id path string true "sql version id"
// @Param data body v1.BatchReleaseWorkflowReqV1 true "batch release workflow request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/sql_version/{sql_version_id}/workflow/batch_release [post]
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
// @Id batchExecuteTasksOnWorkflow
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param sql_version_id path string true "sql version id"
// @Param data body v1.BatchExecuteTasksOnWorkflowReqV1 true "batch execute tasks on workflow request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/sql_version/{sql_version_id}/workflow/batch_execute [post]
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

// @Summary 工单重试
// @Description reject exec failed workflow
// @Tags sql_version
// @Id rejectExecFailedWorkflow
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param sql_version_id path string true "sql version id"
// @Param data body v1.RetryExecWorkflowReqV1 true "reject exec failed workflow request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/sql_version/{sql_version_id}/workflow/execute_retry [post]
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

type WorkflowSqlVersionReqV1 struct {
	SqlVersionID *string `json:"sql_version_id"`
}

// @Summary 工单与版本建立关联
// @Description workflow sql version
// @Tags sql_version
// @Id workflowSqlVersion
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param workflow_id path string true "workflow id"
// @Param data body v1.WorkflowSqlVersionReqV1 true "workflow sql version request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/sql_version/workflow/{workflow_id}/ [post]
func WorkflowSqlVersion(c echo.Context) error {

	// 与第一个阶段进行关联
	return workflowSqlVersion(c)
}
