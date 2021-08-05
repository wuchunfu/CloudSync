package notifyx

import (
	"github.com/wuchunfu/CloudSync/middleware/configx"
	"github.com/wuchunfu/CloudSync/middleware/notifyx/email"
	"github.com/wuchunfu/CloudSync/middleware/notifyx/wechat"
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

			message := &email.MessageType{
				From:        user,
				Subject:     subject,
				ContentType: "text/html",
				Content:     msg,
				To:          []string{to},
			}
			client := email.NewClient(host, port, user, pwd, message)
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
		if setting.NotifyType == "wechat" {
			accessTokenUrl := setting.Wechat.AccessTokenUrl
			sendUrl := setting.Wechat.SendUrl
			corpId := setting.Wechat.CorpId
			corpSecret := setting.Wechat.CorpSecret
			toUser := setting.Wechat.ToUser
			agentId := setting.Wechat.AgentId

			message := &wechat.MessageType{
				MsgType: wechat.TEXT, // 消息类型
				ToUser:  toUser,      // 指定接收消息的成员，成员ID列表（多个接收者用‘|’分隔，最多支持1000个）。特殊情况：指定为”@all”，则向该企业应用的全部成员发送
				AgentId: agentId,     // 企业应用的id，整型。企业内部开发，可在应用的设置页面查看；第三方服务商，可通过接口 获取企业授权信息(https://work.weixin.qq.com/api/doc/90001/90143/90372#10975/%E8%8E%B7%E5%8F%96%E4%BC%81%E4%B8%9A%E6%8E%88%E6%9D%83%E4%BF%A1%E6%81%AF) 获取该参数值
				Text: &wechat.TextType{
					Content: msg,
				},
			}
			client := wechat.NewClient(accessTokenUrl, sendUrl, corpId, corpSecret, message)
			ok, err := client.SendMessage()
			if err != nil {
				log.Printf("send failed:%v", err)
			}
			if ok {
				log.Println("send successfully!")
			} else {
				log.Println("send failed!")
			}
		}
	}
}
