package app

import (
	"allen/trading-pov/models"
	"allen/trading-pov/parser"
	"fmt"
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

func NewEngine(exchange *Exchange, algo Algorithm) *Engine {
	return &Engine{
		exchange:           exchange,
		algo:               algo,
		pendingOrderSlices: make(map[*models.OrderSlice]int),
		pendingOrderPQView: make(map[int]int),
	}
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
	fmt.Printf("Engine: Received client order: %v\n", e.order)
	e.algo.Process(e)
}

func (e *Engine) ReceiveEvent(event []string) {
	evt := parser.ParseEvent(event)
	e.updateStateOnEvent(evt, event[0])
	e.algo.Process(e)
}

func (e *Engine) updateStateOnEvent(evt interface{}, eventType string) {
	if eventType == "T" {
		e.currentTime = evt.(*models.Trade).Time
		e.volume += evt.(*models.Trade).PQ.Quantity
	} else {
		e.currentQuote = evt.(*models.Quote)
		e.currentTime = e.currentQuote.Time
	}
}

// cancel
func (e *Engine) cancelOrderSlice(slice *models.OrderSlice) {
	fmt.Printf("Engine: Cancelled slice: %v\n", slice)
	// send cancel request to exchange
	e.exchange.CancelOrderSlice(slice)
	// clearup pending orders (an order must be pending if can be cancelled)
	e.RemovePendingOrderSlice(slice)
}

// helper to cancel all slices with given price
func (e *Engine) cancelAllSlicesWithPrice(price int) {
	for slice := range e.pendingOrderSlices {
		if slice.Price != price {
			continue
		}
		e.cancelOrderSlice(slice)
	}
}

// new
func (e *Engine) NewOrderSlice(slice *models.OrderSlice) OrderResponse {
	resp := e.exchange.NewOrderSlice(slice)
	fmt.Printf("Engine: New slice: %v, response: %s\n", slice, resp)
	switch resp {
	case ResponseFilled: // Filled immediately
		e.OrderSliceFilled(slice, false)
	case ResponseQueued:
		e.AddPendingOrderSlice(slice)
	}
	return resp
}

// consider reading from a channel, so that the exhcange don't have to include an *Engine field.
// Q: Do our fills also increase the volume traded? Should be yes.
func (e *Engine) OrderSliceFilled(slice *models.OrderSlice, pending bool) {
	e.volume += slice.Quantity
	e.order.QuantityFilled += slice.Quantity
	if pending {
		e.RemovePendingOrderSlice(slice)
	}
}

func (e *Engine) AddPendingOrderSlice(slice *models.OrderSlice) {
	e.pendingOrderSlices[slice] = 1
	e.pendingOrderPQView[slice.Price] += slice.Quantity
}
func (e *Engine) RemovePendingOrderSlice(slice *models.OrderSlice) {
	delete(e.pendingOrderSlices, slice)
	e.pendingOrderPQView[slice.Price] -= slice.Quantity
}
