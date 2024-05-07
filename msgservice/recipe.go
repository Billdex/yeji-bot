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

func QueryRecipe(ctx context.Context, api *openapi.Openapi, msg *qbot.WSGroupAtMessageData) error {
	var (
		order = "稀有度"
		page  = 1
	)

	recipes, err := dao.ListAllRecipes(ctx)
	if err != nil {
		logrus.WithContext(ctx).Errorf("list all recipes fail. err: %v", err)
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
	if len(recipes) == 0 {
		logrus.WithContext(ctx).Errorf("no recipe")
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: "加载菜谱数据失败",
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
		case kit.SliceContains([]string{"图鉴序"}, arg):

		case model.IsRarityStr(arg):
			recipes = filterRecipesByRarity(ctx, recipes, model.RarityToInt(arg))
		case strings.HasPrefix(arg, "来源"):
			recipes, err = filterRecipesByOrigins(ctx, recipes, strings.Split(strings.TrimPrefix(arg, "来源"), "-"))
		case strings.HasPrefix(strings.ToLower(arg), "p"):
			page = kit.ParsePage(arg, 1)
		default:
			recipes, err = filterRecipesByIdOrName(ctx, recipes, arg)
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
	recipes = sortRecipes(recipes, order)

	// 输出结果
	msgReq := generateRecipesMessage(ctx, recipes, order, page)
	msgReq.MsgId = msg.Id
	msgReq.MsgSeq = kit.Seq(ctx)
	_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &msgReq)
	if err != nil {
		logrus.WithContext(ctx).Errorf("post group message fail. err: %v", err)
	}

	return nil
}

// 根据菜谱名或菜谱ID筛选菜谱
func filterRecipesByIdOrName(ctx context.Context, recipes []model.Recipe, arg string) ([]model.Recipe, error) {
	result := make([]model.Recipe, 0, len(recipes))
	numId, err := strconv.Atoi(arg)
	if err != nil {
		re, err := regexp.Compile(strings.ReplaceAll(arg, "%", ".*"))
		if err != nil {
			logrus.WithContext(ctx).Errorf("查询正则格式有误 %v", err)
			return result, errors.New("查询格式有误")
		}
		for _, recipe := range recipes {
			if recipe.Name == arg {
				return []model.Recipe{recipe}, nil
			}
			if re.MatchString(recipe.Name) {
				result = append(result, recipe)
			}
		}
	} else {
		result = kit.SliceFilter(recipes, func(recipe model.Recipe) bool {
			return recipe.RecipeId == numId
		})
	}
	return result, nil
}

// filterRecipesByRarity 根据稀有度筛选菜谱
func filterRecipesByRarity(ctx context.Context, recipes []model.Recipe, rarity int) []model.Recipe {
	if len(recipes) == 0 {
		return recipes
	}
	return kit.SliceFilter(recipes, func(recipe model.Recipe) bool {
		return recipe.Rarity >= rarity
	})
}

// filterRecipesByOrigins 根据来源筛选菜谱
func filterRecipesByOrigins(ctx context.Context, recipes []model.Recipe, origins []string) ([]model.Recipe, error) {
	if len(recipes) == 0 || origins == nil {
		return recipes, nil
	}
	result := make([]model.Recipe, 0)
	copy(result, recipes)
	for _, origin := range origins {
		if origin == "" {
			continue
		}
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
		result = kit.SliceFilter(result, func(recipe model.Recipe) bool {
			for i := range recipe.Origins {
				if re.MatchString(recipe.Origins[i]) {
					return true
				}
			}
			return false
		})
	}

	return result, nil
}

func sortRecipes(recipes []model.Recipe, order string) []model.Recipe {
	if len(recipes) == 0 {
		return recipes
	}
	switch order {
	case "图鉴序":
		sort.Slice(recipes, func(i, j int) bool {
			return recipes[i].RecipeId < recipes[j].RecipeId
		})
	case "稀有度":
		sort.Slice(recipes, func(i, j int) bool {
			return recipes[i].Rarity == recipes[j].Rarity && recipes[i].RecipeId < recipes[j].RecipeId ||
				recipes[i].Rarity > recipes[j].Rarity
		})
	}
	return recipes
}

func generateRecipeMessage(ctx context.Context, recipe model.Recipe) openapi.PostGroupMessageReq {
	// 尝试寻找图片文件，未找到则按照文字格式发送
	if recipe.Img != "" {
		// TODO 发送图片
		return openapi.PostGroupMessageReq{
			Content: "",
			MsgType: openapi.MsgTypeMedia,
			Media:   &openapi.Media{FileInfo: ""},
		}
	}
	logrus.WithContext(ctx).Infof("未找到菜谱 %d %s 图片，以文字格式发送数据", recipe.RecipeId, recipe.Name)
	// 菜谱所需技法数据
	type skillData struct {
		str string
		val int
	}
	var skills = kit.SliceMap(kit.SliceFilter([]skillData{
		{"炒", recipe.Stirfry},
		{"烤", recipe.Bake},
		{"煮", recipe.Boil},
		{"蒸", recipe.Steam},
		{"炸", recipe.Fry},
		{"切", recipe.Cut},
	}, func(skill skillData) bool { return skill.val > 0 }), func(skill skillData) string {
		return fmt.Sprintf("%s*%d", skill.str, skill.val)
	})
	// 食材数据
	var materials = kit.SliceMap(recipe.Materials, func(material model.RecipeMaterial) string {
		return fmt.Sprintf("%s*%d", material.MaterialName, material.Quantity)
	})
	// 贵客礼物数据
	var gifts = kit.SliceMap(recipe.GuestGifts, func(gift model.RecipeGuestGift) string {
		return fmt.Sprintf("%s-%s", gift.GuestName, gift.Antique)
	})
	// 升阶贵客数据
	var guestLevelMap = map[int]string{
		1: "优",
		2: "特",
		3: "神",
	}
	upgradeGuests := make([]string, 0, 3)
	for i := 0; i < 3; i++ {
		if len(recipe.UpgradeGuests) > i && recipe.UpgradeGuests[i] != "" {
			upgradeGuests = append(upgradeGuests, fmt.Sprintf("  %s-%s", guestLevelMap[i+1], recipe.UpgradeGuests[i]))
		} else {
			upgradeGuests = append(upgradeGuests, fmt.Sprintf("  %s-未知", guestLevelMap[i+1]))
		}
	}
	var msg string
	msg += fmt.Sprintf("%s %s %s\n", recipe.GalleryId, recipe.Name, strings.Repeat("🔥", recipe.Rarity))
	msg += fmt.Sprintf("💰: %d(+%d) --- %d/h\n", recipe.Price, recipe.ExPrice, recipe.GoldEfficiency)
	msg += fmt.Sprintf("来源: %s\n", strings.Join(recipe.Origins, ","))
	msg += fmt.Sprintf("单时间: %s\n", kit.FormatRecipeTime(recipe.Time))
	msg += fmt.Sprintf("总时间: %s (%d份)\n", kit.FormatRecipeTime(recipe.TotalTime), recipe.Limit)
	msg += fmt.Sprintf("技法: %s\n", strings.Join(skills, " "))
	msg += fmt.Sprintf("食材: %s\n", strings.Join(materials, ","))
	msg += fmt.Sprintf("耗材效率: %d/h\n", recipe.MaterialEfficiency)
	msg += fmt.Sprintf("可解锁: %s\n", recipe.Unlock)
	msg += fmt.Sprintf("可合成: %s\n", strings.Join(recipe.Combos, ","))
	msg += fmt.Sprintf("神级符文: %s\n", recipe.Gift)
	msg += fmt.Sprintf("贵客礼物: %s\n", strings.Join(gifts, ","))
	msg += fmt.Sprintf("升阶贵客: \n%s", strings.Join(upgradeGuests, "\n"))
	return openapi.PostGroupMessageReq{
		Content: msg,
		MsgType: openapi.MsgTypeText,
	}
}

func generateRecipesMessage(ctx context.Context, recipes []model.Recipe, order string, page int) openapi.PostGroupMessageReq {
	if len(recipes) == 0 {
		return openapi.PostGroupMessageReq{
			Content: "本店没有这道菜呢！",
			MsgType: openapi.MsgTypeText,
		}
	} else if len(recipes) == 1 {
		return generateRecipeMessage(ctx, recipes[0])
	} else {
		msg := kit.PaginationOutput(recipes, page, 10, "这里有你想点的菜吗", func(recipe model.Recipe) string {
			return fmt.Sprintf("%s %s %s", recipe.GalleryId, recipe.Name, getRecipeInfoWithOrder(recipe, order))
		})
		return openapi.PostGroupMessageReq{
			Content: msg,
			MsgType: openapi.MsgTypeText,
		}
	}
}

// 根据排序参数获取菜谱需要输出的信息
func getRecipeInfoWithOrder(recipe model.Recipe, order string) string {
	switch order {
	case "单时间":
		return kit.FormatRecipeTime(recipe.Time)
	case "总时间":
		return kit.FormatRecipeTime(recipe.TotalTime)
	case "单价", "售价":
		return fmt.Sprintf("💰%d", recipe.Price)
	case "金币效率":
		return fmt.Sprintf("💰%d/h", recipe.GoldEfficiency)
	case "耗材效率":
		return fmt.Sprintf("🥗%d/h", recipe.MaterialEfficiency)
	case "稀有度":
		return strings.Repeat("🔥", recipe.Rarity)
	default:
		return ""
	}
}
