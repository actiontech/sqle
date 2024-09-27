//go:build enterprise
// +build enterprise

package v1

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"github.com/actiontech/sqle/sqle/api/controller"
	dms "github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/notification"
	"github.com/actiontech/sqle/sqle/pkg/im"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/actiontech/sqle/sqle/server/sqlversion"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/labstack/echo/v4"
)

func createSqlVersion(c echo.Context) error {
	// TODO 权限校验
	req := new(CreateSqlVersionReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	versionStages := make([]*model.SqlVersionStage, 0, len(req.SqlVersionStage))
	for _, stage := range req.SqlVersionStage {
		stageDeps := make([]*model.SqlVersionStagesDependency, 0)
		for _, dep := range stage.CreateStagesInstanceDep {
			stageInstID, err := strconv.ParseUint(dep.StageInstanceID, 10, 64)
			if err != nil {
				return err
			}
			var nextStageInstID uint64
			if dep.NextStageInstanceID != "" {
				nextStageInstID, err = strconv.ParseUint(dep.NextStageInstanceID, 10, 64)
				if err != nil {
					return err
				}
			}
			stageDeps = append(stageDeps, &model.SqlVersionStagesDependency{
				StageInstanceID:     stageInstID,
				NextStageInstanceID: nextStageInstID,
			})
		}
		versionStages = append(versionStages, &model.SqlVersionStage{
			Name:                       stage.Name,
			StageSequence:              stage.StageSequence,
			SqlVersionStagesDependency: stageDeps,
		})
	}
	sqlVersion := &model.SqlVersion{
		Version:         req.Version,
		Description:     req.Desc,
		Status:          model.SqlVersionStatusReleased,
		ProjectId:       model.ProjectUID(projectUid),
		SqlVersionStage: versionStages,
	}
	s := model.GetStorage()
	err = s.BatchSaveSqlVersion(sqlVersion)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return controller.JSONBaseErrorReq(c, nil)
}

func getSqlVersionList(c echo.Context) error {
	// TODO 权限校验
	req := new(GetSqlVersionListReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	limit, offset := controller.GetLimitAndOffset(req.PageIndex, req.PageSize)

	userId := controller.GetUserID(c)

	up, err := dms.NewUserPermission(userId, projectUid)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	data := map[string]interface{}{
		"filter_by_created_at_from": req.FilterByCreatedAtFrom,
		"filter_by_created_at_to":   req.FilterByCreatedAtTo,
		"filter_by_lock_time_from":  req.FilterByLockTimeFrom,
		"filter_by_lock_time_to":    req.FilterByLockTimeTo,
		"filter_by_version_status":  req.FilterByVersionStatus,
		"fuzzy_search":              req.FuzzySearch,
		"filter_by_project_id":      projectUid,
		"current_user_id":           userId,
		"current_user_is_admin":     up.IsAdmin(),
		"limit":                     limit,
		"offset":                    offset,
	}
	if !up.IsAdmin() {

	}
	s := model.GetStorage()

	sqlVersions, count, err := s.GetSqlVersionByReq(data)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	resData := make([]*SqlVersionResV1, len(sqlVersions))
	for i, v := range sqlVersions {

		resData[i] = &SqlVersionResV1{
			VersionID: v.Id,
			Version:   v.Version.String,
			Desc:      v.Desc.String,
			Status:    v.Status.String,
			LockTime:  v.LockTime,
			CreatedAt: v.CreatedAt,
		}
	}
	return c.JSON(http.StatusOK, &GetSqlVersionListResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      resData,
		TotalNums: count,
	})
}

func getSqlVersionDetail(c echo.Context) error {
	// TODO 权限校验
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	sqlVersionId := c.Param("sql_version_id")
	s := model.GetStorage()
	version, exist, err := s.GetSqlVersionDetailByVersionId(sqlVersionId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewDataNotExistErr("sql version not found"))
	}
	stages := make([]SqlVersionStageDetail, 0, len(version.SqlVersionStage))
	for _, stage := range version.SqlVersionStage {
		stageInstances := make([]VersionStageInstance, 0, len(stage.SqlVersionStagesDependency))
		for _, dep := range stage.SqlVersionStagesDependency {
			instanceIdNameMap, err := dms.GetInstanceIdNameMapByIds(c.Request().Context(), []uint64{dep.StageInstanceID})
			if err != nil {
				return controller.JSONBaseErrorReq(c, err)
			}
			stageInstances = append(stageInstances, VersionStageInstance{
				InstanceID:   strconv.FormatUint(dep.StageInstanceID, 10),
				InstanceName: instanceIdNameMap[dep.StageInstanceID],
			})
		}
		workflows := make([]WorkflowDetailWithInstance, 0, len(stage.WorkflowVersionStage))
		for _, workflowStage := range stage.WorkflowVersionStage {
			workflow, err := dms.GetWorkflowDetailByWorkflowId(projectUid, workflowStage.WorkflowID, s.GetWorkflowDetailWithoutInstancesByWorkflowID)
			if err != nil {
				return controller.JSONBaseErrorReq(c, err)
			}

			workflowInstances := make([]VersionStageInstance, 0, len(workflow.Record.InstanceRecords))
			for _, workflowInstance := range workflow.Record.InstanceRecords {
				workflowInstances = append(workflowInstances, VersionStageInstance{
					InstanceID:     workflowInstance.Instance.GetIDStr(),
					InstanceName:   workflowInstance.Instance.Name,
					InstanceSchema: workflowInstance.Task.Schema,
				})
			}

			workflows = append(workflows, WorkflowDetailWithInstance{
				Name:                  workflow.Subject,
				WorkflowId:            workflow.WorkflowId,
				Desc:                  workflow.Desc,
				WorkflowSequence:      workflowStage.WorkflowSequence,
				Status:                workflow.Record.Status,
				WorkflowReleaseStatus: workflowStage.WorkflowReleaseStatus,
				WorkflowExecTime:      workflowStage.WorkflowExecTime,
				WorkflowInstances:     &workflowInstances,
			})
		}

		versionStage := SqlVersionStageDetail{
			StageID:         stage.ID,
			StageName:       stage.Name,
			StageSequence:   stage.StageSequence,
			StageInstances:  &stageInstances,
			WorkflowDetails: &workflows,
		}
		stages = append(stages, versionStage)
	}
	resData := &SqlVersionDetailResV1{
		SqlVersionID:          version.ID,
		Version:               version.Version,
		Status:                version.Status,
		SqlVersionDesc:        version.Description,
		SqlVersionStageDetail: &stages,
	}
	return c.JSON(http.StatusOK, &GetSqlVersionDetailResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    resData,
	})
}

