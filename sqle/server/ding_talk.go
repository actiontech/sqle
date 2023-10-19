package server

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/model"
	imPkg "github.com/actiontech/sqle/sqle/pkg/im"
	"github.com/actiontech/sqle/sqle/pkg/im/dingding"
	"github.com/sirupsen/logrus"
)

type DingTalkJob struct {
	BaseJob
}

func NewDingTalkJob(entry *logrus.Entry) ServerJob {
	entry = entry.WithField("job", "ding_talk")
	j := &DingTalkJob{}
	j.BaseJob = *NewBaseJob(entry, 60*time.Second, j.dingTalkRotation)
	return j
}

func (j *DingTalkJob) dingTalkRotation(entry *logrus.Entry) {
	st := model.GetStorage()

	ims, err := st.GetAllIMConfig()
	if err != nil {
		entry.Errorf("get all im config failed, error: %v", err)
	}

	for _, im := range ims {
		switch im.Type {
		case model.ImTypeDingTalk:
			d := &dingding.DingTalk{
				AppKey:    im.AppKey,
				AppSecret: im.AppSecret,
			}

			dingTalkInstances, err := st.GetDingTalkInstByStatus(model.ApproveStatusInitialized)
			if err != nil {
				entry.Errorf("get ding talk status error: %v", err)
				continue
			}

			for _, dingTalkInstance := range dingTalkInstances {
				approval, err := d.GetApprovalDetail(dingTalkInstance.ApproveInstanceCode)
				if err != nil {
					entry.Errorf("get ding talk approval detail error: %v", err)
					continue
				}

				switch *approval.Result {
				case model.ApproveStatusAgree:
					workflow, exist, err := st.GetWorkflowDetailById(strconv.Itoa(int(dingTalkInstance.WorkflowId)))
					if err != nil {
						entry.Errorf("get workflow detail error: %v", err)
						continue
					}
					if !exist {
						entry.Errorf("workflow not exist, id: %d", dingTalkInstance.WorkflowId)
						continue
					}
					if workflow.Record.Status == model.WorkflowStatusCancel {
						entry.Errorf("workflow has canceled skip, id: %d", dingTalkInstance.WorkflowId)
						continue
					}

					instanceIds := make([]uint64, 0, len(workflow.Record.InstanceRecords))
					for _, item := range workflow.Record.InstanceRecords {
						instanceIds = append(instanceIds, item.InstanceId)
					}

					instances, err := dms.GetInstancesInProjectByIds(context.Background(), string(workflow.ProjectId), instanceIds)
					if err != nil {
						entry.Errorf("get instance error, %v", err)
						continue
					}
					instanceMap := map[uint64]*model.Instance{}
					for _, instance := range instances {
						instanceMap[instance.ID] = instance
					}
					for i, item := range workflow.Record.InstanceRecords {
						if instance, ok := instanceMap[item.InstanceId]; ok {
							workflow.Record.InstanceRecords[i].Instance = instance
						}
					}

					nextStep := workflow.NextStep()

					userId := *approval.OperationRecords[1].UserId
					user, err := getUserByUserId(d, userId, nil /* TODO workflow.CurrentStep().Assignees*/)
					if err != nil {
						entry.Errorf("get user by user id error: %v", err)
						continue
					}

					if err := ApproveWorkflowProcess(workflow, user, st); err != nil {
						entry.Errorf("approve workflow process error: %v", err)
						continue
					}

					dingTalkInstance.Status = model.ApproveStatusAgree
					if err := st.Save(&dingTalkInstance); err != nil {
						entry.Errorf("save ding talk instance error: %v", err)
						continue
					}

					if nextStep.Template.Typ != model.WorkflowStepTypeSQLExecute {
						imPkg.CreateApprove(strconv.Itoa(int(workflow.ID)))
					}

				case model.ApproveStatusRefuse:
					workflow, exist, err := st.GetWorkflowDetailById(strconv.Itoa(int(dingTalkInstance.WorkflowId)))
					if err != nil {
						entry.Errorf("get workflow detail error: %v", err)
						continue
					}
					if !exist {
						entry.Errorf("workflow not exist, id: %d", dingTalkInstance.WorkflowId)
						continue
					}
					if workflow.Record.Status == model.WorkflowStatusCancel {
						entry.Errorf("workflow has canceled skip, id: %d", dingTalkInstance.WorkflowId)
						continue
					}

					instanceIds := make([]uint64, 0, len(workflow.Record.InstanceRecords))
					for _, item := range workflow.Record.InstanceRecords {
						instanceIds = append(instanceIds, item.InstanceId)
					}

					instances, err := dms.GetInstancesInProjectByIds(context.Background(), string(workflow.ProjectId), instanceIds)
					if err != nil {
						entry.Errorf("notify workflow error, %v", err)
						continue
					}
					instanceMap := map[uint64]*model.Instance{}
					for _, instance := range instances {
						instanceMap[instance.ID] = instance
					}
					for i, item := range workflow.Record.InstanceRecords {
						if instance, ok := instanceMap[item.InstanceId]; ok {
							workflow.Record.InstanceRecords[i].Instance = instance
						}
					}

					var reason string
					if approval.OperationRecords[1] != nil && approval.OperationRecords[1].Remark != nil {
						reason = *approval.OperationRecords[1].Remark
					} else {
						reason = "审批拒绝"
					}

					userId := *approval.OperationRecords[1].UserId
					user, err := getUserByUserId(d, userId, nil /*TODO workflow.CurrentStep().Assignees*/)
					if err != nil {
						entry.Errorf("get user by user id error: %v", err)
						continue
					}

					if err := RejectWorkflowProcess(workflow, reason, user, st); err != nil {
						entry.Errorf("reject workflow process error: %v", err)
						continue
					}

					dingTalkInstance.Status = model.ApproveStatusRefuse
					if err := st.Save(&dingTalkInstance); err != nil {
						entry.Errorf("save ding talk instance error: %v", err)
						continue
					}
				default:
					// ding talk rotation, no action
				}
			}
		}
	}
}

func getUserByUserId(d *dingding.DingTalk, userId string, assignees []*model.User) (*model.User, error) {
	phone, err := d.GetMobileByUserID(userId)
	if err != nil {
		return nil, fmt.Errorf("get user mobile error: %v", err)
	}

	for _, assignee := range assignees {
		if assignee.Phone == phone {
			return assignee, nil
		}
	}

	return nil, fmt.Errorf("user not found, phone: %s", phone)
}
