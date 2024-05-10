package model

import (
	"strings"
)

type Recipe struct {
	RecipeId           int               `gorm:"column:recipe_id"`                      // 菜谱ID
	Name               string            `gorm:"column:name"`                           // 菜名
	GalleryId          string            `gorm:"column:gallery_id"`                     // 图鉴ID
	Rarity             int               `gorm:"column:rarity"`                         // 稀有度
	Origins            []string          `gorm:"column:origins;serializer:json"`        // 来源列表
	Stirfry            int               `gorm:"column:stirfry"`                        // 炒技法
	Bake               int               `gorm:"column:bake"`                           // 烤技法
	Boil               int               `gorm:"column:boil"`                           // 煮技法
	Steam              int               `gorm:"column:steam"`                          // 蒸技法
	Fry                int               `gorm:"column:fry"`                            // 炸技法
	Cut                int               `gorm:"column:cut"`                            // 切技法knife
	Condiment          string            `gorm:"column:condiment"`                      // 调料
	Price              int               `gorm:"column:price"`                          // 价格
	ExPrice            int               `gorm:"column:ex_price"`                       // 熟练加价
	GoldEfficiency     int               `gorm:"column:gold_efficiency"`                // 金币效率
	MaterialEfficiency int               `gorm:"column:material_efficiency"`            // 耗材效率
	Gift               string            `gorm:"column:gift"`                           // 第一次做到神级送的礼物
	GuestGifts         []RecipeGuestGift `gorm:"column:guest_gifts;serializer:json"`    // 贵客礼物列表
	UpgradeGuests      []string          `gorm:"column:upgrade_guests;serializer:json"` // 升阶贵客列表
	Time               int               `gorm:"column:time"`                           // 每份时间(秒)
	Limit              int               `gorm:"column:limit"`                          // 每组数量
	TotalTime          int               `gorm:"column:total_time"`                     // 每组时间(秒)
	Unlock             string            `gorm:"column:unlock"`                         // 可解锁
	Combos             []string          `gorm:"column:combos;serializer:json"`         // 可合成菜谱
	Img                string            `gorm:"column:img"`                            // 图片
	Materials          []RecipeMaterial  `gorm:"column:materials;serializer:json"`      // 所需食材数据
}

func (Recipe) TableName() string {
	return "recipe"
}

// NeedSkill 判断菜谱是否需要某个技法
func (recipe Recipe) NeedSkill(skill string) bool {
	switch strings.ToLower(skill) {
	case "炒", "stirfry":
		return recipe.Stirfry > 0
	case "烤", "bake":
		return recipe.Bake > 0
	case "煮", "boil":
		return recipe.Boil > 0
	case "蒸", "steam":
		return recipe.Steam > 0
	case "炸", "fry":
		return recipe.Fry > 0
	case "切", "cut", "knife":
		return recipe.Cut > 0
	}
	return false
}

// UsedMaterial 判断菜谱是否包含某个食材
func (recipe Recipe) UsedMaterial(material string) bool {
	for i := range recipe.Materials {
		if recipe.Materials[i].MaterialName == material {
			return true
		}
	}
	return false
}

// UsedMaterials 判断菜谱是否包含某些食材
func (recipe Recipe) UsedMaterials(materials []string) bool {
	materialMap := make(map[string]struct{}, len(materials))
	for i := range materials {
		materialMap[materials[i]] = struct{}{}
	}
	for i := range recipe.Materials {
		if _, ok := materialMap[recipe.Materials[i].MaterialName]; ok {
			return true
		}
	}
	return false
}

// HasMaterialFrom 判断菜谱是否包含某个来源的食材
func (recipe Recipe) HasMaterialFrom(origin string) bool {
	for i := range recipe.Materials {
		if recipe.Materials[i].Origin == origin {
			return true
		}
	}
	return false
}

// HasMaterialFromIn 判断菜谱是否包含某些来源的食材
func (recipe Recipe) HasMaterialFromIn(origins []string) bool {
	originMap := make(map[string]struct{}, len(origins))
	for i := range origins {
		originMap[origins[i]] = struct{}{}
	}
	for i := range recipe.Materials {
		if _, ok := originMap[recipe.Materials[i].Origin]; ok {
			return true
		}
	}
	return false
}

// RecipeGuestGift 菜谱贵客礼物数据
type RecipeGuestGift struct {
	RecipeId   int    `json:"-"`          // 菜谱 id
	RecipeName string `json:"-"`          // 菜谱名称
	GuestId    int    `json:"-"`          // 贵客 id
	GuestName  string `json:"guest_name"` // 贵客名称
	Antique    string `json:"antique"`    // 符文名称
}

// RecipeMaterial 菜谱食材数据
type RecipeMaterial struct {
	RecipeId     int    `json:"-"`      // 菜谱 id
	RecipeName   string `json:"-"`      // 菜谱名称
	MaterialId   int    `json:"m_id"`   // 食材 id
	MaterialName string `json:"m_name"` // 食材名称
	Quantity     int    `json:"qty"`    // 消耗数量
	Origin       string `json:"origin"` // 食材来源
	Efficiency   int    `json:"-"`      // 食材消耗效率
}
