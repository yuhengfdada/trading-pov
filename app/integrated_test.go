package app

import (
	"testing"
)

// passive following
func TestFollow(t *testing.T) {
	lines := setup(t, "follow.csv")

	order := makeFIXMsg("1", "1000", "10")
	engine.Order(order)

	sendEvents(t, lines)
}

// behind
func TestBehindMin(t *testing.T) {
	lines := setup(t, "behindmin.csv")

	order := makeFIXMsg("1", "1000", "10")
	engine.Order(order)

	engine.setVolume(5000)
	engine.setOrderFilledQuantity(100)

	sendEvents(t, lines)
}

// ahead
func TestBreachMax(t *testing.T) {
	lines := setup(t, "breachmax.csv")
	order := makeFIXMsg("1", "1000", "10")
	engine.Order(order)

	engine.setVolume(1000)
	engine.setOrderFilledQuantity(200)

	sendEvents(t, lines)
}

// trade
func TestTrade(t *testing.T) {
	lines := setup(t, "followAndTrade.csv")

	order := makeFIXMsg("1", "1000", "10")
	engine.Order(order)

	sendEvents(t, lines)
}

// passive fills
func TestFills(t *testing.T) {
	lines := setup(t, "followAndFill.csv")

	order := makeFIXMsg("1", "1000", "10")
	engine.Order(order)

	for _, line := range lines {
		if line[1] == "20000" {
			engine.setVolume(5000)
		}
		exchange.ReceiveEvent(line)
		engine.ReceiveEvent(line)
	}
}

// behind a lot, ordering all levels of ask
func TestLargeBehind(t *testing.T) {
	lines := setup(t, "largebehind.csv")
	order := makeFIXMsg("1", "100000", "100")
	engine.Order(order)

	engine.setVolume(20000)

	sendEvents(t, lines)
}

func TestRealData(t *testing.T) {
	lines := setup(t, "market_data.csv")

	order := makeFIXMsg("1", "10000", "10")
	engine.Order(order)

	sendEvents(t, lines)
}

func TestRealDataLargeOrder(t *testing.T) {
	lines := setup(t, "market_data.csv")

	order := makeFIXMsg("1", "400000", "10")
	engine.Order(order)

	sendEvents(t, lines)
}

func TestRealDataLateOrder(t *testing.T) {
	lines := setup(t, "market_data.csv")

	order := makeFIXMsg("1", "400000", "10")

	for _, line := range lines {
		if line[1] == "9169924" {
			engine.Order(order)
		}
		exchange.ReceiveEvent(line)
		engine.ReceiveEvent(line)
	}
}

// test bad input for robustness
func TestBadInput(t *testing.T) {
	lines := setup(t, "bad_numbers.csv")
	sendEvents(t, lines)
}
