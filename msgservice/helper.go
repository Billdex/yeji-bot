package msgservice

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"yeji-bot/bot/openapi"
	"yeji-bot/bot/qbot"
)

func introHelp(ctx context.Context, api *openapi.Openapi, msg *qbot.WSGroupAtMessageData) (err error) {
	logrus.Info("send help message")
	const preChar = "/"
	sb := strings.Builder{}
	sb.WriteString("【爆炒江湖查询机器人】\n")
	sb.WriteString(fmt.Sprintf("使用方式『%s功能名 参数』\n", preChar))
	sb.WriteString(fmt.Sprintf("示例「%s厨师 羽十六」\n", preChar))
	//sb.WriteString("\n")
	//sb.WriteString("详情请看说明文档:\n")
	//sb.WriteString("http://bcjhbot.billdex.cn\n")
	//sb.WriteString("数据来源: L图鉴网/\n")
	//sb.WriteString("https://foodgame.github.io")
	_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
		Content: sb.String(),
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
