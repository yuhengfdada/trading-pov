package app

import "allen/trading-pov/models"

type Algorithm interface {
	Process(e *Engine) *models.Execution
}
