package models

type Order struct {
	QuantityTotal  int
	QuantityFilled int
	TargetRate     float64
	MinRate        float64
	MaxRate        float64
}

type OrderSlice struct {
	TimeStamp int
	Quantity  int
	Price     float64
}
