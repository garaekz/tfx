package formfx

import "errors"

var (
	// ErrConfigNotSet is returned when a configuration is not set.
	ErrConfigNotSet = errors.New("formfx: configuration not set")
	// ErrOutOfBounds is returned when an index is out of bounds.
	ErrOutOfBounds = errors.New("formfx: index out of bounds")
	// ErrInvalidOption is returned when an option is invalid.
	ErrInvalidOption = errors.New("formfx: invalid option")
	// ErrInvalidKeyHandler is returned when a key handler is not set.
	ErrInvalidKeyHandler = errors.New("formfx: key handler not set")
	// ErrInvalidRenderer is returned when a renderer is not set.
	ErrInvalidRenderer = errors.New("formfx: renderer not set")
	// ErrCanceled is returned when an operation is canceled.
	ErrCanceled = errors.New("formfx: operation canceled")
	// ErrInvalidConfigType is returned when the provided configuration type is not supported.
	ErrInvalidConfigType = errors.New("formfx: invalid configuration type")
	// ErrInvalidPromptType is returned when the prompt type is not recognized.
	ErrInvalidPromptType = errors.New("formfx: invalid prompt type")
)
