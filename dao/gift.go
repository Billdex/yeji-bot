package dao

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"regexp"
	"sort"
	"strings"
	"yeji-bot/data/model"
)

var (
	cacheAntiqueGiftsMap = make(map[string][]model.Gift)
	cacheGuestGiftsMap   = make(map[string][]model.Gift)
)

func ReloadGifts(ctx context.Context) error {
	recipes, err := ListAllRecipes(ctx)
	if err != nil {
		return err
	}
	antiqueGiftsMap := make(map[string][]model.Gift)
	guestGiftsMap := make(map[string][]model.Gift)
	for _, recipe := range recipes {
		for _, gift := range recipe.GuestGifts {
			antiqueGiftsMap[gift.Antique] = append(antiqueGiftsMap[gift.Antique], model.Gift{
				RecipeId:        recipe.RecipeId,
				RecipeName:      recipe.Name,
				RecipeTotalTime: recipe.TotalTime,
				GuestName:       gift.GuestName,
				Antique:         gift.Antique,
			})
			guestGiftsMap[gift.GuestName] = append(guestGiftsMap[gift.GuestName], model.Gift{
				RecipeId:        recipe.RecipeId,
				RecipeName:      recipe.Name,
				RecipeTotalTime: recipe.TotalTime,
				GuestName:       gift.GuestName,
				Antique:         gift.Antique,
			})
		}
	}
	var sortFn = func(dataMap map[string][]model.Gift) {
		for key := range dataMap {
			gifts := dataMap[key]
			sort.Slice(dataMap[key], func(i, j int) bool {
				if gifts[i].RecipeTotalTime == gifts[j].RecipeTotalTime {
					return gifts[i].RecipeId < gifts[j].RecipeId
				}
				return gifts[i].RecipeTotalTime < gifts[j].RecipeTotalTime
			})
			dataMap[key] = gifts
		}
	}
	sortFn(antiqueGiftsMap)
	sortFn(guestGiftsMap)

	cacheAntiqueGiftsMap = antiqueGiftsMap
	cacheGuestGiftsMap = guestGiftsMap
	return nil
}

func MatchGuestName(ctx context.Context, guestName string) ([]string, error) {
	if len(cacheGuestGiftsMap) == 0 {
		err := ReloadGifts(ctx)
		if err != nil {
			return nil, err
		}
	}
	re, err := regexp.Compile(strings.ReplaceAll(guestName, "%", ".*"))
	if err != nil {
		logrus.WithContext(ctx).Errorf("regexp compile fail. raw str: %s, err: %v", guestName, err)
		return nil, errors.New("贵客查询格式有误")
	}
	names := make([]string, 0, 10)
	for guest := range cacheGuestGiftsMap {
		// 如果完全匹配则直接返回
		if guest == guestName {
			return []string{guestName}, nil
		}
		if re.MatchString(guest) {
			names = append(names, guest)
		}
	}
	sort.Strings(names)
	return names, nil
}

func ListGuestGifts(ctx context.Context, guestName string) ([]model.Gift, error) {
	if len(cacheGuestGiftsMap) == 0 {
		err := ReloadGifts(ctx)
		if err != nil {
			return nil, err
		}
	}
	gifts, has := cacheGuestGiftsMap[guestName]
	if !has || len(gifts) == 0 {
		return nil, errors.New("没有找到贵客数据")
	}
	return gifts, nil
}

// MatchGiftName 匹配贵客礼物名称
func MatchGiftName(ctx context.Context, giftName string) ([]string, error) {
	if len(cacheAntiqueGiftsMap) == 0 {
		err := ReloadGifts(ctx)
		if err != nil {
			return nil, err
		}
	}
	re, err := regexp.Compile(strings.ReplaceAll(giftName, "%", ".*"))
	if err != nil {
		logrus.WithContext(ctx).Errorf("regexp compile fail. raw str: %s, err: %v", giftName, err)
		return nil, errors.New("符文礼物查询格式有误")
	}
	names := make([]string, 0, 10)
	for antique := range cacheAntiqueGiftsMap {
		// 如果完全匹配则直接返回
		if antique == giftName {
			return []string{giftName}, nil
		}
		if re.MatchString(antique) {
			names = append(names, antique)
		}
	}
	sort.Strings(names)
	return names, nil
}

func ListAntiqueGifts(ctx context.Context, giftName string) ([]model.Gift, error) {
	if len(cacheAntiqueGiftsMap) == 0 {
		err := ReloadGifts(ctx)
		if err != nil {
			return nil, err
		}
	}
	gifts, has := cacheAntiqueGiftsMap[giftName]
	if !has || len(gifts) == 0 {
		return nil, errors.New("没有找到符文礼物数据")
	}
	return gifts, nil
}
