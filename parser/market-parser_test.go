package parser

import (
	"encoding/csv"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEvent(t *testing.T) {
	f, _ := os.Open("../datasets/numbers.csv")
	r := csv.NewReader(f)
	lines, _ := r.ReadAll()
	for _, line := range lines {
		res := ParseEvent(line)
		if line[0] == "T" {
			assert.Equal(t, "*models.Trade", reflect.TypeOf(res).String())
		} else {
			assert.Equal(t, "*models.Quote", reflect.TypeOf(res).String())
		}
	}
}
