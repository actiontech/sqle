package misc

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"

	"gopkg.in/gomail.v2"
)

func SendEmailIfConfigureSMTP(workflowId string) error {
	s := model.GetStorage()
	smtpC, exist, err := s.GetSMTPConfiguration()
	if err != nil {
		return err
	}
	if !exist {
		return nil
	}

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
	var emails []string
	for _, user := range users {
		if user.Email != "" {
			emails = append(emails, user.Email)
		}
	}
	// no user has configured email, don't send.
	if len(emails) == 0 {
		return nil
	}
	message := gomail.NewMessage()
	message.SetHeader("From", smtpC.Username)
	message.SetHeader("To", emails...)
	message.SetHeader("Subject", `SQL工单审批请求`)
	body := fmt.Sprintf(`
您有一个SQL工单待%v:
- 工单主题: %v
- 工单描述: %v
- 申请人: %v
`, model.GetWorkflowStepTypeDesc(workflow.CurrentStep().Template.Typ),
		workflow.Subject, workflow.Desc, workflow.CreateUserName())
	message.SetBody("text/html",
		strings.Replace(body, "\n", "<br/>\n", -1))

	port, _ := strconv.Atoi(smtpC.Port)
	dialer := gomail.NewDialer(smtpC.Host, port, smtpC.Username, smtpC.Password)
	if err := dialer.DialAndSend(message); err != nil {
		log.NewEntry().Errorf("send email to %v error: %v", emails, err)
		return err
	}

	return nil
}
