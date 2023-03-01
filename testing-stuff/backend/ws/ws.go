package ws

import (
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

// NewThreadSafeWriter creates a new ThreadSafeWriter
func NewThreadSafeWriter(c *websocket.Conn) ThreadSafeWriter {
	return ThreadSafeWriter{lock: &sync.Mutex{}, c: c}
}

func (t *ThreadSafeWriter) Write(message []byte) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	return writeTextMessage(t.c, message)
}

func (t ThreadSafeWriter) WriteJSON(message interface{}) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	return writeJSONMessage(t.c, message)
}
