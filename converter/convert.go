package converter

import (
	"strconv"
	"strings"
)

func StringToIntSlice(s string) []int64 {
	var list []int64
	for _, v := range strings.Split(strings.TrimSpace(s), ",") {
		i, _ := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		if i != 0 {
			list = append(list, i)
		}
	}
	return list
}

func StringToStringSlice(s string) []string {
	var list []string
	for _, v := range strings.Split(strings.TrimSpace(s), ",") {
		if vv := strings.TrimSpace(v); len(vv) != 0 {
			list = append(list, vv)
		}
	}
	return list
}

func StringSliceToInt64Slice(s []string) []int64 {
	var list []int64
	for _, v := range s {
		i, _ := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		list = append(list, i)
	}
	return list
}

func Int64SliceToStringSlice(s []int64) []string {
	var list []string
	for _, v := range s {
		list = append(list, strconv.FormatInt(v, 10))
	}
	return list
}

func Int64SliceToString(s []int64) string {
	return strings.Join(Int64SliceToStringSlice(s), ",")
}

func ContainsInt64(s []int64, t int64) bool {
	for _, v := range s {
		if v == t {
			return true
		}
	}
	return false
}

func UniqueInt64(s []int64, positive bool) []int64 {
	m := make(map[int64]struct{}, len(s))
	i := 0
	for _, t := range s {
		if _, ok := m[t]; ok {
			continue
		}
		if positive && t <= 0 {
			continue
		}
		m[t] = struct{}{}
		s[i] = t
		i++
	}
	return s[:i]
}

func UniqueString(s []string) []string {
	m := make(map[string]struct{}, len(s))
	i := 0
	for _, t := range s {
		if _, ok := m[t]; ok {
			continue
		}
		m[t] = struct{}{}
		s[i] = t
		i++
	}
	return s[:i]
}
