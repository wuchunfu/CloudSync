package email

import (
	"testing"
)

func TestEmail(t *testing.T) {
	subject := "email test"
	host := "smtp.163.com"
	port := 25
	from := "xxx@163.com"
	pwd := "xxxxxx"
	to := []string{"xxx@qq.com"}
	cc := []string{""}
	attach := ""
	msg := "test msg!"

	message := &MessageType{
		From:        from,
		Subject:     subject,
		To:          to,
		Cc:          cc,
		Attach:      attach,
		ContentType: "text/html",
		Content:     msg,
	}
	client := NewClient(host, port, from, pwd, message)
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
