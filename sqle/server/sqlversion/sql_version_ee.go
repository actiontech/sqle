//go:build enterprise
// +build enterprise

package sqlversion

import (
	"fmt"
	"sort"

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
	Workflows   Workflows   // 该阶段纳管的工单
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
	workflows := make(Workflows, 0, len(modelStage.WorkflowVersionStage))
	for _, workflow := range modelStage.WorkflowVersionStage {
		workflows = append(workflows, &Workflow{
			ID:         workflow.ID,
			WorkflowID: workflow.WorkflowID,
			Sequence:   workflow.WorkflowSequence,
			Status:     workflow.WorkflowReleaseStatus,
		})
	}
	workflows.SortBySequence()
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

func (s Stage) CheckWorkflowExistInStage(workflow *model.Workflow) error {

	workflowIdMap := make(map[string]struct{})

	for _, workflow := range s.Workflows {
		workflowIdMap[workflow.WorkflowID] = struct{}{}
	}

	if _, exists := workflowIdMap[workflow.WorkflowId]; exists {
		return fmt.Errorf("can not associate workflow to stage, workflow already exist in this stage. stage name: %v, workflow subject: %v", s.Name, workflow.Subject)
	}
	return nil
}

func CheckWorkflowHasBoundWithStage(workflowID string) error {
	db := model.GetStorage()
	relation, exist, err := db.GetWorkflowVersionRelationByWorkflowId(workflowID)
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("workflow can only be bound with a stage, this workflow has bound with another stage of version, version id %v stage id %v", relation.SqlVersionID, relation.SqlVersionStageID)
	}
	return nil
}

// 数据源，数据源为SQL版本阶段中涉及的数据源，以关联关系的方式保存，SQL版本阶段:数据源 = 1:n
type Instance struct {
	ID             uint
	InstanceID     uint64 // 该阶段的数据源ID
	NextInstanceID uint64 // 下一阶段的数据源ID
}

type Workflows []*Workflow

// SQL工单，SQL工单被纳管至SQL版本管理中，以关联关系的方式保存，SQL版本阶段:工单 = 1:n
type Workflow struct {
	ID          uint
	WorkflowID  string // 工单ID
	Sequence    int    // 该阶段中工单的排序
	Status      string // 该工单的状态
	Subject     string
	Description string
	Status      string
	workflow    *model.Workflow
}

func (w Workflows) SortBySequence() {
	sort.Slice(w, func(i, j int) bool {
		return w[i].Sequence < w[j].Sequence
	})
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

	modelWorkflows, err := db.GetWorkflowsThatCanBeAssociatedToStage(instanceIdRange)
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
			Status:      modelWorkflow.Status,
		})
	}
	return workflows, nil
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

func BatchAssociateWorkflowsWithStage(projectUid string, versionID, stageID uint, workflowIds []string) error {
	db := model.GetStorage()
	modelStage, err := db.GetStageOfSQLVersion(versionID, stageID)
	if err != nil {
		return err
	}
	stage := ToServiceStage(modelStage)
	for _, workflowID := range workflowIds {
		// check if instance of workflow are entirely belongs to stage
		instanceIds, err := db.GetInstanceIdsByWorkflowID(workflowID)
		if err != nil {
			return err
		}
		if len(instanceIds) == 0 {
			return fmt.Errorf("the workflow does not use any instance")
		}
		err = stage.CheckStageContainsInstances(instanceIds)
		if err != nil {
			return err
		}
		// TODO At present, a workflow only supports binding to one stage, so it is only necessary to check whether the work order has been bound to a stage and annotate and retain the original code to detect whether the work order exists in this stage.
		// // check if workflow exist
		// _, exist, err := db.GetWorkflowByProjectAndWorkflowId(projectUid, workflowID)
		// if err != nil {
		// 	return err
		// }
		// if !exist {
		// 	return fmt.Errorf("can not associate a non-existent workflow with stage, workflow id: %v", workflowID)
		// }
		// // check if workflow exist in this stage
		// err = stage.CheckWorkflowExistInStage(workflow)
		// if err != nil {
		// 	return err
		// }
		// check if workflow has bound to other stage
		err = CheckWorkflowHasBoundWithStage(workflowID)
		if err != nil {
			return err
		}
	}

	return db.BatchCreateWorkflowVersionRelation(modelStage, workflowIds)
}

func CheckWorkflowExecutable(projectUid, workflowId string) (executable bool, reason string, err error) {

	db := model.GetStorage()
	modelStage, exist, err := db.GetStageOfTheWorkflow(workflowId)
	if err != nil {
		return false, "", err
	}
	if !exist {
		return true, "", nil
	}
	stage := ToServiceStage(modelStage)
	for _, workflow := range stage.Workflows {
		// 按顺序检查该工单之前的工单
		if workflow.WorkflowID == workflowId {
			return true, "", nil
		}
		if workflow.workflow == nil {
			workflow.workflow, exist, err = db.GetWorkflowByWorkflowId(workflow.WorkflowID)
			if err != nil {
				return false, fmt.Sprintf("when checking if workflow can execute, find a workflow before this workflow failed, workflow id is %v, err: %v", workflow.WorkflowID, err), err
			}
			if !exist {
				return false, fmt.Sprintf("when checking if workflow can execute, can not find a workflow before this workflow, workflow id is %v", workflow.WorkflowID), nil
			}
		}
		// 当发现该工单之前的工单的状态不是执行完成或者取消这两种状态时，不允许执行或定时执行该工单
		switch workflow.workflow.Record.Status {
		case model.WorkflowStatusFinish, model.WorkflowStatusCancel:
			continue
		default:
			return false, fmt.Sprintf("can not execute or scheduled execute workflow that bind with stage of sql version, before this workflow. there were still workflow with a %v status, sql version id is %v", workflow.workflow.Record.Status, modelStage.SqlVersionID), nil
		}
	}
	return true, "the previous workflows are executed as expected", nil
}