func updateSqlVersion(c echo.Context) error {
	return nil
}

func lockSqlVersion(c echo.Context) error {
	return nil
}

func deleteSqlVersion(c echo.Context) error {
	return nil
}

func getDependenciesBetweenStageInstance(c echo.Context) error {
	// TODO 权限校验

	// projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"))
	// if err != nil {
	// 	return controller.JSONBaseErrorReq(c, err)
	// }
	stageId := c.Param("sql_version_stage_id")
	s := model.GetStorage()
	dependencies, err := s.GetStageDependenciesByStageId(stageId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	resData := make([]*DepBetweenStageInstance, 0, len(dependencies))
	for _, dep := range dependencies {
		instanceIdNameMap, err := dms.GetInstanceIdNameMapByIds(c.Request().Context(), []uint64{dep.StageInstanceID, dep.NextStageInstanceID})
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		depInst := &DepBetweenStageInstance{
			StageInstanceID:   strconv.FormatUint(dep.StageInstanceID, 10),
			StageInstanceName: instanceIdNameMap[dep.StageInstanceID],
		}
		if dep.NextStageInstanceID != 0 {
			depInst.NextStageInstanceID = strconv.FormatUint(dep.NextStageInstanceID, 10)
			depInst.NextStageInstanceName = instanceIdNameMap[dep.NextStageInstanceID]
		}
		resData = append(resData, depInst)
	}
	return c.JSON(http.StatusOK, &GetDepBetweenStageInstanceResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    resData,
	})
}

