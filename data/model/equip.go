package model

type Equip struct {
	EquipId   int      `gorm:"column:equip_id"`                // 厨具ID
	Name      string   `gorm:"column:name"`                    // 厨具名称
	GalleryId string   `gorm:"column:gallery_id"`              // 图鉴ID
	Origins   []string `gorm:"column:origins;serializer:json"` // 列表来源
	Rarity    int      `gorm:"column:rarity"`                  // 稀有度
	Skills    []int    `gorm:"column:skills;serializer:json"`  // 技能效果组
	Img       string   `gorm:"column:img"`                     // 图片地址

	SkillDescs []string `gorm:"-"` // 技能效果描述组
}

func (Equip) TableName() string {
	return "equip"
}
