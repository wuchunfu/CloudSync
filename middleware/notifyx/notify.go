package notifyx

import (
	"fmt"
	"github.com/wuchunfu/CloudSync/middleware/configx"
	"github.com/wuchunfu/CloudSync/middleware/notifyx/email"
)

func SendMessage(setting *configx.Notify, msg string) {
	if setting.IsEnable {
		// 邮件发送
		if setting.NotifyType == "email" {
			newEmail := email.NewEmail(setting)
			err := newEmail.Send(msg)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