func batchReleaseWorkflows(c echo.Context) error {
	req := new(BatchReleaseWorkflowReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	sqlVersionId, err := strconv.ParseInt(c.Param("sql_version_id"), 10, 64)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	s := model.GetStorage()
	for _, releaseWorkflow := range req.ReleaseWorkflows {

		workflow, err := dms.GetWorkflowDetailByWorkflowId(projectUid, releaseWorkflow.WorkFlowID, s.GetWorkflowDetailWithoutInstancesByWorkflowID)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		tasks := make([]*model.Task, 0)
		for _, instRecords := range workflow.Record.InstanceRecords {
			// 根据阶段间数据源对应关系，获发布到取下一阶段数据源
			targetInst, targetSchema, err := getReleaseTargetInstanceByRelation(c, projectUid, strconv.FormatUint(instRecords.InstanceId, 10), instRecords.Task.Schema, releaseWorkflow.TargetReleaseInstances)
			if err != nil {
				return controller.JSONBaseErrorReq(c, err)
			}
			// 根据当前工单的task构建发布工单的task信息
			task := buildNewTaskByOriginalTask(uint64(user.ID), targetSchema, targetInst, instRecords.Task)
			tasks = append(tasks, task)
		}
		taskIds, err := batchCreateTask(tasks)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		// 获取发布下一阶段的id和根据阶段名生成的工单名称
		nextSubject, nextSatgeId, err := genNextStageForWorkflow(s, uint(sqlVersionId), releaseWorkflow.WorkFlowID)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		stageWorkflow, err := s.GetStageWorkflowByWorkflowId(uint(sqlVersionId), releaseWorkflow.WorkFlowID)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		err = createWorkFlow(c, s, nextSubject, workflow.Desc, projectUid, uint(sqlVersionId), nextSatgeId, stageWorkflow.WorkflowSequence, user, taskIds)
		if err != nil {
			return controller.JSONBaseErrorReq(c, errors.New(errors.SQLVersionNotAllTasksExecutedSuccess, fmt.Errorf("workflow %s release fail and stop release", workflow.Subject)))
		}
		err = s.UpdateWorkflowReleaseStatus(releaseWorkflow.WorkFlowID, model.WorkflowReleaseStatusHaveBeenReleased, uint(sqlVersionId))
	}
	return controller.JSONBaseErrorReq(c, nil)
}

func getReleaseTargetInstanceByRelation(c echo.Context, projectUid, originalInstanceId, originalSchema string, instanceRelations []TargetReleaseInstance) (targetInstance *model.Instance, targetSchema string, err error) {
	targetInstId, targetSchema, err := getInstanceIdAndSchemaByRelation(originalInstanceId, originalSchema, instanceRelations)
	if err != nil {
		return nil, "", err
	}
	// TODO 改为批量获取数据源
	targetInstance, exist, err := dms.GetInstancesById(c.Request().Context(), targetInstId)
	if err != nil {
		return nil, "", err
	}
	if !exist {
		return nil, "", ErrInstanceNoAccess
	}
	can, err := CheckCurrentUserCanAccessInstances(c.Request().Context(), projectUid, controller.GetUserID(c), []*model.Instance{targetInstance})
	if err != nil {
		return nil, "", err
	}
	if !can {
		return nil, "", ErrInstanceNoAccess
	}
	return targetInstance, targetSchema, nil
}

func buildNewTaskByOriginalTask(userId uint64, schema string, instance *model.Instance, oldTask *model.Task) *model.Task {
	task := &model.Task{
		Schema:          schema,
		InstanceId:      instance.ID,
		Instance:        instance,
		CreateUserId:    userId,
		ExecuteSQLs:     []*model.ExecuteSQL{},
		SQLSource:       oldTask.SQLSource,
		DBType:          instance.DbType,
		ExecMode:        oldTask.ExecMode,
		FileOrderMethod: oldTask.FileOrderMethod,
	}
	for _, execSql := range oldTask.ExecuteSQLs {
		task.ExecuteSQLs = append(task.ExecuteSQLs, &model.ExecuteSQL{
			BaseSQL: model.BaseSQL{
				Number:      execSql.Number,
				Content:     execSql.Content,
				SourceFile:  execSql.SourceFile,
				StartLine:   execSql.StartLine,
				SQLType:     execSql.SQLType,
				ExecBatchId: execSql.ExecBatchId,
			},
		})
	}

	return task
}

func batchCreateTask(tasks []*model.Task) ([]uint, error) {
	s := model.GetStorage()
	taskIds := make([]uint, 0, len(tasks))
	for _, task := range tasks {
		// if task instance is not nil, gorm will update instance when save task.
		tmpInst := *task.Instance
		task.Instance = nil

		err := convertSQLSourceEncodingFromTask(task)
		if err != nil {
			return nil, err
		}
		taskGroup := model.TaskGroup{Tasks: []*model.Task{task}}
		err = s.Save(&taskGroup)
		if err != nil {
			return nil, err
		}
		task.Instance = &tmpInst
		task, err = server.GetSqled().AddTaskWaitResult(fmt.Sprintf("%d", task.ID), server.ActionTypeAudit)
		if err != nil {
			return nil, err
		}
		taskIds = append(taskIds, task.ID)
	}
	return taskIds, nil
}

func getInstanceIdAndSchemaByRelation(originalInstanceId, originalSchema string, instanceRelations []TargetReleaseInstance) (instance string, schema string, err error) {
	for _, instRelation := range instanceRelations {
		if originalInstanceId == instRelation.InstanceID {
			if originalSchema == "" {
				return instRelation.TargetInstanceID, "", nil
			} else if originalSchema == instRelation.InstanceSchema {
				return instRelation.TargetInstanceID, instRelation.TargetInstanceSchema, nil
			}
		}
	}
	return "", "", errors.New(errors.DataNotExist, fmt.Errorf("release target data source not found"))
}

func genNextStageForWorkflow(s *model.Storage, sqlVersionID uint, workflowId string) (nextSubject string, nextSatgeId uint, err error) {
	firstStageWorkflow, err := s.GetWorkflowOfFirstStage(sqlVersionID, workflowId)
	if err != nil {
		return "", 0, err
	}
	nextStage, err := s.GetWorkflowOfNextStage(sqlVersionID, workflowId)
	if err != nil {
		return "", 0, err
	}
	nextSatgeId = nextStage.ID
	nextSubject = firstStageWorkflow.Subject + "_" + nextStage.Name
	return nextSubject, nextSatgeId, nil
}

func createWorkFlow(c echo.Context, s *model.Storage, subject, desc, projectUid string, sqlVersionId, nextSatgeId uint, workflowStageSequence int, user *model.User, taskIds []uint) error {
	// dms-todo: 与 dms 生成uid保持一致
	// TODO 抽取与创建工单接口共同的方法
	workflowId, err := utils.GenUid()
	if err != nil {
		return err
	}

	tasks, foundAllTasks, err := s.GetTasksByIds(taskIds)
	if err != nil {
		return err
	}
	if !foundAllTasks {
		return errors.NewTaskNoExistOrNoAccessErr()
	}

	instanceIds := make([]uint64, 0, len(tasks))
	for _, task := range tasks {
		instanceIds = append(instanceIds, task.InstanceId)
	}

	instances, err := dms.GetInstancesInProjectByIds(c.Request().Context(), projectUid, instanceIds)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	instanceMap := map[uint64]*model.Instance{}
	for _, instance := range instances {
		instanceMap[instance.ID] = instance
	}

	for _, task := range tasks {
		if instance, ok := instanceMap[task.InstanceId]; ok {
			task.Instance = instance
		}
	}

	workflowTemplate, exist, err := s.GetWorkflowTemplateByProjectId(model.ProjectUID(projectUid))
	if err != nil {
		return err
	}
	if !exist {
		return errors.New(errors.DataNotExist, fmt.Errorf("the task instance is not bound workflow template"))
	}

	stepTemplates, err := s.GetWorkflowStepsByTemplateId(workflowTemplate.ID)
	if err != nil {
		return err
	}

	memberWithPermissions, _, err := dmsobject.ListMembersInProject(c.Request().Context(), controller.GetDMSServerAddress(), dmsV1.ListMembersForInternalReq{
		ProjectUid: projectUid,
		PageSize:   999,
		PageIndex:  1,
	})
	if err != nil {
		return err
	}

	err = s.CreateWorkflowV2(subject, workflowId, desc, user, tasks, stepTemplates, model.ProjectUID(projectUid), &sqlVersionId, &nextSatgeId, &workflowStageSequence, func(tasks []*model.Task) (auditWorkflowUsers, canExecUser [][]*model.User) {
		auditWorkflowUsers = make([][]*model.User, len(tasks))
		executorWorkflowUsers := make([][]*model.User, len(tasks))
		for i, task := range tasks {
			auditWorkflowUsers[i], err = GetCanOpInstanceUsers(memberWithPermissions, task.Instance, []dmsV1.OpPermissionType{dmsV1.OpPermissionTypeAuditWorkflow})
			if err != nil {
				return
			}
			executorWorkflowUsers[i], err = GetCanOpInstanceUsers(memberWithPermissions, task.Instance, []dmsV1.OpPermissionType{dmsV1.OpPermissionTypeExecuteWorkflow})
			if err != nil {
				return
			}
		}
		return auditWorkflowUsers, executorWorkflowUsers
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	workflow, exist, err := s.GetLastWorkflow()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, fmt.Errorf("should exist at least one workflow after create workflow")))
	}

	go notification.NotifyWorkflow(string(workflow.ProjectId), workflow.WorkflowId, notification.WorkflowNotifyTypeCreate)

	go im.CreateApprove(string(workflow.ProjectId), workflow.WorkflowId)
	return nil
}

