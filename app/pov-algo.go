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

	qFilled := util.RoundFloat(x)
	vTraded := util.RoundFloat(y)
	behindThres := util.RoundFloat(y * e.order.MinRate)
	aheadThres := util.RoundFloat(y * e.order.MaxRate)

	fmt.Printf("POV-Algo: Cumulative quantity: %v, Volume traded: %.0f, Behind threshold: %v, Ahead threshold: %v\n", qFilled, vTraded, behindThres, aheadThres)
	if qFilled < behindThres { // behind
		fmt.Printf("POV-Algo: We are behind. Creating aggressive slices...\n")
		quantityToOrder := util.Min(int(math.Round(behindThres-qFilled)), e.order.QuantityTotal-e.order.QuantityFilled)
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

	} else if qFilled > aheadThres { // ahead
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
		quantityPending := e.pendingOrderPQView[price]
		quantityNeeded := util.Min(quantityThreshold-quantityPending, quantityLeft-quantityPending)

		// 1. pending order exceeds (this price's tot Qty * Target Rate), due to quote changes
		// 2. pending order exceeds quantityLeft, due to new fills
		if quantityNeeded < 0 {
			e.cancelAllSlicesWithPrice(price)
			quantityNeeded = quantityThreshold
		}
		// We're not ordering 0 Qty, but we still update the quantity left to order
		if quantityNeeded == 0 || quantityLeft == 0 {
			newQuantityPending := e.pendingOrderPQView[price]
			quantityLeft -= newQuantityPending
			continue
		}
		if quantityLeft < quantityNeeded { // Qty left to order is too little
			e.NewOrderSlice(
				&models.OrderSlice{
					TimeStamp: e.currentTime,
					Quantity:  quantityLeft,
					Price:     price,
				})
		} else { // Order a full level
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
	// cancel those orders at no-more-existent prices (due to quote changes)
	e.cancelNoMoreExistentPriceSlices()
}
