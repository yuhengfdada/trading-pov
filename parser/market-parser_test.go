package parser

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEvent(t *testing.T) {
	f, err := os.Open("../datasets/bad_numbers.csv")
	if err != nil {
		t.Error(err)
	}
	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	if err != nil {
		t.Error(err)
	}
	for _, line := range lines {
		res, err := ParseEvent(line)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if line[0] == "T" {
			assert.Equal(t, "*models.Trade", reflect.TypeOf(res).String())
		} else {
			assert.Equal(t, "*models.Quote", reflect.TypeOf(res).String())
		}
	}
}
