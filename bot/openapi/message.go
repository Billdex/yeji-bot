package openapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// 发送消息协议
// https://bot.q.qq.com/wiki/develop/api-v2/server-inter/message/send-receive/send.html

type MsgType int

const (
	MsgTypeText     MsgType = 0 // 文本消息
	MsgTypeMarkdown MsgType = 2 // markdown 格式消息
	MsgTypeArk      MsgType = 3 // ark 消息
	MsgTypeEmbed    MsgType = 4 // embed 消息
	MsgTypeMedia    MsgType = 7 // 富媒体消息, 包含图片、视频
)

type PostGroupMessageReq struct {
	Content  string    `json:"content"`            // 消息文本内容
	MsgType  MsgType   `json:"msg_type"`           // 消息类型
	Markdown *Markdown `json:"markdown,omitempty"` // markdown 格式消息，MsgType 为  MsgTypeMarkdown(2) 时可用
	Keyboard *Keyboard `json:"keyboard,omitempty"` // 消息按钮，放在文字内容底部
	Media    *Media    `json:"media,omitempty"`    // 富媒体消息，MsgType 为  MsgTypeMedia(7) 时可用
	MsgId    string    `json:"msg_id"`             // 前置收到的用户发送过来的消息 ID，用于发送被动消息（回复）
	MsgSeq   int       `json:"msg_seq"`            // 回复消息的序号，与 msg_id 联合使用，避免相同消息id回复重复发送，不填默认是 1。相同的 msg_id + msg_seq 重复发送会失败。
}

type Markdown struct {
	Content string `json:"content"` // markdown 内容
}

// Keyboard 消息按钮 https://bot.q.qq.com/wiki/develop/api-v2/server-inter/message/trans/msg-btn.html
type Keyboard struct {
	Content KeyboardContent `json:"content"` // markdown 内容
}

type KeyboardContent struct {
	Rows []KeyboardButtonsRow `json:"rows"`
}

type KeyboardButtonsRow struct {
	Buttons []KeyboardButton `json:"buttons"`
}

type KeyboardButton struct {
	Id         string                   `json:"id"`
	RenderData KeyboardButtonRenderData `json:"render_data"` // 按钮渲染数据
}

type KeyboardButtonRenderDataStyle int

const (
	KeyboardButtonRenderDataStyleGray KeyboardButtonRenderDataStyle = 0 // 按钮样式 灰色线框
	KeyboardButtonRenderDataStyleBlue KeyboardButtonRenderDataStyle = 1 // 按钮样式 蓝色现况
)

type KeyboardButtonRenderData struct {
	Label        string                        `json:"label"`         // 按钮文案
	VisitedLabel string                        `json:"visited_label"` // 点击后的按钮文案
	Style        KeyboardButtonRenderDataStyle `json:"style"`         // 按钮样式
	Action       KeyboardButtonAction          `json:"action"`        // 按钮响应动作
}

type KeyboardButtonActionType int

const (
	KeyboardButtonActionTypeLink        = 0 // 跳转按钮：http 或 小程序 客户端识别 scheme
	KeyboardButtonActionTypeCallback    = 1 // 回调按钮：回调后台接口, data 传给后台
	KeyboardButtonActionTypeInstruction = 2 // 指令按钮：自动在输入框插入 @bot data
)

type KeyboardButtonActionAnchor int

const (
	KeyboardButtonActionAnchorQImage = 1 // 唤起手 Q 选图器
)

type KeyboardButtonAction struct {
	Type          KeyboardButtonActionType       `json:"type"`           // 按钮类型
	Permission    KeyboardButtonActionPermission `json:"permission"`     // 按钮权限
	Data          string                         `json:"data"`           // 操作相关的数据
	Reply         bool                           `json:"reply"`          // 指令按钮可用 指令是否带引用回复本消息
	Enter         bool                           `json:"enter"`          // 指令按钮可用 点击按钮后直接自动发送 data
	Anchor        KeyboardButtonActionAnchor     `json:"anchor"`         // 指令按钮可用 设置后后会忽略 action.enter 配置
	UnSupportTips string                         `json:"unsupport_tips"` // 客户端不支持本 action 时候弹出的 toast文案
}

type KeyboardButtonActionPermissionType int

const (
	KeyboardButtonActionPermissionTypeSpecialUser = 0 // 指定用户可操作
	KeyboardButtonActionPermissionTypeAdmin       = 1 // 仅管理者可操作
	KeyboardButtonActionPermissionTypeAllUser     = 2 // 所有人可操作
)

type KeyboardButtonActionPermission struct {
	Type           KeyboardButtonActionPermissionType `json:"type"`             // 权限类型
	SpecifyUserIds []string                           `json:"specify_user_ids"` // 有权限的用户 id 列表
}

type Media struct {
	FileInfo string `json:"file_info,omitempty"`
}

func (p *PostGroupMessageReq) Method() string {
	return http.MethodPost
}

func (p *PostGroupMessageReq) URI() string {
	return postGroupMessageURI
}

func (p *PostGroupMessageReq) Marshal() ([]byte, error) {
	return json.Marshal(*p)
}

type PostGroupMessageResp struct {
	Id        string    `json:"id"`        // 消息唯一 ID
	Timestamp time.Time `json:"timestamp"` // 发送时间
}

func (p *PostGroupMessageResp) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, p)
}

func (o *Openapi) PostGroupMessage(ctx context.Context, groupOpenid string, req *PostGroupMessageReq) (*PostGroupMessageResp, error) {
	resp := new(PostGroupMessageResp)
	if req.MsgType == MsgTypeText || req.MsgType == MsgTypeMarkdown {
		// 回复消息会带 at, 为保证文本格式首行换行
		req.Content = "\n" + req.Content
	}
	err := o.request(ctx, req, resp, map[string]string{
		"group_openid": groupOpenid,
	})
	return resp, err
}
