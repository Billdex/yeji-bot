package model

import (
	"strings"
	"yeji-bot/pkg/kit"
)

// Chef 厨师数据对应的数据库模型
type Chef struct {
	ChefId        int      `gorm:"column:chef_id"`                        // 厨师 id
	Name          string   `gorm:"column:name"`                           // 厨师名字
	Gender        int      `gorm:"column:gender"`                         // 性别
	Rarity        int      `gorm:"column:rarity"`                         // 稀有度
	Origins       []string `gorm:"column:origins;serializer:json"`        // 来源
	GalleryId     string   `gorm:"column:gallery_id"`                     // 图鉴id
	Stirfry       int      `gorm:"column:stirfry"`                        // 炒技法
	Bake          int      `gorm:"column:bake"`                           // 烤技法
	Boil          int      `gorm:"column:boil"`                           // 煮技法
	Steam         int      `gorm:"column:steam"`                          // 蒸技法
	Fry           int      `gorm:"column:fry"`                            // 炸技法
	Cut           int      `gorm:"column:cut"`                            // 切技法knife
	Meat          int      `gorm:"column:meat"`                           // 肉类采集
	Flour         int      `gorm:"column:flour"`                          // 面类采集
	Fish          int      `gorm:"column:fish"`                           // 水产采集
	Vegetable     int      `gorm:"column:vegetable"`                      // 蔬菜采集
	Sweet         int      `gorm:"column:sweet"`                          // 甜
	Sour          int      `gorm:"column:sour"`                           // 酸
	Spicy         int      `gorm:"column:spicy"`                          // 辣
	Salty         int      `gorm:"column:salty"`                          // 咸
	Bitter        int      `gorm:"column:bitter"`                         // 苦
	Tasty         int      `gorm:"column:tasty"`                          // 鲜
	SkillId       int      `gorm:"column:skill_id"`                       // 技能id
	UltimateGoals []int    `gorm:"column:ultimate_goals;serializer:json"` // 修炼任务id数组
	UltimateSkill int      `gorm:"column:ultimate_skill"`                 // 修炼效果id
	Img           string   `gorm:"column:img"`                            // 图片地址
	DiskInfo      []int    `gorm:"column:disk_info;serializer:json"`      // 心法盘孔位数据
	DiskLevel     int      `gorm:"column:disk_level"`                     // 心法盘最大等级

	SkillDesc         string `gorm:"-"` // 技能描述
	UltimateSkillDesc string `gorm:"-"` // 修炼技能描述
}

func (Chef) TableName() string {
	return "chef"
}

// DiskInfoFmt 心法盘数据
func (c Chef) DiskInfoFmt() string {
	return strings.Join(kit.SliceMap(c.DiskInfo, func(d int) string {
		switch d {
		case 1:
			return "红"
		case 2:
			return "绿"
		case 3:
			return "蓝"
		}
		return ""
	}), "|")
}
