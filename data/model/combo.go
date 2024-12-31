package model

// ComboRecipe 后厨合成菜谱数据
type ComboRecipe struct {
	RecipeId      int      `gorm:"column:recipe_id"`                       // 合成菜谱 id
	RecipeName    string   `gorm:"column:recipe_name"`                     // 菜谱名称
	NeedRecipeIds []int    `gorm:"column:need_recipe_ids;serializer:json"` // 所需菜谱 id 列表
	NeedRecipes   []Recipe `gorm:"-"`                                      // 所需菜谱列表
}

func (ComboRecipe) TableName() string {
	return "combo_recipe"
}
