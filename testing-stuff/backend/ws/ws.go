package ws

import (
	"errors"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// ThreadSafeWriter is a wrapper around a websocket connection
type ThreadSafeWriter struct {
	isConnected bool
	lock        *sync.Mutex

	// TODO: figure out if there is a generic interface for this.
	//   otherwise, leave it as it is
	c *websocket.Conn
}

func NewDisconnectedThreadSafeWriter() ThreadSafeWriter {
	return ThreadSafeWriter{isConnected: false, lock: &sync.Mutex{}, c: nil}
}

// NewThreadSafeWriter creates a new ThreadSafeWriter
func NewThreadSafeWriter(c *websocket.Conn) ThreadSafeWriter {
	return ThreadSafeWriter{isConnected: true, lock: &sync.Mutex{}, c: c}
}

// IsConnected returns whether or not the connection is connected
func (t ThreadSafeWriter) IsConnected() bool {
	return t.isConnected
}

// Write writes a message to the websocket connection (assuming there is a
// connection)
func (t *ThreadSafeWriter) Write(message []byte) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if !t.isConnected {
		fmt.Println("Not connected not writing")
		return errors.New("cannot be written to a disconnected connection")
	}

	return writeTextMessage(t.c, message)
}

// WriteJSON writes a JSON message to the websocket connection (assuming there
// is a connection)s
func (t ThreadSafeWriter) WriteJSON(message interface{}) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if !t.isConnected {
		fmt.Println("Not connected not writing")
		return errors.New("cannot be written to a disconnected connection")
	}

	return writeJSONMessage(t.c, message)
}
