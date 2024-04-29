package model

type Skill struct {
	SkillId     int           `gorm:"column:skill_id"`                // 技能ID
	Description string        `gorm:"column:description"`             // 技能描述
	Effects     []SkillEffect `gorm:"column:effects;serializer:json"` // 技能效果详情
}

func (Skill) TableName() string {
	return "skill"
}

// SkillEffect 技能效果详情
type SkillEffect struct {
	Calculation string  `json:"calculation"` // 计算方式
	Type        string  `json:"type"`        // 效果类型
	Condition   string  `json:"condition"`   // 满足条件
	Tag         int     `json:"tag"`         // 对厨师生效的性别 1:男 2:女
	Value       float64 `json:"value"`       // 效果数值
}
