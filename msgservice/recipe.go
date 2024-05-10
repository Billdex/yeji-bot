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
		order = "单时间"
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
		case kit.SliceContains([]string{"图鉴序", "时间", "单时间", "总时间", "单价", "售价", "金币效率", "耗材效率", "稀有度"}, arg):
			order = arg
		case model.IsRarityStr(arg):
			recipes = filterRecipesByRarity(ctx, recipes, model.RarityToInt(arg), strings.HasPrefix(arg, "仅"))
		case strings.HasSuffix(arg, "技法"):
			recipes = filterRecipesBySkill(ctx, recipes, strings.TrimSuffix(arg, "技法"))
		case strings.HasPrefix(arg, "技法"):
			recipes = filterRecipesBySkills(ctx, recipes, strings.Split(strings.TrimPrefix(arg, "技法"), "-"))
		case kit.SliceContains([]string{"甜味", "酸味", "辣味", "咸味", "苦味", "鲜味"}, arg):
			recipes = filterRecipesByCondiment(ctx, recipes, strings.TrimSuffix(arg, "味"))
		case kit.HasPrefixIn(arg, []string{"食材", "材料"}):
			recipes, err = filterRecipesByMaterials(ctx, recipes, strings.Split(kit.TrimPrefixIn(arg, []string{"食材", "材料"}), "-"))
		case kit.HasPrefixIn(arg, []string{"贵客", "稀有客人", "客人", "贵宾", "宾客", "稀客"}):
			recipes, err = filterRecipesByGuests(ctx, recipes, strings.Split(kit.TrimPrefixIn(arg, []string{"贵客", "稀有客人", "客人", "贵宾", "宾客", "稀客"}), "-"))
		case kit.HasPrefixIn(arg, []string{"符文", "礼物"}):
			recipes, err = filterRecipesByGifts(ctx, recipes, strings.Split(kit.TrimPrefixIn(arg, []string{"符文", "礼物"}), "-"))
		case kit.HasPrefixIn(arg, []string{"神级符文", "神级奖励"}):
			recipes, err = filterRecipesByUpgradeGift(ctx, recipes, strings.TrimPrefix(kit.TrimPrefixIn(arg, []string{"神级符文", "神级奖励"}), "-"))
		case kit.HasPrefixIn(arg, []string{"调料", "调味", "味道"}):
			recipes = filterRecipesByCondiment(ctx, recipes, strings.TrimPrefix(kit.TrimPrefixIn(arg, []string{"调料", "调味", "味道"}), "-"))
		case strings.HasPrefix(arg, "来源"):
			recipes, err = filterRecipesByOrigins(ctx, recipes, strings.Split(strings.TrimPrefix(arg, "来源"), "-"))
		case kit.HasPrefixIn(arg, []string{"$", "＄", "￥"}):
			var price int
			price, err = strconv.Atoi(arg[1:])
			if err != nil {
				err = fmt.Errorf("售价筛选格式错误")
			} else {
				recipes = filterRecipesByPrice(ctx, recipes, price)
			}
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
			return result, errors.New("菜谱查询格式有误")
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
func filterRecipesByRarity(ctx context.Context, recipes []model.Recipe, rarity int, mustEqual bool) []model.Recipe {
	if len(recipes) == 0 {
		return recipes
	}
	return kit.SliceFilter(recipes, func(recipe model.Recipe) bool {
		return (mustEqual && recipe.Rarity == rarity) || (!mustEqual && recipe.Rarity >= rarity)
	})
}

// 根据技法筛选菜谱
func filterRecipesBySkill(ctx context.Context, recipes []model.Recipe, skill string) []model.Recipe {
	if len(recipes) == 0 || skill == "" {
		return recipes
	}
	return kit.SliceFilter(recipes, func(recipe model.Recipe) bool {
		return recipe.NeedSkill(skill)
	})
}

// 根据技法列表筛选菜谱
func filterRecipesBySkills(ctx context.Context, recipes []model.Recipe, skills []string) []model.Recipe {
	if len(recipes) == 0 || len(skills) == 0 {
		return recipes
	}
	result := make([]model.Recipe, len(recipes))
	copy(result, recipes)
	for _, skill := range skills {
		if skill == "" {
			continue
		}
		result = filterRecipesBySkill(ctx, result, skill)
	}
	return result
}

