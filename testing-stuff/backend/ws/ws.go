package ws

import (
	"backend/writer"
	"io"
	"sync"

	"github.com/gorilla/websocket"
)

// ThreadSafeWriter is a wrapper around a websocket connection
type ThreadSafeWriter struct {
	lock *sync.Mutex

	// TODO: figure out if there is a generic interface for this.
	//   otherwise, leave it as it is
	c *websocket.Conn
}

var _ io.Closer = &ThreadSafeWriter{}
var _ writer.Writer = &ThreadSafeWriter{}

// NewThreadSafeWriter creates a new ThreadSafeWriter
func NewThreadSafeWriter(c *websocket.Conn) ThreadSafeWriter {
	return ThreadSafeWriter{lock: &sync.Mutex{}, c: c}
}

// Close closes the writer
func (t *ThreadSafeWriter) Close() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.c == nil {
		return nil
	}
	t.c = nil

	return t.c.Close()
}

// Write writes a message to the websocket connection (assuming there is a
// connection)
func (t ThreadSafeWriter) Write(message []byte) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	return writeTextMessage(t.c, message)
}

// WriteJSON writes a JSON message to the websocket connection (assuming there
// is a connection)s
func (t ThreadSafeWriter) WriteJSON(message interface{}) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	return writeJSONMessage(t.c, message)
}
