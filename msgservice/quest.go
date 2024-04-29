package msgservice

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"yeji-bot/bot/openapi"
	"yeji-bot/bot/qbot"
	"yeji-bot/dao"
	"yeji-bot/pkg/seq"
)

// QueryQuest 任务数据查询
func QueryQuest(ctx context.Context, api *openapi.Openapi, msg *qbot.WSGroupAtMessageData) (err error) {
	args := strings.Split(msg.Content, " ")
	if len(args) == 0 {
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: "任务查询格式有误",
			MsgType: openapi.MsgTypeText,
			MsgId:   msg.Id,
			MsgSeq:  seq.Seq(ctx),
		})
		if err != nil {
			logrus.WithContext(ctx).Errorf("post group message fail. err: %v", err)
		}
		return nil
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		logrus.WithContext(ctx).Errorf("解析任务格式失败. raw: %s, id: %d, err: %v", msg.Content, id, err)
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: "任务查询格式有误",
			MsgType: openapi.MsgTypeText,
			MsgId:   msg.Id,
			MsgSeq:  seq.Seq(ctx),
		})
		if err != nil {
			logrus.WithContext(ctx).Errorf("post group message fail. err: %v", err)
		}
		return nil
	}
	var limit = 1
	if len(args) > 1 {
		limit, err = strconv.Atoi(args[1])
		if err != nil {
			logrus.WithContext(ctx).Errorf("解析任务数量失败. raw: %s, id: %d, err: %v", msg.Content, id, err)
			limit = 1
		}
	}
	quests, err := dao.ListMainQuestsWithLimit(ctx, id, limit)
	if err != nil {
		logrus.WithContext(ctx).Errorf("获取任务数据失败. id: %d, err: %v", id, err)
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: "获取任务数据失败",
			MsgType: openapi.MsgTypeText,
			MsgId:   msg.Id,
			MsgSeq:  seq.Seq(ctx),
		})
		if err != nil {
			logrus.WithContext(ctx).Errorf("post group message fail. err: %v", err)
		}
		return nil
	}
	if len(quests) == 0 {
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: "任务不存在",
			MsgType: openapi.MsgTypeText,
			MsgId:   msg.Id,
			MsgSeq:  seq.Seq(ctx),
		})
		if err != nil {
			logrus.WithContext(ctx).Errorf("post group message fail. err: %v", err)
		}
		return nil
	}

	contentMsg := "[主线任务]"
	for _, quest := range quests {
		contentMsg += fmt.Sprintf("\n[%d] %s\n奖励: %s", quest.QuestId, quest.Goal, quest.RewardsStr())
	}

	_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
		Content: contentMsg,
		MsgType: openapi.MsgTypeText,
		MsgId:   msg.Id,
		MsgSeq:  seq.Seq(ctx),
	})
	if err != nil {
		logrus.WithContext(ctx).Errorf("post group message fail. err: %v", err)
	}

	return nil
}
