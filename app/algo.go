package app

type Algorithm interface {
	Process(e *Engine)
}
