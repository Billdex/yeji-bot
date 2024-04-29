package msgservice

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
	"yeji-bot/bot/openapi"
	"yeji-bot/bot/qbot"
	"yeji-bot/dao"
	"yeji-bot/data/gamedata"
	"yeji-bot/data/model"
	"yeji-bot/pkg/seq"
)

const (
	foodGameUrlbase = "https://foodgame.github.io"
	bcjhUrlBase     = "https://www.bcjh.xyz"

	dataURI = "/data/data.min.json"
)

var (
	updateState = false
	updateMux   sync.Mutex
)

// UpdateData 更新游戏数据
func UpdateData(ctx context.Context, api *openapi.Openapi, msg *qbot.WSGroupAtMessageData) (err error) {
	updateMux.Lock()
	if updateState {
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: "数据正在更新中",
			MsgType: openapi.MsgTypeText,
			MsgId:   msg.Id,
			MsgSeq:  seq.Seq(ctx),
		})
		if err != nil {
			logrus.WithContext(ctx).Errorf("post group message fail. err: %v", err)
		}
		updateMux.Unlock()
		return
	}
	updateState = true
	updateMux.Unlock()
	defer func() { updateState = false }()

	var baseURL string
	switch msg.Content {
	case "图鉴网":
		baseURL = foodGameUrlbase
	case "白菜菊花":
		baseURL = bcjhUrlBase
	default:
		if strings.HasPrefix(msg.Content, "http://") || strings.HasPrefix(msg.Content, "https://") {
			baseURL = msg.Content
		} else {
			baseURL = bcjhUrlBase
		}
	}
	_, _ = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
		Content: "开始导入数据",
		MsgType: openapi.MsgTypeText,
		MsgId:   msg.Id,
		MsgSeq:  seq.Seq(ctx),
	})
	updateStart := time.Now()
	contentMsg := ""
	// 获取图鉴网数据
	stepStart := time.Now()
	gameData, err := requestData(baseURL)
	if err != nil {
		logrus.WithContext(ctx).Errorf("requestData fail. err: %v", err)
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: "获取数据失败",
			MsgType: openapi.MsgTypeText,
			MsgId:   msg.Id,
			MsgSeq:  seq.Seq(ctx),
		})
		if err != nil {
			logrus.WithContext(ctx).Errorf("PostGroupMessage fail. err: %v", err)
		}
		return nil
	}
	stepTime := time.Since(stepStart).Round(time.Millisecond).String()
	logrus.WithContext(ctx).Infof("获取图鉴网数据完毕，耗时: %s", stepTime)
	contentMsg += fmt.Sprintf("获取图鉴网数据耗时: %s\n", stepTime)

	// 更新厨师数据
	stepStart = time.Now()
	err = updateChefs(ctx, gameData.Chefs)
	if err != nil {
		logrus.WithContext(ctx).Errorf("更新厨师数据失败. err: %v", err)
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: "更新厨师数据失败",
			MsgType: openapi.MsgTypeText,
			MsgId:   msg.Id,
			MsgSeq:  seq.Seq(ctx),
		})
		if err != nil {
			logrus.WithContext(ctx).Errorf("PostGroupMessage fail. err: %v", err)
		}
		return nil
	}
	stepTime = time.Since(stepStart).Round(time.Millisecond).String()
	logrus.WithContext(ctx).Infof("更新厨师数据完毕，耗时: %s", stepTime)
	contentMsg += fmt.Sprintf("更新厨师数据耗时: %s\n", stepTime)

	// 更新菜谱数据
	stepStart = time.Now()
	err = updateRecipes(ctx, gameData.Recipes, gameData.Materials, gameData.Combos, gameData.Guests)
	if err != nil {
		logrus.WithContext(ctx).Errorf("更新菜谱数据失败. err: %v", err)
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: "更新菜谱数据失败",
			MsgType: openapi.MsgTypeText,
			MsgId:   msg.Id,
			MsgSeq:  seq.Seq(ctx),
		})
		if err != nil {
			logrus.WithContext(ctx).Errorf("PostGroupMessage fail. err: %v", err)
		}
		return nil
	}
	stepTime = time.Since(stepStart).Round(time.Millisecond).String()
	logrus.WithContext(ctx).Infof("更新菜谱数据完毕，耗时: %s", stepTime)
	contentMsg += fmt.Sprintf("更新菜谱数据耗时: %s\n", stepTime)

	// 更新厨具数据
	stepStart = time.Now()
	err = updateEquips(ctx, gameData.Equips)
	if err != nil {
		logrus.WithContext(ctx).Errorf("更新厨具数据失败. err: %v", err)
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: "更新厨具数据失败",
			MsgType: openapi.MsgTypeText,
			MsgId:   msg.Id,
			MsgSeq:  seq.Seq(ctx),
		})
		if err != nil {
			logrus.WithContext(ctx).Errorf("PostGroupMessage fail. err: %v", err)
		}
		return nil
	}
	stepTime = time.Since(stepStart).Round(time.Millisecond).String()
	logrus.WithContext(ctx).Infof("更新厨具数据完毕，耗时: %s", stepTime)
	contentMsg += fmt.Sprintf("更新厨具数据耗时: %s\n", stepTime)

	// 更新食材数据
	stepStart = time.Now()
	err = updateMaterials(ctx, gameData.Materials)
	if err != nil {
		logrus.WithContext(ctx).Errorf("更新食材数据失败. err: %v", err)
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: "更新食材数据失败",
			MsgType: openapi.MsgTypeText,
			MsgId:   msg.Id,
			MsgSeq:  seq.Seq(ctx),
		})
		if err != nil {
			logrus.WithContext(ctx).Errorf("PostGroupMessage fail. err: %v", err)
		}
		return nil
	}
	stepTime = time.Since(stepStart).Round(time.Millisecond).String()
	logrus.WithContext(ctx).Infof("更新食材数据完毕，耗时: %s", stepTime)
	contentMsg += fmt.Sprintf("更新食材数据耗时: %s\n", stepTime)

	// 更新技能信息
	stepStart = time.Now()
	err = updateSkills(ctx, gameData.Skills)
	if err != nil {
		logrus.WithContext(ctx).Errorf("更新技能数据失败. err: %v", err)
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: "更新技能数据失败",
			MsgType: openapi.MsgTypeText,
			MsgId:   msg.Id,
			MsgSeq:  seq.Seq(ctx),
		})
		if err != nil {
			logrus.WithContext(ctx).Errorf("PostGroupMessage fail. err: %v", err)
		}
		return nil
	}
	stepTime = time.Since(stepStart).Round(time.Millisecond).String()
	logrus.WithContext(ctx).Infof("更新技能数据完毕，耗时: %s", stepTime)
	contentMsg += fmt.Sprintf("更新技能数据耗时: %s\n", stepTime)

	// 更新任务数据
	stepStart = time.Now()
	err = updateQuests(ctx, gameData.Quests)
	if err != nil {
		logrus.WithContext(ctx).Errorf("更新任务数据失败. err: %v", err)
		_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
			Content: "更新任务数据失败",
			MsgType: openapi.MsgTypeText,
			MsgId:   msg.Id,
			MsgSeq:  seq.Seq(ctx),
		})
		if err != nil {
			logrus.WithContext(ctx).Errorf("PostGroupMessage fail. err: %v", err)
		}
		return nil
	}
	stepTime = time.Since(stepStart).Round(time.Millisecond).String()
	logrus.WithContext(ctx).Infof("更新任务数据完毕，耗时: %s", stepTime)
	contentMsg += fmt.Sprintf("更新任务数据耗时: %s\n", stepTime)

	logrus.WithContext(ctx).Infof("更新数据完毕，耗时: %s", time.Since(updateStart).Round(time.Millisecond).String())
	contentMsg = strings.TrimSuffix(fmt.Sprintf("更新数据完毕，累计耗时: %s\n%s", time.Since(updateStart).Round(time.Millisecond).String(), contentMsg), "\n")
	_, err = api.PostGroupMessage(ctx, msg.GroupOpenid, &openapi.PostGroupMessageReq{
		Content: contentMsg,
		MsgType: openapi.MsgTypeText,
		MsgId:   msg.Id,
		MsgSeq:  seq.Seq(ctx),
	})
	if err != nil {
		logrus.WithContext(ctx).Errorf("PostGroupMessage fail. err: %v", err)
	}

	return nil
}

