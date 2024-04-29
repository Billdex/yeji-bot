package dao

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"yeji-bot/data/model"
)

var (
	skillsMap = make(map[int]model.Skill)
)

func ReloadSkills(ctx context.Context) ([]model.Skill, error) {
	list := make([]model.Skill, 0)
	err := DB.WithContext(ctx).Find(&list).Error
	if err != nil {
		logrus.WithContext(ctx).Errorf("db.Find(skills) fail. err: %v", err)
		return nil, errors.New("加载技能数据失败")
	}
	tmpSkillsMap := make(map[int]model.Skill)
	for _, skill := range list {
		tmpSkillsMap[skill.SkillId] = skill
	}
	skillsMap = tmpSkillsMap
	return list, nil
}

func GetSkillsMapByIds(ctx context.Context, skillIds []int) (map[int]model.Skill, error) {
	if len(skillsMap) == 0 {
		_, err := ReloadSkills(ctx)
		if err != nil {
			return nil, err
		}
	}
	m := make(map[int]model.Skill)
	for _, skillId := range skillIds {
		skill, ok := skillsMap[skillId]
		if !ok {
			return nil, errors.New("技能不存在")
		}
		m[skillId] = skill
	}
	return m, nil
}

func GetSkillById(ctx context.Context, skillId int) (model.Skill, error) {
	if len(skillsMap) == 0 {
		_, err := ReloadSkills(ctx)
		if err != nil {
			return model.Skill{}, err
		}
	}
	skill, ok := skillsMap[skillId]
	if !ok {
		return model.Skill{}, errors.New("技能不存在")
	}
	return skill, nil
}
