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
	sb.WriteString("使用方式『/功能名 查询参数』\n")
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

func chefHelperStr() string {
	sb := strings.Builder{}
	sb.WriteString("【厨师信息查询】\n")
	sb.WriteString("提供游戏厨师数据查询\n")
	sb.WriteString("示例:「/厨师 099」 「/厨师 羽十六」")
	return sb.String()
}

func recipeHelperStr() string {
	sb := strings.Builder{}
	sb.WriteString("【菜谱信息查询】\n")
	sb.WriteString("提供游戏菜谱数据查询\n")
	sb.WriteString("示例:「/菜谱 100」 「/菜谱 糖番茄」")
	return sb.String()
}

func equipHelperStr() string {
	sb := strings.Builder{}
	sb.WriteString("【厨具信息查询】\n")
	sb.WriteString("提供游戏厨具数据查询\n")
	sb.WriteString("示例:「/厨具 003」 「/厨具 金烤叉」")
	return sb.String()
}

func materialHelperStr() string {
	sb := strings.Builder{}
	sb.WriteString("【食材消耗效率查询】\n")
	sb.WriteString("根据食材名称查询对应菜谱列表，并按照对应食材的每小时消耗效率升序排序。\n")
	sb.WriteString("结果过多可使用「p」参数分页\n")
	sb.WriteString("示例:「/食材 梅花肉」 「/食材 豆腐 p2」")
	return sb.String()
}

func guestHelperStr() string {
	sb := strings.Builder{}
	sb.WriteString("【贵客信息查询】\n")
	sb.WriteString("根据贵客名查询对应菜谱列表，并按照一组的时间升序排序。\n")
	sb.WriteString("结果过多可使用「p」参数分页\n")
	sb.WriteString("示例:「/贵客 木良」 「/贵客 木优 p2」")
	return sb.String()
}

func antiqueHelperStr() string {
	sb := strings.Builder{}
	sb.WriteString("【符文信息查询】\n")
	sb.WriteString("提供根据符文名查询对应菜谱的功能, 并按照一组的时间升序排序。\n")
	sb.WriteString("结果过多可使用「p」参数分页\n")
	sb.WriteString("示例:「/符文 五香果」 「/符文 一昧真火 p2」")
	return sb.String()
}

func questHelperStr() string {
	sb := strings.Builder{}
	sb.WriteString("【任务信息查询】\n")
	sb.WriteString("提供游戏主线任务查询, 可一次查询最多五条数据。\n")
	sb.WriteString("示例「/任务 100」 「/任务 100 5」")
	return sb.String()
}

func comboRecipeHelperStr() string {
	sb := strings.Builder{}
	sb.WriteString("【后厨合成菜谱查询】\n")
	sb.WriteString("提供后厨菜谱的所需前置菜谱与来源信息\n")
	sb.WriteString("示例「/后厨 年夜饭」 「/后厨 沙县轻食」")
	return sb.String()
}
