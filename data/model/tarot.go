package model

// Tarot 抽签分数
type Tarot struct {
	Id          int64  `gorm:"column:id"`
	Score       int64  `gorm:"column:score"`       // 抽签分数
	Description string `gorm:"column:description"` // 签文内容
}

func (Tarot) TableName() string {
	return "tarot"
}

func (t Tarot) LevelStr() string {
	switch {
	case t.Score == 0:
		return "不知道吉不吉"
	case 0 < t.Score && t.Score < 15:
		return "小小吉"
	case 15 <= t.Score && t.Score < 40:
		return "小吉"
	case 40 <= t.Score && t.Score < 60:
		return "中吉"
	case 60 <= t.Score && t.Score < 85:
		return "大吉"
	case 85 <= t.Score && t.Score < 100:
		return "大大吉"
	case t.Score == 100:
		return "超吉"
	default:
		return "?"
	}
}
