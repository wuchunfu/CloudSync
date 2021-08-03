package notifyx

import (
	"github.com/wuchunfu/CloudSync/middleware/configx"
	"github.com/wuchunfu/CloudSync/middleware/notifyx/email"
	"log"
)

func SendMessage(setting *configx.Notify, msg string) {
	if setting.IsEnable {
		// 邮件发送
		if setting.NotifyType == "email" {
			subject := setting.Email.EmailSubject
			host := setting.Email.EmailHost
			port := setting.Email.EmailPort
			user := setting.Email.EmailUser
			pwd := setting.Email.EmailPwd
			to := setting.Email.ToEmail
			message := email.NewEmailMessage(user, subject, "text/html", msg, "", []string{to}, []string{})
			client := email.NewEmailClient(host, port, user, pwd, message)
			ok, err := client.SendMessage()
			if err != nil {
				log.Printf("send failed, error: %v", err)
			}
			if ok {
				log.Println("send successfully")
			} else {
				log.Println("send failed!")
			}
		}
	}
}
