package parser

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestParseEvent(t *testing.T) {
	f, _ := os.Open("../numbers.csv")
	r := csv.NewReader(f)
	lines, _ := r.ReadAll()
	for _, line := range lines {
		res := ParseEvent(line)
		fmt.Println(reflect.TypeOf(res))
	}
}
