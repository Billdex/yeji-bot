package msgservice

import (
	"context"
	"github.com/sirupsen/logrus"
	"strings"
	"yeji-bot/bot/openapi"
	"yeji-bot/bot/qbot"
)

func IntroHelp(ctx context.Context, api *openapi.Openapi, msg *qbot.WSGroupAtMessageData) (err error) {
	_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
		Content: introHelperStr(),
		MsgType: openapi.MsgTypeText,
		Media:   openapi.Media{},
		MsgId:   msg.Id,
		MsgSeq:  msg.Seq,
	})
	if err != nil {
		logrus.Errorf("send help message fail. err: %v", err)
	}
	return nil
}

func introHelperStr() string {
	sb := strings.Builder{}
	sb.WriteString("【爆炒江湖叶姬小助手】\n")
	sb.WriteString("使用方式『/功能名 参数』\n")
	sb.WriteString("示例「/厨师 羽十六」\n")
	sb.WriteString("目前提供以下数据查询功能: \n")
	sb.WriteString("厨师  菜谱  厨具\n")
	sb.WriteString("食材  贵客  符文\n")
	sb.WriteString("任务\n")
	sb.WriteString("\n")
	sb.WriteString("更多功能开发中...")

	//sb.WriteString("\n")
	//sb.WriteString("详情请看说明文档:\n")
	//sb.WriteString("http://bcjhbot.billdex.cn\n")
	//sb.WriteString("数据来源: L图鉴网/\n")
	//sb.WriteString("https://foodgame.github.io")
	return sb.String()
}

func questHelperStr() string {
	sb := strings.Builder{}
	sb.WriteString("【任务信息查询】\n")
	sb.WriteString("目前提供主线任务查询\n")
	sb.WriteString("示例「/任务 100」\n")
	sb.WriteString("\n")
	sb.WriteString("可一次查询最多五条数据\n")
	sb.WriteString("示例「/任务 100 5」")
	return sb.String()
}
