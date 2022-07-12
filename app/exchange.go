package app

import (
	"allen/trading-pov/models"
	"allen/trading-pov/parser"
	"allen/trading-pov/util"
	"fmt"
	"time"
)

type OrderResponse string

const (
	ResponseFilled = "filled"
	ResponseQueued = "queued"
)

type Exchange struct {
	engine             *Engine
	pendingOrderSlices map[*models.OrderSlice]int
	currentTime        int
	currentQuote       *models.Quote
}

func NewExchange() *Exchange {
	return &Exchange{
		engine:             nil,
		pendingOrderSlices: make(map[*models.OrderSlice]int),
		currentTime:        0,
	}
}

// the exchange needs to know the engine to call it back
func (exch *Exchange) SetEngine(e *Engine) {
	exch.engine = e
}

func (exch *Exchange) ReceiveEvent(event []string) {
	fmt.Printf("Exchange: Current Event: %v\n", util.EventToString(event))
	evt := parser.ParseEvent(event)
	exch.updateStateOnEvent(evt, event[0])
	for slice := range exch.pendingOrderSlices {
		if exch.meetFillCriteria(slice) {
			exch.engine.OrderSliceFilled(slice, true)
			delete(exch.pendingOrderSlices, slice) // deletion won't affect the iteration
		}
	}
}

func (exch *Exchange) updateStateOnEvent(evt interface{}, eventType string) {
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

// for engine to call
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
