package dingtalkbot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strconv"
	"time"
)

// 自定义机器人接入 文档: https://developers.dingtalk.com/document/robots/custom-robot-access
// 自定义机器人安全设置 文档: https://developers.dingtalk.com/document/robots/customize-robot-security-settings
const (
	// TEXT 文本消息
	TEXT = "text"

	// LINK link消息
	LINK = "link"

	// MARKDOWN markdown消息
	MARKDOWN = "markdown"

	// ACTIONCARD actionCard消息
	ACTIONCARD = "actionCard"

	// FEEDCARD feedCard消息
	FEEDCARD = "feedCard"
)

// AtType @ 类型
type AtType struct {
	AtUserIds []string `json:"atUserIds"` // 是否必须: 否 , 被@人的用户userid。
	AtMobiles []string `json:"atMobiles"` // 是否必须: 否 , 被@人的手机号。说明 消息内容content中要带上"@手机号"，跟atMobiles参数结合使用，才有@效果，如上示例。
	IsAtAll   bool     `json:"isAtAll"`   // 是否必须: 否 , 是否@所有人, @所有人是true，否则为false。
}

// TextType 文本类型
type TextType struct {
	Content string `json:"content"` // 是否必须: 是 , 消息内容。
}

// LinkType link类型
type LinkType struct {
	Title      string `json:"title"`      // 是否必须: 是 , 消息标题。
	Text       string `json:"text"`       // 是否必须: 是 , 消息内容。如果太长只会部分展示。
	PicUrl     string `json:"picUrl"`     // 是否必须: 是 , 点击消息跳转的URL。
	MessageUrl string `json:"messageUrl"` // 是否必须: 是 , 图片URL。
}

// MarkdownType markdown类型
type MarkdownType struct {
	Title string `json:"title"` // 是否必须: 是 , 首屏会话透出的展示内容。
	Text  string `json:"text"`  // 是否必须: 是 , markdown格式的消息内容。
}

// ActionCardAllType 整体跳转actionCard类型
type ActionCardAllType struct {
	Title          string `json:"title"`          // 是否必须: 是 , 首屏会话透出的展示内容。
	Text           string `json:"text"`           // 是否必须: 是 , markdown格式的消息内容。
	SingleTitle    string `json:"singleTitle"`    // 是否必须: 是 , 单个按钮的标题。
	SingleURL      string `json:"singleURL"`      // 是否必须: 是 , 单个按钮的跳转链接。
	BtnOrientation string `json:"btnOrientation"` // 是否必须: 是 , 按钮排列顺序。0：按钮竖直排列 1：按钮横向排列
}

type Btns struct {
	Title     string `json:"title"`     // 是否必须: 是 , 按钮标题。
	ActionURL string `json:"actionURL"` // 是否必须: 是 , 点击按钮触发的URL。
}

// ActionCardType 独立跳转actionCard类型
type ActionCardType struct {
	Title          string  `json:"title"`          // 是否必须: 是 , 首屏会话透出的展示内容。
	Text           string  `json:"text"`           // 是否必须: 是 , markdown格式的消息内容。
	BtnOrientation string  `json:"btnOrientation"` // 是否必须: 否 , 按钮排列顺序。0：按钮竖直排列 1：按钮横向排列
	Btns           []*Btns `json:"btns"`           // 是否必须: 否 , 按钮。
}

type Links struct {
	Title      string `json:"title"`      // 是否必须: 是 , 单条信息文本。
	MessageURL string `json:"messageURL"` // 是否必须: 是 , 点击单条信息到跳转链接。
	PicURL     string `json:"picURL"`     // 是否必须: 是 , 单条信息后面图片的URL。
}

// FeedCardType feedCard类型
type FeedCardType struct {
	Links []*Links `json:"links"` // 是否必须: 是 , 连接。
}

// MessageType 微信消息
type MessageType struct {
	MsgType       string             `json:"msgtype"`    // 是否必须: 是 , 消息类型, 支持: text,link,markdown,actionCard,feedCard
	At            *AtType            `json:"at"`         // @ 类型
	Text          *TextType          `json:"text"`       // 文本类型
	Link          *LinkType          `json:"link"`       // link类型
	Markdown      *MarkdownType      `json:"markdown"`   // markdown类型
	ActionCardAll *ActionCardAllType `json:"actionCard"` // 整体跳转ActionCard类型
	ActionCard    *ActionCardType    `json:"actionCard"` // 独立跳转ActionCard类型
	FeedCard      *FeedCardType      `json:"feedCard"`   // FeedCard类型
}

// ErrorType 发送消息返回错误
type ErrorType struct {
	ErrCode int    `json:"errcode"` // 出错返回码，为0表示成功，非0表示调用失败
	ErrMsg  string `json:"errmsg"`  // 返回码提示语
}

// ClientType 客户端
type ClientType struct {
	WebHookUrl  string // https://oapi.dingtalk.com/robot/send
	AccessToken string // bf29f17ef2972180bacad9adf19412f7728e5b336fcb4c152af5be8a88888888
	SignToken   string // SEC3383b31ef94081d10e5ae8c009923d738c4b92539746b40c3aec2c8e88888888
	Message     *MessageType
}

func NewClient(webHookUrl string, accessToken string, signToken string, message *MessageType) *ClientType {
	return &ClientType{
		WebHookUrl:  webHookUrl,
		AccessToken: accessToken,
		SignToken:   signToken,
		Message:     message,
	}
}

func Sign(timestamp int64, secret string) string {
	signContent := fmt.Sprintf("%d\n%s", timestamp, secret)
	signByte := []byte(signContent)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(signByte)
	signData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(signData)
}

func GetSign(secret string) (int64, string) {
	// 该值要求毫秒
	timestamp := time.Now().UnixNano() / 1000 / 1000
	sign := Sign(timestamp, secret)
	return timestamp, sign
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

	timestamp, signToken := GetSign(client.SignToken)

	httpClient := resty.New()
	headers := map[string]string{
		"Content-Type": "application/json;charset=utf-8",
	}
	params := map[string]string{
		"access_token": client.AccessToken,
		"sign":         signToken,
		"timestamp":    strconv.FormatInt(timestamp, 10),
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
