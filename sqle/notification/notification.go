package notification

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	v1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"github.com/actiontech/dms/pkg/dms-common/i18nPkg"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/dms"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Notification interface {
	NotificationSubject() i18nPkg.I18nStr
	NotificationBody() i18nPkg.I18nStr
}

type Notifier interface {
	Notify(Notification, []*model.User) error
}

type WorkflowNotifyConfig struct {
	SQLEUrl     *string
	ProjectName string
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
	WorkflowNotifyTypeCancel
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
	case WorkflowNotifyTypeCancel:
		return "cancel"
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

func GetWorkflowStepTypeDesc(s string) *i18n.Message {
	switch s {
	case model.WorkflowStepTypeSQLExecute:
		return locale.NotifyWorkflowStepTypeSQLExecute
	default:
		return locale.NotifyWorkflowStepTypeSQLAudit
	}
}

func (w *WorkflowNotification) NotificationSubject() i18nPkg.I18nStr {
	switch w.notifyType {
	case WorkflowNotifyTypeApprove, WorkflowNotifyTypeCreate:
		return locale.Bundle.LocalizeAllWithArgs(locale.NotifyWorkflowNotifyTypeWaiting, GetWorkflowStepTypeDesc(w.workflow.CurrentStep().Template.Typ))
	case WorkflowNotifyTypeReject:
		return locale.Bundle.LocalizeAll(locale.NotifyWorkflowNotifyTypeReject)
	case WorkflowNotifyTypeExecuteSuccess:
		return locale.Bundle.LocalizeAll(locale.NotifyWorkflowNotifyTypeExecuteSuccess)
	case WorkflowNotifyTypeExecuteFail:
		return locale.Bundle.LocalizeAll(locale.NotifyWorkflowNotifyTypeExecuteFail)
	case WorkflowNotifyTypeCancel:
		return locale.Bundle.LocalizeAll(locale.NotifyWorkflowNotifyTypeCancel)
	default:
		return locale.Bundle.LocalizeAll(locale.NotifyWorkflowNotifyTypeDefault)
	}
}

func (w *WorkflowNotification) NotificationBody() i18nPkg.I18nStr {
	bodyStr := make([]i18nPkg.I18nStr, 0)
	bodyStr = append(bodyStr, locale.Bundle.LocalizeAllWithArgs(locale.NotifyWorkflowBodyHead,
		w.workflow.Subject,
		w.config.ProjectName,
		w.workflow.WorkflowId,
		w.workflow.Desc,
		dms.GetUserNameWithDelTag(w.workflow.CreateUserId),
		w.workflow.CreatedAt,
	))

	s := model.GetStorage()
	taskIds := w.workflow.GetTaskIds()
	tasks, _, err := s.GetTasksByIds(taskIds)
	if err != nil || len(tasks) <= 0 {
		bodyStr = append(bodyStr, locale.Bundle.LocalizeAll(locale.NotifyWorkflowBodyWorkFlowErr))
		return locale.Bundle.JoinI18nStr(bodyStr, "\n")
	}

	if w.config.SQLEUrl != nil {
		link := fmt.Sprintf("%v/project/%v/exec-workflow/%v",
			strings.TrimRight(*w.config.SQLEUrl, "/"),
			w.workflow.ProjectId,
			w.workflow.WorkflowId,
		)
		bodyStr = append(bodyStr, locale.Bundle.LocalizeAllWithArgs(locale.NotifyWorkflowBodyLink, link))
	} else {
		bodyStr = append(bodyStr, locale.Bundle.LocalizeAllWithArgs(locale.NotifyWorkflowBodyLink, locale.NotifyWorkflowBodyConfigUrl))
	}

	instanceIds := make([]uint64, 0, len(tasks))
	for _, task := range tasks {
		instanceIds = append(instanceIds, task.InstanceId)
	}

	instances, err := dms.GetInstancesInProjectByIds(context.Background(), string(w.workflow.ProjectId), instanceIds)
	if err != nil {
		bodyStr = append(bodyStr, locale.Bundle.LocalizeAllWithArgs(locale.NotifyWorkflowBodyInstanceErr, err))
		return locale.Bundle.JoinI18nStr(bodyStr, "\n")
	}

	instanceMap := map[uint64]*model.Instance{}
	for _, instance := range instances {
		instanceMap[instance.ID] = instance
	}

	// Add a blank line before task details for better readability
	bodyStr = append(bodyStr, i18nPkg.ConvertStr2I18nAsDefaultLang(""))

	for _, t := range tasks {
		if instance, ok := instanceMap[t.InstanceId]; ok {
			t.Instance = instance
		}

		bodyStr = append(bodyStr, i18nPkg.ConvertStr2I18nAsDefaultLang("────────────────"), w.buildNotifyBody(t))
	}

	return locale.Bundle.JoinI18nStr(bodyStr, "\n")
}

func (w *WorkflowNotification) buildNotifyBody(task *model.Task) i18nPkg.I18nStr {
	instanceName := task.InstanceName()
	score := task.Score
	// passRate := task.PassRate
	schema := task.Schema
	executeStartAt := task.ExecStartAt
	executeEndAt := task.ExecEndAt

	var res []i18nPkg.I18nStr
	res = append(res, locale.Bundle.LocalizeAllWithArgs(locale.NotifyWorkflowBodyInstanceAndSchema, instanceName, schema))

	switch w.notifyType {
	case WorkflowNotifyTypeExecuteSuccess, WorkflowNotifyTypeExecuteFail:
		res = append(res, locale.Bundle.LocalizeAllWithArgs(locale.NotifyWorkflowBodyStartEnd, executeStartAt, executeEndAt))
	case WorkflowNotifyTypeReject:
		var reason string
		for _, step := range w.workflow.Record.Steps {
			if step.State == model.WorkflowStatusReject {
				reason = step.Reason
				break
			}
		}
		res = append(res, locale.Bundle.LocalizeAllWithArgs(locale.NotifyWorkflowBodyReason, reason))
	default:
		// res = append(res, locale.Bundle.LocalizeAllWithArgs(locale.NotifyWorkflowBodyReport, score, passRate*100))
		res = append(res, locale.Bundle.LocalizeAllWithArgs(locale.NotifyWorkflowBodyReport, score))
	}
	return locale.Bundle.JoinI18nStr(res, "\n")
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
	// if workflow is cancelled, the creator needs to be notified.
	case WorkflowNotifyTypeCancel:
		return []string{
			w.workflow.CreateUserId,
		}
		// if workflow is executed, the creator and executor needs to be notified.
	case WorkflowNotifyTypeExecuteSuccess, WorkflowNotifyTypeExecuteFail:
		// 获取该工单对应数据源上有工单审核权限的所有用户
		auditUsers, err := w.getAuditUsersForWorkflowInstances()
		if err != nil {
			log.NewEntry().Errorf("get audit users for workflow instances error: %v", err)
			return []string{w.workflow.CreateUserId}
		}
		return auditUsers
	default:
		return []string{}
	}
}

