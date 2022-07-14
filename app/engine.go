package app

import (
	"allen/trading-pov/models"
	"allen/trading-pov/parser"
	"allen/trading-pov/util"
	"fmt"
)

type Engine struct {
	algo               Algorithm
	order              *models.Order
	pendingOrderSlices map[*models.OrderSlice]int // A set of pending order slices.
	pendingOrderPQView map[float64]int            // Records the total quantity pending at every price level. Key: price, Value: quantity
	volume             int
	currentTime        int
	currentQuote       *models.Quote
}

func NewEngine(algo Algorithm) *Engine {
	return &Engine{
		algo:               algo,
		pendingOrderSlices: make(map[*models.OrderSlice]int),
		pendingOrderPQView: make(map[float64]int),
		volume:             0,
	}
}

func (e *Engine) Order(FIXMsg string) {
	fixOrder, err := parser.ParseFIX(FIXMsg)
	if err != nil {
		fmt.Println(err)
		return
	}
	e.order = &models.Order{
		QuantityTotal:  fixOrder.Quantity,
		QuantityFilled: 0,
		TargetRate:     fixOrder.POVTargetProp,
		MinRate:        util.RoundFloat(fixOrder.POVTargetProp * 8 / 10),
		MaxRate:        util.RoundFloat(fixOrder.POVTargetProp * 12 / 10),
	}
	fmt.Printf("Engine: Received client order: %v\n\n", util.OrderToString(e.order))
}

func (e *Engine) ReceiveEvent(event []string) *models.Execution {
	if e.order != nil && e.order.QuantityFilled == e.order.QuantityTotal {
		fmt.Printf("Engine: Yay, completely filled client order!\n")
		return nil
	}
	evt, err := parser.ParseEvent(event)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Printf("Engine: Received event: %v\n", util.EventToString(event))

	e.updateStateOnEvent(evt, event[0])
	if e.order == nil {
		return nil
	}
	res := e.algo.Process(e)
	if res == nil || (res.SlicesToCancel == nil && res.SlicesToOrder == nil) {
		fmt.Printf("Engine: Quantity to fill after this execution: %v\n", e.order.QuantityTotal-e.order.QuantityFilled)
		fmt.Printf("Engine: Pending order slices after this execution: %v\n\n", util.MapToString(e.pendingOrderPQView))
	}
	return res
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

func (e *Engine) ReceiveReport(report *models.ExecutionReport) {
	if report == nil || e.order == nil {
		return
	}
	if report.SlicesCancelled == nil && report.SlicesFilled == nil && report.SlicesQueued == nil {
		return
	}
	fmt.Printf("#### Got Response From the Exchange ####\n")

	for _, slice := range report.SlicesFilled {
		e.orderSliceFilled(slice)
	}
	for _, slice := range report.SlicesCancelled {
		e.orderSliceCancelled(slice)
	}
	for _, slice := range report.SlicesQueued {
		e.orderSliceQueued(slice)
	}

	fmt.Printf("Engine: Quantity to fill after this execution: %v\n", e.order.QuantityTotal-e.order.QuantityFilled)
	fmt.Printf("Engine: Pending order slices after this execution: %v\n\n", util.MapToString(e.pendingOrderPQView))
}

func (e *Engine) orderSliceCancelled(slice *models.OrderSlice) {
	fmt.Printf("Engine: Cancelled slice: %v@%v, ordered at timestamp %v\n", slice.Quantity, slice.Price, slice.TimeStamp)
	e.removePendingOrderSlice(slice)
}
func (e *Engine) orderSliceQueued(slice *models.OrderSlice) {
	fmt.Printf("Engine: Queued slice: %v@%v, ordered at timestamp %v\n", slice.Quantity, slice.Price, slice.TimeStamp)
	e.addPendingOrderSlice(slice)
}

// cancel
func (e *Engine) cancelOrderSlice(slice *models.OrderSlice, execution *models.Execution) {
	fmt.Printf("Engine: Cancelling slice: %v@%v, ordered at timestamp %v\n", slice.Quantity, slice.Price, slice.TimeStamp)
	execution.SlicesToCancel = append(execution.SlicesToCancel, slice)
}

// helper to cancel all slices with given price
func (e *Engine) cancelAllSlicesWithPrice(price float64, execution *models.Execution) {
	for slice := range e.pendingOrderSlices {
		if slice.Price != price {
			continue
		}
		e.cancelOrderSlice(slice, execution)
	}
}

// helper to cancel all slices with price no-more-existent
func (e *Engine) cancelNoMoreExistentPriceSlices(execution *models.Execution) {
	for slice := range e.pendingOrderSlices {
		ok := false
		for _, pq := range e.currentQuote.Bids {
			if pq.Price == slice.Price {
				ok = true
				break
			}
		}
		if !ok {
			e.cancelOrderSlice(slice, execution)
		}
	}
}

// new
func (e *Engine) newOrderSlice(slice *models.OrderSlice, execution *models.Execution) {
	fmt.Printf("Engine: New slice: %v@%v\n", slice.Quantity, slice.Price)
	execution.SlicesToOrder = append(execution.SlicesToOrder, slice)
}

// filled
func (e *Engine) orderSliceFilled(slice *models.OrderSlice) {
	e.order.QuantityFilled += slice.Quantity
	fmt.Printf("Engine: Slice filled!\n")
	fmt.Printf("Filled: %v@%v, Cumulative Quantity: %v\n", slice.Quantity, slice.Price, e.order.QuantityFilled)
	if _, ok := e.pendingOrderSlices[slice]; ok { // the slice was pending
		e.removePendingOrderSlice(slice)
	}
}

// helpers to manipulate pending order slice states
func (e *Engine) addPendingOrderSlice(slice *models.OrderSlice) {
	e.pendingOrderSlices[slice] = 1
	e.pendingOrderPQView[slice.Price] += slice.Quantity
}
func (e *Engine) removePendingOrderSlice(slice *models.OrderSlice) {
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
