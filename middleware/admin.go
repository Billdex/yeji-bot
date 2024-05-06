package middleware

import (
	"context"
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
				return nil
			}
			return handler(ctx, api, msg)
		}
	}
}
