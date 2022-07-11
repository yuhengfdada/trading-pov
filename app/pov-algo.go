package app

import (
	"allen/trading-pov/models"
	"allen/trading-pov/util"
	"fmt"
	"math"
)

type POVAlgorithm struct {
}

func NewPOVAlgorithm() *POVAlgorithm {
	return &POVAlgorithm{}
}

func (algo *POVAlgorithm) Process(e *Engine) {

	x := float64(e.order.QuantityFilled)
	y := float64(e.volume)

	if util.RoundFloat(x) < util.RoundFloat(y*e.order.MinRate) { // behind
		fmt.Printf("POV-Algo: We are behind. Creating aggressive slices...\n")
		quantityToOrder := util.Min(int(math.Round(y*e.order.MinRate-x)), e.order.QuantityTotal-e.order.QuantityFilled)
		// only create one slice @ best ask for now. Will change to use 3 levels.
		for _, pq := range e.currentQuote.Asks {
			if quantityToOrder <= 0 {
				break
			}
			resp := e.NewOrderSlice(
				&models.OrderSlice{
					TimeStamp: e.currentTime,
					Quantity:  util.Min(quantityToOrder, pq.Quantity),
					Price:     pq.Price,
				})
			if resp != ResponseFilled {
				panic("should be filled")
			}
			quantityToOrder -= pq.Quantity
		}

	} else if util.RoundFloat(x) > util.RoundFloat(y*e.order.MaxRate) { // ahead
		fmt.Printf("POV-Algo: We are ahead. Cancelling all slices...\n")
		for slice := range e.pendingOrderSlices {
			e.cancelOrderSlice(slice)
		}
		return
	}
	// follow
	fmt.Printf("POV-Algo: Rebalancing passive order slices...\n")
	quantityLeft := e.order.QuantityTotal - e.order.QuantityFilled
	for _, pq := range e.currentQuote.Bids {
		price := pq.Price
		quantityThreshold := int(math.Round(float64(pq.Quantity) * e.order.TargetRate))
		quantityPending := e.pendingOrderPQView[price] // TODO: Add nil check
		quantityNeeded := quantityThreshold - quantityPending
		// 1. pending order exceeds (this price's tot Qty * Target Rate), due to quote changes
		// 2. pending order exceeds quantityLeft, due to new fills
		if quantityNeeded < 0 || quantityPending > quantityLeft {
			e.cancelAllSlicesWithPrice(price)
			quantityNeeded = quantityThreshold
		}
		// We're not ordering 0 Qty
		if quantityNeeded == 0 || quantityLeft == 0 || quantityLeft == quantityPending {
			newQuantityPending := e.pendingOrderPQView[price]
			quantityLeft -= newQuantityPending
			continue
		}
		if quantityLeft < quantityNeeded { // Qty left to order is too little (but not ordering 0 Qty)
			e.NewOrderSlice(
				&models.OrderSlice{
					TimeStamp: e.currentTime,
					Quantity:  quantityLeft,
					Price:     price,
				})
		} else { // Order a full level (but not ordering 0 Qty)
			e.NewOrderSlice(
				&models.OrderSlice{
					TimeStamp: e.currentTime,
					Quantity:  quantityNeeded,
					Price:     price,
				})
		}
		newQuantityPending := e.pendingOrderPQView[price]
		quantityLeft -= newQuantityPending
	}
	// cancel those orders at no-more-existent prices
	e.cancelNoMoreExistentPriceSlices()
}
