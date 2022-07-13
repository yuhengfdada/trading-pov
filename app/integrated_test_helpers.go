package app

import (
	"allen/trading-pov/util"
	"encoding/csv"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	exchange *Exchange
	algo     Algorithm
	engine   *Engine
)

func setup(t *testing.T, dataset string) [][]string {
	f, err := os.Open("../datasets/" + dataset)
	if err != nil {
		t.Error(err)
	}
	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	if err != nil {
		t.Error(err)
	}
	exchange = NewExchange()
	algo = NewPOVAlgorithm()
	engine = NewEngine(exchange, algo)
	exchange.SetEngine(engine)
	return lines
}

func makeFIXMsg(buy, quantity, percentage string) string {
	res := "54=" + buy + "; 40=1; 38=" + quantity + "; 6404=" + percentage
	return res
}

func sendEvents(t *testing.T, lines [][]string) {
	for _, line := range lines {
		exchange.ReceiveEvent(line)
		engine.ReceiveEvent(line)
		checkInvariants(t, engine, exchange)
	}
}

func checkInvariants(t *testing.T, e *Engine, exch *Exchange) {
	if e.order == nil || e.order.QuantityFilled == e.order.QuantityTotal {
		return
	}
	// 1. Pending order slices should:
	//   1.1 match the PQ view.
	//   1.2 should not match fill criteria.
	//   1.3 two pending slice sets in exchange and engine should be the same.
	//   1.4 pending slices should not meet fill criteria.
	PQView := make(map[float64]int)

	for slice := range e.pendingOrderSlices {
		PQView[slice.Price] += slice.Quantity
		if slice.Price >= e.currentQuote.Asks[0].Price {
			t.Fail()
		}
		if _, ok := exch.pendingOrderSlices[slice]; !ok {
			t.Fail()
		}
	}
	for slice := range exch.pendingOrderSlices {
		if _, ok := e.pendingOrderSlices[slice]; !ok {
			t.Fail()
		}
		if exch.meetFillCriteria(slice) {
			t.Fail()
		}
	}
	assert.Equal(t, util.MapToString(e.pendingOrderPQView), util.MapToString(PQView))
}
