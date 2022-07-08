package parser

import (
	"allen/trading-pov/models"
	"fmt"
)

func ParseEvent(event []string) interface{} {
	if event[0] == "T" {
		t := models.NewTrade(event[1], event[2], event[3])
		fmt.Println(t)
		return t
	} else {
		q := models.NewQuote(event[1], event[2], event[3])
		fmt.Println(q)
		return q
	}
}
