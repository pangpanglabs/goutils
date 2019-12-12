package converter

import (
	"strconv"
	"strings"
)

func StringToIntSlice(s string) []int64 {
	var list []int64
	for _, v := range strings.Split(strings.TrimSpace(s), ",") {
		i, _ := strconv.ParseInt(v, 10, 64)
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
