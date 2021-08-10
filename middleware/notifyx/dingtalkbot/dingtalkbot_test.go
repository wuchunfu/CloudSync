package dingtalkbot

import (
	"testing"
)

// https://help.aliyun.com/document_detail/121918.html
// 自定义机器人接入 文档: https://developers.dingtalk.com/document/robots/custom-robot-access
// 自定义机器人安全设置 文档: https://developers.dingtalk.com/document/robots/customize-robot-security-settings
func TestDingTalkBot(t *testing.T) {
	webHookUrl := "https://oapi.dingtalk.com/robot/send"
	accessToken := "bf29f17ef2972180bacad9adf19412f7728e5b336fcb4c152af5be8a88888888"
	signToken := "SEC3383b31ef94081d10e5ae8c009923d738c4b92539746b40c3aec2c8e88888888"
	atUserIds := []string{""}
	atMobiles := []string{""}
	isAtAll := false
	msg := "test msg!"

	message := &MessageType{
		MsgType: TEXT,
		Text: &TextType{
			Content: msg,
		},
		At: &AtType{
			AtUserIds: atUserIds,
			AtMobiles: atMobiles,
			IsAtAll:   isAtAll,
		},
	}
	client := NewClient(webHookUrl, accessToken, signToken, message)
	ok, err := client.SendMessage()
	if ok {
		t.Log("send successfully!")
	} else {
		t.Fatalf("send faild, error:%v", err)
	}
}
