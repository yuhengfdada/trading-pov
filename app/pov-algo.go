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

func (algo *POVAlgorithm) Process(e *Engine) *models.Execution {
	execution := &models.Execution{}

	x := float64(e.order.QuantityFilled)
	y := float64(e.volume)

	qFilled := util.RoundFloat(x)
	vTraded := util.RoundFloat(y)
	behindThres := util.RoundFloat(y * e.order.MinRate)
	aheadThres := util.RoundFloat(y * e.order.MaxRate)

	newlyFilledQuantity := 0

	fmt.Printf("POV-Algo: Cumulative quantity: %v, Volume traded: %.0f, Behind threshold: %v, Ahead threshold: %v\n", qFilled, vTraded, behindThres, aheadThres)
	if qFilled < behindThres { // behind
		quantityToOrder := util.Min(int(math.Round(behindThres-qFilled)), e.order.QuantityTotal-e.order.QuantityFilled)
		newlyFilledQuantity = algo.behind(e, execution, quantityToOrder)
	} else if qFilled > aheadThres { // ahead
		algo.ahead(e, execution)
		return execution
	}
	// follow
	algo.follow(e, execution, int(qFilled)+newlyFilledQuantity)
	return execution
}

func (algo *POVAlgorithm) follow(e *Engine, execution *models.Execution, quantityFilled int) {
	fmt.Printf("POV-Algo: Rebalancing passive order slices...\n")
	quantityLeft := e.order.QuantityTotal - quantityFilled

	tempPendingOrderPQView := make(map[float64]int)
	for k, v := range e.pendingOrderPQView {
		tempPendingOrderPQView[k] = v
	}

	for _, pq := range e.currentQuote.Bids {
		price := pq.Price
		quantityThreshold := int(math.Round(float64(pq.Quantity) * e.order.TargetRate))
		quantityPending := tempPendingOrderPQView[price]
		quantityNeeded := util.Min(quantityThreshold-quantityPending, quantityLeft-quantityPending)

		// 1. pending order exceeds (this price's tot Qty * Target Rate), due to quote changes
		// 2. pending order exceeds quantityLeft, due to new fills
		// 100 100 200 9.9    200
		//
		if quantityNeeded < 0 {
			e.cancelAllSlicesWithPrice(price, execution)
			tempPendingOrderPQView[price] = 0
			quantityNeeded = quantityThreshold
		}
		// We're not ordering 0 Qty, but we still update the quantity left to order
		if quantityNeeded == 0 || quantityLeft == 0 {
			newQuantityPending := tempPendingOrderPQView[price]
			quantityLeft -= newQuantityPending
			continue
		}
		if quantityLeft < quantityNeeded { // Qty left to order is too little
			e.newOrderSlice(
				&models.OrderSlice{
					TimeStamp: e.currentTime,
					Quantity:  quantityLeft,
					Price:     price,
				}, execution)
			tempPendingOrderPQView[price] += quantityLeft
		} else { // Order a full level
			e.newOrderSlice(
				&models.OrderSlice{
					TimeStamp: e.currentTime,
					Quantity:  quantityNeeded,
					Price:     price,
				}, execution)
			tempPendingOrderPQView[price] += quantityNeeded
		}
		newQuantityPending := tempPendingOrderPQView[price]
		quantityLeft -= newQuantityPending
	}
	// cancel those orders at no-more-existent prices (due to quote changes)
	e.cancelNoMoreExistentPriceSlices(execution)
}

func (algo *POVAlgorithm) behind(e *Engine, execution *models.Execution, quantityToOrder int) int {
	fmt.Printf("POV-Algo: We are behind. Creating aggressive slices...\n")
	res := 0
	for _, pq := range e.currentQuote.Asks {
		if quantityToOrder <= 0 {
			break
		}
		e.newOrderSlice(
			&models.OrderSlice{
				TimeStamp: e.currentTime,
				Quantity:  util.Min(quantityToOrder, pq.Quantity),
				Price:     pq.Price,
			}, execution)
		res += util.Min(quantityToOrder, pq.Quantity)
		quantityToOrder -= pq.Quantity
	}
	return res
}

func (algo *POVAlgorithm) ahead(e *Engine, execution *models.Execution) {
	fmt.Printf("POV-Algo: We are ahead. Cancelling all slices...\n")
	for slice := range e.pendingOrderSlices {
		e.cancelOrderSlice(slice, execution)
	}
}
