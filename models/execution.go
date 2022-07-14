package models

type Execution struct {
	SlicesToOrder  []*OrderSlice
	SlicesToCancel []*OrderSlice
}
type ExecutionReport struct {
	SlicesFilled    []*OrderSlice
	SlicesCancelled []*OrderSlice
	SlicesQueued    []*OrderSlice
}
