package app

import (
	"allen/trading-pov/models"
	"allen/trading-pov/parser"
)

type Engine struct {
	exchange           *Exchange
	algo               Algorithm
	order              *models.Order
	pendingOrderSlices map[*models.OrderSlice]int // TODO: Add pending slices when creating orders
	pendingOrderPQView map[int]int
	volume             int
	currentTime        int
	currentQuote       *models.Quote
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Order(FIXMsg string) {
	fixOrder := parser.ParseFIX(FIXMsg)
	e.order = &models.Order{
		QuantityTotal:  fixOrder.Quantity,
		QuantityFilled: 0,
		TargetRate:     fixOrder.POVTargetProp,
		MinRate:        e.order.TargetRate * 0.8,
		MaxRate:        e.order.TargetRate * 1.2,
	}
	e.algo.Process(e)
}

func (e *Engine) ReceiveEvent(event []string) {
	evt := parser.ParseEvent(event)
	e.updateState(evt, event[0])
	e.algo.Process(e)
}

func (e *Engine) updateState(evt interface{}, eventType string) {
	if eventType == "T" {
		e.currentTime = evt.(*models.Trade).Time
		e.volume += evt.(*models.Trade).PQ.Quantity
	} else {
		e.currentQuote = evt.(*models.Quote)
		e.currentTime = e.currentQuote.Time
	}
}

// consider reading from a channel, so that the exhcange don't have to include an *Engine field.
// Q: Do our fills also increase the volume traded? Should be yes.
func (e *Engine) OrderSliceFilled(slice *models.OrderSlice, pending bool) {
	e.volume += slice.Quantity
	e.order.QuantityFilled += slice.Quantity
	if pending {
		delete(e.pendingOrderSlices, slice)
	}
}

func (e *Engine) cancelAllSlicesWithPrice(price int) {

}
