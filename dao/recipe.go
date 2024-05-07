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
	cacheRecipeList         = make([]model.Recipe, 0)
	cacheRecipeMaterialsMap = make(map[string][]model.RecipeMaterial)
)

func ReloadRecipes(ctx context.Context) ([]model.Recipe, error) {
	var recipes []model.Recipe
	err := DB.WithContext(ctx).Find(&recipes).Error
	if err != nil {
		logrus.WithContext(ctx).Errorf("db.Find(recipes) fail. err: %v", err)
		return nil, errors.New("加载菜谱数据失败")
	}

	recipeMaterialsMap := make(map[string][]model.RecipeMaterial)
	for _, recipe := range recipes {
		for _, material := range recipe.Materials {
			if material.Quantity == 0 || recipe.Time == 0 {
				continue
			}
			recipeMaterialsMap[material.MaterialName] = append(recipeMaterialsMap[material.MaterialName], model.RecipeMaterial{
				RecipeId:     recipe.RecipeId,
				RecipeName:   recipe.Name,
				MaterialId:   material.MaterialId,
				MaterialName: material.MaterialName,
				Quantity:     material.Quantity,
				Efficiency:   material.Quantity * 3600 / recipe.Time,
			})
		}
	}
	for materialName := range recipeMaterialsMap {
		recipeMaterials := recipeMaterialsMap[materialName]
		sort.Slice(recipeMaterials, func(i, j int) bool {
			return recipeMaterials[i].Efficiency > recipeMaterials[j].Efficiency
		})
		recipeMaterialsMap[materialName] = recipeMaterials
	}
	cacheRecipeMaterialsMap = recipeMaterialsMap
	cacheRecipeList = recipes
	return recipes, nil
}

func ListAllRecipes(ctx context.Context) ([]model.Recipe, error) {
	if len(cacheRecipeList) == 0 {
		_, err := ReloadRecipes(ctx)
		if err != nil {
			return nil, err
		}
	}
	return cacheRecipeList, nil
}

func MatchRecipeMaterialName(ctx context.Context, materialName string) ([]string, error) {
	if len(cacheRecipeMaterialsMap) == 0 {
		_, err := ReloadRecipes(ctx)
		if err != nil {
			return nil, err
		}
	}

	re, err := regexp.Compile(strings.ReplaceAll(materialName, "%", ".*"))
	if err != nil {
		logrus.WithContext(ctx).Errorf("regexp compile fail. raw str: %s, err: %v", materialName, err)
		return nil, errors.New("查询格式有误")
	}
	names := make([]string, 0, 10)
	for material := range cacheRecipeMaterialsMap {
		// 如果完全匹配则直接返回
		if material == materialName {
			return []string{materialName}, nil
		}
		if re.MatchString(material) {
			names = append(names, material)
		}
	}
	sort.Strings(names)
	return names, nil
}

func ListRecipeMaterials(ctx context.Context, materialName string) ([]model.RecipeMaterial, error) {
	if len(cacheRecipeMaterialsMap) == 0 {
		_, err := ReloadRecipes(ctx)
		if err != nil {
			return nil, err
		}
	}
	recipeMaterials, has := cacheRecipeMaterialsMap[materialName]
	if !has || len(recipeMaterials) == 0 {
		return nil, errors.New("没有找到对应食材数据")
	}
	return recipeMaterials, nil
}
