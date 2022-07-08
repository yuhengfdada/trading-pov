package models

type FIXOrder struct {
	Buy           bool
	Quantity      int
	POVTargetProp float64
}
