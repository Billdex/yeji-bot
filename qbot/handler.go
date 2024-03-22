package qbot

import "yeji-bot/qbot/openapi"

// DefaultHandler 未知消息进行默认统一处理的 handler
type DefaultHandler func(api *openapi.Openapi, event *WSPayload) error

// ReadyHandler 完成连接后首条消息触发
type ReadyHandler func(api *openapi.Openapi, event *WSPayload, data *WSReadyData) error

// GroupAtMessageHandler 收到群组@消息触发
type GroupAtMessageHandler func(api *openapi.Openapi, event *WSPayload, data *WSGroupAtMessageData) error

type EventHandlers struct {
	DefaultHandler DefaultHandler

	ReadyHandler ReadyHandler // 处理 ready 事件

	GroupAtMessageHandler GroupAtMessageHandler
}
