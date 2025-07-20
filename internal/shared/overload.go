package shared

import "fmt"

// Overload tries to cast the first value in 'have' into type T.
// If none is provided, returns the fallback.
func Overload[T any](have []any, fallback T) T {
	if len(have) == 0 {
		return fallback
	}
	if len(have) > 1 {
		panic("overload: only one argument expected")
	}
	if casted, ok := have[0].(T); ok {
		return casted
	}
	if ptr, ok := have[0].(*T); ok {
		return *ptr
	}
	panic(fmt.Sprintf("overload: expected type %T, got %T", fallback, have[0]))
}

// OverloadWithOptions merges overload fallback + functional options.
func OverloadWithOptions[T any](have []any, fallback T, opts ...Option[T]) T {
	instance := Overload(have, fallback)
	ApplyOptions(&instance, opts...)
	return instance
}
