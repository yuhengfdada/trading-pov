package app

import (
	"allen/trading-pov/models"
	"allen/trading-pov/parser"
	"allen/trading-pov/util"
	"fmt"
	"sort"
	"time"
)

type OrderResponse string

const (
	ResponseFilled = "filled"
	ResponseQueued = "queued"
)

type Exchange struct {
	pendingOrderSlices map[*models.OrderSlice]int
	currentTime        int
	currentQuote       *models.Quote
}

func NewExchange() *Exchange {
	return &Exchange{
		pendingOrderSlices: make(map[*models.OrderSlice]int),
		currentTime:        0,
	}
}

func (exch *Exchange) ReceiveEvent(event []string) *models.ExecutionReport {
	evt, err := parser.ParseEvent(event)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Printf("Exchange: Current Event: %v\n", util.EventToString(event))

	exch.updateStateOnEvent(evt, event[0])

	if exch.currentQuote == nil {
		return nil
	}
	return exch.process()
}

func (exch *Exchange) updateStateOnEvent(evt interface{}, eventType string) {
	if eventType == "T" {
		exch.currentTime = evt.(*models.Trade).Time
	} else {
		exch.currentQuote = evt.(*models.Quote)
		exch.currentTime = exch.currentQuote.Time
	}
}

func (exch *Exchange) process() *models.ExecutionReport {
	report := &models.ExecutionReport{}
	if exch.pendingOrderSlices == nil {
		return report
	}

	// To calculate fills correctly even if asks are not enough in quantity
	var copyAsks []models.PriceQuantity
	copyAsks = append(copyAsks, exchange.currentQuote.Asks...)
	// To simulate real life, slices are filled in FIFO order.
	var sortedSlices []*models.OrderSlice
	for slice := range exch.pendingOrderSlices {
		sortedSlices = append(sortedSlices, slice)
	}
	sort.Slice(sortedSlices, func(i, j int) bool {
		return sortedSlices[i].TimeStamp < sortedSlices[j].TimeStamp
	})

	for _, slice := range sortedSlices {
		if exch.meetFillCriteria(slice, copyAsks) {
			report.SlicesFilled = append(report.SlicesFilled, slice)
			delete(exch.pendingOrderSlices, slice) // deletion won't affect the iteration
		}
	}
	return report
}

func (exch *Exchange) meetFillCriteria(slice *models.OrderSlice, asks []models.PriceQuantity) bool {
	timesup := time.Duration(exch.currentTime-slice.TimeStamp)*time.Millisecond >= time.Duration(3)*time.Minute
	// marketable.
	for idx := 0; idx < len(asks); idx++ {
		if slice.Price < asks[idx].Price {
			break
		}
		if slice.Quantity > asks[idx].Quantity {
			continue
		}
		asks[idx].Quantity -= slice.Quantity
		return true
	}
	// after 3 mins, still the best bid
	if timesup && slice.Price == exch.currentQuote.Bids[0].Price {
		return true
	}
	return false
}

func (exch *Exchange) ReceiveExecutions(execution *models.Execution) *models.ExecutionReport {
	if execution == nil || exch.currentQuote == nil {
		return nil
	}
	report := &models.ExecutionReport{}

	var copyAsks []models.PriceQuantity
	copyAsks = append(copyAsks, exchange.currentQuote.Asks...)

	for _, slice := range execution.SlicesToOrder {
		resp := exch.newOrderSlice(slice, copyAsks)
		switch resp {
		case ResponseFilled:
			report.SlicesFilled = append(report.SlicesFilled, slice)
		case ResponseQueued:
			report.SlicesQueued = append(report.SlicesQueued, slice)
		}
	}
	for _, slice := range execution.SlicesToCancel {
		exch.cancelOrderSlice(slice)
		report.SlicesCancelled = append(report.SlicesCancelled, slice)
	}
	return report
}

func (exch *Exchange) newOrderSlice(slice *models.OrderSlice, asks []models.PriceQuantity) OrderResponse {
	if exch.meetFillCriteria(slice, asks) {
		return ResponseFilled
	}
	exch.pendingOrderSlices[slice] = 1
	return ResponseQueued
}
func (exch *Exchange) cancelOrderSlice(slice *models.OrderSlice) {
	delete(exch.pendingOrderSlices, slice)
}
