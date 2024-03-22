package openapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type MsgType int

const (
	MsgTypeText     MsgType = 0 // 文本消息
	MsgTypeMarkdown MsgType = 2 // markdown 格式消息
	MsgTypeArk      MsgType = 3 // ark 消息
	MsgTypeEmbed    MsgType = 4 // embed 消息
	MsgTypeMedia    MsgType = 7 // 富媒体消息, 包含图片、视频
)

type PostGroupMessageReq struct {
	Content string  `json:"content"`
	MsgType MsgType `json:"msg_type"`
	Media   Media   `json:"media"`
	MsgId   string  `json:"msg_id"`
	MsgSeq  int     `json:"msg_seq"`
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
