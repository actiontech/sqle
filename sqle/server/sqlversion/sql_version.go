package sqlversion

import (
	"time"

	"github.com/actiontech/sqle/sqle/model"
)

type SqlVersion struct {
	ID              uint
	Version         string             // 版本号
	Desc            string             // 版本信息描述
	Status          string             // 版本状态
	LockTime        *time.Time         // 锁定时间
	ProjectId       model.ProjectUID   // 项目ID，关联了版本所属的项目 项目:版本 = 1:n
	SqlVersionStage []*SqlVersionStage // 版本阶段，一个版本对应多个阶段
}

type SqlVersionStage struct {
	ID                         uint
	Name                       string                        // 版本阶段的名称
	SqlVersionID               uint                          // 版本ID，关联了版本和阶段，版本:阶段=1:n
	StageSequence              int                           // 版本阶段的排序
	SqlVersionStagesDependency []*SqlVersionStagesDependency // 版本阶段
	WorkflowReleaseStage       []*WorkflowVersionStage       // 一个版本阶段对应多个工单
}

// sql 版本的某一个阶段的数据源依赖，以及下一阶段的关联关系
type SqlVersionStagesDependency struct {
	ID                  uint
	SqlVersionStageID   uint // 版本阶段ID，关联了版本阶段，版本阶段
	NextStageID         uint // 下一阶段的阶段ID
	StageInstanceID     uint // 该阶段的数据源ID
	NextStageInstanceID uint // 该阶段对应下一阶段的数据源ID
}

// 工单、版本、版本阶段 关联关系
type WorkflowVersionStage struct {
	ID                uint
	WorkflowID        string // 工单ID
	SqlVersionID      uint   // 版本ID
	SqlVersionStageID uint   // 版本阶段ID
	WorkflowSequence  int    // 该阶段中工单的排序
}

// todo 完成该函数@WinfredLin
func AttachWorkflowWithTheFirstStageOfSqlVersion(sqlVersionID, workflowID string) error {
	// s := model.GetStorage()
	// currentStages, err := s.GetStagesOfSqlVersion(sqlVersionID)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// todo 完成转换函数@WinfredLin
func toServiceStages(modelStages []model.SqlVersionStage) []SqlVersionStage {
	serviceStage := make([]SqlVersionStage, 0, len(modelStages))
	for _, stage := range modelStages {
		stageDependency := make([]*SqlVersionStagesDependency, 0, len(stage.SqlVersionStagesDependency))

		serviceStage = append(serviceStage, SqlVersionStage{
			SqlVersionID:               stage.SqlVersionID,
			Name:                       stage.Name,
			StageSequence:              stage.StageSequence,
			SqlVersionStagesDependency: stageDependency,
		})
	}
	return serviceStage
}
