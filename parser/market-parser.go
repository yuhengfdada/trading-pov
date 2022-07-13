package parser

import (
	"allen/trading-pov/models"
	"errors"
	"strings"
)

func ParseEvent(event []string) (interface{}, error) {
	if event == nil {
		return nil, errors.New("invalid event")
	}
	errMsg := "invalid event: " + strings.Join(event, ",")
	if len(event) != 4 {
		return nil, errors.New(errMsg)
	}
	if event[0] == "T" {
		t, err := models.NewTrade(event[1], event[2], event[3])
		if err != nil {
			return nil, errors.New(errMsg)
		}
		return t, nil
	} else {
		q, err := models.NewQuote(event[1], event[2], event[3])
		if err != nil {
			return nil, errors.New(errMsg)
		}
		return q, nil
	}
}
