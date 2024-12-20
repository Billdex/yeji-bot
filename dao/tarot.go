package dao

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"yeji-bot/data/model"
)

var (
	cacheTarotList = make([]model.Tarot, 0)
)

func ReloadTarots(ctx context.Context) error {
	list := make([]model.Tarot, 0)
	err := DB.WithContext(ctx).Find(&list).Error
	if err != nil {
		logrus.WithContext(ctx).Errorf("db.Find(tarot) fail. err: %v", err)
		return errors.New("加载签文数据失败")
	}
	cacheTarotList = list
	return nil
}

func ListAllTarots(ctx context.Context) ([]model.Tarot, error) {
	if len(cacheTarotList) == 0 {
		err := ReloadTarots(ctx)
		if err != nil {
			return nil, err
		}
	}
	return cacheTarotList, nil
}
