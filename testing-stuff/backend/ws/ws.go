package ws

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 60 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
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

func (t *ThreadSafeWriter) Write(message []byte) {
	t.lock.Lock()
	defer t.lock.Unlock()
	writeTextMessage(t.c, message)
}
