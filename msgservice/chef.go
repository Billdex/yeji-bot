package msgservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"yeji-bot/bot/openapi"
	"yeji-bot/bot/qbot"
	"yeji-bot/dao"
	"yeji-bot/data/model"
	"yeji-bot/pkg/kit"
)

func QueryChef(ctx context.Context, api *openapi.Openapi, msg *qbot.WSGroupAtMessageData) (err error) {
	var (
		order = "稀有度"
		page  = 1
	)

	chefs, err := dao.ListAllChefs(ctx)
	if err != nil {
		logrus.WithContext(ctx).Errorf("list all chefs fail. err: %v", err)
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: "查询厨师失败",
			MsgType: openapi.MsgTypeText,
			MsgId:   msg.Id,
			MsgSeq:  kit.Seq(ctx),
		})
		if err != nil {
			logrus.WithContext(ctx).Errorf("post group message fail. err: %v", err)
		}
		return nil
	}
	if len(chefs) == 0 {
		logrus.WithContext(ctx).Errorf("no chef")
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: "加载厨师数据失败",
			MsgType: openapi.MsgTypeText,
			MsgId:   msg.Id,
			MsgSeq:  kit.Seq(ctx),
		})
		if err != nil {
			logrus.WithContext(ctx).Errorf("post group message fail. err: %v", err)
		}
		return nil
	}
	args := strings.Split(msg.Content, " ")
	for _, arg := range args {
		if arg == "" {
			continue
		}
		switch arg {
		default:
			if strings.HasPrefix(arg, "来源") {
				origin := strings.Split(arg, "-")
				if len(origin) > 1 {
					chefs, err = filterChefsByOrigin(ctx, chefs, strings.Join(origin[1:], "-"))
				}
			} else if strings.HasPrefix(arg, "技能") {
				skill := strings.Split(arg, "-")
				if len(skill) > 1 {
					chefs, err = filterChefsBySkill(chefs, strings.Join(skill[1:], "-"))
				}
			} else if strings.ToLower(arg[:1]) == "p" {
				var pageNum int
				pageNum, err = strconv.Atoi(strings.Trim(arg[1:], "-"))
				if err != nil {
					err = errors.New("分页参数有误")
				} else {
					if pageNum > 0 {
						page = pageNum
					}
				}
			} else {
				chefs, err = filterChefsByIdOrName(ctx, chefs, arg)
			}

		}
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
	}

	// 对菜谱查询结果排序
	chefs, err = orderChefs(ctx, chefs, order)
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
		return
	}

	msgReq := generateChefsMessage(ctx, chefs, order, page)
	msgReq.MsgId = msg.Id
	_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &msgReq)
	if err != nil {
		logrus.WithContext(ctx).Errorf("post group message fail. err: %v", err)
	}

	return nil
}

// 根据厨师名或厨师ID筛选厨师
func filterChefsByIdOrName(ctx context.Context, chefs []model.Chef, arg string) ([]model.Chef, error) {
	result := make([]model.Chef, 0)
	numId, err := strconv.Atoi(arg)
	if err != nil {
		re, err := regexp.Compile(strings.ReplaceAll(arg, "%", ".*"))
		if err != nil {
			logrus.WithContext(ctx).Errorf("查询正则格式有误", err)
			return nil, errors.New("查询格式有误")
		}
		for i := range chefs {
			if chefs[i].Name == arg {
				return []model.Chef{chefs[i]}, nil
			}
			if re.MatchString(chefs[i].Name) {
				result = append(result, chefs[i])
			}
		}
	} else {
		for i := range chefs {
			if chefs[i].ChefId == (numId+2)/3*3 {
				result = append(result, chefs[i])
			}
		}
	}
	return result, nil
}

// filterChefsByOrigin 根据来源筛选厨师
func filterChefsByOrigin(ctx context.Context, chefs []model.Chef, origin string) ([]model.Chef, error) {
	if len(chefs) == 0 || origin == "" {
		return chefs, nil
	}
	result := make([]model.Chef, 0)
	pattern := strings.ReplaceAll(origin, "%", ".*")
	// 单独增加未入坑礼包查询
	if origin == "仅礼包" || origin == "在售礼包" || origin == "未入坑礼包" {
		pattern = "^限时礼包$"
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		logrus.WithContext(ctx).Errorf("查询正则格式有误. raw: %s, err: %v", pattern, err)
		return nil, errors.New("来源筛选格式有误")
	}
	for i := range chefs {
		for j := range chefs[i].Origins {
			if re.MatchString(chefs[i].Origins[j]) {
				result = append(result, chefs[i])
				break
			}
		}
	}

	return result, nil
}

// filterChefs根据厨师技能筛选厨师
func filterChefsBySkill(chefs []model.Chef, skill string) ([]model.Chef, error) {
	if len(chefs) == 0 || skill == "" {
		return chefs, nil
	}
	// 处理某些技能关键词
	if kit.SliceContains([]string{"炒光环", "烤光环", "煮光环", "蒸光环", "炸光环", "切光环", "光环"}, skill) {
		skill = "场上所有厨师" + strings.TrimSuffix(skill, "光环")
	}
	if kit.SliceContains([]string{"贵客", "贵宾", "客人", "宾客", "稀客"}, skill) {
		skill = "稀有客人"
	}
	if strings.HasPrefix(skill, "采集") {
		skill = "探索" + strings.TrimLeft(skill, "采集")
	}
	pattern := strings.ReplaceAll(skill, "%", ".*")
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("技能描述格式有误 %v", err)
	}
	result := make([]model.Chef, 0)
	for i := range chefs {
		if re.MatchString(chefs[i].SkillDesc) || re.MatchString(chefs[i].UltimateSkillDesc) {
			result = append(result, chefs[i])
		}
	}
	return result, nil
}

