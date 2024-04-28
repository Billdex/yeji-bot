package msgservice

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
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
	err = updateRecipes(ctx, gameData.Recipes, gameData.Combos, gameData.Guests)
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
func updateRecipes(ctx context.Context, recipesData []gamedata.RecipeData, combosData []gamedata.ComboData, guestsData []gamedata.GuestData) error {
	if len(recipesData) == 0 {
		return errors.New("菜谱数据为空")
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
				MaterialId: materialData.MaterialId,
				Quantity:   materialData.Quantity,
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

// 更新厨具信息
func updateEquips(ctx context.Context, equipsData []gamedata.EquipData) error {
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
