package form

// Option is a functional option for configuring form components.
type Option[T any] func(*T)

// ApplyOptions applies a list of options to a configuration struct.
func ApplyOptions[T any](cfg *T, opts ...Option[T]) {
	for _, opt := range opts {
		opt(cfg)
	}
}
