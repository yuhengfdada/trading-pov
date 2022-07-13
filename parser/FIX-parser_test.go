package parser

import (
	"allen/trading-pov/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFIX(t *testing.T) {
	msg := "54=1; 40=1; 38=10000; 6404=10"
	exp := &models.FIXOrder{
		Buy:           true,
		Quantity:      10000,
		POVTargetProp: 0.1,
	}
	act, err := ParseFIX(msg)
	if err != nil {
		t.Fail()
	}
	assert.Equal(t, exp, act)
	msg = "54=0; 40=1; 38=10000; 6404=10"
	_, err = ParseFIX(msg)
	assert.Error(t, err)
	msg = "54=1; 40=1; 38=0000; 6404=10"
	_, err = ParseFIX(msg)
	assert.Error(t, err)
	msg = "54=0; 40=1; 38=10000; 6404=100"
	_, err = ParseFIX(msg)
	assert.Error(t, err)
	msg = "40=1; 38=10000; 6404=100"
	_, err = ParseFIX(msg)
	assert.Error(t, err)
	msg = "54=0; 40=1; 38=asas; 6404=100"
	_, err = ParseFIX(msg)
	assert.Error(t, err)
}
