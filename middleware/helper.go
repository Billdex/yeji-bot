package middleware

import (
	"context"
	"yeji-bot/bot/openapi"
	"yeji-bot/bot/qbot"
	"yeji-bot/bot/scheduler"
)

func Helper(fn func() string) scheduler.GroupAtMessageHandlerMiddleware {
	return func(handler scheduler.GroupAtMessageHandlerFunc) scheduler.GroupAtMessageHandlerFunc {
		return func(ctx context.Context, api *openapi.Openapi, msg *qbot.WSGroupAtMessageData) (err error) {
			if msg.Content == "" {
				_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
					Content: fn(),
					MsgType: openapi.MsgTypeText,
					MsgId:   msg.Id,
					MsgSeq:  msg.Seq,
				})
				return nil
			}
			return handler(ctx, api, msg)
		}
	}
}
