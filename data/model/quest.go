package model

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strings"
)

// Quest 任务数据
type Quest struct {
	QuestId     int              `gorm:"column:quest_id"`                   // 任务 ID
	QuestIdDisp decimal.Decimal  `gorm:"column:quest_id_disp"`              // 系列任务编号
	Type        string           `gorm:"column:type"`                       // 系列类型
	Goal        string           `gorm:"column:goal"`                       // 任务目标
	Rewards     []QuestReward    `gorm:"column:rewards;serializer:json"`    // 任务奖励
	Conditions  []QuestCondition `gorm:"column:conditions;serializer:json"` // 任务条件
}

func (Quest) TableName() string {
	return "quest"
}

func (q Quest) RewardsStr() string {
	if len(q.Rewards) > 0 {
		rewardStrList := make([]string, 0, len(q.Rewards))
		for _, reward := range q.Rewards {
			if reward.Quantity == "" {
				rewardStrList = append(rewardStrList, reward.Name)
			} else {
				rewardStrList = append(rewardStrList, fmt.Sprintf("%s*%s", reward.Name, reward.Quantity))
			}
		}
		return strings.Join(rewardStrList, ", ")
	}
	return "无"
}

type QuestReward struct {
	Name     string `json:"name"`     // 奖励名称
	Quantity string `json:"quantity"` // 奖励数量
}

type QuestCondition struct {
	RecipeId     int    `json:"recipeId"`
	Rank         int    `json:"rank"`
	Num          int    `json:"num"`
	GoldEff      bool   `json:"goldEff"`
	MaterialId   int    `json:"materialId"`
	Guest        string `json:"guest"`
	AnyGuest     bool   `json:"anyGuest"`
	Skill        string `json:"skill"`
	MaterialEff  bool   `json:"materialEff"`
	NewGuest     bool   `json:"newGuest"`
	Rarity       int    `json:"rarity"`
	Price        int    `json:"price"`
	Category     string `json:"category"`
	Condiment    string `json:"condiment"`
	CondimentEff bool   `json:"condimentEff"`
}
