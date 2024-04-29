package model

type Material struct {
	MaterialId int    `gorm:"column:material_id"` // 食材ID
	Name       string `gorm:"column:name"`        // 食材名
	Origin     string `gorm:"column:origin"`      // 食材来源
}

func (Material) TableName() string {
	return "material"
}
