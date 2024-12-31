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

// QueryComboRecipe 查询后厨合成菜谱
func QueryComboRecipe(ctx context.Context, api *openapi.Openapi, msg *qbot.WSGroupAtMessageData) (err error) {
	allComboRecipes, err := dao.ListAllComboRecipes(ctx)
	if err != nil {
		logrus.WithContext(ctx).Errorf("list all combo recipes fail. err: %v", err)
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
	if len(allComboRecipes) == 0 {
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: "加载后厨合成菜谱数据失败",
			MsgType: openapi.MsgTypeText,
			MsgId:   msg.Id,
			MsgSeq:  kit.Seq(ctx),
		})
		if err != nil {
			logrus.WithContext(ctx).Errorf("post group message fail. err: %v", err)
		}
		return nil
	}

	matchedComboRecipes := make([]model.ComboRecipe, 0, len(allComboRecipes))
	for _, comboRecipe := range allComboRecipes {
		if strings.Contains(comboRecipe.RecipeName, msg.Content) {
			matchedComboRecipes = append(matchedComboRecipes, comboRecipe)
		}
	}
	if len(matchedComboRecipes) == 0 {
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: fmt.Sprintf("%s不是后厨合成菜", msg.Content),
			MsgType: openapi.MsgTypeText,
			MsgId:   msg.Id,
			MsgSeq:  kit.Seq(ctx),
		})
		if err != nil {
			logrus.WithContext(ctx).Errorf("post group message fail. err: %v", err)
			return nil
		}
	} else if len(matchedComboRecipes) > 1 {
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: fmt.Sprintf("想查询哪道菜呢?\n%s",
				strings.Join(
					kit.SliceMap(matchedComboRecipes, func(combo model.ComboRecipe) string { return combo.RecipeName }),
					"\n",
				),
			),
			MsgType: openapi.MsgTypeText,
			MsgId:   msg.Id,
			MsgSeq:  kit.Seq(ctx),
		})
		if err != nil {
			logrus.WithContext(ctx).Errorf("post group message fail. err: %v", err)
			return nil
		}
	} else {
		// 匹配一道菜时查询对应所需菜谱信息
		comboRecipe := matchedComboRecipes[0]
		recipeMap, err := dao.ListRecipesMapByRecipeIds(ctx, comboRecipe.NeedRecipeIds)
		logrus.Infof("need recipe Ids: %v, recipeMap: %v", comboRecipe.NeedRecipeIds, recipeMap)
		if err != nil {
			logrus.WithContext(ctx).Errorf("list combo need recipes fail. err: %v", err)
			_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
				Content: err.Error(),
				MsgType: openapi.MsgTypeText,
				MsgId:   msg.Id,
				MsgSeq:  kit.Seq(ctx),
			})
			if err != nil {
				logrus.WithContext(ctx).Errorf("post group message fail. err: %v", err)
				return nil
			}
			return nil
		}
		content := fmt.Sprintf("合成「%s」需要以下菜谱", comboRecipe.RecipeName)
		for _, recipeId := range comboRecipe.NeedRecipeIds {
			recipe := recipeMap[recipeId]
			content += fmt.Sprintf("\n%s %s %s", recipe.GalleryId, recipe.Name, recipe.Origins)
		}
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: content,
			MsgType: openapi.MsgTypeText,
			MsgId:   msg.Id,
			MsgSeq:  kit.Seq(ctx),
		})
		if err != nil {
			logrus.WithContext(ctx).Errorf("post group message fail. err: %v", err)
			return nil
		}
	}

	return nil
}
