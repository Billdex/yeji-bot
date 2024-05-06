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

func QueryMaterial(ctx context.Context, api *openapi.Openapi, msg *qbot.WSGroupAtMessageData) (err error) {
	args := strings.Split(msg.Content, " ")

	materialNames, err := dao.MatchRecipeMaterialName(ctx, args[0])
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

	if len(materialNames) == 0 {
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: fmt.Sprintf("厨师长说没有用%s做过菜", msg.Content),
			MsgType: openapi.MsgTypeText,
			MsgId:   msg.Id,
			MsgSeq:  kit.Seq(ctx),
		})
		if err != nil {
			logrus.WithContext(ctx).Errorf("post group message fail. err: %v", err)
		}
		return nil
	}

	if len(materialNames) > 1 {
		contentMsg := "想查询哪个食材呢？"
		for i := 0; i < 10 && i < len(materialNames); i++ {
			contentMsg += fmt.Sprintf("\n%s", materialNames[i])
		}
		if len(materialNames) > 10 {
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

	recipeMaterials, err := dao.ListRecipeMaterials(ctx, materialNames[0])
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

	contentMsg := kit.PaginationOutput(recipeMaterials, page, 10, fmt.Sprintf("食材%s消耗数据", materialNames[0]), func(recipeMaterial model.RecipeMaterial) string {
		return fmt.Sprintf("%d %s %s", recipeMaterial.RecipeId, recipeMaterial.RecipeName, fmt.Sprintf("🥗%d/h", recipeMaterial.Efficiency))
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
