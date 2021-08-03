package email

import (
	"testing"
)

func TestEmail(t *testing.T) {
	subject := "email test"
	host := "smtp.163.com"
	port := 25
	user := "xxx@163.com"
	pwd := "xxxxxx"
	to := "xxx@qq.com"
	msg := "test msg!"
	message := NewEmailMessage(user, subject, "text/html", msg, "", []string{to}, []string{})
	client := NewEmailClient(host, port, user, pwd, message)
	ok, err := client.SendMessage()
	if err != nil {
		t.Errorf("send failed, error: %v", err)
	}
	if ok {
		t.Log("send successfully")
	} else {
		t.Log("send failed!")
	}
}
