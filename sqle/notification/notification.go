package notification

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	v1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/dms"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
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

type WorkflowNotifyConfig struct {
	SQLEUrl *string
}

func Notify(notification Notification, userIds []string) error {
	return dmsobject.Notify(context.TODO(), controller.GetDMSServerAddress(), v1.NotificationReq{
		Notification: &v1.Notification{
			NotificationSubject: notification.NotificationSubject(),
			NotificationBody:    notification.NotificationBody(),
			UserUids:            userIds,
		},
	})
}

type WorkflowNotifyType int

const (
	WorkflowNotifyTypeCreate WorkflowNotifyType = iota
	WorkflowNotifyTypeApprove
	WorkflowNotifyTypeReject
	WorkflowNotifyTypeExecuteSuccess
	WorkflowNotifyTypeExecuteFail
)

func getWorkflowNotifyTypeAction(wt WorkflowNotifyType) string {
	switch wt {
	case WorkflowNotifyTypeCreate:
		return "create"
	case WorkflowNotifyTypeApprove:
		return "approve"
	case WorkflowNotifyTypeReject:
		return "reject"
	case WorkflowNotifyTypeExecuteSuccess:
		return "exec_success"
	case WorkflowNotifyTypeExecuteFail:
		return "exec_failed"
	}
	return "unknown"
}

type WorkflowNotification struct {
	notifyType WorkflowNotifyType
	workflow   *model.Workflow
	config     WorkflowNotifyConfig
}

