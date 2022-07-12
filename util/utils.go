package util

import (
	"allen/trading-pov/models"
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
	if event == nil || len(event) < 4 {
		return ""
	}
	res := ""
	if event[0] == "Q" {
		res += "Quote@"
		res += event[1]
		res += ", Bids: "
		res += event[2]
		res += ", Asks: "
		res += event[3]
	} else {
		res += "Trade@"
		res += event[1]
		res += ", Price: "
		res += event[2]
		res += ", Quantity: "
		res += event[3]
	}

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
	if len(res) > 2 {
		res = res[:len(res)-2]
	}
	return res
}

func OrderToString(order *models.Order) string {
	res := ""
	res += "Total Quantity: "
	res += strconv.Itoa(order.QuantityTotal)
	res += ", Target Rate: "
	res += fmt.Sprintf("%f", order.TargetRate)
	res += ", Min Rate: "
	res += fmt.Sprintf("%f", order.MinRate)
	res += ", Max Rate: "
	res += fmt.Sprintf("%f", order.MaxRate)
	return res
}
