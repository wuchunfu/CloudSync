package wechatbot

import "testing"

// 企业微信 企业内部开发 客户端API 群机器人 文档: https://work.weixin.qq.com/api/doc/90000/90136/91770
func TestWechatBot(t *testing.T) {
	webHookUrl := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send"
	key := "xxxxxx"
	mentionedList := []string{"xxx"}
	mentionedMobileList := []string{""}
	msg := "test msg!"

	message := &MessageType{
		MsgType: TEXT, // 消息类型, 支持: text,markdown,image,news,file
		Text: &TextType{
			Content:             msg,
			MentionedList:       mentionedList,
			MentionedMobileList: mentionedMobileList,
		},
	}
	client := NewClient(webHookUrl, key, message)
	ok, err := client.SendMessage()
	if ok {
		t.Log("send successfully!")
	} else {
		t.Fatalf("send faild, error:%v", err)
	}
}
