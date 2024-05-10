package dao

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"yeji-bot/data/model"
)

var (
	cacheChefList = make([]model.Chef, 0)
)

func ReloadChefs(ctx context.Context) error {
	list := make([]model.Chef, 0)
	err := DB.WithContext(ctx).Find(&list).Error
	if err != nil {
		logrus.WithContext(ctx).Errorf("db.Find(chefs) fail. err: %v", err)
		return errors.New("加载厨师数据失败")
	}
	cacheChefList = list
	return nil
}

func ListAllChefs(ctx context.Context) ([]model.Chef, error) {
	if len(cacheChefList) == 0 {
		err := ReloadChefs(ctx)
		if err != nil {
			return nil, err
		}
	}
	return cacheChefList, nil
}
