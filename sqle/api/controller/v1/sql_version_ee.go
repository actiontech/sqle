//go:build enterprise
// +build enterprise

package v1

import (
	"net/http"
	"strconv"

	"github.com/actiontech/sqle/sqle/api/controller"
	dms "github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
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
	return nil
}

func batchExecuteTasksOnWorkflow(c echo.Context) error {
	return nil
}
func retryExecWorkflow(c echo.Context) error {
	return nil
}
func batchAssociateWorkflowsWithVersion(c echo.Context) error {
	return nil
}

func getWorkflowsThatCanBeAssociatedToVersion(c echo.Context) error {
	// projectID := c.Param("project_name")
	// versionID := c.Param("sql_version_id")
	// stageID := c.Param("sql_version_stage_id")

	return nil
}
