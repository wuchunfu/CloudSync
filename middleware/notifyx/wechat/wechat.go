package wechat

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/wuchunfu/CloudSync/middleware/database/kvdb"
	"log"
	"sync"
	"time"
)

// 企业微信 企业内部开发 服务端API 消息推送 文档: https://work.weixin.qq.com/api/doc/90000/90135/90235
// 企业微信 企业内部开发 服务端API 开发指南 文档: https://work.weixin.qq.com/api/doc/90000/90135/90664
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

var (
	Cache = kvdb.NewKvDb("./token.db", "token")
)

// TextType 文本消息
// touser、toparty、totag不能同时为空，后面不再强调。
type TextType struct {
	Content string `json:"content"` // 是否必须: 是 , 消息内容，最长不超过2048个字节，超过将截断（支持id转译）, 其中text参数的content字段可以支持换行、以及A标签，即可打开自定义的网页（可参考以上示例代码）(注意：换行符请用转义过的\n)
}

// ImageType 图片消息
type ImageType struct {
	MediaId string `json:"media_id"` // 是否必须: 是 , 图片媒体文件id，可以调用上传临时素材接口获取
}

// VoiceType 语音消息
type VoiceType struct {
	MediaId string `json:"media_id"` // 是否必须: 是 , 语音文件id，可以调用 上传临时素材(https://work.weixin.qq.com/api/doc/90001/90143/90372#10112) 接口获取
}

// VideoType 视频消息
type VideoType struct {
	MediaId     string `json:"media_id"`    // 是否必须: 是 , 视频媒体文件id，可以调用 上传临时素材(https://work.weixin.qq.com/api/doc/90001/90143/90372#10112) 接口获取
	Title       string `json:"title"`       // 是否必须: 否 , 视频消息的标题，不超过128个字节，超过会自动截断
	Description string `json:"description"` // 是否必须: 否 , 视频消息的描述，不超过512个字节，超过会自动截断
}

// FileType 文件消息
type FileType struct {
	MediaId string `json:"media_id"` // 是否必须: 是 , 文件id，可以调用上传临时素材接口获取
}

// TextCardType 文本卡片消息
type TextCardType struct {
	Title       string `json:"title"`       // 是否必须: 是 , 标题，不超过128个字节，超过会自动截断（支持id转译）
	Description string `json:"description"` // 是否必须: 是 , 描述，不超过512个字节，超过会自动截断（支持id转译）
	Url         string `json:"url"`         // 是否必须: 是 , 点击后跳转的链接。最长2048字节，请确保包含了协议头(http/https)
	BtnTxt      string `json:"btntxt"`      // 是否必须: 否 , 按钮文字。 默认为“详情”， 不超过4个文字，超过自动截断。
}

// NewType 单个图文消息
type NewType struct {
	Title       string `json:"title"`       // 是否必须: 是 , 标题，不超过128个字节，超过会自动截断（支持id转译）
	Description string `json:"description"` // 是否必须: 否 , 描述，不超过512个字节，超过会自动截断（支持id转译）
	Url         string `json:"url"`         // 是否必须: 是 , 点击后跳转的链接。 最长2048字节，请确保包含了协议头(http/https)
	PicUrl      string `json:"picurl"`      // 是否必须: 否 , 图文消息的图片链接，支持JPG、PNG格式，较好的效果为大图 1068*455，小图150*150。
}

// NewsType 图文消息
type NewsType struct {
	Articles []NewType `json:"articles"` // 是否必须: 是 , 图文消息，一个图文消息支持1到8条图文
}

// MpNewType 图文消息（mpnews）,
// mpnews类型的图文消息，跟普通的图文消息一致，唯一的差异是图文内容存储在企业微信。
// 多次发送mpnews，会被认为是不同的图文，阅读、点赞的统计会被分开计算。
type MpNewType struct {
	Title            string `json:"title"`              // 是否必须: 是 , 标题，不超过128个字节，超过会自动截断（支持id转译）
	ThumbMediaId     string `json:"thumb_media_id"`     // 是否必须: 是 , 图文消息缩略图的media_id, 可以通过 素材管理(https://work.weixin.qq.com/api/doc/90001/90143/90372#10112) 接口获得。此处thumb_media_id即上传接口返回的media_id
	Author           string `json:"author"`             // 是否必须: 否 , 图文消息的作者，不超过64个字节
	ContentSourceUrl string `json:"content_source_url"` // 是否必须: 否 , 图文消息点击“阅读原文”之后的页面链接
	Content          string `json:"content"`            // 是否必须: 是 , 图文消息的内容，支持html标签，不超过666 K个字节（支持id转译）
	Digest           string `json:"digest"`             // 是否必须: 否 , 图文消息的描述，不超过512个字节，超过会自动截断（支持id转译）
}

// MpNewsType 图文消息（mpnews）,
// mpnews类型的图文消息，跟普通的图文消息一致，唯一的差异是图文内容存储在企业微信。
// 多次发送mpnews，会被认为是不同的图文，阅读、点赞的统计会被分开计算。
type MpNewsType struct {
	Articles []MpNewType `json:"articles"` // 是否必须: 是 , 图文消息，一个图文消息支持1到8条图文
}

