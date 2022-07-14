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
		eventExecutionReport := exchange.ReceiveEvent(line)
		engine.ReceiveReport(eventExecutionReport)

		executions := engine.ReceiveEvent(line)
		orderExecutionReport := exchange.ReceiveExecutions(executions)
		engine.ReceiveReport(orderExecutionReport)

		checkInvariants(t, engine, exchange)
	}
}

// Ask quantity is not enough. Do not fill, even there is a slice with best ask price.
// Notice how the 2500@10 slice did not get filled at the beginning of time 20000
func TestPassiveFillNotEnoughAsk(t *testing.T) {
	lines := setup(t, "notEnoughAsk.csv")
	order := makeFIXMsg("1", "100000", "50")
	engine.Order(order)

	engine.setVolume(200000)

	sendEvents(t, lines)
}

func TestPassiveFillPartialNotEnoughAsk(t *testing.T) {
	lines := setup(t, "notEnoughAsk.csv")
	order := makeFIXMsg("1", "100000", "50")
	engine.Order(order)

	sendEvents(t, lines)
}

// behind a lot, ordering all levels of ask
func TestLargeBehind(t *testing.T) {
	lines := setup(t, "largebehind.csv")
	order := makeFIXMsg("1", "100000", "50")
	engine.Order(order)

	engine.setVolume(200000)

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

// test bad input for robustness
func TestBadInput(t *testing.T) {
	lines := setup(t, "bad_numbers.csv")
	sendEvents(t, lines)
}

func TestBadOrder(t *testing.T) {
	setup(t, "bad_numbers.csv")
	order := makeFIXMsg("1", "-400000", "10")
	engine.Order(order)
	order = makeFIXMsg("2", "40000", "10")
	engine.Order(order)
	order = makeFIXMsg("1", "400000", "0")
	engine.Order(order)
}