// 根据调料筛选菜谱
func filterRecipesByCondiment(ctx context.Context, recipes []model.Recipe, condiment string) []model.Recipe {
	if len(recipes) == 0 || condiment == "" {
		return recipes
	}
	condimentMap := map[string]string{
		"甜": "Sweet",
		"酸": "Sour",
		"辣": "Spicy",
		"咸": "Salty",
		"苦": "Bitter",
		"鲜": "Tasty",
	}
	return kit.SliceFilter(recipes, func(recipe model.Recipe) bool {
		return recipe.Condiment == condimentMap[condiment]
	})
}

// filterRecipesByMaterials 根据食材筛选菜谱
func filterRecipesByMaterial(ctx context.Context, recipes []model.Recipe, material string) ([]model.Recipe, error) {
	if len(recipes) == 0 || material == "" {
		return recipes, nil
	}
	var materialOrigins []string
	// 符合下列特征的关键词视为根据来源筛选食材
	switch material {
	case "鱼类", "水产", "水产类", "海鲜", "海鲜类":
		materialOrigins = []string{"池塘"}
	case "蔬菜", "蔬菜类", "菜类":
		materialOrigins = []string{"菜棚", "菜地", "森林"}
	case "肉类":
		materialOrigins = []string{"牧场", "鸡舍", "猪圈"}
	case "面类", "加工", " 加工类":
		materialOrigins = []string{"作坊"}
	case "菜棚", "菜地", "森林", "牧场", "鸡舍", "猪圈", "作坊", "池塘":
		materialOrigins = []string{material}
	}
	if len(materialOrigins) > 0 {
		return kit.SliceFilter(recipes, func(recipe model.Recipe) bool {
			return recipe.HasMaterialFromIn(materialOrigins)
		}), nil
	}

	// 先查出具体食材名称
	materialNames, err := dao.MatchRecipeMaterialName(ctx, material)
	if err != nil {
		return recipes, err
	}

	if len(materialNames) == 0 {
		return nil, fmt.Errorf("厨师长说没有用%s做过菜", material)
	}
	return kit.SliceFilter(recipes, func(recipe model.Recipe) bool {
		return recipe.UsedMaterials(materialNames)
	}), nil
}

