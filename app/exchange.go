package app

import (
	"allen/trading-pov/models"
	"allen/trading-pov/parser"
	"time"
)

type OrderResponse int

const (
	ResponseFilled = 0
	ResponseQueued = 1
)

type Exchange struct {
	engine             *Engine
	pendingOrderSlices map[*models.OrderSlice]int
	currentTime        int
	currentQuote       *models.Quote
}

func NewExchange() *Exchange {
	return &Exchange{}
}

func (exch *Exchange) SetEngine(e *Engine) {
	exch.engine = e
}

func (exch *Exchange) ReceiveEvent(event []string) {
	evt := parser.ParseEvent(event)
	exch.updateState(evt, event[0])
	for slice := range exch.pendingOrderSlices {
		if exch.meetFillCriteria(slice) {
			exch.engine.OrderSliceFilled(slice, true)
			delete(exch.pendingOrderSlices, slice) // deletion won't affect the iteration
		}
	}
}

func (exch *Exchange) updateState(evt interface{}, eventType string) {
	if eventType == "T" {
		exch.currentTime = evt.(*models.Trade).Time
	} else {
		exch.currentQuote = evt.(*models.Quote)
		exch.currentTime = exch.currentQuote.Time
	}
}

func (exch *Exchange) meetFillCriteria(slice *models.OrderSlice) bool {
	timesup := time.Duration(exch.currentTime-slice.TimeStamp)*time.Millisecond >= time.Duration(3)*time.Minute
	// marketable
	if slice.Price >= exch.currentQuote.Asks[0].Price {
		return true
	}
	// after 3 mins, still the best bid
	if timesup && slice.Price == exch.currentQuote.Bids[0].Price {
		return true
	}
	return false
}

func (exch *Exchange) NewOrderSlice(slice *models.OrderSlice) OrderResponse {
	if exch.meetFillCriteria(slice) {
		return ResponseFilled
	}
	exch.pendingOrderSlices[slice] = 1
	return ResponseQueued
}

func (exch *Exchange) CancelOrderSlice(slice *models.OrderSlice) {
	delete(exch.pendingOrderSlices, slice)
}
