package app

import (
	"allen/trading-pov/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReceiveEvent(t *testing.T) {
	algo = NewPOVAlgorithm()
	engine = NewEngine(algo)
	engine.ReceiveEvent([]string{"Q", "10000", "10.0 5000 9.9 4000 9.8 2000", "10.1 2000 10.2 10000"})

	assert.Equal(t, engine.currentTime, 10000)
	assert.Equal(t, engine.volume, 0)
	assert.Equal(t, len(engine.currentQuote.Asks), 2)
	assert.Equal(t, len(engine.currentQuote.Bids), 3)
	assert.Equal(t, engine.currentQuote.Bids[0].Price, 10.0)
	assert.Equal(t, engine.currentQuote.Bids[0].Quantity, 5000)
	assert.Equal(t, engine.currentQuote.Asks[0].Price, 10.1)
	assert.Equal(t, engine.currentQuote.Asks[0].Quantity, 2000)

	engine.ReceiveEvent([]string{"T", "20000", "10.0", "2000"})
	assert.Equal(t, engine.currentTime, 20000)
	assert.Equal(t, engine.volume, 2000)
}

func TestReceiveReport(t *testing.T) {
	algo = NewPOVAlgorithm()
	engine = NewEngine(algo)
	engine.order = &models.Order{}

	report := &models.ExecutionReport{}
	report.SlicesFilled = append(report.SlicesFilled, &models.OrderSlice{Quantity: 111, Price: 11.1})
	engine.pendingOrderPQView[22.2] = 444
	report.SlicesCancelled = append(report.SlicesCancelled, &models.OrderSlice{Quantity: 222, Price: 22.2})
	pendingSlice := &models.OrderSlice{Quantity: 333, Price: 33.3}
	report.SlicesQueued = append(report.SlicesQueued, pendingSlice)
	engine.ReceiveReport(report)

	assert.Equal(t, engine.order.QuantityFilled, 111)
	assert.Equal(t, engine.pendingOrderPQView[22.2], 222)
	assert.Contains(t, engine.pendingOrderSlices, pendingSlice)
}
