package qbot

import "time"

// EventType 事件类型
type EventType string

// WSPayload websocket 消息结构
type WSPayload struct {
	WSPayloadBase
	Data       interface{} `json:"d,omitempty"`
	RawMessage []byte      `json:"-"` // 原始的 message 数据
}

// WSPayloadBase 基础消息结构，排除了 data
type WSPayloadBase struct {
	OPCode OPCode    `json:"op"`
	Seq    int       `json:"s,omitempty"`
	Type   EventType `json:"t,omitempty"`
}

// 以下为发送到 websocket 的 data

// WSIdentityData 鉴权数据
type WSIdentityData struct {
	Token      string             `json:"token"`
	Intents    Intent             `json:"intents"`
	Shard      []uint32           `json:"shard"` // array of two integers (shard_id, num_shards)
	Properties IdentityProperties `json:"properties,omitempty"`
}

type IdentityProperties struct {
	Os      string `json:"$os,omitempty"`
	Browser string `json:"$browser,omitempty"`
	Device  string `json:"$device,omitempty"`
}

// WSResumeData 重连数据
type WSResumeData struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Seq       int    `json:"seq"`
}

// 以下为会收到的事件data

// WSHelloData hello 返回
type WSHelloData struct {
	HeartbeatInterval int `json:"heartbeat_interval"`
}

// WSReadyData ready，鉴权后返回
type WSReadyData struct {
	Version   int    `json:"version"`
	SessionID string `json:"session_id"`
	User      struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Bot      bool   `json:"bot"`
	} `json:"user"`
	Shard []uint32 `json:"shard"`
}

type WSGroupAtMessageData struct {
	WSPayloadBase
	Id     string `json:"id"`
	Author struct {
		Id           string `json:"id"`
		MemberOpenid string `json:"member_openid"`
	} `json:"author"`
	Content     string    `json:"content"`
	GroupId     string    `json:"group_id"`
	GroupOpenid string    `json:"group_openid"`
	Timestamp   time.Time `json:"timestamp"`

	CmdName string `json:"-"` // 指令名称
}
