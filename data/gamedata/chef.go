package gamedata

// ChefData 厨师数据
type ChefData struct {
	ChefId        int    `json:"chefId"`
	Name          string `json:"name"`          // 厨师名称
	Tags          []int  `json:"tags"`          // 标签列表，第0位是性别
	Rarity        int    `json:"rarity"`        // 稀有度
	Origin        string `json:"origin"`        // 来源
	GalleryId     string `json:"galleryId"`     // 图鉴id
	Stirfry       int    `json:"stirfry"`       // 炒技法
	Bake          int    `json:"bake"`          // 烤技法
	Boil          int    `json:"boil"`          // 煮技法
	Steam         int    `json:"steam"`         // 蒸技法
	Fry           int    `json:"fry"`           // 炸技法
	Cut           int    `json:"knife"`         // 切技法
	Meat          int    `json:"meat"`          // 肉类采集
	Flour         int    `json:"creation"`      // 面类/作坊采集
	Fish          int    `json:"fish"`          // 鱼类采集
	Vegetable     int    `json:"veg"`           // 蔬菜采集
	Sweet         int    `json:"sweet"`         // 甜
	Sour          int    `json:"sour"`          // 咸
	Spicy         int    `json:"spicy"`         // 辣
	Salty         int    `json:"salty"`         // 咸
	Bitter        int    `json:"bitter"`        // 苦
	Tasty         int    `json:"tasty"`         // 鲜
	SkillId       int    `json:"skill"`         // 技能id
	UltimateGoals []int  `json:"ultimateGoal"`  // 修炼任务列表
	UltimateSkill int    `json:"ultimateSkill"` // 修炼技能id
	Disk          int    `json:"disk"`          // 心法盘id
}

// DiskData 心法盘数据
type DiskData struct {
	DiskId   int   `json:"diskId"`
	Info     []int `json:"info"`     // 孔位数据
	MaxLevel int   `json:"maxLevel"` // 最高等级
}
