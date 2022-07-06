package main

import (
	"encoding/csv"
	"os"
	"testing"
)

func MainTest(t *testing.T) {
	f, _ := os.Open("numbers2.csv")
	r := csv.NewReader(f)
	EntryPoint(r)
}
