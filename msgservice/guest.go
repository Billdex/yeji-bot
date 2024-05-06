package msgservice

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"yeji-bot/bot/openapi"
	"yeji-bot/bot/qbot"
	"yeji-bot/dao"
	"yeji-bot/data/model"
	"yeji-bot/pkg/kit"
)

func QueryGuest(ctx context.Context, api *openapi.Openapi, msg *qbot.WSGroupAtMessageData) (err error) {
	args := strings.Split(msg.Content, " ")

	guestNames, err := dao.MatchGuestName(ctx, args[0])
	if err != nil {
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: err.Error(),
			MsgType: openapi.MsgTypeText,
			MsgId:   msg.Id,
			MsgSeq:  kit.Seq(ctx),
		})
		if err != nil {
			logrus.WithContext(ctx).Errorf("post group message fail. err: %v", err)
		}
		return nil
	}

	if len(guestNames) == 0 {
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: fmt.Sprintf("唔, %s未曾光临本店呢", msg.Content),
			MsgType: openapi.MsgTypeText,
			MsgId:   msg.Id,
			MsgSeq:  kit.Seq(ctx),
		})
		if err != nil {
			logrus.WithContext(ctx).Errorf("post group message fail. err: %v", err)
		}
		return nil
	}

	if len(guestNames) > 1 {
		contentMsg := "想查询哪位贵客呢？"
		for i := 0; i < 10 && i < len(guestNames); i++ {
			contentMsg += fmt.Sprintf("\n%s", guestNames[i])
		}
		if len(guestNames) > 10 {
			contentMsg += "\n......"
		}
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: contentMsg,
			MsgType: openapi.MsgTypeText,
			MsgId:   msg.Id,
			MsgSeq:  kit.Seq(ctx),
		})
		if err != nil {
			logrus.WithContext(ctx).Errorf("post group message fail. err: %v", err)
		}
		return nil
	}

	guestGifts, err := dao.ListGuestGifts(ctx, guestNames[0])
	if err != nil {
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: err.Error(),
			MsgType: openapi.MsgTypeText,
			MsgId:   msg.Id,
			MsgSeq:  kit.Seq(ctx),
		})
		if err != nil {
			logrus.WithContext(ctx).Errorf("post group message fail. err: %v", err)
		}
		return nil
	}

	page := 1
	if len(args) > 1 {
		page = kit.ParsePage(strings.TrimPrefix(msg.Content, args[0]), 1)
	}

	contentMsg := kit.PaginationOutput(guestGifts, page, 10, fmt.Sprintf("%s喜欢这些菜谱", guestNames[0]), func(gift model.Gift) string {
		return fmt.Sprintf("%s-%s %s", gift.RecipeName, gift.Antique, kit.FormatRecipeTime(gift.RecipeTotalTime))
	})
	_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
		Content: contentMsg,
		MsgType: openapi.MsgTypeText,
		MsgId:   msg.Id,
		MsgSeq:  kit.Seq(ctx),
	})
	if err != nil {
		logrus.WithContext(ctx).Errorf("post group message fail. err: %v", err)
	}

	return nil
}
