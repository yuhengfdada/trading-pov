package util

import (
	"allen/trading-pov/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundFloat(t *testing.T) {
	assert.Equal(t, RoundFloat(0.080000000003), 0.08)
	assert.Equal(t, RoundFloat(0.11999999993), 0.12)
}

func TestMapToString(t *testing.T) {
	mp := make(map[float64]int)
	MapToString(mp)
	MapToString(nil)
	mp[10] = 500
	assert.Equal(t, "500@10.000000", MapToString(mp))
}

func TestEventToString(t *testing.T) {
	EventToString([]string{"111"})
	EventToString(nil)
	assert.Equal(t, "Quote@10000, Bids: 10.0 5000 9.9 4000 9.8 2000, Asks: 10.1 2000 10.2 10000", EventToString([]string{"Q", "10000", "10.0 5000 9.9 4000 9.8 2000", "10.1 2000 10.2 10000"}))
	assert.Equal(t, "Trade@20000, Price: 10.0, Quantity: 2000", EventToString([]string{"T", "20000", "10.0", "2000"}))
}

func TestOrderToString(t *testing.T) {
	order := &models.Order{
		QuantityTotal:  1000,
		QuantityFilled: 0,
		TargetRate:     0.1,
		MinRate:        0.08,
		MaxRate:        0.12,
	}
	exp := "Total Quantity: 1000, Target Rate: 0.100000, Min Rate: 0.080000, Max Rate: 0.120000"
	assert.Equal(t, exp, OrderToString(order))
}