// 从图鉴网爬取数据
func requestData(baseURL string) (gamedata.GameData, error) {
	var gameData gamedata.GameData
	r, err := http.Get(baseURL + dataURI)
	if err != nil {
		return gameData, err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return gameData, err
	}
	err = json.Unmarshal(body, &gameData)
	return gameData, err
}

// updateChefs 更新厨师数据
func updateChefs(ctx context.Context, chefsData []gamedata.ChefData) error {
	if len(chefsData) == 0 {
		return errors.New("厨师数据为空")
	}
	chefs := make([]model.Chef, 0)
	for _, chefData := range chefsData {
		chef := model.Chef{
			ChefId:        chefData.ChefId,
			Name:          chefData.Name,
			Rarity:        chefData.Rarity,
			Origins:       strings.Split(chefData.Origin, "<br>"),
			GalleryId:     chefData.GalleryId,
			Stirfry:       chefData.Stirfry,
			Bake:          chefData.Bake,
			Boil:          chefData.Boil,
			Steam:         chefData.Steam,
			Fry:           chefData.Fry,
			Cut:           chefData.Cut,
			Meat:          chefData.Meat,
			Flour:         chefData.Flour,
			Fish:          chefData.Fish,
			Vegetable:     chefData.Vegetable,
			Sweet:         chefData.Sweet,
			Sour:          chefData.Sour,
			Spicy:         chefData.Spicy,
			Salty:         chefData.Salty,
			Bitter:        chefData.Bitter,
			Tasty:         chefData.Tasty,
			SkillId:       chefData.SkillId,
			UltimateGoals: chefData.UltimateGoals,
			UltimateSkill: chefData.UltimateSkill,
		}
		if len(chefData.Tags) > 0 {
			chef.Gender = chefData.Tags[0]
		}
		chefs = append(chefs, chef)
	}
	err := dao.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {
		err = tx.Exec(`truncate table chef`).Error
		if err != nil {
			logrus.WithContext(ctx).Errorf("try truncate table chef fail. err: %v", err)
			return err
		}
		err = tx.Create(&chefs).Error
		if err != nil {
			logrus.WithContext(ctx).Errorf("try insert all chefs data fail. err: %v", err)
			return err
		}

		return nil
	})

	return err
}