// getAuditUsersForWorkflowInstances 获取工单对应数据源上有工单审核权限的所有用户
func (w *WorkflowNotification) getAuditUsersForWorkflowInstances() ([]string, error) {
	instanceIds := w.workflow.GetInstanceIds()
	// 获取实例信息
	instances, err := dms.GetInstancesInProjectByIds(context.Background(), string(w.workflow.ProjectId), instanceIds)
	if err != nil {
		return nil, fmt.Errorf("get instances by ids error: %v", err)
	}
	// 获取项目成员和权限信息
	memberWithPermissions, _, err := dmsobject.ListMembersInProject(context.Background(), controller.GetDMSServerAddress(), v1.ListMembersForInternalReq{
		ProjectUid: string(w.workflow.ProjectId),
		PageSize:   999,
		PageIndex:  1,
	})
	if err != nil {
		return nil, fmt.Errorf("list members in project error: %v", err)
	}

	// 收集所有有审核权限的用户ID
	auditUserIds := make(map[string]struct{})
	for _, instance := range instances {
		// 获取该实例上有审核权限的用户
		userIds, err := w.getCanOpInstanceUsers(memberWithPermissions, instance, []v1.OpPermissionType{v1.OpPermissionTypeAuditWorkflow})
		if err != nil {
			log.NewEntry().Errorf("get can op instance users error: %v", err)
			continue
		}

		// 将用户ID添加到集合中
		for _, user := range userIds {
			auditUserIds[user] = struct{}{}
		}
	}
	auditUserIds[w.workflow.CreateUserId] = struct{}{}
	// 转换为字符串切片
	result := make([]string, 0, len(auditUserIds))
	for userId := range auditUserIds {
		result = append(result, userId)
	}
	return result, nil
}

