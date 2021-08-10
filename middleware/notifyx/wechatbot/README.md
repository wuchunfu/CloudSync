# 群机器人配置说明

## 如何使用群机器人

- 在终端某个群组添加机器人之后，创建者可以在机器人详情页看的该机器人特有的 **webhookurl**。开发者可以按以下说明a向这个地址发起HTTP POST 请求，即可实现给该群组发送消息。下面举个简单的例子.
  假设 **webhook** 是：

  ```bash
  https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=693a91f6-7xxx-4bc4-97a0-0ec2sifa5aaa
  ```

  > 特别特别要注意：一定要 **保护好机器人的webhook地址** ，避免泄漏！不要分享到github、博客等可被公开查阅的地方，否则坏人就可以用你的机器人来发垃圾消息了。

以下是用curl工具往群组推送文本消息的示例（注意要将url替换成你的机器人webhook地址，content必须是utf8编码）：

```bash
curl 'https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=693axxx6-7aoc-4bc4-97a0-0ec2sifa5aaa' \
   -H 'Content-Type: application/json' \
   -d '
   {
        "msgtype": "text",
        "text": {
            "content": "hello world"
        }
   }'
```

- 当前自定义机器人支持文本（text）、markdown（markdown）、图片（image）、图文（news）四种消息类型。
- 机器人的text/markdown类型消息支持在content中使用<@userid>扩展语法来@群成员

## 请求示例

```go
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
```
