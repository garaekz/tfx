package form

import "errors"

// ErrCanceled is returned when a form input operation is canceled by the user.
var ErrCanceled = errors.New("form: canceled")