// getCanOpInstanceUsers 获取实例上有指定权限的用户
func (w *WorkflowNotification) getCanOpInstanceUsers(memberWithPermissions []*v1.ListMembersForInternalItem, instance *model.Instance, opPermissions []v1.OpPermissionType) (userIds []string, err error) {
	opMapUsers := make(map[string]struct{}, 0)
	for _, memberWithPermission := range memberWithPermissions {
		for _, memberOpPermission := range memberWithPermission.MemberOpPermissionList {
			if w.CanOperationInstance([]v1.OpPermissionItem{memberOpPermission}, opPermissions, instance) {
				userId := memberWithPermission.User.Uid
				if _, ok := opMapUsers[userId]; !ok {
					opMapUsers[userId] = struct{}{}
					userIds = append(userIds, userId)
				}
			}
		}
	}
	return userIds, nil
}

func (w *WorkflowNotification) CanOperationInstance(userOpPermissions []v1.OpPermissionItem, needOpPermissionTypes []v1.OpPermissionType, instance *model.Instance) bool {
	for _, userOpPermission := range userOpPermissions {
		// 全局操作权限用户
		if userOpPermission.OpPermissionType == v1.OpPermissionTypeGlobalManagement {
			return true
		}

		// 项目管理员可以访问所有数据源
		if userOpPermission.OpPermissionType == v1.OpPermissionTypeProjectAdmin {
			return true
		}

		// 动作权限(创建、审核、上线工单等)
		hasPrivilege := false
		for _, needOpPermissionType := range needOpPermissionTypes {
			if needOpPermissionType == userOpPermission.OpPermissionType {
				hasPrivilege = true
				break
			}
		}
		if !hasPrivilege {
			continue
		}
		// 对象权限(指定数据源)
		if userOpPermission.RangeType == v1.OpRangeTypeDBService {
			for _, id := range userOpPermission.RangeUids {
				if id == instance.GetIDStr() {
					return true
				}
			}
		}
	}
	return false
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
	project, err := dms.GetProjectByID(string(workflow.ProjectId))
	if err != nil {
		log.NewEntry().Errorf("notify workflow get project by id error: %v", err)
	}
	config.ProjectName = project.Name
	wn := NewWorkflowNotification(workflow, wt, config)
	userIds := wn.notifyUser()
	// workflow has been finished.
	if len(userIds) == 0 {
		return
	}

	err = Notify(wn, userIds)
	if err != nil {
		log.NewEntry().Errorf("notify workflow error, %v", err)
	}
}

