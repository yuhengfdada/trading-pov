package main

import (
	"allen/trading-pov/app"
	"encoding/csv"
	"os"
	"testing"
)

var (
	exchange *app.Exchange
	algo     app.Algorithm
	engine   *app.Engine
)

func setup(t *testing.T, dataset string) [][]string {
	f, err := os.Open("datasets/" + dataset)
	if err != nil {
		t.Error(err)
	}
	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	if err != nil {
		t.Error(err)
	}
	exchange = app.NewExchange()
	algo = app.NewPOVAlgorithm()
	engine = app.NewEngine(exchange, algo)
	exchange.SetEngine(engine)
	return lines
}

func makeFIXMsg(buy, quantity, percentage string) string {
	res := "54=" + buy + "; 40=1; 38=" + quantity + "; 6404=" + percentage
	return res
}

func TestFollow(t *testing.T) {
	lines := setup(t, "follow.csv")

	order := makeFIXMsg("1", "1000", "10")
	engine.Order(order)

	for _, line := range lines {
		exchange.ReceiveEvent(line)
		engine.ReceiveEvent(line)
	}
}

func TestBehindMin(t *testing.T) {
	lines := setup(t, "behindmin.csv")

	order := makeFIXMsg("1", "1000", "10")
	engine.Order(order)

	engine.SetVolume(5000)
	engine.SetOrderFilledQuantity(100)

	for _, line := range lines {
		exchange.ReceiveEvent(line)
		engine.ReceiveEvent(line)
	}
}

func TestBreachMax(t *testing.T) {
	lines := setup(t, "breachmax.csv")
	order := makeFIXMsg("1", "1000", "10")
	engine.Order(order)

	engine.SetVolume(1000)
	engine.SetOrderFilledQuantity(200)

	for _, line := range lines {
		exchange.ReceiveEvent(line)
		engine.ReceiveEvent(line)
	}
}
