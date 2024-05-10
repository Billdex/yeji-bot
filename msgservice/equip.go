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

func QueryEquip(ctx context.Context, api *openapi.Openapi, msg *qbot.WSGroupAtMessageData) (err error) {
	var (
		order = "稀有度"
		page  = 1
	)

	equips, err := dao.ListAllEquips(ctx)
	if err != nil {
		logrus.WithContext(ctx).Errorf("list all equips fail. err: %v", err)
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
	if len(equips) == 0 {
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: "加载厨具数据失败",
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
		switch {
		case kit.SliceContains([]string{"图鉴序", "稀有度"}, arg):
			order = arg
		case model.IsRarityStr(arg):
			equips = filterEquipsByRarity(ctx, equips, model.RarityToInt(arg))
		default:
			equips, err = filterEquipsByIdOrName(ctx, equips, arg)
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

	// 排序
	equips = sortEquips(ctx, equips, order)

	// 输出结果
	msgReq := generateEquipsMessage(ctx, equips, order, page)
	msgReq.MsgId = msg.Id
	msgReq.MsgSeq = kit.Seq(ctx)
	_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &msgReq)
	if err != nil {
		logrus.WithContext(ctx).Errorf("post group message fail. err: %v", err)
	}

	return nil
}

func filterEquipsByIdOrName(ctx context.Context, equips []model.Equip, arg string) ([]model.Equip, error) {
	if len(equips) == 0 {
		return equips, nil
	}
	result := make([]model.Equip, 0, len(equips))
	numId, err := strconv.Atoi(arg)
	if err != nil {
		re, err := regexp.Compile(strings.ReplaceAll(arg, "%", ".*"))
		if err != nil {
			logrus.WithContext(ctx).Errorf("查询正则格式有误 raw: %s, err: %v", arg, err)
			return nil, errors.New("厨具查询格式有误")
		}
		for i := range equips {
			if equips[i].Name == arg {
				return []model.Equip{equips[i]}, nil
			}
			if re.MatchString(equips[i].Name) {
				result = append(result, equips[i])
			}
		}
	} else {
		result = kit.SliceFilter(equips, func(equip model.Equip) bool {
			return equip.EquipId == numId
		})
	}
	return result, nil
}

// filterEquipsByRarity 根据稀有度筛选厨具
func filterEquipsByRarity(ctx context.Context, equips []model.Equip, rarity int) []model.Equip {
	if len(equips) == 0 {
		return equips
	}
	return kit.SliceFilter(equips, func(equip model.Equip) bool {
		return equip.Rarity >= rarity
	})
}

// sortEquips 根据排序参数排序厨具
func sortEquips(ctx context.Context, equips []model.Equip, order string) []model.Equip {
	if len(equips) == 0 {
		return equips
	}
	switch order {
	case "图鉴序":
		sort.Slice(equips, func(i, j int) bool {
			return equips[i].EquipId < equips[j].EquipId
		})
	case "稀有度":
		sort.Slice(equips, func(i, j int) bool {
			return equips[i].Rarity == equips[j].Rarity && equips[i].EquipId < equips[j].EquipId ||
				equips[i].Rarity > equips[j].Rarity
		})
	}
	return equips
}

// getEquipInfoWithOrder 根据排序参数获取厨具需要输出的信息
func getEquipInfoWithOrder(equip model.Equip, order string) string {
	switch order {
	case "图鉴序", "稀有度":
		return strings.Repeat("🔥", equip.Rarity)
	default:
		return ""
	}
}

func generateEquipMessage(ctx context.Context, equip model.Equip) openapi.PostGroupMessageReq {
	if equip.Img != "" {
		// TODO 发送图片
		return openapi.PostGroupMessageReq{
			Content: "",
			MsgType: openapi.MsgTypeMedia,
			Media:   &openapi.Media{FileInfo: ""},
		}
	}
	logrus.WithContext(ctx).Infof("未找到厨具 %d %s 图鉴图片, 以文字格式发送数据", equip.EquipId, equip.Name)

	mSkills, err := dao.GetSkillsMapByIds(ctx, equip.Skills)
	if err != nil {
		logrus.WithContext(ctx).Error("查询技能数据出错!", err)
		return openapi.PostGroupMessageReq{
			Content: "哎呀，系统开小差了",
			MsgType: openapi.MsgTypeText,
		}
	}
	skills := make([]model.Skill, 0, len(equip.Skills))
	for _, skill := range equip.Skills {
		skills = append(skills, mSkills[skill])
	}
	sort.Slice(skills, func(i, j int) bool {
		return skills[i].SkillId < skills[j].SkillId
	})
	skillDescs := kit.SliceMap(skills, func(skill model.Skill) string {
		return fmt.Sprintf("%s", skill.Description)
	})
	var msg string
	msg += fmt.Sprintf("%s %s\n", equip.GalleryId, equip.Name)
	msg += fmt.Sprintf("%s\n", strings.Repeat("🔥", equip.Rarity))
	msg += fmt.Sprintf("来源: %s\n", strings.Join(equip.Origins, ","))
	msg += fmt.Sprintf("效果: %s", strings.Join(skillDescs, ","))
	return openapi.PostGroupMessageReq{
		Content: msg,
		MsgType: openapi.MsgTypeText,
	}
}

func generateEquipsMessage(ctx context.Context, equips []model.Equip, order string, page int) openapi.PostGroupMessageReq {
	if len(equips) == 0 {
		return openapi.PostGroupMessageReq{
			Content: "是没见过的厨具呢！",
			MsgType: openapi.MsgTypeText,
		}
	} else if len(equips) == 1 {
		return generateEquipMessage(ctx, equips[0])
	} else {
		msg := kit.PaginationOutput(equips, page, 10, "仓库里翻到了这些厨具", func(equip model.Equip) string {
			return fmt.Sprintf("%s %s %s", equip.GalleryId, equip.Name, getEquipInfoWithOrder(equip, order))
		})
		return openapi.PostGroupMessageReq{
			Content: msg,
			MsgType: openapi.MsgTypeText,
		}
	}
}
