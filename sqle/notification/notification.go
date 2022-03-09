package notification

import (
	"fmt"

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
			log.NewEntry().Error(err)
		}
	}
	return nil
}

func NotifyWorkflow(workflowId string) error {
	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowDetailById(workflowId)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("workflow not exits")
	}
	users := workflow.CurrentAssigneeUser()
	// workflow has been finished.
	if len(users) == 0 {
		return nil
	}
	return Notify(workflow, users)
}
