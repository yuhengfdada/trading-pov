package models

import (
	"strconv"
	"strings"
)

type PriceQuantity struct {
	Price    int
	Quantity int
}
type Trade struct {
	Time int
	PQ   PriceQuantity
}

type Quote struct {
	Time int
	Bids []PriceQuantity
	Asks []PriceQuantity
}

func NewPriceQuantity(price, quantity int) PriceQuantity {
	return PriceQuantity{
		Price:    price,
		Quantity: quantity,
	}
}

func NewTrade(time, price, quantity string) *Trade {
	t, _ := strconv.Atoi(time)
	p, _ := strconv.Atoi(price)
	q, _ := strconv.Atoi(quantity)
	return &Trade{
		Time: t,
		PQ:   NewPriceQuantity(p, q),
	}
}

func NewQuote(time, bids, asks string) *Quote {
	t, _ := strconv.Atoi(time)
	quote := &Quote{
		Time: t,
	}
	FillPQs(&quote.Bids, bids)
	FillPQs(&quote.Asks, asks)
	return quote
}

func FillPQs(arr *[]PriceQuantity, str string) {
	nums := strings.Split(str, " ")
	for i := 0; i < len(nums); i += 2 {
		p, _ := strconv.Atoi(nums[i])
		q, _ := strconv.Atoi(nums[i+1])
		*arr = append(*arr, NewPriceQuantity(p, q))
	}
}
