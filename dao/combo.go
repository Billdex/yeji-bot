package dao

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"yeji-bot/data/model"
	"yeji-bot/pkg/kit"
)

var (
	cacheComboRecipeList = make([]model.ComboRecipe, 0)
)

func ReloadComboRecipes(ctx context.Context) error {
	var comboRecipes []model.ComboRecipe
	err := DB.WithContext(ctx).Find(&comboRecipes).Error
	if err != nil {
		logrus.WithContext(ctx).Errorf("db.Find(combo_recipe) fail. err: %v", err)
		return errors.New("加载合成菜谱数据失败")
	}

	allRecipes, err := ListAllRecipes(ctx)
	if err != nil {
		logrus.WithContext(ctx).Errorf("ListAllRecipes() fail. err: %v", err)
		return err
	}
	recipeMap := kit.SliceToMap(allRecipes, func(recipe model.Recipe) int {
		return recipe.RecipeId
	})
	for i := range comboRecipes {
		for _, needRecipeId := range comboRecipes[i].NeedRecipeIds {
			needRecipe, ok := recipeMap[needRecipeId]
			if !ok {
				continue
			}
			comboRecipes[i].NeedRecipes = append(comboRecipes[i].NeedRecipes, needRecipe)
		}
	}
	cacheComboRecipeList = comboRecipes
	return nil
}

func ListAllComboRecipes(ctx context.Context) ([]model.ComboRecipe, error) {
	if len(cacheComboRecipeList) == 0 {
		err := ReloadComboRecipes(ctx)
		if err != nil {
			return nil, err
		}
	}
	return cacheComboRecipeList, nil
}
