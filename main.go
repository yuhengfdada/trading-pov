package main

import (
	"allen/trading-pov/app"
	"allen/trading-pov/models"
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	f, _ := os.Open("numbers.csv")
	r := csv.NewReader(f)
	EntryPoint(r)
}

func EntryPoint(r *csv.Reader) {
	engine := app.NewEngine()

	lines, _ := r.ReadAll()
	for _, line := range lines {
		fmt.Println(line)
		if line[0] == "T" {
			t := models.NewTrade(line[1], line[2], line[3])
			fmt.Println(t)
			engine.HandleTrade(t)
		} else {
			q := models.NewQuote(line[1], line[2], line[3])
			fmt.Println(q)
			engine.HandleQuote(q)
		}
	}
}
