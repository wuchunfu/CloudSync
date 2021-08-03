package email

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

// MessageType 内容
type MessageType struct {
	From        string
	To          []string
	Cc          []string
	Subject     string
	ContentType string
	Content     string
	Attach      string
}

// ClientType 发送客户端
type ClientType struct {
	Host     string
	Port     int
	Username string
	Password string
	Message  *MessageType
}

// NewEmailMessage 返回消息对象
// from: 发件人
// subject: 标题
// contentType: 内容的类型 text/plain text/html
// attach: 附件
// to: 收件人
// cc: 抄送人
func NewEmailMessage(from string, subject string, contentType string, content string, attach string, to []string, cc []string) *MessageType {
	return &MessageType{
		From:        from,
		Subject:     subject,
		ContentType: contentType,
		Content:     content,
		To:          to,
		Cc:          cc,
		Attach:      attach,
	}
}

// NewEmailClient 返回一个邮件客户端
// host smtp地址
// username 用户名
// password 密码
// port 端口
func NewEmailClient(host string, port int, username string, password string, message *MessageType) *ClientType {
	return &ClientType{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Message:  message,
	}
}

// SendMessage 发送邮件
func (client *ClientType) SendMessage() (bool, error) {
	dialer := gomail.NewDialer(client.Host, client.Port, client.Username, client.Password)
	if 587 == client.Port || 465 == client.Port || 994 == client.Port {
		dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	dm := gomail.NewMessage()
	dm.SetHeader("From", client.Message.From)
	dm.SetHeader("To", client.Message.To...)

	if len(client.Message.Cc) != 0 {
		dm.SetHeader("Cc", client.Message.Cc...)
	}

	dm.SetHeader("Subject", client.Message.Subject)
	dm.SetBody(client.Message.ContentType, client.Message.Content)

	if client.Message.Attach != "" {
		dm.Attach(client.Message.Attach)
	}

	if err := dialer.DialAndSend(dm); err != nil {
		return false, err
	}
	return true, nil
}
