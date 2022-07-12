package util

import (
	"fmt"
	"sort"
	"strconv"
)

func RoundFloat(num float64) float64 {
	return float64(int((num+0.000001)*10000)) / 10000
}

func Min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func EventToString(event []string) string {
	res := ""
	if event[0] == "Q" {
		res += "Quote@"
	} else {
		res += "Trade@"
	}
	res += event[1]
	res += ", Bids: "
	res += event[2]
	res += ", Asks: "
	res += event[3]
	return res
}

func MapToString(mp map[float64]int) string {
	res := ""
	var keys []float64
	for k := range mp {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(keys)))
	for _, k := range keys {
		res += strconv.Itoa(mp[k])
		res += "@"
		res += fmt.Sprintf("%f", k)
		res += ", "
	}
	res = res[:len(res)-2]
	return res
}
