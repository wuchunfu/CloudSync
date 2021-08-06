package wechat

import (
	"testing"
)

// 企业微信 企业内部开发 服务端API 消息推送 文档: https://work.weixin.qq.com/api/doc/90000/90135/90235
// 企业微信 企业内部开发 服务端API 开发指南 文档: https://work.weixin.qq.com/api/doc/90000/90135/90664
// 企业微信 企业内部开发 客户端API 群机器人 文档: https://work.weixin.qq.com/api/doc/90000/90136/91770
// corpId 和 appSecret 参数获取: https://work.weixin.qq.com/api/doc/90000/90135/90665
func TestWechat(t *testing.T) {
	accessTokenUrl := "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
	sendUrl := "https://qyapi.weixin.qq.com/cgi-bin/message/send"
	corpId := "xxxxxx"
	corpSecret := "xxxxxx"
	toUser := "xxx"
	agentId := 1000000
	msg := "test msg!"

	message := &MessageType{
		MsgType: TEXT,    // 消息类型，支持: text,image,voice,video,file,textcard,news,mpnews,markdown
		ToUser:  toUser,  // 指定接收消息的成员，成员ID列表（多个接收者用‘|’分隔，最多支持1000个）。特殊情况：指定为”@all”，则向该企业应用的全部成员发送
		AgentId: agentId, // 企业应用的id，整型。企业内部开发，可在应用的设置页面查看；第三方服务商，可通过接口 获取企业授权信息(https://work.weixin.qq.com/api/doc/90001/90143/90372#10975/%E8%8E%B7%E5%8F%96%E4%BC%81%E4%B8%9A%E6%8E%88%E6%9D%83%E4%BF%A1%E6%81%AF) 获取该参数值
		Text: &TextType{
			Content: msg,
		},
	}
	client := NewClient(accessTokenUrl, sendUrl, corpId, corpSecret, message)
	ok, err := client.SendMessage()
	if ok {
		t.Log("send successfully!")
	} else {
		t.Fatalf("send faild, error:%v", err)
	}
}
