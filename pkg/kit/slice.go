package kit

import "strings"

func SliceMap[T any, R any](arr []T, f func(T) R) []R {
	res := make([]R, len(arr))
	for i := range arr {
		res[i] = f(arr[i])
	}
	return res
}

func SliceFilter[T any](arr []T, f func(T) bool) []T {
	res := make([]T, 0, len(arr))
	for i := range arr {
		if f(arr[i]) {
			res = append(res, arr[i])
		}
	}
	return res
}

func SliceContains[T comparable](arr []T, item T) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

func HasSuffixIn(str string, suffixes []string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(str, suffix) {
			return true
		}
	}
	return false
}
