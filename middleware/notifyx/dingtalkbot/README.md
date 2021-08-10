# 自定义机器人接入

## 场景介绍

企业内部有较多系统支撑着公司的核心业务流程，譬如CRM系统、交易系统、监控报警系统等等。通过钉钉的自定义机器人，可以将这些系统事件同步到钉钉的聊天群。

> **说明**
>
> 当前机器人尚不支持应答机制，该机制指的是群里成员在聊天@机器人的时候，钉钉回调指定的服务地址，即Outgoing机器人。

## 步骤一：获取自定义机器人Webhook

- 1、选择需要添加机器人的群聊，然后依次单击 **群设置** > **智能群助手**。

-

- 2、在机器人管理页面选择 **自定义** 机器人，输入机器人名字并选择要发送消息的群，同时可以为机器人设置机器人头像。

- 3、完成必要的 [安全设置](https://developers.dingtalk.com/document/robots/customize-robot-security-settings#topic-2101465) (https://developers.dingtalk.com/document/robots/customize-robot-security-settings?spm=ding_open_doc.document.0.0.40745e59XVZzIr#topic-2101465)，勾选 **我已阅读并同意《自定义机器人服务及免责条款**》，然后单击**完成**。

- 4、完成安全设置后，复制出机器人的 **Webhook** 地址，可用于向这个群发送消息，格式如下：

  ```bash
  https://oapi.dingtalk.com/robot/send?access_token=XXXXXX
  ```

  > **注意**
  >
  > 请保管好此Webhook 地址，不要公布在外部网站上，泄露后有安全风险。

## 步骤二：使用自定义机器人

获取到 **Webhook** 地址后，用户可以向这个地址发起 **HTTP POST** 请求，即可实现给该钉钉群发送消息。

> **注意**
>
> - 发起POST请求时，必须将字符集编码设置成UTF-8。
> - 每个机器人每分钟最多发送20条。消息发送太频繁会严重影响群成员的使用体验，大量发消息的场景 (譬如系统监控报警) 可以将这些信息进行整合，通过markdown消息以摘要的形式发送到群里。

当前自定义机器人支持**文本 (text)、链接 (link)、markdown(****markdown)、ActionCard、FeedCard**消息类型，请根据自己的使用场景选择合适的消息类型，达到最好的展示样式。详情参考：[消息类型及数据格式](https://developers.dingtalk.com/document/robots/custom-robot-access#section-e4x-4y8-9k0) (https://developers.dingtalk.com/document/robots/custom-robot-access#section-e4x-4y8-9k0)。

自定义机器人发送消息时，可以通过手机号码指定“被@人列表”。在“被@人列表”里面的人员收到该消息时，会有@消息提醒。免打扰会话仍然通知提醒，首屏出现“有人@你”。

## 步骤三：测试自定义机器人

通过以下方法，可以快速验证自定义机器人是否可以正常工作：

使用命令行工具curl。

> **说明**
>
> 为避免出错，将以下命令逐行复制到命令行，需要将xxxxxxxx替换为真实access_token；若测试出错，请检查复制的命令是否和测试命令一致，多特殊字符会报错。

```bash
curl 'https://oapi.dingtalk.com/robot/send?access_token=xxxxxxxx' \
 -H 'Content-Type: application/json' \
 -d '{"msgtype": "text","text": {"content":"我就是我, 是不一样的烟火"}}'
```

## 常见问题

当出现以下错误时，表示消息校验未通过，请查看机器人的安全设置。

```bash
// 消息内容中不包含任何关键词
{
  "errcode":310000,
  "errmsg":"keywords not in content"
}

// timestamp 无效
{
  "errcode":310000,
  "errmsg":"invalid timestamp"
}

// 签名不匹配
{
  "errcode":310000,
  "errmsg":"sign not match"
}

// IP地址不在白名单
{
  "errcode":310000,
  "errmsg":"ip X.X.X.X not in whitelist"
}
```