func NewWorkflowNotification(w *model.Workflow, notifyType WorkflowNotifyType, config WorkflowNotifyConfig) *WorkflowNotification {
	return &WorkflowNotification{
		notifyType: notifyType,
		workflow:   w,
		config:     config,
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
	tasks, _, err := s.GetTasksByIds(taskIds)
	if err != nil || len(tasks) <= 0 {
		return fmt.Sprintf(`
- 工单主题: %v
- 工单ID: %v
- 工单描述: %v
- 申请人: %v
- 创建时间: %v
- 读取工单任务内容失败，请通过SQLE界面确认工单状态
`,
			w.workflow.Subject,
			w.workflow.WorkflowId,
			w.workflow.Desc,
			dms.GetUserNameWithDelTag(w.workflow.CreateUserId),
			w.workflow.CreatedAt)
	}

	buf := bytes.Buffer{}
	head := fmt.Sprintf(`
- 工单主题: %v
- 工单ID: %v
- 工单描述: %v
- 申请人: %v
- 创建时间: %v`,
		w.workflow.Subject,
		w.workflow.WorkflowId,
		w.workflow.Desc,
		dms.GetUserNameWithDelTag(w.workflow.CreateUserId),
		w.workflow.CreatedAt)
	buf.WriteString(head)

	if w.config.SQLEUrl != nil {
		buf.WriteString(fmt.Sprintf("\n- 工单链接: %v/project/%v/order/%v",
			strings.TrimRight(*w.config.SQLEUrl, "/"),
			w.workflow.ProjectId,
			w.workflow.WorkflowId,
		))
	} else {
		buf.WriteString("\n- 工单链接: 请在系统设置-全局配置中补充全局url")
	}

	instanceIds := make([]uint64, 0, len(tasks))
	for _, task := range tasks {
		instanceIds = append(instanceIds, task.InstanceId)
	}

	instances, err := dms.GetInstancesInProjectByIds(context.Background(), string(w.workflow.ProjectId), instanceIds)
	if err != nil {
		buf.WriteString(fmt.Sprintf("\n 获取数据源实例失败: %v\n", err))
		return buf.String()
	}

	instanceMap := map[uint64]*model.Instance{}
	for _, instance := range instances {
		instanceMap[instance.ID] = instance
	}

	for _, t := range tasks {
		if instance, ok := instanceMap[t.InstanceId]; ok {
			t.Instance = instance
		}

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

func (w *WorkflowNotification) notifyUser() []string {
	switch w.notifyType {
	case WorkflowNotifyTypeApprove, WorkflowNotifyTypeCreate:
		return w.workflow.CurrentAssigneeUser()

	// if workflow is rejected, the creator needs to be notified.
	case WorkflowNotifyTypeReject:
		return []string{
			w.workflow.CreateUserId,
		}
		// if workflow is executed, the creator and executor needs to be notified.
	case WorkflowNotifyTypeExecuteSuccess, WorkflowNotifyTypeExecuteFail:
		users := []string{
			w.workflow.CreateUserId,
		}
		if executeUser := w.workflow.FinalStep().OperationUserId; executeUser != "" {
			users = append(users, executeUser)
		}
		return users
	default:
		return []string{}
	}
}

func notifyWorkflowWebhook(workflow *model.Workflow, wt WorkflowNotifyType) {
	// dms-todo 使用projectid代替name
	err := workflowSendRequest(getWorkflowNotifyTypeAction(wt), workflow)
	if err != nil {
		log.NewEntry().Errorf("workflow webhook failed: %v", err)
	}
}

func notifyWorkflow(sqleUrl string, workflow *model.Workflow, wt WorkflowNotifyType) {
	config := WorkflowNotifyConfig{}
	if len(sqleUrl) > 0 {
		config.SQLEUrl = &sqleUrl
	}
	wn := NewWorkflowNotification(workflow, wt, config)
	userIds := wn.notifyUser()
	// workflow has been finished.
	if len(userIds) == 0 {
		return
	}

	err := Notify(wn, userIds)
	if err != nil {
		log.NewEntry().Errorf("notify workflow error, %v", err)
	}
}

func NotifyWorkflow(projectId, workflowId string, wt WorkflowNotifyType) {
	s := model.GetStorage()
	workflow, err := dms.GetWorkflowDetailByWorkflowId(projectId, workflowId, s.GetWorkflowDetailWithoutInstancesByWorkflowID)
	if err != nil {
		log.NewEntry().Error("notify workflow error, workflow not exits")
		return
	}

	go func() { notifyWorkflowWebhook(workflow, wt) }()

	sqleUrl, err := s.GetSqleUrl()
	if err != nil {
		log.NewEntry().Errorf("get sqle url error, %v", err)
		return
	}
	notifyWorkflow(sqleUrl, workflow, wt)
}

type AuditPlanNotification struct {
	auditPlan *model.AuditPlan
	report    *model.AuditPlanReportV2
	config    AuditPlanNotifyConfig
}

type AuditPlanNotifyConfig struct {
	SQLEUrl     *string
	ProjectName *string
}

func NewAuditPlanNotification(auditPlan *model.AuditPlan, report *model.AuditPlanReportV2, config AuditPlanNotifyConfig) *AuditPlanNotification {
	return &AuditPlanNotification{
		auditPlan: auditPlan,
		report:    report,
		config:    config,
	}
}

func (a *AuditPlanNotification) NotificationSubject() string {
	return fmt.Sprintf("SQLE扫描任务[%v]扫描结果[%v]", a.auditPlan.Name, a.report.AuditLevel)
}

func (a *AuditPlanNotification) NotificationBody() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf(`
- 扫描任务: %v
- 审核时间: %v
- 审核类型: %v
- 数据源: %v
- 数据库名: %v
- 审核得分: %v
- 审核通过率：%v
- 审核结果等级: %v`,
		a.auditPlan.Name,
		a.report.CreatedAt.Format(time.RFC3339),
		a.auditPlan.Type,
		a.auditPlan.InstanceName,
		a.auditPlan.InstanceDatabase,
		a.report.Score,
		a.report.PassRate,
		a.report.AuditLevel,
	))

	if a.config.SQLEUrl != nil && a.auditPlan.ProjectId != "" {
		builder.WriteString(fmt.Sprintf("\n- 扫描任务链接: %v/project/%v/auditPlan/detail/%v/report/%v",
			strings.TrimRight(*a.config.SQLEUrl, "/"),
			a.auditPlan.ProjectId,
			a.auditPlan.Name,
			a.report.ID,
		))
	}

	return builder.String()
}

type TestNotify struct {
}

func (t *TestNotify) NotificationSubject() string {
	return "SQLE notification test"
}

func (t *TestNotify) NotificationBody() string {
	return "This is a SQLE test notification\nIf you receive this message, it only means that the message can be pushed"
}

func NotifyAuditPlan(auditPlanId uint, report *model.AuditPlanReportV2) error {
	s := model.GetStorage()
	ap, _, err := s.GetAuditPlanById(auditPlanId)
	if err != nil {
		return err
	}
	// ap.CreateUser, _, err = s.GetUserByID(ap.CreateUserID)
	// if err != nil {
	// 	return err
	// }
	url, err := s.GetSqleUrl()
	if err != nil {
		return err
	}

	config := AuditPlanNotifyConfig{}
	if len(url) > 0 {
		config.SQLEUrl = &url

		// dms-todo: 从 dms 获取 project 名称，但最终考虑将告警移走.
		// project, _, err := s.GetProjectByID(ap.ProjectId)
		// if err != nil {
		// 	return err
		// }
		// config.ProjectName = &project.Name
	}

	if driverV2.RuleLevelLessOrEqual(ap.NotifyLevel, report.AuditLevel) {
		n := NewAuditPlanNotification(ap, report, config)
		return GetAuditPlanNotifier().Notify(n, ap)
	}

	return nil
}

var stdAuditPlanNotifier = NewAuditPlanNotifier()

func GetAuditPlanNotifier() *AuditPlanNotifier {
	return stdAuditPlanNotifier
}

type AuditPlanNotifier struct {
	lastSend map[string] /*audit plan name*/ time.Time /*last send time*/
	mutex    *sync.RWMutex
	// emailNotifier *EmailNotifier
}

func NewAuditPlanNotifier() *AuditPlanNotifier {
	return &AuditPlanNotifier{
		lastSend: map[string]time.Time{},
		mutex:    &sync.RWMutex{},
		// emailNotifier: &EmailNotifier{},
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
		user, err := dms.GetUser(context.TODO(), auditPlan.CreateUserID, controller.GetDMSServerAddress())
		if err != nil {
			log.NewEntry().Errorf("get user error, %v", err)
			return err
		}
		err = n.sendEmail(notification, user)
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
	// dms-todo 只发送邮件告警
	return dmsobject.Notify(context.TODO(), controller.GetDMSServerAddress(), v1.NotificationReq{
		Notification: &v1.Notification{
			NotificationSubject: notification.NotificationSubject(),
			NotificationBody:    notification.NotificationBody(),
			UserUids:            []string{user.GetIDStr()},
		},
	})
}

func (n *AuditPlanNotifier) updateRecord(auditPlanName string) {
	n.mutex.Lock()
	n.lastSend[auditPlanName] = time.Now()
	n.mutex.Unlock()
}