// MarkdownType markdown消息
// 目前仅支持 markdown语法的子集(https://work.weixin.qq.com/api/doc/90001/90143/90372#10167/%E6%94%AF%E6%8C%81%E7%9A%84markdown%E8%AF%AD%E6%B3%95)
// 微工作台（原企业号）不支持展示markdown消息
type MarkdownType struct {
	Content string `json:"content"` // 是否必须: 是 , markdown内容，最长不超过2048个字节，必须是utf8编码
}

// MessageType 微信消息
type MessageType struct {
	ToUser                 string        `json:"touser"`                   // 是否必须: 否 , 指定接收消息的成员，成员ID列表（多个接收者用'|'分隔，最多支持1000个）。特殊情况：指定为"@all"，则向该企业应用的全部成员发送
	ToParty                string        `json:"toparty"`                  // 是否必须: 否 , 指定接收消息的部门，部门ID列表，多个接收者用'|'分隔，最多支持100个。当touser为"@all"时忽略本参数
	ToTag                  string        `json:"totag"`                    // 是否必须: 否 , 指定接收消息的标签，标签ID列表，多个接收者用'|'分隔，最多支持100个。当touser为"@all"时忽略本参数
	MsgType                string        `json:"msgtype"`                  // 是否必须: 是 , 消息类型，text,image,voice,video,file,textcard,news,mpnews,markdown
	AgentId                int           `json:"agentid"`                  // 是否必须: 是 , 企业应用的id，整型。企业内部开发，可在应用的设置页面查看；第三方服务商，可通过接口 获取企业授权信息(https://work.weixin.qq.com/api/doc/90001/90143/90372#10975/%E8%8E%B7%E5%8F%96%E4%BC%81%E4%B8%9A%E6%8E%88%E6%9D%83%E4%BF%A1%E6%81%AF) 获取该参数值
	Safe                   int           `json:"safe"`                     // 是否必须: 否 , 表示是否是保密消息，0表示可对外分享，1表示不能分享且内容显示水印，2表示仅限在企业内分享，默认为0；注意仅mpnews类型的消息支持safe值为2，其他消息类型不支持
	EnableIdTrans          int           `json:"enable_id_trans"`          // 是否必须: 否 , 表示是否开启id转译，0表示否，1表示是，默认0。仅第三方应用需要用到，企业自建应用可以忽略。
	EnableDuplicateCheck   int           `json:"enable_duplicate_check"`   // 是否必须: 否 , 表示是否开启重复消息检查，0表示否，1表示是，默认0
	DuplicateCheckInterval int           `json:"duplicate_check_interval"` // 是否必须: 否 , 表示是否重复消息检查的时间间隔，默认1800s，最大不超过4小时
	Text                   *TextType     `json:"text"`                     // 文本消息
	Image                  *ImageType    `json:"image"`                    // 图片消息
	Voice                  *VoiceType    `json:"voice"`                    // 语音消息
	Video                  *VideoType    `json:"video"`                    // 视频消息
	File                   *FileType     `json:"file"`                     // 文件消息
	TextCard               *TextCardType `json:"textcard"`                 // 文本卡片消息
	News                   *NewsType     `json:"news"`                     // 图文消息
	MpNews                 *MpNewsType   `json:"mpnews"`                   // 图文消息(mpnews)
	Markdown               *MarkdownType `json:"markdown"`                 // markdown消息
}

// ErrorType 发送消息返回错误
type ErrorType struct {
	ErrCode int    `json:"errcode"` // 出错返回码，为0表示成功，非0表示调用失败
	ErrMsg  string `json:"errmsg"`  // 返回码提示语
}

// ResultType 发送消息返回结果
type ResultType struct {
	ErrorType
	InvalidUser  string `json:"invaliduser"`
	InvalidParty string `json:"infvalidparty"`
	InvalidTag   string `json:"invalidtag"`
}

// AccessTokenType 微信企业号 AccessTokenType
type AccessTokenType struct {
	ErrorType
	AccessToken   string `json:"access_token"` // 获取到的凭证，最长为512字节
	ExpiresIn     int    `json:"expires_in"`   // 凭证的有效时间（秒）
	ExpiresInTime time.Time
	Mux           sync.Mutex
}

// ClientType 客户端
type ClientType struct {
	AccessTokenUrl string // https://qyapi.weixin.qq.com/cgi-bin/gettoken
	SendUrl        string // https://qyapi.weixin.qq.com/cgi-bin/message/send
	CorpId         string // 企业ID，获取方式参考：术语说明-corpid(https://work.weixin.qq.com/api/doc/90000/90135/91039#14953/corpid)
	CorpSecret     string // 应用的凭证密钥，获取方式参考：术语说明-secret(https://work.weixin.qq.com/api/doc/90000/90135/91039#14953/secret)
	Message        *MessageType
	AccessToken    *AccessTokenType
}

