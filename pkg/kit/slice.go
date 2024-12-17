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

func HasPrefixIn(str string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(str, prefix) {
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

func WhichPrefixIn(str string, prefixes []string) string {
	for _, prefix := range prefixes {
		if strings.HasPrefix(str, prefix) {
			return prefix
		}
	}
	return ""
}

func TrimPrefixIn(str string, prefixes []string) string {
	for _, prefix := range prefixes {
		if strings.HasPrefix(str, prefix) {
			return strings.TrimPrefix(str, prefix)
		}
	}
	return str
}

func SliceToMap[T any, key comparable](arr []T, keyFunc func(T) key) map[key]T {
	m := make(map[key]T, len(arr))
	for i := range arr {
		m[keyFunc(arr[i])] = arr[i]
	}
	return m
}