func NotifyWorkflow(projectId, workflowId string, wt WorkflowNotifyType) {
	logger := log.NewEntry()
	s := model.GetStorage()
	// 确认推送功能是否开启
	config, err := s.GetReportPushConfigInProjectByType(projectId, model.TypeWorkflow)
	if err != nil {
		logger.Errorf("get report push config failed: %v", err)
		return
	}
	if !config.Enabled {
		return
	}
	workflow, err := dms.GetWorkflowDetailByWorkflowId(projectId, workflowId, s.GetWorkflowDetailWithoutInstancesByWorkflowID)
	if err != nil {
		logger.Error("notify workflow error, workflow not exits")
		return
	}

	go func() { notifyWorkflowWebhook(workflow, wt) }()

	sqleUrl, err := dms.GetSqleUrl(context.TODO())
	if err != nil {
		logger.Errorf("get sqle url error, %v", err)
		return
	}
	notifyWorkflow(sqleUrl, workflow, wt)
	// 更新最新推送时间
	config.ReportPushConfigRecord.ReportPushConfigID = config.ID
	config.ReportPushConfigRecord.LastPushTime = time.Now()
	err = s.Save(&config.ReportPushConfigRecord)
	if err != nil {
		logger.Errorf("update report push config time failed: %v", err)
	}
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

func (a *AuditPlanNotification) NotificationSubject() i18nPkg.I18nStr {
	return locale.Bundle.LocalizeAllWithArgs(locale.NotifyAuditPlanSubject, a.auditPlan.Name, a.report.AuditLevel)
}

func (a *AuditPlanNotification) NotificationBody() i18nPkg.I18nStr {
	var linkInBody i18nPkg.I18nStr
	if a.config.SQLEUrl != nil && a.auditPlan.ProjectId != "" {
		link := fmt.Sprintf("%v/project/%v/auditPlan/detail/%v/report/%v",
			strings.TrimRight(*a.config.SQLEUrl, "/"),
			a.auditPlan.ProjectId,
			a.auditPlan.Name,
			a.report.ID,
		)
		linkInBody = locale.Bundle.LocalizeAllWithArgs(locale.NotifyAuditPlanBodyLink, link)
	}

	body := locale.Bundle.LocalizeAllWithArgs(locale.NotifyAuditPlanBody,
		a.auditPlan.Name,
		a.report.CreatedAt.Format(time.RFC3339),
		a.auditPlan.Type,
		a.auditPlan.InstanceName,
		a.auditPlan.InstanceDatabase,
		a.report.Score,
		a.report.PassRate,
		a.report.AuditLevel,
		linkInBody,
	)
	return body
}

type TestNotify struct {
}

func (t *TestNotify) NotificationSubject() i18nPkg.I18nStr {
	return i18nPkg.ConvertStr2I18nAsDefaultLang("SQLE notification test")
}

func (t *TestNotify) NotificationBody() i18nPkg.I18nStr {
	return i18nPkg.ConvertStr2I18nAsDefaultLang("This is a SQLE test notification\nIf you receive this message, it only means that the message can be pushed")
}

func getAPNotifyConfig() (AuditPlanNotifyConfig, error) {
	config := AuditPlanNotifyConfig{}

	url, err := dms.GetSqleUrl(context.TODO())
	if err != nil {
		return config, err
	}

	if len(url) > 0 {
		config.SQLEUrl = &url

		// dms-todo: 从 dms 获取 project 名称，但最终考虑将告警移走.
		// project, _, err := s.GetProjectByID(ap.ProjectId)
		// if err != nil {
		// 	return err
		// }
		// config.ProjectName = &project.Name
	}
	return config, nil
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
	config, err := getAPNotifyConfig()
	if err != nil {
		return err
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

func NotifyAuditPlanWebhook(auditPlan *model.AuditPlan, report *model.AuditPlanReportV2) {
	config, err := getAPNotifyConfig()
	if err != nil {
		log.NewEntry().Errorf("audit plan webhook failed: %v", err)
		return
	}

	err = auditPlanSendRequest(auditPlan, report, config)
	if err != nil {
		log.NewEntry().Errorf("audit plan webhook failed: %v", err)
	}
}
