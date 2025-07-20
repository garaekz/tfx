package share

// Option is a functional setter for any struct T.
// Example: func WithText(txt string) Option[Config]
type Option[T any] func(*T)

// ApplyOptions applies a set of options to a given instance.
func ApplyOptions[T any](target *T, opts ...Option[T]) {
	for _, opt := range opts {
		opt(target)
	}
}
