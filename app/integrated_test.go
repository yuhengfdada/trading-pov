package app

import (
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

	order := makeFIXMsg("1", "10000", "10")
	engine.Order(order)

	sendEvents(lines)
}

func TestRealDataLargeOrder(t *testing.T) {
	lines := setup(t, "market_data.csv")

	order := makeFIXMsg("1", "400000", "10")
	engine.Order(order)

	sendEvents(lines)
}
