package app

import (
	"allen/trading-pov/models"
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

	if x < y*e.order.MinRate { // behind
		quantity := int(math.Round(y*e.order.MinRate - x))
		// only create one slice @ best ask for now. Will change to use 3 levels.
		resp := e.NewOrderSlice(
			&models.OrderSlice{
				TimeStamp: e.currentTime,
				Quantity:  quantity,
				Price:     e.currentQuote.Asks[0].Price,
			})
		if resp != ResponseFilled {
			panic("should be filled")
		}
	} else if x > y*e.order.MaxRate { // ahead
		for slice := range e.pendingOrderSlices {
			e.exchange.CancelOrderSlice(slice)
			delete(e.pendingOrderSlices, slice)
		}
	} else { // follow
		quantityLeft := e.order.QuantityTotal - e.order.QuantityFilled
		for _, pq := range e.currentQuote.Bids {
			price := pq.Price
			quantityThreshold := int(math.Round(float64(pq.Quantity) * e.order.TargetRate))
			quantityPending := e.pendingOrderPQView[price] // TODO: Add nil check
			quantityNeeded := quantityThreshold - quantityPending
			if quantityNeeded < 0 {
				e.cancelAllSlicesWithPrice(price)
				quantityNeeded = quantityThreshold
			}
			if quantityLeft <= quantityNeeded {
				e.NewOrderSlice(
					&models.OrderSlice{
						TimeStamp: e.currentTime,
						Quantity:  quantityLeft,
						Price:     price,
					})
				break
			}
			e.NewOrderSlice(
				&models.OrderSlice{
					TimeStamp: e.currentTime,
					Quantity:  quantityNeeded,
					Price:     price,
				})
			quantityLeft -= quantityNeeded
		}
	}
}
