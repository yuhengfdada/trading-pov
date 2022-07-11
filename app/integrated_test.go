package app

import (
	"allen/trading-pov/models"
	"testing"
)

func TestFollow(t *testing.T) {
	lines := setup(t, "follow.csv")

	order := makeFIXMsg("1", "1000", "10")
	engine.Order(order)

	sendEvents(lines)
}

func TestBehindMin(t *testing.T) {
	lines := setup(t, "behindmin.csv")

	order := makeFIXMsg("1", "1000", "10")
	engine.Order(order)

	engine.SetVolume(5000)
	engine.SetOrderFilledQuantity(100)

	sendEvents(lines)
}

func TestBreachMax(t *testing.T) {
	lines := setup(t, "breachmax.csv")
	order := makeFIXMsg("1", "1000", "10")
	engine.Order(order)

	engine.SetVolume(1000)
	engine.SetOrderFilledQuantity(200)

	sendEvents(lines)
}

func TestTrade(t *testing.T) {
	lines := setup(t, "followAndTrade.csv")

	order := makeFIXMsg("1", "1000", "10")
	engine.Order(order)

	sendEvents(lines)
}

func TestFills(t *testing.T) {
	lines := setup(t, "followAndFill.csv")

	order := makeFIXMsg("1", "1000", "10")
	engine.Order(order)

	for _, line := range lines {
		if line[1] == "20000" {
			engine.SetVolume(5000)
		}
		exchange.ReceiveEvent(line)
		engine.ReceiveEvent(line)
	}
}

func TestRealData(t *testing.T) {
	lines := setup(t, "market_data.csv")

	order := makeFIXMsg("1", "20000", "10")
	engine.Order(order)

	for _, line := range lines {
		exchange.ReceiveEvent(line)
		engine.ReceiveEvent(line)
	}
}

func TestXxx(t *testing.T) {
	lines := setup(t, "mock.csv")
	slice1 := &models.OrderSlice{
		TimeStamp: 5532248,
		Quantity:  190,
		Price:     52.5,
	}
	slice2 := &models.OrderSlice{
		TimeStamp: 5532248,
		Quantity:  1400,
		Price:     52.55,
	}
	slice3 := &models.OrderSlice{
		TimeStamp: 5532248,
		Quantity:  850,
		Price:     52.6,
	}
	engine.pendingOrderSlices[slice1] = 1
	engine.pendingOrderSlices[slice2] = 1
	engine.pendingOrderSlices[slice3] = 1

	engine.pendingOrderPQView[52.5] = 190
	engine.pendingOrderPQView[52.55] = 1400
	engine.pendingOrderPQView[52.6] = 850

	exchange.pendingOrderSlices[slice1] = 1
	exchange.pendingOrderSlices[slice2] = 1
	exchange.pendingOrderSlices[slice3] = 1

	order := makeFIXMsg("1", "2440", "10")
	engine.Order(order)

	for _, line := range lines {
		exchange.ReceiveEvent(line)
		engine.ReceiveEvent(line)
	}
}
