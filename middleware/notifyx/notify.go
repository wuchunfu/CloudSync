package notifyx

import (
	"github.com/wuchunfu/CloudSync/middleware/configx"
	"github.com/wuchunfu/CloudSync/middleware/notifyx/dingtalkbot"
	"github.com/wuchunfu/CloudSync/middleware/notifyx/email"
	"github.com/wuchunfu/CloudSync/middleware/notifyx/wechat"
	"github.com/wuchunfu/CloudSync/middleware/notifyx/wechatbot"
	"log"
)

func SendMessage(setting *configx.Notify, msg string) {
	if setting.IsEnable {
		// 邮件发送
		if setting.NotifyType == "email" {
			subject := setting.Email.Subject
			host := setting.Email.Host
			port := setting.Email.Port
			form := setting.Email.Form
			password := setting.Email.Password
			to := setting.Email.To

			message := &email.MessageType{
				From:        form,
				Subject:     subject,
				ContentType: "text/html",
				Content:     msg,
				To:          to,
			}
			client := email.NewClient(host, port, form, password, message)
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
		// 企业微信发送
		if setting.NotifyType == "wechat" {
			accessTokenUrl := setting.Wechat.AccessTokenUrl
			sendUrl := setting.Wechat.SendUrl
			corpId := setting.Wechat.CorpId
			corpSecret := setting.Wechat.CorpSecret
			toUser := setting.Wechat.ToUser
			agentId := setting.Wechat.AgentId

			message := &wechat.MessageType{
				MsgType: wechat.TEXT, // 消息类型，支持: text,image,voice,video,file,textcard,news,mpnews,markdown
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
		// 企业微信机器人发送
		if setting.NotifyType == "wechatbot" {
			webHookUrl := setting.WechatBot.WebHookUrl
			key := setting.WechatBot.Key
			mentionedList := setting.WechatBot.MentionedList
			mentionedMobileList := setting.WechatBot.MentionedMobileList

			message := &wechatbot.MessageType{
				MsgType: wechatbot.TEXT, // 消息类型, 支持: text,markdown,image,news,file
				Text: &wechatbot.TextType{
					Content:             msg,
					MentionedList:       mentionedList,
					MentionedMobileList: mentionedMobileList,
				},
			}
			client := wechatbot.NewClient(webHookUrl, key, message)
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
		// 钉钉机器人发送
		if setting.NotifyType == "dingtalkbot" {
			webHookUrl := setting.DingTalkBot.WebHookUrl
			accessToken := setting.DingTalkBot.AccessToken
			signToken := setting.DingTalkBot.SignToken
			atUserIds := setting.DingTalkBot.AtUserIds
			atMobiles := setting.DingTalkBot.AtMobiles
			isAtAll := setting.DingTalkBot.IsAtAll

			message := &dingtalkbot.MessageType{
				MsgType: dingtalkbot.TEXT, // 消息类型, 支持: text,markdown,image,news,file
				Text: &dingtalkbot.TextType{
					Content: msg,
				},
				At: &dingtalkbot.AtType{
					AtUserIds: atUserIds,
					AtMobiles: atMobiles,
					IsAtAll:   isAtAll,
				},
			}
			client := dingtalkbot.NewClient(webHookUrl, accessToken, signToken, message)
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
