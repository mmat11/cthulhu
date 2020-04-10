package task

import "strconv"

func int64FromStr(s string) int64 {
	i, _ := strconv.Atoi(s)
	return int64(i)
}