// 更新菜谱信息
func updateRecipes(ctx context.Context, recipesData []gamedata.RecipeData, materialsData []gamedata.MaterialData, combosData []gamedata.ComboData, guestsData []gamedata.GuestData) error {
	if len(recipesData) == 0 {
		return errors.New("菜谱数据为空")
	}

	mMaterials := make(map[int]gamedata.MaterialData)
	for i := range materialsData {
		mMaterials[materialsData[i].MaterialId] = materialsData[i]
	}
	mIdToNameCombo := make(map[int]struct {
		Name   string
		Combos []string
	}, len(combosData))
	mRecipeGuestsData := make(map[string][]model.RecipeGuestGift)
	for i := range recipesData {
		mIdToNameCombo[recipesData[i].RecipeId] = struct {
			Name   string
			Combos []string
		}{Name: recipesData[i].Name, Combos: []string{}}
		mRecipeGuestsData[recipesData[i].Name] = []model.RecipeGuestGift{}
	}
	// 预处理后厨合成菜数据
	for _, combo := range combosData {
		for _, recipeId := range combo.Recipes {
			nameComboData := mIdToNameCombo[recipeId]
			nameComboData.Combos = append(nameComboData.Combos, mIdToNameCombo[combo.RecipeId].Name)
			mIdToNameCombo[recipeId] = nameComboData
		}
	}
	// 预处理贵客礼物数据
	for _, guestData := range guestsData {
		for _, gift := range guestData.Gifts {
			mRecipeGuestsData[gift.Recipe] = append(mRecipeGuestsData[gift.Recipe], model.RecipeGuestGift{
				GuestName: guestData.Name,
				Antique:   gift.Antique,
			})
		}
	}

	// 生成菜谱数据
	recipes := make([]model.Recipe, 0, len(recipesData))
	for _, recipeData := range recipesData {
		guests := make([]string, 0, len(recipeData.Guests))
		for _, guest := range recipeData.Guests {
			guests = append(guests, guest.Guest)
		}
		materials := make([]model.RecipeMaterial, 0, len(recipeData.Materials))
		materialSum := 0
		for _, materialData := range recipeData.Materials {
			materials = append(materials, model.RecipeMaterial{
				MaterialId:   materialData.MaterialId,
				MaterialName: mMaterials[materialData.MaterialId].Name,
				Quantity:     materialData.Quantity,
			})
			materialSum += materialData.Quantity
		}
		recipe := model.Recipe{
			RecipeId:           recipeData.RecipeId,
			Name:               recipeData.Name,
			GalleryId:          recipeData.GalleryId,
			Rarity:             recipeData.Rarity,
			Origins:            strings.Split(recipeData.Origin, "<br>"),
			Stirfry:            recipeData.Stirfry,
			Bake:               recipeData.Bake,
			Boil:               recipeData.Boil,
			Steam:              recipeData.Steam,
			Fry:                recipeData.Fry,
			Cut:                recipeData.Cut,
			Condiment:          recipeData.Condiment,
			Price:              recipeData.Price,
			ExPrice:            recipeData.ExPrice,
			GoldEfficiency:     recipeData.Price * 3600 / recipeData.Time,
			MaterialEfficiency: materialSum * 3600 / recipeData.Time,
			Gift:               recipeData.Gift,
			GuestGifts:         mRecipeGuestsData[recipeData.Name],
			UpgradeGuests:      guests,
			Time:               recipeData.Time,
			Limit:              recipeData.Limit,
			TotalTime:          recipeData.Time * recipeData.Limit,
			Unlock:             recipeData.Unlock,
			Combos:             mIdToNameCombo[recipeData.RecipeId].Combos,
			Materials:          materials,
		}
		recipes = append(recipes, recipe)
	}
	err := dao.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {
		err = tx.Exec(`truncate table recipe`).Error
		if err != nil {
			logrus.WithContext(ctx).Errorf("try truncate table recipe fail. err: %v", err)
			return err
		}
		err = tx.Create(&recipes).Error
		if err != nil {
			logrus.WithContext(ctx).Errorf("try insert all recipes data fail. err: %v", err)
			return err
		}
		return nil
	})

	return err
}