// filterRecipesByMaterials 根据食材列表筛选菜谱
func filterRecipesByMaterials(ctx context.Context, recipes []model.Recipe, materials []string) ([]model.Recipe, error) {
	if len(recipes) == 0 || len(materials) == 0 {
		return recipes, nil
	}
	result := make([]model.Recipe, len(recipes))
	copy(result, recipes)
	var err error
	for _, material := range materials {
		if material == "" {
			continue
		}
		result, err = filterRecipesByMaterial(ctx, result, material)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// filterRecipesByGuest 根据贵客筛选菜谱
func filterRecipesByGuest(ctx context.Context, recipes []model.Recipe, guest string) ([]model.Recipe, error) {
	if len(recipes) == 0 || guest == "" {
		return recipes, nil
	}

	guestNames, err := dao.MatchGuestName(ctx, guest)
	if err != nil {
		return recipes, err
	}
	if len(guestNames) == 0 {
		return nil, fmt.Errorf("%s似乎未曾光临本店呢", guest)
	}
	guestNameMap := map[string]struct{}{}
	for _, name := range guestNames {
		guestNameMap[name] = struct{}{}
	}

	return kit.SliceFilter(recipes, func(recipe model.Recipe) bool {
		for i := range recipe.GuestGifts {
			if _, ok := guestNameMap[recipe.GuestGifts[i].GuestName]; ok {
				return true
			}
		}
		return false
	}), nil
}

// filterRecipesByGuest 根据贵客列表筛选菜谱
func filterRecipesByGuests(ctx context.Context, recipes []model.Recipe, guests []string) ([]model.Recipe, error) {
	if len(recipes) == 0 || len(guests) == 0 {
		return recipes, nil
	}
	result := make([]model.Recipe, len(recipes))
	copy(result, recipes)
	var err error
	for _, guest := range guests {
		if guest == "" {
			continue
		}
		result, err = filterRecipesByGuest(ctx, result, guest)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// filterRecipesByGift 根据符文礼物筛选菜谱
func filterRecipesByGift(ctx context.Context, recipes []model.Recipe, gift string) ([]model.Recipe, error) {
	if len(recipes) == 0 || gift == "" {
		return recipes, nil
	}

	giftNames, err := dao.MatchGiftName(ctx, gift)
	if err != nil {
		return recipes, err
	}
	if len(giftNames) == 0 {
		return []model.Recipe{}, nil
	}

	giftNameMap := map[string]struct{}{}
	for _, name := range giftNames {
		giftNameMap[name] = struct{}{}
	}
	return kit.SliceFilter(recipes, func(recipe model.Recipe) bool {
		for i := range recipe.GuestGifts {
			if _, ok := giftNameMap[recipe.GuestGifts[i].Antique]; ok {
				return true
			}
		}
		return false
	}), nil
}

// filterRecipesByGift 根据符文礼物列表筛选菜谱
func filterRecipesByGifts(ctx context.Context, recipes []model.Recipe, gifts []string) ([]model.Recipe, error) {
	if len(recipes) == 0 || len(gifts) == 0 {
		return recipes, nil
	}
	result := make([]model.Recipe, len(recipes))
	copy(result, recipes)
	var err error
	for _, gift := range gifts {
		if gift == "" {
			continue
		}
		result, err = filterRecipesByGift(ctx, result, gift)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// filterRecipesByUpgradeGift 根据菜谱神级礼物筛选菜谱
func filterRecipesByUpgradeGift(ctx context.Context, recipes []model.Recipe, upgradeGift string) ([]model.Recipe, error) {
	if len(recipes) == 0 || upgradeGift == "" {
		return recipes, nil
	}
	pattern := strings.ReplaceAll(upgradeGift, "%", ".*")
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, errors.New("神级奖励查询格式有误")
	}
	return kit.SliceFilter(recipes, func(recipe model.Recipe) bool {
		return re.MatchString(recipe.Gift)
	}), nil
}

// filterRecipesByOrigins 根据来源筛选菜谱
func filterRecipesByOrigins(ctx context.Context, recipes []model.Recipe, origins []string) ([]model.Recipe, error) {
	if len(recipes) == 0 || len(origins) == 0 {
		return recipes, nil
	}
	result := make([]model.Recipe, len(recipes))
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

// filterRecipesByPrice 根据价格筛选菜谱
func filterRecipesByPrice(ctx context.Context, recipes []model.Recipe, price int) []model.Recipe {
	if len(recipes) == 0 || price == 0 {
		return recipes
	}
	return kit.SliceFilter(recipes, func(recipe model.Recipe) bool {
		return recipe.Price >= price
	})
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
	case "时间", "单时间":
		sort.Slice(recipes, func(i, j int) bool {
			return recipes[i].Time == recipes[j].Time && recipes[i].RecipeId < recipes[j].RecipeId ||
				recipes[i].Time < recipes[j].Time
		})
	case "总时间":
		sort.Slice(recipes, func(i, j int) bool {
			return recipes[i].TotalTime == recipes[j].TotalTime && recipes[i].RecipeId < recipes[j].RecipeId ||
				recipes[i].TotalTime < recipes[j].TotalTime
		})
	case "单价", "售价":
		sort.Slice(recipes, func(i, j int) bool {
			return recipes[i].Price == recipes[j].Price && recipes[i].RecipeId < recipes[j].RecipeId ||
				recipes[i].Price > recipes[j].Price
		})
	case "金币效率":
		sort.Slice(recipes, func(i, j int) bool {
			return recipes[i].GoldEfficiency == recipes[j].GoldEfficiency && recipes[i].RecipeId < recipes[j].RecipeId ||
				recipes[i].GoldEfficiency > recipes[j].GoldEfficiency
		})
	case "耗材效率":
		sort.Slice(recipes, func(i, j int) bool {
			return recipes[i].MaterialEfficiency == recipes[j].MaterialEfficiency && recipes[i].RecipeId < recipes[j].RecipeId ||
				recipes[i].MaterialEfficiency > recipes[j].MaterialEfficiency
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
