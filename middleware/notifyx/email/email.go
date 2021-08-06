package email

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

// MessageType 内容
type MessageType struct {
	From        string   // 发件人
	To          []string // 收件人
	Cc          []string // 抄送人
	Subject     string   // 主题
	Attach      string   // 附件
	ContentType string   // 内容的类型 text/plain text/html
	Content     string   // 内容
}

// ClientType 发送客户端
type ClientType struct {
	Host     string       // smtp地址
	Port     int          // 端口
	Username string       // 用户名
	Password string       // 密码
	Message  *MessageType // 消息
}

// NewClient 返回一个邮件客户端
// host smtp地址
// port 端口
// username 用户名
// password 密码
// message 消息
func NewClient(host string, port int, username string, password string, message *MessageType) *ClientType {
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
