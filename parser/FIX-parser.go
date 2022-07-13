package parser

import (
	"allen/trading-pov/models"
	"errors"
	"strconv"
	"strings"
)

func ParseFIX(FIXMsg string) (*models.FIXOrder, error) {
	errMsg := "bad FIX msg: " + FIXMsg
	var err error
	params := strings.Split(FIXMsg, ";")
	res := &models.FIXOrder{}

	ok := 0
	for _, param := range params {
		param = strings.TrimSpace(param)
		kv := strings.Split(param, "=")

		if len(kv) != 2 {
			return nil, errors.New(errMsg)
		}

		switch kv[0] {
		case "54":
			if kv[1] == "1" {
				res.Buy = true
			} else {
				res.Buy = false
			}
			ok++
		case "38":
			res.Quantity, err = strconv.Atoi(kv[1])
			if err != nil {
				return nil, errors.New(errMsg)
			}
			ok++
		case "6404":
			percentage, err := strconv.Atoi(kv[1])
			if err != nil {
				return nil, errors.New(errMsg)
			}
			res.POVTargetProp = float64(percentage) / 100
			ok++
		}
	}
	if ok != 3 {
		return nil, errors.New(errMsg)
	}
	if !res.Buy || res.POVTargetProp <= 0 || res.POVTargetProp*1.2 > 100 || res.Quantity <= 0 {
		return nil, errors.New(errMsg)
	}
	return res, nil
}
