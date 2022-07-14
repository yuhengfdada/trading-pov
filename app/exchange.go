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
	for slice := range exch.pendingOrderSlices {
		if exch.meetFillCriteria(slice) {
			report.SlicesFilled = append(report.SlicesFilled, slice)
			delete(exch.pendingOrderSlices, slice) // deletion won't affect the iteration
		}
	}
	return report
}

func (exch *Exchange) ReceiveExecutions(execution *models.Execution) *models.ExecutionReport {
	if execution == nil {
		return nil
	}
	report := &models.ExecutionReport{}
	for _, slice := range execution.SlicesToOrder {
		resp := exch.newOrderSlice(slice)
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

func (exch *Exchange) newOrderSlice(slice *models.OrderSlice) OrderResponse {
	if exch.meetFillCriteria(slice) {
		return ResponseFilled
	}
	exch.pendingOrderSlices[slice] = 1
	return ResponseQueued
}
func (exch *Exchange) cancelOrderSlice(slice *models.OrderSlice) {
	delete(exch.pendingOrderSlices, slice)
}
