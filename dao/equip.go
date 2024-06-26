package dao

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"yeji-bot/data/model"
	"yeji-bot/pkg/kit"
)

var (
	cacheEquipList = make([]model.Equip, 0)
)

func ReloadEquips(ctx context.Context) error {
	tmpEquips := make([]model.Equip, len(cacheEquipList))
	err := DB.WithContext(ctx).Find(&tmpEquips).Error
	if err != nil {
		logrus.WithContext(ctx).Errorf("db.Find(equips) fail. err: %v", err)
		return errors.New("加载厨具数据失败")
	}

	skillMap, err := GetSkillsMapByIds(ctx, nil)
	if err != nil {
		return err
	}
	for i := range tmpEquips {
		tmpEquips[i].SkillDescs = kit.SliceMap(tmpEquips[i].Skills, func(skillId int) string {
			return skillMap[skillId].Description
		})
	}

	cacheEquipList = tmpEquips
	return nil
}

func ListAllEquips(ctx context.Context) ([]model.Equip, error) {
	if len(cacheEquipList) == 0 {
		err := ReloadEquips(ctx)
		if err != nil {
			logrus.WithContext(ctx).Errorf("dao.ReloadEquips() fail. err: %v", err)
			return nil, err
		}
	}
	return cacheEquipList, nil
}
