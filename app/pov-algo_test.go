package app

import (
	"allen/trading-pov/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFollowUnit(t *testing.T) {
	algo = NewPOVAlgorithm()
	engine = NewEngine(algo)
	order := makeFIXMsg("1", "1000", "10")
	engine.Order(order)
	execution := engine.ReceiveEvent([]string{"Q", "10000", "10.0 5000 9.9 4000 9.8 2000", "10.1 2000 10.2 10000"})

	assert.Equal(t, models.OrderSlice{TimeStamp: 10000, Quantity: 500, Price: 10}, *(execution.SlicesToOrder[0]))
	assert.Equal(t, models.OrderSlice{TimeStamp: 10000, Quantity: 400, Price: 9.9}, *(execution.SlicesToOrder[1]))
	assert.Equal(t, models.OrderSlice{TimeStamp: 10000, Quantity: 100, Price: 9.8}, *(execution.SlicesToOrder[2]))
}

func TestBehindUnit(t *testing.T) {
	algo = NewPOVAlgorithm()
	engine = NewEngine(algo)
	engine.ReceiveEvent([]string{"Q", "10000", "10.0 5000 9.9 4000 9.8 2000", "10.1 2000 10.2 10000"})
	execution := &models.Execution{}
	algo.behind(engine, execution, 10000)
	assert.Equal(t, models.OrderSlice{TimeStamp: 10000, Quantity: 2000, Price: 10.1}, *(execution.SlicesToOrder[0]))
	assert.Equal(t, models.OrderSlice{TimeStamp: 10000, Quantity: 8000, Price: 10.2}, *(execution.SlicesToOrder[1]))

	execution = &models.Execution{}
	algo.behind(engine, execution, 100000)
	assert.Equal(t, models.OrderSlice{TimeStamp: 10000, Quantity: 2000, Price: 10.1}, *(execution.SlicesToOrder[0]))
	assert.Equal(t, models.OrderSlice{TimeStamp: 10000, Quantity: 10000, Price: 10.2}, *(execution.SlicesToOrder[1]))
}

func TestAheadUnit(t *testing.T) {
	algo = NewPOVAlgorithm()
	engine = NewEngine(algo)
	engine.addPendingOrderSlice(&models.OrderSlice{TimeStamp: 10000, Quantity: 111, Price: 11.1})
	engine.addPendingOrderSlice(&models.OrderSlice{TimeStamp: 10000, Quantity: 222, Price: 22.2})
	execution := &models.Execution{}
	algo.ahead(engine, execution)
	assert.Equal(t, models.OrderSlice{TimeStamp: 10000, Quantity: 111, Price: 11.1}, *(execution.SlicesToCancel[0]))
	assert.Equal(t, models.OrderSlice{TimeStamp: 10000, Quantity: 222, Price: 22.2}, *(execution.SlicesToCancel[1]))
}
