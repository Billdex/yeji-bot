package model

var rarityStrToIntMap = map[string]int{
	"1火": 1, "1星": 1, "一火": 1, "一星": 1,
	"2火": 2, "2星": 2, "二火": 2, "二星": 2, "两火": 2, "两星": 2,
	"3火": 3, "3星": 3, "三火": 3, "三星": 3,
	"4火": 4, "4星": 4, "四火": 4, "四星": 4,
	"5火": 5, "5星": 5, "五火": 5, "五星": 5,
}

func IsRarityStr(rarity string) bool {
	_, ok := rarityStrToIntMap[rarity]
	return ok
}

func RarityToInt(rarity string) int {
	return rarityStrToIntMap[rarity]
}
