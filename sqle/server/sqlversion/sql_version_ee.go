//go:build enterprise
// +build enterprise

package sqlversion

import (
	"fmt"

	"time"

	"github.com/actiontech/sqle/sqle/model"
)

// SQL版本，是项目的资源，项目:版本 = 1:n
type SQLVersion struct {
	ID         uint
	ProjectUID model.ProjectUID // 项目ID，关联了版本所属的项目
	Version    string           // 版本号
	Desc       string           // 版本信息描述
	Status     string           // 版本状态，发布中，已锁定
	LockTime   *time.Time       // 锁定时间，当前时间超过改时间，SQL版本状态为锁定
	Stages     Stages           // 版本阶段，一个版本对应多个阶段
}

type Stages []*Stage

// SQL版本阶段，是SQL版本的资源，SQL版本:SQL版本阶段 = 1:n
type Stage struct {
	ID          uint
	NextStageID uint        // 下一阶段的阶段ID
	Name        string      // 该阶段名称
	Sequence    int         // 该阶段的次序
	Instances   []*Instance // 该阶段关联的数据源
	Workflows   []*Workflow // 该阶段纳管的工单
}

func ToServiceStage(modelStage *model.SqlVersionStage) *Stage {
	instances := make([]*Instance, 0, len(modelStage.SqlVersionStagesDependency))
	for _, dependency := range modelStage.SqlVersionStagesDependency {
		instances = append(instances, &Instance{
			ID:             dependency.ID,
			InstanceID:     dependency.StageInstanceID,
			NextInstanceID: dependency.NextStageInstanceID,
		})
	}
	workflows := make([]*Workflow, 0, len(modelStage.WorkflowVersionStage))
	for _, workflow := range modelStage.WorkflowVersionStage {
		workflows = append(workflows, &Workflow{
			ID:         workflow.ID,
			WorkflowID: workflow.WorkflowID,
			Sequence:   workflow.WorkflowSequence,
		})
	}
	stage := &Stage{
		ID:        modelStage.ID,
		Name:      modelStage.Name,
		Sequence:  modelStage.StageSequence,
		Instances: instances,
		Workflows: workflows,
	}
	return stage
}

// 判断输入的数据源是否属于当前阶段的数据源的子集
func (s Stage) CheckStageContainsInstances(instanceIds []uint64) error {

	instanceMap := make(map[uint64]struct{})

	for _, instance := range s.Instances {
		instanceMap[instance.InstanceID] = struct{}{}
	}

	// 检查传入的 instances 是否都是 Stage 中的子集
	for _, instanceId := range instanceIds {
		if _, exists := instanceMap[instanceId]; !exists {
			return fmt.Errorf("can not attach workflow with sql version, instances of the workflow does not belong entirely to the first stage.")
		}
	}

	return nil
}

// 数据源，数据源为SQL版本阶段中涉及的数据源，以关联关系的方式保存，SQL版本阶段:数据源 = 1:n
type Instance struct {
	ID             uint
	InstanceID     uint64 // 该阶段的数据源ID
	NextInstanceID uint64 // 下一阶段的数据源ID
}

// SQL工单，SQL工单被纳管至SQL版本管理中，以关联关系的方式保存，SQL版本阶段:工单 = 1:n
type Workflow struct {
	ID          uint
	WorkflowID  string // 工单ID
	Sequence    int    // 该阶段中工单的排序
	Subject     string
	Description string
	workflow    *model.Workflow
}

func CheckInstanceInWorkflowCanAssociateToTheFirstStageOfVersion(versionID uint, instanceId []uint64) error {
	db := model.GetStorage()

	workflowInstanceIds := make([]uint64, 0, len(instanceId))
	for _, instanceId := range instanceId {
		workflowInstanceIds = append(workflowInstanceIds, instanceId)
	}

	// get the first stage of sql version
	modelFirstStage, err := db.GetFirstStageOfSQLVersion(versionID)
	if err != nil {
		return fmt.Errorf("when get first stage of sql version error: %w", err)
	}
	firstStage := ToServiceStage(modelFirstStage)
	err = firstStage.CheckStageContainsInstances(workflowInstanceIds)
	if err != nil {
		return err
	}

	return nil
}

func GetWorkflowsThatCanBeAssociatedToVersionStage(versionID, stageID uint) ([]*Workflow, error) {
	db := model.GetStorage()
	modelStage, err := db.GetStageOfSQLVersion(versionID, stageID)
	if err != nil {
		return nil, err
	}
	stage := ToServiceStage(modelStage)
	instanceIdRange := make([]uint64, 0, len(stage.Instances))
	for _, instance := range stage.Instances {
		instanceIdRange = append(instanceIdRange, instance.InstanceID)
	}
	
	excludeWorkflowIds := make([]string, 0, len(stage.Workflows))
	for _, workflow := range stage.Workflows {
		excludeWorkflowIds = append(excludeWorkflowIds, workflow.WorkflowID)
	}

	modelWorkflows, err := db.GetWorkflowsThatCanBeAssociatedToStage(instanceIdRange, excludeWorkflowIds)
	if err != nil {
		return nil, err
	}
	workflows := make([]*Workflow, 0, len(modelWorkflows))
	for _, modelWorkflow := range modelWorkflows {
		workflows = append(workflows, &Workflow{
			ID:          modelWorkflow.ID,
			WorkflowID:  modelWorkflow.WorkFlowID,
			Subject:     modelWorkflow.Subject,
			Description: modelWorkflow.Desc,
		})
	}
	return workflows, nil
}
