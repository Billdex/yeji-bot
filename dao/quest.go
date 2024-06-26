package dao

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"regexp"
	"sort"
	"strings"
	"yeji-bot/data/model"
)

var (
	cacheQuestsMap      = make(map[int]model.Quest)
	cacheQuestTypesMap  = make(map[string][]model.Quest)
	cacheMainQuestsMap  = make(map[int]model.Quest) // 单独缓存主线任务，目前只提供主线查询
	cacheMainQuestsList = make([]model.Quest, 0)
)

func ReloadQuests(ctx context.Context) error {
	list := make([]model.Quest, 0)
	err := DB.WithContext(ctx).Find(&list).Error
	if err != nil {
		logrus.WithContext(ctx).Errorf("db.Find(quests) fail. err: %v", err)
		return errors.New("加载任务数据失败")
	}
	tmpQuestsMap := make(map[int]model.Quest)
	tmpQuestTypesMap := make(map[string][]model.Quest)
	tmpMainQuestsMap := make(map[int]model.Quest)
	tmpMainQuestsList := make([]model.Quest, 0)
	for i := range list {
		tmpQuestsMap[list[i].QuestId] = list[i]
		tmpQuestTypesMap[list[i].Type] = append(tmpQuestTypesMap[list[i].Type], list[i])
		if list[i].Type == "主线任务" {
			tmpMainQuestsMap[list[i].QuestId] = list[i]
			tmpMainQuestsList = append(tmpMainQuestsList, list[i])
		}
	}
	cacheQuestsMap = tmpQuestsMap
	cacheQuestTypesMap = tmpQuestTypesMap
	cacheMainQuestsMap = tmpMainQuestsMap
	sort.Slice(tmpMainQuestsList, func(i, j int) bool {
		return tmpMainQuestsList[i].QuestId < tmpMainQuestsList[j].QuestId
	})
	cacheMainQuestsList = tmpMainQuestsList

	return nil
}

// FindQuestTypeList 查询符合条件的任务类型
func FindQuestTypeList(ctx context.Context, questType string) ([]string, error) {
	pattern := strings.ReplaceAll(questType, "%", ".*")
	re, err := regexp.Compile(pattern)
	if err != nil {
		logrus.WithContext(ctx).Errorf("任务类型查询正则格式有误. err: %v", err)
		return nil, errors.New("任务类型格式有误")
	}
	questTypes := make([]string, 0, len(cacheQuestTypesMap))
	for key := range cacheQuestTypesMap {
		if re.MatchString(key) {
			questTypes = append(questTypes, key)
		}
	}
	return questTypes, nil
}

// GetMainQuestById 根据 id 查询主线任务
func GetMainQuestById(ctx context.Context, questId int) (model.Quest, error) {
	quest, ok := cacheMainQuestsMap[questId]
	if !ok {
		return model.Quest{}, errors.New("任务不存在")
	}
	return quest, nil
}

func FindQuestsByIds(ctx context.Context, questIds []int) ([]model.Quest, error) {
	if len(cacheQuestsMap) == 0 {
		err := ReloadQuests(ctx)
		if err != nil {
			return nil, err
		}
	}
	quests := make([]model.Quest, 0)
	for _, questId := range questIds {
		quest, ok := cacheQuestsMap[questId]
		if !ok {
			return nil, errors.New("任务不存在")
		}
		quests = append(quests, quest)
	}
	return quests, nil
}

func ListMainQuestsWithLimit(ctx context.Context, questId int, limit int) ([]model.Quest, error) {
	if limit < 1 {
		limit = 1
	}
	if limit > 5 {
		limit = 5
	}
	if len(cacheMainQuestsList) == 0 {
		err := ReloadQuests(ctx)
		if err != nil {
			return nil, err
		}
	}
	result := make([]model.Quest, 0, limit)
	for i := range cacheMainQuestsList {
		if cacheMainQuestsList[i].QuestId >= questId {
			result = append(result, cacheMainQuestsList[i])
		}
		if len(result) >= limit {
			break
		}
	}
	return result, nil
}
