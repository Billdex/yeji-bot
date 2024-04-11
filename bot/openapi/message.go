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
	Content string  `json:"content"`  // 消息文本内容
	MsgType MsgType `json:"msg_type"` // 消息类型
	Media   Media   `json:"media"`    // 富媒体消息，MsgType 为  MsgTypeMedia(7) 时可用
	MsgId   string  `json:"msg_id"`   // 前置收到的用户发送过来的消息 ID，用于发送被动消息（回复）
	MsgSeq  int     `json:"msg_seq"`  // 回复消息的序号，与 msg_id 联合使用，避免相同消息id回复重复发送，不填默认是 1。相同的 msg_id + msg_seq 重复发送会失败。
}

type Media struct {
	FileInfo string `json:"file_info"`
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
	err := o.request(ctx, req, resp, map[string]string{
		"group_openid": groupOpenid,
	})
	return resp, err
}