func batchExecuteWorkflows(c echo.Context) error {
	req := new(BatchExecuteWorkflowsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	sqlVersionId, err := strconv.ParseInt(c.Param("sql_version_id"), 10, 64)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	s := model.GetStorage()
	stageWorkflows, err := s.GetStageWorkflowsByWorkflowIds(uint(sqlVersionId), req.WorkflowIDs)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if len(stageWorkflows) != len(req.WorkflowIDs) {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, fmt.Errorf("unfound workflow")))
	}

	// sort by workflow sequence
	sort.Slice(stageWorkflows, func(i, j int) bool {
		return stageWorkflows[i].WorkflowSequence < stageWorkflows[j].WorkflowSequence
	})

	for _, execWorkflow := range stageWorkflows {
		workflow, err := dms.GetWorkflowDetailByWorkflowId(projectUid, execWorkflow.WorkflowID, s.GetWorkflowDetailWithoutInstancesByWorkflowID)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}

		user, err := controller.GetCurrentUser(c, dms.GetUser)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if err := PrepareForWorkflowExecution(c, projectUid, workflow, user); err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		workflowStatusChan, err := server.ExecuteTasksProcess(workflow.WorkflowId, projectUid, user)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		im.UpdateApprove(workflow.WorkflowId, user, model.ApproveStatusAgree, "")

		// 阻塞继续上线，直到获取到工单上线结果(状态)
		if <-workflowStatusChan != model.WorkflowStatusFinish {
			return controller.JSONBaseErrorReq(c, errors.New(errors.SQLVersionNotAllTasksExecutedSuccess, fmt.Errorf("workflow %s execution status is not finished and stop execution", workflow.Subject)))
		}

	}
	return controller.JSONBaseErrorReq(c, nil)
}

