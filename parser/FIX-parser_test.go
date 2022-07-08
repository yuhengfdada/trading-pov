package parser

import (
	"allen/trading-pov/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFIX(t *testing.T) {
	msg := "54=1; 40=1; 38=10000; 6404=10"
	exp := models.FIXOrder{
		Buy:           true,
		Quantity:      10000,
		POVTargetProp: 0.1,
	}
	act := ParseFIX(msg)
	assert.Equal(t, exp, act)
}
