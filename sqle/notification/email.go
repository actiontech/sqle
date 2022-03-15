package notification

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/actiontech/sqle/sqle/model"

	"gopkg.in/gomail.v2"
)

type EmailNotifier struct{}

func (n *EmailNotifier) Notify(notification Notification, users []*model.User) error {
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

	s := model.GetStorage()
	smtpC, exist, err := s.GetSMTPConfiguration()
	if err != nil {
		return err
	}
	if !exist {
		return nil
	}

	message := gomail.NewMessage()
	message.SetHeader("From", smtpC.Username)
	message.SetHeader("To", emails...)
	message.SetHeader("Subject", notification.NotificationSubject())
	body := notification.NotificationBody()
	message.SetBody("text/html", strings.Replace(body, "\n", "<br/>\n", -1))

	port, _ := strconv.Atoi(smtpC.Port)
	dialer := gomail.NewDialer(smtpC.Host, port, smtpC.Username, smtpC.Password)
	if err := dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("send email to %v error: %v", emails, err)
	}
	return nil
}