// updateMaterials 更新食材数据
func updateMaterials(ctx context.Context, materialsData []gamedata.MaterialData) error {
	if len(materialsData) == 0 {
		return errors.New("食材数据为空")
	}
	materials := make([]model.Material, 0, len(materialsData))
	for _, materialData := range materialsData {
		materials = append(materials, model.Material{
			MaterialId: materialData.MaterialId,
			Name:       materialData.Name,
			Origin:     materialData.Origin,
		})
	}
	err := dao.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {
		err = tx.Exec(`truncate table material`).Error
		if err != nil {
			logrus.WithContext(ctx).Errorf("try truncate table material fail. err: %v", err)
			return err
		}
		err = tx.Create(&materials).Error
		if err != nil {
			logrus.WithContext(ctx).Errorf("try insert all materials data fail. err: %v", err)
			return err
		}
		return nil
	})

	return err
}

// 更新厨具信息
func updateEquips(ctx context.Context, equipsData []gamedata.EquipData) error {
	if len(equipsData) == 0 {
		return errors.New("厨具数据为空")
	}
	equips := make([]model.Equip, 0, len(equipsData))
	for _, equipData := range equipsData {
		equips = append(equips, model.Equip{
			EquipId:   equipData.EquipId,
			Name:      equipData.Name,
			GalleryId: equipData.GalleryId,
			Origins:   strings.Split(equipData.Origin, "<br>"),
			Rarity:    equipData.Rarity,
			Skills:    equipData.Skills,
		})
	}
	err := dao.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {
		err = tx.Exec(`truncate table equip`).Error
		if err != nil {
			logrus.WithContext(ctx).Errorf("try truncate table equip fail. err: %v", err)
			return err
		}
		err = tx.Create(&equips).Error
		if err != nil {
			logrus.WithContext(ctx).Errorf("try insert all equips data fail. err: %v", err)
			return err
		}
		return nil
	})
	return err
}

