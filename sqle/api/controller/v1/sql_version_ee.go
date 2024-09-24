//go:build enterprise
// +build enterprise

package v1

import (
	"context"

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
	return nil
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