// 根据排序参数排序厨师
func orderChefs(ctx context.Context, chefs []model.Chef, order string) ([]model.Chef, error) {
	if len(chefs) == 0 {
		return chefs, nil
	}
	switch order {
	case "图鉴序":
		sort.Slice(chefs, func(i, j int) bool {
			return chefs[i].ChefId < chefs[j].ChefId
		})
	case "稀有度":
		sort.Slice(chefs, func(i, j int) bool {
			return chefs[i].Rarity == chefs[j].Rarity && chefs[i].ChefId < chefs[j].ChefId ||
				chefs[i].Rarity > chefs[j].Rarity
		})
	default:
		return nil, errors.New("排序参数有误")
	}
	return chefs, nil
}

// 输出单厨师消息数据
func genChefMessage(ctx context.Context, chef model.Chef) openapi.PostGroupMessageReq {
	// 尝试寻找图片文件，未找到则按照文字格式发送
	var msg string
	if chef.Img != "" {
		// TODO 发送图片
		return openapi.PostGroupMessageReq{
			Content: msg,
			MsgType: openapi.MsgTypeMedia,
			Media:   openapi.Media{FileInfo: ""},
			MsgSeq:  kit.Seq(ctx),
		}
	}
	logrus.WithContext(ctx).Infof("未找到厨师 %d %s 图鉴图片, 以文字格式发送数据", chef.ChefId, chef.Name)
	var gender string
	if chef.Gender == 1 {
		gender = "♂️"
	} else if chef.Gender == 2 {
		gender = "♀️"
	}
	mSkills, err := dao.GetSkillsMapByIds(ctx, []int{chef.SkillId, chef.UltimateSkill})
	if err != nil {
		logrus.WithContext(ctx).Error("查询技能数据出错!", err)
		return openapi.PostGroupMessageReq{
			Content: "哎呀，系统开小差了",
			MsgType: openapi.MsgTypeText,
			MsgSeq:  kit.Seq(ctx),
		}
	}
	ultimateGoals, err := dao.FindQuestsByIds(ctx, chef.UltimateGoals)
	if err != nil {
		logrus.WithContext(ctx).Error("查询厨师修炼效果数据出错!", err)
		return openapi.PostGroupMessageReq{
			Content: "哎呀，系统开小差了",
			MsgType: openapi.MsgTypeText,
			MsgSeq:  kit.Seq(ctx),
		}
	}
	goals := ""
	for i, ultimateGoal := range ultimateGoals {
		goals += fmt.Sprintf("\n%d.%s", i+1, ultimateGoal.Goal)
	}
	msg += fmt.Sprintf("%s %s %s\n", chef.GalleryId, chef.Name, gender)
	msg += fmt.Sprintf("%s\n", strings.Repeat("🔥", chef.Rarity))
	msg += fmt.Sprintf("来源: %s\n", strings.Join(chef.Origins, ","))
	msg += fmt.Sprintf("炒:%d 烤:%d 煮:%d\n", chef.Stirfry, chef.Bake, chef.Boil)
	msg += fmt.Sprintf("蒸:%d 炸:%d 切:%d\n", chef.Steam, chef.Fry, chef.Cut)
	msg += fmt.Sprintf("🍖:%d 🍞:%d 🥕:%d 🐟:%d\n", chef.Meat, chef.Flour, chef.Vegetable, chef.Fish)
	msg += fmt.Sprintf("技能:%s\n", mSkills[chef.SkillId].Description)
	msg += fmt.Sprintf("修炼效果:%s\n", mSkills[chef.UltimateSkill].Description)
	msg += fmt.Sprintf("修炼任务:%s", goals)

	return openapi.PostGroupMessageReq{
		Content: msg,
		MsgType: openapi.MsgTypeText,
		MsgSeq:  kit.Seq(ctx),
	}
}

// 根据来源和排序参数，输出厨师列表消息数据
func generateChefsMessage(ctx context.Context, chefs []model.Chef, order string, page int) openapi.PostGroupMessageReq {
	if len(chefs) == 0 {
		return openapi.PostGroupMessageReq{
			Content: "诶? 似乎查无此厨哦!",
			MsgType: openapi.MsgTypeText,
			MsgSeq:  kit.Seq(ctx),
		}
	} else if len(chefs) == 1 {
		return genChefMessage(ctx, chefs[0])
	} else {
		msg := kit.PaginationOutput(chefs, page, 10, "查询到以下厨师", func(chef model.Chef) string {
			return fmt.Sprintf("%s %s %s", chef.GalleryId, chef.Name, getChefInfoWithOrder(chef, order))
		})
		return openapi.PostGroupMessageReq{
			Content: msg,
			MsgType: openapi.MsgTypeText,
			MsgSeq:  kit.Seq(ctx),
		}
	}
}

// 根据排序参数获取厨师需要输出的信息
func getChefInfoWithOrder(chef model.Chef, order string) string {
	switch order {
	case "图鉴序", "稀有度":
		return strings.Repeat("🔥", chef.Rarity)
	default:
		return ""
	}
}
