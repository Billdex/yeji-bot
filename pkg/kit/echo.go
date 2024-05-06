package kit

import (
	"fmt"
	"strconv"
	"strings"
)

func ParsePage(str string, defaultPage int) int {
	str = strings.ReplaceAll(str, "p", "")
	str = strings.ReplaceAll(str, "P", "")
	str = strings.ReplaceAll(str, "-", "")
	str = strings.TrimSpace(str)
	page, err := strconv.Atoi(str)
	if err != nil {
		return defaultPage
	}
	return page
}

// PaginationOutput 对数据结果进行分页输出, 超出总页数的按照最后一页输出，页数小于 1 的按照第一页输出
func PaginationOutput[T any](list []T, page int, pageSize int, title string, output func(item T) string) string {
	var msg string = title
	maxPage := (len(list)-1)/pageSize + 1
	if len(list) > pageSize {
		if page > maxPage {
			page = maxPage
		}
		msg = fmt.Sprintf("%s (%d/%d)", title, page, maxPage)
	}
	if page <= 0 {
		page = 1
	}
	for i := (page - 1) * pageSize; i < page*pageSize && i < len(list); i++ {
		msg += fmt.Sprintf("\n%s", output(list[i]))
	}
	if page < maxPage {
		msg += "\n......"
	}
	return msg
}

// FormatRecipeTime 格式化时间
func FormatRecipeTime(t int) string {
	if t < 0 {
		return ""
	} else if t == 0 {
		return "0秒"
	} else {
		var time string
		hour := t / 3600
		minute := t % 3600 / 60
		second := t % 3600 % 60
		if hour > 0 {
			time += fmt.Sprintf("%d小时", hour)
		}
		if minute > 0 {
			time += fmt.Sprintf("%d分", minute)
		}
		if second > 0 {
			time += fmt.Sprintf("%d秒", second)
		}
		return time
	}
}
