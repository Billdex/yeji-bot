package middleware

import (
	"context"
	"github.com/sirupsen/logrus"
	"yeji-bot/bot/openapi"
	"yeji-bot/bot/qbot"
	"yeji-bot/bot/scheduler"
	"yeji-bot/config"
)

func MustAdmin() scheduler.GroupAtMessageHandlerMiddleware {
	return func(handler scheduler.GroupAtMessageHandlerFunc) scheduler.GroupAtMessageHandlerFunc {
		return func(ctx context.Context, api *openapi.Openapi, msg *qbot.WSGroupAtMessageData) (err error) {
			isAdmin := false
			for _, adminOpenId := range config.AppConfig.Bot.AdminUsers {
				if adminOpenId == msg.Author.MemberOpenid {
					isAdmin = true
					break
				}
			}
			if !isAdmin {
				_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
					Content: "我不听我不听",
					MsgType: openapi.MsgTypeText,
					MsgId:   msg.Id,
					MsgSeq:  msg.Seq,
				})
				if err != nil {
					logrus.WithContext(ctx).Errorf("send message fail. err: %v", err)
				}
				return nil
			}
			return handler(ctx, api, msg)
		}
	}
}
