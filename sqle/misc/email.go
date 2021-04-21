package misc

import (
	"fmt"
	"strconv"
	"strings"

	"actiontech.cloud/universe/sqle/v4/sqle/log"
	"actiontech.cloud/universe/sqle/v4/sqle/model"

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
	var emails []string
	for _, user := range workflow.CurrentStep().Template.Users {
		emails = append(emails, user.Email)
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
`, model.GetWorkflowStepTypeDesc(workflow.CurrentStep().Template.Typ), workflow.Subject, workflow.Desc, workflow.CreateUser.Name)
	message.SetBody("text/html",
		strings.Replace(body, "\n", "<br/>\n", -1))

	port, _ := strconv.Atoi(smtpC.Port)
	dialer := gomail.NewDialer(smtpC.Host, port, smtpC.Username, smtpC.Password)
	if err := dialer.DialAndSend(message); err != nil {
		log.NewEntry().Errorf("send emial to %v error: %v", emails, err)
		return err
	}

	return nil
}
