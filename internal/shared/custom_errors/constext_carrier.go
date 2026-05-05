package custom_errors

type ContextCarrier interface {
	ContextData() map[string]any
}
