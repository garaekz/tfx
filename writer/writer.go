package writer

import "io"

// Writer defines the minimal interface required by rendering components.
// It combines io.Writer with a Flush method for buffered implementations.
type Writer interface {
	io.Writer
	Flush() error
}
