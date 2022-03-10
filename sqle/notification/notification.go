package notification

import (
	"fmt"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/log"

	"github.com/actiontech/sqle/sqle/model"
)

type Notification interface {
	NotificationSubject() string
	NotificationBody() string
}

type Notifier interface {
	Notify(Notification, []*model.User) error
}

var Notifiers []Notifier

func init() {
	Notifiers = []Notifier{
		&EmailNotifier{},
	}
}

func Notify(notification Notification, users []*model.User) error {
	for _, n := range Notifiers {
		err := n.Notify(notification, users)
		if err != nil {
			return err
		}
	}
	return nil
}

type WorkflowNotifyType int

const (
	WorkflowNotifyTypeCreate WorkflowNotifyType = iota
	WorkflowNotifyTypeApprove
	WorkflowNotifyTypeReject
	WorkflowNotifyTypeExecuteSuccess
	WorkflowNotifyTypeExecuteFail
)

type WorkflowNotification struct {
	notifyType WorkflowNotifyType
	workflow   *model.Workflow
}

func NewWorkflowNotification(w *model.Workflow, notifyType WorkflowNotifyType) *WorkflowNotification {
	return &WorkflowNotification{
		notifyType: notifyType,
		workflow:   w,
	}
}

func GetWorkflowStepTypeDesc(s string) string {
	switch s {
	case model.WorkflowStepTypeSQLExecute:
		return "上线"
	default:
		return "审批"
	}
}

func (w *WorkflowNotification) NotificationSubject() string {
	switch w.notifyType {
	case WorkflowNotifyTypeApprove, WorkflowNotifyTypeCreate:
		return fmt.Sprintf("SQL工单待%s", GetWorkflowStepTypeDesc(w.workflow.CurrentStep().Template.Typ))
	case WorkflowNotifyTypeReject:
		return "SQL工单已被驳回"
	case WorkflowNotifyTypeExecuteSuccess:
		return "SQL工单上线成功"
	case WorkflowNotifyTypeExecuteFail:
		return "SQL工单上线失败"
	default:
		return "SQL工单未知请求"
	}
}

func (w *WorkflowNotification) NotificationBody() string {
	var (
		instanceName   string
		schema         string
		score          int32
		passRate       float64
		executeStartAt *time.Time
		executeEndAt   *time.Time
	)
	s := model.GetStorage()
	task, exist, err := s.GetTaskById(strconv.Itoa(int(w.workflow.Record.TaskId)))
	if err == nil && exist {
		instanceName = task.InstanceName()
		score = task.Score
		passRate = task.PassRate
		schema = task.Schema
		executeStartAt = task.ExecStartAt
		executeEndAt = task.ExecEndAt
	}
	switch w.notifyType {
	case WorkflowNotifyTypeExecuteSuccess, WorkflowNotifyTypeExecuteFail:
		return fmt.Sprintf(`
- 工单主题: %v
- 工单描述: %v
- 申请人: %v
- 创建时间: %v
- 数据源: %v
- schema: %v
- 上线开始时间: %v
- 上线结束时间: %v
`,
			w.workflow.Subject,
			w.workflow.Desc,
			w.workflow.CreateUserName(),
			w.workflow.CreatedAt,
			instanceName,
			schema,
			executeStartAt,
			executeEndAt,
		)
	case WorkflowNotifyTypeReject:
		var reason string
		for _, step := range w.workflow.Record.Steps {
			if step.State == model.WorkflowStatusReject {
				reason = step.Reason
				break
			}
		}
		return fmt.Sprintf(`
- 工单主题: %v
- 工单描述: %v
- 申请人: %v
- 创建时间: %v
- 数据源: %v
- schema: %v
- 驳回原因: %v
`,
			w.workflow.Subject,
			w.workflow.Desc,
			w.workflow.CreateUserName(),
			w.workflow.CreatedAt,
			instanceName,
			schema,
			reason,
		)
	default:
		return fmt.Sprintf(`
- 工单主题: %v
- 工单描述: %v
- 申请人: %v
- 创建时间: %v
- 数据源: %v
- schema: %v
- 工单审核得分: %v
- 工单审核通过率：%v%%
`,
			w.workflow.Subject,
			w.workflow.Desc,
			w.workflow.CreateUserName(),
			w.workflow.CreatedAt,
			instanceName,
			schema,
			score,
			passRate*100,
		)
	}
}

func (w *WorkflowNotification) notifyUser() []*model.User {
	switch w.notifyType {
	case WorkflowNotifyTypeApprove, WorkflowNotifyTypeCreate:
		return w.workflow.CurrentAssigneeUser()

	// if workflow is rejected, the creator needs to be notified.
	case WorkflowNotifyTypeReject:
		return []*model.User{
			w.workflow.CreateUser,
		}
		// if workflow is executed, the creator and executor needs to be notified.
	case WorkflowNotifyTypeExecuteSuccess, WorkflowNotifyTypeExecuteFail:
		users := []*model.User{
			w.workflow.CreateUser,
		}
		if executeUser := w.workflow.FinalStep().OperationUser; executeUser != nil {
			users = append(users, executeUser)
		}
		return users
	default:
		return []*model.User{}
	}
}

func NotifyWorkflow(workflowId string, wt WorkflowNotifyType) {
	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowDetailById(workflowId)
	if err != nil {
		log.NewEntry().Errorf("notify workflow error, %v", err)
	}
	if !exist {
		log.NewEntry().Error("notify workflow error, workflow not exits")
	}

	wn := NewWorkflowNotification(workflow, wt)
	users := wn.notifyUser()
	// workflow has been finished.
	if len(users) == 0 {
		return
	}
	err = Notify(wn, users)
	if err != nil {
		log.NewEntry().Errorf("notify workflow error, %v", err)
	}
}
