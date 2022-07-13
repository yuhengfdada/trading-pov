package models

import (
	"errors"
	"strconv"
	"strings"
)

type PriceQuantity struct {
	Price    float64
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

func NewPriceQuantity(price float64, quantity int) PriceQuantity {
	return PriceQuantity{
		Price:    price,
		Quantity: quantity,
	}
}

func NewTrade(time, price, quantity string) (*Trade, error) {
	t, err := strconv.Atoi(time)
	if err != nil {
		return nil, err
	}
	p, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return nil, err
	}
	q, err := strconv.Atoi(quantity)
	if err != nil {
		return nil, err
	}
	return &Trade{
		Time: t,
		PQ:   NewPriceQuantity(p, q),
	}, nil
}

func NewQuote(time, bids, asks string) (*Quote, error) {
	t, err := strconv.Atoi(time)
	if err != nil {
		return nil, err
	}
	quote := &Quote{
		Time: t,
	}
	err = FillPQs(&quote.Bids, bids)
	if err != nil {
		return nil, err
	}
	err = FillPQs(&quote.Asks, asks)
	if err != nil {
		return nil, err
	}
	return quote, nil
}

func FillPQs(arr *[]PriceQuantity, str string) error {
	nums := strings.Split(str, " ")
	if len(nums)%2 == 1 {
		return errors.New("bad bid or ask")
	}
	for i := 0; i < len(nums); i += 2 {
		p, err := strconv.ParseFloat(nums[i], 64)
		if err != nil {
			return err
		}
		q, err := strconv.Atoi(nums[i+1])
		if err != nil {
			return err
		}
		*arr = append(*arr, NewPriceQuantity(p, q))
	}
	return nil
}
