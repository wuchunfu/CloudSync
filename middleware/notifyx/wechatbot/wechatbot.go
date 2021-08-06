package wechatbot

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
)

// 企业微信 企业内部开发 客户端API 群机器人 文档: https://work.weixin.qq.com/api/doc/90000/90136/91770
const (
	// TEXT 文本消息
	TEXT = "text"

	// IMAGE 图片消息
	IMAGE = "image"

	// VOICE 语音消息
	VOICE = "voice"

	// VIDEO 视频消息
	VIDEO = "video"

	// FILE 文件消息
	FILE = "file"

	// TEXTCARD 文本卡片消息
	TEXTCARD = "textcard"

	// NEWS 图文消息
	NEWS = "news"

	// MPNEWS 图文消息(mpnews)
	MPNEWS = "mpnews"

	// MARKDOWN markdown消息
	MARKDOWN = "markdown"
)

// TextType 文本类型
type TextType struct {
	Content             string   `json:"content"`               // 是否必须: 是 , 文本内容，最长不超过2048个字节，必须是utf8编码
	MentionedList       []string `json:"mentioned_list"`        // 是否必须: 否 , userid的列表，提醒群中的指定成员(@某个成员)，@all表示提醒所有人，如果开发者获取不到userid，可以使用mentioned_mobile_list
	MentionedMobileList []string `json:"mentioned_mobile_list"` // 是否必须: 否 , 手机号列表，提醒手机号对应的群成员(@某个成员)，@all表示提醒所有人
}

// MarkdownType markdown 类型
type MarkdownType struct {
	Content string `json:"content"` // 是否必须: 是 , markdown内容，最长不超过4096个字节，必须是utf8编码
}

// ImageType 图片类型
type ImageType struct {
	Base64 string `json:"base64"` // 是否必须: 是 , 图片内容的base64编码
	Md5    string `json:"md5"`    // 是否必须: 是 , 图片内容（base64编码前）的md5值
}

// NewType 单个图文消息
type NewType struct {
	Title       string `json:"title"`       // 是否必须: 是 , 标题，不超过128个字节，超过会自动截断
	Description string `json:"description"` // 是否必须: 否 , 描述，不超过512个字节，超过会自动截断
	Url         string `json:"url"`         // 是否必须: 是 , 点击后跳转的链接。
	PicUrl      string `json:"picurl"`      // 是否必须: 否 , 图文消息的图片链接，支持JPG、PNG格式，较好的效果为大图 1068*455，小图150*150。
}

// NewsType 图文类型
type NewsType struct {
	Articles []NewType `json:"articles"` // 是否必须: 是 , 图文消息，一个图文消息支持1到8条图文
}

// FileType 文件类型
type FileType struct {
	MediaId []NewType `json:"media_id"` // 是否必须: 是 , 文件id，通过下文的文件上传接口获取
}

// MessageType 微信消息
type MessageType struct {
	MsgType  string        `json:"msgtype"`  // 是否必须: 是 , 消息类型, 支持: text,markdown,image,news,file
	Text     *TextType     `json:"text"`     // 文本类型
	Markdown *MarkdownType `json:"markdown"` // markdown 类型
	Image    *ImageType    `json:"image"`    // 图片类型
	News     *NewsType     `json:"news"`     // 图文类型
	File     *FileType     `json:"file"`     // 文件类型
}

// ErrorType 发送消息返回错误
type ErrorType struct {
	ErrCode int    `json:"errcode"` // 出错返回码，为0表示成功，非0表示调用失败
	ErrMsg  string `json:"errmsg"`  // 返回码提示语
}

// ClientType 客户端
type ClientType struct {
	WebHookUrl string // https://qyapi.weixin.qq.com/cgi-bin/webhook/send
	Key        string // 6b7372c1-53d0-47a9-941f-6adcb888888a
	Message    *MessageType
}

func NewClient(webHookUrl string, key string, message *MessageType) *ClientType {
	return &ClientType{
		WebHookUrl: webHookUrl,
		Key:        key,
		Message:    message,
	}
}

// SendMessage 发送消息
func (client *ClientType) SendMessage() (bool, error) {
	msg, err := json.Marshal(client.Message)
	if err != nil {
		return false, fmt.Errorf("parse message fail: %v", err)
	}

	if client.WebHookUrl == "" {
		return false, fmt.Errorf("miss webhook url")
	}

	httpClient := resty.New()
	headers := map[string]string{
		"Content-Type": "application/json;charset=utf-8",
	}
	params := map[string]string{
		"key": client.Key,
	}
	resp, err := httpClient.R().
		SetHeaders(headers).
		SetBody(msg).
		SetQueryParams(params).
		Post(client.WebHookUrl)
	if err != nil {
		return false, fmt.Errorf("invoke send api fail: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return false, fmt.Errorf("invoke send api fail: %v", resp.Error())
	}

	result := ErrorType{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return false, fmt.Errorf("parse send api response fail: %v", err)
	}

	if result.ErrCode != 0 {
		return false, fmt.Errorf("invoke send api fail, error: %s", result.ErrMsg)
	} else {
		return true, nil
	}
}
