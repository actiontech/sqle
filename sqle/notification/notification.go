package notification

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/actiontech/sqle/sqle/driver"
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

var Notifiers = []Notifier{}

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
	s := model.GetStorage()
	taskIds := w.workflow.GetTaskIds()
	tasks, err := s.GetTasksByIds(taskIds)
	if err != nil || len(tasks) <= 0 {
		return fmt.Sprintf(`
- 工单主题: %v
- 工单描述: %v
- 申请人: %v
- 创建时间: %v
- 读取工单任务内容失败，请通过SQLE界面确认工单状态
`,
			w.workflow.Subject,
			w.workflow.Desc,
			w.workflow.CreateUserName(),
			w.workflow.CreatedAt)
	}

	buf := bytes.Buffer{}
	head := fmt.Sprintf(`
- 工单主题: %v
- 工单描述: %v
- 申请人: %v
- 创建时间: %v
`,
		w.workflow.Subject,
		w.workflow.Desc,
		w.workflow.CreateUserName(),
		w.workflow.CreatedAt)
	buf.WriteString(head)

	for _, t := range tasks {
		buf.WriteString("\n--------------\n")
		buf.WriteString(w.buildNotifyBody(t))
	}

	return buf.String()
}

func (w *WorkflowNotification) buildNotifyBody(task *model.Task) string {
	instanceName := task.InstanceName()
	score := task.Score
	passRate := task.PassRate
	schema := task.Schema
	executeStartAt := task.ExecStartAt
	executeEndAt := task.ExecEndAt

	switch w.notifyType {
	case WorkflowNotifyTypeExecuteSuccess, WorkflowNotifyTypeExecuteFail:
		return fmt.Sprintf(`
- 数据源: %v
- schema: %v
- 上线开始时间: %v
- 上线结束时间: %v
`,
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
- 数据源: %v
- schema: %v
- 驳回原因: %v
`,
			instanceName,
			schema,
			reason,
		)
	default:
		return fmt.Sprintf(`
- 数据源: %v
- schema: %v
- 工单审核得分: %v
- 工单审核通过率：%v%%
`,
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

type AuditPlanNotification struct {
	auditPlan *model.AuditPlan
	report    *model.AuditPlanReportV2
}

func NewAuditPlanNotification(auditPlan *model.AuditPlan, report *model.AuditPlanReportV2) *AuditPlanNotification {
	return &AuditPlanNotification{
		auditPlan: auditPlan,
		report:    report,
	}
}

func (a *AuditPlanNotification) NotificationSubject() string {
	return fmt.Sprintf("SQLE扫描任务[%v]扫描结果[%v]", a.auditPlan.Name, a.report.AuditLevel)
}

func (a *AuditPlanNotification) NotificationBody() string {
	return fmt.Sprintf(`
- 扫描任务: %v
- 审核时间: %v
- 审核类型: %v
- 数据源: %v
- 数据库名: %v
- 审核得分: %v
- 审核通过率：%v
- 审核结果等级: %v
`,
		a.auditPlan.Name,
		a.report.CreatedAt.Format(time.RFC3339),
		a.auditPlan.Type,
		a.auditPlan.InstanceName,
		a.auditPlan.InstanceDatabase,
		a.report.Score,
		a.report.PassRate,
		a.report.AuditLevel,
	)
}

type TestNotify struct {
}

func (t *TestNotify) NotificationSubject() string {
	return "SQLE notification test"
}

func (t *TestNotify) NotificationBody() string {
	return "This is a SQLE test notification\nIf you receive this message, it only means that the message can be pushed"
}

func NotifyAuditPlan(apName string, report *model.AuditPlanReportV2) error {
	s := model.GetStorage()
	ap, _, err := s.GetAuditPlanByName(apName)
	if err != nil {
		return err
	}
	ap.CreateUser, _, err = s.GetUserByID(ap.CreateUserID)
	if err != nil {
		return err
	}

	if driver.RuleLevelLessOrEqual(ap.NotifyLevel, report.AuditLevel) {
		n := NewAuditPlanNotification(ap, report)
		return GetAuditPlanNotifier().Notify(n, ap)
	}

	return nil
}

var stdAuditPlanNotifier = NewAuditPlanNotifier()

func GetAuditPlanNotifier() *AuditPlanNotifier {
	return stdAuditPlanNotifier
}

type AuditPlanNotifier struct {
	lastSend      map[string] /*audit plan name*/ time.Time /*last send time*/
	mutex         *sync.RWMutex
	emailNotifier *EmailNotifier
}

func NewAuditPlanNotifier() *AuditPlanNotifier {
	return &AuditPlanNotifier{
		lastSend:      map[string]time.Time{},
		mutex:         &sync.RWMutex{},
		emailNotifier: &EmailNotifier{},
	}
}

func (n *AuditPlanNotifier) Notify(notification Notification, auditPlan *model.AuditPlan) error {
	if !n.shouldNotify(auditPlan) {
		return nil
	}

	err := n.Send(notification, auditPlan)
	if err != nil {
		return err
	}

	n.updateRecord(auditPlan.Name)
	return nil
}

func (n *AuditPlanNotifier) shouldNotify(auditPlan *model.AuditPlan) bool {
	n.mutex.RLock()
	last := n.lastSend[auditPlan.Name]
	n.mutex.RUnlock()
	return time.Now().After(last.Add(time.Duration(auditPlan.NotifyInterval) * time.Minute))
}

func (n *AuditPlanNotifier) Send(notification Notification, auditPlan *model.AuditPlan) (err error) {
	if auditPlan.EnableEmailNotify {
		err = n.sendEmail(notification, auditPlan.CreateUser)
		if err != nil {
			return err
		}
	}
	if auditPlan.EnableWebHookNotify {
		err = n.sendWebHook(notification, auditPlan.WebHookURL, auditPlan.WebHookTemplate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *AuditPlanNotifier) sendEmail(notification Notification, user *model.User) error {
	return n.emailNotifier.Notify(notification, []*model.User{user})
}

func (n *AuditPlanNotifier) updateRecord(auditPlanName string) {
	n.mutex.Lock()
	n.lastSend[auditPlanName] = time.Now()
	n.mutex.Unlock()
}