// updateSkills 更新技能数据
func updateSkills(ctx context.Context, skillsData []gamedata.SkillData) error {
	if len(skillsData) == 0 {
		return errors.New("技能数据为空")
	}
	skills := make([]model.Skill, 0, len(skillsData))
	for _, skillData := range skillsData {
		skill := model.Skill{
			SkillId:     skillData.SkillId,
			Description: strings.ReplaceAll(skillData.Description, "<br>", ","),
		}
		effects := make([]model.SkillEffect, 0, len(skillData.Effects))
		for _, effectData := range skillData.Effects {
			effects = append(effects, model.SkillEffect{
				Calculation: effectData.Calculation,
				Type:        effectData.Type,
				Condition:   effectData.Condition,
				Tag:         effectData.Tag,
				Value:       effectData.Value,
			})
		}
		skill.Effects = effects
		skills = append(skills, skill)
	}

	err := dao.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {
		err = tx.Exec(`truncate table skill`).Error
		if err != nil {
			logrus.WithContext(ctx).Errorf("try truncate table skill fail. err: %v", err)
			return err
		}
		err = tx.Create(&skills).Error
		if err != nil {
			logrus.WithContext(ctx).Errorf("try insert all skills data fail. err: %v", err)
			return err
		}
		return nil
	})

	return err
}

// 更新任务信息
func updateQuests(ctx context.Context, questsData []gamedata.QuestData) error {
	if len(questsData) == 0 {
		return errors.New("任务数据为空")
	}
	quests := make([]model.Quest, 0, len(questsData))
	for _, questData := range questsData {
		rewards := make([]model.QuestReward, 0, len(questData.Rewards))
		for _, rewardData := range questData.Rewards {
			rewards = append(rewards, model.QuestReward{
				Name:     rewardData.Name,
				Quantity: rewardData.Quantity,
			})
		}
		conditions := make([]model.QuestCondition, 0)
		for _, conditionData := range questData.Conditions {
			conditions = append(conditions, model.QuestCondition{
				RecipeId:     conditionData.RecipeId,
				Rank:         conditionData.Rank,
				Num:          conditionData.Num,
				GoldEff:      conditionData.GoldEff,
				MaterialId:   conditionData.MaterialId,
				Guest:        conditionData.Guest,
				AnyGuest:     conditionData.AnyGuest,
				Skill:        conditionData.Skill,
				MaterialEff:  conditionData.MaterialEff,
				NewGuest:     conditionData.NewGuest,
				Rarity:       conditionData.Rarity,
				Price:        conditionData.Price,
				Category:     conditionData.Category,
				Condiment:    conditionData.Condiment,
				CondimentEff: conditionData.CondimentEff,
			})
		}
		quest := model.Quest{
			QuestId:     questData.QuestId,
			QuestIdDisp: decimal.NewFromFloatWithExponent(questData.QuestIdDisp, 2),
			Type:        questData.Type,
			Goal:        questData.Goal,
			Rewards:     rewards,
			Conditions:  conditions,
		}
		quests = append(quests, quest)
	}

	err := dao.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {
		err = tx.Exec(`truncate table quest`).Error
		if err != nil {
			logrus.WithContext(ctx).Errorf("try truncate table quest fail. err: %v", err)
			return err
		}
		err = tx.Create(&quests).Error
		if err != nil {
			logrus.WithContext(ctx).Errorf("try insert all quest data fail. err: %v", err)
			return err
		}
		return nil
	})

	return err
}
