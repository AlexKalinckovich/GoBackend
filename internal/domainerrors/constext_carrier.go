package domainerrors

type ContextCarrier interface {
	ContextData() map[string]any
}
