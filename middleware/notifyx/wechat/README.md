# 简介
企业微信(微信) 开发库, 目前只支持发送消息,微信通过企业微信中的应用来接收消息,可以用来做报警，要接收报警的微信必须在企业微信的通讯录里面。

> 注意: 微信企业号现已升级为企业微信，微信插件继承原企业号的所有能力。管理员可在企业微信管理后台新建应用、群发通知，成员无需下载客户端，扫码关注微信插件后即可在微信中接收企业通知和使用企业应用。

# 使用微信报警流程
1. 到[**这里**](https://work.weixin.qq.com/wework_admin/register_wx?from=loginpage)注册企业微信, 使用管理员微信账号扫码登录企业微信管理后台
2. 邀请接收报警的人加入企业，首页里面有 邀请方式，可以通过微信扫码的方式加入，确保都在通讯录里面才能接收报警
3. 企业应用 - 创建应用 - 上传一个Logo,填写应用名称，选择部门/成员 这些人就是通过这个应用接收报警
4. 最后让所有接收报警的人，扫描 企业应用 - 微信插件 的二维码，即可接收报警。

# corpId 和 appSecret 参数获取

[corpId 和 appSecret 参数获取](https://work.weixin.qq.com/api/doc/90000/90135/90665)

# 用法

### 发送消息

[点击查看详细的说明](https://work.weixin.qq.com/api/doc#10167)

```go
package wechat

import (
	"testing"
)

func TestWechat(t *testing.T) {
	accessTokenUrl := "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
	sendUrl := "https://qyapi.weixin.qq.com/cgi-bin/message/send"
	corpId := "xxx"
	corpSecret := "xxxxxx"

	message := &MessageType{
		MsgType: TEXT,       // 消息类型
		ToUser:  "xxx", // 指定接收消息的成员，成员ID列表（多个接收者用‘|’分隔，最多支持1000个）。特殊情况：指定为”@all”，则向该企业应用的全部成员发送
		AgentId: 1000000,    // 企业应用的id，整型。企业内部开发，可在应用的设置页面查看；第三方服务商，可通过接口 获取企业授权信息(https://work.weixin.qq.com/api/doc/90001/90143/90372#10975/%E8%8E%B7%E5%8F%96%E4%BC%81%E4%B8%9A%E6%8E%88%E6%9D%83%E4%BF%A1%E6%81%AF) 获取该参数值
		Text: &TextType{
			Content: "test msg!",
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
```