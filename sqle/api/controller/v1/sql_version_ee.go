//go:build enterprise
// +build enterprise

package v1

import (
	"context"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	dms "github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
	dry "github.com/ungerik/go-dry"
	"gorm.io/gorm"
)

func createSqlVersion(c echo.Context) error {
	req := new(CreateSqlVersionReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	projectUid, err := dms.GetPorjectUIDByName(c.Request().Context(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	s := model.GetStorage()

	err = s.Tx(func(txDB *gorm.DB) error {
		sqlVersion := &model.SqlVersion{
			Version:     req.Version,
			Description: req.Desc,
			Status:      model.SqlVersionStatusReleased,
			ProjectId:   model.ProjectUID(projectUid),
		}
		err = txDB.Save(sqlVersion).Error
		if err != nil {
			return err
		}
		// 保存版本阶段
		versionStages := make([]*model.SqlVersionStage, 0, len(req.SqlVersionStage))
		stageDepMap := make(map[int][]CreateStagesInstanceDep)
		for _, stage := range req.SqlVersionStage {
			versionStages = append(versionStages, &model.SqlVersionStage{
				SqlVersionID:  sqlVersion.ID,
				Name:          stage.Name,
				StageSequence: stage.StageSequence,
			})
			deps := make([]CreateStagesInstanceDep, 0)
			for _, stageDep := range stage.CreateStagesInstanceDep {
				deps = append(deps, CreateStagesInstanceDep{
					StageInstanceID:     stageDep.StageInstanceID,
					NextStageInstanceID: stageDep.NextStageInstanceID,
				})
			}
			stageDepMap[stage.StageSequence] = deps
		}
		err = txDB.Save(versionStages).Error
		if err != nil {
			return err
		}

		// 保存阶段依赖关系
		stageDeps := make([]*model.SqlVersionStagesDependency, 0)
		for _, versionStage := range versionStages {
			nextStage, exist, err := s.GetNextSatgeByVersionIdAndSequence(txDB, versionStage.SqlVersionID, versionStage.StageSequence)
			if err != nil {
				return err
			}
			for _, dep := range stageDepMap[versionStage.StageSequence] {
				stageInst, err := getInstanceByStageInstanceID(c.Request().Context(), dep.StageInstanceID)
				if err != nil {
					return err
				}
				nextStageInst, err := getInstanceByStageInstanceID(c.Request().Context(), dep.NextStageInstanceID)
				if err != nil {
					return err
				}
				sqlVersionStagesDep := &model.SqlVersionStagesDependency{}
				if exist {
					sqlVersionStagesDep.SqlVersionStageID = versionStage.ID
					sqlVersionStagesDep.NextStageID = nextStage.ID
					sqlVersionStagesDep.StageInstanceID = stageInst.ID
					sqlVersionStagesDep.NextStageInstanceID = nextStageInst.ID
				} else {
					sqlVersionStagesDep.SqlVersionStageID = versionStage.ID
					sqlVersionStagesDep.StageInstanceID = stageInst.ID
				}
				stageDeps = append(stageDeps, sqlVersionStagesDep)
			}
		}
		err = txDB.Save(stageDeps).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return controller.JSONBaseErrorReq(c, nil)
}

func getInstanceByStageInstanceID(ctx context.Context, instanceID string) (*model.Instance, error) {
	if instanceID == "" || instanceID == "0" {
		return nil, nil
	}
	inst, exist, err := dms.GetInstancesById(ctx, instanceID)
	if !exist {
		return nil, errors.New(errors.DataConflict, ErrInstanceNotExist)
	} else if err != nil {
		return nil, errors.New(errors.DataConflict, err)
	}

	if !dry.StringInSlice(inst.DbType, driver.GetPluginManager().AllDrivers()) {
		return nil, errors.New(errors.DriverNotExist, &driverV2.DriverNotSupportedError{DriverTyp: inst.DbType})
	}
	return inst, nil
}

func getSqlVersionList(c echo.Context) error {

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
		"filter_project_id":         projectUid,
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
	return nil
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
	return nil
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
