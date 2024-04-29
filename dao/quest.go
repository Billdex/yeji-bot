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
	questsMap      = make(map[string][]model.Quest)
	mainQuestsMap  = make(map[int]model.Quest) // 单独缓存主线任务，目前只提供主线查询
	mainQuestsList = make([]model.Quest, 0)
)

func ReloadQuests(ctx context.Context) ([]model.Quest, error) {
	list := make([]model.Quest, 0)
	err := DB.WithContext(ctx).Find(&list).Error
	if err != nil {
		logrus.WithContext(ctx).Errorf("db.Find(quests) fail. err: %v", err)
		return nil, errors.New("加载任务数据失败")
	}
	tmpQuestsMap := make(map[string][]model.Quest)
	tmpMainQuestsMap := make(map[int]model.Quest)
	tmpMainQuestsList := make([]model.Quest, 0)
	for i := range list {
		tmpQuestsMap[list[i].Type] = append(tmpQuestsMap[list[i].Type], list[i])
		if list[i].Type == "主线任务" {
			tmpMainQuestsMap[list[i].QuestId] = list[i]
			tmpMainQuestsList = append(tmpMainQuestsList, list[i])
		}
	}
	questsMap = tmpQuestsMap
	mainQuestsMap = tmpMainQuestsMap
	sort.Slice(tmpMainQuestsList, func(i, j int) bool {
		return tmpMainQuestsList[i].QuestId < tmpMainQuestsList[j].QuestId
	})
	mainQuestsList = tmpMainQuestsList

	return list, nil
}

// FindQuestTypeList 查询符合条件的任务类型
func FindQuestTypeList(ctx context.Context, questType string) ([]string, error) {
	pattern := strings.ReplaceAll(questType, "%", ".*")
	re, err := regexp.Compile(pattern)
	if err != nil {
		logrus.WithContext(ctx).Errorf("任务类型查询正则格式有误. err: %v", err)
		return nil, errors.New("任务类型格式有误")
	}
	questTypes := make([]string, 0, len(questsMap))
	for key := range questsMap {
		if re.MatchString(key) {
			questTypes = append(questTypes, key)
		}
	}
	return questTypes, nil
}

// GetMainQuestById 根据 id 查询主线任务
func GetMainQuestById(ctx context.Context, questId int) (model.Quest, error) {
	quest, ok := mainQuestsMap[questId]
	if !ok {
		return model.Quest{}, errors.New("任务不存在")
	}
	return quest, nil
}

func ListMainQuestsWithLimit(ctx context.Context, questId int, limit int) ([]model.Quest, error) {
	if limit < 1 {
		limit = 1
	}
	if limit > 5 {
		limit = 5
	}
	if len(mainQuestsList) == 0 {
		_, err := ReloadQuests(ctx)
		if err != nil {
			return nil, err
		}
	}
	result := make([]model.Quest, 0, limit)
	for i := range mainQuestsList {
		if mainQuestsList[i].QuestId >= questId {
			result = append(result, mainQuestsList[i])
		}
		if len(result) >= limit {
			break
		}
	}
	return result, nil
}
