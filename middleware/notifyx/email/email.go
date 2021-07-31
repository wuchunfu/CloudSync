package email

import (
	"fmt"
	"github.com/wuchunfu/CloudSync/middleware/configx"
	"gopkg.in/gomail.v2"
)

type Email struct {
	EmailSubject string
	EmailHost    string
	EmailPort    int
	EmailUser    string
	EmailPwd     string
	ToEmail      string
}

func NewEmail(setting *configx.Notify) *Email {
	subject := setting.Email.EmailSubject
	host := setting.Email.EmailHost
	port := setting.Email.EmailPort
	user := setting.Email.EmailUser
	pwd := setting.Email.EmailPwd
	to := setting.Email.ToEmail
	return &Email{EmailSubject: subject, EmailHost: host, EmailPort: port, EmailUser: user, EmailPwd: pwd, ToEmail: to}
}

func (email *Email) Send(body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "<"+email.EmailUser+">")
	m.SetHeader("To", []string{email.ToEmail}...)
	m.SetHeader("Subject", email.EmailSubject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer(email.EmailHost, email.EmailPort, email.EmailUser, email.EmailPwd)
	fmt.Println("正在发送通知...")
	err := d.DialAndSend(m)
	if err != nil {
		fmt.Errorf("邮件发送失败，返回错误: %s", err.Error())
	} else {
		fmt.Println("邮件发送成功")
	}
	return err
}
