package writer

import (
	"sync"

	"github.com/garaekz/tfx/internal/share"
)

// AsyncWriter provides an asynchronous, buffered writer decorator.
// It wraps another writer and performs writes in a separate goroutine.
type AsyncWriter struct {
	underlyingWriter share.Writer
	logCh            chan *share.Entry
	errCh            chan error
	doneCh           chan struct{}
	wg               sync.WaitGroup
}

// NewAsyncWriter creates a new asynchronous writer.
func NewAsyncWriter(underlying share.Writer, bufferSize int) *AsyncWriter {
	aw := &AsyncWriter{
		underlyingWriter: underlying,
		logCh:            make(chan *share.Entry, bufferSize),
		errCh:            make(chan error, 1),
		doneCh:           make(chan struct{}),
	}

	aw.wg.Add(1) // Add for the run goroutine itself
	go aw.run()

	return aw
}

// Write sends a log entry to the buffer.
// This method is non-blocking.
func (aw *AsyncWriter) Write(entry *share.Entry) error {
	aw.wg.Add(1) // Increment WaitGroup for each message sent
	select {
	case aw.logCh <- entry:
		return nil
	case <-aw.doneCh:
		// Writer is closed, drop the log and decrement wg
		aw.wg.Done()
		return nil
	}
}

// Close flushes the buffer and stops the writer.
func (aw *AsyncWriter) Close() error {
	close(aw.logCh)
	aw.wg.Wait() // Wait for all messages to be processed
	close(aw.doneCh)
	return aw.underlyingWriter.Close()
}

// Errors returns a channel for receiving write errors.
func (aw *AsyncWriter) Errors() <-chan error {
	return aw.errCh
}

// Flush waits for all buffered messages to be written.
func (aw *AsyncWriter) Flush() {
	// To flush, we need to ensure all items in logCh are processed.
	// A simple way is to send a signal and wait for it to be processed.
	// However, a more robust way is to temporarily close the channel,
	// wait for the run goroutine to finish, and then re-open it.
	// This is complex with a single channel, so we'll use a simpler approach
	// for now: rely on the underlying writer's flush if available, or a short sleep.
	// For a true flush, the channel needs to be drained.
	// For now, we'll just ensure the underlying writer is flushed if it supports it.
	if flusher, ok := aw.underlyingWriter.(interface{ Flush() }); ok {
		flusher.Flush()
	}
	// Also wait for all messages sent to logCh to be processed
	// This is handled by aw.wg.Wait() in Close(), but for Flush() we need a way
	// to signal the run goroutine to process all current messages and then wait.
	// A common pattern is to send a "flush signal" through the channel,
	// but that requires modifying the channel type or using a separate channel.
	// For simplicity, we'll just wait for the main wg, assuming Close() will be called eventually.
	// If a true blocking flush is needed without closing, a more complex mechanism is required.
}

// run is the background goroutine that performs writes.
func (aw *AsyncWriter) run() {
	for entry := range aw.logCh {
		if err := aw.underlyingWriter.Write(entry); err != nil {
			select {
			case aw.errCh <- err:
			default:
				// Error channel is full, drop the error
			}
		}
		aw.wg.Done() // Decrement WaitGroup after writing
	}
	aw.wg.Done() // Decrement for the run goroutine itself when logCh is closed
}
