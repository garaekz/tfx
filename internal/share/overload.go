package share

import "fmt"

// Overload[T] provides strict one-argument coercion into T.
//
// - If no value is passed, fallback is used.
// - If one value is passed, it MUST be T or *T.
// - Any other type (including nil, map, or unrelated struct) triggers panic.
//
// This is NOT a soft cast. It’s a strict dispatch helper for multipath entry.
func Overload[T any](have []any, fallback T) T {
	if len(have) == 0 {
		return fallback
	}
	if len(have) > 1 {
		panic("overload: expected 0 or 1 argument")
	}

	arg := have[0]
	var zero T // we only use this to format the error message

	switch v := arg.(type) {
	case T:
		return v
	case *T:
		return *v
	default:
		panic(fmt.Sprintf(
			"overload: expected type %T or *%T, got %T",
			zero, &zero, arg,
		))
	}
}

// OverloadWithOptions[T] combines flexible overload + option injection.
//
// - If no value is provided, fallback is used.
// - If one value is provided, it MUST be T or *T, or panics.
// - Any additional functional options are applied to the result.
//
// Equivalent to:
//
//     instance := Overload(have, fallback)
//     ApplyOptions(&instance, opts...)
//

func OverloadWithOptions[T any](have []any, fallback T, opts ...Option[T]) T {
	instance := Overload(have, fallback)
	ApplyOptions(&instance, opts...)
	return instance
}
