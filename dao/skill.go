package dao

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"yeji-bot/data/model"
)

var (
	cacheSkillsMap = make(map[int]model.Skill)
)

func ReloadSkills(ctx context.Context) error {
	list := make([]model.Skill, 0)
	err := DB.WithContext(ctx).Find(&list).Error
	if err != nil {
		logrus.WithContext(ctx).Errorf("db.Find(skills) fail. err: %v", err)
		return errors.New("加载技能数据失败")
	}
	tmpSkillsMap := make(map[int]model.Skill)
	for _, skill := range list {
		tmpSkillsMap[skill.SkillId] = skill
	}
	cacheSkillsMap = tmpSkillsMap
	return nil
}

func GetSkillsMapByIds(ctx context.Context, skillIds []int) (map[int]model.Skill, error) {
	if len(cacheSkillsMap) == 0 {
		err := ReloadSkills(ctx)
		if err != nil {
			return nil, err
		}
	}
	if len(skillIds) == 0 {
		m := make(map[int]model.Skill, len(skillIds))
		for id := range cacheSkillsMap {
			m[id] = cacheSkillsMap[id]
		}
		return m, nil
	}
	m := make(map[int]model.Skill, len(skillIds))
	for _, skillId := range skillIds {
		skill, ok := cacheSkillsMap[skillId]
		if !ok {
			return nil, errors.New("技能不存在")
		}
		m[skillId] = skill
	}
	return m, nil
}

func GetSkillById(ctx context.Context, skillId int) (model.Skill, error) {
	if len(cacheSkillsMap) == 0 {
		err := ReloadSkills(ctx)
		if err != nil {
			return model.Skill{}, err
		}
	}
	skill, ok := cacheSkillsMap[skillId]
	if !ok {
		return model.Skill{}, errors.New("技能不存在")
	}
	return skill, nil
}
