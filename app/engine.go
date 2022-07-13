package app

import (
	"allen/trading-pov/models"
	"allen/trading-pov/parser"
	"allen/trading-pov/util"
	"fmt"
)

type Engine struct {
	exchange           *Exchange
	algo               Algorithm
	order              *models.Order
	pendingOrderSlices map[*models.OrderSlice]int
	pendingOrderPQView map[float64]int
	volume             int
	currentTime        int
	currentQuote       *models.Quote
}

func NewEngine(exchange *Exchange, algo Algorithm) *Engine {
	return &Engine{
		exchange:           exchange,
		algo:               algo,
		pendingOrderSlices: make(map[*models.OrderSlice]int),
		pendingOrderPQView: make(map[float64]int),
		volume:             0,
	}
}

func (e *Engine) Order(FIXMsg string) {
	fixOrder := parser.ParseFIX(FIXMsg)
	e.order = &models.Order{
		QuantityTotal:  fixOrder.Quantity,
		QuantityFilled: 0,
		TargetRate:     fixOrder.POVTargetProp,
		MinRate:        util.RoundFloat(fixOrder.POVTargetProp * 8 / 10),
		MaxRate:        util.RoundFloat(fixOrder.POVTargetProp * 12 / 10),
	}
	fmt.Printf("Engine: Received client order: %v\n\n", util.OrderToString(e.order))
	if e.currentQuote == nil {
		return
	}
	e.algo.Process(e)
	fmt.Printf("Engine: Quantity to fill after this round: %v\n", e.order.QuantityTotal-e.order.QuantityFilled)
	fmt.Printf("Engine: Pending order slices after this round: %v\n\n", e.pendingOrderPQView)
}

func (e *Engine) ReceiveEvent(event []string) {
	if e.order != nil && e.order.QuantityFilled == e.order.QuantityTotal {
		fmt.Printf("Engine: Yay, completely filled client order!\n")
		return
	}
	evt, err := parser.ParseEvent(event)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Engine: Received event: %v\n", util.EventToString(event))

	e.updateStateOnEvent(evt, event[0])
	if e.order == nil {
		return
	}
	e.algo.Process(e)
	fmt.Printf("Engine: Quantity to fill after this round: %v\n", e.order.QuantityTotal-e.order.QuantityFilled)
	fmt.Printf("Engine: Pending order slices after this round: %v\n\n", util.MapToString(e.pendingOrderPQView))
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
	fmt.Printf("Engine: Cancelled slice: %v@%v, ordered at timestamp %v\n", slice.Quantity, slice.Price, slice.TimeStamp)
	// send cancel request to exchange
	e.exchange.CancelOrderSlice(slice)
	// clearup pending orders (an order must be pending if can be cancelled)
	e.RemovePendingOrderSlice(slice)
}

// helper to cancel all slices with given price
func (e *Engine) cancelAllSlicesWithPrice(price float64) {
	for slice := range e.pendingOrderSlices {
		if slice.Price != price {
			continue
		}
		e.cancelOrderSlice(slice)
	}
}

// helper to cancel all slices with price no-more-existent
func (e *Engine) cancelNoMoreExistentPriceSlices() {
	for slice := range e.pendingOrderSlices {
		ok := false
		for _, pq := range e.currentQuote.Bids {
			if pq.Price == slice.Price {
				ok = true
				break
			}
		}
		if !ok {
			e.cancelOrderSlice(slice)
		}
	}
}

// new
func (e *Engine) NewOrderSlice(slice *models.OrderSlice) OrderResponse {
	if slice.Price == 0 {
		return ""
	}
	resp := e.exchange.NewOrderSlice(slice)
	fmt.Printf("Engine: New slice: %v@%v, response: %s\n", slice.Quantity, slice.Price, resp)
	switch resp {
	case ResponseFilled: // Filled immediately
		e.OrderSliceFilled(slice, false)
	case ResponseQueued:
		e.AddPendingOrderSlice(slice)
	}
	return resp
}

// filled
func (e *Engine) OrderSliceFilled(slice *models.OrderSlice, pending bool) {
	e.order.QuantityFilled += slice.Quantity
	fmt.Printf("Engine: Slice filled!\n")
	fmt.Printf("Filled: %v@%v, Cumulative Quantity: %v\n", slice.Quantity, slice.Price, e.order.QuantityFilled)
	if pending {
		e.RemovePendingOrderSlice(slice)
	}
}

// helpers to manipulate pending order slice states
func (e *Engine) AddPendingOrderSlice(slice *models.OrderSlice) {
	e.pendingOrderSlices[slice] = 1
	e.pendingOrderPQView[slice.Price] += slice.Quantity
}
func (e *Engine) RemovePendingOrderSlice(slice *models.OrderSlice) {
	delete(e.pendingOrderSlices, slice)
	e.pendingOrderPQView[slice.Price] -= slice.Quantity
	if e.pendingOrderPQView[slice.Price] == 0 {
		delete(e.pendingOrderPQView, slice.Price)
	}
}

// setters for testing purposes
func (e *Engine) setVolume(volume int) {
	e.volume = volume
}
func (e *Engine) setOrderFilledQuantity(qty int) {
	e.order.QuantityFilled = qty
}
