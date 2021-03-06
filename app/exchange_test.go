package app

import (
	"allen/trading-pov/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMeetFillCriteria(t *testing.T) {
	exchange = NewExchange()
	exchange.currentQuote = &models.Quote{}
	exchange.currentQuote.Bids = append(exchange.currentQuote.Bids, models.PriceQuantity{Price: 100})
	exchange.currentQuote.Asks = append(exchange.currentQuote.Asks, models.PriceQuantity{Price: 1000, Quantity: 100000})
	exchange.currentTime = 180000
	var copyAsks []models.PriceQuantity
	copyAsks = append(copyAsks, exchange.currentQuote.Asks...)
	assert.True(t, exchange.meetFillCriteria(&models.OrderSlice{TimeStamp: 0, Price: 1000}, copyAsks))
	assert.True(t, exchange.meetFillCriteria(&models.OrderSlice{TimeStamp: 0, Price: 100}, copyAsks))
	assert.False(t, exchange.meetFillCriteria(&models.OrderSlice{TimeStamp: 1, Price: 100}, copyAsks))
}

func TestReceiveEventExchange(t *testing.T) {
	exchange = NewExchange()
	exchange.ReceiveEvent([]string{"Q", "10000", "10.0 5000 9.9 4000 9.8 2000", "10.1 2000 10.2 10000"})
	assert.Equal(t, len(exchange.currentQuote.Asks), 2)
	assert.Equal(t, len(exchange.currentQuote.Bids), 3)
	assert.Equal(t, exchange.currentTime, 10000)
	assert.Equal(t, exchange.currentQuote.Bids[0].Price, 10.0)
	assert.Equal(t, exchange.currentQuote.Bids[0].Quantity, 5000)
	assert.Equal(t, exchange.currentQuote.Asks[0].Price, 10.1)
	assert.Equal(t, exchange.currentQuote.Asks[0].Quantity, 2000)
	exchange.ReceiveEvent([]string{"T", "20000", "10.0", "2000"})
	assert.Equal(t, exchange.currentTime, 20000)
}

// When asks are not enough, fill earlier slices first.
func TestProcessPartialFills(t *testing.T) {
	exchange = NewExchange()
	exchange.ReceiveEvent([]string{"Q", "10000", "10.0 5000 9.9 4000 9.8 2000", "10.1 2000 10.2 10000"})
	slice1 := &models.OrderSlice{TimeStamp: 0, Quantity: 100, Price: 10}
	slice2 := &models.OrderSlice{TimeStamp: 1, Quantity: 100, Price: 10}
	slice3 := &models.OrderSlice{TimeStamp: 2, Quantity: 100, Price: 10}
	exchange.pendingOrderSlices[slice1] = 1
	exchange.pendingOrderSlices[slice2] = 1
	exchange.pendingOrderSlices[slice3] = 1
	report := exchange.ReceiveEvent([]string{"Q", "10000", "9.8 2000", "10.0 200"})
	assert.Len(t, report.SlicesFilled, 2)
	assert.Equal(t, report.SlicesFilled[0], slice1)
	assert.Equal(t, report.SlicesFilled[1], slice2)
}

func TestReceiveExecutions(t *testing.T) {
	exchange = NewExchange()
	execution := &models.Execution{}
	exchange.ReceiveEvent([]string{"Q", "10000", "10.0 5000 9.9 4000 9.8 2000", "10.1 2000 10.2 10000"})

	slice1 := &models.OrderSlice{Quantity: 111, Price: 11.1}
	slice2 := &models.OrderSlice{Quantity: 111, Price: 9.9}
	slice3 := &models.OrderSlice{Quantity: 222, Price: 22.2}

	execution.SlicesToOrder = append(execution.SlicesToOrder, slice1)
	execution.SlicesToOrder = append(execution.SlicesToOrder, slice2)
	execution.SlicesToCancel = append(execution.SlicesToCancel, slice3)

	report := exchange.ReceiveExecutions(execution)
	assert.Equal(t, report.SlicesCancelled[0], slice3)
	assert.Equal(t, report.SlicesFilled[0], slice1)
	assert.Equal(t, report.SlicesQueued[0], slice2)
}