func retryExecWorkflow(c echo.Context) error {
	return nil
}

func batchAssociateWorkflowsWithVersion(c echo.Context) error {
	req := new(BatchAssociateWorkflowsWithVersionReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	versionIDStr := c.Param("sql_version_id")
	versionId, err := strconv.Atoi(versionIDStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	stageIDStr := c.Param("sql_version_stage_id")
	stageID, err := strconv.Atoi(stageIDStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = sqlversion.BatchAssociateWorkflowsWithStage(projectUid, uint(versionId), uint(stageID), req.WorkflowIDs)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return controller.JSONBaseErrorReq(c, nil)
}

func getWorkflowsThatCanBeAssociatedToVersion(c echo.Context) error {
	versionIDStr := c.Param("sql_version_id")
	versionId, err := strconv.Atoi(versionIDStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	stageIDStr := c.Param("sql_version_stage_id")
	stageID, err := strconv.Atoi(stageIDStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	workflows, err := sqlversion.GetWorkflowsThatCanBeAssociatedToVersionStage(uint(versionId), uint(stageID))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, &GetWorkflowsThatCanBeAssociatedToVersionResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertWorkflowToAssociateWorkflows(workflows),
	})
}

func convertWorkflowToAssociateWorkflows(workflows []*sqlversion.Workflow) []*AssociateWorkflows {
	ret := make([]*AssociateWorkflows, 0, len(workflows))
	for _, workflow := range workflows {
		ret = append(ret, &AssociateWorkflows{
			WorkflowID:   workflow.WorkflowID,
			WorkflowName: workflow.Subject,
			WorkflowDesc: workflow.Description,
		})
	}
	return ret
}
