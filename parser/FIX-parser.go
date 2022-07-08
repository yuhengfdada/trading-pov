package parser

import (
	"allen/trading-pov/models"
	"strconv"
	"strings"
)

func ParseFIX(FIXMsg string) models.FIXOrder {
	params := strings.Split(FIXMsg, ";")
	res := models.FIXOrder{}
	for _, param := range params {
		param = strings.TrimSpace(param)
		kv := strings.Split(param, "=")
		switch kv[0] {
		case "54":
			if kv[1] == "1" {
				res.Buy = true
			} else {
				res.Buy = false
			}
		case "38":
			res.Quantity, _ = strconv.Atoi(kv[1])
		case "6404":
			percentage, _ := strconv.Atoi(kv[1])
			res.POVTargetProp = float64(percentage) / 100
		}
	}
	return res
}
