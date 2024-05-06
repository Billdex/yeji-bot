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

func QueryAntique(ctx context.Context, api *openapi.Openapi, msg *qbot.WSGroupAtMessageData) error {
	args := strings.Split(msg.Content, " ")

	giftNames, err := dao.MatchGiftName(ctx, args[0])
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

	if len(giftNames) == 0 {
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: "没有找到符文礼物数据",
			MsgType: openapi.MsgTypeText,
			MsgId:   msg.Id,
			MsgSeq:  kit.Seq(ctx),
		})
		if err != nil {
			logrus.WithContext(ctx).Errorf("post group message fail. err: %v", err)
		}
		return nil
	}

	if len(giftNames) > 1 {
		contentMsg := "想查询哪个符文呢？"
		for i := 0; i < 10 && i < len(giftNames); i++ {
			contentMsg += fmt.Sprintf("\n%s", giftNames[i])
		}
		if len(giftNames) > 10 {
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

	antiqueGifts, err := dao.ListAntiqueGifts(ctx, giftNames[0])
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

	contentMsg := kit.PaginationOutput(antiqueGifts, page, 10, fmt.Sprintf("以下菜谱概率获得%s", giftNames[0]), func(gift model.Gift) string {
		return fmt.Sprintf("%s-%s %s", gift.RecipeName, gift.GuestName, kit.FormatRecipeTime(gift.RecipeTotalTime))
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
