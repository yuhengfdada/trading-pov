package main

import (
	"encoding/csv"
	"os"
)

func main() {
	f, _ := os.Open("numbers.csv")
	r := csv.NewReader(f)
	r.ReadAll()
}
