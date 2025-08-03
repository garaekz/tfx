package formfx

import "errors"

// ErrCanceled is returned when a formfx input operation is canceled by the user.
var ErrCanceled = errors.New("formfx: canceled")