func NewClient(accessTokenUrl string, sendUrl string, corpId string, corpSecret string, message *MessageType) *ClientType {
	return &ClientType{
		AccessTokenUrl: accessTokenUrl,
		SendUrl:        sendUrl,
		CorpId:         corpId,
		CorpSecret:     corpSecret,
		Message:        message,
		AccessToken:    new(AccessTokenType),
	}
}

// IsExpire 验证 access_token 是否过期
func (token *AccessTokenType) IsExpire() bool {
	return token.ExpiresInTime.Before(time.Now())
}

func DateTimeFormatter(dt time.Time) string {
	return dt.Format("2006-01-02 15:04:05")
}

// RefreshAccessToken 用于刷新 access_token
// corpid 每个企业都拥有唯一的corpid，获取此信息可在管理后台“我的企业”－“企业信息”下查看“企业ID”（需要有管理员权限）
// corpsecret 每一个应用都有一个独立的访问密钥，为了保证数据的安全，secret务必不能泄漏, 自建应用secret。在管理后台->“应用与小程序”->“应用”->“自建”，点进某个应用，即可看到。
func (client *ClientType) RefreshAccessToken() error {
	client.AccessToken.Mux.Lock()
	defer client.AccessToken.Mux.Unlock()

	accessToken := AccessTokenType{}
	httpClient := resty.New()
	params := map[string]string{
		"corpid":     client.CorpId,
		"corpsecret": client.CorpSecret,
	}
	resp, err := httpClient.R().
		SetQueryParams(params).
		Get(client.AccessTokenUrl)
	if err != nil {
		return fmt.Errorf("invoke api gettoken fail: %v", err)
	}

	jsonErr := json.Unmarshal(resp.Body(), &accessToken)
	if jsonErr != nil {
		return fmt.Errorf("parse gettoken response body fail: %v", jsonErr)
	}

	if accessToken.ExpiresIn == 0 || accessToken.AccessToken == "" {
		return fmt.Errorf("invoke api gettoken fail, ErrCode: %v, ErrMsg: %v", accessToken.ErrCode, accessToken.ErrMsg)
	}

	duration := time.Duration(accessToken.ExpiresIn) * time.Second
	accessToken.ExpiresInTime = time.Now().Add(duration)
	client.AccessToken = &accessToken

	if Cache != nil {
		token, _ := json.Marshal(&accessToken)
		set := map[string][]byte{
			client.CorpId: token,
		}
		err = Cache.Set(set)
		if err != nil {
			return err
		}
	}
	return nil
}

// getAccessTokenFromCache 从缓存中获取 access_token
func (client *ClientType) getAccessTokenFromCache() (string, error) {
	if Cache == nil {
		return "", fmt.Errorf("client cache processor not found")
	}

	token, err := Cache.GetByKey(client.CorpId)
	if err != nil {
		log.Println("get accessToken fail, error: ", err)
	}

	err = json.Unmarshal(token, &client.AccessToken)
	if client.AccessToken.IsExpire() || client.AccessToken.AccessToken == "" {
		err = client.RefreshAccessToken()
	}
	log.Println("accessToken expire time:", DateTimeFormatter(client.AccessToken.ExpiresInTime))
	log.Println("current time:", DateTimeFormatter(time.Now()))
	return client.AccessToken.AccessToken, err
}

// GetAccessToken 获取 access_token
func (client *ClientType) GetAccessToken() (string, error) {
	// 如果设置了 缓存器，从缓存器中获取 token，防止频繁刷新
	if Cache != nil {
		return client.getAccessTokenFromCache()
	}
	var err error
	if client.AccessToken.IsExpire() {
		err = client.RefreshAccessToken()
		log.Println("accessToken expire time:", DateTimeFormatter(client.AccessToken.ExpiresInTime))
		log.Println("current time:", DateTimeFormatter(time.Now()))
	}
	return client.AccessToken.AccessToken, err
}

// SendMessage 发送消息
func (client *ClientType) SendMessage() (bool, error) {
	msg, err := json.Marshal(client.Message)
	if err != nil {
		return false, fmt.Errorf("parse message fail: %v", err)
	}

	accessToken, err := client.GetAccessToken()
	if err != nil {
		return false, fmt.Errorf("invoke api gettoken fail, error: %v", err)
	}

	httpClient := resty.New()
	headers := map[string]string{
		"Content-Type": "application/json;charset=utf-8",
	}
	params := map[string]string{
		"access_token": accessToken,
	}
	resp, err := httpClient.R().
		SetHeaders(headers).
		SetQueryParams(params).
		SetBody(msg).
		Post(client.SendUrl)
	if err != nil {
		return false, fmt.Errorf("invoke send api fail: %v", err)
	}

	result := ResultType{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return false, fmt.Errorf("parse send api response fail: %v", err)
	}

	if result.InvalidUser != "" || result.InvalidParty != "" || result.InvalidTag != "" {
		return false, fmt.Errorf("invoke send api partial fail, invalid user: %s, invalid party: %s, invalid tag: %s", result.InvalidUser, result.InvalidParty, result.InvalidTag)
	}

	if result.ErrCode != 0 {
		return false, fmt.Errorf("invoke send api return ErrCode = %d", result.ErrCode)
	} else {
		return true, nil
	}
}
