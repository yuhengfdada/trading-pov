package app

import "allen/trading-pov/models"

type Engine struct {
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) HandleQuote(quote *models.Quote) {

}

func (e *Engine) HandleTrade(trade *models.Trade) {

}
