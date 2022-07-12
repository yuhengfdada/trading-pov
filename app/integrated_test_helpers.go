package app

import (
	"encoding/csv"
	"os"
	"testing"
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

func sendEvents(lines [][]string) {
	for _, line := range lines {
		exchange.ReceiveEvent(line)
		engine.ReceiveEvent(line)
	}
}
