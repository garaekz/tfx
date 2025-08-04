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

// OverloadWithOptions[T] handles multipath with strict rules:
//
// - If no args: use fallback
// - Can mix: config struct (T or *T) + functional options (Option[T])
// - Only one config struct allowed
// - All types must match T
func OverloadWithOptions[T any](args []any, fallback T) T {
	if len(args) == 0 {
		return fallback
	}

	var holder *T
	var options []Option[T]

	for _, arg := range args {
		switch v := arg.(type) {
		case T:
			if holder != nil {
				panic("OverloadWithOptions: multiple config structs not allowed")
			}
			temp := v
			holder = &temp
		case *T:
			if holder != nil {
				panic("OverloadWithOptions: multiple config structs not allowed")
			}
			holder = v
		case Option[T]:
			options = append(options, v)
		default:
			panic(fmt.Sprintf(
				"OverloadWithOptions: expected type %T, *%T, or Option[%T], got %T",
				fallback, &fallback, fallback, arg,
			))
		}
	}

	// Determine the base instance
	var instance T
	if holder != nil {
		instance = *holder
	} else {
		instance = fallback
	}

	// Apply all functional options
	ApplyOptions(&instance, options...)
	return instance
}
